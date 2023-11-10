// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
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
	chk.Vector(tst, "A to a", 1e-15, a.data, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})

	chk.Scalar(tst, "Get(1,1)", 1e-17, a.Get(1, 1), 6)
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
	chk.VectorC(tst, "A to a", 1e-15, a.data, []complex128{1 + 0.1i, 5 + 0.5i, 9 + 0.9i, 2, 6, 0, 3, 7, -1, 4 - 0.4i, 8 - 0.8i, -2 + 1i})

	chk.ScalarC(tst, "Get(1,1)", 1e-17, a.Get(1, 1), 6)
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
}

func TestDaxpy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Daxpy01")

	α := 0.5
	x := []float64{20, 10, 30, 123, 123}
	y := []float64{-15, -5, -24, 666, 666, 666}
	n, incx, incy := 3, 1, 1
	err := Daxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}

	chk.Vector(tst, "x", 1e-15, x, []float64{20, 10, 30, 123, 123})
	chk.Vector(tst, "y", 1e-15, y, []float64{-5, 0, -9, 666, 666, 666})
}

func TestZaxpy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zaxpy01")

	α := 1.0 + 0i
	x := []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i}
	y := []complex128{-15 + 1.5i, -5 - 2i, -24 + 1i, 666 - 0.5i, 666 + 5i}
	n, incx, incy := len(x), 1, 1
	err := Zaxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}

	chk.VectorC(tst, "x", 1e-15, x, []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i})
	chk.VectorC(tst, "y", 1e-15, y, []complex128{5 + 2.5i, 5, 6 + 2.5i, 543, 543 + 5.5i})

	α = 0.5 + 1i
	err = Zaxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "y", 1e-15, y, []complex128{14.0 + 23.i, 8.0 + 11.i, 19.5 + 33.25i, 481.0 - 122.75i, 481.0 - 117.25i})
}

func TestDgemv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemv01")

	// allocate
	a := NewMatrix([][]float64{
		{0.1, 0.2, 0.3},
		{1.0, 0.2, 0.3},
		{2.0, 0.2, 0.3},
		{3.0, 0.2, 0.3},
	})
	chk.Vector(tst, "a.data", 1e-15, a.data, []float64{0.1, 1, 2, 3, 0.2, 0.2, 0.2, 0.2, 0.3, 0.3, 0.3, 0.3})

	// perform mv
	m, n := a.M(), a.N()
	α, β := 0.5, 2.0
	x := []float64{20, 10, 30}
	y := []float64{3, 1, 2, 4}
	lda, incx, incy := m, 1, 1
	err := Dgemv(false, m, n, α, a, lda, x, incx, β, y, incy)
	if err != nil {
		tst.Errorf("Dgemv failed:\n%v\n", err)
		return
	}
	chk.Vector(tst, "y", 1e-15, y, []float64{12.5, 17.5, 29.5, 43.5})

	// perform mv with transpose
	err = Dgemv(true, m, n, α, a, lda, y, incy, β, x, incx)
	if err != nil {
		tst.Errorf("Dgemv failed:\n%v\n", err)
		return
	}
	chk.Vector(tst, "x", 1e-15, x, []float64{144.125, 30.3, 75.45})

	// check that a is unmodified
	chk.Vector(tst, "a.data", 1e-15, a.data, []float64{0.1, 1, 2, 3, 0.2, 0.2, 0.2, 0.2, 0.3, 0.3, 0.3, 0.3})
}

