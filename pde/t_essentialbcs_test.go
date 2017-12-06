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

func TestEssentialBcs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("EssentialBcs01. Using Grid")

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
	e.AddUsingTag(10, 0, 123.0, nil) // left ⇒ 0:ux
	e.AddUsingTag(20, 1, 456.0, nil) // bottom ⇒ 1:uy

	// print bcs
	io.Pf("%v\n", e.Print())

	// check
	nodes := e.Nodes()
	chk.Ints(tst, "nodes", nodes, []int{0, 1, 2, 3, 6})
	vals0 := make([]float64, 9)
	vals1 := make([]float64, 9)
	var available0, available1 bool
	for n := 0; n < 9; n++ {
		vals0[n], available0 = e.Value(n, 0, 0)
		vals1[n], available1 = e.Value(n, 1, 0)
		if !available0 {
			vals0[n] = -1
		}
		if !available1 {
			vals1[n] = -1
		}
	}
	chk.Array(tst, "vals0", 1e-14, vals0, []float64{123, -1, -1, 123, -1, -1, 123, -1, -1})
	chk.Array(tst, "vals1", 1e-14, vals1, []float64{456, 456, 456, -1, -1, -1, -1, -1, -1})
}

func TestEssentialBcs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("EssentialBcs02. Using Mesh")

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
	nodes := e.Nodes()
	chk.Ints(tst, "nodes", nodes, []int{0, 1, 2, 3, 6})
	vals0 := make([]float64, 9)
	vals1 := make([]float64, 9)
	var available0, available1 bool
	for n := 0; n < 9; n++ {
		vals0[n], available0 = e.Value(n, 0, 0)
		vals1[n], available1 = e.Value(n, 1, 0)
		if !available0 {
			vals0[n] = -1
		}
		if !available1 {
			vals1[n] = -1
		}
	}
	chk.Array(tst, "vals0", 1e-14, vals0, []float64{123, -1, -1, 123, -1, -1, 123, -1, -1})
	chk.Array(tst, "vals1", 1e-14, vals1, []float64{100, 101, 102, -1, -1, -1, -1, -1, -1})

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

func TestEssentialBcs03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("EssentialBcs03. panic in AddUsingTag. wrong dof")

	defer chk.RecoverTstPanicIsOK(tst)

	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{2, 2}, []int{3, 3}) // 3x3 grid ⇒ 9 equations
	e := NewBoundaryCondsGrid(g, 2)
	e.AddUsingTag(20, 2, 456.0, nil) // 2:wrong dof
}

func TestEssentialBcs04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("EssentialBcs04. panic in AddUsingTag. wrong tag")

	defer chk.RecoverTstPanicIsOK(tst)

	g := new(gm.Grid)
	g.RectGenUniform([]float64{0, 0}, []float64{2, 2}, []int{3, 3}) // 3x3 grid ⇒ 9 equations
	e := NewBoundaryCondsGrid(g, 2)
	e.AddUsingTag(200, 1, 456.0, nil) // 200:wrong tag
}

func TestEssentialBcs05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("EssentialBcs05. panic in AddUsingTag. no grid, no mesh")

	defer chk.RecoverTstPanicIsOK(tst)

	e := new(BoundaryConds)
	e.maxNdof = 1
	e.AddUsingTag(0, 0, 0, nil)
}
