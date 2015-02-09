// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gm (geometry and meshes) implements auxiliary functions for
// handling geometry and mesh structures
package gm

import (
    "math"
    "strings"
    "code.google.com/p/gosl/la"
    "code.google.com/p/gosl/num"
    "code.google.com/p/gosl/plt"
    "code.google.com/p/gosl/utl"
)

const (
    ZTOL = 1e-14
    STOL = 1e-14
)

// Bspline holds B-spline data
type Bspline struct {
    // essential
    T        []float64   // array of knots: e.g: T = [t_0, t_1, ..., t_{m-1}]
    p        int         // order/degreee (has to call SetOrder to change this)
    m        int         // number of knots == len(T)
    tmin     float64     // minimum T
    tmax     float64     // maximum T
    // optional
    Q        [][]float64 // control points (has to call SetControl to change this)
    okQ      bool        // flag telling that Q was properly set
    // auxiliary
    span     int         // the span computed by CalcBasis or CalcBasisAndDerivs to be used in GetBasis or GetDeriv
    le       []float64   // left
    ri       []float64   // right
    ndu      [][]float64 // basis functions and knots differences
    der      [][]float64 // derivatives
    daux     [][]float64 // auxiliary array for computing derivatives
}

// Init initialises B-spline
func (o *Bspline) Init(T []float64, p int) {

    // check
    if len(T) < 2*(p+1) {
        utl.Panic(_bspline_err00, 2*(p+1), p, len(T))
    }

    // essential
    o.T, o.p, o.m  = T, p, len(T)
    o.tmin, o.tmax = la.VecMinMax(T)

    // auxiliary
    o.le   = make([]float64, o.p+1)
    o.ri   = make([]float64, o.p+1)
    o.ndu  = la.MatAlloc(o.p+1, o.p+1)
    o.der  = la.MatAlloc(o.p+1, o.p+1)
    o.daux = la.MatAlloc(2, o.p+1)
}

// SetOrder sets B-spline order (p)
func (o *Bspline) SetOrder(p int) {
    o.p = p
}

// NumBasis returns the number (n) of basis functions == number of control points
func (o *Bspline) NumBasis() int {
    return o.m - o.p - 1
}

// SetControl sets B-spline control points
func (o *Bspline) SetControl(Q [][]float64) {
    if len(Q) != o.NumBasis() {
        utl.Panic(_bspline_err01, o.p, o.NumBasis())
    }
    o.Q, o.okQ = Q, true
}

// CalcBasis computes all non-zero basis functions N[i] @ t
// Note: use GetBasis to get a particular basis function value
func (o *Bspline) CalcBasis(t float64) {
    // check
    if t < o.tmin || t > o.tmax {
        utl.Panic(_bspline_err02, "CalcBasis", t, o.tmin, o.tmax)
    }
    // using basis_funs (Piegl & Tiller, algorithm A2.2)
    o.span = o.find_span(t)
    o.basis_funs(t, o.span)
}

// CalcBasisAndDerivs computes all non-zero basis functions N[i] and corresponding
// first order derivatives of basis functions w.r.t t => dR[i]dt @ t
// Note: use GetBasis to get a particular basis function value
//       use GetDeriv to get a particular derivative
func (o *Bspline) CalcBasisAndDerivs(t float64) {
    // check
    if t < o.tmin || t > o.tmax {
        utl.Panic(_bspline_err02, "CalcBasisAndDerivs", t, o.tmin, o.tmax)
    }
    // using ders_basis_funs (Piegl & Tiller, algorithm A2.3)
    o.span = o.find_span(t)
    o.ders_basis_funs(t, o.span, 1)
}

// GetBasis returns the basis function N[i] just computed by CalcBasis or CalcBasisAndDerivs
func (o *Bspline) GetBasis(i int) float64 {
    j := i + o.p - o.span
    if j >= 0 && j <= o.p {
        return o.ndu[j][o.p]
    }
    return 0
}

// GetDeriv returns the derivative dN[i]dt just computed by CalcBasisAndDerivs
func (o *Bspline) GetDeriv(i int) float64 {
    j := i + o.p - o.span
    if j >= 0 && j <= o.p {
        return o.der[1][j]
    }
    return 0
}

