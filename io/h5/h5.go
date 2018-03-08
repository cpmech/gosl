// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

// Package h5 wraps the HDF5 library. HDF5 is a data model, library, and file format for storing and
// managing data.
package h5

/*
#include "hdf5.h"
#include "hdf5_hl.h"
#include "stdlib.h"

hid_t H5Tdouble() { return H5T_NATIVE_DOUBLE; }

// constants from H5Fpublic.h
unsigned int H5Frdwr() { return H5F_ACC_RDWR; }
unsigned int H5Ftrunc() { return H5F_ACC_TRUNC; }

*/
import "C"

import (
	"bytes"
	"encoding/gob"
	"os"
	"path"
	"strings"
	"unsafe"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// File represents a HDF5 file
type File struct {
	Gob   bool          // use GOB instead of HDF5
	Cs    int           // chunk_size
	dir   string        // directory name
	fname string        // file name
	furl  string        // furl = join(dir,fname)
	h     C.hid_t       // handle
	b     *bytes.Buffer // buffer for GOB
	enc   *gob.Encoder  // encoder in case of writing
	dec   *gob.Decoder  // decoder in case of reading
	read  bool          // reading file instead of writing?
}

// Create creates a new file, deleting existent one
// NOTE: dirOut will be created if non-existent
func Create(dirOut, filename string, useGob bool) (o *File) {
	furl := path.Join(dirOut, filename)
	os.MkdirAll(dirOut, 0777)
	if useGob {
		o = new(File)
		o.Gob = true
		o.dir = dirOut
		o.fname = filename
		o.furl = furl
		o.b = new(bytes.Buffer)
		o.enc = gob.NewEncoder(o.b)
		return
	}
	cfn := C.CString(furl)
	defer C.free(unsafe.Pointer(cfn))
	o = new(File)
	o.Cs = 1
	o.dir = dirOut
	o.fname = filename
	o.furl = furl
	o.h = C.H5Fcreate(cfn, C.H5Ftrunc(), C.H5P_DEFAULT, C.H5P_DEFAULT)
	if o.h < 0 {
		chk.Panic("failed to create file <%s>", o.furl)
	}
	return
}

// Open opens an existent file for read only
func Open(dirIn, filename string, isGob bool) (o *File) {
	if isGob {
		o = Create(dirIn, filename, true)
		b := io.ReadFile(o.furl)
		o.b = bytes.NewBuffer(b)
		o.dec = gob.NewDecoder(o.b)
		o.read = true
		return
	}
	furl := path.Join(dirIn, filename)
	cfn := C.CString(furl)
	defer C.free(unsafe.Pointer(cfn))
	o = new(File)
	o.dir = dirIn
	o.fname = filename
	o.furl = furl
	o.h = C.H5Fopen(cfn, C.H5Frdwr(), C.H5P_DEFAULT)
	o.read = true
	if o.h < 0 {
		chk.Panic("failed to open file <%s>", o.furl)
	}
	return
}

// Close closes file
func (o *File) Close() {
	if o.Gob {
		if !o.read {
			io.WriteFileD(o.dir, o.fname, o.b)
		}
		return
	}
	st := C.H5Fclose(o.h)
	if st < 0 {
		chk.Panic("failed to close file <%s>", o.furl)
	}
}

// VecPut puts a vector with name described in path into HDF5 file
//  NOTE: path = "/myvec"  or   path = "/group/myvec"
func (o *File) VecPut(path string, v []float64) {
	if len(v) < 1 {
		chk.Panic("cannot put empty vector in HDF file. path = %q", path)
	}
	o.putArray(path, []int{len(v)}, v)
}

// VecRead reads a vector from file
func (o *File) VecRead(path string) (v []float64) {
	_, v = o.getArray(path, false) // ismat=false
	return
}

// VecReadInto reads a vector from file into existent pre-allocated variable
func (o *File) VecReadInto(v *[]float64, path string) (dims []int) {
	dims = o.getArrayInto(v, path, false) // ismat=false
	return
}

func matSerialize(path string, a [][]float64) (m, n int, aser []float64) {
	m, n = len(a), len(a[0])
	aser = make([]float64, m*n)
	for i := 0; i < m; i++ {
		if len(a[i]) != n {
			chk.Panic("all rows in matrix must have the same size. path = %q", path)
		}
		for j := 0; j < n; j++ {
			aser[j+i*n] = a[i][j]
		}
	}
	return
}

func matUnserialize(dims []int, aser []float64) (a [][]float64) {
	a = utl.Alloc(dims[0], dims[1])
	for i := 0; i < dims[0]; i++ {
		for j := 0; j < dims[1]; j++ {
			a[i][j] = aser[j+i*dims[1]]
		}
	}
	return
}

// MatPut puts a matrix with name described in path into HDF5 file
//  NOTE: path = "/mymat"  or   path = "/group/mymat"
func (o *File) MatPut(path string, a [][]float64) {
	if len(a) < 1 {
		chk.Panic("cannot put empty matrix in HDF file. path = %q", path)
	}
	if len(a[0]) < 1 {
		chk.Panic("cannot put empty matrix in HDF file. path = %q", path)
	}
	m, n, aser := matSerialize(path, a)
	o.putArray(path, []int{m, n}, aser)
}

// MatRead reads a matrix from file
func (o *File) MatRead(path string) (a [][]float64) {
	dims, aser := o.getArray(path, true) // ismat=true
	return matUnserialize(dims, aser)
}

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

// Deep3Put puts a deep slice with 3 levels and name described in path into HDF5 file
//  NOTE: path = "/mydeep3"  or   path = "/group/mydeep3"
func (o *File) Deep3Put(path string, a [][][]float64) {
	I, P, S := utl.Deep3Serialize(a)
	o.putArray(path+"/S", []int{len(S)}, S)
	o.putArrayIntNoGroups(path+"/I", I)
	o.putArrayIntNoGroups(path+"/P", P)
}

// Deep3Read reads a deep slice with 3 levels from file
func (o *File) Deep3Read(path string) (a [][][]float64) {
	_, S := o.getArray(path+"/S", false) // ismat=false
	_, I := o.getArrayInt(path+"/I", false)
	_, P := o.getArrayInt(path+"/P", false)
	a = utl.Deep3Deserialize(I, P, S, false)
	return
}

// VarVecPut puts a variable length vector
func (o *File) VarVecPut(path string, v []float64) {
	if o.Gob {
		chk.Panic("this method is not available with o.Gob == true yet")
	}
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		pt := C.H5PTcreate_fl(o.h, cp, C.H5Tdouble(), C.hsize_t(o.Cs), -1)
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

// VarVecAppend appends to a variable length vector
func (o *File) VarVecAppend(path string, v []float64) {
	if o.Gob {
		chk.Panic("this method is not available with o.Gob == true yet")
	}
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	pt := C.H5PTopen(o.h, cpth)
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

// TabAppend appends a row to table
func (o *File) TabAppend(path string, r []float64) {
	if o.Gob {
		chk.Panic("this method is not available with o.Gob == true yet")
	}
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	pt := C.H5PTopen(o.h, cpth)
	if pt == C.H5I_INVALID_HID {
		chk.Panic("cannot open table in path %q", path)
	}
	st := C.H5PTappend(pt, 1, unsafe.Pointer(&r[0]))
	if st < 0 {
		chk.Panic("cannot append data to table in path=%q", path)
	}
	st = C.H5PTclose(pt)
	if st < 0 {
		chk.Panic("cannot close table in path=%q", path)
	}
}

// TabPut puts a table
func (o *File) TabPut(path string, keys []string, a [][]float64) {
	if o.Gob {
		chk.Panic("this method is not available with o.Gob == true yet")
	}
	if len(a) < 1 {
		chk.Panic("cannot put empty table in HDF file. path=%q", path)
	}
	if len(a[0]) < 1 {
		chk.Panic("cannot put empty table in HDF file. path=%q", path)
	}
	var allkeys string
	for _, key := range keys {
		allkeys += " " + key
	}
	allkeys = strings.TrimSpace(allkeys)
	sncol, skeys, kkeys := C.CString("ncol"), C.CString("keys"), C.CString(allkeys)
	defer C.free(unsafe.Pointer(sncol))
	defer C.free(unsafe.Pointer(skeys))
	defer C.free(unsafe.Pointer(kkeys))
	m, n, aser := matSerialize(path, a)
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		ncol := []int{n}
		hid := C.H5Tarray_create(C.H5Tdouble(), 1, (*C.hsize_t)(unsafe.Pointer(&ncol[0])))
		if hid == C.H5I_INVALID_HID {
			chk.Panic("cannot create data type for table in path=%q", path)
		}
		pt := C.H5PTcreate_fl(o.h, cp, hid, C.hsize_t(o.Cs), -1)
		if pt == C.H5I_INVALID_HID {
			chk.Panic("cannot create table in path=%q", path)
		}
		st := C.H5LTset_attribute_long(o.h, cp, sncol, (*C.long)(unsafe.Pointer(&ncol[0])), 1)
		if st < 0 {
			chk.Panic("cannot set attibute ncol to table in path=%q", path)
		}
		st = C.H5LTset_attribute_string(o.h, cp, skeys, kkeys)
		if st < 0 {
			chk.Panic("cannot set attibute keys to table in path=%q", path)
		}
		st = C.H5PTappend(pt, C.size_t(m), unsafe.Pointer(&aser[0]))
		if st < 0 {
			chk.Panic("cannot append data to table in path=%q", path)
		}
		st = C.H5PTcreate_index(pt)
		if st < 0 {
			chk.Panic("cannot create index in table of path=%q", path)
		}
		st = C.H5PTclose(pt)
		if st < 0 {
			chk.Panic("cannot close table in path=%q", path)
		}
		return 0
	})
}

// TabRead reads a table
func (o *File) TabRead(path string) (keys []string, a [][]float64) {
	if o.Gob {
		chk.Panic("this method is not available with o.Gob == true yet")
	}
	o.filterPath(path)
	cpth, sncol, skeys, kkeys := C.CString(path), C.CString("ncol"), C.CString("keys"), C.CString("")
	defer C.free(unsafe.Pointer(cpth))
	defer C.free(unsafe.Pointer(sncol))
	defer C.free(unsafe.Pointer(skeys))
	defer C.free(unsafe.Pointer(kkeys))
	rank := 2
	dims := make([]int, rank)
	st := C.H5LTget_dataset_info(o.h, cpth, (*C.hsize_t)(unsafe.Pointer(&dims[0])), nil, nil)
	if st < 0 {
		chk.Panic("cannot read dimensions with path=%q and file <%s>", "TabRead", path, o.furl)
	}
	ncol := []int{0}
	st = C.H5LTget_attribute_long(o.h, cpth, sncol, (*C.long)(unsafe.Pointer(&ncol[0])))
	if st < 0 {
		chk.Panic("cannot read attibute ncol from table in path=%q", path)
	}
	st = C.H5LTget_attribute_string(o.h, cpth, skeys, kkeys)
	if st < 0 {
		chk.Panic("cannot read attibute keys from table in path=%q", path)
	}
	hid := C.H5Tarray_create(C.H5Tdouble(), 1, (*C.hsize_t)(unsafe.Pointer(&ncol[0])))
	if hid == C.H5I_INVALID_HID {
		chk.Panic("cannot create data type for table in path=%q", path)
	}
	dims[1] = ncol[0]
	aser := make([]float64, dims[0]*dims[1])
	st = C.H5LTread_dataset(o.h, cpth, hid, unsafe.Pointer(&aser[0]))
	if st < 0 {
		chk.Panic("cannot read dataset with path=%q in file=<%s>\n  ncol=%v  dims=%v  keys=%v\n", path, o.furl, ncol, dims, C.GoString(kkeys))
	}
	keys = strings.Split(strings.TrimSpace(C.GoString(kkeys)), " ")
	a = matUnserialize(dims, aser)
	return
}

// StrSetAttr sets a string attibute
func (o *File) StrSetAttr(path, key, val string) {
	if o.Gob {
		if o.read {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.enc.Encode("StrSetAttr")
		o.enc.Encode(path)
		o.enc.Encode(key)
		o.enc.Encode(val)
		return
	}
	ckey, cval := C.CString(key), C.CString(val)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cval))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		st := C.H5LTset_attribute_string(o.h, cp, ckey, cval)
		if st < 0 {
			chk.Panic("cannot set attibute key to attr in path=%q", path)
		}
		return 0
	})
}

