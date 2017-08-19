// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestDoPri802(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DoPri802. Dormand-Prince8(5,3). Van de Pol")

	// problem
	p := ProbVanDerPol(1e-3, false)
	p.Y[0] = 2.0
	p.Y[1] = 0.0
	p.Xf = 0.2

	// configuration
	conf, err := NewConfig("dopri8", "", nil)
	status(tst, err)
	conf.SetTol(1e-9, 1e-9)
	conf.Mmin = 0.333
	conf.Mmax = 6.0
	conf.PredCtrl = false
	conf.NmaxSS = 2000

	// output handler
	conf.SetStepOut(true, nil)
	out := NewOutput(p.Ndim, conf)

	// solver
	sol, err := NewSolver(p.Ndim, conf, out, p.Fcn, p.Jac, nil)
	status(tst, err)
	defer sol.Free()

	// solve ODE
	err = sol.Solve(p.Y, 0.0, p.Xf)
	status(tst, err)

	// print stat
	sol.Stat.Print(false)

	// check Stat
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 1924)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 163)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 130)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 33)

	// check results
	H := out.GetStepH()
	X := out.GetStepX()
	Y0 := out.GetStepY(0)
	Y1 := out.GetStepY(1)
	_, d, err := io.ReadTable("data/dr_dop853.txt")
	status(tst, err)
	chk.Array(tst, "h", 1e-6, H, d["h"])
	chk.Array(tst, "x", 1e-6, X, d["x"])
	chk.Array(tst, "y0", 1e-6, Y0, d["y0"])
	chk.Array(tst, "y1", 1e-6, Y1, d["y1"])
}
