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
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// rootSolTest runs root solution test
//  Note: xguess is the trial solution for Newton's method (not Brent's)
func rootSolTest(tst *testing.T, xa, xb, xguess, tolcmp float64, ffcnA fun.Ss, ffcnB fun.Vv, JfcnB fun.Mv) (xbrent float64) {

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
	//xguess := 0.001 // ===> this one fails (Newton)
	xguess := 0.03
	xsol := rootSolTest(tst, xa, xb, xguess, 1e-7, ffcnA, ffcnB, JfcnB)

	if chk.Verbose {
		xx := utl.LinSpace(-0.5, 0.5, 101)
		yy := utl.GetMapped(xx, func(x float64) float64 { return ffcnA(x) })
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Plot(xx, yy, &plt.A{C: "b", Ls: "-", L: "curve1", NoClip: true})
		plt.PlotOne(xsol, ffcnA(xsol), &plt.A{C: "r", M: "o", NoClip: true})
		plt.Gll("$x$", "$y$", nil)
		plt.Cross(0, 0, nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "brent01")
	}
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
	xbrent := rootSolTest(tst, xa, xb, xguess, 1e-7, ffcnA, ffcnB, JfcnB)
	chk.Float64(tst, "xsol", 1e-14, xbrent, 2.09455148154233)

	if chk.Verbose {
		xx := utl.LinSpace(1, 3, 101)
		yy := utl.GetMapped(xx, func(x float64) float64 { return ffcnA(x) })
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Plot(xx, yy, &plt.A{C: "b", Ls: "-", L: "curve1", NoClip: true})
		plt.PlotOne(xbrent, ffcnA(xbrent), &plt.A{C: "r", M: "o", NoClip: true})
		plt.Gll("$x$", "$y$", nil)
		plt.Cross(0, 0, nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "brent02")
	}
}

func TestBrent03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Brent03. minimum finding")

	ffcn := func(x float64) (res float64) {
		return x*x*x - 2.0*x - 5.0
	}

	o := NewBrent(ffcn)
	//xa, xb := -2.0, 2.0 // ===> MinWithDerivs fails
	//xa, xb := -1.5, 1.5 // ===> MinWithDerivs fails
	xa, xb := -1.4, 1.4
	x := o.Min(xa, xb)
	y := ffcn(x)
	xcor := math.Sqrt(2.0 / 3.0)
	io.Pforan("x      = %v (correct=%g)\n", x, xcor)
	io.Pforan("f(x)   = %v\n", y)
	io.Pforan("nit    = %v\n", o.It)
	io.Pforan("nfeval = %v\n", o.NFeval)
	chk.Float64(tst, "xcorrect", 1e-8, x, xcor)

	Jfcn := func(x float64) (dfdx float64) {
		return 3*x*x - 2
	}

	xd := o.MinWithDerivs(xa, xb, Jfcn)
	io.Pf("xd     = %v (correct=%g)\n", xd, xcor)
	io.Pf("f(xd)  = %v\n", ffcn(xd))
	io.Pf("nit    = %v\n", o.It)
	io.Pf("nfeval = %v\n", o.NFeval)
	io.Pf("nJeval = %v\n", o.NJeval)

	if chk.Verbose {
		xx := utl.LinSpace(-2, 2, 101)
		yy := utl.GetMapped(xx, func(x float64) float64 { return ffcn(x) })
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Plot(xx, yy, &plt.A{C: "b", Ls: "-", L: "curve1", NoClip: true})
		plt.PlotOne(x, y, &plt.A{C: "r", M: "o", NoClip: true})
		plt.Gll("$x$", "$y$", nil)
		plt.Cross(0, 0, nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "brent03")
	}
}
