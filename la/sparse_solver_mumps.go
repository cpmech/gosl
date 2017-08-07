// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

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

// Mumps wraps the MUMPS solver
type Mumps struct {

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
func (o *Mumps) Init(t *Triplet, symmetric, verbose bool, ordering, scaling string, comm *mpi.Communicator) (err error) {

	// check
	if o.initialised {
		return chk.Err("solver must be initialised just once\n")
	}
	if t.pos == 0 {
		return chk.Err("triplet must have at least one item for initialisation\n")
	}

	// set comm
	o.comm = comm

	// allocate data
	if C.NumData == C.NumMaxData {
		return chk.Err("number of MUMPS data available reached. can only allocate %d structures\n", C.NumMaxData)
	}
	o.data = &C.AllData[C.NumData]
	C.NumData++

	// initialise data
	o.data.comm_fortran = -987654 // use Fortran communicator by default
	o.data.par = 1                // host also works
	o.data.sym = 0                // 0=unsymmetric, 1=sym(pos-def), 2=symmetric(undef)
	if symmetric {
		o.data.sym = 2
	}
	o.data.job = -1 // initialisation code
	C.dmumps_c(o.data)
	if o.data.info[1-1] < 0 {
		return chk.Err("init failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
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
	if verbose {
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
	o.data.icntl[5-1] = 0     // assembled matrix (needed for distributed matrix)
	o.data.icntl[6-1] = 7     // automatic (default) permuting strategy for diagonal terms
	o.data.icntl[14-1] = 5000 // % increase of working space
	o.data.icntl[18-1] = 3    // distributed matrix
	o.data.icntl[23-1] = 2000 // max 2000Mb per processor // TODO: check this

	// set ordering and scaling
	ord, sca, err := mumOrderingScaling(ordering, scaling)
	if err != nil {
		return
	}
	o.data.icntl[7-1] = C.int(ord) // ordering
	o.data.icntl[8-1] = C.int(sca) // scaling

	// analysis step
	o.data.job = 1     // analysis code
	C.dmumps_c(o.data) // analyse
	if o.data.info[1-1] < 0 {
		return chk.Err("analysis failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.initialised = true
	return
}

// Free clears extra memory allocated by MUMPS
func (o *Mumps) Free() {
	if o.initialised {
		o.data.job = -2    // finalisation code
		C.dmumps_c(o.data) // do finalize
	}
}

// Fact performs the factorisation
func (o *Mumps) Fact() (err error) {

	// check
	if !o.initialised {
		return chk.Err("linear solver must be initialised first\n")
	}

	// factorisation
	o.data.job = 2     // factorisation code
	C.dmumps_c(o.data) // factorise
	if o.data.info[1-1] < 0 {
		return chk.Err("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.factorised = true
	return
}

// Solve solves sparse linear systems using MUMPS or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
//   bIsDistr -- this flag tells that the right-hand-side vector 'b' is distributed.
//
func (o *Mumps) Solve(x, b Vector, bIsDistr bool) (err error) {

	// check
	if !o.factorised {
		return chk.Err("factorisation must be performed first\n")
	}

	// set RHS in processor # 0
	if bIsDistr { // b is distributed => must join
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
		return chk.Err("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// broadcast from root
	o.comm.BcastFromRoot(x)
	return
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
func (o *MumpsC) Init(t *TripletC, symmetric, verbose bool, ordering, scaling string, comm *mpi.Communicator) (err error) {

	// check
	if o.initialised {
		return chk.Err("solver must be initialised just once\n")
	}
	if t.pos == 0 {
		return chk.Err("triplet must have at least one item for initialisation\n")
	}

	// set comm
	o.comm = comm

	// allocate data
	if C.NumDataC == C.NumMaxData {
		return chk.Err("number of MUMPS data available reached. can only allocate %d structures\n", C.NumMaxData)
	}
	o.data = &C.AllDataC[C.NumDataC]
	C.NumDataC++

	// initialise data
	o.data.comm_fortran = -987654 // use Fortran communicator by default
	o.data.par = 1                // host also works
	o.data.sym = 0                // 0=unsymmetric, 1=sym(pos-def), 2=symmetric(undef)
	if symmetric {
		o.data.sym = 2
	}
	o.data.job = -1 // initialisation code
	C.zmumps_c(o.data)
	if o.data.info[1-1] < 0 {
		return chk.Err("init failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
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
	if verbose {
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
	o.data.icntl[5-1] = 0     // assembled matrix (needed for distributed matrix)
	o.data.icntl[6-1] = 7     // automatic (default) permuting strategy for diagonal terms
	o.data.icntl[14-1] = 5000 // % increase of working space
	o.data.icntl[18-1] = 3    // distributed matrix
	o.data.icntl[23-1] = 2000 // max 2000Mb per processor // TODO: check this

	// set ordering and scaling
	ord, sca, err := mumOrderingScaling(ordering, scaling)
	if err != nil {
		return
	}
	o.data.icntl[7-1] = C.int(ord) // ordering
	o.data.icntl[8-1] = C.int(sca) // scaling

	// analysis step
	o.data.job = 1     // analysis code
	C.zmumps_c(o.data) // analyse
	if o.data.info[1-1] < 0 {
		return chk.Err("analysis failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.initialised = true
	return
}

// Free clears extra memory allocated by MUMPS
func (o *MumpsC) Free() {
	if o.initialised {
		o.data.job = -2    // finalisation code
		C.zmumps_c(o.data) // do finalize
	}
}

// Fact performs the factorisation
func (o *MumpsC) Fact() (err error) {

	// check
	if !o.initialised {
		return chk.Err("linear solver must be initialised first\n")
	}

	// factorisation
	o.data.job = 2     // factorisation code
	C.zmumps_c(o.data) // factorise
	if o.data.info[1-1] < 0 {
		return chk.Err("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// success
	o.factorised = true
	return
}

// Solve solves sparse linear systems using MUMPS or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
//   bIsDistr -- this flag tells that the right-hand-side vector 'b' is distributed.
//
func (o *MumpsC) Solve(x, b VectorC, bIsDistr bool) (err error) {

	// check
	if !o.factorised {
		return chk.Err("factorisation must be performed first\n")
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
		return chk.Err("solver failed: %v\n", mumErr(o.data.info[1-1], o.data.info[2-1]))
	}

	// broadcast from root
	o.comm.BcastFromRootC(x)
	return
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// mumOrderingScaling sets the ordering and scaling methods for MUMPS
func mumOrderingScaling(ordering, scaling string) (ord, sca int, err error) {

	// ordering
	if ordering == "" {
		ordering = "amf"
	}
	switch ordering {
	case "amd":
		ord = 0
	case "amf":
		ord = 2
	case "scotch":
		ord = 3
	case "pord":
		ord = 4
	case "metis":
		ord = 5
	case "qamd":
		ord = 6
	case "auto":
		ord = 7
	default:
		err = chk.Err("ordering scheme %s is not available\n", ordering)
		return
	}

	// scaling
	if scaling == "" {
		scaling = "rcit"
	}
	switch scaling {
	case "no":
		sca = 0 // no scaling
	case "diag":
		sca = 1 // diagonal scaling
	case "rcit":
		sca = 7 // row/col iterative
	case "rrcit":
		sca = 8 // rigorous row/col it
	case "auto":
		sca = 77 // automatic
	default:
		err = chk.Err("scaling scheme %s is not available\n", scaling)
		return
	}
	return
}

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
	default:
		return io.Sf("MUMPS Error # %d: unknown error", info)
	}
	return ""
}

// add solvers to database /////////////////////////////////////////////////////////////////////////

func init() {
	spSolverDB["mumps"] = func() SparseSolver { return new(Mumps) }
	spSolverDBc["mumps"] = func() SparseSolverC { return new(MumpsC) }
}
