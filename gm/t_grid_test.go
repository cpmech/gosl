// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/la"
	"gosl/plt"
	"gosl/utl"
)

func TestGrid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid01. rectangular uniform 2D")

	// grid
	g := new(Grid)
	g.RectGenUniform([]float64{-6, -3}, []float64{6, 3}, []int{5, 4})

	// accessors
	io.Pf("\naccessors\n")
	xx, yy := g.Meshgrid2d()
	chk.Int(tst, "Ndim", g.Ndim(), 2)
	chk.Ints(tst, "Npts", []int{g.Npts(0), g.Npts(1)}, []int{5, 4})
	chk.Int(tst, "Size", g.Size(), 20)
	chk.Array(tst, "Umin", 1e-14, []float64{g.Umin(0), g.Umin(1)}, []float64{-1, -1})
	chk.Array(tst, "Umax", 1e-14, []float64{g.Umax(0), g.Umax(1)}, []float64{+1, +1})
	chk.Array(tst, "Xmin", 1e-17, []float64{g.Xmin(0), g.Xmin(1)}, []float64{-6, -3})
	chk.Array(tst, "Xmax", 1e-17, []float64{g.Xmax(0), g.Xmax(1)}, []float64{+6, +3})
	chk.Array(tst, "Xlen", 1e-17, []float64{g.Xlen(0), g.Xlen(1)}, []float64{12, 6})
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

	// metrics accessors
	io.Pf("\nmetrics accessors\n")
	p := 0
	det := 36.0 * 9.0 // determinant of CovarMatrix
	chk.Array(tst, "U(0,0,p)", 1e-14, g.U(0, 0, p), []float64{-1, -1})
	chk.Array(tst, "X(0,0,p)", 1e-14, g.X(0, 0, p), []float64{-6, -3})
	chk.Array(tst, "U(3,1,p)", 1e-14, g.U(3, 1, p), []float64{0.5, -1 + 2.0/3.0})
	chk.Array(tst, "X(3,1,p)", 1e-14, g.X(3, 1, p), []float64{3, -1})
	for _, n := range []int{0, 2} {
		for _, m := range []int{0, 3} {
			t := io.Sf("%d,%d,p", m, n)
			chk.Array(tst, "g_0("+t+")", 1e-14, g.CovarBasis(m, n, p, 0), []float64{6, 0}) // (xmax-xmin)/2
			chk.Array(tst, "g_1("+t+")", 1e-14, g.CovarBasis(m, n, p, 1), []float64{0, 3}) // (ymax-ymin)/2
			chk.Array(tst, "g_2("+t+")", 1e-14, g.CovarBasis(m, n, p, 2), nil)
			chk.Deep2(tst, "g_ij("+t+")", 1e-14, g.CovarMatrix(m, n, p).GetDeep2(), [][]float64{
				{36, 0}, // g0⋅g0 = 6*6
				{0, 9},  // g1⋅g1 = 3*3
			})
			chk.Deep2(tst, "g^ij("+t+")", 1e-14, g.ContraMatrix(m, n, p).GetDeep2(), [][]float64{
				{9.0 / det, -0},
				{-0, 36.0 / det},
			})
			chk.Array(tst, "g^0("+t+")", 1e-14, g.ContraBasis(m, n, p, 0), []float64{6 * 9 / det, 0})
			chk.Array(tst, "g^1("+t+")", 1e-14, g.ContraBasis(m, n, p, 1), []float64{0, 3 * 36 / det})
			chk.Float64(tst, "det(g)("+t+")", 1e-14, g.DetCovarMatrix(m, n, p), det)
			chk.Float64(tst, "Γ("+t+"; 0,0,0)", 1e-14, g.GammaS(m, n, p, 0, 0, 0), 0)
			chk.Float64(tst, "L("+t+"; 0)", 1e-14, g.Lcoeff(0, 0, 0, 0), 0)
		}
	}

	// node accessors
	io.Pf("\nnode accessors\n")
	idx := 0
	for n := 0; n < 4; n++ {
		for m := 0; m < 5; m++ {
			if g.IndexMNPtoI(m, n, -1) != idx {
				tst.Errorf("MNP(%d,%d,-1) to I failed\n", m, n)
				return
			}
			M, N, P := g.IndexItoMNP(idx)
			if M != m || N != n || P != 0 {
				tst.Errorf("I to MNP(%d,%d,-1) failed\n", m, n)
				return
			}
			idx++
		}
	}
	chk.Array(tst, "Node( 0)", 1e-17, g.Node(0), []float64{-6, -3})
	chk.Array(tst, "Node( 7)", 1e-17, g.Node(7), []float64{0, -1})
	chk.Array(tst, "Node( 9)", 1e-17, g.Node(9), []float64{6, -1})
	chk.Array(tst, "Node(10)", 1e-17, g.Node(10), []float64{-6, 1})
	chk.Array(tst, "Node(15)", 1e-17, g.Node(15), []float64{-6, 3})
	chk.Array(tst, "Node(19)", 1e-17, g.Node(19), []float64{6, 3})
	V := g.MapMeshgrid2d([]float64{
		100, 101, 102, 103, 104,
		200, 201, 202, 203, 204,
		300, 301, 302, 303, 304,
		400, 401, 402, 403, 404,
	})
	chk.Deep2(tst, "V", 1e-17, V, [][]float64{
		{100, 101, 102, 103, 104},
		{200, 201, 202, 203, 204},
		{300, 301, 302, 303, 304},
		{400, 401, 402, 403, 404},
	})

	// boundaries and tags
	io.Pf("\nboundaries and tags\n")
	chk.Ints(tst, "Edge(0) => xmin", g.Edge(0), []int{0, 5, 10, 15})
	chk.Ints(tst, "Edge(1) => xmax", g.Edge(1), []int{4, 9, 14, 19})
	chk.Ints(tst, "Edge(2) => ymin", g.Edge(2), []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "Edge(3) => ymax", g.Edge(3), []int{15, 16, 17, 18, 19})
	chk.Ints(tst, "EdgeGivenTag(10) => xmin", g.EdgeGivenTag(10), g.Edge(0))
	chk.Ints(tst, "EdgeGivenTag(11) => xmax", g.EdgeGivenTag(11), g.Edge(1))
	chk.Ints(tst, "EdgeGivenTag(20) => ymin", g.EdgeGivenTag(20), g.Edge(2))
	chk.Ints(tst, "EdgeGivenTag(21) => ymax", g.EdgeGivenTag(21), g.Edge(3))
	if g.EdgeGivenTag(123) != nil {
		tst.Error("g.EdgeGivenTag(123) should be nil\n")
		return
	}
	if g.Face(0) != nil {
		tst.Error("g.Face(0) should be nil\n")
		return
	}
	if g.FaceGivenTag(100) != nil {
		tst.Error("g.FaceGivenTag(100) should be nil\n")
		return
	}
	chk.Ints(tst, "Boundary(10) => xmin", g.Boundary(10), g.Edge(0))
	chk.Ints(tst, "Boundary(11) => xmax", g.Boundary(11), g.Edge(1))
	chk.Ints(tst, "Boundary(20) => ymin", g.Boundary(20), g.Edge(2))
	chk.Ints(tst, "Boundary(21) => ymax", g.Boundary(21), g.Edge(3))
	if g.Boundary(100) != nil {
		tst.Error("g.Boundary(100) should be nil\n")
		return
	}

	// unit normal
	io.Pl()
	N := la.NewVector(g.Ndim())
	for _, I := range []int{0, 5, 15} {
		g.UnitNormal(N, 10, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{-1, 0})
	}
	for _, I := range []int{4, 14, 19} {
		g.UnitNormal(N, 11, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{+1, 0})
	}
	for _, I := range []int{0, 2, 4} {
		g.UnitNormal(N, 20, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{0, -1})
	}
	for _, I := range []int{15, 17, 19} {
		g.UnitNormal(N, 21, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{0, +1})
	}

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

	// grid
	g := new(Grid)
	g.RectGenUniform([]float64{-2, -2, -2}, []float64{-1, 2, 0}, []int{3, 4, 2})

	// accessors
	io.Pf("\naccessors\n")
	xx, yy, zz := g.Meshgrid3d()
	chk.Int(tst, "Ndim", g.Ndim(), 3)
	chk.Ints(tst, "Npts", []int{g.Npts(0), g.Npts(1), g.Npts(2)}, []int{3, 4, 2})
	chk.Int(tst, "Size", g.Size(), 24)
	chk.Array(tst, "Umin", 1e-14, []float64{g.Umin(0), g.Umin(1), g.Umin(2)}, []float64{-1, -1, -1})
	chk.Array(tst, "Umax", 1e-14, []float64{g.Umax(0), g.Umax(1), g.Umax(2)}, []float64{+1, +1, +1})
	chk.Array(tst, "Xmin", 1e-17, []float64{g.Xmin(0), g.Xmin(1), g.Xmin(2)}, []float64{-2, -2, -2})
	chk.Array(tst, "Xmax", 1e-17, []float64{g.Xmax(0), g.Xmax(1), g.Xmax(2)}, []float64{-1, 2, 0})
	chk.Array(tst, "Xlen", 1e-17, []float64{g.Xlen(0), g.Xlen(1), g.Xlen(2)}, []float64{1, 4, 2})
	chk.Deep3(tst, "xx", 1e-17, xx, [][][]float64{
		{ // k=0
			{-2, -1.5, -1}, // j=0
			{-2, -1.5, -1}, // j=1
			{-2, -1.5, -1}, // j=2
			{-2, -1.5, -1}, // j=3
		},
		{ // k=1
			{-2, -1.5, -1}, // j=0
			{-2, -1.5, -1}, // j=1
			{-2, -1.5, -1}, // j=2
			{-2, -1.5, -1}, // j=3
		},
	})
	a := -2 + 4.0/3.0
	b := a + 4.0/3.0
	chk.Deep3(tst, "yy", 1e-14, yy, [][][]float64{
		{ //k=0
			{-2, -2, -2}, // j=0
			{a, a, a},    // j=1
			{b, b, b},    // j=2
			{2, 2, 2},    // j=3
		},
		{ // k=1
			{-2, -2, -2}, // j=0
			{a, a, a},    // j=1
			{b, b, b},    // j=2
			{2, 2, 2},    // j=3
		},
	})
	chk.Deep3(tst, "zz", 1e-17, zz, [][][]float64{
		{ //k=0
			{-2, -2, -2}, // j=0
			{-2, -2, -2}, // j=1
			{-2, -2, -2}, // j=2
			{-2, -2, -2}, // j=3
		},
		{ // k=1
			{0, 0, 0}, // j=0
			{0, 0, 0}, // j=1
			{0, 0, 0}, // j=2
			{0, 0, 0}, // j=3
		},
	})

	// metrics accessors
	io.Pf("\nmetrics accessors\n")
	det := 0.25 * 4.0 * 1.0 // determinant of CovarMatrix
	chk.Array(tst, "U(0,0,0)", 1e-14, g.U(0, 0, 0), []float64{-1, -1, -1})
	chk.Array(tst, "X(0,0,0)", 1e-14, g.X(0, 0, 0), []float64{-2, -2, -2})
	chk.Array(tst, "U(1,2,1)", 1e-14, g.U(1, 2, 1), []float64{0, -1.0 + 4.0/3.0, +1})
	chk.Array(tst, "X(1,2,1)", 1e-14, g.X(1, 2, 1), []float64{-1.5, b, 0})
	for _, p := range []int{0, 1} {
		for _, n := range []int{0, 2} {
			for _, m := range []int{0, 1} {
				t := io.Sf("%d,%d,%d", m, n, p)
				chk.Array(tst, "g0("+t+")", 1e-14, g.CovarBasis(m, n, p, 0), []float64{0.5, 0, 0}) // (xmax-xmin)/2
				chk.Array(tst, "g1("+t+")", 1e-14, g.CovarBasis(m, n, p, 1), []float64{0, 2, 0})   // (ymax-ymin)/2
				chk.Array(tst, "g2("+t+")", 1e-14, g.CovarBasis(m, n, p, 2), []float64{0, 0, 1})   // (zmax-zmin)/2
				chk.Deep2(tst, "g_ij("+t+")", 1e-14, g.CovarMatrix(m, n, p).GetDeep2(), [][]float64{
					{0.25, 0, 0}, // g0⋅g0 = 0.5*0.5
					{0, 4, 0},    // g1⋅g1 = 2*2
					{0, 0, 1},    // g2⋅g2 = 1*1
				})
				io.Pf("%v\n", g.ContraMatrix(m, n, p).Print(""))
				chk.Deep2(tst, "g^ij("+t+")", 1e-14, g.ContraMatrix(m, n, p).GetDeep2(), [][]float64{
					{1.0 / 0.25, 0, 0},
					{0, 1.0 / 4.0, 0},
					{0, 0, 1.0 / 1.0},
				})
				chk.Array(tst, "g^0("+t+")", 1e-14, g.ContraBasis(m, n, p, 0), []float64{0.5 / 0.25, 0, 0})
				chk.Array(tst, "g^1("+t+")", 1e-14, g.ContraBasis(m, n, p, 1), []float64{0, 2 / 4.0, 0})
				chk.Array(tst, "g^2("+t+")", 1e-14, g.ContraBasis(m, n, p, 2), []float64{0, 0, 1})
				chk.Float64(tst, "det(g)("+t+")", 1e-14, g.DetCovarMatrix(m, n, p), det)
				chk.Float64(tst, "Γ("+t+"; 0,0,0)", 1e-14, g.GammaS(m, n, p, 0, 0, 0), 0)
				chk.Float64(tst, "L("+t+"; 0)", 1e-14, g.Lcoeff(0, 0, 0, 0), 0)
			}
		}
	}

	// node accessors
	io.Pf("\nnode accessors\n")
	idx := 0
	for p := 0; p < 2; p++ {
		for n := 0; n < 4; n++ {
			for m := 0; m < 3; m++ {
				if g.IndexMNPtoI(m, n, p) != idx {
					tst.Errorf("MNP(%d,%d,%d) to I failed\n", m, n, p)
					return
				}
				M, N, P := g.IndexItoMNP(idx)
				if M != m || N != n || P != p {
					tst.Errorf("I to MNP(%d,%d,%d) failed\n", m, n, p)
					return
				}
				idx++
			}
		}
	}
	chk.Array(tst, "Node( 0)", 1e-15, g.Node(0), []float64{-2, -2, -2})
	chk.Array(tst, "Node( 7)", 1e-15, g.Node(7), []float64{-1.5, b, -2})
	chk.Array(tst, "Node( 9)", 1e-15, g.Node(9), []float64{-2, 2, -2})
	chk.Array(tst, "Node(10)", 1e-15, g.Node(14), []float64{-1, -2, 0})
	chk.Array(tst, "Node(15)", 1e-15, g.Node(18), []float64{-2, b, 0})
	chk.Array(tst, "Node(19)", 1e-15, g.Node(23), []float64{-1, 2, 0})
	V := g.MapMeshgrid3d([]float64{
		100, 101, 102, // i:0→2, j:0, k:0
		110, 111, 112, // i:0→2, j:1, k:0
		120, 121, 122, // i:0→2, j:2, k:0
		130, 131, 132, // i:0→2, j:3, k:0
		200, 201, 202, // i:0→2, j:0, k:1
		210, 211, 212, // i:0→2, j:1, k:1
		220, 221, 222, // i:0→2, j:2, k:1
		230, 231, 232, // i:0→2, j:3, k:1
	})
	chk.Deep3(tst, "V", 1e-17, V, [][][]float64{
		{ //k=0
			{100, 101, 102}, // j=0
			{110, 111, 112}, // j=1
			{120, 121, 122}, // j=2
			{130, 131, 132}, // j=3
		},
		{ // k=1
			{200, 201, 202}, // j=0
			{210, 211, 212}, // j=1
			{220, 221, 222}, // j=2
			{230, 231, 232}, // j=3
		},
	})

	// boundaries and tags
	io.Pf("\nboundaries and tags\n")
	if g.Edge(0) != nil {
		tst.Error("g.Edge(0) should be nil\n")
		return
	}
	if g.EdgeGivenTag(10) != nil {
		tst.Error("g.EdgeGivenTag(10) should be nil\n")
		return
	}
	chk.Ints(tst, "Face(0) => xmin", g.Face(0), []int{0, 3, 6, 9, 12, 15, 18, 21})
	chk.Ints(tst, "Face(1) => xmax", g.Face(1), []int{2, 5, 8, 11, 14, 17, 20, 23})
	chk.Ints(tst, "Face(2) => ymin", g.Face(2), []int{0, 1, 2, 12, 13, 14})
	chk.Ints(tst, "Face(3) => ymax", g.Face(3), []int{9, 10, 11, 21, 22, 23})
	chk.Ints(tst, "Face(4) => zmin", g.Face(4), []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
	chk.Ints(tst, "Face(5) => zmax", g.Face(5), []int{12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23})
	chk.Ints(tst, "FaceGivenTag(100) => xmin", g.FaceGivenTag(100), g.Face(0))
	chk.Ints(tst, "FaceGivenTag(101) => xmax", g.FaceGivenTag(101), g.Face(1))
	chk.Ints(tst, "FaceGivenTag(200) => ymin", g.FaceGivenTag(200), g.Face(2))
	chk.Ints(tst, "FaceGivenTag(201) => ymax", g.FaceGivenTag(201), g.Face(3))
	chk.Ints(tst, "FaceGivenTag(300) => zmin", g.FaceGivenTag(300), g.Face(4))
	chk.Ints(tst, "FaceGivenTag(301) => zmax", g.FaceGivenTag(301), g.Face(5))
	chk.Ints(tst, "FaceGivenTag(123) => nil", g.FaceGivenTag(123), nil)
	if g.FaceGivenTag(123) != nil {
		tst.Error("g.FaceGivenTag(123) should be nil\n")
		return
	}
	chk.Ints(tst, "Boundary(100) => xmin", g.Boundary(100), g.Face(0))
	chk.Ints(tst, "Boundary(101) => xmax", g.Boundary(101), g.Face(1))
	chk.Ints(tst, "Boundary(200) => ymin", g.Boundary(200), g.Face(2))
	chk.Ints(tst, "Boundary(201) => ymax", g.Boundary(201), g.Face(3))
	chk.Ints(tst, "Boundary(300) => zmin", g.Boundary(300), g.Face(4))
	chk.Ints(tst, "Boundary(301) => zmax", g.Boundary(301), g.Face(5))
	if g.Boundary(10) != nil {
		tst.Error("g.Boundary(10) should be nil\n")
		return
	}

	// unit normal
	io.Pl()
	N := la.NewVector(g.Ndim())
	for _, I := range []int{0, 9, 21} {
		g.UnitNormal(N, 100, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{-1, 0, 0})
	}
	for _, I := range []int{2, 14, 23} {
		g.UnitNormal(N, 101, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{+1, 0, 0})
	}
	for _, I := range []int{0, 12, 14} {
		g.UnitNormal(N, 200, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{0, -1, 0})
	}
	for _, I := range []int{9, 21, 23} {
		g.UnitNormal(N, 201, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{0, +1, 0})
	}
	for _, I := range []int{0, 6, 11} {
		g.UnitNormal(N, 300, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{0, 0, -1})
	}
	for _, I := range []int{12, 20, 23} {
		g.UnitNormal(N, 301, I)
		chk.Array(tst, "unit normal", 1e-15, N, []float64{0, 0, +1})
	}

	// plot
	if chk.Verbose {
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp.Draw()
		gp.Bases(0.5)
		plt.Grid(&plt.A{C: "grey"})
		plt.Triad(3, "x", "y", "z", &plt.A{C: "orange"}, nil)
		plt.Default3dView(-2, 2, -2, 2, -2, 2, true)
		plt.Save("/tmp/gosl/gm", "grid02")
	}
}

