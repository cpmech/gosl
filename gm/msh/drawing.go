// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// DrawArgs holds drawing arguments
type DrawArgs struct {
	OnlyLins     bool           // only 'lin' cells
	WithNodes    bool           // draw nodes
	WithEdges    bool           // draw edges
	WithCells    bool           // draw cells
	WithFaces    bool           // draw faces
	WithIdsVerts bool           // with ids of nodes
	WithIdsCells bool           // with ids of cells
	WithIdsEdges bool           // with ids of edges
	WithIdsFaces bool           // with ids of faces
	ArgsVerts    *plt.A         // arguments for nodes
	ArgsEdges    *plt.A         // arguments for edges
	ArgsCells    map[int]*plt.A // arguments for cells [cellId] => A; if len==1, use the same for all
	ArgsLins     map[int]*plt.A // arguments for lins [cellId] => A; if len==1, use the same for all
	ArgsIdsCells *plt.A         // arguments for the ids of cells
	ArgsIdsVerts *plt.A         // arguments for the ids of vertices
}

// NewArgs returns a new set of drawing arguments
func NewArgs() (o *DrawArgs) {
	o = new(DrawArgs)
	o.Default()
	return
}

// Default sets default argument values
func (o *DrawArgs) Default() {
	o.WithCells = true
	o.WithEdges = false
	o.ArgsVerts = &plt.A{C: "k", NoClip: true}
	o.ArgsEdges = &plt.A{C: "#480085", NoClip: true}
	o.ArgsCells = map[int]*plt.A{-1: &plt.A{Fc: "#dce1f4", Ec: "k", Closed: true, NoClip: true}}
	o.ArgsLins = map[int]*plt.A{-1: &plt.A{C: "#41045a", NoClip: true}}
	o.ArgsIdsCells = &plt.A{C: "k", Fsz: 7, Ha: "center", Va: "center", NoClip: true}
	o.ArgsIdsVerts = &plt.A{C: "r", Fsz: 7, Ha: "left", Va: "bottom", NoClip: true}
}

// Draw draws mesh. Arguments A may be nil (defaults will be selected)
func (o *Mesh) Draw(a *DrawArgs) {

	// auxiliary
	type triple struct{ a, b, c int }   // points on edge
	edgesdrawn := make(map[triple]bool) // edges drawn already
	var tri triple

	// arguments
	if a == nil {
		a = NewArgs()
	}

	// loop over cells
	for _, cell := range o.Cells {

		// skip disabled cells
		if cell.Disabled {
			continue
		}

		// lin cell
		lincell := TypeIndexToKind[cell.TypeIndex] == KindLin
		if !lincell && a.OnlyLins {
			continue
		}

		// draw cells
		var X [][]float64 // all coordinates
		if a.WithCells {
			aa := getargs(cell.Id, a.ArgsCells)
			X = o.ExtractCellCoords(cell.Id)
			if cell.Gndim > 2 {
				// TODO
			} else {
				plt.Polyline(X, aa)
			}
		}

		// draw edges
		if a.WithEdges {
			elv := EdgeLocalVerts[cell.TypeIndex]
			for _, lvids := range elv {

				// set triple of nodes
				tri.a = cell.V[lvids[0]]
				tri.b = cell.V[lvids[1]]
				nv := len(lvids)
				if nv > 2 {
					tri.c = cell.V[lvids[2]]
				} else {
					tri.c = len(o.Verts) // indicator of not-available
				}
				utl.IntSort3(&tri.a, &tri.b, &tri.c)

				// draw edge if not drawn yet
				if _, drawn := edgesdrawn[tri]; !drawn {
					x := make([]float64, nv)
					y := make([]float64, nv)
					x[0] = o.Verts[tri.a].X[0]
					y[0] = o.Verts[tri.a].X[1]
					var z []float64
					ndim := len(o.Verts[tri.a].X)
					if ndim > 2 {
						z = make([]float64, nv)
						z[0] = o.Verts[tri.a].X[2]
					}
					if nv > 2 {
						x[1] = o.Verts[tri.c].X[0]
						y[1] = o.Verts[tri.c].X[1]
						x[2] = o.Verts[tri.b].X[0]
						y[2] = o.Verts[tri.b].X[1]
						if ndim > 2 {
							z[1] = o.Verts[tri.c].X[2]
							z[2] = o.Verts[tri.b].X[2]
						}
					} else {
						x[1] = o.Verts[tri.b].X[0]
						y[1] = o.Verts[tri.b].X[1]
						if ndim > 2 {
							z[1] = o.Verts[tri.b].X[2]
						}
					}
					if ndim > 2 {
						plt.Plot3dLine(x, y, z, a.ArgsEdges)
					} else {
						plt.Plot(x, y, a.ArgsEdges)
					}
					edgesdrawn[tri] = true
				}
			}
		}

		// add middle node
		if a.WithNodes {
			if cell.TypeKey == "qua9" {
				vid := cell.V[8]
				x := o.Verts[vid].X[0]
				y := o.Verts[vid].X[1]
				plt.PlotOne(x, y, a.ArgsVerts)
			}
		}

		// linear cells
		if lincell {

			// collect vertices
			nv := len(cell.V)
			x := make([]float64, nv)
			y := make([]float64, nv)
			for i, vid := range cell.V {
				x[i] = o.Verts[vid].X[0]
				y[i] = o.Verts[vid].X[1]
			}
			aa := getargs(cell.Id, a.ArgsLins)
			plt.Plot(x, y, aa)
		}

		// cell ids
		if a.WithIdsCells {
			if X == nil {
				X = o.ExtractCellCoords(cell.Id)
			}
			xc := make([]float64, cell.Gndim)
			for _, x := range X {
				for i := 0; i < cell.Gndim; i++ {
					xc[i] += x[i]
				}
			}
			for i := 0; i < cell.Gndim; i++ {
				xc[i] /= float64(len(X))
			}
			txt := io.Sf("%d", cell.Id)
			if o.Ndim > 2 {
				z := 0.0
				if cell.Gndim > 2 {
					z = xc[2]
				}
				plt.Text3d(xc[0], xc[1], z, txt, a.ArgsIdsCells)
			} else {
				plt.Text(xc[0], xc[1], txt, a.ArgsIdsCells)
			}
		}
	}

	// loop over vertices
	if a.WithIdsVerts {
		for _, v := range o.Verts {
			if a.WithIdsVerts {
				txt := io.Sf("%d", v.Id)
				if o.Ndim > 2 {
					z := 0.0
					if len(v.X) > 2 {
						z = v.X[2]
					}
					plt.Text3d(v.X[0], v.X[1], z, txt, a.ArgsIdsVerts)
				} else {
					plt.Text(v.X[0], v.X[1], txt, a.ArgsIdsVerts)
				}
			}
		}
	}

	// set figure
	if o.Ndim > 2 {
		plt.Default3dView(o.Xmin[0], o.Xmax[0], o.Xmin[1], o.Xmax[1], o.Xmin[2], o.Xmax[2], true)
	} else {
		plt.Equal()
		plt.AxisRange(o.Xmin[0], o.Xmax[0], o.Xmin[1], o.Xmax[1])
	}
}

// getargs returns arguments
func getargs(id int, argsMap map[int]*plt.A) *plt.A {
	args := argsMap[id]
	if args == nil {
		return argsMap[-1]
	}
	return args
}
