package main

type TicketBlock struct {
	Index        int      `json:"index"`
	Timestamp    int64    `json:"timestamp"`
	Tickets      []Ticket `json:"tickets"`
	PreviousHash string   `json:"previousHash"`
	Nonce        int      `json:"nonce"`
	Hash         string   `json:"hash"`
}

//
