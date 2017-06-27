// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine,!heroku

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
	"time"
	"unsafe"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// LinSolUmfpack holds UMFPACK data
type LinSolUmfpack struct {
	linSolData

	// umfpack data
	info  []float64
	ctrl  []float64
	uinfo *C.double
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

	// derived
	is_initialised bool
	factorised     bool
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
		return chk.Err("triplet must have at least one item before calling this method\n")
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

	// control and info
	o.ctrl = make([]float64, C.UMFPACK_CONTROL)
	o.uctrl = (*C.double)(unsafe.Pointer(&o.ctrl[0]))
	C.umfpack_dl_defaults(o.uctrl)
	if o.verb {
		o.info = make([]float64, C.UMFPACK_INFO)
		o.uinfo = (*C.double)(unsafe.Pointer(&o.info[0]))
		o.ctrl[C.UMFPACK_PRL] = 6 // change the default print level; otherwise, nothing will print
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.InitR = %v\n", o.name, time.Now().Sub(o.tini))
	}

	// success
	o.is_initialised = true
	return
}

// InitC initialises a LinSolUmfpack data structure for Complex systems. It also performs some initial analyses.
func (o *LinSolUmfpack) InitC(tC *TripletC, symmetric, verbose, timing bool) (err error) {

	// check
	o.tC = tC
	if tC.pos == 0 {
		return chk.Err("triplet must have at least one item before calling this method\n")
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
		return chk.Err("length of x (%d) and z (%d) slices must be equal to the length of i (%d) slice when using UMFPACK complex solver\n", len(o.tC.x), len(o.tC.z), len(o.tC.i))
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

	// success
	o.is_initialised = true
	return
}

// Fact performs symbolic/numeric factorisation. This method also converts the triplet form
// to the column-compressed form, including the summation of duplicated entries
func (o *LinSolUmfpack) Fact() (err error) {

	// check
	if !o.is_initialised {
		return chk.Err("linear solver must be initialised first\n")
	}

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// message
	if o.verb {
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolUmfpack.Fact . . . . . . . . . . . . . . . \n\n")
	}

	// free memory
	if o.factorised {
		o.Free()
	}

	// factorisation
	if o.cmplx {

		// UMFPACK: convert triplet to column-compressed format
		st := C.umfpack_zl_triplet_to_col(C.LONG(o.tC.m), C.LONG(o.tC.n), C.LONG(o.tC.pos), o.ti, o.tj, o.tx, o.tz, o.ap, o.ai, o.ax, o.az, nil)
		if st != C.UMFPACK_OK {
			return chk.Err("umfpack_zl_triplet_to_col failed (UMFPACK error: %s)\n", Uerr2Text[int(st)])
		}

		// UMFPACK: symbolic factorisation
		st = C.umfpack_zl_symbolic(C.LONG(o.tC.m), C.LONG(o.tC.n), o.ap, o.ai, o.ax, o.az, &o.usymb, o.uctrl, nil)
		if st != C.UMFPACK_OK {
			return chk.Err("umfpack_zl_symbolic failed (UMFPACK error: %s)\n", Uerr2Text[int(st)])
		}

		// UMFPACK: numeric factorisation
		st = C.umfpack_zl_numeric(o.ap, o.ai, o.ax, o.az, o.usymb, &o.unum, o.uctrl, nil)
		if st != C.UMFPACK_OK {
			return chk.Err("umfpack_zl_numeric failed (UMFPACK error: %s)\n", Uerr2Text[int(st)])
		}

	} else {

		// UMFPACK: convert triplet to column-compressed format
		st := C.umfpack_dl_triplet_to_col(C.LONG(o.tR.m), C.LONG(o.tR.n), C.LONG(o.tR.pos), o.ti, o.tj, o.tx, o.ap, o.ai, o.ax, nil)
		if st != C.UMFPACK_OK {
			return chk.Err("umfpack_dl_triplet_to_col failed (UMFPACK error: %s)\n", Uerr2Text[int(st)])
		}

		// UMFPACK: symbolic factorisation
		st = C.umfpack_dl_symbolic(C.LONG(o.tR.m), C.LONG(o.tR.n), o.ap, o.ai, o.ax, &o.usymb, o.uctrl, o.uinfo)
		if st != C.UMFPACK_OK {
			return chk.Err("umfpack_dl_symbolic failed (UMFPACK error: %s)\n", Uerr2Text[int(st)])
		}
		if o.verb {
			C.umfpack_dl_report_info(o.uctrl, o.uinfo)
		}

		// UMFPACK: numeric factorisation
		st = C.umfpack_dl_numeric(o.ap, o.ai, o.ax, o.usymb, &o.unum, o.uctrl, o.uinfo)
		if st != C.UMFPACK_OK {
			return chk.Err("umfpack_dl_numeric failed (UMFPACK error: %s)\n", Uerr2Text[int(st)])
		}
		if o.verb {
			C.umfpack_dl_report_info(o.uctrl, o.uinfo)
		}
	}

	// set flag
	o.factorised = true

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.Fact  = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// SolveR solves the linear Real system A.x = b
func (o *LinSolUmfpack) SolveR(xR, bR Vector, dummy bool) (err error) {

	// check
	if !o.is_initialised {
		return chk.Err("linear solver must be initialised first\n")
	}
	if o.cmplx {
		return chk.Err("this method must be called with Real matrices\n")
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
	st := C.umfpack_dl_solve(C.UMFPACK_A, o.ap, o.ai, o.ax, pxR, pbR, o.unum, o.uctrl, o.uinfo)
	if st != C.UMFPACK_OK {
		return chk.Err("umfpack_dl_solve failed (UMFPACK error: %s)\n", Uerr2Text[int(st)])
	}
	if o.verb {
		C.umfpack_dl_report_info(o.uctrl, o.uinfo)
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.Solve = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// SolveC solves the linear Complex system A.x = b
func (o *LinSolUmfpack) SolveC(xR, xC, bR, bC Vector, dummy bool) (err error) {

	// check
	if !o.is_initialised {
		return chk.Err("linear solver must be initialised first\n")
	}
	if !o.cmplx {
		return chk.Err("this method must be called with Complex matrices\n")
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
		chk.Err("umfpack_zl_solve failed (UMFPACK error: %s)\n", Uerr2Text[int(st)])
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.Solve = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// Free deletes temporary data structures
func (o *LinSolUmfpack) Free() {

	// exit if not initialised
	if !o.is_initialised {
		return
	}

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// message
	if o.verb {
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolUmfpack.Free . . . . . . . . . . . . . . . \n\n")
	}

	// free memory
	if o.cmplx {
		C.umfpack_zl_free_symbolic(&o.usymb)
		C.umfpack_zl_free_numeric(&o.unum)
	} else {
		C.umfpack_dl_free_symbolic(&o.usymb)
		C.umfpack_dl_free_numeric(&o.unum)
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolUmfpack.Free   = %v\n", o.name, time.Now().Sub(o.tini))
	}
}

// SetOrdScal sets the ordering and scaling methods
//  Note: this method is not available for UMFPACK
func (o *LinSolUmfpack) SetOrdScal(ordering, scaling string) (err error) {
	return
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
