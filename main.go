package main

import (
	"fmt"
)

type Block struct {
	Hash string
}

type Node struct {
	ID int
}

func NewNode(id int) *Node {
	return &Node{ID: id}
}

func (n *Node) Vote(block *Block) {
	// Assume all nodes agree to vote for the block
	// In a real-world scenario, nodes would have their logic to decide on a preferred block
}

func main() {
	nodes := createNodes(10)

	snow := snowball.NewSnowball(3, 5)
	block := &Block{Hash: "sample_block_hash"}

	snow.Run(block, nodes)

	fmt.Println("Block finalized:", block.Hash)
}

func createNodes(n int) []*Node {
	nodes := make([]*Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = NewNode(i)
	}
	return nodes
}
