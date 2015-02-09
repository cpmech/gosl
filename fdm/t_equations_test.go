// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
	"testing"

	"github.com/cpmech/gosl/utl"
)

func TestEquations01(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	utl.TTitle("TestEquations 01")

	var e Equations
	e.Init(9, []int{0, 6, 3})
	utl.Pf("%v\n", e)
	utl.CompareInts(tst, "RF1", e.RF1, []int{1, 2, 4, 5, 7, 8})
	utl.CompareInts(tst, "FR1", e.FR1, []int{-1, 0, 1, -1, 2, 3, -1, 4, 5})
	utl.CompareInts(tst, "RF2", e.RF2, []int{0, 3, 6})
	utl.CompareInts(tst, "FR2", e.FR2, []int{0, -1, -1, 1, -1, -1, 2, -1, -1})

	e.Init(9, []int{0, 1, 2})
	utl.Pf("%v\n", e)
	utl.CompareInts(tst, "RF1", e.RF1, []int{3, 4, 5, 6, 7, 8})
	utl.CompareInts(tst, "FR1", e.FR1, []int{-1, -1, -1, 0, 1, 2, 3, 4, 5})
	utl.CompareInts(tst, "RF2", e.RF2, []int{0, 1, 2})
	utl.CompareInts(tst, "FR2", e.FR2, []int{0, 1, 2, -1, -1, -1, -1, -1, -1})
}
