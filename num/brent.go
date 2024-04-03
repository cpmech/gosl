// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
)

// Brent implements Brent's method for finding the roots of an equation
type Brent struct {

	// configuration
	MaxIt   int     // max iterations
	Tol     float64 // tolerance
	Verbose bool    // show messages

	// statistics
	NumFeval int // number of calls to Ffcn (function evaluations)
	NumJeval int // number of calls to Jfcn (Jacobian/derivatives)
	NumIter  int // number of iterations from last call to Solve

	// internal
	ffcn  fun.Ss  // y = f(x) function
	Jfcn  fun.Ss  // Jfcn(x) = dy/dx [optional / may be nil]
	gsr   float64 // gold section ratio
	sqeps float64 // sqrt(EPS)
}

// NewBrent returns a new Brent structure
//
//	ffcn -- function f(x)
//	Jfcn -- derivative df(x)/dx [optional / may be nil]
func NewBrent(ffcn, Jfcn fun.Ss) (o *Brent) {
	o = new(Brent)
	o.MaxIt = 100
	o.Tol = 1e-10
	o.ffcn = ffcn
	o.Jfcn = Jfcn
	o.gsr = (3.0 - math.Sqrt(5.0)) / 2.0
	o.sqeps = math.Sqrt(MACHEPS)
	return
}

// Root solves y(x) = 0 for x in [xa, xb] with f(xa) * f(xb) < 0
//
//	Based on ZEROIN C math library: http://www.netlib.org/c/
//	By: Oleg Keselyov <oleg@ponder.csci.unt.edu, oleg@unt.edu> May 23, 1991
//
//	 G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
//	 computations. M., Mir, 1980, p.180 of the Russian edition
//
//	 The function makes use of the bissection procedure combined with
//	 the linear or quadric inverse interpolation.
//	 At every step program operates on three abscissae - a, b, and c.
//	 b - the last and the best approximation to the root
//	 a - the last but one approximation
//	 c - the last but one or even earlier approximation than a that
//	     1) |f(b)| <= |f(c)|
//	     2) f(b) and f(c) have opposite signs, i.e. b and c confine
//	        the root
//	 At every step Zeroin selects one of the two new approximations, the
//	 former being obtained by the bissection procedure and the latter
//	 resulting in the interpolation (if a,b, and c are all different
//	 the quadric interpolation is utilized, otherwise the linear one).
//	 If the latter (i.e. obtained by the interpolation) point is
//	 reasonable (i.e. lies within the current interval [b,c] not being
//	 too close to the boundaries) it is accepted. The bissection result
//	 is used in the other case. Therefore, the range of uncertainty is
//	 ensured to be reduced at least by the factor 1.6
func (o *Brent) Root(xa, xb float64) (res float64) {

	// basic variables and function evaluation
	a := xa // the last but one approximation
	b := xb // the last and the best approximation to the root
	c := a  // the last but one or even earlier approximation than a that
	fa := o.ffcn(a)
	fb := o.ffcn(b)
	o.NumFeval = 2
	o.NumJeval = 0
	fc := fa

	// check input
	if fa*fb >= -MACHEPS {
		chk.Panic("root must be bracketed: xa=%g, xb=%g, fa=%g, fb=%b => fa * fb >= 0", xa, xb, fa, fb)
	}

	// message
	if o.Verbose {
		io.Pf("%4s%23s%23s%23s\n", "it", "x", "f(x)", "err")
		io.Pf("%50s%23.1e\n", "", o.Tol)
	}

	// solve
	var prevStep float64   // distance from the last but one to the last approximation
	var tolAct float64     // actual tolerance
	var p, q float64       // interpol. step is calculated in the form p/q (divisions are delayed)
	var newStep float64    // step at this iteration
	var t1, cb, t2 float64 // auxiliary variables
	for o.NumIter = 0; o.NumIter < o.MaxIt; o.NumIter++ {

		// distance
		prevStep = b - a

		// swap data for b to be the best approximation
		if math.Abs(fc) < math.Abs(fb) {
			a = b
			b = c
			c = a
			fa = fb
			fb = fc
			fc = fa
		}
		tolAct = 2.0*MACHEPS*math.Abs(b) + o.Tol/2.0
		newStep = (c - b) / 2.0

		// converged?
		if o.Verbose {
			io.Pf("%4d%23.15e%23.15e%23.15e\n", o.NumIter, b, fb, math.Abs(newStep))
		}
		if math.Abs(newStep) <= tolAct || fb == 0.0 {
			return b
		}

		// decide if the interpolation can be tried
		if math.Abs(prevStep) >= tolAct && math.Abs(fa) > math.Abs(fb) {
			// if prev_step was large enough and was in true direction, interpolation may be tried
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
			if p < (0.75*cb*q-math.Abs(tolAct*q)/2.0) && p < math.Abs(prevStep*q/2.0) {
				newStep = p / q
			}
		}

		// adjust the step to be not less than tolerance
		if math.Abs(newStep) < tolAct {
			if newStep > 0.0 {
				newStep = tolAct
			} else {
				newStep = -tolAct
			}
		}

		// save the previous approximation
		a = b
		fa = fb

		// do step to a new approximation
		b += newStep
		fb = o.ffcn(b)
		o.NumFeval++

		// adjust c for it to have a sign opposite to that of b
		if (fb > 0.0 && fc > 0.0) || (fb < 0.0 && fc < 0.0) {
			c = a
			fc = fa
		}
	}

	// did not converge
	chk.Panic("fail to converge after %d iterations", o.NumIter)
	return
}

