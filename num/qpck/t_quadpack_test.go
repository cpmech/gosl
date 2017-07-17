// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qpck

import (
	. "math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestAgs01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Ags01a.")

	y := func(x float64) (res float64) {
		return Sqrt(1.0 + Pow(Sin(x), 3.0))
	}

	var fid int32
	A, abserr, neval, last, err := Agse(fid, y, 0, 1, 0, 0, nil, nil, nil, nil, nil)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	io.Pforan("A      = %v\n", A)
	io.Pforan("abserr = %v\n", abserr)
	io.Pforan("neval  = %v\n", neval)
	io.Pforan("last   = %v\n", last)
	chk.Float64(tst, "A", 1e-12, A, 1.08268158558)
}

func TestAgs01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Ags01b. goroutines")

	y := func(x float64) (res float64) {
		return Sqrt(1.0 + Pow(Sin(x), 3.0))
	}

	// channels
	nch := 3
	done := make(chan int, nch)

	// run all
	for ich := 0; ich < nch; ich++ {
		go func(fid int) {
			A, _, _, _, _ := Agse(int32(fid), y, 0, 1, 0, 0, nil, nil, nil, nil, nil)
			chk.Float64(tst, "A", 1e-12, A, 1.08268158558)
			done <- 1
		}(ich)
	}

	// wait
	for i := 0; i < nch; i++ {
		<-done
	}
}

// auxiliary function to run test

func runQ(tst *testing.T, name string, y fType, a, b, correct, tol float64) {
	res, _, _, _, err := Agse(0, y, a, b, 0, 0, nil, nil, nil, nil, nil)
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.AnaNum(tst, name, tol, res, correct, chk.Verbose)
}

// auxiliary function to run test (unbounded cases)
func runU(tst *testing.T, name string, y fType, bound float64, infCode int32, correct, tol float64) {
	res, _, _, _, err := Agie(0, y, bound, infCode, 0, 0, nil, nil, nil, nil, nil)
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.AnaNum(tst, name, tol, res, correct, chk.Verbose)
}

// auxiliary function to run test (with points cases)
func runP(tst *testing.T, name string, y fType, a, b, correct, tol float64, ptsAndBuf2 []float64) {
	res, _, _, _, err := Agpe(0, y, a, b, ptsAndBuf2, 0, 0, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.AnaNum(tst, name, tol, res, correct, chk.Verbose)
}

// auxiliary function to run test (with oscillatory factors)
func runO(tst *testing.T, name string, y fType, a, b, omega, correct, tol float64, isSin bool) {
	var integr int32 = 1 // cos(omega*x)
	if isSin {
		integr = 2 // sin(omega*x)
	}
	res, _, _, _, err := Awoe(0, y, a, b, omega, integr, 0, 0, 0, 0, nil, nil, nil, nil, nil, nil, 0, nil)
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.AnaNum(tst, name, tol, res, correct, chk.Verbose)
}

func TestAgs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Ags02. some functions")

	// 1. typical function with two extra arguments
	runQ(tst, "function # 1", func(x float64) float64 {
		// Bessel function integrand
		n, z := 2.0, 1.8
		return Cos(n*x-z*Sin(x)) / Pi
	}, 0, Pi, 0.30614353532540296487, 1e-16)

	// 2. infinite integration limits --- Euler's constant
	var infCode int32 = 1 // (bound,+infinity)
	runU(tst, "function # 2", func(x float64) float64 {
		// Euler's constant integrand
		if x == 0 {
			tst.Errorf("must compute f(0) = %v\n", -Exp(-x)*Log(x))
		}
		return -Exp(-x) * Log(x)
	}, 0, infCode, 5.772156649008392e-01, 1e-15) // exact: 0.577215664901532860606512, 1e-12)

	// 3. singular points in region of integration.
	Aref := 1.0 - Cos(2.5) + Exp(-2.5) - Exp(-5.0)
	runP(tst, "function # 3", func(x float64) float64 {
		if x > 0 && x < 2.5 {
			return Sin(x)
		} else if x >= 2.5 && x <= 5.0 {
			return Exp(-x)
		}
		return 0.0
	}, 0, 10, Aref, 1e-15, []float64{2.5, 5.0, 0, 0})

	// 4. sine weighted integral (finite limits)
	ome := Pow(2.0, 3.4)
	Aref = (20*Sin(ome) - ome*Cos(ome) + ome*Exp(-20)) / (Pow(20, 2) + Pow(ome, 2))
	runO(tst, "function # 4", func(x float64) float64 {
		a := 20.0
		return Exp(a * (x - 1))
	}, 0, 1, ome, Aref, 1e-16, true) // true => sin
}
