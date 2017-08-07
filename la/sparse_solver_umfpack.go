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
	"github.com/cpmech/gosl/mpi"
)

// Umfpack wraps the UMFPACK solver
type Umfpack struct {

	// umfpack data
	info  []float64
	ctrl  []float64
	uinfo *C.double
	uctrl *C.double
	usymb unsafe.Pointer
	unum  unsafe.Pointer

	// data
	apData []int
	aiData []int
	axData []float64

	// pointers
	t  *Triplet
	ti *C.LONG
	tj *C.LONG
	tx *C.double
	ap *C.LONG
	ai *C.LONG
	ax *C.double

	// derived
	initialised bool
	factorised  bool
	symbFact    bool
	numeFact    bool
}

// Init initialises umfpack for sparse linear systems with real numbers
func (o *Umfpack) Init(t *Triplet, symmetric, verbose bool, ordering, scaling string, dummy *mpi.Communicator) (err error) {

	// check
	if t.pos == 0 {
		return chk.Err("triplet must have at least one item for initialisation\n")
	}

	// allocate data
	o.apData = make([]int, t.n+1)
	o.aiData = make([]int, t.pos)
	o.axData = make([]float64, t.pos)

	// pointers
	o.t = t
	o.ti = (*C.LONG)(unsafe.Pointer(&t.i[0]))
	o.tj = (*C.LONG)(unsafe.Pointer(&t.j[0]))
	o.tx = (*C.double)(unsafe.Pointer(&t.x[0]))
	o.ap = (*C.LONG)(unsafe.Pointer(&o.apData[0]))
	o.ai = (*C.LONG)(unsafe.Pointer(&o.aiData[0]))
	o.ax = (*C.double)(unsafe.Pointer(&o.axData[0]))

	// info and control
	o.info = make([]float64, C.UMFPACK_INFO)
	o.ctrl = make([]float64, C.UMFPACK_CONTROL)
	o.uinfo = (*C.double)(unsafe.Pointer(&o.info[0]))
	o.uctrl = (*C.double)(unsafe.Pointer(&o.ctrl[0]))
	C.umfpack_dl_defaults(o.uctrl)

	// flags
	if symmetric {
		o.ctrl[C.UMFPACK_STRATEGY] = C.UMFPACK_STRATEGY_SYMMETRIC
	}
	if verbose {
		o.ctrl[C.UMFPACK_PRL] = 6
	}

	// success
	o.initialised = true
	return
}

// Free clears extra memory allocated by UMFPACK
func (o *Umfpack) Free() {
	if o.symbFact {
		C.umfpack_dl_free_symbolic(&o.usymb)
		o.symbFact = false
	}
	if o.numeFact {
		C.umfpack_dl_free_numeric(&o.unum)
		o.numeFact = false
	}
}

// Fact performs the factorisation
func (o *Umfpack) Fact() (err error) {

	// check
	if !o.initialised {
		return chk.Err("linear solver must be initialised first\n")
	}
	o.factorised = false

	// convert triplet to column-compressed format
	code := C.umfpack_dl_triplet_to_col(C.LONG(o.t.m), C.LONG(o.t.n), C.LONG(o.t.pos), o.ti, o.tj, o.tx, o.ap, o.ai, o.ax, nil)
	if code != C.UMFPACK_OK {
		return chk.Err("conversion failed (UMFPACK error: %s)\n", umfErr(code))
	}

	// symbolic factorisation
	if o.symbFact {
		C.umfpack_dl_free_symbolic(&o.usymb)
		o.symbFact = false
	}
	code = C.umfpack_dl_symbolic(C.LONG(o.t.m), C.LONG(o.t.n), o.ap, o.ai, o.ax, &o.usymb, o.uctrl, o.uinfo)
	if code != C.UMFPACK_OK {
		return chk.Err("symbolic factorised failed (UMFPACK error: %s)\n", umfErr(code))
	}
	o.symbFact = true

	// numeric factorisation
	if o.numeFact {
		C.umfpack_dl_free_numeric(&o.unum)
		o.numeFact = false
	}
	code = C.umfpack_dl_numeric(o.ap, o.ai, o.ax, o.usymb, &o.unum, o.uctrl, o.uinfo)
	if code != C.UMFPACK_OK {
		return chk.Err("numeric factorisation failed (UMFPACK error: %s)\n", umfErr(code))
	}
	o.numeFact = true

	// success
	o.factorised = true
	return
}

