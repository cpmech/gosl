// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
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
}

// Solve solves ODE problem using standard parameters
// NOTE: this solver doesn't change o.Y
func (o *Problem) Solve(method io.Enum, fixedStp, numJac bool) (stat *Stat, out *Output, err error) {

	// current y vector
	y := la.NewVector(o.Ndim)
	y.Apply(1, o.Y)

	// configuration
	conf, err := NewConfig(method, "", nil)
	conf.SaveXY = true
	if fixedStp {
		conf.FixedStp = o.Dx
	}

	// allocate solver
	jac := o.Jac
	if numJac {
		jac = nil
	}
	sol, err := NewSolver(conf, o.Ndim, o.Fcn, jac, nil, nil)
	if err != nil {
		return
	}
	defer sol.Free()

	// solve ODE
	err = sol.Solve(y, 0.0, o.Xf)

	// results
	stat = sol.Stat
	out = sol.Out
	return
}

// Plot plots Y[i] versus x series
func (o *Problem) Plot(label string, out *Output, npts int, withAna bool, argsAna, argsNum *plt.A) {
	if argsAna == nil {
		argsAna = &plt.A{C: "grey", Ls: "-", Lw: 5, L: "ana", NoClip: true}
	}
	if argsNum == nil {
		argsNum = &plt.A{C: "r", M: ".", Ls: "-", L: label, NoClip: true}
	}
	argsNum.L = label
	if withAna {
		X := utl.LinSpace(0, o.Xf, npts)
		Y := utl.GetMapped(X, func(x float64) float64 { return o.Yana(x) })
		plt.Plot(X, Y, argsAna)
	}
	n := out.IdxSave
	Xn := out.Xvalues[:n]
	Yn := out.ExtractTimeSeries(0)
	plt.Plot(Xn, Yn, argsNum)
}

// problems database //////////////////////////////////////////////////////////////////////////////

// ProbHwEq11 returns the Hairer-Wanner problem from VII-p2 Eq.(1.1)
func ProbHwEq11() (o *Problem) {

	o = new(Problem)
	λ := -50.0
	o.Dx = 1.875 / 50.0
	o.Xf = 1.5
	o.Y = la.Vector([]float64{0.0})
	o.Ndim = len(o.Y)

	o.Yana = func(x float64) float64 {
		return -λ * (math.Sin(x) - λ*math.Cos(x) + λ*math.Exp(λ*x)) / (λ*λ + 1.0)
	}

	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = λ*y[0] - λ*math.Cos(x)
		return nil
	}

	o.Jac = func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(1, 1, 1)
		}
		dfdy.Start()
		dfdy.Put(0, 0, λ)
		return nil
	}
	return
}

// ProbVanDerPol returns the Van der Pol' Equation as given in Hairer-Wanner VII-p5 Eq.(1.5)
func ProbVanDerPol() (o *Problem) {

	o = new(Problem)
	eps := 1.0e-6
	o.Xf = 2.0
	o.Y = la.Vector([]float64{2.0, -0.6})
	o.Ndim = len(o.Y)

	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = y[1]
		f[1] = ((1.0-y[0]*y[0])*y[1] - y[0]) / eps
		return nil
	}

	o.Jac = func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(2, 2, 4)
		}
		dfdy.Start()
		dfdy.Put(0, 0, 0.0)
		dfdy.Put(0, 1, 1.0)
		dfdy.Put(1, 0, (-2.0*y[0]*y[1]-1.0)/eps)
		dfdy.Put(1, 1, (1.0-y[0]*y[0])/eps)
		return nil
	}
	return
}

// ProbRobertson returns the Robertson's Equation as given in Hairer-Wanner VII-p3 Eq.(1.4)
func ProbRobertson() (o *Problem) {

	o = new(Problem)
	o.Xf = 0.3
	o.Y = la.Vector([]float64{1.0, 0.0, 0.0})
	o.Ndim = len(o.Y)

	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = -0.04*y[0] + 1.0e4*y[1]*y[2]
		f[1] = 0.04*y[0] - 1.0e4*y[1]*y[2] - 3.0e7*y[1]*y[1]
		f[2] = 3.0e7 * y[1] * y[1]
		return nil
	}

	o.Jac = func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
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
		return nil
	}
	return
}

