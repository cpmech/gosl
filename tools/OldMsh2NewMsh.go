// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"encoding/json"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/gm/msh"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// OldVert holds vertex data
type OldVert struct {
	ID  int       // id
	Tag int       // tag
	C   []float64 // coordinates (size==2 or 3)
}

// OldCell holds cell data
type OldCell struct {
	ID     int    // id
	Tag    int    // tag
	Geo    int    // geometry type (gemlab code)
	Type   string // geometry type (string)
	Part   int    // partition id
	Verts  []int  // vertices
	FTags  []int  // edge (2D) or face (3D) tags
	STags  []int  // seam tags (for 3D only; it is actually a 3D edge tag)
	JlinID int    // joint line id
	JsldID int    // joint solid id
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

	// verts: find largest strings
	var ndim int
	L := make([]int, 5)
	for _, v := range old.Verts {
		L[0] = utl.Imax(L[0], len(io.Sf("%d", v.ID)))
		L[1] = utl.Imax(L[1], len(io.Sf("%d", v.Tag)))
		for j, x := range v.C {
			L[2+j] = utl.Imax(L[2+j], len(io.Sf("%g", x)))
		}
		ndim = len(v.C)
	}
	S := make([]string, 5)
	for i, l := range L {
		S[i] = io.Sf("%d", l)
	}

	// write vertices
	buf := new(bytes.Buffer)
	io.Ff(buf, "{\n  \"verts\":[\n")
	for i, v := range old.Verts {
		if i > 0 {
			io.Ff(buf, ",\n")
		}
		io.Ff(buf, "    {\"i\":%"+S[0]+"d, \"t\":%"+S[1]+"d, \"x\":[", v.ID, v.Tag)
		for j, x := range v.C {
			if j > 0 {
				io.Ff(buf, ", ")
			}
			io.Ff(buf, "%"+S[2+j]+"g", x)
		}
		io.Ff(buf, "] }")
	}

	// cells: find largest strings
	n := 30
	L = make([]int, n*2)
	for _, c := range old.Cells {
		L[0] = utl.Imax(L[0], len(io.Sf("%d", c.ID)))
		L[1] = utl.Imax(L[1], len(io.Sf("%d", c.Tag)))
		L[2] = utl.Imax(L[2], len(io.Sf("%d", c.Part)))
		for j, v := range c.Verts {
			L[3+j] = utl.Imax(L[3+j], len(io.Sf("%d", v)))
		}
	}
	S = make([]string, n*2)
	for i, l := range L {
		S[i] = io.Sf("%d", l)
	}
	io.Ff(buf, "\n  ],")

	// write cells
	io.Ff(buf, "\n  \"cells\":[\n")
	for i, c := range old.Cells {
		if i > 0 {
			io.Ff(buf, ",\n")
		}
		io.Ff(buf, "    {\"i\":%"+S[0]+"d, \"t\":%"+S[1]+"d, \"p\":%"+S[2]+"d, \"y\":%q, \"v\":[", c.ID, c.Tag, c.Part, c.Type)
		for j, v := range c.Verts {
			if j > 0 {
				io.Ff(buf, ", ")
			}
			io.Ff(buf, "%"+S[3+j]+"d", v)
		}
		io.Ff(buf, "]")
		if len(c.FTags) > 0 {
			io.Ff(buf, ", ")
			if ndim == 2 {
				io.Ff(buf, "\"et\":[")
			} else {
				io.Ff(buf, "\"ft\":[")
			}
			for j, t := range c.FTags {
				if j > 0 {
					io.Ff(buf, ", ")
				}
				io.Ff(buf, "%d", t)
			}
			io.Ff(buf, "]")
		}
		io.Ff(buf, " }")
	}
	io.Ff(buf, "\n  ]\n}")
	io.WriteFileVD("/tmp/gosl", fnkey+"-new.msh", buf)

	// check
	m, err := msh.Read("/tmp/gosl/" + fnkey + "-new.msh")
	if err != nil {
		chk.Panic("cannot read new mesh:\n%v", err)
	}
	m.Check()
}
