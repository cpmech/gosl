// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/utl"
)

func Test_mat01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("mat01. Matrix functions")

	// MatAlloc
	utl.Pfblue2("\nfunc MatAlloc(m, n int) (mat [][]float64)\n")
	a := MatAlloc(3, 5)
	a[0][0] = 1
	a[0][1] = 2
	a[0][2] = 3
	a[0][3] = 4
	a[0][4] = 5
	a[1][0] = 0.1
	a[1][1] = 0.2
	a[1][2] = 0.3
	a[1][3] = 0.4
	a[1][4] = 0.5
	a[2][0] = 10
	a[2][1] = 20
	a[2][2] = 30
	a[2][3] = 40
	a[2][4] = 50
	PrintMat("a", a, "%5g", false)
	utl.CheckMatrix(tst, "a", 1e-17, a, [][]float64{
		{1, 2, 3, 4, 5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{10, 20, 30, 40, 50},
	})

	aclone := MatClone(a)
	utl.CheckMatrix(tst, "aclone", 1e-17, a, aclone)

	// MatFill
	utl.Pfblue2("\nfunc MatFill(a [][]float64, s float64)\n")
	b := MatAlloc(5, 3)
	MatFill(b, 2)
	PrintMat("b", b, "%5g", false)
	utl.CheckMatrix(tst, "b", 1e-17, b, [][]float64{
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
	})

	// MatScale
	utl.Pfblue2("\nfunc MatScale(a [][]float64, alp float64)\n")
	c := MatAlloc(5, 3)
	MatFill(c, 2)
	MatScale(c, 1.0/4.0)
	PrintMat("c", c, "%5g", false)
	utl.CheckMatrix(tst, "c", 1e-17, c, [][]float64{
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
	})

	// MatCopy
	utl.Pfblue2("\nfunc MatCopy(a [][]float64, alp float64, b [][]float64)\n")
	d := MatAlloc(3, 5)
	MatCopy(d, 1, a)
	PrintMat("d", d, "%5g", false)
	utl.CheckMatrix(tst, "d", 1e-17, d, [][]float64{
		{1, 2, 3, 4, 5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{10, 20, 30, 40, 50},
	})

	// MatSetDiag
	utl.Pfblue2("\nfunc MatSetDiag(a [][]float64, s float64)\n")
	e := MatAlloc(3, 3)
	MatSetDiag(e, 1)
	PrintMat("e", e, "%5g", false)
	utl.CheckMatrix(tst, "e", 1e-17, e, [][]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	})

	// MatMaxDiff
	utl.Pfblue2("\nfunc MatMaxDiff(a, b [][]float64) (maxdiff float64)\n")
	f := [][]float64{
		{1.1, 2.2, 3.3, 4.4, 5.5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{1, 2, 3, 4, 5},
	}
	fma := [][]float64{
		{f[0][0] - a[0][0], f[0][1] - a[0][1], f[0][2] - a[0][2], f[0][3] - a[0][3], f[0][4] - a[0][4]},
		{f[1][0] - a[1][0], f[1][1] - a[1][1], f[1][2] - a[1][2], f[1][3] - a[1][3], f[1][4] - a[1][4]},
		{f[2][0] - a[2][0], f[2][1] - a[2][1], f[2][2] - a[2][2], f[2][3] - a[2][3], f[2][4] - a[2][4]},
	}
	PrintMat("a", a, "%5g", false)
	PrintMat("f", f, "%5g", false)
	PrintMat("f-a", fma, "%5g", false)
	maxdiff := MatMaxDiff(f, a)
	utl.Pf("max(f-a) = %v\n", maxdiff)
	utl.CheckScalar(tst, "maxdiff", 1e-17, maxdiff, 45)

	// MatGetCol
	utl.Pfblue2("\nfunc MatGetCol(j int, a [][]float64) (col []float64)\n")
	col := MatGetCol(0, a)
	PrintMat("a", a, "%5g", false)
	PrintVec("a[:][0]", col, "%5g", false)
	utl.CheckVector(tst, "a[:][0]", 1e-17, col, []float64{1, 0.1, 10})

	// MatNormF
	utl.Pfblue2("\nfunc MatNormF(a [][]float64) (res float64)\n")
	Pll = true
	NCPU = 3
	A := [][]float64{
		{-3, 5, 7},
		{2, 6, 4},
		{0, 2, 8},
	}
	PrintMat("A", A, "%5g", false)
	normFA := MatNormF(A)
	utl.Pf("normF(A) = %g\n", normFA)
	utl.CheckScalar(tst, "normF(A)", 1e-17, normFA, 1.438749456993816e+01)

	// MatNormI
	utl.Pfblue2("\nfunc MatNormI(a [][]float64) (res float64)\n")
	PrintMat("A", A, "%5g", false)
	normIA := MatNormI(A)
	utl.Pf("normI(A) = %g\n", normIA)
	utl.CheckScalar(tst, "normI(A)", 1e-17, normIA, 15.0)

	// MatInv
	utl.Pfblue2("\nfunc MatInv(ai, a [][]float64, tol float64) (det float64, ok bool)\n")
	g := [][]float64{
		{1, 2, 3},
		{0, 4, 5},
		{1, 0, 6},
	}
	gi := MatAlloc(3, 3)
	detg, err := MatInv(gi, g, 1e-17)
	if err != nil {
		utl.Panic("%v", err.Error())
	}
	gi22 := [][]float64{
		{gi[0][0] * 22.0, gi[0][1] * 22.0, gi[0][2] * 22.0},
		{gi[1][0] * 22.0, gi[1][1] * 22.0, gi[1][2] * 22.0},
		{gi[2][0] * 22.0, gi[2][1] * 22.0, gi[2][2] * 22.0},
	}
	PrintMat("g", g, "%5g", false)
	PrintMat("gi * 22", gi22, "%5g", false)
	utl.Pf("det(g) = %g\n", detg)
	utl.CheckScalar(tst, "det(g)", 1e-17, detg, 22)
	utl.CheckMatrix(tst, "gi", 1e-17, gi, [][]float64{
		{12.0 / 11.0, -6.0 / 11.0, -1.0 / 11.0},
		{5.0 / 22.0, 3.0 / 22.0, -5.0 / 22.0},
		{-2.0 / 11.0, 1.0 / 11.0, 2.0 / 11.0},
	})
}
