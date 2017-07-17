// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_strpair01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("strpair01")

	str := "keyA: value :"
	key, val := ExtractStrPair(str, ":")
	Pforan("key = %q\n", key)
	Pforan("val = %q\n", val)
	if key != "keyA" {
		tst.Errorf("key %q is incorrect\n", key)
	}
	if val != "value" {
		tst.Errorf("val %q is incorrect\n", val)
	}
}

func Test_keycode01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("keycode01")

	str := "!typeA:keycodeA !typeB:keycodeB!typeC  : keycodeC"
	if resA, ok := Keycode(str, "typeA"); !ok {
		tst.Errorf("[1;31mcannot find   typeA[0m")
	} else {
		if resA != "keycodeA" {
			tst.Errorf("[1;31mresA != keycodeA[0m")
		}
	}
	if resB, ok := Keycode(str, "typeB"); !ok {
		tst.Errorf("[1;31mcannot find   typeB[0m")
	} else {
		if resB != "keycodeB" {
			tst.Errorf("[1;31mresB != keycodeB[0m")
		}
	}
	if resC, ok := Keycode(str, "typeC"); !ok {
		tst.Errorf("[1;31mcannot find   typeC[0m")
	} else {
		if resC != "keycodeC" {
			tst.Errorf("[1;31mresC != keycodeC[0m")
		}
	}
	if resD, ok := Keycode(str, "typeD"); ok {
		tst.Errorf("[1;31mmust not find typeD[0m")
	} else {
		if resD != "" {
			tst.Errorf("[1;31mresD != \"\"    [0m")
		}
	}

	res, ok := Keycode("", "")
	if res != "" {
		tst.Errorf("[1;31merror when handling empty string[0m")
	}
	if ok {
		tst.Errorf("[1;31merror when handling empty string[0m")
	}

	res, ok = Keycode("!", "")
	if res != "" {
		tst.Errorf("[1;31merror when handling '!' string[0m")
	}
	if ok {
		tst.Errorf("[1;31merror when handling '!' string[0m")
	}

	res, ok = Keycode("!keyA !keyB", "keyA")
	if res != "" {
		tst.Errorf("[1;31merror when handling '!keyA !keyB' string[0m")
	}
	if !ok {
		tst.Errorf("[1;31merror when handling '!keyA !keyB' string[0m")
	}
}

func Test_keycode02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("keycode02")

	str := "!flagA !typeA:keycodeA !typeB:keycodeB!typeC  : keycodeC !FlagB"
	res := Keycodes(str)
	chk.Strings(tst, "keycodes", res, []string{"flagA", "typeA", "typeB", "typeC", "FlagB"})
}

func Test_parsing02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("parsing02")

	keys := []string{"ux", "uy", "uz", "pl", "sl"}
	skeys := JoinKeys(keys)
	kkeys := SplitKeys(skeys)
	Pforan("keys  = %+v\n", keys)
	Pforan("skeys = %q\n", skeys)
	Pforan("kkeys = %+v\n", kkeys)
	chk.Strings(tst, "keys", keys, kkeys)

	k1, k2 := []string{"ex1", "ex2"}, []string{"more1", "more2"}
	allkeys := JoinKeys3(keys, k1, k2, "")
	Pfpink("allkeys = %q\n", allkeys)
	chk.String(tst, allkeys, "ux uy uz pl sl ex1 ex2 more1 more2")

	allkeys = JoinKeys3(keys, k1, k2, ",")
	Pfpink("allkeys = %q\n", allkeys)
	chk.String(tst, allkeys, "ux uy uz pl sl, ex1 ex2, more1 more2")

	r0, r1, r2 := SplitKeys3(allkeys)
	Pfblue2("r0 = %v\n", r0)
	Pfblue2("r1 = %v\n", r1)
	Pfblue2("r2 = %v\n", r2)
	chk.Strings(tst, "r0", r0, keys)
	chk.Strings(tst, "r1", r1, k1)
	chk.Strings(tst, "r2", r2, k2)

	allkeys = JoinKeys3(keys, []string{}, k2, ",")
	Pfpink("allkeys = %q\n", allkeys)
	chk.String(tst, allkeys, "ux uy uz pl sl,, more1 more2")

	allkeys = JoinKeys3(keys, []string{}, []string{}, ",")
	Pfpink("allkeys = %q\n", allkeys)
	chk.String(tst, allkeys, "ux uy uz pl sl,,")

	allkeys = JoinKeys3([]string{}, k1, k2, ",")
	Pfpink("allkeys = %q\n", allkeys)
	chk.String(tst, allkeys, ", ex1 ex2, more1 more2")

	r0, r1, r2 = SplitKeys3("a ,b c,")
	Pfblue2("r0 = %v\n", r0)
	Pfblue2("r1 = %v\n", r1)
	Pfblue2("r2 = %v\n", r2)
	chk.Strings(tst, "r0", r0, []string{"a"})
	chk.Strings(tst, "r1", r1, []string{"b", "c"})
	chk.Strings(tst, "r2", r2, []string{})
}

