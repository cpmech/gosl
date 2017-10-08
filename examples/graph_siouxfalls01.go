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

	// load graph data from FLOW file
	g := graph.ReadGraphTable("../graph/data/SiouxFalls.flow", false)

	// compute paths
	g.ShortestPaths("FW")

	// print shortest path from 0 to 20
	io.Pf("dist from = %v\n", g.Path(0, 20))
	io.Pf("must be: [0, 2, 11, 12, 23, 20]\n")

	// data for drawing: ids of vertices along columns in plot grid
	columns := [][]int{
		{0, 2, 11, 12},
		{3, 10, 13, 22, 23},
		{4, 8, 9, 14, 21, 20},
		{1, 5, 7, 15, 16, 18, 19},
		{6, 17},
	}

	// data for drawing: y-coordinates of vertices in plot
	Y := [][]float64{
		{7, 6, 4, 0},          // col0
		{6, 4, 2, 1, 0},       // col1
		{6, 5, 4, 2, 1, 0},    // col2
		{7, 6, 5, 4, 3, 2, 0}, // col3
		{5, 4},                // col4
	}

	// data for drawing: set vertices in graph structure
	scalex := 1.8
	scaley := 1.3
	nv := 24
	g.Verts = make([][]float64, nv)
	for j, col := range columns {
		x := float64(j) * scalex
		for i, vid := range col {
			g.Verts[vid] = []float64{x, Y[j][i] * scaley}
		}
	}

	// data for drawing: set vertex labels
	vlabels := make(map[int]string)
	for i := 0; i < nv; i++ {
		vlabels[i] = io.Sf("%d", i)
	}

	// data for drawing: set edge labels
	ne := 76
	elabels := make(map[int]string)
	for i := 0; i < ne; i++ {
		elabels[i] = io.Sf("%d", i)
	}

	// plot
	plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150, Prop: 1.1})
	g.Draw(vlabels, elabels, 0, 0, 0, nil, nil, nil, nil)
	plt.Equal()
	plt.AxisOff()
	plt.Save("/tmp/gosl/graph", "graph_siouxfalls01")
}
