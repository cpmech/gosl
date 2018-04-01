// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
)

// facObjectivesT defines a structure to implement a factory of Objective functions to be minimized
type facObjectivesT struct{}

// FactoryObjectives holds Objective functions to be minimized
var FactoryObjectives = facObjectivesT{}

// unconstrained multi-dimensional problems ////////////////////////////////////////////////////////

// Rosenbrock returns the classical Rosenbrock function
//
//   See https://en.wikipedia.org/wiki/Rosenbrock_function
//
//   Input:
//     a -- parameter a, a=0 ⇒ function is symmetric and minimum is at origin
//     b -- parameter b
//
//   NOTE: the commonly used values are a=1 and b=100
//
//   Output:
//     ffcn -- objective function f({x}) with len(x) = ndim
//     Jfcn -- Jacobian function == derivative df/dx
//     ndim -- ndim = len(x)
//     hasSol -- has known solution
//     xmin -- known solution
//     fmin -- known solution
//
func (o facObjectivesT) Rosenbrock(a, b float64) (ffcn fun.Sv, Jfcn fun.Vv, ndim int, hasSol bool, xmin la.Vector, fmin float64) {
	ffcn = func(x la.Vector) float64 {
		return fun.Pow2(a-x[0]) + b*fun.Pow2(x[1]-x[0]*x[0])
	}
	Jfcn = func(g, x la.Vector) {
		g[0] = -4.0*b*x[0]*(x[1]-x[0]*x[0]) - 2.0*(a-x[0])
		g[1] = 2.0 * b * (x[1] - x[0]*x[0])
	}
	ndim = 2
	hasSol = true
	xmin = la.NewVectorSlice([]float64{a, a * a})
	fmin = 0.0
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
//   Output:
//     ffcn -- objective function f({x}) with len(x) = ndim
//     ndim -- ndim = len(x)
//     hasSol -- has known solution
//     xmin -- known solution
//     fmin -- known solution
//
func (o facObjectivesT) RosenbrockMulti(N int) (ffcn fun.Sv, Jfcn fun.Vv, ndim int, hasSol bool, xmin la.Vector, fmin float64) {
	if N < 3 {
		chk.Panic("RosenbrockMulti requires N ≥ 3\n")
	}
	ffcn = func(x la.Vector) float64 {
		sum := 0.0
		for i := 1; i < len(x); i++ {
			sum += 100.0*fun.Pow2(x[i]-x[i-1]*x[i-1]) + fun.Pow2(1.0-x[i-1])
		}
		return sum
	}
	Jfcn = func(g, x la.Vector) {
		n := len(x)
		for j := 1; j < n-1; j++ {
			g[j] = 200.0*(x[j]-x[j-1]*x[j-1]) - 400.0*x[j]*(x[j+1]-x[j]*x[j]) - 2.0*(1.0-x[j])
		}
		g[0] = -400.0*x[0]*(x[1]-x[0]*x[0]) - 2.0*(1.0-x[0])
		g[n-1] = 200.0 * (x[n-1] - x[n-2]*x[n-2])
	}
	ndim = N
	hasSol = true
	xmin = la.NewVector(N)
	xmin.Fill(1.0)
	fmin = 0.0
	return
}
