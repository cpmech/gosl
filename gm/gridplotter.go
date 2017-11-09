// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/plt"
)

// GridPlotter assists in drawing grids
type GridPlotter struct {

	// input
	G         *CurvGrid // the grid
	OnlyBry   bool      // draw boundary only
	WithVids  bool      // with IDs of vertices
	WithVerts bool      // draw points
	WithBases bool      // with (covariant) basis vectors (g_0, g_1, g_2)
	Npts      []int     // number of points (to reduce the number of lines) [nil => all points]
	ArgsVids  *plt.A    // arguments for vertex ids [may be nil]
	ArgsVerts *plt.A    // arguments for vertices
	ArgsEdges *plt.A    // arguments for edges [may be nil]
	ArgsG0    *plt.A    // arguments for (covariant) basis vector g_0 [may be nil]
	ArgsG1    *plt.A    // arguments for (covariant) basis vector g_1 [may be nil]
	ArgsG2    *plt.A    // arguments for (covariant) basis vector g_2 [may be nil]

	// generated
	X2d, Y2d      [][]float64   // 2D meshgrid
	X3d, Y3d, Z3d [][][]float64 // 3D meshgrid
}

// Draw draws grid
func (o *GridPlotter) Draw() {

	// arguments
	if o.ArgsVids == nil {
		o.ArgsVids = &plt.A{C: plt.C(2, 0), Fsz: 7}
	}
	if o.ArgsVerts == nil {
		o.ArgsVerts = &plt.A{Fc: "none", Ec: "k", Lw: 0.8, NoClip: true}
	}
	if o.ArgsEdges == nil {
		o.ArgsEdges = &plt.A{C: "#427ce5", Lw: 0.8, NoClip: true}
	}

	// generate coordinates
	if o.G.ndim == 2 {
		o.X2d, o.Y2d = o.G.Meshgrid2d()
	} else {
		o.X3d, o.Y3d, o.Z3d = o.G.Meshgrid3d()
	}

	// draw edges
	if o.G.ndim == 2 {
		plt.Grid2d(o.X2d, o.Y2d, o.WithVerts, o.ArgsEdges, o.ArgsVerts)
	} else {
		plt.Grid3d(o.X3d, o.Y3d, o.Z3d, o.WithVids, o.ArgsEdges, o.ArgsVids)
	}
}
