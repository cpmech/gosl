// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"fmt"
	"testing"
)

func TestIntAssert(tst *testing.T) {

	//Verbose = true
	defer func() {
		if err := recover(); err != nil {
			if Verbose {
				fmt.Printf("OK, caught the following message:\n\n\t%v\n", err)
			}
		} else {
			tst.Errorf("\n\tTEST FAILED. Assert should have panicked\n")
		}
	}()

	PrintTitle("IntAssert")

	PrintOk("the next error message is")
	IntAssert(2, 1)
}

func TestDblAssert(tst *testing.T) {

	//Verbose = true
	defer func() {
		if err := recover(); err != nil {
			if Verbose {
				fmt.Printf("OK, caught the following message:\n\n\t%v\n", err)
			}
		} else {
			tst.Errorf("\n\tTEST FAILED. Assert should have panicked\n")
		}
	}()

	PrintTitle("DblAssert")

	PrintOk("the next error message is")
	DblAssert(2, 1)
}

func TestIntAssertLessthan(tst *testing.T) {

	//Verbose = true
	defer func() {
		if err := recover(); err != nil {
			if Verbose {
				fmt.Printf("OK, caught the following message:\n\n\t%v\n", err)
			}
		} else {
			tst.Errorf("\n\tTEST FAILED. Assert should have panicked\n")
		}
	}()

	PrintTitle("IntAssertLessthan")

	PrintOk("the next error message is")
	IntAssertLessThan(1, 1)
}

func TestIntAssertLessthanOrEqualTo(tst *testing.T) {

	//Verbose = true
	defer func() {
		if err := recover(); err != nil {
			if Verbose {
				fmt.Printf("OK, caught the following message:\n\n\t%v\n", err)
			}
		} else {
			tst.Errorf("\n\tTEST FAILED. Assert should have panicked\n")
		}
	}()

	PrintTitle("IntAssertLessthanOrEqualTo")

	PrintOk("the next error message is")
	IntAssertLessThanOrEqualTo(2, 1)
}

func TestStrAssert(tst *testing.T) {

	//Verbose = true
	defer func() {
		if err := recover(); err != nil {
			if Verbose {
				fmt.Printf("OK, caught the following message:\n\n\t%v\n", err)
			}
		} else {
			tst.Errorf("\n\tTEST FAILED. Assert should have panicked\n")
		}
	}()

	PrintTitle("StrAssert")

	PrintOk("the next error message is")
	StrAssert("rambo", "terminator")
}

func myfunction() {}

// TODO: this test doesn't work on MacOS
func testFcnName(tst *testing.T) {

	//Verbose = true
	PrintTitle("FcnName")

	name := GetFunctionName(myfunction)
	if Verbose {
		fmt.Printf("name = %v\n", name)
	}
	if name != "github.com/cpmech/gosl/chk.myfunction" {
		tst.Errorf("function name is incorrect\n")
	}

	fcn := func() {}
	name = GetFunctionName(fcn)
	if Verbose {
		fmt.Printf("name = %v\n", name)
	}
	//if !strings.HasPrefix(name, "github.com/cpmech/gosl/chk.Test_FcnName.func") {
	//tst.Errorf("function name is incorrect\n")
	//}
}

func TestScalarC(tst *testing.T) {

	//Verbose = true

	PrintTitle("ScalarC")

	a := 123.123 + 456.456i
	b := 123.123 + 456.456i
	t1 := new(testing.T)
	Complex128(t1, "|a-b|", 1e-17, a, b)
	if t1.Failed() {
		tst.Errorf("t1 should not have failed\n")
		return
	}

	b = 123.1231 + 456.456i
	t2 := new(testing.T)
	Complex128(t2, "|a-b|", 1e-17, a, b)
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}
}

func TestInts(tst *testing.T) {

	//Verbose = true

	PrintTitle("Ints")

	// std int -----------------------------------------

	a := []int{1, 2, 3, 4}
	t1 := new(testing.T)
	Ints(t1, "a", a, []int{1, 2, 3, 3})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Ints(t2, "a", a, []int{1, 2, 3, 4})
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}

	// 32 int -----------------------------------------

	b := []int32{1, 2, 3, 4}
	t3 := new(testing.T)
	Int32s(t3, "b", b, []int32{1, 2, 3, 3})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	Int32s(t4, "b", b, []int32{1, 2, 3, 4})
	if t4.Failed() {
		tst.Errorf("t4 should not have failed\n")
		return
	}

	// 64 int -----------------------------------------

	c := []int64{1, 2, 3, 4}
	t5 := new(testing.T)
	Int64s(t5, "c", c, []int64{1, 2, 3, 3})
	if !t5.Failed() {
		tst.Errorf("t5 should have failed\n")
		return
	}

	t6 := new(testing.T)
	Int64s(t6, "c", c, []int64{1, 2, 3, 4})
	if t6.Failed() {
		tst.Errorf("t6 should not have failed\n")
		return
	}
}
