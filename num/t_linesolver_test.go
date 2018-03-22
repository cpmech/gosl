// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
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

	io.Pf(". . .solution is @ halfway in x + n\n")
	x := la.NewVectorSlice([]float64{0, 0})
	n := la.NewVectorSlice([]float64{1, 1})
	line := NewLineSolver(2, ffcn, Jfcn)
	λroot := line.Root(x, n)
	chk.Float64(tst, "λroot", 1e-11, λroot, 0.5)
	chk.Float64(tst, "g(λroot)", 1e-11, line.G(λroot), 0.0)

	if chk.Verbose {
		useBracket := true
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		line.PlotC(0, 1, λroot, -1, 1, -1, 1, 41)
		plt.Subplot(2, 1, 2)
		line.PlotG(λroot, 0, 1, 101, useBracket)
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

	io.Pf(". . .solution is @ halfway in x + n\n")
	x := la.NewVectorSlice([]float64{-1, -1})
	n := la.NewVectorSlice([]float64{2, 2})
	line := NewLineSolver(2, ffcn, Jfcn)
	λmin := line.Min(x, n)
	chk.Float64(tst, "λmin", 1e-14, λmin, 0.5)
	chk.Float64(tst, "g(λmin)", 1e-14, line.G(λmin), -0.5)

	io.Pf("\n. . .solution is @ x + λ\n")
	x[0], x[1] = -0.5, -0.5
	n[0], n[1] = 0.5, 0.5
	λmin = line.Min(x, n)
	chk.Float64(tst, "λmin", 1e-14, λmin, 1.0)
	chk.Float64(tst, "g(λmin)", 1e-14, line.G(λmin), -0.5)

	io.Pf("\n. . .solution is @ x\n")
	x[0], x[1] = 0.0, 0.0
	n[0], n[1] = 0.5, 0.5
	λmin = line.Min(x, n)
	chk.Float64(tst, "λmin", 1e-8, λmin, 0.0)
	chk.Float64(tst, "g(λmin)", 1e-14, line.G(λmin), -0.5)

	io.Pf("\n. . .negative lambda\n")
	x[0], x[1] = 0.5, 0.5
	n[0], n[1] = 0.5, 0.5
	λmin = line.Min(x, n)
	chk.Float64(tst, "λmin", 1e-8, λmin, -1.0)
	chk.Float64(tst, "g(λmin)", 1e-14, line.G(λmin), -0.5)

	if chk.Verbose {
		useBracket := true
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		line.PlotC(0, 1, λmin, -1, 1, -1, 1, 41)
		plt.Subplot(2, 1, 2)
		line.PlotG(λmin, 0, 1, 101, useBracket)
		plt.Save("/tmp/gosl/num", "linesolver02")
	}
}
