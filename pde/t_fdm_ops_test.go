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
	"github.com/cpmech/gosl/utl"
)

func TestFdmLaplace01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FdmLaplace01. Full Auu matrix.")

	// grid
	g, err := gm.NewGrid([]float64{0, 0}, []float64{2, 2}, []int{2, 2}) // 2x2 divs ⇒ 3x3 grid ⇒ 9 equations
	status(tst, err)

	// equations
	e := la.NewEquations(g.N, nil)

	// operator
	op, err := NewFdmOperator("laplacian", dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}})
	status(tst, err)

	// assemble
	op.Assemble(g, e)
	Duu := e.Auu.ToDense()
	io.Pf("%v\n", Duu.Print("%4g"))

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

func TestFdmLaplace02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FdmLaplace02. Auu without borders")

	// grid
	g, err := gm.NewGrid([]float64{0, 0}, []float64{3, 3}, []int{3, 3}) // 3x3 divs ⇒ 4x4 grid ⇒ 16 equations
	status(tst, err)

	// equations
	e := la.NewEquations(g.N, utl.IntUnique(g.Edge...))

	// operator
	op, err := NewFdmOperator("laplacian", dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}})
	status(tst, err)

	// assemble
	op.Assemble(g, e)
	Duu := e.Auu.ToDense()
	io.Pf("%v\n", Duu.Print("%4g"))

	// check
	chk.Deep2(tst, "Auu", 1e-17, Duu.GetDeep2(), [][]float64{
		{+4, -1, -1, +0}, // 0 ⇒ node (1,1) 5
		{-1, +4, +0, -1}, // 1 ⇒ node (2,1) 6
		{-1, +0, +4, -1}, // 2 ⇒ node (1,2) 9
		{+0, -1, -1, +4}, // 3 ⇒ node (2,2) 10
	})

	// solve problem
	//    ∂²u     ∂²u
	//    ———  +  ——— = 0    with   u(x,0)=1   u(3,y)=2   u(x,3)=2   u(0,y)=1
	//    ∂x²     ∂y²               (bottom)   (right)    (top)      (left)

	// set BCS
	for _, I := range g.Edge[0] { // bottom
		e.Xk[e.FtoK[I]] = 1.0
	}
	for _, I := range g.Edge[1] { // right
		e.Xk[e.FtoK[I]] = 2.0
	}
	for _, I := range g.Edge[2] { // top
		e.Xk[e.FtoK[I]] = 2.0
	}
	for _, I := range g.Edge[3] { // left
		e.Xk[e.FtoK[I]] = 1.0
	}

	// fix RHS: bu -= Auk⋅xk
	la.SpMatVecMulAdd(e.Bu, -1.0, e.Auk.ToMatrix(nil), e.Xk)

	// solve system
	err = la.SolveRealLinSysSPD(e.Xu, Duu, e.Bu)
	status(tst, err)

	// joint parts
	x := la.NewVector(g.N)
	e.JoinVector(x, e.Xu, e.Xk)

	// check
	io.Pf("x = %v\n", x)
	chk.Array(tst, "x", 1e-15, x, []float64{1, 1, 1, 2, 1, 1.25, 1.5, 2, 1, 1.5, 1.75, 2, 1, 2, 2, 2})
}
