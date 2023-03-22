package ticket

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"example.com/ticket_system/qr"
)

type Ticket struct {
	ID     int
	Seller string
	Price  float64
}

func NewTicket(id int, seller string, price float64) *Ticket {
	return &Ticket{
		ID:     id,
		Seller: seller,
		Price:  price,
	}
}

func (t *Ticket) WriteToQR() ([]byte, error) {
	// Write the ticket information to a QR code
	data := t.serialize()
	qrBytes, err := qr.GenerateQRCode(data, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return qrBytes, nil
}

func (t *Ticket) VerifyQR(qrData string) bool {
	// Verify the integrity of the ticket information in the QR code
	hash := t.calculateHash()
	dataHash := qr.calculateDataHash(qrData)
	return hash == dataHash
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
