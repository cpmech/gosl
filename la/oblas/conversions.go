// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

import (
	"strings"

	"github.com/cpmech/gosl/io"
)

// SliceToColMajor converts nested slice into an array representing a col-major matrix
//
//   Example:
//             _      _
//            |  0  3  |
//        a = |  1  4  |            ⇒     data = [0, 1, 2, 3, 4, 5]
//            |_ 2  5 _|(m x n)
//
//        data[i+j*m] = a[i][j]
//
//   NOTE: make sure to have at least 1x1 item
func SliceToColMajor(a [][]float64) (data []float64) {
	m, n := len(a), len(a[0])
	data = make([]float64, m*n)
	k := 0
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			data[k] = a[i][j]
			k += 1
		}
	}
	return
}

// ColMajorToSlice converts col-major matrix to nested slice
func ColMajorToSlice(m, n int, data []float64) (a [][]float64) {
	a = make([][]float64, m)
	for i := 0; i < m; i++ {
		a[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			a[i][j] = data[i+j*m]
		}
	}
	return
}

// PrintColMajor prints matrix (without commas or brackets)
func PrintColMajor(m, n int, data []float64, nfmt string) (l string) {
	if nfmt == "" {
		nfmt = "%g "
	}
	for i := 0; i < m; i++ {
		if i > 0 {
			l += "\n"
		}
		for j := 0; j < n; j++ {
			l += io.Sf(nfmt, data[i+j*m])
		}
	}
	return
}

// PrintColMajorGo prints matrix in Go format
func PrintColMajorGo(m, n int, data []float64, nfmt string) (l string) {
	if nfmt == "" {
		nfmt = "%10g"
	}
	l = "[][]float64{\n"
	for i := 0; i < m; i++ {
		l += "    {"
		for j := 0; j < n; j++ {
			if j > 0 {
				l += ","
			}
			l += io.Sf(nfmt, data[i+j*m])
		}
		l += "},\n"
	}
	l += "}"
	return
}

// PrintColMajorPy prints matrix in Python format
func PrintColMajorPy(m, n int, data []float64, nfmt string) (l string) {
	if nfmt == "" {
		nfmt = "%10g"
	}
	l = "np.matrix([\n"
	for i := 0; i < m; i++ {
		l += "    ["
		for j := 0; j < n; j++ {
			if j > 0 {
				l += ","
			}
			l += io.Sf(nfmt, data[i+j*m])
		}
		l += "],\n"
	}
	l += "], dtype=float)"
	return
}

// complex ///////////////////////////////////////////////////////////////////////////////////////

// SliceToColMajorC converts nested slice into an array representing a col-major matrix of
// complex numbers.
//
//   Example:
//             _            _
//            |  0+0i  3+3i  |
//        a = |  1+1i  4+4i  |          ⇒   data = [0+0i, 1+1i, 2+2i, 3+3i, 4+4i, 5+5i]
//            |_ 2+2i  5+5i _|(m x n)
//
//        data[i+j*m] = a[i][j]
//
//   NOTE: make sure to have at least 1x1 item
func SliceToColMajorC(a [][]complex128) (data []complex128) {
	m, n := len(a), len(a[0])
	data = make([]complex128, m*n)
	k := 0
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			data[k] = a[i][j]
			k += 1
		}
	}
	return
}

// ColMajorCtoSlice converts col-major matrix to nested slice
func ColMajorCtoSlice(m, n int, data []complex128) (a [][]complex128) {
	a = make([][]complex128, m)
	for i := 0; i < m; i++ {
		a[i] = make([]complex128, n)
		for j := 0; j < n; j++ {
			a[i][j] = data[i+j*m]
		}
	}
	return
}

// PrintColMajorC prints matrix (without commas or brackets).
// NOTE: if non-empty, nfmtI must have '+' e.g. %+g
func PrintColMajorC(m, n int, data []complex128, nfmtR, nfmtI string) (l string) {
	if nfmtR == "" {
		nfmtR = "%g"
	}
	if nfmtI == "" {
		nfmtI = "%+g"
	}
	if !strings.ContainsRune(nfmtI, '+') {
		nfmtI = strings.Replace(nfmtI, "%", "%+", -1)
	}
	for i := 0; i < m; i++ {
		if i > 0 {
			l += "\n"
		}
		for j := 0; j < n; j++ {
			if j > 0 {
				l += ", "
			}
			v := data[i+j*m]
			l += io.Sf(nfmtR, real(v)) + io.Sf(nfmtI, imag(v)) + "i"
		}
	}
	return
}

// PrintColMajorCgo prints matrix in Go format
// NOTE: if non-empty, nfmtI must have '+' e.g. %+g
func PrintColMajorCgo(m, n int, data []complex128, nfmtR, nfmtI string) (l string) {
	if nfmtR == "" {
		nfmtR = "%g"
	}
	if nfmtI == "" {
		nfmtI = "%+g"
	}
	if !strings.ContainsRune(nfmtI, '+') {
		nfmtI = strings.Replace(nfmtI, "%", "%+", -1)
	}
	l = "[][]complex128{\n"
	for i := 0; i < m; i++ {
		l += "    {"
		for j := 0; j < n; j++ {
			if j > 0 {
				l += ","
			}
			v := data[i+j*m]
			l += io.Sf(nfmtR, real(v)) + io.Sf(nfmtI, imag(v)) + "i"
		}
		l += "},\n"
	}
	l += "}"
	return
}

// PrintColMajorCpy prints matrix in Python format
// NOTE: if non-empty, nfmtI must have '+' e.g. %+g
func PrintColMajorCpy(m, n int, data []complex128, nfmtR, nfmtI string) (l string) {
	if nfmtR == "" {
		nfmtR = "%g"
	}
	if nfmtI == "" {
		nfmtI = "%+g"
	}
	if !strings.ContainsRune(nfmtI, '+') {
		nfmtI = strings.Replace(nfmtI, "%", "%+", -1)
	}
	l = "np.matrix([\n"
	for i := 0; i < m; i++ {
		l += "    ["
		for j := 0; j < n; j++ {
			if j > 0 {
				l += ","
			}
			v := data[i+j*m]
			l += io.Sf(nfmtR, real(v)) + io.Sf(nfmtI, imag(v)) + "j"
		}
		l += "],\n"
	}
	l += "], dtype=complex)"
	return
}
