// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import "github.com/cpmech/gosl/chk"

// MatToColMaj returns a vector representation of a column-major matrix
func MatToColMaj(a [][]float64) (am []float64) {
	m, n := len(a), len(a[0])
	am = make([]float64, m*n) // column-major matrix
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			am[i+j*m] = a[i][j]
		}
	}
	return
}

// ColMajToMatNew returns a new matrix from a column-major representation
func ColMajToMatNew(am []float64, m, n int) (a [][]float64) {
	a = make([][]float64, m)
	for i := 0; i < m; i++ {
		a[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			a[i][j] = am[i+j*m]
		}
	}
	return
}

// ColMajToMat returns a matrix from a column-major representation
func ColMajToMat(a [][]float64, am []float64) {
	m, n := len(a), len(a[0])
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			a[i][j] = am[i+j*m]
		}
	}
	return
}

// RCtoComplex converts two slices into a slice of complex numbers
func RCtoComplex(r, c []float64) (rc []complex128) {
	if len(r) != len(c) {
		chk.Panic("length of r and c slices must be equal to each other. %d != %d", len(r), len(c))
	}
	rc = make([]complex128, len(r))
	for i := 0; i < len(r); i++ {
		rc[i] = complex(r[i], c[i])
	}
	return
}

// ComplexToRC converts a slice of complex numbers into two slices
func ComplexToRC(rc []complex128) (r, c []float64) {
	r = make([]float64, len(rc))
	c = make([]float64, len(rc))
	for i := 0; i < len(rc); i++ {
		r[i] = real(rc[i])
		c[i] = imag(rc[i])
	}
	return
}

// RCpairsToComplex converts pairs such as [real,imag, real,imag, real,imag, ...] to slice of complex numbers
//   Note: len(pairs) must be even
func RCpairsToComplex(pairs []float64) (v []complex128) {
	n := len(pairs)
	if n%2 != 0 {
		chk.Panic("len(pairs) must be even. %d is invalid", n)
	}
	v = make([]complex128, n/2)
	for i := 0; i < n/2; i++ {
		v[i] = complex(pairs[i*2], pairs[i*2+1])
	}
	return
}

// ComplexToRCpairs converts slice of complex numbers to pairs such as [real,imag, real,imag, real,imag, ...]
func ComplexToRCpairs(v []complex128) (pairs []float64) {
	n := len(v)
	pairs = make([]float64, n*2)
	for i := 0; i < n; i++ {
		pairs[i*2] = real(v[i])
		pairs[i*2+1] = imag(v[i])
	}
	return
}
