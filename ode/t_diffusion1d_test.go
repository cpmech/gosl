// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

func TestDiffusion1D(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Test Diffusion 1D (cooling)")

	// solution parameters
	silent := false
	fixstp := true
	//fixstp := false
	//method := "FwEuler"
	method := "BwEuler"
	//method := "Dopri5"
	//method := "Radau5"
	//numjac := true
	numjac := false
	rtol := 1e-4
	atol := rtol

	// timestep
	t0, tf, dt := 0.0, 0.2, 0.03

	// problem data
	kx := 1.0 // conductivity
	N := 6    // number of nodes
	//Nb   := N + 2             // augmented system dimension
	xmax := 1.0               // length
	dx := xmax / float64(N-1) // spatial step size
	dxx := dx * dx
	mol := []float64{kx / dxx, -2.0 * kx / dxx, kx / dxx}

	// function
	fcn := func(f []float64, t float64, y []float64, args ...interface{}) error {
		for i := 0; i < N; i++ {
			f[i] = 0
			if i == 0 || i == N-1 {
				continue // skip presc node
			}
			for p, j := range []int{i - 1, i, i + 1} {
				if j < 0 {
					j = i + 1
				} //  left boundary
				if j == N {
					j = i - 1
				} //  right boundary
				f[i] += mol[p] * y[j]
			}
		}
		//io.Pfgrey("y = %v\n", y)
		//io.Pfcyan("f = %v\n", f)
		return nil
	}

	// Jacobian
	jac := func(dfdy *la.Triplet, t float64, y []float64, args ...interface{}) error {
		//chk.Panic("jac is not available")
		if dfdy.Max() == 0 {
			//dfdy.Init(Nb, Nb, 3*N)
			dfdy.Init(N, N, 3*N)
		}
		dfdy.Start()
		for i := 0; i < N; i++ {
			if i == 0 || i == N-1 {
				dfdy.Put(i, i, 0.0)
				continue
			}
			for p, j := range []int{i - 1, i, i + 1} {
				if j < 0 {
					j = i + 1
				} //  left boundary
				if j == N {
					j = i - 1
				} //  right boundary
				dfdy.Put(i, j, mol[p])
			}
		}
		return nil
	}

	// initial values
	x := utl.LinSpace(0.0, xmax, N)
	y := make([]float64, N)
	//y := make([]float64, Nb)
	for i := 0; i < N; i++ {
		y[i] = 4.0*x[i] - 4.0*x[i]*x[i]
	}

	// debug
	f0 := make([]float64, N)
	//f0 := make([]float64, Nb)
	fcn(f0, 0, y)
	if false {
		io.Pforan("y0 = %v\n", y)
		io.Pforan("f0 = %v\n", f0)
		var J la.Triplet
		jac(&J, 0, y)
		la.PrintMat("J", J.ToMatrix(nil).ToDense(), "%8.3f", false)
	}
	//chk.Panic("stop")

	/*
	   // constraints
	   var A la.Triplet
	   A.Init(2, N, 2)
	   A.Put(0,   0, 1.0)
	   A.Put(1, N-1, 1.0)
	   io.Pfcyan("A = %+v\n", A)
	   Am := A.ToMatrix(nil)
	   c  := make([]float64, 2)
	   la.SpMatVecMul(c, 1, Am, y) // c := Am*y
	   la.PrintMat("A", Am.ToDense(), "%3g", false)
	   io.Pfcyan("c = %v  ([0, 0] => consistent)\n", c)
	*/

	/*
	   // mass matrix
	   var M la.Triplet
	   M.Init(Nb, Nb, N + 4)
	   for i := 0; i < N; i++ {
	       M.Put(i, i, 1.0)
	   }
	   M.PutMatAndMatT(&A)
	   Mm := M.ToMatrix(nil)
	   la.PrintMat("M", Mm.ToDense(), "%3g", false)
	*/

	// output
	var b0, b1, b2 bytes.Buffer
	fmt.Fprintf(&b0, "from gosl import *\n")
	fmt.Fprintf(&b1, "T = array([")
	fmt.Fprintf(&b2, "U = array([")
	out := func(first bool, dt, t float64, y []float64, args ...interface{}) error {
		fmt.Fprintf(&b1, "%23.15E,", t)
		fmt.Fprintf(&b2, "[")
		for i := 0; i < N; i++ {
			fmt.Fprintf(&b2, "%23.15E,", y[i])
		}
		fmt.Fprintf(&b2, "],")
		return nil
	}
	defer func() {
		fmt.Fprintf(&b1, "])\n")
		fmt.Fprintf(&b2, "])\n")
		fmt.Fprintf(&b2, "X = linspace(0.0, %g, %d)\n", xmax, N)
		fmt.Fprintf(&b2, "tt, xx = meshgrid(T, X)\n")
		fmt.Fprintf(&b2, "ax = PlotSurf(tt, xx, vstack(transpose(U)), 't', 'x', 'u', 0.0, 1.0)\n")
		fmt.Fprintf(&b2, "ax.view_init(20.0, 30.0)\n")
		fmt.Fprintf(&b2, "show()\n")
		io.WriteFileD("/tmp/gosl", "plot_diffusion_1d.py", &b0, &b1, &b2)
	}()

	// ode solver
	var Jfcn Cb_jac
	var osol ODE
	if !numjac {
		Jfcn = jac
	}
	osol.Init(method, N, fcn, Jfcn, nil, out, silent)
	//osol.Init(method, Nb, fcn, Jfcn, &M, out, silent)
	osol.SetTol(atol, rtol)

	// constant Jacobian
	if method == "BwEuler" {
		osol.CteTg = true
		osol.Verbose = true
	}

	// run
	wallt0 := time.Now()
	if !fixstp {
		dt = tf - t0
	}
	osol.Solve(y, t0, tf, dt, fixstp)
	io.Pfmag("elapsed time = %v\n", time.Now().Sub(wallt0))
}
