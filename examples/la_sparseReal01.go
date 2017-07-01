// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func main() {

	// input matrix in Triplet format
	// including repeated positions. e.g. (0,0)
	A := new(la.Triplet)
	A.Init(5, 5, 13)  // 5 x 5 matrix with 13 non-zero entries
	A.Put(0, 0, +1.0) // 0  << repeated
	A.Put(0, 0, +1.0) // 1  << repeated
	A.Put(1, 0, +3.0) // 2
	A.Put(0, 1, +3.0) // 3
	A.Put(2, 1, -1.0) // 4
	A.Put(4, 1, +4.0) // 5
	A.Put(1, 2, +4.0) // 6
	A.Put(2, 2, -3.0) // 7
	A.Put(3, 2, +1.0) // 8
	A.Put(4, 2, +2.0) // 9
	A.Put(2, 3, +2.0) // 10
	A.Put(1, 4, +6.0) // 11
	A.Put(4, 4, +1.0) // 12

	// right-hand-side
	b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}

	// allocate solver
	o := la.NewSparseSolver("umfpack")
	defer o.Free()

	// initialise solver
	symmetric, verbose := false, false
	err := o.Init(A, symmetric, verbose, "", "", nil)
	if err != nil {
		io.Pf("Init failed:\n%v\n", err)
		return
	}

	// factorise
	err = o.Fact()
	if err != nil {
		io.Pf("Fact failed:\n%v\n", err)
		return
	}

	// solve
	x := la.NewVector(len(b))
	err = o.Solve(x, b, false) // x := inv(A) * b
	if err != nil {
		io.Pf("Solve failed:\n%v\n", err)
		return
	}

	// print solution
	xCorrect := []float64{1, 2, 3, 4, 5}
	io.Pf("x = %v\n\n", x)
	io.Pf("xCorrect = %v\n", xCorrect)
}
