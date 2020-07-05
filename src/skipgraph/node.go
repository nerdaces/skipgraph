package skipgraph

import (
	"math"
	"math/rand"
)

// Node is struct of a node for Skip Graph
type Node struct {
	key       int
	mv        []byte
	maxLevel  int
	neighbors [][]*Node
}

func initNode() *Node {
	return &Node{
		key:       getKey(),
		mv:        getMV(),
		maxLevel:  getMaxLevel(),
		neighbors: getDoubleList(),
	}
}

// return a key selected randomly in [0,2^30)
func getKey() int {
	maxKey := int(math.Pow(2, 30))
	return rand.Intn(maxKey)
}

// return a membership vector decided randomly
// MV is defined as 32 bit strings here
// whose alphabet size is just 2, and alphabet is {0, 1}
func getMV() []byte {
	alphabet := [2]byte{0, 1}
	b := make([]byte, DefaultMaxLevel)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return b
}

// return the maximum level that each node could have
// if you change the premise here that each node has common maximum level(DefaultMaxLevel),
// then you have to change the return value (also delete constant and change struct)
func getMaxLevel() int {
	return DefaultMaxLevel
}

// return a double linked list for neighbors of each nodes
func getDoubleList() [][]*Node {
	n := make([][]*Node, 2)
	for i := 0; i < 2; i++ {
		n[i] = make([]*Node, DefaultMaxLevel)
	}
	return n
}
