// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/plt"
)

func main() {

	mpi.Start(false)
	defer func() {
		mpi.Stop(false)
	}()

	if mpi.Rank() == 0 {
		chk.PrintTitle("ode02: Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation")
	}
	if mpi.Size() != 2 {
		chk.Panic(">> error: this test requires 2 MPI processors\n")
		return
	}

	eps := 1.0e-6
	w := make([]float64, 2) // workspace
	fcn := func(f []float64, dx, x float64, y []float64, args ...interface{}) error {
		f[0], f[1] = 0, 0
		switch mpi.Rank() {
		case 0:
			f[0] = y[1]
		case 1:
			f[1] = ((1.0-y[0]*y[0])*y[1] - y[0]) / eps
		}
		// join all f
		mpi.AllReduceSum(f, w)
		return nil
	}
	jac := func(dfdy *la.Triplet, dx, x float64, y []float64, args ...interface{}) error {
		if dfdy.Max() == 0 {
			dfdy.Init(2, 2, 4)
		}
		dfdy.Start()
		if false { // per column
			switch mpi.Rank() {
			case 0:
				dfdy.Put(0, 0, 0.0)
				dfdy.Put(1, 0, (-2.0*y[0]*y[1]-1.0)/eps)
			case 1:
				dfdy.Put(0, 1, 1.0)
				dfdy.Put(1, 1, (1.0-y[0]*y[0])/eps)
			}
		} else { // per row
			switch mpi.Rank() {
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
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	numjac := false
	xa, xb := 0.0, 2.0
	ya := []float64{2.0, -0.6}
	ndim := len(ya)

	// structure to hold numerical results
	res := ode.Results{Method: method}

	// allocate ODE object
	var o ode.Solver
	o.Distr = true
	if numjac {
		o.Init(method, ndim, fcn, nil, nil, ode.SimpleOutput, silent)
	} else {
		o.Init(method, ndim, fcn, jac, nil, ode.SimpleOutput, silent)
	}

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
		o.Solve(y, xa, xb, 0.05, fixstp, &res)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp, &res)
	}

	// plot
	if mpi.Rank() == 0 {
		io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))
		ode.Plot("/tmp/gosl/ode", "vdpolA_mpi", &res, nil, xa, xb, func() {
			_, T, err := io.ReadTable("data/vdpol_radau5_for.dat")
			if err != nil {
				chk.Panic("%v", err)
			}
			plt.Subplot(3, 1, 1)
			plt.Plot(T["x"], T["y0"], &plt.A{C: "k", M: "+", L: "reference"})
			plt.Subplot(3, 1, 2)
			plt.Plot(T["x"], T["y1"], &plt.A{C: "k", M: "+", L: "reference"})
		})
	}
}
