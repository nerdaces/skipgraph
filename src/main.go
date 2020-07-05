package main

import (
	"fmt"

	"github.com/nerdaces/skipgraph/src/skipgraph"
)

func main() {
	// rand.Seed(time.Now().UnixNano())

	graph := skipgraph.SetGraph()
	// printNodesInfo(graph)
	// printNodesKeyMV(graph)

	// test for simple search
	ave := skipgraph.TestForSearch(graph, skipgraph.Search)
	fmt.Println(ave)

	// test for DSG search
	ave = skipgraph.TestForSearch(graph, skipgraph.DSGSearch)
	fmt.Println(ave)

	// test for sequential
	ave = skipgraph.TestForSequential(graph, skipgraph.Sequential)
	fmt.Println(ave)

	// test for sequential2
	ave = skipgraph.TestForSequential(graph, skipgraph.Sequential2)
	fmt.Println(ave)

	// test for search from left using normal search
	aveLeft, aveRight := skipgraph.TestForSearchFromLeft(graph, skipgraph.Search)
	fmt.Println(aveLeft, aveRight)

	// test for search from left using DSG search
	aveLeft, aveRight = skipgraph.TestForSearchFromLeft(graph, skipgraph.DSGSearch)
	fmt.Println(aveLeft, aveRight)

	// test for search from left and get the distribution of path length
	l1 := skipgraph.TestForSearchFromLeftDist(graph, skipgraph.Search)
	l2 := skipgraph.TestForSearchFromLeftDist(graph, skipgraph.DSGSearch)
	fmt.Println(l1)
	fmt.Println(l2)
	fmt.Println(skipgraph.Ave(l1))
	fmt.Println(skipgraph.Ave(l2))
}
