// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
)

func TestEssentialBcs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("EssentialBcs01")

	// grid      6        7        8
	//          o--------o--------o
	//          |        |        |
	//          |        |        |
	//          |3       |4       |5
	//          o--------o--------o
	//          |        |        |
	//          |        |        |
	//          |0       |1       |2
	//          o--------o--------o
	g := new(gm.Grid)
	err := g.GenUniform([]float64{0, 0}, []float64{2, 2}, []int{2, 2}, false) // 2x2 divs ⇒ 3x3 grid ⇒ 9 equations
	status(tst, err)

	// essential boundary conditions
	ebcs := NewEssentialBcs()
	err = ebcs.SetInGrid(g, 10, "u", 123.0, nil) // left
	status(tst, err)
	err = ebcs.SetInGrid(g, 20, "u", 123.0, nil) // bottom
	status(tst, err)

	// print bcs
	io.Pf("%v\n", ebcs.Print())

	// check
	chk.Ints(tst, "left boundary", ebcs.All[0].GetNodesSorted(), []int{0, 3, 6})
	chk.Ints(tst, "right boundary", ebcs.All[1].GetNodesSorted(), []int{0, 1, 2})
}
