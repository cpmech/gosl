// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
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
		f := func(x []float64) float64 { return c[0]*x[0] + c[1]*x[1] }
		g := func(x []float64) (res []float64) {
			res = make([]float64, len(b))
			la.MatVecMul(res, 1, Ad, x)
			for i := 0; i < len(b); i++ {
				res[i] -= b[i]
			}
			return
		}
		np := 41
		vmin, vmax := -2.0, 2.0
		xx, yy := utl.MeshGrid2D(vmin, vmax, vmin, vmax, np, np)
		zz := la.MatAlloc(np, np)
		nl := len(b)
		ww := utl.Deep3alloc(nl, np, np)
		for i := 0; i < np; i++ {
			for j := 0; j < np; j++ {
				xtmp := []float64{xx[i][j], yy[i][j]}
				zz[i][j] = f(xtmp)
				res := g(xtmp)
				for k := 0; k < nl; k++ {
					ww[k][i][j] = res[k]
				}
			}
		}
		plt.SetForEps(0.8, 300)
		plt.Contour(xx, yy, zz, "")
		for i := 0; i < nl; i++ {
			plt.ContourSimple(xx, yy, ww[i], "zorder=5, levels=[0], colors=['yellow'], linewidths=[2], clip_on=0")
		}
		plt.PlotOne(x[0], x[1], "'r*',label='optimum', zorder=10")
		plt.Gll("$x$", "$y$", "leg_out=1")
		plt.Cross("clr='grey'")
		plt.SaveD("/tmp/gosl", "test_linipm01.eps")
	}
}
