// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
)

func TestDeriv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deriv01. 1st derivs using Cen5, Fwd4, Bwd4")

	names := []string{
		"x²",       // 1
		"exp(x)",   // 2
		"exp(-x²)", // 3
		"1/x",      // 4
		"x⋅√x",     // 5
		"sin(1/x)", // 6
	}

	fcns := []fun.Ss{
		func(x float64) float64 { // 1
			return x * x
		},
		func(x float64) float64 { // 2
			return math.Exp(x)
		},
		func(x float64) float64 { // 3
			return math.Exp(-x * x)
		},
		func(x float64) float64 { // 4
			return 1.0 / x
		},
		func(x float64) float64 { // 5
			return x * math.Sqrt(x)
		},
		func(x float64) float64 { // 6
			return math.Sin(1.0 / x)
		},
	}

	danas := []fun.Ss{
		func(x float64) float64 { // 1
			return 2 * x
		},
		func(x float64) float64 { // 2
			return math.Exp(x)
		},
		func(x float64) float64 { // 3
			return -2 * x * math.Exp(-x*x)
		},
		func(x float64) float64 { // 4
			return -1.0 / (x * x)
		},
		func(x float64) float64 { // 5
			return 1.5 * math.Sqrt(x)
		},
		func(x float64) float64 { // 6
			return -math.Cos(1.0/x) / (x * x)
		},
	}

	xvals := [][]float64{
		{-0.5, 0, 0.5},   // 1
		{-0.5, 0, 0.5},   // 2
		{-0.5, 0, 0.5},   // 3
		{-0.5, 0.2, 0.5}, // 4
		{0.2, 0.5, 1.0},  // 5
		{-0.5, 0.2, 0.5}, // 6
	}

	//                    1      2      3     4      5      6
	tols := []float64{1e-15, 1e-10, 1e-11, 1e-9, 1e-11, 1e-10}

	//                 1     2     3     4     5     6
	hs := []float64{1e-1, 1e-1, 1e-1, 1e-1, 1e-1, 1e-1}

	// check
	smethods := []string{"DerivCen5", "DerivFwd4", "DerivBwd4"}
	methods := []func(float64, float64, fun.Ss) float64{
		DerivCen5, DerivFwd4, DerivBwd4,
	}
	for idx, method := range methods {
		if idx > 0 {
			tols = []float64{1e-7, 1e-7, 1e-7, 1e-7, 1e-7, 1e-6}
		}
		io.Pf("\ncheck %s:\n", smethods[idx])
		for i, f := range fcns {
			for _, x := range xvals[i] {
				dnum := method(x, hs[i], f)
				dana := danas[i](x)
				chk.Float64(tst, "    "+names[i], tols[i], dnum, dana)
			}
		}
	}
}

func TestDeriv02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deriv02. 1st deriv using Cen5")

	// scalar field
	fcn := func(x, y float64) float64 {
		return -math.Pow(math.Pow(math.Cos(x), 2.0)+math.Pow(math.Cos(y), 2.0), 2.0)
	}

	// gradient. u=dfdx, v=dfdy
	grad := func(x, y float64) (u, v float64) {
		m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
		u = 4.0 * math.Cos(x) * math.Sin(x) * m
		v = 4.0 * math.Cos(y) * math.Sin(y) * m
		return
	}

	// grid size
	xmin, xmax, N := -math.Pi/2.0, math.Pi/2.0, 7
	dx := (xmax - xmin) / float64(N-1)

	// step size for numerical differentiation
	h := 0.01

	// tolerance
	tol := 1e-9

	// loop over grid points
	for i := 0; i < N; i++ {
		x := xmin + float64(i)*dx
		for j := 0; j < N; j++ {
			y := xmin + float64(i)*dx

			// scalar and vector field
			f := fcn(x, y)
			u, v := grad(x, y)

			// numerical dfdx @ (x,y)
			unum := DerivCen5(x, h, func(xvar float64) float64 {
				return fcn(xvar, y)
			})

			// numerical dfdy @ (x,y)
			vnum := DerivCen5(y, h, func(yvar float64) float64 {
				return fcn(x, yvar)
			})

			// output
			if chk.Verbose {
				io.Pforan("x=%+.3f y=%+.3f f=%+.3f u=%+.3f v=%+.3f unum=%+.3f vnum=%+.3f\n", x, y, f, u, v, unum, vnum)
			}

			// check
			if math.Abs(unum-u) > tol {
				tst.Errorf("x=%v y=%v f=%v u=%v v=%v unum=%v vnum=%v error=%v\n", x, y, f, u, v, unum, vnum, math.Abs(unum-u))
				return
			}
			if math.Abs(vnum-v) > tol {
				tst.Errorf("x=%v y=%v f=%v u=%v v=%v unum=%v vnum=%v error=%v\n", x, y, f, u, v, unum, vnum, math.Abs(vnum-v))
				return
			}
		}
	}
}

