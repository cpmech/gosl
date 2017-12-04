// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

// check NURBS derivatives /////////////////////////////////////////////////////////////////////////

func checkNurbsBasisDerivs(tst *testing.T, b *Nurbs, npts int, tol float64, verb bool) {
	δ := 1e-1 // used to avoid using central differences @ boundaries of t in [0,5]
	dana := make([]float64, 2)
	uu := make([]float64, 2)
	for _, u := range utl.LinSpace(b.b[0].tmin+δ, b.b[0].tmax-δ, npts) {
		for _, v := range utl.LinSpace(b.b[1].tmin+δ, b.b[1].tmax-δ, npts) {
			uu[0], uu[1] = u, v
			b.CalcBasisAndDerivs(uu)
			for i := 0; i < b.n[0]; i++ {
				for j := 0; j < b.n[1]; j++ {
					l := i + j*b.n[0]
					b.GetDerivL(dana, l)
					chk.DerivScaVec(tst, io.Sf("dR%d(%.3f,%.3f)", l, uu[0], uu[1]), tol, dana, uu, 1e-1, verb, func(x []float64) float64 {
						return b.RecursiveBasis(x, l)
					})
				}
			}
		}
	}
}

func checkNurbsCurveDerivs(tst *testing.T, curve *Nurbs, uvals []float64, verb bool) {
	ndim := 2
	u := la.NewVector(1)
	x := la.NewVector(ndim)
	C := la.NewVector(ndim)
	tmp := la.NewVector(ndim)
	dCdu := la.NewMatrix(ndim, curve.gnd)
	dxdr, ddxdrr := la.NewVector(ndim), la.NewVector(ndim)
	for i := 0; i < len(uvals); i++ {
		u[0] = uvals[i]
		curve.PointAndDerivs(x, dxdr, nil, nil, ddxdrr, nil, nil, nil, nil, nil, u, ndim)
		curve.PointAndFirstDerivs(dCdu, C, u, ndim)
		chk.Array(tst, io.Sf("x      (%.2f)", u[0]), 1e-13, x, C)
		chk.Array(tst, io.Sf("dC/du0 (%.2f)", u[0]), 1e-13, dxdr, dCdu.GetCol(0))
		chk.DerivVecSca(tst, io.Sf("dx/dr  (%.2f)", u[0]), 1e-7, dxdr, u[0], 1e-6, verb, func(xx []float64, r float64) {
			curve.Point(xx, []float64{r}, ndim)
		})
		chk.DerivVecSca(tst, io.Sf("d²x/dr²(%.2f)", u[0]), 1e-7, ddxdrr, u[0], 1e-6, verb, func(xx []float64, r float64) {
			curve.PointAndDerivs(tmp, xx, nil, nil, nil, nil, nil, nil, nil, nil, []float64{r}, ndim)
		})
		if verb {
			io.Pl()
		}
	}
}

