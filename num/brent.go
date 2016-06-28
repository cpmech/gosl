// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// Brent implements Brent's method for finding the roots of an equation
type Brent struct {
	MaxIt  int     // max iterations
	Tol    float64 // tolerance
	Ffcn   Cb_yxe  // y = f(x) function
	NFeval int     // number of calls to Ffcn (function evaluations)
	It     int     // number of iterations from last call to Solve
	sqeps  float64 // sqrt(EPS)
	gsr    float64 // gold section ratio
}

// Init intialises Brent structure
func (o *Brent) Init(ffcn Cb_yxe) {
	o.MaxIt = 30
	o.Tol = 1e-14
	o.Ffcn = ffcn
	o.gsr = (3.0 - math.Sqrt(5.0)) / 2.0
	o.sqeps = math.Sqrt(EPS)
}

// Solve solves y(x) = 0 for x in [xa, xb] with f(xa) * f(xb) < 0
//
//  Based on ZEROIN C math library: http://www.netlib.org/c/
//  By: Oleg Keselyov <oleg@ponder.csci.unt.edu, oleg@unt.edu> May 23, 1991
//
//   G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
//   computations. M., Mir, 1980, p.180 of the Russian edition
//
//   The function makes use of the bissection procedure combined with
//   the linear or quadric inverse interpolation.
//   At every step program operates on three abscissae - a, b, and c.
//   b - the last and the best approximation to the root
//   a - the last but one approximation
//   c - the last but one or even earlier approximation than a that
//       1) |f(b)| <= |f(c)|
//       2) f(b) and f(c) have opposite signs, i.e. b and c confine
//          the root
//   At every step Zeroin selects one of the two new approximations, the
//   former being obtained by the bissection procedure and the latter
//   resulting in the interpolation (if a,b, and c are all different
//   the quadric interpolation is utilized, otherwise the linear one).
//   If the latter (i.e. obtained by the interpolation) point is
//   reasonable (i.e. lies within the current interval [b,c] not being
//   too close to the boundaries) it is accepted. The bissection result
//   is used in the other case. Therefore, the range of uncertainty is
//   ensured to be reduced at least by the factor 1.6
//
func (o *Brent) Solve(xa, xb float64, silent bool) (res float64, err error) {

	// basic variables and function evaluation
	a := xa // the last but one approximation
	b := xb // the last and the best approximation to the root
	c := a  // the last but one or even earlier approximation than a that
	fa, erra := o.Ffcn(a)
	fb, errb := o.Ffcn(b)
	o.NFeval = 2
	if erra != nil {
		return 0, chk.Err(_brent_err1, "a", xa, erra.Error())
	}
	if errb != nil {
		return 0, chk.Err(_brent_err1, "b", xb, errb.Error())
	}
	fc := fa

	// check input
	if fa*fb >= -EPS {
		return 0, chk.Err(_brent_err2, xa, xb, fa, fb)
	}

	// message
	if !silent {
		io.Pfpink("%4s%23s%23s%23s\n", "it", "x", "f(x)", "err")
		io.Pfpink("%50s%23.1e\n", "", o.Tol)
	}

	// solve
	var prev_step float64  // distance from the last but one to the last approximation
	var tol_act float64    // actual tolerance
	var p, q float64       // interpol. step is calculated in the form p/q (divisions are delayed)
	var new_step float64   // step at this iteration
	var t1, cb, t2 float64 // auxiliary variables
	for o.It = 0; o.It < o.MaxIt; o.It++ {

		// distance
		prev_step = b - a

		// swap data for b to be the best approximation
		if math.Abs(fc) < math.Abs(fb) {
			a = b
			b = c
			c = a
			fa = fb
			fb = fc
			fc = fa
		}
		tol_act = 2.0*EPS*math.Abs(b) + o.Tol/2.0
		new_step = (c - b) / 2.0

		// converged?
		if !silent {
			io.Pfyel("%4d%23.15e%23.15e%23.15e\n", o.It, b, fb, math.Abs(new_step))
		}
		if math.Abs(new_step) <= tol_act || fb == 0.0 {
			return b, nil
		}

		// decide if the interpolation can be tried
		if math.Abs(prev_step) >= tol_act && math.Abs(fa) > math.Abs(fb) {
			// if prev_step was large enough and was in true direction, interpolatiom may be tried
			cb = c - b

			// with two distinct points, linear interpolation must be applied
			if a == c {
				t1 = fb / fa
				p = cb * t1
				q = 1.0 - t1

				// otherwise, quadric inverse interpolation is applied
			} else {
				q = fa / fc
				t1 = fb / fc
				t2 = fb / fa
				p = t2 * (cb*q*(q-t1) - (b-a)*(t1-1.0))
				q = (q - 1.0) * (t1 - 1.0) * (t2 - 1.0)
			}

			// p was calculated with the opposite sign;
			// make p positive and assign possible minus to q
			if p > 0.0 {
				q = -q
			} else {
				p = -p
			}

			// if b+p/q falls in [b,c] and isn't too large, it is accepted
			// if p/q is too large then the bissection procedure can reduce [b,c] range to more extent
			if p < (0.75*cb*q-math.Abs(tol_act*q)/2.0) && p < math.Abs(prev_step*q/2.0) {
				new_step = p / q
			}
		}

		// adjust the step to be not less than tolerance
		if math.Abs(new_step) < tol_act {
			if new_step > 0.0 {
				new_step = tol_act
			} else {
				new_step = -tol_act
			}
		}

		// save the previous approximation
		a = b
		fa = fb

		// do step to a new approximation
		b += new_step
		fb, errb = o.Ffcn(b)
		o.NFeval += 1
		if errb != nil {
			return 0, chk.Err(_brent_err1, "", b, errb.Error())
		}

		// adjust c for it to have a sign opposite to that of b
		if (fb > 0.0 && fc > 0.0) || (fb < 0.0 && fc < 0.0) {
			c = a
			fc = fa
		}
	}

	// did not converge
	return fb, chk.Err(_brent_err3, "Solve", o.It)
}

