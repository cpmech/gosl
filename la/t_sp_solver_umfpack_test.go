// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestSpUmfpack01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack01a. real")

	// input matrix data into Triplet
	var t Triplet
	t.Init(5, 5, 13)
	t.Put(0, 0, +1.0) // << duplicated
	t.Put(0, 0, +1.0) // << duplicated
	t.Put(1, 0, +3.0)
	t.Put(0, 1, +3.0)
	t.Put(2, 1, -1.0)
	t.Put(4, 1, +4.0)
	t.Put(1, 2, +4.0)
	t.Put(2, 2, -3.0)
	t.Put(3, 2, +1.0)
	t.Put(4, 2, +2.0)
	t.Put(2, 3, +2.0)
	t.Put(1, 4, +6.0)
	t.Put(4, 4, +1.0)

	// run test
	b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}
	xCorrect := []float64{1, 2, 3, 4, 5}
	TestSpSolver(tst, "umfpack", false, &t, b, xCorrect, 1e-14, 1e-13, chk.Verbose, false, nil)
}

func TestSpUmfpack01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack01b. real. go-routines")

	// input matrix data into Triplet
	var t Triplet
	t.Init(5, 5, 13)
	t.Put(0, 0, +1.0) // << duplicated
	t.Put(0, 0, +1.0) // << duplicated
	t.Put(1, 0, +3.0)
	t.Put(0, 1, +3.0)
	t.Put(2, 1, -1.0)
	t.Put(4, 1, +4.0)
	t.Put(1, 2, +4.0)
	t.Put(2, 2, -3.0)
	t.Put(3, 2, +1.0)
	t.Put(4, 2, +2.0)
	t.Put(2, 3, +2.0)
	t.Put(1, 4, +6.0)
	t.Put(4, 4, +1.0)

	// run test
	b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}
	xCorrect := []float64{1, 2, 3, 4, 5}
	nch := 2
	done := make(chan int, nch)
	for i := 0; i < nch; i++ {
		go func() {
			TestSpSolver(tst, "umfpack", false, &t, b, xCorrect, 1e-14, 1e-13, false, false, nil)
			done <- 1
		}()
	}
	for i := 0; i < nch; i++ {
		<-done
	}
}

func TestSpUmfpack02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack02. real")

	// input matrix data into Triplet
	var t Triplet
	t.Init(10, 10, 64)
	for i := 0; i < 10; i++ {
		j := i
		if i > 0 {
			j = i - 1
		}
		for ; j < 10; j++ {
			val := 10.0 - float64(j)
			if i > j {
				val -= 1.0
			}
			t.Put(i, j, val)
		}
	}

	// run test
	b := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	xCorrect := []float64{-1, 8, -65, 454, -2725, 13624, -54497, 163490, -326981, 326991}
	tolx := 1e-4 // << for macOS, for Linux and Windows, tol = 1e-5
	tol := 1e-9  // TODO: check why tests fails with 1e-10 @ office but not @ home
	TestSpSolver(tst, "umfpack", false, &t, b, xCorrect, tolx, tol, false, false, nil)
}

func TestSpUmfpack03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack03. complex (without imaginary part)")

	// input matrix data into Triplet
	var t TripletC
	t.Init(5, 5, 13)
	t.Put(0, 0, +1.0+0i) // << duplicated
	t.Put(0, 0, +1.0+0i) // << duplicated
	t.Put(1, 0, +3.0+0i)
	t.Put(0, 1, +3.0+0i)
	t.Put(2, 1, -1.0+0i)
	t.Put(4, 1, +4.0+0i)
	t.Put(1, 2, +4.0+0i)
	t.Put(2, 2, -3.0+0i)
	t.Put(3, 2, +1.0+0i)
	t.Put(4, 2, +2.0+0i)
	t.Put(2, 3, +2.0+0i)
	t.Put(1, 4, +6.0+0i)
	t.Put(4, 4, +1.0+0i)

	// run test
	b := []complex128{8.0, 45.0, -3.0, 3.0, 19.0}
	xCorrect := []complex128{1, 2, 3, 4, 5}
	TestSpSolverC(tst, "umfpack", false, &t, b, xCorrect, 1e-14, 1e-13, true, false, nil)
}

func TestSpUmfpack04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack04. complex (without imaginary part)")

	// input matrix data into Triplet
	var t TripletC
	t.Init(10, 10, 64)
	for i := 0; i < 10; i++ {
		j := i
		if i > 0 {
			j = i - 1
		}
		for ; j < 10; j++ {
			val := 10.0 - float64(j)
			if i > j {
				val -= 1.0
			}
			t.Put(i, j, complex(val, 0))
		}
	}

	// run test
	b := []complex128{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	xCorrect := []complex128{-1, 8, -65, 454, -2725, 13624, -54497, 163490, -326981, 326991}
	TestSpSolverC(tst, "umfpack", false, &t, b, xCorrect, 2.7e-5, 1e-9, true, false, nil)
}

