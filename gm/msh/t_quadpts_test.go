// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_quadpts01(tst *testing.T) {

	// TODO: implement this test

	//verbose()
	chk.PrintTitle("quadpts01. quadrature points")

	for name, pts := range IntPoints {

		io.Pfyel("--------------------------------- %-6s---------------------------------\n", name)

		for n, p := range pts {
			io.Pforan("%2d: %v\n\n", n, p)
		}
	}
}
