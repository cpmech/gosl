// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func Test_cubiceq01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cubiceq01. y(x) = x³ - 3x² - 144x + 432")

	a, b, c := -3.0, -144.0, 432.0
	x1, x2, x3, nx := EqCubicSolveReal(a, b, c)
	io.Pforan("\na=%v b=%v c=%v\n", a, b, c)
	io.Pfcyan("nx=%v\n", nx)
	io.Pfcyan("x1=%v x2=%v x3=%v\n", x1, x2, x3)
	chk.IntAssert(nx, 3)
	chk.Float64(tst, "x1", 1e-17, x1, -12)
	chk.Float64(tst, "x2", 1e-17, x2, 12)
	chk.Float64(tst, "x3", 1e-14, x3, 3)
}

func Test_cubiceq02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cubiceq02. y(x) = x³ + x²")

	a, b, c := 1.0, 0.0, 0.0
	x1, x2, x3, nx := EqCubicSolveReal(a, b, c)
	io.Pforan("\na=%v b=%v c=%v\n", a, b, c)
	io.Pfcyan("nx=%v\n", nx)
	io.Pfcyan("x1=%v x2=%v x3=%v\n", x1, x2, x3)
	chk.IntAssert(nx, 2)
	chk.Float64(tst, "x1", 1e-17, x1, -1)
	chk.Float64(tst, "x2", 1e-17, x2, 0)
}

func Test_cubiceq03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cubiceq03. y(x) = x³ + c")

	doplot := chk.Verbose
	np := 41
	var X, Y []float64
	if doplot {
		X = utl.LinSpace(-2, 2, np)
		Y = make([]float64, np)
		plt.Reset(false, nil)
	}

	a, b := 0.0, 0.0
	colors := []string{"red", "green", "blue"}
	for k, c := range []float64{-1, 0, 1} {
		x1, x2, x3, nx := EqCubicSolveReal(a, b, c)
		io.Pforan("\na=%v b=%v c=%v\n", a, b, c)
		io.Pfcyan("nx=%v\n", nx)
		io.Pfcyan("x1=%v x2=%v x3=%v\n", x1, x2, x3)
		chk.IntAssert(nx, 1)
		chk.Float64(tst, "x1", 1e-17, x1, -c)
		if doplot {
			for i, x := range X {
				Y[i] = x*x*x + a*x*x + b*x + c
			}
			plt.Plot(X, Y, &plt.A{C: colors[k], L: io.Sf("c=%g", c)})
			plt.PlotOne(x1, 0, &plt.A{C: colors[k], M: "o"})
			plt.Cross(0, 0, nil)
			plt.Gll("x", "y", nil)
		}
	}
	if doplot {
		plt.Save("/tmp/gosl", "fig_cubiceq03")
	}
}

func Test_cubiceq04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cubiceq03. y(x) = x³ - 3xr²/4 - r³cos(3α)/4")

	doplot := chk.Verbose
	r := 1.0
	np := 41
	var X, Y []float64
	if doplot {
		X = utl.LinSpace(-1.2*r, 1.2*r, np)
		Y = make([]float64, np)
		plt.Reset(false, nil)
	}

	π := math.Pi
	a := 0.0
	b := -3.0 * r * r / 4.0
	colors := []string{"red", "green", "blue"}
	for k, α := range []float64{0, π / 6.0, π / 3.0} {
		c := -math.Pow(r, 3.0) * math.Cos(3.0*α) / 4.0
		for i, x := range X {
			Y[i] = x*x*x + a*x*x + b*x + c
		}
		x1, x2, x3, nx := EqCubicSolveReal(a, b, c)
		io.Pforan("\na=%v b=%v c=%v\n", a, b, c)
		io.Pfcyan("nx=%v\n", nx)
		io.Pfcyan("x1=%v x2=%v x3=%v\n", x1, x2, x3)
		if k == 0 {
			chk.IntAssert(nx, 2)
			chk.Float64(tst, "x1", 1e-17, x1, r)
			chk.Float64(tst, "x2", 1e-17, x2, -r/2.0)
		}
		if k == 1 {
			chk.IntAssert(nx, 3)
			chk.Float64(tst, "x1", 1e-15, x1, r*math.Cos(α+2.0*π/3.0))
			chk.Float64(tst, "x2", 1e-17, x2, r*math.Cos(α))
			chk.Float64(tst, "x3", 1e-15, x3, r*math.Cos(α-2.0*π/3.0))
		}
		if k == 2 {
			chk.IntAssert(nx, 2)
			chk.Float64(tst, "x1", 1e-17, x1, -r)
			chk.Float64(tst, "x2", 1e-17, x2, r/2.0)
		}
		if doplot {
			switch nx {
			case 1:
				plt.Plot([]float64{x1}, []float64{0}, &plt.A{C: colors[k], M: "o"})
			case 2:
				plt.Plot([]float64{x1, x2}, []float64{0, 0}, &plt.A{C: colors[k], M: "o"})
			case 3:
				plt.Plot([]float64{x1, x2, x3}, []float64{0, 0, 0}, &plt.A{C: colors[k], M: "o"})
			}
			plt.Plot(X, Y, &plt.A{C: colors[k], L: plt.TexPiRadFmt(α)})
		}
	}
	if doplot {
		plt.Circle(0, 0, r, &plt.A{Ec: "black", Fc: "none"})
		plt.Equal()
		plt.Gll("x", "y", nil)
		plt.Save("/tmp/gosl", "fig_cubiceq04")
	}
}