func checkNurbsSurfDerivs(tst *testing.T, surf *Nurbs, uvals, vvals []float64, verb bool, tol0, tol1, tol2, tol3 float64) {
	ndim := 2
	u := la.NewVector(2)
	U := la.NewVector(2)
	x := la.NewVector(ndim)
	C := la.NewVector(ndim)
	tmp1 := la.NewVector(ndim)
	tmp2 := la.NewVector(ndim)
	dCdu := la.NewMatrix(ndim, surf.gnd)
	dxdr, dxds, ddxdrr, ddxdss, ddxdrs := la.NewVector(ndim), la.NewVector(ndim), la.NewVector(ndim), la.NewVector(ndim), la.NewVector(ndim)
	for i := 0; i < len(uvals); i++ {
		for j := 0; j < len(vvals); j++ {
			u[0], u[1] = uvals[i], vvals[j]
			if verb {
				io.Pf("\nu = %f\n", u)
			}
			surf.PointAndDerivs(x, dxdr, dxds, nil, ddxdrr, ddxdss, nil, ddxdrs, nil, nil, u, ndim)
			surf.PointAndFirstDerivs(dCdu, C, u, ndim)
			chk.Array(tst, "x     ", tol0, x, C)
			chk.Array(tst, "dC/du0", tol0, dxdr, dCdu.GetCol(0))
			chk.Array(tst, "dC/du1", tol0, dxds, dCdu.GetCol(1))
			chk.DerivVecSca(tst, "dx/dr    ", tol1, dxdr, u[0], 1e-6, verb, func(xx []float64, r float64) {
				U[0], U[1] = r, u[1]
				surf.Point(xx, U, ndim)
			})
			chk.DerivVecSca(tst, "dx/ds    ", tol1, dxds, u[1], 1e-6, verb, func(xx []float64, s float64) {
				U[0], U[1] = u[0], s
				surf.Point(xx, U, ndim)
			})
			chk.DerivVecSca(tst, "d²x/dr²  ", tol2, ddxdrr, u[0], 1e-6, verb, func(xx []float64, r float64) {
				U[0], U[1] = r, u[1]
				surf.PointAndDerivs(tmp1, xx, tmp2, nil, nil, nil, nil, nil, nil, nil, U, ndim)
			})
			chk.DerivVecSca(tst, "d²x/ds²  ", tol2, ddxdss, u[1], 1e-6, verb, func(xx []float64, s float64) {
				U[0], U[1] = u[0], s
				surf.PointAndDerivs(tmp1, tmp2, xx, nil, nil, nil, nil, nil, nil, nil, U, ndim)
			})
			chk.DerivVecSca(tst, "d²x/drds ", tol2, ddxdrs, u[1], 1e-6, verb, func(xx []float64, s float64) {
				U[0], U[1] = u[0], s
				surf.PointAndDerivs(tmp1, xx, tmp2, nil, nil, nil, nil, nil, nil, nil, U, ndim)
			})
			ddx0drs := num.SecondDerivMixedO4v1(u[0], u[1], 1e-3, func(r, s float64) float64 {
				U[0], U[1] = r, s
				surf.Point(tmp1, U, 2)
				return tmp1[0]
			})
			ddx1drs := num.SecondDerivMixedO4v1(u[0], u[1], 1e-3, func(r, s float64) float64 {
				U[0], U[1] = r, s
				surf.Point(tmp1, U, 2)
				return tmp1[1]
			})
			chk.Array(tst, "d²x/drds (FDM)", tol3, ddxdrs, []float64{ddx0drs, ddx1drs})
		}
	}
}

