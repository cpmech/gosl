// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestMatrix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix01. real")

	A := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 0, -1, -2},
	}

	a := NewMatrix(A)
	chk.Vector(tst, "A to a", 1e-15, a.Data, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})

	chk.Scalar(tst, "Get(0,0)", 1e-17, a.Get(0, 0), 1)
	chk.Scalar(tst, "Get(0,1)", 1e-17, a.Get(0, 1), 2)
	chk.Scalar(tst, "Get(0,2)", 1e-17, a.Get(0, 2), 3)
	chk.Scalar(tst, "Get(0,3)", 1e-17, a.Get(0, 3), 4)

	chk.Scalar(tst, "Get(1,0)", 1e-17, a.Get(1, 0), 5)
	chk.Scalar(tst, "Get(1,1)", 1e-17, a.Get(1, 1), 6)
	chk.Scalar(tst, "Get(1,2)", 1e-17, a.Get(1, 2), 7)
	chk.Scalar(tst, "Get(1,3)", 1e-17, a.Get(1, 3), 8)

	chk.Scalar(tst, "Get(2,0)", 1e-17, a.Get(2, 0), 9)
	chk.Scalar(tst, "Get(2,1)", 1e-17, a.Get(2, 1), 0)
	chk.Scalar(tst, "Get(2,2)", 1e-17, a.Get(2, 2), -1)
	chk.Scalar(tst, "Get(2,3)", 1e-17, a.Get(2, 3), -2)

	Aback := a.GetSlice()
	chk.Matrix(tst, "a to A", 1e-15, Aback, A)

	l := a.Print("")
	chk.String(tst, l, "1 2 3 4 \n5 6 7 8 \n9 0 -1 -2 ")

	l = a.PrintGo("%2g")
	lCorrect := "[][]float64{\n    { 1, 2, 3, 4},\n    { 5, 6, 7, 8},\n    { 9, 0,-1,-2},\n}"
	chk.String(tst, l, lCorrect)

	l = a.PrintPy("%2g")
	lCorrect = "np.matrix([\n    [ 1, 2, 3, 4],\n    [ 5, 6, 7, 8],\n    [ 9, 0,-1,-2],\n], dtype=float)"
	chk.String(tst, l, lCorrect)

	a.Add(1, 2, -7)
	chk.Scalar(tst, "Get(1,2)", 1e-17, a.Get(1, 2), 0)

	b := a.GetCopy()
	if b.M != a.M {
		tst.Errorf("b.M should be equal to a.M\n")
		return
	}
	if b.N != a.N {
		tst.Errorf("b.N should be equal to a.N\n")
		return
	}
	chk.Vector(tst, "b.Data", 1e-17, b.Data, a.Data)
}

func TestMatrix02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix02. complex")

	A := [][]complex128{
		{1 + 0.1i, 2, 3, 4 - 0.4i},
		{5 + 0.5i, 6, 7, 8 - 0.8i},
		{9 + 0.9i, 0, -1, -2 + 1i},
	}

	a := NewMatrixC(A)
	chk.VectorC(tst, "A to a", 1e-15, a.Data, []complex128{1 + 0.1i, 5 + 0.5i, 9 + 0.9i, 2, 6, 0, 3, 7, -1, 4 - 0.4i, 8 - 0.8i, -2 + 1i})

	chk.ScalarC(tst, "Get(0,0)", 1e-17, a.Get(0, 0), 1+0.1i)
	chk.ScalarC(tst, "Get(0,1)", 1e-17, a.Get(0, 1), 2)
	chk.ScalarC(tst, "Get(0,2)", 1e-17, a.Get(0, 2), 3)
	chk.ScalarC(tst, "Get(0,3)", 1e-17, a.Get(0, 3), 4-0.4i)

	chk.ScalarC(tst, "Get(1,0)", 1e-17, a.Get(1, 0), 5+0.5i)
	chk.ScalarC(tst, "Get(1,1)", 1e-17, a.Get(1, 1), 6)
	chk.ScalarC(tst, "Get(1,2)", 1e-17, a.Get(1, 2), 7)
	chk.ScalarC(tst, "Get(1,3)", 1e-17, a.Get(1, 3), 8-0.8i)

	chk.ScalarC(tst, "Get(2,0)", 1e-17, a.Get(2, 0), 9+0.9i)
	chk.ScalarC(tst, "Get(2,1)", 1e-17, a.Get(2, 1), 0)
	chk.ScalarC(tst, "Get(2,2)", 1e-17, a.Get(2, 2), -1)
	chk.ScalarC(tst, "Get(2,3)", 1e-17, a.Get(2, 3), -2+1i)

	Aback := a.GetSlice()
	chk.MatrixC(tst, "a to A", 1e-15, Aback, A)

	l := a.Print("%g", "")
	chk.String(tst, l, "1+0.1i, 2+0i, 3+0i, 4-0.4i\n5+0.5i, 6+0i, 7+0i, 8-0.8i\n9+0.9i, 0+0i, -1+0i, -2+1i")

	l = a.PrintGo("%2g", "%+4.1f")
	lCorrect := "[][]complex128{\n    { 1+0.1i, 2+0.0i, 3+0.0i, 4-0.4i},\n    { 5+0.5i, 6+0.0i, 7+0.0i, 8-0.8i},\n    { 9+0.9i, 0+0.0i,-1+0.0i,-2+1.0i},\n}"
	chk.String(tst, l, lCorrect)

	l = a.PrintPy("%2g", "%4.1f")
	lCorrect = "np.matrix([\n    [ 1+0.1j, 2+0.0j, 3+0.0j, 4-0.4j],\n    [ 5+0.5j, 6+0.0j, 7+0.0j, 8-0.8j],\n    [ 9+0.9j, 0+0.0j,-1+0.0j,-2+1.0j],\n], dtype=complex)"
	chk.String(tst, l, lCorrect)

	a.Add(1, 3, -8+0.8i)
	chk.ScalarC(tst, "Get(1,3)", 1e-17, a.Get(1, 3), 0)

	b := a.GetCopy()
	if b.M != a.M {
		tst.Errorf("b.M should be equal to a.M\n")
		return
	}
	if b.N != a.N {
		tst.Errorf("b.N should be equal to a.N\n")
		return
	}
	chk.VectorC(tst, "b.Data", 1e-17, b.Data, a.Data)
}

