package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

const (
	DefaultMaxLevel = 32
	nodeNum         = 10000
	left            = 0
	right           = 1
	iter1           = 5
	iter2           = 100
	targetRange     = 10
)

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

func setGraph() []*Node {
	g := make([]*Node, nodeNum)
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

func getKey() int {
	maxKey := int(math.Pow(2, 30))
	// maxKey := 300
	// rand.Seed(time.Now().UnixNano())
	return rand.Intn(maxKey)
}

func getMV() []byte {
	alphabet := [2]byte{0, 1}
	b := make([]byte, DefaultMaxLevel)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return b
}

func getMaxLevel() int {
	return DefaultMaxLevel
}

func getDoubleList() [][]*Node {
	n := make([][]*Node, 2)
	for i := 0; i < 2; i++ {
		n[i] = make([]*Node, DefaultMaxLevel)
	}
	return n
}

func sortNodes(g []*Node) {
	sort.SliceStable(g, func(i, j int) bool { return g[i].key < g[j].key })
}

func search(startNode *Node, searchKey, level int, flagRange bool) (v *Node, pathLength int) {
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
	} else {
		v = nil
		pathLength = 0
		return
	}
}

func dsgSearch(startNode *Node, searchKey, level int, flagRange bool) (v *Node, pathLength int) {
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
	} else {
		v = nil
		pathLength = 0
		return
	}
}

func closeToLeft(searchKey, leftKey, rightKey int) bool {
	mid := (leftKey + rightKey) / 2
	if mid > searchKey {
		return true
	} else {
		return false
	}
}

func closeToRight(searchKey, leftKey, rightKey int) bool {
	mid := (leftKey + rightKey) / 2
	if mid < searchKey {
		return true
	} else {
		return false
	}
}

// have to fix the calculation of pathLength
func rangeSearch(startNode *Node, bottomKey, topKey, level int) (firstNode, lastNode *Node, avePathLength float64) {
	if bottomKey > topKey {
		fmt.Println("Error: bottomKey must be smaller than topKey.")
		return
	}
	firstNode, pathForFirst := search(startNode, bottomKey, level, true)
	lastNode, pathForLast := search(startNode, topKey, level, true)

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

// have to fix the calculation of pathLength
func dsgRangeSearch(startNode *Node, bottomKey, topKey, level int) (firstNode, lastNode *Node, pathLength float64) {
	if bottomKey > topKey {
		fmt.Println("Error: bottomKey must be smaller than topKey.")
		return
	}
	firstNode, pathForFirst := dsgSearch(startNode, bottomKey, level, true)
	lastNode, pathForLast := dsgSearch(startNode, topKey, level, true)

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

func inRange(targetKey, bottomKey, topKey int) bool {
	if bottomKey <= targetKey && targetKey <= topKey {
		return true
	}
	return false
}

func errorForRange() (_, _ *Node, _ float64) {
	fmt.Println("Error: Such a node does not exist in the range.")
	return nil, nil, 0
}

// sequential from only left node of the range
func sequential(startNode *Node, bottomKey, topKey, level int) (targetNodes []*Node, avePathLength float64) {
	if bottomKey > topKey {
		fmt.Println("Error: bottomKey must be smaller than topKey.")
		return
	}

	firstNode, _, _ := dsgRangeSearch(startNode, bottomKey, topKey, level)
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

// sequential from both side of nodes of the range
func sequential2(startNode *Node, bottomKey, topKey, level int) (targetNodes []*Node, avePathLength float64) {
	if bottomKey > topKey {
		fmt.Println("Error: bottomKey must be smaller than topKey.")
		return
	}

	firstNode, lastNode, _ := dsgRangeSearch(startNode, bottomKey, topKey, level)
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

func swap(a, b int) (int, int) {
	return b, a
}

func testForSearch(g []*Node, searchFunc func(*Node, int, int, bool) (*Node, int)) float64 {
	sum := 0
	for i := 0; i < iter1; i++ {
		for j := 0; j < iter2; j++ {
			randNode := g[rand.Intn(nodeNum)]
			randKey := g[rand.Intn(nodeNum)].key
			_, path := searchFunc(randNode, randKey, randNode.maxLevel-1, false)
			sum += path
		}
	}
	var sumf, iter1f, iter2f float64 = float64(sum), float64(iter1), float64(iter2)
	return sumf / (iter1f * iter2f)
}

func testForRangeSearch(g []*Node, searchFunc func(*Node, int, int, int) (*Node, *Node, float64)) float64 {
	var sum float64 = 0
	e := 0
	for i := 0; i < iter1; i++ {
		for j := 0; j < iter2; j++ {
			randNode := g[rand.Intn(nodeNum)]
			randStartKey := g[rand.Intn(nodeNum)].key
			randEndKey := g[rand.Intn(nodeNum)].key
			if randStartKey > randEndKey {
				randStartKey, randEndKey = swap(randStartKey, randEndKey)
			}
			_, _, path := searchFunc(randNode, randStartKey, randEndKey, randNode.maxLevel-1)
			if path != 0 {
				sum += path
			} else {
				e++
			}
		}
	}
	var iter1f, iter2f, ef float64 = float64(iter1), float64(iter2), float64(e)
	return sum / (iter1f*iter2f - ef)
}

func testForSequential(g []*Node, searchFunc func(*Node, int, int, int) ([]*Node, float64)) float64 {
	var sum float64 = 0
	e := 0
	for i := 0; i < iter1; i++ {
		for j := 0; j < iter2; j++ {
			randNode := g[rand.Intn(nodeNum)]
			startIndex := rand.Intn(nodeNum - targetRange)
			randStartKey := g[startIndex].key
			randEndKey := g[startIndex+targetRange-1].key
			_, path := searchFunc(randNode, randStartKey, randEndKey, randNode.maxLevel-1)
			if path != 0 {
				sum += path
			} else {
				e++
			}
		}
	}
	var iter1f, iter2f, ef float64 = float64(iter1), float64(iter2), float64(e)
	return sum / (iter1f*iter2f - ef)
}

func main() {
	graph := setGraph()
	// printNodesInfo(graph)
	// printNodesKeyMV(graph)

	// test for simple search
	ave := testForSearch(graph, search)
	fmt.Println(ave)

	// test for DSG search
	ave = testForSearch(graph, dsgSearch)
	fmt.Println(ave)

	// sequential
	ave = testForSequential(graph, sequential)
	fmt.Println(ave)

	// sequential2
	ave = testForSequential(graph, sequential2)
	fmt.Println(ave)
}