// StrReadAttr reads string attribute
func (o *File) StrReadAttr(path, key string) (val string) {
	if o.Gob {
		var cmd string
		o.dec.Decode(&cmd)
		if cmd != "StrSetAttr" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.dec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		var rkey string
		o.dec.Decode(&rkey)
		if rkey != key {
			chk.Panic("cannot read key: %s != %s\n(r/w commands need to be called in the same order)", key, rkey)
		}
		o.dec.Decode(&val)
		return
	}
	o.filterPath(path)
	val = strings.Repeat(" ", 2048)
	cpth, ckey, cval := C.CString(path), C.CString(key), C.CString(val)
	defer C.free(unsafe.Pointer(cpth))
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cval))
	st := C.H5LTget_attribute_string(o.h, cpth, ckey, cval)
	if st < 0 {
		chk.Panic("cannot read attibute %q from val in path=%q", key, path)
	}
	return strings.TrimSpace(C.GoString(cval))
}

// IntSetAttr sets int attibute
func (o *File) IntSetAttr(path, key string, val int) {
	if o.Gob {
		if o.read {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.enc.Encode("IntSetAttr")
		o.enc.Encode(path)
		o.enc.Encode(key)
		o.enc.Encode(val)
		return
	}
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	vals := []int{val}
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		st := C.H5LTset_attribute_long(o.h, cp, ckey, (*C.long)(unsafe.Pointer(&vals[0])), 1)
		if st < 0 {
			chk.Panic("cannot set attibute %q to val in path=%q", key, path)
		}
		return 0
	})
}