func Test_mat01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mat01. Matrix functions")

	// MatAlloc
	io.Pfblue2("\nfunc MatAlloc(m, n int) (mat [][]float64)\n")
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
	chk.Matrix(tst, "a", 1e-17, a, [][]float64{
		{1, 2, 3, 4, 5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{10, 20, 30, 40, 50},
	})

	aclone := MatClone(a)
	chk.Matrix(tst, "aclone", 1e-17, a, aclone)

	// MatFill
	io.Pfblue2("\nfunc MatFill(a [][]float64, s float64)\n")
	b := MatAlloc(5, 3)
	MatFill(b, 2)
	PrintMat("b", b, "%5g", false)
	chk.Matrix(tst, "b", 1e-17, b, [][]float64{
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
	})

	// MatScale
	io.Pfblue2("\nfunc MatScale(a [][]float64, alp float64)\n")
	c := MatAlloc(5, 3)
	MatFill(c, 2)
	MatScale(c, 1.0/4.0)
	PrintMat("c", c, "%5g", false)
	chk.Matrix(tst, "c", 1e-17, c, [][]float64{
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5},
	})

	// MatCopy
	io.Pfblue2("\nfunc MatCopy(a [][]float64, alp float64, b [][]float64)\n")
	d := MatAlloc(3, 5)
	MatCopy(d, 1, a)
	PrintMat("d", d, "%5g", false)
	chk.Matrix(tst, "d", 1e-17, d, [][]float64{
		{1, 2, 3, 4, 5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{10, 20, 30, 40, 50},
	})

	// MatSetDiag
	io.Pfblue2("\nfunc MatSetDiag(a [][]float64, s float64)\n")
	e := MatAlloc(3, 3)
	MatSetDiag(e, 1)
	PrintMat("e", e, "%5g", false)
	chk.Matrix(tst, "e", 1e-17, e, [][]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	})

	// MatMaxDiff
	io.Pfblue2("\nfunc MatMaxDiff(a, b [][]float64) (maxdiff float64)\n")
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
	io.Pf("max(f-a) = %v\n", maxdiff)
	chk.Scalar(tst, "maxdiff", 1e-17, maxdiff, 45)

	// MatLargest
	io.Pfblue2("\nfunc MatLargest(a [][]float64, den float64) (largest float64)\n")
	PrintMat("a", a, "%5g", false)
	largest := MatLargest(a, 1)
	io.Pf("largest(a) = %v\n", largest)
	chk.Scalar(tst, "largest(a)", 1e-17, largest, 50)

	// MatGetCol
	io.Pfblue2("\nfunc MatGetCol(j int, a [][]float64) (col []float64)\n")
	col := MatGetCol(0, a)
	PrintMat("a", a, "%5g", false)
	PrintVec("a[:][0]", col, "%5g", false)
	chk.Vector(tst, "a[:][0]", 1e-17, col, []float64{1, 0.1, 10})

	// MatNormF
	io.Pfblue2("\nfunc MatNormF(a [][]float64) (res float64)\n")
	A := [][]float64{
		{-3, 5, 7},
		{2, 6, 4},
		{0, 2, 8},
	}
	PrintMat("A", A, "%5g", false)
	normFA := MatNormF(A)
	io.Pf("normF(A) = %g\n", normFA)
	chk.Scalar(tst, "normF(A)", 1e-17, normFA, 1.438749456993816e+01)

	// MatNormI
	io.Pfblue2("\nfunc MatNormI(a [][]float64) (res float64)\n")
	PrintMat("A", A, "%5g", false)
	normIA := MatNormI(A)
	io.Pf("normI(A) = %g\n", normIA)
	chk.Scalar(tst, "normI(A)", 1e-17, normIA, 15.0)
}
