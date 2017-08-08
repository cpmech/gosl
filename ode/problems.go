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
	Yana YanaF     // analytical solution
	Fcn  Func      // the function f(x,y) = dy/df
	Jac  JacF      // df/dy function
	Dx   float64   // timestep for fixedStep tests
	Xf   float64   // final x
	Y    la.Vector // initial (current) y vector
}

// Solve solves ODE problem using standard parameters
// NOTE: this solver doesn't change o.Y
func (o *Problem) Solve(method io.Enum, fixedStp, numJac bool) (stat *Stat, out *Output, err error) {

	// current y vector
	ndim := len(o.Y)
	y := la.NewVector(ndim)
	y.Apply(1, o.Y)

	// configuration
	conf, err := NewConfig(method, "", nil, nil)
	conf.SaveXY = true
	if fixedStp {
		conf.FixedStp = o.Dx
	}

	// allocate solver
	jac := o.Jac
	if numJac {
		jac = nil
	}
	sol, err := NewSolver(conf, ndim, o.Fcn, jac, nil, nil)
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

// ProbHwEq11 returns the Hairer-Wanner problem from VII-p2 Eq.(1.1)
func ProbHwEq11() (o *Problem) {

	o = new(Problem)
	λ := -50.0
	o.Dx = 1.875 / 50.0
	o.Xf = 1.5
	o.Y = la.Vector([]float64{0.0})

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
