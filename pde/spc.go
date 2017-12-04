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

// SpcLaplacian implements the Spectral Collocation (SPC) Laplacian operator (2D or 3D)
//
//              ∂²u        ∂²u        ∂²u
//    L{u} = kx ———  +  ky ———  +  kz ———
//              ∂x²        ∂y²        ∂z²
//
type SpcLaplacian struct {
	Kx       float64       // isotropic coefficient x
	Ky       float64       // isotropic coefficient y
	Kz       float64       // isotropic coefficient z
	LagInt   fun.LagIntSet // Lagrange interpolators [ndim]
	Grid     *gm.Grid      // grid
	Source   fun.Svs       // source term function s({x},t)
	EssenBcs *EssentialBcs // essential boundary conditions
	Eqs      *la.Equations // equations
	bcsReady bool          // boundary conditions are set
}

// NewSpcLaplacian creates a new SPC Laplacian operator with given parameters
func NewSpcLaplacian(params dbf.Params, lis fun.LagIntSet, grid *gm.Grid, source fun.Svs) (o *SpcLaplacian) {
	o = new(SpcLaplacian)
	err := params.ConnectSetOpt(
		[]*float64{&o.Kx, &o.Ky, &o.Kz},
		[]string{"kx", "ky", "kz"},
		[]bool{false, false, true},
		"SpcLaplacian",
	)
	if err != "" {
		chk.Panic(err)
	}
	o.LagInt = lis
	o.Grid = grid
	o.Source = source
	o.EssenBcs = NewEssentialBcsGrid(grid, 1) // 1:maxNdof
	o.bcsReady = false
	return
}

// AddBc adds essential or natural boundary condition
//   essential -- essential BC; otherwise natural boundary condition
//   tag       -- edge or face tag in grid
//   cvalue    -- constant value [optional]; or
//   fvalue    -- function value [optional]
func (o *SpcLaplacian) AddBc(essential bool, tag int, cvalue float64, fvalue dbf.T) {
	o.bcsReady = false
	if essential {
		o.EssenBcs.AddUsingTag(tag, 0, cvalue, fvalue)
		return
	}
	chk.Panic("TODO: Implement natural boundary condition\n")
}

// Assemble assembles operator into A matrix from [A] ⋅ {u} = {b}
//  reactions -- prepare for computation of RHS
func (o *SpcLaplacian) Assemble(reactions bool) {
	nx := o.LagInt[0].N + 1
	ny := o.LagInt[1].N + 1
	if !o.bcsReady {
		nnz := (nx*nx)*ny + (ny*ny)*nx
		o.Eqs = la.NewEquations(o.Grid.Size(), o.EssenBcs.Nodes())
		o.Eqs.Alloc([]int{nnz, nnz, nnz, nnz}, reactions, true) // TODO: optimise nnz
		for _, li := range o.LagInt {
			li.CalcD2()
		}
		o.bcsReady = true
	}
	o.Eqs.Start()
	if o.Grid.Ndim() == 2 {
		for k := 0; k < ny; k++ {
			for i := 0; i < nx; i++ {
				for j := 0; j < nx; j++ {
					o.Eqs.Put(i+k*nx, j+k*nx, o.Kx*o.LagInt[0].D2.Get(i, j))
				}
			}
		}
		for k := 0; k < nx; k++ {
			for i := 0; i < ny; i++ {
				for j := 0; j < ny; j++ {
					o.Eqs.Put(i*nx+k, j*nx+k, o.Ky*o.LagInt[1].D2.Get(i, j))
				}
			}
		}
		return
	}
	chk.Panic("TODO: Implement Assemble() in 3D\n")
}

// SolveSteady solves steady problem
//   Solves: [K]⋅{u} = {f} represented by [A]⋅{x} = {b}
func (o *SpcLaplacian) SolveSteady(reactions bool) (u, f []float64) {
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
func (o *SpcLaplacian) calcXk(I int, t float64) float64 {
	val, available := o.EssenBcs.Value(I, 0, t)
	if available {
		return val
	}
	return 0
}

// calcBu calculates RHS vector (e.g. source) corresponding to known values of {u} (CalcBu in la.Equations)
//  I -- node number
//  t -- time
func (o *SpcLaplacian) calcBu(I int, t float64) float64 {
	if o.Source != nil {
		return o.Source(o.Grid.Node(I), t)
	}
	return 0
}
