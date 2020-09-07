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

// NewPageRank return pagerank
func NewPageRank(nodeNum int64, edges [][2]int64) *PageRank {
	pagerank := new(PageRank)
	pagerank.NodeNum = nodeNum
	pagerank.DirectedNetwork = make([][]bool, nodeNum)
	pagerank.H = make([][]float64, nodeNum)
	pagerank.G = make([][]float64, nodeNum)
	pagerank.Directed = make(map[int64]bool)
	for node := range pagerank.DirectedNetwork {
		pagerank.DirectedNetwork[node] = make([]bool, nodeNum)
		pagerank.H[node] = make([]float64, nodeNum)
		pagerank.G[node] = make([]float64, nodeNum)
	}
	for i := 0; int64(i) < nodeNum; i++ {
		pagerank.Directed[int64(i)] = false
	}
	for edge := range edges {
		pagerank.DirectedNetwork[edges[edge][0]-1][edges[edge][1]-1] = true
		pagerank.Directed[edges[edge][1]-1] = true
	}
	return pagerank
}

// HyperLinkMatrix calclate HyperLinkMatrix based on DirectedNetwork, return HyperLinkMatrix
func (o *PageRank) HyperLinkMatrix() [][]float64 {
	for from := range o.DirectedNetwork {
		refNum := 0
		var directeds []int64
		for to := range o.DirectedNetwork[from] {
			if o.DirectedNetwork[from][to] {
				refNum++
				directeds = append(directeds, int64(to))
			}
		}
		for node := range directeds {
			o.H[from][directeds[node]] = 1.0 / float64(refNum)
		}
	}
	return o.H
}

// GoogleMatrix calclate GoogleMatrix based on random surfer, etc..
func (o *PageRank) GoogleMatrix(randomsurfer float64) [][]float64 {
	if randomsurfer < 0.0 || randomsurfer > 1.0 {
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
		a[i] = 1.0
		if o.Directed[int64(i)] {
			a[i] = 0.0
		}
	}
	e[0] = 1
	for i := 1; int64(i) < o.NodeNum; i *= 2 {
		copy(e[i:], e[:i])
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

func multi(m, n [][]float64) [][]float64 {
	lencol := len(m[0])
	lenrow := len(n)
	if lencol != lenrow {
		return nil
	}
	res := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		res[i] = make([]float64, len(n[0]))
		for j := 0; j < len(n[0]); j++ {
			var val float64
			for k := 0; k < lenrow; k++ {
				val += m[i][k] * n[k][j]
			}
			res[i][j] = val
		}
	}
	return res
}

// CalcPageRank calclate each nodes pagerank
func (o *PageRank) CalcPageRank(randomsurfer float64, repetition int) []float64 {
	o.HyperLinkMatrix()
	o.GoogleMatrix(0.5)
	pie := make([]float64, o.NodeNum)
	pie[0] = 1 / float64(o.NodeNum)
	for i := 1; int64(i) < o.NodeNum; i *= 2 {
		copy(pie[i:], pie[:i])
	}
	base := make([][]float64, 1)
	base[0] = pie
	for i := 0; i < repetition; i++ {
		multi(base, o.G)
	}
	return base[0]
}
