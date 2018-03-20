// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"github.com/cpmech/gosl/chk"
)

// NaiveQuickSort is a Naive implementation of the quicksort algorithm (Hoare version)
//
//  compare -- returns:  +1 if a > b
//                        0 if a == b
//                       -1 if a < b
//
//  Note: recursive version
//
//  TODO: replace []int with []Adata
//
func NaiveQuickSort(A []int, compare func(a, b int) int) {
	naiveQuickSort(A, 0, len(A)-1, compare)
}

// naiveQuickSort performs the quick sort operations recursively
func naiveQuickSort(A []int, lo, hi int, compare func(a, b int) int) {
	if lo < hi {
		p := partition(A, lo, hi, compare)
		naiveQuickSort(A, lo, p, compare)
		naiveQuickSort(A, p+1, hi, compare)
	}
}

// partition partitions and modifies array by reordering the two parts to make them sorted
func partition(A []int, lo, hi int, compare func(a, b int) int) int {
	pivot := A[lo]
	i := lo - 1
	j := hi + 1
	for {
		for {
			i++
			if A[i] >= pivot {
				break
			}
		}
		for {
			j--
			if A[j] <= pivot {
				break
			}
		}
		if i >= j {
			return j
		}
		A[i], A[j] = A[j], A[i]
	}
	chk.Panic("partition failed\n")
	return 0 // unreachable
}