// Solve solves sparse linear systems using UMFPACK or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
func (o *Umfpack) Solve(x, b Vector, dummy bool) (err error) {

	// check
	if !o.factorised {
		return chk.Err("factorisation must be performed first\n")
	}

	// pointers
	px := (*C.double)(unsafe.Pointer(&x[0]))
	pb := (*C.double)(unsafe.Pointer(&b[0]))

	// solve
	code := C.umfpack_dl_solve(C.UMFPACK_A, o.ap, o.ai, o.ax, px, pb, o.unum, o.uctrl, o.uinfo)
	if code != C.UMFPACK_OK {
		err = chk.Err("solve failed (UMFPACK error: %s)\n", umfErr(code))
	}
	return
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// UmfpackC wraps the UMFPACK solver (complex version)
type UmfpackC struct {

	// umfpack data
	info  []float64
	ctrl  []float64
	uinfo *C.double
	uctrl *C.double
	usymb unsafe.Pointer
	unum  unsafe.Pointer

	// data
	apData []int
	aiData []int
	axData []float64

	// pointers
	t  *TripletC
	ti *C.LONG
	tj *C.LONG
	tx *C.double
	ap *C.LONG
	ai *C.LONG
	ax *C.double

	// derived
	initialised bool
	factorised  bool
	symbFact    bool
	numeFact    bool
}

// Init initialises umfpack for sparse linear systems with real numbers
func (o *UmfpackC) Init(t *TripletC, symmetric, verbose bool, ordering, scaling string, dummy *mpi.Communicator) (err error) {

	// check
	if t.pos == 0 {
		return chk.Err("triplet must have at least one item for initialisation\n")
	}

	// allocate data
	o.apData = make([]int, t.n+1)
	o.aiData = make([]int, t.pos)
	o.axData = make([]float64, t.pos)

	// pointers
	o.t = t
	o.ti = (*C.LONG)(unsafe.Pointer(&t.i[0]))
	o.tj = (*C.LONG)(unsafe.Pointer(&t.j[0]))
	o.tx = (*C.double)(unsafe.Pointer(&t.x[0]))
	o.ap = (*C.LONG)(unsafe.Pointer(&o.apData[0]))
	o.ai = (*C.LONG)(unsafe.Pointer(&o.aiData[0]))
	o.ax = (*C.double)(unsafe.Pointer(&o.axData[0]))

	// info and control
	o.info = make([]float64, C.UMFPACK_INFO)
	o.ctrl = make([]float64, C.UMFPACK_CONTROL)
	o.uinfo = (*C.double)(unsafe.Pointer(&o.info[0]))
	o.uctrl = (*C.double)(unsafe.Pointer(&o.ctrl[0]))
	C.umfpack_zl_defaults(o.uctrl)

	// flags
	if symmetric {
		o.ctrl[C.UMFPACK_STRATEGY] = C.UMFPACK_STRATEGY_SYMMETRIC
	}
	if verbose {
		o.ctrl[C.UMFPACK_PRL] = 6
	}

	// success
	o.initialised = true
	return
}

// Free clears extra memory allocated by UMFPACK
func (o *UmfpackC) Free() {
	if o.symbFact {
		C.umfpack_zl_free_symbolic(&o.usymb)
		o.symbFact = false
	}
	if o.numeFact {
		C.umfpack_zl_free_numeric(&o.unum)
		o.numeFact = false
	}
}

// Fact performs the factorisation
func (o *UmfpackC) Fact() (err error) {

	// check
	if !o.initialised {
		return chk.Err("linear solver must be initialised first\n")
	}
	o.factorised = false

	// convert triplet to column-compressed format
	code := C.umfpack_zl_triplet_to_col(C.LONG(o.t.m), C.LONG(o.t.n), C.LONG(o.t.pos), o.ti, o.tj, o.tx, nil, o.ap, o.ai, o.ax, nil, nil)
	if code != C.UMFPACK_OK {
		return chk.Err("conversion failed (UMFPACK error: %s)\n", umfErr(code))
	}

	// symbolic factorisation
	if o.symbFact {
		C.umfpack_zl_free_symbolic(&o.usymb)
		o.symbFact = false
	}
	code = C.umfpack_zl_symbolic(C.LONG(o.t.m), C.LONG(o.t.n), o.ap, o.ai, o.ax, nil, &o.usymb, o.uctrl, o.uinfo)
	if code != C.UMFPACK_OK {
		return chk.Err("symbolic factorised failed (UMFPACK error: %s)\n", umfErr(code))
	}
	o.symbFact = true

	// numeric factorisation
	if o.numeFact {
		C.umfpack_zl_free_numeric(&o.unum)
		o.numeFact = false
	}
	code = C.umfpack_zl_numeric(o.ap, o.ai, o.ax, nil, o.usymb, &o.unum, o.uctrl, o.uinfo)
	if code != C.UMFPACK_OK {
		return chk.Err("numeric factorisation failed (UMFPACK error: %s)\n", umfErr(code))
	}
	o.numeFact = true

	// success
	o.factorised = true
	return
}

// Solve solves sparse linear systems using UMFPACK or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
func (o *UmfpackC) Solve(x, b VectorC, dummy bool) (err error) {

	// check
	if !o.factorised {
		return chk.Err("factorisation must be performed first\n")
	}

	// pointers
	px := (*C.double)(unsafe.Pointer(&x[0]))
	pb := (*C.double)(unsafe.Pointer(&b[0]))

	// solve
	code := C.umfpack_zl_solve(C.UMFPACK_A, o.ap, o.ai, o.ax, nil, px, nil, pb, nil, o.unum, o.uctrl, o.uinfo)
	if code != C.UMFPACK_OK {
		err = chk.Err("solve failed (UMFPACK error: %s)\n", umfErr(code))
	}
	return
}

// add solvers to database /////////////////////////////////////////////////////////////////////////

func init() {
	spSolverDB["umfpack"] = func() SparseSolver { return new(Umfpack) }
	spSolverDBc["umfpack"] = func() SparseSolverC { return new(UmfpackC) }
}
