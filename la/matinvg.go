// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine,!heroku

package la

/*
#cgo CFLAGS: -O3
#cgo linux   LDFLAGS: -lm -llapack -lgfortran -lblas
#cgo darwin  LDFLAGS: -lm -llapack -lblas
#include "connectlapack.h"
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

//  NOTE: __IMPORTANT__ this package works on 64bit machines only,
//                      since Go's int is converted to C's long

// MatSvd returns the singular value decomposition of a matrix such that:
//  Note:
//  a = u * s * v
//  u[m][m], s[n], vt[n][n] must be pre-allocated
func MatSvd(u [][]float64, s []float64, vt [][]float64, a [][]float64, tol float64) (err error) {
	if len(a) < 1 {
		return chk.Err(_matinvg_err1)
	}
	m, n := len(a), len(a[0])
	um, vtm := make([]float64, m*m), make([]float64, n*n)
	am := MatToColMaj(a)
	status := C.lapack_svd((*C.double)(unsafe.Pointer(&um[0])), (*C.double)(unsafe.Pointer(&s[0])), (*C.double)(unsafe.Pointer(&vtm[0])), (C.long)(m), (C.long)(n), (*C.double)(unsafe.Pointer(&am[0])))
	if status != 0 {
		return chk.Err(_matinvg_err3, m, n, "lapack_svd", status)
	}
	ColMajToMat(u, um)
	ColMajToMat(vt, vtm)
	return
}

// MatInvG returns the matrix inverse of 'a' in 'ai'. 'a' can be of any size,
// even non-square; in this case, the pseudo-inverse is returned
func MatInvG(ai, a [][]float64, tol float64) (err error) {
	if len(a) < 1 {
		return chk.Err(_matinvg_err1)
	}
	m, n := len(a), len(a[0])
	if m == n && m < 4 { // call simple function
		_, err = MatInv(ai, a, tol)
		return
	}
	am := MatToColMaj(a)
	ami := make([]float64, n*m) // column-major inverse matrix
	if m == n {                 // general matrix inverse
		status := C.lapack_square_inverse((*C.double)(unsafe.Pointer(&ami[0])), (C.long)(m), (*C.double)(unsafe.Pointer(&am[0])))
		if status != 0 {
			return chk.Err(_matinvg_err2, m, n, "lapack_square_inverse", status)
		}
	} else { // pseudo inverse
		status := C.lapack_pseudo_inverse((*C.double)(unsafe.Pointer(&ami[0])), (C.long)(m), (C.long)(n), (*C.double)(unsafe.Pointer(&am[0])), (C.double)(tol))
		if status != 0 {
			return chk.Err(_matinvg_err2, m, n, "lapack_pseudo_inverse", status)
		}
	}
	ColMajToMat(ai, ami)
	return
}

// MatCondG returns the condition number of a square matrix using the inverse of this matrix; thus
// it is not as efficient as it could be, e.g. by using the SV decomposition.
//  normtype -- Type of norm to use:
//    "F" or "" => Frobenius
//    "I"       => Infinite
func MatCondG(a [][]float64, normtype string, tol float64) (res float64, err error) {
	if len(a) < 1 {
		return 0, chk.Err(_matinvg_err1)
	}
	m, n := len(a), len(a[0])
	ai := MatAlloc(m, n)
	err = MatInvG(ai, a, tol)
	if err != nil {
		return 0, chk.Err(_matinvg_err4, err.Error())
	}
	if normtype == "I" {
		res = MatNormI(a) * MatNormI(ai)
	} else {
		res = MatNormF(a) * MatNormF(ai)
	}
	return
}

// error messages
const (
	_matinvg_err1 = "cannot handle nil matrix. len(a) = 0"
	_matinvg_err2 = "inverse of (%d x %d) matrix failed with %s.status = %d"
	_matinvg_err3 = "SV decomposition of (%d x %d) matrix failed with %s.status = %d"
	_matinvg_err4 = "cannot compute inverse:\n%v"
)
