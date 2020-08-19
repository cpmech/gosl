// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
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

	params.SetValue("klx", 0.001)
	chk.Float64(tst, "klx", 1e-15, params.GetValue("klx"), 0.001)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.SetValue("invalid", 666)
}

func TestParams04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params04")

	params := NewParams(
		&P{N: "a", V: 123, Min: 0, Max: math.MaxFloat64},
		&P{N: "b", V: 456, Min: 0, Max: math.MaxFloat64},
	)

	params.CheckLimits()

	values := params.CheckAndGetValues([]string{"a", "b"})
	chk.Array(tst, "values", 1e-15, values, []float64{123, 456})

	var a, b float64
	params.CheckAndSetVariables([]string{"a", "b"}, []*float64{&a, &b})
	chk.Float64(tst, "a", 1e-15, a, 123)
	chk.Float64(tst, "b", 1e-15, b, 456)

	params.SetBool("b", -100)
	chk.Float64(tst, "b-", 1e-15, params.GetValue("b"), -1.0)

	params.SetBool("b", +100)
	chk.Float64(tst, "b+", 1e-15, params.GetValue("b"), +1.0)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.SetBool("invalid", 666)
}

func TestParams05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params05")

	params := NewParams(
		&P{N: "P0", V: 0},
		&P{N: "P1", V: 1},
		&P{N: "P2", V: 2},
		&P{N: "P3", V: 3},
	)
	io.Pf("%v\n", params)

	values, found := params.GetValues([]string{"P1", "P2", "P3", "P4"})
	chk.Array(tst, "values", 1e-15, values, []float64{1, 2, 3, 0})
	chk.Bools(tst, "found", found, []bool{true, true, true, false})

	res := params.GetValue("P1")
	chk.Float64(tst, "P1", 1e-15, res, 1.0)

	res = params.GetValue("P2")
	chk.Float64(tst, "2", 1e-15, res, 2.0)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	res = params.GetValue("invalid")
}

func TestParams06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params06")

	params := NewParams(
		&P{N: "P0", V: 0},
		&P{N: "P1", V: 1},
		&P{N: "P2", V: 2},
		&P{N: "P3", V: 3},
	)
	io.Pf("%v\n", params)

	values, found := params.GetValues([]string{"P1", "P2", "P3", "P4"})
	chk.Array(tst, "values", 1e-15, values, []float64{1, 2, 3, 0})
	chk.Bools(tst, "found", found, []bool{true, true, true, false})

	res := params.GetBool("P1")
	if res == false {
		tst.Error("res should be true\n")
		return
	}

	res = params.GetBool("P0")
	if res == true {
		tst.Error("res should be false\n")
		return
	}

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	res = params.GetBool("invalid")
}

func TestParams07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params07. String and CheckLimits")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
		&P{N: "b", V: 1, Min: -2, Max: +2},
	)
	chk.String(tst, params.String(), `{"n":"a", "v":1, "min":-2, "max":2, "s":0, "d":"", "u":"", "adj":0, "dep":0, "extra":"", "inact":false, "setdef":false},
{"n":"b", "v":1, "min":-2, "max":2, "s":0, "d":"", "u":"", "adj":0, "dep":0, "extra":"", "inact":false, "setdef":false}`)

	params.CheckLimits()

	params.SetValue("a", -10) // out of range
	res := params.GetValue("a")
	chk.Float64(tst, "a", 1e-15, res, -10)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckLimits()
}

func TestParams08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params08. CheckLimits")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
	)

	params.CheckLimits()

	params.SetValue("a", +10) // out of range
	res := params.GetValue("a")
	chk.Float64(tst, "a", 1e-15, res, +10)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckLimits()
}

func TestParams09(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params09. CheckAndGetValues. Panic # 1")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
	)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckAndGetValues([]string{"invalid"})
}

func TestParams10(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params09. CheckAndGetValues. Panic # 2")

	params := NewParams(
		&P{N: "a", V: -10, Min: -2, Max: +2},
	)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckAndGetValues([]string{"a"})
}

func TestParams11(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params09. CheckAndGetValues. Panic # 3")

	params := NewParams(
		&P{N: "a", V: +10, Min: -2, Max: +2},
	)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckAndGetValues([]string{"a"})
}

func TestParams12(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params12. CheckAndSetVariables. Panic # 1")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
	)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckAndSetVariables([]string{"a"}, nil)
}

func TestParams13(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params13. CheckAndSetVariables. Panic # 2")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
	)
	a := 0.0

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckAndSetVariables([]string{"invalid"}, []*float64{&a})
}

func TestParams14(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params14. CheckAndSetVariables. Panic # 3")

	params := NewParams(
		&P{N: "a", V: -10, Min: -2, Max: +2},
	)
	a := 0.0

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckAndSetVariables([]string{"a"}, []*float64{&a})
}

func TestParams15(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params15. CheckAndSetVariables. Panic # 4")

	params := NewParams(
		&P{N: "a", V: +10, Min: -2, Max: +2},
	)
	a := 0.0

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckAndSetVariables([]string{"a"}, []*float64{&a})
}

func TestParams16(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params16. CheckAndSetVariables. Panic # 5")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
	)

	io.Pf("\n>>> the following Panic is OK <<<\n")
	defer chk.RecoverTstPanicIsOK(tst)
	params.CheckAndSetVariables([]string{"a"}, []*float64{nil})
}

func TestParams17(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params17. ConnectSet. Error")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
	)

	a := 0.0
	vars := []*float64{&a}
	names := []string{"b"}

	io.Pf("\n>>> the following Panic is OK <<<\n")
	errmsg := params.ConnectSet(vars, names, "TestParams17")
	chk.String(tst, errmsg, `cannot find parameter named "b" as requested by "TestParams17"
`)
}

func TestParams18(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params18. GetValueOrDefault")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
	)

	res := params.GetValueOrDefault("a", -2)
	chk.Float64(tst, "a", 1e-15, res, +1)

	res = params.GetValueOrDefault("invalid", -2)
	chk.Float64(tst, "a", 1e-15, res, -2)
}

func TestParams19(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params19. GetBoolOrDefault")

	params := NewParams(
		&P{N: "a", V: 0, Min: -2, Max: +2}, // V=0 => false
		&P{N: "b", V: 1, Min: -2, Max: +2}, // V=1 => true
	)

	res := params.GetBoolOrDefault("a", true)
	if res != false {
		tst.Errorf("res(a) should be false\n")
		return
	}

	res = params.GetBoolOrDefault("b", false)
	if res != true {
		tst.Errorf("res(b) should be true\n")
		return
	}

	res = params.GetBoolOrDefault("invalid", true)
	if res != true {
		tst.Errorf("res(invalid) should be true\n")
		return
	}
}

func TestParams20(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Params20. GetIntOrDefault")

	params := NewParams(
		&P{N: "a", V: 1, Min: -2, Max: +2},
	)

	res := params.GetIntOrDefault("a", -2)
	chk.Int(tst, "a", res, +1)

	res = params.GetIntOrDefault("invalid", -2)
	chk.Int(tst, "a", res, -2)
}
