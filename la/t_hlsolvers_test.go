// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math/cmplx"
	"testing"

	"github.com/cpmech/gosl/chk"
)

func checkResid(tst *testing.T, a *Matrix, x, b Vector, tolNorm float64) {
	r := NewVector(len(x))
	r.Apply(-1, b)           // r := -b
	MatVecMulAdd(r, 1, a, x) // r += 1*a*x
	resid := r.Norm()
	if resid > tolNorm {
		tst.Errorf("residual is too large: %g\n", resid)
		return
	}
}

func checkResidC(tst *testing.T, a *MatrixC, x, b VectorC, tolNorm float64) {
	r := NewVectorC(len(x))
	r.Apply(-1, b)            // r = -b
	MatVecMulAddC(r, 1, a, x) // r += 1*a*x
	resid := cmplx.Abs(r.Norm())
	if resid > tolNorm {
		tst.Errorf("residual is too large: %g\n", resid)
		return
	}
}

func TestHLsolver01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("HLsolver01. real")

	// sparse matrix
	var A Triplet
	A.Init(5, 5, 13)
	A.Put(0, 0, +1.0) // << repeated
	A.Put(0, 0, +1.0) // << repeated
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

	// right-hand-side
	b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}

	// solve
	x, err := SolveRealLinSys(&A, b)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check
	xCorrect := []float64{1, 2, 3, 4, 5}
	chk.Vector(tst, "x", 1e-14, x, xCorrect)
	checkResid(tst, A.GetDenseMatrix(), x, b, 1e-13)
}

func TestHLsolver02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("HLsolver02. complex")

	// flag indicating to store (real,complex) values in monolithic form => 1D array
	xzmono := false

	// input matrix in Complex Triplet format
	var A TripletC
	A.Init(5, 5, 16, xzmono) // 5 x 5 matrix with 16 non-zeros

	// first column
	A.Put(0, 0, 19.73, 0) // i=0, j=0, real=19.73, complex=0
	A.Put(1, 0, 0, -0.51) // i=1, j=0, real=0, complex=-0.51

	// second column
	A.Put(0, 1, 12.11, -1) // i=0, j=1, real=12.11, complex=-1
	A.Put(1, 1, 32.3, 7)
	A.Put(2, 1, 0, -0.51)

	// third column
	A.Put(0, 2, 0, 5)
	A.Put(1, 2, 23.07, 0)
	A.Put(2, 2, 70, 7.3)
	A.Put(3, 2, 1, 1.1)

	// fourth column
	A.Put(1, 3, 0, 1)
	A.Put(2, 3, 3.95, 0)
	A.Put(3, 3, 50.17, 0)
	A.Put(4, 3, 0, -9.351)

	// fifth column
	A.Put(2, 4, 19, 31.83)
	A.Put(3, 4, 45.51, 0)
	A.Put(4, 4, 55, 0)

	// right-hand-side
	b := []complex128{
		+77.38 + 8.82i,
		157.48 + 19.8i,
		1175.62 + 20.69i,
		912.12 - 801.75i,
		550.00 - 1060.4i,
	}

	// solve
	x, err := SolveComplexLinSys(&A, b)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check
	x_correct := []complex128{
		+3.3 - 1.00i,
		+1.0 + 0.17i,
		+5.5 + 0.00i,
		+9.0 + 0.00i,
		10.0 - 17.75i,
	}
	chk.VectorC(tst, "x", 1e-3, x, x_correct)
	checkResidC(tst, A.GetDenseMatrix(), x, b, 1e-12)
}
