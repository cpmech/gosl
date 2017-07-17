// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"fmt"
	"testing"
)

func TestET(tst *testing.T) {

	//Verbose = true
	PrintTitle("ET (error test)")

	t1 := new(testing.T)
	ET(t1, nil)
	if t1.Failed() {
		tst.Errorf("t1 should NOT have failed\n")
	}

	err := Err("stop")
	t2 := new(testing.T)
	ET(t2, err)
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
	}
}

func TestEP(tst *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			if Verbose {
				fmt.Printf("OK, caught the following message:\n\n\t%v\n", err)
			}
		} else {
			tst.Errorf("\n\tTEST FAILED. test should have panicked\n")
		}
	}()

	//Verbose = true
	PrintTitle("EP (error panic)")
	EP(Err("stop"))
}

func TestErr01(tst *testing.T) {

	//Verbose = true
	PrintTitle("Err01")

	err := Err("TEST_ERR = %d", 666)
	if err.Error() != "TEST_ERR = 666" {
		panic("Err failed")
	}

	PrintOk("hello from PrintOk => ")
	PrintFail("hello from PrintFail => ")
}

func TestCallerInfo(tst *testing.T) {
	//Verbose = true
	PrintTitle("CallerInfo")
	CallerInfo(100)
}

func TestPanicErr(tst *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			if Verbose {
				fmt.Printf("OK, caught the following message:\n\n\t%v\n", err)
			}
		} else {
			tst.Errorf("\n\tTEST FAILED. test should have panicked\n")
		}
	}()

	//Verbose = true
	PrintTitle("PanicErr")
	PanicErr(Err("stop"))
}
