package skipgraph

import (
	"bytes"
	"fmt"
	"sort"
)

// DefaultMaxlevel is the maximum level of each node
const (
	DefaultMaxLevel = 32
	nodeNum         = 10000
	left            = 0
	right           = 1
	iter1           = 5
	iter2           = 100
	iterD1          = 10
	iterD2          = 100
	targetRange     = 100
)

// SetGraph sets Skip Graph from scratch, then return the graph
// At first, create Nodes and sort them
// After that, set neighbors for each Nodes at each level
func SetGraph() []*Node {
	g := make([]*Node, nodeNum)
	// rand.Seed(time.Now().UnixNano())
	for i := range g {
		g[i] = initNode()
	}

	sortNodes(g)

	level := 0
	g[0].neighbors[right][level] = g[1]
	for i := 1; i < nodeNum-1; i++ {
		g[i].neighbors[left][level] = g[i-1]
		g[i].neighbors[right][level] = g[i+1]
	}
	g[nodeNum-1].neighbors[left][level] = g[nodeNum-2]

	for level < DefaultMaxLevel {
		level++
		for i := range g {
			prev := g[i].neighbors[left][level-1]
			next := g[i].neighbors[right][level-1]
			for prev != nil {
				if bytes.Equal(g[i].mv[:level], prev.mv[:level]) {
					g[i].neighbors[left][level] = prev
					break
				}
				prev = prev.neighbors[left][level-1]
			}
			for next != nil {
				if bytes.Equal(g[i].mv[:level], next.mv[:level]) {
					g[i].neighbors[right][level] = next
					break
				}
				next = next.neighbors[right][level-1]
			}
		}
	}

	return g
}

// receive a graph and sort the Nodes by key
func sortNodes(g []*Node) {
	sort.SliceStable(g, func(i, j int) bool { return g[i].key < g[j].key })
}

// Search is the normal search function of Skip Graph, one-way search
// receive start node, target key, level and flagRange,
// then return target node and path length from start node to target node
// search from start node to the node which has the target key from the designated level
// flagRange=true : if you want to get only target node unless the target key does not exist in the graph
// -> if it does, return nil and 0
// flagRange=false: if you want to get the neighborhood unless the target key does not exist in the graph
// -> return the neighborhood node of the target key
func Search(startNode *Node, searchKey, level int, flagRange bool) (v *Node, pathLength int) {
	pathLength = 0
	v = startNode
	if v.key < searchKey {
		for level >= 0 {
			if v.key == searchKey {
				return
			} else if v.neighbors[right][level] != nil && v.neighbors[right][level].key <= searchKey {
				v = v.neighbors[right][level]
				pathLength++
			} else {
				level--
			}
		}
	} else {
		for level >= 0 {
			if v.key == searchKey {
				return
			} else if v.neighbors[left][level] != nil && v.neighbors[left][level].key >= searchKey {
				v = v.neighbors[left][level]
				pathLength++
			} else {
				level--
			}
		}
	}
	if flagRange {
		return
	}
	v = nil
	pathLength = 0
	return
}

// DSGSearch is the search function of Detouring Skip Graph, 2 method below are reflected on normal search
// 1. search from max level every time when a node receives a query
// 2. detour if there is shortcut
// receive start node, target key, level and flagRange,
// then return target node and path length from start node to target node
func DSGSearch(startNode *Node, searchKey, level int, flagRange bool) (v *Node, pathLength int) {
	pathLength = 0
	v = startNode
	for level >= 0 {
		if v.key < searchKey {
			for level >= 0 {
				if v.key == searchKey {
					return
				}
				frontNode := v.neighbors[right][level]
				if frontNode != nil && frontNode.key <= searchKey {
					backNode := v.neighbors[right][level+1]
					if backNode != nil && closeToRight(searchKey, frontNode.key, backNode.key) {
						v = backNode
						level = v.maxLevel - 1
						pathLength++
						break
					} else {
						v = frontNode
						level = v.maxLevel - 1
						pathLength++
					}
				} else {
					level--
				}
			}
		} else {
			for level >= 0 {
				if v.key == searchKey {
					return
				}
				frontNode := v.neighbors[left][level]
				if frontNode != nil && frontNode.key >= searchKey {
					backNode := v.neighbors[left][level+1]
					if backNode != nil && closeToLeft(searchKey, backNode.key, frontNode.key) {
						v = backNode
						level = v.maxLevel - 1
						pathLength++
						break
					} else {
						v = frontNode
						level = v.maxLevel - 1
						pathLength++
					}
				} else {
					level--
				}
			}
		}
	}
	if flagRange {
		return
	}
	v = nil
	pathLength = 0
	return

}

// compare the average of leftKey and rightKey and searchKey
// if former is bigger, return true
// otherwise, return false
func closeToLeft(searchKey, leftKey, rightKey int) bool {
	mid := (leftKey + rightKey) / 2
	if mid > searchKey {
		return true
	}
	return false

}

// compare the average of leftKey and rightKey and searchKey
// if former is smaller, return true
// otherwise, return false
func closeToRight(searchKey, leftKey, rightKey int) bool {
	mid := (leftKey + rightKey) / 2
	if mid < searchKey {
		return true
	}
	return false

}

