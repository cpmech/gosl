// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package msh defines mesh data structures and implements interpolation functions for finite
// element analyses (FEA)
package msh

import (
	"encoding/json"
	"math"
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// Vertex holds vertex data (e.g. from msh file)
type Vertex struct {

	// input
	ID  int       `json:"i"` // identifier
	Tag int       `json:"t"` // tag
	X   []float64 `json:"x"` // coordinates (size==2 or 3)

	// auxiliary
	Entity interface{} `json:"-"` // any entity attached to this vertex
}

// VertexSet defines a set of vertices
type VertexSet []*Vertex

// Len returns the length of vertex set
func (o VertexSet) Len() int { return len(o) }

// Swap swaps two entries in vertex set
func (o VertexSet) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

// Less compares ides in vertex set
func (o VertexSet) Less(i, j int) bool { return o[i].ID < o[j].ID }

// IDs returns the IDs of vertices in VertexSet
func (o VertexSet) IDs() (ids []int) {
	ids = make([]int, len(o))
	for i, vertex := range o {
		ids[i] = vertex.ID
	}
	return
}

// Cell holds cell data (in e.g. from msh file)
type Cell struct {

	// input
	ID       int    `json:"i"`  // identifier
	Tag      int    `json:"t"`  // tag
	Part     int    `json:"p"`  // partition id
	Disabled bool   `json:"d"`  // cell is disabled
	TypeKey  string `json:"y"`  // geometry type; e.g. "lin2"
	V        []int  `json:"v"`  // vertices
	EdgeTags []int  `json:"et"` // edge tags (2D or 3D)
	FaceTags []int  `json:"ft"` // face tags (3D only)
	NurbsID  int    `json:"b"`  // id of NURBS (or something else) that this cell belongs to
	Span     []int  `json:"s"`  // span in NURBS

	// derived
	TypeIndex int        `json:"-"` // type index of cell. converted from TypeKey
	Gndim     int        `json:"-"` // geometry ndim
	X         *la.Matrix `json:"-"` // all vertex coordinates [nverts][ndim]
}

// CellSet defines a set of cells
type CellSet []*Cell

// Mesh defines mesh data
type Mesh struct {

	// input
	Verts VertexSet `json:"verts"` // vertices
	Cells CellSet   `json:"cells"` // cells

	// derived
	Ndim int       // max space dimension among all vertices
	Xmin []float64 // min(x) among all vertices [ndim]
	Xmax []float64 // max(x) among all vertices [ndim]

	// auxiliary
	Tmaps *TagMaps // map of tags
}

// BoundaryData holds ID of edge or face and pointer to Cell at boundary (edge or face)
type BoundaryData struct {
	LocalID int   // edge local id (edgeId) OR face local id (faceId)
	Cell    *Cell // cell
}

// BoundaryDataSet defines a set of BoundaryData
type BoundaryDataSet []*BoundaryData

// TagMaps holds data for finding information using on tags
type TagMaps struct {
	VertexTag2verts map[int]VertexSet       // vertex tag => set of vertices
	CellTag2cells   map[int]CellSet         // cell tag => set of cells
	CellType2cells  map[int]CellSet         // cell type => set of cells
	CellPart2cells  map[int]CellSet         // partition number => set of cells
	EdgeTag2cells   map[int]BoundaryDataSet // edge tag => set of cells {cell,boundaryId}
	EdgeTag2verts   map[int]VertexSet       // edge tag => vertices on tagged edge [unique]
	FaceTag2cells   map[int]BoundaryDataSet // face tag => set of cells {cell,boundaryId}
	FaceTag2verts   map[int]VertexSet       // face tag => vertices on tagged edge [unique]
}

// NewMesh creates mesh from json string
func NewMesh(jsonString string) (o *Mesh) {
	o = new(Mesh)
	err := json.Unmarshal([]byte(jsonString), o)
	if err != nil {
		chk.Panic("%v\n", err)
	}
	o.CheckAndCalcDerivedVars()
	return
}

// Read reads mesh and call CheckAndCalcDerivedVars
func Read(fn string) (o *Mesh) {
	o = new(Mesh)
	b := io.ReadFile(fn)
	err := json.Unmarshal(b, o)
	if err != nil {
		chk.Panic("%v\n", err)
	}
	o.CheckAndCalcDerivedVars()
	return
}

// CheckAndCalcDerivedVars checks input data and computes derived quantities such as the max space
// dimension, min(x) and max(x) among all vertices, cells' gndim, etc.
// This function will set o.Ndim, o.Xmin and o.Xmax.
// This function will also generate the maps of tags.
func (o *Mesh) CheckAndCalcDerivedVars() {

	// check for at least one vertex
	if len(o.Verts) < 1 {
		chk.Panic("at least 1 vertex is required in mesh\n")
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
		if id != vert.ID {
			chk.Panic("vertex ids must be sequential. vertex %d must be %d\n", vert.ID, id)
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
		if id != cell.ID {
			chk.Panic("cell ids must be sequential. cell %d must be %d\n", cell.ID, id)
		}
		tindex, ok := TypeKeyToIndex[cell.TypeKey]
		if !ok {
			chk.Panic("cannot find cell type key %q in database\n", cell.TypeKey)
		}
		cell.TypeIndex = tindex
		cell.Gndim = GeomNdim[cell.TypeIndex]
		nv := NumVerts[cell.TypeIndex]
		if len(cell.V) != nv {
			chk.Panic("number of vertices for cell %d is incorrect. %d != %d\n", cell.ID, len(cell.V), nv)
		}
		nEtags := len(cell.EdgeTags)
		if nEtags > 0 {
			lv := EdgeLocalVerts[cell.TypeIndex]
			if nEtags != len(lv) {
				chk.Panic("number of edge tags for cell %d is incorrect. %d != %d\n", cell.ID, nEtags, len(lv))
			}
		}
		cell.X = o.ExtractCellCoords(cell.ID)
	}

	// new tag maps
	o.Tmaps = new(TagMaps)
	o.Tmaps.VertexTag2verts = make(map[int]VertexSet)
	o.Tmaps.CellTag2cells = make(map[int]CellSet)
	o.Tmaps.CellType2cells = make(map[int]CellSet)
	o.Tmaps.CellPart2cells = make(map[int]CellSet)
	o.Tmaps.EdgeTag2cells = make(map[int]BoundaryDataSet)
	o.Tmaps.EdgeTag2verts = make(map[int]VertexSet)
	o.Tmaps.FaceTag2cells = make(map[int]BoundaryDataSet)
	o.Tmaps.FaceTag2verts = make(map[int]VertexSet)

	// loop over vertices
	for _, vert := range o.Verts {
		if vert.Tag != 0 {
			o.Tmaps.VertexTag2verts[vert.Tag] = append(o.Tmaps.VertexTag2verts[vert.Tag], vert)
		}
	}

	// loop over cells
	for _, cell := range o.Cells {

		// basic data
		edgeLocVerts := EdgeLocalVerts[cell.TypeIndex]
		faceLocVerts := FaceLocalVerts[cell.TypeIndex]

		// cell data
		o.Tmaps.CellTag2cells[cell.Tag] = append(o.Tmaps.CellTag2cells[cell.Tag], cell)
		o.Tmaps.CellType2cells[cell.TypeIndex] = append(o.Tmaps.CellType2cells[cell.TypeIndex], cell)
		o.Tmaps.CellPart2cells[cell.Part] = append(o.Tmaps.CellPart2cells[cell.Part], cell)

		// check edge tags
		if len(cell.EdgeTags) > 0 {
			if len(cell.EdgeTags) != len(edgeLocVerts) {
				chk.Panic("number of edge tags in \"et\" list for cell # %d is incorrect. %d != %d\n", cell.ID, len(cell.EdgeTags), len(edgeLocVerts))
			}
		}

		// check face tags
		if len(cell.FaceTags) > 0 {
			if len(cell.FaceTags) != len(faceLocVerts) {
				chk.Panic("number of face tags in \"ft\" list for cell # %d is incorrect. %d != %d\n", cell.ID, len(cell.FaceTags), len(faceLocVerts))
			}
		}

		// edge tags => cells, verts
		o.setBryTagMaps(&o.Tmaps.EdgeTag2cells, &o.Tmaps.EdgeTag2verts, cell, cell.EdgeTags, edgeLocVerts)

		// face tags => cells, verts
		if len(faceLocVerts) > 0 {
			o.setBryTagMaps(&o.Tmaps.FaceTag2cells, &o.Tmaps.FaceTag2verts, cell, cell.FaceTags, faceLocVerts)
		}
	}

	// sort entries in EdgeTag2verts
	for edgeTag, vertsOnEdge := range o.Tmaps.EdgeTag2verts {
		sort.Sort(vertsOnEdge)
		o.Tmaps.EdgeTag2verts[edgeTag] = vertsOnEdge
	}

	// sort entries in FaceTag2verts
	for faceTag, vertsOnFace := range o.Tmaps.FaceTag2verts {
		sort.Sort(vertsOnFace)
		o.Tmaps.FaceTag2verts[faceTag] = vertsOnFace
	}
	return
}

// ExtractCellCoords extracts cell coordinates
//   X -- matrix with coordinates [nverts][gndim]
func (o *Mesh) ExtractCellCoords(cellID int) (X *la.Matrix) {
	c := o.Cells[cellID]
	X = la.NewMatrix(len(c.V), c.Gndim)
	for m, v := range c.V {
		for i := 0; i < c.Gndim; i++ {
			X.Set(m, i, o.Verts[v].X[i])
		}
	}
	return
}

// Boundary returns a list of indices of nodes on edge (2D) or face (3D) of boundary
//   NOTE: will return empty list if tag is not available
func (o *Mesh) Boundary(tag int) []int {
	if o.Ndim == 2 {
		if vset, ok := o.Tmaps.EdgeTag2verts[tag]; ok {
			return vset.IDs()
		}
		return nil
	}
	if vset, ok := o.Tmaps.FaceTag2verts[tag]; ok {
		return vset.IDs()
	}
	return nil
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// setBryTagMaps sets maps of boundary tags
func (o *Mesh) setBryTagMaps(cellBryMap *map[int]BoundaryDataSet, vertBryMap *map[int]VertexSet, cell *Cell, tagList []int, locVerts [][]int) {

	// loop over each tag attached to a side of the cell
	for localID, tag := range tagList {

		// there is a tag (i.e. it's nonzero)
		if tag != 0 {

			// set edgeTag => cells map
			(*cellBryMap)[tag] = append((*cellBryMap)[tag], &BoundaryData{localID, cell})

			// loop over local edges of cell
			for _, locVid := range locVerts[localID] {

				// find vertex
				vid := cell.V[locVid] // local vertex id => global vertex id (vid)
				vert := o.Verts[vid]  // pointer to vertex

				// find whether this edgeTag is present in the map or not
				if vertsOnEdge, ok := (*vertBryMap)[tag]; ok {

					// find whether this vertex is in the slice attached to edgeTag or not
					found := false
					for _, v := range vertsOnEdge {
						if vert.ID == v.ID {
							found = true
							break
						}
					}

					// add vertex to (unique) slice attached to edgeTag
					if !found {
						(*vertBryMap)[tag] = append(vertsOnEdge, vert)
					}

					// edgeTag is not in the map => create new slice with the first vertex in it
				} else {
					(*vertBryMap)[tag] = []*Vertex{vert}
				}
			}
		}
	}
}
