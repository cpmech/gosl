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

func TestGrid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid01. rectangular uniform 2D")

	g := new(Grid)
	g.RectGenUniform([]float64{-6, -3}, []float64{6, 3}, []int{5, 4})

	chk.Int(tst, "ndim", g.Ndim(), 2)
	chk.Int(tst, "size", g.Size(), 20)
	chk.Int(tst, "nx", g.Npts(0), 5)
	chk.Int(tst, "ny", g.Npts(1), 4)

	chk.Float64(tst, "Lx", 1e-15, g.Xlength(0), 12.0)
	chk.Float64(tst, "Ly", 1e-15, g.Xlength(1), 6.0)

	min := []float64{g.Xmin(0), g.Xmin(1)}
	max := []float64{g.Xmax(0), g.Xmax(1)}
	del := []float64{g.Xlength(0), g.Xlength(1)}

	chk.Array(tst, "Min", 1e-17, min, []float64{-6, -3})
	chk.Array(tst, "Max", 1e-17, max, []float64{+6, +3})
	chk.Array(tst, "Del", 1e-17, del, []float64{12, 6})

	chk.Ints(tst, "B", g.Edge(0), []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "R", g.Edge(1), []int{4, 9, 14, 19})
	chk.Ints(tst, "T", g.Edge(2), []int{15, 16, 17, 18, 19})
	chk.Ints(tst, "L", g.Edge(3), []int{0, 5, 10, 15})

	chk.Ints(tst, "Tag # 10: L", g.EdgeT(10), []int{0, 5, 10, 15})
	chk.Ints(tst, "Tag # 11: R", g.EdgeT(11), []int{4, 9, 14, 19})
	chk.Ints(tst, "Tag # 20: B", g.EdgeT(20), []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "Tag # 21: T", g.EdgeT(21), []int{15, 16, 17, 18, 19})

	xx, yy := g.Meshgrid2d()
	chk.Deep2(tst, "xx", 1e-17, xx, [][]float64{
		{-6, -3, 0, 3, 6},
		{-6, -3, 0, 3, 6},
		{-6, -3, 0, 3, 6},
		{-6, -3, 0, 3, 6},
	})
	chk.Deep2(tst, "yy", 1e-17, yy, [][]float64{
		{-3, -3, -3, -3, -3},
		{-1, -1, -1, -1, -1},
		{+1, +1, +1, +1, +1},
		{+3, +3, +3, +3, +3},
	})

	chk.Array(tst, "x[0]", 1e-17, g.Node(0), []float64{-6, -3})
	chk.Array(tst, "x[7]", 1e-17, g.Node(7), []float64{0, -1})
	chk.Array(tst, "x[9]", 1e-17, g.Node(9), []float64{6, -1})
	chk.Array(tst, "x[15]", 1e-17, g.Node(15), []float64{-6, 3})
	chk.Array(tst, "x[19]", 1e-17, g.Node(19), []float64{6, 3})

	// plot
	if chk.Verbose {
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp.Draw()
		gp.Bases(1)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.SetXnticks(19)
		plt.SetYnticks(15)
		plt.Save("/tmp/gosl/gm", "grid01")
	}
}

func TestGrid02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid02. rectangular uniform 3D")

	//TODO
}

func TestGrid03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid03. rectangular uniform (RectSet2D)")

	g := new(Grid)
	g.RectSet2d([]float64{1, 2, 4, 8, 16}, []float64{0, 3, 4, 7})

	chk.Int(tst, "ndim", g.Ndim(), 2)
	chk.Int(tst, "size", g.Size(), 20)
	chk.Int(tst, "nx", g.Npts(0), 5)
	chk.Int(tst, "ny", g.Npts(1), 4)

	min := []float64{g.Xmin(0), g.Xmin(1)}
	max := []float64{g.Xmax(0), g.Xmax(1)}
	del := []float64{g.Xlength(0), g.Xlength(1)}

	chk.Array(tst, "Min", 1e-17, min, []float64{1, 0})
	chk.Array(tst, "Max", 1e-17, max, []float64{16, 7})
	chk.Array(tst, "Del", 1e-17, del, []float64{15, 7})

	chk.Ints(tst, "B", g.Edge(0), []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "R", g.Edge(1), []int{4, 9, 14, 19})
	chk.Ints(tst, "T", g.Edge(2), []int{15, 16, 17, 18, 19})
	chk.Ints(tst, "L", g.Edge(3), []int{0, 5, 10, 15})

	chk.Ints(tst, "Tag # 10: L", g.EdgeT(10), []int{0, 5, 10, 15})
	chk.Ints(tst, "Tag # 11: R", g.EdgeT(11), []int{4, 9, 14, 19})
	chk.Ints(tst, "Tag # 20: B", g.EdgeT(20), []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "Tag # 21: T", g.EdgeT(21), []int{15, 16, 17, 18, 19})

	xx, yy := g.Meshgrid2d()
	chk.Deep2(tst, "xx", 1e-17, xx, [][]float64{
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
	})
	chk.Deep2(tst, "yy", 1e-17, yy, [][]float64{
		{0, 0, 0, 0, 0},
		{3, 3, 3, 3, 3},
		{4, 4, 4, 4, 4},
		{7, 7, 7, 7, 7},
	})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp := GridPlotter{G: g, WithVids: true}
		gp.Draw()
		gp.Bases(1)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.SetXnticks(17)
		plt.SetYnticks(17)
		plt.Save("/tmp/gosl/gm", "grid03")
	}
}

