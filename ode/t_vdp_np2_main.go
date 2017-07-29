// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/plt"
)

func main() {

	mpi.Start()
	defer mpi.Stop()

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

	comm := mpi.NewCommunicator(nil)

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

	xa, xb := 0.0, 2.0
	ndim := 2
	y := la.Vector([]float64{2.0, -0.6})

	sol := ode.NewSolver(ode.Radau5kind, ndim, fcn, jac, nil, nil)
	sol.SaveXY = true
	sol.Distr = true // <<<<<<< distributed mode

	rtol := 1e-4
	atol := rtol
	sol.IniH = 1.0e-4
	sol.SetTol(atol, rtol)

	sol.Solve(y, xa, xb, xb-xa, false)

	if mpi.WorldRank() == 0 {
		tst := new(testing.T)
		chk.Int(tst, "number of F evaluations ", sol.Nfeval, 2233)
		chk.Int(tst, "number of J evaluations ", sol.Njeval, 160)
		chk.Int(tst, "total number of steps   ", sol.Nsteps, 280)
		chk.Int(tst, "number of accepted steps", sol.Naccepted, 241)
		chk.Int(tst, "number of rejected steps", sol.Nrejected, 7)
		chk.Int(tst, "number of decompositions", sol.Ndecomp, 251)
		chk.Int(tst, "number of lin solutions ", sol.Nlinsol, 663)
		chk.Int(tst, "max number of iterations", sol.Nitmax, 6)
		chk.Int(tst, "IdxSave", sol.IdxSave, sol.Naccepted+1)

		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T, err := io.ReadTable("data/vdpol_radau5_for.dat")
		chk.EP(err)
		s := sol.IdxSave
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(sol.Xvalues[:s], sol.Yvalues[j][:s], &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(sol.Xvalues[1:s], sol.Hvalues[1:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "vdp_np2")
	}
}
