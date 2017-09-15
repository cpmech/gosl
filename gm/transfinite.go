// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

// Transfinite maps a reference square [-1,1]×[-1,1] into a curve-bounded quadrilateral
//
//                                   Γ1(η) _,B
//       c———————b                C--.___-'  |
//       |       |                 \         | Γ0(ξ)
//   η   |       |      y     Γ2(ξ) \        |
//   ↑   |       |      ↑           /      __A
//   |   d———————a      |          / ____,'
//   |                  |         D-'   Γ3(η)
//   +————→ξ            +————→x
//
type Transfinite struct {
	Ndim int         // space dimension
	Γ    []fun.Vs    // the boundary functions
	C    []la.Vector // "corner" points
	X    []la.Vector // points on arbitrary positions along edges/faces
}

// NewTransfinite allocates a new structure
func NewTransfinite(ndim int, Γ []fun.Vs) (o *Transfinite) {
	o = new(Transfinite)
	o.Ndim = ndim
	o.Γ = Γ
	if o.Ndim == 2 {
		if len(Γ) != 4 {
			chk.Panic("in 2D, four boundary functions Γ are required\n")
		}
		o.C = make([]la.Vector, 4)
		o.X = make([]la.Vector, 4)
		for i := 0; i < len(o.C); i++ {
			o.C[i] = la.NewVector(o.Ndim)
			o.X[i] = la.NewVector(o.Ndim)
		}
		o.Γ[0](o.C[0], -1)
		o.Γ[0](o.C[1], +1)
		o.Γ[2](o.C[2], +1)
		o.Γ[2](o.C[3], -1)
	} else if o.Ndim == 3 {
		if len(Γ) != 6 {
			chk.Panic("in 3D, six boundary functions Γ are required\n")
		}
	} else {
		chk.Panic("space dimension (ndim) must be 2 or 3\n")
	}
	return
}

// QuadMap computes "real" position x(ξ,η)
//  Input:
//    r -- the "reference" coordinates {ξ,η}
//
//  Output:
//    x -- the "real" coordinates {x,y}
//
func (o *Transfinite) QuadMap(x, r la.Vector) {
	if o.Ndim == 2 {
		ξ, η := r[0], r[1]
		o.Γ[0](o.X[0], ξ)
		o.Γ[1](o.X[1], η)
		o.Γ[2](o.X[2], ξ)
		o.Γ[3](o.X[3], η)
		for i := 0; i < o.Ndim; i++ {
			x[i] = 0.5*((1-ξ)*o.X[3][i]+(1+ξ)*o.X[1][i]+(1-η)*o.X[0][i]+(1+η)*o.X[2][i]) -
				0.25*((1-ξ)*((1-η)*o.C[0][i]+(1+η)*o.C[3][i])+(1+ξ)*((1-η)*o.C[1][i]+(1+η)*o.C[2][i]))
		}
		return
	}
	chk.Panic("QuadMap function is not ready for 3D yet\n")
}

// Draw draws figure formed by Γ
func (o *Transfinite) Draw(npts []int, args, argsBry *plt.A) {

	// auxiliary
	if len(npts) != o.Ndim {
		npts = make([]int, o.Ndim)
	}
	for i := 0; i < o.Ndim; i++ {
		if npts[i] < 3 {
			npts[i] = 3
		}
	}
	if args == nil {
		args = &plt.A{C: plt.C(0, 0), NoClip: true}
	}
	if argsBry == nil {
		argsBry = &plt.A{C: plt.C(0, 0), Lw: 2, NoClip: true}
	}
	x := la.NewVector(o.Ndim)
	r := la.NewVector(o.Ndim)

	// draw 0-lines
	x0 := make([]float64, npts[0])
	y0 := make([]float64, npts[0])
	for j := 0; j < npts[1]; j++ {
		r[1] = -1 + 2*float64(j)/float64(npts[1]-1)
		for i := 0; i < npts[0]; i++ {
			r[0] = -1 + 2*float64(i)/float64(npts[0]-1)
			o.QuadMap(x, r)
			x0[i] = x[0]
			y0[i] = x[1]
		}
		plt.Plot(x0, y0, args)
	}

	// draw 1-lines
	x1 := make([]float64, npts[1])
	y1 := make([]float64, npts[1])
	for i := 0; i < npts[0]; i++ {
		r[0] = -1 + 2*float64(i)/float64(npts[0]-1)
		for j := 0; j < npts[1]; j++ {
			r[1] = -1 + 2*float64(j)/float64(npts[1]-1)
			o.QuadMap(x, r)
			x1[j] = x[0]
			y1[j] = x[1]
		}
		plt.Plot(x1, y1, args)
	}

	// draw Γ0(ξ)
	for i := 0; i < npts[0]; i++ {
		ξ := -1 + 2*float64(i)/float64(npts[0]-1)
		o.Γ[0](x, ξ)
		x0[i] = x[0]
		y0[i] = x[1]
	}
	plt.Plot(x0, y0, argsBry)

	// draw Γ1(η)
	for j := 0; j < npts[1]; j++ {
		η := -1 + 2*float64(j)/float64(npts[1]-1)
		o.Γ[1](x, η)
		x1[j] = x[0]
		y1[j] = x[1]
	}
	plt.Plot(x1, y1, argsBry)

	// draw Γ2(ξ)
	for i := 0; i < npts[0]; i++ {
		ξ := -1 + 2*float64(i)/float64(npts[0]-1)
		o.Γ[2](x, ξ)
		x0[i] = x[0]
		y0[i] = x[1]
	}
	plt.Plot(x0, y0, argsBry)

	// draw Γ3(η)
	for j := 0; j < npts[1]; j++ {
		η := -1 + 2*float64(j)/float64(npts[1]-1)
		o.Γ[3](x, η)
		x1[j] = x[0]
		y1[j] = x[1]
	}
	plt.Plot(x1, y1, argsBry)
}
