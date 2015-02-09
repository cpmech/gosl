// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

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
