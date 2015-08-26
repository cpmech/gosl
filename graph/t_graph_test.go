// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_graph01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("graph01")

	//           [10]
	//      0 ––––––––→ 3      numbers in parentheses
	//      |    (1)    ↑      indicate edge ids
	//   [5]|(0)        |
	//      |        (3)|[1]
	//      ↓    (2)    |      numbers in brackets
	//      1 ––––––––→ 2      indicate weights
	//           [3]

	var G Graph
	G.Init(
		// edge:  0       1       2       3
		[][]int{{0, 1}, {0, 3}, {1, 2}, {2, 3}},
		[]float64{5, 10, 3, 1}, // weights
		nil, nil,
	)

	chk.IntAssert(len(G.Shares), 4)   // nverts
	chk.IntAssert(len(G.Key2edge), 4) // nedges
	chk.IntAssert(len(G.Dist), 4)     // nverts
	chk.IntAssert(len(G.Next), 4)     // nverts

	shares := [][]int{
		{0, 1}, // edges sharing node 0
		{0, 2}, // edges sharing node 1
		{2, 3}, // edges sharing node 2
		{1, 3}, // edges sharing node 3
	}
	for k, share := range shares {
		chk.Ints(tst, io.Sf("edges sharing %d", k), G.Shares[k], share)
	}

	chk.IntAssert(G.Key2edge[G.HashEdgeKey(0, 1)], 0) // (0,1) → edge 0
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(0, 3)], 1) // (0,3) → edge 1
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(1, 2)], 2) // (1,2) → edge 2
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(2, 3)], 3) // (2,3) → edge 3

	edg, err := G.GetEdge(-1, 1)
	if err == nil {
		tst.Errorf("GetEdge should have failed with (-1,1)\n")
		return
	}
	edg, err = G.GetEdge(0, 1)
	if err != nil {
		tst.Errorf("GetEdge failed:\n%v", err)
		return
	}
	chk.IntAssert(edg, 0)
	edg, err = G.GetEdge(0, 3)
	if err != nil {
		tst.Errorf("GetEdge failed:\n%v", err)
		return
	}
	chk.IntAssert(edg, 1)
	edg, err = G.GetEdge(1, 2)
	if err != nil {
		tst.Errorf("GetEdge failed:\n%v", err)
		return
	}
	chk.IntAssert(edg, 2)
	edg, err = G.GetEdge(2, 3)
	if err != nil {
		tst.Errorf("GetEdge failed:\n%v", err)
		return
	}
	chk.IntAssert(edg, 3)

	err = G.ShortestPaths("FW")
	if err != nil {
		tst.Errorf("ShortestPaths failed:\n%v", err)
		return
	}
	inf := GRAPH_INF
	pth := G.Path(0, 3)
	io.Pforan("dist =\n%v", G.StrDistMatrix())
	io.Pforan("path from 0 to 3 = %v\n", pth)
	chk.Ints(tst, "0 → 3", pth, []int{0, 1, 2, 3})
	chk.Matrix(tst, "dist", 1e-17, G.Dist, [][]float64{
		{0, 5, 8, 9},
		{inf, 0, 3, 4},
		{inf, inf, 0, 1},
		{inf, inf, inf, 0},
	})

	G.WeightsE[3] = 13
	err = G.ShortestPaths("FW")
	if err != nil {
		tst.Errorf("ShortestPaths failed:\n%v", err)
		return
	}
	pth = G.Path(0, 3)
	io.Pf("\n")
	io.Pfcyan("dist =\n%v", G.StrDistMatrix())
	io.Pfcyan("path from 0 to 3 = %v\n", pth)
	chk.Ints(tst, "0 → 3", pth, []int{0, 3})
	chk.Matrix(tst, "dist", 1e-17, G.Dist, [][]float64{
		{0, 5, 8, 10},
		{inf, 0, 3, 16},
		{inf, inf, 0, 13},
		{inf, inf, inf, 0},
	})
}

