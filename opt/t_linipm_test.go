// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func Test_linipm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("linipm01")

	// linear programming problem
	//   min  -4*x0 - 5*x1
	//   s.t.  2*x0 +   x1 ≤ 3
	//           x0 + 2*x1 ≤ 3
	//         x0,x1 ≥ 0
	// standard:
	//         2*x0 +   x1 + x2     = 3
	//           x0 + 2*x1     + x3 = 3
	//         x0,x1,x2,x3 ≥ 0
	var A la.Triplet
	A.Init(2, 4, 6)
	A.Put(0, 0, 2.0)
	A.Put(0, 1, 1.0)
	A.Put(0, 2, 1.0)
	A.Put(1, 0, 1.0)
	A.Put(1, 1, 2.0)
	A.Put(1, 3, 1.0)
	Ad := A.ToMatrix(nil).ToDense()
	c := []float64{-4, -5, 0, 0}
	b := []float64{3, 3}

	// print LP
	la.PrintMat("A", Ad, "%6g", false)
	la.PrintVec("b", b, "%6g", false)
	la.PrintVec("c", c, "%6g", false)
	io.Pf("\n")

	// solve LP
	var ipm LinIpm
	defer ipm.Clean()
	ipm.Init(&A, b, c, nil)
	err := ipm.Solve(chk.Verbose)
	if err != nil {
		tst.Errorf("ipm failed:\n%v", err)
		return
	}

	// check
	io.Pf("\n")
	io.Pforan("x = %v\n", ipm.X)
	io.Pfcyan("λ = %v\n", ipm.L)
	io.Pforan("s = %v\n", ipm.S)
	x := ipm.X[:2]
	bres := make([]float64, 2)
	la.MatVecMul(bres, 1, Ad, x)
	io.Pforan("bres = %v\n", bres)
	chk.Vector(tst, "x", 1e-9, x, []float64{1, 1})
	chk.Vector(tst, "A*x=b", 1e-8, bres, b)

	// plot
	if true && chk.Verbose {
		f := func(x []float64) float64 {
			return c[0]*x[0] + c[1]*x[1]
		}
		g := func(x []float64, i int) float64 {
			return Ad[i][0]*x[0] + Ad[i][1]*x[1] - b[i]
		}
		np := 41
		vmin, vmax := []float64{-2.0, -2.0}, []float64{2.0, 2.0}
		PlotTwoVarsContour("/tmp/gosl", "test_linipm01", x, np, nil, true, vmin, vmax, f,
			func(x []float64) float64 { return g(x, 0) },
			func(x []float64) float64 { return g(x, 1) },
		)
	}
}

func Test_linipm02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("linipm02")

	// linear program
	//   min   2*x0 +   x1
	//   s.t.   -x0 +   x1 ≤ 1
	//           x0 +   x1 ≥ 2   →  -x0 - x1 ≤ -2
	//           x0 - 2*x1 ≤ 4
	//         x1 ≥ 0
	// standard (step 1) add slack
	//   s.t.   -x0 +   x1 + x2           = 1
	//          -x0 -   x1      + x3      = -2
	//           x0 - 2*x1           + x4 = 4
	// standard (step 2)
	//    replace x0 := x0_ - x5
	//    because it's unbounded
	//    min  2*x0_ +   x1                - 2*x5
	//    s.t.  -x0_ +   x1 + x2           +   x5 = 1
	//          -x0_ -   x1      + x3      +   x5 = -2
	//           x0_ - 2*x1           + x4 -   x5 = 4
	//         x0_,x1,x2,x3,x4,x5 ≥ 0
	var A la.Triplet
	A.Init(3, 6, 12)
	A.Put(0, 0, -1)
	A.Put(0, 1, 1)
	A.Put(0, 2, 1)
	A.Put(0, 5, 1)
	A.Put(1, 0, -1)
	A.Put(1, 1, -1)
	A.Put(1, 3, 1)
	A.Put(1, 5, 1)
	A.Put(2, 0, 1)
	A.Put(2, 1, -2)
	A.Put(2, 4, 1)
	A.Put(2, 5, -1)
	Ad := A.ToMatrix(nil).ToDense()
	c := []float64{2, 1, 0, 0, 0, -2}
	b := []float64{1, -2, 4}

	// print LP
	la.PrintMat("A", Ad, "%6g", false)
	la.PrintVec("b", b, "%6g", false)
	la.PrintVec("c", c, "%6g", false)
	io.Pf("\n")

	// solve LP
	var ipm LinIpm
	defer ipm.Clean()
	ipm.Init(&A, b, c, nil)
	err := ipm.Solve(chk.Verbose)
	if err != nil {
		tst.Errorf("ipm failed:\n%v", err)
		return
	}

	// check
	io.Pf("\n")
	bres := make([]float64, len(b))
	la.MatVecMul(bres, 1, Ad, ipm.X)
	io.Pforan("bres = %v\n", bres)
	chk.Vector(tst, "A*x=b", 1e-8, bres, b)

	// fix and check x
	x := ipm.X[:2]
	x[0] -= ipm.X[5]
	io.Pforan("x = %v\n", x)
	chk.Vector(tst, "x", 1e-8, x, []float64{0.5, 1.5})

	// plot
	if true && chk.Verbose {
		f := func(x []float64) float64 {
			return c[0]*x[0] + c[1]*x[1]
		}
		g := func(x []float64, i int) float64 {
			return Ad[i][0]*x[0] + Ad[i][1]*x[1] - b[i]
		}
		np := 41
		vmin, vmax := []float64{-2.0, -2.0}, []float64{2.0, 2.0}
		PlotTwoVarsContour("/tmp/gosl", "test_linipm02", x, np, nil, true, vmin, vmax, f,
			func(x []float64) float64 { return g(x, 0) },
			func(x []float64) float64 { return g(x, 1) },
		)
	}
}
