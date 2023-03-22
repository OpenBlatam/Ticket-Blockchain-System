package main

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/golang/protobuf/proto"
)

type Ticket struct {
	ID     int32
	Seller string
	Price  float64
}

func (t *Ticket) serialize() ([]byte, error) {
	pb := &ticketpb.Ticket{
		Id:     t.ID,
		Seller: t.Seller,
		Price:  t.Price,
	}
	return proto.Marshal(pb)
}

func (t *Ticket) calculateHash() string {
	record, err := t.serialize()
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(record)
	return hex.EncodeToString(hash[:])
}

func generateQRCode(data []byte, level qrcode.RecoveryLevel, size int) ([]byte, error) {
	qr, err := qrcode.New(string(data), level)
	if err != nil {
		return nil, err
	}
	qr.DisableBorder = true
	return qr.PNG(size)
}

func (t *Ticket) WriteToQR() ([]byte, error) {
	data, err := t.serialize()
	if err != nil {
		return nil, err
	}
	qrBytes, err := generateQRCode(data, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return qrBytes, nil
}

func (t *Ticket) VerifyQR(qrData []byte) bool {
	hash := t.calculateHash()
	dataHash := calculateDataHash(qrData)
	return hash == dataHash
}

func calculateDataHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
