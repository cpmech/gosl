// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_hc01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("hc01. hypercube")

	Init(111)

	x := LatinIHS(2, 10, 5)
	io.Pforan("x = %v\n", x)

	xcor := [][]int{
		{2, 9, 5, 8, 1, 4, 3, 10, 7, 6},
		{3, 10, 1, 2, 7, 4, 9, 6, 5, 8},
	}
	chk.IntMat(tst, "x", x, xcor)
}
