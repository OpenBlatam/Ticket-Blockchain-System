package Database

import (
	"fmt"
	"github.com/ethereum/eth-go/ethutil"
	"github.com/syndtr/goleveldb/leveldb"
	"path"
)

type TicketDatabase struct {
	db *level.db.DB
}

func NewTicketDatabase


