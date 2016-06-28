// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

// SolveRealLinSys is a high-level function to generically solve a linear system with
// real variables:
//  A.x = b  =>  x = inv(A)*b
//  Input:
//    A -- [n][n] sparse matrix in triplet format
//    b -- [n] right-hand-side vector
//  Output:
//    x -- [n] solution vector
func SolveRealLinSys(A *Triplet, b []float64) (x []float64, err error) {

	// allocate solver
	lis := GetSolver("umfpack")
	defer lis.Clean()

	// info
	symmetric := false
	verbose := false
	timing := false

	// initialise solver (R)eal
	err = lis.InitR(A, symmetric, verbose, timing)
	if err != nil {
		return
	}

	// factorise
	err = lis.Fact()
	if err != nil {
		return
	}

	// solve (R)eal
	var dummy bool
	x = make([]float64, len(b))
	err = lis.SolveR(x, b, dummy) // x := inv(a) * b
	return
}

// SolveComplexLinSys is a high-level function to generically solve a linear system with
// complex variables:
//  A.x = b  =>  x = inv(A)*b
//  Input:
//    A -- [n][n] sparse matrix in triplet format
//    b -- [n] right-hand-side vector
//  Output:
//    x -- [n] solution vector
func SolveComplexLinSys(A *TripletC, b []complex128) (x []complex128, err error) {

	// allocate solver
	lis := GetSolver("umfpack")
	defer lis.Clean()

	// info
	symmetric := false
	verbose := false
	timing := false

	// initialise solver (C)omplex
	err = lis.InitC(A, symmetric, verbose, timing)
	if err != nil {
		return
	}

	// factorise
	err = lis.Fact()
	if err != nil {
		return
	}

	// auxiliary variables
	bR, bC := ComplexToRC(b)      // real and complex components of b
	xR := make([]float64, len(b)) // real compoments of x
	xC := make([]float64, len(b)) // complex compoments of x

	// solve (C)omplex
	var dummy bool
	err = lis.SolveC(xR, xC, bR, bC, dummy) // x := inv(A) * b
	if err != nil {
		return
	}

	// join solution vector
	x = RCtoComplex(xR, xC)
	return
}
