// Copyright 2015 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func printThickLine(n int) (l string) {
	for i := 0; i < n; i++ {
		l += "="
	}
	return l + "\n"
}

func printThinLine(n int) (l string) {
	for i := 0; i < n; i++ {
		l += "-"
	}
	return l + "\n"
}

func printSpaces(n int) (l string) {
	for i := 0; i < n; i++ {
		l += " "
	}
	return
}
