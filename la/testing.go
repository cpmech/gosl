// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math/cmplx"
	"testing"

	"github.com/cpmech/gosl/chk"
)

// TestSolverResidual check the residual of a linear system solution
func TestSolverResidual(tst *testing.T, a *Matrix, x, b Vector, tolNorm float64) {
	r := NewVector(len(x))
	r.Apply(-1, b)           // r := -b
	MatVecMulAdd(r, 1, a, x) // r += 1*a*x
	resid := r.Norm()
	if resid > tolNorm {
		tst.Errorf("residual is too large: %g\n", resid)
		return
	}
}

// TestSolverResidualC check the residual of a linear system solution (complex version)
func TestSolverResidualC(tst *testing.T, a *MatrixC, x, b VectorC, tolNorm float64) {
	r := NewVectorC(len(x))
	r.Apply(-1, b)            // r = -b
	MatVecMulAddC(r, 1, a, x) // r += 1*a*x
	resid := cmplx.Abs(r.Norm())
	if resid > tolNorm {
		tst.Errorf("residual is too large: %g\n", resid)
		return
	}
}

// TestSpSolver tests a sparse solver
func TestSpSolver(tst *testing.T, solverKind string, symmetric bool, t *Triplet, b, xCorrect Vector,
	tolX, tolRes float64, verbose bool) {

	// allocate solver
	o := NewSparseSolver(solverKind)
	defer o.Free()

	// initialize solver
	args := NewSparseConfig()
	if symmetric {
		if solverKind == "mumps" {
			args.SetMumpsSymmetry(true, false)
		} else {
			args.SetUmfpackSymmetry()
		}
	}
	args.Verbose = verbose
	o.Init(t, args)

	// factorise
	o.Fact()

	// solve
	x := NewVector(len(b))
	o.Solve(x, b) // x := inv(A) * b

	// check
	chk.Array(tst, "x", tolX, x, xCorrect)
	TestSolverResidual(tst, t.ToDense(), x, b, tolRes)
}

// TestSpSolverC tests a sparse solver (complex version)
func TestSpSolverC(tst *testing.T, solverKind string, symmetric bool, t *TripletC, b, xCorrect VectorC,
	tolX, tolRes float64, verbose bool) {

	// allocate solver
	o := NewSparseSolverC(solverKind)
	defer o.Free()

	// initialize solver
	args := NewSparseConfig()
	if symmetric {
		if solverKind == "mumps" {
			args.SetMumpsSymmetry(true, false)
		} else {
			args.SetUmfpackSymmetry()
		}
	}
	args.Verbose = verbose
	o.Init(t, args)

	// factorise
	o.Fact()

	// solve
	x := NewVectorC(len(b))
	o.Solve(x, b) // x := inv(A) * b

	// check
	chk.ArrayC(tst, "x", tolX, x, xCorrect)
	TestSolverResidualC(tst, t.ToDense(), x, b, tolRes)
}
