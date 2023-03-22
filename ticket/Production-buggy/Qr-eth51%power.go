package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/protobuf/proto"
)

type Ticket struct {
	ID     int32
	Seller string
	Price  float64
	Status int32
}

func (t *Ticket) serialize() ([]byte, error) {
	pb := &ticketpb.Ticket{
		Id:     t.ID,
		Seller: t.Seller,
		Price:  t.Price,
		Status: t.Status,
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

func generateQRCodeToFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	qr, err := qrcode.New(string(data), qrcode.Medium)
	if err != nil {
		return err
	}
	qr.DisableBorder = true

	err = png.Encode(file, qr.Image(256))
	if err != nil {
		return err
	}
	return nil
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
func main() {
	// Connect to the Ethereum network
	client, err := ethclient.Dial("<INFURA_PROJECT_URL>")
	if err != nil {
		log.Fatal(err)
	}

	// Load the ticket contract
	contractAddress := common.HexToAddress("<CONTRACT_ADDRESS>")
	instance, err := ticket.NewTicket(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Get the account to sign transactions
	privateKey, err := crypto.HexToECDSA("<PRIVATE_KEY>")
	if err != nil {
		log.Fatal(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)

	// Call the smart contract function to add a new ticket
	newTicket := &Ticket{
		ID:     1,
		Seller: "John",
		Price:  10.5,
		Status: 0,
	}
	data, err := newTicket.serialize()
	if err != nil {
		log.Fatal(err)
	}
	tx, err := instance.AddTicket(auth, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())

	// Generate a QR code for the new ticket
	qrCodeData, err := newTicket.WriteToQR()
	if err != nil {
		log.Fatal(err)
	}
	err = generateQRCodeToFile(qrCodeData, "ticket.png")
	if err != nil {
		log.Fatal(err)
	}

	// Verify the ticket using the QR code data
	verified := newTicket.VerifyQR(qrCodeData)
	fmt.Printf("Ticket verified: %t\n", verified)
}
