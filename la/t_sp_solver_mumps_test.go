// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/mpi"
)

func switchMPI() {
	if !mpi.IsOn() {
		mpi.Start()
	}
}

func TestSpSolver01aM(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpSolver01aM. real")

	switchMPI()
	comm := mpi.NewCommunicator(nil)

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
	bIsDistr := false
	xCorrect := []float64{1, 2, 3, 4, 5}
	TestSpSolver(tst, "mumps", false, &t, b, xCorrect, 1e-14, 1e-13, false, bIsDistr, comm)
}

func TestSpSolver02M(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpSolver02M. real")

	switchMPI()
	comm := mpi.NewCommunicator(nil)

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
	bIsDistr := false
	xCorrect := []float64{-1, 8, -65, 454, -2725, 13624, -54497, 163490, -326981, 326991}
	TestSpSolver(tst, "mumps", false, &t, b, xCorrect, 1e-4, 1e-9, false, bIsDistr, comm)
}

func TestSpSolver03M(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpSolver03M. complex (without imaginary part)")

	switchMPI()
	comm := mpi.NewCommunicator(nil)

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
	bIsDistr := false
	xCorrect := []complex128{1, 2, 3, 4, 5}
	TestSpSolverC(tst, "mumps", false, &t, b, xCorrect, 1e-14, 1e-13, false, bIsDistr, comm)
}

func TestSpSolver04M(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpSolver04M. complex (without imaginary part)")

	switchMPI()
	comm := mpi.NewCommunicator(nil)

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
	bIsDistr := false
	xCorrect := []complex128{-1, 8, -65, 454, -2725, 13624, -54497, 163490, -326981, 326991}
	TestSpSolverC(tst, "mumps", false, &t, b, xCorrect, 1e-4, 1e-9, false, bIsDistr, comm)
}

func TestSpSolver05M(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpSolver05M. complex")

	switchMPI()
	comm := mpi.NewCommunicator(nil)

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
	bIsDistr := false
	TestSpSolverC(tst, "mumps", false, &t, b, xCorrect, 1e-14, 1e-13, false, bIsDistr, comm)
}

func TestSpSolver06M(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpSolver06M. complex")

	switchMPI()
	comm := mpi.NewCommunicator(nil)

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
	bIsDistr := false
	TestSpSolverC(tst, "mumps", false, &t, b, xCorrect, 1e-3, 1e-12, false, bIsDistr, comm)
}
