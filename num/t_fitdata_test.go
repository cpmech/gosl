// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"testing"

	"gosl/chk"
	"gosl/io"
)

func TestLinFit01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinFit01")

	// data
	x := []float64{1, 2, 3, 4}
	y := []float64{1, 2, 3, 4}
	a, b, σa, σb, χ2 := LinFitSigma(x, y)
	io.Pforan("a=%v b=%v σa=%v σb=%v χ2=%v\n", a, b, σa, σb, χ2)
	chk.Float64(tst, "a", 1e-17, a, 0)
	chk.Float64(tst, "b", 1e-17, b, 1)
	chk.Float64(tst, "σa", 1e-17, σa, 0)
	chk.Float64(tst, "σb", 1e-17, σb, 0)
	chk.Float64(tst, "χ2", 1e-17, χ2, 0)
}

func TestLinFit02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinFit02")

	// data
	x := []float64{1, 2, 3, 4}
	y := []float64{6, 5, 7, 10}
	a, b := LinFit(x, y)
	io.Pforan("a=%v b=%v\n", a, b)
	chk.Float64(tst, "a", 1e-17, a, 3.5)
	chk.Float64(tst, "b", 1e-17, b, 1.4)
}
