// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_bestsq01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bestsq01")

	for i := 1; i <= 12; i++ {
		nrow, ncol := BestSquare(i)
		io.Pforan("nrow, ncol, nrow*ncol = %2d, %2d, %2d\n", nrow, ncol, nrow*ncol)
		if nrow*ncol != i {
			chk.Panic("BestSquare failed")
		}
	}
}
