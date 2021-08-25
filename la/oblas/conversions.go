// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

import (
	"math"
	"strings"

	"github.com/cpmech/gosl/chk"
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
			k++
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
			k++
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

// complex arrays //////////////////////////////////////////////////////////////////////////////////

// GetJoinComplex joins real and imag parts of array
func GetJoinComplex(vReal, vImag []float64) (v []complex128) {
	v = make([]complex128, len(vReal))
	for i := 0; i < len(vReal); i++ {
		v[i] = complex(vReal[i], vImag[i])
	}
	return
}

// GetSplitComplex splits real and imag parts of array
func GetSplitComplex(v []complex128) (vReal, vImag []float64) {
	vReal = make([]float64, len(v))
	vImag = make([]float64, len(v))
	for i := 0; i < len(v); i++ {
		vReal[i] = real(v[i])
		vImag[i] = imag(v[i])
	}
	return
}

// JoinComplex joins real and imag parts of array
func JoinComplex(v []complex128, vReal, vImag []float64) {
	for i := 0; i < len(vReal); i++ {
		v[i] = complex(vReal[i], vImag[i])
	}
	return
}

// SplitComplex splits real and imag parts of array
func SplitComplex(vReal, vImag []float64, v []complex128) {
	for i := 0; i < len(v); i++ {
		vReal[i] = real(v[i])
		vImag[i] = imag(v[i])
	}
	return
}

// extraction //////////////////////////////////////////////////////////////////////////////////////

// ExtractRow extracts i row from (m,n) col-major matrix
func ExtractRow(i, m, n int, A []float64) (rowi []float64) {
	rowi = make([]float64, n)
	for j := 0; j < n; j++ {
		rowi[j] = A[i+j*m]
	}
	return
}

// ExtractCol extracts j column from (m,n) col-major matrix
func ExtractCol(j, m, n int, A []float64) (colj []float64) {
	colj = make([]float64, m)
	for i := 0; i < m; i++ {
		colj[i] = A[i+j*m]
	}
	return
}

// ExtractRowC extracts i row from (m,n) col-major matrix (complex version)
func ExtractRowC(i, m, n int, A []complex128) (rowi []complex128) {
	rowi = make([]complex128, n)
	for j := 0; j < n; j++ {
		rowi[j] = A[i+j*m]
	}
	return
}

// ExtractColC extracts j column from (m,n) col-major matrix (complex version)
func ExtractColC(j, m, n int, A []complex128) (colj []complex128) {
	colj = make([]complex128, m)
	for i := 0; i < m; i++ {
		colj[i] = A[i+j*m]
	}
	return
}

// eigenvector matrices ////////////////////////////////////////////////////////////////////////////

// EigenvecsBuild builds complex eigenvectors created by Dgeev function
//  INPUT:
//   wr, wi -- real and imag parts of eigenvalues
//   v      -- left or right eigenvectors from Dgeev
//  OUTPUT:
//   vv -- complex version of left or right eigenvector [pre-allocated]
//  NOTE (no checks made)
//   n = len(wr) = len(wi)
//   n * n = len(v)
//   n * n = len(vv)
func EigenvecsBuild(vv []complex128, wr, wi, v []float64) {
	n := len(wr)
	dj := 1                      // increment for next conjugate pair
	for j := 0; j < n; j += dj { // loop over columns == eigenvalues
		if math.Abs(wi[j]) > 0.0 { // eigenvalue is complex
			if j > n-2 {
				chk.Panic("last eigenvalue cannot be complex\n")
			}
			for i := 0; i < n; i++ { // loop over rows
				p := i + j*n
				q := i + (j+1)*n
				vv[p] = complex(v[p], v[q])
				vv[q] = complex(v[p], -v[q])
			}
			dj = 2
		} else {
			for i := 0; i < n; i++ { // loop over rows
				p := i + j*n
				vv[p] = complex(v[p], 0.0)
			}
			dj = 1
		}
	}
}

// EigenvecsBuildBoth builds complex left and right eigenvectors created by Dgeev function
//  INPUT:
//   wr, wi -- real and imag parts of eigenvalues
//   vl, vr -- left and right eigenvectors from Dgeev
//  OUTPUT:
//   vvl, vvr -- complex version of left and right eigenvectors [pre-allocated]
//  NOTE (no checks made)
//   n = len(wr) = len(wi)
//   n * n = len(vl) = len(vr)
//   n * n = len(vvl) = len(vvr)
func EigenvecsBuildBoth(vvl, vvr []complex128, wr, wi, vl, vr []float64) {
	n := len(wr)
	dj := 1                      // increment for next conjugate pair
	for j := 0; j < n; j += dj { // loop over columns == eigenvalues
		if math.Abs(wi[j]) > 0.0 { // eigenvalue is complex
			if j > n-2 {
				chk.Panic("last eigenvalue cannot be complex\n")
			}
			for i := 0; i < n; i++ { // loop over rows
				p := i + j*n
				q := i + (j+1)*n
				vvl[p] = complex(vl[p], vl[q])
				vvr[p] = complex(vr[p], vr[q])
				vvl[q] = complex(vl[p], -vl[q])
				vvr[q] = complex(vr[p], -vr[q])
			}
			dj = 2
		} else {
			for i := 0; i < n; i++ { // loop over rows
				p := i + j*n
				vvl[p] = complex(vl[p], 0.0)
				vvr[p] = complex(vr[p], 0.0)
			}
			dj = 1
		}
	}
}
