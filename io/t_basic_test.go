// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_basic01(tst *testing.T) {

	//chk.Verbose = true
	chk.PrintTitle("basic01")

	v0 := Atob("1")
	if !v0 {
		chk.Panic("Atob failed: v0")
	}

	v1 := Atob("true")
	if !v1 {
		chk.Panic("Atob failed: v1")
	}

	v2 := Atob("0")
	if v2 {
		chk.Panic("Atob failed: v2")
	}

	v3 := Itob(0)
	if v3 {
		chk.Panic("Itob failed: v3")
	}

	v4 := Itob(-1)
	if !v4 {
		chk.Panic("Itob failed: v4")
	}

	v5 := Btoi(true)
	if v5 != 1 {
		chk.Panic("Btoi failed: v5")
	}

	v6 := Btoi(false)
	if v6 != 0 {
		chk.Panic("Btoi failed: v6")
	}

	v7 := Btoa(true)
	if v7 != "true" {
		chk.Panic("Btoa failed: v7")
	}

	v8 := Btoa(false)
	if v8 != "false" {
		chk.Panic("Btoa failed: v8")
	}
}