// IntReadAttr reads int attribute
func (o *File) IntReadAttr(path, key string) (val int) {
	if o.Gob {
		var cmd string
		o.dec.Decode(&cmd)
		if cmd != "IntSetAttr" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.dec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		var rkey string
		o.dec.Decode(&rkey)
		if rkey != key {
			chk.Panic("cannot read key: %s != %s\n(r/w commands need to be called in the same order)", key, rkey)
		}
		o.dec.Decode(&val)
		return
	}
	o.filterPath(path)
	cpth, ckey := C.CString(path), C.CString(key)
	defer C.free(unsafe.Pointer(cpth))
	defer C.free(unsafe.Pointer(ckey))
	vals := []int{0}
	st := C.H5LTget_attribute_long(o.h, cpth, ckey, (*C.long)(unsafe.Pointer(&vals[0])))
	if st < 0 {
		chk.Panic("cannot read attibute %q from val in path=%q", key, path)
	}
	return vals[0]
}

// IntsSetAttr sets ints attibute
func (o *File) IntsSetAttr(path, key string, vals []int) {
	if o.Gob {
		if o.read {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.enc.Encode("IntsSetAttr")
		o.enc.Encode(path)
		o.enc.Encode(key)
		o.enc.Encode(vals)
		return
	}
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	n := C.size_t(len(vals))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		st := C.H5LTset_attribute_long(o.h, cp, ckey, (*C.long)(unsafe.Pointer(&vals[0])), n)
		if st < 0 {
			chk.Panic("cannot set attibute %q to vals in path=%q", key, path)
		}
		return 0
	})
}

