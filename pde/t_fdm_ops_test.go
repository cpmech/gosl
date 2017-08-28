// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func TestFdmLaplace01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FdmLaplace01")

	// grid
	g := new(gm.Grid)
	g.Init([]float64{0, 0}, []float64{2, 2}, []int{2, 2}) // 2x2 divs ⇒ 3x3 grid ⇒ 9 equations

	// equations
	e := new(la.Equations)
	e.Init(g.N, nil)

	// operator
	op := NewFdmOperator("laplacian")
	err := op.Init(dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}})
	status(tst, err)

	// assemble
	op.Assemble(g, e)
	Duu := e.Auu.ToDense()
	io.Pf("%v\n", Duu.Print("%+5g"))

	// check
	chk.Deep2(tst, "Auu", 1e-17, Duu.GetDeep2(), [][]float64{
		{+4, -2, +0, -2, +0, +0, +0, +0, +0}, // 0
		{-1, +4, -1, +0, -2, +0, +0, +0, +0}, // 1
		{+0, -2, +4, +0, +0, -2, +0, +0, +0}, // 2
		{-1, +0, +0, +4, -2, +0, -1, +0, +0}, // 3
		{+0, -1, +0, -1, +4, -1, +0, -1, +0}, // 4
		{+0, +0, -1, +0, -2, +4, +0, +0, -1}, // 5
		{+0, +0, +0, -2, +0, +0, +4, -2, +0}, // 6
		{+0, +0, +0, +0, -2, +0, -1, +4, -1}, // 7
		{+0, +0, +0, +0, +0, -2, +0, -2, +4}, // 8
	})
}
