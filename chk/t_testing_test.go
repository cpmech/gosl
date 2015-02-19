// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"fmt"
	"testing"
)

func Test_IntAssert(tst *testing.T) {

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

func Test_DblAssert(tst *testing.T) {

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

func Test_IntAssertLessthan(tst *testing.T) {

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

func Test_IntAssertLessthanOrEqualTo(tst *testing.T) {

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

func Test_StrAssert(tst *testing.T) {

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
