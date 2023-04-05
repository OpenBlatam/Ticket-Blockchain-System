package database

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/cache"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type TicketDatabase struct {
	db               *leveldb.DB
	writeBuffer      chan *leveldb.Batch
	writeBufferBytes *atomic.Value
	readCache        *cache.Cache
	writeAheadLog    *WriteAheadLog
	backgroundTask   *BackgroundTask
	keyBufferPool    *sync.Pool
	valueBufferPool  *sync.Pool
	redisPool        *redis.Pool
	frequencyCache   *redisCache
}

func NewTicketDatabase(file string, cacheSize, writeBufferSize int, useWriteAheadLog bool, redisURL string) (*TicketDatabase, error) {
	opts := &opt.Options{
		Filter:                filter.NewBloomFilter(10),
		BlockCache:            cache.NewLRU(cacheSize),
		WriteBuffer:           writeBufferSize,
		CompactionL0:          &opt.LevelCompaction{TargetFileSize: 64 * 1024 * 1024},
		CompactionL0Trigger:   8,
		CompactionTotalSize:   512 * 1024 * 1024,
		CompactionConcurrency: 4,
		Compression:           opt.SnappyCompression,
	}
	db, err := leveldb.OpenFile(file, opts)
	if err != nil {
		return nil, err
	}

	var writeAheadLog *WriteAheadLog
	if useWriteAheadLog {
		writeAheadLog, err = NewWriteAheadLog(file + ".wal")
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	backgroundTask := NewBackgroundTask(db, writeAheadLog)

	td := &TicketDatabase{
		db:               db,
		writeBuffer:      make(chan *leveldb.Batch, 16),
		writeBufferBytes: &atomic.Value{},
		readCache:        cache.NewLRU(cacheSize),
		writeAheadLog:    writeAheadLog,
		backgroundTask:   backgroundTask,
		keyBufferPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 64)
			},
		},
		valueBufferPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024)
			},
		},
	}

	if redisURL != "" {
		td.redisPool = &redis.Pool{
			MaxIdle:     10,
			MaxActive:   100,
			IdleTimeout: 5 * time.Minute,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(redisURL)
			},
		}

		td.frequencyCache = newRedisCache(td.redisPool, "frequency", 24*time.Hour)
	}

	for i := 0; i < 16; i++ {
		go td.writeWorker()
	}

	go td.backgroundTask.Start()

	return td, nil
}

func (td *TicketDatabase) Put(key, value []byte) {
	// Get key buffer from pool
	keyLen := len(key)
	keyBuf := td.keyBufferPool.Get().([]byte)[:0]
	keyBuf = append(keyBuf, key...)

	// Get value buffer from pool
	valueLen := len(value)
	valueBuf := td.valueBufferPool.Get().([]byte)[:0]
	valueBuf = append(valueBuf, value...)

	// Add write to batch
	writeBatch := td.writeBufferBytes.Load().(*leveldb.Batch)
	writeBatch.Put(keyBuf, valueBuf)

	// Return key and value buffers to pool
	td.keyBufferPool.Put(keyBuf[:0])
	td.valueBufferPool.Put(valueBuf[:0])

	// Flush write buffer if it reaches a certain size
	if writeBatch.ValueSize()+keyLen+valueLen > 16*1024*1024 {
		td.flushWriteBuffer()
	}
}

func (td *TicketDatabase) Get(key []byte) ([]byte, error) {
	// Check read cache first
	if val, ok := td.readCache.Get(string(key)); ok {
		return val.([]byte), nil
	} // Check frequency cache
	if td.frequencyCache != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if val, err := td.frequencyCache.get(ctx, string(key)); err == nil {
			td.readCache.Add(string(key), val)
			return val, nil
		}
	}

	// Otherwise, check LevelDB database
	val, err := td.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Update read cache
	td.readCache.Add(string(key), val)

	return val, nil
}

func (td *TicketDatabase) Delete(key []byte) error {
	// Add write to batch
	writeBatch := td.writeBufferBytes.Load().(*leveldb.Batch)
	writeBatch.Delete(key)

	// Remove key from read cache
	td.readCache.Remove(string(key))

	return nil
}

func (td *TicketDatabase) writeWorker() {
	for writeBatch := range td.writeBuffer {
		err := td.db.Write(writeBatch, &opt.WriteOptions{Sync: false})
		if err != nil {
			fmt.Println("Error flushing write buffer:", err)
		}
		// Write to write-ahead log, if enabled
		if td.writeAheadLog != nil {
			err = td.writeAheadLog.Write(writeBatch)
			if err != nil {
				fmt.Println("Error writing to write-ahead log:", err)
			}
		}

		// Clear batch and update byte count
		writeBatch.Reset()
		td.writeBufferBytes.Store(writeBatch)
	}
	// Write to write-ahead log, if enabled
	if td.writeAheadLog != nil {
		err = td.writeAheadLog.Write(writeBatch)
		if err != nil {
			fmt.Println("Error writing to write-ahead log:", err)
		}
	}

	// Clear batch and update byte count
	writeBatch.Reset()
	td.writeBufferBytes.Store(writeBatch)
}
}

func (td *TicketDatabase) flushWriteBuffer() {
	// Swap write buffer
	writeBatch := td.writeBufferBytes.Load().(*leveldb.Batch)
	newWriteBatch := new(leveldb.Batch)
	td.writeBufferBytes.Store(newWriteBatch)

	// Send old write buffer to write worker
	td.writeBuffer <- writeBatch
}

func (td *TicketDatabase) Close() error {
	close(td.writeBuffer)
	td.backgroundTask.Stop()
	err := td.db.Close()
	if err != nil {
		return err
	}

	if td.writeAheadLog != nil {
		err = td.writeAheadLog.Close()
		if err != nil {
			return err
		}
	}

	if td.redisPool != nil {
		td.redisPool.Close()
	}

	return nil
	err := td.db.Close()
	if err != nil {
		return err
	}

	if td.writeAheadLog != nil {
		err = td.writeAheadLog.Close()
		if err != nil {
			return err
		}
	}

	if td.redisPool != nil {
		td.redisPool.Close()
	}

	return nil
}