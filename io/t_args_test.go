// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_args01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("args01")

	fn, fnk := ArgToFilename(0, "simulation", ".sim", true)
	Pforan("fn  = %v\n", fn)
	Pforan("fnk = %v\n", fnk)
	chk.String(tst, fn, "simulation.sim")
	chk.String(tst, fnk, "simulation")

	resFloat := ArgToFloat(1, 456)
	chk.Float64(tst, "456", 1e-17, resFloat, 456)

	resInt := ArgToInt(1, 123)
	if resInt != 123 {
		tst.Errorf("test failed: resInt != 123\n")
		return
	}

	resBool := ArgToBool(1, true)
	if !resBool {
		tst.Errorf("test failed: resBool != true\n")
		return
	}

	resString := ArgToString(1, "myname")
	chk.String(tst, resString, "myname")

	tab := ArgsTable(
		"1 2 3 INPUT PARAMETERS 3 2 1",
		"first argument", "first", true,
		"second argument", "second", "string",
		"third argument", "third", 123,
		"fourth argument", "fourth", 666.0,
	)

	Pf("\n%v\n", tab)

	chk.String(tst, tab,
		`   1 2 3 INPUT PARAMETERS 3 2 1
===================================
     description      key     value
-----------------------------------
  first argument    first      true
 second argument   second    string
  third argument    third       123
 fourth argument   fourth       666
===================================
`)
}
