// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

// Problem defines the data for an ODE problems (e.g. for testing)
type Problem struct {
	Yana YanaF       // analytical solution
	Fcn  Func        // the function f(x,y) = dy/df
	Jac  JacF        // df/dy function
	Dx   float64     // timestep for fixedStep tests
	Xf   float64     // final x
	Y    la.Vector   // initial (current) y vector
	Ndim int         // dimension == len(Y)
	M    *la.Triplet // "mass" matrix
	Ytmp la.Vector   // to use with Yana
}

// Solve solves ODE problem using standard parameters
// NOTE: this solver doesn't change o.Y
func (o *Problem) Solve(method string, fixedStp, numJac bool) (y la.Vector, stat *Stat, out *Output) {

	// current y vector
	y = la.NewVector(o.Ndim)
	y.Apply(1, o.Y)

	// configuration
	conf := NewConfig(method, "")
	if fixedStp {
		conf.SetFixedH(o.Dx, o.Xf)
	}
	conf.SetStepOut(true, nil)

	// allocate solver
	jac := o.Jac
	if numJac {
		jac = nil
	}
	sol := NewSolver(o.Ndim, conf, o.Fcn, jac, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(y, 0.0, o.Xf)

	// set output
	stat = sol.Stat
	out = sol.Out

	// set auxiliary variable
	o.Ytmp = la.NewVector(o.Ndim)
	return
}

// ConvergenceTest runs convergence test
//
//	yExact -- is the exact (reference) y @ xf
func (o *Problem) ConvergenceTest(tst *testing.T, dxmin, dxmax float64, ndx int, yExact la.Vector,
	methods []string, orders, tols []float64) {

	// constants
	dxs := utl.LinSpace(dxmin, dxmax, ndx)
	U := make([]float64, ndx)
	V := make([]float64, ndx)
	lu := make([]float64, ndx)
	lv := make([]float64, ndx)

	// try methods
	for im, method := range methods {

		// run for many dx
		for idx, dx := range dxs {

			// solve problem
			o.Dx = dx
			y, stat, _ := o.Solve(method, true, false)

			// global error
			diff := la.VecMaxDiff(y, yExact)
			U[idx] = float64(stat.Nfeval)
			V[idx] = diff

			// log-log values
			lu[idx] = math.Log10(U[idx])
			lv[idx] = math.Log10(V[idx])
		}

		// calc convergence rate
		_, m := num.LinFit(lu, lv)
		chk.AnaNum(tst, "slope m", tols[im], m, -orders[im], chk.Verbose)
	}
}

// CalcYana computes component idxY of analytical solution @ x, if available
func (o *Problem) CalcYana(idxY int, x float64) float64 {
	if o.Yana == nil {
		chk.Panic("analytical solution is not available\n")
	}
	if len(o.Ytmp) != o.Ndim {
		o.Ytmp = la.NewVector(o.Ndim)
	}
	o.Yana(o.Ytmp, x)
	return o.Ytmp[idxY]
}

// problems database //////////////////////////////////////////////////////////////////////////////

// ProbHwEq11 returns the Hairer-Wanner problem from VII-p2 Eq.(1.1)
func ProbHwEq11() (o *Problem) {

	o = new(Problem)
	λ := -50.0
	o.Dx = 1.875 / 50.0
	o.Xf = 1.5
	o.Y = la.NewVectorSlice([]float64{0.0})
	o.Ndim = len(o.Y)

	o.Yana = func(res []float64, x float64) {
		res[0] = -λ * (math.Sin(x) - λ*math.Cos(x) + λ*math.Exp(λ*x)) / (λ*λ + 1.0)
	}

	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) {
		f[0] = λ*y[0] - λ*math.Cos(x)
	}

	o.Jac = func(dfdy *la.Triplet, dx, x float64, y la.Vector) {
		if dfdy.Max() == 0 {
			dfdy.Init(1, 1, 1)
		}
		dfdy.Start()
		dfdy.Put(0, 0, λ)
	}
	return
}

// ProbVanDerPol returns the Van der Pol' Equation as given in Hairer-Wanner VII-p5 Eq.(1.5)
//
//	eps  -- ε coefficient; use 0 for default value [=1e-6]
//	stationary -- use eps=1 and compute period and amplitude such that
//	              y = [A, 0] is a stationary point
func ProbVanDerPol(eps float64, stationary bool) (o *Problem) {

	o = new(Problem)
	o.Xf = 2.0
	o.Y = la.NewVectorSlice([]float64{2.0, -0.6})
	o.Ndim = len(o.Y)

	if eps < 1e-16 {
		eps = 1e-6
	}

	if stationary {
		eps = 1.0
		T := 6.6632868593231301896996820305
		A := 2.00861986087484313650940188
		o.Y[0] = A
		o.Y[1] = 0.0
		o.Xf = T
	}

	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) {
		f[0] = y[1]
		f[1] = ((1.0-y[0]*y[0])*y[1] - y[0]) / eps
	}

	o.Jac = func(dfdy *la.Triplet, dx, x float64, y la.Vector) {
		if dfdy.Max() == 0 {
			dfdy.Init(2, 2, 4)
		}
		dfdy.Start()
		dfdy.Put(0, 0, 0.0)
		dfdy.Put(0, 1, 1.0)
		dfdy.Put(1, 0, (-2.0*y[0]*y[1]-1.0)/eps)
		dfdy.Put(1, 1, (1.0-y[0]*y[0])/eps)
	}
	return
}

