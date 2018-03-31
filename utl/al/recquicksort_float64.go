// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

// Float64RecQuickSort is a Recursive (Naive) implementation of the quicksort algorithm (Hoare version)
//
//  compare -- returns:  +1 if a > b
//                        0 if a == b
//                       -1 if a < b
//
//  Note: recursive version
//
func Float64RecQuickSort(A []float64, compare func(a, b float64) int) {
	recursiveFloat64RecQuickSort(A, 0, len(A)-1, compare)
}

// Float64RecQuickSortNonOpt is the non-optimal version of Float64RecQuickSort
func Float64RecQuickSortNonOpt(A []float64, compare func(a, b float64) int) {
	recursiveFloat64RecQuickSortNonOpt(A, 0, len(A)-1, compare)
}

// recursiveFloat64RecQuickSort performs the quick sort operations recursively
func recursiveFloat64RecQuickSort(A []float64, lo, hi int, compare func(a, b float64) int) {
	if lo < hi {
		p := partitionFloat64(A, lo, hi, compare)
		if (p - lo) < (hi - p) {
			recursiveFloat64RecQuickSort(A, lo, p, compare)
			recursiveFloat64RecQuickSort(A, p+1, hi, compare)
		} else {
			recursiveFloat64RecQuickSort(A, p+1, hi, compare)
			recursiveFloat64RecQuickSort(A, lo, p, compare)
		}
	}
}

// recursiveFloat64RecQuickSortNonOpt performs the quick sort operations recursively
func recursiveFloat64RecQuickSortNonOpt(A []float64, lo, hi int, compare func(a, b float64) int) {
	if lo < hi {
		p := partitionFloat64(A, lo, hi, compare)
		recursiveFloat64RecQuickSort(A, lo, p, compare)
		recursiveFloat64RecQuickSort(A, p+1, hi, compare)
	}
}

// partitionFloat64 partitions and modifies array by reordering the two parts to make them sorted
func partitionFloat64(A []float64, lo, hi int, compare func(a, b float64) int) int {
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
