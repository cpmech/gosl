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

// IntPut puts a slice of integers into file
//  NOTE: path = "/myint"  or   path = "/group/myint"
func (o *File) IntPut(path string, v []int) {
	if len(v) < 1 {
		chk.Panic("cannot put empty slice in HDF file. path = %q", path)
	}
	o.putArrayInt(path, []int{len(v)}, v)
}

// IntRead reads a slice of integers from file
func (o *File) IntRead(path string) (v []int) {
	_, v = o.getArrayInt(path, false) // ismat=false
	return
}

// IntsSetAttr sets ints attibute
func (o *File) IntsSetAttr(path, key string, vals []int) {
	if o.useGob {
		if o.gobReading {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.gobEnc.Encode("IntsSetAttr")
		o.gobEnc.Encode(path)
		o.gobEnc.Encode(key)
		o.gobEnc.Encode(vals)
		return
	}
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	n := C.size_t(len(vals))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		st := C.H5LTset_attribute_long(o.hdfHandle, cp, ckey, (*C.long)(unsafe.Pointer(&vals[0])), n)
		if st < 0 {
			chk.Panic("cannot set attibute %q to vals in path=%q", key, path)
		}
		return 0
	})
}

// IntsReadAttr reads ints attribute
func (o *File) IntsReadAttr(path, key string) (vals []int) {
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "IntsSetAttr" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.gobDec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		var rkey string
		o.gobDec.Decode(&rkey)
		if rkey != key {
			chk.Panic("cannot read key: %s != %s\n(r/w commands need to be called in the same order)", key, rkey)
		}
		o.gobDec.Decode(&vals)
		return
	}
	o.filterPath(path)
	cpth, ckey := C.CString(path), C.CString(key)
	defer C.free(unsafe.Pointer(cpth))
	defer C.free(unsafe.Pointer(ckey))
	var rank int
	st := C.H5LTget_attribute_ndims(o.hdfHandle, cpth, ckey, (*C.int)(unsafe.Pointer(&rank))) //unsafe.Pointer(&rank[0])))
	if st < 0 {
		chk.Panic("cannot read attibute %q from rank in path=%q", key, path)
	}
	if rank != 1 {
		chk.Panic("cannot read attibute %q because rank == %d != 1. path=%q", key, rank, path)
	}
	var typeClass C.H5T_class_t
	var typeSize C.size_t
	dims := make([]int, rank)
	st = C.H5LTget_attribute_info(o.hdfHandle, cpth, ckey, (*C.hsize_t)(unsafe.Pointer(&dims[0])), &typeClass, &typeSize)
	if st < 0 {
		chk.Panic("cannot read attibute %q from dims in path=%q", key, path)
	}
	vals = make([]int, dims[0])
	st = C.H5LTget_attribute_long(o.hdfHandle, cpth, ckey, (*C.long)(unsafe.Pointer(&vals[0])))
	if st < 0 {
		chk.Panic("cannot read attibute %q from vals in path=%q", key, path)
	}
	return
}

// auxiliary methods ///////////////////////////////////////////////////////////////////////////

// putArrayInt puts an array of integers into file
func (o *File) putArrayInt(path string, dims []int, dat []int) {
	if o.useGob {
		if o.gobReading {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.gobEnc.Encode("putArrayInt")
		o.gobEnc.Encode(path)
		o.gobEnc.Encode(len(dims))
		o.gobEnc.Encode(dims)
		o.gobEnc.Encode(dat)
		return
	}
	rnk := C.int(len(dims))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		return C.H5LTmake_dataset_long(o.hdfHandle, cp, rnk, (*C.hsize_t)(unsafe.Pointer(&dims[0])), (*C.long)(unsafe.Pointer(&dat[0])))
	})
}

// putArrayIntNoGroups puts integers into file without creating groups
func (o *File) putArrayIntNoGroups(path string, dat []int) {
	if o.useGob {
		o.putArrayInt(path, []int{len(dat)}, dat)
		return
	}
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	dims := []int{len(dat)}
	st := C.H5LTmake_dataset_long(o.hdfHandle, cpth, 1, (*C.hsize_t)(unsafe.Pointer(&dims[0])), (*C.long)(unsafe.Pointer(&dat[0])))
	if st < 0 {
		chk.Panic("cannot put int array with path=%q in file <%s>", path, o.furl)
	}
}

// getArrayInt gets an array of integers from file
func (o *File) getArrayInt(path string, ismat bool) (dims, dat []int) {
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "putArrayInt" {
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
		dat = make([]int, dims[0]*dims[1])
	} else {
		dat = make([]int, dims[0])
	}
	st = C.H5LTread_dataset_long(o.hdfHandle, cpth, (*C.long)(unsafe.Pointer(&dat[0])))
	if st < 0 {
		chk.Panic("cannot read dataset with path=%q in file=<%s>", path, o.furl)
	}
	return
}
