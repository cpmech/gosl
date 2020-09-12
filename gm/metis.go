// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

/*
#cgo linux CFLAGS: -O2
#cgo linux LDFLAGS: -L/usr/lib -L/usr/local/lib -lmetis

#include <metis.h>
*/
import "C"

import (
	"unsafe"

	"gosl/chk"
	"gosl/utl"
)

// MetisShares returns a map of shares owned by vertices
// i.e. each vertex is shared by a number of edges, so,
// we return the map vertexId => [edgeIds...] attached to this vertexId
//
// Example:
//               0        1
//           0-------1--------2
//           |       |        |
//          2|      3|       4|
//           |       |        |
//           3-------4--------5
//               5       6
//
// Input. edges = {0,1}, {1,2}, {0,3}, {1,4}, {2,5}, {3,4}, {4,5}
//
// Output. shares = 0:(0,2), 1:(0,1,3), 2:(1,4)
//                  3:(2,5), 4:(3,5,6), 5:(4,6)
// (notation. vertexID:(firstEdgeID, secondEdgeID))
//
// NOTE: (1) the pairs or triples will have sorted edgeIDs
//       (2) len(shares) = number_of_vertices
//
func MetisShares(edges [][2]int) (shares map[int][]int) {
	shares = make(map[int][]int) // [nverts] edges sharing a vertex
	for k, edge := range edges {
		i, j := edge[0], edge[1]
		utl.IntIntsMapAppend(shares, i, k)
		utl.IntIntsMapAppend(shares, j, k)
	}
	return
}

// MetisAdjacency returns adjacency list as a compressed storage format for METIS
// shares is the map returned by MetisShares
func MetisAdjacency(edges [][2]int, shares map[int][]int) (xadj, adjncy []int32) {

	// total size of array
	nv := len(shares)
	szadj := 0
	for vid := 0; vid < nv; vid++ {
		szadj += len(shares[vid]) // = number of connected vertices
	}

	// data for METIS
	xadj = make([]int32, nv+1)
	adjncy = make([]int32, szadj)
	k := 0
	for vid := 0; vid < nv; vid++ {
		list := shares[vid]
		for _, eid := range list {
			otherVid := edges[eid][0]
			if otherVid == vid {
				otherVid = edges[eid][1]
			}
			adjncy[k] = int32(otherVid)
			k++
		}
		xadj[1+vid] = xadj[vid] + int32(len(list))
	}
	return
}

// MetisPartitionLowLevel performs graph partitioning using METIS
func MetisPartitionLowLevel(npart, nvert int, xadj, adjncy []int32, recursive bool) (objval int32, parts []int32) {

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

// MetisPartition performs graph partitioning using METIS
// This function computes the shares, adjacency list, and partition by calling the other 3 functions
func MetisPartition(edges [][2]int, npart int, recursive bool) (shares map[int][]int, objval int32, parts []int32) {
	shares = MetisShares(edges)
	xadj, adjncy := MetisAdjacency(edges, shares)
	nvert := len(shares)
	objval, parts = MetisPartitionLowLevel(npart, nvert, xadj, adjncy, recursive)
	return
}
