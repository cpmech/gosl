// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"gosl/chk"
	"gosl/io"
)

func TestInterpCubic01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("InterpCubic01. Interp with cubic poly using 4 points")

	// test set
	ycor := func(x float64) float64 { return Pow3(x) - 3*Pow2(x) - 144*x + 432 }
	dcor := func(x float64) float64 { return 3*Pow2(x) - 6*x - 144 }
	x0, y0 := -12.0, 0.0
	x1, y1 := -6.0, 972.0
	x2, y2 := 1.0, 286.0
	x3, y3 := 12.0, 0.0

	// intepolator
	interp := NewInterpCubic()
	interp.Fit4points(x0, y0, x1, y1, x2, y2, x3, y3)

	// check model and derivatives
	for _, x := range []float64{-10, 0, 1, 8} {
		y := interp.F(x)
		yd := interp.G(x)
		chk.Float64(tst, io.Sf("y(%g)", x), 1e-15, y, ycor(x))
		chk.Float64(tst, io.Sf("y'(%g)", x), 1e-15, yd, dcor(x))
	}

	// check critical points
	xmin, xmax, xifl, hasMin, hasMax, hasIfl := interp.Critical()
	if !hasMin {
		tst.Errorf("hasMin should be true\n")
		return
	}
	if !hasMax {
		tst.Errorf("hasMax should be true\n")
		return
	}
	if !hasIfl {
		tst.Errorf("hasIfl should be true\n")
		return
	}
	chk.Float64(tst, "xmin", 1e-15, xmin, 8)
	chk.Float64(tst, "xmax", 1e-15, xmax, -6)
	chk.Float64(tst, "xifl", 1e-15, xifl, 1)
}

func TestInterpCubic02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("InterpCubic02. Inflection only")

	// test set
	ycor := func(x float64) float64 { return Pow3(x-25) + 450 }
	dcor := func(x float64) float64 { return 3 * Pow2(x-25) }
	x0, x1, x2, x3 := 22.0, 23.0, 24.0, 25.0
	y0, y1, y2, y3 := ycor(x0), ycor(x1), ycor(x2), ycor(x3)

	// intepolator
	interp := NewInterpCubic()
	interp.Fit4points(x0, y0, x1, y1, x2, y2, x3, y3)

	// check model and derivatives
	for _, x := range []float64{-10, 0, 5, 100} {
		y := interp.F(x)
		yd := interp.G(x)
		chk.Float64(tst, io.Sf("y(%g)", x), 1e-15, y, ycor(x))
		chk.Float64(tst, io.Sf("y'(%g)", x), 1e-15, yd, dcor(x))
	}

	// check critical points
	_, _, xifl, hasMin, hasMax, hasIfl := interp.Critical()
	if hasMin {
		tst.Errorf("hasMin should be false\n")
		return
	}
	if hasMax {
		tst.Errorf("hasMax should be false\n")
		return
	}
	if !hasIfl {
		tst.Errorf("hasIfl should be true\n")
		return
	}
	chk.Float64(tst, "xifl", 1e-15, xifl, 25)
}

func TestInterpCubic03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("InterpCubic03. Strictly monotonic")

	// test set
	ycor := func(x float64) float64 { return Pow3(x) + Pow2(x) + x + 1 }
	dcor := func(x float64) float64 { return 3*Pow2(x) + 2*x + 1 }
	x0, x1, x2, x3 := -2.0, -1.0, 1.0, 2.0
	y0, y1, y2, y3 := ycor(x0), ycor(x1), ycor(x2), ycor(x3)

	// intepolator
	interp := NewInterpCubic()
	interp.Fit4points(x0, y0, x1, y1, x2, y2, x3, y3)

	// check model and derivatives
	for _, x := range []float64{-10, 0, 5, 100} {
		y := interp.F(x)
		yd := interp.G(x)
		chk.Float64(tst, io.Sf("y(%g)", x), 1e-15, y, ycor(x))
		chk.Float64(tst, io.Sf("y'(%g)", x), 1e-15, yd, dcor(x))
	}

	// check critical points
	_, _, _, hasMin, hasMax, hasIfl := interp.Critical()
	if hasMin {
		tst.Errorf("hasMin should be false\n")
		return
	}
	if hasMax {
		tst.Errorf("hasMax should be false\n")
		return
	}
	if hasIfl {
		tst.Errorf("hasIfl should be false\n")
		return
	}
}

func TestInterpCubic04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("InterpCubic01. Interp with cubic poly using 3 points and deriv")

	// test set
	ycor := func(x float64) float64 { return Pow3(x) - 3*Pow2(x) - 144*x + 432 }
	dcor := func(x float64) float64 { return 3*Pow2(x) - 6*x - 144 }
	x0, y0 := -12.0, 0.0
	x1, y1 := -6.0, 972.0
	x2, y2 := 1.0, 286.0
	x3, d3 := 8.0, 0.0

	// intepolator
	interp := NewInterpCubic()
	interp.Fit3pointsD(x0, y0, x1, y1, x2, y2, x3, d3)

	// check model and derivatives
	for _, x := range []float64{-10, 0, 1, 8} {
		y := interp.F(x)
		yd := interp.G(x)
		chk.Float64(tst, io.Sf("y(%g)", x), 1e-15, y, ycor(x))
		chk.Float64(tst, io.Sf("y'(%g)", x), 1e-15, yd, dcor(x))
	}

	// check critical points
	xmin, xmax, xifl, hasMin, hasMax, hasIfl := interp.Critical()
	if !hasMin {
		tst.Errorf("hasMin should be true\n")
		return
	}
	if !hasMax {
		tst.Errorf("hasMax should be true\n")
		return
	}
	if !hasIfl {
		tst.Errorf("hasIfl should be true\n")
		return
	}
	chk.Float64(tst, "xmin", 1e-15, xmin, 8)
	chk.Float64(tst, "xmax", 1e-15, xmax, -6)
	chk.Float64(tst, "xifl", 1e-15, xifl, 1)
}
