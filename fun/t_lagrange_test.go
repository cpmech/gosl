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

func TestLagCardinal01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCardinal01. Lagrange cardinal polynomials")

	xa, xb := -1.0, 1.0
	p := 4
	X := utl.LinSpace(xa, xb, p+1)
	for i := 0; i < p+1; i++ {
		for j, x := range X {
			lpi := LagrangeCardinal(p, i, x, X)
			ana := 1.0
			if i != j {
				ana = 0
			}
			chk.AnaNum(tst, io.Sf("L^%d_%d(X[%d])", p, i, j), 1e-17, lpi, ana, chk.Verbose)
		}
		io.Pl()
	}

	if chk.Verbose {
		xx := utl.LinSpace(xa, xb, 201)
		yy := make([]float64, len(xx))
		plt.Reset(true, nil)
		for n := 0; n < p+1; n++ {
			for k, x := range xx {
				yy[k] = LagrangeCardinal(p, n, x, X)
			}
			plt.Plot(xx, yy, &plt.A{NoClip: true})
		}
		Y := make([]float64, len(X))
		plt.Plot(X, Y, &plt.A{C: "k", Ls: "none", M: "o", Void: true, NoClip: true})
		plt.Gll("x", "y", nil)
		plt.Cross(0, 0, &plt.A{C: "grey"})
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/fun", "lagcardinal01")
	}
}
