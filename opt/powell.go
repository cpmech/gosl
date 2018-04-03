// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
)

// Powell implements the multidimensional minimization by Powell's method (no derivatives required)
//
//   NOTE: Check Convergence to see how to set convergence parameters,
//         max iteration number, or to enable and access history of iterations
//
//   REFERENCES:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes:
//       The Art of Scientific Computing. Third Edition. Cambridge University Press. 1235p.
//
type Powell struct {

	// merge properties
	Convergence // auxiliary object to check convergence

	// access
	Umat *la.Matrix // matrix whose columns contain the directions u

	// internal
	line *num.LineSolver // line solver wrapping Brent's method
	xcpy la.Vector       // copy of initial x
	xext la.Vector       // auxiliary "extrapolated" point
	uave la.Vector       // average direction moved
}

// NewPowell returns a new multidimensional optimizer using Powell's method (no derivatives required)
//   ndim -- length(x)
//   Ffcn -- objective function: y = f({x})
func NewPowell(ndim int, Ffcn fun.Sv) (o *Powell) {
	o = new(Powell)
	o.InitConvergence(Ffcn, nil)
	o.line = num.NewLineSolver(ndim, o.Ffcn, nil)
	o.xcpy = la.NewVector(ndim)
	o.xext = la.NewVector(ndim)
	o.uave = la.NewVector(ndim)
	o.Umat = la.NewMatrix(ndim, ndim)
	return
}

// Min solves minimization problem
//
//  Input:
//    x -- [ndim] initial starting point (will be modified)
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
	o.NumFeval = 0
	ndim := len(x)
	fmin = o.Ffcn(x)

	// history
	var λhist float64
	if o.UseHist {
		o.InitHist(x)
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
		for jdir := 0; jdir < ndim; jdir++ {

			// minimize along direction jdir
			u := o.Umat.GetCol(jdir)              // direction
			fold := fmin                          // save fmin
			λhist, fmin = o.line.MinUpdateX(x, u) // x := x @ min

			// update jdel if jdir gives the largest decrease
			if fold-fmin > delF {
				delF = fold - fmin
				jdel = jdir
			}

			// history
			if o.UseHist {
				o.uhist.Apply(λhist, u)
				o.Hist.Append(fmin, x, o.uhist)
			}
		}

		// exit point
		if o.Fconvergence(fx, fmin) {
			return
		}

		// compute extrapolated point, compute average direction, and save starting point
		for i := 0; i < ndim; i++ {
			o.xext[i] = 2.0*x[i] - o.xcpy[i] // xext := 2⋅x - x0  extrapolated point
			o.uave[i] = x[i] - o.xcpy[i]     // uave := x - x0    average direction moved
			o.xcpy[i] = x[i]                 // xcpy := x         save the old starting point
		}

		// function value at extrapolated point
		fext := o.Ffcn(o.xext)

		// move to the minimum of the new direction, and save the new direction
		if fext < fx {
			if 2*(fx-2*fmin+fext)*fun.Pow2(fx-fmin-delF) < fun.Pow2(fx-fext)*delF {

				// minimize along average direction
				λhist, fmin = o.line.MinUpdateX(x, o.uave)

				// save average direction
				for i := 0; i < ndim; i++ {
					o.Umat.Set(i, jdel, o.Umat.Get(i, ndim-1))
					o.Umat.Set(i, ndim-1, o.uave[i])
				}

				// history
				if o.UseHist {
					o.uhist.Apply(λhist, o.uave)
					o.Hist.Append(fmin, x, o.uhist)
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