func TestZgemv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgemv01")

	// allocate
	a := NewMatrixC([][]complex128{
		{0.1 + 3i, 0.2, 0.3 - 0.3i},
		{1.0 + 2i, 0.2, 0.3 - 0.4i},
		{2.0 + 1i, 0.2, 0.3 - 0.5i},
		{3.0 + 0.1i, 0.2, 0.3 - 0.6i},
	})
	m, n := a.M(), a.N()
	chk.VectorC(tst, "a.data", 1e-15, a.data, []complex128{0.1 + 3i, 1 + 2i, 2 + 1i, 3 + 0.1i, 0.2, 0.2, 0.2, 0.2, 0.3 - 0.3i, 0.3 - 0.4i, 0.3 - 0.5i, 0.3 - 0.6i})

	// perform mv
	α, β := 0.5+1i, 2.0+1i
	x := []complex128{20, 10, 30}
	y := []complex128{3, 1, 2, 4}
	lda, incx, incy := m, 1, 1
	err := Zgemv(false, m, n, α, a, lda, x, incx, β, y, incy)
	if err != nil {
		tst.Errorf("Zgemv failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "y", 1e-15, y, []complex128{-38.5 + 41.5i, -10.5 + 46i, 24.5 + 55.5i, 59.5 + 67i})

	// perform mv with transpose
	err = Zgemv(true, m, n, α, a, lda, y, incy, β, x, incx)
	if err != nil {
		tst.Errorf("Zgemv failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "x", 1e-13, x, []complex128{-248.875 + 82.5i, -18.5 + 38i, 83.85 + 154.7i})

	// check that a is unmodified
	chk.VectorC(tst, "a.data", 1e-15, a.data, []complex128{0.1 + 3i, 1 + 2i, 2 + 1i, 3 + 0.1i, 0.2, 0.2, 0.2, 0.2, 0.3 - 0.3i, 0.3 - 0.4i, 0.3 - 0.5i, 0.3 - 0.6i})
}

func TestDgemm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemm01")

	// allocate matrices
	a := NewMatrix([][]float64{
		{1, 2, 0, 1, -1},
		{2, 3, -1, 1, +1},
		{1, 2, 0, 4, -1},
		{4, 0, 3, 1, +1},
	})
	b := NewMatrix([][]float64{
		{1, 0, 0},
		{0, 0, 3},
		{0, 0, 1},
		{1, 0, 1},
		{0, 2, 0},
	})
	c := NewMatrix([][]float64{
		{+0.50, 0, +0.25},
		{+0.25, 0, -0.25},
		{-0.25, 0, +0.00},
		{-0.25, 0, +0.00},
	})

	// sizes
	m, k := a.M(), a.N()
	n := b.N()

	// run dgemm
	transA, transB := false, false
	alpha, beta := 0.5, 2.0
	lda, ldb, ldc := m, k, m
	err := Dgemm(transA, transB, m, n, k, alpha, a, lda, b, ldb, beta, c, ldc)
	if err != nil {
		tst.Errorf("Dgemm failed:\n%v\n", err)
		return
	}

	// check
	chk.Matrix(tst, "0.5⋅a⋅b + 2⋅c", 1e-17, c.GetSlice(), [][]float64{
		{2, -1, 4},
		{2, 1, 4},
		{2, -1, 5},
		{2, 1, 2},
	})
}

func TestZgemm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgemm01")

	// allocate matrices
	a := NewMatrixC([][]complex128{
		{1, 2, 0 + 1i, 1, -1},
		{2, 3, -1 - 1i, 1, +1},
		{1, 2, 0 + 1i, 4, -1},
		{4, 0, 3 - 1i, 1, +1},
	})
	b := NewMatrixC([][]complex128{
		{1, 0, 0 + 1i},
		{0, 0, 3 - 1i},
		{0, 0, 1 + 1i},
		{1, 0, 1 - 1i},
		{0, 2, 0 + 1i},
	})
	c := NewMatrixC([][]complex128{
		{+0.50, 1i, +0.25},
		{+0.25, 1i, -0.25},
		{-0.25, 1i, +0.00},
		{-0.25, 1i, +0.00},
	})

	// sizes
	m, k := a.M(), a.N()
	n := b.N()

	// run dgemm
	transA, transB := false, false
	alpha, beta := 0.5-2i, 2.0-4i
	lda, ldb, ldc := m, k, m
	err := Zgemm(transA, transB, m, n, k, alpha, a, lda, b, ldb, beta, c, ldc)
	if err != nil {
		tst.Errorf("Zgemm failed:\n%v\n", err)
		return
	}

	// check
	chk.MatrixC(tst, "(0.5-2i)⋅a⋅b + (2-4i)⋅c", 1e-17, c.GetSlice(), [][]complex128{
		{2 - 6i, 3 + 6i, -0.5 - 14i},
		{2 - 7i, 5 - 2i, -1.5 - 20.5i},
		{2 - 9i, 3 + 6i, -5.5 - 20.5i},
		{2 - 9i, 5 - 2i, 14.5 - 7i},
	})
}

func TestDgesv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgesv01")

	// matrix
	amat := [][]float64{
		{2, 3, 0, 0, 0},
		{3, 0, 4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, 0, 1, 0, 0},
		{0, 4, 2, 0, 1},
	}
	n := 5
	a := NewMatrixMN(n, n)
	a.SetFromSlice(amat)

	// right-hand-side
	b := []float64{8, 45, -3, 3, 19}

	// solution
	xCorrect := []float64{1, 2, 3, 4, 5}

	// run test
	nrhs := 1
	lda, ldb := n, n
	ipiv := make([]int32, n)
	err := Dgesv(n, nrhs, a, lda, ipiv, b, ldb)
	if err != nil {
		tst.Errorf("Dgesv failed:\n%v\n", err)
		return
	}
	chk.Vector(tst, "x = A⁻¹ b", 1e-15, b, xCorrect)
}

