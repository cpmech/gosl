// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"strings"

	"github.com/cpmech/gosl/io"
)

func main() {

	// catch errors
	defer func() {
		if err := recover(); err != nil {
			io.PfRed("ERROR: %v\n", err)
		}
	}()

	// input data
	filename, fnkey := io.ArgToFilename(0, "data/equations.txt", "", true)
	io.Pf("\n%s\n", io.ArgsTable("INPUT ARGUMENTS",
		"file with equations", "filename", filename,
	))

	io.Pforan("fnkey = %v\n", fnkey)

	// open file
	f, err := io.OpenFileR(filename)
	if err != nil {
		return
	}

	// constants
	MAXNIDX := 100
	INDICES := []string{"i", "j", "k", "l", "m", "n", "r", "s", "p", "q"}

	// read file
	var buf bytes.Buffer
	io.ReadLinesFile(f, func(idx int, line string) (stop bool) {

		// indices
		l := line
		for i := 0; i < MAXNIDX; i++ {
			l = strings.Replace(l, io.Sf("[%d]", i), io.Sf("_{%d}", i), -1)
		}
		for _, idx := range INDICES {
			l = strings.Replace(l, io.Sf("[%s]", idx), io.Sf("_{%s}", idx), -1)
		}

		// constants
		l = strings.Replace(l, "math.Sqrt2", "\\sqrt{2}", -1)
		l = strings.Replace(l, "SQ2", "\\sqrt{2}", -1)

		// functions
		l = strings.Replace(l, "math.Sqrt", "\\sqrt", -1)
		l = strings.Replace(l, "math.Pow", "\\pow", -1)
		l = strings.Replace(l, "math.Exp", "\\exp", -1)
		l = strings.Replace(l, "math.Sin", "\\sin", -1)
		l = strings.Replace(l, "math.Cos", "\\cos", -1)

		// star
		l = strings.Replace(l, "*", " \\, ", -1)

		// colon-equal
		l = strings.Replace(l, ":=", "=", -1)

		// add to results
		io.Ff(&buf, "%s\n", l)
		return
	})

	// write file
	io.WriteFileVD("/tmp/gosl", fnkey+".tex", &buf)
}
