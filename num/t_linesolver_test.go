// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestLineSolver01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LineSolver01. Root")

	ffcn := func(x la.Vector) float64 {
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	Jfcn := func(dfdx, x la.Vector) {
		dfdx[0] = 2.0 * x[0]
		dfdx[1] = 2.0 * x[1]
	}

	x := la.NewVectorSlice([]float64{0, 0})
	n := la.NewVectorSlice([]float64{1, 1})

	line := NewLineSolver(2, ffcn, Jfcn)
	λroot := line.Root(x, n)
	chk.Float64(tst, "λroot", 1e-11, λroot, 0.5)

	if chk.Verbose {

		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})

		// contour
		plt.Subplot(2, 1, 1)
		nx, ny := 41, 41
		xmin, xmax, ymin, ymax := -1.0, +1.0, -1.0, +1.0
		xvec := la.NewVector(2)
		xx, yy, zz := utl.MeshGrid2dF(xmin, xmax, ymin, ymax, nx, ny, func(u, v float64) float64 {
			xvec[0], xvec[1] = u, v
			return ffcn(xvec)
		})
		plt.ContourF(xx, yy, zz, nil)
		plt.ContourL(xx, yy, zz, &plt.A{Colors: []string{"y"}, Levels: []float64{0}, Lw: 3})
		plt.PlotOne(x[0]+λroot*n[0], x[1]+λroot*n[1], &plt.A{C: "y", M: "o", NoClip: true})
		plt.DrawArrow2d(x, n, false, 1, nil)
		plt.Equal()
		plt.Gll("$x_0$", "$x_1$", nil)
		plt.HideTRborders()

		// scalar function along n
		plt.Subplot(2, 1, 2)
		line.Set(x, n)
		ll := utl.LinSpace(0, 1, 101)
		gg := utl.GetMapped(ll, line.G)
		plt.Plot(ll, gg, &plt.A{C: plt.C(0, 0), L: "$g(\\lambda)$", NoClip: true})
		plt.PlotOne(λroot, line.G(λroot), &plt.A{C: "y", M: "o", NoClip: true})
		plt.Cross(0, 0, nil)
		plt.Gll("$\\lambda$", "$g(\\lambda)$", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "linesolver01")
	}
}

func TestLineSolver02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LineSolver02. Min")

	ffcn := func(x la.Vector) float64 {
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	Jfcn := func(dfdx, x la.Vector) {
		dfdx[0] = 2.0 * x[0]
		dfdx[1] = 2.0 * x[1]
	}

	x := la.NewVectorSlice([]float64{-1, -1})
	n := la.NewVectorSlice([]float64{2, 2})

	line := NewLineSolver(2, ffcn, Jfcn)
	λmin := line.Min(x, n)
	chk.Float64(tst, "λmin", 1e-11, λmin, 0.5)

	if chk.Verbose {

		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})

		// contour
		plt.Subplot(2, 1, 1)
		nx, ny := 41, 41
		xmin, xmax, ymin, ymax := -1.0, +1.0, -1.0, +1.0
		xvec := la.NewVector(2)
		xx, yy, zz := utl.MeshGrid2dF(xmin, xmax, ymin, ymax, nx, ny, func(u, v float64) float64 {
			xvec[0], xvec[1] = u, v
			return ffcn(xvec)
		})
		plt.ContourF(xx, yy, zz, nil)
		plt.ContourL(xx, yy, zz, &plt.A{Colors: []string{"y"}, Levels: []float64{0}, Lw: 3})
		plt.PlotOne(x[0]+λmin*n[0], x[1]+λmin*n[1], &plt.A{C: "y", M: "o", NoClip: true})
		plt.DrawArrow2d(x, n, false, 1, nil)
		plt.Equal()
		plt.AxisRange(-1, 1, -1, 1)
		plt.Gll("$x_0$", "$x_1$", nil)
		plt.HideTRborders()

		// scalar function along n
		plt.Subplot(2, 1, 2)
		line.Set(x, n)
		ll := utl.LinSpace(0, 1, 101)
		gg := utl.GetMapped(ll, line.G)
		plt.Plot(ll, gg, &plt.A{C: plt.C(0, 0), L: "$g(\\lambda)$", NoClip: true})
		plt.PlotOne(λmin, line.G(λmin), &plt.A{C: "y", M: "o", NoClip: true})
		plt.Cross(0, 0, nil)
		plt.Gll("$\\lambda$", "$g(\\lambda)$", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "linesolver02")
	}
}