func TestGrid03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid03. rectangular uniform (RectSet2d)")

	// grid
	g := new(Grid)
	g.RectSet2d([]float64{1, 2, 4, 8, 16}, []float64{0, 3, 4, 7})

	// check
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
	chk.Array(tst, "Node( 0)", 1e-15, g.Node(0), []float64{1, 0})
	chk.Array(tst, "Node( 8)", 1e-15, g.Node(8), []float64{8, 3})
	chk.Array(tst, "Node(14)", 1e-15, g.Node(14), []float64{16, 4})
	chk.Array(tst, "Node(19)", 1e-15, g.Node(19), []float64{16, 7})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp := GridPlotter{G: g, WithVids: true}
		gp.Draw()
		gp.Bases(1)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.SetXnticks(19)
		plt.SetYnticks(17)
		plt.Save("/tmp/gosl/gm", "grid03")
	}
}

func TestGrid04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid04. rectangular uniform (RectSet3d)")

	// grid
	g := new(Grid)
	g.RectSet3d([]float64{1, 2, 4, 8}, []float64{0, 3, 4}, []float64{-1, -0.5})

	// check
	chk.Ints(tst, "Face(0) => xmin", g.Face(0), []int{0, 4, 8, 12, 16, 20})
	chk.Ints(tst, "Face(0) => xmax", g.Face(1), []int{3, 7, 11, 15, 19, 23})
	chk.Ints(tst, "Face(0) => ymin", g.Face(2), []int{0, 1, 2, 3, 12, 13, 14, 15})
	chk.Ints(tst, "Face(0) => ymax", g.Face(3), []int{8, 9, 10, 11, 20, 21, 22, 23})
	chk.Ints(tst, "Face(0) => zmin", g.Face(4), []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
	chk.Ints(tst, "Face(0) => zmax", g.Face(5), []int{12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23})
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
	chk.Array(tst, "Node( 0)", 1e-17, g.Node(0), []float64{1, 0, -1})
	chk.Array(tst, "Node( 1)", 1e-17, g.Node(1), []float64{2, 0, -1})
	chk.Array(tst, "Node( 6)", 1e-17, g.Node(6), []float64{4, 3, -1})
	chk.Array(tst, "Node( 8)", 1e-17, g.Node(8), []float64{1, 4, -1})
	chk.Array(tst, "Node(11)", 1e-17, g.Node(11), []float64{8, 4, -1})
	chk.Array(tst, "Node(12)", 1e-17, g.Node(12), []float64{1, 0, -0.5})
	chk.Array(tst, "Node(17)", 1e-17, g.Node(17), []float64{2, 3, -0.5})
	chk.Array(tst, "Node(19)", 1e-17, g.Node(19), []float64{8, 3, -0.5})
	chk.Array(tst, "Node(22)", 1e-17, g.Node(22), []float64{4, 4, -0.5})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp := GridPlotter{G: g, WithVids: true}
		gp.Draw()
		gp.Bases(0.5)
		plt.Grid(&plt.A{C: "grey"})
		plt.DefaultTriad(1)
		plt.Default3dView(g.Xmin(0), g.Xmax(0), g.Xmin(1), g.Xmax(1), g.Xmin(2), g.Xmax(2), true)
		plt.Save("/tmp/gosl/gm", "grid04")
	}
}

