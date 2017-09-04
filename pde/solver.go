// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// Solver solvers a PDE by calling specific operators
type Solver struct {
	Op     Operator        // differential operator
	Eqs    *la.Equations   // equations numbering in linear system
	Grid   *gm.Grid        // grid structure
	Ebcs   *EssentialBcs   // essential boundary conditions
	U      la.Vector       // vector of unknowns
	F      la.Vector       // right hand-side [if reactions=true]
	matAuk *la.CCMatrix    // c-c-mat form of Auk matrix
	matAku *la.CCMatrix    // c-c-mat form of Aku matrix [if reactions=true]
	matAkk *la.CCMatrix    // c-c-mat form of Akk matrix [if reactions=true]
	buCopy la.Vector       // copy of Bu vector if reactions == true
	linsol la.SparseSolver // linear solver
}

// NewGridSolver returns a new grid-based (e.g. FDM, SPC) solver
//  method   -- "fdm", "spc"(spectral-collocation), "fem"
//  gtype    -- grid type: "uni", "cgl" (Chebyshev-Gauss-Lobato)
//  operator -- differential operator; e.g. "laplacian"
//  params   -- parameters for operator; e.g. "kx" and "ky"
//  xmin     -- Grid: [ndim] min/initial coordinates of the whole space (box/cube)
//  xmax     -- Grid: [ndim] max/final coordinates of the whole space (box/cube)
//  ndiv     -- Grid: [ndim] number of divisions for xmax-xmin
func NewGridSolver(method, gtype, operator string, params dbf.Params, xmin, xmax []float64, ndiv []int) (o *Solver, err error) {

	// check lengths
	ndim := len(xmin)
	if len(xmax) != ndim {
		return nil, chk.Err("len(xmax) must be equal to len(xmin) == ndim. %d != %d\n", len(xmax), ndim)
	}
	if len(ndiv) != ndim {
		return nil, chk.Err("len(ndiv) must be equal to len(xmin) == ndim. %d != %d\n", len(ndiv), ndim)
	}

	// new solver and operator
	o = new(Solver)
	o.Op, err = NewOperator(method+"."+operator, params)
	if err != nil {
		return
	}

	// grid
	o.Grid, err = o.Op.InitWithGrid(gtype, xmin, xmax, ndiv)
	return
}

// SetBcs sets boundary conditions
func (o *Solver) SetBcs(ebcs *EssentialBcs) (err error) {

	// collect known equations
	o.Ebcs = ebcs
	knownEqs := o.Ebcs.GetNodesList()

	// init equations structure
	o.Eqs, err = la.NewEquations(o.Grid.Size(), knownEqs)
	if err != nil {
		return
	}
	o.U = la.NewVector(o.Eqs.N)
	o.F = la.NewVector(o.Eqs.N)

	// assemble matrices
	o.Op.Assemble(o.Eqs)
	if o.Eqs.Nk > 0 {
		o.matAuk = o.Eqs.Auk.ToMatrix(nil)
	}

	// init linear solver
	o.linsol = la.NewSparseSolver("umfpack")
	o.linsol.Init(o.Eqs.Auu, true, false, "", "", nil)
	err = o.linsol.Fact()
	return
}

// Solve solves problem
//  The solution will be saved in U (and F if reactions == true)
func (o *Solver) Solve(reactions bool) (err error) {

	// check
	if o.Eqs == nil {
		return chk.Err("please set boundary conditions first\n")
	}
	if o.Eqs.Nk == 0 {
		reactions = false
	}

	// auxiliary
	bu := o.Eqs.Bu
	xu := o.Eqs.Xu
	xk := o.Eqs.Xk

	// set RHS vector
	bu.Fill(0)
	err = o.Op.SourceTerm(o.Eqs)
	if err != nil {
		return
	}

	// set known part of RHS reactions vector
	if reactions {
		if o.matAku == nil {
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
	if o.Eqs.Nk > 0 {
		la.SpMatVecMulAdd(bu, -1.0, o.matAuk, xk)
	}

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

// Ugrid2d returns the U results converted to grid-shape
//   U -- [ny][nx] results at grid nodes
func (o *Solver) Ugrid2d() (uu [][]float64) {
	nx, ny := o.Grid.Npts(0), o.Grid.Npts(1)
	uu = utl.Alloc(ny, nx)
	for _, I := range o.Eqs.UtoF {
		m, n := I%nx, I/nx
		uu[n][m] = o.U[I]
	}
	for _, I := range o.Eqs.KtoF {
		m, n := I%nx, I/nx
		uu[n][m] = o.U[I]
	}
	return
}
