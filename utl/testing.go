// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"fmt"
	"math"
	"runtime"
	"testing"
)

var Tassert = true // activates/deactivates asserts

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
	if UseColors {
		Pf("[1;35mfile = %s:%d[0m\n", file, line)
		Pf("[1;35mfunc = %s[0m\n", fname)
		return
	}
	Pf("file = %s:%d\n", file, line)
	Pf("func = %s\n", fname)
}

// IntAssert asserts that a is equal to b (ints)
func IntAssert(a, b int) {
	if Tassert {
		if a != b {
			CallerInfo(3)
			CallerInfo(2)
			PanicSimple("Assert failed: ints are different: %d != %d", a, b)
		}
	}
}

// IntAssertLessThan asserts that a < b (ints)
func IntAssertLessThan(a, b int) {
	if Tassert {
		if a < b {
			return
		}
		CallerInfo(3)
		CallerInfo(2)
		PanicSimple("Assert failed: %d < %d is incorrect", a, b)
	}
}

// IntAssertLessThanOrEqualTo asserts that a ≤ b (ints)
func IntAssertLessThanOrEqualTo(a, b int) {
	if Tassert {
		if a <= b {
			return
		}
		CallerInfo(3)
		CallerInfo(2)
		PanicSimple("Assert failed: %d ≤ %d is incorrect", a, b)
	}
}

// DblAssert asserts that a is equal to b (floats)
func DblAssert(a, b float64) {
	if Tassert {
		if a != b {
			CallerInfo(3)
			CallerInfo(2)
			PanicSimple("Assert failed: float64 are different: %g != %g", a, b)
		}
	}
}

// StrAssert asserts that a is equal to b (strings)
func StrAssert(a, b string) {
	if Tassert {
		if a != b {
			CallerInfo(3)
			CallerInfo(2)
			PanicSimple("Assert failed: strings are different: %s != %s", a, b)
		}
	}
}

// TTitle returns the Test Title
func TTitle(title string) {
	if Tsilent {
		fmt.Printf("   . . . testing . . .   %s\n", title)
	} else {
		fmt.Printf("\n========================= %s =========================\n", title)
	}
}

// ErrorMsg returns a formatted error message
func ErrorMsg(tst *testing.T, msg string, tol, maxdiff float64) {
	if maxdiff > tol {
		fmt.Printf("%s [1;31merror |maxdiff| = %g[0m\n", msg, maxdiff)
		tst.Errorf("[1;31m%s failed with |maxdiff| = %g[0m", msg, maxdiff)
	} else {
		if !Tsilent {
			fmt.Printf("%s [1;32mOK[0m\n", msg)
		}
	}
}

// CheckScalar compares two scalars
func CheckScalar(tst *testing.T, msg string, tol, res, correct float64) {
	ErrorMsg(tst, msg, tol, math.Abs(res-correct))
}

// CheckString compares two strings
func CheckString(tst *testing.T, str, correct string) {
	if str != correct {
		fmt.Printf("[1;31merror %s != %s[0m\n", str, correct)
		tst.Errorf("[1;31mstring failed with: %s != %s[0m", str, correct)
	} else {
		if !Tsilent {
			fmt.Printf("%s == %s [1;32mOK[0m\n", str, correct)
		}
	}
}

// AnaNum formats the output of analytical versus numerical comparisons
func AnaNum(msg string, tol, ana, num float64, verbose bool) (e error) {
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
		Pf("%s %23v %23v %s%23v[0m\n", msg, ana, num, clr, err)
	}
	if err > tol {
		e = Err("[1;31m%s failed with |maxdiff| = %g[0m", msg, err)
	}
	return
}

// CheckAnaNum compares analytical versus numerical values
func CheckAnaNum(tst *testing.T, msg string, tol, ana, num float64, verbose bool) {
	e := AnaNum(msg, tol, ana, num, verbose)
	if e != nil {
		tst.Errorf("%v", e.Error())
	}
}

