// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func main() {

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

	// flag indicating to store (real,complex) values in monolithic form => 1D array
	xzmono := false

	// input matrix in Complex Triplet format
	var t la.TripletC
	t.Init(5, 5, 16, xzmono) // 5 x 5 matrix with 16 non-zeros

	// first column
	t.Put(0, 0, 19.73, 0) // i=0, j=0, real=19.73, complex=0
	t.Put(1, 0, 0, -0.51) // i=1, j=0, real=0, complex=-0.51

	// second column
	t.Put(0, 1, 12.11, -1) // i=0, j=1, real=12.11, complex=-1
	t.Put(1, 1, 32.3, 7)
	t.Put(2, 1, 0, -0.51)

	// third column
	t.Put(0, 2, 0, 5)
	t.Put(1, 2, 23.07, 0)
	t.Put(2, 2, 70, 7.3)
	t.Put(3, 2, 1, 1.1)

	// fourth column
	t.Put(1, 3, 0, 1)
	t.Put(2, 3, 3.95, 0)
	t.Put(3, 3, 50.17, 0)
	t.Put(4, 3, 0, -9.351)

	// fifth column
	t.Put(2, 4, 19, 31.83)
	t.Put(3, 4, 45.51, 0)
	t.Put(4, 4, 55, 0)

	// right-hand-side
	b := []complex128{
		77.38 + 8.82i,
		157.48 + 19.8i,
		1175.62 + 20.69i,
		912.12 - 801.75i,
		550 - 1060.4i,
	}

	// allocate solver
	lis := la.GetSolver("umfpack")
	defer lis.Clean()

	// info
	symmetric := false
	verbose := false
	timing := false

	// initialise solver (C)omplex
	err := lis.InitC(&t, symmetric, verbose, timing)
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	// factorise
	err = lis.Fact()
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	// auxiliary variables
	bR, bC := la.ComplexToRC(b)   // real and complex components of b
	xR := make([]float64, len(b)) // real compoments of x
	xC := make([]float64, len(b)) // complex compoments of x

	// solve (C)omplex
	var dummy bool
	err = lis.SolveC(xR, xC, bR, bC, dummy) // x := inv(A) * b
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	// join solution vector
	x := la.RCtoComplex(xR, xC)

	// output
	A := t.ToMatrix(nil)
	io.Pforan("A.x = b\n")
	la.PrintMatC("A", A.ToDense(), "(%5g", "%+6gi) ", false)
	la.PrintVecC("b", b, "(%g", "%+gi) ", false)
	la.PrintVecC("x", x, "(%.3f", "%+.3fi) ", false)
}