func TestZgesv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgesv01. low accuracy")

	// NOTE: zgesv performs badly with this problem
	//       the best tolerance that can be selected is 0.00032
	//       the same problem happens in python (probably using lapack as well)
	tol := 0.0032

	// matrix
	a := NewMatrixC([][]complex128{
		{19.730 + 0.000i, 12.110 - 1.000i, 0.000 + 5.000i, 0.000 + 0.000i, 0.000 + 0.000i},
		{0.000 - 0.510i, 32.300 + 7.000i, 23.070 + 0.000i, 0.000 + 1.000i, 0.000 + 0.000i},
		{0.000 + 0.000i, 0.000 - 0.510i, 70.000 + 7.300i, 3.950 + 0.000i, 19.000 + 31.830i},
		{0.000 + 0.000i, 0.000 + 0.000i, 1.000 + 1.100i, 50.170 + 0.000i, 45.510 + 0.000i},
		{0.000 + 0.000i, 0.000 + 0.000i, 0.000 + 0.000i, 0.000 - 9.351i, 55.000 + 0.000i},
	})

	// right-hand-side
	b := []complex128{
		77.38 + 8.82i,
		157.48 + 19.8i,
		1175.62 + 20.69i,
		912.12 - 801.75i,
		550 - 1060.4i,
	}

	// solution
	xCorrect := []complex128{
		3.3 - 1i,
		1 + 0.17i,
		5.5,
		9,
		10 - 17.75i,
	}

	// run test
	n := a.N()
	nrhs := 1
	lda, ldb := n, n
	ipiv := make([]int32, n)
	err := Zgesv(n, nrhs, a, lda, ipiv, b, ldb)
	if err != nil {
		tst.Errorf("Zgesv failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "x = A⁻¹ b", tol, b, xCorrect)

	// compare with python results
	xPython := []complex128{
		3.299687426933794e+00 - 1.000372829305209e+00i,
		9.997606020636992e-01 + 1.698383755401385e-01i,
		5.500074759292877e+00 - 4.556001293922560e-05i,
		8.999787912842375e+00 - 6.662818244209770e-05i,
		1.000001132800243e+01 - 1.774987242230929e+01i,
	}
	chk.VectorC(tst, "x = A⁻¹ b", 1e-13, b, xPython)
}

func checksvd(tst *testing.T, amat, uCorrect, vtCorrect [][]float64, sCorrect []float64, tolu, tols, tolv, tolusv float64) {

	// allocate matrix
	m, n := len(amat), len(amat[0])
	a := NewMatrixMN(m, n)
	a.SetFromSlice(amat)

	// compute dimensions
	minMN := imin(m, n)
	maxMN := imax(m, n)
	lda := m
	ldu := m
	ldvt := n
	lwork := 2 * imax(3*minMN+maxMN, 5*minMN)

	// allocate output arrays
	s := make([]float64, minMN)
	u := NewMatrixMN(m, m)
	vt := NewMatrixMN(n, n)
	work := make([]float64, lwork)

	// perform SVD
	jobu := 'A'
	jobvt := 'A'
	err := Dgesvd(jobu, jobvt, m, n, a, lda, s, u, ldu, vt, ldvt, work, lwork)
	if err != nil {
		tst.Errorf("Dgesvd failed:\n%v\n", err)
		return
	}

	// compare results
	umat := u.GetSlice()
	vtmat := vt.GetSlice()
	chk.Matrix(tst, "u", tolu, umat, uCorrect)
	chk.Vector(tst, "s", tols, s, sCorrect)
	chk.Matrix(tst, "vt", tolv, vtmat, vtCorrect)

	// check SVD
	usv := make([][]float64, m)
	for i := 0; i < m; i++ {
		usv[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			for k := 0; k < minMN; k++ {
				usv[i][j] += umat[i][k] * s[k] * vtmat[k][j]
			}
		}
	}
	chk.Matrix(tst, "u⋅s⋅vt", tolusv, amat, usv)
}

func TestDgesvd02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgesvd02")

	// allocate matrices
	s33 := math.Sqrt(3.0) / 3.0
	amat := [][]float64{
		{-s33, -s33, 1},
		{+s33, -s33, 1},
		{-s33, +s33, 1},
		{+s33, +s33, 1},
	}
	uCorrect := [][]float64{
		{-0.5, -0.5, -0.5, +0.5},
		{-0.5, -0.5, +0.5, -0.5},
		{-0.5, +0.5, -0.5, -0.5},
		{-0.5, +0.5, +0.5, +0.5},
	}
	sCorrect := []float64{2, 2.0 / math.Sqrt(3.0), 2.0 / math.Sqrt(3.0)}
	vtCorrect := [][]float64{
		{+0, +0, -1},
		{+0, +1, +0},
		{+1, +0, +0},
	}

	// check
	checksvd(tst, amat, uCorrect, vtCorrect, sCorrect, 1e-15, 1e-15, 1e-17, 1e-15)
}

func checksvdC(tst *testing.T, amat, uCorrect, vtCorrect [][]complex128, sCorrect []float64, tolu, tols, tolv, tolusv float64) {

	// allocate matrix
	m, n := len(amat), len(amat[0])
	a := NewMatrixCmn(m, n)
	a.SetFromSlice(amat)

	// compute dimensions
	minMN := imin(m, n)
	maxMN := imax(m, n)
	lda := m
	ldu := m
	ldvt := n
	lwork := 2 * (2*minMN + maxMN)
	lrwork := 2 * (5 * minMN)

	// allocate output arrays
	s := make([]float64, minMN)
	u := NewMatrixCmn(m, m)
	vt := NewMatrixCmn(n, n)
	work := make([]complex128, lwork)
	rwork := make([]float64, lrwork)

	// perform SVD
	jobu := 'A'
	jobvt := 'A'
	err := Zgesvd(jobu, jobvt, m, n, a, lda, s, u, ldu, vt, ldvt, work, lwork, rwork)
	if err != nil {
		tst.Errorf("Zgesvd failed:\n%v\n", err)
		return
	}

	// compare results
	umat := u.GetSlice()
	vtmat := vt.GetSlice()
	chk.MatrixC(tst, "u", tolu, umat, uCorrect)
	chk.Vector(tst, "s", tols, s, sCorrect)
	chk.MatrixC(tst, "vt", tolv, vtmat, vtCorrect)

	// check SVD
	usv := make([][]complex128, m)
	for i := 0; i < m; i++ {
		usv[i] = make([]complex128, n)
		for j := 0; j < n; j++ {
			for k := 0; k < minMN; k++ {
				usv[i][j] += umat[i][k] * complex(s[k], 0) * vtmat[k][j]
			}
		}
	}
	chk.MatrixC(tst, "u⋅s⋅vt", tolusv, amat, usv)
}

