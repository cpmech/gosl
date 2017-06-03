// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
)

// Integrator implements methods to perform integrations over mesh cells
type Integrator struct {
	Nv  int         // number of vertices == len(X)
	Dim int         // space dimension == len(X[0])
	X   [][]float64 // coordinates of corners of polyhedron [nverts][ndim]
}

// NewIntegrator returns a new integrator
func NewIntegrator(cellId int, mesh *Mesh) (o *Integrator, err error) {
	o = new(Integrator)
	o.X = mesh.ExtractCellCoords(cellId)
	if len(o.X) < 2 {
		err = chk.Err("at least two vertices are required. len(X)=%d is invalid", len(o.X))
		return
	}
	o.Nv = len(o.X)
	o.Dim = len(o.X[0])
	return
}

// IntegrateSv integrates scalar function of vector argument over Cell
//
//   Computes:
//
//           ⌠⌠⌠   →       ⌠⌠⌠   → →     →        n-1   →  →      →
//     res = │││ f(x) dΩ = │││ f(x(r))⋅J(r) dΩr ≈  Σ  f(xi(ri))⋅J(ri)⋅wi
//           ⌡⌡⌡           ⌡⌡⌡                    i=0
//              Ω             Ωr
//   where:
//            → →    m-1   m →     m
//            x(r) ≈  Σ   S (r) ⋅ x                   J = det(Jmat)
//                   i=0
//
//           →     _                           _
//          dx    |  ∂x0/∂r0  ∂x0/∂r1  ∂x0/∂r2  |                ∂xi
//   Jmat = —— =  |  ∂x1/∂r0  ∂x1/∂r1  ∂x1/∂r2  |     Jmat[ij] = ———
//           →    |_ ∂x2/∂r0  ∂x2/∂r1  ∂x2/∂r2 _|                ∂rj
//          dr
//
//   and:
//            m -- number of cell nodes
//            n -- number of integration points
//   Input:
//     pts -- (Gauss) integration points. May be nil => a default set will be selected then
func (o *Integrator) IntegrateSv(f fun.Sv) (res float64, err error) {
	// TODO
	return
}
