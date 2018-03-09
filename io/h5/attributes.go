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
	"strings"
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

// String //////////////////////////////////////////////////////////////////////////////////////

// SetStringAttribute sets a string attibute
func (o *File) SetStringAttribute(path, key, val string) {

	// GOB
	if o.useGob {
		if o.gobReading {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.gobEnc.Encode("SetStringAttribute")
		o.gobEnc.Encode(path)
		o.gobEnc.Encode(key)
		o.gobEnc.Encode(val)
		return
	}

	// HDF5
	ckey, cval := C.CString(key), C.CString(val)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cval))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		st := C.H5LTset_attribute_string(o.hdfHandle, cp, ckey, cval)
		if st < 0 {
			chk.Panic("cannot set attibute key to attr in path=%q", path)
		}
		return 0
	})
}

// GetStringAttribute gets string attribute
func (o *File) GetStringAttribute(path, key string) (val string) {

	// GOB
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "SetStringAttribute" {
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
		o.gobDec.Decode(&val)
		return
	}

	// HDF5
	o.filterPath(path)
	val = strings.Repeat(" ", 2048)
	cpth, ckey, cval := C.CString(path), C.CString(key), C.CString(val)
	defer C.free(unsafe.Pointer(cpth))
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cval))
	st := C.H5LTget_attribute_string(o.hdfHandle, cpth, ckey, cval)
	if st < 0 {
		chk.Panic("cannot read attibute %q from val in path=%q", key, path)
	}
	return strings.TrimSpace(C.GoString(cval))
}

// Int /////////////////////////////////////////////////////////////////////////////////////////

// SetIntAttribute sets int attibute
func (o *File) SetIntAttribute(path, key string, val int) {

	// GOB
	if o.useGob {
		if o.gobReading {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.gobEnc.Encode("SetIntAttribute")
		o.gobEnc.Encode(path)
		o.gobEnc.Encode(key)
		o.gobEnc.Encode(val)
		return
	}

	// HDF5
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	vals := []int{val}
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		st := C.H5LTset_attribute_long(o.hdfHandle, cp, ckey, (*C.long)(unsafe.Pointer(&vals[0])), 1)
		if st < 0 {
			chk.Panic("cannot set attibute %q to val in path=%q", key, path)
		}
		return 0
	})
}

// GetIntAttribute gets int attribute
func (o *File) GetIntAttribute(path, key string) (val int) {

	// GOB
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "SetIntAttribute" {
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
		o.gobDec.Decode(&val)
		return
	}

	// HDF5
	o.filterPath(path)
	cpth, ckey := C.CString(path), C.CString(key)
	defer C.free(unsafe.Pointer(cpth))
	defer C.free(unsafe.Pointer(ckey))
	vals := []int{0}
	st := C.H5LTget_attribute_long(o.hdfHandle, cpth, ckey, (*C.long)(unsafe.Pointer(&vals[0])))
	if st < 0 {
		chk.Panic("cannot read attibute %q from val in path=%q", key, path)
	}
	return vals[0]
}

// Ints /////////////////////////////////////////////////////////////////////////////////////////

// SetIntsAttribute sets slice-of-ints attibute
func (o *File) SetIntsAttribute(path, key string, vals []int) {

	// GOB
	if o.useGob {
		if o.gobReading {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.gobEnc.Encode("SetIntsAttribute")
		o.gobEnc.Encode(path)
		o.gobEnc.Encode(key)
		o.gobEnc.Encode(vals)
		return
	}

	// HDF5
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

// GetIntsAttribute gets slice-of-ints attribute
func (o *File) GetIntsAttribute(path, key string) (vals []int) {

	// GOB
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "SetIntsAttribute" {
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

	// HDF5
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
