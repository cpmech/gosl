// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"gosl/chk"
	"gosl/plt"
	"gosl/utl"
)

func TestMultiInterp01(t *testing.T) {

	//verbose()
	chk.PrintTitle("MultiInterp01. bilinear interpolation")

	// checking f(x,y) = x^2 + 2y^2
	f := []float64{
		0.00, 0.25, 1.00, 4.00,
		2.00, 2.25, 3.00, 6.00,
		8.00, 8.25, 9.00, 12.00,
	}

	xx := []float64{0.0, 0.5, 1.0, 2.0}
	yy := []float64{0.0, 1.0, 2.0}

	o := NewBiLinear(f, xx, yy)

	for i, x := range xx {
		for j, y := range yy {
			chk.Float64(t, "P(x,y)", 1e-17, o.P(x, y), f[i+j*len(xx)])
		}
	}

	fref := []float64{1.125, 1.625, 3.2, 8.4}
	xref := []float64{0.25, 0.75, 1.2, 1.2}
	yref := []float64{0.5, 0.5, 0.8, 1.8}

	for i := 0; i < len(fref); i++ {
		chk.Float64(t, "P(xref,yref)", 1e-17, o.P(xref[i], yref[i]), fref[i])
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		xmin, xmax, ymin, ymax := 0.0, 2.0, 0.0, 2.0
		nx, ny := 41, 41
		uu, vv, ww := utl.MeshGrid2dF(xmin, xmax, ymin, ymax, nx, ny, func(x, y float64) float64 {
			return x*x + 2*y*y
		})
		rr, ss, tt := utl.MeshGrid2dF(xmin, xmax, ymin, ymax, nx, ny, func(x, y float64) float64 {
			return o.P(x, y)
		})
		plt.Wireframe(uu, vv, ww, &plt.A{C: plt.C(0, 0), Rstride: 5, Cstride: 5})
		plt.Wireframe(rr, ss, tt, &plt.A{C: plt.C(1, 0), Rstride: 5, Cstride: 5})
		for i, x := range xx {
			for j, y := range yy {
				plt.Plot3dPoint(x, y, f[i+j*len(xx)], &plt.A{C: "r", M: "o"})
			}
		}
		plt.Save("/tmp/gosl/fun", "multiinterp01")
	}
}
