// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tri wraps Triangle to perform mesh generation a Delaunay triangulation
package tri

/*
#cgo CFLAGS: -Wno-pointer-to-int-cast -Wno-int-to-pointer-cast
#cgo LDFLAGS: -lm
#include "triangle.h"
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

// Delaunay computes 2D Delaunay triangulation using Triangle
//  Input:
//    X = { x0, x1, x2, ... Npoints }
//    Y = { y0, y1, y2, ... Npoints }
//  Ouptut:
//    Verts = { { x0, y0 }, { x1, y1 }, { x2, y2 } ... Nvertices }
//    Cells = { { id0, id1, id2 }, { id0, id1, id2 } ... Ncellls }
func Delaunay(X, Y []float64, verbose bool) (Verts [][]float64, Cells [][]int) {

	// input
	chk.IntAssert(len(X), len(Y))
	n := len(X)
	verb := 0
	if verbose {
		verb = 1
	}

	// perform triangulation
	var T C.struct_triangulateio
	defer func() { C.tiofree(&T) }()
	res := C.delaunay2d(
		&T,
		C.int(n),
		(*C.double)(unsafe.Pointer(&X[0])),
		(*C.double)(unsafe.Pointer(&Y[0])),
		C.int(verb),
	)
	if res != 0 {
		chk.Panic("Delaunay2d failed: Triangle returned %d code\n", res)
	}

	// output
	nverts := int(T.numberofpoints)
	ncells := int(T.numberoftriangles)
	Verts = utl.Alloc(nverts, 2)
	Cells = utl.IntAlloc(ncells, 3)
	for i := 0; i < nverts; i++ {
		Verts[i][0] = float64(C.getpoint(C.int(i), 0, &T))
		Verts[i][1] = float64(C.getpoint(C.int(i), 1, &T))
	}
	for i := 0; i < ncells; i++ {
		Cells[i][0] = int(C.getcorner(C.int(i), 0, &T))
		Cells[i][1] = int(C.getcorner(C.int(i), 1, &T))
		Cells[i][2] = int(C.getcorner(C.int(i), 2, &T))
	}
	return
}
