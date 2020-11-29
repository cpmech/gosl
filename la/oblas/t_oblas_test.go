// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

func TestMatrix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix01. real")

	A := [][]float64{
		{1, 2, +3, +4},
		{5, 6, +7, +8},
		{9, 0, -1, -2},
	}
	m, n := len(A), len(A[0])

	a := SliceToColMajor(A)
	chk.Array(tst, "A to a", 1e-15, a, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})

	Aback := ColMajorToSlice(m, n, a)
	chk.Deep2(tst, "a to A", 1e-15, Aback, A)

	l := PrintColMajor(m, n, a, "")
	chk.String(tst, l, "1 2 3 4 \n5 6 7 8 \n9 0 -1 -2 ")

	l = PrintColMajorGo(m, n, a, "%2g")
	lCorrect := "[][]float64{\n    { 1, 2, 3, 4},\n    { 5, 6, 7, 8},\n    { 9, 0,-1,-2},\n}"
	chk.String(tst, l, lCorrect)

	l = PrintColMajorPy(m, n, a, "%2g")
	lCorrect = "np.matrix([\n    [ 1, 2, 3, 4],\n    [ 5, 6, 7, 8],\n    [ 9, 0,-1,-2],\n], dtype=float)"
	chk.String(tst, l, lCorrect)
}

func TestMatrix02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix02. complex")

	A := [][]complex128{
		{1 + 0.1i, 2, +3, +4 - 0.4i},
		{5 + 0.5i, 6, +7, +8 - 0.8i},
		{9 + 0.9i, 0, -1, -2 + 1.0i},
	}
	m, n := len(A), len(A[0])

	a := SliceToColMajorC(A)
	chk.ArrayC(tst, "A to a", 1e-15, a, []complex128{1 + 0.1i, 5 + 0.5i, 9 + 0.9i, 2, 6, 0, 3, 7, -1, 4 - 0.4i, 8 - 0.8i, -2 + 1i})

	Aback := ColMajorCtoSlice(m, n, a)
	chk.Deep2c(tst, "a to A", 1e-15, Aback, A)

	l := PrintColMajorC(m, n, a, "%g", "")
	chk.String(tst, l, "1+0.1i, 2+0i, 3+0i, 4-0.4i\n5+0.5i, 6+0i, 7+0i, 8-0.8i\n9+0.9i, 0+0i, -1+0i, -2+1i")

	l = PrintColMajorCgo(m, n, a, "%2g", "%+4.1f")
	lCorrect := "[][]complex128{\n    { 1+0.1i, 2+0.0i, 3+0.0i, 4-0.4i},\n    { 5+0.5i, 6+0.0i, 7+0.0i, 8-0.8i},\n    { 9+0.9i, 0+0.0i,-1+0.0i,-2+1.0i},\n}"
	chk.String(tst, l, lCorrect)

	l = PrintColMajorCpy(m, n, a, "%2g", "%4.1f")
	lCorrect = "np.matrix([\n    [ 1+0.1j, 2+0.0j, 3+0.0j, 4-0.4j],\n    [ 5+0.5j, 6+0.0j, 7+0.0j, 8-0.8j],\n    [ 9+0.9j, 0+0.0j,-1+0.0j,-2+1.0j],\n], dtype=complex)"
	chk.String(tst, l, lCorrect)
}

func TestDdot01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Ddot01")

	x := []float64{20, 10, 30, 123, 123}
	y := []float64{-15, -5, -24, 666, 666, 666}
	n, incx, incy := 3, 1, 1
	res := Ddot(n, x, incx, y, incy)
	chk.Float64(tst, "x⋅y", 1e-15, res, -1070)
}

func TestDscal01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dscal01")

	α := 0.5
	x := []float64{20, 10, 30, 123, 123}
	n, incx := 3, 1
	Dscal(n, α, x, incx)
	chk.Array(tst, "x", 1e-15, x, []float64{10, 5, 15, 123, 123})
}

func TestDaxpy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Daxpy01")

	α := 0.5
	x := []float64{20, 10, 30, 123, 123}
	y := []float64{-15, -5, -24, 666, 666, 666}
	n, incx, incy := 3, 1, 1
	Daxpy(n, α, x, incx, y, incy)

	chk.Array(tst, "x", 1e-15, x, []float64{20, 10, 30, 123, 123})
	chk.Array(tst, "y", 1e-15, y, []float64{-5, 0, -9, 666, 666, 666})
}

func TestZaxpy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zaxpy01")

	α := 1.0 + 0i
	x := []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i}
	y := []complex128{-15 + 1.5i, -5 - 2i, -24 + 1i, 666 - 0.5i, 666 + 5i}
	n, incx, incy := len(x), 1, 1
	Zaxpy(n, α, x, incx, y, incy)

	chk.ArrayC(tst, "x", 1e-15, x, []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i})
	chk.ArrayC(tst, "y", 1e-15, y, []complex128{5 + 2.5i, 5, 6 + 2.5i, 543, 543 + 5.5i})

	α = 0.5 + 1i
	Zaxpy(n, α, x, incx, y, incy)
	chk.ArrayC(tst, "y", 1e-15, y, []complex128{14.0 + 23.i, 8.0 + 11.i, 19.5 + 33.25i, 481.0 - 122.75i, 481.0 - 117.25i})
}

func TestDgemv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemv01")

	// allocate
	a := SliceToColMajor([][]float64{
		{0.1, 0.2, 0.3},
		{1.0, 0.2, 0.3},
		{2.0, 0.2, 0.3},
		{3.0, 0.2, 0.3},
	})
	chk.Array(tst, "a", 1e-15, a, []float64{0.1, 1, 2, 3, 0.2, 0.2, 0.2, 0.2, 0.3, 0.3, 0.3, 0.3})

	// perform mv
	m, n := 4, 3
	α, β := 0.5, 2.0
	x := []float64{20, 10, 30}
	y := []float64{3, 1, 2, 4}
	lda, incx, incy := m, 1, 1
	Dgemv(false, m, n, α, a, lda, x, incx, β, y, incy)
	chk.Array(tst, "y", 1e-15, y, []float64{12.5, 17.5, 29.5, 43.5})

	// perform mv with transpose
	Dgemv(true, m, n, α, a, lda, y, incy, β, x, incx)
	chk.Array(tst, "x", 1e-15, x, []float64{144.125, 30.3, 75.45})

	// check that a is unmodified
	chk.Array(tst, "a", 1e-15, a, []float64{0.1, 1, 2, 3, 0.2, 0.2, 0.2, 0.2, 0.3, 0.3, 0.3, 0.3})
}