// IntsReadAttr reads ints attribute
func (o *File) IntsReadAttr(path, key string) (vals []int) {
	if o.Gob {
		var cmd string
		o.dec.Decode(&cmd)
		if cmd != "IntsSetAttr" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.dec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		var rkey string
		o.dec.Decode(&rkey)
		if rkey != key {
			chk.Panic("cannot read key: %s != %s\n(r/w commands need to be called in the same order)", key, rkey)
		}
		o.dec.Decode(&vals)
		return
	}
	o.filterPath(path)
	cpth, ckey := C.CString(path), C.CString(key)
	defer C.free(unsafe.Pointer(cpth))
	defer C.free(unsafe.Pointer(ckey))
	var rank int
	st := C.H5LTget_attribute_ndims(o.h, cpth, ckey, (*C.int)(unsafe.Pointer(&rank))) //unsafe.Pointer(&rank[0])))
	if st < 0 {
		chk.Panic("cannot read attibute %q from rank in path=%q", key, path)
	}
	if rank != 1 {
		chk.Panic("cannot read attibute %q because rank == %d != 1. path=%q", key, rank, path)
	}
	var typeClass C.H5T_class_t
	var typeSize C.size_t
	dims := make([]int, rank)
	st = C.H5LTget_attribute_info(o.h, cpth, ckey, (*C.hsize_t)(unsafe.Pointer(&dims[0])), &typeClass, &typeSize)
	if st < 0 {
		chk.Panic("cannot read attibute %q from dims in path=%q", key, path)
	}
	vals = make([]int, dims[0])
	st = C.H5LTget_attribute_long(o.h, cpth, ckey, (*C.long)(unsafe.Pointer(&vals[0])))
	if st < 0 {
		chk.Panic("cannot read attibute %q from vals in path=%q", key, path)
	}
	return
}

