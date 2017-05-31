// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chk contains functions for checking and testing computations
package chk

import (
	"fmt"
	"math"
	"runtime"
	"testing"
)

var (
	// AssertOn activates or deactivates asserts
	AssertOn = true

	// Verbose turn on verbose mode
	Verbose = false

	// ColorsOn turn on use of colours on console
	ColorsOn = true
)

// CallerInfo returns the file and line positions where an error occurred
//  idx -- use idx=2 to get the caller of Panic
func CallerInfo(idx int) {
	pc, file, line, ok := runtime.Caller(idx)
	if !ok {
		file, line = "?", 0
	}
	var fname string
	f := runtime.FuncForPC(pc)
	if f != nil {
		fname = f.Name()
	}
	if Verbose {
		fmt.Printf("file = %s:%d\n", file, line)
		fmt.Printf("func = %s\n", fname)
	}
}

// PanicSimple panicks without calling CallerInfo
func PanicSimple(msg string, prm ...interface{}) {
	panic(fmt.Sprintf(msg, prm...))
}

// Panic panicks
func Panic(msg string, prm ...interface{}) {
	CallerInfo(4)
	CallerInfo(3)
	CallerInfo(2)
	panic(fmt.Sprintf(msg, prm...))
}

// Err returns a new error
func Err(msg string, prm ...interface{}) error {
	return fmt.Errorf(msg, prm...)
}

// PrintOk prints "OK" in green (if ColorsOn==true)
func PrintOk(msg string, prm ...interface{}) {
	if Verbose {
		fmt.Printf(msg, prm...)
		if ColorsOn {
			fmt.Println(" [1;32mOK[0m")
			return
		}
		fmt.Println(" OK")
	}
}

// PrintFail prints "FAIL" in red (if ColorsOn==true)
func PrintFail(msg string, prm ...interface{}) {
	if Verbose {
		fmt.Printf(msg, prm...)
		if ColorsOn {
			fmt.Println(" [1;31mFAIL[0m")
			return
		}
		fmt.Println(" FAIL")
	}
}

// PrintTitle returns the Test Title
func PrintTitle(title string) {
	if Verbose {
		fmt.Printf("\n================= %s =================\n", title)
		return
	}
	fmt.Printf("   . . . testing . . .   %s\n", title)
}

// CheckAndPrint returns a formatted error message
func CheckAndPrint(tst *testing.T, msg string, tol, diff float64) {
	if math.IsNaN(diff) || math.IsInf(diff, 0) {
		tst.Errorf("[1;31m%s failed with NaN or Inf: %v[0m", msg, diff)
		return
	}
	if diff > tol {
		if Verbose {
			fmt.Printf("%s [1;31merror |diff| = %g[0m\n", msg, diff)
		}
		tst.Errorf("[1;31m%s failed with |diff| = %g[0m", msg, diff)
		return
	}
	PrintOk(msg)
}

// PrintAnaNum formats the output of analytical versus numerical comparisons
func PrintAnaNum(msg string, tol, ana, num float64, verbose bool) (e error) {
	diff := math.Abs(ana - num)
	if math.IsNaN(diff) || math.IsInf(diff, 0) {
		e = Err("[1;31m%s failed with NaN or Inf: %v[0m", msg, diff)
		return
	}
	if verbose {
		clr := "[1;32m" // green
		if diff > tol {
			clr = "[1;31m" // red
		}
		fmt.Printf("%s %23v %23v %s%23v[0m\n", msg, ana, num, clr, diff)
	}
	if diff > tol {
		e = Err("[1;31m%s failed with |diff| = %g[0m", msg, diff)
	}
	return
}