func TestGrid05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid05. 2d ring (using Transfinite)")

	// mapping
	a, b := 1.0, 2.0
	trf := FactoryTfinite.Surf2dQuarterRing(a, b)

	// coordinates
	R := utl.LinSpace(-1, 1, 5)
	S := utl.LinSpace(-1, 1, 5)

	// curvgrid
	g := new(Grid)
	g.SetTransfinite2d(trf, R, S)

	// check limits
	chk.Array(tst, "umin", 1e-15, g.umin, []float64{-1, -1, -1})
	chk.Array(tst, "umax", 1e-15, g.umax, []float64{+1, +1, -1})
	chk.Array(tst, "xmin", 1e-15, g.xmin, []float64{0, 0, 0})
	chk.Array(tst, "xmax", 1e-15, g.xmax, []float64{b, b, 0})

	// check metrics
	π := math.Pi
	A := (b - a) / 2.0 // dρ/dr
	B := π / 4.0       // dα/ds
	p := 0             // z-index
	for n := 0; n < g.npts[1]; n++ {
		for m := 0; m < g.npts[0]; m++ {
			mtr := g.mtr[p][n][m]
			ρ := a + (1.0+mtr.U[0])*A // cylindrical coordinates
			α := (1.0 + mtr.U[1]) * B // cylindrical coordinates
			cα, sα := math.Cos(α), math.Sin(α)
			chk.Array(tst, "x", 1e-14, mtr.X, []float64{ρ * cα, ρ * sα})
			chk.Array(tst, "g_0", 1e-14, mtr.CovG0, []float64{cα * A, sα * A})
			chk.Array(tst, "g_1", 1e-14, mtr.CovG1, []float64{-ρ * sα * B, ρ * cα * B})
			chk.Deep2(tst, "CovGmat", 1e-14, mtr.CovGmat.GetDeep2(), [][]float64{
				{A * A, 0.0},
				{0.0, ρ * ρ * B * B},
			})
			chk.Deep2(tst, "CntGmat", 1e-14, mtr.CntGmat.GetDeep2(), [][]float64{
				{1.0 / (A * A), 0.0},
				{0.0, 1.0 / (ρ * ρ * B * B)},
			})
			chk.Array(tst, "g^0", 1e-14, mtr.CntG0, []float64{cα * A / (A * A), sα * A / (A * A)})
			chk.Array(tst, "g^1", 1e-14, mtr.CntG1, []float64{-ρ * sα * B / (ρ * ρ * B * B), ρ * cα * B / (ρ * ρ * B * B)})
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

	// unit normal
	io.Pl()
	N := la.NewVector(g.Ndim())
	g.UnitNormal(N, 10, 0)
	chk.Array(tst, "unit normal (0)", 1e-15, N, []float64{-1, 0})
	g.UnitNormal(N, 10, 10)
	chk.Array(tst, "unit normal (10)", 1e-15, N, []float64{-1.0 / math.Sqrt2, -1.0 / math.Sqrt2})
	g.UnitNormal(N, 10, 20)
	chk.Array(tst, "unit normal (20)", 1e-15, N, []float64{0, -1})
	g.UnitNormal(N, 11, 4)
	chk.Array(tst, "unit normal (4)", 1e-15, N, []float64{+1, 0})
	g.UnitNormal(N, 11, 14)
	chk.Array(tst, "unit normal (14)", 1e-15, N, []float64{+1.0 / math.Sqrt2, +1.0 / math.Sqrt2})
	g.UnitNormal(N, 11, 24)
	chk.Array(tst, "unit normal (24)", 1e-15, N, []float64{0, +1})
	for _, I := range []int{0, 2, 4} {
		g.UnitNormal(N, 20, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, -1})
	}
	for _, I := range []int{20, 22, 24} {
		g.UnitNormal(N, 21, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{-1, 0})
	}

	// check derivatives
	checkGridTfiniteDerivs2d(tst, trf, g, 1e-10, 1e-9, 1e-8, chk.Verbose)

	// check interface functions
	io.Pl()
	chk.Int(tst, "Ndim()", g.Ndim(), 2)
	chk.Int(tst, "Npts(0)", g.Npts(0), len(R))
	chk.Int(tst, "Size()", g.Size(), len(R)*len(S))
	chk.Float64(tst, "Umin(0)", 1e-14, g.Umin(0), -1)
	chk.Float64(tst, "Umax(0)", 1e-14, g.Umax(0), +1)
	chk.Float64(tst, "Xmin(0)", 1e-14, g.Xmin(0), 0)
	chk.Float64(tst, "Xmax(0)", 1e-14, g.Xmax(0), b)
	chk.Array(tst, "U(0,0,0)", 1e-14, g.U(0, 0, 0), []float64{-1, -1})
	chk.Array(tst, "X(0,0,0)", 1e-14, g.X(0, 0, 0), []float64{a, 0})
	chk.Array(tst, "g_0(0,0,0)", 1e-14, g.CovarBasis(0, 0, 0, 0), []float64{A, 0})
	chk.Array(tst, "g_1(0,0,0)", 1e-14, g.CovarBasis(0, 0, 0, 1), []float64{0, a * B})
	chk.Array(tst, "g_2(0,0,0)", 1e-14, g.CovarBasis(0, 0, 0, 2), nil)
	chk.Array(tst, "g^0(0,0,0)", 1e-14, g.ContraBasis(0, 0, 0, 0), []float64{A / (A * A), 0})
	chk.Array(tst, "g^1(0,0,0)", 1e-14, g.ContraBasis(0, 0, 0, 1), []float64{0, a * B / (a * a * B * B)})
	chk.Array(tst, "g^2(0,0,0)", 1e-14, g.ContraBasis(0, 0, 0, 2), nil)
	chk.Deep2(tst, "g_ij(0,0,0)", 1e-14, g.CovarMatrix(0, 0, 0).GetDeep2(), [][]float64{
		{A * A, 0},
		{0, a * a * B * B},
	})
	chk.Deep2(tst, "g^ij(0,0,0)", 1e-14, g.ContraMatrix(0, 0, 0).GetDeep2(), [][]float64{
		{1.0 / (A * A), 0},
		{0, 1.0 / (a * a * B * B)},
	})
	chk.Float64(tst, "det(g)(0,0,0)", 1e-14, g.DetCovarMatrix(0, 0, 0), A*A*a*a*B*B)
	chk.Float64(tst, "Γ(0,0,0; 0,1,1)", 1e-14, g.GammaS(0, 0, 0, 0, 1, 1), -a*B*B/A)
	chk.Float64(tst, "L(0,0,0; 0)", 1e-14, g.Lcoeff(0, 0, 0, 0), -1.0/(a*A))

	// plot
	if chk.Verbose {
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{11, 21}, false, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		gp.Draw()
		gp.Bases(0.15)
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "grid05")
	}
}

