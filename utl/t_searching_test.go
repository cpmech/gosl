// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"gosl/chk"
	"gosl/io"
)

func Test_search01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("search01")

	a := []string{"66", "644", "666", "653", "10", "0", "1", "1", "1"}
	idx := StrIndexSmall(a, "666")
	io.Pf("a = %v\n", a)
	io.Pf("idx of '666' = %v\n", idx)
	if idx != 2 {
		tst.Errorf("idx is wrong")
		return
	}
	idx = StrIndexSmall(a, "1")
	io.Pf("idx of '1'   = %v\n", idx)
	if idx != 6 {
		tst.Errorf("idx is wrong")
		return
	}

	b := []int{66, 644, 666, 653, 10, 0, 1, 1, 1}
	idx = IntIndexSmall(b, 666)
	io.Pf("b = %v\n", b)
	io.Pf("idx of 666 = %v\n", idx)
	if idx != 2 {
		tst.Errorf("idx is wrong")
		return
	}
	idx = IntIndexSmall(b, 1)
	io.Pf("idx of 1   = %v\n", idx)
	if idx != 6 {
		tst.Errorf("idx is wrong")
		return
	}

	// TODO: finish this
	//jdx, bsorted := IntIndexSorted(b, 666, nil)
	//io.Pforan("\njdx = %v\n", jdx)
	//io.Pforan("bsorted = %v\n", bsorted)
	//chk.Ints(tst, "bsorted", bsorted, []int{0, 1, 1, 1, 10, 66, 644, 653, 666})
}
