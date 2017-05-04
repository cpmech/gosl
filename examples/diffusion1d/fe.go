// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// constants
	doPlot := true
	problem := 1
	epsDtCrit := 0.0
	tf := 10.0
	nx := 21

	// problem
	p := Problem(problem)

	// grid
	dx := 1.0 / float64(nx-1)      // grid size
	X := utl.LinSpace(0, p.Xf, nx) // grid points

	// time discretisation
	dtCrit := 0.5 * dx * dx / p.Kx
	dt := dtCrit + epsDtCrit

	// dUbdt = f(t,Ub)
	m := nx - 2 // problem dimension: inner points = nx - end points
	fcn, jac := p.Version1(m, dx, nx < 11)

	// initial values
	Ub := la.VecGetMapped(m, func(i int) float64 { return p.U0(X[i+1]) })

	// ode solver
	var o ode.Solver
	o.Init("FwEuler", m, fcn, jac, nil, nil)
	o.SaveXY = true
	o.NmaxSS = int(tf/dt) + 1

	// solve problem
	fixedStep := true
	T0 := time.Now()
	err := o.Solve(Ub, 0, tf, dt, fixedStep)
	if err != nil {
		chk.Panic("%v", err)
	}
	elapsedTime := time.Now().Sub(T0)

	// compute error
	emax := p.Error(dx, &o)

	// results
	o.Stat()
	io.Pf("\ndt = %v\n", dt)
	io.Pf("max(error) = %e\n", emax)
	io.Pf("elapsed time = %v\n", elapsedTime)

	// plot
	if doPlot {
		io.Pf("\n")
		p.Results(dx, &o, 30, 30)
		p.PlotRes3d("/tmp/gosl", "fe-version1-res")
		p.PlotErr3d("/tmp/gosl", "fe-version1-err")
	}
}
