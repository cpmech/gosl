// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func runPowellTest(tst *testing.T, fnkey string, p *Problem, x0 la.Vector, tolf, tolx, α float64) (sol *Powell) {

	// wrap functions
	nFeval, nGeval := 0, 0
	FfcnWrapped := func(x la.Vector) float64 {
		nFeval++
		return p.Ffcn(x)
	}
	ndim := len(x0)

	// solve using Gradient-Descent
	xmin := x0.GetCopy()
	sol = NewPowell(ndim, FfcnWrapped)
	sol.UseHist = true
	reuseUmat := false
	fmin := sol.Min(xmin, reuseUmat)

	// check
	name := "Powell"
	io.Pforan("%s: NumIter = %v\n", name, sol.NumIter)
	chk.Int(tst, io.Sf("%s: NumFeval", name), sol.NumFeval, nFeval)
	chk.Int(tst, io.Sf("%s: NumGeval", name), sol.NumGeval, nGeval)
	chk.Float64(tst, io.Sf("%s: fmin", name), tolf, fmin, p.Fref)
	chk.Array(tst, io.Sf("%s: xmin", name), tolx, xmin, p.Xref)
	io.Pl()

	// plot
	if chk.Verbose {
		if ndim > 2 {
			plt.Reset(true, &plt.A{WidthPt: 600, Dpi: 150, Prop: 0.8})
			sol.Hist.PlotAll3d("Powell", xmin)
		} else {
			plt.Reset(true, &plt.A{WidthPt: 300, Dpi: 150, Prop: 1.5})
			sol.Hist.PlotAll2d("Powell", xmin)
		}
		plt.Save("/tmp/gosl/opt", fnkey)
		io.Pl()
	}
	return
}

func TestPowell01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell01. Very simple bi-dimensional optimization")

	// problem
	p := Factory.SimpleParaboloid()

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 1})

	// run test
	α := 0.6
	runPowellTest(tst, "powell01", p, x0, 1e-10, 1e-5, α)
}

func TestPowell02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell02. quadratic optimization in 2D")

	// problem
	p := Factory.SimpleQuadratic2d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1.5, -0.75})

	// run test
	α := 0.2
	runPowellTest(tst, "powell02", p, x0, 1e-13, 1e-6, α)
}

func TestPowell03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell02. quadratic optimization in 3D")

	// problem
	p := Factory.SimpleQuadratic3d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 2, 3})

	// run test
	α := 0.1
	runPowellTest(tst, "powell03", p, x0, 1e-10, 1e-5, α)
}