// ProbHwAmplifier returns the Hairer-Wanner VII-p376 Transistor Amplifier
// NOTE: from E. Hairer's website, not the equation in the book
func ProbHwAmplifier() (o *Problem) {

	// problem
	o = new(Problem)
	o.Xf = 0.05
	//o.Xf = 0.0123 // OK
	//o.Xf = 0.01235 // !OK

	// constants
	ue, ub, uf, α, β := 0.1, 6.0, 0.026, 0.99, 1.0E-6
	r0, r1, r2, r3, r4, r5 := 1000.0, 9000.0, 9000.0, 9000.0, 9000.0, 9000.0
	r6, r7, r8, r9 := 9000.0, 9000.0, 9000.0, 9000.0
	w := 2.0 * 3.141592654 * 100.0

	// initial values
	o.Y = la.Vector([]float64{0.0,
		ub,
		ub / (r6/r5 + 1.0),
		ub / (r6/r5 + 1.0),
		ub,
		ub / (r2/r1 + 1.0),
		ub / (r2/r1 + 1.0),
		0.0,
	})
	o.Ndim = len(o.Y)

	// right-hand side of the amplifier problem
	o.Fcn = func(f la.Vector, dx, x float64, y la.Vector) error {
		uet := ue * math.Sin(w*x)
		fac1 := β * (math.Exp((y[3]-y[2])/uf) - 1.0)
		fac2 := β * (math.Exp((y[6]-y[5])/uf) - 1.0)
		f[0] = y[0] / r9
		f[1] = (y[1]-ub)/r8 + α*fac1
		f[2] = y[2]/r7 - fac1
		f[3] = y[3]/r5 + (y[3]-ub)/r6 + (1.0-α)*fac1
		f[4] = (y[4]-ub)/r4 + α*fac2
		f[5] = y[5]/r3 - fac2
		f[6] = y[6]/r1 + (y[6]-ub)/r2 + (1.0-α)*fac2
		f[7] = (y[7] - uet) / r0
		return nil
	}

	// Jacobian of the amplifier problem
	o.Jac = func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		fac14 := β * math.Exp((y[3]-y[2])/uf) / uf
		fac27 := β * math.Exp((y[6]-y[5])/uf) / uf
		if dfdy.Max() == 0 {
			dfdy.Init(8, 8, 16)
		}
		nu := 2
		dfdy.Start()
		dfdy.Put(2+0-nu, 0, 1.0/r9)
		dfdy.Put(2+1-nu, 1, 1.0/r8)
		dfdy.Put(1+2-nu, 2, -α*fac14)
		dfdy.Put(0+3-nu, 3, α*fac14)
		dfdy.Put(2+2-nu, 2, 1.0/r7+fac14)
		dfdy.Put(1+3-nu, 3, -fac14)
		dfdy.Put(2+3-nu, 3, 1.0/r5+1.0/r6+(1.0-α)*fac14)
		dfdy.Put(3+2-nu, 2, -(1.0-α)*fac14)
		dfdy.Put(2+4-nu, 4, 1.0/r4)
		dfdy.Put(1+5-nu, 5, -α*fac27)
		dfdy.Put(0+6-nu, 6, α*fac27)
		dfdy.Put(2+5-nu, 5, 1.0/r3+fac27)
		dfdy.Put(1+6-nu, 6, -fac27)
		dfdy.Put(2+6-nu, 6, 1.0/r1+1.0/r2+(1.0-α)*fac27)
		dfdy.Put(3+5-nu, 5, -(1.0-α)*fac27)
		dfdy.Put(2+7-nu, 7, 1.0/r0)
		return nil
	}

	// "mass" matrix
	c1, c2, c3, c4, c5 := 1.0e-6, 2.0e-6, 3.0e-6, 4.0e-6, 5.0e-6
	o.M = new(la.Triplet)
	o.M.Init(8, 8, 14)
	o.M.Start()
	nu := 1
	o.M.Put(1+0-nu, 0, -c5)
	o.M.Put(0+1-nu, 1, c5)
	o.M.Put(2+0-nu, 0, c5)
	o.M.Put(1+1-nu, 1, -c5)
	o.M.Put(1+2-nu, 2, -c4)
	o.M.Put(1+3-nu, 3, -c3)
	o.M.Put(0+4-nu, 4, c3)
	o.M.Put(2+3-nu, 3, c3)
	o.M.Put(1+4-nu, 4, -c3)
	o.M.Put(1+5-nu, 5, -c2)
	o.M.Put(1+6-nu, 6, -c1)
	o.M.Put(0+7-nu, 7, c1)
	o.M.Put(2+6-nu, 6, c1)
	o.M.Put(1+7-nu, 7, -c1)
	return
}
