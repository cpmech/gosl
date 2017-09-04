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

func TestSpc01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Spc01. Auu matrix after homogeneous bcs")

	// problem data
	params := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	xmin := []float64{0, 0}
	xmax := []float64{4, 4}
	ndiv := []int{4, 4} // 4x4 divs ⇒ 5x5 grid ⇒ 25 equations

	// spectral-collocation solver
	spc, err := NewGridSolver("spc", "cgl", "laplacian", params, nil, xmin, xmax, ndiv)
	status(tst, err)

	// essential boundary conditions
	ebcs := NewEssentialBcs()
	L, R, B, T := 10, 11, 20, 21 // left, right, bottom, top
	ebcs.SetInGrid(spc.Grid, L, "u", 0.0, nil)
	ebcs.SetInGrid(spc.Grid, R, "u", 0.0, nil)
	ebcs.SetInGrid(spc.Grid, B, "u", 0.0, nil)
	ebcs.SetInGrid(spc.Grid, T, "u", 0.0, nil)

	// set bcs
	spc.SetBcs(ebcs)

	// check
	Duu := spc.Eqs.Auu.ToDense()
	io.Pf("Auu =\n%v\n", Duu.Print("%10.5f"))
	chk.Deep2(tst, "Auu", 1e-14, Duu.GetDeep2(), [][]float64{
		{+28, -6., +2., -6., +0., +0., +2., +0., +0.},
		{-4., +20, -4., +0., -6., +0., +0., +2., +0.},
		{+2., -6., +28, +0., +0., -6., +0., +0., +2.},
		{-4., +0., +0., +20, -6., +2., -4., +0., +0.},
		{+0., -4., +0., -4., +12, -4., +0., -4., +0.},
		{+0., +0., -4., +2., -6., +20, +0., +0., -4.},
		{+2., +0., +0., -6., +0., +0., +28, -6., +2.},
		{+0., +2., +0., +0., -6., +0., -4., +20, -4.},
		{+0., +0., +2., +0., +0., -6., +2., -6., +28},
	})
	err = spc.Solve(true)
	status(tst, err)
}

func TestSpc02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Spc02. simple Dirichlet problem (unif-grid / Laplace)")

	// solve problem
	//    ∂²u     ∂²u
	//    ———  +  ——— = 0    with   u(x,0)=1   u(3,y)=2   u(x,3)=2   u(0,y)=1
	//    ∂x²     ∂y²               (bottom)   (right)    (top)      (left)

	// problem data
	params := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	xmin := []float64{0, 0}
	xmax := []float64{3, 3}
	ndiv := []int{3, 3} // 3x3 divs ⇒ 4x4 grid ⇒ 16 equations

	// spectral-collocation solver
	spc, err := NewGridSolver("spc", "uni", "laplacian", params, nil, xmin, xmax, ndiv)
	status(tst, err)

	// essential boundary conditions
	ebcs := NewEssentialBcs()
	L, R, B, T := 10, 11, 20, 21 // left, right, bottom, top
	ebcs.SetInGrid(spc.Grid, L, "u", 1.0, nil)
	ebcs.SetInGrid(spc.Grid, R, "u", 2.0, nil)
	ebcs.SetInGrid(spc.Grid, B, "u", 1.0, nil)
	ebcs.SetInGrid(spc.Grid, T, "u", 2.0, nil)

	// set bcs
	spc.SetBcs(ebcs)
	chk.Ints(tst, "UtoF", spc.Eqs.UtoF, []int{5, 6, 9, 10})
	chk.Ints(tst, "KtoF", spc.Eqs.KtoF, []int{0, 1, 2, 3, 4, 7, 8, 11, 12, 13, 14, 15})

	// check
	Duu := spc.Eqs.Auu.ToDense()
	io.Pf("Auu =\n%v\n", Duu.Print("%6g"))
	chk.Deep2(tst, "Auu", 1e-17, Duu.GetDeep2(), [][]float64{
		{+9.00, -2.25, -2.25, +0.00}, // 0 ⇒ node (1,1) 5
		{-2.25, +9.00, +0.00, -2.25}, // 1 ⇒ node (2,1) 6
		{-2.25, +0.00, +9.00, -2.25}, // 2 ⇒ node (1,2) 9
		{+0.00, -2.25, -2.25, +9.00}, // 3 ⇒ node (2,2) 10
	})

	// solve problem
	err = spc.Solve(true)
	status(tst, err)
	chk.Array(tst, "Xk", 1e-17, spc.Eqs.Xk, []float64{1, 1, 1, 1, 1, 2, 1, 2, 2, 2, 2, 2})
	chk.Array(tst, "U", 1e-15, spc.U, []float64{1, 1, 1, 1, 1, 1.25, 1.5, 2, 1, 1.5, 1.75, 2, 2, 2, 2, 2})

	// check
	eqFull, err := la.NewEquations(spc.Grid.Size(), nil)
	status(tst, err)
	spc.Op.Assemble(eqFull)
	K := eqFull.Auu.ToMatrix(nil)
	Fref := la.NewVector(spc.Eqs.N)
	la.SpMatVecMul(Fref, 1.0, K, spc.U)
	chk.Array(tst, "F", 1e-14, spc.F, Fref)

	// get results over grid
	uu := spc.Ugrid2d()
	chk.Deep2(tst, "uu", 1e-15, uu, [][]float64{
		{1.00, 1.00, 1.00, 1.00},
		{1.00, 1.25, 1.50, 2.00},
		{1.00, 1.50, 1.75, 2.00},
		{2.00, 2.00, 2.00, 2.00},
	})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		spc.Grid.Draw(true, nil, &plt.A{C: plt.C(1, 0), Fsz: 7})
		xx, yy := spc.Grid.Mesh2d()
		plt.ContourL(xx, yy, uu, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.HideAllBorders()
		err = plt.Save("/tmp/gosl/pde", "spc02")
		status(tst, err)
	}
}
