// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestSpSolver01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpSolver01. real")

	// input matrix data into Triplet
	A := new(Triplet)
	A.Init(5, 5, 13)
	A.Put(0, 0, +1.0) // << duplicated
	A.Put(0, 0, +1.0) // << duplicated
	A.Put(1, 0, +3.0)
	A.Put(0, 1, +3.0)
	A.Put(2, 1, -1.0)
	A.Put(4, 1, +4.0)
	A.Put(1, 2, +4.0)
	A.Put(2, 2, -3.0)
	A.Put(3, 2, +1.0)
	A.Put(4, 2, +2.0)
	A.Put(2, 3, +2.0)
	A.Put(1, 4, +6.0)
	A.Put(4, 4, +1.0)

	// solve
	b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}
	x, err := SpSolve(A, b)
	if err != nil {
		tst.Errorf("%v\n")
		return
	}

	// check
	chk.Array(tst, "x", 1e-14, x, []float64{1, 2, 3, 4, 5})
	TestSolverResidual(tst, A.GetDenseMatrix(), x, b, 1e-13)
}

func TestSpSolver02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpSolver02. complex")

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
	A := new(TripletC)
	A.Init(5, 5, 16) // 5 x 5 matrix with 16 non-zeros

	// first column
	A.Put(0, 0, 19.73+0.00i)
	A.Put(1, 0, +0.00-0.51i)

	// second column
	A.Put(0, 1, 12.11-1.00i)
	A.Put(1, 1, 32.30+7.00i)
	A.Put(2, 1, +0.00-0.51i)

	// third column
	A.Put(0, 2, +0.00+5.0i)
	A.Put(1, 2, 23.07+0.0i)
	A.Put(2, 2, 70.00+7.3i)
	A.Put(3, 2, +1.00+1.1i)

	// fourth column
	A.Put(1, 3, +0.00+1.000i)
	A.Put(2, 3, +3.95+0.000i)
	A.Put(3, 3, 50.17+0.000i)
	A.Put(4, 3, +0.00-9.351i)

	// fifth column
	A.Put(2, 4, 19.00+31.83i)
	A.Put(3, 4, 45.51+0.00i)
	A.Put(4, 4, 55.00+0.00i)

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

	// solve
	x, err := SpSolveC(A, b)
	if err != nil {
		tst.Errorf("%v\n")
		return
	}

	// check
	chk.ArrayC(tst, "x", 1e-3, x, xCorrect)
	TestSolverResidualC(tst, A.GetDenseMatrix(), x, b, 1e-12)
}
