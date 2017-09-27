// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package graph

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func TestAdjacency01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Adjacency01")

	// new graph
	g := new(Graph)
	g.Init([][]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 4},
		{5, 6}, {6, 7}, {7, 8}, {8, 9},
		{10, 11}, {11, 12}, {12, 13}, {13, 14},
		{0, 5}, {1, 6}, {2, 7}, {3, 8}, {4, 9},
		{5, 10}, {6, 11}, {7, 12}, {8, 13}, {9, 14},
	}, nil, nil, nil)

	// adjacency list
	xadj, adjncy := g.GetAdjacency()
	chk.Int32s(tst, "xadj", xadj, []int32{0, 2, 5, 8, 11, 13, 16, 20, 24, 28, 31, 33, 36, 39, 42, 44})
	chk.Int32s(tst, "adjncy", adjncy, []int32{1, 5, 0, 2, 6, 1, 3, 7, 2, 4, 8, 3, 9, 6, 0, 10, 5, 7, 1, 11, 6, 8, 2, 12, 7, 9, 3, 13, 8, 4, 14, 11, 5, 10, 12, 6, 11, 13, 7, 12, 14, 8, 13, 9})

	// partition graph
	nv := g.Nverts()
	npart := 2
	recursive := false
	objval, parts := MetisPartition(npart, nv, xadj, adjncy, recursive)
	io.Pf("objval = %v\n", objval)
	io.Pf("parts = %v\n", parts)
	chk.Int32(tst, "objval", objval, 7)
	chk.Int32s(tst, "parts", parts, []int32{0, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1})

	// plot
	if chk.Verbose {
		g.Verts = [][]float64{
			{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2},
			{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
			{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0},
		}
		p := Plotter{G: g, Parts: parts}
		plt.Reset(true, &plt.A{Prop: 0.7})
		p.Draw()
		plt.AxisOff()
		plt.Equal()
		plt.AxisRange(-0.5, 4.5, -1, 3)
		plt.Save("/tmp/gosl/graph", "adjacency01")
	}
}