// RecursiveBasis computes one particular basis function N[i] recursively (not efficient)
func (o *Bspline) RecursiveBasis(t float64, i int) float64 {
    // check
    if t < o.tmin || t > o.tmax {
        utl.Panic(_bspline_err02, "RecursiveBasis", t, o.tmin, o.tmax)
    }
    // using Cox-DeBoor formula
    return o.recursiveN(t, i, o.p)
}

// NumericalDeriv computes a particular derivative dN[i]dt @ t using numerical differentiation
// Note: it uses RecursiveBasis and therefore is highly non-efficient
func (o *Bspline) NumericalDeriv(t float64, i int) float64 {
    // check
    if t < o.tmin || t > o.tmax {
        utl.Panic(_bspline_err02, "NumericalDeriv", t, o.tmin, o.tmax)
    }
    // derivatives
    f := func(x float64, args ...interface{}) float64 {
        return o.RecursiveBasis(x, i)
    }
    return num.DerivRange(f, t, o.tmin, o.tmax)
}

// Point returns the x-y-z coordinates of a point on B-spline
// option =  0 : use CalcBasis
//           1 : use RecursiveBasis
func (o *Bspline) Point(t float64, option int) (C []float64) {
    // check
    if !o.okQ {
        utl.Panic(_bspline_err03, "Point")
    }
    // compute point on curve
    ncp := len(o.Q[0]) // number of components in Q
    C    = make([]float64, ncp)
    switch option {
    case 0: // recursive
        for i, q := range o.Q {
            for j := 0; j < ncp; j++ {
                C[j] += o.RecursiveBasis(t, i) * q[j]
            }
        }
    case 1: // Piegl & Tiller: A3.1 p82
        span := o.find_span(t)
        o.basis_funs(t, span)
        for i := 0; i <= o.p; i++ {
            for j := 0; j < ncp; j++ {
                C[j] += o.ndu[i][o.p] * o.Q[span-o.p+i][j]
            }
        }
    }
    return
}

// Elements returns the indices of nonzero spans
func (o *Bspline) Elements() (spans [][]int) {
    nspans := 0
    for i := 0; i < o.m-1; i++ {
        l := o.T[i+1] - o.T[i]
        if math.Abs(l) > STOL {
            nspans += 1
        }
    }
    spans = utl.IntsAlloc(nspans, 2)
    ispan := 0
    for i := 0; i < o.m-1; i++ {
        l := o.T[i+1] - o.T[i]
        if math.Abs(l) > STOL {
            spans[ispan][0] = i
            spans[ispan][1] = i+1
            ispan += 1
        }
    }
    return
}

// Draw draws curve and control points
// option =  0 : use CalcBasis
//           1 : use RecursiveBasis
func (o *Bspline) Draw2D(curveArgs, ctrlArgs string, npts, option int) {
    if !o.okQ {
        utl.Panic(_bspline_err03, "Draw")
    }
    tt := utl.LinSpace(o.tmin, o.tmax, npts)
    xx := make([]float64, npts)
    yy := make([]float64, npts)
    for i, t := range tt {
        C := o.Point(t, option)
        xx[i], yy[i] = C[0], C[1]
    }
    qx := make([]float64, o.NumBasis())
    qy := make([]float64, o.NumBasis())
    for i := 0; i < o.NumBasis(); i++ {
        qx[i], qy[i] = o.Q[i][0], o.Q[i][1]
    }
    lbls := []string{"Nonly", "recN"}
    plt.Plot(xx, yy, utl.Sf("'k-', clip_on=0, label=r'%s'",lbls[option]) + curveArgs)
    plt.Plot(qx, qy, "'r-', clip_on=0, label=r'ctrl', marker='.'" + ctrlArgs)
    plt.Gll("$x$", "$y$", "leg=1, leg_out=1, leg_ncol=2, leg_hlen=1.5, leg_fsz=7")
}

