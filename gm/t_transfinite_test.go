// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestTransfinite01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite01")

	π := math.Pi
	e0 := []float64{1, 0}
	e1 := []float64{0, 1}

	trf := NewTransfinite(2, []fun.Vs{
		func(x la.Vector, ξ float64) {
			for i := 0; i < len(x); i++ {
				x[i] = (2 + ξ) * e0[i]
			}
		},
		func(x la.Vector, η float64) {
			for i := 0; i < len(x); i++ {
				θ := π * (η + 1) / 4.0
				x[i] = 3*math.Cos(θ)*e0[i] + 3*math.Sin(θ)*e1[i]
			}
		},
		func(x la.Vector, ξ float64) {
			for i := 0; i < len(x); i++ {
				x[i] = (2 + ξ) * e1[i]
			}
		},
		func(x la.Vector, η float64) {
			for i := 0; i < len(x); i++ {
				θ := π * (η + 1) / 4.0
				x[i] = math.Cos(θ)*e0[i] + math.Sin(θ)*e1[i]
			}
		},
	})

	chk.Array(tst, "C0", 1e-17, trf.C[0], []float64{1, 0})
	chk.Array(tst, "C1", 1e-17, trf.C[1], []float64{3, 0})
	chk.Array(tst, "C2", 1e-17, trf.C[2], []float64{0, 3})
	chk.Array(tst, "C3", 1e-17, trf.C[3], []float64{0, 1})

	a := 1.0 / math.Sqrt(2)
	b := 2.0 / math.Sqrt(2)
	c := 3.0 / math.Sqrt(2)
	x := la.NewVector(2)

	trf.QuadMap(x, []float64{-1, -1})
	chk.Array(tst, "x(-1,-1)", 1e-17, x, []float64{1, 0})

	trf.QuadMap(x, []float64{0, -1})
	chk.Array(tst, "x( 0,-1)", 1e-17, x, []float64{2, 0})

	trf.QuadMap(x, []float64{+1, -1})
	chk.Array(tst, "x(+1,-1)", 1e-17, x, []float64{3, 0})

	trf.QuadMap(x, []float64{-1, 0})
	chk.Array(tst, "x(-1, 0)", 1e-15, x, []float64{a, a})

	trf.QuadMap(x, []float64{0, 0})
	chk.Array(tst, "x( 0, 0)", 1e-15, x, []float64{b, b})

	trf.QuadMap(x, []float64{+1, 0})
	chk.Array(tst, "x(+1, 0)", 1e-15, x, []float64{c, c})

	trf.QuadMap(x, []float64{-1, +1})
	chk.Array(tst, "x(-1,+1)", 1e-15, x, []float64{0, 1})

	trf.QuadMap(x, []float64{0, +1})
	chk.Array(tst, "x( 0,+1)", 1e-15, x, []float64{0, 2})

	trf.QuadMap(x, []float64{+1, +1})
	chk.Array(tst, "x(+1,+1)", 1e-15, x, []float64{0, 3})

	trf.Γ[0](x, -1)
	io.Pf("%v\n", x)

	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{21, 21}, nil, nil)
		plt.Arc(0, 0, 1, 0, 90, &plt.A{C: plt.C(2, 0), NoClip: true, Z: 10})
		plt.Arc(0, 0, 3, 0, 90, &plt.A{C: plt.C(2, 0), NoClip: true, Z: 10})
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite01")
	}
}

func TestTransfinite02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite02")

	π := math.Pi
	e0 := []float64{1, 0}
	e1 := []float64{0, 1}

	trf := NewTransfinite(2, []fun.Vs{
		func(x la.Vector, ξ float64) {
			for i := 0; i < len(x); i++ {
				x[i] = (2 + ξ) * e0[i]
			}
		},
		func(x la.Vector, η float64) {
			for i := 0; i < len(x); i++ {
				x[i] = 1.5*(1-η)*e0[i] + 1.5*(1+η)*e1[i]
			}
		},
		func(x la.Vector, ξ float64) {
			for i := 0; i < len(x); i++ {
				x[i] = (2 + ξ) * e1[i]
			}
		},
		func(x la.Vector, η float64) {
			for i := 0; i < len(x); i++ {
				θ := π * (η + 1) / 4.0
				x[i] = math.Cos(θ)*e0[i] + math.Sin(θ)*e1[i]
			}
		},
	})

	chk.Array(tst, "C0", 1e-17, trf.C[0], []float64{1, 0})
	chk.Array(tst, "C1", 1e-17, trf.C[1], []float64{3, 0})
	chk.Array(tst, "C2", 1e-17, trf.C[2], []float64{0, 3})
	chk.Array(tst, "C3", 1e-17, trf.C[3], []float64{0, 1})

	a := 1.0 / math.Sqrt(2)
	c := 1.5
	b := (a + c) / 2.0
	x := la.NewVector(2)

	trf.QuadMap(x, []float64{-1, -1})
	chk.Array(tst, "x(-1,-1)", 1e-17, x, []float64{1, 0})

	trf.QuadMap(x, []float64{0, -1})
	chk.Array(tst, "x( 0,-1)", 1e-17, x, []float64{2, 0})

	trf.QuadMap(x, []float64{+1, -1})
	chk.Array(tst, "x(+1,-1)", 1e-17, x, []float64{3, 0})

	trf.QuadMap(x, []float64{-1, 0})
	chk.Array(tst, "x(-1, 0)", 1e-15, x, []float64{a, a})

	trf.QuadMap(x, []float64{0, 0})
	chk.Array(tst, "x( 0, 0)", 1e-15, x, []float64{b, b})

	trf.QuadMap(x, []float64{+1, 0})
	chk.Array(tst, "x(+1, 0)", 1e-15, x, []float64{c, c})

	trf.QuadMap(x, []float64{-1, +1})
	chk.Array(tst, "x(-1,+1)", 1e-15, x, []float64{0, 1})

	trf.QuadMap(x, []float64{0, +1})
	chk.Array(tst, "x( 0,+1)", 1e-15, x, []float64{0, 2})

	trf.QuadMap(x, []float64{+1, +1})
	chk.Array(tst, "x(+1,+1)", 1e-15, x, []float64{0, 3})

	trf.Γ[0](x, -1)
	io.Pf("%v\n", x)

	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{21, 21}, nil, nil)
		plt.Arc(0, 0, 1, 0, 90, &plt.A{C: plt.C(2, 0), NoClip: true, Z: 10})
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite02")
	}
}
