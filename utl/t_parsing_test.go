// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
    "math"
    "testing"
)

func Test_keycode01(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("keycode01")

    str := "!typeA:keycodeA !typeB:keycodeB!typeC  : keycodeC"
    if resA, ok := Keycode(str, "typeA"); !ok { tst.Errorf("[1;31mcannot find   typeA[0m") } else { if resA != "keycodeA" { tst.Errorf("[1;31mresA != keycodeA[0m") } }
    if resB, ok := Keycode(str, "typeB"); !ok { tst.Errorf("[1;31mcannot find   typeB[0m") } else { if resB != "keycodeB" { tst.Errorf("[1;31mresB != keycodeB[0m") } }
    if resC, ok := Keycode(str, "typeC"); !ok { tst.Errorf("[1;31mcannot find   typeC[0m") } else { if resC != "keycodeC" { tst.Errorf("[1;31mresC != keycodeC[0m") } }
    if resD, ok := Keycode(str, "typeD");  ok { tst.Errorf("[1;31mmust not find typeD[0m") } else { if resD != ""         { tst.Errorf("[1;31mresD != \"\"    [0m") } }

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

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("keycode02")

    str := "!flagA !typeA:keycodeA !typeB:keycodeB!typeC  : keycodeC !FlagB"
	res := Keycodes(str)
	CompareStrs(tst, "keycodes", res, []string{"flagA", "typeA", "typeB", "typeC", "FlagB"})
}

func Test_parsing02(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("parsing02")

    keys  := []string{"ux", "uy", "uz", "pl", "sl"}
    skeys := JoinKeys(keys)
    kkeys := SplitKeys(skeys)
    Pforan("keys  = %+v\n", keys)
    Pforan("skeys = %q\n",  skeys)
    Pforan("kkeys = %+v\n", kkeys)
    CompareStrs(tst, "keys", keys, kkeys)

    k1, k2 := []string{"ex1", "ex2"}, []string{"more1", "more2"}
    allkeys := JoinKeys3(keys, k1, k2, "")
    Pfpink("allkeys = %q\n", allkeys)
    CheckString(tst, allkeys, "ux uy uz pl sl ex1 ex2 more1 more2")

    allkeys = JoinKeys3(keys, k1, k2, ",")
    Pfpink("allkeys = %q\n", allkeys)
    CheckString(tst, allkeys, "ux uy uz pl sl, ex1 ex2, more1 more2")

    r0, r1, r2 := SplitKeys3(allkeys)
    Pfblue2("r0 = %v\n", r0)
    Pfblue2("r1 = %v\n", r1)
    Pfblue2("r2 = %v\n", r2)
    CompareStrs(tst, "r0", r0, keys)
    CompareStrs(tst, "r1", r1, k1)
    CompareStrs(tst, "r2", r2, k2)

    allkeys = JoinKeys3(keys, []string{}, k2, ",")
    Pfpink("allkeys = %q\n", allkeys)
    CheckString(tst, allkeys, "ux uy uz pl sl,, more1 more2")

    allkeys = JoinKeys3(keys, []string{}, []string{}, ",")
    Pfpink("allkeys = %q\n", allkeys)
    CheckString(tst, allkeys, "ux uy uz pl sl,,")

    allkeys = JoinKeys3([]string{}, k1, k2, ",")
    Pfpink("allkeys = %q\n", allkeys)
    CheckString(tst, allkeys, ", ex1 ex2, more1 more2")

    r0, r1, r2 = SplitKeys3("a ,b c,")
    Pfblue2("r0 = %v\n", r0)
    Pfblue2("r1 = %v\n", r1)
    Pfblue2("r2 = %v\n", r2)
    CompareStrs(tst, "r0", r0, []string{"a"})
    CompareStrs(tst, "r1", r1, []string{"b","c"})
    CompareStrs(tst, "r2", r2, []string{})
}

