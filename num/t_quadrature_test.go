// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	. "math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestQuadGen01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("QuadGen01. using QUADPACK general function")

	f := func(x float64) float64 { return Sqrt(1.0 + Pow(Sin(x), 3.0)) }
	A, err := QuadGen(0, 1, 0, f)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	io.Pforan("A  = %v\n", A)
	chk.Scalar(tst, "A", 1e-12, A, 1.08268158558)
}

func TestQuadCs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("QuadCs01. using QUADPACK oscillatory function")

	ω := Pow(2.0, 3.4)
	f := func(x float64) float64 { return Exp(20.0 * (x - 1)) }
	A, err := QuadCs(0, 1, ω, true, 0, f)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	io.Pforan("A  = %v\n", A)
	Aref := (20*Sin(ω) - ω*Cos(ω) + ω*Exp(-20)) / (Pow(20, 2) + Pow(ω, 2))
	chk.Scalar(tst, "A", 1e-16, A, Aref)
}
