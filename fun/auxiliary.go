// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

// max returns the max between two floats
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// min returns the min between two floats
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// imax returns the max between two ints
func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// imin returns the min between two ints
func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// iabs performs the absolute operation with ints
func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
