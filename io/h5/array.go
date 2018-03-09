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
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

// PutArray puts an array with name described in path into HDF5 file
//  Input:
//    path -- HDF5 path such as "/myvec" or "/group/myvec"
//    v    -- slice of float64
func (o *File) PutArray(path string, v []float64) {
	if len(v) < 1 {
		chk.Panic("cannot put empty vector in HDF file. path = %q", path)
	}
	o.putArray(path, []int{len(v)}, v)
}

// GetArray gets an array from file. Memory will be allocated
func (o *File) GetArray(path string) (v []float64) {
	_, v = o.getArray(path, false) // ismat=false
	return
}

// ReadArray reads an array from file into existent pre-allocated memory
//  Input:
//    path -- HDF5 path such as "/myvec" or "/group/myvec"
//  Output:
//    array -- values in pre-allocated array => must know dimension
//    dims  -- dimensions (for confirmation)
func (o *File) ReadArray(v []float64, path string) (dims []int) {
	dims = o.readArray(&v, path, false) // ismat=false
	return
}

// PutFloat64 puts one float64 into file
//  Input:
//    path -- HDF5 path such as "/myvec" or "/group/myvec"
//    val  -- value
//  Note: this is a convenience function wrapping PutArray
func (o *File) PutFloat64(path string, val float64) {
	o.putArray(path, []int{1}, []float64{val})
}

// GetFloat64 gets one float64 from file
//  Note: this is a convenience function wrapping GetArray
func (o *File) GetFloat64(path string) float64 {
	_, v := o.getArray(path, false)
	if len(v) != 1 {
		chk.Panic("failed to get ONE integer\n")
	}
	return v[0]
}

// auxiliary methods ///////////////////////////////////////////////////////////////////////////

// putArray puts an array into file
func (o *File) putArray(path string, dims []int, dat []float64) {

	// GOB
	if o.useGob {
		if o.gobReading {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.gobEnc.Encode("putArray")
		o.gobEnc.Encode(path)
		o.gobEnc.Encode(len(dims))
		o.gobEnc.Encode(dims)
		o.gobEnc.Encode(dat)
		return
	}

	// HDF5
	rnk := C.int(len(dims))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		return C.H5LTmake_dataset_double(o.hdfHandle, cp, rnk, (*C.hsize_t)(unsafe.Pointer(&dims[0])), (*C.double)(unsafe.Pointer(&dat[0])))
	})
}

// getArray gets an array from file
func (o *File) getArray(path string, ismat bool) (dims []int, dat []float64) {

	// GOB
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "putArray" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.gobDec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		var length int
		_, dims, length = o.deGobRnkDims()
		dat = make([]float64, length)
		o.gobDec.Decode(&dat)
		return
	}

	// HDF5
	o.filterPath(path)
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	rank := 1
	if ismat {
		rank = 2
	}
	dims = make([]int, rank)
	st := C.H5LTget_dataset_info(o.hdfHandle, cpth, (*C.hsize_t)(unsafe.Pointer(&dims[0])), nil, nil)
	if st < 0 {
		chk.Panic("cannot read dimensions with path=%q and file <%s>", path, o.furl)
	}
	if len(dims) != rank {
		chk.Panic("size of dims=%d is incorrectly read: %d != %d. path=%q. file <%s>", dims, len(dims), rank, path, o.furl)
	}
	if ismat {
		dat = make([]float64, dims[0]*dims[1])
	} else {
		dat = make([]float64, dims[0])
	}
	st = C.H5LTread_dataset_double(o.hdfHandle, cpth, (*C.double)(unsafe.Pointer(&dat[0])))
	if st < 0 {
		chk.Panic("cannot read dataset with path=%q in file=<%s>", path, o.furl)
	}
	return
}

// readArray gets an array from file and store in pre-allocated variable
func (o *File) readArray(dat *[]float64, path string, ismat bool) (dims []int) {

	// GOB
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "putArray" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.gobDec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		_, dims, _ = o.deGobRnkDims()
		o.gobDec.Decode(dat)
		return
	}

	// HDF5
	o.filterPath(path)
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	rank := 1
	if ismat {
		rank = 2
	}
	dims = make([]int, rank)
	st := C.H5LTget_dataset_info(o.hdfHandle, cpth, (*C.hsize_t)(unsafe.Pointer(&dims[0])), nil, nil)
	if st < 0 {
		chk.Panic("cannot read dimensions with path=%q and file <%s>", path, o.furl)
	}
	if len(dims) != rank {
		chk.Panic("size of dims=%d is incorrectly read: %d != %d. path=%q. file <%s>", dims, len(dims), rank, path, o.furl)
	}
	if ismat {
		if len(*dat) != dims[0]*dims[1] {
			chk.Panic("size of pre-allocated array with matrix data is incorrect. %d != %d. path=%q. file <%s>", len(*dat), dims[0]*dims[1], path, o.furl)
		}
	} else {
		if len(*dat) != dims[0] {
			chk.Panic("size of pre-allocated array with vector data is incorrect. %d != %d. path=%q. file <%s>", len(*dat), dims[0], path, o.furl)
		}
	}
	st = C.H5LTread_dataset_double(o.hdfHandle, cpth, (*C.double)(unsafe.Pointer(&(*dat)[0])))
	if st < 0 {
		chk.Panic("cannot read dataset with path=%q in file=<%s>", path, o.furl)
	}
	return
}
