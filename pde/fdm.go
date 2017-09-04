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

// FdmLaplacian implements the (negative) FDM Laplacian operator (2D or 3D)
//
//                ∂²u        ∂²u        ∂²u
//    L{u} = - kx ———  -  ky ———  -  kz ———
//                ∂x²        ∂y²        ∂z²
//
type FdmLaplacian struct {
	kx float64  // isotropic coefficient x
	ky float64  // isotropic coefficient y
	kz float64  // isotropic coefficient z
	s  dbf.T    // source term function
	g  *gm.Grid // grid
}

// add to database
func init() {
	operatorDB["fdm.laplacian"] = func(params dbf.Params, source dbf.T) (Operator, error) {
		return newFdmLaplacian(params, source)
	}
}

// newFdmLaplacian creates a new Laplacian operator with given parameters
func newFdmLaplacian(params dbf.Params, source dbf.T) (o *FdmLaplacian, err error) {
	o = new(FdmLaplacian)
	e := params.ConnectSetOpt(
		[]*float64{&o.kx, &o.ky, &o.kz},
		[]string{"kx", "ky", "kz"},
		[]bool{false, false, true},
		"FdmLaplacian",
	)
	if e != "" {
		err = chk.Err(e)
	}
	o.s = source
	return
}

// InitWithGrid initialises operator with new grid
func (o *FdmLaplacian) InitWithGrid(gtype string, xmin, xmax []float64, ndiv []int) (g *gm.Grid, err error) {
	g = new(gm.Grid)
	err = g.GenUniform(xmin, xmax, ndiv, false)
	o.g = g
	return
}

// Assemble assembles operator into A matrix from [A] ⋅ {u} = {b}
func (o *FdmLaplacian) Assemble(e *la.Equations) (err error) {
	if e.Auu == nil {
		e.Alloc([]int{5 * e.Nu, 5 * e.Nu, 5 * e.Nk, 5 * e.Nk}, true, true)
	}
	e.Start()
	if o.g.Ndim() == 2 {
		nx := o.g.Npts(0)
		ny := o.g.Npts(1)
		dx := o.g.Length(0) / float64(nx-1)
		dy := o.g.Length(1) / float64(ny-1)
		dx2 := dx * dx
		dy2 := dy * dy
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
	}
	return
}

// SourceTerm assembles the source term vector
func (o *FdmLaplacian) SourceTerm(e *la.Equations) (err error) {
	return
}
