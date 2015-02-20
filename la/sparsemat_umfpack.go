// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine,!heroku

package la

/*
#cgo linux   CFLAGS: -O3 -I/usr/include/suitesparse
#cgo windows CFLAGS: -O3 -IC:/Gosl/include
#cgo linux   LDFLAGS: -llapack -lgfortran -lblas -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig
#cgo windows LDFLAGS: -llapack -lgfortran -lblas -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig -LC:/Gosl/lib
#ifdef WIN32
#define LONG long long
#else
#define LONG long
#endif
#include <umfpack.h>
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

// ToMatrix converts a sparse matrix in triplet form to column-compressed form using Umfpack's
// routines. "realloc_a" indicates whether the internal "a" matrix must be reallocated or not,
// for instance, in case the structure of the triplet has changed.
//  INPUT:
//   a -- a previous CCMatrix to be filled in; otherwise, "nil" tells to allocate a new one
//  OUTPUT:
//   the previous "a" matrix or a pointer to a new one
func (t *Triplet) ToMatrix(a *CCMatrix) *CCMatrix {
	if t.pos < 1 {
		chk.Panic(_sparsemat_umfpack_err1, t.pos)
	}
	if a == nil {
		a = new(CCMatrix)
		a.m, a.n, a.nnz = t.m, t.n, t.pos
		a.p = make([]int, a.n+1)
		a.i = make([]int, a.nnz)
		a.x = make([]float64, a.nnz)
	}
	Ti := (*C.LONG)(unsafe.Pointer(&t.i[0]))
	Tj := (*C.LONG)(unsafe.Pointer(&t.j[0]))
	Tx := (*C.double)(unsafe.Pointer(&t.x[0]))
	Ap := (*C.LONG)(unsafe.Pointer(&a.p[0]))
	Ai := (*C.LONG)(unsafe.Pointer(&a.i[0]))
	Ax := (*C.double)(unsafe.Pointer(&a.x[0]))
	status := C.umfpack_dl_triplet_to_col(C.LONG(a.m), C.LONG(a.n), C.LONG(a.nnz), Ti, Tj, Tx, Ap, Ai, Ax, nil)
	if status != C.UMFPACK_OK {
		chk.Panic(_sparsemat_umfpack_err2, Uerr2Text[int(status)])
	}
	return a
}

// ToMatrix converts a sparse matrix in triplet form with complex numbers to column-compressed form.
// "realloc_a" indicates whether the internal "a" matrix must be reallocated or not, for instance,
// in case the structure of the triplet has changed.
//  INPUT:
//   a -- a previous CCMatrixC to be filled in; otherwise, "nil" tells to allocate a new one
//  OUTPUT:
//   the previous "a" matrix or a pointer to a new one
func (t *TripletC) ToMatrix(a *CCMatrixC) *CCMatrixC {
	if t.pos < 1 {
		chk.Panic(_sparsemat_umfpack_err3, t.pos)
	}
	if a == nil {
		a = new(CCMatrixC)
		a.m, a.n, a.nnz = t.m, t.n, t.pos
		a.p = make([]int, a.n+1)
		a.i = make([]int, a.nnz)
		a.x = make([]float64, a.nnz)
		a.z = make([]float64, a.nnz)
	}
	Ap := (*C.LONG)(unsafe.Pointer(&a.p[0]))
	Ai := (*C.LONG)(unsafe.Pointer(&a.i[0]))
	Ax := (*C.double)(unsafe.Pointer(&a.x[0]))
	Az := (*C.double)(unsafe.Pointer(&a.z[0]))
	Ti := (*C.LONG)(unsafe.Pointer(&t.i[0]))
	Tj := (*C.LONG)(unsafe.Pointer(&t.j[0]))
	var Tx, Tz *C.double
	if t.xz != nil {
		x := make([]float64, t.pos)
		z := make([]float64, t.pos)
		for k := 0; k < t.pos; k++ {
			x[k], z[k] = t.xz[k*2], t.xz[k*2+1]
		}
		Tx = (*C.double)(unsafe.Pointer(&x[0]))
		Tz = (*C.double)(unsafe.Pointer(&z[0]))
	} else {
		Tx = (*C.double)(unsafe.Pointer(&t.x[0]))
		Tz = (*C.double)(unsafe.Pointer(&t.z[0]))
	}
	status := C.umfpack_zl_triplet_to_col(C.LONG(a.m), C.LONG(a.n), C.LONG(a.nnz), Ti, Tj, Tx, Tz, Ap, Ai, Ax, Az, nil)
	if status != C.UMFPACK_OK {
		chk.Panic(_sparsemat_umfpack_err4, Uerr2Text[int(status)])
	}
	return a
}

// error messages
var (
	_sparsemat_umfpack_err1 = "sparsemat_umfpack.go: la.Triplet.ToMatrix: conversion can only be made for non-empty triplets. error: (pos = %d)"
	_sparsemat_umfpack_err2 = "sparsemat_umfpack.go: la.Triplet.ToMatrix: umfpack_dl_triplet_to_col failed (UMFPACK error: %s)"
	_sparsemat_umfpack_err3 = "sparsemat_umfpack.go: la.TripletC.ToMatrix: conversion can only be made for non-empty triplets. error: (pos = %d)"
	_sparsemat_umfpack_err4 = "sparsemat_umfpack.go: la.TripletC.ToMatrix: umfpack_zl_triplet_to_col failed (UMFPACK error: %s)"
)
