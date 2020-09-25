// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hdf

/*
#include "hdf5.h"
#include "hdf5_hl.h"
#include "stdlib.h"

hid_t H5Tdouble() { return H5T_NATIVE_DOUBLE; }
*/
import "C"

import (
	"gosl/chk"
	"unsafe"
)

// VarArray ///////////////////////////////////////////////////////////////////////////////

// PutVarArray puts a variable array with name described in path into HDF5 file
//  Input:
//    path -- HDF5 path such as "/myvec" or "/group/myvec"
//    v    -- slice of float64
func (o *File) PutVarArray(path string, v []float64) {

	// GOB
	if o.useGob {
		chk.Panic("this method is not available with useGob == true yet")
	}

	// HDF5
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		pt := C.H5PTcreate_fl(o.hdfHandle, cp, C.H5Tdouble(), C.hsize_t(o.chunkSize), -1)
		if pt == C.H5I_INVALID_HID {
			chk.Panic("cannot create []float64 to path=%q", path)
			return -1
		}
		n := len(v)
		if n > 0 {
			st := C.H5PTappend(pt, C.size_t(n), unsafe.Pointer(&v[0]))
			if st < 0 {
				chk.Panic("cannot append data to vector to path=%q", path)
			}
		}
		st := C.H5PTcreate_index(pt)
		if st < 0 {
			chk.Panic("cannot create index in vector to path=%q", path)
		}
		st = C.H5PTclose(pt)
		if st < 0 {
			chk.Panic("cannot close vector in path=%q", path)
		}
		return 0
	})
}

// AppendToArray appends values to a variable array
func (o *File) AppendToArray(path string, v []float64) {

	// GOB
	if o.useGob {
		chk.Panic("this method is not available with useGob == true yet")
	}

	// HDF5
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	pt := C.H5PTopen(o.hdfHandle, cpth)
	if pt == C.H5I_INVALID_HID {
		chk.Panic("cannot open vector in path %q", path)
	}
	st := C.H5PTappend(pt, C.size_t(len(v)), unsafe.Pointer(&v[0]))
	if st < 0 {
		chk.Panic("cannot append data to vector in path=%q", path)
	}
	st = C.H5PTclose(pt)
	if st < 0 {
		chk.Panic("cannot close vector in path=%q", path)
	}
}