func Test_parsing03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("parsing03")

	keys := []string{"ux", "uy", "uz", "pl", "sl"}
	k1, k2, k3 := []string{"ex1", "ex2"}, []string{"more1", "more2"}, []string{"a", "b", "c"}
	allkeys := JoinKeys4(keys, k1, k2, k3, "")
	Pfpink("allkeys = %q\n", allkeys)
	chk.String(tst, allkeys, "ux uy uz pl sl ex1 ex2 more1 more2 a b c")

	allkeys = JoinKeys4(keys, k1, k2, k3, ",")
	Pfpink("allkeys = %q\n", allkeys)
	chk.String(tst, allkeys, "ux uy uz pl sl, ex1 ex2, more1 more2, a b c")

	r0, r1, r2, r3 := SplitKeys4(allkeys)
	Pfblue2("r0 = %v\n", r0)
	Pfblue2("r1 = %v\n", r1)
	Pfblue2("r2 = %v\n", r2)
	Pfblue2("r3 = %v\n", r3)
	chk.Strings(tst, "r0", r0, keys)
	chk.Strings(tst, "r1", r1, k1)
	chk.Strings(tst, "r2", r2, k2)
	chk.Strings(tst, "r3", r3, k3)

	allkeys = JoinKeys4(keys, []string{}, k2, []string{}, ",")
	Pfpink("allkeys = %q\n", allkeys)
	chk.String(tst, allkeys, "ux uy uz pl sl,, more1 more2,")
}

func Test_parsing04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("parsing04")

	keys := []string{"ux", "uy", "uz"}
	skeys := JoinKeysPre("R", keys)
	Pforan("keys  = %+v\n", keys)
	Pforan("skeys = %q\n", skeys)
	chk.String(tst, skeys, "Rux Ruy Ruz")
}

func Test_parsing05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("parsing05. split by spaces quoted")

	str := " arg1 arg2 '   hello    world' '' \"  another hello world   \" "
	res := SplitSpacesQuoted(str)
	Pforan("res = %q\n", res)
	chk.String(tst, res[0], "arg1")
	chk.String(tst, res[1], "arg2")
	chk.String(tst, res[2], "'   hello    world'")
	chk.String(tst, res[3], "''")
	chk.String(tst, res[4], "\"  another hello world   \"")
}

func Test_parsing06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("parsing06. extract within brackets")

	str := "(arg1, (arg2.1, arg2.2),  arg3, arg4, (arg5.1,arg5.2,  arg5.3 ) )"
	res := SplitWithinParentheses(str)
	Pforan("%q\n", str)
	for _, r := range res {
		Pf("%q\n", r)
	}
	chk.String(tst, res[0], "arg1")
	chk.String(tst, res[1], "arg2.1  arg2.2")
	chk.String(tst, res[2], "arg3")
	chk.String(tst, res[3], "arg4")
	chk.String(tst, res[4], "arg5.1 arg5.2   arg5.3 ")
}

func Test_parsing07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("parsing07. float64 fields")

	str := "1.0   1111.11 2.0   3.0"
	res := SplitFloats(str)
	Pforan("floats = %v\n", res)
	chk.Array(tst, "res", 1e-17, res, []float64{1, 1111.11, 2, 3})

	str = "1   1111 2  3"
	ints := SplitInts(str)
	Pforan("ints = %v\n", ints)
	chk.Ints(tst, "res", ints, []int{1, 1111, 2, 3})
}
