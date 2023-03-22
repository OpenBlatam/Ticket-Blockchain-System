package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"snowball"
	"time"
)

type Ticket struct {
	ID     int
	Seller string
	Price  float64
}

type TicketBlock struct {
	Index         int
	Timestamp     int64
	Tickets       []*Ticket
	PrevBlockHash string
	Hash          string
}

type Node struct {
	ID     int
	Blocks []*TicketBlock
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

func NewNode(id int) *Node {
	return &Node{ID: id}
}

func (n *Node) Vote(block *TicketBlock) {
	// Assume all nodes agree to vote for the block
	n.Blocks = append(n.Blocks, block)
}

func main() {
	// Initialize a simple ticket blockchain
	genesisBlock := NewTicketBlock(0, "", []*Ticket{})
	nodes := createNodes(10)

	snow := snowball.NewSnowball(3, 5)
	block1 := NewTicketBlock(1, genesisBlock.Hash, []*Ticket{
		{ID: 1, Seller: "Alice", Price: 100.0},
	})
	snow.Run(block1, nodes)

	block2 := NewTicketBlock(2, block1.Hash, []*Ticket{
		{ID: 2, Seller: "Bob", Price: 120.0},
	})
	snow.Run(block2, nodes)

	fmt.Println("Blockchain:")
	for _, block := range []*TicketBlock{genesisBlock, block1, block2} {
		fmt.Printf("Index: %d, Hash: %s, PrevHash: %s\n", block.Index, block.Hash, block.PrevBlockHash)
	}
}

func createNodes(n int) []*Node {
	nodes := make([]*Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = NewNode(i)
	}
	return nodes
}
