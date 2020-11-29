// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
)

func main() {

	// start mpi
	mpi.Start()
	defer mpi.Stop()

	// message
	chk.Verbose = true

	// check number of processors
	if mpi.WorldRank() == 0 {
		chk.Verbose = true
		chk.PrintTitle("Hairer-Wanner VII-p2 Eq.(1.1) (mumps)")
	}
	if mpi.WorldSize() != 1 {
		if mpi.WorldRank() == 0 {
			io.Pf("ERROR: this test can only be run with 1 processor\n")
		}
		return
	}

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
}
