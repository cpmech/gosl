// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/opt"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// objective function: Rosenbrock
	N := 5 // 5D
	problem := opt.Factory.RosenbrockMulti(N)

	// initial point
	x0 := la.NewVectorSlice([]float64{1.3, 0.7, 0.8, 1.9, 1.2})

	// solve using Brent's method as Line Search
	xmin1 := x0.GetCopy()
	sol1 := opt.NewConjGrad(N, problem.Ffcn, problem.Gfcn)
	sol1.UseBrent = true
	sol1.UseHist = true
	t0 := time.Now()
	fmin1 := sol1.Min(xmin1)
	dt := time.Now().Sub(t0)

	// stat
	io.Pf("Brent: fmin     = %g  (fref = %g)\n", fmin1, problem.Fref)
	io.Pf("Brent: xmin     = %.9f\n", xmin1)
	io.Pf("Brent: NumIter  = %d\n", sol1.NumIter)
	io.Pf("Brent: NumFeval = %d\n", sol1.NumFeval)
	io.Pf("Brent: NumGeval = %d\n", sol1.NumGeval)
	io.Pf("Brent: ElapsedT = %v\n", dt)

	// solve using Wolfe's method as Line Search
	xmin2 := x0.GetCopy()
	sol2 := opt.NewConjGrad(N, problem.Ffcn, problem.Gfcn)
	sol2.UseBrent = false
	sol2.UseHist = true
	t0 = time.Now()
	fmin2 := sol2.Min(xmin2)
	dt = time.Now().Sub(t0)

	// stat
	io.Pl()
	io.Pf("Wolfe: fmin     = %g  (fref = %g)\n", fmin2, problem.Fref)
	io.Pf("Wolfe: xmin     = %.9f\n", xmin2)
	io.Pf("Wolfe: NumIter  = %d\n", sol2.NumIter)
	io.Pf("Wolfe: NumFeval = %d\n", sol2.NumFeval)
	io.Pf("Wolfe: NumGeval = %d\n", sol2.NumGeval)
	io.Pf("Wolfe: ElapsedT = %v\n", dt)

	// compare histories
	io.Pl()
	plt.Reset(true, &plt.A{WidthPt: 600, Dpi: 150, Prop: 0.8})
	opt.CompareHistory3d("Brent", "Wolfe", sol1.Hist, sol2.Hist, xmin1, xmin2)
	plt.Save("/tmp/gosl", "opt_conjgrad01")
}