func Test_parsing03(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("parsing03")

    keys       := []string{"ux", "uy", "uz", "pl", "sl"}
    k1, k2, k3 := []string{"ex1", "ex2"}, []string{"more1", "more2"}, []string{"a", "b", "c"}
    allkeys := JoinKeys4(keys, k1, k2, k3, "")
    Pfpink("allkeys = %q\n", allkeys)
    CheckString(tst, allkeys, "ux uy uz pl sl ex1 ex2 more1 more2 a b c")

    allkeys = JoinKeys4(keys, k1, k2, k3, ",")
    Pfpink("allkeys = %q\n", allkeys)
    CheckString(tst, allkeys, "ux uy uz pl sl, ex1 ex2, more1 more2, a b c")

    r0, r1, r2, r3 := SplitKeys4(allkeys)
    Pfblue2("r0 = %v\n", r0)
    Pfblue2("r1 = %v\n", r1)
    Pfblue2("r2 = %v\n", r2)
    Pfblue2("r3 = %v\n", r3)
    CompareStrs(tst, "r0", r0, keys)
    CompareStrs(tst, "r1", r1, k1)
    CompareStrs(tst, "r2", r2, k2)
    CompareStrs(tst, "r3", r3, k3)

    allkeys = JoinKeys4(keys, []string{}, k2, []string{}, ",")
    Pfpink("allkeys = %q\n", allkeys)
    CheckString(tst, allkeys, "ux uy uz pl sl,, more1 more2,")
}

func Test_parsing04(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("parsing04")

    keys  := []string{"ux", "uy", "uz"}
    skeys := JoinKeysPre("R", keys)
    Pforan("keys  = %+v\n", keys)
    Pforan("skeys = %q\n",  skeys)
    CheckString(tst, skeys, "Rux Ruy Ruz")
}

func TestEval01(tst *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    TTitle("Eval 01")

    corr := ((5.0*3.0/2.0)+(2.0-3.0+4.0)-(5.0/3.0))*2.0
    expr := "((5*3/2)+(2-3+4)-(5/3))*2"
    res  := Eval(expr, nil, nil)
    Pf("%v = %v\n", expr, res)
    CheckScalar(tst, expr, 1e-20, res, corr)
}

func TestEval02(tst *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    TTitle("Eval 02")

    corr := (((5.0*3.0/2.0)+(2.0-3.0+4.0)-(5.0/3.0))*2.0 + math.Sin(math.Pi/2.0) + math.Pow(3.0, 3.0/2.0)) / 2.0
    expr := "(((5*3/2)+(2-3+4)-(5/3))*2 + sin(pi/2) + pow(3, 3/2)) / 2"
    res  := Eval(expr, nil, nil)
    Pf("%v = %v\n", expr, res)
    CheckScalar(tst, expr, 1e-20, res, corr)
}

func TestEval03(tst *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    TTitle("Eval 03")

    expr     := "(2+(5+(3+b*3)+1)*2)/3"
    aex, err := ParseE(expr)
    if err != nil {
        Panic("cannot parse expression: %v", err)
    }
    res, err := EvalE(aex, nil, nil)
    if err == nil {
        Panic("EvalE must return a non-nil error, since 'b' is not available")
    }
    Pf("%v = %v\n", expr, res)
    Pf("err = %v\n", UnColor(err.Error()))
    PfGreen("OK\n")
}

func TestEval04(tst *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    TTitle("Eval 04")

    vars     := map[string]float64{"b": 5}
    expr     := "(2+(5+(3+b*3)+1)*2)/3"
    corr     := 50.0 / 3.0
    aex, err := ParseE(expr)
    if err != nil {
        Panic("cannot parse expression: %v", err)
    }
    res, err := EvalE(aex, vars, nil)
    if err != nil {
        Panic("cannot evaluate expression: %v", err)
    }
    Pf("%v = %v\n", expr, res)
    CheckScalar(tst, expr, 1e-20, res, corr)

    expr += " + dorival(8)"
    aex, err = ParseE(expr)
    if err != nil {
        Panic("cannot parse expression: %v", err)
    }
    res, err = EvalE(aex, vars, nil)
    if err == nil {
        Panic("EvalE must return a non-nil error, since 'dorival' is not available")
    }
    Pf("%v = %v\n", expr, res)
    Pf("err = %v\n", UnColor(err.Error()))
    PfGreen("OK\n")
}

func TestEval05(tst *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    TTitle("Eval 05")

    vars := map[string]float64{"b":-1, "β":200, "t":0.02}
    expr := "b*β*β*exp(β*t)"
    corr := -2183926.0013257693
    res  := Eval(expr, vars, nil)
    Pf("%v = %v\n", expr, res)
    CheckScalar(tst, expr, 1e-20, res, corr)
}

func TestEval06(tst *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    TTitle("Eval 06")

    expr := "-(2*(3+4))"
    corr := -(2.0*(3.0+4.0))
    res  := Eval(expr, nil, nil)
    Pf("%v = %v\n", expr, res)
    CheckScalar(tst, expr, 1e-20, res, corr)
}
