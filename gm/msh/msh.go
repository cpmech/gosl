// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"encoding/json"
	"math"
	"path/filepath"
	"strings"

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

// Vert holds vertex data
type Vert struct {

	// input
	Id  int       // id
	Tag int       // tag
	C   []float64 // coordinates (size==2 or 3)

	// derived
	SharedBy []int // cells sharing this vertex
}

// Cell holds cell data
type Cell struct {

	// input
	Id       int    // id
	Tag      int    // tag
	Type     string // geometry type (string)
	Part     int    // partition id
	Verts    []int  // vertices
	FTags    []int  // edge (2D) or face (3D) tags
	STags    []int  // seam tags (for 3D only; it is actually a 3D edge tag)
	JlinId   int    // joint line id
	JsldId   int    // joint solid id
	Disabled bool   // cell is disabled

	// derived
	GoroutineId    int     // go routine id
	FaceLocalVerts [][]int // local ids of face vertices [nfaces][...]
	Neighbours     []*Cell // neighbour cells
}

// Mesh holds a mesh for FE analyses
type Mesh struct {

	// input
	Verts []*Vert // vertices
	Cells []*Cell // cells

	// derived
	FnamePath  string  // complete filename path
	Ndim       int     // space dimension
	Xmin, Xmax float64 // min and max x-coordinate
	Ymin, Ymax float64 // min and max x-coordinate
	Zmin, Zmax float64 // min and max x-coordinate

	// derived: maps
	VertTag2verts map[int][]*Vert      // vertex tag => set of vertices
	CellTag2cells map[int][]*Cell      // cell tag => set of cells
	FaceTag2cells map[int][]CellFaceId // face tag => set of cells
	FaceTag2verts map[int][]int        // face tag => vertices on tagged face
	SeamTag2cells map[int][]CellSeamId // seam tag => set of cells
	Ctype2cells   map[string][]*Cell   // cell type => set of cells
	Part2cells    map[int][]*Cell      // partition number => set of cells
}

// CellFaceId structure
type CellFaceId struct {
	C   *Cell // cell
	Fid int   // face id
}

// CellSeamId structure
type CellSeamId struct {
	C   *Cell // cell
	Sid int   // seam id
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
	o.Xmin = o.Verts[0].C[0]
	o.Ymin = o.Verts[0].C[1]
	if len(o.Verts[0].C) > 2 {
		o.Zmin = o.Verts[0].C[2]
	}
	o.Xmax = o.Xmin
	o.Ymax = o.Ymin
	o.Zmax = o.Zmin
	o.VertTag2verts = make(map[int][]*Vert)
	for i, v := range o.Verts {

		// check vertex id
		if v.Id != i {
			err = chk.Err("vertices ids must coincide with order in \"verts\" list. %d != %d\n", v.Id, i)
			return
		}

		// ndim
		nd := len(v.C)
		if nd < 2 || nd > 4 {
			err = chk.Err("number of space dimensions must be 2, 3 or 4 (NURBS). %d is invalid\n", nd)
			return
		}
		if nd == 3 {
			if math.Abs(v.C[2]) > TOL_ZMIN_FOR_3D {
				o.Ndim = 3
			}
		}

		// tags
		if v.Tag < 0 {
			verts := o.VertTag2verts[v.Tag]
			o.VertTag2verts[v.Tag] = append(verts, v)
		}

		// limits
		o.Xmin = utl.Min(o.Xmin, v.C[0])
		o.Xmax = utl.Max(o.Xmax, v.C[0])
		o.Ymin = utl.Min(o.Ymin, v.C[1])
		o.Ymax = utl.Max(o.Ymax, v.C[1])
		if nd > 2 {
			o.Zmin = utl.Min(o.Zmin, v.C[2])
			o.Zmax = utl.Max(o.Zmax, v.C[2])
		}
	}

	// derived: maps
	o.CellTag2cells = make(map[int][]*Cell)
	o.FaceTag2cells = make(map[int][]CellFaceId)
	o.FaceTag2verts = make(map[int][]int)
	o.SeamTag2cells = make(map[int][]CellSeamId)
	o.Ctype2cells = make(map[string][]*Cell)
	o.Part2cells = make(map[int][]*Cell)
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
		c.GoroutineId = goroutineId

		// face local vertices
		var ok bool
		c.FaceLocalVerts, ok = FaceLocalVerts[c.Type]
		if !ok {
			err = chk.Err("cannot handle type %q\n", c.Type)
			return
		}

		// face tags
		cells := o.CellTag2cells[c.Tag]
		o.CellTag2cells[c.Tag] = append(cells, c)
		for i, ftag := range c.FTags {
			if ftag < 0 {
				pairs := o.FaceTag2cells[ftag]
				o.FaceTag2cells[ftag] = append(pairs, CellFaceId{c, i})
				for _, l := range c.FaceLocalVerts[i] {
					utl.IntIntsMapAppend(&o.FaceTag2verts, ftag, o.Verts[c.Verts[l]].Id)
				}
			}
		}

		// seam tags
		if o.Ndim == 3 {
			for i, stag := range c.STags {
				if stag < 0 {
					pairs := o.SeamTag2cells[stag]
					o.SeamTag2cells[stag] = append(pairs, CellSeamId{c, i})
				}
			}
		}

		// cell type => cells
		cells = o.Ctype2cells[c.Type]
		o.Ctype2cells[c.Type] = append(cells, c)

		// partition => cells
		cells = o.Part2cells[c.Part]
		o.Part2cells[c.Part] = append(cells, c)

		// set SharedBy information on vertices
		for _, vid := range c.Verts {
			if utl.IntIndexSmall(o.Verts[vid].SharedBy, c.Id) < 0 {
				o.Verts[vid].SharedBy = append(o.Verts[vid].SharedBy, c.Id)
			}
		}
	}

	// remove duplicates
	for ftag, verts := range o.FaceTag2verts {
		o.FaceTag2verts[ftag] = utl.IntUnique(verts)
	}
	return
}

