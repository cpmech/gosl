// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package msh implements functions to generate meshes
package msh

import (
	"encoding/json"
	"math"
	"path/filepath"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// constants
const (
	TOL_ZMIN_FOR_3D      = 1e-7
	TOL_COINCIDENT_VERTS = 1e-5
)

// Vertex holds vertex data (in .msh file)
type Vertex struct {

	// input
	Id  int       // identifier
	Tag int       // tag
	X   []float64 // coordinates (size==2 or 3)

	// auxiliary
	Entity interface{} // any entity attached to this vertex

	// derived
	C []int // cells sharing this vertex
}

// Cell holds cell data (in .msh file)
type Cell struct {

	// input
	Id       int    // identifier
	Tag      int    // tag
	Type     string // geometry type
	Part     int    // partition id
	V        []int  // vertices
	EdgeTags []int  // edge tags (2D or 3D)
	FaceTags []int  // face tags (3D only)
	Disabled bool   // cell is disabled

	// auxiliary
	Entity interface{} // any entity attached to this vertex

	// derived
	GoroutineId    int     // go routine id
	EdgeLocalVerts [][]int // local ids of vertices on edges [nedges][...]
	Neighbours     []*Cell // neighbour cells
}

// Mesh holds mesh data (in .msh file)
type Mesh struct {

	// input
	Verts []*Vertex // vertices
	Cells []*Cell   // cells

	// derived
	FnamePath  string  // complete filename path
	Ndim       int     // space dimension
	Xmin, Xmax float64 // min and max x-coordinate
	Ymin, Ymax float64 // min and max y-coordinate
	Zmin, Zmax float64 // min and max z-coordinate

	// derived: maps
	Tag2verts     map[int][]*Vertex    // vertex tag => set of vertices
	Tag2cells     map[int][]*Cell      // cell tag => set of cells
	Type2cells    map[string][]*Cell   // cell type => set of cells
	Part2cells    map[int][]*Cell      // partition number => set of cells
	EdgeTag2cells map[int][]CellIdPair // edge tag => set of cells
	EdgeTag2verts map[int][]int        // edge tag => vertices on tagged edge

	// derived: sets
	Edges EdgeSet // all edges
}

// Three holds 3 indices (for defining edges)
type Three struct{ A, B, C int }

// Four holds 4 indices (for defining faces)
type Four struct{ A, B, C, D int }

// Edge defines an edge
type Edge struct {
	V []int // vertices on edge
	C []int // cells connected to this edge
}

// Face defines a face
type Face struct {
	V []int // vertices on face
	C []int // cells connected to this face
}

// EdgeSet defines a set of Edges
type EdgeSet map[Three]*Edge

// FaceSet defines a set of Faces
type FaceSet map[Four]*Face

// CellIdPair structure
type CellIdPair struct {
	C  *Cell
	Id int
}

// ReadMsh reads mesh
func ReadMsh(dir, fn string, goroutineId int) (o *Mesh, err error) {

	// new mesh
	o = new(Mesh)

	// read file
	o.FnamePath = filepath.Join(dir, fn)
	b, err := io.ReadFile(o.FnamePath)
	if err != nil {
		return
	}

	// decode
	err = json.Unmarshal(b, &o)
	if err != nil {
		return
	}

	// compute derived quantities
	err = o.CalcDerived(goroutineId)
	return
}

// CalcDerived computes derived quantities
func (o *Mesh) CalcDerived(goroutineId int) (err error) {

	// check
	if len(o.Verts) < 2 {
		err = chk.Err("at least 2 vertices are required in mesh\n")
		return
	}
	if len(o.Cells) < 1 {
		err = chk.Err("at least 1 cell is required in mesh\n")
		return
	}

	// vertex related derived data
	o.Ndim = 2
	o.Xmin = o.Verts[0].X[0]
	o.Ymin = o.Verts[0].X[1]
	if len(o.Verts[0].X) > 2 {
		o.Zmin = o.Verts[0].X[2]
	}
	o.Xmax = o.Xmin
	o.Ymax = o.Ymin
	o.Zmax = o.Zmin
	o.Tag2verts = make(map[int][]*Vertex)
	for i, v := range o.Verts {

		// check vertex id
		if v.Id != i {
			err = chk.Err("vertices ids must coincide with order in \"verts\" list. %d != %d\n", v.Id, i)
			return
		}

		// ndim
		nd := len(v.X)
		if nd < 2 || nd > 4 {
			err = chk.Err("number of space dimensions must be 2, 3 or 4 (NURBS). %d is invalid\n", nd)
			return
		}
		if nd == 3 {
			if math.Abs(v.X[2]) > TOL_ZMIN_FOR_3D {
				o.Ndim = 3
			}
		}

		// tags
		if v.Tag < 0 {
			o.Tag2verts[v.Tag] = append(o.Tag2verts[v.Tag], v)
		}

		// limits
		o.Xmin = utl.Min(o.Xmin, v.X[0])
		o.Xmax = utl.Max(o.Xmax, v.X[0])
		o.Ymin = utl.Min(o.Ymin, v.X[1])
		o.Ymax = utl.Max(o.Ymax, v.X[1])
		if nd > 2 {
			o.Zmin = utl.Min(o.Zmin, v.X[2])
			o.Zmax = utl.Max(o.Zmax, v.X[2])
		}
	}

	// derived: maps and sets
	o.Tag2cells = make(map[int][]*Cell)
	o.Type2cells = make(map[string][]*Cell)
	o.Part2cells = make(map[int][]*Cell)
	o.EdgeTag2cells = make(map[int][]CellIdPair)
	o.Edges = make(map[Three]*Edge)
	for i, c := range o.Cells {

		// check id and tag
		if c.Id != i {
			err = chk.Err("cells ids must coincide with order in \"verts\" list. %d != %d\n", c.Id, i)
			return
		}
		if c.Tag >= 0 {
			err = chk.Err("cells tags must be negative. %d is incorrect\n", c.Tag)
			return
		}

		// cell data
		c.GoroutineId = goroutineId
		o.Tag2cells[c.Tag] = append(o.Tag2cells[c.Tag], c)
		o.Type2cells[c.Type] = append(o.Type2cells[c.Type], c)
		o.Part2cells[c.Part] = append(o.Part2cells[c.Part], c)

		// edge local vertices
		var ok bool
		c.EdgeLocalVerts, ok = EdgeLocalVerts[c.Type]
		if !ok {
			err = chk.Err("cannot handle type %q\n", c.Type)
			return
		}

		// edge tags
		for i, tag := range c.EdgeTags {
			if tag < 0 {
				o.EdgeTag2cells[tag] = append(o.EdgeTag2cells[tag], CellIdPair{c, i})
				for _, l := range c.EdgeLocalVerts[i] {
					utl.IntIntsMapAppend(&o.EdgeTag2verts, tag, o.Verts[c.V[l]].Id)
				}
			}
		}

		// face tags
		//if o.Ndim == 3 {
		//}

		// set cells sharing this vertex
		for _, vid := range c.V {
			if utl.IntIndexSmall(o.Verts[vid].C, c.Id) < 0 {
				o.Verts[vid].C = append(o.Verts[vid].C, c.Id)
			}
		}

		// edges
		for _, L := range c.EdgeLocalVerts {
			key := Three{-1, c.V[L[0]], c.V[L[1]]}
			if len(L) > 2 {
				key.A = c.V[L[2]]
			}
			utl.IntSort3(&key.A, &key.B, &key.C)
			if e, ok := o.Edges[key]; ok {
				if utl.IntIndexSmall(e.C, c.Id) < 0 {
					e.C = append(e.C, c.Id)
				}
			} else {
				verts := []int{c.V[L[0]], c.V[L[1]]}
				if len(L) > 2 {
					verts = append(verts, c.V[L[2]])
				}
				o.Edges[key] = &Edge{verts, []int{c.Id}}
			}
		}
	}

	// remove duplicates
	for tag, verts := range o.EdgeTag2verts {
		o.EdgeTag2verts[tag] = utl.IntUnique(verts)
	}
	return
}

// Draw2d draws 2D mesh
//  lwds -- linewidths: maps cid => lwd. Use <nil> for default lwd
func (o *Mesh) Draw2d(vids, setup bool, lwds map[int]float64, fmt *plt.Fmt) {

	// format
	if fmt == nil {
		fmt = &plt.Fmt{C: "k", M: ".", Ms: 1}
	}
	args := fmt.GetArgs("") + ",clip_on=0"

	// loop over edges
	for _, e := range o.Edges {
		X := make([]float64, len(e.V))
		Y := make([]float64, len(e.V))
		X[0] = o.Verts[e.V[0]].X[0]
		Y[0] = o.Verts[e.V[0]].X[1]
		for j, v := range e.V[1:] {
			X[1+j] = o.Verts[v].X[0]
			Y[1+j] = o.Verts[v].X[1]
		}
		plt.Plot(X, Y, args)
	}

	// vertex ids
	if vids {
		for _, v := range o.Verts {
			plt.Text(v.X[0], v.X[1], io.Sf("%d", v.Id), "")
		}
	}

	// set up
	if setup {
		plt.Equal()
		plt.AxisRange(o.Xmin, o.Xmax, o.Ymin, o.Ymax)
		plt.AxisOff()
	}
}