func Test_graph02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("graph02")

	//             [3]
	//      4 –––––––––––→ 5 .  [4]      numbers in parentheses
	//      ↑      (0)     |  `.         indicate edge ids
	//      |           (4)| (6)`.v
	//      |              |       3
	//  [11]|(1)        [7]|  (5),^      numbers in brackets
	//      |              |   ,' [9]    indicate weights
	//      |   (2)    (3) ↓ ,'
	//      1 ←–––– 0 ––––→ 2
	//          [6]    [8]

	var G Graph
	G.Init(
		// edge:  0       1       2       3       4       5       6
		[][]int{{4, 5}, {1, 4}, {0, 1}, {0, 2}, {5, 2}, {2, 3}, {5, 3}},
		[]float64{3, 11, 6, 8, 7, 9, 4},
		nil, nil,
	)

	chk.IntAssert(len(G.Shares), 6)   // nverts
	chk.IntAssert(len(G.Key2edge), 7) // nedges
	chk.IntAssert(len(G.Dist), 6)     // nverts
	chk.IntAssert(len(G.Next), 6)     // nverts

	shares := [][]int{
		{2, 3},    // edges sharing node 0
		{1, 2},    // edges sharing node 1
		{3, 4, 5}, // edges sharing node 2
		{5, 6},    // edges sharing node 3
		{0, 1},    // edges sharing node 4
		{0, 4, 6}, // edges sharing node 5
	}
	for k, share := range shares {
		chk.Ints(tst, io.Sf("edges sharing %d", k), G.Shares[k], share)
	}

	chk.IntAssert(G.Key2edge[G.HashEdgeKey(4, 5)], 0) // (4,5) → edge 0
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(1, 4)], 1) // (1,4) → edge 1
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(0, 1)], 2) // (0,1) → edge 2
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(0, 2)], 3) // (0,2) → edge 3
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(5, 2)], 4) // (5,2) → edge 4
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(2, 3)], 5) // (2,3) → edge 5
	chk.IntAssert(G.Key2edge[G.HashEdgeKey(5, 3)], 6) // (5,3) → edge 6

	err := G.ShortestPaths("FW")
	if err != nil {
		tst.Errorf("ShortestPaths failed:\n%v", err)
		return
	}
	inf := GRAPH_INF
	pth := G.Path(1, 3)
	io.Pforan("dist =\n%v", G.StrDistMatrix())
	io.Pforan("path from 1 to 3 = %v\n", pth)
	chk.Ints(tst, "1 → 3", pth, []int{1, 4, 5, 3})
	chk.Matrix(tst, "dist", 1e-17, G.Dist, [][]float64{
		{0, 6, 8, 17, 17, 20},
		{inf, 0, 21, 18, 11, 14},
		{inf, inf, 0, 9, inf, inf},
		{inf, inf, inf, 0, inf, inf},
		{inf, inf, 10, 7, 0, 3},
		{inf, inf, 7, 4, inf, 0},
	})
}

func Test_graph03(tst *testing.T) {

	verbose()
	chk.PrintTitle("graph03. Sioux Falls")

	//G := ReadGraphTable("data/SiouxFalls_net1.txt", true)
	G := ReadGraphTable("data/SiouxFalls_net.txt", true)

	err := G.ShortestPaths("FW")
	if err != nil {
		tst.Errorf("ShortestPaths failed:\n%v", err)
		return
	}

	pth := G.Path(0, 22)
	io.Pforan("1 → 23 = %v\n", pth)
	chk.Ints(tst, "1 → 23", pth, []int{0, 2, 11, 12, 23, 22})

	pth = G.Path(0, 20)
	io.Pforan("1 → 21 = %v\n", pth)
	chk.Ints(tst, "1 → 21", pth, []int{0, 2, 11, 12, 23, 20})

	pth = G.Path(2, 21)
	io.Pforan("3 → 22 = %v\n", pth)
	chk.Ints(tst, "3 → 22", pth, []int{2, 11, 12, 23, 22, 21})

	columns := [][]int{
		{1, 3, 12, 13},
		{4, 11, 14, 23, 24},
		{5, 9, 10, 15, 22, 21},
		{2, 6, 8, 16, 17, 19, 20},
		{7, 18},
	}
	Y := [][]float64{
		{7, 6, 4, 0},          // col0
		{6, 4, 2, 1, 0},       // col1
		{6, 5, 4, 2, 1, 0},    // col2
		{7, 6, 5, 4, 3, 2, 0}, // col3
		{5, 4},                // col4
	}

	r := 0.25
	W := 0.15
	dwt := 0.12
	aws := 12.0
	scalex := 1.8
	scaley := 1.3
	nv := 24
	G.Verts = make([][]float64, nv)
	for j, col := range columns {
		x := float64(j) * scalex
		for i, vidp1 := range col {
			vid := vidp1 - 1
			G.Verts[vid] = []float64{x, Y[j][i] * scaley}
		}
	}

	ne := 76
	elabels := make(map[int]string)
	for i := 0; i < ne; i++ {
		elabels[i] = io.Sf("%d", i+1)
	}

	vlabels := make(map[int]string)
	for i := 0; i < nv; i++ {
		vlabels[i] = io.Sf("%d", i+1)
	}

	vfsz := 7.0
	vclr := "red"
	efsz := 7.0
	eclr := "blue"
	plt.SetForEps(1.2, 350)
	if chk.Verbose {
		G.Draw("/tmp/graph", "siouxfalls.eps", r, W, dwt, aws, vlabels, vfsz, vclr, elabels, efsz, eclr)
	}
}
