// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hb

import "gosl/chk"

// PutInts puts a slice of integers into file
//  Input:
//    path -- HDF5 path such as "/myvec" or "/group/myvec"
//    v    -- slice of integers
func (o *File) PutInts(path string, v []int) {
	if len(v) < 1 {
		chk.Panic("cannot put empty slice in HDF file. path = %q", path)
	}
	o.putInts(path, []int{len(v)}, v)
}

// GetInts gets a slice of ints from file. Memory will be allocated
func (o *File) GetInts(path string) (v []int) {
	_, v = o.getInts(path, false, false)
	return
}

// PutInt puts one integer into file
//  Input:
//    path -- HDF5 path such as "/myvec" or "/group/myvec"
//    val  -- value
//  Note: this is a convenience function wrapping PutInts
func (o *File) PutInt(path string, val int) {
	o.putInts(path, []int{1}, []int{val})
}

// GetInt gets one integer from file
//  Note: this is a convenience function wrapping GetInts
func (o *File) GetInt(path string) int {
	_, v := o.getInts(path, true, false)
	if len(v) != 1 {
		chk.Panic("failed to get ONE integer\n")
	}
	return v[0]
}

// auxiliary methods ///////////////////////////////////////////////////////////////////////////

// putInts puts an array of integers into file
func (o *File) putInts(path string, dims []int, dat []int) {
	if o.gobReading {
		chk.Panic("cannot put %q because file is open for READONLY", path)
	}
	o.gobEnc.Encode("putInts")
	o.gobEnc.Encode(path)
	o.gobEnc.Encode(len(dims))
	o.gobEnc.Encode(dims)
	o.gobEnc.Encode(dat)
	return
}

// putIntsNoGroup puts integers into file without creating groups
func (o *File) putIntsNoGroup(path string, dat []int) {
	o.putInts(path, []int{len(dat)}, dat)
	return
}

// getInts gets an array of integers from file
func (o *File) getInts(path string, isScalar, isMatrix bool) (dims, dat []int) {
	var cmd string
	o.gobDec.Decode(&cmd)
	if cmd != "putInts" {
		chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
	}
	var rpath string
	o.gobDec.Decode(&rpath)
	if rpath != path {
		chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
	}
	var length int
	_, dims, length = o.deGobRnkDims()
	dat = make([]int, length)
	o.gobDec.Decode(&dat)
	return
}
