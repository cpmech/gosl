// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"fmt"
	"math"
	"runtime"
	"testing"
)

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
			fmt.Println("[1;31mFAIL[0m")
			return
		}
		fmt.Println("FAIL")
	}
}

// PrintTitle returns the Test Title
func PrintTitle(title string) {
	if Verbose {
		fmt.Printf("\n=== %s =================\n", title)
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

// PrintAnaNumC formats the output of analytical versus numerical comparisons (complex version)
func PrintAnaNumC(msg string, tol float64, ana, num complex128, verbose bool) (e error) {
	diffR := math.Abs(real(ana) - real(num))
	diffC := math.Abs(imag(ana) - imag(num))
	if math.IsNaN(diffR) || math.IsInf(diffR, 0) {
		e = Err("[1;31m%s (real part) failed with NaN or Inf: %v[0m", msg, diffR)
		return
	}
	if math.IsNaN(diffC) || math.IsInf(diffC, 0) {
		e = Err("[1;31m%s (imag part) failed with NaN or Inf: %v[0m", msg, diffC)
		return
	}
	if verbose {
		clrR := "[1;32m" // green
		clrC := "[1;32m" // green
		if diffR > tol {
			clrR = "[1;31m" // red
		}
		if diffC > tol {
			clrC = "[1;31m" // red
		}
		f := "%" + fmt.Sprintf("%d", len(msg)) + "s"
		fmt.Printf(f+" %23v  %23v  %s%23v[0m\n", msg, real(ana), real(num), clrR, diffR)
		fmt.Printf(f+" %23vi %23vi %s%23v[0m\n", "", imag(ana), imag(num), clrC, diffC)
	}
	if diffR > tol || diffC > tol {
		e = Err("[1;31m%s failed with |diffR| = %g  |diffC| = %g[0m", msg, diffR, diffC)
	}
	return
}

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