func checkNurbsSolidDerivs(tst *testing.T, solid *Nurbs, uvals, vvals, wvals []float64, verb bool, tol0, tol1, tol2, tol3 float64) {
	ndim := 3
	u := la.NewVector(3)
	U := la.NewVector(3)
	x := la.NewVector(ndim)
	tmp1 := la.NewVector(ndim)
	tmp2 := la.NewVector(ndim)
	tmp3 := la.NewVector(ndim)
	dxdr, dxds, dxdt := la.NewVector(ndim), la.NewVector(ndim), la.NewVector(ndim)
	ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst := la.NewVector(ndim), la.NewVector(ndim), la.NewVector(ndim), la.NewVector(ndim), la.NewVector(ndim), la.NewVector(ndim)
	for i := 0; i < len(uvals); i++ {
		for j := 0; j < len(vvals); j++ {
			for k := 0; k < len(wvals); k++ {
				u[0], u[1], u[2] = uvals[i], vvals[j], wvals[k]
				if verb {
					io.Pf("\nu = %f\n", u)
				}
				solid.PointAndDerivs(x, dxdr, dxds, dxdt, ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst, u, ndim)
				chk.DerivVecSca(tst, "dx/dr    ", tol1, dxdr, u[0], 1e-6, verb, func(xx []float64, r float64) {
					U[0], U[1], U[2] = r, u[1], u[2]
					solid.Point(xx, U, ndim)
				})
				chk.DerivVecSca(tst, "dx/ds    ", tol1, dxds, u[1], 1e-6, verb, func(xx []float64, s float64) {
					U[0], U[1], U[2] = u[0], s, u[2]
					solid.Point(xx, U, ndim)
				})
				chk.DerivVecSca(tst, "dx/dt    ", tol1, dxdt, u[2], 1e-6, verb, func(xx []float64, t float64) {
					U[0], U[1], U[2] = u[0], u[1], t
					solid.Point(xx, U, ndim)
				})
				chk.DerivVecSca(tst, "d²x/dr²  ", tol2, ddxdrr, u[0], 1e-6, verb, func(xx []float64, r float64) {
					U[0], U[1], U[2] = r, u[1], u[2]
					solid.PointAndDerivs(tmp1, xx, tmp2, tmp3, nil, nil, nil, nil, nil, nil, U, ndim)
				})
				chk.DerivVecSca(tst, "d²x/ds²  ", tol2, ddxdss, u[1], 1e-6, verb, func(xx []float64, s float64) {
					U[0], U[1], U[2] = u[0], s, u[2]
					solid.PointAndDerivs(tmp1, tmp2, xx, tmp3, nil, nil, nil, nil, nil, nil, U, ndim)
				})
				chk.DerivVecSca(tst, "d²x/dt²  ", tol2, ddxdtt, u[2], 1e-6, verb, func(xx []float64, t float64) {
					U[0], U[1], U[2] = u[0], u[1], t
					solid.PointAndDerivs(tmp1, tmp2, tmp3, xx, nil, nil, nil, nil, nil, nil, U, ndim)
				})
				chk.DerivVecSca(tst, "d²x/drds ", tol2, ddxdrs, u[1], 1e-6, verb, func(xx []float64, s float64) {
					U[0], U[1], U[2] = u[0], s, u[2]
					solid.PointAndDerivs(tmp1, xx, tmp2, tmp3, nil, nil, nil, nil, nil, nil, U, ndim)
				})
				chk.DerivVecSca(tst, "d²x/drdt ", tol2, ddxdrt, u[2], 1e-6, verb, func(xx []float64, t float64) {
					U[0], U[1], U[2] = u[0], u[1], t
					solid.PointAndDerivs(tmp1, xx, tmp2, tmp3, nil, nil, nil, nil, nil, nil, U, ndim)
				})
				chk.DerivVecSca(tst, "d²x/dsdt ", tol2, ddxdst, u[2], 1e-6, verb, func(xx []float64, t float64) {
					U[0], U[1], U[2] = u[0], u[1], t
					solid.PointAndDerivs(tmp1, tmp2, xx, tmp3, nil, nil, nil, nil, nil, nil, U, ndim)
				})
				ddx0drs := num.SecondDerivMixedO4v1(u[0], u[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = r, s, u[2]
					solid.Point(tmp1, U, ndim)
					return tmp1[0]
				})
				ddx1drs := num.SecondDerivMixedO4v1(u[0], u[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = r, s, u[2]
					solid.Point(tmp1, U, ndim)
					return tmp1[1]
				})
				ddx2drs := num.SecondDerivMixedO4v1(u[0], u[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = r, s, u[2]
					solid.Point(tmp1, U, ndim)
					return tmp1[2]
				})
				chk.Array(tst, "d²x/drds (FDM)", tol3, ddxdrs, []float64{ddx0drs, ddx1drs, ddx2drs})
				ddx0drt := num.SecondDerivMixedO4v1(u[0], u[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = r, u[1], t
					solid.Point(tmp1, U, ndim)
					return tmp1[0]
				})
				ddx1drt := num.SecondDerivMixedO4v1(u[0], u[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = r, u[1], t
					solid.Point(tmp1, U, ndim)
					return tmp1[1]
				})
				ddx2drt := num.SecondDerivMixedO4v1(u[0], u[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = r, u[1], t
					solid.Point(tmp1, U, ndim)
					return tmp1[2]
				})
				chk.Array(tst, "d²x/drdt (FDM)", tol3, ddxdrt, []float64{ddx0drt, ddx1drt, ddx2drt})
				ddx0dst := num.SecondDerivMixedO4v1(u[0], u[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = u[0], s, t
					solid.Point(tmp1, U, ndim)
					return tmp1[0]
				})
				ddx1dst := num.SecondDerivMixedO4v1(u[0], u[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = u[0], s, t
					solid.Point(tmp1, U, ndim)
					return tmp1[1]
				})
				ddx2dst := num.SecondDerivMixedO4v1(u[0], u[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = u[0], s, t
					solid.Point(tmp1, U, ndim)
					return tmp1[2]
				})
				chk.Array(tst, "d²x/dsdt (FDM)", tol3, ddxdst, []float64{ddx0dst, ddx1dst, ddx2dst})
			}
		}
	}
}

