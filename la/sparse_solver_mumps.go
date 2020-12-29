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
	"github.com/cpmech/gosl/mpi"
)

// sparseSolverMumps wraps the MUMPS solver
type sparseSolverMumps struct {

	// internal
	comm *mpi.Communicator
	t    *Triplet
	mi   []int32
	mj   []int32

	// MUMPS data
	data *C.DMUMPS_STRUC_C

	// derived
	initialised bool
	factorised  bool
}

// Init initialises mumps for sparse linear systems with real numbers
// args may be nil
func (o *sparseSolverMumps) Init(t *Triplet, args *SparseConfig) {

	// check
	if o.initialised {
		chk.Panic("solver must be initialised just once\n")
	}
	if t.pos == 0 {
		chk.Panic("triplet must have at least one item for initialisation\n")
	}
	if args == nil {
		chk.Panic("the MUMPS solver requires args with an MPI communicator")
	}
	if args.communicator == nil {
		chk.Panic("the MUMPS solver requires an MPI communicator")
	}

	// set comm
	o.comm = args.communicator

	// allocate data
	if C.NumData == C.NumMaxData {
		chk.Panic("number of MUMPS data available reached. can only allocate %d structures\n", C.NumMaxData)
	}
	o.data = &C.AllData[C.NumData]
	C.NumData++

	// initialise data
	o.data.comm_fortran = -987654 // use Fortran communicator by default
	o.data.par = 1                // host also works
	o.data.sym = 0                // 0=unsymmetric, 1=sym positive definite, 2=general symmetric
	if args.symmetric {
		o.data.sym = 2
	}
	if args.symPosDef {
		o.data.sym = 1
	}
	o.data.job = -1 // initialisation code
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
	o.data.nz_loc = C.int(o.t.pos)
	o.data.irn_loc = (*C.int)(unsafe.Pointer(&o.mi[0]))
	o.data.jcn_loc = (*C.int)(unsafe.Pointer(&o.mj[0]))
	o.data.a_loc = (*C.double)(unsafe.Pointer(&o.t.x[0]))

	// control
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
	o.data.icntl[5-1] = 0  // assembled matrix (needed for distributed matrix)
	o.data.icntl[6-1] = 0  // must be zero because matrix is distributed
	o.data.icntl[18-1] = 3 // distributed matrix
	o.data.icntl[14-1] = C.int(args.MumpsIncreaseOfWorkingSpacePct)
	o.data.icntl[23-1] = C.int(args.MumpsMaxMemoryPerProcessor)

	// set ordering and scaling
	o.data.icntl[7-1] = C.int(args.mumpsOrdering) // ordering
	o.data.icntl[8-1] = C.int(args.mumpsScaling)  // scaling

	// analysis step
	o.data.job = 1     // analysis code
	C.dmumps_c(o.data) // analyse
	if o.data.info[1-1] < 0 {
		chk.Panic("analysis failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.initialised = true
}

// Free clears extra memory allocated by MUMPS
func (o *sparseSolverMumps) Free() {
	if o.initialised {
		o.data.job = -2    // finalisation code
		C.dmumps_c(o.data) // do finalize
	}
}

// Fact performs the factorisation
func (o *sparseSolverMumps) Fact() {

	// check
	if !o.initialised {
		chk.Panic("linear solver must be initialised first\n")
	}

	// factorisation
	o.data.job = 2     // factorisation code
	C.dmumps_c(o.data) // factorise
	if o.data.info[1-1] < 0 {
		chk.Panic("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.factorised = true
}

// Solve solves sparse linear systems using MUMPS or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
//   bIsDistr -- this flag tells that the right-hand-side vector 'b' is distributed.
//
func (o *sparseSolverMumps) Solve(x, b Vector, bIsDistr bool) {

	// check
	if !o.factorised {
		chk.Panic("factorisation must be performed first\n")
	}

	// set RHS in processor # 0
	if bIsDistr { // b is distributed => must join
		x.Fill(0)
		o.comm.ReduceSum(x, b) // x := join(b)
	} else {
		if o.comm.Rank() == 0 {
			x.Apply(1, b) // x := b   or   copy(x, b)
		}
	}

	// only proc # 0 needs the RHS
	if o.comm.Rank() == 0 {
		o.data.rhs = (*C.double)(unsafe.Pointer(&x[0]))
	}

	// solve
	o.data.job = 3     // solution code
	C.dmumps_c(o.data) // solve
	if o.data.info[1-1] < 0 {
		chk.Panic("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// broadcast from root
	o.comm.BcastFromRoot(x)
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// MumpsC wraps the MUMPS solver (complex version)
type MumpsC struct {

	// internal
	comm *mpi.Communicator
	t    *TripletC
	mi   []int32
	mj   []int32

	// MUMPS data
	data *C.ZMUMPS_STRUC_C

	// derived
	initialised bool
	factorised  bool
}

// Init initialises mumps for sparse linear systems with real numbers
// args may be nil
func (o *MumpsC) Init(t *TripletC, args *SparseConfig) {

	// check
	if o.initialised {
		chk.Panic("solver must be initialised just once\n")
	}
	if t.pos == 0 {
		chk.Panic("triplet must have at least one item for initialisation\n")
	}

	// default arguments
	if args == nil {
		args = NewSparseConfig(nil)
	}

	// set comm
	o.comm = args.communicator

	// allocate data
	if C.NumDataC == C.NumMaxData {
		chk.Panic("number of MUMPS data available reached. can only allocate %d structures\n", C.NumMaxData)
	}
	o.data = &C.AllDataC[C.NumDataC]
	C.NumDataC++

	// initialise data
	o.data.comm_fortran = -987654 // use Fortran communicator by default
	o.data.par = 1                // host also works
	o.data.sym = 0                // 0=unsymmetric, 1=sym positive definite, 2=general symmetric
	if args.symmetric {
		o.data.sym = 2
	}
	if args.symPosDef {
		o.data.sym = 1
	}
	o.data.job = -1 // initialisation code
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
	o.data.nz_loc = C.int(o.t.pos)
	o.data.irn_loc = (*C.int)(unsafe.Pointer(&o.mi[0]))
	o.data.jcn_loc = (*C.int)(unsafe.Pointer(&o.mj[0]))
	o.data.a_loc = (*C.ZMUMPS_COMPLEX)(unsafe.Pointer(&o.t.x[0]))

	// control
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
	o.data.icntl[5-1] = 0  // assembled matrix (needed for distributed matrix)
	o.data.icntl[6-1] = 7  // automatic (default) permuting strategy for diagonal terms
	o.data.icntl[18-1] = 3 // distributed matrix
	o.data.icntl[14-1] = C.int(args.MumpsIncreaseOfWorkingSpacePct)
	o.data.icntl[23-1] = C.int(args.MumpsMaxMemoryPerProcessor)

	// set ordering and scaling
	o.data.icntl[7-1] = C.int(args.mumpsOrdering) // ordering
	o.data.icntl[8-1] = C.int(args.mumpsScaling)  // scaling

	// analysis step
	o.data.job = 1     // analysis code
	C.zmumps_c(o.data) // analyse
	if o.data.info[1-1] < 0 {
		chk.Panic("analysis failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.initialised = true
}

// Free clears extra memory allocated by MUMPS
func (o *MumpsC) Free() {
	if o.initialised {
		o.data.job = -2    // finalisation code
		C.zmumps_c(o.data) // do finalize
	}
}

// Fact performs the factorisation
func (o *MumpsC) Fact() {

	// check
	if !o.initialised {
		chk.Panic("linear solver must be initialised first\n")
	}

	// factorisation
	o.data.job = 2     // factorisation code
	C.zmumps_c(o.data) // factorise
	if o.data.info[1-1] < 0 {
		chk.Panic("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.factorised = true
}

// Solve solves sparse linear systems using MUMPS or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
//   bIsDistr -- this flag tells that the right-hand-side vector 'b' is distributed.
//
func (o *MumpsC) Solve(x, b VectorC, bIsDistr bool) {

	// check
	if !o.factorised {
		chk.Panic("factorisation must be performed first\n")
	}

	// set RHS in processor # 0
	if bIsDistr { // b is distributed => must join
		o.comm.ReduceSumC(x, b) // x := join(b)
	} else {
		if o.comm.Rank() == 0 {
			x.Apply(1, b) // x := b   or   copy(x, b)
		}
	}

	// only proc # 0 needs the RHS
	if o.comm.Rank() == 0 {
		o.data.rhs = (*C.ZMUMPS_COMPLEX)(unsafe.Pointer(&x[0]))
	}

	// solve
	o.data.job = 3     // solution code
	C.zmumps_c(o.data) // solve
	if o.data.info[1-1] < 0 {
		chk.Panic("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// broadcast from root
	o.comm.BcastFromRootC(x)
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
	spSolverDBc["mumps"] = func() SparseSolverC { return new(MumpsC) }
}
