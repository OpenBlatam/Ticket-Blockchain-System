package main

import (
	"time"
)

// 1. Define parameters: k (sample size), α (confidence), and β (finalization)
// 2. For each round:
//     a. Query k randomly selected nodes for their preferred block
//     b. Tally the votes
//     c. If a block has α more votes than the other, increase its confidence counter
//     d. If a block's confidence counter reaches β, finalize the decision and end the consensus

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
	// Vote for the block
	s.votes[block.Hash]++

	// If the block has reached the threshold, finalize it
	if s.votes[block.Hash] >= s.beta {
		s.finalize <- block
	}
}

func (s *Snowball) Run(block *TicketBlock, nodes []*Node) {
	// Simulate the voting process
	// the chatgpt simulate a lot! the snow folder for proto implemnetation in better
	for {
		select {
		case <-s.finalize:
			println("Block finalized:", block.Hash)
			return
		default:
			for _, node := range nodes {
				node.Vote(block)
			}
			time.Sleep(time.Duration(s.k) * time.Millisecond)
		}
	}
}
