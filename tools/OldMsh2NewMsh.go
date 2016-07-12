// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/gm/msh"
	"github.com/cpmech/gosl/io"
)

// OldVert holds vertex data
type OldVert struct {
	Id  int       // id
	Tag int       // tag
	C   []float64 // coordinates (size==2 or 3)
}

// OldCell holds cell data
type OldCell struct {
	Id     int    // id
	Tag    int    // tag
	Geo    int    // geometry type (gemlab code)
	Type   string // geometry type (string)
	Part   int    // partition id
	Verts  []int  // vertices
	FTags  []int  // edge (2D) or face (3D) tags
	STags  []int  // seam tags (for 3D only; it is actually a 3D edge tag)
	JlinId int    // joint line id
	JsldId int    // joint solid id
}

// OldMesh holds a mesh for FE analyses
type OldMesh struct {
	Verts []*OldVert // vertices
	Cells []*OldCell // cells
}

func main() {

	// catch errors
	defer func() {
		if err := recover(); err != nil {
			io.PfRed("ERROR: %v\n", err)
		}
	}()

	// input data
	mshfn, fnkey := io.ArgToFilename(0, "data/sgm57", ".msh", true)

	// old mesh
	var old OldMesh

	// read file
	b, err := io.ReadFile(mshfn)
	if err != nil {
		chk.Panic("%v", err)
	}

	// decode
	err = json.Unmarshal(b, &old)
	if err != nil {
		chk.Panic("%v", err)
	}

	// new mesh
	var m msh.Mesh
	//var buf bytes.Buffer
	nv := len(old.Verts)
	nc := len(old.Cells)
	m.Verts = make([]*msh.Vertex, nv)
	m.Cells = make([]*msh.Cell, nc)
	ndim := 2
	for i, v := range old.Verts {
		m.Verts[i] = new(msh.Vertex)
		m.Verts[i].Id = v.Id
		m.Verts[i].Tag = v.Tag
		m.Verts[i].X = v.C
		if len(v.C) == 3 {
			ndim = 3
		}
	}
	for i, c := range old.Cells {
		m.Cells[i] = new(msh.Cell)
		m.Cells[i].Id = c.Id
		m.Cells[i].Tag = c.Tag
		m.Cells[i].Part = c.Part
		m.Cells[i].Disabled = false
		m.Cells[i].Type = c.Type
		m.Cells[i].V = c.Verts
		if ndim == 2 {
			m.Cells[i].EdgeTags = c.FTags
		} else {
			m.Cells[i].FaceTags = c.FTags
		}
	}
	io.Pforan("m= %v\n", m)

	// encode
	res, err := json.Marshal(&m)
	if err != nil {
		chk.Panic("%v", err)
	}
	buf := bytes.NewBuffer(res)
	io.WriteFileD("/tmp/gosl", fnkey+"-new.msh", buf)
}
