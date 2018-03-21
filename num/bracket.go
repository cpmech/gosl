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

// Bracket implements routines to bracket roots or optima
//
//   A root of a function is known to be bracketed by a pair of points, a and b,
//   when the function has opposite sign at those two points.
//
//   A minimum is known to be bracketed only when there is a triplet of points,
//   a < b < c (or c < b < a), such that f(b) is less than both f(a) and f(c)
//
//   REFERENCES:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes:
//       The Art of Scientific Computing. Third Edition. Cambridge University Press. 1235p.
//
type Bracket struct {

	// configuration
	MaxIt   int  // max iterations
	Verbose bool // show messages

	// statistics
	NFeval int // number of calls to Ffcn (function evaluations)
	It     int // number of iterations from last call to Solve

	// internal
	ffcn   fun.Ss  // y = f(x) function
	gold   float64 // golden ratio: default ratio by which successive intervals are magnified
	glimit float64 // maximum magnification allowed fora parabolic-fit step
	tiny   float64 // small number to prevent division by zero
}

// NewBracket returns a new bracket-er object
func NewBracket(ffcn fun.Ss) (o *Bracket) {
	o = new(Bracket)
	o.MaxIt = 500
	o.ffcn = ffcn
	o.gold = (1.0 + math.Sqrt(5.0)) / 2.0
	o.tiny = math.Sqrt(MACHEPS)
	return
}

// Min brackets minimum
//
//  Given a function and given distinct initial points a0 and b0, search in the downhill direction
//  (defined by the function as evaluated at the initial points) and return new points a, b, c
//  that bracket a minimum of the function.
//
//  Returns also the function values at the three points, fa, fb, and fc
//
func (o *Bracket) Min(a0, b0 float64) (a, b, c, fa, fb, fc float64) {

	// check
	if a0 == b0 {
		chk.Panic("a0=%g must be different than b0=%g\n", a0, b0)
	}

	// sort output
	exitCase := 0
	defer func() {
		if c < a {
			a, c = c, a
			fa, fc = fc, fa
		}
		if o.Verbose {
			txtCase := "main exit point"
			if exitCase == 2 {
				txtCase = "min between u and c"
			}
			if exitCase == 3 {
				txtCase = "min between a and u"
			}
			io.Pf("exit case = %q\n\ta=%g b=%g c=%g\n\tfa=%g fb=%g fc=%g\n\tnIterations=%d nFeval=%d\n", txtCase, a, b, c, fa, fb, fc, o.It+1, o.NFeval)
		}
	}()

	// initialization
	a = a0
	b = b0
	fa = o.ffcn(a)
	fb = o.ffcn(b)
	if fb > fa { // switch roles of a and b so that we can go downhill in the direction from a to b.
		swap(&a, &b)
		swap(&fb, &fa)
	}
	c = b + o.gold*(b-a) // first guess for c
	fc = o.ffcn(c)
	o.NFeval = 3

	// auxiliary
	var r, q, del, den, u, fu, ulim, aux float64

	// search
	for o.It = 0; o.It < o.MaxIt; o.It++ {

		// exit point
		if fb <= fc {
			exitCase = 1
			return
		}

		// parabolic extrapolation
		r = (b - a) * (fb - fc) // compute u by parabolic extrapolation from a,b,c
		q = (b - c) * (fb - fa) // tiny is used to prevent any possible division by zero.
		del = q - r
		if math.Abs(del) < o.tiny {
			del = signp(del) * o.tiny
		}
		den = 2.0 * del
		u = b - ((b-c)*q-(b-a)*r)/den
		ulim = b + o.glimit*(c-b)

		// parabolic u is between b and c: try it
		if (b-u)*(u-c) > 0.0 {
			fu = o.ffcn(u)
			o.NFeval++

			// got a minimum between u and c
			if fu < fc {
				a = b
				b = u
				fa = fb
				fb = fu
				exitCase = 2
				return // exit point

				// got a minimum between between a and u
			} else if fu > fb {
				c = u
				fc = fu
				exitCase = 3
				return // exit point
			}

			// parabolic fit was no use. Use default magnification
			u = c + o.gold*(c-b)
			fu = o.ffcn(u)
			o.NFeval++

			// parabolic fit is between c and its allowed limit
		} else if (c-u)*(u-ulim) > 0.0 {
			fu = o.ffcn(u)
			o.NFeval++
			if fu < fc {
				aux = u + o.gold*(u-c)
				shft3(&b, &c, &u, aux)
				shft3(&fb, &fc, &fu, fu) // TODO: check why [1] calls ffcn(u) again here
			}

			// limit parabolic u to maximum allowed value
		} else if (u-ulim)*(ulim-c) >= 0.0 {
			u = ulim
			fu = o.ffcn(u)
			o.NFeval++

			// reject parabolic u, use default magnification
		} else {
			u = c + o.gold*(c-b)
			fu = o.ffcn(u)
			o.NFeval++
		}

		// eliminate oldest point and continue
		shft3(&a, &b, &c, u)
		shft3(&fa, &fb, &fc, fu)
	}

	// check
	chk.Panic("fail to converge after %d iterations", o.It)
	return
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

func swap(a, b *float64) {
	*a, *b = *b, *a
}

func shft2(a, b *float64, c float64) {
	*a = *b
	*b = c
}

func shft3(a, b, c *float64, d float64) {
	*a = *b
	*b = *c
	*c = d
}

func mov3(a, b, c *float64, d, e, f float64) {
	*a = d
	*b = e
	*c = f
}

func signp(x float64) float64 {
	if x < 0.0 {
		return -1.0
	}
	return +1.0
}
