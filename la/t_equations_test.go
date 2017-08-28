// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestEqs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs01")

	// some prescribed
	var e Equations
	e.Init(9, []int{0, 6, 3})
	e.Stat(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{1, 2, 4, 5, 7, 8})
	chk.Ints(tst, "FtoU", e.FtoU, []int{-1, 0, 1, -1, 2, 3, -1, 4, 5})
	chk.Ints(tst, "KtoF", e.KtoF, []int{0, 3, 6})
	chk.Ints(tst, "FtoK", e.FtoK, []int{0, -1, -1, 1, -1, -1, 2, -1, -1})

	// some prescribed
	io.Pl()
	e.Init(9, []int{0, 2, 1})
	e.Stat(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{3, 4, 5, 6, 7, 8})
	chk.Ints(tst, "FtoU", e.FtoU, []int{-1, -1, -1, 0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "KtoF", e.KtoF, []int{0, 1, 2})
	chk.Ints(tst, "FtoK", e.FtoK, []int{0, 1, 2, -1, -1, -1, -1, -1, -1})

	// none prescribed
	io.Pl()
	e.Init(5, nil)
	e.Stat(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "FtoU", e.FtoU, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "KtoF", e.KtoF, nil)
	chk.Ints(tst, "FtoK", e.FtoK, []int{-1, -1, -1, -1, -1})

	// all prescribed
	io.Pl()
	e.Init(5, []int{0, 1, 2, 3, 4})
	e.Stat(true)
	chk.Ints(tst, "UtoF", e.UtoF, nil)
	chk.Ints(tst, "FtoU", e.FtoU, []int{-1, -1, -1, -1, -1})
	chk.Ints(tst, "KtoF", e.KtoF, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "FtoK", e.FtoK, []int{0, 1, 2, 3, 4})
}