// ProbRobertson returns the Robertson's Equation as given in Hairer-Wanner VII-p3 Eq.(1.4)
func ProbRobertson() (o *Problem) {

	o = new(Problem)
	o.Xf = 0.3
	o.Y = la.NewVectorSlice([]float64{1.0, 0.0, 0.0})
	o.Ndim = len(o.Y)

	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) {
		f[0] = -0.04*y[0] + 1.0e4*y[1]*y[2]
		f[1] = 0.04*y[0] - 1.0e4*y[1]*y[2] - 3.0e7*y[1]*y[1]
		f[2] = 3.0e7 * y[1] * y[1]
	}

	o.Jac = func(dfdy *la.Triplet, dx, x float64, y la.Vector) {
		if dfdy.Max() == 0 {
			dfdy.Init(3, 3, 9)
		}
		dfdy.Start()
		dfdy.Put(0, 0, -0.04)
		dfdy.Put(0, 1, 1.0e4*y[2])
		dfdy.Put(0, 2, 1.0e4*y[1])
		dfdy.Put(1, 0, 0.04)
		dfdy.Put(1, 1, -1.0e4*y[2]-6.0e7*y[1])
		dfdy.Put(1, 2, -1.0e4*y[1])
		dfdy.Put(2, 0, 0.0)
		dfdy.Put(2, 1, 6.0e7*y[1])
		dfdy.Put(2, 2, 0.0)
	}
	return
}

// ProbArenstorf returns the Arenstorf orbit problem
func ProbArenstorf() (o *Problem) {
	o = new(Problem)
	o.Xf = 17.0652165601579625588917206249
	o.Y = la.NewVectorSlice([]float64{
		0.994,
		0.0,
		0.0,
		-2.00158510637908252240537862224,
	})
	o.Ndim = len(o.Y)
	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) {
		amu := 0.012277471
		amup := 1.0 - amu
		f[0] = y[2]
		f[1] = y[3]
		r1 := (y[0]+amu)*(y[0]+amu) + y[1]*y[1]
		r1 = r1 * math.Sqrt(r1)
		r2 := (y[0]-amup)*(y[0]-amup) + y[1]*y[1]
		r2 = r2 * math.Sqrt(r2)
		f[2] = y[0] + 2*y[3] - amup*(y[0]+amu)/r1 - amu*(y[0]-amup)/r2
		f[3] = y[1] - 2*y[2] - amup*y[1]/r1 - amu*y[1]/r2
	}
	return
}

// ProbSimpleNdim2 returns a simple 2-dim problem
func ProbSimpleNdim2() (o *Problem) {
	o = new(Problem)
	o.Yana = func(res []float64, x float64) {
		e2x := math.Exp(2.0 * x)
		res[0] = -0.5*e2x + x*x + 2*x - 0.5
		res[1] = +0.5*e2x + x*x - 0.5
	}
	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) {
		f[0] = +y[0] - y[1] + 2.0
		f[1] = -y[0] + y[1] + 4.0*x
	}
	o.Y = la.NewVectorSlice([]float64{-1.0, 0.0})
	o.Ndim = len(o.Y)
	o.Dx = 0.1
	o.Xf = 1.0
	return
}

// ProbSimpleNdim4a returns a simple 4-dim problem (a)
func ProbSimpleNdim4a() (o *Problem) {
	o = new(Problem)
	o.Yana = func(res []float64, x float64) {
		res[0] = math.Exp(math.Sin(x * x))
		res[1] = math.Exp(5.0 * math.Sin(x*x))
		res[2] = math.Sin(x*x) + 1.0
		res[3] = math.Cos(x * x)
	}
	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) {
		f[0] = 2.0 * x * y[3] * y[0]
		f[1] = 10.0 * x * y[3] * fun.PowP(y[0], 5)
		f[2] = 2.0 * x * y[3]
		f[3] = -2.0 * x * (y[2] - 1)
	}
	o.Y = la.NewVectorSlice([]float64{1, 1, 1, 1})
	o.Ndim = len(o.Y)
	o.Dx = 0.1
	o.Xf = 2.8
	return
}

// ProbSimpleNdim4b returns a simple 4-dim problem (b)
func ProbSimpleNdim4b() (o *Problem) {
	o = new(Problem)
	o.Yana = func(res []float64, x float64) {
		res[0] = math.Exp(math.Sin(x * x))
		res[1] = math.Exp(5.0 * math.Sin(x*x))
		res[2] = math.Sin(x*x) + 1.0
		res[3] = math.Cos(x * x)
	}
	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) {
		f[0] = 2.0 * x * math.Pow(y[1], 1.0/5.0) * y[3]
		f[1] = 10.0 * x * math.Exp(5.0*(y[2]-1.0)) * y[3]
		f[2] = 2.0 * x * y[3]
		f[3] = -2.0 * x * math.Log(y[0])
		if y[0] < 0 || y[1] < 0 {
			io.Pf("x = %v\n", x)
			io.Pf("y = %v\n", y)
			io.Pf("f = %v\n", f)
			chk.Panic("y0 and y1 cannot be negative\n")
		}
	}
	o.Y = la.NewVectorSlice([]float64{1, 1, 1, 1})
	o.Ndim = len(o.Y)
	o.Dx = 0.1
	o.Xf = 2.8
	return
}
