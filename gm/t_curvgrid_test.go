// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestCurvGrid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("CurvGrid01. 2d ring")

	// mapping
	a, b := 1.0, 2.0
	trf := FactoryTfinite.Surf2dQuarterRing(a, b)

	// coordinates
	R := utl.LinSpace(-1, 1, 5)
	S := utl.LinSpace(-1, 1, 5)

	// curvgrid
	cg := new(CurvGrid)
	cg.SetTransfinite2d(trf, R, S)

	// check limits
	chk.Array(tst, "umin", 1e-15, cg.umin, []float64{-1, -1, -1})
	chk.Array(tst, "umax", 1e-15, cg.umax, []float64{+1, +1, -1})
	chk.Array(tst, "xmin", 1e-15, cg.xmin, []float64{0, 0, 0})
	chk.Array(tst, "xmax", 1e-15, cg.xmax, []float64{b, b, 0})

	// check metrics
	π := math.Pi
	A := (b - a) / 2.0 // dρ/dr
	B := π / 4.0       // dα/ds
	p := 0             // z-index
	for n := 0; n < cg.npts[1]; n++ {
		for m := 0; m < cg.npts[0]; m++ {
			mtr := cg.mtr[p][n][m]
			ρ := a + (1.0+mtr.U[0])*A // cylindrical coordinates
			α := (1.0 + mtr.U[1]) * B // cylindrical coordinates
			cα, sα := math.Cos(α), math.Sin(α)
			chk.Array(tst, "x      ", 1e-14, mtr.X, []float64{ρ * cα, ρ * sα})
			chk.Array(tst, "CovG0  ", 1e-14, mtr.CovG0, []float64{cα * A, sα * A})
			chk.Array(tst, "CovG1  ", 1e-14, mtr.CovG1, []float64{-ρ * sα * B, ρ * cα * B})
			chk.Deep2(tst, "CovGmat", 1e-14, mtr.CovGmat.GetDeep2(), [][]float64{
				{A * A, 0.0},
				{0.0, ρ * ρ * B * B},
			})
			chk.Deep2(tst, "CntGmat", 1e-14, mtr.CntGmat.GetDeep2(), [][]float64{
				{1.0 / (A * A), 0.0},
				{0.0, 1.0 / (ρ * ρ * B * B)},
			})
			chk.Deep3(tst, "GammaS", 1e-14, mtr.GammaS, [][][]float64{
				{
					{0, 0},
					{0, -ρ * B * B / A},
				},
				{
					{0, A / ρ},
					{A / ρ, 0},
				},
			})
			chk.Array(tst, "L", 1e-14, mtr.L, []float64{-1.0 / (ρ * A), 0})
		}
	}

	// check interface functions
	io.Pl()
	chk.Int(tst, "Ndim()", cg.Ndim(), 2)
	chk.Int(tst, "Npts(0)", cg.Npts(0), len(R))
	chk.Int(tst, "Size()", cg.Size(), len(R)*len(S))
	chk.Float64(tst, "Umin(0)", 1e-14, cg.Umin(0), -1)
	chk.Float64(tst, "Umax(0)", 1e-14, cg.Umax(0), +1)
	chk.Float64(tst, "Xmin(0)", 1e-14, cg.Xmin(0), 0)
	chk.Float64(tst, "Xmax(0)", 1e-14, cg.Xmax(0), b)
	chk.Array(tst, "U(0,0,0)", 1e-14, cg.U(0, 0, 0), []float64{-1, -1})
	chk.Array(tst, "X(0,0,0)", 1e-14, cg.X(0, 0, 0), []float64{a, 0})
	chk.Array(tst, "g0(0,0,0)", 1e-14, cg.CovarBasis(0, 0, 0, 0), []float64{A, 0})
	chk.Array(tst, "g1(0,0,0)", 1e-14, cg.CovarBasis(0, 0, 0, 1), []float64{0, a * B})
	chk.Array(tst, "g2(0,0,0)", 1e-14, cg.CovarBasis(0, 0, 0, 2), nil)
	chk.Deep2(tst, "g_ij(0,0,0)", 1e-14, cg.CovarMatrix(0, 0, 0).GetDeep2(), [][]float64{
		{A * A, 0},
		{0, a * a * B * B},
	})
	chk.Deep2(tst, "g^ij(0,0,0)", 1e-14, cg.ContraMatrix(0, 0, 0).GetDeep2(), [][]float64{
		{1.0 / (A * A), 0},
		{0, 1.0 / (a * a * B * B)},
	})
	chk.Float64(tst, "det(g)(0,0,0)", 1e-14, cg.DetCovarMatrix(0, 0, 0), A*A*a*a*B*B)
	chk.Float64(tst, "Γ(0,0,0; 0,1,1)", 1e-14, cg.GammaS(0, 0, 0, 0, 1, 1), -a*B*B/A)
	chk.Float64(tst, "L(0,0,0; 0)", 1e-14, cg.Lcoeff(0, 0, 0, 0), -1.0/(a*A))

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{11, 21}, false, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		cg.DrawBases(0.15, nil, nil, nil)
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "curvgrid01")
	}
}

