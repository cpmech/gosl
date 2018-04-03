// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/gm/tri"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// set Planar Straight Line Graph (PSLG)
	p := &tri.Input{
		[]*tri.Point{
			{11, 0.0, 0.0},
			{22, 1.0, 0.0},
			{33, 1.0, 1.0},
			{44, 0.0, 1.0},
		},
		[]*tri.Segment{
			{100, 0, 1},
			{200, 1, 2},
			{300, 2, 3},
			{400, 3, 0},
		},
		[]*tri.Region{
			{10, 1.0, 0.5, 0.5},
		},
		[]*tri.Hole{},
	}

	// configuration
	globalMaxArea := 0.1
	globalMinAngle := 20.0

	// generate
	m := p.Generate(globalMaxArea, globalMinAngle, false, true, "")

	// plot
	plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
	tri.DrawMesh(m, true, nil, nil, nil, nil)
	plt.Gll("x", "y", nil)
	plt.Equal()
	plt.HideAllBorders()
	plt.Save("/tmp/gosl", "tri_generate01")
}