// -------- Auxiliary methods ------------------------------------------------------------------------

// filterPath checks path syntax and return a list split by '/'
func (o *File) filterPath(path string) []string {
	if len(path) < 1 {
		chk.Panic("path must be contain at least 1 character, including '/'. path=%q is invalid. file =<%s>", path, o.furl)
	}
	if path[0] != '/' {
		chk.Panic("first character of path must be '/'. path=%q is invalid. file =<%s>", path, o.furl)
	}
	return strings.Split(path, "/")[1:]
}

func (o *File) hierarchCreate(path string, docreate func(cp *C.char) C.herr_t) {
	res := o.filterPath(path)
	pth := ""
	for i := 0; i < len(res); i++ {
		pth += "/" + res[i]
		cpth := C.CString(pth)
		defer C.free(unsafe.Pointer(cpth))
		if i < len(res)-1 { // create group
			st := C.H5Lexists(o.h, cpth, C.H5P_DEFAULT)
			if st < 0 {
				chk.Panic("cannot check whether path=%q exists or not", path)
			}
			if st == 1 { // group exists
				continue
			}
			gid := C.H5Gcreate2(o.h, cpth, C.H5P_DEFAULT, C.H5P_DEFAULT, C.H5P_DEFAULT)
			if gid < 0 {
				chk.Panic("cannot create group with path=%q in file <%s>", path, o.furl)
			}
			C.H5Gclose(gid)
		} else { // create dataset of other structures
			st := docreate(cpth)
			if st < 0 {
				chk.Panic("cannot create dataset/structure with path=%q in file <%s>", path, o.furl)
			}
		}
	}
}

// putArray puts an array into file
func (o *File) putArray(path string, dims []int, dat []float64) {
	if o.Gob {
		if o.read {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.enc.Encode("putArray")
		o.enc.Encode(path)
		o.enc.Encode(len(dims))
		o.enc.Encode(dims)
		o.enc.Encode(dat)
		return
	}
	rnk := C.int(len(dims))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		return C.H5LTmake_dataset_double(o.h, cp, rnk, (*C.hsize_t)(unsafe.Pointer(&dims[0])), (*C.double)(unsafe.Pointer(&dat[0])))
	})
}

