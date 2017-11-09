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
	s   fun.Svs               // source term function s({x},t)
	g   *gm.CurvGrid          // grid
	lip []*fun.LagrangeInterp // Lagrange interpolators [ndim]
}

// add to database
func init() {
	operatorDB["spc.laplacian"] = func(params dbf.Params, source fun.Svs) Operator {
		return newSpcLaplacian(params, source)
	}
}

// newSpcLaplacian creates a new Laplacian operator with given parameters
func newSpcLaplacian(params dbf.Params, source fun.Svs) (o *SpcLaplacian) {
	o = new(SpcLaplacian)
	e := params.ConnectSetOpt(
		[]*float64{&o.kx, &o.ky, &o.kz},
		[]string{"kx", "ky", "kz"},
		[]bool{false, false, true},
		"SpcLaplacian",
	)
	if e != "" {
		chk.Panic(e)
	}
	o.s = source
	return
}

// InitWithGrid initialises operator with new grid
func (o *SpcLaplacian) InitWithGrid(gtype string, xmin, xmax []float64, ndiv []int) (g *gm.CurvGrid) {

	// Lagrange interpolators
	ndim := len(xmin)
	o.lip = make([]*fun.LagrangeInterp, ndim)
	for i := 0; i < ndim; i++ {

		// allocate
		o.lip[i] = fun.NewLagrangeInterp(ndiv[i], gtype)

		// compute D2 matrix
		o.lip[i].CalcD2()
	}

	// new grid
	g = new(gm.CurvGrid)
	if ndim == 2 {
		g.RectSet2d(o.lip[0].X, o.lip[1].X)
	} else {
		g.RectSet3d(o.lip[0].X, o.lip[1].X, o.lip[2].X)
	}
	o.g = g

	// TODO:
	//  map [-1,+1] to xmin and xmax
	return
}

// Assemble assembles operator into A matrix from [A] ⋅ {u} = {b}
func (o *SpcLaplacian) Assemble(e *la.Equations) {
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
func (o *SpcLaplacian) SourceTerm(e *la.Equations, reactions bool) {
	if o.s == nil {
		return
	}
	for i, I := range e.UtoF {
		e.Bu[i] = o.s(o.g.Node(I), 0)
	}
	if reactions {
		for i, I := range e.KtoF {
			e.Bk[i] = o.s(o.g.Node(I), 0)
		}
	}
	return
}
