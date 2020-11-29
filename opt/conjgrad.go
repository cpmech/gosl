// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

// ConjGrad implements the multidimensional minimization by the Fletcher-Reeves-Polak-Ribiere method.
//
//   NOTE: Check Convergence to see how to set convergence parameters,
//         max iteration number, or to enable and access history of iterations
//
//   REFERENCES:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes:
//       The Art of Scientific Computing. Third Edition. Cambridge University Press. 1235p.
//
type ConjGrad struct {

	// merge properties
	Convergence // auxiliary object to check convergence

	// configuration
	UseBrent    bool // use Brent method insted of LineSearch (Wolfe conditions)
	UseFRmethod bool // use Fletcher-Reeves method instead of Polak-Ribiere
	CheckJfcn   bool // check Jacobian function at all points during minimization

	// internal
	u    la.Vector // direction vector for line minimization
	g    la.Vector // conjugate direction vector
	h    la.Vector // conjugate direction vector
	tmp  la.Vector // auxiliary vector
	zero float64   // constant to prevent division by zero

	// line solver
	lines *LineSearch     // line search
	lineb *num.LineSolver // line solver wrapping Brent's method
}

// add optimizer to database
func init() {
	nlsMakersDB["conjgrad"] = func(prob *Problem) NonLinSolver { return NewConjGrad(prob) }
}

// NewConjGrad returns a new multidimensional optimizer using ConjGrad's method (no derivatives required)
func NewConjGrad(prob *Problem) (o *ConjGrad) {
	o = new(ConjGrad)
	o.InitConvergence(prob.Ffcn, prob.Gfcn)
	o.lines = NewLineSearch(prob.Ndim, o.Ffcn, o.Gfcn)
	o.lineb = num.NewLineSolver(prob.Ndim, o.Ffcn, o.Gfcn)
	o.u = la.NewVector(prob.Ndim)
	o.g = la.NewVector(prob.Ndim)
	o.h = la.NewVector(prob.Ndim)
	o.tmp = la.NewVector(prob.Ndim)
	o.zero = 1e-18
	return
}

// Min solves minimization problem
//
//  Input:
//    x -- [ndim] initial starting point (will be modified)
//    params -- [may be nil] optional parameters. e.g. "alpha", "maxit". Example:
//                 params := utl.NewParams(
//                     &utl.P{N: "brent", V: 1},
//                     &utl.P{N: "maxit", V: 1000},
//                     &utl.P{N: "maxitls", V: 20},
//                     &utl.P{N: "maxitzoom", V: 20},
//                     &utl.P{N: "ftol", V: 1e-2},
//                     &utl.P{N: "gtol", V: 1e-2},
//                     &utl.P{N: "hist", V: 1},
//                     &utl.P{N: "verb", V: 1},
//                 )
//
//  Output:
//    fmin -- f(x@min) minimum f({x}) found
//    x -- [modify input] position of minimum f({x})
//
func (o *ConjGrad) Min(x la.Vector, params utl.Params) (fmin float64) {

	// set parameters
	o.Convergence.SetParams(params)
	o.UseBrent = params.GetBoolOrDefault("brent", o.UseBrent)
	o.lines.SetParams(params)

	// line search function and counters
	linesearch := o.lines.Wolfe
	if o.UseBrent {
		linesearch = func(x, u la.Vector, dum1 bool, dum2 float64) (λ, fmin float64) { return o.lineb.MinUpdateX(x, u) }
	}

	// initializations
	ndim := len(x)
	fx := o.Ffcn(x) // fx := f(x)
	o.Gfcn(o.u, x)  // u := dy/dx
	for j := 0; j < ndim; j++ {
		o.u[j] = -o.u[j] // u := -dy/dx
		o.g[j] = o.u[j]  // g := -dy/dx
		o.h[j] = o.u[j]  // h := g
	}
	fmin = fx

	// history
	var λhist float64
	if o.UseHist {
		o.InitHist(x)
	}

	// auxiliary
	var nume, deno, γ float64

	// estimate old f(x)
	fold := fx + o.u.Norm()/2.0 // TODO: find reference to this

	// iterations
	for o.NumIter = 0; o.NumIter < o.MaxIt; o.NumIter++ {

		// exit point # 1: old gradient is exactly zero
		deno = la.VecDot(o.g, o.g)
		if math.Abs(deno) < o.zero {
			return
		}

		// line minimization
		λhist, fmin = linesearch(x, o.u, true, fold) // x := x @ min

		// update fold
		fold = fx

		// history
		if o.UseHist {
			o.uhist.Apply(λhist, o.u)
			o.Hist.Append(fmin, x, o.uhist)
		}

		// exit point # 2: converged on f
		if o.Fconvergence(fx, fmin) {
			return
		}

		// update fx and gradient dy/dx
		fx = fmin
		o.Gfcn(o.u, x) // u := dy/dx

		// check Jacobian @ x
		if o.CheckJfcn {
			o.checkJacobian(x)
		}

		// exit point # 3: converged on dy/dx (new)
		if o.Gconvergence(fx, x, o.u) {
			return
		}

		// compute scaling factor, noting that, now:
		//   u = -gNew
		//   g =  gOld
		if o.UseFRmethod {
			nume = la.VecDot(o.u, o.u) // nume := gNew ⋅ gNew  [Equation 10.8.5 page 517 of Ref 1]
		} else {
			la.VecAdd(o.tmp, 1, o.u, 1, o.g) // tmp := u + g = -gNew + gOld
			nume = la.VecDot(o.tmp, o.u)     // nume := (gOld - gNew) ⋅ (-gNew) = (gNew - gOld) ⋅ gNew  [Equation 10.8.7 page 517 of Ref 1]
			nume = utl.Max(nume, 0)          // avoid negative values
		}

		// update directions
		γ = nume / deno
		for j := 0; j < ndim; j++ {
			o.g[j] = -o.u[j]           // g := -dy/dx = gNew
			o.u[j] = o.g[j] + γ*o.h[j] // u := gNew + γ⋅hOld = hNew
			o.h[j] = o.u[j]            // h := hNew
		}
	}

	// did not converge
	chk.Panic("fail to converge after %d iterations\n", o.NumIter)
	return
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// checkJacobian checks Jacobian at intermediate point x
func (o *ConjGrad) checkJacobian(x la.Vector) {
	ndim := len(x)
	tolJ := 1e-12
	for k := 0; k < ndim; k++ {
		dfdxk := num.DerivCen5(x[k], 1e-3, func(xk float64) float64 {
			copy(o.tmp, x)
			o.tmp[k] = xk
			o.NumFeval++
			return o.Ffcn(o.tmp)
		})
		diff := math.Abs(o.u[k] - dfdxk)
		if diff > tolJ {
			chk.Panic("Jacobian function is incorrect. diff = %v\n", diff)
		}
	}
}
