// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
)

// LineSolver finds the scalar λ that zeroes or minimizes f(x+λ⋅n)
type LineSolver struct {

	// inernal
	ffcn    fun.Sv    // scalar function of vector: y = f({x})
	Jfcn    fun.Vv    // vector function of vector: {J} = df/d{x} @ {x} [optional / may be nil]
	y       la.Vector // {y} = {x} + λ⋅{n}
	dfdx    la.Vector // derivative df/d{x}
	bracket *Bracket  // bracket

	// minimizer
	Brent *Brent // scalar minimizer

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
	o.Brent = NewBrent(o.G, o.H)
	return
}

// Root finds the scalar λ that zeroes f(x+λ⋅n)
func (o *LineSolver) Root(x, n la.Vector) (λ float64) {
	o.Set(x, n)
	return o.Brent.Root(0, 1)
}

// Min finds the scalar λ that minimizes f(x+λ⋅n)
func (o *LineSolver) Min(x, n la.Vector) (λ float64) {
	o.Set(x, n)
	λmin, _, λmax, _, _, _ := o.bracket.Min(0, 1)
	return o.Brent.Min(λmin, λmax)
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
