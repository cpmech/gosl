// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
    "math"
    "os"
    "testing"
    "code.google.com/p/gosl/la"
    "code.google.com/p/gosl/plt"
    "code.google.com/p/gosl/utl"
)

const (
    DBL_EPSILON = 1.0e-15
    DBL_MIN     = math.SmallestNonzeroFloat64
)

const (
    NaN = iota
    Inf
    Equal
    NotEqual
)

// PlotYxe plots the function y(x) implemented by Cb_yxe
func PlotYxe(ffcn Cb_yxe, dirout, fname string, xsol, xa, xb float64, np int, xsolLbl, args string, save, show bool, extra func()) (err error) {
    if !save && !show {
        return
    }
    x := utl.LinSpace(xa, xb, np)
    y := make([]float64, np)
    for i := 0; i < np; i++ {
        y[i], err = ffcn(x[i])
        if err != nil {
            return
        }
    }
    var ysol float64
    ysol, err = ffcn(xsol)
    if err != nil {
        return
    }
    plt.Cross()
    plt.Plot(x, y, args)
    plt.PlotOne(xsol, ysol, utl.Sf("'ro', label='%s'", xsolLbl))
    if extra != nil {
        extra()
    }
    plt.Gll("x", "y(x)", "")
    if save {
        os.MkdirAll(dirout, 0777)
        plt.Save(dirout + "/" + fname)
    }
    if show {
        plt.Show()
    }
    return
}

func MinComp(tol, expected float64) float64 {
    return min(tol, math.Abs(expected)) + DBL_EPSILON
}

func TestAbs(result, expected, absolute_error float64, test_description string) (status int) {
    switch {
    case math.IsNaN(result) || math.IsNaN(expected):
        status = NaN

    case math.IsInf(result, 0) || math.IsInf(expected, 0):
        status = Inf

    case (expected > 0 && expected < DBL_MIN) || (expected < 0 && expected > -DBL_MIN):
        status = NotEqual

    default:
        if math.Abs(result-expected) > absolute_error {
            status = NotEqual
        } else {
            status = Equal
        }
    }
    if test_description != "" {
        utl.Pf(test_description)
        switch status {
        case NaN:
            utl.Pf(" [1;31mNaN[0m\n  %v observed\n  %v expected.  diff = %v\n", result, expected, result-expected)
        case Inf:
            utl.Pf(" [1;31mInf[0m\n  %v observed\n  %v expected.  diff = %v\n", result, expected, result-expected)
        case Equal:
            utl.Pf(" [1;32mOk[0m\n  %v observed\n  %v expected.  diff = %v\n", result, expected, result-expected)
        case NotEqual:
            utl.Pf(" [1;31mError[0m\n  %v observed\n  %v expected.  diff = %v\n", result, expected, result-expected)
        }
    }
    return
}

func CompareJac(tst *testing.T, ffcn Cb_f, Jfcn Cb_J, x []float64, tol float64, distr bool) {
    n := len(x)
    // numerical
    fx := make([]float64, n)
    w  := make([]float64, n) // workspace
    ffcn(fx, x)
    var Jnum la.Triplet
    Jnum.Init(n, n, n*n)
    Jacobian(&Jnum, ffcn, x, fx, w, distr)
    jn := Jnum.ToMatrix(nil)
    // analytical
    var Jana la.Triplet
    Jana.Init(n, n, n*n)
    Jfcn(&Jana, x)
    ja := Jana.ToMatrix(nil)
    // compare
    //la.PrintMat(fmt.Sprintf("Jana(%d)",mpi.Rank()), ja.ToDense(), "%13.6f", false)
    //la.PrintMat(fmt.Sprintf("Jnum(%d)",mpi.Rank()), jn.ToDense(), "%13.6f", false)
    max_diff := la.MatMaxDiff(jn.ToDense(), ja.ToDense())
    if max_diff > tol {
        tst.Errorf("[1;31mmax_diff = %g[0m\n", max_diff)
    } else {
        utl.Pf("[1;32mmax_diff = %g[0m\n", max_diff)
    }
}

// auxiliary functions
func max(a, b float64) float64 { if a > b { return a }; return b }
func min(a, b float64) float64 { if a < b { return a }; return b }