// PlotBasis plots basis functions in I
// option =  0 : use CalcBasis
//           1 : use CalcBasisAndDerivs
//           2 : use RecursiveBasis
func (o *Bspline) PlotBasis(args string, npts, option int) {
    nmks := 10
    tt   := utl.LinSpace(o.tmin, o.tmax, npts)
    I    := utl.IntRange(o.NumBasis())
    f    := make([]float64, len(tt))
    lbls := []string{"Nonly", "N\\&dN", "recN"}
    var cmd string
    for _, i := range I {
        for j, t := range tt {
            switch option {
            case 0:
                o.CalcBasis(t)
                f[j] = o.GetBasis(i)
            case 1:
                o.CalcBasisAndDerivs(t)
                f[j] = o.GetBasis(i)
            case 2:
                f[j] = o.RecursiveBasis(t, i)
            }
        }
        if strings.Contains(args, "marker") {
            cmd = utl.Sf("label=r'%s:%d', color=GetClr(%d, 2) %s", lbls[option], i, i, args)
        } else {
            cmd = utl.Sf("label=r'%s:%d', marker=(None if %d %%2 == 0 else GetMrk(%d/2,1)), markevery=(%d-1)/%d, clip_on=0, color=GetClr(%d, 2) %s", lbls[option], i, i, i, npts, nmks, i, args)
        }
        plt.Plot(tt, f, cmd)
    }
    plt.Gll("$t$", utl.Sf("$N_{i,%d}$", o.p), utl.Sf("leg=1, leg_out=1, leg_ncol=%d, leg_hlen=1.5, leg_fsz=7", o.NumBasis()))
    o.plt_ticks_spans()
}

// PlotDerivs plots derivatives of basis functions in I
// option =  0 : use CalcBasisAndDerivs
//           1 : use NumericalDeriv
func (o *Bspline) PlotDerivs(args string, npts, option int) {
    nmks := 10
    tt   := utl.LinSpace(o.tmin, o.tmax, npts)
    I    := utl.IntRange(o.NumBasis())
    f    := make([]float64, len(tt))
    lbls := []string{"N\\&dN", "numD"}
    var cmd string
    for _, i := range I {
        for j, t := range tt {
            switch option {
            case 0:
                o.CalcBasisAndDerivs(t)
                f[j] = o.GetDeriv(i)
            case 1:
                f[j] = o.NumericalDeriv(t, i)
            }
        }
        if strings.Contains(args, "marker") {
            cmd = utl.Sf("label=r'%s:%d', color=GetClr(%d, 2) %s", lbls[option], i, i, args)
        } else {
            cmd = utl.Sf("label=r'%s:%d', marker=(None if %d %%2 == 0 else GetMrk(%d/2,1)), markevery=(%d-1)/%d, clip_on=0, color=GetClr(%d, 2) %s", lbls[option], i, i, i, npts, nmks, i, args)
        }
        plt.Plot(tt, f, cmd)
    }
    plt.Gll("$t$", utl.Sf(`$\frac{\mathrm{d}N_{i,%d}}{\mathrm{d}t}$`, o.p), utl.Sf("leg=1, leg_out=1, leg_ncol=%d, leg_hlen=1.5, leg_fsz=7", o.NumBasis()))
    o.plt_ticks_spans()
}

// auxiliary methods /////////////////////////////////////////////////////////////////////////////////

// find_span returns the span where t falls in
func (o *Bspline) find_span(t float64) int {
    // Piegl & Tiller: A2.1 p68
    n := o.NumBasis()
    if t >= o.T[n] {
        return n-1
    }
    if t <= o.T[o.p] {
        return o.p
    }
    low, high, mid := o.p, n, (o.p + n)/2
    for t < o.T[mid] || t >= o.T[mid+1] {
        if t < o.T[mid] {
            high = mid
        } else {
            low = mid
        }
        mid = (low + high)/2
    }
    return mid
}

// recursiveN computes basis functions using Cox-DeBoors recursive formula
func (o *Bspline) recursiveN(t float64, i int, p int) float64 {
    if math.Abs(t-o.tmax) < ZTOL {
        t = o.tmax - ZTOL // clean rubbish. e.g. 1.000000000000002
    }
    if p == 0 {
        if t < o.T[i]   { return 0.0 }
        if t < o.T[i+1] { return 1.0 }
        return 0.0
    } else {
        d1 := o.T[i+p]   - o.T[i]
        d2 := o.T[i+p+1] - o.T[i+1]
        var N1, N2 float64
        if math.Abs(d1) < ZTOL {
            N1, d1 = 0.0, 1.0
        } else {
            N1 = o.recursiveN(t, i, p-1)
        }
        if math.Abs(d2) < ZTOL {
            N2, d2 = 0.0, 1.0
        } else {
            N2 = o.recursiveN(t, i+1, p-1)
        }
        return (t - o.T[i]) * N1 / d1 + (o.T[i+p+1] - t) * N2 / d2
    }
}

