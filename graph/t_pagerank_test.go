package graph

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestPageRank01(tst *testing.T) {
	chk.PrintTitle("pagerank01")

	// define new directed graph
	edges := make([][2]int64, 17)
	edges[0] = [2]int64{1, 2}
	edges[1] = [2]int64{1, 3}
	edges[2] = [2]int64{2, 4}
	edges[3] = [2]int64{3, 2}
	edges[4] = [2]int64{3, 5}
	edges[5] = [2]int64{4, 2}
	edges[6] = [2]int64{4, 5}
	edges[7] = [2]int64{4, 6}
	edges[8] = [2]int64{5, 6}
	edges[9] = [2]int64{5, 7}
	edges[10] = [2]int64{5, 8}
	edges[11] = [2]int64{6, 8}
	edges[12] = [2]int64{7, 1}
	edges[13] = [2]int64{7, 5}
	edges[14] = [2]int64{7, 8}
	edges[15] = [2]int64{8, 6}
	edges[16] = [2]int64{8, 7}

	pg := NewPageRank(8, edges)

	// directed graph list
	chk.IntAssert(int(pg.NodeNum), 8)

	hyperlinkMatrix := pg.HyperLinkMatrix()
	hyperLinkMatrixAns := [][]float64{
		{0.0, 0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0},
		{0.0, 0.5, 0.0, 0.0, 0.5, 0.0, 0.0, 0.0},
		{0.0, 0.3333333333333333, 0.0, 0.0, 0.3333333333333333, 0.3333333333333333, 0.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 0.0, 0.3333333333333333, 0.3333333333333333, 0.3333333333333333},
		{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0},
		{0.3333333333333333, 0.0, 0.0, 0.0, 0.3333333333333333, 0.0, 0.0, 0.3333333333333333},
		{0.0, 0.0, 0.0, 0.0, 0.0, 0.5, 0.5, 0.0},
	}
	for i := 0; i < len(hyperlinkMatrix); i++ {
		// check HyperLink Matrix is correct or not.
		chk.Float64s(tst, "Hyper Link Matrix", hyperlinkMatrix[i], hyperLinkMatrixAns[i], 3)
	}

	googleMatrix := pg.GoogleMatrix(0.5)
	googleMatrixAns := [][]float64{
		{0.0625, 0.3125, 0.3125, 0.0625, 0.0625, 0.0625, 0.0625, 0.0625},
		{0.0625, 0.0625, 0.0625, 0.5625, 0.0625, 0.0625, 0.0625, 0.0625},
		{0.0625, 0.3125, 0.0625, 0.0625, 0.3125, 0.0625, 0.0625, 0.0625},
		{0.0625, 0.229167, 0.0625, 0.0625, 0.229167, 0.229167, 0.0625, 0.0625},
		{0.0625, 0.0625, 0.0625, 0.0625, 0.0625, 0.229167, 0.229167, 0.229167},
		{0.0625, 0.0625, 0.0625, 0.0625, 0.0625, 0.0625, 0.0625, 0.5625},
		{0.229167, 0.0625, 0.0625, 0.0625, 0.229167, 0.0625, 0.0625, 0.229167},
		{0.0625, 0.0625, 0.0625, 0.0625, 0.0625, 0.3125, 0.3125, 0.0625},
	}
	for i := 0; i < len(googleMatrix); i++ {
		// check Google Matrix is correct or not.
		chk.Float64s(tst, "Google Matrix", googleMatrix[i], googleMatrixAns[i], 3)
	}

	res := pg.CalcPageRank(0.5, 50)
	pageRankAns := []float64{
		0.08387937453462392,
		0.1251861504095308,
		0.08346984363365594,
		0.12509307520476534,
		0.12559568131049878,
		0.14912509307520466,
		0.12827624720774378,
		0.17937453462397604,
	}

	for i := 0; i < len(pageRankAns); i++ {
		// check pagerank value is correct or not
		chk.Float64assert(res[int64(i)], pageRankAns[i])
	}
}