// putArrayInt puts an array of integers into file
func (o *File) putArrayInt(path string, dims []int, dat []int) {
	if o.Gob {
		if o.read {
			chk.Panic("cannot put %q because file is open for READONLY", path)
		}
		o.enc.Encode("putArrayInt")
		o.enc.Encode(path)
		o.enc.Encode(len(dims))
		o.enc.Encode(dims)
		o.enc.Encode(dat)
		return
	}
	rnk := C.int(len(dims))
	o.hierarchCreate(path, func(cp *C.char) C.herr_t {
		return C.H5LTmake_dataset_long(o.h, cp, rnk, (*C.hsize_t)(unsafe.Pointer(&dims[0])), (*C.long)(unsafe.Pointer(&dat[0])))
	})
}

// putArrayIntNoGroups puts integers into file without creating groups
func (o *File) putArrayIntNoGroups(path string, dat []int) {
	if o.Gob {
		o.putArrayInt(path, []int{len(dat)}, dat)
		return
	}
	cpth := C.CString(path)
	defer C.free(unsafe.Pointer(cpth))
	dims := []int{len(dat)}
	st := C.H5LTmake_dataset_long(o.h, cpth, 1, (*C.hsize_t)(unsafe.Pointer(&dims[0])), (*C.long)(unsafe.Pointer(&dat[0])))
	if st < 0 {
		chk.Panic("cannot put int array with path=%q in file <%s>", path, o.furl)
	}
}

func (o *File) deGobRnkDims() (rnk int, dims []int, length int) {
	o.dec.Decode(&rnk)
	dims = make([]int, rnk)
	o.dec.Decode(&dims)
	if rnk == 1 {
		length = dims[0]
	} else if rnk == 2 {
		length = dims[0] * dims[1]
	} else {
		chk.Panic("rank must be 1 or 2. rnk == %v is invalid", rnk)
	}
	return
}

// getArray gets an array from file
func (o *File) getArray(path string, ismat bool) (dims []int, dat []float64) {
	if o.Gob {
		var cmd string
		o.dec.Decode(&cmd)
		if cmd != "putArray" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.dec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		var length int
		_, dims, length = o.deGobRnkDims()
		dat = make([]float64, length)
		o.dec.Decode(&dat)
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
	st := C.H5LTget_dataset_info(o.h, cpth, (*C.hsize_t)(unsafe.Pointer(&dims[0])), nil, nil)
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
	st = C.H5LTread_dataset_double(o.h, cpth, (*C.double)(unsafe.Pointer(&dat[0])))
	if st < 0 {
		chk.Panic("cannot read dataset with path=%q in file=<%s>", path, o.furl)
	}
	return
}

// getArrayInt gets an array of integers from file
func (o *File) getArrayInt(path string, ismat bool) (dims, dat []int) {
	if o.Gob {
		var cmd string
		o.dec.Decode(&cmd)
		if cmd != "putArrayInt" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.dec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		var length int
		_, dims, length = o.deGobRnkDims()
		dat = make([]int, length)
		o.dec.Decode(&dat)
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
	st := C.H5LTget_dataset_info(o.h, cpth, (*C.hsize_t)(unsafe.Pointer(&dims[0])), nil, nil)
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
	st = C.H5LTread_dataset_long(o.h, cpth, (*C.long)(unsafe.Pointer(&dat[0])))
	if st < 0 {
		chk.Panic("cannot read dataset with path=%q in file=<%s>", path, o.furl)
	}
	return
}

// getArrayInto gets an array from file and store in pre-allocated variable
func (o *File) getArrayInto(dat *[]float64, path string, ismat bool) (dims []int) {
	if o.Gob {
		var cmd string
		o.dec.Decode(&cmd)
		if cmd != "putArray" {
			chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
		}
		var rpath string
		o.dec.Decode(&rpath)
		if rpath != path {
			chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
		}
		_, dims, _ = o.deGobRnkDims()
		o.dec.Decode(dat)
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
	st := C.H5LTget_dataset_info(o.h, cpth, (*C.hsize_t)(unsafe.Pointer(&dims[0])), nil, nil)
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
	st = C.H5LTread_dataset_double(o.h, cpth, (*C.double)(unsafe.Pointer(&(*dat)[0])))
	if st < 0 {
		chk.Panic("cannot read dataset with path=%q in file=<%s>", path, o.furl)
	}
	return
}
