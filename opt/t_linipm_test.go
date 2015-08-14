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

	verbose()
	chk.PrintTitle("linipm01")

	// linear programming problem
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
	if true {
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
