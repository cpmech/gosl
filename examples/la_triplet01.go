// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func main() {

	// generate the following matrix to view with VisMatrix
	//   0 2 0 0
	//   1 0 4 0
	//   0 0 0 5
	//   0 3 0 6
	a := new(la.Triplet)
	a.Init(4, 4, 6)
	a.Put(1, 0, 1)
	a.Put(0, 1, 2)
	a.Put(3, 1, 3)
	a.Put(1, 2, 4)
	a.Put(2, 3, 5)
	a.Put(3, 3, 6)

	a.WriteSmat("/tmp/gosl", "triplet01", 0)
	io.Pf("Now, run:\nvismatrix /tmp/gosl/triplet01.smat\n")
}
