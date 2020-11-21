// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"gosl/chk"
	"gosl/mpi"
)

// The SparseConfig structure holds configuration arguments for sparse solvers
type SparseConfig struct {
	// external
	Symmetric bool   // indicates symmetric system. NOTE: when using MUMPS, only the upper or lower part of the matrix must be provided
	SymPosDef bool   // indicates symmetric-positive-defined system
	Verbose   bool   // run on Verbose mode
	Guess     Vector // initial guess for iterative solvers [may be nil]

	// MUMPS control parameters (check MUMPS solver manual)
	MumpsIncreaseOfWorkingSpacePct int // ICNTL(14) default = 100%
	MumpsMaxMemoryPerProcessor     int // ICNTL(23) default = 2000Mb
	mumpsOrdering                  int // ICNTL(7) default = "" == "auto"
	mumpsScaling                   int // Scaling type (check MUMPS solver) [may be empty]

	// internal
	communicator *mpi.Communicator // MPI communicator for parallel solvers [may be nil]
}

// NewSparseConfig returns a new SparseConfig
// comm may be nil
func NewSparseConfig(comm *mpi.Communicator) (o *SparseConfig) {
	o = new(SparseConfig)
	o.communicator = comm
	o.MumpsIncreaseOfWorkingSpacePct = 100
	o.MumpsMaxMemoryPerProcessor = 2000
	o.SetMumpsOrdering("")
	o.SetMumpsScaling("")
	return
}

// SetMumpsOrdering sets ordering for MUMPS solver
// ordering -- "" or "amf" [default]
//             "amf", "scotch", "pord", "metis", "qamd", "auto"
// ICNTL(7)
//   0: "amd" Approximate Minimum Degree (AMD)
//   2: "amf" Approximate Minimum Fill (AMF)
//   3: "scotch" SCOTCH5 package is used if previously installed by the user otherwise treated as 7.
//   4: "pord" PORD6 is used if previously installed by the user otherwise treated as 7.
//   5: "metis" Metis7 package is used if previously installed by the user otherwise treated as 7.
//   6: "qamd" Approximate Minimum Degree with automatic quasi-dense row detection (QAMD) is used.
//   7: "auto" automatic choice by the software during analysis phase. This choice will depend on the
//       ordering packages made available, on the matrix (type and size), and on the number of processors.
func (o *SparseConfig) SetMumpsOrdering(ordering string) {
	switch ordering {
	case "amd":
		o.mumpsOrdering = 0
	case "", "amf":
		o.mumpsOrdering = 2
	case "scotch":
		o.mumpsOrdering = 3
	case "pord":
		o.mumpsOrdering = 4
	case "metis":
		o.mumpsOrdering = 5
	case "qamd":
		o.mumpsOrdering = 6
	case "auto":
		o.mumpsOrdering = 7
	default:
		chk.Panic("ordering scheme %s is not available\n", ordering)
	}
}

// SetMumpsScaling sets scaling for MUMPS solver
// scaling -- "" or "rcit" [default]
//            "no", "diag", "col", "rcinf", "rcit", "rrcit", "auto"
// ICNTL(8)
//   0: "no" No scaling applied/computed.
//   1: "diag" Diagonal scaling computed during the numerical factorization phase,
//   3: "col" Column scaling computed during the numerical factorization phase,
//   4: "rcinf" Row and column scaling based on infinite row/column norms, computed during the numerical
//      factorization phase,
//   7: "rcit" Simultaneous row and column iterative scaling based on [41] and [15] computed during the
//      numerical factorization phase.
//   8: "rrcit" Similar to 7 but more rigorous and expensive to compute; computed during the numerical
//      factorization phase.
//   77: "auto" Automatic choice of the value of ICNTL(8) done during analy
func (o *SparseConfig) SetMumpsScaling(scaling string) {
	switch scaling {
	case "no":
		o.mumpsScaling = 0 // no scaling
	case "diag":
		o.mumpsScaling = 1 // diagonal scaling
	case "col":
		o.mumpsScaling = 3 // column
	case "rcinf":
		o.mumpsScaling = 4 // row and column based on inf norms
	case "", "rcit":
		o.mumpsScaling = 7 // row/col iterative
	case "rrcit":
		o.mumpsScaling = 8 // rigorous row/col it
	case "auto":
		o.mumpsScaling = 77 // automatic
	default:
		chk.Panic("scaling scheme %s is not available\n", scaling)
	}
}
