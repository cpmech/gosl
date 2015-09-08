// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
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
	var A la.Triplet
	A.Init(5, 5, 13)
	A.Put(0, 0, 1.0) // << repeated
	A.Put(0, 0, 1.0) // << repeated
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

	// allocate solver
	lis := la.GetSolver("umfpack")
	defer lis.Clean()

	// info
	symmetric := false
	verbose := false
	timing := false

	// initialise solver (R)eal
	err := lis.InitR(&A, symmetric, verbose, timing)
	if err != nil {
		io.Pfred("solver failed:\n%v", err)
		return
	}

	// factorise
	err = lis.Fact()
	if err != nil {
		io.Pfred("solver failed:\n%v", err)
		return
	}

	// solve (R)eal
	var dummy bool
	x := make([]float64, len(b))
	err = lis.SolveR(x, b, dummy) // x := inv(a) * b
	if err != nil {
		io.Pfred("solver failed:\n%v", err)
		return
	}

	// output
	la.PrintMat("a", A.ToMatrix(nil).ToDense(), "%5g", false)
	la.PrintVec("b", b, "%v ", false)
	la.PrintVec("x", x, "%v ", false)
}
