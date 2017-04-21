// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func main() {

	// a hypotetical matrix obtained from a "numerical" solution
	// notice "noise" on the last component
	Anumerical := [][]float64{
		{76, 35, 64, 8, 35},
		{92, 37, 16, 0, 46},
		{64, 36, 60, 27, 42},
		{41, 20, 21, 35, 45},
		{4, 48, 37, 87, 9 + 1e-9},
	}

	// a hypotetical matrix obtained from a "closed-form" solution
	Aanalytical := [][]float64{
		{76, 35, 64, 8, 35},
		{92, 37, 16, 0, 46},
		{64, 36, 60, 27, 42},
		{41, 20, 21, 35, 45},
		{4, 48, 37, 87, 9},
	}

	// allocate testing structure, just for this example
	tst := &testing.T{}

	// tolerance for comparison
	//tolerance := 1e-10 // this makes the test to fail
	tolerance := 1e-8 // this allows test to pass

	// compare matrices
	chk.Matrix(tst, "A", tolerance, Anumerical, Aanalytical)

	// note that this is not the a common way of using testing.T
	// usually, the "tst" variable comes from a unit test
	if tst.Failed() {
		io.PfRed("test failed\n")
	} else {
		io.PfGreen("OK\n")
	}
}
