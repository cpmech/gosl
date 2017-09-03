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
	Operator  FdmOperator     // differential operator
	Equations *la.Equations   // equations numbering in linear system
	Grid      *gm.Grid        // grid structure
	Ebcs      *EssentialBcs   // essential boundary conditions
	U         la.Vector       // vector of unknowns
	F         la.Vector       // right hand-side [if reactions=true]
	denseAuk  *la.CCMatrix    // dense form of Auk matrix
	denseAku  *la.CCMatrix    // dense form of Aku matrix [if reactions=true]
	denseAkk  *la.CCMatrix    // dense form of Akk matrix [if reactions=true]
	buCopy    la.Vector       // copy of Bu vector if reactions == true
	linsol    la.SparseSolver // linear solver
}

// NewFdmSolver returns a new FDM solver
//  operator -- differential operator; e.g. "laplacian"
//  params   -- parameters for operator; e.g. "kx" and "ky"
//  xmin     -- Grid: [ndim] min/initial coordinates of the whole space (box/cube)
//  xmax     -- Grid: [ndim] max/final coordinates of the whole space (box/cube)
//  ndiv     -- Grid: [ndim] number of divisions for xmax-xmin
func NewFdmSolver(operator string, params dbf.Params, xmin []float64, xmax []float64, ndiv []int) (o *FdmSolver, err error) {
	o = new(FdmSolver)
	o.Operator, err = NewFdmOperator(operator, params)
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
	var knownEqs []int
	for _, bc := range o.Ebcs.All {
		nodes := bc.GetNodesSorted()
		knownEqs = append(knownEqs, nodes...)
	}

	// init equations structure
	o.Equations, err = la.NewEquations(o.Grid.Size(), knownEqs)
	if err != nil {
		return
	}
	o.U = la.NewVector(o.Equations.N)

	// assemble matrices
	o.Operator.Assemble(o.Grid, o.Equations)
	o.denseAuk = o.Equations.Auk.ToMatrix(nil)

	// init linear solver
	o.linsol = la.NewSparseSolver("umfpack")
	o.linsol.Init(o.Equations.Auu, true, false, "", "", nil)
	err = o.linsol.Fact()
	return
}

// Solve solves problem
func (o *FdmSolver) Solve(reactions bool) (err error) {

	// check
	if o.Equations == nil {
		return chk.Err("please set boundary conditions first\n")
	}

	// auxiliary
	bu := o.Equations.Bu
	xu := o.Equations.Xu
	xk := o.Equations.Xk

	// set RHS vector
	bu.Fill(0)

	// set known part of RHS reactions vector
	if reactions {
		if o.denseAku == nil {
			o.F = la.NewVector(o.Equations.N)
			o.buCopy = la.NewVector(o.Equations.Nu)
			o.denseAku = o.Equations.Aku.ToMatrix(nil)
			o.denseAkk = o.Equations.Akk.ToMatrix(nil)
		}
		copy(o.buCopy, bu)
	}

	// set vector of known values
	for _, bc := range o.Ebcs.All {
		for n := range bc.Nodes {
			i := o.Equations.FtoK[n]
			xk[i] = bc.Value(0)
		}
	}

	// fix RHS vector: bu -= Auk⋅xk
	la.SpMatVecMulAdd(bu, -1.0, o.denseAuk, xk)

	// solve system
	err = o.linsol.Solve(xu, bu, false)
	if err != nil {
		return
	}

	// collect results
	o.Equations.JoinVector(o.U, xu, xk)
	if reactions {
		xu := o.Equations.Xu
		bk := o.Equations.Bk
		la.SpMatVecMul(bk, 1.0, o.denseAku, xu)
		la.SpMatVecMulAdd(bk, 1.0, o.denseAkk, xk)
		o.Equations.JoinVector(o.F, o.buCopy, bk)
	}
	return
}
