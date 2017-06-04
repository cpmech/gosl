// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// GenRing2d generates mesh of TypeQua8 representing a 2D ring
//   ndivR -- number of divisions along radius
//   ndivθ -- number of divisions along θ
//   rMin  -- minimum radius
//   rMax  -- maximum radius
//   θmax  -- maximum θ
func GenRing2d(ndivR, ndivθ int, rMin, rMax, θmax float64) (o *Mesh, err error) {

	// check
	if ndivR < 2 {
		err = chk.Err("Nr (number of lines/columns for each radii) must be greater than or equal to 2")
		return
	}
	if ndivθ < 2 {
		err = chk.Err("Nth (number of lines/rows for each theta) must be greater than or equal to 2")
		return
	}

	// constants
	nc := (ndivR - 1) * (ndivθ - 1)
	nv := ndivR*ndivθ + ndivR*(ndivθ-1) + (ndivR-1)*ndivθ
	dr := (rMax - rMin) / float64(ndivR-1)
	dθ := θmax / float64(ndivθ-1)

	// allocate geometry
	o = new(Mesh)
	o.Verts = make([]*Vertex, nv)
	o.Cells = make([]*Cell, nc)

	// function to compute vertex tag
	vtag := func(j int) int {
		if j == 0 {
			return 100
		}
		if j == ndivθ-1 {
			return 200
		}
		return 0
	}

	// generate mesh
	ivert := 0
	icell := 0
	for i := 0; i < ndivR-1; i++ {

		// current radius
		r := rMin + float64(i)*dr

		// vertices: first column
		if i == 0 {
			for j := 0; j < ndivθ; j++ {
				θ := θmax - float64(j)*dθ
				x := r * math.Cos(θ)
				y := r * math.Sin(θ)
				o.Verts[ivert] = &Vertex{Id: ivert, Tag: vtag(j), X: []float64{x, y}}
				ivert++
				if j != ndivθ-1 { // intermediate nodes
					θ -= dθ / 2.0
					x = r * math.Cos(θ)
					y = r * math.Sin(θ)
					o.Verts[ivert] = &Vertex{Id: ivert, Tag: 0, X: []float64{x, y}}
					ivert++
				}
			}
		}

		// vertices: middle column
		r += dr / 2.0
		for j := 0; j < ndivθ; j++ {
			th := θmax - float64(j)*dθ
			x := r * math.Cos(th)
			y := r * math.Sin(th)
			o.Verts[ivert] = &Vertex{Id: ivert, Tag: vtag(j), X: []float64{x, y}}
			ivert++
		}

		// vertices: last column
		r += dr / 2.0
		for j := 0; j < ndivθ; j++ {
			th := θmax - float64(j)*dθ
			x := r * math.Cos(th)
			y := r * math.Sin(th)
			o.Verts[ivert] = &Vertex{Id: ivert, Tag: vtag(j), X: []float64{x, y}}
			ivert++
			if j != ndivθ-1 {
				th -= dθ / 2.0
				x = r * math.Cos(th)
				y = r * math.Sin(th)
				o.Verts[ivert] = &Vertex{Id: ivert, Tag: 0, X: []float64{x, y}}
				ivert++
			}
		}

		// set cells
		for j := 0; j < ndivθ-1; j++ {
			a := (3*ndivθ-1)*i + 2*j
			b := (3*ndivθ-1)*i + (2*ndivθ - 1) + j
			c := a + 3*ndivθ - 1
			o.Cells[icell] = &Cell{Id: icell, Tag: -1, TypeKey: "qua8", V: []int{a + 2, c + 2, c, a, b + 1, c + 1, b, a + 1}}
			if i == 0 {
				//SetBryTag(idx_cell, 3, -10) // TODO
			}
			if i == ndivR-2 {
				//SetBryTag(idx_cell, 1, -20) // TODO
			}
			if j == ndivθ-2 {
				//SetBryTag(idx_cell, 0, -30) // TODO
			}
			if j == 0 {
				//SetBryTag(idx_cell, 2, -40) // TODO
			}
			icell++
		}
	}

	// calculate derived
	err = o.CheckAndCalcDerivedVars()
	return
}
