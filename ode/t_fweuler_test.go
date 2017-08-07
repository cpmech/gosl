// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestFwEuler01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FwEuler01. Forward-Euler")

	dx, xf, y, yana, fcn, _ := eq11data()

	conf, err := NewConfig(FwEulerKind, nil, "", 0, 0)
	status(tst, err)
	conf.SaveXY = true
	conf.FixedStp = dx

	sol, err := NewSolver(conf, 1, fcn, nil, nil, nil)
	status(tst, err)

	err = sol.Solve(y, 0.0, xf)
	status(tst, err)

	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 40)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 0)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 0)

	chk.Float64(tst, "yFin", 0.004753, y[0], yana(xf))

	if chk.Verbose {
		eq11plotOne("fweuler01", "FwEuler", xf, yana, sol)
	}
}
