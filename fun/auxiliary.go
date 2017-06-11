// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

// IsPowerOfTwo checks if n is power of 2; i.e. 2⁰, 2¹, 2², 2³, 2⁴, ...
func IsPowerOfTwo(n int) bool {
	if n < 1 {
		return false
	}
	return n&(n-1) == 0
}

// Swap swaps two float64 numbers
func Swap(a, b *float64) {
	*a, *b = *b, *a
}

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

// setvzero sets v := 0
func setvzero(v []float64) {
	for i := 0; i < len(v); i++ {
		v[i] = 0
	}
}
