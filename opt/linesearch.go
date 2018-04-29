// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// LineSearch finds the scalar 'a' that gives a substantial reduction of f({x}+a⋅{u})
//
//   REFERENCES:
//   [1] Nocedal, J. and Wright, S. (2006) Numerical Optimization.
//       Springer Series in Operations Research. 2nd Edition. Springer. 664p
//
type LineSearch struct {

	// configuration
	MaxIt     int     // max iterations
	MaxItZoom int     // max iterations in zoom routine
	MaxAlpha  float64 // max 'a'
	MulAlpha  float64 // multiplier to increase 'a'; e.g. 1.5
	Coef1     float64 // "sufficient decrease" coefficient (c1); typically = 1e-4 (Fig 3.3, page 33 of [1])
	Coef2     float64 // "curvature condition" coefficient (c2); typically = 0.1 for CG methods and 0.9 for Newton or quasi-Newton methods (Fig 3.4, page 34 of [1])
	CoefQuad  float64 // coefficient for limiting 'a' in quadratic interpolation in zoom
	CoefCubic float64 // coefficient for limiting 'a' in cubic interpolation in zoom

	// statistics and History (for debugging)
	NumFeval    int // number of calls to Ffcn (function evaluations)
	NumJeval    int // number of calls to Jfcn (Jacobian evaluations)
	NumIter     int // number of iterations from last call to Find
	NumIterZoom int // number of iterations from last call to zoom

	// internal
	ffcn    fun.Sv           // scalar function of vector: y = f({x})
	Jfcn    fun.Vv           // vector function of vector: {J} = df/d{x} @ {x}
	xnew    la.Vector        // {xnew} = {x} + a⋅{p}
	dfdx    la.Vector        // derivative df/d{x}
	interp2 *fun.InterpQuad  // quadratic intepolator
	interp3 *fun.InterpCubic // cubic intepolator

	// pointers
	x la.Vector // starting point
	u la.Vector // direction
}

// NewLineSearch returns a new LineSearch object
//   ndim -- length(x)
//   ffcn -- function y = f({x})
//   Jfcn -- Jacobian {J} = df/d{x} @ {x}
func NewLineSearch(ndim int, ffcn fun.Sv, Jfcn fun.Vv) (o *LineSearch) {
	o = new(LineSearch)
	o.MaxIt = 10
	o.MaxItZoom = 10
	o.MaxAlpha = 100
	o.MulAlpha = 2
	o.Coef1 = 1e-4
	o.Coef2 = 0.4 // 0.1 for CG and 0.9 for Newton/quasi-Newton. Using 0.4 though as in SciPy
	o.CoefQuad = 0.1
	o.CoefCubic = 0.2
	o.ffcn = ffcn
	o.Jfcn = Jfcn
	o.xnew = la.NewVector(ndim)
	o.dfdx = la.NewVector(ndim)
	o.interp2 = fun.NewInterpQuad()
	o.interp3 = fun.NewInterpCubic()
	o.interp3.TolDen = 1e-20
	return
}

// SetParams sets parameters
//   Example:
//             o.SetParams(dbf.NewParams(
//                 &dbf.P{N: "maxitls", V: 10},
//                 &dbf.P{N: "maxitzoom", V: 10},
//                 &dbf.P{N: "maxalpha", V: 100},
//                 &dbf.P{N: "mulalpha", V: 2},
//                 &dbf.P{N: "coef1", V: 1e-4},
//                 &dbf.P{N: "coef2", V: 0.4},
//                 &dbf.P{N: "coefquad", V: 0.1},
//                 &dbf.P{N: "coefcubic", V: 0.2},
//             ))
func (o *LineSearch) SetParams(params dbf.Params) {
	o.MaxIt = params.GetIntOrDefault("maxitls", o.MaxIt)
	o.MaxItZoom = params.GetIntOrDefault("maxitzoom", o.MaxItZoom)
	o.MaxAlpha = params.GetValueOrDefault("maxalpha", o.MaxAlpha)
	o.MulAlpha = params.GetValueOrDefault("mulalpha", o.MulAlpha)
	o.Coef1 = params.GetValueOrDefault("coef1", o.Coef1)
	o.Coef2 = params.GetValueOrDefault("coef2", o.Coef2)
	o.CoefQuad = params.GetValueOrDefault("coefquad", o.CoefQuad)
	o.CoefCubic = params.GetValueOrDefault("coefcubic", o.CoefCubic)
}