// check grid derivatives //////////////////////////////////////////////////////////////////////////

func checkGridNurbsDerivs2d(tst *testing.T, nrb *Nurbs, g *Grid, tol1, tol2, tol3 float64, verb bool) {
	x := la.NewVector(2)
	U := la.NewVector(2)
	Γ00 := la.NewVector(2)
	Γ11 := la.NewVector(2)
	Γ01 := la.NewVector(2)
	p := 0
	rs2uv := func(r, s float64) (u, v float64) {
		return nrb.UfromR(0, r), nrb.UfromR(1, s)
	}
	for n := 0; n < g.npts[1]; n++ {
		for m := 0; m < g.npts[0]; m++ {
			mtr := g.mtr[p][n][m]
			if verb {
				io.Pf("\nx = %v\n", mtr.X)
			}
			chk.DerivVecSca(tst, "g0 ", tol1, mtr.CovG0, mtr.U[0], 1e-3, verb, func(xx []float64, r float64) {
				U[0], U[1] = rs2uv(r, mtr.U[1])
				nrb.Point(xx, U, 2)
			})
			chk.DerivVecSca(tst, "g1 ", tol1, mtr.CovG1, mtr.U[1], 1e-3, verb, func(xx []float64, s float64) {
				U[0], U[1] = rs2uv(mtr.U[0], s)
				nrb.Point(xx, U, 2)
			})
			ddx0drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
				U[0], U[1] = rs2uv(r, mtr.U[1])
				nrb.Point(x, U, 2)
				return x[0]
			})
			ddx1drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
				U[0], U[1] = rs2uv(r, mtr.U[1])
				nrb.Point(x, U, 2)
				return x[1]
			})
			ddx0dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
				U[0], U[1] = rs2uv(mtr.U[0], s)
				nrb.Point(x, U, 2)
				return x[0]
			})
			ddx1dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
				U[0], U[1] = rs2uv(mtr.U[0], s)
				nrb.Point(x, U, 2)
				return x[1]
			})
			ddx0drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
				U[0], U[1] = rs2uv(r, s)
				nrb.Point(x, U, 2)
				return x[0]
			})
			ddx1drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
				U[0], U[1] = rs2uv(r, s)
				nrb.Point(x, U, 2)
				return x[1]
			})
			Γ00[0], Γ00[1] = ddx0drr, ddx1drr
			Γ11[0], Γ11[1] = ddx0dss, ddx1dss
			Γ01[0], Γ01[1] = ddx0drs, ddx1drs
			cntG0, cntG1 := mtr.GetContraVectors2d()
			Γ000 := la.VecDot(Γ00, cntG0)
			Γ011 := la.VecDot(Γ11, cntG0)
			Γ001 := la.VecDot(Γ01, cntG0)
			Γ100 := la.VecDot(Γ00, cntG1)
			Γ111 := la.VecDot(Γ11, cntG1)
			Γ101 := la.VecDot(Γ01, cntG1)
			chk.Deep3(tst, "GammaS", tol3, mtr.GammaS, [][][]float64{
				{
					{Γ000, Γ001},
					{Γ001, Γ011},
				},
				{
					{Γ100, Γ101},
					{Γ101, Γ111},
				},
			})
			if tst.Failed() {
				return
			}
		}
	}
}

