// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_cubiceq01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("cubiceq01")

	x1, x2, x3, nx := EqCubicSolveReal(-3, -144, 432)
	io.Pforan("x1 = %v\n", x1)
	io.Pforan("x2 = %v\n", x2)
	io.Pforan("x3 = %v\n", x3)
	io.Pforan("nx = %v\n", nx)
	chk.IntAssert(nx, 3)
	chk.Scalar(tst, "x1", 1e-17, x1, -12)
	chk.Scalar(tst, "x2", 1e-17, x2, 12)
	chk.Scalar(tst, "x3", 1e-14, x3, 3)
}
