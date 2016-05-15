// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package msh implements functions to operate on meshes
package msh

/*
#cgo LDFLAGS: -lm
#include "connecttriangle.h"
*/
import "C"
import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

func Delaunay2d(X, Y []float64, verbose bool) (M *Mesh, err error) {

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
	M = new(Mesh)
	M.Verts = make([]*Vert, nverts)
	M.Cells = make([]*Cell, ncells)
	for i := 0; i < nverts; i++ {
		M.Verts[i] = &Vert{
			Id: i,
			C: []float64{
				float64(C.getpoint((C.long)(i), 0, &T)),
				float64(C.getpoint((C.long)(i), 1, &T)),
			},
		}
	}
	for i := 0; i < ncells; i++ {
		M.Cells[i] = &Cell{
			Id:   i,
			Tag:  -1,
			Type: "tri3",
			Verts: []int{
				int(C.getcorner((C.long)(i), 0, &T)),
				int(C.getcorner((C.long)(i), 1, &T)),
				int(C.getcorner((C.long)(i), 2, &T)),
			},
		}
	}
	return
}