// Wolfe finds the scalar 'a' that gives a substantial reduction of f({x}+a⋅{u}) (Wolfe conditions)
//
//  Input:
//    x -- initial point
//    u -- direction
//
//  Output:
//    a -- scale parameter
//    f -- f @ a
//    x -- x + a⋅u  [update input x]
//
//  Reference: Algorithm 3.5, page 60 of [1]
//
func (o *LineSearch) Wolfe(x, u la.Vector, useFold bool, fold float64) (a, f float64) {

	// update x
	defer func() {
		la.VecAdd(x, 1, o.x, a, o.u) // xnew := x + a⋅u
	}()

	// set pointers needed by F and G functions
	o.Set(x, u)

	// compute initial F and G
	o.NumFeval = 0
	o.NumJeval = 0
	fini := o.F(0)
	gini := o.G(0)

	// auxiliary
	a0 := 0.0
	f0 := fini
	a1 := 1.0
	var f1, g1 float64

	// estimate a1
	if useFold && gini != 0.0 {
		a1 = utl.Min(1.0, 1.01*2*(fini-fold)/gini)
	}

	// iterations
	for o.NumIter = 1; o.NumIter <= o.MaxIt; o.NumIter++ {

		// compute F(a1)
		f1 = o.F(a1)

		// exit point
		if f1 > fini+o.Coef1*a1*gini || (f1 >= f0 && o.NumIter > 1) {
			a, f = o.zoom(fini, gini, a0, a1, f0, f1)
			return
		}

		// compute G(a1)
		g1 = o.G(a1)

		// exit point
		if math.Abs(g1) <= -o.Coef2*gini {
			a = a1
			f = f1
			return
		}

		// exit point
		if g1 >= 0 {
			a, f = o.zoom(fini, gini, a1, a0, f1, f0)
			return
		}

		// update a0 and a1
		a0 = a1
		f0 = f1
		a1 = utl.Min(o.MulAlpha*a1, o.MaxAlpha)
	}

	// check
	if o.NumIter > o.MaxIt {
		chk.Panic("failed to converge after %d iterations\n", o.NumIter)
	}
	return
}

// zoom generates 'a' between alo and ahi to satisfy the 3 conditions in page 61 of [1]
// Reference: Algorithm 3.6, page 61 of [1]
func (o *LineSearch) zoom(fini, gini, alo, ahi, flo, fhi float64) (a, f float64) {

	// auxiliary
	glo := o.G(alo)
	var g, aprev, fprev float64
	var interpOk bool
	var da, achk, bchk, dchk float64
	var err error

	// iterations
	for o.NumIterZoom = 0; o.NumIterZoom < o.MaxItZoom; o.NumIterZoom++ {

		// variables for checking range
		da = ahi - alo
		achk, bchk = alo, ahi
		if da < 0 {
			achk, bchk = ahi, alo
		}

		// cubic interpolation because we have previous values and the gap=ahi-alo is not small
		interpOk = false
		if o.NumIterZoom > 0 {
			a = cubicmin(alo, flo, glo, ahi, fhi, aprev, fprev)
			dchk = o.CoefCubic * da
			if a >= achk+dchk && a <= bchk-dchk { // accept only if a didn't change much
				interpOk = true
			}
		}

		// quadratic interpolation and the gap=ahi-alo is not small
		if !interpOk {
			err = o.interp2.Fit2pointsD(alo, flo, ahi, fhi, alo, glo)
			if err == nil {
				a, _ = o.interp2.Optimum()
				dchk = o.CoefQuad * da
				if a >= achk+dchk && a <= bchk-dchk { // accept only if a didn't change much
					interpOk = true
				}
			}
		}

		// bi-section because interpolation failed or the gap is too small
		if !interpOk {
			a = alo + 0.5*(ahi-alo)
		}

		// evaluate F
		f = o.F(a)

		// update
		if f > fini+o.Coef1*a*gini || f >= flo {
			aprev = ahi
			fprev = fhi
			ahi = a
			fhi = f
		} else {

			// evaluate G
			g = o.G(a)

			// exit condition
			if math.Abs(g) <= -o.Coef2*gini {
				return
			}

			// swap hi;lo
			if g*(ahi-alo) >= 0 {
				aprev = ahi
				fprev = fhi
				ahi = alo
				fhi = flo
			} else {
				aprev = alo
				fprev = flo
			}

			// next 'lo' values
			alo = a
			flo = f
			glo = g
		}
	}

	// failure
	chk.Panic("zoom did not converge after %d iterations\n", o.NumIterZoom)
	return
}

