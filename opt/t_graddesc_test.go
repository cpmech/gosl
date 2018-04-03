// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func checkGradDesc(tst *testing.T, sol *GradDesc, nFevalRef, nGevalRef int, fmin, fminRef, tolf, tolx float64, xmin, xminRef []float64) {
	name := "GradDesc"
	io.Pforan("%s: NumIter = %v\n", name, sol.NumIter)
	chk.Int(tst, io.Sf("%s: NumFeval", name), sol.NumFeval, nFevalRef)
	chk.Int(tst, io.Sf("%s: NumJeval", name), sol.NumGeval, nGevalRef)
	chk.Float64(tst, io.Sf("%s: fmin", name), tolf, fmin, fminRef)
	chk.Array(tst, io.Sf("%s: xmin", name), tolx, xmin, xminRef)
	io.Pl()
}

func runGradDesctest(tst *testing.T, fnkey string, ffcn fun.Sv, Jfcn fun.Vv, x0 la.Vector, α, fminRef, tolf, tolx float64, xminRef []float64) (sol *GradDesc) {

	// wrap functions to compute nfeval and nJeval
	nFeval, nGeval := 0, 0
	FfcnWrapped := func(x la.Vector) float64 {
		nFeval++
		return ffcn(x)
	}
	GfcnWrapped := func(g, x la.Vector) {
		nGeval++
		Jfcn(g, x)
	}
	ndim := len(x0)

	// solve using Gradient-Descent
	nFeval, nGeval = 0, 0
	xmin := x0.GetCopy()
	sol = NewGradDesc(ndim, FfcnWrapped, GfcnWrapped)
	sol.Alpha = α
	sol.UseHist = true
	fmin := sol.Min(xmin)
	checkGradDesc(tst, sol, nFeval, nGeval, fmin, fminRef, tolf, tolx, xmin, xminRef)

	// plot
	if chk.Verbose {
		if ndim > 2 {
			plt.Reset(true, &plt.A{WidthPt: 600, Dpi: 150, Prop: 0.8})
			sol.Hist.PlotAll3d("GradDesc", xmin)
		} else {
			plt.Reset(true, &plt.A{WidthPt: 300, Dpi: 150, Prop: 1.5})
			sol.Hist.PlotAll2d("GradDesc", xmin)
		}
		plt.Save("/tmp/gosl/opt", fnkey)
		io.Pl()
	}
	return
}

func TestGradDesc01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GradDesc01. Very simple bi-dimensional optimization")

	// problem
	Ffcn, Gfcn, _, fref, xref := FactoryObjectives.SimpleParaboloid()

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 1})

	// run test
	α := 0.6
	runGradDesctest(tst, "graddesc01", Ffcn, Gfcn, x0, α, fref, 1e-10, 1e-5, xref)
}

func TestGradDesc02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GradDesc02. quadratic optimization in 2D")

	// problem
	Ffcn, Gfcn, _, fref, xref := FactoryObjectives.SimpleQuadratic2d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1.5, -0.75})

	// run test
	α := 0.2
	runGradDesctest(tst, "graddesc02", Ffcn, Gfcn, x0, α, fref, 1e-13, 1e-6, xref)
}

func TestGradDesc03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GradDesc02. quadratic optimization in 3D")

	// problem
	Ffcn, Gfcn, _, fref, xref := FactoryObjectives.SimpleQuadratic3d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 2, 3})

	// run test
	α := 0.1
	runGradDesctest(tst, "graddesc03", Ffcn, Gfcn, x0, α, fref, 1e-10, 1e-5, xref)
}
