package main

type Ticket struct {
	ID           string `json:"id"`
	Event        string `json:"event"`
	TicketHolder string `json:"ticketHolder"`
	Signature    string `json:"signature"`
}

// Add functions related to Ticket struct here
