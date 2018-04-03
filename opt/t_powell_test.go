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

func checkPowell(tst *testing.T, sol *Powell, nfevalRef, nJevalRef int, fmin, fminRef, tolf, tolx float64, xmin, xminRef []float64) {
	name := "Powell"
	io.Pforan("%s: NumIter = %v\n", name, sol.NumIter)
	chk.Int(tst, io.Sf("%s: NumFeval", name), sol.NumFeval, nfevalRef)
	chk.Float64(tst, io.Sf("%s: fmin", name), tolf, fmin, fminRef)
	chk.Array(tst, io.Sf("%s: xmin", name), tolx, xmin, xminRef)
	io.Pl()
}

func TestPowell01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell01. Simple bi-dimensional optimization")

	// function f({x})
	nfeval := 0
	ffcn := func(x la.Vector) float64 {
		nfeval++
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	// initial point => xmin
	x := la.NewVectorSlice([]float64{1, 1})

	// solver
	solver := NewPowell(2, ffcn)
	solver.History = true
	fmin := solver.Min(x, false)
	chk.Int(tst, "NumFeval", solver.NumFeval, nfeval)
	chk.Float64(tst, "fmin", 1e-15, fmin, -0.5)
	chk.Array(tst, "xmin", 1e-10, x, []float64{0, 0})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})

		plt.Subplot(2, 1, 1)
		solver.Hist.RangeXi = []float64{-1.5, 1.5}
		solver.Hist.RangeXj = []float64{-1.5, 1.5}
		//solver.Hist.PlotX(0, 1, x)

		plt.Subplot(2, 1, 2)
		solver.Hist.PlotF(nil)
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
	nfeval := 0
	ffcn := func(x la.Vector) float64 {
		nfeval++
		la.MatVecMul(tmp, 1, A, x)
		return la.VecDot(x, tmp) // xᵀ A x
	}

	// initial point => xmin
	x := la.NewVector(n)
	x[0] = 1.5
	x[1] = -0.75

	// solver
	solver := NewPowell(n, ffcn)
	solver.History = true
	nfeval = 0
	fmin := solver.Min(x, false)
	chk.Int(tst, "NumFeval", solver.NumFeval, nfeval)
	chk.Float64(tst, "fmin", 1e-15, fmin, 0.0)
	chk.Array(tst, "xmin", 1e-15, x, []float64{0, 0})
	io.Pforan("num f eval = %v\n", solver.NumFeval)
	io.Pforan("iterations = %v\n", solver.NumIter)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})

		plt.Subplot(2, 1, 1)
		solver.Hist.RangeXi = []float64{-1.5, 1.5}
		solver.Hist.RangeXj = []float64{-1.5, 1.5}
		//solver.Hist.PlotX(0, 1, x)

		plt.Subplot(2, 1, 2)
		solver.Hist.PlotF(nil)
		plt.Save("/tmp/gosl/opt", "powell02")
	}
}

func TestPowell03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Powell03. Simple three-dimensional optimization")

	// function f({x})
	A := la.NewMatrixDeep2([][]float64{
		{5, 3, 1},
		{3, 4, 2},
		{1, 2, 3},
	})
	tmp := la.NewVector(A.M)
	nfeval := 0
	ffcn := func(x la.Vector) float64 {
		nfeval++
		la.MatVecMul(tmp, 1, A, x)
		return la.VecDot(x, tmp) // xᵀ A x
	}

	// initial point => xmin
	x := la.NewVectorSlice([]float64{1, 2, 3})

	// solver
	solver := NewPowell(len(x), ffcn)
	solver.History = true
	nfeval = 0
	fmin := solver.Min(x, false)
	chk.Int(tst, "NumFeval", solver.NumFeval, nfeval)
	chk.Float64(tst, "fmin", 1e-15, fmin, 0.0)
	chk.Array(tst, "xmin", 1e-9, x, []float64{0, 0, 0})
	io.Pforan("num f eval = %v\n", solver.NumFeval)
	io.Pforan("iterations = %v\n", solver.NumIter)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 3.0})

		plt.Subplot(4, 1, 1)
		solver.Hist.GapXi = 0.5
		solver.Hist.GapXj = 0.5
		//solver.Hist.PlotX(0, 1, x)

		plt.Subplot(4, 1, 2)
		solver.Hist.GapXi = 0.5
		solver.Hist.GapXj = 0.5
		//solver.Hist.PlotX(1, 2, x)

		plt.Subplot(4, 1, 3)
		solver.Hist.GapXi = 0.5
		solver.Hist.GapXj = 0.5
		//solver.Hist.PlotX(2, 0, x)

		plt.Subplot(4, 1, 4)
		solver.Hist.PlotF(nil)
		plt.Save("/tmp/gosl/opt", "powell03")
	}
}
