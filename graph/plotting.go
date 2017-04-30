// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Draw draws grid
//  r   -- radius of circles
//  W   -- width of paths
//  dwt -- δwt for positioning text (w = W/2)
//  arrow_scale -- scale for arrows. use 0 for default value
func (o *Graph) Draw(r, W, dwt, arrow_scale float64,
	verts_lbls map[int]string, verts_fsz float64, verts_clr string,
	edges_lbls map[int]string, edges_fsz float64, edges_clr string) {
	if len(o.Verts) < 1 {
		return
	}
	xmin, ymin := o.Verts[0][0], o.Verts[0][1]
	xmax, ymax := xmin, ymin
	var lbl string
	for i, v := range o.Verts {
		if verts_lbls == nil {
			lbl = io.Sf("%d", i)
		} else {
			lbl = verts_lbls[i]
		}
		plt.Text(v[0], v[1], lbl, &plt.A{C: verts_clr, Fsz: verts_fsz, Va: "center", Ha: "center"})
		plt.Circle(v[0], v[1], r, &plt.A{Fc: "none", Ec: "k", Lw: 0.8, NoClip: true})
		xmin, ymin = utl.Min(xmin, v[0]), utl.Min(ymin, v[1])
		xmax, ymax = utl.Max(xmax, v[0]), utl.Max(ymax, v[1])
	}
	if W > 2*r {
		W = 1.8 * r
	}
	w := W / 2.0
	xa, xb := make([]float64, 2), make([]float64, 2)
	xc, dx := make([]float64, 2), make([]float64, 2)
	mu, nu := make([]float64, 2), make([]float64, 2)
	xi, xj := make([]float64, 2), make([]float64, 2)
	l := math.Sqrt(r*r - w*w)
	var L float64
	if arrow_scale <= 0 {
		arrow_scale = 20
	}
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
			xc[i] = (xi[i]+xj[i])/2.0 - dwt*nu[i]
		}
		plt.Arrow(xi[0], xi[1], xj[0], xj[1], &plt.A{Style: "->", Scale: arrow_scale})
		if edges_lbls == nil {
			lbl = io.Sf("%d", k)
		} else {
			lbl = edges_lbls[k]
		}
		plt.Text(xc[0], xc[1], lbl, &plt.A{C: edges_clr, Fsz: edges_fsz, Va: "center", Ha: "center"})
	}
	xmin -= r
	xmax += r
	ymin -= r
	ymax += r
	plt.AxisRange(xmin, xmax, ymin, ymax)
	return
}