// Min finds the minimum of f(x) in [xa, xb]
//
//	Based on ZEROIN C math library: http://www.netlib.org/c/
//	By: Oleg Keselyov <oleg@ponder.csci.unt.edu, oleg@unt.edu> May 23, 1991
//
//	 G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
//	 computations. M., Mir, 1980, p.202 of the Russian edition
//
//	 The function makes use of the "gold section" procedure combined with
//	 the parabolic interpolation.
//	 At every step program operates three abscissae - x,v, and w.
//	 x - the last and the best approximation to the minimum location,
//	     i.e. f(x) <= f(a) or/and f(x) <= f(b)
//	     (if the function f has a local minimum in (a,b), then the both
//	     conditions are fulfilled after one or two steps).
//	 v,w are previous approximations to the minimum location. They may
//	 coincide with a, b, or x (although the algorithm tries to make all
//	 u, v, and w distinct). Points x, v, and w are used to construct
//	 interpolating parabola whose minimum will be treated as a new
//	 approximation to the minimum location if the former falls within
//	 [a,b] and reduces the range enveloping minimum more efficient than
//	 the gold section procedure.
//	 When f(x) has a second derivative positive at the minimum location
//	 (not coinciding with a or b) the procedure converges superlinearly
//	 at a rate order about 1.324
//
//	 The function always obtains a local minimum which coincides with
//	 the global one only if a function under investigation being
//	 unimodular. If a function being examined possesses no local minimum
//	 within the given range, Fminbr returns 'a' (if f(a) < f(b)), otherwise
//	 it returns the right range boundary value b.
func (o *Brent) Min(xa, xb float64) (xAtMin float64) {

	// check
	if xb < xa {
		chk.Panic("xa(%g) must be smaller than xb(%g)", xa, xb)
	}

	// first step: always gold section
	v := xa + o.gsr*(xb-xa)
	fv := o.ffcn(v)
	o.NumFeval = 1
	o.NumJeval = 0
	x, w, fx, fw := v, v, fv, fv

	// solve
	var rng float64         // range over which the minimum is seeked for
	var midRng float64      // middle range
	var tolAct float64      // actual tolerance
	var newStep float64     // step at one iteration
	var tmp float64         // temporary
	var p, q, t, ft float64 // auxiliary
	for o.NumIter = 0; o.NumIter < o.MaxIt; o.NumIter++ {

		// auxiliary variables
		rng = xb - xa
		midRng = (xa + xb) / 2.0
		tolAct = o.sqeps*math.Abs(x) + o.Tol/3.0

		// converged?
		if o.Verbose {
			io.Pf("%4d%23.15e%23.15e%23.15e\n", o.NumIter, x, fx, math.Abs(x-midRng)+rng/2.0)
		}
		if math.Abs(x-midRng)+rng/2.0 <= 2.0*tolAct {
			return x
		}

		// Obtain the gold section step
		tmp = xa - x
		if x < midRng {
			tmp = xb - x
		}
		newStep = o.gsr * tmp

		// decide if the interpolation can be tried
		if math.Abs(x-w) >= tolAct { // if x and w are distinct interpolation may be tried

			t = (x - w) * (fx - fv)
			q = (x - v) * (fx - fw)
			p = (x-v)*q - (x-w)*t
			q = 2.0 * (q - t)

			if q > 0.0 { // q was calculated with the op positive sign; make q positive and assign possible minus to p
				p = -p
			} else {
				q = -q
			}

			// x+p/q falls in [a,b] not too close to a and b, and isn't too large
			if math.Abs(p) < math.Abs(newStep*q) && p > q*(xa-x+2.0*tolAct) && p < q*(xb-x-2.0*tolAct) {
				newStep = p / q // it is accepted
				// if p/q is too large then the gold section procedure can reduce [a,b] rng to more extent
			}
		}

		// adjust the step to be not less than tolerance
		if math.Abs(newStep) < tolAct {
			if newStep > 0.0 {
				newStep = tolAct
			} else {
				newStep = -tolAct
			}
		}

		// obtain the next approximation to min and reduce the enveloping rng
		t = x + newStep // tentative point for the min  */
		ft = o.ffcn(t)
		o.NumFeval++

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
	chk.Panic("fail to converge after %d iterations", o.NumIter)
	return
}

// MinUseD finds minimum and uses information about derivatives
//
//	Given a function and deriva funcd that computes a function and also its derivative function df, and
//	given a bracketing triplet of abscissas ax, bx, cx [such that bx is between ax and cx, and
//	f(bx) is less than both f(ax) and f(cx)], this routine isolates the minimum to a fractional
//	precision of about tol using a modification of Brent’s method that uses derivatives. The
//	abscissa of the minimum is returned as xAtMin, and the minimum function value is returned
//	as min, the returned function value.
//
//	REFERENCES:
//	[1] Press WH, Teukolsky SA, Vetterling WT, Flannery BP (2007) Numerical Recipes:
//	    The Art of Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func (o *Brent) MinUseD(xa, xb float64) (xAtMin float64) {

	// check
	if o.Jfcn == nil {
		chk.Panic("MinUseD requires derivative function Jfcn\n")
	}
	if xb < xa {
		chk.Panic("xa(%g) must be smaller than xb(%g)", xa, xb)
	}

	// bracket
	bracket := NewBracket(o.ffcn)
	ax, bx, cx, _, fb, _ := bracket.Min(xa, xb)

	// auxiliary
	var ok1, ok2 bool // will be used as flags for whether proposed steps are acceptable or not.
	var a, b, d, d1, d2, du, dv, dw, dx, e float64
	var fu, fv, fw, fx, olde, tol1, tol2, u, u1, u2, v, w, x, xm float64

	// initialization
	a, b = cx, cx
	if ax < cx {
		a = ax
	}
	if ax > cx {
		b = ax
	}
	x = bx
	fx = fb
	dx = o.Jfcn(x)
	w, v = x, x
	fw, fv = fx, fx
	dw, dv = dx, dx

	// constants
	zeps := MACHEPS * 1e-3

	// stat
	o.NumFeval = bracket.NumFeval
	o.NumJeval = 1

	// iterations
	for o.NumIter = 0; o.NumIter < o.MaxIt; o.NumIter++ {
		xm = 0.5 * (a + b)
		tol1 = o.Tol*math.Abs(x) + zeps
		tol2 = 2.0 * tol1
		if math.Abs(x-xm) <= (tol2 - 0.5*(b-a)) {
			return x
		}
		if math.Abs(e) > tol1 {
			d1 = 2.0 * (b - a)
			d2 = d1
			if dw != dx {
				d1 = (w - x) * dx / (dx - dw)
			}
			if dv != dx {
				d2 = (v - x) * dx / (dx - dv)
			}
			u1 = x + d1
			u2 = x + d2
			ok1 = (a-u1)*(u1-b) > 0.0 && dx*d1 <= 0.0
			ok2 = (a-u2)*(u2-b) > 0.0 && dx*d2 <= 0.0
			olde = e
			e = d
			if ok1 || ok2 {
				if ok1 && ok2 {
					d = d2
					if math.Abs(d1) < math.Abs(d2) {
						d = d1
					}
				} else if ok1 {
					d = d1
				} else {
					d = d2
				}
				if math.Abs(d) <= math.Abs(0.5*olde) {
					u = x + d
					if u-a < tol2 || b-u < tol2 {
						d = sgn(tol1, xm-x)
					}
				} else {
					e = b - x
					if dx >= 0.0 {
						e = a - x
					}
					d = 0.5 * e
				}
			} else {
				e = b - x
				if dx >= 0.0 {
					e = a - x
				}
				d = 0.5 * e
			}
		} else {
			e = b - x
			if dx >= 0.0 {
				e = a - x
			}
			d = 0.5 * e
		}
		if math.Abs(d) >= tol1 {
			u = x + d
			fu = o.ffcn(u)
			o.NumFeval++
		} else {
			u = x + sgn(tol1, d)
			fu = o.ffcn(u)
			o.NumFeval++
			if fu > fx {
				return x
			}
		}
		du = o.Jfcn(u)
		o.NumJeval++
		if fu <= fx {
			if u >= x {
				a = x
			} else {
				b = x
			}
			mov3(&v, &fv, &dv, w, fw, dw)
			mov3(&w, &fw, &dw, x, fx, dx)
			mov3(&x, &fx, &dx, u, fu, du)
		} else {
			if u < x {
				a = u
			} else {
				b = u
			}
			if fu <= fw || w == x {
				mov3(&v, &fv, &dv, w, fw, dw)
				mov3(&w, &fw, &dw, u, fu, du)
			} else if fu < fv || v == x || v == w {
				mov3(&v, &fv, &dv, u, fu, du)
			}
		}
	}

	// did not converge
	chk.Panic("fail to converge after %d iterations", o.NumIter)
	return
}
