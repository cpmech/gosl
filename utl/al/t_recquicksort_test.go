// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

// int ///////////////////////////////////////////////////////////////////////////

func TestIntRecQuickSort01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("IntRecQuickSort01. Using unsorted input")

	A := []int{0, 1, 4, 5, -1, 3, 8, 2}
	IntRecQuickSort(A, IntComparator)
	chk.Ints(tst, "A.sorted", A, []int{-1, 0, 1, 2, 3, 4, 5, 8})

	A = []int{0, 1, 4, 5, -1, 3, 8, 2}
	IntRecQuickSortNonOpt(A, IntComparator)
	chk.Ints(tst, "A.sorted (nonopt)", A, []int{-1, 0, 1, 2, 3, 4, 5, 8})
}

func TestIntRecQuickSort02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("IntRecQuickSort02. Using asc. sorted input")

	A := []int{-1, 0, 1, 2, 3, 4, 5, 8}
	IntRecQuickSort(A, IntComparator)
	chk.Ints(tst, "A.sorted", A, []int{-1, 0, 1, 2, 3, 4, 5, 8})

	A = []int{-1, 0, 1, 2, 3, 4, 5, 8}
	IntRecQuickSortNonOpt(A, IntComparator)
	chk.Ints(tst, "A.sorted (nonopt)", A, []int{-1, 0, 1, 2, 3, 4, 5, 8})
}

func TestIntRecQuickSort03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("IntRecQuickSort03. Using desc. sorted input")

	A := []int{100, 99, 988, 5, 4, 3, 2, 1, 0, -1}
	IntRecQuickSort(A, IntComparator)
	chk.Ints(tst, "A.sorted", A, []int{-1, 0, 1, 2, 3, 4, 5, 99, 100, 988})

	A = []int{100, 99, 988, 5, 4, 3, 2, 1, 0, -1}
	IntRecQuickSortNonOpt(A, IntComparator)
	chk.Ints(tst, "A.sorted (nonopt)", A, []int{-1, 0, 1, 2, 3, 4, 5, 99, 100, 988})
}

// float64 ///////////////////////////////////////////////////////////////////////

func TestFloat64RecQuickSort01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Float64RecQuickSort01. Using unsorted input")

	A := []float64{0, 1, 4, 5, -1, 3, 8, 2}
	Float64RecQuickSort(A, Float64Comparator)
	chk.Array(tst, "A.sorted", 1e-15, A, []float64{-1, 0, 1, 2, 3, 4, 5, 8})
}

func TestFloat64RecQuickSort02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Float64RecQuickSort02. Using asc. sorted input")

	A := []float64{-1, 0, 1, 2, 3, 4, 5, 8}
	Float64RecQuickSort(A, Float64Comparator)
	chk.Array(tst, "A.sorted", 1e-15, A, []float64{-1, 0, 1, 2, 3, 4, 5, 8})
}

func TestFloat64RecQuickSort03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Float64RecQuickSort03. Using desc. sorted input")

	A := []float64{100, 99, 988, 5, 4, 3, 2, 1, 0, -1}
	Float64RecQuickSort(A, Float64Comparator)
	chk.Array(tst, "A.sorted", 1e-15, A, []float64{-1, 0, 1, 2, 3, 4, 5, 99, 100, 988})
}

// string ///////////////////////////////////////////////////////////////////////

func TestStringRecQuickSort01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("StringRecQuickSort01. Using unsorted input")

	A := []string{"xyz", "bbc", "xxy", "bba", "aab", "zzx", "aba", "xxz", "ccb", "abc", "zzy", "bca", "cca", "zxy", "cab", "yzx"}
	StringRecQuickSort(A, StringComparator)
	chk.Strings(tst, "A.sorted", A, []string{"aab", "aba", "abc", "bba", "bbc", "bca", "cab", "cca", "ccb", "xxy", "xxz", "xyz", "yzx", "zxy", "zzx", "zzy"})
}

func TestStringRecQuickSort02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("StringRecQuickSort02. Using asc. sorted input")

	A := []string{"aab", "aba", "abc", "bba", "bbc", "bca", "cab", "cca", "ccb", "xxy", "xxz", "xyz", "yzx", "zxy", "zzx", "zzy"}
	StringRecQuickSort(A, StringComparator)
	chk.Strings(tst, "A.sorted", A, []string{"aab", "aba", "abc", "bba", "bbc", "bca", "cab", "cca", "ccb", "xxy", "xxz", "xyz", "yzx", "zxy", "zzx", "zzy"})
}

func TestStringRecQuickSort03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("StringRecQuickSort03. Using desc. sorted input")

	A := []string{"zzy", "zzx", "zxy", "yzx", "xyz", "xxz", "xxy", "ccb", "cca", "cab", "bca", "bbc", "bba", "abc", "aba", "aab"}
	StringRecQuickSort(A, StringComparator)
	chk.Strings(tst, "A.sorted", A, []string{"aab", "aba", "abc", "bba", "bbc", "bca", "cab", "cca", "ccb", "xxy", "xxz", "xyz", "yzx", "zxy", "zzx", "zzy"})
}
