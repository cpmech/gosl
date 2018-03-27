// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

// Comparator is a function that compares two entities in order to position them in a sorted list
//  returns:  +1 if a > b
//             0 if a == b
//            -1 if a < b

// IntComparator compares ints
func IntComparator(a, b int) int {
	if a > b {
		return +1
	}
	if a < b {
		return -1
	}
	return 0
}

// Float64Comparator compares float64s
func Float64Comparator(a, b float64) int {
	if a > b {
		return +1
	}
	if a < b {
		return -1
	}
	return 0
}

// StringComparator compares strings lexicographically (lexicographical order)
func StringComparator(a, b string) int {
	if a > b {
		return +1
	}
	if a < b {
		return -1
	}
	return 0
}