// basis_funs computes basis functions using Piegl-Tiller algorithm A2.2/2.3 p70/p72
func (o *Bspline) basis_funs(t float64, span int) {
    // Piegl & Tiller: A2.3 p72
    // compute basis functions and knot differences
    var temp, saved float64
    o.ndu[0][0] = 1
    for j := 1; j <= o.p; j++ {
        o.le[j] = t - o.T[span+1-j]
        o.ri[j] = o.T[span+j] - t
        saved   = 0
        for r := 0; r < j; r++ {
            o.ndu[j][r] = o.ri[r+1] + o.le[j-r]
            temp        = o.ndu[r][j-1] / o.ndu[j][r]
            o.ndu[r][j] = saved + o.ri[r+1] * temp
            saved       = o.le[j-r] * temp
        }
        o.ndu[j][j] = saved
    }
}

// ders_basis_funs computes derivatives of basis functions using Piegl-Tiller algorithm A2.3 p72
func (o *Bspline) ders_basis_funs(t float64, span, upto int) {
    // compute and load the basis functions
    o.basis_funs(t, span)
    for j := 0; j <= o.p; j++ {
        o.der[0][j] = o.ndu[j][o.p]
    }
    // compute the derivatives (Eq 2.9)
    var d float64
    var s1, s2, rk, pk, j1, j2 int
    for r := 0; r <= o.p; r++ {
        s1, s2 = 0, 1 // alternate rows in array
        o.daux[0][0] = 1
        // loop to compute k-th derivative
        for k := 1; k <= upto; k++ {
            d, rk, pk = 0, r-k, o.p-k
            if r >= k {
                o.daux[s2][0] = o.daux[s1][0] / o.ndu[pk+1][rk]
                d             = o.daux[s2][0] * o.ndu[rk][pk]
            }
            if rk >= -1 {
                j1 = 1
            } else {
                j1 = -rk
            }
            if r-1 <= pk {
                j2 = k-1
            } else {
                j2 = o.p-r
            }
            for j := j1; j <= j2; j++ {
                o.daux[s2][j] = (o.daux[s1][j] - o.daux[s1][j-1]) / o.ndu[pk+1][rk+j]
                d            += o.daux[s2][j] * o.ndu[rk+j][pk]
            }
            if r <= pk {
                o.daux[s2][k] = -o.daux[s1][k-1] / o.ndu[pk+1][r]
                d            += o.daux[s2][k] * o.ndu[r][pk]
            }
            o.der[k][r] = d
            s1, s2 = s2, s1
        }
    }
    // multiply through by the correct factors
    d = float64(o.p)
    for k := 1; k <= upto; k++ {
        for j := 0; j <= o.p; j++ {
            o.der[k][j] *= d
        }
        d *= float64(o.p-k)
    }
}

// plt_ticks_spans adds ticks indicating spans
func (o *Bspline) plt_ticks_spans() {
    lbls := make(map[float64]string, 0)
    for i, t := range o.T {
        if _, ok := lbls[t]; !ok {
            lbls[t] = utl.Sf("'[%d", i)
        } else {
            lbls[t] += utl.Sf(",%d", i)
        }
    }
    for t, l := range lbls {
        plt.AnnotateXlabels(t, utl.Sf("%s]'", l), "")
    }
}

// error messages
var (
    _bspline_err00 = "bspline.go: Init: at least %d knots are required to define clamped B-spline of order p==%d. m==%d is invalid"
    _bspline_err01 = "bspline.go: Set_Q: B-spline of order %d needs %d control points"
    _bspline_err02 = "bspline.go: %s: t must be within [%g, %g]. t=%g is incorrect"
    _bspline_err03 = "bspline.go: %s: Q must be set before calling this method"
)
