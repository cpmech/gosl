// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_constants(tst *testing.T) {
	//verbose()
	chk.PrintTitle("constants. sqrt(2). sqrt(3)")
	io.Pf("SQ2     = %.16f\n", SQ2)
	io.Pf("sqrt(2) = %.16f\n", math.Sqrt(2.0))
	io.Pf("\n")
	io.Pf("SQ3     = %.16f\n", SQ3)
	io.Pf("sqrt(3) = %.16f\n", math.Sqrt(3.0))
	io.Pf("\n")
	io.Pf("SQ5     = %.16f\n", SQ5)
	io.Pf("sqrt(5) = %.16f\n", math.Sqrt(5.0))
	io.Pf("\n")
	io.Pf("SQ6     = %.16f\n", SQ6)
	io.Pf("sqrt(6) = %.16f\n", math.Sqrt(6.0))
	io.Pf("\n")
	io.Pf("SQ7     = %.16f\n", SQ7)
	io.Pf("sqrt(7) = %.16f\n", math.Sqrt(7.0))
	io.Pf("\n")
	io.Pf("SQ8     = %.16f\n", SQ8)
	io.Pf("sqrt(8) = %.16f\n", math.Sqrt(8.0))
	io.Pf("\n")
	chk.Float64(tst, "sqrt(2)", 1e-15, SQ2, math.Sqrt(2.0))
	chk.Float64(tst, "sqrt(3)", 1e-15, SQ3, math.Sqrt(3.0))
	chk.Float64(tst, "sqrt(5)", 1e-15, SQ5, math.Sqrt(5.0))
	chk.Float64(tst, "sqrt(6)", 1e-15, SQ6, math.Sqrt(6.0))
	chk.Float64(tst, "sqrt(7)", 1e-15, SQ7, math.Sqrt(7.0))
	chk.Float64(tst, "sqrt(8)", 1e-15, SQ8, math.Sqrt(8.0))
}
