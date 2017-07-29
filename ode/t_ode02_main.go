// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"
	"time"

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

	comm := mpi.NewCommunicator(nil)

	if comm.Rank() == 0 {
		chk.PrintTitle("ode02: Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation")
	}
	if comm.Size() != 2 {
		chk.Panic(">> error: this test requires 2 MPI processors\n")
		return
	}

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
		// join all f
		comm.AllReduceSum(f, w)
		return nil
	}
	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(2, 2, 4)
		}
		dfdy.Start()
		if false { // per column
			switch comm.Rank() {
			case 0:
				dfdy.Put(0, 0, 0.0)
				dfdy.Put(1, 0, (-2.0*y[0]*y[1]-1.0)/eps)
			case 1:
				dfdy.Put(0, 1, 1.0)
				dfdy.Put(1, 1, (1.0-y[0]*y[0])/eps)
			}
		} else { // per row
			switch comm.Rank() {
			case 0:
				dfdy.Put(0, 0, 0.0)
				dfdy.Put(0, 1, 1.0)
			case 1:
				dfdy.Put(1, 0, (-2.0*y[0]*y[1]-1.0)/eps)
				dfdy.Put(1, 1, (1.0-y[0]*y[0])/eps)
			}
		}
		return nil
	}

	// method and flags
	fixstp := false
	method := ode.Radau5kind
	numjac := false
	xa, xb := 0.0, 2.0
	ya := []float64{2.0, -0.6}
	ndim := len(ya)

	// allocate ODE object
	var o *ode.Solver
	if numjac {
		o = ode.NewSolver(method, ndim, fcn, nil, nil, nil)
	} else {
		o = ode.NewSolver(method, ndim, fcn, jac, nil, nil)
	}
	o.SaveXY = true
	o.Distr = true

	// tolerances and initial step size
	rtol := 1e-4
	atol := rtol
	o.IniH = 1.0e-4
	o.SetTol(atol, rtol)
	//o.NmaxSS = 2

	// solve problem
	y := make([]float64, ndim)
	copy(y, ya)
	t0 := time.Now()
	if fixstp {
		o.Solve(y, xa, xb, 0.05, fixstp)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp)
	}

	// check
	if comm.Rank() == 0 {
		chk.Verbose = true
		tst := new(testing.T)
		chk.Int(tst, "number of F evaluations ", o.Nfeval, 2233)
		chk.Int(tst, "number of J evaluations ", o.Njeval, 160)
		chk.Int(tst, "total number of steps   ", o.Nsteps, 280)
		chk.Int(tst, "number of accepted steps", o.Naccepted, 241)
		chk.Int(tst, "number of rejected steps", o.Nrejected, 7)
		chk.Int(tst, "number of decompositions", o.Ndecomp, 251)
		chk.Int(tst, "number of lin solutions ", o.Nlinsol, 663)
		chk.Int(tst, "max number of iterations", o.Nitmax, 6)
	}

	// plot
	if comm.Rank() == 0 {
		io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T, err := io.ReadTable("data/vdpol_radau5_for.dat")
		if err != nil {
			chk.Panic("%v", err)
		}
		s := o.IdxSave
		for j := 0; j < ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			plt.Subplot(ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(o.Xvalues[:s], o.Yvalues[j][:s], &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(ndim+1, 1, ndim+1)
		plt.Plot(o.Xvalues[1:s], o.Hvalues[1:s], &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "vdpolA_mpi")
	}
}