func TestDeriv03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deriv03. d²f/dx² using 2nd order Cen3 and Cen5")

	names := []string{
		"x²",       // 1
		"exp(x)",   // 2
		"exp(-x²)", // 3
		"1/x",      // 4
		"x⋅√x",     // 5
		"sin(1/x)", // 6
	}

	fcns := []fun.Ss{
		func(x float64) float64 { // 1
			return x * x
		},
		func(x float64) float64 { // 2
			return math.Exp(x)
		},
		func(x float64) float64 { // 3
			return math.Exp(-x * x)
		},
		func(x float64) float64 { // 4
			return 1.0 / x
		},
		func(x float64) float64 { // 5
			return x * math.Sqrt(x)
		},
		func(x float64) float64 { // 6
			return math.Sin(1.0 / x)
		},
	}

	ddanas := []fun.Ss{ // d²f/dx²
		func(x float64) float64 { // 1
			return 2
		},
		func(x float64) float64 { // 2
			return math.Exp(x)
		},
		func(x float64) float64 { // 3
			return 4*x*x*math.Exp(-x*x) - 2*math.Exp(-x*x)
		},
		func(x float64) float64 { // 4
			return 2.0 / (x * x * x)
		},
		func(x float64) float64 { // 5
			return 0.75 / math.Sqrt(x)
		},
		func(x float64) float64 { // 6
			x3 := x * x * x
			x4 := x * x3
			return 2.0*math.Cos(1.0/x)/x3 - math.Sin(1.0/x)/x4
		},
	}

	xvals := [][]float64{
		{-0.5, 0, 0.5},   // 1
		{-0.5, 0, 0.5},   // 2
		{-0.5, 0, 0.5},   // 3
		{-0.5, 0.2, 0.5}, // 4
		{0.2, 0.5, 1.0},  // 5
		{-0.5, 0.2, 0.5}, // 6
	}

	tols := [][]float64{
		{1e-10, 1e-6, 1e-6, 1e-2, 1e-5, 1e-2},
		{1e-10, 1e-9, 1e-9, 1e-6, 1e-9, 1e-4},
	}

	hs := []float64{1e-3, 1e-3, 1e-3, 1e-3, 1e-3, 1e-3}

	// check
	smethods := []string{"SecondDerivCen3", "SecondDerivCen5"}
	methods := []func(float64, float64, fun.Ss) float64{
		SecondDerivCen3, SecondDerivCen5,
	}
	for idx, method := range methods {
		io.Pf("\ncheck %s:\n", smethods[idx])
		for i, f := range fcns {
			for _, x := range xvals[i] {
				dnum := method(x, hs[i], f)
				dana := ddanas[i](x)
				chk.Float64(tst, "    "+names[i], tols[idx][i], dnum, dana)
			}
		}
	}
}

func TestDeriv04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deriv04. mixed derivatives d²f/dxdy²")

	names := []string{
		"2x³+3xy+2y²", // 1
		"exp(xy)",     // 2
		"x⋅sin(y)+y⋅cos(x)+xy⋅sin(xy)+xy⋅cos(xy)", // 3
	}

	fcns := []fun.Sss{
		func(x, y float64) float64 { // 1
			return 2*x*x*x + 3*x*y + 2*y*y
		},
		func(x, y float64) float64 { // 2
			return math.Exp(x * y)
		},
		func(x, y float64) float64 { // 3
			return x*math.Sin(y) + y*math.Cos(x) + x*y*math.Sin(x*y) + x*y*math.Cos(x*y)
		},
	}

	dxdanas := []fun.Sss{ // d²f/dxdy
		func(x, y float64) float64 { // 1
			return 3
		},
		func(x, y float64) float64 { // 2
			return x*y*math.Exp(x*y) + math.Exp(x*y)
		},
		func(x, y float64) float64 { // 3
			return -x*x*y*y*math.Sin(x*y) - 3*x*y*math.Sin(x*y) + math.Sin(x*y) - x*x*y*y*cos(x*y) + 3*x*y*math.Cos(x*y) + math.Cos(x*y) + math.Cos(y) - math.Sin(x)
		},
	}

	xvals := [][]float64{
		{-0.5, 0, 0.5}, // 1
		{-0.5, 0, 0.5}, // 2
		{-0.5, 0, 0.5}, // 3
	}

	yvals := [][]float64{
		{-0.5, 0, 0.5}, // 1
		{-0.5, 0, 0.5}, // 2
		{-0.5, 0, 0.5}, // 3
	}

	tols := [][]float64{
		{1e-10, 1e-6, 1e-5},   // method 1
		{1e-10, 1e-10, 1e-10}, // method 2
		{1e-9, 1e-9, 1e-10},   // method 3
	}

	hs := []float64{1e-3, 1e-3, 1e-3}

	// check
	smethods := []string{"SecondDerivMixedO2", "SecondDerivMixedO4v1", "SecondDerivMixedO4v2"}
	methods := []func(float64, float64, float64, fun.Sss) float64{
		SecondDerivMixedO2, SecondDerivMixedO4v1, SecondDerivMixedO4v2,
	}
	for idx, method := range methods {
		io.Pf("\ncheck %s:\n", smethods[idx])
		for i, f := range fcns {
			for _, x := range xvals[i] {
				for _, y := range yvals[i] {
					dnum := method(x, y, hs[i], f)
					dana := dxdanas[i](x, y)
					chk.Float64(tst, "    "+names[i], tols[idx][i], dnum, dana)
				}
			}
		}
	}
}