func TestDgetrf01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgetrf01. Dgetrf and Dgetri")

	// matrix
	amat := [][]float64{
		{1, 2, 0, 1},
		{2, 3, -1, 1},
		{1, 2, 0, 4},
		{4, 0, 3, 1},
	}
	m, n := 4, 4
	a := NewMatrixMN(m, n)
	a.SetFromSlice(amat)

	// run dgetrf
	lda := m
	ipiv := make([]int32, imin(m, n))
	err := Dgetrf(m, n, a, lda, ipiv)
	if err != nil {
		tst.Errorf("Dgetrf failed:\n%v\n", err)
		return
	}

	// check LU
	chk.Matrix(tst, "lu", 1e-17, a.GetSlice(), [][]float64{
		{+4.0e+00, +0.000000000000000e+00, +3.000000000000000e+00, +1.000000000000000e+00},
		{+5.0e-01, +3.000000000000000e+00, -2.500000000000000e+00, +5.000000000000000e-01},
		{+2.5e-01, +6.666666666666666e-01, +9.166666666666665e-01, +3.416666666666667e+00},
		{+2.5e-01, +6.666666666666666e-01, +1.000000000000000e+00, -3.000000000000000e+00},
	})

	// run dgetri
	lwork := 2 * n
	work := make([]float64, lwork)
	err = Dgetri(n, a, lda, ipiv, work, lwork)
	if err != nil {
		tst.Errorf("Dgetri failed:\n%v\n", err)
		return
	}

	// compare inverse
	ai := a.GetSlice()
	chk.Matrix(tst, "inv(a)", 1e-15, ai, [][]float64{
		{-8.484848484848487e-01, +5.454545454545455e-01, +3.030303030303039e-02, +1.818181818181818e-01},
		{+1.090909090909091e+00, -2.727272727272728e-01, -1.818181818181817e-01, -9.090909090909091e-02},
		{+1.242424242424243e+00, -7.272727272727273e-01, -1.515151515151516e-01, +9.090909090909088e-02},
		{-3.333333333333333e-01, +0.000000000000000e+00, +3.333333333333333e-01, +0.000000000000000e+00},
	})

	// check inverse
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			res := 0.0
			for k := 0; k < m; k++ {
				res += amat[i][k] * ai[k][j]
			}
			if i == j {
				chk.Scalar(tst, "diag(a⋅a⁻¹)=diag(I)=1", 1e-15, res, 1)
			} else {
				chk.Scalar(tst, "diag(a⋅a⁻¹)=offdiag(I)=0", 1e-15, res, 0)
			}
		}
	}
}

func TestZgetrf01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgetrf01. Zgetrf and Zgetri")

	// matrix
	amat := [][]complex128{
		{1 + 1i, 2, +0, 1 - 1i},
		{2 + 1i, 3, -1, 1 - 1i},
		{1 + 1i, 2, +0, 4 - 1i},
		{4 + 1i, 0, +3, 1 - 1i},
	}
	m, n := 4, 4
	a := NewMatrixCmn(m, n)
	a.SetFromSlice(amat)

	// run
	lda := m
	ipiv := make([]int32, imin(m, n))
	err := Zgetrf(m, n, a, lda, ipiv)
	if err != nil {
		tst.Errorf("Zgetrf failed:\n%v\n", err)
		return
	}

	// check LU
	chk.MatrixC(tst, "lu", 1e-15, a.GetSlice(), [][]complex128{
		{+4.000000000000000e+00 + 1.000000000000000e+00i, +0.000000000000000e+00, +3.000000000000000e+00 + 0.000000000000000e+00i, +1.000000000000000e+00 - 1.000000000000000e+00i},
		{+5.294117647058824e-01 + 1.176470588235294e-01i, +3.000000000000000e+00, -2.588235294117647e+00 - 3.529411764705882e-01i, +3.529411764705882e-01 - 5.882352941176471e-01i},
		{+2.941176470588235e-01 + 1.764705882352941e-01i, +6.666666666666666e-01, +8.431372549019609e-01 - 2.941176470588235e-01i, +3.294117647058823e+00 - 4.901960784313725e-01i},
		{+2.941176470588235e-01 + 1.764705882352941e-01i, +6.666666666666666e-01, +1.000000000000000e+00 + 0.000000000000000e+00i, -3.000000000000000e+00 + 0.000000000000000e+00i},
	})

	// run zgetri
	lwork := 2 * n
	work := make([]complex128, lwork)
	err = Zgetri(n, a, lda, ipiv, work, lwork)
	if err != nil {
		tst.Errorf("Zgetri failed:\n%v\n", err)
		return
	}

	// compare inverse
	ai := a.GetSlice()
	chk.MatrixC(tst, "inv(a)", 1e-15, ai, [][]complex128{
		{-8.442622950819669e-01 - 4.644808743169393e-02i, +5.409836065573769e-01 + 4.918032786885240e-02i, +3.278688524590156e-02 - 2.732240437158467e-02i, +1.803278688524591e-01 + 1.639344262295081e-02i},
		{+1.065573770491803e+00 + 2.786885245901638e-01i, -2.459016393442623e-01 - 2.950819672131146e-01i, -1.967213114754096e-01 + 1.639344262295082e-01i, -8.196721311475419e-02 - 9.836065573770497e-02i},
		{+1.221311475409836e+00 + 2.322404371584698e-01i, -7.049180327868851e-01 - 2.459016393442622e-01i, -1.639344262295082e-01 + 1.366120218579235e-01i, +9.836065573770481e-02 - 8.196721311475411e-02i},
		{-3.333333333333333e-01 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +3.333333333333333e-01 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i},
	})

	// check inverse
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			res := 0.0 + 0.0i
			for k := 0; k < m; k++ {
				res += amat[i][k] * ai[k][j]
			}
			if i == j {
				chk.ScalarC(tst, "diag(a⋅a⁻¹)=diag(I)=1", 1e-15, res, 1)
			} else {
				chk.ScalarC(tst, "diag(a⋅a⁻¹)=offdiag(I)=0", 1e-15, res, 0)
			}
		}
	}
}

