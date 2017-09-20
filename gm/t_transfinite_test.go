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

	π := math.Pi
	e0 := []float64{1, 0}
	e1 := []float64{0, 1}

	trf := NewTransfinite(2, []fun.Vs{
		func(x la.Vector, ξ float64) {
			for i := 0; i < 2; i++ {
				x[i] = (2 + ξ) * e0[i]
			}
		},
		func(x la.Vector, η float64) {
			θ := π * (η + 1) / 4.0
			for i := 0; i < 2; i++ {
				x[i] = 3*math.Cos(θ)*e0[i] + 3*math.Sin(θ)*e1[i]
			}
		},
		func(x la.Vector, ξ float64) {
			for i := 0; i < 2; i++ {
				x[i] = (2 + ξ) * e1[i]
			}
		},
		func(x la.Vector, η float64) {
			θ := π * (η + 1) / 4.0
			for i := 0; i < 2; i++ {
				x[i] = math.Cos(θ)*e0[i] + math.Sin(θ)*e1[i]
			}
		},
	}, []fun.Vs{
		func(dxdξ la.Vector, ξ float64) {
			for i := 0; i < 2; i++ {
				dxdξ[i] = e0[i]
			}
		},
		func(dxdη la.Vector, η float64) {
			θ := π * (η + 1) / 4.0
			dθdη := π / 4.0
			for i := 0; i < 2; i++ {
				dxdη[i] = (-3*math.Sin(θ)*e0[i] + 3*math.Cos(θ)*e1[i]) * dθdη
			}
		},
		func(dxdξ la.Vector, ξ float64) {
			for i := 0; i < 2; i++ {
				dxdξ[i] = e1[i]
			}
		},
		func(dxdη la.Vector, η float64) {
			θ := π * (η + 1) / 4.0
			dθdη := π / 4.0
			for i := 0; i < 2; i++ {
				dxdη[i] = (-math.Sin(θ)*e0[i] + math.Cos(θ)*e1[i]) * dθdη
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

	trf.Point(x, []float64{-1, -1})
	chk.Array(tst, "x(-1,-1)", 1e-17, x, []float64{1, 0})

	trf.Point(x, []float64{0, -1})
	chk.Array(tst, "x( 0,-1)", 1e-17, x, []float64{2, 0})

	trf.Point(x, []float64{+1, -1})
	chk.Array(tst, "x(+1,-1)", 1e-17, x, []float64{3, 0})

	trf.Point(x, []float64{-1, 0})
	chk.Array(tst, "x(-1, 0)", 1e-15, x, []float64{a, a})

	trf.Point(x, []float64{0, 0})
	chk.Array(tst, "x( 0, 0)", 1e-15, x, []float64{b, b})

	trf.Point(x, []float64{+1, 0})
	chk.Array(tst, "x(+1, 0)", 1e-15, x, []float64{c, c})

	trf.Point(x, []float64{-1, +1})
	chk.Array(tst, "x(-1,+1)", 1e-15, x, []float64{0, 1})

	trf.Point(x, []float64{0, +1})
	chk.Array(tst, "x( 0,+1)", 1e-15, x, []float64{0, 2})

	trf.Point(x, []float64{+1, +1})
	chk.Array(tst, "x(+1,+1)", 1e-15, x, []float64{0, 3})

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
	}, []fun.Vs{
		func(dxdξ la.Vector, ξ float64) {
			for i := 0; i < 2; i++ {
				dxdξ[i] = e0[i]
			}
		},
		func(dxdη la.Vector, η float64) {
			θ := π * (η + 1) / 4.0
			dθdη := π / 4.0
			for i := 0; i < 2; i++ {
				dxdη[i] = (-3*math.Sin(θ)*e0[i] + 3*math.Cos(θ)*e1[i]) * dθdη
			}
		},
		func(dxdξ la.Vector, ξ float64) {
			for i := 0; i < 2; i++ {
				dxdξ[i] = e1[i]
			}
		},
		func(dxdη la.Vector, η float64) {
			θ := π * (η + 1) / 4.0
			dθdη := π / 4.0
			for i := 0; i < 2; i++ {
				dxdη[i] = (-math.Sin(θ)*e0[i] + math.Cos(θ)*e1[i]) * dθdη
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

	trf.Point(x, []float64{-1, -1})
	chk.Array(tst, "x(-1,-1)", 1e-17, x, []float64{1, 0})

	trf.Point(x, []float64{0, -1})
	chk.Array(tst, "x( 0,-1)", 1e-17, x, []float64{2, 0})

	trf.Point(x, []float64{+1, -1})
	chk.Array(tst, "x(+1,-1)", 1e-17, x, []float64{3, 0})

	trf.Point(x, []float64{-1, 0})
	chk.Array(tst, "x(-1, 0)", 1e-15, x, []float64{a, a})

	trf.Point(x, []float64{0, 0})
	chk.Array(tst, "x( 0, 0)", 1e-15, x, []float64{b, b})

	trf.Point(x, []float64{+1, 0})
	chk.Array(tst, "x(+1, 0)", 1e-15, x, []float64{c, c})

	trf.Point(x, []float64{-1, +1})
	chk.Array(tst, "x(-1,+1)", 1e-15, x, []float64{0, 1})

	trf.Point(x, []float64{0, +1})
	chk.Array(tst, "x( 0,+1)", 1e-15, x, []float64{0, 2})

	trf.Point(x, []float64{+1, +1})
	chk.Array(tst, "x(+1,+1)", 1e-15, x, []float64{0, 3})

	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{21, 21}, nil, nil)
		plt.Arc(0, 0, 1, 0, 90, &plt.A{C: plt.C(2, 0), NoClip: true, Z: 10})
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
	u := []float64{0} // knots
	trf := NewTransfinite(2, []fun.Vs{

		// Γ0
		func(x la.Vector, ξ float64) {
			u[0] = (1 + ξ) / 2.0
			for i := 0; i < len(x); i++ {
				curve0.Point(x, u, 2)
			}
		},

		// Γ1
		func(x la.Vector, η float64) {
			x[0] = 3
			x[1] = 1.5 * (1 + η) * e1[1]
		},

		// Γ2
		func(x la.Vector, ξ float64) {
			x[0] = 1.5 * (1 + ξ) * e0[0]
			x[1] = 3
		},

		// Γ3
		func(x la.Vector, η float64) {
			x[0] = 0
			x[1] = 1.5 * (1 + η) * e1[1]
		},
	}, []fun.Vs{

		// dΓ0/dξ
		func(dxdξ la.Vector, ξ float64) {
			u[0] = (1 + ξ) / 2.0
			dCdu := la.NewMatrix(2, curve0.Gnd())
			C := la.NewVector(2)
			curve0.PointDeriv(dCdu, C, u, 2)
			for i := 0; i < 2; i++ {
				dxdξ[i] = dCdu.Get(i, 0) * 0.5
			}
		},

		// dΓ1/dη
		func(dxdη la.Vector, η float64) {
			dxdη[0] = 0
			dxdη[1] = 1.5 * e1[1]
		},

		// dΓ2/dξ
		func(dxdξ la.Vector, ξ float64) {
			dxdξ[0] = 1.5 * e0[0]
			dxdξ[1] = 0
		},

		// dΓ3/dη
		func(dxdη la.Vector, η float64) {
			dxdη[0] = 0
			dxdη[1] = 1.5 * e1[1]
		},
	})

	// auxiliary
	xtmp := la.NewVector(2)
	dxdξ := la.NewVector(2)
	dxdη := la.NewVector(2)
	ξs := utl.LinSpace(-1, 1, 5)
	ηs := utl.LinSpace(-1, 1, 5)

	// check: dΓ0/dξ
	//verb := chk.Verbose
	verb := false
	for _, ξ := range ξs {
		trf.Γd[0](dxdξ, ξ)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dΓ0_%d/dξ", i), 1e-10, dxdξ[i], ξ, 1e-3, verb, func(s float64) float64 {
				trf.Γ[0](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dΓ1/dη
	io.Pl()
	for _, η := range ηs {
		trf.Γd[1](dxdη, η)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dΓ1_%d/dη", i), 1e-12, dxdη[i], η, 1e-3, verb, func(s float64) float64 {
				trf.Γ[1](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dΓ2/dξ
	io.Pl()
	for _, ξ := range ξs {
		trf.Γd[2](dxdξ, ξ)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dΓ2_%d/dξ", i), 1e-12, dxdξ[i], ξ, 1e-3, verb, func(s float64) float64 {
				trf.Γ[2](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dΓ3/dη
	io.Pl()
	for _, η := range ηs {
		trf.Γd[3](dxdη, η)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dΓ3_%d/dη", i), 1e-12, dxdη[i], η, 1e-3, verb, func(s float64) float64 {
				trf.Γ[3](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check derivs
	dxdr := la.NewMatrix(2, 2)
	x := la.NewVector(2)
	r := la.NewVector(2)
	for _, η := range ηs {
		for _, ξ := range ξs {
			r[0] = ξ
			r[1] = η
			trf.Derivs(dxdr, x, r)
			chk.DerivVecVec(tst, "dx/dr", 1e-9, dxdr.GetDeep2(), r, 1e-3, verb, func(xx, rr []float64) {
				trf.Point(xx, rr)
			})
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		curve0.DrawElems(2, 41, false, &plt.A{C: plt.C(2, 0), Z: 10}, nil)
		trf.Draw([]int{21, 21}, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		for _, η := range ηs {
			for _, ξ := range ξs {
				r[0] = ξ
				r[1] = η
				trf.Derivs(dxdr, x, r)
				DrawArrow2dM(x, dxdr, 0, true, 0.3, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				DrawArrow2dM(x, dxdr, 1, true, 0.3, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite03")
	}
}
