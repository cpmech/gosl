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

func TestDoPri501(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DoPri501. Dormand-Prince5. Hw Eq (1.1)")

	// problem
	p := ProbHwEq11()

	// configuration
	conf := NewConfig("dopri5", "")
	conf.SetStepOut(true, nil)

	// solver
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(p.Y, 0.0, p.Xf)

	// check Stat
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 242)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 0)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 40)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 0)

	// check results
	chk.Float64(tst, "yFin", 0.0179276654, p.Y[0], p.CalcYana(0, p.Xf))
}

func TestDoPri502(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DoPri502. Dormand-Prince5. Arenstorf orbit")

	// problem
	p := ProbArenstorf()

	// configuration
	conf := NewConfig("dopri5", "")
	conf.SetTols(1e-7, 1e-7)
	conf.Mmin = 0.2
	conf.Mmax = 10.0
	conf.PredCtrl = false

	// step output
	io.Pf("\n%5s%15s%12s%15s%15s%15s%15s\n", "s", "h", "x", "y0", "y1", "y2", "y3")
	conf.SetStepOut(true, func(istep int, h, x float64, y la.Vector) (stop bool) {
		io.Pf("%5d%15.7E%12.7f%15.7E%15.7E%15.7E%15.7E\n", istep, h, x, y[0], y[1], y[2], y[3])
		return false
	})

	// dense output function
	ss := make([]int, 10)
	xx := make([]float64, 10)
	yy0 := make([]float64, 10)
	yy1 := make([]float64, 10)
	yy2 := make([]float64, 10)
	yy3 := make([]float64, 10)
	iout := 0
	conf.SetDenseOut(true, 2.0, p.Xf, func(istep int, h, x float64, y la.Vector, xout float64, yout la.Vector) (stop bool) {
		xold := x - h
		dx := xout - xold
		io.Pforan("%5d%15.7E%12.7f%15.7E%15.7E%15.7E%15.7E\n", istep, dx, xout, yout[0], yout[1], yout[2], yout[3])
		ss[iout] = istep
		xx[iout] = xout
		yy0[iout] = yout[0]
		yy1[iout] = yout[1]
		yy2[iout] = yout[2]
		yy3[iout] = yout[3]
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
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 1430)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 238)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 217)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 21)

	// check results: step
	_, d := io.ReadTable("data/dr_dopri5.txt")
	chk.Array(tst, "step: y0", 1e-6, sol.Out.GetStepY(0), d["y0"])
	chk.Array(tst, "step: y1", 1e-6, sol.Out.GetStepY(1), d["y1"])
	chk.Array(tst, "step: y2", 1e-6, sol.Out.GetStepY(2), d["y2"])
	chk.Array(tst, "step: y3", 1e-6, sol.Out.GetStepY(3), d["y3"])

	// check results: dense
	_, dd := io.ReadTable("data/dr_dopri5_dense.txt")
	chk.Array(tst, "dense: y0", 1e-7, yy0, dd["y0"])
	chk.Array(tst, "dense: y1", 1e-7, yy1, dd["y1"])
	chk.Array(tst, "dense: y2", 1e-7, yy2, dd["y2"])
	chk.Array(tst, "dense: y3", 1e-7, yy3, dd["y3"])
}

func TestDoPri503(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DoPri503. Dormand-Prince5(4). Van de Pol")

	// problem
	p := ProbVanDerPol(1e-3, false)
	p.Y[0] = 2.0
	p.Y[1] = 0.0
	p.Xf = 0.2

	// configuration
	conf := NewConfig("dopri5", "")
	conf.SetTols(1e-7, 1e-7)
	conf.Mmin = 0.2
	conf.Mmax = 10.0
	conf.PredCtrl = false
	conf.NmaxSS = 2000
	conf.MfirstRej = 0

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
		_ = dx
		//io.Pforan("%6d%15.7E%12.7f%15.7E%15.7E\n", istep, dx, xout, yout[0], yout[1])
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
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 1274)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 212)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 211)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)

	// check results: dense
	_, dd := io.ReadTable("data/dr2_dopri5_dense.txt")
	chk.Array(tst, "dense: y0", 1e-7, yy0, dd["y0"])
	chk.Array(tst, "dense: y1", 1e-7, yy1, dd["y1"])
}
