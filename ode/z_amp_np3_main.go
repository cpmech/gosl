// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
)

func main() {

	// start mpi
	mpi.Start()
	defer mpi.Stop()

	// check number of processors
	if mpi.WorldRank() == 0 {
		chk.Verbose = true
		chk.PrintTitle("Hairer-Wanner VII-p376 Transistor Amplifier")
	}
	if mpi.WorldSize() != 3 {
		if mpi.WorldRank() == 0 {
			io.Pf("ERROR: this test needs 3 processors (run with mpi -np 3)\n")
		}
		return
	}

	// communicator
	comm := mpi.NewCommunicator(nil)

	// constants
	ue, ub, uf, α, β := 0.1, 6.0, 0.026, 0.99, 1.0e-6
	r0, r1, r2, r3, r4, r5 := 1000.0, 9000.0, 9000.0, 9000.0, 9000.0, 9000.0
	r6, r7, r8, r9 := 9000.0, 9000.0, 9000.0, 9000.0
	w := 2.0 * 3.141592654 * 100.0
	xf := 0.05

	// initial values
	y := la.NewVectorSlice([]float64{0.0,
		ub,
		ub / (r6/r5 + 1.0),
		ub / (r6/r5 + 1.0),
		ub,
		ub / (r2/r1 + 1.0),
		ub / (r2/r1 + 1.0),
		0.0,
	})
	ndim := len(y)

	// right-hand side of the amplifier problem
	fcn := func(f la.Vector, dx, x float64, y la.Vector) {
		uet := ue * math.Sin(w*x)
		fac1 := β * (math.Exp((y[3]-y[2])/uf) - 1.0)
		fac2 := β * (math.Exp((y[6]-y[5])/uf) - 1.0)
		f[0] = y[0] / r9
		f[1] = (y[1]-ub)/r8 + α*fac1
		f[2] = y[2]/r7 - fac1
		f[3] = y[3]/r5 + (y[3]-ub)/r6 + (1.0-α)*fac1
		f[4] = (y[4]-ub)/r4 + α*fac2
		f[5] = y[5]/r3 - fac2
		f[6] = y[6]/r1 + (y[6]-ub)/r2 + (1.0-α)*fac2
		f[7] = (y[7] - uet) / r0
	}

	// Jacobian of the amplifier problem
	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) {
		fac14 := β * math.Exp((y[3]-y[2])/uf) / uf
		fac27 := β * math.Exp((y[6]-y[5])/uf) / uf
		if dfdy.Max() == 0 {
			dfdy.Init(8, 8, 16)
		}
		nu := 2
		dfdy.Start()
		switch comm.Rank() {
		case 0:
			dfdy.Put(2+0-nu, 0, 1.0/r9)
			dfdy.Put(2+1-nu, 1, 1.0/r8)
			dfdy.Put(1+2-nu, 2, -α*fac14)
			dfdy.Put(0+3-nu, 3, α*fac14)
			dfdy.Put(2+2-nu, 2, 1.0/r7+fac14)
		case 1:
			dfdy.Put(1+3-nu, 3, -fac14)
			dfdy.Put(2+3-nu, 3, 1.0/r5+1.0/r6+(1.0-α)*fac14)
			dfdy.Put(3+2-nu, 2, -(1.0-α)*fac14)
			dfdy.Put(2+4-nu, 4, 1.0/r4)
			dfdy.Put(1+5-nu, 5, -α*fac27)
		case 2:
			dfdy.Put(0+6-nu, 6, α*fac27)
			dfdy.Put(2+5-nu, 5, 1.0/r3+fac27)
			dfdy.Put(1+6-nu, 6, -fac27)
			dfdy.Put(2+6-nu, 6, 1.0/r1+1.0/r2+(1.0-α)*fac27)
			dfdy.Put(3+5-nu, 5, -(1.0-α)*fac27)
			dfdy.Put(2+7-nu, 7, 1.0/r0)
		}
	}

	// "mass" matrix
	c1, c2, c3, c4, c5 := 1.0e-6, 2.0e-6, 3.0e-6, 4.0e-6, 5.0e-6
	M := new(la.Triplet)
	M.Init(8, 8, 14)
	M.Start()
	nu := 1
	switch comm.Rank() {
	case 0:
		M.Put(1+0-nu, 0, -c5)
		M.Put(0+1-nu, 1, c5)
		M.Put(2+0-nu, 0, c5)
		M.Put(1+1-nu, 1, -c5)
		M.Put(1+2-nu, 2, -c4)
		M.Put(1+3-nu, 3, -c3)
	case 1:
		M.Put(0+4-nu, 4, c3)
		M.Put(2+3-nu, 3, c3)
		M.Put(1+4-nu, 4, -c3)
	case 2:
		M.Put(1+5-nu, 5, -c2)
		M.Put(1+6-nu, 6, -c1)
		M.Put(0+7-nu, 7, c1)
		M.Put(2+6-nu, 6, c1)
		M.Put(1+7-nu, 7, -c1)
	}

	// configurations
	conf := ode.NewConfig("radau5", "mumps", comm)
	conf.SetStepOut(true, nil)
	conf.IniH = 1.0e-6 // initial step size

	// set tolerances
	atol, rtol := 1e-11, 1e-5
	conf.SetTols(atol, rtol)

	// ODE solver
	sol := ode.NewSolver(ndim, conf, fcn, jac, M)
	defer sol.Free()

	// run
	sol.Solve(y, 0.0, xf)

	// only root
	if mpi.WorldRank() == 0 {

		// check
		tst := new(testing.T)
		chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2559)
		chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 212)
		chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 266)
		chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 217)
		chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 23)
		chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 265)
		chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 780)
		chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 6)
		chk.String(tst, sol.Stat.LsKind, "mumps")
	}
}
