// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package msh implements functions to generate meshes
package msh

import (
	"encoding/json"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// ThreeIds holds 3 indices (e.g. for defining edges)
type ThreeIds struct{ A, B, C int }

// FourIds holds 4 indices (e.g. for defining faces)
type FourIds struct{ A, B, C, D int }

// Edge defines an edge
type Edge struct {
	V []*Vertex // vertices on edge
	C []*Cell   // cells connected to this edge
}

// Face defines a face
type Face struct {
	V []*Vertex // vertices on face
	C []*Cell   // cells connected to this face
}

// EdgeSet defines a set of Edges
type EdgeSet map[ThreeIds]*Edge

// FaceSet defines a set of Faces
type FaceSet map[FourIds]*Face

// Vertex holds vertex data (in .msh file)
type Vertex struct {

	// input
	Id  int       // identifier
	Tag int       // tag
	X   []float64 // coordinates (size==2 or 3)

	// auxiliary
	Entity interface{} // any entity attached to this vertex

	// derived
	SharedByCells []*Cell   // cells sharing this vertex
	Neighbours    []*Vertex // neighbour vertices
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
	Edges      []*Edge // edges on this cell
	Faces      []*Face // faces on this cell
	Neighbours []*Cell // neighbour cells
}

// Mesh defines mesh data
type Mesh struct {

	// input
	Verts []*Vertex // vertices
	Cells []*Cell   // cells

	// derived
	Edges EdgeSet // all edges
	Faces FaceSet // all faces
}

// CellBryIdPair structure
type CellBryIdPair struct {
	C     *Cell // cell
	BryId int   // edge local id (edgeId) OR face local id (faceId)
}

// TagMaps holds data for finding information based on tags
type TagMaps struct {
	VertTag2verts  map[int][]*Vertex       // vertex tag => set of vertices
	CellTag2cells  map[int][]*Cell         // cell tag => set of cells
	CellType2cells map[string][]*Cell      // cell type => set of cells
	CellPart2cells map[int][]*Cell         // partition number => set of cells
	EdgeTag2cells  map[int][]CellBryIdPair // edge tag => set of cells {cell,boundaryId}
	EdgeTag2verts  map[int][]*Vertex       // edge tag => vertices on tagged edge
}

// Read reads mesh
func Read(fn string) (o *Mesh, err error) {

	// new mesh
	o = new(Mesh)

	// read file
	b, err := io.ReadFile(fn)
	if err != nil {
		return
	}

	// decode
	err = json.Unmarshal(b, &o)
	if err != nil {
		return
	}
	return
}

// CalcLimits calculates limits and space dimension
func (o *Mesh) CalcLimits() (Ndim int, Min []float64, Max []float64, err error) {

	// check
	if len(o.Verts) < 1 {
		err = chk.Err("at least 1 vertex is required in mesh\n")
		return
	}

	// allocate slices
	Ndim = len(o.Verts[0].X)
	for i := 0; i < Ndim; i++ {
		Min = append(Min, o.Verts[0].X[i])
		Max = append(Max, o.Verts[0].X[i])
	}

	// loop over vertices
	for _, vert := range o.Verts {
		for i := 0; i < Ndim; i++ {
			if vert.X[i] < Min[i] {
				Min[i] = vert.X[i]
			}
			if vert.X[i] > Max[i] {
				Max[i] = vert.X[i]
			}
		}
	}
	return
}

// GetTagMaps finds tagged entities
func (o *Mesh) GetTagMaps() (m *TagMaps, err error) {

	// new tag maps
	m = new(TagMaps)
	m.VertTag2verts = make(map[int][]*Vertex)
	m.CellTag2cells = make(map[int][]*Cell)
	m.CellType2cells = make(map[string][]*Cell)
	m.CellPart2cells = make(map[int][]*Cell)
	m.EdgeTag2cells = make(map[int][]CellBryIdPair)
	m.EdgeTag2verts = make(map[int][]*Vertex)

	// loop over vertices
	for _, vert := range o.Verts {
		if vert.Tag < 0 {
			m.VertTag2verts[vert.Tag] = append(m.VertTag2verts[vert.Tag], vert)
		}
	}

	// loop over cells
	for _, cell := range o.Cells {

		// basic data
		var ok bool
		var geomNdim int
		if geomNdim, ok = GeomNdim[cell.Type]; !ok {
			err = chk.Err("cell type %q is not available in factory of shapes (in GeomNdim)")
			return
		}
		var edgeLocVerts [][]int
		if edgeLocVerts, ok = EdgeLocalVerts[cell.Type]; !ok {
			err = chk.Err("cell type %q is not available in factory of shapes (in EdgeLocalVerts)")
			return
		}
		var faceLocVerts [][]int
		if geomNdim == 3 {
			if faceLocVerts, ok = FaceLocalVerts[cell.Type]; !ok {
				err = chk.Err("cell type %q is not available in factory of shapes (in FaceLocalVerts)")
				return
			}
		}
		if cell.Tag >= 0 {
			err = chk.Err("cells tags must be negative. %d is incorrect\n", cell.Tag)
			return
		}

		// cell data
		m.CellTag2cells[cell.Tag] = append(m.CellTag2cells[cell.Tag], cell)
		m.CellType2cells[cell.Type] = append(m.CellType2cells[cell.Type], cell)
		m.CellPart2cells[cell.Part] = append(m.CellPart2cells[cell.Part], cell)

		// edge tags
		if len(cell.EdgeTags) > 0 {
			if len(cell.EdgeTags) != len(edgeLocVerts) {
				err = chk.Err("number of tags in \"edgetags\" list for cell # %d is incorrect", cell.Id)
				return
			}
		}
		for edgeId, edgeTag := range cell.EdgeTags {
			if edgeTag < 0 {
				m.EdgeTag2cells[edgeTag] = append(m.EdgeTag2cells[edgeTag], CellBryIdPair{cell, edgeId})
				for _, locVid := range edgeLocVerts[edgeId] {
					vid := cell.V[locVid]
					vert := o.Verts[vid]
					if vertsOnEdge, ok := m.EdgeTag2verts[edgeTag]; ok {
						m.EdgeTag2verts[edgeTag] = append(vertsOnEdge, vert)
					} else {
						m.EdgeTag2verts[edgeTag] = []*Vertex{vert}
					}
				}
			}
		}

		// face tags
		//for faceId, faceTag := range cell.FaceTags {
		//if faceTag < 0 {
		//io.Pforan("locVid=%d faceTag=%d\n", locVid, faceTag)
		//_ = faceLocVerts
		//}
		//}
		_ = faceLocVerts
		/*
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
		*/
	}

	/*
		// remove duplicates
		for tag, verts := range o.EdgeTag2verts {
			o.EdgeTag2verts[tag] = utl.IntUnique(verts)
		}
		return
	*/
	return
}

/*
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
*/
