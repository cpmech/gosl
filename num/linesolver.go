// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// LineSolver finds the scalar λ that zeroes or minimizes f(x+λ⋅n)
type LineSolver struct {

	// inernal
	ffcn    fun.Sv    // scalar function of vector: y = f({x})
	Jfcn    fun.Vv    // vector function of vector: {J} = df/d{x} @ {x} [optional / may be nil]
	y       la.Vector // {y} = {x} + λ⋅{n}
	dfdx    la.Vector // derivative df/d{x}
	bracket *Bracket  // bracket
	solver  *Brent    // scalar minimizer

	// Stat
	NumFeval int // number of function evalutions
	NumJeval int // number of Jacobian evaluations

	// pointers
	x la.Vector // starting point
	n la.Vector // direction
}

// NewLineSolver returns a new LineSolver object
//   size -- length(x)
//   ffcn -- scalar function of vector: y = f({x})
//   Jfcn -- vector function of vector: {J} = df/d{x} @ {x} [optional / may be nil]
func NewLineSolver(size int, ffcn fun.Sv, Jfcn fun.Vv) (o *LineSolver) {
	o = new(LineSolver)
	o.ffcn = ffcn
	o.Jfcn = Jfcn
	o.y = la.NewVector(size)
	o.dfdx = la.NewVector(size)
	o.bracket = NewBracket(o.G)
	o.solver = NewBrent(o.G, o.H)
	return
}

// Root finds the scalar λ that zeroes f(x+λ⋅n)
func (o *LineSolver) Root(x, n la.Vector) (λ float64) {
	o.Set(x, n)
	λ = o.solver.Root(0, 1)
	o.NumFeval = o.solver.NumFeval
	o.NumJeval = o.solver.NumJeval
	return
}

// Min finds the scalar λ that minimizes f(x+λ⋅n)
func (o *LineSolver) Min(x, n la.Vector) (λ float64) {
	o.Set(x, n)
	λmin, _, λmax, _, _, _ := o.bracket.Min(0, 1)
	λ = o.solver.Min(λmin, λmax)
	o.NumFeval = o.solver.NumFeval
	o.NumJeval = o.solver.NumJeval
	return
}

// MinUpdateX finds the scalar λ that minimizes f(x+λ⋅n), updates x and returns fmin = f({x})
//  Input:
//    x -- initial point
//    n -- direction
//  Output:
//    λ -- scale parameter
//    x -- x @ minimum
//    fmin -- f({x})
func (o *LineSolver) MinUpdateX(x, n la.Vector) (λ, fmin float64) {
	λ = o.Min(x, n)
	la.VecAdd(o.x, 1, x, λ, n) // x := x + λ⋅n
	fmin = o.ffcn(o.x)
	o.NumFeval = o.solver.NumFeval + 1
	o.NumJeval = o.solver.NumJeval
	return
}

// Set sets x and n vectors as required by G(λ) and H(λ) functions
func (o *LineSolver) Set(x, n la.Vector) {
	o.x = x
	o.n = n
}

// G implements g(λ) := f({y}(λ)) where {y}(λ) := {x} + λ⋅{n}
func (o *LineSolver) G(λ float64) float64 {
	la.VecAdd(o.y, 1, o.x, λ, o.n) // xpn := x + λ⋅n
	return o.ffcn(o.y)
}

// H implements h(λ) = dg/dλ = df/d{y} ⋅ d{y}/dλ where {y} == {x} + λ⋅{n}
func (o *LineSolver) H(λ float64) float64 {
	la.VecAdd(o.y, 1, o.x, λ, o.n) // y := x + λ⋅n
	o.Jfcn(o.dfdx, o.y)            // dfdx @ y
	return la.VecDot(o.dfdx, o.n)  // dfdx ⋅ n
}

// PlotC plots contour for current x and n vectors
//   i, j -- the indices in x[i] and x[j] to plot x[j] versus x[i]
func (o *LineSolver) PlotC(i, j int, λ, xmin, xmax, ymin, ymax float64, npts int) {

	// auxiliary
	x2d := []float64{o.x[i], o.x[j]}
	n2d := []float64{o.n[i], o.n[j]}
	xvec := la.NewVector(len(o.x))
	copy(xvec, o.x)

	// contour
	xx, yy, zz := utl.MeshGrid2dF(xmin, xmax, ymin, ymax, npts, npts, func(u, v float64) float64 {
		xvec[i], xvec[j] = u, v
		return o.ffcn(xvec)
	})
	plt.ContourF(xx, yy, zz, nil)
	plt.PlotOne(x2d[0]+λ*n2d[0], x2d[1]+λ*n2d[1], &plt.A{C: "y", M: "o", NoClip: true})
	plt.DrawArrow2d(x2d, n2d, false, 1, nil)
	plt.Gll(io.Sf("$x_{%d}$", i), io.Sf("$x_{%d}$", j), nil)
	plt.HideTRborders()
}

// PlotG plots g(λ) curve for current x and n vectors
func (o *LineSolver) PlotG(λ, λmin, λmax float64, npts int, useBracket bool) {

	// auxiliary
	if useBracket {
		λmin, _, λmax, _, _, _ = o.bracket.Min(0, 1)
	}
	ll := utl.LinSpace(λmin, λmax, npts)

	// scalar function along n
	gg := utl.GetMapped(ll, o.G)
	plt.Plot(ll, gg, &plt.A{C: plt.C(0, 0), L: "$g(\\lambda)$", NoClip: true})
	plt.PlotOne(λmin, o.G(λmin), &plt.A{C: "r", M: "|", Ms: 40, NoClip: true})
	plt.PlotOne(λmax, o.G(λmax), &plt.A{C: "r", M: "|", Ms: 40, NoClip: true})
	plt.PlotOne(λ, o.G(λ), &plt.A{C: "y", M: "o", NoClip: true})
	plt.Cross(0, 0, nil)
	plt.Gll("$\\lambda$", "$g(\\lambda)$", nil)
	plt.HideTRborders()
}
