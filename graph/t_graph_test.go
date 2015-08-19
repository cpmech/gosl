// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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

	var graph Graph
	graph.Init(
		// edge:  0       1       2       3
		[][]int{{0, 1}, {0, 3}, {1, 2}, {2, 3}},
		[]float64{5, 10, 3, 1}, // weights
		nil, nil,
	)

	chk.IntAssert(len(graph.Shares), 4)   // nverts
	chk.IntAssert(len(graph.Key2edge), 4) // nedges
	chk.IntAssert(len(graph.Dist), 4)     // nverts
	chk.IntAssert(len(graph.Next), 4)     // nverts

	shares := [][]int{
		{0, 1}, // edges sharing node 0
		{0, 2}, // edges sharing node 1
		{2, 3}, // edges sharing node 2
		{1, 3}, // edges sharing node 3
	}
	for k, share := range shares {
		chk.Ints(tst, io.Sf("edges sharing %d", k), graph.Shares[k], share)
	}

	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(0, 1)], 0) // (0,1) → edge 0
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(0, 3)], 1) // (0,3) → edge 1
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(1, 2)], 2) // (1,2) → edge 2
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(2, 3)], 3) // (2,3) → edge 3

	err := graph.ShortestPaths("FW")
	if err != nil {
		tst.Errorf("ShortestPaths failed:\n%v", err)
		return
	}
	inf := GRAPH_INF
	pth := graph.Path(0, 3)
	io.Pforan("dist =\n%v", graph.StrDistMatrix())
	io.Pforan("path from 0 to 3 = %v\n", pth)
	chk.Ints(tst, "0 → 3", pth, []int{0, 1, 2, 3})
	chk.Matrix(tst, "dist", 1e-17, graph.Dist, [][]float64{
		{0, 5, 8, 9},
		{inf, 0, 3, 4},
		{inf, inf, 0, 1},
		{inf, inf, inf, 0},
	})

	graph.WeightsE[3] = 13
	err = graph.ShortestPaths("FW")
	if err != nil {
		tst.Errorf("ShortestPaths failed:\n%v", err)
		return
	}
	pth = graph.Path(0, 3)
	io.Pf("\n")
	io.Pfcyan("dist =\n%v", graph.StrDistMatrix())
	io.Pfcyan("path from 0 to 3 = %v\n", pth)
	chk.Ints(tst, "0 → 3", pth, []int{0, 3})
	chk.Matrix(tst, "dist", 1e-17, graph.Dist, [][]float64{
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

	var graph Graph
	graph.Init(
		// edge:  0       1       2       3       4       5       6
		[][]int{{4, 5}, {1, 4}, {0, 1}, {0, 2}, {5, 2}, {2, 3}, {5, 3}},
		[]float64{3, 11, 6, 8, 7, 9, 4},
		nil, nil,
	)

	chk.IntAssert(len(graph.Shares), 6)   // nverts
	chk.IntAssert(len(graph.Key2edge), 7) // nedges
	chk.IntAssert(len(graph.Dist), 6)     // nverts
	chk.IntAssert(len(graph.Next), 6)     // nverts

	shares := [][]int{
		{2, 3},    // edges sharing node 0
		{1, 2},    // edges sharing node 1
		{3, 4, 5}, // edges sharing node 2
		{5, 6},    // edges sharing node 3
		{0, 1},    // edges sharing node 4
		{0, 4, 6}, // edges sharing node 5
	}
	for k, share := range shares {
		chk.Ints(tst, io.Sf("edges sharing %d", k), graph.Shares[k], share)
	}

	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(4, 5)], 0) // (4,5) → edge 0
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(1, 4)], 1) // (1,4) → edge 1
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(0, 1)], 2) // (0,1) → edge 2
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(0, 2)], 3) // (0,2) → edge 3
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(5, 2)], 4) // (5,2) → edge 4
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(2, 3)], 5) // (2,3) → edge 5
	chk.IntAssert(graph.Key2edge[graph.HashEdgeKey(5, 3)], 6) // (5,3) → edge 6

	err := graph.ShortestPaths("FW")
	if err != nil {
		tst.Errorf("ShortestPaths failed:\n%v", err)
		return
	}
	inf := GRAPH_INF
	pth := graph.Path(1, 3)
	io.Pforan("dist =\n%v", graph.StrDistMatrix())
	io.Pforan("path from 1 to 3 = %v\n", pth)
	chk.Ints(tst, "1 → 3", pth, []int{1, 4, 5, 3})
	chk.Matrix(tst, "dist", 1e-17, graph.Dist, [][]float64{
		{0, 6, 8, 17, 17, 20},
		{inf, 0, 21, 18, 11, 14},
		{inf, inf, 0, 9, inf, inf},
		{inf, inf, inf, 0, inf, inf},
		{inf, inf, 10, 7, 0, 3},
		{inf, inf, 7, 4, inf, 0},
	})
}
