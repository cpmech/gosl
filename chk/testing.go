// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"fmt"
	"math"
	"testing"
)

// Scalar compares two scalars
func Scalar(tst *testing.T, msg string, tol, res, correct float64) {
	CheckAndPrint(tst, msg, tol, math.Abs(res-correct))
}

// AnaNum compares analytical versus numerical values
func AnaNum(tst *testing.T, msg string, tol, ana, num float64, verbose bool) {
	e := PrintAnaNum(msg, tol, ana, num, verbose)
	if e != nil {
		tst.Errorf("%v", e.Error())
	}
}

// String compares two strings
func String(tst *testing.T, str, correct string) {
	if str != correct {
		fmt.Printf("[1;31merror %q != %q[0m\n", str, correct)
		tst.Errorf("[1;31mstring failed with: %q != %q[0m", str, correct)
		return
	}
	PrintOk(fmt.Sprintf("%s == %s", str, correct))
}

// Int compares two ints
func Int(tst *testing.T, msg string, val, correct int) {
	if val != correct {
		fmt.Printf("[1;31m%s: error %d != %d[0m\n", msg, val, correct)
		tst.Errorf("[1;31m%s failed with: %d != %d[0m", msg, val, correct)
		return
	}
	PrintOk(fmt.Sprintf("%s: %d == %d", msg, val, correct))
}

// Ints compares two slices of integers
func Ints(tst *testing.T, msg string, a, b []int) {
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
	PrintOk(msg)
}

// Bools compare two slices of bools
func Bools(tst *testing.T, msg string, a, b []bool) {
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
	PrintOk(msg)
}

// Strings compare two slices of strings
func Strings(tst *testing.T, msg string, a, b []string) {
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
	PrintOk(msg)
}

// Matrix compares two matrices
func Matrix(tst *testing.T, msg string, tol float64, res, correct [][]float64) {
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
	CheckAndPrint(tst, msg, tol, maxdiff)
}

// Vector compares two vectors
func Vector(tst *testing.T, msg string, tol float64, res, correct []float64) {
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
	CheckAndPrint(tst, msg, tol, maxdiff)
}

// MatrixC compares two matrices of complex nummbers
func MatrixC(tst *testing.T, msg string, tol float64, res, correct [][]complex128) {
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
		return
	}
	PrintOk(msg)
}

// VectorC compares two vectors of complex nummbers
func VectorC(tst *testing.T, msg string, tol float64, res, correct []complex128) {
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
		return
	}
	PrintOk(msg)
}

// StrMat compares matrices of strings
func StrMat(tst *testing.T, msg string, res, correct [][]string) {
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
	PrintOk(msg)
}

// IntMat compares matrices of ints
func IntMat(tst *testing.T, msg string, res, correct [][]int) {
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
	PrintOk(msg)
}

// Deep3 compares two deep3 slices
func Deep3(tst *testing.T, msg string, tol float64, a, b [][][]float64) {
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
				if math.Abs(a[i][j][k]-b[i][j][k]) > tol {
					fmt.Printf("%s [1;31merror %v != %v[0m\n", msg, a[i][j][k], b[i][j][k])
					tst.Errorf("[1;31m%s failed: slices are different: %d,%d,%d component %v != %v\n%v != \n%v[0m", msg, i, j, k, a[i][j][k], b[i][j][k], a, b)
					return
				}
			}
		}
	}
	PrintOk(msg)
}

// Deep4 compares two deep4 slices
func Deep4(tst *testing.T, msg string, tol float64, a, b [][][][]float64) {
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
					if math.Abs(a[i][j][k][l]-b[i][j][k][l]) > tol {
						fmt.Printf("%s [1;31merror %v != %v[0m\n", msg, a[i][j][k][l], b[i][j][k][l])
						tst.Errorf("[1;31m%s failed: slices are different: %d,%d,%d,%d component %v != %v\n%v != \n%v[0m", msg, i, j, k, l, a[i][j][k][l], b[i][j][k][l], a, b)
						return
					}
				}
			}
		}
	}
	PrintOk(msg)
}
