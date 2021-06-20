// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
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
	conf := NewConfig("dopri8", "")
	conf.SetTols(1e-9, 1e-9)
	conf.Mmin = 0.333
	conf.Mmax = 6.0
	conf.PredCtrl = false
	conf.NmaxSS = 2000

	// step output
	io.Pf("\n%6s%15s%15s%15s%15s\n", "s", "h", "x", "y0", "y1")
	conf.SetStepOut(true, func(istep int, h, x float64, y la.Vector) (stop bool) {
		io.Pf("%6d%15.7E%12.7f%15.7E%15.7E\n", istep, h, x, y[0], y[1])
		return false
	})

	// dense output function
	ss := make([]int, 11)
	xx := make([]float64, 11)
	yy0 := make([]float64, 11)
	yy1 := make([]float64, 11)
	iout := 0
	conf.SetDenseOut(true, 0.02, p.Xf, func(istep int, h, x float64, y la.Vector, xout float64, yout la.Vector) (stop bool) {
		xold := x - h
		dx := xout - xold
		io.Pforan("%6d%15.7E%12.7f%15.7E%15.7E\n", istep, dx, xout, yout[0], yout[1])
		ss[iout] = istep
		xx[iout] = xout
		yy0[iout] = yout[0]
		yy1[iout] = yout[1]
		iout++
		return
	})

	// solver
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(p.Y, 0.0, p.Xf)

	// print stat
	sol.Stat.Print()

	// check Stat
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2314)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 163)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 130)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 33)

	// check results: setps
	_, d := io.ReadTable("data/dr_dop853.txt")
	chk.Array(tst, "h", 1e-6, sol.Out.GetStepH(), d["h"])
	chk.Array(tst, "x", 1e-6, sol.Out.GetStepX(), d["x"])
	chk.Array(tst, "y0", 1e-6, sol.Out.GetStepY(0), d["y0"])
	chk.Array(tst, "y1", 1e-6, sol.Out.GetStepY(1), d["y1"])

	// check results: dense
	_, dd := io.ReadTable("data/dr_dop853_dense.txt")
	chk.Array(tst, "dense: y0", 1e-7, yy0, dd["y0"])
	chk.Array(tst, "dense: y1", 1e-7, yy1, dd["y1"])
}
