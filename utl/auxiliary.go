// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

// BestSquare finds the best square for given size=Nrows * Ncolumns
func BestSquare(size int) (nrow, ncol int) {
	nrow = -1 // not found
	ncol = -1 // not found
	for x := 1; x <= size; x++ {
		if (x * x) >= size {
			if (x * x) == size {
				nrow = x
				ncol = x
				return
			} else {
				for y := x; y >= 1; y-- {
					if (x * y) == size {
						nrow = x
						ncol = y
						return
					}
				}
			}
		}
	}
	return
}

// Imin returns the minimum between two integers
func Imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Imin returns the maximum between two integers
func Imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum between two float point numbers
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum between two float point numbers
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
