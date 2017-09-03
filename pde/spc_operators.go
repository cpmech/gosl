// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/la"
)

// database ////////////////////////////////////////////////////////////////////////////////////////

// SpcOperator defines the interface for SPC (spectral collocation) operators such as
// the Laplacian and so on
type SpcOperator interface {
	Assemble(L []*fun.LagrangeInterp, e *la.Equations)
}

// spcOperatorMaker defines a function that makes (allocates) SpcOperators
type spcOperatorMaker func(params dbf.Params) (SpcOperator, error)

// spcOperatorDB implemetns a database of SpcOperators
var spcOperatorDB = make(map[string]spcOperatorMaker)

// NewSpcOperator finds a SpcOperator in database or panic
func NewSpcOperator(kind string, params dbf.Params) (SpcOperator, error) {
	if maker, ok := spcOperatorDB[kind]; ok {
		return maker(params)
	}
	return nil, chk.Err("cannot find SpcOperator named %q in database", kind)
}

// implementation: Laplacian ///////////////////////////////////////////////////////////////////////

// SpcLaplacian implements the (negative) FDM Laplacian operator (2D or 3D)
//
//                ∂²u        ∂²u        ∂²u
//    L{u} = - kx ———  -  ky ———  -  kz ———
//                ∂x²        ∂y²        ∂z²
//
type SpcLaplacian struct {
	kx float64 // isotropic coefficient x
	ky float64 // isotropic coefficient y
	kz float64 // isotropic coefficient z
}

// add to database
func init() {
	spcOperatorDB["laplacian"] = func(params dbf.Params) (SpcOperator, error) {
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

// Assemble assembles operator into A matrix from [A] ⋅ {u} = {b}
func (o *SpcLaplacian) Assemble(L []*fun.LagrangeInterp, e *la.Equations) {
	ndim := len(L)
	if ndim != 2 { // TODO
		return
	}
	nx := L[0].N + 1
	ny := L[1].N + 1
	if e.Auu == nil {
		nnz := (nx*nx)*ny + (ny*ny)*nx
		e.Alloc([]int{nnz, nnz, nnz, nnz}, true, true) // TODO: optimise nnz
	}
	e.Start()
	for k := 0; k < ny; k++ {
		for i := 0; i < nx; i++ {
			for j := 0; j < nx; j++ {
				e.Put(i+k*nx, j+k*nx, L[0].D2.Get(i, j))
			}
		}
	}
	for k := 0; k < nx; k++ {
		for i := 0; i < ny; i++ {
			for j := 0; j < ny; j++ {
				e.Put(i*nx+k, j*nx+k, L[1].D2.Get(i, j))
			}
		}
	}
}