func checkGridTfiniteDerivs2d(tst *testing.T, trf *Transfinite, g *Grid, tol1, tol2, tol3 float64, verb bool) {
	x := la.NewVector(2)
	U := la.NewVector(2)
	Γ00 := la.NewVector(2)
	Γ11 := la.NewVector(2)
	Γ01 := la.NewVector(2)
	p := 0
	for n := 0; n < g.npts[1]; n++ {
		for m := 0; m < g.npts[0]; m++ {
			mtr := g.mtr[p][n][m]
			if verb {
				io.Pf("\nx = %v\n", mtr.X)
			}
			chk.DerivVecSca(tst, "g0 ", tol1, mtr.CovG0, mtr.U[0], 1e-3, verb, func(xx []float64, r float64) {
				U[0], U[1] = r, mtr.U[1]
				trf.Point(xx, U)
			})
			chk.DerivVecSca(tst, "g1 ", tol1, mtr.CovG1, mtr.U[1], 1e-3, verb, func(xx []float64, s float64) {
				U[0], U[1] = mtr.U[0], s
				trf.Point(xx, U)
			})
			ddx0drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
				U[0], U[1] = r, mtr.U[1]
				trf.Point(x, U)
				return x[0]
			})
			ddx1drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
				U[0], U[1] = r, mtr.U[1]
				trf.Point(x, U)
				return x[1]
			})
			ddx0dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
				U[0], U[1] = mtr.U[0], s
				trf.Point(x, U)
				return x[0]
			})
			ddx1dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
				U[0], U[1] = mtr.U[0], s
				trf.Point(x, U)
				return x[1]
			})
			ddx0drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
				U[0], U[1] = r, s
				trf.Point(x, U)
				return x[0]
			})
			ddx1drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
				U[0], U[1] = r, s
				trf.Point(x, U)
				return x[1]
			})
			Γ00[0], Γ00[1] = ddx0drr, ddx1drr
			Γ11[0], Γ11[1] = ddx0dss, ddx1dss
			Γ01[0], Γ01[1] = ddx0drs, ddx1drs
			cntG0, cntG1 := mtr.GetContraVectors2d()
			Γ000 := la.VecDot(Γ00, cntG0)
			Γ011 := la.VecDot(Γ11, cntG0)
			Γ001 := la.VecDot(Γ01, cntG0)
			Γ100 := la.VecDot(Γ00, cntG1)
			Γ111 := la.VecDot(Γ11, cntG1)
			Γ101 := la.VecDot(Γ01, cntG1)
			chk.Deep3(tst, "GammaS", tol3, mtr.GammaS, [][][]float64{
				{
					{Γ000, Γ001},
					{Γ001, Γ011},
				},
				{
					{Γ100, Γ101},
					{Γ101, Γ111},
				},
			})
			if tst.Failed() {
				return
			}
		}
	}
}