// return the 2 nodes that have bottomKey and topKey respectively, and the average path length from start node to those nodes
// using Search function
func searchStartEnd(startNode *Node, bottomKey, topKey, level int) (firstNode, lastNode *Node, avePathLength float64) {
	if bottomKey > topKey {
		fmt.Println("Error: bottomKey must be smaller than topKey.")
		return
	}
	firstNode, pathForFirst := Search(startNode, bottomKey, level, true)
	lastNode, pathForLast := Search(startNode, topKey, level, true)

	if bottomKey == topKey {
		avePathLength = float64(pathForFirst) + float64(pathForLast)
		return
	}

	rightNode := firstNode.neighbors[right][0]
	leftNode := lastNode.neighbors[left][0]

	if startNode.key < bottomKey {
		if firstNode.key == bottomKey {
		} else if rightNode != nil && inRange(rightNode.key, bottomKey, topKey) {
			firstNode = rightNode
			pathForFirst++
		} else {
			firstNode, lastNode, avePathLength = errorForRange()
			return
		}
	} else if startNode.key > topKey {
		if lastNode.key == topKey {
		} else if lastNode != nil && inRange(leftNode.key, bottomKey, topKey) {
			lastNode = leftNode
			pathForLast++
		} else {
			firstNode, lastNode, avePathLength = errorForRange()
			return
		}
	}

	avePathLength = (float64(pathForFirst) + float64(pathForFirst)) / 2
	return
}

// return the 2 nodes that have bottomKey and topKey respectively, and the average path length from start node to those nodes
// using DSGSearch function
func dsgSearchStartEnd(startNode *Node, bottomKey, topKey, level int) (firstNode, lastNode *Node, pathLength float64) {
	if bottomKey > topKey {
		fmt.Println("Error: bottomKey must be smaller than topKey.")
		return
	}
	firstNode, pathForFirst := DSGSearch(startNode, bottomKey, level, true)
	lastNode, pathForLast := DSGSearch(startNode, topKey, level, true)

	if bottomKey == topKey {
		pathLength = float64(pathForFirst)
		return
	}

	rightNode := firstNode.neighbors[right][0]
	leftNode := lastNode.neighbors[left][0]

	if firstNode.key == bottomKey {
	} else if firstNode.key < bottomKey {
		if rightNode != nil && inRange(rightNode.key, bottomKey, topKey) {
			firstNode = rightNode
			pathForFirst++
		} else {
			firstNode, lastNode, pathLength = errorForRange()
			return
		}
	}

	if lastNode.key == topKey {
	} else if lastNode.key > topKey {
		if leftNode != nil && inRange(lastNode.key, bottomKey, topKey) {
			lastNode = leftNode
			pathForLast++
		} else {
			firstNode, lastNode, pathLength = errorForRange()
			return
		}
	}

	pathLength = (float64(pathForFirst) + float64(pathForFirst)) / 2
	return
}

// bottomKey <= targetKey <= topKey : true
// otherwise : false
func inRange(targetKey, bottomKey, topKey int) bool {
	if bottomKey <= targetKey && targetKey <= topKey {
		return true
	}
	return false
}

// print error message and return nil and 0
func errorForRange() (_, _ *Node, _ float64) {
	fmt.Println("Error: Such a node does not exist in the range.")
	return nil, nil, 0
}

// Sequential is a sequential range search function from only left end node of the range
// return the nodes list which contains all nodes of the range and average path length from start node
func Sequential(startNode *Node, bottomKey, topKey, level int) (targetNodes []*Node, avePathLength float64) {
	if bottomKey > topKey {
		fmt.Println("Error: bottomKey must be smaller than topKey.")
		return
	}

	firstNode, _, _ := dsgSearchStartEnd(startNode, bottomKey, topKey, level)
	v := firstNode

	var pathLength []int
	currentPathLength, totalPathLength := 0, 0
	for v != nil && inRange(v.key, bottomKey, topKey) {
		targetNodes = append(targetNodes, v)
		pathLength = append(pathLength, currentPathLength)
		totalPathLength += currentPathLength
		currentPathLength++
		v = v.neighbors[right][0]
	}

	avePathLength = float64(totalPathLength) / float64(len(targetNodes))
	return
}

// Sequential2 is a sequential range search function from both end of nodes of the range
// return the nodes list which contains all nodes of the range and average path length from start node
func Sequential2(startNode *Node, bottomKey, topKey, level int) (targetNodes []*Node, avePathLength float64) {
	if bottomKey > topKey {
		fmt.Println("Error: bottomKey must be smaller than topKey.")
		return
	}

	firstNode, lastNode, _ := dsgSearchStartEnd(startNode, bottomKey, topKey, level)
	f, l := firstNode, lastNode

	var pathLength []int
	lPathLength, rPathLength, totalPathLength := 0, 0, 0
	for f.key < l.key && f != nil && l != nil {
		targetNodes = append(targetNodes, f, l)
		pathLength = append(pathLength, lPathLength, rPathLength)
		totalPathLength += lPathLength + rPathLength
		lPathLength++
		rPathLength++
		f = f.neighbors[right][0]
		l = l.neighbors[left][0]
	}

	if f == l {
		targetNodes = append(targetNodes, f)
		pathLength = append(pathLength, lPathLength)
	}

	avePathLength = float64(totalPathLength) / float64(len(targetNodes))
	return
}
