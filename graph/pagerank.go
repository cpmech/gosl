package graph

// PageRank implements
type PageRank struct {
	DirectedNetwork [][]bool
	PageRank        [][]float64
	Rank            []int64
	NodeNum         int64
}

// Init return pagerank
func Init(nodeNum int64, edges [][2]int64) *PageRank {
	pagerank := new(PageRank)
	pagerank.NodeNum = nodeNum
	pagerank.DirectedNetwork = make([][]bool, 0, nodeNum)
	for node := range pagerank.DirectedNetwork {
		pagerank.DirectedNetwork[node] = make([]bool, 0, nodeNum)
	}
	for edge := range edges {
		pagerank.DirectedNetwork[edges[edge][0]][edges[edge][1]] = true
	}
	return pagerank
}

// HyperLinkMatrix calclate HyperLinkMatrix based on DirectedNetwork, return HyperLinkMatrix
func (o *PageRank) HyperLinkMatrix() [][]float64 {
	for from := range o.DirectedNetwork {
		refNum := 0
		for to := range o.DirectedNetwork[from] {
			if o.DirectedNetwork[from][to] {
				refNum++
			}
		}
		for to := range o.DirectedNetwork {
			if from != to && refNum > 0 {
				o.PageRank[from][to] = float64(1 / refNum)
			} else {
				o.PageRank[from][to] = 0
			}
		}
	}
	return o.PageRank
}
