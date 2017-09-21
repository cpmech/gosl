// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// GenQuadRegionHL generates 2D region made of quads (high-level version of GenQuadRegion)
//   NOTE: see GenQuadRegion for more details
func GenQuadRegionHL(ctype, ndivR, ndivS int, xmin, xmax, ymin, ymax float64) (o *Mesh) {
	return GenQuadRegion(ctype, ndivR, ndivS, false, func(i, j, nr, ns int) (x, y float64) {
		dx := (xmax - xmin) / float64(nr-1)
		dy := (ymax - ymin) / float64(ns-1)
		x = xmin + float64(i)*dx
		y = ymin + float64(j)*dy
		return
	})
}

// GenRing2d generates mesh of quads representing a 2D ring
//   ctype -- one of Type{Qua4,Qua8,Qua9,Qua12,Qua16,Qua17}
//   ndivR -- number of divisions along radius
//   ndivA -- number of divisions along alpha
//   r     -- minimum radius
//   R     -- maximum radius
//   alpha -- maximum alpha
//   NOTE: a circular region is created if maxA=2⋅π
func GenRing2d(ctype int, ndivR, ndivA int, r, R, alpha float64) (o *Mesh) {
	circle := math.Abs(alpha-2.0*math.Pi) < 1e-14
	return GenQuadRegion(ctype, ndivR, ndivA, circle, func(i, j, nr, ns int) (x, y float64) {
		dr := (R - r) / float64(nr-1)
		da := alpha / float64(ns-1)
		a := float64(j) * da
		l := r + float64(i)*dr
		x = l * math.Cos(a)
		y = l * math.Sin(a)
		return
	})
}

