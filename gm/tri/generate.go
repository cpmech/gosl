// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tri

import "C"

/*
#include "triangle.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"gosl/chk"
	"gosl/io"
)

// Point holds input data defining a "point"
type Point struct {
	Tag int     // tag
	X   float64 // x-coordinate
	Y   float64 // y-coordinate
}

// Segment holds input data defining a "segment"
type Segment struct {
	Tag int // tag
	L   int // left point
	R   int // right point
}

// Region holds input data defining a "region"
type Region struct {
	Tag     int     // tag of region
	MaxArea float64 // max area constraint for triangulation of region
	X       float64 // x-coordinate of a point inside region
	Y       float64 // y-coordinate of a point inside region
}

// Hole holds input data defining a "hole"
type Hole struct {
	X float64 // x-coordinate of a point inside hole
	Y float64 // y-coordinate of a point inside hole
}

// Input holds a Planar Straight Line Graph (PSLG)
type Input struct {
	Points   []*Point   // list of points
	Segments []*Segment // list of segments
	Regions  []*Region  // list of regions
	Holes    []*Hole    // list of holes
}

// Generate generates unstructured mesh of triangles
//  globalMaxArea  -- imposes a maximum triangle area constraint; fixed area constraint that applies to every triangle
//  globalMinAngle -- quality mesh generation with no angles smaller than specified value in degrees
//                    globalMinAngle must be in [0, 60]
//  o2             -- generates quadratic triangles
//  verbose        -- shows Triangle messages
//  extraSwitches  -- extra arguments to be passed to Triangle
func (o *Input) Generate(globalMaxArea, globalMinAngle float64, o2, verbose bool, extraSwitches string) (m *Mesh) {

	// check
	if globalMinAngle > 60 {
		chk.Panic("globalMinAngle must not be greater than 60°\n")
	}

	// sizes
	npoints := C.int(len(o.Points))
	nsegments := C.int(len(o.Segments))
	nregions := C.int(len(o.Regions))
	nholes := C.int(len(o.Holes))

	// input data
	var tin C.struct_triangulateio
	C.tioalloc(&tin, npoints, nsegments, nregions, nholes)
	defer func() { C.tiofree(&tin) }()
	for i, p := range o.Points {
		C.setpoint(&tin, C.int(i), C.int(p.Tag), C.double(p.X), C.double(p.Y))
	}
	for i, s := range o.Segments {
		C.setsegment(&tin, C.int(i), C.int(s.Tag), C.int(s.L), C.int(s.R))
	}
	for i, r := range o.Regions {
		C.setregion(&tin, C.int(i), C.int(r.Tag), C.double(r.MaxArea), C.double(r.X), C.double(r.Y))
	}
	for i, h := range o.Holes {
		C.sethole(&tin, C.int(i), C.double(h.X), C.double(h.Y))
	}

	// parameters
	prms := "pzA"
	if !verbose {
		prms += "Q"
	}
	if globalMaxArea > 0 {
		prms += io.Sf("a%g", globalMaxArea+1e-14)
	}
	if globalMinAngle > 0 {
		prms += io.Sf("q%g", globalMinAngle)
	} else {
		prms += "q"
	}
	if o2 {
		prms += "o2"
	}
	prms += extraSwitches
	switches := C.CString(prms)
	defer C.free(unsafe.Pointer(switches))
	if verbose {
		io.Pf("... triangle switches = %q ...\n", prms)
	}

	// generate
	var tout C.struct_triangulateio
	C.triangulate(switches, &tin, &tout, nil)

	// output
	m = new(Mesh)
	m.Verts = make([]*Vertex, tout.numberofpoints)
	m.Cells = make([]*Cell, tout.numberoftriangles)
	for i := 0; i < int(tout.numberofpoints); i++ {
		v := new(Vertex)
		v.ID = i
		if i < len(o.Points) {
			v.Tag = o.Points[i].Tag
		}
		v.X = []float64{
			float64(C.getpoint(C.int(i), 0, &tout)),
			float64(C.getpoint(C.int(i), 1, &tout)),
		}
		m.Verts[i] = v
	}
	for i := 0; i < int(tout.numberoftriangles); i++ {
		c := new(Cell)
		c.ID = i
		c.Tag = int(C.getcelltag(C.int(i), &tout))
		c.V = make([]int, int(tout.numberofcorners))
		c.EdgeTags = make([]int, 3)
		for j := 0; j < int(tout.numberofcorners); j++ {
			c.V[j] = int(C.getcorner(C.int(i), C.int(j), &tout))
		}
		for j := 0; j < 3; j++ {
			c.EdgeTags[j] = int(C.getedgetag(C.int(i), C.int(j), &tout))
		}
		m.Cells[i] = c
	}
	return
}
