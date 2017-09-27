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
	"github.com/cpmech/gosl/utl"
)

func TestTransfinite01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite01")

	// new mapping
	rin, rou := 2.0, 6.0 // radii
	trf := FactoryTfinite.Surf2dQuarterRing(rin, rou)

	// check corners
	chk.Array(tst, "C0", 1e-17, trf.C[0], []float64{rin, 0})
	chk.Array(tst, "C1", 1e-17, trf.C[1], []float64{rou, 0})
	chk.Array(tst, "C2", 1e-17, trf.C[2], []float64{0, rou})
	chk.Array(tst, "C3", 1e-17, trf.C[3], []float64{0, rin})

	// auxiliary
	a := rin / math.Sqrt(2)
	b := 0.5 * (rin + rou) / math.Sqrt(2)
	c := rou / math.Sqrt(2)
	x := la.NewVector(2)

	// check points
	trf.Point(x, []float64{-1, -1})
	chk.Array(tst, "x(-1,-1)", 1e-17, x, []float64{rin, 0})

	trf.Point(x, []float64{0, -1})
	chk.Array(tst, "x( 0,-1)", 1e-17, x, []float64{0.5 * (rin + rou), 0})

	trf.Point(x, []float64{+1, -1})
	chk.Array(tst, "x(+1,-1)", 1e-17, x, []float64{rou, 0})

	trf.Point(x, []float64{-1, 0})
	chk.Array(tst, "x(-1, 0)", 1e-15, x, []float64{a, a})

	trf.Point(x, []float64{0, 0})
	chk.Array(tst, "x( 0, 0)", 1e-15, x, []float64{b, b})

	trf.Point(x, []float64{+1, 0})
	chk.Array(tst, "x(+1, 0)", 1e-15, x, []float64{c, c})

	trf.Point(x, []float64{-1, +1})
	chk.Array(tst, "x(-1,+1)", 1e-15, x, []float64{0, rin})

	trf.Point(x, []float64{0, +1})
	chk.Array(tst, "x( 0,+1)", 1e-15, x, []float64{0, 0.5 * (rin + rou)})

	trf.Point(x, []float64{+1, +1})
	chk.Array(tst, "x(+1,+1)", 1e-15, x, []float64{0, rou})

	// check derivs
	dxdu := la.NewMatrix(2, 2)
	u := la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 5)
	svals := utl.LinSpace(-1, 1, 5)
	verb := false
	for _, s := range svals {
		for _, r := range rvals {
			u[0] = r
			u[1] = s
			trf.Derivs(dxdu, x, u)
			chk.DerivVecVec(tst, "dx/dr", 1e-9, dxdu.GetDeep2(), u, 1e-3, verb, func(xx, rr []float64) {
				trf.Point(xx, rr)
			})
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{21, 21}, false, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		plt.Arc(0, 0, rin, 0, 90, &plt.A{C: plt.C(5, 9), NoClip: true, Z: 10})
		plt.Arc(0, 0, rou, 0, 90, &plt.A{C: plt.C(5, 9), NoClip: true, Z: 10})
		for _, s := range svals {
			for _, r := range rvals {
				u[0] = r
				u[1] = s
				trf.Derivs(dxdu, x, u)
				DrawArrow2dM(x, dxdu, 0, true, 0.3, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				DrawArrow2dM(x, dxdu, 1, true, 0.3, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite01")
	}
}

func TestTransfinite02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite02")

	// new mapping
	a, b := 2.0, 8.0
	trf := FactoryTfinite.Surf2dQuarterPerfLozenge(a, b)

	// auxiliary
	c := 0.5 * (a + b)
	A := a / math.Sqrt(2)
	B := b / 2.0
	C := (A + B) / 2.0
	x := la.NewVector(2)

	// check corners
	chk.Array(tst, "C0", 1e-17, trf.C[0], []float64{a, 0})
	chk.Array(tst, "C1", 1e-17, trf.C[1], []float64{b, 0})
	chk.Array(tst, "C2", 1e-17, trf.C[2], []float64{0, b})
	chk.Array(tst, "C3", 1e-17, trf.C[3], []float64{0, a})

	// check points
	trf.Point(x, []float64{-1, -1})
	chk.Array(tst, "x(-1,-1)", 1e-17, x, []float64{a, 0})

	trf.Point(x, []float64{0, -1})
	chk.Array(tst, "x( 0,-1)", 1e-17, x, []float64{c, 0})

	trf.Point(x, []float64{+1, -1})
	chk.Array(tst, "x(+1,-1)", 1e-17, x, []float64{b, 0})

	trf.Point(x, []float64{-1, 0})
	chk.Array(tst, "x(-1, 0)", 1e-15, x, []float64{A, A})

	trf.Point(x, []float64{0, 0})
	chk.Array(tst, "x( 0, 0)", 1e-15, x, []float64{C, C})

	trf.Point(x, []float64{+1, 0})
	chk.Array(tst, "x(+1, 0)", 1e-15, x, []float64{B, B})

	trf.Point(x, []float64{-1, +1})
	chk.Array(tst, "x(-1,+1)", 1e-15, x, []float64{0, a})

	trf.Point(x, []float64{0, +1})
	chk.Array(tst, "x( 0,+1)", 1e-15, x, []float64{0, c})

	trf.Point(x, []float64{+1, +1})
	chk.Array(tst, "x(+1,+1)", 1e-15, x, []float64{0, b})

	// check derivs
	dxdu := la.NewMatrix(2, 2)
	u := la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 5)
	svals := utl.LinSpace(-1, 1, 5)
	verb := false
	for _, s := range svals {
		for _, r := range rvals {
			u[0] = r
			u[1] = s
			trf.Derivs(dxdu, x, u)
			chk.DerivVecVec(tst, "dx/dr", 1e-9, dxdu.GetDeep2(), u, 1e-3, verb, func(xx, rr []float64) {
				trf.Point(xx, rr)
			})
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{11, 11}, false, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		for _, s := range svals {
			for _, r := range rvals {
				u[0] = r
				u[1] = s
				trf.Derivs(dxdu, x, u)
				DrawArrow2dM(x, dxdu, 0, true, 0.5, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				DrawArrow2dM(x, dxdu, 1, true, 0.5, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite02")
	}
}

func TestTransfinite03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite03")

	// boundary functions
	curve0 := FactoryNurbs.Curve2dExample1()
	e0 := []float64{1, 0}
	e1 := []float64{0, 1}
	knot := []float64{0}
	trf := NewTransfinite(2, []fun.Vs{

		// B0
		func(x la.Vector, r float64) {
			knot[0] = (1 + r) / 2.0
			for i := 0; i < len(x); i++ {
				curve0.Point(x, knot, 2)
			}
		},

		// B1
		func(x la.Vector, s float64) {
			x[0] = 3
			x[1] = 1.5 * (1 + s) * e1[1]
		},

		// B2
		func(x la.Vector, r float64) {
			x[0] = 1.5 * (1 + r) * e0[0]
			x[1] = 3
		},

		// B3
		func(x la.Vector, s float64) {
			x[0] = 0
			x[1] = 1.5 * (1 + s) * e1[1]
		},
	}, []fun.Vs{

		// dB0/dr
		func(dxdr la.Vector, r float64) {
			knot[0] = (1 + r) / 2.0
			dCdu := la.NewMatrix(2, curve0.Gnd())
			C := la.NewVector(2)
			curve0.PointDeriv(dCdu, C, knot, 2)
			for i := 0; i < 2; i++ {
				dxdr[i] = dCdu.Get(i, 0) * 0.5
			}
		},

		// dB1/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = 0
			dxds[1] = 1.5 * e1[1]
		},

		// dB2/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 1.5 * e0[0]
			dxdr[1] = 0
		},

		// dB3/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = 0
			dxds[1] = 1.5 * e1[1]
		},
	})

	// auxiliary
	xtmp := la.NewVector(2)
	dxdr := la.NewVector(2)
	dxds := la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 9)
	svals := utl.LinSpace(-1, 1, 9)

	// check: dB0/dr
	//verb := chk.Verbose
	verb := false
	for _, r := range rvals {
		trf.Bd[0](dxdr, r)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB0_%d/dr", i), 1e-10, dxdr[i], r, 1e-3, verb, func(s float64) float64 {
				trf.B[0](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB1/ds
	io.Pl()
	for _, s := range svals {
		trf.Bd[1](dxds, s)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB1_%d/ds", i), 1e-12, dxds[i], s, 1e-3, verb, func(s float64) float64 {
				trf.B[1](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB2/dr
	io.Pl()
	for _, r := range rvals {
		trf.Bd[2](dxdr, r)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB2_%d/dr", i), 1e-12, dxdr[i], r, 1e-3, verb, func(s float64) float64 {
				trf.B[2](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB3/ds
	io.Pl()
	for _, s := range svals {
		trf.Bd[3](dxds, s)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB3_%d/ds", i), 1e-12, dxds[i], s, 1e-3, verb, func(s float64) float64 {
				trf.B[3](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check derivs
	dxdu := la.NewMatrix(2, 2)
	x := la.NewVector(2)
	u := la.NewVector(2)
	for _, s := range svals {
		for _, r := range rvals {
			u[0] = r
			u[1] = s
			trf.Derivs(dxdu, x, u)
			chk.DerivVecVec(tst, "dx/dr", 1e-9, dxdu.GetDeep2(), u, 1e-3, verb, func(xx, rr []float64) {
				trf.Point(xx, rr)
			})
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Prop: 1, Eps: true})
		//curve0.DrawElems(2, 41, false, &plt.A{C: plt.C(2, 0), Z: 10}, nil)
		trf.Draw([]int{21, 21}, false, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		for _, s := range svals {
			for _, r := range rvals {
				u[0] = r
				u[1] = s
				trf.Derivs(dxdu, x, u)
				DrawArrow2dM(x, dxdu, 0, true, 0.13, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				DrawArrow2dM(x, dxdu, 1, true, 0.13, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.AxisOff()
		plt.Equal()
		plt.AxisRange(-0.1, 3.2, -0.1, 3.2)
		plt.Save("/tmp/gosl/gm", "transfinite03")
	}
}
