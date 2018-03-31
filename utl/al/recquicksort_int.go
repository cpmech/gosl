// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

// IntRecQuickSort is a Recursive (Naive) implementation of the quicksort algorithm (Hoare version)
//
//  compare -- returns:  +1 if a > b
//                        0 if a == b
//                       -1 if a < b
//
//  Note: recursive version
//
func IntRecQuickSort(A []int, compare func(a, b int) int) {
	recursiveIntRecQuickSort(A, 0, len(A)-1, compare)
}

// IntRecQuickSortNonOpt is the non-optimal version of IntRecQuickSort
func IntRecQuickSortNonOpt(A []int, compare func(a, b int) int) {
	recursiveIntRecQuickSortNonOpt(A, 0, len(A)-1, compare)
}

// recursiveIntRecQuickSort performs the quick sort operations recursively
func recursiveIntRecQuickSort(A []int, lo, hi int, compare func(a, b int) int) {
	if lo < hi {
		p := partitionInt(A, lo, hi, compare)
		if (p - lo) < (hi - p) {
			recursiveIntRecQuickSort(A, lo, p, compare)
			recursiveIntRecQuickSort(A, p+1, hi, compare)
		} else {
			recursiveIntRecQuickSort(A, p+1, hi, compare)
			recursiveIntRecQuickSort(A, lo, p, compare)
		}
	}
}

// recursiveIntRecQuickSortNonOpt performs the quick sort operations recursively
func recursiveIntRecQuickSortNonOpt(A []int, lo, hi int, compare func(a, b int) int) {
	if lo < hi {
		p := partitionInt(A, lo, hi, compare)
		recursiveIntRecQuickSort(A, lo, p, compare)
		recursiveIntRecQuickSort(A, p+1, hi, compare)
	}
}

// partitionInt partitions and modifies array by reordering the two parts to make them sorted
func partitionInt(A []int, lo, hi int, compare func(a, b int) int) int {
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
