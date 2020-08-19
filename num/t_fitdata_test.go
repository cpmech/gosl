// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/plt"
	"gosl/utl"
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

	// plot
	if chk.Verbose {
		model := func(x float64) float64 { return a + b*x }
		xx := utl.LinSpace(0, 5, 11)
		yy := utl.GetMapped(xx, func(x float64) float64 { return model(x) })
		plt.Reset(true, nil)
		plt.Plot(x, y, &plt.A{L: "data", C: plt.C(0, 0), M: plt.M(0, 0), NoClip: true})
		plt.Plot(xx, yy, &plt.A{L: "model", C: plt.C(1, 0), NoClip: true})
		plt.Gll("$x$", "$y$", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "linfit01")
	}
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

	// plot
	if chk.Verbose {
		model := func(x float64) float64 { return a + b*x }
		xx := utl.LinSpace(0, 5, 11)
		yy := utl.GetMapped(xx, func(x float64) float64 { return model(x) })
		plt.Reset(true, nil)
		plt.Plot(x, y, &plt.A{L: "data", C: plt.C(0, 0), M: plt.M(0, 0), NoClip: true})
		plt.Plot(xx, yy, &plt.A{L: "model", C: plt.C(1, 0), NoClip: true})
		plt.Gll("$x$", "$y$", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/num", "linfit02")
	}
}