func TestZgemv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgemv01")

	// allocate
	a := SliceToColMajorC([][]complex128{
		{0.1 + 3.0i, 0.2, 0.3 - 0.3i},
		{1.0 + 2.0i, 0.2, 0.3 - 0.4i},
		{2.0 + 1.0i, 0.2, 0.3 - 0.5i},
		{3.0 + 0.1i, 0.2, 0.3 - 0.6i},
	})
	m, n := 4, 3
	chk.ArrayC(tst, "a", 1e-15, a, []complex128{0.1 + 3i, 1 + 2i, 2 + 1i, 3 + 0.1i, 0.2, 0.2, 0.2, 0.2, 0.3 - 0.3i, 0.3 - 0.4i, 0.3 - 0.5i, 0.3 - 0.6i})

	// perform mv
	α, β := 0.5+1i, 2.0+1i
	x := []complex128{20, 10, 30}
	y := []complex128{3, 1, 2, 4}
	lda, incx, incy := m, 1, 1
	Zgemv(false, m, n, α, a, lda, x, incx, β, y, incy)
	chk.ArrayC(tst, "y", 1e-15, y, []complex128{-38.5 + 41.5i, -10.5 + 46i, 24.5 + 55.5i, 59.5 + 67i})

	// perform mv with transpose
	Zgemv(true, m, n, α, a, lda, y, incy, β, x, incx)
	chk.ArrayC(tst, "x", 1e-13, x, []complex128{-248.875 + 82.5i, -18.5 + 38i, 83.85 + 154.7i})

	// check that a is unmodified
	chk.ArrayC(tst, "a", 1e-15, a, []complex128{0.1 + 3i, 1 + 2i, 2 + 1i, 3 + 0.1i, 0.2, 0.2, 0.2, 0.2, 0.3 - 0.3i, 0.3 - 0.4i, 0.3 - 0.5i, 0.3 - 0.6i})
}

func TestDger01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dger01: a = 100 + 0.5⋅u⋅vᵀ")

	// allocate matrix
	a := SliceToColMajor([][]float64{ // 4 x 3
		{100, 100, 100},
		{100, 100, 100},
		{100, 100, 100},
		{100, 100, 100},
	})
	u := []float64{1, 2, 3, 4}
	v := []float64{4, 3, 2}

	// sizes
	m := 4 // m = nrow(a) = len(u)
	n := 3 // n = ncol(a) = len(v)

	// run dger
	alpha := 0.5
	lda := 4
	Dger(m, n, alpha, u, 1, v, 1, a, lda)

	// check
	chk.Deep2(tst, "100 + 0.5⋅u⋅vᵀ", 1e-17, ColMajorToSlice(4, 3, a), [][]float64{
		{102, 101.5, 101},
		{104, 103.0, 102},
		{106, 104.5, 103},
		{108, 106.0, 104},
	})
}

func TestDgemm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemm01: 0.5⋅a⋅b + 2⋅c")

	// allocate matrices
	a := SliceToColMajor([][]float64{ // 4 x 5
		{1, 2, +0, 1, -1},
		{2, 3, -1, 1, +1},
		{1, 2, +0, 4, -1},
		{4, 0, +3, 1, +1},
	})
	b := SliceToColMajor([][]float64{ // 5 x 3
		{1, 0, 0},
		{0, 0, 3},
		{0, 0, 1},
		{1, 0, 1},
		{0, 2, 0},
	})
	c := SliceToColMajor([][]float64{ // 4 x 3
		{+0.50, 0, +0.25},
		{+0.25, 0, -0.25},
		{-0.25, 0, +0.00},
		{-0.25, 0, +0.00},
	})

	// sizes
	m := 4 // m = nrow(a) = a.M = nrow(c)
	k := 5 // k = ncol(a) = a.N = nrow(b)
	n := 3 // n = ncol(b) = b.N = ncol(c)

	// run dgemm
	transA, transB := false, false
	alpha, beta := 0.5, 2.0
	lda, ldb, ldc := 4, 5, 4
	Dgemm(transA, transB, m, n, k, alpha, a, lda, b, ldb, beta, c, ldc)

	// check
	chk.Deep2(tst, "0.5⋅a⋅b + 2⋅c", 1e-17, ColMajorToSlice(4, 3, c), [][]float64{
		{2, -1, 4},
		{2, +1, 4},
		{2, -1, 5},
		{2, +1, 2},
	})
}

func TestDgemm02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemm02: 0.5⋅a⋅bᵀ + 2⋅c")

	// allocate matrices
	a := SliceToColMajor([][]float64{ // 4 x 5
		{1, 2, +0, 1, -1},
		{2, 3, -1, 1, +1},
		{1, 2, +0, 4, -1},
		{4, 0, +3, 1, +1},
	})
	b := SliceToColMajor([][]float64{ // 3 x 5
		{1, 0, 0, 1, 0},
		{0, 0, 0, 0, 2},
		{0, 3, 1, 1, 0},
	})
	c := SliceToColMajor([][]float64{ // 4 x 3
		{+0.50, 0, +0.25},
		{+0.25, 0, -0.25},
		{-0.25, 0, +0.00},
		{-0.25, 0, +0.00},
	})

	// sizes
	m := 4 // m = nrow(a)        = a.M = nrow(c)
	k := 5 // k = ncol(a)        = a.N = nrow(trans(b))
	n := 3 // n = ncol(trans(b)) = b.M = ncol(c)

	// run dgemm
	transA, transB := false, true
	alpha, beta := 0.5, 2.0
	lda, ldb, ldc := 4, 3, 4
	Dgemm(transA, transB, m, n, k, alpha, a, lda, b, ldb, beta, c, ldc)

	// check
	chk.Deep2(tst, "0.5⋅a⋅bᵀ + 2⋅c", 1e-17, ColMajorToSlice(4, 3, c), [][]float64{
		{2, -1, 4},
		{2, +1, 4},
		{2, -1, 5},
		{2, +1, 2},
	})
}

func TestDgemm03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemm03: 0.5⋅aᵀ⋅b + 2⋅c")

	// allocate matrices
	a := SliceToColMajor([][]float64{ // 5 x 4
		{+1, +2, +1, +4},
		{+2, +3, +2, +0},
		{+0, -1, +0, +3},
		{+1, +1, +4, +1},
		{-1, +1, -1, +1},
	})
	b := SliceToColMajor([][]float64{ // 5 x 3
		{1, 0, 0},
		{0, 0, 3},
		{0, 0, 1},
		{1, 0, 1},
		{0, 2, 0},
	})
	c := SliceToColMajor([][]float64{ // 4 x 3
		{+0.50, 0, +0.25},
		{+0.25, 0, -0.25},
		{-0.25, 0, +0.00},
		{-0.25, 0, +0.00},
	})

	// sizes
	m := 4 // m = nrow(trans(a)) = a.N = nrow(c)
	k := 5 // k = ncol(trans(a)) = a.M = nrow(trans(b))
	n := 3 // n = ncol(b)        = b.N = ncol(c)

	// run dgemm
	transA, transB := true, false
	alpha, beta := 0.5, 2.0
	lda, ldb, ldc := 5, 5, 4
	Dgemm(transA, transB, m, n, k, alpha, a, lda, b, ldb, beta, c, ldc)

	// check
	chk.Deep2(tst, "0.5⋅aᵀ⋅b + 2⋅c", 1e-17, ColMajorToSlice(4, 3, c), [][]float64{
		{2, -1, 4},
		{2, +1, 4},
		{2, -1, 5},
		{2, +1, 2},
	})
}