// CheckMatrix compares two matrices
func CheckMatrix(tst *testing.T, msg string, tol float64, res, correct [][]float64) {
	zero := false
	if len(correct) == 0 {
		zero = true
	} else {
		if len(res) != len(correct) {
			fmt.Printf("%s [1;31merror[0m\n", msg)
			tst.Errorf("[1;31m%s failed: res and correct matrices have different lengths. %d != %d[0m", msg, len(res), len(correct))
			return
		}
	}
	var diff, maxdiff float64
	for i := 0; i < len(res); i++ {
		if !zero {
			if len(res[i]) != len(correct[i]) {
				fmt.Printf("%s [1;31m error[0m\n", msg)
				tst.Errorf("[1;31m%s failed: matrices have different number of columns[0m", msg)
			}
		}
		for j := 0; j < len(res[i]); j++ {
			if math.IsNaN(res[i][j]) {
				tst.Errorf("[1;31m%s failed: NaN detected => %v[0m", msg, res[i][j])
			}
			if zero {
				diff = math.Abs(res[i][j])
			} else {
				diff = math.Abs(res[i][j] - correct[i][j])
			}
			if diff > maxdiff {
				maxdiff = diff
			}
		}
	}
	ErrorMsg(tst, msg, tol, maxdiff)
}

// CheckVector compares two vectors
func CheckVector(tst *testing.T, msg string, tol float64, res, correct []float64) {
	zero := false
	if len(correct) == 0 {
		zero = true
	} else {
		if len(res) != len(correct) {
			fmt.Printf("%s [1;31merror[0m\n", msg)
			tst.Errorf("[1;31m%s failed: res and correct vectors have different lengths. %d != %d[0m", msg, len(res), len(correct))
			return
		}
	}
	var diff, maxdiff float64
	for i := 0; i < len(res); i++ {
		if math.IsNaN(res[i]) {
			tst.Errorf("[1;31m%s failed: NaN detected => %v[0m", msg, res[i])
		}
		if zero {
			diff = math.Abs(res[i])
		} else {
			diff = math.Abs(res[i] - correct[i])
		}
		if diff > maxdiff {
			maxdiff = diff
		}
	}
	ErrorMsg(tst, msg, tol, maxdiff)
}

// CheckMatrixC compares two matrices of complex nummbers
func CheckMatrixC(tst *testing.T, msg string, tol float64, res, correct [][]complex128) {
	zero := false
	if len(correct) == 0 {
		zero = true
	} else {
		if len(res) != len(correct) {
			fmt.Printf("%s [1;31merror[0m\n", msg)
			tst.Errorf("[1;31m%s failed: res and correct matrices have different lengths. %d != %d[0m", msg, len(res), len(correct))
			return
		}
	}
	var diff, maxdiff float64
	var diffz, maxdiffz float64
	for i := 0; i < len(res); i++ {
		if !zero {
			if len(res[i]) != len(correct[i]) {
				fmt.Printf("%s [1;31m error[0m\n", msg)
				tst.Errorf("[1;31m%s failed: matrices have different number of columns[0m", msg)
			}
		}
		for j := 0; j < len(res[i]); j++ {
			if zero {
				diff = math.Abs(real(res[i][j]))
				diffz = math.Abs(imag(res[i][j]))
			} else {
				diff = math.Abs(real(res[i][j]) - real(correct[i][j]))
				diffz = math.Abs(imag(res[i][j]) - imag(correct[i][j]))
			}
			if diff > maxdiff {
				maxdiff = diff
			}
			if diffz > maxdiffz {
				maxdiffz = diffz
			}
		}
	}
	if maxdiff > tol || maxdiffz > tol {
		fmt.Printf("%s [1;31merror |maxdiff| = %g,  |maxdiffz| = %g[0m\n", msg, maxdiff, maxdiffz)
		tst.Errorf("[1;31m%s failed with |maxdiff| = %g,  |maxdiffz| = %g[0m", msg, maxdiff, maxdiffz)
	} else {
		if !Tsilent {
			fmt.Printf("%s [1;32mOK[0m\n", msg)
		}
	}
}

