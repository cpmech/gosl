// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"
	"os"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// title
	io.Pf("Hairer-Wanner VII-p2 Eq.(1.1)")

	// constants
	λ := -50.0
	dx := 1.875 / 50.0
	xf := 1.5
	atol := 1e-4
	rtol := 1e-4
	numJac := false
	saveStep := true
	saveDense := false

	// ode function: f(x,y) = dy/dx
	fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = λ*y[0] - λ*math.Cos(x)
		return nil
	}

	// Jacobian function: J = df/dy
	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(1, 1, 1)
		}
		dfdy.Start()
		dfdy.Put(0, 0, λ)
		return nil
	}

	// initial values
	ndim := 1
	y := la.NewVector(ndim)
	y[0] = 0.0

	// FwEuler
	io.Pf("\n------------ Forward-Euler ------------------\n")
	fixedStep := true
	stat1, out1, err := ode.Solve("fweuler", fcn, jac, y.GetCopy(), xf, dx, atol, rtol, numJac, fixedStep, saveStep, saveDense)
	stat1.Print(false)
	status(err)

	// BwEuler
	io.Pf("\n------------ Backward-Euler ------------------\n")
	fixedStep = true
	stat2, out2, err := ode.Solve("bweuler", fcn, jac, y.GetCopy(), xf, dx, atol, rtol, numJac, fixedStep, saveStep, saveDense)
	stat2.Print(false)
	status(err)

	// MoEuler
	io.Pf("\n------------ Modified-Euler ------------------\n")
	fixedStep = false
	stat3, out3, err := ode.Solve("moeuler", fcn, jac, y.GetCopy(), xf, dx, atol, rtol, numJac, fixedStep, saveStep, saveDense)
	stat3.Print(true)
	status(err)

	// DoPri5
	io.Pf("\n------------ Dormand-Prince5 -----------------\n")
	fixedStep = false
	stat4, out4, err := ode.Solve("dopri5", fcn, jac, y.GetCopy(), xf, dx, atol, rtol, numJac, fixedStep, saveStep, saveDense)
	stat4.Print(true)
	status(err)

	// DoPri8
	io.Pf("\n------------ Dormand-Prince8 -----------------\n")
	fixedStep = false
	stat5, out5, err := ode.Solve("dopri8", fcn, jac, y.GetCopy(), xf, dx, atol, rtol, numJac, fixedStep, saveStep, saveDense)
	stat5.Print(true)
	status(err)

	// Radau5
	io.Pf("\n------------ Radau5 --------------------------\n")
	fixedStep = false
	stat6, out6, err := ode.Solve("radau5", fcn, jac, y.GetCopy(), xf, dx, atol, rtol, numJac, fixedStep, saveStep, saveDense)
	stat6.Print(true)
	status(err)

	// analytical solution
	npts := 201
	xx := utl.LinSpace(0, xf, npts)
	yy := utl.GetMapped(xx, func(x float64) float64 {
		return -λ * (math.Sin(x) - λ*math.Cos(x) + λ*math.Exp(λ*x)) / (λ*λ + 1.0)
	})

	// plot
	plt.Reset(true, &plt.A{WidthPt: 500, Prop: 1.7})

	plt.Subplot(3, 1, 1)
	plt.Plot(xx, yy, &plt.A{C: "grey", Ls: "-", Lw: 5, L: "analytical", NoClip: true})
	plt.Plot(out1.GetStepX(), out1.GetStepY(0), &plt.A{L: "fweuler", C: "k", M: ".", Ls: ":"})
	plt.Plot(out2.GetStepX(), out2.GetStepY(0), &plt.A{L: "bweuler", C: "r", M: ".", Ls: "--"})
	plt.Gll("$x$", "$y$", nil)

	plt.Subplot(3, 1, 2)
	plt.Plot(xx, yy, &plt.A{C: "grey", Ls: "-", Lw: 5, L: "analytical", NoClip: true})
	plt.Plot(out3.GetStepX(), out3.GetStepY(0), &plt.A{L: "moeuler", C: "k", M: ".", Ls: ":"})
	plt.Plot(out4.GetStepX(), out4.GetStepY(0), &plt.A{L: "dopri5", C: "r", M: ".", Ls: "none"})
	plt.Gll("$x$", "$y$", nil)

	plt.Subplot(3, 1, 3)
	plt.Plot(xx, yy, &plt.A{C: "grey", Ls: "-", Lw: 5, L: "analytical", NoClip: true})
	plt.Plot(out5.GetStepX(), out5.GetStepY(0), &plt.A{L: "dopri8", C: "k", M: "o", Ls: ":"})
	plt.Plot(out6.GetStepX(), out6.GetStepY(0), &plt.A{L: "radau5", C: "r", M: "o", Ls: "-", Void: true})
	plt.Gll("$x$", "$y$", nil)

	plt.Save("/tmp/gosl", "ode_hweq11")
}

func status(err error) {
	if err != nil {
		io.Pf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
