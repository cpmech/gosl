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

// Draw draws grid
//  vertsLabels -- labels of vertices. may be nil
//  edgesLabels -- labels of edges. may be nil
//  radius -- radius of circles. can be 0
//  width -- distance between two edges. can be 0
//  gap -- gap between text and edges; e.g. width/2. can be 0
//  argsVertsTxt -- plt arguments for vertex ids. may be nil
//  argsEdgesTxt -- plt arguments for edge ids. may be nil
//  argsVerts -- plt arguments for vertices. may be nil
//  argsEdges -- plt arguments for edges. may be nil
func (o *Graph) Draw(vertsLabels map[int]string, edgesLabels map[int]string, radius, width, gap float64,
	argsVertsTxt, argsEdgesTxt, argsVerts, argsEdges *plt.A) {

	// check
	if len(o.Verts) < 1 {
		chk.Panic("vertices are required to draw graph\n")
	}

	// plot arguments
	if argsVertsTxt == nil {
		argsVertsTxt = &plt.A{C: "#d70d0d", Fsz: 8, Va: "center", Ha: "center"}
	}
	if argsEdgesTxt == nil {
		argsEdgesTxt = &plt.A{C: "#2732c6", Fsz: 7, Va: "center", Ha: "center"}
	}
	if argsVerts == nil {
		argsVerts = &plt.A{Fc: "none", Ec: "k", Lw: 0.8, NoClip: true}
	}
	if argsEdges == nil {
		argsEdges = &plt.A{Style: "->", Scale: 12}
	}

	// constants
	if radius < 1e-10 {
		radius = 0.25
	}
	if width < 1e-10 {
		width = 0.15
	}
	if gap < 1e-10 {
		gap = 0.12
	}

	// draw vertex data and find limits
	xmin, ymin := o.Verts[0][0], o.Verts[0][1]
	xmax, ymax := xmin, ymin
	var lbl string
	for i, v := range o.Verts {
		if vertsLabels == nil {
			lbl = io.Sf("%d", i)
		} else {
			lbl = vertsLabels[i]
		}
		plt.Text(v[0], v[1], lbl, argsVertsTxt)
		plt.Circle(v[0], v[1], radius, argsVerts)
		xmin, ymin = utl.Min(xmin, v[0]), utl.Min(ymin, v[1])
		xmax, ymax = utl.Max(xmax, v[0]), utl.Max(ymax, v[1])
	}

	// distance between edges
	if width > 2*radius {
		width = 1.8 * radius
	}
	w := width / 2.0
	l := math.Sqrt(radius*radius - w*w)

	// draw edges
	xa, xb := make([]float64, 2), make([]float64, 2)
	xc, dx := make([]float64, 2), make([]float64, 2)
	mu, nu := make([]float64, 2), make([]float64, 2)
	xi, xj := make([]float64, 2), make([]float64, 2)
	var L float64
	for k, e := range o.Edges {
		L = 0.0
		for i := 0; i < 2; i++ {
			xa[i] = o.Verts[e[0]][i]
			xb[i] = o.Verts[e[1]][i]
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
			xc[i] = (xi[i]+xj[i])/2.0 - gap*nu[i]
		}
		plt.Arrow(xi[0], xi[1], xj[0], xj[1], argsEdges)
		if edgesLabels == nil {
			lbl = io.Sf("%d", k)
		} else {
			lbl = edgesLabels[k]
		}
		plt.Text(xc[0], xc[1], lbl, argsEdgesTxt)
	}
	xmin -= radius
	xmax += radius
	ymin -= radius
	ymax += radius
	plt.AxisRange(xmin, xmax, ymin, ymax)
	return
}
