// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"encoding/json"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/gm/msh"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func check3x3grid(tst *testing.T, e *BoundaryConds, checkNormals bool) {
	chk.Ints(tst, "n2i", e.n2i, []int{0, 3, 4, 1, -1, -1, 2, -1, -1})
	chk.IntDeep2(tst, "tags", e.tags, [][]int{{10, 20}, {10}, {10}, {20}, {20}})
	chk.Int(tst, "len(fcns)", len(e.fcns), 5)
	chk.Ints(tst, "nodes", e.Nodes(), []int{0, 1, 2, 3, 6})
	tags0 := make([][]int, 9)
	tags1 := make([][]int, 9)
	vals0 := make([]float64, 9)
	vals1 := make([]float64, 9)
	normals := make([][]float64, 9)
	N := la.NewVector(2)
	var available0, available1 bool
	for n := 0; n < 9; n++ {
		tags0[n], vals0[n], available0 = e.Value(n, 0, 0)
		tags1[n], vals1[n], available1 = e.Value(n, 1, 0)
		if !available0 {
			vals0[n] = -1
		}
		if !available1 {
			vals1[n] = -1
		}
		if checkNormals {
			e.NormalGrid(N, n)
			normals[n] = N.GetCopy()
		}
	}
	chk.IntDeep2(tst, "tags0", tags0, [][]int{{10, 20}, {}, {}, {10}, {}, {}, {10}, {}, {}})
	chk.IntDeep2(tst, "tags1", tags1, [][]int{{10, 20}, {20}, {20}, {}, {}, {}, {}, {}, {}})
	chk.Array(tst, "vals0", 1e-14, vals0, []float64{123, -1, -1, 123, -1, -1, 123, -1, -1})
	chk.Array(tst, "vals1", 1e-14, vals1, []float64{100, 101, 102, -1, -1, -1, -1, -1, -1})
	if checkNormals {
		chk.Deep2(tst, "normals", 1e-14, normals, [][]float64{
			{-1, -1}, {0, -1}, {0, -1},
			{-1, 0}, {0, 0}, {0, 0},
			{-1, 0}, {0, 0}, {0, 0},
		})
	}
}

func TestBryConds01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BryConds01. Using Grid")

	// grid             [21]
	//          6         7         8
	//          o---------o---------o
	//          |         |         |
	//          |         |         |
	//          |3        |4        |5
	//    [10]  o---------o---------o   [11]
	//          |         |         |
	//          |         |         |
	//          |0        |1        |2
	//          o---------o---------o
	//                  [20]
	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{2, 2}, []int{3, 3}) // 3x3 grid ⇒ 9 equations

	// essential boundary conditions
	e := NewBoundaryCondsGrid(g, 2)
	e.AddUsingTag(10, 0, 123.0, nil)                                                    // left ⇒ 0:ux
	e.AddUsingTag(20, 1, 0, func(x la.Vector, t float64) float64 { return 100 + x[0] }) // bottom ⇒ 1:uy

	// print bcs
	io.Pf("%v\n", e.Print())

	// check
	check3x3grid(tst, e, true)
}

func TestBryConds02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BryConds02. Using Mesh")

	// mesh             [21]
	//          6         7         8
	//          o---------o---------o
	//          |         |         |
	//          |    2    |    3    |
	//          |3        |4        |5
	//    [10]  o---------o---------o   [11]
	//          |         |         |
	//          |    0    |    1    |
	//          |0        |1        |2
	//          o---------o---------o
	//                  [20]
	dat := `
{
  "verts" : [
    { "i":0, "t":0, "x":[0, 0] },
    { "i":1, "t":0, "x":[1, 0] },
    { "i":2, "t":0, "x":[2, 0] },
    { "i":3, "t":0, "x":[0, 1] },
    { "i":4, "t":0, "x":[1, 1] },
    { "i":5, "t":0, "x":[2, 1] },
    { "i":6, "t":0, "x":[0, 2] },
    { "i":7, "t":0, "x":[1, 2] },
    { "i":8, "t":0, "x":[2, 2] }
  ],
  "cells" : [
    { "i":0, "t":1, "p":0, "y":"qua4", "v":[0,1,4,3], "et":[20,  0,  0, 10] },
    { "i":1, "t":1, "p":0, "y":"qua4", "v":[1,2,5,4], "et":[20, 11,  0,  0] },
    { "i":2, "t":1, "p":0, "y":"qua4", "v":[3,4,7,6], "et":[ 0,  0, 21, 10] },
    { "i":3, "t":1, "p":0, "y":"qua4", "v":[4,5,8,7], "et":[ 0, 11, 21,  0] }
  ]
}`

	// mesh
	m := new(msh.Mesh)
	err := json.Unmarshal([]byte(dat), m)
	if err != nil {
		chk.Panic("%v\n", err)
	}
	m.CheckAndCalcDerivedVars()

	// essential boundary conditions
	e := NewBoundaryCondsMesh(m, 2)
	e.AddUsingTag(10, 0, 0, func(x la.Vector, t float64) float64 { return 123 })        // left ⇒ 0:ux
	e.AddUsingTag(20, 1, 0, func(x la.Vector, t float64) float64 { return 100 + x[0] }) // bottom ⇒ 1:uy

	// print bcs
	io.Pf("%v\n", e.Print())

	// check
	check3x3grid(tst, e, false)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		args := msh.NewArgs()
		args.WithIdsVerts = true
		args.WithIdsCells = true
		args.WithTagsEdges = true
		m.Draw(args)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/pde", "bcs02")
	}
}

func TestBryConds03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BryConds03. panic in AddUsingTag. wrong dof")

	defer chk.RecoverTstPanicIsOK(tst)

	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{2, 2}, []int{3, 3}) // 3x3 grid ⇒ 9 equations
	e := NewBoundaryCondsGrid(g, 2)
	e.AddUsingTag(20, 2, 456.0, nil) // 2:wrong dof
}

func TestBryConds04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BryConds04. panic in AddUsingTag. wrong tag")

	defer chk.RecoverTstPanicIsOK(tst)

	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{2, 2}, []int{3, 3}) // 3x3 grid ⇒ 9 equations
	e := NewBoundaryCondsGrid(g, 2)
	e.AddUsingTag(200, 1, 456.0, nil) // 200:wrong tag
}

func TestBryConds05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("BryConds05. panic in AddUsingTag. no grid, no mesh")

	defer chk.RecoverTstPanicIsOK(tst)

	e := new(BoundaryConds)
	e.ndof = 1
	e.AddUsingTag(0, 0, 0, nil)
}
