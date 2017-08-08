// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
)

func TestRadau501a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Radau501a. Radau5 (analytical Jacobian)")

	// problem
	p := ProbHwEq11()
	ndim := len(p.Y)

	// configuration
	conf, err := NewConfig(Radau5kind, "", nil, nil)
	status(tst, err)
	conf.SaveXY = true

	// solver
	sol, err := NewSolver(conf, ndim, p.Fcn, p.Jac, nil, nil)
	status(tst, err)
	defer sol.Free()

	// solve ODE
	err = sol.Solve(p.Y, 0.0, p.Xf)
	status(tst, err)

	// check Stat
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 1)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	// check results
	chk.Float64(tst, "yFin", 2.88898538383e-5, p.Y[0], p.Yana(p.Xf))

	// plot
	if chk.Verbose {
		plt.Reset(true, nil)
		p.Plot("Radau5,Jana", sol.Out, 101, true, nil, nil)
		plt.Save("/tmp/gosl/ode", "radau501a")
	}
}
