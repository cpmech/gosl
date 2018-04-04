// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestParams01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params01")

	var params Params
	params = []*P{
		{N: "a", V: 1.0},
		{N: "b", V: 2.0},
		{N: "c", V: 3.0},
	}
	io.Pforan("%v\n", params)

	var a, b, c, A, B float64
	e := params.Connect(&a, "a", "test call")
	e += params.Connect(&b, "b", "test call")
	e += params.Connect(&c, "c", "test call")
	e += params.Connect(&A, "a", "test call")
	e += params.Connect(&B, "b", "test call")
	if e != "" {
		tst.Errorf("connect failed: %v\n", e)
		return
	}

	chk.Float64(tst, "a", 1e-15, a, 1.0)
	chk.Float64(tst, "b", 1e-15, b, 2.0)
	chk.Float64(tst, "c", 1e-15, c, 3.0)
	chk.Float64(tst, "A", 1e-15, A, 1.0)
	chk.Float64(tst, "B", 1e-15, B, 2.0)

	prm := params.Find("a")
	if prm == nil {
		tst.Error("cannot find parameter 'a'\n")
		return
	}
	prm.Set(123)
	chk.Float64(tst, "a", 1e-15, a, 123)
	chk.Float64(tst, "A", 1e-15, A, 123)
}

func TestParams02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params02")

	var params Params
	params = []*P{
		{N: "klx", V: 1.0},
		{N: "kly", V: 2.0},
		{N: "klz", V: 3.0},
	}
	io.Pf("%v\n", params)

	var klx, kly, klz float64
	errMsg := params.ConnectSet([]*float64{&klx, &kly, &klz}, []string{"klx", "kly", "klz"}, "TestParams02")
	if errMsg != "" {
		tst.Errorf("connect set failed: %v\n", errMsg)
		return
	}

	chk.Float64(tst, "klx", 1e-15, klx, 1.0)
	chk.Float64(tst, "kly", 1e-15, kly, 2.0)
	chk.Float64(tst, "klz", 1e-15, klz, 3.0)

	params[1].Set(2.2)
	chk.Float64(tst, "kly", 1e-15, kly, 2.2)

	var dummy float64
	errMsg = params.ConnectSetOpt(
		[]*float64{&klx, &kly, &dummy},
		[]string{"klx", "kly", "dummy"},
		[]bool{false, false, true},
		"TestParams02",
	)
	if errMsg != "" {
		tst.Errorf("connect set failed: %v\n", errMsg)
		return
	}

	errMsg = params.ConnectSetOpt(
		[]*float64{&klx, &kly, &dummy},
		[]string{"klx", "kly", "dummy"},
		[]bool{false, false, false},
		"TestParams02",
	)
	if errMsg == "" {
		tst.Errorf("errMsg should not be empty\n")
		return
	}
}

func TestParams03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params03")

	var params Params
	params = []*P{
		{N: "klx", V: 1.0},
		{N: "kly", V: 2.0},
		{N: "klz", V: 3.0},
	}
	io.Pforan("%v\n", params)

	values, found := params.GetValues([]string{"klx", "kly", "klz"})
	chk.Array(tst, "values", 1e-15, values, []float64{1, 2, 3})
	chk.Bools(tst, "found", found, []bool{true, true, true})

	values, found = params.GetValues([]string{"klx", "klY", "klz"})
	chk.Array(tst, "values", 1e-15, values, []float64{1, 0, 3})
	chk.Bools(tst, "found", found, []bool{true, false, true})
}

func TestParams04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params04")

	var params Params
	params = []*P{
		{N: "a", V: 123, Min: 0, Max: math.MaxFloat64},
		{N: "b", V: 456, Min: 0, Max: math.MaxFloat64},
	}

	params.CheckLimits()

	values := params.CheckAndGetValues([]string{"a", "b"})
	chk.Array(tst, "values", 1e-15, values, []float64{123, 456})

	var a, b float64
	params.CheckAndSetVars([]string{"a", "b"}, []*float64{&a, &b})
	chk.Float64(tst, "a", 1e-15, a, 123)
	chk.Float64(tst, "b", 1e-15, b, 456)
}
