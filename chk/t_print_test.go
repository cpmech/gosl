// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"math"
	"testing"
)

func TestCallerInfo(tst *testing.T) {
	prev := Verbose
	//Verbose = true
	CallerInfo(100)
	Verbose = prev
}

func TestPrintOk(tst *testing.T) {
	prevV := Verbose
	prevC := ColorsOn
	//Verbose = true
	ColorsOn = false
	PrintTitle("PrintOk: verbose is ON for this test")
	PrintOk("x")
	ColorsOn = true
	PrintOk("x")
	Verbose = prevV
	ColorsOn = prevC
}

func TestPrintFail(tst *testing.T) {
	prevV := Verbose
	prevC := ColorsOn
	//Verbose = true
	ColorsOn = false
	PrintTitle("PrintFail: verbose is ON for this test")
	PrintFail("the message 'FAIL' is ok ⇒ ")
	ColorsOn = true
	PrintFail("the message 'FAIL' is ok ⇒ ")
	Verbose = prevV
	ColorsOn = prevC
}

func TestCheckAndPrint(tst *testing.T) {
	prevV := Verbose
	//Verbose = true
	PrintTitle("CheckAndPrint: verbose is ON for this test")
	t1 := new(testing.T)
	CheckAndPrint(t1, "x", 0, math.Sqrt(-1))
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
	}
	t2 := new(testing.T)
	CheckAndPrint(t2, "x", 0, math.Inf(+1))
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
	}
	t3 := new(testing.T)
	CheckAndPrint(t3, "x", 0, 1)
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
	}
	Verbose = prevV
}

func TestPrintAnaNum(tst *testing.T) {
	prevV := Verbose
	//Verbose = true
	PrintTitle("PrintAnaNum: verbose is ON for this test")
	err := PrintAnaNum("x", 0, math.Sqrt(-1), math.Sqrt(-1), Verbose)
	if err == nil {
		tst.Errorf("error message should not be nil\n")
	}
	err = PrintAnaNum("x", 0, math.Inf(+1), math.Inf(+1), Verbose)
	if err == nil {
		tst.Errorf("error message should not be nil\n")
	}
	err = PrintAnaNum("x", 0, 1, 2, Verbose)
	if err == nil {
		tst.Errorf("error message should not be nil\n")
	}
	Verbose = prevV
}
