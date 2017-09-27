// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Plotter draws graphs
type Plotter struct {
	G            *Graph         // the graph
	Parts        []int32        // [nverts] partitions
	VertsLabels  map[int]string // [nverts] labels of vertices. may be nil
	EdgesLabels  map[int]string // [nedges] labels of edges. may be nil
	FcParts      []string       // [nverts] face colors of circles indicating partitions. may be nil
	Radius       float64        // radius of circles. can be 0
	Width        float64        // distance between two edges. can be 0
	Gap          float64        // gap between text and edges; e.g. width/2. can be 0
	ArgsVertsTxt *plt.A         // arguments for vertex ids. may be nil
	ArgsEdgesTxt *plt.A         // arguments for edge ids. may be nil
	ArgsVerts    *plt.A         // arguments for vertices. may be nil
	ArgsEdges    *plt.A         // arguments for edges. may be nil
}

// Draw draws graph
func (o *Plotter) Draw() {

	// check
	nverts := len(o.G.Verts)
	if nverts < 1 {
		chk.Panic("vertices are required to draw graph\n")
	}

	// plot arguments
	if o.ArgsVertsTxt == nil {
		o.ArgsVertsTxt = &plt.A{C: "#d70d0d", Fsz: 8, Va: "center", Ha: "center"}
	}
	if o.ArgsEdgesTxt == nil {
		o.ArgsEdgesTxt = &plt.A{C: "#2732c6", Fsz: 7, Va: "center", Ha: "center"}
	}
	if o.ArgsVerts == nil {
		o.ArgsVerts = &plt.A{Fc: "none", Ec: "k", Lw: 0.8, NoClip: true}
	}
	if o.ArgsEdges == nil {
		o.ArgsEdges = &plt.A{Style: "->", Scale: 12}
	}
	if len(o.Parts) == nverts {
		nparts := 0
		for i := 0; i < nverts; i++ {
			nparts = utl.Imax(nparts, int(o.Parts[i]))
		}
		nparts++
		if len(o.FcParts) != nparts {
			o.FcParts = make([]string, nparts)
			for i := 0; i < nparts; i++ {
				o.FcParts[i] = plt.C(i, 9)
			}
		}
	}

	// constants
	if o.Radius < 1e-10 {
		o.Radius = 0.25
	}
	if o.Width < 1e-10 {
		o.Width = 0.15
	}
	if o.Gap < 1e-10 {
		o.Gap = 0.12
	}

	// draw vertex data and find limits
	xmin, ymin := o.G.Verts[0][0], o.G.Verts[0][1]
	xmax, ymax := xmin, ymin
	var lbl string
	for i, v := range o.G.Verts {
		if o.VertsLabels == nil {
			lbl = io.Sf("%d", i)
		} else {
			lbl = o.VertsLabels[i]
		}
		plt.Text(v[0], v[1], lbl, o.ArgsVertsTxt)
		fcOrig := o.ArgsVerts.Fc
		if len(o.Parts) == nverts {
			o.ArgsVerts.Fc = o.FcParts[o.Parts[i]]
		}
		plt.Circle(v[0], v[1], o.Radius, o.ArgsVerts)
		o.ArgsVerts.Fc = fcOrig
		xmin, ymin = utl.Min(xmin, v[0]), utl.Min(ymin, v[1])
		xmax, ymax = utl.Max(xmax, v[0]), utl.Max(ymax, v[1])
	}

	// distance between edges
	if o.Width > 2*o.Radius {
		o.Width = 1.8 * o.Radius
	}
	w := o.Width / 2.0
	l := math.Sqrt(o.Radius*o.Radius - w*w)

	// draw edges
	xa, xb := make([]float64, 2), make([]float64, 2)
	xc, dx := make([]float64, 2), make([]float64, 2)
	mu, nu := make([]float64, 2), make([]float64, 2)
	xi, xj := make([]float64, 2), make([]float64, 2)
	var L float64
	for k, e := range o.G.Edges {
		L = 0.0
		for i := 0; i < 2; i++ {
			xa[i] = o.G.Verts[e[0]][i]
			xb[i] = o.G.Verts[e[1]][i]
			xc[i] = (xa[i] + xb[i]) / 2.0
			dx[i] = xb[i] - xa[i]
			L += dx[i] * dx[i]
		}
		L = math.Sqrt(L)
		mu[0], mu[1] = dx[0]/L, dx[1]/L
		nu[0], nu[1] = -dx[1]/L, dx[0]/L
		for i := 0; i < 2; i++ {
			xi[i] = xa[i] + l*mu[i] - w*nu[i]
			xj[i] = xb[i] - l*mu[i] - w*nu[i]
			xc[i] = (xi[i]+xj[i])/2.0 - o.Gap*nu[i]
		}
		plt.Arrow(xi[0], xi[1], xj[0], xj[1], o.ArgsEdges)
		if o.EdgesLabels == nil {
			lbl = io.Sf("%d", k)
		} else {
			lbl = o.EdgesLabels[k]
		}
		plt.Text(xc[0], xc[1], lbl, o.ArgsEdgesTxt)
	}
	xmin -= o.Radius
	xmax += o.Radius
	ymin -= o.Radius
	ymax += o.Radius
	plt.AxisRange(xmin, xmax, ymin, ymax)
}
