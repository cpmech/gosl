package graph

// PageRank implements
type PageRank struct {
	DirectedNetwork [][]bool
	H               [][]float64
	G               [][]float64
	Rank            []int64
	Directed        map[int64]bool
	NodeNum         int64
}

// Init return pagerank
func Init(nodeNum int64, edges [][2]int64) *PageRank {
	pagerank := new(PageRank)
	pagerank.NodeNum = nodeNum
	pagerank.DirectedNetwork = make([][]bool, nodeNum)
	pagerank.H = make([][]float64, nodeNum)
	pagerank.G = make([][]float64, nodeNum)
	for node := range pagerank.DirectedNetwork {
		pagerank.DirectedNetwork[node] = make([]bool, nodeNum)
		pagerank.H[node] = make([]float64, nodeNum)
		pagerank.G[node] = make([]float64, nodeNum)
	}
	for i := 0; int64(i) < nodeNum; i++ {
		pagerank.Directed[int64(i)] = false
	}
	for edge := range edges {
		pagerank.DirectedNetwork[edges[edge][0]][edges[edge][1]-1] = true
		pagerank.Directed[edges[edge][1]-1] = true
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
				o.H[from][to] = float64(1 / refNum)
			} else {
				o.H[from][to] = 0
			}
		}
	}
	return o.H
}

// GoogleMatrix calclate GoogleMatrix based on random surfer, etc..
func (o *PageRank) GoogleMatrix(randomsurfer float64) [][]float64 {
	if randomsurfer >= 0 && randomsurfer <= 1 {
		return nil
	}
	S := o.calcS()
	for row := range S {
		for col := range S[row] {
			o.G[row][col] = randomsurfer*S[row][col] + (1.0-randomsurfer)/float64(o.NodeNum)
		}
	}
	return o.G
}

func (o *PageRank) calcS() [][]float64 {
	a := make([]float64, o.NodeNum)
	e := make([]float64, o.NodeNum)
	for i := 0; int64(i) < o.NodeNum; i++ {
		a[i] = 0
		if o.Directed[int64(i)] {
			a[i] = 1
		}
	}
	e[0] = 1
	for i := 0; int64(i) < o.NodeNum; i *= 2 {
		copy(e[:i], e[i:])
	}
	coefficient := outer(a, e)
	return add(o.H, coefficient)
}

func outer(m, n []float64) [][]float64 {
	mlen := len(m)
	nlen := len(n)
	res := make([][]float64, nlen)
	for i := range res {
		res[i] = make([]float64, mlen)
	}
	for col := range m {
		for row := range n {
			res[row][col] = m[col] * n[row]
		}
	}
	return res
}

func add(m, n [][]float64) [][]float64 {
	res := make([][]float64, len(m))
	for row := range m {
		res[row] = make([]float64, len(m[row]))
		for col := range m[row] {
			res[row][col] = m[row][col] + n[row][col]
		}
	}
	return res
}

// CalcPageRank calclate each nodes rank
func (o *PageRank) CalcPageRank() {

}
