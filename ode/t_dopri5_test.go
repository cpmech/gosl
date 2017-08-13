// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestDoPri501(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DoPri501. Dormand-Prince5. Hw Eq (1.1)")

	// problem
	p := ProbHwEq11()

	// configuration
	conf, err := NewConfig("dopri5", "", nil)
	status(tst, err)
	conf.SetStepOut(true, nil)

	// output handler
	out := NewOutput(p.Ndim, conf)

	// solver
	sol, err := NewSolver(p.Ndim, conf, out, p.Fcn, p.Jac, nil)
	status(tst, err)
	defer sol.Free()

	// solve ODE
	err = sol.Solve(p.Y, 0.0, p.Xf)
	status(tst, err)

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

	// plot
	if chk.Verbose {
		plt.Reset(true, nil)
		p.Plot("DoPri5", 0, out, 101, true, nil, nil)
		plt.Save("/tmp/gosl/ode", "dopri501")
	}
}

func TestDoPri502(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DoPri502. Dormand-Prince5. Arenstorf orbit")

	// problem
	p := ProbArenstorf()

	// configuration
	conf, err := NewConfig("dopri5", "", nil)
	status(tst, err)
	conf.SetTol(1e-7, 1e-7)
	conf.Mmin = 0.2
	conf.Mmax = 10.0
	conf.PredCtrl = false

	// output handler
	io.Pf("%6s%10s%10s%14s%14s%14s%14s\n", "istep", "h", "x", "y0", "y1", "y2", "y3")
	conf.SetStepOut(true, func(istep int, h, x float64, y la.Vector) (stop bool, err error) {
		io.Pf("%6d%10.6f%10.6f%14.6E%14.6E%14.6E%14.6E\n", istep, h, x, y[0], y[1], y[2], y[3])
		return false, nil
	})
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
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 1430)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 238)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 217)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 21)

	// check results
	_, d, err := io.ReadTable("data/dr_dopri5.txt")
	status(tst, err)
	chk.Array(tst, "y0", 1e-6, out.GetStepY(0), d["y0"])
	chk.Array(tst, "y1", 1e-6, out.GetStepY(1), d["y1"])
	chk.Array(tst, "y2", 1e-6, out.GetStepY(2), d["y2"])
	chk.Array(tst, "y3", 1e-6, out.GetStepY(3), d["y3"])
}
