// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
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
