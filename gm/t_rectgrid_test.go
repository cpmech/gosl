// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
)

func TestRectGrid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("RectGrid01")

	g := new(CurvGrid)
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
		gp := GridPlotter{
			G:        g,
			WithVids: true,
		}
		plt.Reset(true, &plt.A{WidthPt: 500})
		gp.Draw()
		gp.Bases(1)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.SetXnticks(12)
		plt.SetYnticks(12)
		plt.Save("/tmp/gosl/gm", "rectgrid01")
	}
}

func TestRectGrid02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("RectGrid02")

	g := new(RectGrid)
	g.Set2d([]float64{1, 2, 4, 8, 16}, []float64{0, 3, 4, 7}, true)

	chk.Int(tst, "ndim", g.Ndim(), 2)
	chk.Int(tst, "size", g.Size(), 20)
	chk.Int(tst, "nx", g.Npts(0), 5)
	chk.Int(tst, "ny", g.Npts(1), 4)

	xx, yy := g.Mesh2d()
	min := []float64{g.Min(0), g.Min(1)}
	max := []float64{g.Max(0), g.Max(1)}
	del := []float64{g.Length(0), g.Length(1)}

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

	chk.Deep2(tst, "x2d", 1e-17, xx, [][]float64{
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
	})
	chk.Deep2(tst, "y2d", 1e-17, yy, [][]float64{
		{0, 0, 0, 0, 0},
		{3, 3, 3, 3, 3},
		{4, 4, 4, 4, 4},
		{7, 7, 7, 7, 7},
	})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		g.Draw(true, nil, nil)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.SetXnticks(17)
		plt.SetYnticks(17)
		plt.Save("/tmp/gosl/gm", "rectgrid02")
	}
}

func TestGrid03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("RectGrid03")

	g := new(RectGrid)
	g.Set3d([]float64{1, 2, 4, 8}, []float64{0, 3, 4}, []float64{-1, -0.5}, true)

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

	xx, yy, zz := g.Mesh3d()

	chk.Deep3(tst, "X3d", 1e-17, xx, [][][]float64{
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
	chk.Deep3(tst, "Y3d", 1e-17, yy, [][][]float64{
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
	chk.Deep3(tst, "Z3d", 1e-17, zz, [][][]float64{
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

	min := []float64{g.Min(0), g.Min(1), g.Min(2)}
	max := []float64{g.Max(0), g.Max(1), g.Max(2)}
	del := []float64{g.Length(0), g.Length(1), g.Length(2)}

	chk.Array(tst, "Min", 1e-17, min, []float64{1, 0, -1})
	chk.Array(tst, "Max", 1e-17, max, []float64{8, 4, -0.5})
	chk.Array(tst, "Del", 1e-17, del, []float64{7, 4, 0.5})

	chk.Array(tst, "x[0]", 1e-17, g.GetNode(0), []float64{1, 0, -1})
	chk.Array(tst, "x[1]", 1e-17, g.GetNode(1), []float64{2, 0, -1})
	chk.Array(tst, "x[6]", 1e-17, g.GetNode(6), []float64{4, 3, -1})
	chk.Array(tst, "x[8]", 1e-17, g.GetNode(8), []float64{1, 4, -1})
	chk.Array(tst, "x[11]", 1e-17, g.GetNode(11), []float64{8, 4, -1})
	chk.Array(tst, "x[12]", 1e-17, g.GetNode(12), []float64{1, 0, -0.5})
	chk.Array(tst, "x[17]", 1e-17, g.GetNode(17), []float64{2, 3, -0.5})
	chk.Array(tst, "x[19]", 1e-17, g.GetNode(19), []float64{8, 3, -0.5})
	chk.Array(tst, "x[22]", 1e-17, g.GetNode(22), []float64{4, 4, -0.5})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500})
		g.Draw(true, nil, &plt.A{Fsz: 6})
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.DefaultTriad(1)
		plt.Default3dView(g.Min(0), g.Max(0), g.Min(1), g.Max(1), g.Min(2), g.Max(2), true)
		plt.Save("/tmp/gosl/gm", "rectgrid03")
	}
}
