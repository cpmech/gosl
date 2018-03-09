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

// StrSetAttr sets a string attibute
func (o *File) StrSetAttr(path, key, val string) {
	if o.useGob {
		if o.gobReading {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.gobEnc.Encode("StrSetAttr")
		o.gobEnc.Encode(path)
		o.gobEnc.Encode(key)
		o.gobEnc.Encode(val)
		return
	}
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

// StrReadAttr reads string attribute
func (o *File) StrReadAttr(path, key string) (val string) {
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "StrSetAttr" {
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

// IntSetAttr sets int attibute
func (o *File) IntSetAttr(path, key string, val int) {
	if o.useGob {
		if o.gobReading {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.gobEnc.Encode("IntSetAttr")
		o.gobEnc.Encode(path)
		o.gobEnc.Encode(key)
		o.gobEnc.Encode(val)
		return
	}
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

// IntReadAttr reads int attribute
func (o *File) IntReadAttr(path, key string) (val int) {
	if o.useGob {
		var cmd string
		o.gobDec.Decode(&cmd)
		if cmd != "IntSetAttr" {
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