// Min finds the minimum of f(x) in [xa, xb]
//
//  Based on ZEROIN C math library: http://www.netlib.org/c/
//  By: Oleg Keselyov <oleg@ponder.csci.unt.edu, oleg@unt.edu> May 23, 1991
//
//   G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
//   computations. M., Mir, 1980, p.202 of the Russian edition
//
//   The function makes use of the "gold section" procedure combined with
//   the parabolic interpolation.
//   At every step program operates three abscissae - x,v, and w.
//   x - the last and the best approximation to the minimum location,
//       i.e. f(x) <= f(a) or/and f(x) <= f(b)
//       (if the function f has a local minimum in (a,b), then the both
//       conditions are fulfiled after one or two steps).
//   v,w are previous approximations to the minimum location. They may
//   coincide with a, b, or x (although the algorithm tries to make all
//   u, v, and w distinct). Points x, v, and w are used to construct
//   interpolating parabola whose minimum will be treated as a new
//   approximation to the minimum location if the former falls within
//   [a,b] and reduces the range enveloping minimum more efficient than
//   the gold section procedure.
//   When f(x) has a second derivative positive at the minimum location
//   (not coinciding with a or b) the procedure converges superlinearly
//   at a rate order about 1.324
//
//   The function always obtains a local minimum which coincides with
//   the global one only if a function under investigation being
//   unimodular. If a function being examined possesses no local minimum
//   within the given range, Fminbr returns 'a' (if f(a) < f(b)), otherwise
//   it returns the right range boundary value b.
//
func (o *Brent) Min(xa, xb float64, silent bool) (res float64, err error) {

	// check
	if xb < xa {
		return 0, chk.Err(_brent_err4, xa, xb)
	}

	// first step: always gold section
	v := xa + o.gsr*(xb-xa)
	fv, errv := o.Ffcn(v)
	o.NFeval = 1
	if errv != nil {
		return 0, chk.Err(_brent_err5, errv.Error())
	}
	x, w, fx, fw := v, v, fv, fv

	// solve
	var rng float64         // range over which the minimum is seeked for
	var mid_rng float64     // middle range
	var tol_act float64     // actual tolerance
	var new_step float64    // step at one iteration
	var tmp float64         // temporary
	var p, q, t, ft float64 // auxiliary
	for o.It = 0; o.It < o.MaxIt; o.It++ {

		// auxiliary variables
		rng = xb - xa
		mid_rng = (xa + xb) / 2.0
		tol_act = o.sqeps*math.Abs(x) + o.Tol/3.0

		// converged?
		if !silent {
			io.Pfyel("%4d%23.15e%23.15e%23.15e\n", o.It, x, fx, math.Abs(x-mid_rng)+rng/2.0)
		}
		if math.Abs(x-mid_rng)+rng/2.0 <= 2.0*tol_act {
			return x, nil
		}

		// Obtain the gold section step
		tmp = xa - x
		if x < mid_rng {
			tmp = xb - x
		}
		new_step = o.gsr * tmp

		// decide if the interpolation can be tried
		if math.Abs(x-w) >= tol_act { // if x and w are distinct interpolatiom may be tried

			t = (x - w) * (fx - fv)
			q = (x - v) * (fx - fw)
			p = (x-v)*q - (x-w)*t
			q = 2.0 * (q - t)

			if q > 0.0 { // q was calculated with the op posite sign; make q positive and assign possible minus to p
				p = -p
			} else {
				q = -q
			}

			// x+p/q falls in [a,b] not too close to a and b, and isn't too large
			if math.Abs(p) < math.Abs(new_step*q) && p > q*(xa-x+2.0*tol_act) && p < q*(xb-x-2.0*tol_act) {
				new_step = p / q // it is accepted
				// if p/q is too large then the gold section procedure can reduce [a,b] rng to more extent
			}
		}

		// adjust the step to be not less than tolerance
		if math.Abs(new_step) < tol_act {
			if new_step > 0.0 {
				new_step = tol_act
			} else {
				new_step = -tol_act
			}
		}

		// obtain the next approximation to min and reduce the enveloping rng
		t = x + new_step // tentative point for the min  */
		ft, err = o.Ffcn(t)
		if err != nil {
			return 0, chk.Err(_brent_err5, err.Error())
		}

		// t is a better approximation
		if ft <= fx {
			if t < x { // reduce the range so that t would fall within it
				xb = x
			} else {
				xa = x
			}
			// assign the best approx to x
			v = w
			w = x
			x = t
			fv = fw
			fw = fx
			fx = ft

			// x remains the better approx
		} else {
			if t < x { // reduce the range enclosing x
				xa = t
			} else {
				xb = t
			}
			if ft <= fw || w == x {
				v = w
				w = t
				fv = fw
				fw = ft
			} else if ft <= fv || v == x || v == w {
				v = t
				fv = ft
			}
		}
	}

	// did not converge
	return x, chk.Err(_brent_err3, "Min", o.It)
}

// error messages
var (
	_brent_err1 = "brent.go: Brent.Solve: f%s(%g) failed:\n%v"
	_brent_err2 = "brent.go: Brent.Solve: root must be bracketed: xa=%g, xb=%g, fa=%g, fb=%b => fa * fb >= 0"
	_brent_err3 = "brent.go: Brent.%s: dit not converge after %d iterations"
	_brent_err4 = "min1d.go: Brent.FindMin: xa(%g) must be smaller than xb(%g)"
	_brent_err5 = "min1d.go: Brent.FindMin: f(%g) failed:\n%v"
)
