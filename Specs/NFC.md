package main

import (
   "log"
	"github.com/go-piv/piv-go/piv"
)

func ReadNFC() (*Ticket, error) {
// Connect to the NFC reader
reader := piv.New()
err := reader.Connect()
if err != nil {
return nil, err
}
defer reader.Close()

	// Read the data from the NFC tag
	data, err := reader.ReadNFC()
	if err != nil {
		return nil, err
	}

	// Parse the data and create a new ticket
	var ticketData struct {
		Event        string `json:"event"`
		TicketHolder string `json:"ticketHolder"`
	}
	err = json.Unmarshal(data, &ticketData)
	if err != nil {
		return nil, err
	}
	ticket, err := NewTicket(ticketData.Event, ticketData.TicketHolder)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