func checkUplo(tst *testing.T, testname string, c, cLo, cUp *Matrix, tol float64) {
	n := c.n
	maxdiff := 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				diff := math.Abs(cLo.Get(i, j) - c.Get(i, j))
				if diff > tol {
					maxdiff = diff
				}
				diff = math.Abs(cUp.Get(i, j) - c.Get(i, j))
				if diff > tol {
					maxdiff = diff
				}
			} else {
				diff := math.Abs(cLo.Get(i, j) + cUp.Get(i, j) - c.Get(i, j))
				if diff > tol {
					maxdiff = diff
				}
			}
		}
	}
	if maxdiff > 0 {
		tst.Errorf("checkUplo failed in test %q. maxdiff=%g\n", testname, maxdiff)
	}
}

func TestDsyrk01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dsyrk01")

	// c matrices
	c := NewMatrix([][]float64{
		{+3, +0, -3, +0},
		{+0, +3, +1, +2},
		{-3, +1, +4, +1},
		{+0, +2, +1, +3},
	})
	cUp := NewMatrix([][]float64{
		{+3, +0, -3, +0},
		{+0, +3, +1, +2},
		{+0, +0, +4, +1},
		{+0, +0, +0, +3},
	})
	cLo := NewMatrix([][]float64{
		{+3, +0, +0, +0},
		{+0, +3, +0, +0},
		{-3, +1, +4, +0},
		{+0, +2, +1, +3},
	})

	// check cUp and cLo
	checkUplo(tst, "Dsyrk01", c, cLo, cUp, 1e-17)

	// a matrix
	a := NewMatrix([][]float64{
		{+1, +2, +1, +1, -1, +0},
		{+2, +2, +1, +0, +0, +0},
		{+3, +1, +3, +1, +2, -1},
		{+1, +0, +1, -1, +0, +0},
	})

	// sizes
	n := c.N()
	k := a.N()

	// constants
	alpha, beta := 3.0, -1.0

	// run dsyrk with up(c)
	up, trans := true, false
	lda, ldc := n, n
	err := Dsyrk(up, trans, n, k, alpha, a, lda, beta, cUp, ldc)
	if err != nil {
		tst.Errorf("Dsyrk failed:\n%v\n", err)
		return
	}

	// compare resulting up(c) matrix
	chk.Matrix(tst, "using up(c): c := 3⋅a⋅aᵀ - c", 1e-17, cUp.GetSlice(), [][]float64{
		{21, 21, 24, +3},
		{+0, 24, 32, +7},
		{+0, +0, 71, 14},
		{+0, +0, +0, +6},
	})

	// run dsyrk with lo(c)
	up = false
	err = Dsyrk(up, trans, n, k, alpha, a, lda, beta, cLo, ldc)
	if err != nil {
		tst.Errorf("Dsyrk failed:\n%v\n", err)
		return
	}

	// compare resulting up(c) matrix
	chk.Matrix(tst, "using lo(c): c := 3⋅a⋅aᵀ - c", 1e-17, cLo.GetSlice(), [][]float64{
		{21, +0, +0, +0},
		{21, 24, +0, +0},
		{24, 32, 71, +0},
		{+3, +7, 14, +6},
	})
}

func TestDsyrk02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dsyrk02")

	// c matrices
	c := NewMatrix([][]float64{
		{+3, 0, -3, 0, 0, 0},
		{+0, 3, +1, 2, 2, 2},
		{-3, 1, +4, 1, 1, 1},
		{+0, 2, +1, 3, 3, 3},
		{+0, 2, +1, 3, 4, 3},
		{+0, 2, +1, 3, 3, 4},
	})
	cUp := NewMatrix([][]float64{
		{+3, 0, -3, 0, 0, 0},
		{+0, 3, +1, 2, 2, 2},
		{+0, 0, +4, 1, 1, 1},
		{+0, 0, +0, 3, 3, 3},
		{+0, 0, +0, 0, 4, 3},
		{+0, 0, +0, 0, 0, 4},
	})
	cLo := NewMatrix([][]float64{
		{+3, 0, +0, 0, 0, 0},
		{+0, 3, +0, 0, 0, 0},
		{-3, 1, +4, 0, 0, 0},
		{+0, 2, +1, 3, 0, 0},
		{+0, 2, +1, 3, 4, 0},
		{+0, 2, +1, 3, 3, 4},
	})

	// check cUp and cLo
	checkUplo(tst, "Dsyrk02", c, cLo, cUp, 1e-17)

	// a matrix
	a := NewMatrix([][]float64{
		{+1, +2, +1, +1, -1, +0},
		{+2, +2, +1, +0, +0, +0},
		{+3, +1, +3, +1, +2, -1},
		{+1, +0, +1, -1, +0, +0},
	})

	// sizes
	n := c.N()
	k := a.M() // m now

	// constants
	alpha, beta := 3.0, +1.0

	// run dsyrk with up(c)
	up, trans := true, true
	lda, ldc := k, n
	err := Dsyrk(up, trans, n, k, alpha, a, lda, beta, cUp, ldc)
	if err != nil {
		tst.Errorf("Dsyrk failed:\n%v\n", err)
		return
	}

	// compare resulting up(c) matrix
	chk.Matrix(tst, "using up(c): c := 3⋅a⋅aᵀ + c", 1e-17, cUp.GetSlice(), [][]float64{
		{48, 27, 36, +9, 15, -9},
		{+0, 30, 22, 11, +2, -1},
		{+0, +0, 40, 10, 16, -8},
		{+0, +0, +0, 12, +6, +0},
		{+0, +0, +0, +0, 19, -3},
		{+0, +0, +0, +0, +0, +7},
	})

	// run dsyrk with lo(c)
	up = false
	err = Dsyrk(up, trans, n, k, alpha, a, lda, beta, cLo, ldc)
	if err != nil {
		tst.Errorf("Dsyrk failed:\n%v\n", err)
		return
	}

	// compare resulting up(c) matrix
	chk.Matrix(tst, "using lo(c): c := 3⋅a⋅aᵀ + c", 1e-17, cLo.GetSlice(), [][]float64{
		{48, +0, +0, +0, +0, +0},
		{27, 30, +0, +0, +0, +0},
		{36, 22, 40, +0, +0, +0},
		{+9, 11, 10, 12, +0, +0},
		{15, +2, 16, +6, 19, +0},
		{-9, -1, -8, +0, -3, +7},
	})
}

