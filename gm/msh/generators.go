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
	if ndivR < 1 {
		err = chk.Err("number of divisions along radius must be greater than zero")
		return
	}
	if ndivθ < 1 {
		err = chk.Err("number of divisions along θ must be greater than zero")
		return
	}

	// number of lines along each direction
	nr := ndivR + 1
	nθ := ndivθ + 1

	// flag
	fullcircle := math.Abs(θmax-2.0*math.Pi) < 1e-14
	if fullcircle {
		nθ = ndivθ
	}

	// constants
	nc := ndivR * ndivθ
	nv := (2*nr-1)*nθ + nr*ndivθ
	dr := (rMax - rMin) / float64(ndivR)
	dθ := θmax / float64(ndivθ)

	// allocate geometry
	o = new(Mesh)
	o.Verts = make([]*Vertex, nv)
	o.Cells = make([]*Cell, nc)

	// function to compute vertex tag
	vtag := func(i, j int) int {
		if i == 0 {
			if j == 0 {
				return 41
			}
			if j == nθ-1 {
				return 43
			}
			return 4
		}
		if j == 0 {
			if i == nr-1 {
				return 21
			}
			return 1
		}
		if i == nr-1 {
			if j == nθ-1 {
				return 23
			}
			return 2
		}
		if j == nθ-1 {
			return 3
		}
		return 0
	}

	// function to compute edge tag
	etag := func(i, j int) []int {
		if nc == 1 {
			return []int{10, 20, 30, 40}
		}
		if i > 0 && i < ndivR-1 && j > 0 && j < ndivθ-1 {
			return nil
		}
		tags := []int{0, 0, 0, 0}
		if j == 0 {
			tags[0] = 10
		}
		if i == ndivR-1 {
			tags[1] = 20
		}
		if j == ndivθ-1 {
			tags[2] = 30
		}
		if i == 0 {
			tags[3] = 40
		}
		return tags
	}

	// set vertices
	ivert := 0
	icell := 0
	for j := 0; j < nθ; j++ {
		θ := float64(j) * dθ
		for i := 0; i < nr; i++ {
			r := rMin + float64(i)*dr
			x := r * math.Cos(θ)
			y := r * math.Sin(θ)
			o.Verts[ivert] = &Vertex{Id: ivert, Tag: vtag(i, j), X: []float64{x, y}}
			ivert++
			if i < nr-1 {
				r += dr / 2.0
				x = r * math.Cos(θ)
				y = r * math.Sin(θ)
				o.Verts[ivert] = &Vertex{Id: ivert, Tag: vtag(-1, j), X: []float64{x, y}}
				ivert++
			}
		}
		midrow := j < nθ-1
		if fullcircle {
			midrow = j < nθ
		}
		if midrow {
			θ += dθ / 2.0
			for i := 0; i < nr; i++ {
				r := rMin + float64(i)*dr
				x := r * math.Cos(θ)
				y := r * math.Sin(θ)
				o.Verts[ivert] = &Vertex{Id: ivert, Tag: vtag(i, -1), X: []float64{x, y}}
				ivert++
			}
		}
	}

	// set cells
	for j := 0; j < ndivθ; j++ {
		for i := 0; i < ndivR; i++ {
			a := (3*nr-1)*j + 2*i            // left-lower
			b := (3*nr-1)*j + (2*nr - 1) + i // left-centre
			c := (3*nr-1)*(j+1) + 2*i        // left-upper
			if fullcircle && j == ndivθ-1 {
				c = 2 * i // left-upper
			}
			o.Cells[icell] = &Cell{Id: icell, Tag: -1, TypeKey: "qua8", V: []int{
				a, a + 2, c + 2, c, a + 1, b + 1, c + 1, b,
			}, EdgeTags: etag(i, j)}
			icell++
		}
	}

	// calculate derived
	err = o.CheckAndCalcDerivedVars()
	return
}
