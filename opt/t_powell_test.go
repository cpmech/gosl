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

func TestPowell01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell01. Simple bi-dimensional optimization")

	// function f({x})
	ffcn := func(x la.Vector) float64 {
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	// initial point
	x0 := la.NewVectorSlice([]float64{1, 1})

	// solver
	solver := NewPowell(2, ffcn)
	solver.History = true
	fmin, xmin := solver.Min(x0, false)
	fminV2, xminV2 := solver.MinVersion2(x0, false)
	chk.Float64(tst, "fmin", 1e-15, fmin, -0.5)
	chk.Float64(tst, "fminV2", 1e-15, fminV2, -0.5)
	chk.Array(tst, "xmin", 1e-10, xmin, []float64{0, 0})
	chk.Array(tst, "xminV@", 1e-10, xminV2, []float64{0, 0})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})
		solver.Hist.Plot(0, 1, -1.1, 1.1, -1.1, 1.1, 41)
		plt.Save("/tmp/gosl/opt", "powell01")
	}
}

func TestPowell02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell02. Simple bi-dimensional optimization")

	// function f({x})
	n := 2
	A := la.NewMatrix(n, n)
	A.Set(0, 0, 3.0)
	A.Set(0, 1, 1.0)
	A.Set(1, 0, 1.0)
	A.Set(1, 1, 2.0)
	tmp := la.NewVector(n)
	ffcn := func(x la.Vector) float64 {
		la.MatVecMul(tmp, 1, A, x)
		return la.VecDot(x, tmp) // xᵀ A x
	}

	// initial point
	x0 := la.NewVector(n)
	x0[0] = 1.5
	x0[1] = -0.75

	// solver
	solver := NewPowell(n, ffcn)
	solver.History = true
	fmin, xmin := solver.Min(x0, false)
	chk.Float64(tst, "fmin", 1e-15, fmin, 0.0)
	chk.Array(tst, "xmin", 1e-15, xmin, []float64{0, 0})
	io.Pforan("num f eval = %v\n", solver.Nfeval)
	io.Pforan("iterations = %v\n", solver.Niter)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})
		ximin, ximax, xjmin, xjmax := -2.0, 2.0, -2.0, 2.0
		solver.Hist.Plot(0, 1, ximin, ximax, xjmin, xjmax, 41)
		plt.Subplot(2, 1, 1)
		plt.Equal()
		plt.AxisRange(ximin, ximax, xjmin, xjmax)
		plt.Save("/tmp/gosl/opt", "powell02")
	}
}
