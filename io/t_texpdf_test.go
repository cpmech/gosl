// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_texreport01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("texreport01")

	l1 := TexNum("", 123.456, true)
	l2 := TexNum("", 123.456e-8, true)
	l3 := TexNum("%.1e", 123.456e+8, true)
	l4 := TexNum("%.2e", 123.456e-8, false)
	Pforan("l1 = %v\n", l1)
	Pforan("l2 = %v\n", l2)
	Pforan("l3 = %v\n", l3)
	Pforan("l4 = %v\n", l4)
	chk.String(tst, l1, "123.456")
	chk.String(tst, l2, "1.23456\\cdot 10^{-6}")
	chk.String(tst, l3, "1.2\\cdot 10^{10}")
	chk.String(tst, l4, "1.23e-06")
}
