// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
    "math"
    "testing"
    "code.google.com/p/gosl/plt"
    "code.google.com/p/gosl/utl"
)

// run_rootsol_test runs root solution test
//  Note: xguess is the trial solution for Newton's method (not Brent's)
func run_rootsol_test(tst *testing.T, xa, xb, xguess, tolcmp float64, ffcnA Cb_yxe, ffcnB Cb_f, JfcnB Cb_Jd, fname string, save, show bool) (xbrent float64) {

    // Brent
    utl.Pfcyan("\n       - - - - - - - using Brent's method - - -- - - - \n")
    var o Brent
    o.Init(ffcnA)
    var err error
    xbrent, err = o.Solve(xa, xb, false)
    if err != nil {
        utl.Panic("%v", err)
    }
    var ybrent float64
    ybrent, err = ffcnA(xbrent)
    if err != nil {
        utl.Panic("%v", err)
    }
    utl.Pforan("x      = %v\n", xbrent)
    utl.Pforan("f(x)   = %v\n", ybrent)
    utl.Pforan("nfeval = %v\n", o.NFeval)
    utl.Pforan("nit    = %v\n", o.It)
    if math.Abs(ybrent) > 1e-10 {
        utl.Panic("Brent failed: f(x) = %g > 1e-10\n", ybrent)
    }

    // Newton
    utl.Pfcyan("\n       - - - - - - - using Newton's method - - -- - - - \n")
    var p NlSolver
    p.Init(1, ffcnB, nil, JfcnB, true, false, nil)
    xnewt := []float64{xguess}
    var cnd float64
    cnd, err = p.CheckJ(xnewt, 1e-6, true, false)
    utl.Pforan("cond(J) = %v\n", cnd)
    if err != nil {
        utl.Panic("%v", err.Error())
    }
    err = p.Solve(xnewt, false)
    if err != nil {
        utl.Panic("%v", err.Error())
    }
    var ynewt float64
    ynewt, err = ffcnA(xnewt[0])
    if err != nil {
        utl.Panic("%v", err)
    }
    utl.Pforan("x      = %v\n", xnewt[0])
    utl.Pforan("f(x)   = %v\n", ynewt)
    utl.Pforan("nfeval = %v\n", p.NFeval)
    utl.Pforan("nJeval = %v\n", p.NJeval)
    utl.Pforan("nit    = %v\n", p.It)
    if math.Abs(ynewt) > 1e-9 {
        utl.Panic("Newton failed: f(x) = %g > 1e-10\n", ynewt)
    }

    // compare Brent's and Newton's solutions
    PlotYxe(ffcnA, "results", fname, xbrent, xa, xb, 101, "Brent", "'b-'", save, show, func() {
        plt.PlotOne(xnewt[0], ynewt, "'g+', ms=15, label='Newton'")
    })
    utl.CheckScalar(tst, "xbrent - xnewt", tolcmp, xbrent, xnewt[0])
    return
}

func Test_brent01(tst *testing.T) {

    prevTs := utl.Tsilent
    defer func() {
        utl.Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //utl.Tsilent = false
    utl.TTitle("brent01. root finding")

    ffcnA := func(x float64) (res float64, err error) {
        res = math.Pow(x,3.0) - 0.165*math.Pow(x,2.0) + 3.993e-4
        return
    }

    ffcnB := func(fx, x []float64) (err error) {
        fx[0], err = ffcnA(x[0])
        return
    }

    JfcnB := func(dfdx [][]float64, x []float64) (err error) {
        dfdx[0][0] = 3.0 * x[0]*x[0] - 2.0 * 0.165 * x[0]
        return
    }

    xa, xb := 0.0, 0.11
    //xguess := 0.001 // ===> this one converges to the right-hand solution
    xguess := 0.03
    //save   := true
    save   := false
    run_rootsol_test(tst, xa, xb, xguess, 1e-7, ffcnA, ffcnB, JfcnB, "brent01.png", save, false)
}

func Test_brent02(tst *testing.T) {

    prevTs := utl.Tsilent
    defer func() {
        utl.Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //utl.Tsilent = false
    utl.TTitle("brent02. root finding")

    ffcnA := func(x float64) (res float64, err error) {
        return x*x*x - 2.0*x - 5.0, nil
    }

    ffcnB := func(fx, x []float64) (err error) {
        fx[0], err = ffcnA(x[0])
        return
    }

    JfcnB := func(dfdx [][]float64, x []float64) (err error) {
        dfdx[0][0] = 3.0 * x[0]*x[0] - 2.0
        return
    }

    xa, xb := 2.0, 3.0
    xguess := 2.1
    //save   := true
    save   := false
    xbrent := run_rootsol_test(tst, xa, xb, xguess, 1e-7, ffcnA, ffcnB, JfcnB, "brent02.png", save, false)
    utl.CheckScalar(tst, "xsol", 1e-14, xbrent, 2.09455148154233)
}

func Test_brent03(tst *testing.T) {

    prevTs := utl.Tsilent
    defer func() {
        utl.Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //utl.Tsilent = false
    utl.TTitle("brent03. minimum finding")

    ffcn := func(x float64) (res float64, err error) {
        return x*x*x - 2.0*x - 5.0, nil
    }

    var o Brent
    o.Init(ffcn)
    xa, xb := 0.0, 1.0
    x, err := o.Min(xa, xb, false)
    if err != nil {
        utl.Panic("%v", err)
    }
    y, err := ffcn(x)
    if err != nil {
        utl.Panic("%v", err)
    }
    xcor := math.Sqrt(2.0/3.0)
    utl.Pforan("x      = %v (correct=%g)\n", x, xcor)
    utl.Pforan("f(x)   = %v\n", y)
    utl.Pforan("nfeval = %v\n", o.NFeval)
    utl.Pforan("nit    = %v\n", o.It)

    //save := true
    save := false
    PlotYxe(ffcn, "results", "brent03.png", x, -1, 3, 101, "Brent", "'b-'", save, false, nil)
    utl.CheckScalar(tst, "xcorrect", 1e-8, x, xcor)
}
