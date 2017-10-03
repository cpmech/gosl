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
//                                             B[3](r(x,y)) _,'\
//              B[3](r)                                  _,'    \ B[1](s(x,y))
//             ┌───────┐                              _,'        \
//             │       │                             \            \
//      B[0](s)│       │B[1](s)     ⇒                 \         _,'
//   s         │       │                  B[0](s(x,y)) \     _,'
//   │         └───────┘               y                \ _,'  B[2](r(x,y))
//   └──r       B[2](r)                │                 '
//                                     └──x
type Transfinite struct {

	// input
	ndim int      // space dimension
	b    []fun.Vs // the boundary functions
	bd   []fun.Vs // derivatives of boundary functions

	// workspase
	p0, p1, p2, p3         la.Vector // corner points
	p4, p5, p6, p7         la.Vector // corner points
	b0s, b1s, b2r, b3r     la.Vector // 2d function evaluations
	bd0s, bd1s, bd2r, bd3r la.Vector // 2d derivatives evaluations
}

// NewTransfinite allocates a new structure
//  B  -- boundary functions x(s) = B(s)
//  Bd -- derivative functions dxds(s) = B'(s)
func NewTransfinite(ndim int, B, Bd []fun.Vs) (o *Transfinite) {
	o = new(Transfinite)
	o.ndim = ndim
	o.b = B
	o.bd = Bd
	o.p0 = la.NewVector(o.ndim)
	o.p1 = la.NewVector(o.ndim)
	o.p2 = la.NewVector(o.ndim)
	o.p3 = la.NewVector(o.ndim)
	o.p4 = la.NewVector(o.ndim)
	o.p5 = la.NewVector(o.ndim)
	o.p6 = la.NewVector(o.ndim)
	o.p7 = la.NewVector(o.ndim)
	if o.ndim == 2 {
		if len(B) != 4 || len(Bd) != 4 {
			chk.Panic("in 2D, four boundary functions B are required\n")
		}
		o.b0s = la.NewVector(o.ndim)
		o.b1s = la.NewVector(o.ndim)
		o.b2r = la.NewVector(o.ndim)
		o.b3r = la.NewVector(o.ndim)
		o.bd0s = la.NewVector(o.ndim)
		o.bd1s = la.NewVector(o.ndim)
		o.bd2r = la.NewVector(o.ndim)
		o.bd3r = la.NewVector(o.ndim)
		o.b[0](o.p0, -1)
		o.b[0](o.p3, +1)
		o.b[1](o.p1, -1)
		o.b[1](o.p2, +1)
	} else if o.ndim == 3 {
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
	if o.ndim == 2 {
		r, s := u[0], u[1]
		o.b[0](o.b0s, s)
		o.b[1](o.b1s, s)
		o.b[2](o.b2r, r)
		o.b[3](o.b3r, r)
		for i := 0; i < o.ndim; i++ {
			x[i] = 0 + // trick to enforce alignment
				+(1.0-r)*o.b0s[i]/2.0 + (1.0+r)*o.b1s[i]/2.0 +
				+(1.0-s)*o.b2r[i]/2.0 + (1.0+s)*o.b3r[i]/2.0 +
				-(1.0-r)*(1.0-s)*o.p0[i]/4.0 - (1.0+r)*(1.0-s)*o.p1[i]/4.0 +
				-(1.0+r)*(1.0+s)*o.p2[i]/4.0 - (1.0-r)*(1.0+s)*o.p3[i]/4.0
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
func (o *Transfinite) Derivs(dxdu *la.Matrix, u la.Vector) {
	if o.ndim == 2 {
		r, s := u[0], u[1]
		o.b[0](o.b0s, s)
		o.b[1](o.b1s, s)
		o.b[2](o.b2r, r)
		o.b[3](o.b3r, r)
		o.bd[0](o.bd0s, s)
		o.bd[1](o.bd1s, s)
		o.bd[2](o.bd2r, r)
		o.bd[3](o.bd3r, r)
		var dxidr, dxids float64
		for i := 0; i < o.ndim; i++ {

			dxidr = 0 + // trick to enforce alignment
				-o.b0s[i]/2.0 + o.b1s[i]/2.0 +
				+(1.0-s)*o.bd2r[i]/2.0 + (1.0+s)*o.bd3r[i]/2.0 +
				+(1.0-s)*o.p0[i]/4.0 - (1.0-s)*o.p1[i]/4.0 +
				-(1.0+s)*o.p2[i]/4.0 + (1.0+s)*o.p3[i]/4.0

			dxids = 0 + // trick to enforce alignment
				+(1.0-r)*o.bd0s[i]/2.0 + (1.0+r)*o.bd1s[i]/2.0 +
				-o.b2r[i]/2.0 + o.b3r[i]/2.0 +
				+(1.0-r)*o.p0[i]/4.0 + (1.0+r)*o.p1[i]/4.0 +
				-(1.0+r)*o.p2[i]/4.0 - (1.0-r)*o.p3[i]/4.0

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
	if len(npts) != o.ndim {
		npts = make([]int, o.ndim)
	}
	for i := 0; i < o.ndim; i++ {
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
	x := la.NewVector(o.ndim)
	u := la.NewVector(o.ndim)
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
		o.b[0](x, r)
		x0[i] = x[0]
		y0[i] = x[1]
	}
	plt.Plot(x0, y0, argsBry)

	// draw B1(s)
	for j := 0; j < npts[1]; j++ {
		s := -1 + 2*float64(j)/float64(npts[1]-1)
		o.b[1](x, s)
		x1[j] = x[0]
		y1[j] = x[1]
	}
	plt.Plot(x1, y1, argsBry)

	// draw B2(r)
	for i := 0; i < npts[0]; i++ {
		r := -1 + 2*float64(i)/float64(npts[0]-1)
		o.b[2](x, r)
		x0[i] = x[0]
		y0[i] = x[1]
	}
	plt.Plot(x0, y0, argsBry)

	// draw B3(s)
	for j := 0; j < npts[1]; j++ {
		s := -1 + 2*float64(j)/float64(npts[1]-1)
		o.b[3](x, s)
		x1[j] = x[0]
		y1[j] = x[1]
	}
	plt.Plot(x1, y1, argsBry)
}
