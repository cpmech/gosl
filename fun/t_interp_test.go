// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Interp01. linear interpolation")

	xx := []float64{0, 1, 2, 3, 4, 5}
	yy := []float64{0.50, 0.20, 0.20, 0.05, 0.01, 0.00}

	o, err := NewInterpolator(LinearInterpKind, 1, xx, yy)
	if err != nil {
		tst.Errorf("%v\n", err)
	}

	for i, x := range xx {
		chk.Float64(tst, "P(xi)", 1e-17, o.P(x), yy[i])
	}

	xref := []float64{1.0 / 3.0, 2.5, 2.0 / 3.0, 1.1, 1.5, 3.5, 4.5}
	yref := []float64{0.4, 0.125, 0.3, 0.2, 0.2, 0.03, 0.005}
	for i, x := range xref {
		chk.Float64(tst, "P(xref)", 1e-16, o.P(x), yref[i])
	}

	if chk.Verbose {
		X := utl.LinSpace(-0.5, 5.5, 101)
		Y := utl.GetMapped(X, func(x float64) float64 { return o.P(x) })
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Plot(xx, yy, &plt.A{C: "b", Ls: "-", M: ".", L: "data", NoClip: true})
		plt.Plot(X, Y, &plt.A{C: "r", Ls: ":", M: "+", L: "interp", NoClip: true})
		plt.Plot(xref, yref, &plt.A{C: "g", Ls: "none", M: "o", L: "interp", NoClip: true})
		plt.Gll("x", "y", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "interp01")
	}
}

func TestInterp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Interp02. polynomial interpolation")

	xx := []float64{0, 1, 2, 3, 4, 5}
	yy := []float64{0.50, 0.20, 0.20, 0.05, 0.01, 0.00}

	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Plot(xx, yy, &plt.A{C: "k", Ls: "-", M: ".", L: "data", NoClip: true})
	}

	for _, p := range []int{1, 2, 3} {

		o, err := NewInterpolator(PolyInterpKind, p, xx, yy)
		if err != nil {
			tst.Errorf("%v\n", err)
		}

		for i, x := range xx {
			chk.Float64(tst, "P(xi)", 1e-17, o.P(x), yy[i])
		}

		if o.m == 2 {
			xref := []float64{1.0 / 3.0, 2.5, 2.0 / 3.0, 1.1, 1.5, 3.5, 4.5}
			yref := []float64{0.4, 0.125, 0.3, 0.2, 0.2, 0.03, 0.005}
			for i, x := range xref {
				chk.Float64(tst, "P(xref)", 1e-16, o.P(x), yref[i])
			}
		}

		if chk.Verbose {
			X := utl.LinSpace(-0.5, 5.5, 101)
			Y := utl.GetMapped(X, func(x float64) float64 { return o.P(x) })
			plt.Plot(X, Y, &plt.A{Ls: "-", L: io.Sf("p=%d", p), NoClip: true})
			plt.Gll("x", "y", nil)
			plt.HideTRborders()
		}
	}
	if chk.Verbose {
		plt.Save("/tmp/gosl/fun", "interp02")
	}
}