func TestGrid06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid06. 3d ring (using Transfinite)")

	// mapping
	a, b, h := 2.0, 3.0, 2.0 // radii and thickness
	trf := FactoryTfinite.SolidQuarterRing(a, b, h)

	// coordinates
	npts := 3
	R := utl.LinSpace(-1, 1, npts)
	S := utl.LinSpace(-1, 1, npts)
	T := utl.LinSpace(-1, 1, npts)

	// curvgrid
	g := new(Grid)
	g.SetTransfinite3d(trf, R, S, T)

	// check limits
	chk.Array(tst, "umin", 1e-15, g.umin, []float64{-1, -1, -1})
	chk.Array(tst, "umax", 1e-15, g.umax, []float64{+1, +1, +1})
	chk.Array(tst, "xmin", 1e-15, g.xmin, []float64{0, 0, 0})
	chk.Array(tst, "xmax", 1e-15, g.xmax, []float64{h, b, b})

	// check
	π := math.Pi
	A := (b - a) / 2.0 // dρ/dr
	B := π / 4.0       // dα/ds
	for p := 0; p < g.npts[2]; p++ {
		for n := 0; n < g.npts[1]; n++ {
			for m := 0; m < g.npts[0]; m++ {
				mtr := g.mtr[p][n][m]
				x0 := h * float64(m) / float64(g.npts[0]-1)
				ρ := a + (1.0+mtr.U[1])*A // cylindrical coordinates
				α := (1.0 + mtr.U[2]) * B // cylindrical coordinates
				cα, sα := math.Cos(α), math.Sin(α)
				chk.Array(tst, "x", 1e-14, mtr.X, []float64{x0, ρ * cα, ρ * sα})
				chk.Array(tst, "g_0", 1e-14, mtr.CovG0, []float64{1, 0, 0})
				chk.Array(tst, "g_1", 1e-14, mtr.CovG1, []float64{0, cα * A, sα * A})
				chk.Array(tst, "g_2", 1e-14, mtr.CovG2, []float64{0, -ρ * sα * B, ρ * cα * B})
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
				chk.Array(tst, "g^0", 1e-14, mtr.CntG0, []float64{1, 0, 0})
				chk.Array(tst, "g^1", 1e-14, mtr.CntG1, []float64{0, cα * A / (A * A), sα * A / (A * A)})
				chk.Array(tst, "g^2", 1e-14, mtr.CntG2, []float64{0, -ρ * sα * B / (ρ * ρ * B * B), ρ * cα * B / (ρ * ρ * B * B)})
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

	// unit normal
	dv := 1.0 / math.Sqrt2
	io.Pl()
	N := la.NewVector(g.Ndim())
	// n0-planes
	for _, I := range []int{0, 15, 24} {
		g.UnitNormal(N, 100, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{-1, 0, 0})
	}
	for _, I := range []int{2, 11, 26} {
		g.UnitNormal(N, 101, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{+1, 0, 0})
	}
	// n1-plane: inner cylindrical face
	for _, I := range []int{0, 1, 2} {
		g.UnitNormal(N, 200, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, -1, 0})
	}
	for _, I := range []int{9, 10, 11} {
		g.UnitNormal(N, 200, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, -dv, -dv})
	}
	for _, I := range []int{18, 19, 20} {
		g.UnitNormal(N, 200, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, 0, -1})
	}
	// n1-plane: outer cylindrical face
	for _, I := range []int{6, 7, 8} {
		g.UnitNormal(N, 201, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, +1, 0})
	}
	for _, I := range []int{15, 16, 17} {
		g.UnitNormal(N, 201, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, +dv, +dv})
	}
	for _, I := range []int{24, 25, 26} {
		g.UnitNormal(N, 201, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, 0, +1})
	}
	// n2-planes
	for _, I := range []int{0, 4, 8} {
		g.UnitNormal(N, 300, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, 0, -1})
	}
	for _, I := range []int{18, 22, 26} {
		g.UnitNormal(N, 301, I)
		chk.Array(tst, io.Sf("unit normal (%d)", I), 1e-15, N, []float64{0, -1, 0})
	}

	// check derivatives
	checkGridTfiniteDerivs3d(tst, trf, g, 1e-10, 1e-9, 1e-8, chk.Verbose)

	// check interface functions
	io.Pl()
	chk.Int(tst, "Ndim()", g.Ndim(), 3)
	chk.Int(tst, "Npts(0)", g.Npts(0), len(R))
	chk.Int(tst, "Size()", g.Size(), len(R)*len(S)*len(T))
	chk.Float64(tst, "Umin(2)", 1e-14, g.Umin(2), -1)
	chk.Float64(tst, "Umax(2)", 1e-14, g.Umax(2), +1)
	chk.Float64(tst, "Xmin(2)", 1e-14, g.Xmin(2), 0)
	chk.Float64(tst, "Xmax(2)", 1e-14, g.Xmax(2), b)
	chk.Array(tst, "U(0,0,0)", 1e-14, g.U(0, 0, 0), []float64{-1, -1, -1})
	chk.Array(tst, "X(0,0,0)", 1e-14, g.X(0, 0, 0), []float64{0, a, 0})
	chk.Array(tst, "g_0(0,0,0)", 1e-14, g.CovarBasis(0, 0, 0, 0), []float64{1, 0, 0})
	chk.Array(tst, "g_1(0,0,0)", 1e-14, g.CovarBasis(0, 0, 0, 1), []float64{0, A, 0})
	chk.Array(tst, "g_2(0,0,0)", 1e-14, g.CovarBasis(0, 0, 0, 2), []float64{0, 0, a * B})
	chk.Array(tst, "g^0(0,0,0)", 1e-14, g.ContraBasis(0, 0, 0, 0), []float64{1, 0, 0})
	chk.Array(tst, "g^1(0,0,0)", 1e-14, g.ContraBasis(0, 0, 0, 1), []float64{0, A / (A * A), 0})
	chk.Array(tst, "g^2(0,0,0)", 1e-14, g.ContraBasis(0, 0, 0, 2), []float64{0, 0, a * B / (a * a * B * B)})
	chk.Deep2(tst, "g_ij(0,0,0)", 1e-14, g.CovarMatrix(0, 0, 0).GetDeep2(), [][]float64{
		{1.0, 0.0, 0.0},
		{0.0, A * A, 0.0},
		{0.0, 0.0, a * a * B * B},
	})
	chk.Deep2(tst, "g^ij(0,0,0)", 1e-14, g.ContraMatrix(0, 0, 0).GetDeep2(), [][]float64{
		{1.0, 0.0, 0.0},
		{0.0, 1.0 / (A * A), 0.0},
		{0.0, 0.0, 1.0 / (a * a * B * B)},
	})
	chk.Float64(tst, "det(g)(0,0,0)", 1e-14, g.DetCovarMatrix(0, 0, 0), A*A*a*a*B*B)
	chk.Float64(tst, "Γ(0,0,0; 1,2,2)", 1e-14, g.GammaS(0, 0, 0, 1, 2, 2), -a*B*B/A)
	chk.Float64(tst, "L(0,0,0; 1)", 1e-14, g.Lcoeff(0, 0, 0, 1), -1.0/(a*A))

	// plot
	if chk.Verbose {
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 400})
		trf.Draw([]int{5, 5, 11}, true, &plt.A{C: "#7d8891"}, &plt.A{C: plt.C(3, 9), Lw: 2})
		gp.Draw()
		gp.Bases(0.20)
		plt.Default3dView(0, 3, 0, 3, 0, 3, true)
		plt.Save("/tmp/gosl/gm", "grid06")
	}
}

