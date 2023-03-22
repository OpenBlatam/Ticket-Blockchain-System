package ticket

import (
	"crypto/sha256"
	"encoding/hex"
	"example.com/ticket_system/nfc"
	"fmt"
)

type Ticket struct {
	ID     int
	Seller string
	Price  float64
	Tag    *nfc.NFCTag
}

func NewTicket(id int, seller string, price float64, tag *nfc.NFCTag) *Ticket {
	return &Ticket{
		ID:     id,
		Seller: seller,
		Price:  price,
		Tag:    tag,
	}
}

func (t *Ticket) WriteToTag() {
	// Write the ticket information to the NFC tag
	data := t.serialize()
	t.Tag.Write(data)
}

func (t *Ticket) ReadFromTag() {
	// Read the ticket information from the NFC tag
	data := t.Tag.Read()
	t.deserialize(data)
}

func (t *Ticket) Verify() bool {
	// Verify the integrity of the ticket information on the NFC tag
	hash := t.calculateHash()
	return t.Tag.Verify(hash)
}

func (t *Ticket) serialize() string {
	// Serialize the ticket information into a string
	return fmt.Sprintf("%d|%s|%.2f", t.ID, t.Seller, t.Price)
}

func (t *Ticket) deserialize(data string) {
	// Deserialize the ticket information from a string
	// Implement deserialization logic here
}

func (t *Ticket) calculateHash() string {
	h := sha256.New()
	h.Write([]byte(t.serialize()))
	return hex.EncodeToString(h.Sum(nil))
}
