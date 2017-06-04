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
	WithVerts    bool           // draw vertices
	WithEdges    bool           // draw edges
	WithCells    bool           // draw cells
	WithFaces    bool           // draw faces
	WithIdsVerts bool           // with ids of vertices
	WithIdsCells bool           // with ids of cells
	WithIdsEdges bool           // with ids of edges
	WithIdsFaces bool           // with ids of faces
	ArgsVerts    *plt.A         // arguments for vertices
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
func (o *Mesh) Draw(args *DrawArgs) {

	// auxiliary
	type triple struct{ a, b, c int }   // points on edge
	edgesdrawn := make(map[triple]bool) // edges drawn already
	var tri triple

	// arguments
	if args == nil {
		args = NewArgs()
	}

	// loop over cells
	for _, cell := range o.Cells {

		// skip disabled cells
		if cell.Disabled {
			continue
		}

		// lin cell
		lincell := TypeIndexToKind[cell.TypeIndex] == KindLin
		if !lincell && args.OnlyLins {
			continue
		}

		// draw cells
		if args.WithCells {
			aa := getargs(cell.Id, args.ArgsCells)
			if cell.Gndim > 2 {
				// TODO
			} else {
				k := 0
				xx := make([][]float64, len(cell.V))
				for _, lvids := range EdgeLocalVertsD[cell.TypeIndex] { // loop over edges
					for ivert := 0; ivert < len(lvids)-1; ivert++ { // loop over verts on edge
						lv := lvids[ivert]
						xx[k] = o.Verts[cell.V[lv]].X
						k++
					}
				}
				plt.Polyline(xx, aa)
			}
		}

		// draw edges
		if args.WithEdges {
			for _, lvids := range EdgeLocalVertsD[cell.TypeIndex] {

				// set triple of vertices
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
					var z []float64
					ndim := len(o.Verts[tri.a].X)
					if ndim > 2 {
						z = make([]float64, nv)
					}
					for i, lv := range lvids {
						v := o.Verts[cell.V[lv]]
						x[i], y[i] = v.X[0], v.X[1]
						if ndim > 2 {
							z[i] = v.X[2]
						}
					}
					if ndim > 2 {
						plt.Plot3dLine(x, y, z, args.ArgsEdges)
					} else {
						plt.Plot(x, y, args.ArgsEdges)
					}
					edgesdrawn[tri] = true
				}
			}
		}

		// add middle node
		if args.WithVerts {
			if cell.TypeKey == "qua9" {
				vid := cell.V[8]
				x := o.Verts[vid].X[0]
				y := o.Verts[vid].X[1]
				plt.PlotOne(x, y, args.ArgsVerts)
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
			aa := getargs(cell.Id, args.ArgsLins)
			plt.Plot(x, y, aa)
		}

		// cell ids
		if args.WithIdsCells {
			xc := make([]float64, cell.Gndim)
			for _, x := range cell.X {
				for i := 0; i < cell.Gndim; i++ {
					xc[i] += x[i]
				}
			}
			for i := 0; i < cell.Gndim; i++ {
				xc[i] /= float64(len(cell.X))
			}
			txt := io.Sf("%d", cell.Id)
			if o.Ndim > 2 {
				z := 0.0
				if cell.Gndim > 2 {
					z = xc[2]
				}
				plt.Text3d(xc[0], xc[1], z, txt, args.ArgsIdsCells)
			} else {
				plt.Text(xc[0], xc[1], txt, args.ArgsIdsCells)
			}
		}
	}

	// loop over vertices
	if args.WithIdsVerts {
		for _, v := range o.Verts {
			if args.WithIdsVerts {
				txt := io.Sf("%d", v.Id)
				if o.Ndim > 2 {
					z := 0.0
					if len(v.X) > 2 {
						z = v.X[2]
					}
					plt.Text3d(v.X[0], v.X[1], z, txt, args.ArgsIdsVerts)
				} else {
					plt.Text(v.X[0], v.X[1], txt, args.ArgsIdsVerts)
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