func checkUploC(tst *testing.T, testname string, c, cLo, cUp *MatrixC, tolR, tolI float64) {
	n := c.n
	maxdiffR, maxdiffI := 0.0, 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				diffR := math.Abs(real(cLo.Get(i, j)) - real(c.Get(i, j)))
				diffI := math.Abs(imag(cLo.Get(i, j)) - imag(c.Get(i, j)))
				if diffR > tolR {
					maxdiffR = diffR
				}
				if diffI > tolI {
					maxdiffI = diffI
				}
				diffR = math.Abs(real(cUp.Get(i, j)) - real(c.Get(i, j)))
				diffI = math.Abs(imag(cUp.Get(i, j)) - imag(c.Get(i, j)))
				if diffR > tolR {
					maxdiffR = diffR
				}
				if diffI > tolI {
					maxdiffI = diffI
				}
			} else {
				diffR := math.Abs(real(cLo.Get(i, j)) + real(cUp.Get(i, j)) - real(c.Get(i, j)))
				diffI := math.Abs(imag(cLo.Get(i, j)) + imag(cUp.Get(i, j)) - imag(c.Get(i, j)))
				if diffR > tolR {
					maxdiffR = diffR
				}
				if diffI > tolI {
					maxdiffI = diffI
				}
			}
		}
	}
	if maxdiffR > 0 || maxdiffI > 0 {
		tst.Errorf("checkUploC failed in test %q. maxdiffR=%g maxdiffI=%g\n", testname, maxdiffR, maxdiffI)
	}
}

func TestZsyrk01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zsyrk01")

	// c matrices
	c := NewMatrixC([][]complex128{
		{+3 + 1i, 0 + 0i, -2 + 0i, 0 + 0i},
		{-1 + 0i, 3 + 0i, +0 + 0i, 2 + 0i},
		{-4 + 0i, 1 + 0i, +3 + 0i, 1 + 0i},
		{-1 + 0i, 2 + 0i, +0 + 0i, 3 - 1i},
	})
	cUp := NewMatrixC([][]complex128{
		{+3 + 1i, 0 + 0i, -2 + 0i, 0 + 0i},
		{+0 + 0i, 3 + 0i, +0 + 0i, 2 + 0i},
		{+0 + 0i, 0 + 0i, +3 + 0i, 1 + 0i},
		{+0 + 0i, 0 + 0i, +0 + 0i, 3 - 1i},
	})
	cLo := NewMatrixC([][]complex128{
		{+3 + 1i, 0 + 0i, +0 + 0i, 0 + 0i},
		{-1 + 0i, 3 + 0i, +0 + 0i, 0 + 0i},
		{-4 + 0i, 1 + 0i, +3 + 0i, 0 + 0i},
		{-1 + 0i, 2 + 0i, +0 + 0i, 3 - 1i},
	})

	// check cUp and cLo
	checkUploC(tst, "Zsyrk02", c, cLo, cUp, 1e-17, 1e-17)

	// a matrix
	a := NewMatrixC([][]complex128{
		{+1 - 1i, +2, +1, +1, -1, +0 + 0i},
		{+2 + 0i, +2, +1, +0, +0, +0 + 1i},
		{+3 + 1i, +1, +3, +1, +2, -1 + 0i},
		{+1 + 0i, +0, +1, -1, +0, +0 + 1i},
	})

	// sizes
	n := c.N()
	k := a.N()

	// constants
	alpha, beta := 3.0+0i, +1.0+0i

	// run zsyrk with up(c)
	up, trans := true, false
	lda, ldc := n, n
	err := Zsyrk(up, trans, n, k, alpha, a, lda, beta, cUp, ldc)
	if err != nil {
		tst.Errorf("Zsyrk failed:\n%v\n", err)
		return
	}

	// compare resulting up(c) matrix
	chk.MatrixC(tst, "using up(c): c := 3⋅a⋅aᵀ + c", 1e-17, cUp.GetSlice(), [][]complex128{
		{24 - 5i, 21 - 6i, 22 - 6i, +3 - 3i},
		{+0 + 0i, 27 + 0i, 33 + 3i, +8 + 0i},
		{+0 + 0i, +0 + 0i, 75 + 18i, 16 + 0i},
		{+0 + 0i, +0 + 0i, +0 + 0i, +9 - 1i},
	})

	// run zsyrk with lo(c)
	up = false
	err = Zsyrk(up, trans, n, k, alpha, a, lda, beta, cLo, ldc)
	if err != nil {
		tst.Errorf("Zsyrk failed:\n%v\n", err)
		return
	}

	// compare resulting up(c) matrix
	chk.MatrixC(tst, "using lo(c): c := 3⋅a⋅aᵀ + c", 1e-17, cLo.GetSlice(), [][]complex128{
		{24 - 5i, +0 + 0i, +0 + 0i, +0 + 0i},
		{20 - 6i, 27 + 0i, +0 + 0i, +0 + 0i},
		{20 - 6i, 34 + 3i, 75 + 18i, +0 + 0i},
		{+2 - 3i, +8 + 0i, 15 + 0i, +9 - 1i},
	})
}

