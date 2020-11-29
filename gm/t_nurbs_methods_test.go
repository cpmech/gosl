// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestNurbsMethods01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("NurbsMethods01")

	nrb := FactoryNurbs.Surf2dExample1()
	chk.Int(tst, "Gnd()", nrb.Gnd(), 2)
	chk.Int(tst, "Ord(0)", nrb.Ord(0), 2)
	chk.Int(tst, "Ord(1)", nrb.Ord(1), 1)
	chk.Float64(tst, "U(0,0)", 1e-15, nrb.U(0, 0), 0)
	chk.Float64(tst, "Udelta(0)", 1e-15, nrb.Udelta(0), 3)
	chk.Float64(tst, "Udelta(1)", 1e-15, nrb.Udelta(1), 1)
	chk.Float64(tst, "UfromR(0,0.0)", 1e-15, nrb.UfromR(0, 0.0), 1.5)
	chk.Float64(tst, "UfromR(1,0.0)", 1e-15, nrb.UfromR(1, 0.0), 0.5)
	chk.Int(tst, "NumBasis(0)", nrb.NumBasis(0), 5)
	chk.Int(tst, "NumBasis(1)", nrb.NumBasis(1), 2)
	chk.IntDeep2(tst, "NonZeroSpans(0)", nrb.NonZeroSpans(0), [][]int{{2, 3}, {3, 4}, {4, 5}})
	chk.IntDeep2(tst, "NonZeroSpans(1)", nrb.NonZeroSpans(1), [][]int{{1, 2}})
}
