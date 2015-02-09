// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"fmt"
	"testing"
)

func Test_IntAssert(tst *testing.T) {

	defer func() {
		err := recover()
		if err == nil {
			fmt.Println("[1;31mtest failed: assert should have failed[0m")
			tst.Fail()
		}
		Pf("the error message is: %v\n", err)
	}()

	//Tsilent = false
	TTitle("IntAssert")

	Pf("the next error message is OK\n")
	IntAssert(2, 1)
}

func Test_DblAssert(tst *testing.T) {

	defer func() {
		err := recover()
		if err == nil {
			fmt.Println("[1;31mtest failed: assert should have failed[0m")
			tst.Fail()
		}
		Pf("the error message is: %v\n", err)
	}()

	//Tsilent = false
	TTitle("DblAssert")

	Pf("the next error message is OK\n")
	DblAssert(2, 1)
}

func Test_IntAssertLessthan(tst *testing.T) {

	defer func() {
		err := recover()
		if err == nil {
			fmt.Println("[1;31mtest failed: assert should have failed[0m")
			tst.Fail()
		}
		Pf("the error message is: %v\n", err)
	}()

	//Tsilent = false
	TTitle("IntAssertLessthan")

	Pf("the next error message is OK\n")
	IntAssertLessThan(1, 1)
}

func Test_IntAssertLessthanOrEqualTo(tst *testing.T) {

	defer func() {
		err := recover()
		if err == nil {
			fmt.Println("[1;31mtest failed: assert should have failed[0m")
			tst.Fail()
		}
		Pf("the error message is: %v\n", err)
	}()

	//Tsilent = false
	TTitle("IntAssertLessthanOrEqualTo")

	Pf("the next error message is OK\n")
	IntAssertLessThanOrEqualTo(2, 1)
}

func Test_StrAssert(tst *testing.T) {

	defer func() {
		err := recover()
		if err == nil {
			fmt.Println("[1;31mtest failed: assert should have failed[0m")
			tst.Fail()
		}
		Pf("the error message is: %v\n", err)
	}()

	//Tsilent = false
	TTitle("StrAssert")

	Pf("the next error message is OK\n")
	StrAssert("rambo", "terminator")
}
