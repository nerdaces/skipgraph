package skipgraph

import "math/rand"

// TestForSearch returns the average path length as a result of test
// receive a graph and a search function (e.g. Search, DSGSearch)
// test iter2 * iter1 times
func TestForSearch(g []*Node, searchFunc func(*Node, int, int, bool) (*Node, int)) float64 {
	sum := 0
	for i := 0; i < iter1; i++ {
		for j := 0; j < iter2; j++ {
			randNode := g[rand.Intn(nodeNum)]
			randKey := g[rand.Intn(nodeNum)].key
			_, path := searchFunc(randNode, randKey, randNode.maxLevel-1, false)
			sum += path
		}
	}
	ave := float64(sum) / (float64(iter1) * float64(iter2))
	return ave
}

// TestForSearchFromLeft returns the 2 average path lengths below as a result of test
// 1. average path length from left end node (graph[0]) to left side nodes (10% nodes from left)
// 1. average path length from left end node (graph[0]) to right side nodes (10% nodes from right)
// receive a graph and a search function (e.g. Search, DSGSearch)
// test iter2 * iter1 times
func TestForSearchFromLeft(g []*Node, searchFunc func(*Node, int, int, bool) (*Node, int)) (float64, float64) {
	sumLeft, sumRight := 0, 0
	startNode := g[0] //mid:nodeNum/2
	var dif int = nodeNum / 10
	for i := 0; i < iter1; i++ {
		for j := 0; j < iter2; j++ {
			random := rand.Intn(dif)
			_, path := searchFunc(startNode, g[random].key, startNode.maxLevel-1, false)
			sumLeft += path
			_, path = searchFunc(startNode, g[nodeNum-random].key, startNode.maxLevel-1, false)
			sumRight += path
		}
	}
	aveLeft := float64(sumLeft) / (float64(iter1) * float64(iter2))
	aveRight := float64(sumRight) / (float64(iter1) * float64(iter2))
	return aveLeft, aveRight
}

// TestForSearchFromLeftDist returns list of average path lengths
// each average path length is that of (100 / iterD1)% nodes
// for example, if iterD1=1-0, _[1]　indicates that of 10% to 20% nodes from left end node
// test iterD2 times for each range
func TestForSearchFromLeftDist(g []*Node, searchFunc func(*Node, int, int, bool) (*Node, int)) []float64 {
	sum := 0
	startNode := g[0] //mid:nodeNum/2
	var ave []float64
	var dif int = nodeNum / iterD1
	min := -dif
	for i := 0; i < iterD1; i++ {
		min += dif
		for j := 0; j < iterD2; j++ {
			random := rand.Intn(dif) + min
			_, path := searchFunc(startNode, g[random].key, startNode.maxLevel-1, false)
			sum += path
		}
		temp := float64(sum) / float64(iterD2)
		ave = append(ave, temp)
		sum = 0
	}
	return ave
}

// TestForRangeSearch returns the average path length as a result of test
// receive a graph and a search function (e.g. searchStartEnd, DSGSearchStartEnd)
// test iter2 * iter1 times
func TestForRangeSearch(g []*Node, searchFunc func(*Node, int, int, int) (*Node, *Node, float64)) float64 {
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

// TestForSequential returns the average path length as a result of test
// receive a graph and a search function (e.g. Sequential, Sequential2)
// test iter2 * iter1 times
func TestForSequential(g []*Node, searchFunc func(*Node, int, int, int) ([]*Node, float64)) float64 {
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

// just swap 2 variables and return it
func swap(a, b int) (int, int) {
	return b, a
}

// Difference returns difference of the each element in the list
func Difference(a, b []float64) []float64 {
	var dif []float64
	for i := range a {
		dif = append(dif, a[i]-b[i])
	}
	return dif
}

// Ave returns average of the elements in the list
func Ave(a []float64) float64 {
	var sum float64 = 0
	for i := range a {
		sum += a[i]
	}
	ave := sum / float64(len(a))
	return ave
}
