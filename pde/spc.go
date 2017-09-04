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

// SpcLaplacian implements the (negative) FDM Laplacian operator (2D or 3D)
//
//                ∂²u        ∂²u        ∂²u
//    L{u} = - kx ———  -  ky ———  -  kz ———
//                ∂x²        ∂y²        ∂z²
//
type SpcLaplacian struct {
	kx  float64               // isotropic coefficient x
	ky  float64               // isotropic coefficient y
	kz  float64               // isotropic coefficient z
	lip []*fun.LagrangeInterp // Lagrange interpolators [ndim]
}

// add to database
func init() {
	operatorDB["spc.laplacian"] = func(params dbf.Params) (Operator, error) {
		return newSpcLaplacian(params)
	}
}

// newSpcLaplacian creates a new Laplacian operator with given parameters
func newSpcLaplacian(params dbf.Params) (o *SpcLaplacian, err error) {
	o = new(SpcLaplacian)
	e := params.ConnectSetOpt(
		[]*float64{&o.kx, &o.ky, &o.kz},
		[]string{"kx", "ky", "kz"},
		[]bool{false, false, true},
		"SpcLaplacian",
	)
	if e != "" {
		err = chk.Err(e)
	}
	return
}

// InitWithGrid initialises operator with new grid
func (o *SpcLaplacian) InitWithGrid(gtype string, xmin, xmax []float64, ndiv []int) (g *gm.Grid, err error) {

	// Lagrange interpolators
	ndim := len(xmin)
	o.lip = make([]*fun.LagrangeInterp, ndim)
	for i := 0; i < ndim; i++ {

		// allocate
		o.lip[i], err = fun.NewLagrangeInterp(ndiv[i], gtype)
		if err != nil {
			return
		}

		// compute D2 matrix
		err = o.lip[i].CalcD2()
		if err != nil {
			return
		}
	}

	// new grid
	g = new(gm.Grid)
	if ndim == 2 {
		err = g.Set2d(o.lip[0].X, o.lip[1].X, false)
	} else {
		err = g.Set3d(o.lip[0].X, o.lip[1].X, o.lip[2].X, false)
	}

	// TODO:
	//  map [-1,+1] to xmin and xmax
	return
}

// Assemble assembles operator into A matrix from [A] ⋅ {u} = {b}
func (o *SpcLaplacian) Assemble(e *la.Equations) (err error) {
	ndim := len(o.lip)
	if ndim != 2 { // TODO
		return
	}
	nx := o.lip[0].N + 1
	ny := o.lip[1].N + 1
	if e.Auu == nil {
		nnz := (nx*nx)*ny + (ny*ny)*nx
		e.Alloc([]int{nnz, nnz, nnz, nnz}, true, true) // TODO: optimise nnz
	}
	e.Start()
	for k := 0; k < ny; k++ {
		for i := 0; i < nx; i++ {
			for j := 0; j < nx; j++ {
				e.Put(i+k*nx, j+k*nx, -o.kx*o.lip[0].D2.Get(i, j))
			}
		}
	}
	for k := 0; k < nx; k++ {
		for i := 0; i < ny; i++ {
			for j := 0; j < ny; j++ {
				e.Put(i*nx+k, j*nx+k, -o.ky*o.lip[1].D2.Get(i, j))
			}
		}
	}
	return
}

// SourceTerm assembles the source term vector
func (o *SpcLaplacian) SourceTerm(e *la.Equations) (err error) {
	return
}