func TestDgemm04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemm04: 0.5⋅aᵀ⋅bᵀ + 2⋅c")

	// allocate matrices
	a := SliceToColMajor([][]float64{ // 5 x 4
		{+1, +2, +1, +4},
		{+2, +3, +2, +0},
		{+0, -1, +0, +3},
		{+1, +1, +4, +1},
		{-1, +1, -1, +1},
	})
	b := SliceToColMajor([][]float64{ // 3 x 5
		{1, 0, 0, 1, 0},
		{0, 0, 0, 0, 2},
		{0, 3, 1, 1, 0},
	})
	c := SliceToColMajor([][]float64{ // 4 x 3
		{+0.50, 0, +0.25},
		{+0.25, 0, -0.25},
		{-0.25, 0, +0.00},
		{-0.25, 0, +0.00},
	})

	// sizes
	m := 4 // m = nrow(trans(a)) = a.N = nrow(c)
	k := 5 // k = ncol(trans(a)) = a.M = nrow(trans(b))
	n := 3 // n = ncol(trans(b)) = b.M = ncol(c)

	// run dgemm
	transA, transB := true, true
	alpha, beta := 0.5, 2.0
	lda, ldb, ldc := 5, 3, 4
	Dgemm(transA, transB, m, n, k, alpha, a, lda, b, ldb, beta, c, ldc)

	// check
	chk.Deep2(tst, "0.5⋅aᵀ⋅bᵀ + 2⋅c", 1e-17, ColMajorToSlice(4, 3, c), [][]float64{
		{2, -1, 4},
		{2, +1, 4},
		{2, -1, 5},
		{2, +1, 2},
	})
}

func TestZgemm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgemm01")

	// allocate matrices
	a := SliceToColMajorC([][]complex128{
		{1, 2, +0 + 1i, 1, -1},
		{2, 3, -1 - 1i, 1, +1},
		{1, 2, +0 + 1i, 4, -1},
		{4, 0, +3 - 1i, 1, +1},
	})
	b := SliceToColMajorC([][]complex128{
		{1, 0, 0 + 1i},
		{0, 0, 3 - 1i},
		{0, 0, 1 + 1i},
		{1, 0, 1 - 1i},
		{0, 2, 0 + 1i},
	})
	c := SliceToColMajorC([][]complex128{
		{+0.50, 1i, +0.25},
		{+0.25, 1i, -0.25},
		{-0.25, 1i, +0.00},
		{-0.25, 1i, +0.00},
	})

	// sizes
	m, k := 4, 5
	n := 3

	// run dgemm
	transA, transB := false, false
	alpha, beta := 0.5-2i, 2.0-4i
	lda, ldb, ldc := m, k, m
	Zgemm(transA, transB, m, n, k, alpha, a, lda, b, ldb, beta, c, ldc)

	// check
	chk.Deep2c(tst, "(0.5-2i)⋅a⋅b + (2-4i)⋅c", 1e-17, ColMajorCtoSlice(4, 3, c), [][]complex128{
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
		{2, +3, +0, 0, 0},
		{3, +0, +4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, +0, +1, 0, 0},
		{0, +4, +2, 0, 1},
	}
	n := 5
	a := SliceToColMajor(amat)

	// right-hand-side
	b := []float64{8, 45, -3, 3, 19}

	// solution
	xCorrect := []float64{1, 2, 3, 4, 5}

	// run test
	nrhs := 1
	lda, ldb := n, n
	ipiv := make([]int32, n)
	Dgesv(n, nrhs, a, lda, ipiv, b, ldb)
	chk.Array(tst, "x = A⁻¹ b", 1e-14, b, xCorrect)

	// check ipiv
	chk.Int32s(tst, "ipiv", ipiv, []int32{2, 5, 5, 5, 5})
}

func TestZgesv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgesv01. low accuracy")

	// NOTE: zgesv performs badly with this problem
	//       the best tolerance that can be selected is 0.00038
	//       the same problem happens in python (probably using lapack as well)
	tol := 0.00049

	// matrix
	a := SliceToColMajorC([][]complex128{
		{19.730 + 0.000i, 12.110 - 1.000i, +0.000 + 5.000i, +0.000 + 0.000i, +0.000 + 0.000i},
		{+0.000 - 0.510i, 32.300 + 7.000i, 23.070 + 0.000i, +0.000 + 1.000i, +0.000 + 0.000i},
		{+0.000 + 0.000i, +0.000 - 0.510i, 70.000 + 7.300i, +3.950 + 0.000i, 19.000 + 31.830i},
		{+0.000 + 0.000i, +0.000 + 0.000i, +1.000 + 1.100i, 50.170 + 0.000i, 45.510 + 0.000i},
		{+0.000 + 0.000i, +0.000 + 0.000i, +0.000 + 0.000i, +0.000 - 9.351i, 55.000 + 0.000i},
	})

	// right-hand-side
	b := []complex128{
		+77.38 + 8.82i,
		157.48 + 19.8i,
		1175.62 + 20.69i,
		912.12 - 801.75i,
		550.00 - 1060.4i,
	}

	// solution
	xCorrect := []complex128{
		+3.3 - 1.00i,
		+1.0 + 0.17i,
		+5.5 + 0.00i,
		+9.0 + 0.00i,
		10.0 - 17.75i,
	}

	// run test
	n := 5
	nrhs := 1
	lda, ldb := n, n
	ipiv := make([]int32, n)
	Zgesv(n, nrhs, a, lda, ipiv, b, ldb)
	chk.ArrayC(tst, "x = A⁻¹ b (comparison)", tol, b, xCorrect)

	// compare with python results
	xPython := []complex128{
		3.299687426933794e+00 - 1.000372829305209e+00i,
		9.997606020636992e-01 + 1.698383755401385e-01i,
		5.500074759292877e+00 - 4.556001293922560e-05i,
		8.999787912842375e+00 - 6.662818244209770e-05i,
		1.000001132800243e+01 - 1.774987242230929e+01i,
	}
	chk.ArrayC(tst, "x = A⁻¹ b", 1e-13, b, xPython)

	// check ipiv
	chk.Int32s(tst, "ipiv", ipiv, []int32{1, 2, 3, 4, 5})
}

