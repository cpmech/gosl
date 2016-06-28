// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_hlsolver01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("hlsolver01. real")

	// sparse matrix
	var A Triplet
	A.Init(5, 5, 13)
	A.Put(0, 0, 1.0)
	A.Put(0, 0, 1.0)
	A.Put(1, 0, 3.0)
	A.Put(0, 1, 3.0)
	A.Put(2, 1, -1.0)
	A.Put(4, 1, 4.0)
	A.Put(1, 2, 4.0)
	A.Put(2, 2, -3.0)
	A.Put(3, 2, 1.0)
	A.Put(4, 2, 2.0)
	A.Put(2, 3, 2.0)
	A.Put(1, 4, 6.0)
	A.Put(4, 4, 1.0)

	// right-hand-side
	b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}

	// solve
	x, err := SolveRealLinSys(&A, b)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check
	x_correct := []float64{1, 2, 3, 4, 5}
	chk.Vector(tst, "x", 1e-14, x, x_correct)
	CheckResidR(tst, 1e-13, A.ToMatrix(nil).ToDense(), x, b)
}

func Test_hlsolver02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("hlsolver02. complex")

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
		77.38 + 8.82i,
		157.48 + 19.8i,
		1175.62 + 20.69i,
		912.12 - 801.75i,
		550 - 1060.4i,
	}

	// solve
	x, err := SolveComplexLinSys(&A, b)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check
	x_correct := []complex128{
		3.3 - 1i,
		1 + 0.17i,
		5.5,
		9,
		10 - 17.75i,
	}
	chk.VectorC(tst, "x", 1e-3, x, x_correct)
	CheckResidC(tst, 1e-12, A.ToMatrix(nil).ToDense(), x, b)
}
