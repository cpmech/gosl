// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/chk"

// imax returns the max between two integers
func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// imin returns the min between two integers
func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
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

// checkprm checks if parameter value is correct
func checkprm(name string, val, minval, maxval float64, usemin, usemax bool) (err error) {
	if usemin {
		if val < minval {
			return chk.Err("%q parameter: wrong value: %g < %g", name, val, minval)
		}
	}
	if usemax {
		if val > maxval {
			return chk.Err("%q parameter: wrong value: %g > %g", name, val, maxval)
		}
	}
	return
}

// setvzero sets v := 0
func setvzero(v []float64) {
	for i := 0; i < len(v); i++ {
		v[i] = 0
	}
}
