package snowball

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

type TicketBlock struct {
	Index         int
	Timestamp     int64
	Tickets       []*Ticket
	PrevBlockHash string
	Hash          string
}

func NewTicketBlock(index int, prevBlockHash string, tickets []*Ticket) *TicketBlock {
	block := &TicketBlock{
		Index:         index,
		Timestamp:     time.Now().Unix(),
		Tickets:       tickets,
		PrevBlockHash: prevBlockHash,
	}
	block.Hash = block.calculateHash()
	return block
}

func (b *TicketBlock) calculateHash() string {
	record := fmt.Sprintf("%d%d%s", b.Index, b.Timestamp, b.PrevBlockHash)
	hash := sha256.New()
	hash.Write([]byte(record))
	return hex.EncodeToString(hash.Sum(nil))
}

type Ticket struct {
	ID     int
	Seller string
	Price  float64
}

type Node struct {
	ID     int
	Blocks []*TicketBlock
}

func NewNode(id int) *Node {
	return &Node{ID: id}
}

func (n *Node) Vote(block *TicketBlock) {
	// Assume all nodes agree to vote for the block
	n.Blocks = append(n.Blocks, block)
}

type Snowball struct {
	k        int
	beta     int
	votes    map[string]int
	finalize chan *TicketBlock
}

func NewSnowball(k, beta int) *Snowball {
	return &Snowball{
		k:        k,
		beta:     beta,
		votes:    make(map[string]int),
		finalize: make(chan *TicketBlock),
	}
}

func (s *Snowball) Vote(block *TicketBlock) {
	s.votes[block.Hash]++
	if s.votes[block.Hash] >= s.beta {
		s.finalize <- block
	}
}

func (s *Snowball) Run(block *TicketBlock, nodes []*Node) {
	for {
		select {
		case <-s.finalize:
			println("Block finalized:", block.Hash)
			return
		default:
			s.voteRound(block, nodes)
			time.Sleep(time.Duration(s.k) * time.Millisecond)
		}
	}
}

func (s *Snowball) voteRound(block *TicketBlock, nodes []*Node) {
	nodeCount := len(nodes)
	for i := 0; i < s.k; i++ {
		nodeIdx := rand.Intn(nodeCount)
		nodes[nodeIdx].Vote(block)
	}
}