// CheckVectorC compares two vectors of complex nummbers
func CheckVectorC(tst *testing.T, msg string, tol float64, res, correct []complex128) {
	zero := false
	if len(correct) == 0 {
		zero = true
	} else {
		if len(res) != len(correct) {
			fmt.Printf("%s [1;31merror[0m\n", msg)
			tst.Errorf("[1;31m%s failed: res and correct matrices have different lengths. %d != %d[0m", msg, len(res), len(correct))
			return
		}
	}
	var diff, maxdiff float64
	var diffz, maxdiffz float64
	for i := 0; i < len(res); i++ {
		if zero {
			diff = math.Abs(real(res[i]))
			diffz = math.Abs(imag(res[i]))
		} else {
			diff = math.Abs(real(res[i]) - real(correct[i]))
			diffz = math.Abs(imag(res[i]) - imag(correct[i]))
		}
		if diff > maxdiff {
			maxdiff = diff
		}
		if diffz > maxdiffz {
			maxdiffz = diffz
		}
	}
	if maxdiff > tol || maxdiffz > tol {
		fmt.Printf("%s [1;31merror |maxdiff| = %g,  |maxdiffz| = %g[0m\n", msg, maxdiff, maxdiffz)
		tst.Errorf("[1;31m%s failed with |maxdiff| = %g,  |maxdiffz| = %g[0m", msg, maxdiff, maxdiffz)
	} else {
		if !Tsilent {
			fmt.Printf("%s [1;32mOK[0m\n", msg)
		}
	}
}

// CompareInts compares two slices of integers
func CompareInts(tst *testing.T, msg string, a, b []int) {
	if len(a) != len(b) {
		fmt.Printf("%s [1;31merror len(a)=%d != len(b)=%d[0m\n", msg, len(a), len(b))
		tst.Errorf("[1;31m%s failed: slices have different lengths: %v != %v[0m", msg, a, b)
		return
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			fmt.Printf("%s [1;31merror %d != %d[0m\n", msg, a[i], b[i])
			tst.Errorf("[1;31m%s failed: slices are different: %dth comp %v != %v\n%v != \n%v[0m", msg, i, a[i], b[i], a, b)
			return
		}
	}
	if !Tsilent {
		fmt.Printf("%s [1;32mOK[0m\n", msg)
	}
}

// CompareDbls compares two slices of doubles (floats)
func CompareDbls(tst *testing.T, msg string, a, b []float64) {
	if len(a) != len(b) {
		fmt.Printf("%s [1;31merror len(a)=%d != len(b)=%d[0m\n", msg, len(a), len(b))
		tst.Errorf("[1;31m%s failed: slices have different lengths: %v != %v[0m", msg, a, b)
		return
	}
	for i := 0; i < len(a); i++ {
		if math.Abs(a[i]-b[i]) > 1.0e-16 {
			fmt.Printf("%s [1;31merror %v != %v[0m\n", msg, a[i], b[i])
			tst.Errorf("[1;31m%s failed: slices are different: %dth comp %v != %v\n%v != \n%v[0m", msg, i, a[i], b[i], a, b)
			return
		}
	}
	if !Tsilent {
		fmt.Printf("%s [1;32mOK[0m\n", msg)
	}
}

// CompareDeep3 compares two deep3 slices
func CompareDeep3(tst *testing.T, msg string, a, b [][][]float64) {
	if len(a) != len(b) {
		fmt.Printf("%s [1;31merror len(a)=%d != len(b)=%d[0m\n", msg, len(a), len(b))
		tst.Errorf("[1;31m%s failed: slices have different lengths: %v != %v[0m", msg, a, b)
		return
	}
	for i := 0; i < len(a); i++ {
		if len(a[i]) != len(b[i]) {
			fmt.Printf("%s [1;31merror len(a[%d])=%d != len(b[%d])=%d[0m\n", msg, i, len(a[i]), i, len(b[i]))
			tst.Errorf("[1;31m%s failed: subslices have different lengths[0m", msg)
			return
		}
		for j := 0; j < len(a[i]); j++ {
			if len(a[i][j]) != len(b[i][j]) {
				fmt.Printf("%s [1;31merror len(a[%d][%d])=%d != len(b[%d][%d])=%d[0m\n", msg, i, j, len(a[i][j]), i, j, len(b[i][j]))
				tst.Errorf("[1;31m%s failed: subsubslices have different lengths[0m", msg)
				return
			}
			for k := 0; k < len(a[i][j]); k++ {
				if math.Abs(a[i][j][k]-b[i][j][k]) > 1.0e-16 {
					fmt.Printf("%s [1;31merror %v != %v[0m\n", msg, a[i][j][k], b[i][j][k])
					tst.Errorf("[1;31m%s failed: slices are different: %d,%d,%d component %v != %v\n%v != \n%v[0m", msg, i, j, k, a[i][j][k], b[i][j][k], a, b)
					return
				}
			}
		}
	}
	if !Tsilent {
		fmt.Printf("%s [1;32mOK[0m\n", msg)
	}
}

