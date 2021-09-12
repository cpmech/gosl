// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hb

import (
	"github.com/cpmech/gosl/chk"
)

// PutArray puts an array with name described by path
//  Input:
//    path -- path such as "/myvec" or "/group/myvec"
//    v    -- slice of float64
func (o *File) PutArray(path string, v []float64) {
	if len(v) < 1 {
		chk.Panic("cannot put empty vector in file. path = %q", path)
	}
	o.putArray(path, []int{len(v)}, v)
}

// GetArray gets an array from file. Memory will be allocated
func (o *File) GetArray(path string) (v []float64) {
	_, v = o.getArray(path, false, false)
	return
}

// ReadArray reads an array from file into existent pre-allocated memory
//  Input:
//    path -- path such as "/myvec" or "/group/myvec"
//  Output:
//    array -- values in pre-allocated array => must know dimension
//    dims  -- dimensions (for confirmation)
func (o *File) ReadArray(v []float64, path string) (dims []int) {
	dims = o.readArray(&v, path, false) // ismat=false
	return
}

// PutFloat64 puts one float64 into file
//  Input:
//    path -- path such as "/myvec" or "/group/myvec"
//    val  -- value
//  Note: this is a convenience function wrapping PutArray
func (o *File) PutFloat64(path string, val float64) {
	o.putArray(path, []int{1}, []float64{val})
}

// GetFloat64 gets one float64 from file
//  Note: this is a convenience function wrapping GetArray
func (o *File) GetFloat64(path string) float64 {
	_, v := o.getArray(path, true, false)
	if len(v) != 1 {
		chk.Panic("failed to get ONE integer\n")
	}
	return v[0]
}

// auxiliary methods ///////////////////////////////////////////////////////////////////////////

// putArray puts an array into file
func (o *File) putArray(path string, dims []int, dat []float64) {
	if o.gobReading {
		chk.Panic("cannot put %q because file is open for READ", path)
	}
	o.gobEnc.Encode("putArray") // 1. command
	o.gobEnc.Encode(path)       // 2. path
	o.gobEnc.Encode(len(dims))  // 3. size of dims
	o.gobEnc.Encode(dims)       // 4. dims
	o.gobEnc.Encode(dat)        // 5. data
}

// getArray gets an array from file
func (o *File) getArray(path string, isScalar, isMatrix bool) (dims []int, dat []float64) {
	if !o.gobReading {
		chk.Panic("cannot read %q because file is set for WRITE", path)
	}
	// 1. command
	var cmd string
	o.gobDec.Decode(&cmd)
	if cmd != "putArray" {
		chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
	}
	// 2. path
	var rpath string
	o.gobDec.Decode(&rpath)
	if rpath != path {
		chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
	}
	// 3. size of dims and 4. dims
	var length int
	_, dims, length = o.deGobRnkDims()
	dat = make([]float64, length)
	// 5. data
	o.gobDec.Decode(&dat)
	return
}

// readArray gets an array from file and store in pre-allocated variable
func (o *File) readArray(dat *[]float64, path string, ismat bool) (dims []int) {
	// 1. command
	var cmd string
	o.gobDec.Decode(&cmd)
	if cmd != "putArray" {
		chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
	}
	// 2. path
	var rpath string
	o.gobDec.Decode(&rpath)
	if rpath != path {
		chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
	}
	// 4. dims
	_, dims, _ = o.deGobRnkDims()
	// 5. data
	o.gobDec.Decode(dat)
	return
}