func TestGrid04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid04. rectangular uniform (RectSet3D)")

	g := new(Grid)
	g.RectSet3d([]float64{1, 2, 4, 8}, []float64{0, 3, 4}, []float64{-1, -0.5})

	chk.Int(tst, "ndim", g.Ndim(), 3)
	chk.Int(tst, "size", g.Size(), 24)
	chk.Int(tst, "nx", g.Npts(0), 4)
	chk.Int(tst, "ny", g.Npts(1), 3)
	chk.Int(tst, "nz", g.Npts(2), 2)

	chk.Ints(tst, "xmin", g.Face(0), []int{0, 4, 8, 12, 16, 20})
	chk.Ints(tst, "xmax", g.Face(1), []int{3, 7, 11, 15, 19, 23})
	chk.Ints(tst, "ymin", g.Face(2), []int{0, 1, 2, 3, 12, 13, 14, 15})
	chk.Ints(tst, "ymax", g.Face(3), []int{8, 9, 10, 11, 20, 21, 22, 23})
	chk.Ints(tst, "zmin", g.Face(4), []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
	chk.Ints(tst, "zmax", g.Face(5), []int{12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23})

	chk.Ints(tst, "Tag # 10: xmin", g.Boundary(100), g.Face(0))
	chk.Ints(tst, "Tag # 11: xmax", g.Boundary(101), g.Face(1))
	chk.Ints(tst, "Tag # 20: ymin", g.Boundary(200), g.Face(2))
	chk.Ints(tst, "Tag # 21: ymax", g.Boundary(201), g.Face(3))
	chk.Ints(tst, "Tag # 30: zmin", g.Boundary(300), g.Face(4))
	chk.Ints(tst, "Tag # 31: zmax", g.Boundary(301), g.Face(5))

	xx, yy, zz := g.Meshgrid3d()

	chk.Deep3(tst, "xx", 1e-17, xx, [][][]float64{
		{
			{1, 2, 4, 8},
			{1, 2, 4, 8},
			{1, 2, 4, 8},
		},
		{
			{1, 2, 4, 8},
			{1, 2, 4, 8},
			{1, 2, 4, 8},
		},
	})
	chk.Deep3(tst, "yy", 1e-17, yy, [][][]float64{
		{
			{0, 0, 0, 0},
			{3, 3, 3, 3},
			{4, 4, 4, 4},
		},
		{
			{0, 0, 0, 0},
			{3, 3, 3, 3},
			{4, 4, 4, 4},
		},
	})
	chk.Deep3(tst, "zz", 1e-17, zz, [][][]float64{
		{
			{-1, -1, -1, -1},
			{-1, -1, -1, -1},
			{-1, -1, -1, -1},
		},
		{
			{-0.5, -0.5, -0.5, -0.5},
			{-0.5, -0.5, -0.5, -0.5},
			{-0.5, -0.5, -0.5, -0.5},
		},
	})

	min := []float64{g.Xmin(0), g.Xmin(1), g.Xmin(2)}
	max := []float64{g.Xmax(0), g.Xmax(1), g.Xmax(2)}
	del := []float64{g.Xlength(0), g.Xlength(1), g.Xlength(2)}

	chk.Array(tst, "Min", 1e-17, min, []float64{1, 0, -1})
	chk.Array(tst, "Max", 1e-17, max, []float64{8, 4, -0.5})
	chk.Array(tst, "Del", 1e-17, del, []float64{7, 4, 0.5})

	chk.Array(tst, "x[0]", 1e-17, g.Node(0), []float64{1, 0, -1})
	chk.Array(tst, "x[1]", 1e-17, g.Node(1), []float64{2, 0, -1})
	chk.Array(tst, "x[6]", 1e-17, g.Node(6), []float64{4, 3, -1})
	chk.Array(tst, "x[8]", 1e-17, g.Node(8), []float64{1, 4, -1})
	chk.Array(tst, "x[11]", 1e-17, g.Node(11), []float64{8, 4, -1})
	chk.Array(tst, "x[12]", 1e-17, g.Node(12), []float64{1, 0, -0.5})
	chk.Array(tst, "x[17]", 1e-17, g.Node(17), []float64{2, 3, -0.5})
	chk.Array(tst, "x[19]", 1e-17, g.Node(19), []float64{8, 3, -0.5})
	chk.Array(tst, "x[22]", 1e-17, g.Node(22), []float64{4, 4, -0.5})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp := GridPlotter{G: g, WithVids: true}
		gp.Draw()
		gp.Bases(0.5)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.DefaultTriad(1)
		plt.Default3dView(g.Xmin(0), g.Xmax(0), g.Xmin(1), g.Xmax(1), g.Xmin(2), g.Xmax(2), true)
		plt.Save("/tmp/gosl/gm", "grid04")
	}
}

func TestGrid05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid05. 2d ring")

	// mapping
	a, b := 1.0, 2.0
	trf := FactoryTfinite.Surf2dQuarterRing(a, b)

	// coordinates
	R := utl.LinSpace(-1, 1, 5)
	S := utl.LinSpace(-1, 1, 5)

	// curvgrid
	cg := new(Grid)
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
		plt.Save("/tmp/gosl/gm", "grid05")
	}
}

func TestGrid06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid06. 3d ring")

	// mapping
	a, b, h := 2.0, 3.0, 2.0 // radii and thickness
	trf := FactoryTfinite.Surf3dQuarterRing(a, b, h)

	// coordinates
	npts := 3
	R := utl.LinSpace(-1, 1, npts)
	S := utl.LinSpace(-1, 1, npts)
	T := utl.LinSpace(-1, 1, npts)

	// curvgrid
	cg := new(Grid)
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
		plt.Save("/tmp/gosl/gm", "grid06")
	}
}