func checksvd(tst *testing.T, amat, uCorrect, vtCorrect [][]float64, sCorrect []float64, tolu, tols, tolv, tolusv float64) {

	// allocate matrix
	m, n := len(amat), len(amat[0])
	a := SliceToColMajor(amat)

	// compute dimensions
	minMN := utl.Imin(m, n)
	lda := m
	ldu := m
	ldvt := n

	// allocate output arrays
	s := make([]float64, minMN)
	u := make([]float64, m*m)
	vt := make([]float64, n*n)
	superb := make([]float64, minMN)

	// perform SVD
	jobu := 'A'
	jobvt := 'A'
	Dgesvd(jobu, jobvt, m, n, a, lda, s, u, ldu, vt, ldvt, superb)

	// compare results
	umat := ColMajorToSlice(m, m, u)
	vtmat := ColMajorToSlice(n, n, vt)
	if uCorrect != nil {
		chk.Deep2(tst, "u", tolu, umat, uCorrect)
	}
	chk.Array(tst, "s", tols, s, sCorrect)
	if vtCorrect != nil {
		chk.Deep2(tst, "vt", tolv, vtmat, vtCorrect)
	}

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
	chk.Deep2(tst, "u⋅s⋅vt", tolusv, amat, usv)
}

func TestDgesvd01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgesvd01")

	// allocate matrices
	amat := [][]float64{
		{1, 0, 0, 0, 2},
		{0, 0, 3, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 2, 0, 0, 0},
	}
	uCorrect := [][]float64{
		{0, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 0, -1},
		{0, 0, 1, 0},
	}
	sCorrect := []float64{3, math.Sqrt(5.0), 2, 0}
	s2 := math.Sqrt(0.2)
	s8 := math.Sqrt(0.8)
	vtCorrect := [][]float64{
		{0, 0, 1, 0, 0},
		{s2, 0, 0, 0, s8},
		{0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0},
		{-s8, 0, 0, 0, s2},
	}

	// check
	checksvd(tst, amat, uCorrect, vtCorrect, sCorrect, 1e-17, 1e-17, 1e-15, 1e-15)
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

func TestDgesvd03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgesvd03")

	// allocate matrices
	amat := [][]float64{
		{64, 2, 3, 61, 60, 6},
		{9, 55, 54, 12, 13, 51},
		{17, 47, 46, 20, 21, 43},
		{40, 26, 27, 37, 36, 30},
		{32, 34, 35, 29, 28, 38},
		{41, 23, 22, 44, 45, 19},
		{49, 15, 14, 52, 53, 11},
		{8, 58, 59, 5, 4, 62},
	}
	sCorrect := []float64{+2.251695779937001e+02, +1.271865289052834e+02, +1.175789144211322e+01, +1.277237188369868e-14, +6.934703857768031e-15, +5.031833747507930e-15}

	// check
	checksvd(tst, amat, nil, nil, sCorrect, 1e-15, 1e-13, 1e-15, 1e-13)
}

func checksvdC(tst *testing.T, amat, uCorrect, vtCorrect [][]complex128, sCorrect []float64, tolu, tols, tolv, tolusv float64) {

	// allocate matrix
	m, n := len(amat), len(amat[0])
	a := SliceToColMajorC(amat)

	// compute dimensions
	minMN := utl.Imin(m, n)
	lda := m
	ldu := m
	ldvt := n

	// allocate output arrays
	s := make([]float64, minMN)
	u := make([]complex128, m*m)
	vt := make([]complex128, n*n)
	superb := make([]float64, minMN)

	// perform SVD
	jobu := 'A'
	jobvt := 'A'
	Zgesvd(jobu, jobvt, m, n, a, lda, s, u, ldu, vt, ldvt, superb)

	// compare results
	umat := ColMajorCtoSlice(m, m, u)
	vtmat := ColMajorCtoSlice(n, n, vt)
	if uCorrect != nil {
		chk.Deep2c(tst, "u", tolu, umat, uCorrect)
	}
	chk.Array(tst, "s", tols, s, sCorrect)
	if vtCorrect != nil {
		chk.Deep2c(tst, "vt", tolv, vtmat, vtCorrect)
	}

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
	chk.Deep2c(tst, "u⋅s⋅vt", tolusv, amat, usv)
}

func TestZgesvd01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgesvd01")

	// allocate matrices
	amat := [][]complex128{
		{+0.000000000000000e+00 + 0.000000000000000e+00i, +7.071067811865475e-01 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, -7.071067811865475e-01 + 0.000000000000000e+00i},
		{+7.071067811865475e-01 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00 + 7.071067811865475e-01i, +0.000000000000000e+00 + 0.000000000000000e+00i},
		{+0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00 + 7.071067811865475e-01i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00 + 7.071067811865475e-01i},
		{-7.071067811865475e-01 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00 + 7.071067811865475e-01i, +0.000000000000000e+00 + 0.000000000000000e+00i},
	}
	sCorrect := []float64{1, 1, 1, 1}

	// check
	checksvdC(tst, amat, nil, nil, sCorrect, 1e-16, 1e-15, 1e-17, 1e-15)
}

func TestZgesvd02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgesvd02")

	// allocate matrices
	amat := [][]complex128{
		{0, 3, 2, 1},
		{1, 1i, 1i, 1i},
		{2, 2, 2i, 2i},
		{3, 3, 3, 3i},
	}
	sCorrect := []float64{+7.578301582272183e+00, +3.008108139593885e+00, +1.854745532331560e+00, +2.838125418935204e-01}

	// check
	checksvdC(tst, amat, nil, nil, sCorrect, 1e-16, 1e-15, 1e-16, 1e-14)
}

