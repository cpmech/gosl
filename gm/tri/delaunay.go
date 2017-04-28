// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package tri wraps Triangle to perform mesh generation a Delaunay triangulation
package tri

/*
#cgo LDFLAGS: -lm
#include "connecttriangle.h"
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
//    V = { { x0, y0 }, { x1, y1 }, { x2, y2 } ... Nvertices }
//    C = { { id0, id1, id2 }, { id0, id1, id2 } ... Ncellls }
func Delaunay(X, Y []float64, verbose bool) (V [][]float64, C [][]int, err error) {

	// input
	chk.IntAssert(len(X), len(Y))
	n := len(X)
	verb := 0
	if verbose {
		verb = 1
	}

	// perform triangulation
	var T C.triangulateio
	defer func() { C.trifree(&T) }()
	res := C.delaunay2d(
		&T,
		(C.long)(n),
		(*C.double)(unsafe.Pointer(&X[0])),
		(*C.double)(unsafe.Pointer(&Y[0])),
		(C.long)(verb),
	)
	if res != 0 {
		chk.Err("Delaunay2d failed: Triangle returned %d code\n", res)
	}

	// output
	nverts := int(T.numberofpoints)
	ncells := int(T.numberoftriangles)
	V = utl.Alloc(nverts, 2)
	C = utl.IntAlloc(ncells, 3)
	for i := 0; i < nverts; i++ {
		V[i][0] = float64(C.getpoint((C.long)(i), 0, &T))
		V[i][1] = float64(C.getpoint((C.long)(i), 1, &T))
	}
	for i := 0; i < ncells; i++ {
		C[i][0] = int(C.getcorner((C.long)(i), 0, &T))
		C[i][1] = int(C.getcorner((C.long)(i), 1, &T))
		C[i][2] = int(C.getcorner((C.long)(i), 2, &T))
	}
	return
}
