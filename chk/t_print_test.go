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

func TestTstFail(tst *testing.T) {
	prevV := Verbose
	prevC := ColorsOn
	//Verbose = true
	ColorsOn = false
	PrintTitle("TstFail: verbose is ON for this test")
	t1 := new(testing.T)
	TstFail(t1, "the message 'FAIL' is ok ⇒ ")
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
	}
	ColorsOn = true
	t2 := new(testing.T)
	TstFail(t2, "the message 'FAIL' is ok ⇒ ")
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
	}
	Verbose = prevV
	ColorsOn = prevC
}

func TestTstDiff(tst *testing.T) {
	prevV := Verbose
	//Verbose = true
	PrintTitle("TstDiff: verbose is ON for this test")
	t1 := new(testing.T)
	TstDiff(t1, "x", 0, math.Sqrt(-1), math.Sqrt(-1), false)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
	}
	t2 := new(testing.T)
	TstDiff(t2, "x", 0, math.Inf(+1), math.Inf(+1), true)
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
	}
	t3 := new(testing.T)
	TstDiff(t3, "x", 0, 1, 0, false)
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
	}
	Verbose = prevV
}

func TestTstDiffC(tst *testing.T) {
	prevV := Verbose
	//Verbose = true
	PrintTitle("TstDiffC: verbose is ON for this test")
	t1 := new(testing.T)
	TestDiffC(t1, "x", 0, complex(math.Sqrt(-1), 1), complex(math.Sqrt(-1), 1), false)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
	}
	t2 := new(testing.T)
	TestDiffC(t2, "x", 0, complex(math.Inf(+1), 1), complex(math.Inf(+1), 1), true)
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
	}
	t3 := new(testing.T)
	TestDiffC(t3, "x", 0, 1+1i, 0, false)
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

func TestPrintAnaNumC(tst *testing.T) {
	prevV := Verbose
	//Verbose = true
	PrintTitle("PrintAnaNumC: verbose is ON for this test")

	a := 123 + 456i
	b := 123 + 456i
	err := PrintAnaNumC("x", 0, a, b, Verbose)
	if err != nil {
		tst.Errorf("error message should be nil\n")
	}

	a = complex(math.NaN(), 0)
	b = 1i
	err = PrintAnaNumC("x", 0, a, b, Verbose)
	if err == nil {
		tst.Errorf("error message should NOT be nil\n")
	}

	a = 1i
	b = complex(0, math.NaN())
	err = PrintAnaNumC("x", 0, a, b, Verbose)
	if err == nil {
		tst.Errorf("error message should NOT be nil\n")
	}

	a = complex(math.Inf(+1), 0)
	b = complex(math.Inf(+1), 0)
	err = PrintAnaNumC("x", 0, a, b, Verbose)
	if err == nil {
		tst.Errorf("error message should NOT be nil\n")
	}

	err = PrintAnaNumC("x", 0, 1, 2, Verbose)
	if err == nil {
		tst.Errorf("error message should NOT be nil\n")
	}

	err = PrintAnaNumC("x", 0, 1i, 2i, Verbose)
	if err == nil {
		tst.Errorf("error message should NOT be nil\n")
	}
	Verbose = prevV
}
