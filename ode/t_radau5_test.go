// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestRadau501(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Radau501. Radau5 (analytical Jacobian)")

	_, xf, y, yana, fcn, jac := eq11data()

	conf, err := NewConfig(Radau5kind, nil, "", 0, 0)
	status(tst, err)
	conf.SaveXY = true

	sol, err := NewSolver(conf, 1, fcn, jac, nil, nil)
	status(tst, err)

	err = sol.Solve(y, 0.0, xf)
	status(tst, err)

	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 1)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	chk.Float64(tst, "yFin", 2.88898538383e-5, y[0], yana(xf))

	if chk.Verbose {
		eq11plotOne("radau501", "Radau5,Jana", xf, yana, sol)
	}
}
