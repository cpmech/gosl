// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"os"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// start mpi
	mpi.Start()
	defer mpi.Stop()

	//check number of processors
	if mpi.WorldRank() == 0 {
		chk.Verbose = true
		chk.PrintTitle("Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation (Distr=true)")
	}
	if mpi.WorldSize() != 2 {
		if mpi.WorldRank() == 0 {
			io.Pf("ERROR: this test needs 2 processors (run with mpi -np 2)\n")
		}
		return
	}

	// communicator
	comm := mpi.NewCommunicator(nil)

	// dy/dx function
	eps := 1.0e-6
	w := la.NewVector(2) // workspace
	fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
		w.Fill(0)
		switch comm.Rank() {
		case 0:
			w[0] = y[1]
		case 1:
			w[1] = ((1.0-y[0]*y[0])*y[1] - y[0]) / eps
		}
		comm.AllReduceSum(f, w)
		return nil
	}

	// Jacobian
	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(2, 2, 4)
		}
		dfdy.Start()
		switch comm.Rank() {
		case 0:
			dfdy.Put(0, 0, 0.0)
			dfdy.Put(0, 1, 1.0)
		case 1:
			dfdy.Put(1, 0, (-2.0*y[0]*y[1]-1.0)/eps)
			dfdy.Put(1, 1, (1.0-y[0]*y[0])/eps)
		}
		return nil
	}

	// initial values
	xb := 2.0
	ndim := 2
	y := la.Vector([]float64{2.0, -0.6})

	// configurations
	conf := ode.NewConfig("radau5", "", comm)
	conf.SetStepOut(true, nil)
	conf.SetTol(1e-4)

	// solver
	sol, err := ode.NewSolver(ndim, conf, fcn, jac, nil)
	status(err)

	// solve
	err = sol.Solve(y, 0, xb)
	status(err)

	// only root
	if mpi.WorldRank() == 0 {

		//check
		tst := new(testing.T)
		chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2233)
		chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 160)
		chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 280)
		chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 241)
		chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 7)
		chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 251)
		chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 663)
		chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 6)

		// plot
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T := io.ReadTable("data/vdpol_radau5_for.dat")
		X := sol.Out.GetStepX()
		H := sol.Out.GetStepH()
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := sol.Out.GetStepY(j)
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(X, Yj, &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(X, H, &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "vdp_np2")
	}
}

func status(err error) {
	if err != nil {
		io.Pf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
