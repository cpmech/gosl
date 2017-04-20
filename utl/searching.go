// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"sort"
	"time"

	"github.com/cpmech/gosl/chk"
)

// IntMinMax returns the maximum and minimum elements in v
//  NOTE: this is not efficient and should be used for small slices only
func IntMinMax(v []int) (mi, ma int) {
	mi, ma = v[0], v[0]
	for i := 1; i < len(v); i++ {
		if v[i] < mi {
			mi = v[i]
		}
		if v[i] > ma {
			ma = v[i]
		}
	}
	return
}

// MinMax returns the maximum and minimum elements in v
//  NOTE: this is not efficient and should be used for small slices only
func MinMax(v []float64) (mi, ma float64) {
	mi, ma = v[0], v[0]
	for i := 1; i < len(v); i++ {
		if v[i] < mi {
			mi = v[i]
		}
		if v[i] > ma {
			ma = v[i]
		}
	}
	return
}

// Sum sums all items in v
//  NOTE: this is not efficient and should be used for small slices only
func Sum(v []float64) (sum float64) {
	for i := 0; i < len(v); i++ {
		sum += v[i]
	}
	return
}

// DurSum sums all seconds in v
//  NOTE: this is not efficient and should be used for small slices only
func DurSum(v []time.Duration) (seconds float64) {
	for _, t := range v {
		seconds += t.Seconds()
	}
	return
}

// StrIndexSmall finds the index of an item in a slice of strings
//  NOTE: this function is not efficient and should be used with small slices only; say smaller than 20
func StrIndexSmall(a []string, val string) int {
	for idx, str := range a {
		if str == val {
			return idx
		}
	}
	return -1 // not found
}

// IntIndexSmall finds the index of an item in a slice of ints
//  NOTE: this function is not efficient and should be used with small slices only; say smaller than 20
func IntIndexSmall(a []int, val int) int {
	for idx, item := range a {
		if item == val {
			return idx
		}
	}
	return -1 // not found
}

// IntIndexSorted finds the index of an item in a slice of ints by means of
// searching a sorted list (given or computed here)
//  Input:
//   a         -- the list
//   val       -- value to be searched for
//   aSortedIn -- an existing sorted 'a' list or <nil>
//  Output:
//   idx        -- position in 'a' where val is located if found; otherwise returns -1
//   aSortedOut -- sorted list if the input aSortedIn is <nil>
//  Note: If you were to sort the array with MergeSort or any other O(nlogn) algorithm then
//        the complexity would be O(nlogn).
//        O(logn) < O(n) < O(nlogn)
//
//  TODO: this function is finished yet
//
func IntIndexSorted(a []int, val int, aSortedIn []int) (idx int, aSortedOut []int) {
	chk.Panic("this function is not implemented yet")
	aSortedOut = aSortedIn
	if aSortedIn == nil {
		aSortedOut = make([]int, len(a))
		copy(aSortedOut, a)
		sort.Ints(aSortedOut)
	}
	idx = -1 // not found
	idxInSorted := sort.SearchInts(aSortedOut, val)
	if idxInSorted < len(aSortedOut) && aSortedOut[idxInSorted] == val { // found in sorted list
		// TODO: convert index in sorted list into index in original list
		return
	}
	return
}

// IntFilter filters out components in slice
//  NOTE: this function is not efficient and should be used with small slices only
func IntFilter(a []int, out func(idx int) bool) (res []int) {
	for i := 0; i < len(a); i++ {
		if out(i) {
			continue
		}
		res = append(res, a[i])
	}
	return
}

// IntNegOut filters out negative components in slice
//  NOTE: this function is not efficient and should be used with small slices only
func IntNegOut(a []int) (res []int) {
	for i := 0; i < len(a); i++ {
		if a[i] < 0 {
			continue
		}
		res = append(res, a[i])
	}
	return
}