func TestCurvGrid02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("CurvGrid02. 3d ring")

	// mapping
	a, b, h := 2.0, 3.0, 2.0 // radii and thickness
	trf := FactoryTfinite.Surf3dQuarterRing(a, b, h)

	// coordinates
	npts := 3
	R := utl.LinSpace(-1, 1, npts)
	S := utl.LinSpace(-1, 1, npts)
	T := utl.LinSpace(-1, 1, npts)

	// curvgrid
	cg := new(CurvGrid)
	cg.SetTransfinite3d(trf, R, S, T)

	// check limits
	chk.Array(tst, "umin", 1e-15, cg.umin, []float64{-1, -1, -1})
	chk.Array(tst, "umax", 1e-15, cg.umax, []float64{+1, +1, +1})
	chk.Array(tst, "xmin", 1e-15, cg.xmin, []float64{0, 0, 0})
	chk.Array(tst, "xmax", 1e-15, cg.xmax, []float64{h, b, b})

	// check
	π := math.Pi
	A := (b - a) / 2.0 // dρ/dr
	B := π / 4.0       // dα/ds
	for p := 0; p < cg.npts[2]; p++ {
		for n := 0; n < cg.npts[1]; n++ {
			for m := 0; m < cg.npts[0]; m++ {
				mtr := cg.mtr[p][n][m]
				x0 := h * float64(m) / float64(cg.npts[0]-1)
				ρ := a + (1.0+mtr.U[1])*A // cylindrical coordinates
				α := (1.0 + mtr.U[2]) * B // cylindrical coordinates
				cα, sα := math.Cos(α), math.Sin(α)
				chk.Array(tst, "x      ", 1e-14, mtr.X, []float64{x0, ρ * cα, ρ * sα})
				chk.Array(tst, "CovG0  ", 1e-14, mtr.CovG0, []float64{1, 0, 0})
				chk.Array(tst, "CovG1  ", 1e-14, mtr.CovG1, []float64{0, cα * A, sα * A})
				chk.Array(tst, "CovG2  ", 1e-14, mtr.CovG2, []float64{0, -ρ * sα * B, ρ * cα * B})
				chk.Deep2(tst, "CovGmat", 1e-14, mtr.CovGmat.GetDeep2(), [][]float64{
					{1.0, 0.0, 0.0},
					{0.0, A * A, 0.0},
					{0.0, 0.0, ρ * ρ * B * B},
				})
				chk.Deep2(tst, "CntGmat", 1e-14, mtr.CntGmat.GetDeep2(), [][]float64{
					{1.0, 0.0, 0.0},
					{0.0, 1.0 / (A * A), 0.0},
					{0.0, 0.0, 1.0 / (ρ * ρ * B * B)},
				})
				chk.Deep3(tst, "GammaS", 1e-14, mtr.GammaS, [][][]float64{
					{
						{0, 0, 0},
						{0, 0, 0},
						{0, 0, 0},
					},
					{
						{0, 0, 0},
						{0, 0, 0},
						{0, 0, -ρ * B * B / A},
					},
					{
						{0, 0, 0},
						{0, 0, A / ρ},
						{0, A / ρ, 0},
					},
				})
				chk.Array(tst, "L", 1e-14, mtr.L, []float64{0, -1.0 / (ρ * A), 0})
			}
		}
	}

	// check interface functions
	io.Pl()
	chk.Int(tst, "Ndim()", cg.Ndim(), 3)
	chk.Int(tst, "Npts(0)", cg.Npts(0), len(R))
	chk.Int(tst, "Size()", cg.Size(), len(R)*len(S)*len(T))
	chk.Float64(tst, "Umin(2)", 1e-14, cg.Umin(2), -1)
	chk.Float64(tst, "Umax(2)", 1e-14, cg.Umax(2), +1)
	chk.Float64(tst, "Xmin(2)", 1e-14, cg.Xmin(2), 0)
	chk.Float64(tst, "Xmax(2)", 1e-14, cg.Xmax(2), b)
	chk.Array(tst, "U(0,0,0)", 1e-14, cg.U(0, 0, 0), []float64{-1, -1, -1})
	chk.Array(tst, "X(0,0,0)", 1e-14, cg.X(0, 0, 0), []float64{0, a, 0})
	chk.Array(tst, "g0(0,0,0)", 1e-14, cg.CovarBasis(0, 0, 0, 0), []float64{1, 0, 0})
	chk.Array(tst, "g1(0,0,0)", 1e-14, cg.CovarBasis(0, 0, 0, 1), []float64{0, A, 0})
	chk.Array(tst, "g2(0,0,0)", 1e-14, cg.CovarBasis(0, 0, 0, 2), []float64{0, 0, a * B})
	chk.Deep2(tst, "g_ij(0,0,0)", 1e-14, cg.CovarMatrix(0, 0, 0).GetDeep2(), [][]float64{
		{1.0, 0.0, 0.0},
		{0.0, A * A, 0.0},
		{0.0, 0.0, a * a * B * B},
	})
	chk.Deep2(tst, "g^ij(0,0,0)", 1e-14, cg.ContraMatrix(0, 0, 0).GetDeep2(), [][]float64{
		{1.0, 0.0, 0.0},
		{0.0, 1.0 / (A * A), 0.0},
		{0.0, 0.0, 1.0 / (a * a * B * B)},
	})
	chk.Float64(tst, "det(g)(0,0,0)", 1e-14, cg.DetCovarMatrix(0, 0, 0), A*A*a*a*B*B)
	chk.Float64(tst, "Γ(0,0,0; 1,2,2)", 1e-14, cg.GammaS(0, 0, 0, 1, 2, 2), -a*B*B/A)
	chk.Float64(tst, "L(0,0,0; 1)", 1e-14, cg.Lcoeff(0, 0, 0, 1), -1.0/(a*A))

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400})
		trf.Draw([]int{5, 5, 11}, true, &plt.A{C: "#7d8891"}, &plt.A{C: plt.C(3, 9), Lw: 2})
		cg.DrawBases(0.20, nil, nil, nil)
		plt.Default3dView(0, 3, 0, 3, 0, 3, true)
		plt.Save("/tmp/gosl/gm", "curvgrid02")
	}
}
