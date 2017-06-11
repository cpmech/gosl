// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func TestFourierTrans01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierTrans01. FFT")

	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	err := Dft1d(x, false)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	X := la.RCpairsToComplex(x)
	io.Pf("X = %v\n", X)

	y := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}
	Y := dft1dslow(y)
	io.Pf("Y = %v\n", Y)

	chk.VectorC(tst, "X", 1e-14, Y, X)
}
