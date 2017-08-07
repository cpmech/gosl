// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestBwEuler01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BwEuler01a. Backward-Euler (analytical Jacobian)")

	dx, xf, y, yana, fcn, jac := eq11data()

	conf, err := NewConfig(BwEulerKind, nil, "", 0, 0)
	status(tst, err)
	conf.SaveXY = true
	conf.FixedStp = dx

	sol, err := NewSolver(conf, 1, fcn, jac, nil, nil)
	status(tst, err)

	err = sol.Solve(y, 0.0, xf)
	status(tst, err)

	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 80)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 40)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 40)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 40)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	chk.Float64(tst, "yFin", 1e-4, y[0], yana(xf))

	if chk.Verbose {
		eq11plotOne("bweuler01a", "BwEuler,Jana", xf, yana, sol)
	}
}

func TestBwEuler01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BwEuler01b. Backward-Euler (numerical Jacobian)")

	dx, xf, y, yana, fcn, _ := eq11data()

	conf, err := NewConfig(BwEulerKind, nil, "", 0, 0)
	status(tst, err)
	conf.SaveXY = true
	conf.FixedStp = dx

	sol, err := NewSolver(conf, 1, fcn, nil, nil, nil)
	status(tst, err)

	err = sol.Solve(y, 0.0, xf)
	status(tst, err)

	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 120)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 40)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 40)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 40)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	chk.Float64(tst, "yFin", 1e-4, y[0], yana(xf))

	if chk.Verbose {
		eq11plotOne("bweuler01b", "BwEuler,Jnum", xf, yana, sol)
	}
}
