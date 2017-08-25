// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"fmt"
	"testing"
)

func TestIntAssert01(tst *testing.T) {

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

	PrintTitle("IntAssert01")

	PrintOk("the next error message is")
	IntAssert(2, 1)
}

func TestIntAssert02(tst *testing.T) {

	//Verbose = true

	PrintTitle("IntAssert02")

	IntAssert(1, 1)
}

func TestDblAssert01(tst *testing.T) {

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

	PrintTitle("DblAssert01")

	PrintOk("the next error message is")
	Float64assert(2, 1)
}

func TestIntAssertLessthan01(tst *testing.T) {

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

	PrintTitle("IntAssertLessthan01")

	PrintOk("the next error message is")
	IntAssertLessThan(1, 1)
}

func TestIntAssertLessthan02(tst *testing.T) {

	//Verbose = true

	PrintTitle("IntAssertLessthan02")

	IntAssertLessThan(1, 2)
}

func TestIntAssertLessthanOrEqualTo01(tst *testing.T) {

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

	PrintTitle("IntAssertLessthanOrEqualTo01")

	PrintOk("the next error message is")
	IntAssertLessThanOrEqualTo(2, 1)
}

func TestIntAssertLessthanOrEqualTo02(tst *testing.T) {

	//Verbose = true

	PrintTitle("IntAssertLessthanOrEqualTo02")

	IntAssertLessThanOrEqualTo(1, 2)
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
