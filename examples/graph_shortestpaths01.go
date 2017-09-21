// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/graph"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	//           [10]
	//      0 ––––––––→ 3      numbers in parentheses
	//      |    (1)    ↑      indicate edge ids
	//   [5]|(0)        |
	//      |        (3)|[1]
	//      ↓    (2)    |      numbers in brackets
	//      1 ––––––––→ 2      indicate weights
	//           [3]

	// initialise graph
	var g graph.Graph
	g.Init(
		// edge:  0       1       2       3
		[][]int{{0, 1}, {0, 3}, {1, 2}, {2, 3}},

		// weights:
		[]float64{5, 10, 3, 1},

		// vertices (coordinates, for drawings):
		[][]float64{
			{0, 1}, // x,y vertex 0
			{0, 0}, // x,y vertex 1
			{1, 0}, // x,y vertex 2
			{1, 1}, // x,y vertex 3
		},

		// weights:
		nil,
	)

	// compute paths
	g.ShortestPaths("FW")

	// print shortest path from 0 to 2
	io.Pf("dist from = %v\n", g.Path(0, 2))

	// print shortest path from 0 to 3
	io.Pf("dist from = %v\n", g.Path(0, 3))

	// print distance matrix
	io.Pf("dist =\n%v", g.StrDistMatrix())

	// constants for plot
	radius, width, gap := 0.05, 1e-8, 0.05

	// plot
	plt.Reset(true, &plt.A{WidthPt: 250, Dpi: 150, Prop: 1.0})
	g.Draw(nil, nil, radius, width, gap, nil, nil, nil, nil)
	plt.Equal()
	plt.AxisOff()
	plt.Save("/tmp/gosl/graph", "graph_shortestpath01")
}
