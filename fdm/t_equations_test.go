// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestEquations01(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	chk.PrintTitle("TestEquations 01")

	var e Equations
	e.Init(9, []int{0, 6, 3})
	io.Pf("%v\n", e)
	chk.Ints(tst, "RF1", e.RF1, []int{1, 2, 4, 5, 7, 8})
	chk.Ints(tst, "FR1", e.FR1, []int{-1, 0, 1, -1, 2, 3, -1, 4, 5})
	chk.Ints(tst, "RF2", e.RF2, []int{0, 3, 6})
	chk.Ints(tst, "FR2", e.FR2, []int{0, -1, -1, 1, -1, -1, 2, -1, -1})

	e.Init(9, []int{0, 1, 2})
	io.Pf("%v\n", e)
	chk.Ints(tst, "RF1", e.RF1, []int{3, 4, 5, 6, 7, 8})
	chk.Ints(tst, "FR1", e.FR1, []int{-1, -1, -1, 0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "RF2", e.RF2, []int{0, 1, 2})
	chk.Ints(tst, "FR2", e.FR2, []int{0, 1, 2, -1, -1, -1, -1, -1, -1})
}
