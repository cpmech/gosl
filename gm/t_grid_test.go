// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestGrid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid01")

	g, err := NewUniformGrid([]float64{-6, -3}, []float64{6, 3}, []int{4, 3})
	status(tst, err)

	chk.Int(tst, "N", g.N, 20)
	chk.Int(tst, "nx", g.Npts[0], 5)
	chk.Int(tst, "ny", g.Npts[1], 4)

	chk.Float64(tst, "Lx", 1e-15, g.Xdel[0], 12.0)
	chk.Float64(tst, "Ly", 1e-15, g.Xdel[1], 6.0)
	chk.Float64(tst, "Dx", 1e-15, g.Size[0], 3.0)
	chk.Float64(tst, "Dy", 1e-15, g.Size[1], 2.0)

	chk.Ints(tst, "B", g.Edge[0], []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "R", g.Edge[1], []int{4, 9, 14, 19})
	chk.Ints(tst, "T", g.Edge[2], []int{15, 16, 17, 18, 19})
	chk.Ints(tst, "L", g.Edge[3], []int{0, 5, 10, 15})

	chk.Ints(tst, "Tag # 10: L", g.GetNodesWithTag(10), []int{0, 5, 10, 15})
	chk.Ints(tst, "Tag # 11: R", g.GetNodesWithTag(11), []int{4, 9, 14, 19})
	chk.Ints(tst, "Tag # 20: B", g.GetNodesWithTag(20), []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "Tag # 21: T", g.GetNodesWithTag(21), []int{15, 16, 17, 18, 19})

	// plot
	if chk.Verbose {
		X, Y, F, err := g.Eval2d(func(x la.Vector) (float64, error) {
			return x[0]*x[0] + x[1]*x[1], nil
		})
		status(tst, err)
		plt.Reset(true, &plt.A{WidthPt: 500})
		g.Draw(true, true, true, true, nil, nil, nil, nil, nil)
		plt.Grid(&plt.A{C: "grey"})
		plt.ContourL(X, Y, F, nil)
		plt.Equal()
		plt.HideAllBorders()
		plt.SetXnticks(12)
		plt.SetYnticks(12)
		err = plt.Save("/tmp/gosl/gm", "grid01")
		status(tst, err)
	}
}

func TestGrid02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid02")

	g, err := NewGrid([]int{4, 3})
	status(tst, err)

	chk.Int(tst, "N", g.N, 20)
	chk.Int(tst, "nx", g.Npts[0], 5)
	chk.Int(tst, "ny", g.Npts[1], 4)

	chk.Ints(tst, "B", g.Edge[0], []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "R", g.Edge[1], []int{4, 9, 14, 19})
	chk.Ints(tst, "T", g.Edge[2], []int{15, 16, 17, 18, 19})
	chk.Ints(tst, "L", g.Edge[3], []int{0, 5, 10, 15})

	chk.Ints(tst, "Tag # 10: L", g.GetNodesWithTag(10), []int{0, 5, 10, 15})
	chk.Ints(tst, "Tag # 11: R", g.GetNodesWithTag(11), []int{4, 9, 14, 19})
	chk.Ints(tst, "Tag # 20: B", g.GetNodesWithTag(20), []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "Tag # 21: T", g.GetNodesWithTag(21), []int{15, 16, 17, 18, 19})

	err = g.SetCoords2d([]float64{1, 2, 4, 8, 16}, []float64{0, 3, 4, 7})
	status(tst, err)
	chk.Deep2(tst, "X2d", 1e-17, g.X2d, [][]float64{
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
		{1, 2, 4, 8, 16},
	})
	chk.Deep2(tst, "Y2d", 1e-17, g.Y2d, [][]float64{
		{0, 0, 0, 0, 0},
		{3, 3, 3, 3, 3},
		{4, 4, 4, 4, 4},
		{7, 7, 7, 7, 7},
	})

	chk.Array(tst, "Min", 1e-17, g.Min, []float64{1, 0})
	chk.Array(tst, "Max", 1e-17, g.Max, []float64{16, 7})
	chk.Array(tst, "Del", 1e-17, g.Del, []float64{15, 7})

	// plot
	if chk.Verbose {
		status(tst, err)
		plt.Reset(true, &plt.A{WidthPt: 500})
		g.Draw(true, nil, nil)
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.SetXnticks(17)
		plt.SetYnticks(17)
		err = plt.Save("/tmp/gosl/gm", "grid02")
		status(tst, err)
	}
}

func TestGrid03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Grid03")

	g, err := NewGrid([]int{3, 2, 1})
	status(tst, err)

	chk.Int(tst, "N", g.N, 24)
	chk.Int(tst, "nx", g.Npts[0], 4)
	chk.Int(tst, "ny", g.Npts[1], 3)
	chk.Int(tst, "nz", g.Npts[2], 2)

	chk.Ints(tst, "xmin", g.Face[0], []int{0, 4, 8, 12, 16, 20})
	chk.Ints(tst, "xmax", g.Face[1], []int{3, 7, 11, 15, 19, 23})
	chk.Ints(tst, "ymin", g.Face[2], []int{0, 1, 2, 3, 12, 13, 14, 15})
	chk.Ints(tst, "ymax", g.Face[3], []int{8, 9, 10, 11, 20, 21, 22, 23})
	chk.Ints(tst, "zmin", g.Face[4], []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
	chk.Ints(tst, "zmax", g.Face[5], []int{12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23})

	chk.Ints(tst, "Tag # 10: xmin", g.GetNodesWithTag(10), g.Face[0])
	chk.Ints(tst, "Tag # 11: xmax", g.GetNodesWithTag(11), g.Face[1])
	chk.Ints(tst, "Tag # 20: ymin", g.GetNodesWithTag(20), g.Face[2])
	chk.Ints(tst, "Tag # 21: ymax", g.GetNodesWithTag(21), g.Face[3])
	chk.Ints(tst, "Tag # 30: zmin", g.GetNodesWithTag(30), g.Face[4])
	chk.Ints(tst, "Tag # 31: zmax", g.GetNodesWithTag(31), g.Face[5])

	err = g.SetCoords3d([]float64{1, 2, 4, 8}, []float64{0, 3, 4}, []float64{-1, -0.5})
	status(tst, err)
	chk.Deep3(tst, "X3d", 1e-17, g.X3d, [][][]float64{
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
	chk.Deep3(tst, "Y3d", 1e-17, g.Y3d, [][][]float64{
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
	chk.Deep3(tst, "Z3d", 1e-17, g.Z3d, [][][]float64{
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

	chk.Array(tst, "Min", 1e-17, g.Min, []float64{1, 0, -1})
	chk.Array(tst, "Max", 1e-17, g.Max, []float64{8, 4, -0.5})
	chk.Array(tst, "Del", 1e-17, g.Del, []float64{7, 4, 0.5})

	// plot
	if chk.Verbose {
		status(tst, err)
		plt.Reset(true, &plt.A{WidthPt: 500})
		g.Draw(true, nil, &plt.A{Fsz: 6})
		plt.Grid(&plt.A{C: "grey"})
		plt.Equal()
		plt.HideAllBorders()
		plt.DefaultTriad(1)
		plt.Default3dView(g.Min[0], g.Max[0], g.Min[1], g.Max[1], g.Min[2], g.Max[2], true)
		err = plt.Save("/tmp/gosl/gm", "grid03")
		status(tst, err)
	}
}
