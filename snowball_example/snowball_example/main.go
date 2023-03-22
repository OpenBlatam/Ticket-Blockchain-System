package main

import (
	"example.com/snowball_example/snowball"
	"fmt"
)

func main() {
	// Initialize a simple ticket blockchain
	genesisBlock := snowball.NewTicketBlock(0, "", []*snowball.Ticket{})
	nodes := createNodes(10)

	snow := snowball.NewSnowball(3, 5)
	block1 := snowball.NewTicketBlock(1, genesisBlock.Hash, []*snowball.Ticket{
		{ID: 1, Seller: "Alice", Price: 100.0},
	})
	snow.Run(block1, nodes)

	block2 := snowball.NewTicketBlock(2, block1.Hash, []*snowball.Ticket{
		{ID: 2, Seller: "Bob", Price: 120.0},
	})
	snow.Run(block2, nodes)

	fmt.Println("Blockchain:")
	for _, block := range []*snowball.TicketBlock{genesisBlock, block1, block2} {
		fmt.Printf("Index: %d, Hash: %s, PrevHash: %s\n", block.Index, block.Hash, block.PrevBlockHash)
	}
}

func createNodes(n int) []*snowball.Node {
	nodes := make([]*snowball.Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = snowball.NewNode(i)
	}
	return nodes
}
