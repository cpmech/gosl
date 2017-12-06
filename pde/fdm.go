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

// FdmLaplacian implements the Finite Difference (FDM) Laplacian operator (2D or 3D)
//
//              ∂²u        ∂²u        ∂²u
//    L{u} = kx ———  +  ky ———  +  kz ———
//              ∂x²        ∂y²        ∂z²
//
type FdmLaplacian struct {
	Kx       float64        // isotropic coefficient x
	Ky       float64        // isotropic coefficient y
	Kz       float64        // isotropic coefficient z
	Grid     *gm.Grid       // grid
	Source   fun.Svs        // source term function s({x},t)
	EssenBcs *BoundaryConds // essential boundary conditions
	Eqs      *la.Equations  // equations
	bcsReady bool           // boundary conditions are set
}

// NewFdmLaplacian creates a new FDM Laplacian operator with given parameters
func NewFdmLaplacian(params dbf.Params, grid *gm.Grid, source fun.Svs) (o *FdmLaplacian) {
	o = new(FdmLaplacian)
	err := params.ConnectSetOpt(
		[]*float64{&o.Kx, &o.Ky, &o.Kz},
		[]string{"kx", "ky", "kz"},
		[]bool{false, false, true},
		"FdmLaplacian",
	)
	if err != "" {
		chk.Panic(err)
	}
	o.Grid = grid
	o.Source = source
	o.EssenBcs = NewBoundaryCondsGrid(grid, 1) // 1:maxNdof
	o.bcsReady = false
	return
}

// AddEbc adds essential boundary condition given tag of edge or face
//   tag    -- edge or face tag in grid
//   cvalue -- constant value [optional]; or
//   fvalue -- function value [optional]
func (o *FdmLaplacian) AddEbc(tag int, cvalue float64, fvalue fun.Svs) {
	o.bcsReady = false
	o.EssenBcs.AddUsingTag(tag, 0, cvalue, fvalue)
}

// SetHbc sets homogeneous boundary conditions; i.e. all boundaries with zero EBC
func (o *FdmLaplacian) SetHbc() {
	if o.Grid.Ndim() == 2 {
		o.AddEbc(10, 0.0, nil)
		o.AddEbc(11, 0.0, nil)
		o.AddEbc(20, 0.0, nil)
		o.AddEbc(21, 0.0, nil)
		return
	}
	o.AddEbc(100, 0.0, nil)
	o.AddEbc(101, 0.0, nil)
	o.AddEbc(200, 0.0, nil)
	o.AddEbc(201, 0.0, nil)
	o.AddEbc(300, 0.0, nil)
	o.AddEbc(301, 0.0, nil)
}

// Assemble assembles operator into A matrix from [A] ⋅ {u} = {b}
//  reactions -- prepare for computation of RHS
func (o *FdmLaplacian) Assemble(reactions bool) {
	if !o.bcsReady {
		o.Eqs = la.NewEquations(o.Grid.Size(), o.EssenBcs.Nodes())
		o.Eqs.Alloc([]int{5 * o.Eqs.Nu, 5 * o.Eqs.Nu, 5 * o.Eqs.Nk, 5 * o.Eqs.Nk}, reactions, true)
		o.bcsReady = true
	}
	o.Eqs.Start()
	if o.Grid.Ndim() == 2 {
		nx := o.Grid.Npts(0)
		ny := o.Grid.Npts(1)
		dx := o.Grid.Xlen(0) / float64(nx-1)
		dy := o.Grid.Xlen(1) / float64(ny-1)
		dx2 := dx * dx
		dy2 := dy * dy
		α := -2.0 * (o.Kx/dx2 + o.Ky/dy2)
		β := o.Kx / dx2
		γ := o.Ky / dy2
		mol := []float64{α, β, β, γ, γ}
		jays := make([]int, 5)
		for I := 0; I < o.Eqs.N; I++ { // loop over all Nx*Ny equations
			col := I % nx    // grid column number
			row := I / nx    // grid row number
			jays[0] = I      // current node
			jays[1] = I - 1  // left node
			jays[2] = I + 1  // right node
			jays[3] = I - nx // bottom node
			jays[4] = I + nx // top node
			if col == 0 {
				jays[1] = jays[2]
			}
			if col == nx-1 {
				jays[2] = jays[1]
			}
			if row == 0 {
				jays[3] = jays[4]
			}
			if row == ny-1 {
				jays[4] = jays[3]
			}
			for k, J := range jays { // loop over non-zero columns
				o.Eqs.Put(I, J, mol[k])
			}
		}
		return
	}
	chk.Panic("TODO: Implement Assemble() in 3D\n")
}

// SolveSteady solves steady problem
//   Solves: [K]⋅{u} = {f} represented by [A]⋅{x} = {b}
func (o *FdmLaplacian) SolveSteady(reactions bool) (u, f []float64) {
	o.Eqs.SolveOnce(o.calcXk, o.calcBu)
	u = make([]float64, o.Grid.Size())
	o.Eqs.JoinVector(u, o.Eqs.Xu, o.Eqs.Xk)
	if reactions {
		f = make([]float64, o.Grid.Size())
		if o.Eqs.Nk > 0 { // need to calc Bu again because it was modified
			for i, I := range o.Eqs.UtoF {
				o.Eqs.Bu[i] = o.calcBu(I, 0)
			}
		}
		o.Eqs.JoinVector(f, o.Eqs.Bu, o.Eqs.Bk)
	}
	return
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// calcXk calculates known {u} values (CalcXk in la.Equations)
//  I -- node number
//  t -- time
func (o *FdmLaplacian) calcXk(I int, t float64) float64 {
	_, val, available := o.EssenBcs.Value(I, 0, t)
	if available {
		return val
	}
	return 0
}

// calcBu calculates RHS vector (e.g. source) corresponding to known values of {u} (CalcBu in la.Equations)
//  I -- node number
//  t -- time
func (o *FdmLaplacian) calcBu(I int, t float64) float64 {
	if o.Source != nil {
		return o.Source(o.Grid.Node(I), t)
	}
	return 0
}
