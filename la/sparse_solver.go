// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import "github.com/cpmech/gosl/chk"

// SparseSolver solves sparse linear systems using UMFPACK or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
type SparseSolver interface {
	Init(t *Triplet, symmetric, verbose bool, ordering, scaling string) error
	Free()
	Fact() error
	Solve(x, b Vector, sumBtoRoot bool) error
}

// spSolverMaker defines a function that makes spSolvers
type spSolverMaker func() SparseSolver

// spSolverDB implements a database of SparseSolver makers
var spSolverDB map[string]spSolverMaker = make(map[string]spSolverMaker)

// NewSparseSolver finds a SparseSolver in database or panic
//   kind -- "umfpack" or "mumps"
func NewSparseSolver(kind string) SparseSolver {
	if maker, ok := spSolverDB[kind]; ok {
		return maker()
	}
	chk.Panic("cannot find SparseSolver named %q in database", kind)
	return nil
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// SparseSolverC solves sparse linear systems using UMFPACK or MUMPS (complex version)
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
type SparseSolverC interface {
	Init(t *TripletC, symmetric, verbose bool, ordering, scaling string) error
	Free()
	Fact() error
	Solve(x, b VectorC, sumBtoRoot bool) error
}

// spSolverMakerC defines a function that makes spSolvers (complex version)
type spSolverMakerC func() SparseSolverC

// spSolverDBc implements a database of SparseSolver makers (complex version)
var spSolverDBc map[string]spSolverMakerC = make(map[string]spSolverMakerC)

// NewSparseSolverC finds a SparseSolver in database or panic
func NewSparseSolverC(kind string) SparseSolverC {
	if maker, ok := spSolverDBc[kind]; ok {
		return maker()
	}
	chk.Panic("cannot find SparseSolverC named %q in database", kind)
	return nil
}
