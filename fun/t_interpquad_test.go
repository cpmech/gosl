// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/plt"
	"gosl/utl"
)

func TestInterpQuad01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("InterpQuad01. Interp with quad poly using 3 points")

	// test set
	ycor := func(x float64) float64 { return 1.0 + Pow2(x-1.0) }
	dcor := func(x float64) float64 { return 2.0 * (x - 1.0) }
	x0, y0 := 0.0, 2.0
	x1, y1 := 2.0, 2.0
	x2, y2 := 3.0, 5.0

	// intepolator
	interp := NewInterpQuad()
	interp.Fit3points(x0, y0, x1, y1, x2, y2)

	// check model and derivatives
	for _, x := range []float64{-10, 0, 1, 8} {
		y := interp.F(x)
		yd := interp.G(x)
		chk.Float64(tst, io.Sf("y(%g)", x), 1e-15, y, ycor(x))
		chk.Float64(tst, io.Sf("y'(%g)", x), 1e-15, yd, dcor(x))
	}

	// check optimum
	xopt, fopt := interp.Optimum()
	chk.Float64(tst, "xopt", 1e-15, xopt, 1.0)
	chk.Float64(tst, "fopt", 1e-15, fopt, 1.0)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		xx := utl.LinSpace(-2, 4, 21)
		yy := utl.GetMapped(xx, func(x float64) float64 { return interp.F(x) })
		plt.Plot(xx, yy, &plt.A{C: plt.C(0, 0), M: ".", Ls: "-", NoClip: true})
		plt.PlotOne(x0, y0, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.PlotOne(x1, y1, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.PlotOne(x2, y2, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.Gll("$x$", "$y$", nil)
		plt.Cross(0, 0, nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "interpquad01")
	}
}

func TestInterpQuad02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("InterpQuad02. Interp with quad poly using 3 points (flipped)")

	// test set (flipped compared to previous test)
	ycor := func(x float64) float64 { return 3.0 - Pow2(x-1.0) }
	dcor := func(x float64) float64 { return -2.0 * (x - 1.0) }
	x0, y0 := 0.0, 2.0
	x1, y1 := 2.0, 2.0
	x2, y2 := 3.0, -1.0

	// intepolator
	interp := NewInterpQuad()
	interp.Fit3points(x0, y0, x1, y1, x2, y2)

	// check model and derivatives
	for _, x := range []float64{-10, 0, 1, 8} {
		y := interp.F(x)
		yd := interp.G(x)
		chk.Float64(tst, io.Sf("y(%g)", x), 1e-15, y, ycor(x))
		chk.Float64(tst, io.Sf("y'(%g)", x), 1e-15, yd, dcor(x))
	}

	// check optimum
	xopt, fopt := interp.Optimum()
	chk.Float64(tst, "xopt", 1e-15, xopt, 1.0)
	chk.Float64(tst, "fopt", 1e-15, fopt, 3.0)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		xx := utl.LinSpace(-2, 4, 21)
		yy := utl.GetMapped(xx, func(x float64) float64 { return interp.F(x) })
		plt.Plot(xx, yy, &plt.A{C: plt.C(0, 0), M: ".", Ls: "-", NoClip: true})
		plt.PlotOne(x0, y0, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.PlotOne(x1, y1, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.PlotOne(x2, y2, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.Gll("$x$", "$y$", nil)
		plt.Cross(0, 0, nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "interpquad02")
	}
}

func TestInterpQuad03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("InterpQuad03. Interp with quad poly using 2 points and deriv")

	// test set
	ycor := func(x float64) float64 { return 1.0 + Pow2(x-1.0) }
	dcor := func(x float64) float64 { return 2.0 * (x - 1.0) }
	x0, y0 := 0.0, 2.0
	x1, y1 := 2.0, 2.0
	x2, d2 := -1.0, -4.0

	// intepolator
	interp := NewInterpQuad()
	interp.Fit2pointsD(x0, y0, x1, y1, x2, d2)

	// check model and derivatives
	for _, x := range []float64{-10, 0, 1, 8} {
		y := interp.F(x)
		yd := interp.G(x)
		chk.Float64(tst, io.Sf("y(%g)", x), 1e-15, y, ycor(x))
		chk.Float64(tst, io.Sf("y'(%g)", x), 1e-15, yd, dcor(x))
	}

	// check optimum
	xopt, fopt := interp.Optimum()
	chk.Float64(tst, "xopt", 1e-15, xopt, 1.0)
	chk.Float64(tst, "fopt", 1e-15, fopt, 1.0)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		xx := utl.LinSpace(-2, 4, 21)
		yy := utl.GetMapped(xx, func(x float64) float64 { return interp.F(x) })
		plt.Plot(xx, yy, &plt.A{C: plt.C(0, 0), M: ".", Ls: "-", NoClip: true})
		plt.PlotOne(x0, y0, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.PlotOne(x1, y1, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.PlotOne(x2, ycor(x2), &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.SlopeInd(d2, x2, ycor(x2), 1.0, io.Sf("%g", d2), true, false, false, nil, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.Cross(0, 0, nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "interpquad03")
	}
}

func TestInterpQuad04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("InterpQuad04. Interp with quad poly using 3 points (flipped)")

	// test set (flipped compared to previous test)
	ycor := func(x float64) float64 { return 3.0 - Pow2(x-1.0) }
	dcor := func(x float64) float64 { return -2.0 * (x - 1.0) }
	x0, y0 := 0.0, 2.0
	x1, y1 := 2.0, 2.0
	x2, d2 := -1.0, +4.0

	// intepolator
	interp := NewInterpQuad()
	interp.Fit2pointsD(x0, y0, x1, y1, x2, d2)

	// check model and derivatives
	for _, x := range []float64{-10, 0, 1, 8} {
		y := interp.F(x)
		yd := interp.G(x)
		chk.Float64(tst, io.Sf("y(%g)", x), 1e-15, y, ycor(x))
		chk.Float64(tst, io.Sf("y'(%g)", x), 1e-15, yd, dcor(x))
	}

	// check optimum
	xopt, fopt := interp.Optimum()
	chk.Float64(tst, "xopt", 1e-15, xopt, 1.0)
	chk.Float64(tst, "fopt", 1e-15, fopt, 3.0)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		xx := utl.LinSpace(-2, 4, 21)
		yy := utl.GetMapped(xx, func(x float64) float64 { return interp.F(x) })
		plt.Plot(xx, yy, &plt.A{C: plt.C(0, 0), M: ".", Ls: "-", NoClip: true})
		plt.PlotOne(x0, y0, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.PlotOne(x1, y1, &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.PlotOne(x2, ycor(x2), &plt.A{C: plt.C(2, 0), M: "o", NoClip: true})
		plt.SlopeInd(d2, x2, ycor(x2), 1.0, io.Sf("%g", d2), true, false, false, nil, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.Cross(0, 0, nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "interpquad04")
	}
}
