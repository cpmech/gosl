// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"gosl/plt"
)

// GridPlotter assists in drawing grids
type GridPlotter struct {

	// input
	G         *Grid  // the grid
	OnlyBry   bool   // draw boundary only
	WithVids  bool   // with IDs of vertices
	WithBases bool   // with (covariant) basis vectors (g_0, g_1, g_2)
	Npts      []int  // number of points (to reduce the number of lines) [nil => all points]
	ArgsVids  *plt.A // arguments for vertex ids [may be nil]
	ArgsEdges *plt.A // arguments for edges [may be nil]
	ArgsG0    *plt.A // arguments for (covariant) basis vector g_0 [may be nil]
	ArgsG1    *plt.A // arguments for (covariant) basis vector g_1 [may be nil]
	ArgsG2    *plt.A // arguments for (covariant) basis vector g_2 [may be nil]

	// generated
	X2d, Y2d      [][]float64   // 2D meshgrid
	X3d, Y3d, Z3d [][][]float64 // 3D meshgrid
}

// Draw draws grid
func (o *GridPlotter) Draw() {

	// arguments
	if o.ArgsVids == nil {
		o.ArgsVids = &plt.A{C: "k", Fsz: 6}
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
		plt.Grid2d(o.X2d, o.Y2d, o.WithVids, o.ArgsEdges, o.ArgsVids)
	} else {
		plt.Grid3d(o.X3d, o.Y3d, o.Z3d, o.WithVids, o.ArgsEdges, o.ArgsVids)
	}
}

// Bases draw basis vectors
func (o *GridPlotter) Bases(scale float64) {
	if o.ArgsG0 == nil {
		o.ArgsG0 = &plt.A{C: "r", Scale: 7, Z: 10}
	}
	if o.ArgsG1 == nil {
		o.ArgsG1 = &plt.A{C: "g", Scale: 7, Z: 10}
	}
	if o.ArgsG2 == nil {
		o.ArgsG2 = &plt.A{C: "b", Scale: 7, Z: 10}
	}
	if o.G.ndim == 2 {
		p := 0
		for n := 0; n < o.G.npts[1]; n++ {
			for m := 0; m < o.G.npts[0]; m++ {
				M := o.G.mtr[p][n][m]
				plt.DrawArrow2d(M.X, M.CovG0, true, scale, o.ArgsG0)
				plt.DrawArrow2d(M.X, M.CovG1, true, scale, o.ArgsG1)
			}
		}
		return
	}
	for p := 0; p < o.G.npts[2]; p++ {
		for n := 0; n < o.G.npts[1]; n++ {
			for m := 0; m < o.G.npts[0]; m++ {
				M := o.G.mtr[p][n][m]
				plt.DrawArrow3d(M.X, M.CovG0, true, scale, o.ArgsG0)
				plt.DrawArrow3d(M.X, M.CovG1, true, scale, o.ArgsG1)
				plt.DrawArrow3d(M.X, M.CovG2, true, scale, o.ArgsG2)
			}
		}
	}
}
