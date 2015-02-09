// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func Test_basic01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("basic01")

	v0 := Atob("1")
	if !v0 {
		Panic("Atob failed: v0")
	}

	v1 := Atob("true")
	if !v1 {
		Panic("Atob failed: v1")
	}

	v2 := Atob("0")
	if v2 {
		Panic("Atob failed: v2")
	}

	v3 := Itob(0)
	if v3 {
		Panic("Itob failed: v3")
	}

	v4 := Itob(-1)
	if !v4 {
		Panic("Itob failed: v4")
	}

	v5 := Btoi(true)
	if v5 != 1 {
		Panic("Btoi failed: v5")
	}

	v6 := Btoi(false)
	if v6 != 0 {
		Panic("Btoi failed: v6")
	}

	v7 := Btoa(true)
	if v7 != "true" {
		Panic("Btoa failed: v7")
	}

	v8 := Btoa(false)
	if v8 != "false" {
		Panic("Btoa failed: v8")
	}
}

func Test_basic02(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("basic02")

	v := []int{1, 4, 5, 1, 2, 8, 3}
	var b bytes.Buffer
	WriteIntSlice(&b, v)

	bb, err := ReadFile("data/basic02.dat")
	if err != nil {
		Panic("cannot read data/basic02.dat")
	}
	bb = bb[:len(bb)-1] // remove EOF

	Pf("b  = %v\n", b.Bytes())
	Pf("bb = %v\n", bb)

	if !bytes.Equal(b.Bytes(), bb) {
		Panic("WriteIntSlice failed")
	}
}

func Test_json01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("json01")

	m := map[int][]int{
		-1: []int{1, 4, 5, 1, 2, 8, 3},
		-2: []int{4, 4, 9, 1, 1, 2, 8},
		-3: []int{6, 7, 4, 0, 3, 5, 3},
		-4: []int{0, 1, 2, 3, 4, 9, 0},
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "{\n")
	JsonMapIntSlice(&b, "  ", "mymap", &m)
	fmt.Fprintf(&b, "\n}")
	Pfpink("%v\n", string(b.Bytes()))

	var dict interface{}
	err := json.Unmarshal(b.Bytes(), &dict)
	if err != nil {
		Panic("cannot unmarshal b")
	}
	Pforan("dict = %v\n", dict)

	d := dict.(map[string]interface{})
	for k, v := range d {
		Pfpink("k=%v  v=%+#v\n", k, v)
		switch vv := v.(type) {
		case map[string]interface{}:
			for key, val := range vv {
				Pfblue2("key=%v  val=%+#v\n", key, val)
				var ints []int
				switch vval := val.(type) {
				case []interface{}:
					Pfgrey("  vval = %+#v\n", vval)
					for i, x := range vval {
						switch xx := x.(type) {
						case int:
							Pfgrey2("    i=%v x=%v (int)\n", i, xx)
						case float64:
							Pfgrey2("    i=%v x=%v (float64)\n", i, xx)
							ints = append(ints, int(xx))
						default:
							Panic("cannot find type of %+#v", x)
						}
					}
				default:
					Panic("cannot find type of %+#v", val)
				}
				Pfgreen("ints = %v\n", ints)
				cor, ok := m[Atoi(key)]
				if !ok {
					Panic("key==%v is not in map")
				}
				CompareInts(tst, "ints", ints, cor)
			}
		default:
			Panic("cannot find type of %+#v", v)
		}
	}
}

func do_check_cells_map(tst *testing.T, cells map[int][]int) {
	correct := map[int][]int{
		-1: []int{0, 1, 3, 4, 6, 7, 9, 10, 12, 13},
		-2: []int{2, 5, 8, 11, 14},
	}
	if len(cells) != len(correct) {
		Panic("reading 'cells' failed")
	}
	for k, v := range cells {
		CompareInts(tst, "cells => correct", correct[k], v)
	}
	for k, v := range correct {
		CompareInts(tst, "correct => cells", cells[k], v)
	}
}

func Test_json02(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("json02")

	fn := "data/json01.dat"
	b, err := ReadFile(fn)
	if err != nil {
		Panic("Cannot read file <%s>: %v", fn, err)
	}

	type Data_t struct {
		Cells interface{}
	}
	var d Data_t
	err = json.Unmarshal(b, &d)
	if err != nil {
		Panic("Cannot unmarshal file <%s>: %v", fn, err)
	}

	Pfblue2("d = %+#v\n", d)
	cells := GetMapIntSlice(d.Cells)
	Pforan("cells = %+#v\n", cells)
	do_check_cells_map(tst, cells)
}

func Test_json03(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("json03")

	// using map[string]interface{}
	var res map[string]interface{}
	OpenAndParseJson(&res, "data/json01", ".dat")
	Pforan("res = %v\n", res)
	cells := GetMapIntSlice(res["cells"])
	do_check_cells_map(tst, cells)

	// using Data_t
	type Data_t struct {
		Name string
		Id   int
		Ok   bool
		Nums []float64
	}
	var dat Data_t
	OpenAndParseJson(&dat, "data/json02", ".dat")
	Pforan("dat = %v\n", dat)
	CheckString(tst, dat.Name, "Dorival Pedroso")
	if dat.Id != 666 {
		Panic("dat.Id failed")
	}
	if !dat.Ok {
		Panic("dat.Ok failed")
	}
	CompareDbls(tst, "Nums", dat.Nums, []float64{1.0, 5.0, 11.0, 21.0, 666.66})
}
