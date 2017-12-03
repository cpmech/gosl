// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/gm/tri"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// vertices (points)
	V := [][]float64{
		{0, 0}, {1, 0},
		{1, 1}, {0, 1},
	}

	// cells (triangles)
	C := [][]int{
		{0, 1, 2},
		{2, 3, 0},
	}

	// plot
	plt.Reset(true, &plt.A{WidthPt: 300})
	tri.DrawVC(V, C, &plt.A{C: "#376ac6", Lw: 2, NoClip: true})
	for i, v := range V {
		plt.Text(v[0], v[1], io.Sf("(%d)", i), &plt.A{C: "r", Fsz: 12, NoClip: true})
	}
	plt.Gll("x", "y", nil)
	plt.Equal()
	plt.HideAllBorders()
	plt.Save("/tmp/gosl", "tri_draw01")
}
