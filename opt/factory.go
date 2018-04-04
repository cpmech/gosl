// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
)

// FactoryType defines a structure to implement a factory of Objective functions to be minimized
type FactoryType struct{}

// Factory holds Objective functions to be minimized
var Factory = FactoryType{}

// unconstrained multi-dimensional problems ////////////////////////////////////////////////////////

// SimpleParaboloid returns a simple optimization problem using a paraboloid as the objective function
func (o FactoryType) SimpleParaboloid() (p *Problem) {

	// new problem
	p = new(Problem)
	p.Ndim = 2

	// objective function f({x})
	p.Ffcn = func(x la.Vector) float64 {
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	// gradient function df/d{x}|(x)
	p.Gfcn = func(g, x la.Vector) {
		g[0] = 2.0 * x[0]
		g[1] = 2.0 * x[1]
	}

	// Hessian function d²f/d{x}d{x}!(x)
	p.Hfcn = func(h *la.Matrix, x la.Vector) {
		h.Set(0, 0, 2.0)
		h.Set(0, 1, 0.0)
		h.Set(1, 0, 0.0)
		h.Set(1, 1, 2.0)
	}

	// known solution
	p.Fref = -0.5
	p.Xref = la.NewVectorSlice([]float64{0, 0})
	return
}

// SimpleQuadratic2d returns a simple problem with a quadratic function such that f(x) = xᵀ A x (2D)
func (o FactoryType) SimpleQuadratic2d() (p *Problem) {
	return NewQuadraticProblem([][]float64{
		{3, 1},
		{1, 2},
	})
}

// SimpleQuadratic3d returns a simple problem with a quadratic function such that f(x) = xᵀ A x (3D)
func (o FactoryType) SimpleQuadratic3d() (p *Problem) {
	return NewQuadraticProblem([][]float64{
		{5, 3, 1},
		{3, 4, 2},
		{1, 2, 3},
	})
}

// Rosenbrock2d returns the classical Rosenbrock2d function
//
//   See https://en.wikipedia.org/wiki/Rosenbrock_function
//
//   Input:
//     a -- parameter a, a=0 ⇒ function is symmetric and minimum is at origin
//     b -- parameter b
//
//   NOTE: the commonly used values are a=1 and b=100
//
func (o FactoryType) Rosenbrock2d(a, b float64) (p *Problem) {

	// new problem
	p = new(Problem)
	p.Ndim = 2

	// objective function f({x})
	p.Ffcn = func(x la.Vector) float64 {
		return fun.Pow2(a-x[0]) + b*fun.Pow2(x[1]-x[0]*x[0])
	}

	// gradient function df/d{x}|(x)
	p.Gfcn = func(g, x la.Vector) {
		g[0] = -4.0*b*x[0]*(x[1]-x[0]*x[0]) - 2.0*(a-x[0])
		g[1] = 2.0 * b * (x[1] - x[0]*x[0])
	}

	// known solution
	p.Fref = 0.0
	p.Xref = la.NewVectorSlice([]float64{a, a * a})
	return
}

// RosenbrockMulti returns the multi-variate version of the Rosenbrock function
//
//   See https://en.wikipedia.org/wiki/Rosenbrock_function
//   See https://docs.scipy.org/doc/scipy-0.14.0/reference/tutorial/optimize.html#unconstrained-minimization-of-multivariate-scalar-functions-minimize
//
//   Input:
//     N -- dimension == ndim
//
func (o FactoryType) RosenbrockMulti(N int) (p *Problem) {

	// check
	if N < 3 {
		chk.Panic("RosenbrockMulti requires N ≥ 3\n")
	}

	// new problem
	p = new(Problem)
	p.Ndim = N

	// objective function f({x})
	p.Ffcn = func(x la.Vector) float64 {
		sum := 0.0
		for i := 1; i < len(x); i++ {
			sum += 100.0*fun.Pow2(x[i]-x[i-1]*x[i-1]) + fun.Pow2(1.0-x[i-1])
		}
		return sum
	}

	// gradient function df/d{x}|(x)
	p.Gfcn = func(g, x la.Vector) {
		n := len(x)
		for j := 1; j < n-1; j++ {
			g[j] = 200.0*(x[j]-x[j-1]*x[j-1]) - 400.0*x[j]*(x[j+1]-x[j]*x[j]) - 2.0*(1.0-x[j])
		}
		g[0] = -400.0*x[0]*(x[1]-x[0]*x[0]) - 2.0*(1.0-x[0])
		g[n-1] = 200.0 * (x[n-1] - x[n-2]*x[n-2])
	}

	// known solution
	p.Fref = 0.0
	p.Xref = la.NewVector(N)
	p.Xref.Fill(1.0)
	return
}
