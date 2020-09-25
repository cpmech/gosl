// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hdf

import "gosl/utl"

// PutDeep3 puts a deep slice with 3 levels and name described in path into HDF5 file
//  Input:
//    path -- HDF5 path such as "/myvec" or "/group/myvec"
//    a    -- slice of slices of slices of float64
//  Note: Slice will be serialized
func (o *File) PutDeep3(path string, a [][][]float64) {
	I, P, S := utl.SerializeDeep3(a)
	o.putArray(path+"/S", []int{len(S)}, S)
	o.putIntsNoGroup(path+"/I", I)
	o.putIntsNoGroup(path+"/P", P)
}

// GetDeep3 gets a deep slice with 3 levels from file. Memory will be allocated
func (o *File) GetDeep3(path string) (a [][][]float64) {
	_, S := o.getArray(path+"/S", false, false)
	_, I := o.getInts(path+"/I", false, false)
	_, P := o.getInts(path+"/P", false, false)
	a = utl.DeserializeDeep3(I, P, S, false)
	return
}
