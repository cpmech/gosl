// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/la"
)

// FdmSolver implements solvers based on the finite-differences method (FDM)
type FdmSolver struct {
	Op     FdmOperator     // differential operator
	Eqs    *la.Equations   // equations numbering in linear system
	Grid   *gm.Grid        // grid structure
	Ebcs   *EssentialBcs   // essential boundary conditions
	U      la.Vector       // vector of unknowns
	F      la.Vector       // right hand-side [if reactions=true]
	matAuk *la.CCMatrix    // dense form of Auk matrix
	matAku *la.CCMatrix    // dense form of Aku matrix [if reactions=true]
	matAkk *la.CCMatrix    // dense form of Akk matrix [if reactions=true]
	buCopy la.Vector       // copy of Bu vector if reactions == true
	linsol la.SparseSolver // linear solver
}

// NewFdmSolver returns a new FDM solver
//  operator -- differential operator; e.g. "laplacian"
//  params   -- parameters for operator; e.g. "kx" and "ky"
//  xmin     -- Grid: [ndim] min/initial coordinates of the whole space (box/cube)
//  xmax     -- Grid: [ndim] max/final coordinates of the whole space (box/cube)
//  ndiv     -- Grid: [ndim] number of divisions for xmax-xmin
func NewFdmSolver(operator string, params dbf.Params, xmin []float64, xmax []float64, ndiv []int) (o *FdmSolver, err error) {
	o = new(FdmSolver)
	o.Op, err = NewFdmOperator(operator, params)
	if err != nil {
		return
	}
	o.Grid = new(gm.Grid)
	err = o.Grid.GenUniform(xmin, xmax, ndiv, false)
	return
}

// SetBcs sets boundary conditions
func (o *FdmSolver) SetBcs(ebcs *EssentialBcs) (err error) {

	// collect known equations
	o.Ebcs = ebcs
	knownEqs := o.Ebcs.GetNodesList()

	// init equations structure
	o.Eqs, err = la.NewEquations(o.Grid.Size(), knownEqs)
	if err != nil {
		return
	}
	o.U = la.NewVector(o.Eqs.N)

	// assemble matrices
	o.Op.Assemble(o.Grid, o.Eqs)
	o.matAuk = o.Eqs.Auk.ToMatrix(nil)

	// init linear solver
	o.linsol = la.NewSparseSolver("umfpack")
	o.linsol.Init(o.Eqs.Auu, true, false, "", "", nil)
	err = o.linsol.Fact()
	return
}

// Solve solves problem
//  The solution will be saved in U
func (o *FdmSolver) Solve(reactions bool) (err error) {

	// check
	if o.Eqs == nil {
		return chk.Err("please set boundary conditions first\n")
	}

	// auxiliary
	bu := o.Eqs.Bu
	xu := o.Eqs.Xu
	xk := o.Eqs.Xk

	// set RHS vector
	bu.Fill(0)

	// set known part of RHS reactions vector
	if reactions {
		if o.matAku == nil {
			o.F = la.NewVector(o.Eqs.N)
			o.buCopy = la.NewVector(o.Eqs.Nu)
			o.matAku = o.Eqs.Aku.ToMatrix(nil)
			o.matAkk = o.Eqs.Akk.ToMatrix(nil)
		}
		copy(o.buCopy, bu)
	}

	// set vector of known values
	for _, bc := range o.Ebcs.All {
		for n := range bc.Nodes {
			i := o.Eqs.FtoK[n]
			xk[i] = bc.Value(0)
		}
	}

	// fix RHS vector: bu -= Auk⋅xk
	la.SpMatVecMulAdd(bu, -1.0, o.matAuk, xk)

	// solve system
	err = o.linsol.Solve(xu, bu, false)
	if err != nil {
		return
	}

	// collect results
	o.Eqs.JoinVector(o.U, xu, xk)
	if reactions {
		xu := o.Eqs.Xu
		bk := o.Eqs.Bk
		la.SpMatVecMul(bk, 1.0, o.matAku, xu)
		la.SpMatVecMulAdd(bk, 1.0, o.matAkk, xk)
		o.Eqs.JoinVector(o.F, o.buCopy, bk)
	}
	return
}
