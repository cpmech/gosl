// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
)

// The algorithms below are based on [1]
// REFERENCES:
// [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//     Scientific Computing. Third Edition. Cambridge University Press. 1235p.

// QuadElementary defines the interface for elementary quadrature algorithms with refinement.
type QuadElementary interface {
	Init(f fun.Ss, a, b, eps float64) // The constructor takes as inputs f, the function or functor to be integrated between limits a and b, also input.
	Integrate() (float64, error)      // Returns the integral for the specified input data
}

// Trap structure is used for the trapezoidal integration rule with refinement.
type ElementaryTrapz struct {
	n    int     // current level of refinement.
	a, b float64 // limits
	s    float64 // current value of the integral
	eps  float64 // precision
	f    fun.Ss  // the function
}

// Init initialises Trap structure
func (o *ElementaryTrapz) Init(f fun.Ss, a, b, eps float64) {
	o.n = 0
	o.f = f
	o.a = a
	o.b = b
	o.eps = eps
}

// Next returns the nth stage of refinement of the extended trapezoidal rule. On the first call (n=1),
// R b the routine returns the crudest estimate of a f .x/dx. Subsequent calls set n=2,3,... and
// improve the accuracy by adding 2 n-2 additional interior points.
func (o *ElementaryTrapz) Next() (res float64, err error) {
	var x, sum, del float64
	var it, j, tnm int
	o.n += 1
	var fa, fb, fx float64
	if o.n == 1 {
		fa, err = o.f(o.a)
		if err != nil {
			return
		}
		fb, err = o.f(o.b)
		if err != nil {
			return
		}
		o.s = 0.5 * (o.b - o.a) * (fa + fb)
		return o.s, nil
	} else {
		for it, j = 1, 1; j < o.n-1; j++ {
			it *= 2
		}
		tnm = it
		del = (o.b - o.a) / float64(tnm)

		// This is the spacing of the points to be added.
		x = o.a + 0.5*del

		for sum, j = 0.0, 0; j < it; j, x = j+1, x+del {
			fx, err = o.f(x)
			if err != nil {
				return
			}
			sum += fx
		}
		o.s = 0.5 * (o.s + (o.b-o.a)*sum/float64(tnm))

		// This replaces s by its refined value.
		return o.s, nil
	}
}

// Integrate performs the numerical integration
func (o *ElementaryTrapz) Integrate() (res float64, err error) {
	jmax := 20
	var olds float64
	for j := 0; j < jmax; j++ {
		o.s, err = o.Next()
		if err != nil {
			return
		}
		if j > 5 {
			if math.Abs(o.s-olds) < o.eps*math.Abs(olds) || (o.s == 0 && olds == 0) {
				return o.s, nil
			}
		}
		olds = o.s
	}
	return 0, chk.Err("achieved maximum number of iterations (n=%d)", jmax)
}

// Simp structure implements the Simpson's method for quadrature with refinement.
type ElementarySimpson struct {
	n    int     // current level of refinement.
	a, b float64 // limits
	s    float64 // current value of the integral
	eps  float64 // precision
	f    fun.Ss  // the function
}

// Init initialises Simp structure
func (o *ElementarySimpson) Init(f fun.Ss, a, b, eps float64) {
	o.n = 0
	o.f = f
	o.a = a
	o.b = b
	o.eps = eps
}

// Next returns the nth stage of refinement of the extended trapezoidal rule. On the first call (n=1),
// R b the routine returns the crudest estimate of a f .x/dx. Subsequent calls set n=2,3,... and
// improve the accuracy by adding 2 n-2 additional interior points.
func (o *ElementarySimpson) Next() (res float64, err error) {
	var x, sum, del, fa, fb, fx float64
	var it, j, tnm int
	o.n += 1
	if o.n == 1 {
		fa, err = o.f(o.a)
		if err != nil {
			return
		}
		fb, err = o.f(o.b)
		if err != nil {
			return
		}
		o.s = 0.5 * (o.b - o.a) * (fa + fb)
		return o.s, nil
	} else {
		for it, j = 1, 1; j < o.n-1; j++ {
			it *= 2
		}
		tnm = it
		del = (o.b - o.a) / float64(tnm)

		// This is the spacing of the points to be added.
		x = o.a + 0.5*del

		for sum, j = 0.0, 0; j < it; j, x = j+1, x+del {
			fx, err = o.f(x)
			if err != nil {
				return
			}
			sum += fx
		}
		o.s = 0.5 * (o.s + (o.b-o.a)*sum/float64(tnm))

		// This replaces s by its refined value.
		return o.s, nil
	}
}

// Integrate performs the numerical integration
func (o *ElementarySimpson) Integrate() (res float64, err error) {
	jmax := 20
	var s, st, ost, os float64
	for j := 0; j < jmax; j++ {
		st, err = o.Next()
		if err != nil {
			return
		}
		s = (4*st - ost) / 3
		if j > 5 {
			if math.Abs(s-os) < o.eps*math.Abs(os) || (s == 0 && os == 0) {
				return s, nil
			}
		}
		os = s
		ost = st
	}
	return 0, chk.Err("achieved maximum number of iterations (n=%d)", jmax)
}
