package qr

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image/png"
)

func GenerateQRCode(data string, level qrcode.RecoveryLevel, size int) ([]byte, error) {
	qr, err := qrcode.New(data, level)
	if err != nil {
		return nil, err
	}
	qr.DisableBorder = true
	img := qr.ToImage(size)
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *Ticket) WriteToQR() ([]byte, error) {
	// Write the ticket information to a QR code
	data := t.serialize()
	qrBytes, err := GenerateQRCode(data, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return qrBytes, nil
}

func (t *Ticket) VerifyQR(qrData string) bool {
	// Verify the integrity of the ticket information in the QR code
	hash := t.calculateHash()
	dataHash := calculateDataHash(qrData)
	return hash == dataHash
}

func calculateDataHash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
