// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package msh defines mesh data structures and implements interpolation functions for finite
// element analyses (FEA)
package msh

import (
	"encoding/json"
	"math"
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// EdgeKey implements a key to identify edges
type EdgeKey struct {
	NumVerts int // number of vertices, from A to C
	A, B, C  int // vertices
}

// FaceKey implements a key to identify faces
type FaceKey struct {
	NumVerts   int // number of vertices, from A to D
	A, B, C, D int // vertices
}

// VertSet defines a set of vertices
type VertSet []*Vertex

// EdgeSet defines a set of edges
type EdgeSet []*Edge

// FaceSet defines a set of faces
type FaceSet []*Face

// CellSet defines a set of cells
type CellSet []*Cell

// EdgeSetMap defines a map of Edges
type EdgesMap map[EdgeKey]*Edge

// FaceSet defines a map of Faces
type FacesMap map[FaceKey]*Face

// Vertex holds vertex data (e.g. from msh file)
type Vertex struct {

	// input
	Id  int       `json:"i"` // identifier
	Tag int       `json:"t"` // tag
	X   []float64 `json:"x"` // coordinates (size==2 or 3)

	// auxiliary
	Entity interface{} `json:"-"` // any entity attached to this vertex

	// derived
	SharedByCells CellSet `json:"-"` // cells sharing this vertex
	Neighbours    VertSet `json:"-"` // neighbour vertices
}

// Cell holds cell data (in .msh file)
type Cell struct {

	// input
	Id       int    `json:"i"`  // identifier
	Tag      int    `json:"t"`  // tag
	Part     int    `json:"p"`  // partition id
	Disabled bool   `json:"d"`  // cell is disabled
	TypeKey  string `json:"y"`  // geometry type; e.g. "lin2"
	V        []int  `json:"v"`  // vertices
	EdgeTags []int  `json:"et"` // edge tags (2D or 3D)
	FaceTags []int  `json:"ft"` // face tags (3D only)
	NurbsId  int    `json:"b"`  // id of NURBS (or something else) that this cell belongs to
	Span     []int  `json:"s"`  // span in NURBS

	// auxiliary
	TypeIndex int `json:"-"` // type index of cell. converted from TypeKey

	// derived
	Edges      EdgeSet     `json:"-"` // edges on this cell
	Faces      FaceSet     `json:"-"` // faces on this cell
	Neighbours CellSet     `json:"-"` // neighbour cells
	Gndim      int         `json:"-"` // geometry ndim
	X          [][]float64 `json:"-"` // all vertex coordinates
}

// Mesh defines mesh data
type Mesh struct {

	// input
	Verts VertSet `json:"verts"` // vertices
	Cells CellSet `json:"cells"` // cells

	// derived
	EdgesMap EdgesMap `json:"-"` // all edges
	FacesMap FacesMap `json:"-"` // all faces

	// calculated by CalcDerived
	Ndim int       // max space dimension among all vertices
	Xmin []float64 // min(x) among all vertices [ndim]
	Xmax []float64 // max(x) among all vertices [ndim]
}

// Edge defines an edge
type Edge struct {
	V VertSet // vertices on edge
	C CellSet // cells connected to this edge
}

// Face defines a face
type Face struct {
	V VertSet // vertices on face
	C CellSet // cells connected to this face
}

// BryPair defines a structure to identify bryIds => cells pairs
type BryPair struct {
	C     *Cell // cell
	BryId int   // edge local id (edgeId) OR face local id (faceId)
}

// BryPairSet defines a set of BryPair identifiers
type BryPairSet []*BryPair

// TagMaps holds data for finding information based on tags
type TagMaps struct {
	VertTag2verts  map[int]VertSet    // vertex tag => set of vertices
	CellTag2cells  map[int]CellSet    // cell tag => set of cells
	CellType2cells map[int]CellSet    // cell type => set of cells
	CellPart2cells map[int]CellSet    // partition number => set of cells
	EdgeTag2cells  map[int]BryPairSet // edge tag => set of cells {cell,boundaryId}
	EdgeTag2verts  map[int]VertSet    // edge tag => vertices on tagged edge [unique]
	FaceTag2cells  map[int]BryPairSet // face tag => set of cells {cell,boundaryId}
	FaceTag2verts  map[int]VertSet    // face tag => vertices on tagged edge [unique]
}

// Read reads mesh and call CheckAndDerivedVars
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

	// check
	err = o.CheckAndCalcDerivedVars()
	return
}

