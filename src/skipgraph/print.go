package skipgraph

import "fmt"

func (n *Node) printNodeInfo() {
	fmt.Printf("%+v\n", n)
}

func printNodesInfo(g []*Node) {
	for i := range g {
		fmt.Printf("%+v\n", g[i])
	}
}

func (n *Node) printNodeKeyMV() {
	fmt.Println(n.key)
	fmt.Println(n.mv)
}

func printNodesKeyMV(g []*Node) {
	for i := range g {
		fmt.Println(g[i].key)
		fmt.Println(g[i].mv)
	}
}
