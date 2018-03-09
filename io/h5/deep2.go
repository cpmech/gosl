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

// MatPut puts a matrix with name described in path into HDF5 file
//  NOTE: path = "/mymat"  or   path = "/group/mymat"
func (o *File) MatPut(path string, a [][]float64) {
	m := len(a)
	if m < 1 {
		chk.Panic("cannot put matrix in HDF file. path = %q", path)
	}
	n := len(a[0])
	if n < 1 {
		chk.Panic("cannot put empty matrix in HDF file. path = %q", path)
	}
	aser := utl.SerializeDeep2(a)
	o.putArray(path, []int{m, n}, aser)
}

// MatRead reads a matrix from file
func (o *File) MatRead(path string) (a [][]float64) {
	dims, aser := o.getArray(path, true) // ismat=true
	return utl.DeserializeDeep2(aser, dims[0], dims[1])
}