// CompareDeep4 compares two deep4 slices
func CompareDeep4(tst *testing.T, msg string, a, b [][][][]float64) {
	if len(a) != len(b) {
		fmt.Printf("%s [1;31merror len(a)=%d != len(b)=%d[0m\n", msg, len(a), len(b))
		tst.Errorf("[1;31m%s failed: slices have different lengths: %v != %v[0m", msg, a, b)
		return
	}
	for i := 0; i < len(a); i++ {
		if len(a[i]) != len(b[i]) {
			fmt.Printf("%s [1;31merror len(a[%d])=%d != len(b[%d])=%d[0m\n", msg, i, len(a[i]), i, len(b[i]))
			tst.Errorf("[1;31m%s failed: subslices have different lengths[0m", msg)
			return
		}
		for j := 0; j < len(a[i]); j++ {
			if len(a[i][j]) != len(b[i][j]) {
				fmt.Printf("%s [1;31merror len(a[%d][%d])=%d != len(b[%d][%d])=%d[0m\n", msg, i, j, len(a[i][j]), i, j, len(b[i][j]))
				tst.Errorf("[1;31m%s failed: subsubslices have different lengths[0m", msg)
				return
			}
			for k := 0; k < len(a[i][j]); k++ {
				if len(a[i][j][k]) != len(b[i][j][k]) {
					fmt.Printf("%s [1;31merror len(a[%d][%d][%d])=%d != len(b[%d][%d][%d])=%d[0m\n", msg, i, j, k, len(a[i][j][k]), i, j, k, len(b[i][j][k]))
					tst.Errorf("[1;31m%s failed: subsubsubslices have different lengths[0m", msg)
					return
				}
				for l := 0; l < len(a[i][j][k]); l++ {
					if math.Abs(a[i][j][k][l]-b[i][j][k][l]) > 1.0e-16 {
						fmt.Printf("%s [1;31merror %v != %v[0m\n", msg, a[i][j][k][l], b[i][j][k][l])
						tst.Errorf("[1;31m%s failed: slices are different: %d,%d,%d,%d component %v != %v\n%v != \n%v[0m", msg, i, j, k, l, a[i][j][k][l], b[i][j][k][l], a, b)
						return
					}
				}
			}
		}
	}
	if !Tsilent {
		fmt.Printf("%s [1;32mOK[0m\n", msg)
	}
}

// CompareBools compare two slices of bools
func CompareBools(tst *testing.T, msg string, a, b []bool) {
	if len(a) != len(b) {
		fmt.Printf("%s [1;31merror len(%q)=%d != len(%q)=%d[0m\n", msg, a, len(a), b, len(b))
		tst.Errorf("[1;31m%s failed: slices have different lengths: %v != %v[0m", msg, a, b)
		return
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			fmt.Printf("%s [1;31merror %v != %v[0m\n", msg, a[i], b[i])
			tst.Errorf("[1;31m%s failed: slices are different: %dth comp %v != %v\n%v != \n%v[0m", msg, i, a[i], b[i], a, b)
			return
		}
	}
	if !Tsilent {
		fmt.Printf("%s [1;32mOK[0m\n", msg)
	}
}

