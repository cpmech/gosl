// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// rootSolTest runs root solution test
//  Note: xguess is the trial solution for Newton's method (not Brent's)
func rootSolTest(tst *testing.T, xa, xb, xguess, tolcmp float64, ffcnA fun.Ss, ffcnB fun.Vv, JfcnB fun.Mv, fname string, save, show bool) (xbrent float64) {

	// Brent
	io.Pfcyan("\n       - - - - - - - using Brent's method - - -- - - - \n")
	o := NewBrent(ffcnA)
	xbrent = o.Solve(xa, xb, false)
	var ybrent float64
	ybrent = ffcnA(xbrent)
	io.Pforan("x      = %v\n", xbrent)
	io.Pforan("f(x)   = %v\n", ybrent)
	io.Pforan("nfeval = %v\n", o.NFeval)
	io.Pforan("nit    = %v\n", o.It)
	if math.Abs(ybrent) > 1e-10 {
		tst.Errorf("Brent failed: f(x) = %g > 1e-10\n", ybrent)
		return
	}

	// Newton
	io.Pfcyan("\n       - - - - - - - using Newton's method - - -- - - - \n")
	var p NlSolver
	p.Init(1, ffcnB, nil, JfcnB, true, false, nil)
	xnewt := []float64{xguess}
	var cnd float64
	cnd = p.CheckJ(xnewt, 1e-6, true, !chk.Verbose)
	io.Pforan("cond(J) = %v\n", cnd)
	p.Solve(xnewt, false)
	var ynewt float64
	ynewt = ffcnA(xnewt[0])
	io.Pforan("x      = %v\n", xnewt[0])
	io.Pforan("f(x)   = %v\n", ynewt)
	io.Pforan("nfeval = %v\n", p.NFeval)
	io.Pforan("nJeval = %v\n", p.NJeval)
	io.Pforan("nit    = %v\n", p.It)
	if math.Abs(ynewt) > 1e-9 {
		tst.Errorf("Newton failed: f(x) = %g > 1e-10\n", ynewt)
		return
	}

	// compare Brent's and Newton's solutions
	chk.Float64(tst, "xbrent - xnewt", tolcmp, xbrent, xnewt[0])
	return
}

func TestBrent01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Brent01. root finding")

	ffcnA := func(x float64) (res float64) {
		res = math.Pow(x, 3.0) - 0.165*math.Pow(x, 2.0) + 3.993e-4
		return
	}

	ffcnB := func(fx, x la.Vector) {
		fx[0] = ffcnA(x[0])
		return
	}

	JfcnB := func(dfdx *la.Matrix, x la.Vector) {
		dfdx.Set(0, 0, 3.0*x[0]*x[0]-2.0*0.165*x[0])
		return
	}

	xa, xb := 0.0, 0.11
	//xguess := 0.001 // ===> this one converges to the right-hand solution
	xguess := 0.03
	//save   := true
	save := false
	rootSolTest(tst, xa, xb, xguess, 1e-7, ffcnA, ffcnB, JfcnB, "brent01.png", save, false)
}

func TestBrent02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Brent02. root finding")

	ffcnA := func(x float64) (res float64) {
		return x*x*x - 2.0*x - 5.0
	}

	ffcnB := func(fx, x la.Vector) {
		fx[0] = ffcnA(x[0])
		return
	}

	JfcnB := func(dfdx *la.Matrix, x la.Vector) {
		dfdx.Set(0, 0, 3.0*x[0]*x[0]-2.0)
		return
	}

	xa, xb := 2.0, 3.0
	xguess := 2.1
	//save   := true
	save := false
	xbrent := rootSolTest(tst, xa, xb, xguess, 1e-7, ffcnA, ffcnB, JfcnB, "brent02.png", save, false)
	chk.Float64(tst, "xsol", 1e-14, xbrent, 2.09455148154233)
}

func TestBrent03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Brent03. minimum finding")

	ffcn := func(x float64) (res float64) {
		return x*x*x - 2.0*x - 5.0
	}

	o := NewBrent(ffcn)
	xa, xb := 0.0, 1.0
	x := o.Min(xa, xb, false)
	y := ffcn(x)
	xcor := math.Sqrt(2.0 / 3.0)
	io.Pforan("x      = %v (correct=%g)\n", x, xcor)
	io.Pforan("f(x)   = %v\n", y)
	io.Pforan("nfeval = %v\n", o.NFeval)
	io.Pforan("nit    = %v\n", o.It)

	//save := true
	chk.Float64(tst, "xcorrect", 1e-8, x, xcor)
}
