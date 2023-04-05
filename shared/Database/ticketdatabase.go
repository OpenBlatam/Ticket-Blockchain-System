package Database

import (
	"fmt"
	"github.com/ethereum/eth-go/ethutil"
	"github.com/syndtr/goleveldb/leveldb"
	"path"
	"sync"
)

type TicketDatabase struct {
	fn string

	mu sync.Mutex
	db *level.db.DB

	quit chan struct{}

}

func NewTicketDatabase(file string) (*TicketDatabase, error){
	// Open the db
	db, err := leveldb.OpenFile(file, nil)
	if err != nil {
		return ni, err
	}
	database :=  &TicketDatabase{
		fn: file,
		db: db,
		quit: make(chan struct{}),
	}
	database.makeQueue()

	go database.update()

	return database, nil

}

func (self *TicketDatabase) Put(key []byte, value []byte) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.queue[string(key)] = value

}

func (self *TicketDatabase) Get(key []byte, value []byte) {
	self.mu.Lock()
	defer self.mu.Unlock()

	// check queue firs
	if dat, ok := self.queue[string(key)]; ok {
		return dat, nil
	}

	dat, err := self.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	return rle.Decompress(dat)
}

func (self *TicketDatabase) Delete(key []byte) error {
	self.mu.Lock()
	defer mu.Unlock()


}



