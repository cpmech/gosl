// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package h5

/*
#include "hdf5.h"
#include "hdf5_hl.h"
#include "stdlib.h"
*/
import "C"

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

// PutDeep2 puts a Deep2 slice into file
//  Input:
//    path -- HDF5 path such as "/myvec" or "/group/myvec"
//    a    -- slice of slices of float64
//  Note: Slice will be serialized
func (o *File) PutDeep2(path string, a [][]float64) {
	m := len(a)
	if m < 1 {
		chk.Panic("cannot put empty Deep2 into file. path = %q", path)
	}
	n := len(a[0])
	if n < 1 {
		chk.Panic("cannot put empty Deep2 into file. path = %q", path)
	}
	aser := utl.SerializeDeep2(a)
	o.putArray(path, []int{m, n}, aser)
}

// GetDeep2 gets a Deep2 slice (that was serialized). Memory will be allocated
func (o *File) GetDeep2(path string) (a [][]float64) {
	dims, aser := o.getArray(path, false, true)
	return utl.DeserializeDeep2(aser, dims[0], dims[1])
}

// GetDeep2raw returns the serialized data corresponding to a Deep2 slice
func (o *File) GetDeep2raw(path string) (m, n int, a []float64) {
	var dims []int
	dims, a = o.getArray(path, false, true)
	m, n = dims[0], dims[1]
	return
}
