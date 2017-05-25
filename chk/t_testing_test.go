// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"fmt"
	"math"
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

func myfunction() {}

// TODO: this test doesn't work on MacOS
func test_FcnName(tst *testing.T) {

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

func Test_deriv01(tst *testing.T) {

	//Verbose = true
	PrintTitle("deriv01. DerivScaScaCen")

	f := func(x float64) float64 { return math.Cos(math.Pi * x / 2.0) }

	dfdxAna := -1.0 * math.Pi / 2.0
	xAt := 1.0
	dx := 1e-3

	t1 := new(testing.T)
	DerivScaScaCen(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, f)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivScaScaCen(t2, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, f)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}
}

func Test_deriv02(tst *testing.T) {

	//Verbose = true
	PrintTitle("deriv02. DerivVecScaCen")

	fcn := func(f []float64, x float64) {
		f[0] = math.Cos(math.Pi * x / 2.0)
		f[1] = math.Sin(math.Pi * x / 2.0)
		return
	}

	dfdx := func(x float64) []float64 {
		return []float64{
			-math.Sin(math.Pi*x/2.0) * math.Pi / 2.0,
			+math.Cos(math.Pi*x/2.0) * math.Pi / 2.0,
		}
	}

	dx := 1e-3
	xAt := 1.0
	dfdxAna := dfdx(xAt)

	t1 := new(testing.T)
	DerivVecScaCen(t1, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	DerivVecScaCen(t2, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t2.Failed() {
		tst.Errorf("t2 should not have failed\n")
		return
	}

	xAt = 0.0
	dfdxAna = dfdx(xAt)
	dfdxAna[0] += 0.0001

	t3 := new(testing.T)
	DerivVecScaCen(t3, "dfdx", 1e-15, dfdxAna, xAt, dx, Verbose, fcn)
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	dfdxAna = dfdx(xAt)
	t4 := new(testing.T)
	DerivVecScaCen(t4, "dfdx", 1.5e-11, dfdxAna, xAt, dx, Verbose, fcn)
	if t4.Failed() {
		tst.Errorf("t4 should not have failed\n")
		return
	}
}
