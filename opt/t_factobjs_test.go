// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestFactObjs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FactObjs01. Standard Rosenbrock function")

	ffcn, Jfcn, ndim, hasSol, xmin, fmin := FactoryObjectives.Rosenbrock(1, 100)

	if ndim != 2 {
		tst.Errorf("ndim should be equal to 2\n")
		return
	}
	if !hasSol {
		tst.Errorf("hasSol should be true\n")
		return
	}
	chk.Array(tst, "xmin", 1e-15, xmin, utl.Ones(ndim))
	chk.Float64(tst, "fmin", 1e-15, fmin, 0.0)

	x := la.NewVectorSlice([]float64{-3, -4})
	gAna := la.NewVector(ndim)
	Jfcn(gAna, x)
	chk.DerivScaVec(tst, "Jfcn", 1e-7, gAna, x, 1e-3, chk.Verbose, func(xx []float64) float64 { return ffcn(xx) })

	// plot
	if chk.Verbose {
		xvec := la.NewVector(2)
		xx, yy, zz := utl.MeshGrid2dF(-2.0, 2.0, -0.5, 3.0, 101, 101, func(r, s float64) float64 {
			xvec[0], xvec[1] = r, s
			return ffcn(xvec)
		})
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.ContourF(xx, yy, zz, &plt.A{Nlevels: 200, NoLines: true, NoLabels: true})
		plt.PlotOne(xmin[0], xmin[1], &plt.A{C: "y", M: "o"})
		plt.Gll("$x$", "$y$", nil)
		plt.Save("/tmp/gosl/opt", "factobjs01")
	}
}

func TestFactObjs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FactObjs02. Multi-dimensional Rosenbrock function")

	N := 5
	ffcn, Jfcn, ndim, hasSol, xmin, fmin := FactoryObjectives.RosenbrockMulti(N)

	if ndim != N {
		tst.Errorf("ndim should be equal to %d\n", N)
		return
	}
	if !hasSol {
		tst.Errorf("hasSol should be true\n")
		return
	}
	chk.Array(tst, "xmin", 1e-15, xmin, utl.Ones(ndim))
	chk.Float64(tst, "fmin", 1e-15, fmin, 0.0)

	x := la.NewVectorSlice([]float64{-1, -2, -3, -4, -5})
	gAna := la.NewVector(ndim)
	Jfcn(gAna, x)
	chk.DerivScaVec(tst, "Jfcn", 1e-6, gAna, x, 1e-3, chk.Verbose, func(xx []float64) float64 { return ffcn(xx) })
}
