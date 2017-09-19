// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"math"
	"testing"
)

func TestDeriv01(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deriv01. DerivScaSca")

	f := func(x float64) float64 { return math.Cos(math.Pi * x / 2.0) }

	dfdxAna := -1.0 * math.Pi / 2.0
	xAt := 1.0
	dx := 1e-3

	t1 := new(testing.T)
	DerivScaSca(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, f)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivScaSca(t2, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, f)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}
}

func TestDeriv02(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deriv02. DerivVecSca")

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
	DerivVecSca(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivVecSca(t2, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}

	xAt = 0.0
	dfdxAna = dfdx(xAt)
	dfdxAna[0] += 0.0001

	t3 := new(testing.T)
	DerivVecSca(t3, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	dfdxAna = dfdx(xAt)
	t4 := new(testing.T)
	DerivVecSca(t4, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t4.Failed() {
		tst.Errorf("t4 should not have failed\n")
		return
	}
}

func TestDeriv03(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deriv03. DerivScaVec")

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
	DerivScaVec(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivScaVec(t2, "dfdx", 2.0e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}

	xAt = []float64{0.0, 0.0}
	dfdxAna = dfdx(xAt)

	t3 := new(testing.T)
	DerivScaVec(t3, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if t3.Failed() {
		tst.Errorf("t3 should not have failed\n")
		return
	}
}

func TestDeriv04(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deriv04. DerivVecVec")

	fcn := func(f, x []float64) {
		f[0] = x[0]*x[0]*x[0] + x[1]*x[1] + x[0]*x[1] + x[0] - x[1]
		f[1] = math.Cos(math.Pi*x[0]/2.0) * math.Sin(math.Pi*x[1]/2.0)
		return
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
	DerivVecVec(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivVecVec(t2, "dfdx", 2.0e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t2.Failed() {
		tst.Errorf("t2 should NOT have failed\n")
		return
	}

	xAt = []float64{0.0, 0.0}
	dfdxAna = dfdx(xAt)

	t3 := new(testing.T)
	DerivVecVec(t3, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if t3.Failed() {
		tst.Errorf("t3 should NOT have failed\n")
		return
	}
}

func TestDeriv05(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deriv05. Deriv methods. error messages")

	t1 := new(testing.T)
	DerivScaVec(t1, "", 0, []float64{1, 2, 3}, []float64{1}, 0, false, nil)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed")
	}

	t2 := new(testing.T)
	DerivVecVec(t2, "", 0, nil, []float64{1}, 0, false, nil)
	if !t2.Failed() {
		tst.Errorf("t2 should have failed")
	}

	t3 := new(testing.T)
	DerivVecVec(t3, "", 0, [][]float64{{1, 2}, {2, 1}, {3, 4}}, []float64{1}, 0, false, nil)
	if !t3.Failed() {
		tst.Errorf("t3 should have failed")
	}

	t4 := new(testing.T)
	DerivVecVec(t4, "", 0, [][]float64{{1}, {2, 1}, {3, 4}}, []float64{1}, 1e-3, false, func(f, x []float64) { return })
	if !t4.Failed() {
		tst.Errorf("t4 should have failed")
	}

	t6 := new(testing.T)
	xAt, h := 1.0, 1.0
	DerivScaSca(t6, "", 0, 0, xAt, h, false, func(x float64) float64 {
		if x == 2.0 { // xAt + h
			return 0
		}
		return x
	})
	if !t6.Failed() {
		tst.Errorf("t6 should have failed")
	}

	t7 := new(testing.T)
	DerivScaSca(t7, "", 0, 0, xAt, h, false, func(x float64) float64 {
		if x == 0.5 { // xAt - h/2
			return 0
		}
		return x
	})
	if !t7.Failed() {
		tst.Errorf("t7 should have failed")
	}

	t8 := new(testing.T)
	DerivScaSca(t8, "", 0, 0, xAt, h, false, func(x float64) float64 {
		if x == 1.5 { // xAt + h/2
			return 0
		}
		return x
	})
	if !t8.Failed() {
		tst.Errorf("t8 should have failed")
	}
}
