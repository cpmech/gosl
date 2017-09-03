// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/la"
)

func TestFdm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Fdm01. Auu without borders (FdmSolver)")

	// problem data
	params := dbf.Params{{N: "kx", V: 1}, {N: "ky", V: 1}}
	xmin := []float64{0, 0}
	xmax := []float64{3, 3}
	ndiv := []int{3, 3} // 3x3 divs ⇒ 4x4 grid ⇒ 16 equations

	// fdm solver
	fdm, err := NewFdmSolver("laplacian", params, xmin, xmax, ndiv)
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

	// solve problem
	err = fdm.Solve(true)
	status(tst, err)
	chk.Array(tst, "Xk", 1e-17, fdm.Eqs.Xk, []float64{1, 1, 1, 1, 1, 2, 1, 2, 2, 2, 2, 2})
	chk.Array(tst, "U", 1e-15, fdm.U, []float64{1, 1, 1, 1, 1, 1.25, 1.5, 2, 1, 1.5, 1.75, 2, 2, 2, 2, 2})

	// check
	eqFull, err := la.NewEquations(fdm.Grid.Size(), nil)
	status(tst, err)
	fdm.Op.Assemble(fdm.Grid, eqFull)
	K := eqFull.Auu.ToMatrix(nil)
	Fref := la.NewVector(fdm.Eqs.N)
	la.SpMatVecMul(Fref, 1.0, K, fdm.U)
	chk.Array(tst, "F", 1e-15, fdm.F, Fref)
}
