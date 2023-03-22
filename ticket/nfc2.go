package nfc

import (
	"crypto/sha256"
	"encoding/hex"
)

type NFCTag struct {
	ID   string
	Data string
}

func NewNFCTag(id, data string) *NFCTag {
	return &NFCTag{
		ID:   id,
		Data: data,
	}
}

func (tag *NFCTag) Read() string {
	// Read data from NFC tag
	// Implement NFC reading logic here
	return tag.Data
}

func (tag *NFCTag) Write(data string) {
	// Write data to NFC tag
	// Implement NFC writing logic here
	tag.Data = data
}

func (tag *NFCTag) Verify(hash string) bool {
	// Verify the integrity of data on the NFC tag
	return tag.calculateHash() == hash
}

func (tag *NFCTag) calculateHash() string {
	h := sha256.New()
	h.Write([]byte(tag.ID + tag.Data))
	return hex.EncodeToString(h.Sum(nil))
}
