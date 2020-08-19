// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/plt"
	"gosl/utl"
)

func checkBracket(tst *testing.T, a, b, c, fa, fb, fc float64) {
	if a >= b {
		tst.Errorf("a=%g should be smaller than b=%g\n", a, b)
	}
	if c <= b {
		tst.Errorf("c=%g should be greater than b=%g\n", c, b)
	}
	if fa < fb {
		tst.Errorf("fa=%g should be greater than or equal to fb=%g\n", fa, fb)
	}
	if fc < fb {
		tst.Errorf("fc=%g should be greater than or equal to fb=%g\n", fc, fb)
	}
}

func TestBracket01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Bracket01. quadratic polynomial")

	// function
	ffcn := func(x float64) float64 { return x*x - 1 }

	// bracket minimum
	bracket := NewBracket(ffcn)
	bracket.Verbose = true

	// check
	a, b, c, fa, fb, fc := bracket.Min(0, 1)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-10, 10)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(1, 0)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(10, -10)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-2, -1)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-1, -2)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// plot
	if chk.Verbose {
		xx := utl.LinSpace(-1.2, 1.2, 101)
		yy := utl.GetMapped(xx, func(x float64) float64 { return ffcn(x) })
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Cross(0, 0, nil)
		plt.Plot(xx, yy, &plt.A{C: "b", Ls: "-", L: "curve1", NoClip: true})
		plt.PlotOne(a, fa, &plt.A{C: "r", M: "|", Ms: 40, NoClip: true})
		plt.PlotOne(b, fb, &plt.A{C: "r", M: "*", Ms: 10, NoClip: true})
		plt.PlotOne(c, fc, &plt.A{C: "r", M: "|", Ms: 40, NoClip: true})
		plt.Text(a-0.2, fa-0.2, io.Sf("a=(%.2f,%.2f)", a, fa), &plt.A{Fsz: 7, NoClip: true, Ha: "right", Va: "bottom"})
		plt.Text(b+0.0, fb+0.0, io.Sf("b=(%.2f,%.2f)", b, fb), &plt.A{Fsz: 7, NoClip: true, Ha: "center", Va: "center"})
		plt.Text(c+0.2, fc+0.2, io.Sf("c=(%.2f,%.2f)", c, fc), &plt.A{Fsz: 7, NoClip: true, Ha: "left", Va: "top"})
		//plt.AxisYrange(-1.5, 1.5)
		plt.Gll("$x$", "$y$", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "bracket01")
	}
}

func TestBracket02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Bracket02. sin(5x)")

	// function
	ffcn := func(x float64) float64 { return math.Sin(5 * x) }

	// bracket minimum
	bracket := NewBracket(ffcn)
	bracket.Verbose = true

	// check
	a, b, c, fa, fb, fc := bracket.Min(0, 1)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-10, 10)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(1, 0)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(10, -10)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-2, -1)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-1, -2)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// plot
	if chk.Verbose {
		xx := utl.LinSpace(-4.0, 1.0, 101)
		yy := utl.GetMapped(xx, func(x float64) float64 { return ffcn(x) })
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Cross(0, 0, nil)
		plt.Plot(xx, yy, &plt.A{C: "b", Ls: "-", L: "curve1", NoClip: true})
		plt.PlotOne(a, fa, &plt.A{C: "r", M: "|", Ms: 40, NoClip: true})
		plt.PlotOne(b, fb, &plt.A{C: "r", M: "*", Ms: 10, NoClip: true})
		plt.PlotOne(c, fc, &plt.A{C: "r", M: "|", Ms: 40, NoClip: true})
		plt.Text(a-0.1, fa-0.1, io.Sf("a=(%.2f,%.2f)", a, fa), &plt.A{Fsz: 7, NoClip: true, Ha: "right", Va: "bottom"})
		plt.Text(b+0.0, fb+0.0, io.Sf("b=(%.2f,%.2f)", b, fb), &plt.A{Fsz: 7, NoClip: true, Ha: "center", Va: "center"})
		plt.Text(c+0.1, fc+0.1, io.Sf("c=(%.2f,%.2f)", c, fc), &plt.A{Fsz: 7, NoClip: true, Ha: "left", Va: "top"})
		plt.AxisYrange(-1.2, 1.2)
		plt.Gll("$x$", "$y$", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "bracket02")
	}
}
