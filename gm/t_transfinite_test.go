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
	chk.PrintTitle("Transfinite01. 2d ring")

	// new mapping
	rin, rou := 2.0, 6.0 // radii
	trf := FactoryTfinite.Surf2dQuarterRing(rin, rou)

	// check corners
	chk.Array(tst, "p0", 1e-15, trf.p0, []float64{rin, 0})
	chk.Array(tst, "p1", 1e-15, trf.p1, []float64{rou, 0})
	chk.Array(tst, "p2", 1e-15, trf.p2, []float64{0, rou})
	chk.Array(tst, "p3", 1e-15, trf.p3, []float64{0, rin})

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
	dxDr, dxDs := la.NewVector(2), la.NewVector(2)
	ddxDrr, ddxDss, ddxDrs := la.NewVector(2), la.NewVector(2), la.NewVector(2)
	u, utmp, tmp := la.NewVector(2), la.NewVector(2), la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 5)
	svals := utl.LinSpace(-1, 1, 5)
	verb := chk.Verbose
	for _, s := range svals {
		for _, r := range rvals {
			u[0], u[1] = r, s
			trf.PointAndDerivs(x, dxDr, dxDs, nil, ddxDrr, ddxDss, nil, ddxDrs, nil, nil, u)
			chk.DerivVecSca(tst, "dx/dr   ", 1e-10, dxDr, r, 1e-3, verb, func(xx []float64, ξ float64) {
				utmp[0], utmp[1] = ξ, u[1]
				trf.Point(xx, utmp)
			})
			chk.DerivVecSca(tst, "dx/ds   ", 1e-9, dxDs, s, 1e-3, verb, func(xx []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.Point(xx, utmp)
			})
			chk.DerivVecSca(tst, "d²x/dr² ", 1e-10, ddxDrr, r, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = ξ, u[1]
				trf.PointAndDerivs(tmp, d, tmp, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			chk.DerivVecSca(tst, "d²x/ds² ", 1e-10, ddxDss, s, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.PointAndDerivs(tmp, tmp, d, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			chk.DerivVecSca(tst, "d²x/drds", 1e-10, ddxDrs, s, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.PointAndDerivs(tmp, d, tmp, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			if verb {
				io.Pl()
			}
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
				u[0], u[1] = r, s
				trf.PointAndDerivs(x, dxDr, dxDs, nil, nil, nil, nil, nil, nil, nil, u)
				plt.DrawArrow2d(x, dxDr, true, 0.3, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				plt.DrawArrow2d(x, dxDs, true, 0.3, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite01")
	}
}

func TestTransfinite02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite02. 2d lozenge with hole")

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
	chk.Array(tst, "p0", 1e-15, trf.p0, []float64{a, 0})
	chk.Array(tst, "p1", 1e-15, trf.p1, []float64{b, 0})
	chk.Array(tst, "p2", 1e-15, trf.p2, []float64{0, b})
	chk.Array(tst, "p3", 1e-15, trf.p3, []float64{0, a})

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
	dxDr, dxDs := la.NewVector(2), la.NewVector(2)
	ddxDrr, ddxDss, ddxDrs := la.NewVector(2), la.NewVector(2), la.NewVector(2)
	u, utmp, tmp := la.NewVector(2), la.NewVector(2), la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 5)
	svals := utl.LinSpace(-1, 1, 5)
	verb := chk.Verbose
	for _, s := range svals {
		for _, r := range rvals {
			u[0], u[1] = r, s
			trf.PointAndDerivs(x, dxDr, dxDs, nil, ddxDrr, ddxDss, nil, ddxDrs, nil, nil, u)
			chk.DerivVecSca(tst, "dx/dr   ", 1e-10, dxDr, r, 1e-3, verb, func(xx []float64, ξ float64) {
				utmp[0], utmp[1] = ξ, u[1]
				trf.Point(xx, utmp)
			})
			chk.DerivVecSca(tst, "dx/ds   ", 1e-10, dxDs, s, 1e-3, verb, func(xx []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.Point(xx, utmp)
			})
			chk.DerivVecSca(tst, "d²x/dr² ", 1e-10, ddxDrr, r, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = ξ, u[1]
				trf.PointAndDerivs(tmp, d, tmp, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			chk.DerivVecSca(tst, "d²x/ds² ", 1e-10, ddxDss, s, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.PointAndDerivs(tmp, tmp, d, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			chk.DerivVecSca(tst, "d²x/drds", 1e-10, ddxDrs, s, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.PointAndDerivs(tmp, d, tmp, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			if verb {
				io.Pl()
			}
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{11, 11}, false, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		for _, s := range svals {
			for _, r := range rvals {
				u[0], u[1] = r, s
				trf.PointAndDerivs(x, dxDr, dxDs, nil, nil, nil, nil, nil, nil, nil, u)
				plt.DrawArrow2d(x, dxDr, true, 0.3, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				plt.DrawArrow2d(x, dxDs, true, 0.3, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite02")
	}
}

func TestTransfinite03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite03. 2d square with NURBS")

	// boundary functions
	curve0 := FactoryNurbs.Curve2dExample1()
	knot := []float64{0}
	trf := NewTransfinite2d([]fun.Vs{

		// B0
		func(x la.Vector, s float64) {
			x[0] = 0.0
			x[1] = 1.5 * (1 + s)
		},

		// B1
		func(x la.Vector, s float64) {
			x[0] = 3.0
			x[1] = 1.5 * (1 + s)
		},

		// B2
		func(x la.Vector, r float64) {
			knot[0] = (1 + r) / 2.0
			curve0.Point(x, knot, 2)
		},

		// B3
		func(x la.Vector, r float64) {
			x[0] = 1.5 * (1 + r)
			x[1] = 3.0
		},

		// first order derivatives

	}, []fun.Vs{

		// dB0/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = 0.0
			dxds[1] = 1.5
		},

		// dB1/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = 0.0
			dxds[1] = 1.5
		},

		// dB2/dr
		func(dxdr la.Vector, r float64) {
			knot[0] = (1.0 + r) / 2.0
			dCdu := la.NewMatrix(2, curve0.Gnd())
			C := la.NewVector(2)
			curve0.PointAndFirstDerivs(dCdu, C, knot, 2)
			for i := 0; i < 2; i++ {
				dxdr[i] = dCdu.Get(i, 0) * 0.5
			}
		},

		// dB3/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 1.5
			dxdr[1] = 0.0
		},

		// second order derivatives

	}, []fun.Vs{

		// d²B[0]/ds²
		func(ddxdss la.Vector, s float64) {
			ddxdss[0] = 0.0
			ddxdss[1] = 0.0
		},

		// d²B[1]/ds²
		func(ddxdss la.Vector, s float64) {
			ddxdss[0] = 0.0
			ddxdss[1] = 0.0
		},

		// d²B[2]/dr²
		func(ddxdrr la.Vector, r float64) {
			knot[0] = (1.0 + r) / 2.0
			dCdu := la.NewMatrix(2, curve0.Gnd())
			C := la.NewVector(2)
			// TODO: fix this
			curve0.PointAndFirstDerivs(dCdu, C, knot, 2)
			for i := 0; i < 2; i++ {
				ddxdrr[i] = dCdu.Get(i, 0) * 0.5
			}
		},

		// d²B[3]/dr²
		func(ddxdrr la.Vector, r float64) {
			ddxdrr[0] = 0.0
			ddxdrr[1] = 0.0
		},
	})

	// auxiliary
	xtmp := la.NewVector(2)
	dxdr := la.NewVector(2)
	dxds := la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 9)
	svals := utl.LinSpace(-1, 1, 9)

	// check: dB0/dr
	verb := false
	for _, r := range rvals {
		trf.ed[0](dxdr, r)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB0_%d/dr", i), 1e-10, dxdr[i], r, 1e-3, verb, func(s float64) float64 {
				trf.e[0](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB1/ds
	io.Pl()
	for _, s := range svals {
		trf.ed[1](dxds, s)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB1_%d/ds", i), 1e-12, dxds[i], s, 1e-3, verb, func(s float64) float64 {
				trf.e[1](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB2/dr
	io.Pl()
	for _, r := range rvals {
		trf.ed[2](dxdr, r)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB2_%d/dr", i), 1e-10, dxdr[i], r, 1e-3, verb, func(s float64) float64 {
				trf.e[2](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB3/ds
	io.Pl()
	for _, s := range svals {
		trf.ed[3](dxds, s)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB3_%d/ds", i), 1e-12, dxds[i], s, 1e-3, verb, func(s float64) float64 {
				trf.e[3](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check derivs
	verb = chk.Verbose
	dxDr, dxDs := la.NewVector(2), la.NewVector(2)
	ddxDrr, ddxDss, ddxDrs := la.NewVector(2), la.NewVector(2), la.NewVector(2)
	x, u, utmp, tmp := la.NewVector(2), la.NewVector(2), la.NewVector(2), la.NewVector(2)
	for _, s := range svals {
		for _, r := range rvals {
			u[0], u[1] = r, s
			trf.PointAndDerivs(x, dxDr, dxDs, nil, ddxDrr, ddxDss, nil, ddxDrs, nil, nil, u)
			chk.DerivVecSca(tst, "dx/dr   ", 1e-9, dxDr, r, 1e-3, verb, func(xx []float64, ξ float64) {
				utmp[0], utmp[1] = ξ, u[1]
				trf.Point(xx, utmp)
			})
			chk.DerivVecSca(tst, "dx/ds   ", 1e-10, dxDs, s, 1e-3, verb, func(xx []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.Point(xx, utmp)
			})
			// TODO: implement this
			if false {
				chk.DerivVecSca(tst, "d²x/dr² ", 1e-10, ddxDrr, r, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1] = ξ, u[1]
					trf.PointAndDerivs(tmp, d, tmp, tmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/ds² ", 1e-10, ddxDss, s, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1] = u[0], ξ
					trf.PointAndDerivs(tmp, tmp, d, tmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/drds", 1e-10, ddxDrs, s, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1] = u[0], ξ
					trf.PointAndDerivs(tmp, d, tmp, tmp, nil, nil, nil, nil, nil, nil, utmp)
				})
			}
			if verb {
				io.Pl()
			}
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400})
		//plt.Reset(true, &plt.A{WidthPt: 400, Prop: 1, Eps: true})
		//curve0.DrawElems(2, 41, false, &plt.A{C: plt.C(2, 0), Z: 10}, nil)
		trf.Draw([]int{21, 21}, false, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		for _, s := range svals {
			for _, r := range rvals {
				u[0], u[1] = r, s
				trf.PointAndDerivs(x, dxDr, dxDs, nil, nil, nil, nil, nil, nil, nil, u)
				plt.DrawArrow2d(x, dxDr, true, 0.15, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				plt.DrawArrow2d(x, dxDs, true, 0.15, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.AxisOff()
		plt.Equal()
		plt.AxisRange(-0.1, 3.2, -0.1, 3.2)
		plt.Save("/tmp/gosl/gm", "transfinite03")
	}
}

func TestTransfinite04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite04. 3d cube")

	// new mapping
	trf := FactoryTfinite.SolidCube(1, 1, 1)

	// check corners
	chk.Array(tst, "p0", 1e-15, trf.p0, []float64{0, 0, 0})
	chk.Array(tst, "p1", 1e-15, trf.p1, []float64{1, 0, 0})
	chk.Array(tst, "p2", 1e-15, trf.p2, []float64{1, 1, 0})
	chk.Array(tst, "p3", 1e-15, trf.p3, []float64{0, 1, 0})
	chk.Array(tst, "p4", 1e-15, trf.p4, []float64{0, 0, 1})
	chk.Array(tst, "p5", 1e-15, trf.p5, []float64{1, 0, 1})
	chk.Array(tst, "p6", 1e-15, trf.p6, []float64{1, 1, 1})
	chk.Array(tst, "p7", 1e-15, trf.p7, []float64{0, 1, 1})

	// auxiliary
	verb := chk.Verbose
	rvals := utl.LinSpace(-1, 1, 3)
	svals := utl.LinSpace(-1, 1, 3)
	tvals := utl.LinSpace(-1, 1, 3)

	// check derivs
	dxDr, dxDs, dxDt := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	ddxDrr, ddxDss, ddxDtt := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	ddxDrs, ddxDrt, ddxDst := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	x, u, utmp := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	xtmp, dxDrTmp, dxDsTmp, dxDtTmp := la.NewVector(3), la.NewVector(3), la.NewVector(3), la.NewVector(3)
	for _, t := range tvals {
		for _, s := range svals {
			for _, r := range rvals {
				u[0], u[1], u[2] = r, s, t
				trf.PointAndDerivs(x, dxDr, dxDs, dxDt, ddxDrr, ddxDss, ddxDtt, ddxDrs, ddxDrt, ddxDst, u)
				chk.DerivVecSca(tst, "dx/dr   ", 1e-10, dxDr, r, 1e-3, verb, func(xx []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = ξ, u[1], u[2]
					trf.Point(xx, utmp)
				})
				chk.DerivVecSca(tst, "dx/ds   ", 1e-10, dxDs, s, 1e-3, verb, func(xx []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], ξ, u[2]
					trf.Point(xx, utmp)
				})
				chk.DerivVecSca(tst, "dx/dt   ", 1e-10, dxDt, t, 1e-3, verb, func(xx []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], u[1], ξ
					trf.Point(xx, utmp)
				})
				chk.DerivVecSca(tst, "d²x/dr² ", 1e-10, ddxDrr, r, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = ξ, u[1], u[2]
					trf.PointAndDerivs(xtmp, d, dxDsTmp, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/ds² ", 1e-10, ddxDss, s, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], ξ, u[2]
					trf.PointAndDerivs(xtmp, dxDrTmp, d, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/dt² ", 1e-10, ddxDtt, t, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], u[1], ξ
					trf.PointAndDerivs(xtmp, dxDrTmp, dxDsTmp, d, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/drds", 1e-10, ddxDrs, s, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], ξ, u[2]
					trf.PointAndDerivs(xtmp, d, dxDsTmp, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/drdt", 1e-10, ddxDrt, t, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], u[1], ξ
					trf.PointAndDerivs(xtmp, d, dxDsTmp, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/dsdt", 1e-10, ddxDst, t, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], u[1], ξ
					trf.PointAndDerivs(xtmp, dxDrTmp, d, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				if verb {
					io.Pl()
				}
			}
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400})
		trf.Draw([]int{3, 3, 3}, false, nil, nil)
		plt.Default3dView(0, 1, 0, 1, 0, 1, true)
		//plt.Show()
		plt.Save("/tmp/gosl/gm", "transfinite04")
	}
}

func TestTransfinite05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite05. 3d ring")

	// new mapping
	a, b, h := 2.0, 3.0, 1.0 // radii and thickness
	trf := FactoryTfinite.SolidQuarterRing(a, b, h)

	// check corners
	chk.Array(tst, "p0", 1e-15, trf.p0, []float64{0, a, 0})
	chk.Array(tst, "p1", 1e-15, trf.p1, []float64{h, a, 0})
	chk.Array(tst, "p2", 1e-15, trf.p2, []float64{h, b, 0})
	chk.Array(tst, "p3", 1e-15, trf.p3, []float64{0, b, 0})
	chk.Array(tst, "p4", 1e-15, trf.p4, []float64{0, 0, a})
	chk.Array(tst, "p5", 1e-15, trf.p5, []float64{h, 0, a})
	chk.Array(tst, "p6", 1e-15, trf.p6, []float64{h, 0, b})
	chk.Array(tst, "p7", 1e-15, trf.p7, []float64{0, 0, b})

	// auxiliary
	verb := chk.Verbose
	rvals := utl.LinSpace(-1, 1, 3)
	svals := utl.LinSpace(-1, 1, 3)
	tvals := utl.LinSpace(-1, 1, 3)

	// check derivs
	dxDr, dxDs, dxDt := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	ddxDrr, ddxDss, ddxDtt := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	ddxDrs, ddxDrt, ddxDst := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	x, u, utmp := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	xtmp, dxDrTmp, dxDsTmp, dxDtTmp := la.NewVector(3), la.NewVector(3), la.NewVector(3), la.NewVector(3)
	for _, t := range tvals {
		for _, s := range svals {
			for _, r := range rvals {
				u[0], u[1], u[2] = r, s, t
				trf.PointAndDerivs(x, dxDr, dxDs, dxDt, ddxDrr, ddxDss, ddxDtt, ddxDrs, ddxDrt, ddxDst, u)
				chk.DerivVecSca(tst, "dx/dr   ", 1e-10, dxDr, r, 1e-3, verb, func(xx []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = ξ, u[1], u[2]
					trf.Point(xx, utmp)
				})
				chk.DerivVecSca(tst, "dx/ds   ", 1e-10, dxDs, s, 1e-3, verb, func(xx []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], ξ, u[2]
					trf.Point(xx, utmp)
				})
				chk.DerivVecSca(tst, "dx/dt   ", 1e-10, dxDt, t, 1e-3, verb, func(xx []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], u[1], ξ
					trf.Point(xx, utmp)
				})
				chk.DerivVecSca(tst, "d²x/dr² ", 1e-10, ddxDrr, r, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = ξ, u[1], u[2]
					trf.PointAndDerivs(xtmp, d, dxDsTmp, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/ds² ", 1e-10, ddxDss, s, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], ξ, u[2]
					trf.PointAndDerivs(xtmp, dxDrTmp, d, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/dt² ", 1e-10, ddxDtt, t, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], u[1], ξ
					trf.PointAndDerivs(xtmp, dxDrTmp, dxDsTmp, d, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/drds", 1e-10, ddxDrs, s, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], ξ, u[2]
					trf.PointAndDerivs(xtmp, d, dxDsTmp, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/drdt", 1e-10, ddxDrt, t, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], u[1], ξ
					trf.PointAndDerivs(xtmp, d, dxDsTmp, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				chk.DerivVecSca(tst, "d²x/dsdt", 1e-10, ddxDst, t, 1e-3, verb, func(d []float64, ξ float64) {
					utmp[0], utmp[1], utmp[2] = u[0], u[1], ξ
					trf.PointAndDerivs(xtmp, dxDrTmp, d, dxDtTmp, nil, nil, nil, nil, nil, nil, utmp)
				})
				if verb {
					io.Pl()
				}
			}
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400})
		trf.Draw([]int{5, 5, 11}, true, nil, nil)
		plt.Default3dView(0, 3, 0, 3, 0, 3, true)
		//plt.Show()
		plt.Save("/tmp/gosl/gm", "transfinite05")
	}
}

func TestTransfinite06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite06. quadrilateral")

	// new mapping
	A := []float64{0, 0}
	B := []float64{4, 1}
	C := []float64{3, 3}
	D := []float64{-1, 4}
	trf := FactoryTfinite.Surf2dQuad(A, B, C, D)

	// check corners
	chk.Array(tst, "p0", 1e-15, trf.p0, A)
	chk.Array(tst, "p1", 1e-15, trf.p1, B)
	chk.Array(tst, "p2", 1e-15, trf.p2, C)
	chk.Array(tst, "p3", 1e-15, trf.p3, D)

	// check derivs
	x := la.NewVector(2)
	dxDr, dxDs := la.NewVector(2), la.NewVector(2)
	ddxDrr, ddxDss, ddxDrs := la.NewVector(2), la.NewVector(2), la.NewVector(2)
	u, utmp, tmp := la.NewVector(2), la.NewVector(2), la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 5)
	svals := utl.LinSpace(-1, 1, 5)
	verb := chk.Verbose
	for _, s := range svals {
		for _, r := range rvals {
			u[0], u[1] = r, s
			trf.PointAndDerivs(x, dxDr, dxDs, nil, ddxDrr, ddxDss, nil, ddxDrs, nil, nil, u)
			chk.DerivVecSca(tst, "dx/dr   ", 1e-11, dxDr, r, 1e-3, verb, func(xx []float64, ξ float64) {
				utmp[0], utmp[1] = ξ, u[1]
				trf.Point(xx, utmp)
			})
			chk.DerivVecSca(tst, "dx/ds   ", 1e-11, dxDs, s, 1e-3, verb, func(xx []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.Point(xx, utmp)
			})
			chk.DerivVecSca(tst, "d²x/dr² ", 1e-11, ddxDrr, r, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = ξ, u[1]
				trf.PointAndDerivs(tmp, d, tmp, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			chk.DerivVecSca(tst, "d²x/ds² ", 1e-11, ddxDss, s, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.PointAndDerivs(tmp, tmp, d, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			chk.DerivVecSca(tst, "d²x/drds", 1e-11, ddxDrs, s, 1e-3, verb, func(d []float64, ξ float64) {
				utmp[0], utmp[1] = u[0], ξ
				trf.PointAndDerivs(tmp, d, tmp, tmp, nil, nil, nil, nil, nil, nil, utmp)
			})
			if verb {
				io.Pl()
			}
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{11, 11}, false, &plt.A{C: plt.C(2, 9), NoClip: true}, &plt.A{C: plt.C(3, 9), Lw: 2, NoClip: true})
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite06")
	}
}
