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
	OnlyLins      bool           // only 'lin' cells
	WithVerts     bool           // draw vertices
	WithEdges     bool           // draw edges
	WithCells     bool           // draw cells
	WithFaces     bool           // draw faces
	WithIdsVerts  bool           // with ids of vertices
	WithIdsCells  bool           // with ids of cells
	WithIdsEdges  bool           // with ids of edges
	WithIdsFaces  bool           // with ids of faces
	WithTagsVerts bool           // with tags of vertices
	WithTagsEdges bool           // with tags of edges
	ArgsVerts     *plt.A         // arguments for vertices
	ArgsEdges     *plt.A         // arguments for edges
	ArgsCells     map[int]*plt.A // arguments for cells [cellId] => A; if len==1, use the same for all
	ArgsLins      map[int]*plt.A // arguments for lins [cellId] => A; if len==1, use the same for all
	ArgsIdsCells  *plt.A         // arguments for the ids of cells
	ArgsIdsVerts  *plt.A         // arguments for the ids of vertices
	ArgsTagsVerts *plt.A         // arguments for the tags of vertices
	ArgsTagsEdges *plt.A         // arguments for the tags of edges
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
	o.ArgsVerts = &plt.A{C: "k", M: ".", Ms: 3, Ls: "none", NoClip: true}
	o.ArgsEdges = &plt.A{C: "#480085", NoClip: true}
	o.ArgsCells = map[int]*plt.A{-1: &plt.A{Fc: "#dce1f4", Ec: "k", Closed: true, NoClip: true}}
	o.ArgsLins = map[int]*plt.A{-1: &plt.A{C: "#41045a", NoClip: true}}
	o.ArgsIdsCells = &plt.A{C: "k", Fsz: 6, Ha: "center", Va: "center", NoClip: true}
	o.ArgsIdsVerts = &plt.A{C: "r", Fsz: 6, Ha: "left", Va: "bottom", NoClip: true}
	o.ArgsTagsVerts = &plt.A{C: "g", Fsz: 6, Ha: "right", Va: "bottom", NoClip: true}
	o.ArgsTagsEdges = &plt.A{C: "m", Fsz: 6, Ha: "center", Va: "center", NoClip: true}
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
			ct := cell.TypeIndex
			nEdges := len(EdgeLocalVertsD[ct])
			nvEdge := len(EdgeLocalVertsD[ct][0]) // number of vertices along edge of cell
			nvEtot := nvEdge*nEdges - nEdges      // total number of vertices along all edges
			aa := getargs(cell.ID, args.ArgsCells)
			if cell.Gndim > 2 {
				// TODO
			} else {
				k := 0
				xx := make([][]float64, nvEtot)
				for _, lvids := range EdgeLocalVertsD[ct] { // loop over edges
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
		if args.WithEdges || args.WithIdsEdges {
			for iedge, lvids := range EdgeLocalVertsD[cell.TypeIndex] {

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
						if args.WithEdges {
							plt.Plot3dLine(x, y, z, args.ArgsEdges)
						}
						if args.WithTagsEdges && len(cell.EdgeTags) > 0 {
							tag := cell.EdgeTags[iedge]
							if tag != 0 {
								txt := io.Sf("%d", tag)
								xc, yc, zc := 0.0, 0.0, 0.0
								for i := 0; i < len(x); i++ {
									xc += x[i]
									yc += y[i]
									zc += z[i]
								}
								xc /= float64(len(x))
								yc /= float64(len(x))
								zc /= float64(len(x))
								plt.Text3d(xc, yc, zc, txt, args.ArgsTagsEdges)
							}
						}
					} else {
						if args.WithEdges {
							plt.Plot(x, y, args.ArgsEdges)
						}
						if args.WithTagsEdges && len(cell.EdgeTags) > 0 {
							tag := cell.EdgeTags[iedge]
							if tag != 0 {
								txt := io.Sf("%d", tag)
								xc, yc := 0.0, 0.0
								for i := 0; i < len(x); i++ {
									xc += x[i]
									yc += y[i]
								}
								xc /= float64(len(x))
								yc /= float64(len(x))
								plt.Text(xc, yc, txt, args.ArgsTagsEdges)
							}
						}
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
			aa := getargs(cell.ID, args.ArgsLins)
			plt.Plot(x, y, aa)
		}

		// cell ids
		if args.WithIdsCells {
			xc := make([]float64, cell.Gndim)
			for i := 0; i < cell.X.M; i++ {
				for j := 0; j < cell.Gndim; j++ {
					xc[j] += cell.X.Get(i, j)
				}
			}
			for i := 0; i < cell.Gndim; i++ {
				xc[i] /= float64(cell.X.M)
			}
			txt := io.Sf("%d", cell.ID)
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
	var xv, yv, zv []float64
	if args.WithVerts || args.WithIdsVerts || args.WithTagsVerts {
		xv = make([]float64, len(o.Verts))
		yv = make([]float64, len(o.Verts))
		if o.Ndim > 2 {
			zv = make([]float64, len(o.Verts))
		}
		for i, v := range o.Verts {
			xv[i], yv[i] = v.X[0], v.X[1]
			if len(v.X) > 2 {
				zv[i] = v.X[2]
			}
			if args.WithIdsVerts {
				txt := io.Sf("%d", v.ID)
				if o.Ndim > 2 {
					plt.Text3d(xv[i], yv[i], zv[i], txt, args.ArgsIdsVerts)
				} else {
					plt.Text(xv[i], yv[i], txt, args.ArgsIdsVerts)
				}
			}
			if args.WithTagsVerts && v.Tag != 0 {
				txt := io.Sf("%d", v.Tag)
				if o.Ndim > 2 {
					plt.Text3d(xv[i], yv[i], zv[i], txt, args.ArgsTagsVerts)
				} else {
					plt.Text(xv[i], yv[i], txt, args.ArgsTagsVerts)
				}
			}
		}
	}

	// draw vertices
	if args.WithVerts {
		if o.Ndim > 2 {
			plt.Plot3dPoints(xv, yv, zv, args.ArgsVerts)
		} else {
			plt.Plot(xv, yv, args.ArgsVerts)
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
