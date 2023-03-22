package node

type Node struct {
	snowball *Snowball
}

func NewNode(snowball *Snowball) *Node {
	return &Node{
		snowball: snowball,
	}
}

func (n *Node) Vote(block *TicketBlock) {
	// Simulate random voting behavior (weighted by block confidence)
	if rand.Float64() < float64(n.snowball.votes[block.Hash])/float64(n.snowball.beta) {
		n.snowball.Vote(block)
	}
}