func checkGridNurbsDerivs3d(tst *testing.T, nrb *Nurbs, g *Grid, tol1, tol2, tol3 float64, verb bool) {
	x := la.NewVector(3)
	U := la.NewVector(3)
	Γ00 := la.NewVector(3)
	Γ11 := la.NewVector(3)
	Γ22 := la.NewVector(3)
	Γ01 := la.NewVector(3)
	Γ02 := la.NewVector(3)
	Γ12 := la.NewVector(3)
	rst2uvw := func(r, s, t float64) (u, v, w float64) {
		return nrb.UfromR(0, r), nrb.UfromR(1, s), nrb.UfromR(2, t)
	}
	for p := 0; p < g.npts[2]; p++ {
		for n := 0; n < g.npts[1]; n++ {
			for m := 0; m < g.npts[0]; m++ {
				mtr := g.mtr[p][n][m]
				if verb {
					io.Pf("\nx = %v\n", mtr.X)
				}
				chk.DerivVecSca(tst, "g0 ", tol1, mtr.CovG0, mtr.U[0], 1e-3, verb, func(xx []float64, r float64) {
					U[0], U[1], U[2] = rst2uvw(r, mtr.U[1], mtr.U[2])
					nrb.Point(xx, U, 3)
				})
				chk.DerivVecSca(tst, "g1 ", tol1, mtr.CovG1, mtr.U[1], 1e-3, verb, func(xx []float64, s float64) {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], s, mtr.U[2])
					nrb.Point(xx, U, 3)
				})
				chk.DerivVecSca(tst, "g2 ", tol1, mtr.CovG2, mtr.U[2], 1e-3, verb, func(xx []float64, t float64) {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], mtr.U[1], t)
					nrb.Point(xx, U, 3)
				})

				ddx0drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, mtr.U[1], mtr.U[2])
					nrb.Point(x, U, 3)
					return x[0]
				})
				ddx1drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, mtr.U[1], mtr.U[2])
					nrb.Point(x, U, 3)
					return x[1]
				})
				ddx2drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, mtr.U[1], mtr.U[2])
					nrb.Point(x, U, 3)
					return x[2]
				})

				ddx0dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], s, mtr.U[2])
					nrb.Point(x, U, 3)
					return x[0]
				})
				ddx1dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], s, mtr.U[2])
					nrb.Point(x, U, 3)
					return x[1]
				})
				ddx2dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], s, mtr.U[2])
					nrb.Point(x, U, 3)
					return x[2]
				})

				ddx0dtt := num.SecondDerivCen5(mtr.U[2], 1e-3, func(t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], mtr.U[1], t)
					nrb.Point(x, U, 3)
					return x[0]
				})
				ddx1dtt := num.SecondDerivCen5(mtr.U[2], 1e-3, func(t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], mtr.U[1], t)
					nrb.Point(x, U, 3)
					return x[1]
				})
				ddx2dtt := num.SecondDerivCen5(mtr.U[2], 1e-3, func(t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], mtr.U[1], t)
					nrb.Point(x, U, 3)
					return x[2]
				})

				ddx0drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, s, mtr.U[2])
					nrb.Point(x, U, 3)
					return x[0]
				})
				ddx1drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, s, mtr.U[2])
					nrb.Point(x, U, 3)
					return x[1]
				})
				ddx2drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, s, mtr.U[2])
					nrb.Point(x, U, 3)
					return x[2]
				})

				ddx0drt := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, mtr.U[1], t)
					nrb.Point(x, U, 3)
					return x[0]
				})
				ddx1drt := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, mtr.U[1], t)
					nrb.Point(x, U, 3)
					return x[1]
				})
				ddx2drt := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(r, mtr.U[1], t)
					nrb.Point(x, U, 3)
					return x[2]
				})

				ddx0dst := num.SecondDerivMixedO4v1(mtr.U[1], mtr.U[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], s, t)
					nrb.Point(x, U, 3)
					return x[0]
				})
				ddx1dst := num.SecondDerivMixedO4v1(mtr.U[1], mtr.U[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], s, t)
					nrb.Point(x, U, 3)
					return x[1]
				})
				ddx2dst := num.SecondDerivMixedO4v1(mtr.U[1], mtr.U[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = rst2uvw(mtr.U[0], s, t)
					nrb.Point(x, U, 3)
					return x[2]
				})

				Γ00[0], Γ00[1], Γ00[2] = ddx0drr, ddx1drr, ddx2drr
				Γ11[0], Γ11[1], Γ11[2] = ddx0dss, ddx1dss, ddx2dss
				Γ22[0], Γ22[1], Γ22[2] = ddx0dtt, ddx1dtt, ddx2dtt
				Γ01[0], Γ01[1], Γ01[2] = ddx0drs, ddx1drs, ddx2drs
				Γ02[0], Γ02[1], Γ02[2] = ddx0drt, ddx1drt, ddx2drt
				Γ12[0], Γ12[1], Γ12[2] = ddx0dst, ddx1dst, ddx2dst
				cntG0, cntG1, cntG2 := mtr.GetContraVectors3d()

				Γ000 := la.VecDot(Γ00, cntG0)
				Γ011 := la.VecDot(Γ11, cntG0)
				Γ022 := la.VecDot(Γ22, cntG0)
				Γ001 := la.VecDot(Γ01, cntG0)
				Γ002 := la.VecDot(Γ02, cntG0)
				Γ012 := la.VecDot(Γ12, cntG0)

				Γ100 := la.VecDot(Γ00, cntG1)
				Γ111 := la.VecDot(Γ11, cntG1)
				Γ122 := la.VecDot(Γ22, cntG1)
				Γ101 := la.VecDot(Γ01, cntG1)
				Γ102 := la.VecDot(Γ02, cntG1)
				Γ112 := la.VecDot(Γ12, cntG1)

				Γ200 := la.VecDot(Γ00, cntG2)
				Γ211 := la.VecDot(Γ11, cntG2)
				Γ222 := la.VecDot(Γ22, cntG2)
				Γ201 := la.VecDot(Γ01, cntG2)
				Γ202 := la.VecDot(Γ02, cntG2)
				Γ212 := la.VecDot(Γ12, cntG2)

				chk.Deep3(tst, "GammaS", tol3, mtr.GammaS, [][][]float64{
					{
						{Γ000, Γ001, Γ002},
						{Γ001, Γ011, Γ012},
						{Γ002, Γ012, Γ022},
					},
					{
						{Γ100, Γ101, Γ102},
						{Γ101, Γ111, Γ112},
						{Γ102, Γ112, Γ122},
					},
					{
						{Γ200, Γ201, Γ202},
						{Γ201, Γ211, Γ212},
						{Γ202, Γ212, Γ222},
					},
				})
				if tst.Failed() {
					return
				}
			}
		}
	}
}

