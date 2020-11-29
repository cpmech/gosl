// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"github.com/cpmech/gosl/utl"
)

// EdgeKey holds 3 sorted numbers to identify an edge
type EdgeKey struct {
	A int // id of one vertex on edge
	B int // id of another vertex on edge
	C int // id of a third vertex on edge or the number of mesh vertices if edge has only 2 vertices
}

// Edge holds the vertices and cells attached to an edge
type Edge struct {
	Verts VertexSet       // vertices on edge
	Bdata BoundaryDataSet // cells attached to edge, including which local edge id of cell is attached
}

// EdgesMap is a map of edges
type EdgesMap map[EdgeKey]*Edge

// ExtractEdges find edges in mesh
func (o *Mesh) ExtractEdges() (edges EdgesMap) {

	// new map
	edges = make(map[EdgeKey]*Edge)

	// loop over cells
	var edgeKey EdgeKey
	for _, cell := range o.Cells {

		// loop over edges of cell
		for localEdgeID, localVids := range EdgeLocalVerts[cell.TypeIndex] {

			// set edge key as triple of vertices
			nVertsOnEdge := len(localVids)
			edgeKey.A = cell.V[localVids[0]]
			edgeKey.B = cell.V[localVids[1]]
			if nVertsOnEdge > 2 {
				edgeKey.C = cell.V[localVids[2]]
			} else {
				edgeKey.C = len(o.Verts) // indicator of not-available
			}
			utl.IntSort3(&edgeKey.A, &edgeKey.B, &edgeKey.C)

			// append this cell to list of shared cells of edge
			if edge, ok := edges[edgeKey]; ok {
				edge.Bdata = append(edge.Bdata, &BoundaryData{localEdgeID, cell})

				// new edge
			} else {
				edge = new(Edge)
				edge.Verts = make([]*Vertex, nVertsOnEdge)
				edge.Bdata = []*BoundaryData{{localEdgeID, cell}}
				for j, lvid := range localVids {
					edge.Verts[j] = o.Verts[cell.V[lvid]]
				}
				edges[edgeKey] = edge
			}
		}
	}
	return
}

// Split splits map into two sets: internal and boundary edges
// NOTE: boundary edge is determined by checking if edge is shared by only cell only
func (o *EdgesMap) Split() (internal, boundary EdgesMap) {
	internal = make(map[EdgeKey]*Edge)
	boundary = make(map[EdgeKey]*Edge)
	for ekey, edge := range *o {
		if len(edge.Bdata) == 1 {
			boundary[ekey] = edge
		} else {
			internal[ekey] = edge
		}
	}
	return
}
