// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// enums
var (

	// LinearInterpKind defines the linear interpolator kind
	LinearInterpKind = io.NewEnum("Linear", "fun.interp", "L", "Linear Interpolator")

	// PolyInterpKind defines the polynomial interpolator kind
	PolyInterpKind = io.NewEnum("Polynomial", "fun.interp", "L", "Polynomial Interpolator")
)

// Interpolator implements numeric interpolators
type Interpolator struct {

	// configuration data
	DisableHunt bool // do not use hunt code at all

	// output data
	Dy float64 // error estimate

	// input data
	itype io.Enum   // type of interpolator
	xx    []float64 // x-data values
	yy    []float64 // y-data values

	// derived data
	m       int  // number of points of interpolating formula; e.g. 2 for segments, 3 for 2nd order polynomials
	n       int  // length of xx
	jHunt   int  // temporary j to decide on using hunt
	djHunt  int  // increent of j to decide on using hunt function or locate
	useHunt bool // use hunt code instead of locate
	ascnd   bool // ascending order of x-values

	// implementation
	interp func(j int, x float64) float64
}

// NewInterpolator creates new interpolator of type=Type for data point sets xx and yy (with same lengths)
//   Input:
//     Type -- type of interpolator
//     p    -- order of interpolator
//     xx   -- x-data
//     yy   -- y-data
func NewInterpolator(Type io.Enum, p int, xx, yy []float64) (o *Interpolator, err error) {
	o = new(Interpolator)
	o.itype = Type
	switch Type {
	case LinearInterpKind:
		o.m = 2
		o.interp = o.linInterp
	case PolyInterpKind:
		o.m = p + 1
		o.interp = o.polyInterp
	default:
		return nil, chk.Err("cannot find interpolator type == %q\n", Type)
	}
	err = o.Reset(xx, yy)
	return
}

// Reset re-assigns xx and yy data sets
func (o *Interpolator) Reset(xx, yy []float64) (err error) {
	if len(xx) != len(yy) {
		return chk.Err("lengths of data sets must be the same. %d != %d\n", len(xx), len(yy))
	}
	if len(xx) < 2 {
		return chk.Err("length of data sets must be at least 2. %d is invalid\n", len(xx))
	}
	if len(xx) <= o.m {
		return chk.Err("length of data sets must be at smaller than %d when using %q interpolator\n", o.m, o.itype)
	}
	o.xx = xx
	o.yy = yy
	o.n = len(o.xx)
	o.djHunt = imin(1, int(math.Pow(float64(o.n), 0.25)))
	o.useHunt = false
	o.ascnd = o.xx[o.n-1] >= o.xx[0]
	return
}

// P computes P(x); i.e. performs the interpolation
func (o *Interpolator) P(x float64) float64 {
	var jlo int
	if o.useHunt && !o.DisableHunt {
		jlo = o.hunt(x)
	} else {
		jlo = o.locate(x)
	}
	return o.interp(jlo, x)
}

// locate returns a value j such that x is (insofar as possible) centered in the subrange
// xx[j..j+mm-1], where xx is the stored pointer. The values in xx must be monotonic, either
// increasing or decreasing. The returned value is not less than 0, nor greater than n-1.
func (o *Interpolator) locate(x float64) int {

	// bisection
	jl := 0         // initialize lower
	ju := o.n - 1   // and upper limits.
	for ju-jl > 1 { // if not done yet done
		jm := (ju + jl) >> 1 // compute a midpoint
		if x >= o.xx[jm] == o.ascnd {
			jl = jm // replace the lower limit
		} else {
			ju = jm // replace the upper limit
		}
	}

	// set hunt flag
	if iabs(jl-o.jHunt) > o.djHunt {
		o.useHunt = false // too large, use locate next time
	} else {
		o.useHunt = true // ok, use hunt next time
	}
	o.jHunt = jl

	// results
	return imax(0, imin(o.n-o.m, jl-((o.m-2)>>1)))
}

// hunt returns a value j such that x is (insofar as possible) centered in the subrange
// xx[j..j+mm-1], where xx is the stored pointer. The values in xx must be monotonic, either
// increasing or decreasing. The returned value is not less than 0, nor greater than n-1.
func (o *Interpolator) hunt(x float64) int {

	// hunting
	jl := o.jHunt
	inc := 1
	var ju, jm int
	if jl < 0 || jl > o.n-1 { // input guess not useful. skip hunting
		jl = 0
		ju = o.n - 1
	} else {
		if x >= o.xx[jl] == o.ascnd { // hunt up
			for {
				ju = jl + inc
				if ju >= o.n-1 {
					ju = o.n - 1
					break // off end of table.
				} else if x < o.xx[ju] == o.ascnd {
					break // found bracket.
				} else { // not done, so double the increment and try again.
					jl = ju
					inc += inc
				}
			}
		} else { // hunt down
			ju = jl
			for {
				jl = jl - inc
				if jl <= 0 {
					jl = 0
					break // off end of table.
				} else if x >= o.xx[jl] == o.ascnd {
					break // found bracket.
				} else { // not done, so double the increment and try again.
					ju = jl
					inc += inc
				}
			}
		}
	}

	// hunt is done, so begin the final bisection phase:
	for ju-jl > 1 {
		jm = (ju + jl) >> 1
		if x >= o.xx[jm] == o.ascnd {
			jl = jm
		} else {
			ju = jm
		}
	}

	// set hunt flag
	if iabs(jl-o.jHunt) > o.djHunt {
		o.useHunt = false
	} else {
		o.useHunt = true
	}
	o.jHunt = jl

	// results
	return imax(0, imin(o.n-o.m, jl-((o.m-2)>>1)))
}

// linInterp implements linear interpolator
func (o *Interpolator) linInterp(j int, x float64) float64 {
	if o.xx[j] == o.xx[j+1] { // table is defective, but we can recover.
		return o.yy[j]
	}
	return o.yy[j] + (o.yy[j+1]-o.yy[j])*(x-o.xx[j])/(o.xx[j+1]-o.xx[j])
}

// polyInterp performs a polynomial interpolation. This routine returns an interpolated value y, and
// stores an error estimate dy. The returned value is obtained by m-point polynomial interpolation
// on the subrange xx[jl..jl+m-1].
func (o *Interpolator) polyInterp(jl int, x float64) (y float64) {

	// allocate variables
	xa := o.xx[jl:]
	ya := o.yy[jl:]
	dif := math.Abs(x - xa[0])
	c := make([]float64, o.m)
	d := make([]float64, o.m)

	// find the index ns of the closest table entry,
	var ns int
	var dift float64
	for i := 0; i < o.m; i++ {
		dift = math.Abs(x - xa[i])
		if dift < dif {
			ns = i
			dif = dift
		}
		c[i] = ya[i] // initialize the tableau of c's and d's.
		d[i] = ya[i]
	}

	// initial approximation to y.
	y = ya[ns]
	ns--

	// perform interpolation
	var ho, hp, w, den float64
	for m := 1; m < o.m; m++ { // for each column of the tableau,
		for i := 0; i < o.m-m; i++ { // loop over the current c's and d's and update them.
			ho = xa[i] - x
			hp = xa[i+m] - x
			w = c[i+1] - d[i]
			den = ho - hp
			if den == 0.0 {
				chk.Panic("polyInterp failed because two input x points are identical (within roundoff)")
			}
			den = w / den
			d[i] = hp * den
			c[i] = ho * den
		}
		if 2*(ns+1) < (o.m - m) {
			o.Dy = c[ns+1]
		} else {
			o.Dy = d[ns]
			ns--
		}
		y += o.Dy
	}
	return
}
