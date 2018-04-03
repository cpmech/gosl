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

func checkCG(tst *testing.T, sol *ConjGrad, nfevalRef, nJevalRef int, fmin, fminRef, tolf, tolx float64, xmin, xminRef []float64) {
	name := "Wolfe"
	if sol.UseBrent {
		name = "Brent"
	}
	io.Pforan("%s: NumIter = %v\n", name, sol.NumIter)
	chk.Int(tst, io.Sf("%s: NumFeval", name), sol.NumFeval, nfevalRef)
	chk.Int(tst, io.Sf("%s: NumJeval", name), sol.NumJeval, nJevalRef)
	chk.Float64(tst, io.Sf("%s: fmin", name), tolf, fmin, fminRef)
	chk.Array(tst, io.Sf("%s: xmin", name), tolx, xmin, xminRef)
	io.Pl()
}

func runCGtest(tst *testing.T, fnkey string, ffcn fun.Sv, Jfcn fun.Vv, x0 la.Vector, fminRef, tolf, tolx float64, xminRef []float64) (sol1, sol2 *ConjGrad, sol3 *Powell) {

	// wrap functions to compute nfeval and nJeval
	nfeval, nJeval := 0, 0
	ffcnWrapped := func(x la.Vector) float64 {
		nfeval++
		return ffcn(x)
	}
	JfcnWrapped := func(g, x la.Vector) {
		nJeval++
		Jfcn(g, x)
	}
	ndim := len(x0)

	// solve using Brent
	nfeval, nJeval = 0, 0
	xmin1 := x0.GetCopy()
	sol1 = NewConjGrad(ndim, ffcnWrapped, JfcnWrapped)
	sol1.UseBrent = true
	sol1.History = true
	fmin1 := sol1.Min(xmin1)
	checkCG(tst, sol1, nfeval, nJeval, fmin1, fminRef, tolf, tolx, xmin1, xminRef)

	// solve again using Wolfe's method
	nfeval, nJeval = 0, 0
	xmin2 := x0.GetCopy()
	sol2 = NewConjGrad(ndim, ffcnWrapped, JfcnWrapped)
	sol2.History = true
	fmin2 := sol2.Min(xmin2)
	checkCG(tst, sol2, nfeval, nJeval, fmin2, fminRef, tolf, tolx, xmin2, xminRef)

	// solve with Powell's method
	nfeval, nJeval = 0, 0
	xmin3 := x0.GetCopy()
	sol3 = NewPowell(ndim, ffcnWrapped)
	sol3.History = true
	fmin3 := sol3.Min(xmin3, false)
	checkPowell(tst, sol3, nfeval, nJeval, fmin3, fminRef, tolf, tolx, xmin3, xminRef)

	// plot
	if chk.Verbose {
		if ndim > 2 {
			plt.Reset(true, &plt.A{WidthPt: 600, Dpi: 150, Prop: 0.8})
			CompareHistory3d("Brent", "Wolfe", sol1.Hist, sol2.Hist, xmin1, xmin2)
			plt.Save("/tmp/gosl/opt", fnkey+"a")

			plt.Reset(true, &plt.A{WidthPt: 600, Dpi: 150, Prop: 0.8})
			CompareHistory3d("CG", "Powell", sol2.Hist, sol3.Hist, xmin2, xmin3)
			plt.Save("/tmp/gosl/opt", fnkey+"b")
		} else {
			plt.Reset(true, &plt.A{WidthPt: 300, Dpi: 150, Prop: 1.5})
			CompareHistory2d("Brent", "Wolfe", sol1.Hist, sol2.Hist, xmin1, xmin2)
			plt.Save("/tmp/gosl/opt", fnkey+"a")

			plt.Reset(true, &plt.A{WidthPt: 300, Dpi: 150, Prop: 1.5})
			CompareHistory2d("CG", "Powell", sol2.Hist, sol3.Hist, xmin2, xmin3)
			plt.Save("/tmp/gosl/opt", fnkey+"b")
		}
		io.Pl()
	}
	return
}

func TestConjGrad01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ConjGrad01. Simple bi-dimensional optimization")

	// function f({x})
	nfeval := 0
	ffcn := func(x la.Vector) float64 {
		nfeval++
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	// Jacobian df/dx
	nJeval := 0
	Jfcn := func(g, x la.Vector) {
		nJeval++
		g[0] = 2.0 * x[0]
		g[1] = 2.0 * x[1]
	}

	// initial point and reference
	fref := -0.5
	xref := []float64{0, 0}
	x0 := la.NewVectorSlice([]float64{1, 1})

	// run test
	runCGtest(tst, "conjgrad01", ffcn, Jfcn, x0, fref, 1e-15, 1e-10, xref)
}

func TestConjGrad02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ConjGrad02. Simple bi-dimensional optimization")

	// function f({x})
	A := la.NewMatrixDeep2([][]float64{
		{3, 1},
		{1, 2},
	})
	tmp := la.NewVector(A.M)
	nfeval := 0
	ffcn := func(x la.Vector) float64 {
		nfeval++
		la.MatVecMul(tmp, 1, A, x)
		return la.VecDot(x, tmp) // xᵀ A x
	}

	// Jacobian df/dx
	At := A.GetTranspose()
	AtPlusA := la.NewMatrix(A.M, A.M)
	la.MatAdd(AtPlusA, 1, At, 1, A)
	nJeval := 0
	Jfcn := func(g, x la.Vector) {
		nJeval++
		la.MatVecMul(g, 1, AtPlusA, x) // g := (Aᵀ+A)⋅x
	}

	// initial point and reference
	fref := 0.0
	xref := []float64{0, 0}
	x0 := la.NewVectorSlice([]float64{1.5, -0.75})

	// run test
	runCGtest(tst, "conjgrad02", ffcn, Jfcn, x0, fref, 1e-15, 1e-9, xref)
}

func TestConjGrad03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ConjGrad03. Simple three-dimensional optimization")

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

	// Jacobian df/dx
	At := A.GetTranspose()
	AtPlusA := la.NewMatrix(A.M, A.M)
	la.MatAdd(AtPlusA, 1, At, 1, A)
	nJeval := 0
	Jfcn := func(g, x la.Vector) {
		nJeval++
		la.MatVecMul(g, 1, AtPlusA, x) // g := (Aᵀ+A)⋅x
	}

	// initial point and reference
	fref := 0.0
	xref := []float64{0, 0, 0}
	x0 := la.NewVectorSlice([]float64{1, 2, 3})

	// run test
	runCGtest(tst, "conjgrad03", ffcn, Jfcn, x0, fref, 1e-15, 1e-9, xref)
}

func TestConjGrad04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ConjGrad04. 3D Rosenbrock function")

	// objective function: Rosenbrock
	N := 5
	p := Factory.RosenbrockMulti(N)

	// initial point => xmin
	x0 := la.NewVectorSlice([]float64{1.3, 0.7, 0.8, 1.9, 1.2})

	// run test
	runCGtest(tst, "conjgrad04", p.Ffcn, p.Gfcn, x0, p.Fref, 1e-15, 1e-8, p.Xref)
}