// GenQuadRegion generates 2D region made of quads
//
//   ctype  -- one of Type{Qua4,Qua8,Qua9,Qua12,Qua16,Qua17}
//   ndivR  -- number of divisions (cells) along r (e.g. x)
//   ndivS  -- number of divisions (cells) along s (e.g. y)
//   circle -- connect last row (s=ndivS) with the previous one (s=0)
//
//   f(i,j,nr,ns) -- is a function that computes the (x,y) coordinates of grid nodes
//                   were nr=ndivR+1 and nr=ndivS+1
//
//   example (to generate a rectangle):
//
//        f := func(i, j, nr, ns int) (x, y float64) {
//        	dx := (xmax - xmin) / float64(nr-1)
//        	dy := (ymax - ymin) / float64(ns-1)
//        	x = xmin + float64(i)*dx
//        	y = ymin + float64(j)*dy
//        	return
//        }
//
//   The boundaries are tagged as below
//
//       34      3       23                  30
//         @-----@------@               +-----------+
//         |            |               |           |
//         |            |               |           |
//       4 @  vertices  @ 2          40 |   edges   | 20
//         |            |               |           |
//         |            |               |           |
//         @-----@------@               +-----------+
//       41      1       12                  10
//
func GenQuadRegion(ctype, ndivR, ndivS int, circle bool, f func(i, j, nr, ns int) (x, y float64)) (o *Mesh) {

	// check
	if ndivR < 1 {
		chk.Panic("number of divisions along u must be greater than zero\n")
	}
	if ndivS < 1 {
		chk.Panic("number of divisions along v must be greater than zero\n")
	}

	// type of cell
	typekey := TypeIndexToKey[ctype]

	// data about shape
	nvCell := NumVerts[ctype]                // number of vertices of cell
	nEdges := len(EdgeLocalVertsD[ctype])    // number of edges of cell
	nvEdge := len(EdgeLocalVertsD[ctype][0]) // number of vertices along edge of cell
	nvEtot := nvEdge*nEdges - nEdges         // total number of vertices along all edges

	nvCen := nvCell - nvEtot     // number of central vertices of cell
	nlCen := nvEdge - 2          // number of lines along central part of cell
	nvlCen := make([]int, nlCen) // number of vertices along central line of cell
	for i := 0; i < nlCen; i++ {
		nvlCen[i] = 2
	}

	// compute total number of vertices
	nvertR := make([]int, 1+nlCen)         // number of vertices along r-lines of mesh
	nvertR[0] = nvEdge*ndivR - (ndivR - 1) // number of vertices along r-line
	for i := 0; i < nlCen; i++ {
		nvertR[1+i] = 2*ndivR - (ndivR - 1) // number of vertices along central lines
	}

	// extra central vertices
	if nvCen > 0 {
		switch ctype {
		case TypeQua9:
			nvlCen[0]++            // per cell
			nvertR[1] += 1 * ndivR // total along cen line
		case TypeQua16:
			nvlCen[0] += 2         // per cell
			nvlCen[1] += 2         // per cell
			nvertR[1] += 2 * ndivR // total along cen line
			nvertR[2] += 2 * ndivR // total along cen line
		case TypeQua17:
			nvlCen[1]++            // per cell
			nvertR[2] += 1 * ndivR // total along cen line
		default:
			chk.Panic("cannot handle central vertices of %q cells\n", typekey)
		}
	}

	// total number of vertices
	nvTot := 0
	for _, nr := range nvertR {
		nvTot += nr * ndivS
	}
	if !circle {
		nvTot += nvertR[0] // verts along top "main" line
	}

	// total number of lines along s
	ns := (1+nlCen)*ndivS + 1 // +1 => top line (needed even with circle)

	// allocate vertices
	o = new(Mesh)
	o.Verts = make([]*Vertex, nvTot)

	// generate vertices
	j, iv := 0, 0
	for j < ns-1 {
		for _, nr := range nvertR {
			for i := 0; i < nr; i++ {
				x, y := f(i, j, nr, ns)
				o.Verts[iv] = &Vertex{ID: iv, Tag: qvtag(i, j, nr, ns), X: []float64{x, y}}
				iv++
			}
			j++
		}
	}
	if !circle {
		nr := nvertR[0]
		for i := 0; i < nr; i++ {
			x, y := f(i, ns-1, nr, ns)
			o.Verts[iv] = &Vertex{ID: iv, Tag: qvtag(i, ns-1, nr, ns), X: []float64{x, y}}
			iv++
		}
	}

	// allocate cells
	ncTot := ndivR * ndivS
	o.Cells = make([]*Cell, ncTot)

	// constants
	sumNvr := 0
	for _, nr := range nvertR {
		sumNvr += nr
	}

	// generate cells
	ic := 0
	for j := 0; j < ndivS; j++ {
		for i := 0; i < ndivR; i++ {

			// compute pivot points
			incR := i * (nvEdge - 1) // increment along r
			remR := nvertR[0] - incR // remainder along r +1
			a := incR + sumNvr*j     // lower pivot
			b := incR + sumNvr*(j+1) // upper pivot
			up := func(idxCenLine int) (idx int) {
				idx = a + remR
				for n := idxCenLine - 1; n >= 0; n-- {
					idx += nvertR[1+n]
				}
				idx += (nvlCen[idxCenLine] - 1) * i
				return
			}

			// correct upper pivot
			if circle && j == ndivS-1 {
				b = incR
			}

			// set connectivity
			var v []int
			switch ctype {
			case TypeQua4:
				v = []int{a, a + 1, b + 1, b}
			case TypeQua8:
				c := up(0)
				v = []int{a, a + 2, b + 2, b, a + 1, c + 1, b + 1, c}
			case TypeQua9:
				c := up(0)
				v = []int{a, a + 2, b + 2, b, a + 1, c + 2, b + 1, c, c + 1}
			case TypeQua12:
				c := up(0)
				d := up(1)
				v = []int{a, a + 3, b + 3, b, a + 1, c + 1, b + 2, d, a + 2, d + 1, b + 1, c}
			case TypeQua16:
				c := up(0)
				d := up(1)
				v = []int{a, a + 3, b + 3, b, a + 1, c + 3, b + 2, d, a + 2, d + 3, b + 1, c, c + 1, c + 2, d + 2, d + 1}
			case TypeQua17:
				c := up(0)
				d := up(1)
				e := up(2)
				v = []int{a, a + 4, b + 4, b, a + 1, c + 1, b + 3, e, a + 2, d + 2, b + 2, d, a + 3, e + 1, b + 1, c, d + 1}
			default:
				chk.Panic("cannot handle cell type = %q at this time\n", typekey)
			}

			// set cell
			o.Cells[ic] = &Cell{ID: ic, Tag: -1, TypeKey: typekey, V: v, EdgeTags: qetag(i, j, ndivR, ndivS)}
			ic++
		}
	}

	// results
	o.CheckAndCalcDerivedVars()
	return
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// qvtag returns vertex tags for quad-meshes
//
//   example:
//
//       34    3   3   3    23
//         @---@---@---@---@
//         |               |
//       4 @               @ 2
//         |               |
//       4 @               @ 2
//         |               |
//       4 @               @ 2
//         |               |
//         @---@---@---@---@
//       41    1   1   1    12
//
//   nx and ny are number of lines (rows/cols)
func qvtag(i, j, nx, ny int) int {
	if i == 0 {
		if j == 0 {
			return 41
		}
		if j == ny-1 {
			return 34
		}
		return 4
	}
	if j == 0 {
		if i == nx-1 {
			return 12
		}
		return 1
	}
	if i == nx-1 {
		if j == ny-1 {
			return 23
		}
		return 2
	}
	if j == ny-1 {
		return 3
	}
	return 0
}

// qetag returns edge tags for quad-meshes
//
//   example:
//
//               30
//         +------------+
//         |            |
//         |            |
//      40 |            | 20
//         |            |
//         |            |
//         +------------+
//               10
//
//   nx and ny are number of cells along each direction
func qetag(i, j, nx, ny int) []int {
	if nx*ny == 1 {
		return []int{10, 20, 30, 40}
	}
	if i > 0 && i < nx-1 && j > 0 && j < ny-1 {
		return nil
	}
	tags := []int{0, 0, 0, 0}
	if j == 0 {
		tags[0] = 10
	}
	if i == nx-1 {
		tags[1] = 20
	}
	if j == ny-1 {
		tags[2] = 30
	}
	if i == 0 {
		tags[3] = 40
	}
	return tags
}
