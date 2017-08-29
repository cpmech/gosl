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

// database ////////////////////////////////////////////////////////////////////////////////////////

// FdmOperator defines the interface for FDM (finite-difference method) operators such as
// the Laplacian and so on
type FdmOperator interface {
	Init(params dbf.Params) (err error)
	Assemble(g *gm.Grid, e *la.Equations)
}

// fdmOperatorMaker defines a function that makes (allocates) FdmOperators
type fdmOperatorMaker func() FdmOperator

// fdmOperatorDB implemetns a database of FdmOperators
var fdmOperatorDB = make(map[string]fdmOperatorMaker)

// NewFdmOperator finds a FdmOperator in database or panic
func NewFdmOperator(kind string) (FdmOperator, error) {
	if maker, ok := fdmOperatorDB[kind]; ok {
		return maker(), nil
	}
	return nil, chk.Err("cannot find FdmOperator named %q in database", kind)
}

// implementation: Laplacian ///////////////////////////////////////////////////////////////////////

// FdmLaplacian implements the (negative) FDM Laplacian operator (2D or 3D)
//
//                ∂²u        ∂²u        ∂²u
//    L{u} = - kx ———  -  ky ———  -  kz ———
//                ∂x²        ∂y²        ∂z²
//
type FdmLaplacian struct {
	kx float64 // isotropic coefficient x
	ky float64 // isotropic coefficient y
	kz float64 // isotropic coefficient z
}

// add to database
func init() {
	fdmOperatorDB["laplacian"] = func() FdmOperator { return new(FdmLaplacian) }
}

// Init initialises operator with given parameters
func (o *FdmLaplacian) Init(params dbf.Params) (err error) {
	e := params.ConnectSetOpt(
		[]*float64{&o.kx, &o.ky, &o.kz},
		[]string{"kx", "ky", "kz"},
		[]bool{false, false, true},
		"FdmLaplacian",
	)
	if e != "" {
		err = chk.Err(e)
	}
	return
}

// Assemble assembles operator into A matrix from [A] ⋅ {u} = {b}
func (o *FdmLaplacian) Assemble(g *gm.Grid, e *la.Equations) {
	if e.Auu == nil {
		e.Alloc([]int{5 * e.Nu, 5 * e.Nu, 5 * e.Nk, 5 * e.Nk}, true, true)
	}
	e.Start()
	if g.Ndim == 2 {
		nx := g.Npts[0]
		ny := g.Npts[1]
		dx2 := g.Size[0] * g.Size[0]
		dy2 := g.Size[1] * g.Size[1]
		α := 2.0 * (o.kx/dx2 + o.ky/dy2)
		β := -o.kx / dx2
		γ := -o.ky / dy2
		mol := []float64{α, β, β, γ, γ}
		jays := make([]int, 5)
		for I := 0; I < e.N; I++ { // loop over all Nx*Ny equations
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
				e.Put(I, J, mol[k])
			}
		}
		return
	}
}