func TestGrid07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid07. using 2D NURBS surface (flat)")

	// nurbs
	verts := [][]float64{
		{0, 0, 0, 1},
		{4, 1, 0, 1},
		{-1, 4, 0, 1},
		{3, 3, 0, 1},
	}
	knots := [][]float64{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
	}
	nrb := NewNurbs(2, []int{1, 1}, knots)
	nrb.SetControl(verts, utl.IntRange(len(verts)))

	// coordinates
	R := utl.LinSpace(-1, 1, 3)
	S := utl.LinSpace(-1, 1, 3)

	// grid
	g := new(Grid)
	g.SetNurbsSurf2d(nrb, R, S)

	// check
	verb := chk.Verbose
	checkGridNurbsDerivs2d(tst, nrb, g, 1e-12, 1e-12, 1e-9, verb)

	// plot
	if chk.Verbose {
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.HideAllBorders()
		PlotNurbs("/tmp/gosl/gm", "grid07", nrb, 2, 41, true, true, nil, nil, nil, func() {
			gp.Draw()
			gp.Bases(0.5)
			nrb.DrawSurface(2, 5, 5, false, true, nil, nil)
			plt.AxisOff()
			plt.Equal()
		})
	}
}

func TestGrid08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid08. using 2D NURBS (quarter ring)")

	// nurbs
	nrb := FactoryNurbs.Surf2dQuarterRing(0, 0, 1, 3)

	// coordinates
	R := utl.LinSpace(-1, 1, 3)
	S := utl.LinSpace(-1, 1, 3)

	// grid
	g := new(Grid)
	g.SetNurbsSurf2d(nrb, R, S)

	// check
	verb := chk.Verbose
	checkGridNurbsDerivs2d(tst, nrb, g, 1e-10, 1e-12, 1e-8, verb)

	// plot
	if chk.Verbose {
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.HideAllBorders()
		PlotNurbs("/tmp/gosl/gm", "grid08", nrb, 2, 21, true, true, nil, nil, nil, func() {
			gp.Draw()
			gp.Bases(0.5)
			nrb.DrawSurface(2, 11, 11, false, true, nil, nil)
			plt.AxisOff()
			plt.Equal()
		})
	}
}

