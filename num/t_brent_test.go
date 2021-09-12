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
func rootSolTest(tst *testing.T, xa, xb, xguess, tolcmp float64, ffcnA fun.Ss, ffcnB fun.Vv, JfcnB fun.Mv) (xbrent float64) {

	// Brent
	io.Pfcyan("\n       - - - - - - - using Brent's method - - -- - - - \n")
	o := NewBrent(ffcnA, nil)
	xbrent = o.Root(xa, xb)
	ybrent := ffcnA(xbrent)
	io.Pforan("x      = %v\n", xbrent)
	io.Pforan("f(x)   = %v\n", ybrent)
	io.Pforan("nfeval = %v\n", o.NumFeval)
	io.Pforan("nit    = %v\n", o.NumIter)
	if math.Abs(ybrent) > 1e-10 {
		tst.Errorf("Brent failed: f(x) = %g > 1e-10\n", ybrent)
		return
	}

	// Newton
	io.Pfcyan("\n       - - - - - - - using Newton's method - - -- - - - \n")
	sol := NewNlSolver(1, ffcnB)
	sol.SetJacobianFunction(nil, JfcnB)
	xnewt := []float64{xguess}
	cnd := sol.CheckJ(xnewt, 1e-6, chk.Verbose)
	io.Pforan("cond(J) = %v\n", cnd)
	sol.Solve(xnewt)
	ynewt := ffcnA(xnewt[0])
	io.Pforan("x       = %v\n", xnewt[0])
	io.Pforan("f(x)    = %v\n", ynewt)
	io.Pforan("nfeval  = %v\n", sol.Nfeval)
	io.Pforan("nJeval  = %v\n", sol.Njeval)
	io.Pforan("nit     = %v\n", sol.Niter)
	if math.Abs(ynewt) > 1e-9 {
		tst.Errorf("Newton failed: f(x) = %g > 1e-10\n", ynewt)
		return
	}

	// compare Brent's and Newton's solutions
	chk.Float64(tst, "xbrent - xnewt", tolcmp, xbrent, xnewt[0])
	return
}

func TestBrent01(tst *testing.T) {

	// verbose()
	chk.PrintTitle("Brent01. root finding")

	ffcnA := func(x float64) (res float64) {
		res = math.Pow(x, 3.0) - 0.165*math.Pow(x, 2.0) + 3.993e-4
		return
	}

	ffcnB := func(fx, x la.Vector) {
		fx[0] = ffcnA(x[0])
	}

	JfcnB := func(dfdx *la.Matrix, x la.Vector) {
		dfdx.Set(0, 0, 3.0*x[0]*x[0]-2.0*0.165*x[0])
	}

	xa, xb := 0.0, 0.11
	//xguess := 0.001 // ===> this one fails (Newton)
	xguess := 0.03
	rootSolTest(tst, xa, xb, xguess, 1e-7, ffcnA, ffcnB, JfcnB)
}

func TestBrent02(tst *testing.T) {

	// verbose()
	chk.PrintTitle("Brent02. root finding")

	ffcnA := func(x float64) (res float64) {
		return x*x*x - 2.0*x - 5.0
	}

	ffcnB := func(fx, x la.Vector) {
		fx[0] = ffcnA(x[0])
	}

	JfcnB := func(dfdx *la.Matrix, x la.Vector) {
		dfdx.Set(0, 0, 3.0*x[0]*x[0]-2.0)
	}

	xa, xb := 2.0, 3.0
	xguess := 2.1
	xbrent := rootSolTest(tst, xa, xb, xguess, 1e-7, ffcnA, ffcnB, JfcnB)
	chk.Float64(tst, "xsol", 1e-11, xbrent, 2.09455148154233)
}

func TestBrent03(tst *testing.T) {

	// verbose()
	chk.PrintTitle("Brent03. minimum finding")

	ffcn := func(x float64) (res float64) {
		return x*x*x - 2.0*x - 5.0
	}
	Jfcn := func(x float64) (dfdx float64) {
		return 3*x*x - 2
	}

	o := NewBrent(ffcn, Jfcn)
	//xa, xb := -2.0, 2.0 // ===> fails
	//xa, xb := -1.5, 1.5 // ===> fails
	xa, xb := -1.4, 1.4
	x := o.Min(xa, xb)
	y := ffcn(x)
	xcor := math.Sqrt(2.0 / 3.0)
	io.Pforan("x      = %v (correct=%g)\n", x, xcor)
	io.Pforan("f(x)   = %v\n", y)
	io.Pforan("nit    = %v\n", o.NumIter)
	io.Pforan("nfeval = %v\n", o.NumFeval)
	chk.Float64(tst, "xcorrect", 1.18e-8, x, xcor)

	xd := o.MinUseD(xa, xb)
	io.Pf("xd     = %v (correct=%g)\n", xd, xcor)
	io.Pf("f(xd)  = %v\n", ffcn(xd))
	io.Pf("nit    = %v\n", o.NumIter)
	io.Pf("nfeval = %v\n", o.NumFeval)
	io.Pf("nJeval = %v\n", o.NumJeval)
}
