// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
)

func main() {

	if mpi.Rank() == 0 {
		chk.PrintTitle("Test ODE 02b")
		io.Pfcyan("Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation (MPI)\n")
	}
	if mpi.Size() != 2 {
		chk.Panic(">> error: this test requires 2 MPI processors\n")
		return
	}

	eps := 1.0e-6
	w := make([]float64, 2) // workspace
	fcn := func(f []float64, x float64, y []float64, args ...interface{}) error {
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
	jac := func(dfdy *la.Triplet, x float64, y []float64, args ...interface{}) error {
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

	// data
	silent := false
	fixstp := false
	//method := "Dopri5"
	method := "Radau5"
	xa, xb := 0.0, 2.0
	ya := []float64{2.0, -0.6}
	ndim := len(ya)

	// output
	var b bytes.Buffer
	out := func(first bool, dx, x float64, y []float64, args ...interface{}) error {
		if mpi.Rank() == 0 {
			if first {
				fmt.Fprintf(&b, "%23s %23s %23s %23s\n", "dx", "x", "y0", "y1")
			}
			fmt.Fprintf(&b, "%23.15E %23.15E %23.15E %23.15E\n", dx, x, y[0], y[1])
		}
		return nil
	}
	defer func() {
		if mpi.Rank() == 0 {
			extra := "d2 = Read('data/vdpol_radau5_for.dat')\n" +
				"subplot(3,1,1)\n" +
				"plot(d2['x'],d2['y0'],'k+',label='res',ms=10)\n" +
				"subplot(3,1,2)\n" +
				"plot(d2['x'],d2['y1'],'k+',label='res',ms=10)\n"
			ode.Plot("/tmp/gosl", "vdpolB", method, &b, []int{0, 1}, ndim, nil, xa, xb, true, false, extra)
		}
	}()

	// one run
	var o ode.ODE
	o.Distr = true
	//numjac := true
	numjac := false
	if numjac {
		o.Init(method, ndim, fcn, nil, nil, out, silent)
	} else {
		o.Init(method, ndim, fcn, jac, nil, out, silent)
	}

	// tolerances and initial step size
	rtol := 1e-4
	atol := rtol
	o.SetTol(atol, rtol)
	o.IniH = 1.0e-4

	//o.NmaxSS = 2

	y := make([]float64, ndim)
	copy(y, ya)
	t0 := time.Now()
	if fixstp {
		o.Solve(y, xa, xb, 0.05, fixstp)
	} else {
		o.Solve(y, xa, xb, xb-xa, fixstp)
	}
	if mpi.Rank() == 0 {
		io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))
	}
}