// Draw2d draws 2D mesh
//  lwds -- linewidths: maps cid => lwd. Use <nil> for default lwd
//  ms   -- markersize for nodes
func (o *Mesh) Draw2d(onlyLin, setup bool, lwds map[int]float64, ms int) {

	// auxiliary
	type triple struct{ a, b, c int }   // points on edge
	edgesdrawn := make(map[triple]bool) // edges drawn already
	var tri triple

	// loop over cells
	for _, cell := range o.Cells {

		// skip disabled cells
		if cell.Disabled {
			continue
		}

		// lin cell
		lincell := strings.HasPrefix(cell.Type, "lin")
		if onlyLin && !lincell {
			continue
		}

		// loop edges of cells
		for _, lvids := range cell.FaceLocalVerts {

			// set triple of nodes
			tri.a = cell.Verts[lvids[0]]
			tri.b = cell.Verts[lvids[1]]
			nv := len(lvids)
			if nv > 2 {
				tri.c = cell.Verts[lvids[2]]
			} else {
				tri.c = len(o.Verts) + 1 // indicator of not-available
			}
			utl.IntSort3(&tri.a, &tri.b, &tri.c)

			// draw edge if not drawn yet
			if _, drawn := edgesdrawn[tri]; !drawn {
				x := make([]float64, nv)
				y := make([]float64, nv)
				x[0] = o.Verts[tri.a].C[0]
				y[0] = o.Verts[tri.a].C[1]
				if nv == 3 {
					x[1] = o.Verts[tri.c].C[0]
					y[1] = o.Verts[tri.c].C[1]
					x[2] = o.Verts[tri.b].C[0]
					y[2] = o.Verts[tri.b].C[1]
				} else {
					x[1] = o.Verts[tri.b].C[0]
					y[1] = o.Verts[tri.b].C[1]
				}
				plt.Plot(x, y, io.Sf("'k-o', ms=%d, clip_on=0", ms))
				edgesdrawn[tri] = true
			}
		}

		// add middle node
		if cell.Type == "qua9" {
			vid := cell.Verts[8]
			x := o.Verts[vid].C[0]
			y := o.Verts[vid].C[1]
			plt.PlotOne(x, y, io.Sf("'ko', ms=%d, clip_on=0", ms))
		}

		// linear cells
		if lincell {
			nv := len(cell.Verts)
			x := make([]float64, nv)
			y := make([]float64, nv)
			for i, vid := range cell.Verts {
				x[i] = o.Verts[vid].C[0]
				y[i] = o.Verts[vid].C[1]
			}
			lw := 2.0
			if lwd, ok := lwds[cell.Id]; ok {
				lw = lwd
			}
			plt.Plot(x, y, io.Sf("'-o', ms=%d, clip_on=0, color='#41045a', lw=%g", ms, lw))
		}
	}

	// set up
	if setup {
		plt.Equal()
		plt.AxisRange(o.Xmin, o.Xmax, o.Ymin, o.Ymax)
		plt.AxisOff()
	}
}
