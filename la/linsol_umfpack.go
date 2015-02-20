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
	"time"
	"unsafe"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// LinSolUmfpack holds UMFPACK data
type LinSolUmfpack struct {
	linSolData

	// umfpack data
	uctrl *C.double
	usymb unsafe.Pointer
	unum  unsafe.Pointer

	// pointers
	ti     *C.LONG
	tj     *C.LONG
	tx, tz *C.double
	ap     *C.LONG
	ai     *C.LONG
	ax, az *C.double
}

// factory of allocators
func init() {
	lsAllocators["umfpack"] = func() LinSol { return new(LinSolUmfpack) }
}

// InitR initialises a LinSolUmfpack data structure for Real systems. It also performs some initial analyses.
func (o *LinSolUmfpack) InitR(tR *Triplet, symmetric, verbose, timing bool) (err error) {

	// check
	o.tR = tR
	if tR.pos == 0 {
		return chk.Err(_linsol_umfpack_err01)
	}

	// flags
	o.name = "umfpack"
	o.sym = symmetric
	o.cmplx = false
	o.verb = verbose
	o.ton = timing

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// pointers
	o.ti = (*C.LONG)(unsafe.Pointer(&o.tR.i[0]))
	o.tj = (*C.LONG)(unsafe.Pointer(&o.tR.j[0]))
	o.tx = (*C.double)(unsafe.Pointer(&o.tR.x[0]))
	o.ap = (*C.LONG)(unsafe.Pointer(&make([]int, o.tR.n+1)[0]))
	o.ai = (*C.LONG)(unsafe.Pointer(&make([]int, o.tR.pos)[0]))
	o.ax = (*C.double)(unsafe.Pointer(&make([]float64, o.tR.pos)[0]))

	// control
	o.uctrl = (*C.double)(unsafe.Pointer(&make([]float64, C.UMFPACK_CONTROL)[0]))
	C.umfpack_dl_defaults(o.uctrl)

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.InitR = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// InitC initialises a LinSolUmfpack data structure for Complex systems. It also performs some initial analyses.
func (o *LinSolUmfpack) InitC(tC *TripletC, symmetric, verbose, timing bool) (err error) {

	// check
	o.tC = tC
	if tC.pos == 0 {
		return chk.Err(_linsol_umfpack_err02)
	}

	// flags
	o.name = "umfpack"
	o.sym = symmetric
	o.cmplx = true
	o.verb = verbose
	o.ton = timing

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// check x and z
	if len(o.tC.x) != len(o.tC.i) || len(o.tC.z) != len(o.tC.i) {
		return chk.Err(_linsol_umfpack_err03, len(o.tC.x), len(o.tC.z), len(o.tC.i))
	}

	// pointers
	o.ti = (*C.LONG)(unsafe.Pointer(&o.tC.i[0]))
	o.tj = (*C.LONG)(unsafe.Pointer(&o.tC.j[0]))
	o.tx = (*C.double)(unsafe.Pointer(&o.tC.x[0]))
	o.tz = (*C.double)(unsafe.Pointer(&o.tC.z[0]))
	o.ap = (*C.LONG)(unsafe.Pointer(&make([]int, o.tC.n+1)[0]))
	o.ai = (*C.LONG)(unsafe.Pointer(&make([]int, o.tC.pos)[0]))
	o.ax = (*C.double)(unsafe.Pointer(&make([]float64, o.tC.pos)[0]))
	o.az = (*C.double)(unsafe.Pointer(&make([]float64, o.tC.pos)[0]))

	// control
	o.uctrl = (*C.double)(unsafe.Pointer(&make([]float64, C.UMFPACK_CONTROL)[0]))
	C.umfpack_zl_defaults(o.uctrl)

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.InitC = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// Fact performs symbolic/numeric factorisation. This method also converts the triplet form
// to the column-compressed form, including the summation of duplicated entries
func (o *LinSolUmfpack) Fact() (err error) {

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// message
	if o.verb {
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolUmfpack.Fact . . . . . . . . . . . . . . . \n\n")
	}

	// factorisation
	if o.cmplx {

		// UMFPACK: convert triplet to column-compressed format
		st := C.umfpack_zl_triplet_to_col(C.LONG(o.tC.m), C.LONG(o.tC.n), C.LONG(o.tC.pos), o.ti, o.tj, o.tx, o.tz, o.ap, o.ai, o.ax, o.az, nil)
		if st != C.UMFPACK_OK {
			return chk.Err(_linsol_umfpack_err04, Uerr2Text[int(st)])
		}

		// UMFPACK: symbolic factorisation
		st = C.umfpack_zl_symbolic(C.LONG(o.tC.m), C.LONG(o.tC.n), o.ap, o.ai, o.ax, o.az, &o.usymb, o.uctrl, nil)
		if st != C.UMFPACK_OK {
			return chk.Err(_linsol_umfpack_err05, Uerr2Text[int(st)])
		}

		// UMFPACK: numeric factorisation
		st = C.umfpack_zl_numeric(o.ap, o.ai, o.ax, o.az, o.usymb, &o.unum, o.uctrl, nil)
		if st != C.UMFPACK_OK {
			return chk.Err(_linsol_umfpack_err06, Uerr2Text[int(st)])
		}

	} else {

		// UMFPACK: convert triplet to column-compressed format
		st := C.umfpack_dl_triplet_to_col(C.LONG(o.tR.m), C.LONG(o.tR.n), C.LONG(o.tR.pos), o.ti, o.tj, o.tx, o.ap, o.ai, o.ax, nil)
		if st != C.UMFPACK_OK {
			return chk.Err(_linsol_umfpack_err07, Uerr2Text[int(st)])
		}

		// UMFPACK: symbolic factorisation
		st = C.umfpack_dl_symbolic(C.LONG(o.tR.m), C.LONG(o.tR.n), o.ap, o.ai, o.ax, &o.usymb, o.uctrl, nil)
		if st != C.UMFPACK_OK {
			return chk.Err(_linsol_umfpack_err08, Uerr2Text[int(st)])
		}

		// UMFPACK: numeric factorisation
		st = C.umfpack_dl_numeric(o.ap, o.ai, o.ax, o.usymb, &o.unum, o.uctrl, nil)
		if st != C.UMFPACK_OK {
			return chk.Err(_linsol_umfpack_err09, Uerr2Text[int(st)])
		}

		return

	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.Fact  = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// SolveR solves the linear Real system A.x = b
func (o *LinSolUmfpack) SolveR(xR, bR []float64, dummy bool) (err error) {

	// check
	if o.cmplx {
		return chk.Err(_linsol_umfpack_err10)
	}

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// message
	if o.verb {
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolUmfpack.SolveR . . . . . . . . . . . . . . . \n\n")
	}

	// UMFPACK: pointers
	pxR := (*C.double)(unsafe.Pointer(&xR[0]))
	pbR := (*C.double)(unsafe.Pointer(&bR[0]))

	// UMFPACK: solve
	st := C.umfpack_dl_solve(C.UMFPACK_A, o.ap, o.ai, o.ax, pxR, pbR, o.unum, o.uctrl, nil)
	if st != C.UMFPACK_OK {
		return chk.Err(_linsol_umfpack_err11, Uerr2Text[int(st)])
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.Solve = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// SolveC solves the linear Complex system A.x = b
func (o *LinSolUmfpack) SolveC(xR, xC, bR, bC []float64, dummy bool) (err error) {

	// check
	if !o.cmplx {
		return chk.Err(_linsol_umfpack_err12)
	}

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// message
	if o.verb {
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolUmfpack.SolveC . . . . . . . . . . . . . . . \n\n")
	}

	// UMFPACK: pointers
	pxR := (*C.double)(unsafe.Pointer(&xR[0]))
	pxC := (*C.double)(unsafe.Pointer(&xC[0]))
	pbR := (*C.double)(unsafe.Pointer(&bR[0]))
	pbC := (*C.double)(unsafe.Pointer(&bC[0]))

	// UMFPACK: solve
	st := C.umfpack_zl_solve(C.UMFPACK_A, o.ap, o.ai, o.ax, o.az, pxR, pxC, pbR, pbC, o.unum, o.uctrl, nil)
	if st != C.UMFPACK_OK {
		chk.Err(_linsol_umfpack_err13, Uerr2Text[int(st)])
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.Solve = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// Clean deletes temporary data structures
func (o *LinSolUmfpack) Clean() {

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// message
	if o.verb {
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolUmfpack.Clean . . . . . . . . . . . . . . . \n\n")
	}

	// clean up
	if o.cmplx {
		C.umfpack_zl_free_symbolic(&o.usymb)
		C.umfpack_zl_free_numeric(&o.unum)
	} else {
		C.umfpack_dl_free_symbolic(&o.usymb)
		C.umfpack_dl_free_numeric(&o.unum)
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.Clean   = %v\n", o.name, time.Now().Sub(o.tini))
	}
}

// SetOrdScal sets the ordering and scaling methods for MUMPS
func (o *LinSolUmfpack) SetOrdScal(ordering, scaling string) (err error) {
	return
}

// variables required by Umfpack
var (
	_umfpctrl  [C.UMFPACK_CONTROL]float64 // Umfpack control array
	_umfpctrlz [C.UMFPACK_CONTROL]float64 // Umfpack control array for complex matrices
	_uctrl     *C.double                  // pointer to Umfpack control
	_uctrlz    *C.double                  // pointer to Umfpack control for complex routines
)

// initialise Umfpack control arrays
func init() {
	_uctrl = (*C.double)(unsafe.Pointer(&_umfpctrl[0]))
	_uctrlz = (*C.double)(unsafe.Pointer(&_umfpctrlz[0]))
	C.umfpack_dl_defaults(_uctrl)
	C.umfpack_zl_defaults(_uctrlz)
}

// Umfpack error codes
var (
	Uerr2Text = map[int]string{
		C.UMFPACK_ERROR_out_of_memory:           "out_of_memory (-1)",
		C.UMFPACK_ERROR_invalid_Numeric_object:  "invalid_Numeric_object (-3)",
		C.UMFPACK_ERROR_invalid_Symbolic_object: "invalid_Symbolic_object (-4)",
		C.UMFPACK_ERROR_argument_missing:        "argument_missing (-5)",
		C.UMFPACK_ERROR_n_nonpositive:           "n_nonpositive (-6)",
		C.UMFPACK_ERROR_invalid_matrix:          "invalid_matrix (-8)",
		C.UMFPACK_ERROR_different_pattern:       "different_pattern (-11)",
		C.UMFPACK_ERROR_invalid_system:          "invalid_system (-13)",
		C.UMFPACK_ERROR_invalid_permutation:     "invalid_permutation (-15)",
		C.UMFPACK_ERROR_internal_error:          "internal_error (-911)",
		C.UMFPACK_ERROR_file_IO:                 "file_IO (-17)",
		-18: "ordering_failed (-18)",
		C.UMFPACK_WARNING_singular_matrix:       "singular_matrix (1)",
		C.UMFPACK_WARNING_determinant_underflow: "determinant_underflow (2)",
		C.UMFPACK_WARNING_determinant_overflow:  "determinant_overflow (3)",
	}
)

// error messages
var (
	_linsol_umfpack_err01 = "linsol_umfpack.go: InitR: triplet must have at least one item before calling this method\n"
	_linsol_umfpack_err02 = "linsol_umfpack.go: InitC: triplet must have at least one item before calling this method\n"
	_linsol_umfpack_err03 = "linsol_umfpack.go: InitC: length of x (%d) and z (%d) slices must be equal to the length of i (%d) slice when using UMFPACK complex solver\n"
	_linsol_umfpack_err04 = "linsol_umfpack.go: Fact: umfpack_zl_triplet_to_col failed (UMFPACK error: %s)\n"
	_linsol_umfpack_err05 = "linsol_umfpack.go: Fact: umfpack_zl_symbolic failed (UMFPACK error: %s)\n"
	_linsol_umfpack_err06 = "linsol_umfpack.go: Fact: umfpack_zl_numeric failed (UMFPACK error: %s)\n"
	_linsol_umfpack_err07 = "linsol_umfpack.go: Fact: umfpack_dl_triplet_to_col failed (UMFPACK error: %s)\n"
	_linsol_umfpack_err08 = "linsol_umfpack.go: Fact: umfpack_dl_symbolic failed (UMFPACK error: %s)\n"
	_linsol_umfpack_err09 = "linsol_umfpack.go: Fact: umfpack_dl_numeric failed (UMFPACK error: %s)\n"
	_linsol_umfpack_err10 = "linsol_umfpack.go: SolveR: this method must be called with Real matrices\n"
	_linsol_umfpack_err11 = "linsol_umfpack.go: SolveR: umfpack_dl_solve failed (UMFPACK error: %s)\n"
	_linsol_umfpack_err12 = "linsol_umfpack.go: SolveC: this method must be called with Complex matrices\n"
	_linsol_umfpack_err13 = "linsol_umfpack.go: SolveC: umfpack_zl_solve failed (UMFPACK error: %s)\n"
)