func TestGrid09(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid09. using 2D NURBS (wavy domain)")

	// nurbs
	nrb := FactoryNurbs.Surf2dExample1()

	// coordinates
	R := utl.LinSpace(-1, 1, 3)
	S := utl.LinSpace(-1, 1, 3)

	// grid
	g := new(Grid)
	g.SetNurbsSurf2d(nrb, R, S)

	// check
	verb := chk.Verbose
	checkGridNurbsDerivs2d(tst, nrb, g, 1e-10, 1e-12, 1e-7, verb)

	// plot
	if chk.Verbose {
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.HideAllBorders()
		PlotNurbs("/tmp/gosl/gm", "grid09", nrb, 2, 11, true, true, nil, nil, nil, func() {
			gp.Draw()
			gp.Bases(0.1)
			//nrb.DrawSurface(2, 11, 11, false, true, nil, nil)
			plt.AxisOff()
			plt.Equal()
		})
	}
}

func TestGrid10(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid10. using 3D NURBS (solid hexahedron)")

	// nurbs
	nrb := FactoryNurbs.SolidHex([][]float64{
		{0.0, 0, 0.0}, // 0
		{1.0, 0, 0.0}, // 1
		{0.0, 2, 0.0}, // 2
		{2.0, 2, 0.0}, // 3
		{0.0, 0, 1.0}, // 4
		{0.8, 0, 0.8}, // 5
		{0.0, 2, 2.0}, // 6
		{2.0, 2, 2.0}, // 7
	})

	// coordinates
	R := utl.LinSpace(-1, 1, 3)
	S := utl.LinSpace(-1, 1, 3)
	T := utl.LinSpace(-1, 1, 3)

	// grid
	g := new(Grid)
	g.SetNurbsSolid(nrb, R, S, T)

	// check
	verb := chk.Verbose
	checkGridNurbsDerivs3d(tst, nrb, g, 1e-10, 1e-12, 1e-8, verb)

	// plot
	if chk.Verbose {
		npts := 3
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		gp.Draw()
		nrb.DrawSolid(2, 4, 3, &plt.A{C: plt.C(1, 0)})
		nrb.DrawElems(3, npts, true, &plt.A{C: plt.C(0, 0)}, &plt.A{C: "k", Fsz: 7})
		nrb.DrawCtrl(3, true, nil, nil)
		plt.Triad(0.5, "x", "y", "z", &plt.A{C: "orange"}, &plt.A{C: "green"})
		plt.Default3dView(0, 2, 0, 2, 0, 2, true)
		plt.Save("/tmp/gosl/gm", "grid10")
	}
}

