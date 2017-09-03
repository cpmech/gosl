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
	chk.PrintTitle("Spc01. simple Dirichlet problem")

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
	spc, err := NewGridSolver("spc", "uni", "laplacian", params, xmin, xmax, ndiv)
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
		{-9.00, +2.25, +2.25, +0.00}, // 0 ⇒ node (1,1) 5
		{+2.25, -9.00, +0.00, +2.25}, // 1 ⇒ node (2,1) 6
		{+2.25, +0.00, -9.00, +2.25}, // 2 ⇒ node (1,2) 9
		{+0.00, +2.25, +2.25, -9.00}, // 3 ⇒ node (2,2) 10
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

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		spc.Grid.Draw(true, nil, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.HideAllBorders()
		err = plt.Save("/tmp/gosl/pde", "spc01")
		status(tst, err)
	}
}