// CheckAndCalcDerivedVars checks input data and computes derived quantities such as the max space
// dimension, min(x) and max(x) among all vertices, cells' gndim, etc.
// This function will set o.Ndim, o.Xmin and o.Xmax
func (o *Mesh) CheckAndCalcDerivedVars() (err error) {

	// check for at least one vertex
	if len(o.Verts) < 1 {
		err = chk.Err("at least 1 vertex is required in mesh\n")
		return
	}

	// check vertex data and find max(ndim), Xmin, and Xmax
	o.Xmin = make([]float64, 4)
	o.Xmax = make([]float64, 4)
	for i := 0; i < 4; i++ {
		o.Xmin[i] = math.MaxFloat64
		o.Xmax[i] = math.SmallestNonzeroFloat64
	}
	o.Ndim = len(o.Verts[0].X)
	for id, vert := range o.Verts {
		if id != vert.Id {
			err = chk.Err("vertex ids must be sequential. vertex %d must be %d", vert.Id, id)
			return
		}
		ndim := len(vert.X)
		if ndim > o.Ndim {
			o.Ndim = ndim
		}
		for i := 0; i < ndim; i++ {
			if vert.X[i] < o.Xmin[i] {
				o.Xmin[i] = vert.X[i]
			}
			if vert.X[i] > o.Xmax[i] {
				o.Xmax[i] = vert.X[i]
			}
		}
	}
	o.Xmin = o.Xmin[0:o.Ndim] // re-slice
	o.Xmax = o.Xmax[0:o.Ndim] // re-slice

	// check cell data, set TypeIndex, gndim, and coordinates X
	for id, cell := range o.Cells {
		if id != cell.Id {
			err = chk.Err("cell ids must be sequential. cell %d must be %d", cell.Id, id)
			return
		}
		if tindex, ok := TypeKeyToIndex[cell.TypeKey]; !ok {
			err = chk.Err("cannot find cell type ken %q in database\n", cell.TypeKey)
			return
		} else {
			cell.TypeIndex = tindex
		}
		cell.Gndim = GeomNdim[cell.TypeIndex]
		nv := NumVerts[cell.TypeIndex]
		if len(cell.V) != nv {
			err = chk.Err("number of vertices for cell %d is incorrect. %d != %d", cell.Id, len(cell.V), nv)
			return
		}
		nEtags := len(cell.EdgeTags)
		if nEtags > 0 {
			lv := EdgeLocalVerts[cell.TypeIndex]
			if nEtags != len(lv) {
				err = chk.Err("number of edge tags for cell %d is incorrect. %d != %d", cell.Id, nEtags, len(lv))
				return
			}
		}
		cell.X = o.ExtractCellCoords(cell.Id)
	}
	return
}