func TestZherk01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zherk01")

	// c matrices
	c := NewMatrixC([][]complex128{ // must be Hermitian: c = c^H
		{+4 + 0i, 0 + 1i, -3 + 1i, 0 + 2i},
		{+0 - 1i, 3 + 0i, +1 + 0i, 2 + 0i},
		{-3 - 1i, 1 + 0i, +4 + 0i, 1 - 1i},
		{+0 - 2i, 2 + 0i, +1 + 1i, 4 + 0i},
	})
	cUp := NewMatrixC([][]complex128{
		{+4 + 0i, 0 + 1i, -3 + 1i, 0 + 2i},
		{+0 + 0i, 3 + 0i, +1 + 0i, 2 + 0i},
		{+0 + 0i, 0 + 0i, +4 + 0i, 1 - 1i},
		{+0 + 0i, 0 + 0i, +0 + 0i, 4 + 0i},
	})
	cLo := NewMatrixC([][]complex128{
		{+4 + 0i, 0 + 0i, +0 + 0i, 0 + 0i},
		{+0 - 1i, 3 + 0i, +0 + 0i, 0 + 0i},
		{-3 - 1i, 1 + 0i, +4 + 0i, 0 + 0i},
		{+0 - 2i, 2 + 0i, +1 + 1i, 4 + 0i},
	})

	// check cUp and cLo
	checkUploC(tst, "Zherk01", c, cLo, cUp, 1e-17, 1e-17)

	// a matrix
	a := NewMatrixC([][]complex128{
		{+1 - 1i, +2, +1, +1, -1, +0 + 0i},
		{+2 + 0i, +2, +1, +0, +0, +0 + 1i},
		{+3 + 1i, +1, +3, +1, +2, -1 + 0i},
		{+1 + 0i, +0, +1, -1, +0, +0 + 1i},
	})

	// sizes
	n := c.N()
	k := a.N()

	// constants
	alpha, beta := 3.0, +1.0

	// run zherk with up(c)
	up, trans := true, false
	lda, ldc := n, n
	err := Zherk(up, trans, n, k, alpha, a, lda, beta, cUp, ldc)
	if err != nil {
		tst.Errorf("Zherk failed:\n%v\n", err)
		return
	}

	// compare resulting up(c) matrix
	chk.MatrixC(tst, "using up(c): c := 3⋅a⋅aᵀ + c", 1e-17, cUp.GetSlice(), [][]complex128{
		{31 + 0i, 21 - 5i, 15 - 11i, 3 - 1i},
		{+0 + 0i, 33 + 0i, 34 - 9i, 14 + 0i},
		{+0 + 0i, +0 + 0i, 82 + 0i, 16 + 5i},
		{+0 + 0i, +0 + 0i, +0 + 0i, 16 + 0i},
	})

	// run zherk with lo(c)
	up = false
	err = Zherk(up, trans, n, k, alpha, a, lda, beta, cLo, ldc)
	if err != nil {
		tst.Errorf("Zherk failed:\n%v\n", err)
		return
	}

	// compare resulting up(c) matrix
	chk.MatrixC(tst, "using lo(c): c := 3⋅a⋅aᵀ + c", 1e-17, cLo.GetSlice(), [][]complex128{
		{31 + 0i, +0 + 0i, +0 + 0i, +0 + 0i},
		{21 + 5i, 33 + 0i, +0 + 0i, +0 + 0i},
		{15 + 11i, 34 + 9i, 82 + 0i, +0 + 0i},
		{3 + 1i, 14 + 0i, 16 - 5i, 16 + 0i},
	})
}

