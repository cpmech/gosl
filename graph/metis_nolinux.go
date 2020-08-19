// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux

package graph

import "gosl/chk"

// MetisPartition performs graph partitioning using METIS
func MetisPartition(npart, nvert int, xadj, adjncy []int32, recursive bool) (objval int32, parts []int32) {
	chk.Panic("MetisPartition is only available on Linux at this moment\n")
	return
}
