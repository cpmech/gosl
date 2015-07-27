// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"flag"

	"github.com/cpmech/gosl/chk"
)

// Args0toFilename parses the first argument as a filename
//  Input:
//   fnDefault -- default filename; can be ""
//   ext       -- the file extension to be added
//   check     -- check for null filename
//  Output:
//   filename -- the filename with extension added
//   fnkey    -- filename key == filename without extension
//  Notes:
//   The first first argument may be a file with extention or not.
//  Examples:
//   If the first argument is "simulation.sim" or "simulation" (with ext="sim")
//   then the results are: filename="simulation.sim" and fnkey="simulation"
func Args0toFilename(fnDefault, ext string, check bool) (filename, fnkey string) {
	if !flag.Parsed() {
		flag.Parse()
	}
	filename = fnDefault
	if len(flag.Args()) > 0 {
		filename = flag.Arg(0)
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
