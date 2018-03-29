// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package h5

/*
#include "hdf5.h"
#include "hdf5_hl.h"
#include "stdlib.h"

// constants from H5Tpublic.h
#ifdef WIN32
	int H5LONG() { return H5T_NATIVE_LLONG; }
#else
	int H5LONG() { return H5T_NATIVE_LONG; }
#endif
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

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

	// GOB
	if o.useGob {
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

	// HDF5
	rnk := C.int(len(dims))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		return C.H5LTmake_dataset(o.hdfHandle, cp, rnk, (*C.hsize_t)(unsafe.Pointer(&dims[0])), C.H5LONG(), unsafe.Pointer(&dat[0]))
	})
}

// putIntsNoGroup puts integers into file without creating groups
func (o *File) putIntsNoGroup(path string, dat []int) {

	// GOB
	if o.useGob {
		o.putInts(path, []int{len(dat)}, dat)
		return
	}

	// HDF5
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	dims := []int{len(dat)}
	st := C.H5LTmake_dataset(o.hdfHandle, cpth, 1, (*C.hsize_t)(unsafe.Pointer(&dims[0])), C.H5LONG(), unsafe.Pointer(&dat[0]))
	if st < 0 {
		chk.Panic("cannot put int array with path=%q in file <%s>", path, o.furl)
	}
}

// getInts gets an array of integers from file
func (o *File) getInts(path string, isScalar, isMatrix bool) (dims, dat []int) {

	// GOB
	if o.useGob {
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

	// HDF5
	o.filterPath(path)
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	rank := 1
	if isMatrix {
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
	if isMatrix {
		dat = make([]int, dims[0]*dims[1])
	} else {
		if isScalar {
			if dims[0] == 0 { // TODO: check why Matlab/Octave scalars set this to zero
				dims[0] = 1
			}
		}
		dat = make([]int, dims[0])
	}
	st = C.H5LTread_dataset(o.hdfHandle, cpth, C.H5LONG(), unsafe.Pointer(&dat[0]))
	if st < 0 {
		chk.Panic("cannot read dataset with path=%q in file=<%s>", path, o.furl)
	}
	return
}