func TestSpUmfpack05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack05. complex")

	// data
	n := 10
	b := make([]complex128, n)
	xCorrect := make([]complex128, n)

	// input matrix data into Triplet
	var t TripletC
	t.Init(n, n, n)
	for i := 0; i < n; i++ {

		// Some very fake diagonals. Should take exactly 20 GMRES steps
		ar := 10.0 + float64(i)/(float64(n)/10.0)
		ac := 10.0 - float64(i)/(float64(n)/10.0)
		t.Put(i, i, complex(ar, ac))

		// Let exact solution = 1 + 0.5i
		xCorrect[i] = complex(float64(i+1), float64(i+1)/10.0)

		// Generate RHS to match exact solution
		b[i] = complex(ar*real(xCorrect[i])-ac*imag(xCorrect[i]),
			ar*imag(xCorrect[i])+ac*real(xCorrect[i]))
	}

	// run test
	TestSpSolverC(tst, "umfpack", false, &t, b, xCorrect, 1e-15, 1e-13, true, false, nil)
}

func TestSpUmfpack06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack06. complex")

	// given the following matrix of complex numbers:
	//      _                                                  _
	//     |  19.73    12.11-i      5i        0          0      |
	//     |  -0.51i   32.3+7i    23.07       i          0      |
	// A = |    0      -0.51i    70+7.3i     3.95    19+31.83i  |
	//     |    0        0        1+1.1i    50.17      45.51    |
	//     |_   0        0          0      -9.351i       55    _|
	//
	// and the following vector:
	//      _                  _
	//     |    77.38+8.82i     |
	//     |   157.48+19.8i     |
	// b = |  1175.62+20.69i    |
	//     |   912.12-801.75i   |
	//     |_     550-1060.4i  _|
	//
	// solve:
	//         A.x = b
	//
	// the solution is:
	//      _            _
	//     |     3.3-i    |
	//     |    1+0.17i   |
	// x = |      5.5     |
	//     |       9      |
	//     |_  10-17.75i _|

	// input matrix in Complex Triplet format
	var t TripletC
	t.Init(5, 5, 16) // 5 x 5 matrix with 16 non-zeros

	// first column
	t.Put(0, 0, 19.73+0.00i)
	t.Put(1, 0, +0.00-0.51i)

	// second column
	t.Put(0, 1, 12.11-1.00i)
	t.Put(1, 1, 32.30+7.00i)
	t.Put(2, 1, +0.00-0.51i)

	// third column
	t.Put(0, 2, +0.00+5.0i)
	t.Put(1, 2, 23.07+0.0i)
	t.Put(2, 2, 70.00+7.3i)
	t.Put(3, 2, +1.00+1.1i)

	// fourth column
	t.Put(1, 3, +0.00+1.000i)
	t.Put(2, 3, +3.95+0.000i)
	t.Put(3, 3, 50.17+0.000i)
	t.Put(4, 3, +0.00-9.351i)

	// fifth column
	t.Put(2, 4, 19.00+31.83i)
	t.Put(3, 4, 45.51+0.00i)
	t.Put(4, 4, 55.00+0.00i)

	// right-hand-side
	b := []complex128{
		+77.38 + 8.82i,
		+157.48 + 19.8i,
		1175.62 + 20.69i,
		+912.12 - 801.75i,
		+550.00 - 1060.4i,
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
	TestSpSolverC(tst, "umfpack", false, &t, b, xCorrect, 1e-3, 1e-12, true, false, nil)
}

func TestSpUmfpack07a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack07a. large matrix")

	// matrix and rhs
	N := 512
	A := new(Triplet)
	b := NewVector(N)

	// generator
	gen := func() {
		if A.Max() == 0 {
			A.Init(N, N, N*N)
		}
		A.Start()
		for i := 0; i < N; i++ {
			j := i
			if i > 0 {
				j = i - 1
			}
			for ; j < N; j++ {
				val := float64(N) - float64(j)
				if i > j {
					val -= 1.0
				}
				A.Put(i, j, val)
			}
			b[i] = float64(i)
		}
	}

	// allocate solver
	var sol Umfpack
	defer sol.Free()

	// loop
	for k := 0; k < 3; k++ {

		// generate matrix
		gen()

		// initialise solver
		if k == 0 {
			sol.Init(A, nil)
		}

		// factorise
		sol.Fact()

		// solve
		x := NewVector(len(b))
		sol.Solve(x, b, false) // x := inv(A) * b
	}
}

func TestSpUmfpack07b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpUmfpack07b. large matrix (complex)")

	// matrix and rhs
	N := 512
	A := new(TripletC)
	b := NewVectorC(N)

	// generator
	gen := func() {
		if A.Max() == 0 {
			A.Init(N, N, N*N)
		}
		A.Start()
		for i := 0; i < N; i++ {
			j := i
			if i > 0 {
				j = i - 1
			}
			for ; j < N; j++ {
				val := float64(N) - float64(j)
				if i > j {
					val -= 1.0
				}
				A.Put(i, j, complex(val, 0))
			}
			b[i] = complex(float64(i), 0)
		}
	}

	// allocate solver
	var sol UmfpackC
	defer sol.Free()

	// loop
	for k := 0; k < 3; k++ {

		// generate matrix
		gen()
		chk.Int(tst, "pos", A.pos, (N*N-N)/2+N+N-1)

		// initialise solver
		if k == 0 {
			sol.Init(A, nil)
		}

		// factorise
		sol.Fact()

		// solve
		x := NewVectorC(len(b))
		sol.Solve(x, b, false) // x := inv(A) * b
	}
}
