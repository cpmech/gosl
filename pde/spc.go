// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/la"
)

// SpcSolver implements solvers based on spectral collocation method
type SpcSolver struct {
	Op     SpcOperator           // differential operator
	Interp []*fun.LagrangeInterp // Lagrange interpolators
	Eqs    *la.Equations         // equations numbering in linear system
	Grid   *gm.Grid              // grid structure
	Ebcs   *EssentialBcs         // essential boundary conditions
	U      la.Vector             // vector of unknowns
	F      la.Vector             // right hand-side [if reactions=true]
	matAuk *la.CCMatrix          // c-c-mat form of Auk matrix
	matAku *la.CCMatrix          // c-c-mat form of Aku matrix [if reactions=true]
	matAkk *la.CCMatrix          // c-c-mat form of Akk matrix [if reactions=true]
	buCopy la.Vector             // copy of Bu vector if reactions == true
	linsol la.SparseSolver       // linear solver
}

// NewSpcSolver returns a new FDM solver
//  operator -- differential operator; e.g. "laplacian"
//  gtype    -- grid type: "uni", "cgl" (Chebyshev-Gauss-Lobato)
//  params   -- parameters for operator; e.g. "kx" and "ky"
//  xmin     -- Grid: [ndim] min/initial coordinates of the whole space (box/cube)
//  xmax     -- Grid: [ndim] max/final coordinates of the whole space (box/cube)
//  ndiv     -- Grid: [ndim] number of divisions for xmax-xmin
func NewSpcSolver(operator, gtype string, params dbf.Params, xmin []float64, xmax []float64, ndiv []int) (o *SpcSolver, err error) {

	// check lengths
	ndim := len(xmin)
	if len(xmax) != ndim {
		return nil, chk.Err("len(xmax) must be equal to len(xmin) == ndim. %d != %d\n", len(xmax), ndim)
	}
	if len(ndiv) != ndim {
		return nil, chk.Err("len(ndiv) must be equal to len(xmin) == ndim. %d != %d\n", len(ndiv), ndim)
	}

	// new solver and operator
	o = new(SpcSolver)
	o.Op, err = NewSpcOperator(operator, params)
	if err != nil {
		return
	}

	// Lagrange interpolators
	o.Interp = make([]*fun.LagrangeInterp, ndim)
	for i := 0; i < ndim; i++ {

		// allocate
		o.Interp[i], err = fun.NewLagrangeInterp(ndiv[i], gtype)
		if err != nil {
			return
		}

		// compute D1 matrix
		err = o.Interp[i].CalcD1()
		if err != nil {
			return
		}

		// compute D2 matrix
		err = o.Interp[i].CalcD2()
		if err != nil {
			return
		}
	}

	// grid
	o.Grid = new(gm.Grid)
	if ndim == 2 {
		err = o.Grid.Set2d(o.Interp[0].X, o.Interp[1].X, false)
	} else {
		err = o.Grid.Set3d(o.Interp[0].X, o.Interp[1].X, o.Interp[2].X, false)
	}
	if err != nil {
		return
	}

	// TODO:
	//  map [-1,+1] to xmin and xmax
	return
}

// SetBcs sets boundary conditions
func (o *SpcSolver) SetBcs(ebcs *EssentialBcs) (err error) {

	// collect known equations
	o.Ebcs = ebcs
	knownEqs := o.Ebcs.GetNodesList()

	// total number of equations
	neq := 1
	for _, lag := range o.Interp {
		neq *= (lag.N + 1)
	}

	// init equations structure
	o.Eqs, err = la.NewEquations(neq, knownEqs)
	if err != nil {
		return
	}
	o.U = la.NewVector(o.Eqs.N)

	// assemble matrices
	o.Op.Assemble(o.Interp, o.Eqs)
	o.matAuk = o.Eqs.Auk.ToMatrix(nil)

	// init linear solver
	o.linsol = la.NewSparseSolver("umfpack")
	o.linsol.Init(o.Eqs.Auu, true, false, "", "", nil)
	err = o.linsol.Fact()
	return
}

// Solve solves problem
//  The solution will be saved in U
func (o *SpcSolver) Solve(reactions bool) (err error) {

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
