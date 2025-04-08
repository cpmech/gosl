// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

/*
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

// ToMatrix converts a sparse matrix in triplet form to column-compressed form using Umfpack's routines.
//
//	INPUT:
//	 a -- a previous CCMatrix to be filled in; otherwise, "nil" tells to allocate a new one
//	OUTPUT:
//	 the previous "a" matrix or a pointer to a new one
func (t *Triplet) ToMatrix(a *CCMatrix) *CCMatrix {
	if t.pos < 1 {
		chk.Panic("conversion can only be made for non-empty triplets. error: (pos = %d)", t.pos)
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
		chk.Panic("umfpack_dl_triplet_to_col failed (UMFPACK error: %s)", umfErr((int)(status)))
	}
	return a
}

// ToMatrix converts a sparse matrix in triplet form with complex numbers to column-compressed form.
//
//	INPUT:
//	 a -- a previous CCMatrixC to be filled in; otherwise, "nil" tells to allocate a new one
//	OUTPUT:
//	 the previous "a" matrix or a pointer to a new one
func (t *TripletC) ToMatrix(a *CCMatrixC) *CCMatrixC {
	if t.pos < 1 {
		chk.Panic("conversion can only be made for non-empty triplets. error: (pos = %d)", t.pos)
	}
	if a == nil {
		a = new(CCMatrixC)
		a.m, a.n, a.nnz = t.m, t.n, t.pos
		a.p = make([]int, a.n+1)
		a.i = make([]int, a.nnz)
		a.x = make([]complex128, a.nnz)
	}
	Ap := (*C.LONG)(unsafe.Pointer(&a.p[0]))
	Ai := (*C.LONG)(unsafe.Pointer(&a.i[0]))
	Ax := (*C.double)(unsafe.Pointer(&a.x[0]))
	Ti := (*C.LONG)(unsafe.Pointer(&t.i[0]))
	Tj := (*C.LONG)(unsafe.Pointer(&t.j[0]))
	Tx := (*C.double)(unsafe.Pointer(&t.x[0]))
	status := C.umfpack_zl_triplet_to_col(C.LONG(a.m), C.LONG(a.n), C.LONG(a.nnz), Ti, Tj, Tx, nil, Ap, Ai, Ax, nil, nil)
	if status != C.UMFPACK_OK {
		chk.Panic("umfpack_zl_triplet_to_col failed (UMFPACK error: %s)", umfErr((int)(status)))
	}
	return a
}

// umfErr returns UMFPACK error codes
func umfErr(code int) string {
	switch code {
	case C.UMFPACK_ERROR_out_of_memory:
		return "out_of_memory (-1)"
	case C.UMFPACK_ERROR_invalid_Numeric_object:
		return "invalid_Numeric_object (-3)"
	case C.UMFPACK_ERROR_invalid_Symbolic_object:
		return "invalid_Symbolic_object (-4)"
	case C.UMFPACK_ERROR_argument_missing:
		return "argument_missing (-5)"
	case C.UMFPACK_ERROR_n_nonpositive:
		return "n_nonpositive (-6)"
	case C.UMFPACK_ERROR_invalid_matrix:
		return "invalid_matrix (-8)"
	case C.UMFPACK_ERROR_different_pattern:
		return "different_pattern (-11)"
	case C.UMFPACK_ERROR_invalid_system:
		return "invalid_system (-13)"
	case C.UMFPACK_ERROR_invalid_permutation:
		return "invalid_permutation (-15)"
	case C.UMFPACK_ERROR_internal_error:
		return "internal_error (-911)"
	case C.UMFPACK_ERROR_file_IO:
		return "file_IO (-17)"
	case -18:
		return "ordering_failed (-18)"
	case C.UMFPACK_WARNING_singular_matrix:
		return "singular_matrix (1)"
	case C.UMFPACK_WARNING_determinant_underflow:
		return "determinant_underflow (2)"
	case C.UMFPACK_WARNING_determinant_overflow:
		return "determinant_overflow (3)"
	}
	return "unknown UMFPACK error"
}
