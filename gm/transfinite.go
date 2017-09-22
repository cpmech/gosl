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

// Transfinite maps a reference square [-1,+1]×[-1,+1] into a curve-bounded quadrilateral
//
//                                             B[2](r(x,y)) _,'\
//              B[2](r)                                  _,'    \ B[1](s(x,y))
//             ┌───────┐                              _,'        \
//             │       │                             \            \
//      B[3](s)│       │B[1](s)     ⇒                 \         _,'
//   s         │       │                  B[3](s(x,y)) \     _,'
//   │         └───────┘               y                \ _,'  B[0](r(x,y))
//   └──r       B[0](r)                │                 '
//                                     └──x
type Transfinite struct {

	// input
	Ndim int         // space dimension
	B    []fun.Vs    // the boundary functions
	Bd   []fun.Vs    // derivatives of boundary functions
	C    []la.Vector // "corner" points

	// workspase
	xface  []la.Vector // points at arbitrary positions along edges/faces
	dxface []la.Vector // derivatives at arbitrary positions along edges/faces
}

// NewTransfinite allocates a new structure
//  B  -- boundary functions x(s) = B(s)
//  Bd -- derivative functions dxds(s) = B'(s)
func NewTransfinite(ndim int, B, Bd []fun.Vs) (o *Transfinite) {
	o = new(Transfinite)
	o.Ndim = ndim
	o.B = B
	o.Bd = Bd
	if o.Ndim == 2 {
		if len(B) != 4 || len(Bd) != 4 {
			chk.Panic("in 2D, four boundary functions B are required\n")
		}
		o.C = make([]la.Vector, 4)
		o.xface = make([]la.Vector, 4)
		o.dxface = make([]la.Vector, 4)
		for i := 0; i < len(o.C); i++ {
			o.C[i] = la.NewVector(o.Ndim)
			o.xface[i] = la.NewVector(o.Ndim)
			o.dxface[i] = la.NewVector(o.Ndim)
		}
		o.B[0](o.C[0], -1)
		o.B[0](o.C[1], +1)
		o.B[2](o.C[2], +1)
		o.B[2](o.C[3], -1)
	} else if o.Ndim == 3 {
		if len(B) != 6 {
			chk.Panic("in 3D, six boundary functions B are required\n")
		}
	} else {
		chk.Panic("space dimension (ndim) must be 2 or 3\n")
	}
	return
}

// Point computes "real" position x(r,s,t)
//  Input:
//    u -- the "reference" coordinates {r,s,t} ϵ [-1,+1]×[-1,+1]×[-1,+1]
//  Output:
//    x -- the "real" coordinates {x,y,z}
func (o *Transfinite) Point(x, u la.Vector) {
	if o.Ndim == 2 {
		r, s := u[0], u[1]
		A, B, C, D := o.xface[0], o.xface[1], o.xface[2], o.xface[3]
		m, n, p, q := o.C[0], o.C[1], o.C[2], o.C[3]
		o.B[0](A, r)
		o.B[1](B, s)
		o.B[2](C, r)
		o.B[3](D, s)
		for i := 0; i < o.Ndim; i++ {
			x[i] = 0.5*((1-s)*A[i]+(1+r)*B[i]+(1+s)*C[i]+(1-r)*D[i]) -
				0.25*((1-r)*((1-s)*m[i]+(1+s)*q[i])+(1+r)*((1-s)*n[i]+(1+s)*p[i]))
		}
		return
	}
	chk.Panic("Point function is not ready for 3D yet\n")
}

