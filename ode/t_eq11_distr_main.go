// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	mpi.Start()
	defer mpi.Stop()

	chk.Verbose = true
	chk.PrintTitle("Hairer-Wanner VII-p2 Eq.(1.1) (Distr=true)")

	lam := -50.0

	fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = lam*y[0] - lam*math.Cos(x)
		return nil
	}

	jac := func(dfdy *la.Triplet, dx, x float64, y la.Vector) error {
		if dfdy.Max() == 0 {
			dfdy.Init(1, 1, 1)
		}
		dfdy.Start()
		dfdy.Put(0, 0, lam)
		return nil
	}

	xa, xb := 0.0, 1.5
	ndim := 1
	y := la.NewVector(ndim)

	conf, err := ode.NewConfig(ode.Radau5kind, "", nil, nil)
	status(err)
	conf.SaveXY = true

	sol := ode.NewSolver(ode.Radau5kind, ndim, fcn, jac, nil, nil)
	sol.SaveXY = true
	sol.Distr = true // <<<<<<< distributed mode

	sol.Solve(y, xa, xb, xb-xa, false)

	tst := new(testing.T)
	chk.Int(tst, "number of F evaluations ", sol.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", sol.Njeval, 1)
	chk.Int(tst, "total number of steps   ", sol.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", sol.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", sol.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", sol.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", sol.Nitmax, 2)
	chk.Int(tst, "IdxSave", sol.IdxSave, sol.Naccepted+1)

	X := utl.LinSpace(xa, xb, 101)
	Y := make([]float64, len(X))
	for i := 0; i < len(X); i++ {
		Y[i] = -lam * (math.Sin(X[i]) - lam*math.Cos(X[i]) + lam*math.Exp(lam*X[i])) / (lam*lam + 1.0)
	}
	e := sol.IdxSave
	Y0 := sol.ExtractTimeSeries(0)
	plt.Reset(false, nil)
	plt.Plot(X, Y, &plt.A{C: "grey", Ls: "-", Lw: 10, L: "solution", NoClip: true})
	plt.Plot(sol.Xvalues[:e], Y0, &plt.A{C: "b", M: "o", Ls: "-", L: "Radau5", NoClip: true})
	plt.Gll("$x$", "$y$", nil)
	plt.Save("/tmp/gosl/ode", "eq11_distr")
}
