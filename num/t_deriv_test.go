// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func f1(x float64, args ...interface{}) float64    { return math.Exp(x) }
func df1dx(x float64, args ...interface{}) float64 { return math.Exp(x) }
func f2(x float64, args ...interface{}) float64 {
	if x >= 0.0 {
		return x * math.Sqrt(x)
	}
	return 0.0
}
func df2dx(x float64, args ...interface{}) float64 {
	if x >= 0.0 {
		return 1.5 * math.Sqrt(x)
	}
	return 0.0
}
func f3(x float64, args ...interface{}) float64 {
	if x != 0.0 {
		return math.Sin(1 / x)
	}
	return 0.0
}
func df3dx(x float64, args ...interface{}) float64 {
	if x != 0.0 {
		return -math.Cos(1/x) / (x * x)
	}
	return 0.0
}
func f4(x float64, args ...interface{}) float64    { return math.Exp(-x * x) }
func df4dx(x float64, args ...interface{}) float64 { return -2.0 * x * math.Exp(-x*x) }
func f5(x float64, args ...interface{}) float64    { return x * x }
func df5dx(x float64, args ...interface{}) float64 { return 2.0 * x }
func f6(x float64, args ...interface{}) float64    { return 1.0 / x }
func df6dx(x float64, args ...interface{}) float64 { return -1.0 / (x * x) }

type NumDerivFunc func(f Cb_fx, x float64, h float64, args ...interface{}) (result, abserr float64)

func check(tst *testing.T, deriv NumDerivFunc, f Cb_fx, dfdx Cb_fx, x float64, desc string, args ...interface{}) {
	var tol float64 = 1.0e-6
	expected := dfdx(x)
	result, abserr := deriv(f, x, tol, args)
	if TestAbs(result, expected, MinComp(tol, expected), desc) != Equal {
		tst.Errorf("TestDeriv01 failed with abserr = [1;31m%g[0m", abserr)
	}
}

func TestDeriv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestDeriv 01")

	check(tst, DerivCentral, f1, df1dx, 1.0, "exp(x), x=1, central deriv")
	check(tst, DerivForward, f1, df1dx, 1.0, "exp(x), x=1, forward deriv")
	check(tst, DerivBackward, f1, df1dx, 1.0, "exp(x), x=1, backward deriv")

	check(tst, DerivCentral, f2, df2dx, 0.1, "x^(3/2), x=0.1, central deriv")
	check(tst, DerivForward, f2, df2dx, 0.1, "x^(3/2), x=0.1, forward deriv")
	check(tst, DerivBackward, f2, df2dx, 0.1, "x^(3/2), x=0.1, backward deriv")

	check(tst, DerivCentral, f3, df3dx, 0.45, "sin(1/x), x=0.45, central deriv")
	check(tst, DerivForward, f3, df3dx, 0.45, "sin(1/x), x=0.45, forward deriv")
	check(tst, DerivBackward, f3, df3dx, 0.45, "sin(1/x), x=0.45, backward deriv")

	check(tst, DerivCentral, f4, df4dx, 0.5, "exp(-x^2), x=0.5, central deriv")
	check(tst, DerivForward, f4, df4dx, 0.5, "exp(-x^2), x=0.5, forward deriv")
	check(tst, DerivBackward, f4, df4dx, 0.5, "exp(-x^2), x=0.5, backward deriv")

	check(tst, DerivCentral, f5, df5dx, 0.0, "x^2, x=0, central deriv")
	check(tst, DerivForward, f5, df5dx, 0.0, "x^2, x=0, forward deriv")
	check(tst, DerivBackward, f5, df5dx, 0.0, "x^2, x=0, backward deriv")

	check(tst, DerivCentral, f6, df6dx, 10.0, "1/x, x=10, central deriv")
	check(tst, DerivForward, f6, df6dx, 10.0, "1/x, x=10, forward deriv")
	check(tst, DerivBackward, f6, df6dx, 10.0, "1/x, x=10, backward deriv")
}

func TestDeriv02(tst *testing.T) {

	verbose()
	chk.PrintTitle("TestDeriv 01")

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
			unum, _ := DerivCentral(func(xvar float64, a ...interface{}) float64 {
				return fcn(xvar, y)
			}, x, h)

			// numerical dfdy @ (x,y)
			vnum, _ := DerivCentral(func(yvar float64, a ...interface{}) float64 {
				return fcn(x, yvar)
			}, y, h)

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
