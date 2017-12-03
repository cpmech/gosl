// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tri

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

// DrawVC draws mesh with given vertices in V (coordinates) and cells in C (connectivity)
func DrawVC(V [][]float64, C [][]int, style *plt.A) {
	if style == nil {
		style = &plt.A{C: "b", M: "o", Ms: 2}
	}
	type edgeType struct{ A, B int }
	drawnEdges := make(map[edgeType]bool)
	for _, cell := range C {
		for i := 0; i < 3; i++ {
			a, b := cell[i], cell[(i+1)%3]
			edge := edgeType{a, b}
			if b < a {
				edge.A, edge.B = edge.B, edge.A
			}
			if _, found := drawnEdges[edge]; !found {
				x := []float64{V[a][0], V[b][0]}
				y := []float64{V[a][1], V[b][1]}
				plt.Plot(x, y, style)
				drawnEdges[edge] = true
			}
		}
	}
}

// DrawMesh draws mesh
func DrawMesh(m *Mesh, withIds bool, args, argsVids, argsCids, argsEtags *plt.A) {
	if args == nil {
		args = &plt.A{C: plt.C(0, 0), M: "o", Ms: 2, Lw: 0.7, NoClip: true}
	}
	if argsVids == nil {
		argsVids = &plt.A{C: plt.C(2, 0), Fsz: 6, NoClip: true}
	}
	if argsCids == nil {
		argsCids = &plt.A{C: plt.C(1, 0), Fsz: 6, NoClip: true}
	}
	if argsEtags == nil {
		argsEtags = &plt.A{C: plt.C(4, 0), Fsz: 6, NoClip: true, Ha: "right", Va: "top"}
	}
	type edgeType struct{ A, B int }
	drawnEdges := make(map[edgeType]bool)
	for _, cell := range m.Cells {
		var xm, ym float64
		for i := 0; i < 3; i++ { // loop over sides
			a, b := cell.V[i], cell.V[(i+1)%3]
			edge := edgeType{a, b}
			if b < a {
				edge.A, edge.B = edge.B, edge.A
			}
			if _, found := drawnEdges[edge]; !found {
				x := []float64{m.Verts[a].X[0], m.Verts[b].X[0]}
				y := []float64{m.Verts[a].X[1], m.Verts[b].X[1]}
				plt.Plot(x, y, args)
				drawnEdges[edge] = true
			}
			exm := (m.Verts[a].X[0] + m.Verts[b].X[0]) / 2.0
			eym := (m.Verts[a].X[1] + m.Verts[b].X[1]) / 2.0
			plt.Text(exm, eym, io.Sf("%d", cell.EdgeTags[i]), argsEtags)
			xm += m.Verts[cell.V[i]].X[0]
			ym += m.Verts[cell.V[i]].X[1]
		}
		xm /= 3.0
		ym /= 3.0
		plt.Text(xm, ym, io.Sf("%d(%d)", cell.ID, cell.Tag), argsCids)
		if len(cell.V) > 3 { // middle vertices
			plt.PlotOne(m.Verts[cell.V[3]].X[0], m.Verts[cell.V[3]].X[1], args)
			plt.PlotOne(m.Verts[cell.V[4]].X[0], m.Verts[cell.V[4]].X[1], args)
			plt.PlotOne(m.Verts[cell.V[5]].X[0], m.Verts[cell.V[5]].X[1], args)
		}
	}
	if withIds {
		for _, vert := range m.Verts {
			plt.Text(vert.X[0], vert.X[1], io.Sf("%d(%d)", vert.ID, vert.Tag), argsVids)
		}
	}
}
