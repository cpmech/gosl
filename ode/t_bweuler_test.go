// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"gosl/chk"
	"gosl/plt"
)

func TestBwEuler01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BwEuler01a. Backward-Euler (analytical Jacobian)")

	// problem
	p := ProbHwEq11()

	// configuration
	conf := NewConfig("bweuler", "", nil)
	conf.SetFixedH(p.Dx, p.Xf)
	conf.SetStepOut(true, nil)

	// solver
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(p.Y, 0.0, p.Xf)

	// check Stat
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 80)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 40)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 40)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 40)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	// check results
	chk.Float64(tst, "yFin", 1e-4, p.Y[0], p.CalcYana(0, p.Xf))

	// plot
	if chk.Verbose {
		plt.Reset(true, nil)
		p.Plot("BwEuler,Jana", 0, sol.Out, 101, true, nil, nil)
		plt.Save("/tmp/gosl/ode", "bweuler01a")
	}
}

func TestBwEuler01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BwEuler01b. Backward-Euler (numerical Jacobian)")

	// problem
	p := ProbHwEq11()

	// configuration
	conf := NewConfig("bweuler", "", nil)
	conf.SetFixedH(p.Dx, p.Xf)
	conf.SetStepOut(true, nil)

	// solver
	sol := NewSolver(p.Ndim, conf, p.Fcn, nil, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(p.Y, 0.0, p.Xf)

	// check Stat
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 120)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 40)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 40)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 40)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	// check results
	chk.Float64(tst, "yFin", 1e-4, p.Y[0], p.CalcYana(0, p.Xf))

	// plot
	if chk.Verbose {
		plt.Reset(true, nil)
		p.Plot("BwEuler,Jnum", 0, sol.Out, 101, true, nil, nil)
		plt.Save("/tmp/gosl/ode", "bweuler01b")
	}
}
