// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
)

// LineSearch finds a new point x along the direction dx, from x0, where the function
// has decreased sufficiently. The new function value is returned in fx
//  Input:
//      ffcn      -- f(x) callback
//      dx        -- direction vector
//      x0        -- initial x
//      dφdx0     -- initial dφdx0 = fx * dfdx
//      φ0        -- initial φ = 0.5 * dot(fx,fx)
//      max_it    -- max number of iterations
//      dx_is_mdx -- whether dx is actually -dx ==> IMPORTANT: dx will then be changed dx := -dx
//  Output:
//      x      -- updated x (along dx)
//      fx     -- updated f(x)
//      φ0     -- updated φ = 0.5 * dot(fx,fx)
//      dx     -- changed to -dx if dx_is_mdx == true
//      nFeval -- number of calls to f(x)
//  Local constants:
//      tol_gra_min -- tolerance to consider local minimum
//      mul_dx_max  -- multiplier to control maximum dx
//      slope_max   -- ~0 but < 0
//      α           -- Armijo coefficient
//      ε           -- machine epsilon
func LineSearch(x, fx []float64, ffcn fun.Vv, dx, x0, dφdx0 []float64, φ0 float64, max_it int, dx_is_mdx bool) (nFeval int, err error) {

	// tolerances
	tol_gra_min := 1e-12
	mul_dx_max := 100.0
	slope_max := -MACHEPS
	slope_max = 0.0

	// constants
	α := 1e-4  // Armijo coefficient
	ε := 1e-16 // machine epsilon

	// scale dx if step is too big
	n := len(x0)
	var nrm_x0, nrm_dx float64
	for i := 0; i < n; i++ {
		if dx_is_mdx {
			dx[i] = -dx[i]
		}
		nrm_x0 += x0[i] * x0[i]
		nrm_dx += dx[i] * dx[i]
	}
	nrm_x0, nrm_dx = math.Sqrt(nrm_x0), math.Sqrt(nrm_dx)
	nrm_dx_max := mul_dx_max * max(nrm_x0, float64(n))
	if nrm_dx > nrm_dx_max {
		for i := 0; i < n; i++ {
			dx[i] *= nrm_dx_max / nrm_dx // scale if attempted step is to big
		}
	}

	// descent slope and λ min
	var slope, max_val, tmp float64
	for i := 0; i < n; i++ {
		slope += dφdx0[i] * dx[i]
		tmp = math.Abs(dx[i]) / max(math.Abs(x0[i]), 1.0)
		if tmp > max_val {
			max_val = tmp
		}
	}
	λ_min := ε / max_val

	// check slope on the direction of dx
	if slope > slope_max {
		return nFeval, chk.Err(_linesearch_err1, slope)
	}

	// iterations
	var λ, φ, λ2, φ2, gra, den, r1, r2, a, b, d float64
	λ = 1.0 // always try full step first
	var it int
	for it = 0; it < max_it; it++ {

		// update search
		for i := 0; i < n; i++ {
			x[i] = x0[i] + λ*dx[i]
		}
		err = ffcn(fx, x)
		nFeval += 1
		if err != nil {
			return
		}

		// compute φ
		φ = 0.0
		for i := 0; i < n; i++ {
			φ += fx[i] * fx[i]
		}
		φ *= 0.5

		// dx is too small
		if λ < λ_min {
			// check for spurious convergence (local minimum)
			gra = 0.0
			den = max(φ, 0.5*float64(n))
			for i := 0; i < n; i++ {
				tmp = math.Abs(dφdx0[i]) * max(math.Abs(x[i]), 1.0) / den
				if tmp > gra {
					gra = tmp
				}
			}
			if gra < tol_gra_min {
				return nFeval, chk.Err(_linesearch_err2, λ, λ_min, gra)
			}
			return // converged
		}

		// converged? (sufficient function decrease)
		if φ <= φ0+α*λ*slope {
			return
		}

		// backtrack
		if it == 0 {
			tmp = -0.5 * slope / (φ - φ0 - slope)
		} else {
			r1 = φ - φ0 - λ*slope
			r2 = φ2 - φ0 - λ2*slope
			a = (r1/(λ*λ) - r2/(λ2*λ2)) / (λ - λ2)
			b = (-λ2*r1/(λ*λ) + λ*r2/(λ2*λ2)) / (λ - λ2)
			if math.Abs(a) < ε {
				tmp = -0.5 * slope / b
			} else {
				d = b*b - 3.0*a*slope
				if d < 0.0 {
					tmp = 0.5 * λ
				} else if b <= 0.0 {
					tmp = (-b + math.Sqrt(d)) / (3.0 * a)
				} else {
					tmp = -slope / (b + math.Sqrt(d))
				}
			}
			tmp = min(tmp, 0.5*λ) // make sure tmp is smaller than 0.5*λ
		}

		// save previous values
		λ2, φ2 = λ, φ

		// new λ
		λ = max(tmp, 0.1*λ) // make sure λ is greater than 0.1*λ
	}

	// check convergence
	if it == max_it {
		return nFeval, chk.Err(_linesearch_err3, it+1)
	}

	return
}

// error messages
var (
	_linesearch_err1 = "linesearch.go: LineSearch: slope must be negative (%g is invalid)"
	_linesearch_err2 = "linesearch.go: LineSearch: local mininum reached? λ=%g, λ_min=%g, gra=%g"
	_linesearch_err3 = "linesearch.go: LineSearch: failed to converge after %d iterations"
)