// Set sets x and u vectors as required by G(a) and H(a) functions
func (o *LineSearch) Set(x, u la.Vector) {
	o.x = x
	o.u = u
}

// F implements f(a) := f({xnew}(a,u)) where {xnew}(a,u) := {x} + a⋅{u}
func (o *LineSearch) F(a float64) float64 {
	o.NumFeval++
	la.VecAdd(o.xnew, 1, o.x, a, o.u) // xnew := x + a⋅u
	return o.ffcn(o.xnew)
}

// G implements g(a) = df/da|({xnew}(a,u)) = df/d{xnew}⋅d{xnew}/da where {xnew} == {x} + a⋅{u}
func (o *LineSearch) G(a float64) float64 {
	o.NumJeval++
	la.VecAdd(o.xnew, 1, o.x, a, o.u) // xnew := x + a⋅u
	o.Jfcn(o.dfdx, o.xnew)            // dfdx @ xnew
	return la.VecDot(o.dfdx, o.u)     // dfdx ⋅ u
}

// PlotC plots contour for current x and u vectors
//   i, j -- the indices in x[i] and x[j] to plot xnew[j] versus xnew[i] with xnew = x + a⋅u
func (o *LineSearch) PlotC(i, j int, x, u la.Vector, a, ximin, ximax, xjmin, xjmax float64, npts int) {

	// auxiliary
	la.VecAdd(o.xnew, 1, o.x, a, o.u) // xnew := x + a⋅u
	x2d := []float64{o.x[i], o.x[j]}
	u2d := []float64{o.u[i], o.u[j]}
	xvec := la.NewVector(len(o.x))
	copy(xvec, o.xnew)

	// meshgrid
	xx, yy, zz := utl.MeshGrid2dF(ximin, ximax, xjmin, xjmax, npts, npts, func(r, s float64) float64 {
		xvec[i], xvec[j] = r, s
		return o.ffcn(xvec)
	})

	// contour
	plt.ContourF(xx, yy, zz, nil)
	plt.PlotOne(x[i], x[j], &plt.A{C: "y", M: "o", NoClip: true})
	plt.PlotOne(o.xnew[i], o.xnew[j], &plt.A{C: "y", M: "*", Ms: 10, NoClip: true})
	plt.DrawArrow2d(x2d, u2d, false, 1, &plt.A{C: "k"})
	plt.DrawArrow2d(x2d, u2d, false, a, nil)
	plt.Gll(io.Sf("$x_{%d}$", i), io.Sf("$x_{%d}$", j), nil)
}

// PlotF plots f(a) curve for current x and u vectors
func (o *LineSearch) PlotF(a, amin, amax float64, npts int) {
	ll := utl.LinSpace(amin-0.001, amax+0.001, npts)
	ff := utl.GetMapped(ll, o.F)
	plt.Plot(ll, ff, &plt.A{C: plt.C(0, 0), NoClip: true})
	plt.PlotOne(amin, o.F(amin), &plt.A{C: "r", Mew: 2, M: "|", Ms: 40, NoClip: true})
	plt.PlotOne(amax, o.F(amax), &plt.A{C: "r", Mew: 2, M: "|", Ms: 40, NoClip: true})
	plt.PlotOne(0, o.F(0), &plt.A{C: "y", M: "o", NoClip: true})
	plt.PlotOne(a, o.F(a), &plt.A{C: "y", M: "*", Ms: 10, NoClip: true})
	plt.Cross(0, 0, nil)
	plt.Gll("$a$", "$f(a)$", nil)
}

// temporary ///////////////////////////////////////////////////////////////////////////////////////

func cubicmin(a, fa, fpa, b, fb, c, fc float64) (xmin float64) {
	C := fpa
	db := b - a
	dc := c - a
	denom := fun.Pow2(db*dc) * (db - dc)
	d00 := +fun.Pow2(dc)
	d01 := -fun.Pow2(db)
	d10 := -fun.Pow3(dc)
	d11 := +fun.Pow3(db)
	v0 := fb - fa - C*db
	v1 := fc - fa - C*dc
	A := d01*v1 + d00*v0
	B := d11*v1 + d10*v0
	A /= denom
	B /= denom
	radical := B*B - 3*A*C
	xmin = a + (-B+math.Sqrt(radical))/(3*A)
	return
}