// GetTagMaps finds tagged entities
func (o *Mesh) GetTagMaps() (m *TagMaps, err error) {

	// new tag maps
	m = new(TagMaps)
	m.VertTag2verts = make(map[int]VertSet)
	m.CellTag2cells = make(map[int]CellSet)
	m.CellType2cells = make(map[int]CellSet)
	m.CellPart2cells = make(map[int]CellSet)
	m.EdgeTag2cells = make(map[int]BryPairSet)
	m.EdgeTag2verts = make(map[int]VertSet)
	m.FaceTag2cells = make(map[int]BryPairSet)
	m.FaceTag2verts = make(map[int]VertSet)

	// loop over vertices
	for _, vert := range o.Verts {
		if vert.Tag < 0 {
			m.VertTag2verts[vert.Tag] = append(m.VertTag2verts[vert.Tag], vert)
		}
	}

	// loop over cells
	for _, cell := range o.Cells {

		// basic data
		edgeLocVerts := EdgeLocalVerts[cell.TypeIndex]
		faceLocVerts := FaceLocalVerts[cell.TypeIndex]

		// cell data
		m.CellTag2cells[cell.Tag] = append(m.CellTag2cells[cell.Tag], cell)
		m.CellType2cells[cell.TypeIndex] = append(m.CellType2cells[cell.TypeIndex], cell)
		m.CellPart2cells[cell.Part] = append(m.CellPart2cells[cell.Part], cell)

		// check edge tags
		if len(cell.EdgeTags) > 0 {
			if len(cell.EdgeTags) != len(edgeLocVerts) {
				err = chk.Err("number of edge tags in \"et\" list for cell # %d is incorrect. %d != %d", cell.Id, len(cell.EdgeTags), len(edgeLocVerts))
				return
			}
		}

		// check face tags
		if len(cell.FaceTags) > 0 {
			if len(cell.FaceTags) != len(faceLocVerts) {
				err = chk.Err("number of face tags in \"ft\" list for cell # %d is incorrect. %d != %d", cell.Id, len(cell.FaceTags), len(faceLocVerts))
				return
			}
		}

		// edge tags => cells, verts
		o.setBryTagMaps(&m.EdgeTag2cells, &m.EdgeTag2verts, cell, cell.EdgeTags, edgeLocVerts)

		// face tags => cells, verts
		if len(faceLocVerts) > 0 {
			o.setBryTagMaps(&m.FaceTag2cells, &m.FaceTag2verts, cell, cell.FaceTags, faceLocVerts)
		}
	}

	// sort entries in EdgeTag2verts
	for edgeTag, vertsOnEdge := range m.EdgeTag2verts {
		sort.Sort(vertsOnEdge)
		m.EdgeTag2verts[edgeTag] = vertsOnEdge
	}

	// sort entries in FaceTag2verts
	for faceTag, vertsOnFace := range m.FaceTag2verts {
		sort.Sort(vertsOnFace)
		m.FaceTag2verts[faceTag] = vertsOnFace
	}
	return
}

// methods ////////////////////////////////////////////////////////////////////////////////////////

// ExtractCellCoords extracts cell coordinates
//   X -- matrix with coordinates [nverts][gndim]
func (o *Mesh) ExtractCellCoords(cellId int) (X [][]float64) {
	c := o.Cells[cellId]
	X = make([][]float64, len(c.V))
	for m, v := range c.V {
		X[m] = make([]float64, c.Gndim)
		for i := 0; i < c.Gndim; i++ {
			X[m][i] = o.Verts[v].X[i]
		}
	}
	return
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// setBryTagMaps sets maps of boundary tags
func (o *Mesh) setBryTagMaps(cellBryMap *map[int]BryPairSet, vertBryMap *map[int]VertSet, cell *Cell, tagList []int, locVerts [][]int) {

	// loop over each tag attached to a side of the cell
	for edgeId, edgeTag := range tagList {

		// there is a tag (i.e. it's negative)
		if edgeTag < 0 {

			// set edgeTag => cells map
			(*cellBryMap)[edgeTag] = append((*cellBryMap)[edgeTag], &BryPair{cell, edgeId})

			// loop over local edges of cell
			for _, locVid := range locVerts[edgeId] {

				// find vertex
				vid := cell.V[locVid] // local vertex id => global vertex id (vid)
				vert := o.Verts[vid]  // pointer to vertex

				// find whether this edgeTag is present in the map or not
				if vertsOnEdge, ok := (*vertBryMap)[edgeTag]; ok {

					// find whether this vertex is in the slice attached to edgeTag or not
					found := false
					for _, v := range vertsOnEdge {
						if vert.Id == v.Id {
							found = true
							break
						}
					}

					// add vertex to (unique) slice attached to edgeTag
					if !found {
						(*vertBryMap)[edgeTag] = append(vertsOnEdge, vert)
					}

					// edgeTag is not in the map => create new slice with the first vertex in it
				} else {
					(*vertBryMap)[edgeTag] = []*Vertex{vert}
				}
			}
		}
	}
}
