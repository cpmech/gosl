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

func checkConjGrad(tst *testing.T, sol *ConjGrad, fmin, fref, tolf, tolx float64, xmin, xref []float64) {
	name := "Wolfe"
	if sol.UseBrent {
		name = "Brent"
	}
	io.Pforan("%s: NumIter = %v\n", name, sol.NumIter)
	io.Pf("%s: NumFeval = %v\n", name, sol.NumFeval)
	io.Pf("%s: NumGeval = %v\n", name, sol.NumGeval)
	chk.Float64(tst, io.Sf("%s: fmin", name), tolf, fmin, fref)
	chk.Array(tst, io.Sf("%s: xmin", name), tolx, xmin, xref)
	io.Pl()
}

func runConjGradTest(tst *testing.T, fnkey string, p *Problem, x0 la.Vector, tolf, tolx float64) (sol1, sol2 *ConjGrad) {

	// solve using Brent
	xmin1 := x0.GetCopy()
	sol1 = NewConjGrad(p)
	sol1.UseBrent = true
	sol1.UseHist = true
	fmin1 := sol1.Min(xmin1, nil)
	checkConjGrad(tst, sol1, fmin1, p.Fref, tolf, tolx, xmin1, p.Xref)

	// solve again using Wolfe's method
	xmin2 := x0.GetCopy()
	sol2 = NewConjGrad(p)
	sol2.UseBrent = false
	sol2.UseHist = true
	fmin2 := sol2.Min(xmin2, nil)
	checkConjGrad(tst, sol2, fmin2, p.Fref, tolf, tolx, xmin2, p.Xref)

	// plot
	if chk.Verbose {
		if p.Ndim > 2 {
			plt.Reset(true, &plt.A{WidthPt: 600, Dpi: 150, Prop: 0.8})
			CompareHistory3d("Brent", "Wolfe", sol1.Hist, sol2.Hist, xmin1, xmin2)
			plt.Save("/tmp/gosl/opt", fnkey)
		} else {
			plt.Reset(true, &plt.A{WidthPt: 300, Dpi: 150, Prop: 1.5})
			CompareHistory2d("Brent", "Wolfe", sol1.Hist, sol2.Hist, xmin1, xmin2)
			plt.Save("/tmp/gosl/opt", fnkey)
		}
		io.Pl()
	}
	return
}

func TestConjGrad01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ConjGrad01. Very simple bi-dimensional optimization")

	// problem
	p := Factory.SimpleParaboloid()

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 1})

	// run test
	runConjGradTest(tst, "conjgrad01", p, x0, 1e-15, 1e-10)
}

func TestConjGrad02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ConjGrad02. quadratic optimization in 2D")

	// problem
	p := Factory.SimpleQuadratic2d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1.5, -0.75})

	// run test
	runConjGradTest(tst, "conjgrad02", p, x0, 1e-15, 1e-9)
}

func TestConjGrad03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ConjGard03. quadratic optimization in 3D")

	// problem
	p := Factory.SimpleQuadratic3d()

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 2, 3})

	// run test
	runConjGradTest(tst, "conjgrad03", p, x0, 1e-15, 1e-8)
}

func TestConjGrad04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ConjGrad04. 3D Rosenbrock function")

	// objective function: Rosenbrock
	N := 5
	p := Factory.RosenbrockMulti(N)

	// initial point
	x0 := la.NewVectorSlice([]float64{1.3, 0.7, 0.8, 1.9, 1.2})

	// run test
	runConjGradTest(tst, "conjgrad04", p, x0, 1e-13, 1e-6)
}
