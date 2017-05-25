// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"math"
	"testing"
)

func Test_deriv01(tst *testing.T) {

	//Verbose = true
	PrintTitle("deriv01. DerivScaScaCen")

	f := func(x float64) float64 { return math.Cos(math.Pi * x / 2.0) }

	dfdxAna := -1.0 * math.Pi / 2.0
	xAt := 1.0
	dx := 1e-3

	t1 := new(testing.T)
	DerivScaScaCen(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, f)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivScaScaCen(t2, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, f)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}
}

func Test_deriv02(tst *testing.T) {

	//Verbose = true
	PrintTitle("deriv02. DerivVecScaCen")

	fcn := func(f []float64, x float64) {
		f[0] = math.Cos(math.Pi * x / 2.0)
		f[1] = math.Sin(math.Pi * x / 2.0)
		return
	}

	dfdx := func(x float64) []float64 {
		return []float64{
			-math.Sin(math.Pi*x/2.0) * math.Pi / 2.0,
			+math.Cos(math.Pi*x/2.0) * math.Pi / 2.0,
		}
	}

	dx := 1e-3
	xAt := 1.0
	dfdxAna := dfdx(xAt)

	t1 := new(testing.T)
	DerivVecScaCen(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivVecScaCen(t2, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}

	xAt = 0.0
	dfdxAna = dfdx(xAt)
	dfdxAna[0] += 0.0001

	t3 := new(testing.T)
	DerivVecScaCen(t3, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	dfdxAna = dfdx(xAt)
	t4 := new(testing.T)
	DerivVecScaCen(t4, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t4.Failed() {
		tst.Errorf("t4 should not have failed\n")
		return
	}
}

func Test_deriv03(tst *testing.T) {

	//Verbose = true
	PrintTitle("deriv03. DerivScaVecCen")

	fcn := func(x []float64) float64 {
		return x[0]*x[0]*x[0] + x[1]*x[1] + x[0]*x[1] + x[0] - x[1]
	}

	dfdx := func(x []float64) []float64 {
		return []float64{
			3.0*x[0]*x[0] + x[1] + 1.0,
			2.0*x[1] + x[0] - 1.0,
		}
	}

	dx := 1e-3
	xAt := []float64{0.5, 0.5}
	dfdxAna := dfdx(xAt)

	t1 := new(testing.T)
	DerivScaVecCen(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivScaVecCen(t2, "dfdx", 2.0e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}

	xAt = []float64{0.0, 0.0}
	dfdxAna = dfdx(xAt)

	t3 := new(testing.T)
	DerivScaVecCen(t3, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if t3.Failed() {
		tst.Errorf("t3 should not have failed\n")
		return
	}
}

func Test_deriv04(tst *testing.T) {

	//Verbose = true
	PrintTitle("deriv04. DerivVecVecCen")

	fcn := func(f, x []float64) {
		f[0] = x[0]*x[0]*x[0] + x[1]*x[1] + x[0]*x[1] + x[0] - x[1]
		f[1] = math.Cos(math.Pi*x[0]/2.0) * math.Sin(math.Pi*x[1]/2.0)
	}

	dfdx := func(x []float64) [][]float64 {
		return [][]float64{
			{3.0*x[0]*x[0] + x[1] + 1.0, 2.0*x[1] + x[0] - 1.0},
			{-0.5 * math.Pi * math.Sin(math.Pi*x[0]/2.0) * math.Sin(math.Pi*x[1]/2.0), 0.5 * math.Pi * math.Cos(math.Pi*x[0]/2.0) * math.Cos(math.Pi*x[1]/2.0)},
		}
	}

	dx := 1e-3
	xAt := []float64{0.5, 0.5}
	dfdxAna := dfdx(xAt)

	t1 := new(testing.T)
	DerivVecVecCen(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivVecVecCen(t2, "dfdx", 2.0e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}

	xAt = []float64{0.0, 0.0}
	dfdxAna = dfdx(xAt)

	t3 := new(testing.T)
	DerivVecVecCen(t3, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if t3.Failed() {
		tst.Errorf("t3 should not have failed\n")
		return
	}
}
