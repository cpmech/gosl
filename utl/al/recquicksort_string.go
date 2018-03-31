// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

// StringRecQuickSort is a Recursive (Naive) implementation of the quicksort algorithm (Hoare version)
//
//  compare -- returns:  +1 if a > b
//                        0 if a == b
//                       -1 if a < b
//
//  Note: recursive version
//
func StringRecQuickSort(A []string, compare func(a, b string) int) {
	recursiveStringRecQuickSort(A, 0, len(A)-1, compare)
}

// StringRecQuickSortNonOpt is the non-optimal version of StringRecQuickSort
func StringRecQuickSortNonOpt(A []string, compare func(a, b string) int) {
	recursiveStringRecQuickSortNonOpt(A, 0, len(A)-1, compare)
}

// recursiveStringRecQuickSort performs the quick sort operations recursively
func recursiveStringRecQuickSort(A []string, lo, hi int, compare func(a, b string) int) {
	if lo < hi {
		p := partitionString(A, lo, hi, compare)
		if (p - lo) < (hi - p) {
			recursiveStringRecQuickSort(A, lo, p, compare)
			recursiveStringRecQuickSort(A, p+1, hi, compare)
		} else {
			recursiveStringRecQuickSort(A, p+1, hi, compare)
			recursiveStringRecQuickSort(A, lo, p, compare)
		}
	}
}

// recursiveStringRecQuickSortNonOpt performs the quick sort operations recursively
func recursiveStringRecQuickSortNonOpt(A []string, lo, hi int, compare func(a, b string) int) {
	if lo < hi {
		p := partitionString(A, lo, hi, compare)
		recursiveStringRecQuickSort(A, lo, p, compare)
		recursiveStringRecQuickSort(A, p+1, hi, compare)
	}
}

// partitionString partitions and modifies array by reordering the two parts to make them sorted
func partitionString(A []string, lo, hi int, compare func(a, b string) int) int {
	pivot := A[lo]
	i := lo - 1
	j := hi + 1
	for {
		for {
			i++
			if compare(A[i], pivot) >= 0 { // A[i] >= pivot
				break
			}
		}
		for {
			j--
			if compare(A[j], pivot) <= 0 { // A[j] <= pivot
				break
			}
		}
		if i >= j {
			return j
		}
		A[i], A[j] = A[j], A[i]
	}
}
