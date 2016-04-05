// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!appengine,!heroku

package la

/*
#cgo CFLAGS: -O3
#cgo LDFLAGS: -ldmumps -lzmumps -lmumps_common -lpord
#include <dmumps_c.h>
#include <zmumps_c.h>
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/mpi"
)

// LinSolMumps holds MUMPS data
type LinSolMumps struct {
	linSolData
	m      C.DMUMPS_STRUC_C
	mz     C.ZMUMPS_STRUC_C
	mi, mj []int32
	xRC    []float64

	// derived
	is_initialised bool
}

// factory of allocators
func init() {
	lsAllocators["mumps"] = func() LinSol { return new(LinSolMumps) }
}

// InitR initialises a LinSolMumps data structure for Real systems. It also performs some initial analyses.
func (o *LinSolMumps) InitR(tR *Triplet, symmetric, verbose, timing bool) (err error) {

	// check
	o.tR = tR
	if tR.pos == 0 {
		return chk.Err(_linsol_mumps_err01)
	}

	// flags
	o.name = "mumps"
	o.sym = symmetric
	o.cmplx = false
	o.verb = verbose
	o.ton = timing

	// start time
	if mpi.Rank() != 0 {
		o.verb = false
		o.ton = false
	}
	if o.ton {
		o.tini = time.Now()
	}

	// initialise Mumps
	o.m.comm_fortran = -987654 // use Fortran communicator by default
	o.m.par = 1                // host also works
	o.m.sym = 0                // 0=unsymmetric, 1=sym(pos-def), 2=symmetric(undef)
	if symmetric {
		o.m.sym = 2
	}
	o.m.job = -1     // initialisation code
	C.dmumps_c(&o.m) // initialise
	if o.m.info[1-1] < 0 {
		return chk.Err(_linsol_mumps_err02, mumps_error(o.m.info[1-1], o.m.info[2-1]))
	}

	// convert indices to C.int (not C.long) and
	// increment indices since Mumps is 1-based (FORTRAN)
	o.mi, o.mj = make([]int32, o.tR.pos), make([]int32, o.tR.pos)
	for k := 0; k < tR.pos; k++ {
		o.mi[k] = int32(o.tR.i[k]) + 1
		o.mj[k] = int32(o.tR.j[k]) + 1
	}

	// set pointers
	o.m.n = C.int(o.tR.m)
	o.m.nz_loc = C.int(o.tR.pos)
	o.m.irn_loc = (*C.int)(unsafe.Pointer(&o.mi[0]))
	o.m.jcn_loc = (*C.int)(unsafe.Pointer(&o.mj[0]))
	o.m.a_loc = (*C.double)(unsafe.Pointer(&o.tR.x[0]))

	// control
	if verbose {
		if mpi.Rank() == 0 {
			io.Pfgreen("\n . . . . . . . . . . . . . . LinSolMumps.InitR . . (MUMPS) . . . . . . . . . \n\n")
		}
		o.m.icntl[1-1] = 6 // output stream for error messages
		o.m.icntl[2-1] = 0 // output stream for statistics and warnings
		o.m.icntl[3-1] = 6 // output stream for global information
		o.m.icntl[4-1] = 2 // message level: 2==errors and warnings
	} else {
		o.m.icntl[1-1] = -1 // no output messages
		o.m.icntl[2-1] = -1 // no warnings
		o.m.icntl[3-1] = -1 // no global information
		o.m.icntl[4-1] = -1 // message level
	}
	o.m.icntl[5-1] = 0     // assembled matrix (needed for distributed matrix)
	o.m.icntl[6-1] = 7     // automatic (default) permuting strategy for diagonal terms
	o.m.icntl[14-1] = 5000 // % increase of working space
	o.m.icntl[18-1] = 3    // distributed matrix
	o.m.icntl[23-1] = 2000 // max 2000Mb per processor // TODO: check this
	o.SetOrdScal("", "")

	// analysis step
	o.m.job = 1      // analysis code
	C.dmumps_c(&o.m) // analyse
	if o.m.info[1-1] < 0 {
		return chk.Err(_linsol_mumps_err03, mumps_error(o.m.info[1-1], o.m.info[2-1]))
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolMumps.InitR = %v\n", o.name, time.Now().Sub(o.tini))
	}

	// success
	o.is_initialised = true
	return
}

// InitC initialises a LinSolMumps data structure for Complex systems. It also performs some initial analyses.
func (o *LinSolMumps) InitC(tC *TripletC, symmetric, verbose, timing bool) (err error) {

	// check
	o.tC = tC
	if tC.pos == 0 {
		return chk.Err(_linsol_mumps_err04)
	}

	// flags
	o.name = "mumps"
	o.sym = symmetric
	o.cmplx = true
	o.verb = verbose
	o.ton = timing

	// start time
	if mpi.Rank() != 0 {
		o.verb = false
		o.ton = false
	}
	if o.ton {
		o.tini = time.Now()
	}

	// check xz
	if len(o.tC.xz) != 2*len(o.tC.i) {
		return chk.Err(_linsol_mumps_err05, len(o.tC.xz), len(o.tC.i))
	}

	// initialise Mumps
	o.mz.comm_fortran = -987654 // use Fortran communicator by default
	o.mz.par = 1                // host also works
	o.mz.sym = 0                // 0=unsymmetric, 1=sym(pos-def), 2=symmetric(undef)
	if symmetric {
		o.mz.sym = 2
	}
	o.mz.job = -1     // initialisation code
	C.zmumps_c(&o.mz) // initialise
	if o.mz.info[1-1] < 0 {
		return chk.Err(_linsol_mumps_err06, mumps_error(o.mz.info[1-1], o.mz.info[2-1]))
	}

	// convert indices to C.int (not C.long) and
	// increment indices since Mumps is 1-based (FORTRAN)
	o.mi, o.mj = make([]int32, o.tC.pos), make([]int32, o.tC.pos)
	for k := 0; k < tC.pos; k++ {
		o.mi[k] = int32(o.tC.i[k]) + 1
		o.mj[k] = int32(o.tC.j[k]) + 1
	}

	// set pointers
	o.mz.n = C.int(o.tC.m)
	o.mz.nz_loc = C.int(o.tC.pos)
	o.mz.irn_loc = (*C.int)(unsafe.Pointer(&o.mi[0]))
	o.mz.jcn_loc = (*C.int)(unsafe.Pointer(&o.mj[0]))
	o.mz.a_loc = (*C.ZMUMPS_COMPLEX)(unsafe.Pointer(&o.tC.xz[0]))

	// only proc # 0 needs the RHS
	if mpi.Rank() == 0 {
		o.xRC = make([]float64, 2*o.tC.n)
		o.mz.rhs = (*C.ZMUMPS_COMPLEX)(unsafe.Pointer(&o.xRC[0]))
	}

	// control
	if verbose {
		if mpi.Rank() == 0 {
			io.Pfgreen("\n . . . . . . . . . . . . . . LinSolMumps.InitC . . (MUMPS) . . . . . . . . . \n\n")
		}
		o.mz.icntl[1-1] = 6 // output stream for error messages
		o.mz.icntl[2-1] = 0 // output stream for statistics and warnings
		o.mz.icntl[3-1] = 6 // output stream for global information
		o.mz.icntl[4-1] = 2 // message level: 2==errors and warnings
	} else {
		o.mz.icntl[1-1] = -1 // no output messages
		o.mz.icntl[2-1] = -1 // no warnings
		o.mz.icntl[3-1] = -1 // no global information
		o.mz.icntl[4-1] = -1 // message level
	}
	o.mz.icntl[5-1] = 0     // assembled matrix (needed for distributed matrix)
	o.mz.icntl[6-1] = 7     // automatic (default) permuting strategy for diagonal terms
	o.mz.icntl[14-1] = 5000 // % increase of working space
	o.mz.icntl[18-1] = 3    // distributed matrix
	o.SetOrdScal("", "")

	// analysis step
	o.mz.job = 1      // analysis code
	C.zmumps_c(&o.mz) // analyse
	if o.mz.info[1-1] < 0 {
		return chk.Err(_linsol_mumps_err07, mumps_error(o.mz.info[1-1], o.mz.info[2-1]))
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolMumps.InitC = %v\n", o.name, time.Now().Sub(o.tini))
	}

	// success
	o.is_initialised = true
	return
}

// Fact performs symbolic/numeric factorisation. This method also converts the triplet form
// to the column-compressed form, including the summation of duplicated entries
func (o *LinSolMumps) Fact() (err error) {

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
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolMumps.Fact . . . . . . . . . . . . . . . \n\n")
	}

	// complex
	if o.cmplx {

		// MUMPS: factorisation
		o.mz.job = 2      // factorisation code
		C.zmumps_c(&o.mz) // factorise
		if o.mz.info[1-1] < 0 {
			return chk.Err(_linsol_mumps_err08, "Real", mumps_error(o.mz.info[1-1], o.mz.info[2-1]))
		}

		// real
	} else {

		// MUMPS: factorisation
		o.m.job = 2      // factorisation code
		C.dmumps_c(&o.m) // factorise
		if o.m.info[1-1] < 0 {
			return chk.Err(_linsol_mumps_err08, "Complex", mumps_error(o.m.info[1-1], o.m.info[2-1]))
		}
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolMumps.Fact  = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// SolveR solves the linear Real system A.x = b
//  NOTES:
//    1) sum_b_to_root is a flag for MUMPS; it tells Solve to sum the values in 'b' arrays to the root processor
func (o *LinSolMumps) SolveR(xR, bR []float64, sum_b_to_root bool) (err error) {

	// check
	if !o.is_initialised {
		return chk.Err("linear solver must be initialised first\n")
	}
	if o.cmplx {
		return chk.Err(_linsol_mumps_err09)
	}

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// message
	if o.verb {
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolMumps.SolveR . . . . . . . . . . . . . . . \n\n")
	}

	// MUMPS: set RHS in processor # 0
	if sum_b_to_root {
		mpi.SumToRoot(xR, bR)
	} else {
		if mpi.Rank() == 0 {
			copy(xR, bR) // x := b
		}
	}

	// only proc # 0 needs the RHS
	if mpi.Rank() == 0 {
		o.m.rhs = (*C.double)(unsafe.Pointer(&xR[0]))
	}

	// MUMPS: solve
	o.m.job = 3      // solution code
	C.dmumps_c(&o.m) // solve
	if o.m.info[1-1] < 0 {
		return chk.Err(_linsol_mumps_err10, mumps_error(o.m.info[1-1], o.m.info[2-1]))
	}
	mpi.BcastFromRoot(xR) // broadcast from root

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolMumps.Solve = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// SolveC solves the linear Complex system A.x = b
//  NOTES:
//    1) sum_b_to_root is a flag for MUMPS; it tells Solve to sum the values in 'b' arrays to the root processor
func (o *LinSolMumps) SolveC(xR, xC, bR, bC []float64, sum_b_to_root bool) (err error) {

	// check
	if !o.is_initialised {
		return chk.Err("linear solver must be initialised first\n")
	}
	if !o.cmplx {
		return chk.Err(_linsol_mumps_err11)
	}

	// start time
	if o.ton {
		o.tini = time.Now()
	}

	// message
	if o.verb {
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolMumps.SolveC . . . . . . . . . . . . . . . \n\n")
	}

	// MUMPS: set RHS in processor # 0
	if sum_b_to_root {
		mpi.SumToRoot(xR, bR)
		mpi.SumToRoot(xC, bC)
		// join complex values
		if mpi.Rank() == 0 {
			for i := 0; i < len(xR); i++ {
				o.xRC[i*2], o.xRC[i*2+1] = xR[i], xC[i]
			}
		}
	} else {
		// join complex values
		if mpi.Rank() == 0 {
			for i := 0; i < len(xR); i++ {
				o.xRC[i*2], o.xRC[i*2+1] = bR[i], bC[i]
			}
		}
	}

	// MUMPS: solve
	o.mz.job = 3      // solution code
	C.zmumps_c(&o.mz) // solve
	if o.mz.info[1-1] < 0 {
		return chk.Err(_linsol_mumps_err12, mumps_error(o.mz.info[1-1], o.mz.info[2-1]))
	}

	// MUMPS: split complex values
	if mpi.Rank() == 0 {
		for i := 0; i < len(xR); i++ {
			xR[i], xC[i] = o.xRC[i*2], o.xRC[i*2+1]
		}
	}

	// MUMPS: broadcast from root
	mpi.BcastFromRoot(xR)
	mpi.BcastFromRoot(xC)

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolMumps.Solve = %v\n", o.name, time.Now().Sub(o.tini))
	}
	return
}

// Clean deletes temporary data structures
func (o *LinSolMumps) Clean() {

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
		io.Pfgreen("\n . . . . . . . . . . . . . . LinSolMumps.Clean . . . . . . . . . . . . . . . \n\n")
	}

	// clean up
	if o.cmplx {
		o.mz.job = -2     // finalize code
		C.zmumps_c(&o.mz) // do finalize
	} else {
		o.m.job = -2     // finalize code
		C.dmumps_c(&o.m) // do finalize
	}

	// duration
	if o.ton {
		io.Pfcyan("%s: Time spent in LinSolMumps.Clean   = %v\n", o.name, time.Now().Sub(o.tini))
	}
}

// SetOrdScal sets the ordering and scaling methods for MUMPS
func (o *LinSolMumps) SetOrdScal(ordering, scaling string) (err error) {
	var ord, sca int
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
		return chk.Err(_linsol_mumps_err13, ordering)
	}
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
		return chk.Err(_linsol_mumps_err14, scaling)
	}
	if o.cmplx {
		o.mz.icntl[7-1] = C.int(ord) // ordering
		o.mz.icntl[8-1] = C.int(sca) // scaling
	} else {
		o.m.icntl[7-1] = C.int(ord) // ordering
		o.m.icntl[8-1] = C.int(sca) // scaling
	}
	return
}

func mumps_error(info, infx C.int) string {
	switch info {
	case -3:
		return "linsol.go: Solve(MUMPS): Error # -3: MUMPS was called with an invalid value for JOB"
	case -6:
		return "linsol.go: Solve(MUMPS): Error # -6: singular matrix"
	case -9:
		return io.Sf("linsol.go: Solve(MUMPS): Error # -9: main internal real/complex workarray S too small. info(2)=%v", infx)
	case -10:
		return "linsol.go: Solve(MUMPS): Error # -10: singular matrix"
	case -13:
		return "linsol.go: Solve(MUMPS): Error # -13: out of memory"
	default:
		return io.Sf("linsol.go: Solve(MUMPS): Error # %d: unknown error", info)
	}
	return ""
}

func RunMumpsTestR(t *Triplet, tol_cmp float64, b, x_correct []float64, sum_b_to_root bool) {

	// info
	symmetric := false
	verbose := false
	timing := false

	// allocate solver
	lis := GetSolver("mumps")
	defer lis.Clean()

	// initialise solver
	err := lis.InitR(t, symmetric, verbose, timing)
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// factorise
	err = lis.Fact()
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// solve
	x := make([]float64, len(b))
	err = lis.SolveR(x, b, sum_b_to_root) // x := inv(A) * b
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	if mpi.Rank() == 0 {
		// output
		A := t.ToMatrix(nil)
		io.Pforan("A.x = b\n")
		PrintMat("A", A.ToDense(), "%5g", false)
		PrintVec("x", x, "%g ", false)
		PrintVec("b", b, "%g ", false)

		// check
		err := VecMaxDiff(x, x_correct)
		if err > tol_cmp {
			chk.Panic("test failed: err = %g", err)
		}
		io.Pf("err(x) = %g [1;32mOK[0m\n", err)
	}
}

func RunMumpsTestC(t *TripletC, tol_cmp float64, b, x_correct []complex128, sum_b_to_root bool) {

	// info
	symmetric := false
	verbose := false
	timing := false

	// allocate solver
	lis := GetSolver("mumps")
	defer lis.Clean()

	// initialise solver
	err := lis.InitC(t, symmetric, verbose, timing)
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// factorise
	err = lis.Fact()
	if err != nil {
		chk.Panic("%v", err.Error())
	}

	// solve
	bR, bC := ComplexToRC(b)
	xR := make([]float64, len(b))
	xC := make([]float64, len(b))
	err = lis.SolveC(xR, xC, bR, bC, sum_b_to_root) // x := inv(A) * b
	if err != nil {
		chk.Panic("%v", err.Error())
	}
	x := RCtoComplex(xR, xC)

	if mpi.Rank() == 0 {
		// output
		A := t.ToMatrix(nil)
		io.Pforan("A.x = b\n")
		PrintMatC("A", A.ToDense(), "(%g+", "%gi) ", false)
		PrintVecC("x", x, "(%g+", "%gi) ", false)
		PrintVecC("b", b, "(%g+", "%gi) ", false)

		// check
		xR_correct, xC_correct := ComplexToRC(x_correct)
		errR := VecMaxDiff(xR, xR_correct)
		if errR > tol_cmp {
			chk.Panic("test failed: errR = %g", errR)
		}
		errC := VecMaxDiff(xC, xC_correct)
		if errC > tol_cmp {
			chk.Panic("test failed: errC = %g", errC)
		}
		io.Pf("err(xR) = %g [1;32mOK[0m\n", errR)
		io.Pf("err(xC) = %g [1;32mOK[0m\n", errC)
	}
}

// error messages
var (
	_linsol_mumps_err01 = "triplet must have at least one item before calling this method\n"
	_linsol_mumps_err02 = "init failed: %v\n"
	_linsol_mumps_err03 = "analysis failed: %v\n"
	_linsol_mumps_err04 = "triplet must have at least one item before calling this method\n"
	_linsol_mumps_err05 = "length of xz (%d) slice must be equal to two-times the length of i (%d) slice when using MUMPS complex solver\n"
	_linsol_mumps_err06 = "init failed %v\n"
	_linsol_mumps_err07 = "analysis failed: %v\n"
	_linsol_mumps_err08 = "(%s) failed %v\n"
	_linsol_mumps_err09 = "this method must be called with Real matrices\n"
	_linsol_mumps_err10 = "solver failed: %v\n"
	_linsol_mumps_err11 = "this method must be called with Complex matrices\n"
	_linsol_mumps_err12 = "solver failed: %v\n"
	_linsol_mumps_err13 = "ordering scheme %s is not available\n"
	_linsol_mumps_err14 = "scaling scheme %s is not available\n"
)
