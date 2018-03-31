// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
)

// Powell implements the multidimensional minimization by Powell's method (no derivatives required)
//
//   REFERENCES:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes:
//       The Art of Scientific Computing. Third Edition. Cambridge University Press. 1235p.
//
type Powell struct {

	// configuration
	MaxIt   int     // max iterations
	Ftol    float64 // tolerance on f({x})
	Verbose bool    // show messages
	History bool    // save history

	// statistics and History (for debugging)
	NFeval int      // number of calls to Ffcn (function evaluations)
	It     int      // number of iterations from last call to Solve
	Hist   *History // history of optimization data (for debugging)

	// internal
	size int       // problem dimension
	ffcn fun.Sv    // scalar function of vector: y = f({x})
	x    la.Vector // auxiliary "current" point
	xe   la.Vector // auxiliary "extrapolated" point
	nAve la.Vector // average direction moved
	tiny float64   // small number for convergence check

	// access
	LS   *num.LineSolver // line solver wrapping Brent's method
	Nmat *la.Matrix      // matrix whose columns contain the directions n
}

// NewPowell returns a new multidimensional optimizer using Powell's method (no derivatives required)
//   size -- length(x)
//   ffcn -- scalar function of vector: y = f({x})
func NewPowell(size int, ffcn fun.Sv) (o *Powell) {
	o = new(Powell)
	o.size = size
	o.ffcn = ffcn
	o.MaxIt = 1000
	o.Ftol = 1e-8
	o.tiny = 1e-25
	o.LS = num.NewLineSolver(size, ffcn, nil)
	o.x = la.NewVector(size)
	o.xe = la.NewVector(size)
	o.nAve = la.NewVector(size)
	o.Nmat = la.NewMatrix(size, size)
	return
}

// Min solves minimization problem
//
//  Input:
//    x0 -- [size] initial starting point (will be modified)
//    reuseNmat -- use pre-computed Nmat containing the directions as columns
//
//  Output:
//    o.x -- will hold the point corresponding to the just found fmin
//
func (o *Powell) Min(x0 la.Vector, reuseNmat bool) (fmin float64) {

	// set Nmat with unit vectors
	if !reuseNmat {
		o.Nmat.SetDiag(1) // set diagonal
	}

	// initializations
	o.x.Apply(1, x0)   // x := x0
	fmin = o.ffcn(o.x) // fmin := f({x0})

	// history
	var λhist float64
	var nhist la.Vector
	if o.History {
		o.Hist = NewHistory(o.MaxIt, fmin, o.x, o.ffcn)
		nhist = la.NewVector(o.size)
	}

	// iterations
	for o.It = 0; o.It < o.MaxIt; o.It++ {

		// set iteration values
		fx := fmin  // iteration f({x})
		delJ := 0   // index of largest decrease
		delF := 0.0 // largest function decrease

		// loop over all directions in the set
		for jDir := 0; jDir < o.size; jDir++ {

			// minimize along direction jDir
			n := o.Nmat.GetCol(jDir)              // direction
			fold := fmin                          // save fmin
			λhist, fmin = o.LS.MinUpdateX(o.x, n) // x := x @ min

			// record direction if it corresponds to the largest decrease so far
			if fold-fmin > delF {
				delF = fold - fmin
				delJ = jDir + 1
			}

			// history
			if o.History {
				nhist.Apply(λhist, n)
				o.Hist.Append(fmin, o.x, nhist)
			}
		}

		// exit point
		if 2.0*(fx-fmin) <= o.Ftol*(math.Abs(fx)+math.Abs(fmin))+o.tiny {
			return
		}

		// update
		for i := 0; i < o.size; i++ {
			o.xe[i] = 2.0*o.x[i] - x0[i] // xe := 2⋅x - x0  extrapolated point
			o.nAve[i] = o.x[i] - x0[i]   // nAve := x - x0  average direction moved
			x0[i] = o.x[i]               // save the old starting point
		}

		// function value at extrapolated point
		fe := o.ffcn(o.xe)

		// move to the minimum of the new direction, and save the new direction
		if fe < fx {
			t := 2.0*(fx-2.0*fmin+fe)*math.Pow(fx-fmin-delF, 2) - delF*math.Pow(fx-fe, 2)
			if t < 0.0 {
				fmin = o.LS.Min(o.x, o.nAve)
				for i := 0; i < o.size; i++ {
					o.Nmat.Set(i, delJ-1, o.Nmat.Get(i, o.size-1))
					o.Nmat.Set(i, o.size-1, o.nAve[i])
				}
			}
		}
	}

	// did not converge
	chk.Panic("fail to converge after %d iterations", o.It)
	return
}