func checkGridTfiniteDerivs3d(tst *testing.T, trf *Transfinite, g *Grid, tol1, tol2, tol3 float64, verb bool) {
	x := la.NewVector(3)
	U := la.NewVector(3)
	Γ00 := la.NewVector(3)
	Γ11 := la.NewVector(3)
	Γ22 := la.NewVector(3)
	Γ01 := la.NewVector(3)
	Γ02 := la.NewVector(3)
	Γ12 := la.NewVector(3)
	for p := 0; p < g.npts[2]; p++ {
		for n := 0; n < g.npts[1]; n++ {
			for m := 0; m < g.npts[0]; m++ {
				mtr := g.mtr[p][n][m]
				if verb {
					io.Pf("\nx = %v\n", mtr.X)
				}
				chk.DerivVecSca(tst, "g0 ", tol1, mtr.CovG0, mtr.U[0], 1e-3, verb, func(xx []float64, r float64) {
					U[0], U[1], U[2] = r, mtr.U[1], mtr.U[2]
					trf.Point(xx, U)
				})
				chk.DerivVecSca(tst, "g1 ", tol1, mtr.CovG1, mtr.U[1], 1e-3, verb, func(xx []float64, s float64) {
					U[0], U[1], U[2] = mtr.U[0], s, mtr.U[2]
					trf.Point(xx, U)
				})
				chk.DerivVecSca(tst, "g2 ", tol1, mtr.CovG2, mtr.U[2], 1e-3, verb, func(xx []float64, t float64) {
					U[0], U[1], U[2] = mtr.U[0], mtr.U[1], t
					trf.Point(xx, U)
				})

				ddx0drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
					U[0], U[1], U[2] = r, mtr.U[1], mtr.U[2]
					trf.Point(x, U)
					return x[0]
				})
				ddx1drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
					U[0], U[1], U[2] = r, mtr.U[1], mtr.U[2]
					trf.Point(x, U)
					return x[1]
				})
				ddx2drr := num.SecondDerivCen5(mtr.U[0], 1e-3, func(r float64) float64 {
					U[0], U[1], U[2] = r, mtr.U[1], mtr.U[2]
					trf.Point(x, U)
					return x[2]
				})

				ddx0dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], s, mtr.U[2]
					trf.Point(x, U)
					return x[0]
				})
				ddx1dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], s, mtr.U[2]
					trf.Point(x, U)
					return x[1]
				})
				ddx2dss := num.SecondDerivCen5(mtr.U[1], 1e-3, func(s float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], s, mtr.U[2]
					trf.Point(x, U)
					return x[2]
				})

				ddx0dtt := num.SecondDerivCen5(mtr.U[2], 1e-3, func(t float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], mtr.U[1], t
					trf.Point(x, U)
					return x[0]
				})
				ddx1dtt := num.SecondDerivCen5(mtr.U[2], 1e-3, func(t float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], mtr.U[1], t
					trf.Point(x, U)
					return x[1]
				})
				ddx2dtt := num.SecondDerivCen5(mtr.U[2], 1e-3, func(t float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], mtr.U[1], t
					trf.Point(x, U)
					return x[2]
				})

				ddx0drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = r, s, mtr.U[2]
					trf.Point(x, U)
					return x[0]
				})
				ddx1drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = r, s, mtr.U[2]
					trf.Point(x, U)
					return x[1]
				})
				ddx2drs := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[1], 1e-3, func(r, s float64) float64 {
					U[0], U[1], U[2] = r, s, mtr.U[2]
					trf.Point(x, U)
					return x[2]
				})

				ddx0drt := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = r, mtr.U[1], t
					trf.Point(x, U)
					return x[0]
				})
				ddx1drt := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = r, mtr.U[1], t
					trf.Point(x, U)
					return x[1]
				})
				ddx2drt := num.SecondDerivMixedO4v1(mtr.U[0], mtr.U[2], 1e-3, func(r, t float64) float64 {
					U[0], U[1], U[2] = r, mtr.U[1], t
					trf.Point(x, U)
					return x[2]
				})

				ddx0dst := num.SecondDerivMixedO4v1(mtr.U[1], mtr.U[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], s, t
					trf.Point(x, U)
					return x[0]
				})
				ddx1dst := num.SecondDerivMixedO4v1(mtr.U[1], mtr.U[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], s, t
					trf.Point(x, U)
					return x[1]
				})
				ddx2dst := num.SecondDerivMixedO4v1(mtr.U[1], mtr.U[2], 1e-3, func(s, t float64) float64 {
					U[0], U[1], U[2] = mtr.U[0], s, t
					trf.Point(x, U)
					return x[2]
				})

				Γ00[0], Γ00[1], Γ00[2] = ddx0drr, ddx1drr, ddx2drr
				Γ11[0], Γ11[1], Γ11[2] = ddx0dss, ddx1dss, ddx2dss
				Γ22[0], Γ22[1], Γ22[2] = ddx0dtt, ddx1dtt, ddx2dtt
				Γ01[0], Γ01[1], Γ01[2] = ddx0drs, ddx1drs, ddx2drs
				Γ02[0], Γ02[1], Γ02[2] = ddx0drt, ddx1drt, ddx2drt
				Γ12[0], Γ12[1], Γ12[2] = ddx0dst, ddx1dst, ddx2dst
				cntG0, cntG1, cntG2 := mtr.GetContraVectors3d()

				Γ000 := la.VecDot(Γ00, cntG0)
				Γ011 := la.VecDot(Γ11, cntG0)
				Γ022 := la.VecDot(Γ22, cntG0)
				Γ001 := la.VecDot(Γ01, cntG0)
				Γ002 := la.VecDot(Γ02, cntG0)
				Γ012 := la.VecDot(Γ12, cntG0)

				Γ100 := la.VecDot(Γ00, cntG1)
				Γ111 := la.VecDot(Γ11, cntG1)
				Γ122 := la.VecDot(Γ22, cntG1)
				Γ101 := la.VecDot(Γ01, cntG1)
				Γ102 := la.VecDot(Γ02, cntG1)
				Γ112 := la.VecDot(Γ12, cntG1)

				Γ200 := la.VecDot(Γ00, cntG2)
				Γ211 := la.VecDot(Γ11, cntG2)
				Γ222 := la.VecDot(Γ22, cntG2)
				Γ201 := la.VecDot(Γ01, cntG2)
				Γ202 := la.VecDot(Γ02, cntG2)
				Γ212 := la.VecDot(Γ12, cntG2)

				chk.Deep3(tst, "GammaS", tol3, mtr.GammaS, [][][]float64{
					{
						{Γ000, Γ001, Γ002},
						{Γ001, Γ011, Γ012},
						{Γ002, Γ012, Γ022},
					},
					{
						{Γ100, Γ101, Γ102},
						{Γ101, Γ111, Γ112},
						{Γ102, Γ112, Γ122},
					},
					{
						{Γ200, Γ201, Γ202},
						{Γ201, Γ211, Γ212},
						{Γ202, Γ212, Γ222},
					},
				})
				if tst.Failed() {
					return
				}
			}
		}
	}
}
