// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestFdm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Fdm01. Full Auu matrix.")

	// operator
	op, err := NewOperator("fdm.laplacian", dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}})
	status(tst, err)

	// init operator with grid: 2x2 divs ⇒ 3x3 grid ⇒ 9 equations
	g, err := op.InitWithGrid("uni", []float64{0, 0}, []float64{2, 2}, []int{2, 2})
	status(tst, err)

	// equations
	e, err := la.NewEquations(g.Size(), nil)
	status(tst, err)

	// assemble
	op.Assemble(e)
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

func TestFdm02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Fdm02. simple Dirichlet problem (Laplace)")

	// solve problem
	//    ∂²u     ∂²u
	//    ———  +  ——— = 0    with   u(x,0)=1   u(3,y)=2   u(x,3)=2   u(0,y)=1
	//    ∂x²     ∂y²               (bottom)   (right)    (top)      (left)

	// problem data
	params := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	xmin := []float64{0, 0}
	xmax := []float64{3, 3}
	ndiv := []int{3, 3} // 3x3 divs ⇒ 4x4 grid ⇒ 16 equations

	// fdm solver
	fdm, err := NewGridSolver("fdm", "uni", "laplacian", params, xmin, xmax, ndiv)
	status(tst, err)

	// essential boundary conditions
	ebcs := NewEssentialBcs()
	L, R, B, T := 10, 11, 20, 21 // left, right, bottom, top
	ebcs.SetInGrid(fdm.Grid, L, "u", 1.0, nil)
	ebcs.SetInGrid(fdm.Grid, R, "u", 2.0, nil)
	ebcs.SetInGrid(fdm.Grid, B, "u", 1.0, nil)
	ebcs.SetInGrid(fdm.Grid, T, "u", 2.0, nil)

	// set bcs
	fdm.SetBcs(ebcs)
	chk.Ints(tst, "UtoF", fdm.Eqs.UtoF, []int{5, 6, 9, 10})
	chk.Ints(tst, "KtoF", fdm.Eqs.KtoF, []int{0, 1, 2, 3, 4, 7, 8, 11, 12, 13, 14, 15})

	// check
	Duu := fdm.Eqs.Auu.ToDense()
	io.Pf("Auu =\n%v\n", Duu.Print("%4g"))
	chk.Deep2(tst, "Auu", 1e-17, Duu.GetDeep2(), [][]float64{
		{+4, -1, -1, +0}, // 0 ⇒ node (1,1) 5
		{-1, +4, +0, -1}, // 1 ⇒ node (2,1) 6
		{-1, +0, +4, -1}, // 2 ⇒ node (1,2) 9
		{+0, -1, -1, +4}, // 3 ⇒ node (2,2) 10
	})

	// solve problem
	err = fdm.Solve(true)
	status(tst, err)
	chk.Array(tst, "Xk", 1e-17, fdm.Eqs.Xk, []float64{1, 1, 1, 1, 1, 2, 1, 2, 2, 2, 2, 2})
	chk.Array(tst, "U", 1e-15, fdm.U, []float64{1, 1, 1, 1, 1, 1.25, 1.5, 2, 1, 1.5, 1.75, 2, 2, 2, 2, 2})

	// check
	eqFull, err := la.NewEquations(fdm.Grid.Size(), nil)
	status(tst, err)
	fdm.Op.Assemble(eqFull)
	K := eqFull.Auu.ToMatrix(nil)
	Fref := la.NewVector(fdm.Eqs.N)
	la.SpMatVecMul(Fref, 1.0, K, fdm.U)
	chk.Array(tst, "F", 1e-15, fdm.F, Fref)

	// get results over grid
	uu := fdm.Ugrid2d()
	chk.Deep2(tst, "uu", 1e-15, uu, [][]float64{
		{1.00, 1.00, 1.00, 1.00},
		{1.00, 1.25, 1.50, 2.00},
		{1.00, 1.50, 1.75, 2.00},
		{2.00, 2.00, 2.00, 2.00},
	})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		fdm.Grid.Draw(true, nil, &plt.A{C: plt.C(1, 0), Fsz: 7})
		xx, yy := fdm.Grid.Mesh2d()
		plt.ContourL(xx, yy, uu, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.HideAllBorders()
		err = plt.Save("/tmp/gosl/pde", "fdm02")
		status(tst, err)
	}
}
