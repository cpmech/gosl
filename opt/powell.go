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
	NumFeval int      // number of calls to Ffcn (function evaluations)
	NumIter  int      // number of iterations from last call to Solve
	Hist     *History // history of optimization data (for debugging)

	// access
	Umat *la.Matrix // matrix whose columns contain the directions u

	// internal
	size int       // problem dimension = len(x)
	ffcn fun.Sv    // scalar function of vector: y = f({x})
	xcpy la.Vector // copy of initial x
	xext la.Vector // auxiliary "extrapolated" point
	uave la.Vector // average direction moved
	tiny float64   // small number for convergence check

	// line solver
	line *num.LineSolver // line solver wrapping Brent's method
}

// NewPowell returns a new multidimensional optimizer using Powell's method (no derivatives required)
//   size -- length(x)
//   ffcn -- scalar function of vector: y = f({x})
func NewPowell(size int, ffcn fun.Sv) (o *Powell) {
	o = new(Powell)
	o.size = size
	o.ffcn = ffcn
	o.MaxIt = 200
	o.Ftol = 3e-8
	o.tiny = 1e-15
	o.line = num.NewLineSolver(size, ffcn, nil)
	o.xcpy = la.NewVector(size)
	o.xext = la.NewVector(size)
	o.uave = la.NewVector(size)
	o.Umat = la.NewMatrix(size, size)
	return
}

// Min solves minimization problem
//
//  Input:
//    x -- [size] initial starting point (will be modified)
//    reuseUmat -- use pre-computed Nmat containing the directions as columns;
//                 otherwise, set diagonal matrix
//
//  Output:
//    fmin -- f(x@min) minimum f({x}) found
//    x -- [given as input] position of minimum f({x})
//
func (o *Powell) Min(x la.Vector, reuseUmat bool) (fmin float64) {

	// set Nmat with unit vectors
	if !reuseUmat {
		o.Umat.SetDiag(1) // set diagonal
	}

	// initializations
	fmin = o.ffcn(x)
	o.NumFeval = 1

	// history
	var λhist float64
	var uhist la.Vector
	if o.History {
		o.Hist = NewHistory(o.MaxIt, fmin, x, o.ffcn)
		uhist = la.NewVector(o.size)
	}

	// save initial x
	o.xcpy.Apply(1, x) // xcpy := x

	// iterations
	for o.NumIter = 0; o.NumIter < o.MaxIt; o.NumIter++ {

		// set iteration values
		fx := fmin  // iteration f({x})
		jdel := 0   // index of largest decrease
		delF := 0.0 // largest function decrease

		// loop over all directions in the set
		for jdir := 0; jdir < o.size; jdir++ {

			// minimize along direction jdir
			u := o.Umat.GetCol(jdir)              // direction
			fold := fmin                          // save fmin
			λhist, fmin = o.line.MinUpdateX(x, u) // x := x @ min
			o.NumFeval += o.line.NumFeval

			// update jdel if jdir gives the largest decrease
			if fold-fmin > delF {
				delF = fold - fmin
				jdel = jdir
			}

			// history
			if o.History {
				uhist.Apply(λhist, u)
				o.Hist.Append(fmin, x, uhist)
			}
		}

		// exit point
		if 2.0*(fx-fmin) <= o.Ftol*(math.Abs(fx)+math.Abs(fmin))+o.tiny {
			return
		}

		// compute extrapolated point, compute average direction, and save starting point
		for i := 0; i < o.size; i++ {
			o.xext[i] = 2.0*x[i] - o.xcpy[i] // xext := 2⋅x - x0  extrapolated point
			o.uave[i] = x[i] - o.xcpy[i]     // uave := x - x0    average direction moved
			o.xcpy[i] = x[i]                 // xcpy := x         save the old starting point
		}

		// function value at extrapolated point
		fext := o.ffcn(o.xext)
		o.NumFeval++

		// move to the minimum of the new direction, and save the new direction
		if fext < fx {
			if 2*(fx-2*fmin+fext)*fun.Pow2(fx-fmin-delF) < fun.Pow2(fx-fext)*delF {

				// minimize along average direction
				λhist, fmin = o.line.MinUpdateX(x, o.uave)
				o.NumFeval += o.line.NumFeval

				// save average direction
				for i := 0; i < o.size; i++ {
					o.Umat.Set(i, jdel, o.Umat.Get(i, o.size-1))
					o.Umat.Set(i, o.size-1, o.uave[i])
				}

				// history
				if o.History {
					uhist.Apply(λhist, o.uave)
					o.Hist.Append(fmin, x, uhist)
				}

			} // else: keep the old set of directions for the next basic procedure, because either
			// (i) the decrease along the average direction was not primarily due to any single
			// direction's decrease, or
			// (ii) there is a substantial second derivative along the average
			// direction and we seem to be near to the bottom of its minimum.

		} // else: keep the old set of directions  because the average direction x-xcpy is good
	}

	// did not converge
	chk.Panic("fail to converge after %d iterations\n", o.NumIter)
	return
}