// Derivs calculates derivatives (=metric terms) @ u={r,s,t}
//  Input:
//    u -- the "reference" coordinates {r,s,t} ϵ [-1,+1]×[-1,+1]×[-1,+1]
//  Output:
//    dxdu -- the derivatives [dx/du]ij = dxi/duj
//    x    -- the "real" coordinates {x,y,z}
func (o *Transfinite) Derivs(dxdu *la.Matrix, x, u la.Vector) {
	if o.Ndim == 2 {
		r, s := u[0], u[1]
		A, B, C, D := o.xface[0], o.xface[1], o.xface[2], o.xface[3]
		a, b, c, d := o.dxface[0], o.dxface[1], o.dxface[2], o.dxface[3]
		m, n, p, q := o.C[0], o.C[1], o.C[2], o.C[3]
		o.B[0](A, r)
		o.B[1](B, s)
		o.B[2](C, r)
		o.B[3](D, s)
		o.Bd[0](a, r)
		o.Bd[1](b, s)
		o.Bd[2](c, r)
		o.Bd[3](d, s)
		var dxidr, dxids float64
		for i := 0; i < o.Ndim; i++ {

			x[i] = 0.5*((1-s)*A[i]+(1+r)*B[i]+(1+s)*C[i]+(1-r)*D[i]) -
				0.25*((1-r)*((1-s)*m[i]+(1+s)*q[i])+(1+r)*((1-s)*n[i]+(1+s)*p[i]))

			dxidr = 0.5*((1-s)*a[i]+B[i]+(1+s)*c[i]-D[i]) -
				0.25*((1-s)*(n[i]-m[i])+(1+s)*(p[i]-q[i]))

			dxids = 0.5*(-A[i]+(1+r)*b[i]+C[i]+(1-r)*d[i]) -
				0.25*((1-r)*(q[i]-m[i])+(1+r)*(p[i]-n[i]))

			dxdu.Set(i, 0, dxidr)
			dxdu.Set(i, 1, dxids)
		}
		return
	}
	chk.Panic("Derivs function is not ready for 3D yet\n")
}

// Draw draws figure formed by B
func (o *Transfinite) Draw(npts []int, onlyBry bool, args, argsBry *plt.A) {

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
	u := la.NewVector(o.Ndim)
	x0 := make([]float64, npts[0])
	y0 := make([]float64, npts[0])
	x1 := make([]float64, npts[1])
	y1 := make([]float64, npts[1])

	if !onlyBry {
		// draw 0-lines
		for j := 0; j < npts[1]; j++ {
			u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
			for i := 0; i < npts[0]; i++ {
				u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
				o.Point(x, u)
				x0[i] = x[0]
				y0[i] = x[1]
			}
			plt.Plot(x0, y0, args)
		}

		// draw 1-lines
		for i := 0; i < npts[0]; i++ {
			u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
			for j := 0; j < npts[1]; j++ {
				u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
				o.Point(x, u)
				x1[j] = x[0]
				y1[j] = x[1]
			}
			plt.Plot(x1, y1, args)
		}
	}

	// draw B0(r)
	for i := 0; i < npts[0]; i++ {
		r := -1 + 2*float64(i)/float64(npts[0]-1)
		o.B[0](x, r)
		x0[i] = x[0]
		y0[i] = x[1]
	}
	plt.Plot(x0, y0, argsBry)

	// draw B1(s)
	for j := 0; j < npts[1]; j++ {
		s := -1 + 2*float64(j)/float64(npts[1]-1)
		o.B[1](x, s)
		x1[j] = x[0]
		y1[j] = x[1]
	}
	plt.Plot(x1, y1, argsBry)

	// draw B2(r)
	for i := 0; i < npts[0]; i++ {
		r := -1 + 2*float64(i)/float64(npts[0]-1)
		o.B[2](x, r)
		x0[i] = x[0]
		y0[i] = x[1]
	}
	plt.Plot(x0, y0, argsBry)

	// draw B3(s)
	for j := 0; j < npts[1]; j++ {
		s := -1 + 2*float64(j)/float64(npts[1]-1)
		o.B[3](x, s)
		x1[j] = x[0]
		y1[j] = x[1]
	}
	plt.Plot(x1, y1, argsBry)
}
