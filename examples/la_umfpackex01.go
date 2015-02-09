// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"

	"code.google.com/p/gosl/la"
)

func main() {

	// input matrix in Triplet format
	// including repeated positions. e.g. (0,0)
	var t la.Triplet
	t.Init(5, 5, 13)
	t.Put(0, 0, 1.0) // << repeated
	t.Put(0, 0, 1.0) // << repeated
	t.Put(1, 0, 3.0)
	t.Put(0, 1, 3.0)
	t.Put(2, 1, -1.0)
	t.Put(4, 1, 4.0)
	t.Put(1, 2, 4.0)
	t.Put(2, 2, -3.0)
	t.Put(3, 2, 1.0)
	t.Put(4, 2, 2.0)
	t.Put(2, 3, 2.0)
	t.Put(1, 4, 6.0)
	t.Put(4, 4, 1.0)

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
	err := lis.InitR(&t, symmetric, verbose, timing)
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

	// solve (R)eal
	var dummy bool
	x := make([]float64, len(b))
	err = lis.SolveR(x, b, dummy) // x := inv(a) * b
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	// output
	la.PrintMat("a", t.ToMatrix(nil).ToDense(), "%5g", false)
	la.PrintVec("b", b, "%v ", false)
	la.PrintVec("x", x, "%v ", false)
}
