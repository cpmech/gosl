// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"sort"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_msh01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("msh01")

	// load mesh
	m, err := Read("data/singleq4square1x1.msh")
	if err != nil {
		tst.Errorf("Read failed:\n%v\n", err)
		return
	}
	chk.IntAssert(len(m.Verts), 4)
	chk.IntAssert(len(m.Cells), 1)

	// check coordinates
	X := [][]float64{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}
	io.Pfyel("\nvertices:\n")
	for i, v := range m.Verts {
		io.Pf("%+v\n", v)
		chk.IntAssert(i, v.Id)
		chk.Vector(tst, io.Sf("vertex %2d: X", v.Id), 1e-15, v.X, X[v.Id])
		if v.Id == 0 || v.Id == 2 {
			chk.IntAssert(-1, v.Tag)
		}
	}

	// check connectivity
	V := [][]int{
		{0, 1, 2, 3},
	}
	etags := [][]int{
		{-10, -11, -12, -13},
	}
	io.Pfyel("\ncells:\n")
	for i, c := range m.Cells {
		io.Pf("%+v\n", c)
		chk.IntAssert(i, c.Id)
		chk.IntAssert(-1, c.Tag)
		chk.String(tst, "qua4", c.Type)
		chk.IntAssert(0, c.Part)
		chk.Ints(tst, io.Sf("cell %2d : V", c.Id), c.V, V[c.Id])
		chk.Ints(tst, io.Sf("cell %2d : edgetags", c.Id), c.EdgeTags, etags[c.Id])
	}

	// get map of tags
	d, err := m.GetTagMaps()
	if err != nil {
		tst.Errorf("GetTagMaps failed:\n%v\n", err)
		return
	}

	// check VertTag2verts
	chk.IntAssert(len(d.VertTag2verts), 1)
	var neg1verts []int
	io.Pfyel("\nVertTag2verts:\n")
	for key, val := range d.VertTag2verts {
		io.Pf("%v => %v\n", key, val)
		for _, v := range val {
			if v.Tag == -1 {
				neg1verts = append(neg1verts, v.Id)
			}
		}
	}
	sort.Ints(neg1verts)
	chk.Ints(tst, "-1 vertices", neg1verts, []int{0, 2})

	// check CellTag2cells
	chk.IntAssert(len(d.CellTag2cells), 1)
	var neg1cells []int
	io.Pfyel("\nCellTag2cells:\n")
	for key, val := range d.CellTag2cells {
		io.Pf("%v => %v\n", key, val)
		for _, c := range val {
			if c.Tag == -1 {
				neg1cells = append(neg1cells, c.Id)
			}
		}
	}
	sort.Ints(neg1cells)
	chk.Ints(tst, "-1 cells", neg1cells, []int{0})

	// check CellType2cells
	chk.IntAssert(len(d.CellType2cells), 1)
	var qua4cells []int
	io.Pfyel("\nCellType2cells:\n")
	for key, val := range d.CellType2cells {
		io.Pf("%v => %v\n", key, val)
		for _, c := range val {
			if c.Type == "qua4" {
				qua4cells = append(qua4cells, c.Id)
			}
		}
	}
	sort.Ints(qua4cells)
	chk.Ints(tst, "qua4 cells", qua4cells, []int{0})

	// check CellPart2cells
	chk.IntAssert(len(d.CellPart2cells), 1)
	var part0cells []int
	io.Pfyel("\nCellPart2cells:\n")
	for key, val := range d.CellPart2cells {
		io.Pf("%v => %v\n", key, val)
		for _, c := range val {
			if c.Part == 0 {
				part0cells = append(part0cells, c.Id)
			}
		}
	}
	sort.Ints(part0cells)
	chk.Ints(tst, "part 0 cells", part0cells, []int{0})

	// check EdgeTag2cells
	chk.IntAssert(len(d.EdgeTag2cells), 4)
	var neg10edgeLids []int
	var neg10edgeCells []int
	var neg11edgeLids []int
	var neg11edgeCells []int
	var neg12edgeLids []int
	var neg12edgeCells []int
	var neg13edgeLids []int
	var neg13edgeCells []int
	io.Pfyel("\nEdgeTag2cells:\n")
	for key, val := range d.EdgeTag2cells {
		io.Pf("%v => %v\n", key, val)
		for _, pair := range val {
			id := pair.BryId
			c := pair.C
			if c.EdgeTags[id] == -10 {
				neg10edgeLids = append(neg10edgeLids, id)
				neg10edgeCells = append(neg10edgeCells, c.Id)
			}
			if c.EdgeTags[id] == -11 {
				neg11edgeLids = append(neg11edgeLids, id)
				neg11edgeCells = append(neg11edgeCells, c.Id)
			}
			if c.EdgeTags[id] == -12 {
				neg12edgeLids = append(neg12edgeLids, id)
				neg12edgeCells = append(neg12edgeCells, c.Id)
			}
			if c.EdgeTags[id] == -13 {
				neg13edgeLids = append(neg13edgeLids, id)
				neg13edgeCells = append(neg13edgeCells, c.Id)
			}
		}
	}
	sort.Ints(neg10edgeLids)
	sort.Ints(neg10edgeCells)
	sort.Ints(neg11edgeLids)
	sort.Ints(neg11edgeCells)
	sort.Ints(neg12edgeLids)
	sort.Ints(neg12edgeCells)
	sort.Ints(neg13edgeLids)
	sort.Ints(neg13edgeCells)
	chk.Ints(tst, "-10 edge => lids ", neg10edgeLids, []int{0})
	chk.Ints(tst, "-10 edge => cells", neg10edgeCells, []int{0})
	chk.Ints(tst, "-11 edge => lids ", neg11edgeLids, []int{1})
	chk.Ints(tst, "-11 edge => cells", neg11edgeCells, []int{0})
	chk.Ints(tst, "-12 edge => lids ", neg12edgeLids, []int{2})
	chk.Ints(tst, "-12 edge => cells", neg12edgeCells, []int{0})
	chk.Ints(tst, "-13 edge => lids ", neg13edgeLids, []int{3})
	chk.Ints(tst, "-13 edge => cells", neg13edgeCells, []int{0})

	// check EdgeTag2verts
	chk.IntAssert(len(d.EdgeTag2verts), 4)
	var neg10edgeVerts []int
	var neg11edgeVerts []int
	var neg12edgeVerts []int
	var neg13edgeVerts []int
	io.Pfyel("\nEdgeTag2verts:\n")
	for key, val := range d.EdgeTag2verts {
		io.Pf("%v => %v\n", key, val)
		for _, v := range val {
			if key == -10 {
				neg10edgeVerts = append(neg10edgeVerts, v.Id)
			}
			if key == -11 {
				neg11edgeVerts = append(neg11edgeVerts, v.Id)
			}
			if key == -12 {
				neg12edgeVerts = append(neg12edgeVerts, v.Id)
			}
			if key == -13 {
				neg13edgeVerts = append(neg13edgeVerts, v.Id)
			}
		}
	}
	sort.Ints(neg10edgeVerts)
	sort.Ints(neg11edgeVerts)
	sort.Ints(neg12edgeVerts)
	sort.Ints(neg13edgeVerts)
	chk.Ints(tst, "-10 edge => verts", neg10edgeVerts, []int{0, 1})
	chk.Ints(tst, "-11 edge => verts", neg11edgeVerts, []int{1, 2})
	chk.Ints(tst, "-12 edge => verts", neg12edgeVerts, []int{2, 3})
	chk.Ints(tst, "-13 edge => verts", neg13edgeVerts, []int{0, 3})
}
