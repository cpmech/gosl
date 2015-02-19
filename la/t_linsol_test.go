// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func run_linsol_testR(tst *testing.T, t *Triplet, tol_cmp, tol_res float64, b, x_correct []float64) {

	// info
	symmetric := false
	verbose := false
	timing := false

	// allocate solver
	lis := GetSolver("umfpack")
	defer lis.Clean()

	// initialise solver
	err := lis.InitR(t, symmetric, verbose, timing)
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// factorise
	err = lis.Fact()
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// solve
	var dummy bool
	x := make([]float64, len(b))
	err = lis.SolveR(x, b, dummy) // x := inv(A) * b
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// output
	A := t.ToMatrix(nil)
	io.Pforan("A.x = b\n")
	PrintMat("A", A.ToDense(), "%5g", false)
	PrintVec("x", x, "%g ", false)
	PrintVec("b", b, "%g ", false)

	// check
	chk.Vector(tst, "x", tol_cmp, x, x_correct)
	check_residR(tst, tol_res, A.ToDense(), x, b)
}

func run_linsol_testC(tst *testing.T, t *TripletC, tol_cmp, tol_res float64, b, x_correct []complex128) {

	// info
	symmetric := false
	verbose := false
	timing := false

	// allocate solver
	lis := GetSolver("umfpack")
	defer lis.Clean()

	// initialise solver
	err := lis.InitC(t, symmetric, verbose, timing)
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// factorise
	err = lis.Fact()
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// solve
	var dummy bool
	bR, bC := ComplexToRC(b)
	xR := make([]float64, len(b))
	xC := make([]float64, len(b))
	err = lis.SolveC(xR, xC, bR, bC, dummy) // x := inv(A) * b
	if err != nil {
		chk.Panic("%v", err.Error())
	}
	x := RCtoComplex(xR, xC)

	// output
	A := t.ToMatrix(nil)
	io.Pforan("A.x = b\n")
	PrintMatC("A", A.ToDense(), "(%g+", "%gi) ", false)
	PrintVecC("x", x, "(%g+", "%gi) ", false)
	PrintVecC("b", b, "(%g+", "%gi) ", false)

	// check
	chk.VectorC(tst, "x", tol_cmp, x, x_correct)
	check_residC(tst, tol_res, A.ToDense(), x, b)
}

func Test_linsol01(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("linsol01. real")

	// input matrix data into Triplet
	var t Triplet
	t.Init(5, 5, 13)
	t.Put(0, 0, 1.0)
	t.Put(0, 0, 1.0)
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

	// run test
	b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}
	x_correct := []float64{1, 2, 3, 4, 5}
	run_linsol_testR(tst, &t, 1e-14, 1e-13, b, x_correct)
}

func Test_linsol02(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("linsol02. real")

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
	x_correct := []float64{-1, 8, -65, 454, -2725, 13624, -54497, 163490, -326981, 326991}
	run_linsol_testR(tst, &t, 1e-4, 1e-10, b, x_correct)
}

func Test_linsol03(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("linsol03. complex (but real)")

	// input matrix data into Triplet
	var t TripletC
	t.Init(5, 5, 13, false)
	t.Put(0, 0, 1.0, 0)
	t.Put(0, 0, 1.0, 0)
	t.Put(1, 0, 3.0, 0)
	t.Put(0, 1, 3.0, 0)
	t.Put(2, 1, -1.0, 0)
	t.Put(4, 1, 4.0, 0)
	t.Put(1, 2, 4.0, 0)
	t.Put(2, 2, -3.0, 0)
	t.Put(3, 2, 1.0, 0)
	t.Put(4, 2, 2.0, 0)
	t.Put(2, 3, 2.0, 0)
	t.Put(1, 4, 6.0, 0)
	t.Put(4, 4, 1.0, 0)

	// run test
	b := []complex128{8.0, 45.0, -3.0, 3.0, 19.0}
	x_correct := []complex128{1, 2, 3, 4, 5}
	run_linsol_testC(tst, &t, 1e-14, 1e-13, b, x_correct)
}

func Test_linsol04(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("linsol04. complex (but real)")

	// input matrix data into Triplet
	var t TripletC
	t.Init(10, 10, 64, false)
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
			t.Put(i, j, val, 0)
		}
	}

	// run test
	b := []complex128{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	x_correct := []complex128{-1, 8, -65, 454, -2725, 13624, -54497, 163490, -326981, 326991}
	run_linsol_testC(tst, &t, 1e-4, 1e-9, b, x_correct)
}

func Test_linsol05(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("linsol05. complex (but real)")

	// data
	n := 10
	b := make([]complex128, n)
	x_correct := make([]complex128, n)

	// input matrix data into Triplet
	var t TripletC
	t.Init(n, n, n, false)
	for i := 0; i < n; i++ {

		// Some very fake diagonals. Should take exactly 20 GMRES steps
		ar := 10.0 + float64(i)/(float64(n)/10.0)
		ac := 10.0 - float64(i)/(float64(n)/10.0)
		t.Put(i, i, ar, ac)

		// Let exact solution = 1 + 0.5i
		x_correct[i] = complex(float64(i+1), float64(i+1)/10.0)

		// Generate RHS to match exact solution
		b[i] = complex(ar*real(x_correct[i])-ac*imag(x_correct[i]),
			ar*imag(x_correct[i])+ac*real(x_correct[i]))
	}

	// run text
	run_linsol_testC(tst, &t, 1e-14, 1e-13, b, x_correct)
}
