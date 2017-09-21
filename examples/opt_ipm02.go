// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/opt"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// linear programming problem:
	//
	//   min cᵀx   s.t.   Aᵀx = b, x ≥ 0
	//    x
	//
	// specific problem:
	//
	//   min   2*x0 +   x1
	//   s.t.   -x0 +   x1 ≤ 1
	//           x0 +   x1 ≥ 2   →  -x0 - x1 ≤ -2
	//           x0 - 2*x1 ≤ 4
	//         x1 ≥ 0
	//
	// standard (step 1) add slack
	//   s.t.   -x0 +   x1 + x2           = 1
	//          -x0 -   x1      + x3      = -2
	//           x0 - 2*x1           + x4 = 4
	//
	// standard (step 2)
	//    replace x0 := x0_ - x5
	//    because it's unbounded
	//
	//    min  2*x0_ +   x1                - 2*x5
	//    s.t.  -x0_ +   x1 + x2           +   x5 = 1
	//          -x0_ -   x1      + x3      +   x5 = -2
	//           x0_ - 2*x1           + x4 -   x5 = 4
	//         x0_,x1,x2,x3,x4,x5 ≥ 0

	// coefficients vector
	c := []float64{2, 1, 0, 0, 0, -2}

	// constraints as a sparse matrix
	var T la.Triplet
	T.Init(3, 6, 12) // 3 by 6 matrix, with 12 non-zeros
	T.Put(0, 0, -1)
	T.Put(0, 1, 1)
	T.Put(0, 2, 1)
	T.Put(0, 5, 1)
	T.Put(1, 0, -1)
	T.Put(1, 1, -1)
	T.Put(1, 3, 1)
	T.Put(1, 5, 1)
	T.Put(2, 0, 1)
	T.Put(2, 1, -2)
	T.Put(2, 4, 1)
	T.Put(2, 5, -1)
	Am := T.ToMatrix(nil) // compressed-column matrix

	// right-hand side
	b := []float64{1, -2, 4}

	// solve LP
	var ipm opt.LinIpm
	defer ipm.Free()
	ipm.Init(Am, b, c, nil)
	ipm.Solve(true)

	// print solution
	io.Pl()
	io.Pf("x = %v\n", ipm.X)
	io.Pf("λ = %v\n", ipm.L)
	io.Pf("s = %v\n", ipm.S)

	// check solution
	chk.Verbose = true
	tst := new(testing.T)
	A := Am.ToDense()
	bres := make([]float64, len(b))
	la.MatVecMul(bres, 1, A, ipm.X)
	chk.Array(tst, "A*x=b", 1e-8, bres, b)

	// fix and check x
	x := ipm.X[:2]
	x[0] -= ipm.X[5]
	chk.Array(tst, "x", 1e-8, x, []float64{0.5, 1.5})

	// plotting
	plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
	f := func(x []float64) float64 { return c[0]*x[0] + c[1]*x[1] }
	g := func(x []float64, i int) float64 { return A.Get(i, 0)*x[0] + A.Get(i, 1)*x[1] - b[i] }
	np := 41
	argsF := &plt.A{CmapIdx: 1}
	argsG := &plt.A{Levels: []float64{0}, Colors: []string{"yellow"}, Lw: 2, Fsz: 10}
	vmin, vmax := []float64{-2.0, -2.0}, []float64{2.0, 2.0}
	opt.PlotTwoVarsContour(ipm.X[:2], np, nil, true, vmin, vmax, argsF, argsG, f,
		func(x []float64) float64 { return g(x, 0) },
		func(x []float64) float64 { return g(x, 1) },
	)
	plt.Equal()
	plt.HideAllBorders()
	plt.Gll("$x$", "$y$", &plt.A{LegOut: true})
	plt.Save("/tmp/gosl", "opt_ipm02")
}
