// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

type yanaType func(x float64) float64

// eq11data returns equations for Hairer-Wanner VII-p2 Eq.(1.1)
func eq11data() (dx, xf float64, y la.Vector, yana yanaType, fcn Func, jac JacF) {

	λ := -50.0
	dx = 1.875 / 50.0
	xf = 1.5
	y = la.Vector([]float64{0.0})

	yana = func(x float64) float64 {
		return -λ * (math.Sin(x) - λ*math.Cos(x) + λ*math.Exp(λ*x)) / (λ*λ + 1.0)
	}

	fcn = func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = λ*y[0] - λ*math.Cos(x)
		return nil
	}

	jac = func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(1, 1, 1)
		}
		dfdy.Start()
		dfdy.Put(0, 0, λ)
		return nil
	}
	return
}

func eq11plotOne(fnk, label string, xf float64, yana yanaType, sol *Solver) {
	X := utl.LinSpace(0, xf, 101)
	Y := make([]float64, len(X))
	for i := 0; i < len(X); i++ {
		Y[i] = yana(X[i])
	}
	a := sol.Out.IdxSave
	Xa := sol.Out.Xvalues[:a]
	Ya := sol.Out.ExtractTimeSeries(0)
	plt.Reset(true, nil)
	plt.Plot(X, Y, &plt.A{C: "grey", Ls: "-", Lw: 7, L: "ana", NoClip: true})
	plt.Plot(Xa, Ya, &plt.A{C: "r", M: ".", Ls: ":", L: label, NoClip: true})
	plt.Gll("$x$", "$y$", nil)
	plt.Save("/tmp/gosl/ode", fnk)
}
