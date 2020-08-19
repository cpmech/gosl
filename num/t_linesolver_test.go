// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/la"
	"gosl/plt"
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

func TestLineSolver03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LineSolver03. Min with vertical n")

	nfeval := 0
	ffcn := func(x la.Vector) float64 {
		nfeval++
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	io.Pf(". . .solution is @ halfway in x + n\n")
	x := la.NewVectorSlice([]float64{0.5, -0.5})
	n := la.NewVectorSlice([]float64{0, 1})
	line := NewLineSolver(2, ffcn, nil)
	λmin := line.Min(x, n)
	chk.Int(tst, "NumFeval", line.NumFeval, nfeval)
	chk.Float64(tst, "λroot", 1e-11, λmin, 0.5)
	chk.Float64(tst, "g(λroot)", 1e-11, line.G(λmin), -0.250)

	if chk.Verbose {
		useBracket := true
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		line.PlotC(0, 1, λmin, -1, 1, -1, 1, 41)
		plt.Subplot(2, 1, 2)
		line.PlotG(λmin, 0, 1, 101, useBracket)
		plt.Save("/tmp/gosl/num", "linesolver03")
	}
}

func TestLineSolver04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LineSolver04. Min in two steps")

	ffcn := func(x la.Vector) float64 {
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	io.Pf(". . .solution is @ halfway in x + n\n")
	x := la.NewVectorSlice([]float64{0.5, -0.5})
	n := la.NewVectorSlice([]float64{0, 0.2})
	line := NewLineSolver(2, ffcn, nil)
	λ := line.Min(x, n)
	chk.Float64(tst, "λ", 1e-11, λ, 2.5)
	chk.Float64(tst, "g(λ)", 1e-11, line.G(λ), -0.250)

	useBracket := true
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		line.PlotC(0, 1, λ, -1, 1, -1, 1, 41)
		plt.Subplot(2, 1, 2)
		line.PlotG(λ, -1, 3, 101, useBracket)
	}

	la.VecAdd(x, 1, x, λ, n) // x += λ⋅n
	n[0], n[1] = -0.2, 0.0
	λ = line.Min(x, n)
	chk.Float64(tst, "λ", 1e-11, λ, 2.5)
	chk.Float64(tst, "g(λ)", 1e-11, line.G(λ), -0.5)

	if chk.Verbose {
		plt.Subplot(2, 1, 1)
		x2d := []float64{x[0], x[1]}
		n2d := []float64{n[0], n[1]}
		plt.PlotOne(x2d[0]+λ*n2d[0], x2d[1]+λ*n2d[1], &plt.A{C: "y", M: "o", NoClip: true})
		plt.DrawArrow2d(x2d, n2d, false, 1, nil)
		plt.Subplot(2, 1, 2)
		line.PlotG(λ, -1, 3, 101, useBracket)
		plt.Save("/tmp/gosl/num", "linesolver04")
	}
}
