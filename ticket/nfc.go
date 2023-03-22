package ticket

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func CreateNFCTagSignature(ticket *Ticket, secretKey string) string {
	data := []byte(secretKey + ticketToString(ticket))
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func VerifyNFCTagSignature(ticket *Ticket, signature, secretKey string) bool {
	expectedSignature := CreateNFCTagSignature(ticket, secretKey)
	return expectedSignature == signature
}

func ticketToString(ticket *Ticket) string {
	return fmt.Sprintf("%d:%s:%.2f", ticket.ID, ticket.Seller, ticket.Price)
}