func TestDgetrf01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgetrf01. Dgetrf and Dgetri")

	// matrix
	amat := [][]float64{
		{1, 2, +0, 1},
		{2, 3, -1, 1},
		{1, 2, +0, 4},
		{4, 0, +3, 1},
	}
	m, n := 4, 4
	a := SliceToColMajor(amat)

	// run dgetrf
	lda := m
	ipiv := make([]int32, utl.Imin(m, n))
	Dgetrf(m, n, a, lda, ipiv)

	// check ipiv
	chk.Int32s(tst, "ipiv", ipiv, []int32{4, 2, 3, 4})

	// check LU
	chk.Deep2(tst, "lu", 1e-15, ColMajorToSlice(m, n, a), [][]float64{
		{+4.0e+00, +0.000000000000000e+00, +3.000000000000000e+00, +1.000000000000000e+00},
		{+5.0e-01, +3.000000000000000e+00, -2.500000000000000e+00, +5.000000000000000e-01},
		{+2.5e-01, +6.666666666666666e-01, +9.166666666666665e-01, +3.416666666666667e+00},
		{+2.5e-01, +6.666666666666666e-01, +1.000000000000000e+00, -3.000000000000000e+00},
	})

	// run dgetri
	Dgetri(n, a, lda, ipiv)

	// compare inverse
	ai := ColMajorToSlice(n, m, a)
	chk.Deep2(tst, "inv(a)", 1e-15, ai, [][]float64{
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
				chk.Float64(tst, "diag(a⋅a⁻¹)=diag(I)=1", 1e-15, res, 1)
			} else {
				chk.Float64(tst, "diag(a⋅a⁻¹)=offdiag(I)=0", 1e-15, res, 0)
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
	a := SliceToColMajorC(amat)

	// run
	lda := m
	ipiv := make([]int32, utl.Imin(m, n))
	Zgetrf(m, n, a, lda, ipiv)

	// check ipiv
	chk.Int32s(tst, "ipiv", ipiv, []int32{4, 2, 3, 4})

	// check LU
	chk.Deep2c(tst, "lu", 1e-15, ColMajorCtoSlice(m, n, a), [][]complex128{
		{+4.000000000000000e+00 + 1.000000000000000e+00i, +0.000000000000000e+00, +3.000000000000000e+00 + 0.000000000000000e+00i, +1.000000000000000e+00 - 1.000000000000000e+00i},
		{+5.294117647058824e-01 + 1.176470588235294e-01i, +3.000000000000000e+00, -2.588235294117647e+00 - 3.529411764705882e-01i, +3.529411764705882e-01 - 5.882352941176471e-01i},
		{+2.941176470588235e-01 + 1.764705882352941e-01i, +6.666666666666666e-01, +8.431372549019609e-01 - 2.941176470588235e-01i, +3.294117647058823e+00 - 4.901960784313725e-01i},
		{+2.941176470588235e-01 + 1.764705882352941e-01i, +6.666666666666666e-01, +1.000000000000000e+00 + 0.000000000000000e+00i, -3.000000000000000e+00 + 0.000000000000000e+00i},
	})

	// run zgetri
	Zgetri(n, a, lda, ipiv)

	// compare inverse
	ai := ColMajorCtoSlice(n, m, a)
	chk.Deep2c(tst, "inv(a)", 1e-15, ai, [][]complex128{
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
				chk.Complex128(tst, "diag(a⋅a⁻¹)=diag(I)=1", 1e-15, res, 1)
			} else {
				chk.Complex128(tst, "diag(a⋅a⁻¹)=offdiag(I)=0", 1e-15, res, 0)
			}
		}
	}
}