func TestDpotrf01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dpotrf01")

	// a matrices
	a := NewMatrix([][]float64{
		{+3, +0, -3, +0},
		{+0, +3, +1, +2},
		{-3, +1, +4, +1},
		{+0, +2, +1, +3},
	})
	aUp := NewMatrix([][]float64{
		{+3, +0, -3, +0},
		{+0, +3, +1, +2},
		{+0, +0, +4, +1},
		{+0, +0, +0, +3},
	})
	aLo := NewMatrix([][]float64{
		{+3, +0, +0, +0},
		{+0, +3, +0, +0},
		{-3, +1, +4, +0},
		{+0, +2, +1, +3},
	})

	// check aUp and aLo
	checkUplo(tst, "Dpotrf01", a, aLo, aUp, 1e-17)

	// run dpotrf with up(a)
	n := a.N()
	up := true
	lda := n
	err := Dpotrf(up, n, aUp, lda)
	if err != nil {
		tst.Errorf("Dpotrf failed:\n%v\n", err)
		return
	}

	// check aUp
	chk.Matrix(tst, "chol(aUp)", 1e-15, aUp.GetSlice(), [][]float64{
		{+1.732050807568877e+00, +0.000000000000000e+00, -1.732050807568878e+00, +0.000000000000000e+00},
		{+0.000000000000000e+00, +1.732050807568877e+00, +5.773502691896258e-01, +1.154700538379252e+00},
		{+0.000000000000000e+00, +0.000000000000000e+00, +8.164965809277251e-01, +4.082482904638632e-01},
		{+0.000000000000000e+00, +0.000000000000000e+00, +0.000000000000000e+00, +1.224744871391589e+00},
	})

	// run dpotrf with lo(a)
	up = false
	err = Dpotrf(up, n, aLo, lda)
	if err != nil {
		tst.Errorf("Dpotrf failed:\n%v\n", err)
		return
	}

	// check aLo
	chk.Matrix(tst, "chol(aLo)", 1e-15, aLo.GetSlice(), [][]float64{
		{+1.732050807568877e+00, +0.000000000000000e+00, +0.000000000000000e+00, +0.000000000000000e+00},
		{+0.000000000000000e+00, +1.732050807568877e+00, +0.000000000000000e+00, +0.000000000000000e+00},
		{-1.732050807568878e+00, +5.773502691896258e-01, +8.164965809277251e-01, +0.000000000000000e+00},
		{+0.000000000000000e+00, +1.154700538379252e+00, +4.082482904638632e-01, +1.224744871391589e+00},
	})
}

func TestZpotrf01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zpotrf01")

	// a matrices
	a := NewMatrixC([][]complex128{ // must be Hermitian: a = a^H
		{+4 + 0i, 0 + 1i, -3 + 1i, 0 + 2i},
		{+0 - 1i, 3 + 0i, +1 + 0i, 2 + 0i},
		{-3 - 1i, 1 + 0i, +4 + 0i, 1 - 1i},
		{+0 - 2i, 2 + 0i, +1 + 1i, 4 + 0i},
	})
	aUp := NewMatrixC([][]complex128{
		{+4 + 0i, 0 + 1i, -3 + 1i, 0 + 2i},
		{+0 + 0i, 3 + 0i, +1 + 0i, 2 + 0i},
		{+0 + 0i, 0 + 0i, +4 + 0i, 1 - 1i},
		{+0 + 0i, 0 + 0i, +0 + 0i, 4 + 0i},
	})
	aLo := NewMatrixC([][]complex128{
		{+4 + 0i, 0 + 0i, +0 + 0i, 0 + 0i},
		{+0 - 1i, 3 + 0i, +0 + 0i, 0 + 0i},
		{-3 - 1i, 1 + 0i, +4 + 0i, 0 + 0i},
		{+0 - 2i, 2 + 0i, +1 + 1i, 4 + 0i},
	})

	// check aUp and aLo
	checkUploC(tst, "Zherk01", a, aLo, aUp, 1e-17, 1e-17)

	// run zpotrf with up(a)
	n := a.N()
	up := true
	lda := n
	err := Zpotrf(up, n, aUp, lda)
	if err != nil {
		tst.Errorf("Zpotrf failed:\n%v\n", err)
		return
	}

	// check aUp
	chk.MatrixC(tst, "chol(aUp)", 1e-15, aUp.GetSlice(), [][]complex128{
		{+2, +0.000000000000000e+00 + 5.0e-01i, -1.500000000000000e+00 + 5.000000000000000e-01i, +0.000000000000000e+00 + 1.000000000000000e+00i},
		{+0, +1.658312395177700e+00 + 0.0e+00i, +4.522670168666454e-01 - 4.522670168666454e-01i, +9.045340337332909e-01 + 0.000000000000000e+00i},
		{+0, +0.000000000000000e+00 + 0.0e+00i, +1.044465935734187e+00 + 0.000000000000000e+00i, +8.703882797784884e-02 + 8.703882797784884e-02i},
		{+0, +0.000000000000000e+00 + 0.0e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +1.471960144387974e+00 + 0.000000000000000e+00i},
	})

	// run zpotrf with lo(a)
	up = false
	err = Zpotrf(up, n, aLo, lda)
	if err != nil {
		tst.Errorf("Dpotrf failed:\n%v\n", err)
		return
	}

	// check aLo
	chk.MatrixC(tst, "chol(aLo)", 1e-15, aLo.GetSlice(), [][]complex128{
		{+2.0 + 0.0e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00},
		{+0.0 - 5.0e-01i, +1.658312395177700e+00 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00},
		{-1.5 - 5.0e-01i, +4.522670168666454e-01 + 4.522670168666454e-01i, +1.044465935734187e+00 + 0.000000000000000e+00i, +0.000000000000000e+00},
		{+0.0 - 1.0e+00i, +9.045340337332909e-01 + 0.000000000000000e+00i, +8.703882797784884e-02 - 8.703882797784884e-02i, +1.471960144387974e+00},
	})
}
