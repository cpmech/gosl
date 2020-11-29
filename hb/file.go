// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hb implements a pseudo hierarchical binary (hb) data file format
package hb

import (
	"bytes"
	"encoding/gob"
	"os"
	"path"
	"strings"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// File represents a HDF5 file
type File struct {

	// constants
	dir   string // directory name
	fname string // file name: fnKey + ext
	furl  string // furl = join(dir,fname)

	// go-binary data
	gobBuffer  *bytes.Buffer // buffer for GOB
	gobEnc     *gob.Encoder  // encoder in case of writing
	gobDec     *gob.Decoder  // decoder in case of reading
	gobReading bool          // reading file instead of writing?
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
//
//   Output:
//     returns a new File object where the filename will be:
//       fnameKey + .hb
//
func Create(dirOut, fnameKey string) (o *File) {
	fname, furl := filepath(dirOut, fnameKey)
	o = new(File)
	o.dir = dirOut
	o.fname = fname
	o.furl = furl
	o.gobBuffer = new(bytes.Buffer)
	o.gobEnc = gob.NewEncoder(o.gobBuffer)
	return
}

// Open opens an existent file for read only
//
//   Input:
//     dirIn    -- directory name where the file is located
//                 Note: dirIn may contain environment variables
//     fnameKey -- filename key; e.g. without extension
//
//   Output:
//     returns a new File object where the filename will be:
//       fnameKey + .hb
//
func Open(dirIn, fnameKey string) (o *File) {
	fname, furl := filepath(dirIn, fnameKey)
	o = new(File)
	o.dir = dirIn
	o.fname = fname
	o.furl = furl
	b := io.ReadFile(o.furl)
	o.gobBuffer = bytes.NewBuffer(b)
	o.gobDec = gob.NewDecoder(o.gobBuffer)
	o.gobReading = true
	return
}

// Close closes file
func (o *File) Close() {
	if !o.gobReading {
		io.WriteFileD(o.dir, o.fname, o.gobBuffer)
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

func filepath(dir, fnameKey string) (filename, fileurl string) {
	ext := io.FnExt(fnameKey)
	if ext != "" {
		filename = fnameKey
	} else {
		ext = ".hb"
		filename = fnameKey + ext
	}
	fileurl = path.Join(os.ExpandEnv(dir), filename)
	return
}