func TestGrid11(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid11. using solid NURBS (quarter ring)")

	// nurbs
	nrb := FactoryNurbs.SolidQuarterRing(0, 0, 0, 2, 3, 1)

	// coordinates
	R := utl.LinSpace(-1, 1, 3)
	S := utl.LinSpace(-1, 1, 3)
	T := utl.LinSpace(-1, 1, 3)

	// grid
	g := new(Grid)
	g.SetNurbsSolid(nrb, R, S, T)

	// check
	verb := chk.Verbose
	checkGridNurbsDerivs3d(tst, nrb, g, 1e-10, 1e-12, 1e-7, verb)

	// plot
	if chk.Verbose {
		gp := GridPlotter{G: g, WithVids: true}
		plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
		gp.Draw()
		nrb.DrawSolid(3, 3, 9, &plt.A{C: plt.C(1, 0)})
		nrb.DrawCtrl(3, false, nil, nil)
		plt.Triad(0.5, "x", "y", "z", &plt.A{C: "orange"}, &plt.A{C: "green"})
		plt.Default3dView(0, 3, 0, 3, 0, 3, true)
		plt.Save("/tmp/gosl/gm", "grid11")
	}
}

func TestGrid12(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid12. Quarter ring: NURBS vs Transfinite")

	// nurbs and transfinite
	a, b, h := 2.0, 3.0, 1.0
	nrb := FactoryNurbs.SolidQuarterRing(0, 0, 0, a, b, h)
	trf := FactoryTfinite.SolidQuarterRing(a, b, h)

	// coordinates
	R := utl.LinSpace(-1, 1, 3)
	S := utl.LinSpace(-1, 1, 3)
	T := utl.LinSpace(-1, 1, 7)

	// grids
	gnrb := new(Grid)
	gtrf := new(Grid)
	gnrb.SetNurbsSolid(nrb, R, S, T)
	gtrf.SetTransfinite3d(trf, R, S, T)

	// check
	tolx := 0.05  // yes: there's a difference between Tfinite and NURBS
	tolg1 := 0.01 // yes: there's a difference between Tfinite and NURBS
	tolg2 := 0.25 // yes: there's a difference between Tfinite and NURBS
	for p, t := range T {
		for n, s := range S {
			for m, r := range R {
				io.Pf("\nrst = (%v,%v,%v)\n", r, s, t)
				chk.Array(tst, "x", tolx, gnrb.X(m, n, p), gtrf.X(m, n, p))
				chk.Array(tst, "g0", 1e-15, gnrb.CovarBasis(m, n, p, 0), gtrf.CovarBasis(m, n, p, 0))
				chk.Array(tst, "g1", tolg1, gnrb.CovarBasis(m, n, p, 1), gtrf.CovarBasis(m, n, p, 1))
				chk.Array(tst, "g2", tolg2, gnrb.CovarBasis(m, n, p, 2), gtrf.CovarBasis(m, n, p, 2))
			}
		}
	}

	// plot
	if chk.Verbose {
		gpnrb := GridPlotter{G: gnrb, ArgsEdges: &plt.A{C: plt.C(5, 0), Lw: 4}}
		gptrf := GridPlotter{G: gtrf, ArgsEdges: &plt.A{C: plt.C(0, 0)}}
		plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
		gpnrb.Draw()
		gptrf.Draw()
		plt.Triad(0.5, "x", "y", "z", &plt.A{C: "orange"}, &plt.A{C: "green"})
		plt.Default3dView(0, 3, 0, 3, 0, 3, true)
		plt.Save("/tmp/gosl/gm", "grid12")
	}
}

