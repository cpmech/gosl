// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tri

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/plt"
)

func TestGenerate01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Generate01")

	l := 1.0
	h := math.Sqrt(3.0) * l / 2.0
	A := l * h / 2.0

	p := &Input{
		[]*Point{
			{11, 0.0, 0.0},
			{22, l, 0.0},
			{33, l / 2.0, h},
		},
		[]*Segment{
			{100, 0, 1},
			{200, 1, 2},
			{300, 2, 0},
		},
		[]*Region{
			{10, 1.0, l / 2.0, h / 2.0},
		},
		[]*Hole{},
	}

	globalMaxArea := A / 4.0
	globalMinAngle := 45.0
	o2 := true
	m := p.Generate(globalMaxArea, globalMinAngle, o2, chk.Verbose, "")

	io.Pf("\nVertices:\n")
	for _, v := range m.Verts {
		io.Pf("%+v\n", v)
	}

	io.Pf("\nCells:\n")
	for _, c := range m.Cells {
		io.Pf("%+v\n", c)
	}

	chk.Int(tst, "vert0: tag", m.Verts[0].Tag, 11)
	chk.Int(tst, "vert1: tag", m.Verts[1].Tag, 22)
	chk.Int(tst, "vert2: tag", m.Verts[2].Tag, 33)
	chk.Ints(tst, "cell0: verts", m.Cells[0].V, []int{4, 2, 5, 6, 7, 8})
	chk.Ints(tst, "cell0: etags", m.Cells[0].EdgeTags, []int{200, 300, 0})
	chk.Ints(tst, "cell1: verts", m.Cells[1].V, []int{4, 5, 3, 8, 9, 10})
	chk.Ints(tst, "cell1: etags", m.Cells[1].EdgeTags, []int{0, 0, 0})
	chk.Ints(tst, "cell2: verts", m.Cells[2].V, []int{5, 0, 3, 11, 12, 9})
	chk.Ints(tst, "cell2: etags", m.Cells[2].EdgeTags, []int{300, 100, 0})
	chk.Ints(tst, "cell3: verts", m.Cells[3].V, []int{3, 1, 4, 13, 14, 10})
	chk.Ints(tst, "cell3: etags", m.Cells[3].EdgeTags, []int{100, 200, 0})

	if chk.Verbose {
		io.Pl()
		plt.Reset(true, nil)
		DrawMesh(m, true, nil, nil, nil, nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/tri", "generate01")
	}
}