// CompareStrs compare two slices of strings
func CompareStrs(tst *testing.T, msg string, a, b []string) {
	if len(a) != len(b) {
		fmt.Printf("%s [1;31merror len(%q)=%d != len(%q)=%d[0m\n", msg, a, len(a), b, len(b))
		tst.Errorf("[1;31m%s failed: slices have different lengths: %v != %v[0m", msg, a, b)
		return
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			fmt.Printf("%s [1;31merror %v != %v[0m\n", msg, a[i], b[i])
			tst.Errorf("[1;31m%s failed: slices are different: %dth comp %v != %v\n%v != \n%v[0m", msg, i, a[i], b[i], a, b)
			return
		}
	}
	if !Tsilent {
		fmt.Printf("%s [1;32mOK[0m\n", msg)
	}
}

// CheckStrMat compares matrices of strings
func CheckStrMat(tst *testing.T, msg string, res, correct [][]string) {
	empty := false
	if len(correct) == 0 {
		empty = true
	} else {
		if len(res) != len(correct) {
			fmt.Printf("%s [1;31merror len(res)=%d != len(correct)=%d[0m\n", msg, len(res), len(correct))
			tst.Errorf("[1;31m%s failed: res and correct matrices have different lengths. %d != %d[0m", msg, len(res), len(correct))
			return
		}
	}
	for i := 0; i < len(res); i++ {
		if !empty {
			if len(res[i]) != len(correct[i]) {
				fmt.Printf("%s [1;31merror len(res[%d])=%d != len(correct[%d])=%d[0m\n", msg, i, len(res[i]), i, len(correct[i]))
				tst.Errorf("[1;31m%s failed: string matrices have different number of columns[0m", msg)
				return
			}
		}
		for j := 0; j < len(res[i]); j++ {
			var c string
			if !empty {
				c = correct[i][j]
			}
			if res[i][j] != c {
				fmt.Printf("%s [1;31merror [%d,%d] %v != %v[0m\n", msg, i, j, res[i][j], c)
				tst.Errorf("[1;31m%s failed: different str matrices:\n [%d,%d] item is wrong: %v != %v[0m", msg, i, j, res[i][j], c)
				return
			}
		}
	}
	if !Tsilent {
		fmt.Printf("%s [1;32mOK[0m\n", msg)
	}
}

// CheckIntMat compares matrices of ints
func CheckIntMat(tst *testing.T, msg string, res, correct [][]int) {
	zero := false
	if len(correct) == 0 {
		zero = true
	} else {
		if len(res) != len(correct) {
			fmt.Printf("%s [1;31merror len(res)=%d != len(correct)=%d[0m\n", msg, len(res), len(correct))
			tst.Errorf("[1;31m%s failed: res and correct matrices have different lengths. %d != %d[0m", msg, len(res), len(correct))
			return
		}
	}
	for i := 0; i < len(res); i++ {
		if !zero {
			if len(res[i]) != len(correct[i]) {
				fmt.Printf("%s [1;31merror len(res[%d])=%d != len(correct[%d])=%d[0m\n", msg, i, len(res[i]), i, len(correct[i]))
				tst.Errorf("[1;31m%s failed: matrices have different number of columns[0m", msg)
				return
			}
		}
		for j := 0; j < len(res[i]); j++ {
			var c int
			if !zero {
				c = correct[i][j]
			}
			if res[i][j] != c {
				fmt.Printf("%s [1;31merror [%d,%d] %v != %v[0m\n", msg, i, j, res[i][j], c)
				tst.Errorf("[1;31m%s failed: different int matrices:\n [%d,%d] item is wrong: %v != %v[0m", msg, i, j, res[i][j], c)
				return
			}
		}
	}
	if !Tsilent {
		fmt.Printf("%s [1;32mOK[0m\n", msg)
	}
}

// DiffVecs returns the formatted difference between two slices of floats
func DiffVecs(a, b []float64, msg string, tol float64) {
	for i := 0; i < len(a); i++ {
		fmt.Printf(msg, a[i])
		fmt.Printf(" ")
		fmt.Printf(msg, b[i])
		fmt.Printf("[1;34m =>[0m")
		diff := math.Abs(a[i] - b[i])
		if diff > tol {
			fmt.Printf("[1;31m%22.15e[0m\n", diff)
		} else {
			fmt.Printf("[1;32m%22.15e[0m\n", diff)
		}
	}
}

// max returns the max between a and b (floats)
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// min returns the min between a and b (floats)
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
