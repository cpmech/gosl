// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package graph

/*
#cgo linux CFLAGS: -O2
#cgo linux LDFLAGS: -L/usr/lib -L/usr/local/lib -lmetis

#include <metis.h>
*/
import "C"

import (
	"unsafe"

	"gosl/chk"
)

// MetisPartition performs graph partitioning using METIS
func MetisPartition(npart, nvert int, xadj, adjncy []int32, recursive bool) (objval int32, parts []int32) {

	// output array
	parts = make([]int32, nvert)
	if npart < 2 {
		return
	}
	if npart > nvert {
		chk.Panic("number of partitions must be smaller than the number of vertices. npart=%d is invalid. nvert=%d\n", npart, nvert)
	}

	// set default options
	noptions := int(C.METIS_NOPTIONS)
	options := make([]int32, noptions)
	opt := (*C.idx_t)(unsafe.Pointer(&options[0]))
	C.METIS_SetDefaultOptions(opt)

	// information
	ncon := 1 // number of balancing constraints
	npa := int32(npart)
	nv := int32(nvert)
	np := (*C.idx_t)(unsafe.Pointer(&npa))
	n := (*C.idx_t)(unsafe.Pointer(&nv))
	nc := (*C.idx_t)(unsafe.Pointer(&ncon))
	xa := (*C.idx_t)(unsafe.Pointer(&xadj[0]))
	ad := (*C.idx_t)(unsafe.Pointer(&adjncy[0]))

	// output data
	ov := (*C.idx_t)(unsafe.Pointer(&objval))
	p := (*C.idx_t)(unsafe.Pointer(&parts[0]))

	// call METIS
	var status C.int
	if recursive {
		status = C.METIS_PartGraphRecursive(n, nc, xa, ad, nil, nil, nil, np, nil, nil, opt, ov, p)
	} else {
		status = C.METIS_PartGraphKway(n, nc, xa, ad, nil, nil, nil, np, nil, nil, opt, ov, p)
	}
	if status != C.METIS_OK {
		chk.Panic("METIS failed\n")
	}
	return
}
