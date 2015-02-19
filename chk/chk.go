// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"testing"
)

var (
	AssertOn = true // activates/deactivates asserts
	Verbose  = false
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
	return errors.New(fmt.Sprintf(msg, prm...))
}

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
		fmt.Printf("\n========================= %s =========================\n", title)
		return
	}
	fmt.Printf("   . . . testing . . .   %s\n", title)
}

// PrintMsg returns a formatted error message
func PrintMsg(tst *testing.T, msg string, tol, maxdiff float64) {
	if Verbose {
		if maxdiff > tol {
			fmt.Printf("%s [1;31merror |maxdiff| = %g[0m\n", msg, maxdiff)
			tst.Errorf("[1;31m%s failed with |maxdiff| = %g[0m", msg, maxdiff)
			return
		}
		PrintOk(msg)
	}
}

// PrintAnaNum formats the output of analytical versus numerical comparisons
func PrintAnaNum(msg string, tol, ana, num float64, verbose bool) (e error) {
	err := math.Abs(ana - num)
	if math.IsNaN(err) || math.IsInf(err, 0) {
		e = Err("[1;31m%s failed with NaN or Inf: %v[0m", msg, err)
		return
	}
	if verbose {
		clr := "[1;32m" // green
		if err > tol {
			clr = "[1;31m" // red
		}
		fmt.Printf("%s %23v %23v %s%23v[0m\n", msg, ana, num, clr, err)
	}
	if err > tol {
		e = Err("[1;31m%s failed with |maxdiff| = %g[0m", msg, err)
	}
	return
}
