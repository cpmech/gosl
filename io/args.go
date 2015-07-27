// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"flag"

	"github.com/cpmech/gosl/chk"
)

// ArgToFilename parses an argument as a filename
//  Input:
//   idxArg    -- index of argument; e.g. 0==first, 1==second, etc.
//   fnDefault -- default filename; can be ""
//   ext       -- the file extension to be added; e.g. ".sim"
//   check     -- check for null filename
//  Output:
//   filename -- the filename with extension added
//   fnkey    -- filename key == filename without extension
//  Notes:
//   The first first argument may be a file with extention or not.
//  Examples:
//   If the first argument is "simulation.sim" or "simulation" (with ext=".sim")
//   then the results are: filename="simulation.sim" and fnkey="simulation"
func ArgToFilename(idxArg int, fnDefault, ext string, check bool) (filename, fnkey string) {
	if !flag.Parsed() {
		flag.Parse()
	}
	filename = fnDefault
	if len(flag.Args()) > idxArg {
		filename = flag.Arg(idxArg)
	}
	if FnExt(filename) == "" {
		filename += ext
	}
	fnkey = FnKey(filename)
	if check {
		if filename == "" || fnkey == "" {
			chk.Panic("filename must be given as first argument")
		}
	}
	return
}

// ArgToBool parses an argument as a boolean value
//  Input:
//   idxArg       -- index of argument; e.g. 0==first, 1==second, etc.
//   defaultValue -- default value
func ArgToBool(idxArg int, defaultValue bool) bool {
	if !flag.Parsed() {
		flag.Parse()
	}
	if len(flag.Args()) > idxArg {
		return Atob(flag.Arg(idxArg))
	}
	return defaultValue
}

// ArgToString parses an argument as a string
//  Input:
//   idxArg       -- index of argument; e.g. 0==first, 1==second, etc.
//   defaultValue -- default value
func ArgToString(idxArg int, defaultValue string) string {
	if !flag.Parsed() {
		flag.Parse()
	}
	if len(flag.Args()) > idxArg {
		return flag.Arg(idxArg)
	}
	return defaultValue
}
