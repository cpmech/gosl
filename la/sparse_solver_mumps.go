// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

/*
#include <dmumps_c.h>
#include <zmumps_c.h>

#define NumMaxData 64

DMUMPS_STRUC_C AllData[NumMaxData];
int NumData = 0;

ZMUMPS_STRUC_C AllDataC[NumMaxData];
int NumDataC = 0;
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// sparseSolverMumps wraps the MUMPS solver
type sparseSolverMumps struct {

	// internal
	t  *Triplet
	mi []int32
	mj []int32

	// MUMPS data
	data *C.DMUMPS_STRUC_C

	// derived
	initialized bool
	factorized  bool
}

// Init initializes mumps for sparse linear systems with real numbers
// args may be nil
func (o *sparseSolverMumps) Init(t *Triplet, args *SparseConfig) {

	// check
	if o.initialized {
		chk.Panic("solver must be initialized just once\n")
	}
	if t.pos == 0 {
		chk.Panic("triplet must have at least one item for initialization\n")
	}
	if args == nil {
		chk.Panic("the MUMPS solver requires args")
	}

	// allocate data
	if C.NumData == C.NumMaxData {
		chk.Panic("number of MUMPS data available reached. can only allocate %d structures\n", C.NumMaxData)
	}
	o.data = &C.AllData[C.NumData]
	C.NumData++

	// initialize data
	o.data.par = 1 // host also works
	o.data.sym = 0 // 0=unsymmetric, 1=sym positive definite, 2=general symmetric
	if args.symmetric {
		o.data.sym = 2
	}
	if args.symPosDef {
		o.data.sym = 1
	}
	o.data.job = -1 // initialization code
	C.dmumps_c(o.data)
	if o.data.info[1-1] < 0 {
		chk.Panic("init failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// convert indices to C.int (not C.long) and
	// increment indices since Mumps is 1-based (FORTRAN)
	o.t = t
	o.mi = make([]int32, t.pos)
	o.mj = make([]int32, t.pos)
	for k := 0; k < o.t.pos; k++ {
		o.mi[k] = int32(o.t.i[k]) + 1
		o.mj[k] = int32(o.t.j[k]) + 1
	}

	// set pointers
	o.data.n = C.int(o.t.m)
	o.data.nz = C.int(o.t.pos)
	o.data.irn = (*C.int)(unsafe.Pointer(&o.mi[0]))
	o.data.jcn = (*C.int)(unsafe.Pointer(&o.mj[0]))
	o.data.a = (*C.double)(unsafe.Pointer(&o.t.x[0]))

	// verbose level
	if args.Verbose {
		o.data.icntl[1-1] = 6 // output stream for error messages
		o.data.icntl[2-1] = 0 // output stream for statistics and warnings
		o.data.icntl[3-1] = 6 // output stream for global information
		o.data.icntl[4-1] = 2 // message level: 2==errors and warnings
	} else {
		o.data.icntl[1-1] = -1 // no output messages
		o.data.icntl[2-1] = -1 // no warnings
		o.data.icntl[3-1] = -1 // no global information
		o.data.icntl[4-1] = -1 // message level
	}

	// handle args
	o.data.icntl[7-1] = C.int(args.mumpsOrdering) // ordering
	o.data.icntl[8-1] = C.int(args.mumpsScaling)  // scaling
	o.data.icntl[14-1] = C.int(args.MumpsIncreaseOfWorkingSpacePct)
	o.data.icntl[23-1] = C.int(args.MumpsMaxMemoryPerProcessor)

	// options
	o.data.icntl[5-1] = 0  // assembled matrix (not elemental)
	o.data.icntl[6-1] = 7  // automatic col perm
	o.data.icntl[18-1] = 0 // matrix is centralized on the host
	o.data.icntl[28-1] = 1 // sequential computation
	o.data.icntl[29-1] = 0 // auto => ignored

	// analysis step
	o.data.job = 1     // analysis code
	C.dmumps_c(o.data) // analyse
	if o.data.info[1-1] < 0 {
		chk.Panic("analysis failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.initialized = true
}

// Free clears extra memory allocated by MUMPS
func (o *sparseSolverMumps) Free() {
	if o.initialized {
		o.data.job = -2    // finalisation code
		C.dmumps_c(o.data) // do finalize
	}
}

// Fact performs the factorisation
func (o *sparseSolverMumps) Fact() {

	// check
	if !o.initialized {
		chk.Panic("linear solver must be initialized first\n")
	}

	// factorisation
	o.data.job = 2     // factorisation code
	C.dmumps_c(o.data) // factorize
	if o.data.info[1-1] < 0 {
		chk.Panic("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.factorized = true
}

// Solve solves sparse linear systems using MUMPS or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
//   bIsDistr -- this flag tells that the right-hand-side vector 'b' is distributed.
//
func (o *sparseSolverMumps) Solve(x, b Vector) {

	// check
	if !o.factorized {
		chk.Panic("factorisation must be performed first\n")
	}

	// set RHS
	x.Apply(1, b) // x := b   or   copy(x, b)
	o.data.rhs = (*C.double)(unsafe.Pointer(&x[0]))

	// solve
	o.data.job = 3     // solution code
	C.dmumps_c(o.data) // solve
	if o.data.info[1-1] < 0 {
		chk.Panic("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// MumpsC wraps the MUMPS solver (complex version)
type sparseSolverMumpsC struct {

	// internal
	t  *TripletC
	mi []int32
	mj []int32

	// MUMPS data
	data *C.ZMUMPS_STRUC_C

	// derived
	initialized bool
	factorized  bool
}

// Init initializes mumps for sparse linear systems with real numbers
// args may be nil
func (o *sparseSolverMumpsC) Init(t *TripletC, args *SparseConfig) {

	// check
	if o.initialized {
		chk.Panic("solver must be initialized just once\n")
	}
	if t.pos == 0 {
		chk.Panic("triplet must have at least one item for initialization\n")
	}

	// default arguments
	if args == nil {
		args = NewSparseConfig()
	}

	// allocate data
	if C.NumDataC == C.NumMaxData {
		chk.Panic("number of MUMPS data available reached. can only allocate %d structures\n", C.NumMaxData)
	}
	o.data = &C.AllDataC[C.NumDataC]
	C.NumDataC++

	// initialize data
	o.data.comm_fortran = -987654 // use Fortran communicator by default
	o.data.par = 1                // host also works
	o.data.sym = 0                // 0=unsymmetric, 1=sym positive definite, 2=general symmetric
	if args.symmetric {
		o.data.sym = 2
	}
	if args.symPosDef {
		o.data.sym = 1
	}
	o.data.job = -1 // initialization code
	C.zmumps_c(o.data)
	if o.data.info[1-1] < 0 {
		chk.Panic("init failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// convert indices to C.int (not C.long) and
	// increment indices since Mumps is 1-based (FORTRAN)
	o.t = t
	o.mi = make([]int32, t.pos)
	o.mj = make([]int32, t.pos)
	for k := 0; k < o.t.pos; k++ {
		o.mi[k] = int32(o.t.i[k]) + 1
		o.mj[k] = int32(o.t.j[k]) + 1
	}

	// set pointers
	o.data.n = C.int(o.t.m)
	o.data.nz = C.int(o.t.pos)
	o.data.irn = (*C.int)(unsafe.Pointer(&o.mi[0]))
	o.data.jcn = (*C.int)(unsafe.Pointer(&o.mj[0]))
	o.data.a = (*C.ZMUMPS_COMPLEX)(unsafe.Pointer(&o.t.x[0]))

	// verbose level
	if args.Verbose {
		o.data.icntl[1-1] = 6 // output stream for error messages
		o.data.icntl[2-1] = 0 // output stream for statistics and warnings
		o.data.icntl[3-1] = 6 // output stream for global information
		o.data.icntl[4-1] = 2 // message level: 2==errors and warnings
	} else {
		o.data.icntl[1-1] = -1 // no output messages
		o.data.icntl[2-1] = -1 // no warnings
		o.data.icntl[3-1] = -1 // no global information
		o.data.icntl[4-1] = -1 // message level
	}

	// handle args
	o.data.icntl[7-1] = C.int(args.mumpsOrdering) // ordering
	o.data.icntl[8-1] = C.int(args.mumpsScaling)  // scaling
	o.data.icntl[14-1] = C.int(args.MumpsIncreaseOfWorkingSpacePct)
	o.data.icntl[23-1] = C.int(args.MumpsMaxMemoryPerProcessor)

	// options
	o.data.icntl[5-1] = 0  // assembled matrix (not elemental)
	o.data.icntl[6-1] = 7  // automatic col perm
	o.data.icntl[18-1] = 0 // matrix is centralized on the host
	o.data.icntl[28-1] = 1 // sequential computation
	o.data.icntl[29-1] = 0 // auto => ignored

	// analysis step
	o.data.job = 1     // analysis code
	C.zmumps_c(o.data) // analyse
	if o.data.info[1-1] < 0 {
		chk.Panic("analysis failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.initialized = true
}

// Free clears extra memory allocated by MUMPS
func (o *sparseSolverMumpsC) Free() {
	if o.initialized {
		o.data.job = -2    // finalisation code
		C.zmumps_c(o.data) // do finalize
	}
}

// Fact performs the factorisation
func (o *sparseSolverMumpsC) Fact() {

	// check
	if !o.initialized {
		chk.Panic("linear solver must be initialized first\n")
	}

	// factorisation
	o.data.job = 2     // factorisation code
	C.zmumps_c(o.data) // factorize
	if o.data.info[1-1] < 0 {
		chk.Panic("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.factorized = true
}

// Solve solves sparse linear systems using MUMPS or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
func (o *sparseSolverMumpsC) Solve(x, b VectorC) {

	// check
	if !o.factorized {
		chk.Panic("factorisation must be performed first\n")
	}

	// set RHS
	x.Apply(1, b) // x := b   or   copy(x, b)
	o.data.rhs = (*C.ZMUMPS_COMPLEX)(unsafe.Pointer(&x[0]))

	// solve
	o.data.job = 3     // solution code
	C.zmumps_c(o.data) // solve
	if o.data.info[1-1] < 0 {
		chk.Panic("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// mumErr returns error message from MUMPS
func mumErr(info, infx C.int) string {
	switch info {
	case -3:
		return "MUMPS Error # -3: solver was called with an invalid value for JOB"
	case -6:
		return "MUMPS Error # -6: singular matrix"
	case -9:
		return io.Sf("MUMPS Error # -9: main internal real/complex workarray S too small. info(2)=%v", infx)
	case -10:
		return "MUMPS Error # -10: singular matrix"
	case -13:
		return "MUMPS Error # -13: out of memory"
	case -19:
		return "MUMPS Error # -19: the maximum allowed size of working memory is too small to run the factorization"
	}
	return io.Sf("MUMPS Error # %d: unknown error", info)
}

// add solvers to database /////////////////////////////////////////////////////////////////////////

func init() {
	spSolverDB["mumps"] = func() SparseSolver { return new(sparseSolverMumps) }
	spSolverDBc["mumps"] = func() SparseSolverC { return new(sparseSolverMumpsC) }
}
