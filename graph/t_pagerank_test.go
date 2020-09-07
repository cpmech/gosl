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
	// res := pg.CalcPageRank(0.5, 50)

}
