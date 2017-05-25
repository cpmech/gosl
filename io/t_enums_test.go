// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

var (
	ALPHA = NewEnum("α")
	BETA  = NewEnum("β", "rnd")
	GAMMA = NewEnum("γ", "rnd", "G")
	OMEGA = NewEnum("ω", "rnd", "O", "the omega")
)

func Test_enums01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("enums01")

	names := []string{"α", "β", "γ", "ω"}
	pfixs := []string{"", "rnd", "rnd", "rnd"}
	descs := []string{"", "", "", "the omega"}
	keys := []string{"", "", "G", "O"}
	for i, enum := range enums {
		Pf("enum = %q\n", enum)
		chk.String(tst, enum.Name, names[i])
		chk.String(tst, enum.Prefix, pfixs[i])
		chk.String(tst, enum.Desc, descs[i])
		chk.String(tst, enum.Key, keys[i])
	}

	if ALPHA != 0 {
		tst.Errorf("ALPHA != 0\n")
		return
	}
	if BETA != 1 {
		tst.Errorf("BETA != 1\n")
		return
	}
	if GAMMA != 2 {
		tst.Errorf("GAMMA != 2\n")
		return
	}
	if OMEGA != 3 {
		tst.Errorf("OMEGA != 3\n")
		return
	}

	for key, item := range enumsMap {
		Pforan("key=%q item=%+v\n", key, item)
	}
	enum := EnumsFind("", "α")
	if enum != ALPHA {
		tst.Errorf("cannot find ALPHA\n")
		return
	}
	enum = EnumsFind("", "another")
	if enum != -1 {
		tst.Errorf("EnumsFind(another) should return -1\n")
		return
	}
}
