// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hdf implements a hierarchical data format by wrapping the HDF5 library.
// This library also uses the simpler Go-binary (gob) file format that could be used for smaller datasets.
package hdf

/*
#include "hdf5.h"
#include "hdf5_hl.h"
#include "stdlib.h"

// constants from H5Fpublic.h
unsigned int H5Frdwr() { return H5F_ACC_RDWR; }
unsigned int H5Ftrunc() { return H5F_ACC_TRUNC; }

*/
import "C"

import (
	"bytes"
	"encoding/gob"
	"gosl/chk"
	"gosl/io"
	"os"
	"path"
	"strings"
	"unsafe"
)

// File represents a HDF5 file
type File struct {

	// constants
	useGob bool   // use GOB instead of HDF5
	dir    string // directory name
	fname  string // file name: fnKey + ext
	furl   string // furl = join(dir,fname)

	// GOB
	gobBuffer  *bytes.Buffer // buffer for GOB
	gobEnc     *gob.Encoder  // encoder in case of writing
	gobDec     *gob.Decoder  // decoder in case of reading
	gobReading bool          // reading file instead of writing?

	// HDF5
	chunkSize int     // HDF5 chunk size
	hdfHandle C.hid_t // handle
}

// Filename returns the filename; i.e. fileNameKey + extension
func (o File) Filename() string { return o.fname }

// Filepath returns the full filepath, including directory name
func (o File) Filepath() string { return o.furl }

// Create creates a new file, deleting existent one
//
//   Input:
//     dirOut   -- directory name that will be created if non-existent
//                 Note: dirOut may contain environment variables
//     fnameKey -- filename key; e.g. without extension
//     useGob   -- use Go's own format gob instead of HDF5
//
//   Output:
//     returns a new File object where the filename will be:
//       fnameKey + .h5   if useGob == false, or
//       fnameKey + .gob  if useGob == true
//
func Create(dirOut, fnameKey string, useGob bool) (o *File) {

	// constants
	fname, furl := filepath(dirOut, fnameKey, useGob)
	os.MkdirAll(dirOut, 0777)

	// GOB
	if useGob {
		o = new(File)
		o.useGob = true
		o.dir = dirOut
		o.fname = fname
		o.furl = furl
		o.gobBuffer = new(bytes.Buffer)
		o.gobEnc = gob.NewEncoder(o.gobBuffer)
		return
	}

	// HDF5
	cfn := C.CString(furl)
	defer C.free(unsafe.Pointer(cfn))
	o = new(File)
	o.dir = dirOut
	o.fname = fname
	o.furl = furl
	o.chunkSize = 1
	o.hdfHandle = C.H5Fcreate(cfn, C.H5Ftrunc(), C.H5P_DEFAULT, C.H5P_DEFAULT)
	if o.hdfHandle < 0 {
		chk.Panic("failed to create file <%s>", o.furl)
	}
	return
}

// Open opens an existent file for read only
//
//   Input:
//     dirIn    -- directory name where the file is located
//                 Note: dirIn may contain environment variables
//     fnameKey -- filename key; e.g. without extension
//     useGob   -- use Go's own format gob instead of HDF5
//
//   Output:
//     returns a new File object where the filename will be:
//       fnameKey + .h5   if useGob == false, or
//       fnameKey + .gob  if useGob == true
//
func Open(dirIn, fnameKey string, useGob bool) (o *File) {

	// constants
	fname, furl := filepath(dirIn, fnameKey, useGob)

	// GOB
	if useGob {
		o = Create(dirIn, fnameKey, true)
		b := io.ReadFile(o.furl)
		o.gobBuffer = bytes.NewBuffer(b)
		o.gobDec = gob.NewDecoder(o.gobBuffer)
		o.gobReading = true
		return
	}

	// HDF5
	cfn := C.CString(furl)
	defer C.free(unsafe.Pointer(cfn))
	o = new(File)
	o.dir = dirIn
	o.fname = fname
	o.furl = furl
	o.hdfHandle = C.H5Fopen(cfn, C.H5Frdwr(), C.H5P_DEFAULT)
	o.gobReading = true
	if o.hdfHandle < 0 {
		chk.Panic("failed to open file <%s>", o.furl)
	}
	return
}

// Close closes file
func (o *File) Close() {
	if o.useGob {
		if !o.gobReading {
			io.WriteFileD(o.dir, o.fname, o.gobBuffer)
		}
		return
	}
	st := C.H5Fclose(o.hdfHandle)
	if st < 0 {
		chk.Panic("failed to close file <%s>", o.furl)
	}
}

// auxiliary methods ///////////////////////////////////////////////////////////////////////////

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

// hierarchCreate creates hierarchy
func (o *File) hierarchCreate(path string, docreate func(cp *C.char) C.herr_t) {
	res := o.filterPath(path)
	pth := ""
	for i := 0; i < len(res); i++ {
		pth += "/" + res[i]
		cpth := C.CString(pth)
		defer C.free(unsafe.Pointer(cpth))
		if i < len(res)-1 { // create group
			st := C.H5Lexists(o.hdfHandle, cpth, C.H5P_DEFAULT)
			if st < 0 {
				chk.Panic("cannot check whether path=%q exists or not", path)
			}
			if st == 1 { // group exists
				continue
			}
			gid := C.H5Gcreate2(o.hdfHandle, cpth, C.H5P_DEFAULT, C.H5P_DEFAULT, C.H5P_DEFAULT)
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

func (o *File) deGobRnkDims() (rnk int, dims []int, length int) {
	o.gobDec.Decode(&rnk)
	dims = make([]int, rnk)
	o.gobDec.Decode(&dims)
	if rnk == 1 {
		length = dims[0]
	} else if rnk == 2 {
		length = dims[0] * dims[1]
	} else {
		chk.Panic("rank must be 1 or 2. rnk == %v is invalid", rnk)
	}
	return
}

// auxiliary functions /////////////////////////////////////////////////////////////////////////

func filepath(dir, fnameKey string, useGob bool) (filename, fileurl string) {
	ext := io.FnExt(fnameKey)
	if ext != "" {
		filename = fnameKey
	} else {
		ext = ".h5"
		if useGob {
			ext = ".gob"
		}
		filename = fnameKey + ext
	}
	fileurl = path.Join(os.ExpandEnv(dir), filename)
	return
}
