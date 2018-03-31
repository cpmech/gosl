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
	"github.com/cpmech/gosl/utl"
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
	Nfeval int      // number of calls to Ffcn (function evaluations)
	Niter  int      // number of iterations from last call to Solve
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
//    fmin -- f(x@min) minimum f({x}) found
//    xmin -- x at position of minimum f({x}) [internal slice; no copies created]
//
func (o *Powell) Min(x0 la.Vector, reuseNmat bool) (fmin float64, xmin []float64) {

	// set Nmat with unit vectors
	if !reuseNmat {
		o.Nmat.SetDiag(1) // set diagonal
	}

	// initializations
	o.x.Apply(1, x0)   // x := x0
	fmin = o.ffcn(o.x) // fmin := f({x0})
	o.Nfeval = 1
	o.Niter = 0

	// history
	var λhist float64
	var nhist la.Vector
	if o.History {
		o.Hist = NewHistory(o.MaxIt, fmin, o.x, o.ffcn)
		nhist = la.NewVector(o.size)
	}

	// iterations
	for it := 0; it < o.MaxIt; it++ {

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
			//o.Nfeval += o.LS.Brent.NFeval

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
			o.Niter++
		}

		// exit point
		if 2.0*(fx-fmin) <= o.Ftol*(math.Abs(fx)+math.Abs(fmin))+o.tiny {
			xmin = o.x
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
		o.Nfeval++

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
	chk.Panic("fail to converge\n")
	return
}

// MinVersion2 is the original version
func (o *Powell) MinVersion2(x0 la.Vector, reuseNmat bool) (fmin float64, xmin []float64) {

	ftol := 3.0e-8

	n := len(x0)
	ximat := utl.Alloc(n, n)
	for i := 0; i < n; i++ {
		ximat[i][i] = 1.0
	}

	ITMAX := 200    // Maximum allowed iterations.
	TINY := 1.0e-25 // A small number.
	var fptt float64
	p := x0.GetCopy()
	pt := make([]float64, n)
	ptt := make([]float64, n)
	xi := make([]float64, n)
	fret := o.ffcn(p)
	for j := 0; j < n; j++ {
		pt[j] = p[j] // Save the initial point.
	}
	for iter := 0; ; iter++ {
		fp := fret
		ibig := 0
		del := 0.0               // Will be the biggest function decrease.
		for i := 0; i < n; i++ { // In each iteration, loop over all directions in the set
			for j := 0; j < n; j++ {
				xi[j] = ximat[j][i] // Copy the direction
			}
			fptt = fret
			_, fret = o.LS.MinUpdateX(p, xi) // minimize along it,
			if fptt-fret > del {             // and record it if it is the largest so far.
				del = fptt - fret
				ibig = i + 1
			}
		}
		// Here comes the termination criterion:
		if 2.0*(fp-fret) <= ftol*(math.Abs(fp)+math.Abs(fret))+TINY {
			fmin = fret
			xmin = p
			return
		}
		if iter == ITMAX {
			chk.Panic("powell exceeding maximum iterations")
		}
		for j := 0; j < n; j++ { // Construct the extrapolated point and the average direction moved
			ptt[j] = 2.0*p[j] - pt[j]
			xi[j] = p[j] - pt[j] // old starting point.
			pt[j] = p[j]
		}
		fptt = o.ffcn(ptt) // Function value at extrapoif (fptt < fp) {
		t := 2.0*(fp-2.0*fret+fptt)*math.Pow(fp-fret-del, 2) - del*math.Pow(fp-fptt, 2)
		if t < 0.0 {
			_, fret = o.LS.MinUpdateX(p, xi) // Move to the minimum of the new direction, and save the new direction
			for j := 0; j < n; j++ {
				ximat[j][ibig-1] = ximat[j][n-1]
				ximat[j][n-1] = xi[j]
			}
		}
	}
}