func checkUplo(tst *testing.T, testname string, n int, c, cLo, cUp []float64, tol float64) {
	maxdiff := 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := c[i+j*n]
			vLo := cLo[i+j*n]
			vUp := cUp[i+j*n]
			if i == j {
				diff := math.Abs(vLo - v)
				if diff > tol {
					maxdiff = diff
				}
				diff = math.Abs(vUp - v)
				if diff > tol {
					maxdiff = diff
				}
			} else {
				diff := math.Abs(vLo + vUp - v)
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
	c := SliceToColMajor([][]float64{
		{+3, +0, -3, +0},
		{+0, +3, +1, +2},
		{-3, +1, +4, +1},
		{+0, +2, +1, +3},
	})
	cUp := SliceToColMajor([][]float64{
		{+3, +0, -3, +0},
		{+0, +3, +1, +2},
		{+0, +0, +4, +1},
		{+0, +0, +0, +3},
	})
	cLo := SliceToColMajor([][]float64{
		{+3, +0, +0, +0},
		{+0, +3, +0, +0},
		{-3, +1, +4, +0},
		{+0, +2, +1, +3},
	})

	// n-size
	n := 4 // c.N

	// check cUp and cLo
	checkUplo(tst, "Dsyrk01", n, c, cLo, cUp, 1e-17)

	// a matrix
	a := SliceToColMajor([][]float64{
		{+1, +2, +1, +1, -1, +0},
		{+2, +2, +1, +0, +0, +0},
		{+3, +1, +3, +1, +2, -1},
		{+1, +0, +1, -1, +0, +0},
	})

	// k-size
	k := 6 // a.N

	// constants
	alpha, beta := 3.0, -1.0

	// run dsyrk with up(c)
	up, trans := true, false
	lda, ldc := n, n
	Dsyrk(up, trans, n, k, alpha, a, lda, beta, cUp, ldc)

	// compare resulting up(c) matrix
	chk.Deep2(tst, "using up(c): c := 3⋅a⋅aᵀ - c", 1e-17, ColMajorToSlice(n, n, cUp), [][]float64{
		{21, 21, 24, +3},
		{+0, 24, 32, +7},
		{+0, +0, 71, 14},
		{+0, +0, +0, +6},
	})

	// run dsyrk with lo(c)
	up = false
	Dsyrk(up, trans, n, k, alpha, a, lda, beta, cLo, ldc)

	// compare resulting up(c) matrix
	chk.Deep2(tst, "using lo(c): c := 3⋅a⋅aᵀ - c", 1e-17, ColMajorToSlice(n, n, cLo), [][]float64{
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
	c := SliceToColMajor([][]float64{
		{+3, 0, -3, 0, 0, 0},
		{+0, 3, +1, 2, 2, 2},
		{-3, 1, +4, 1, 1, 1},
		{+0, 2, +1, 3, 3, 3},
		{+0, 2, +1, 3, 4, 3},
		{+0, 2, +1, 3, 3, 4},
	})
	cUp := SliceToColMajor([][]float64{
		{+3, 0, -3, 0, 0, 0},
		{+0, 3, +1, 2, 2, 2},
		{+0, 0, +4, 1, 1, 1},
		{+0, 0, +0, 3, 3, 3},
		{+0, 0, +0, 0, 4, 3},
		{+0, 0, +0, 0, 0, 4},
	})
	cLo := SliceToColMajor([][]float64{
		{+3, 0, +0, 0, 0, 0},
		{+0, 3, +0, 0, 0, 0},
		{-3, 1, +4, 0, 0, 0},
		{+0, 2, +1, 3, 0, 0},
		{+0, 2, +1, 3, 4, 0},
		{+0, 2, +1, 3, 3, 4},
	})

	// n-size
	n := 6 // c.N

	// check cUp and cLo
	checkUplo(tst, "Dsyrk02", n, c, cLo, cUp, 1e-17)

	// a matrix
	a := SliceToColMajor([][]float64{
		{+1, +2, +1, +1, -1, +0},
		{+2, +2, +1, +0, +0, +0},
		{+3, +1, +3, +1, +2, -1},
		{+1, +0, +1, -1, +0, +0},
	})

	// k-size
	k := 4 // a.M (it's m now)

	// constants
	alpha, beta := 3.0, +1.0

	// run dsyrk with up(c)
	up, trans := true, true
	lda, ldc := k, n
	Dsyrk(up, trans, n, k, alpha, a, lda, beta, cUp, ldc)

	// compare resulting up(c) matrix
	chk.Deep2(tst, "using up(c): c := 3⋅a⋅aᵀ + c", 1e-17, ColMajorToSlice(n, n, cUp), [][]float64{
		{48, 27, 36, +9, 15, -9},
		{+0, 30, 22, 11, +2, -1},
		{+0, +0, 40, 10, 16, -8},
		{+0, +0, +0, 12, +6, +0},
		{+0, +0, +0, +0, 19, -3},
		{+0, +0, +0, +0, +0, +7},
	})

	// run dsyrk with lo(c)
	up = false
	Dsyrk(up, trans, n, k, alpha, a, lda, beta, cLo, ldc)

	// compare resulting up(c) matrix
	chk.Deep2(tst, "using lo(c): c := 3⋅a⋅aᵀ + c", 1e-17, ColMajorToSlice(n, n, cLo), [][]float64{
		{48, +0, +0, +0, +0, +0},
		{27, 30, +0, +0, +0, +0},
		{36, 22, 40, +0, +0, +0},
		{+9, 11, 10, 12, +0, +0},
		{15, +2, 16, +6, 19, +0},
		{-9, -1, -8, +0, -3, +7},
	})
}

func checkUploC(tst *testing.T, testname string, n int, c, cLo, cUp []complex128, tolR, tolI float64) {
	maxdiffR, maxdiffI := 0.0, 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := c[i+j*n]
			vLo := cLo[i+j*n]
			vUp := cUp[i+j*n]
			if i == j {
				diffR := math.Abs(real(vLo) - real(v))
				diffI := math.Abs(imag(vLo) - imag(v))
				if diffR > tolR {
					maxdiffR = diffR
				}
				if diffI > tolI {
					maxdiffI = diffI
				}
				diffR = math.Abs(real(vUp) - real(v))
				diffI = math.Abs(imag(vUp) - imag(v))
				if diffR > tolR {
					maxdiffR = diffR
				}
				if diffI > tolI {
					maxdiffI = diffI
				}
			} else {
				diffR := math.Abs(real(vLo) + real(vUp) - real(v))
				diffI := math.Abs(imag(vLo) + imag(vUp) - imag(v))
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
	c := SliceToColMajorC([][]complex128{
		{+3 + 1i, 0 + 0i, -2 + 0i, 0 + 0i},
		{-1 + 0i, 3 + 0i, +0 + 0i, 2 + 0i},
		{-4 + 0i, 1 + 0i, +3 + 0i, 1 + 0i},
		{-1 + 0i, 2 + 0i, +0 + 0i, 3 - 1i},
	})
	cUp := SliceToColMajorC([][]complex128{
		{+3 + 1i, 0 + 0i, -2 + 0i, 0 + 0i},
		{+0 + 0i, 3 + 0i, +0 + 0i, 2 + 0i},
		{+0 + 0i, 0 + 0i, +3 + 0i, 1 + 0i},
		{+0 + 0i, 0 + 0i, +0 + 0i, 3 - 1i},
	})
	cLo := SliceToColMajorC([][]complex128{
		{+3 + 1i, 0 + 0i, +0 + 0i, 0 + 0i},
		{-1 + 0i, 3 + 0i, +0 + 0i, 0 + 0i},
		{-4 + 0i, 1 + 0i, +3 + 0i, 0 + 0i},
		{-1 + 0i, 2 + 0i, +0 + 0i, 3 - 1i},
	})

	// n-size
	n := 4 // c.N

	// check cUp and cLo
	checkUploC(tst, "Zsyrk02", n, c, cLo, cUp, 1e-17, 1e-17)

	// a matrix
	a := SliceToColMajorC([][]complex128{
		{+1 - 1i, +2, +1, +1, -1, +0 + 0i},
		{+2 + 0i, +2, +1, +0, +0, +0 + 1i},
		{+3 + 1i, +1, +3, +1, +2, -1 + 0i},
		{+1 + 0i, +0, +1, -1, +0, +0 + 1i},
	})

	// sizes
	k := 6 // a.N

	// constants
	alpha, beta := 3.0+0i, +1.0+0i

	// run zsyrk with up(c)
	up, trans := true, false
	lda, ldc := n, n
	Zsyrk(up, trans, n, k, alpha, a, lda, beta, cUp, ldc)

	// compare resulting up(c) matrix
	chk.Deep2c(tst, "using up(c): c := 3⋅a⋅aᵀ + c", 1e-17, ColMajorCtoSlice(n, n, cUp), [][]complex128{
		{24 - 5i, 21 - 6i, 22 - 6i, +3 - 3i},
		{+0 + 0i, 27 + 0i, 33 + 3i, +8 + 0i},
		{+0 + 0i, +0 + 0i, 75 + 18i, 16 + 0i},
		{+0 + 0i, +0 + 0i, +0 + 0i, +9 - 1i},
	})

	// run zsyrk with lo(c)
	up = false
	Zsyrk(up, trans, n, k, alpha, a, lda, beta, cLo, ldc)

	// compare resulting up(c) matrix
	chk.Deep2c(tst, "using lo(c): c := 3⋅a⋅aᵀ + c", 1e-17, ColMajorCtoSlice(n, n, cLo), [][]complex128{
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
	c := SliceToColMajorC([][]complex128{ // must be Hermitian: c = c^H
		{+4 + 0i, 0 + 1i, -3 + 1i, 0 + 2i},
		{+0 - 1i, 3 + 0i, +1 + 0i, 2 + 0i},
		{-3 - 1i, 1 + 0i, +4 + 0i, 1 - 1i},
		{+0 - 2i, 2 + 0i, +1 + 1i, 4 + 0i},
	})
	cUp := SliceToColMajorC([][]complex128{
		{+4 + 0i, 0 + 1i, -3 + 1i, 0 + 2i},
		{+0 + 0i, 3 + 0i, +1 + 0i, 2 + 0i},
		{+0 + 0i, 0 + 0i, +4 + 0i, 1 - 1i},
		{+0 + 0i, 0 + 0i, +0 + 0i, 4 + 0i},
	})
	cLo := SliceToColMajorC([][]complex128{
		{+4 + 0i, 0 + 0i, +0 + 0i, 0 + 0i},
		{+0 - 1i, 3 + 0i, +0 + 0i, 0 + 0i},
		{-3 - 1i, 1 + 0i, +4 + 0i, 0 + 0i},
		{+0 - 2i, 2 + 0i, +1 + 1i, 4 + 0i},
	})

	// n-size
	n := 4 // c.N

	// check cUp and cLo
	checkUploC(tst, "Zherk01", n, c, cLo, cUp, 1e-17, 1e-17)

	// a matrix
	a := SliceToColMajorC([][]complex128{
		{+1 - 1i, +2, +1, +1, -1, +0 + 0i},
		{+2 + 0i, +2, +1, +0, +0, +0 + 1i},
		{+3 + 1i, +1, +3, +1, +2, -1 + 0i},
		{+1 + 0i, +0, +1, -1, +0, +0 + 1i},
	})

	// sizes
	k := 6 // a.N

	// constants
	alpha, beta := 3.0, +1.0

	// run zherk with up(c)
	up, trans := true, false
	lda, ldc := n, n
	Zherk(up, trans, n, k, alpha, a, lda, beta, cUp, ldc)

	// compare resulting up(c) matrix
	chk.Deep2c(tst, "using up(c): c := 3⋅a⋅aᵀ + c", 1e-17, ColMajorCtoSlice(n, n, cUp), [][]complex128{
		{31 + 0i, 21 - 5i, 15 - 11i, 3 - 1i},
		{+0 + 0i, 33 + 0i, 34 - 9i, 14 + 0i},
		{+0 + 0i, +0 + 0i, 82 + 0i, 16 + 5i},
		{+0 + 0i, +0 + 0i, +0 + 0i, 16 + 0i},
	})

	// run zherk with lo(c)
	up = false
	Zherk(up, trans, n, k, alpha, a, lda, beta, cLo, ldc)

	// compare resulting up(c) matrix
	chk.Deep2c(tst, "using lo(c): c := 3⋅a⋅aᵀ + c", 1e-17, ColMajorCtoSlice(n, n, cLo), [][]complex128{
		{31 + 0i, +0 + 0i, +0 + 0i, +0 + 0i},
		{21 + 5i, 33 + 0i, +0 + 0i, +0 + 0i},
		{15 + 11i, 34 + 9i, 82 + 0i, +0 + 0i},
		{+3 + 1i, 14 + 0i, 16 - 5i, 16 + 0i},
	})
}

func TestDpotrf01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dpotrf01")

	// a matrices
	a := SliceToColMajor([][]float64{
		{+3, +0, -3, +0},
		{+0, +3, +1, +2},
		{-3, +1, +4, +1},
		{+0, +2, +1, +3},
	})
	aUp := SliceToColMajor([][]float64{
		{+3, +0, -3, +0},
		{+0, +3, +1, +2},
		{+0, +0, +4, +1},
		{+0, +0, +0, +3},
	})
	aLo := SliceToColMajor([][]float64{
		{+3, +0, +0, +0},
		{+0, +3, +0, +0},
		{-3, +1, +4, +0},
		{+0, +2, +1, +3},
	})

	// n-size
	n := 4 // a.N

	// check aUp and aLo
	checkUplo(tst, "Dpotrf01", n, a, aLo, aUp, 1e-17)

	// run dpotrf with up(a)
	up := true
	lda := n
	Dpotrf(up, n, aUp, lda)

	// check aUp
	chk.Deep2(tst, "chol(aUp)", 1e-15, ColMajorToSlice(n, n, aUp), [][]float64{
		{+1.732050807568877e+00, +0.000000000000000e+00, -1.732050807568878e+00, +0.000000000000000e+00},
		{+0.000000000000000e+00, +1.732050807568877e+00, +5.773502691896258e-01, +1.154700538379252e+00},
		{+0.000000000000000e+00, +0.000000000000000e+00, +8.164965809277251e-01, +4.082482904638632e-01},
		{+0.000000000000000e+00, +0.000000000000000e+00, +0.000000000000000e+00, +1.224744871391589e+00},
	})

	// run dpotrf with lo(a)
	up = false
	Dpotrf(up, n, aLo, lda)

	// check aLo
	chk.Deep2(tst, "chol(aLo)", 1e-15, ColMajorToSlice(n, n, aLo), [][]float64{
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
	a := SliceToColMajorC([][]complex128{ // must be Hermitian: a = a^H
		{+4 + 0i, 0 + 1i, -3 + 1i, 0 + 2i},
		{+0 - 1i, 3 + 0i, +1 + 0i, 2 + 0i},
		{-3 - 1i, 1 + 0i, +4 + 0i, 1 - 1i},
		{+0 - 2i, 2 + 0i, +1 + 1i, 4 + 0i},
	})
	aUp := SliceToColMajorC([][]complex128{
		{+4 + 0i, 0 + 1i, -3 + 1i, 0 + 2i},
		{+0 + 0i, 3 + 0i, +1 + 0i, 2 + 0i},
		{+0 + 0i, 0 + 0i, +4 + 0i, 1 - 1i},
		{+0 + 0i, 0 + 0i, +0 + 0i, 4 + 0i},
	})
	aLo := SliceToColMajorC([][]complex128{
		{+4 + 0i, 0 + 0i, +0 + 0i, 0 + 0i},
		{+0 - 1i, 3 + 0i, +0 + 0i, 0 + 0i},
		{-3 - 1i, 1 + 0i, +4 + 0i, 0 + 0i},
		{+0 - 2i, 2 + 0i, +1 + 1i, 4 + 0i},
	})

	// n-size
	n := 4 // a.N

	// check aUp and aLo
	checkUploC(tst, "Zherk01", n, a, aLo, aUp, 1e-17, 1e-17)

	// run zpotrf with up(a)
	up := true
	lda := n
	Zpotrf(up, n, aUp, lda)

	// check aUp
	chk.Deep2c(tst, "chol(aUp)", 1e-15, ColMajorCtoSlice(n, n, aUp), [][]complex128{
		{+2, +0.000000000000000e+00 + 5.0e-01i, -1.500000000000000e+00 + 5.000000000000000e-01i, +0.000000000000000e+00 + 1.000000000000000e+00i},
		{+0, +1.658312395177700e+00 + 0.0e+00i, +4.522670168666454e-01 - 4.522670168666454e-01i, +9.045340337332909e-01 + 0.000000000000000e+00i},
		{+0, +0.000000000000000e+00 + 0.0e+00i, +1.044465935734187e+00 + 0.000000000000000e+00i, +8.703882797784884e-02 + 8.703882797784884e-02i},
		{+0, +0.000000000000000e+00 + 0.0e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +1.471960144387974e+00 + 0.000000000000000e+00i},
	})

	// run zpotrf with lo(a)
	up = false
	Zpotrf(up, n, aLo, lda)

	// check aLo
	chk.Deep2c(tst, "chol(aLo)", 1e-15, ColMajorCtoSlice(n, n, aLo), [][]complex128{
		{+2.0 + 0.0e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00},
		{+0.0 - 5.0e-01i, +1.658312395177700e+00 + 0.000000000000000e+00i, +0.000000000000000e+00 + 0.000000000000000e+00i, +0.000000000000000e+00},
		{-1.5 - 5.0e-01i, +4.522670168666454e-01 + 4.522670168666454e-01i, +1.044465935734187e+00 + 0.000000000000000e+00i, +0.000000000000000e+00},
		{+0.0 - 1.0e+00i, +9.045340337332909e-01 + 0.000000000000000e+00i, +8.703882797784884e-02 - 8.703882797784884e-02i, +1.471960144387974e+00},
	})
}

func TestDgeev01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgeev01")

	adeep2 := [][]float64{
		{+0.35, +0.45, -0.14, -0.17},
		{+0.09, +0.07, -0.54, +0.35},
		{-0.44, -0.33, -0.03, +0.17},
		{+0.25, -0.32, -0.13, +0.11},
	}
	a := SliceToColMajor(adeep2)

	n := 4
	lda := n

	wr := make([]float64, n)   // eigen values (real part)
	wi := make([]float64, n)   // eigen values (imaginary part)
	vl := make([]float64, n*n) // left eigenvectors
	vr := make([]float64, n*n) // right eigenvectors

	ldvl := n
	ldvr := n

	calcVl := true
	calcVr := true

	Dgeev(calcVl, calcVr, n, a, lda, wr, wi, vl, ldvl, vr, ldvr)

	vvl := make([]complex128, n*n)
	vvr := make([]complex128, n*n)
	EigenvecsBuildBoth(vvl, vvr, wr, wi, vl, vr)

	// check eigenvalues
	wRef := []complex128{
		+7.994821225862098e-01,
		-9.941245329507467e-02 + 4.007924719897546e-01i,
		-9.941245329507467e-02 - 4.007924719897546e-01i,
		-1.006572159960587e-01,
	}
	ww := GetJoinComplex(wr, wi)
	chk.ArrayC(tst, "w", 4.5e-16, ww, wRef)

	// check left eigenvectors
	vl0Ref := []complex128{
		-6.244707486379453e-01,
		-5.994889025288728e-01,
		+4.999156725721429e-01,
		-2.708616172576073e-02,
	}
	vl1Ref := []complex128{
		+5.330229831716200e-01,
		-2.666163325181558e-01 + 4.041362636762622e-01i,
		+3.455257668600027e-01 + 3.152853126680209e-01i,
		-2.540814367391268e-01 - 4.451133008385643e-01i,
	}
	vl2Ref := []complex128{
		+5.330229831716200e-01,
		-2.666163325181558e-01 - 4.041362636762622e-01i,
		+3.455257668600027e-01 - 3.152853126680209e-01i,
		-2.540814367391268e-01 + 4.451133008385643e-01i,
	}
	vl3Ref := []complex128{
		+6.641410231734539e-01,
		-1.068153340034493e-01,
		+7.293254091191846e-01,
		+1.248664621625170e-01,
	}
	chk.ArrayC(tst, "vl0", 1e-15, ExtractColC(0, n, n, vvl), vl0Ref)
	chk.ArrayC(tst, "vl1", 1e-15, ExtractColC(1, n, n, vvl), vl1Ref)
	chk.ArrayC(tst, "vl2", 1e-15, ExtractColC(2, n, n, vvl), vl2Ref)
	chk.ArrayC(tst, "vl3", 1e-15, ExtractColC(3, n, n, vvl), vl3Ref)

	// check right eigenvectors
	vr0Ref := []complex128{
		-6.550887675124076e-01,
		-5.236294609021240e-01,
		+5.362184613722345e-01,
		-9.560677820122976e-02,
	}
	vr1Ref := []complex128{
		-1.933015482642217e-01 + 2.546315719275843e-01i,
		+2.518565317267399e-01 - 5.224047347116287e-01i,
		+9.718245844328152e-02 - 3.083837558972283e-01i,
		+6.759540542547480e-01,
	}
	vr2Ref := []complex128{
		-1.933015482642217e-01 - 2.546315719275843e-01i,
		+2.518565317267399e-01 + 5.224047347116287e-01i,
		+9.718245844328152e-02 + 3.083837558972283e-01i,
		+6.759540542547480e-01,
	}
	vr3Ref := []complex128{
		+1.253326972309026e-01,
		+3.320222155717508e-01,
		+5.938377595573312e-01,
		+7.220870298624550e-01,
	}
	chk.ArrayC(tst, "vr0", 1e-15, ExtractColC(0, n, n, vvr), vr0Ref)
	chk.ArrayC(tst, "vr1", 1e-15, ExtractColC(1, n, n, vvr), vr1Ref)
	chk.ArrayC(tst, "vr2", 1e-15, ExtractColC(2, n, n, vvr), vr2Ref)
	chk.ArrayC(tst, "vr3", 1e-15, ExtractColC(3, n, n, vvr), vr3Ref)

	// call Dgeev again without vr
	a2 := SliceToColMajor(adeep2)
	wr2 := make([]float64, n)   // eigen values (real part)
	wi2 := make([]float64, n)   // eigen values (imaginary part)
	vl2 := make([]float64, n*n) // left eigenvectors
	calcVl = true
	calcVr = false
	Dgeev(calcVl, calcVr, n, a2, lda, wr2, wi2, vl2, ldvl, nil, 0)

	// check eigenvalues and left eigenvectors
	vvl2 := make([]complex128, n*n)
	EigenvecsBuild(vvl2, wr2, wi2, vl2)
	ww2 := GetJoinComplex(wr2, wi2)
	chk.ArrayC(tst, "2: w", 4.5e-16, ww2, wRef)
	chk.ArrayC(tst, "2: vl0", 1e-15, ExtractColC(0, n, n, vvl2), vl0Ref)
	chk.ArrayC(tst, "2: vl1", 1e-15, ExtractColC(1, n, n, vvl2), vl1Ref)
	chk.ArrayC(tst, "2: vl2", 1e-15, ExtractColC(2, n, n, vvl2), vl2Ref)
	chk.ArrayC(tst, "2: vl3", 1e-15, ExtractColC(3, n, n, vvl2), vl3Ref)

	// call Dgeev again without vl
	a3 := SliceToColMajor(adeep2)
	wr3 := make([]float64, n)   // eigen values (real part)
	wi3 := make([]float64, n)   // eigen values (imaginary part)
	vr3 := make([]float64, n*n) // right eigenvectors
	calcVl = false
	calcVr = true
	Dgeev(calcVl, calcVr, n, a3, lda, wr3, wi3, nil, 0, vr3, ldvr)

	// check eigenvalues and right eigenvectors
	vvr3 := make([]complex128, n*n)
	EigenvecsBuild(vvr3, wr3, wi3, vr3)
	ww3 := GetJoinComplex(wr3, wi3)
	chk.ArrayC(tst, "3: w", 4.5e-16, ww3, wRef)
	chk.ArrayC(tst, "3: vr0", 1e-15, ExtractColC(0, n, n, vvr3), vr0Ref)
	chk.ArrayC(tst, "3: vr1", 1e-15, ExtractColC(1, n, n, vvr3), vr1Ref)
	chk.ArrayC(tst, "3: vr3", 1e-15, ExtractColC(3, n, n, vvr3), vr3Ref)
	chk.ArrayC(tst, "3: vr3", 1e-15, ExtractColC(3, n, n, vvr3), vr3Ref)

	// call Dgeev again without eigenvectors
	a4 := SliceToColMajor(adeep2)
	wr4 := make([]float64, n) // eigen values (real part)
	wi4 := make([]float64, n) // eigen values (imaginary part)
	calcVl = false
	calcVr = false
	Dgeev(calcVl, calcVr, n, a4, lda, wr4, wi4, nil, 0, nil, 0)

	// check eigenvalues
	ww4 := GetJoinComplex(wr4, wi4)
	chk.ArrayC(tst, "4: w", 4.5e-16, ww4, wRef)
}
