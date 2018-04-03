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

	// objective function
	problem := opt.Factory.SimpleQuadratic2d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1.5, -0.75})

	// ConjGrad
	xmin1 := x0.GetCopy()
	sol1 := opt.NewConjGrad(len(x0), problem.Ffcn, problem.Gfcn)
	sol1.UseHist = true
	t0 := time.Now()
	fmin1 := sol1.Min(xmin1)
	dt := time.Now().Sub(t0)

	// stat
	io.Pf("ConjGrad: fmin     = %g  (fref = %g)\n", fmin1, problem.Fref)
	io.Pf("ConjGrad: xmin     = %.9f  (xref = %g)\n", xmin1, problem.Xref)
	io.Pf("ConjGrad: NumIter  = %d\n", sol1.NumIter)
	io.Pf("ConjGrad: NumFeval = %d\n", sol1.NumFeval)
	io.Pf("ConjGrad: NumGeval = %d\n", sol1.NumGeval)
	io.Pf("ConjGrad: ElapsedT = %v\n", dt)

	// Powell
	xmin2 := x0.GetCopy()
	sol2 := opt.NewPowell(len(x0), problem.Ffcn)
	sol2.UseHist = true
	reuseUmat := false
	t0 = time.Now()
	fmin2 := sol2.Min(xmin2, reuseUmat)
	dt = time.Now().Sub(t0)

	// stat
	io.Pl()
	io.Pf("Powell: fmin     = %g  (fref = %g)\n", fmin2, problem.Fref)
	io.Pf("Powell: xmin     = %.9f  (xref = %g)\n", xmin2, problem.Xref)
	io.Pf("Powell: NumIter  = %d\n", sol2.NumIter)
	io.Pf("Powell: NumFeval = %d\n", sol2.NumFeval)
	io.Pf("Powell: NumGeval = %d\n", sol2.NumGeval)
	io.Pf("Powell: ElapsedT = %v\n", dt)

	// GradDesc
	xmin3 := x0.GetCopy()
	sol3 := opt.NewGradDesc(len(x0), problem.Ffcn, problem.Gfcn)
	sol3.UseHist = true
	sol3.Alpha = 0.2
	t0 = time.Now()
	fmin3 := sol3.Min(xmin3)
	dt = time.Now().Sub(t0)

	// stat
	io.Pl()
	io.Pf("GradDesc: fmin     = %g  (fref = %g)\n", fmin3, problem.Fref)
	io.Pf("GradDesc: xmin     = %.9f  (xref = %g)\n", xmin3, problem.Xref)
	io.Pf("GradDesc: NumIter  = %d\n", sol3.NumIter)
	io.Pf("GradDesc: NumFeval = %d\n", sol3.NumFeval)
	io.Pf("GradDesc: NumGeval = %d\n", sol3.NumGeval)
	io.Pf("GradDesc: ElapsedT = %v\n", dt)

	// plot ConjGrad vs Powell
	io.Pl()
	plt.Reset(true, &plt.A{WidthPt: 300, Dpi: 150, Prop: 1.5})
	opt.CompareHistory2d("ConjGrad", "Powell", sol1.Hist, sol2.Hist, xmin1, xmin2)
	plt.Save("/tmp/gosl", "opt_comparison01a")

	// plot ConjGrad vs GradDesc
	plt.Reset(true, &plt.A{WidthPt: 300, Dpi: 150, Prop: 1.5})
	opt.CompareHistory2d("ConjGrad", "GradDesc", sol1.Hist, sol3.Hist, xmin1, xmin3)
	plt.Save("/tmp/gosl", "opt_comparison01b")
}
