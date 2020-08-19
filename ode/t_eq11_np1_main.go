// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"gosl/chk"
	"gosl/mpi"
	"gosl/ode"
	"gosl/plt"
)

func main() {

	// start mpi
	mpi.Start()
	defer mpi.Stop()

	// message
	chk.Verbose = true
	chk.PrintTitle("Hairer-Wanner VII-p2 Eq.(1.1) (mumps)")

	// communicator
	comm := mpi.NewCommunicator(nil)

	// problem
	p := ode.ProbHwEq11()

	// configuration
	conf := ode.NewConfig("radau5", "mumps", comm)
	conf.SetStepOut(true, nil)

	// solver
	sol := ode.NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(p.Y, 0.0, p.Xf)

	// check Stat
	chk.Verbose = true
	tst := new(testing.T)
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 1)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)
	chk.String(tst, sol.Stat.LsKind, "mumps")

	// check results
	chk.Float64(tst, "yFin", 2.94e-5, p.Y[0], p.CalcYana(0, p.Xf))

	// plot
	plt.Reset(true, nil)
	p.Plot("Radau5,Jana", 0, sol.Out, 101, true, nil, nil)
	plt.Save("/tmp/gosl/ode", "mpi_eq11_np1")
}