func TestGrid13(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid13. rectangular uniform (RectSet2dU)")

	// grid
	R := []float64{-1, -0.5, 0, +0.5, +1}
	S := []float64{-1, 0, +1}
	g := new(Grid)
	g.RectSet2dU([]float64{2, 5}, []float64{4, 6}, R, S)

	// check
	xx, yy := g.Meshgrid2d()
	chk.Deep2(tst, "xx", 1e-17, xx, [][]float64{
		{2, 2.5, 3, 3.5, 4},
		{2, 2.5, 3, 3.5, 4},
		{2, 2.5, 3, 3.5, 4},
	})
	chk.Deep2(tst, "yy", 1e-17, yy, [][]float64{
		{5, 5, 5, 5, 5},
		{5.5, 5.5, 5.5, 5.5, 5.5},
		{6, 6, 6, 6, 6},
	})
	rr := make([]float64, g.npts[0])
	for n := 0; n < g.npts[1]; n++ {
		for m := 0; m < g.npts[0]; m++ {
			rr[m] = g.U(m, n, 0)[0]
		}
		chk.Array(tst, "R", 1e-17, rr, R)
	}
	ss := make([]float64, g.npts[1])
	for m := 0; m < g.npts[0]; m++ {
		for n := 0; n < g.npts[1]; n++ {
			ss[n] = g.U(m, n, 0)[1]
		}
		chk.Array(tst, "S", 1e-17, ss, S)
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp := GridPlotter{G: g, WithVids: true}
		gp.Draw()
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/gm", "grid13")
	}
}

func TestGrid14(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid14. rectangular uniform (RectSet3dU)")

	// grid
	lu := 2.0 / 7.0
	R := []float64{-1, -1 + lu, -1 + 3*lu, 1}
	S := []float64{-1, 0.5, +1}
	T := []float64{-1, +1}
	g := new(Grid)
	g.RectSet3dU([]float64{1, 0, -1}, []float64{8, 4, -0.5}, R, S, T)

	// check
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
	rr := make([]float64, g.npts[0])
	ss := make([]float64, g.npts[1])
	for p := 0; p < g.npts[2]; p++ {
		for n := 0; n < g.npts[1]; n++ {
			for m := 0; m < g.npts[0]; m++ {
				rr[m] = g.U(m, n, p)[0]
			}
			chk.Array(tst, "R", 1e-17, rr, R)
		}
		for m := 0; m < g.npts[0]; m++ {
			for n := 0; n < g.npts[1]; n++ {
				ss[n] = g.U(m, n, p)[1]
			}
			chk.Array(tst, "S", 1e-17, ss, S)
		}
	}
	tt := make([]float64, g.npts[2])
	for n := 0; n < g.npts[1]; n++ {
		for m := 0; m < g.npts[0]; m++ {
			for p := 0; p < g.npts[2]; p++ {
				tt[p] = g.U(m, n, p)[2]
			}
			chk.Array(tst, "T", 1e-17, tt, T)
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp := GridPlotter{G: g, WithVids: true}
		gp.Draw()
		gp.Bases(0.5)
		plt.Grid(&plt.A{C: "grey"})
		plt.DefaultTriad(1)
		plt.Default3dView(g.Xmin(0), g.Xmax(0), g.Xmin(1), g.Xmax(1), g.Xmin(2), g.Xmax(2), true)
		plt.Save("/tmp/gosl/gm", "grid14")
	}
}
