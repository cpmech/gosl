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
	ndim int       // space dimension
	b    []fun.Vs  // 2D: the boundary functions
	bd   []fun.Vs  // 2D: derivatives of boundary functions
	b3d  []fun.Vss // 3D: the boundary functions
	bd3d []fun.Mss // 3D: derivatives of boundary functions

	// workspace for 2d
	p0, p1     la.Vector // corner points
	p2, p3     la.Vector // corner points
	b0s, b1s   la.Vector // 2d function evaluations
	b2r, b3r   la.Vector // 2d function evaluations
	bd0s, bd1s la.Vector // 2d derivatives evaluations
	bd2r, bd3r la.Vector // 2d derivatives evaluations

	// workspace for 3d
	p4, p5     la.Vector // corner points
	p6, p7     la.Vector // corner points
	b0st, b1st la.Vector // 3d function evaluations
	b2rt, b3rt la.Vector // 3d function evaluations
	b4rs, b5rs la.Vector // 3d function evaluations
	b0mt, b0pt la.Vector // 3d function evaluations
	b1mt, b1pt la.Vector // 3d function evaluations
	b0sm, b0sp la.Vector // 3d function evaluations
	b1sm, b1sp la.Vector // 3d function evaluations
	b2rm, b2rp la.Vector // 3d function evaluations
	b3rm, b3rp la.Vector // 3d function evaluations

	// derivatives for 3d
	bdr2rt, bdr3rt *la.Matrix
	bdr4rs, bdr5rs *la.Matrix
	bdr2rm, bdr2rp *la.Matrix
	bdr3rm, bdr3rp *la.Matrix

	bd0st, bd1st *la.Matrix
	bd4rs, bd5rs *la.Matrix
	bd0sm, bd0sp *la.Matrix
	bd1sm, bd1sp *la.Matrix

	bd2rt, bd3rt *la.Matrix
	bd0mt, bd0pt *la.Matrix
	bd1mt, bd1pt *la.Matrix
}

// NewTransfinite2d allocates a new structure
//  B  -- boundary functions x(s) = B(s)
//  Bd -- derivative functions dxds(s) = B'(s)
func NewTransfinite2d(ndim int, B, Bd []fun.Vs) (o *Transfinite) {
	if len(B) != 4 || len(Bd) != 4 {
		chk.Panic("in 2D, four boundary functions B are required\n")
	}
	o = new(Transfinite)
	o.ndim = ndim
	o.b = B
	o.bd = Bd
	o.p0, o.p1 = la.NewVector(2), la.NewVector(2)
	o.p2, o.p3 = la.NewVector(2), la.NewVector(2)
	o.b0s, o.b1s = la.NewVector(2), la.NewVector(2)
	o.b2r, o.b3r = la.NewVector(2), la.NewVector(2)
	o.bd0s, o.bd1s = la.NewVector(2), la.NewVector(2)
	o.bd2r, o.bd3r = la.NewVector(2), la.NewVector(2)
	o.b[0](o.p0, -1)
	o.b[0](o.p3, +1)
	o.b[1](o.p1, -1)
	o.b[1](o.p2, +1)
	return
}

// NewTransfinite3d allocates a new structure
//  B  -- boundary functions x(s) = B(r,s)
//  Bd -- derivative functions dxds(s) = B'(r,s)
func NewTransfinite3d(B []fun.Vss, Bd []fun.Mss) (o *Transfinite) {
	if len(B) != 6 {
		chk.Panic("in 3D, six boundary functions B are required\n")
	}
	o = new(Transfinite)
	o.ndim = 3
	o.b3d = B
	o.bd3d = Bd
	o.p0, o.p1 = la.NewVector(3), la.NewVector(3)
	o.p2, o.p3 = la.NewVector(3), la.NewVector(3)
	o.p4, o.p5 = la.NewVector(3), la.NewVector(3)
	o.p6, o.p7 = la.NewVector(3), la.NewVector(3)
	o.b0st, o.b1st = la.NewVector(3), la.NewVector(3)
	o.b2rt, o.b3rt = la.NewVector(3), la.NewVector(3)
	o.b4rs, o.b5rs = la.NewVector(3), la.NewVector(3)
	o.b0mt, o.b0pt = la.NewVector(3), la.NewVector(3)
	o.b1mt, o.b1pt = la.NewVector(3), la.NewVector(3)
	o.b0sm, o.b0sp = la.NewVector(3), la.NewVector(3)
	o.b1sm, o.b1sp = la.NewVector(3), la.NewVector(3)
	o.b2rm, o.b2rp = la.NewVector(3), la.NewVector(3)
	o.b3rm, o.b3rp = la.NewVector(3), la.NewVector(3)
	o.b3d[4](o.p0, -1.0, -1.0)
	o.b3d[4](o.p1, +1.0, -1.0)
	o.b3d[4](o.p2, +1.0, +1.0)
	o.b3d[4](o.p3, -1.0, +1.0)
	o.b3d[5](o.p4, -1.0, -1.0)
	o.b3d[5](o.p5, +1.0, -1.0)
	o.b3d[5](o.p6, +1.0, +1.0)
	o.b3d[5](o.p7, -1.0, +1.0)

	o.bdr2rt, o.bdr3rt = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	o.bdr4rs, o.bdr5rs = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	o.bdr2rm, o.bdr2rp = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	o.bdr3rm, o.bdr3rp = la.NewMatrix(3, 3), la.NewMatrix(3, 3)

	o.bd0st, o.bd1st = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	o.bd4rs, o.bd5rs = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	o.bd0sm, o.bd0sp = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	o.bd1sm, o.bd1sp = la.NewMatrix(3, 3), la.NewMatrix(3, 3)

	o.bd2rt, o.bd3rt = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	o.bd0mt, o.bd0pt = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	o.bd1mt, o.bd1pt = la.NewMatrix(3, 3), la.NewMatrix(3, 3)
	return
}

// Point computes "real" position x(r,s,t)
//  Input:
//    u -- the "reference" coordinates {r,s,t} ϵ [-1,+1]×[-1,+1]×[-1,+1]
//  Output:
//    x -- the "real" coordinates {x,y,z}
func (o *Transfinite) Point(x, u la.Vector) {

	// 2D
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

	// 3D
	r, s, t := u[0], u[1], u[2]
	m, p := -1.0, +1.0
	o.b3d[0](o.b0st, s, t)
	o.b3d[1](o.b1st, s, t)
	o.b3d[2](o.b2rt, r, t)
	o.b3d[3](o.b3rt, r, t)
	o.b3d[4](o.b4rs, r, s)
	o.b3d[5](o.b5rs, r, s)

	o.b3d[0](o.b0mt, m, t)
	o.b3d[0](o.b0pt, p, t)
	o.b3d[1](o.b1mt, m, t)
	o.b3d[1](o.b1pt, p, t)

	o.b3d[0](o.b0sm, s, m)
	o.b3d[0](o.b0sp, s, p)
	o.b3d[1](o.b1sm, s, m)
	o.b3d[1](o.b1sp, s, p)

	o.b3d[2](o.b2rm, r, m)
	o.b3d[2](o.b2rp, r, p)
	o.b3d[3](o.b3rm, r, m)
	o.b3d[3](o.b3rp, r, p)
	for i := 0; i < o.ndim; i++ {
		x[i] = 0 + // trick to force alignment

			+(1.0-r)*o.b0st[i]/2.0 + (1.0+r)*o.b1st[i]/2.0 +
			+(1.0-s)*o.b2rt[i]/2.0 + (1.0+s)*o.b3rt[i]/2.0 +
			+(1.0-t)*o.b4rs[i]/2.0 + (1.0+t)*o.b5rs[i]/2.0 +

			-(1.0-r)*(1.0-s)*o.b0mt[i]/4.0 - (1.0-r)*(1.0+s)*o.b0pt[i]/4.0 +
			-(1.0+r)*(1.0-s)*o.b1mt[i]/4.0 - (1.0+r)*(1.0+s)*o.b1pt[i]/4.0 +

			-(1.0-r)*(1.0-t)*o.b0sm[i]/4.0 - (1.0-r)*(1.0+t)*o.b0sp[i]/4.0 +
			-(1.0+r)*(1.0-t)*o.b1sm[i]/4.0 - (1.0+r)*(1.0+t)*o.b1sp[i]/4.0 +

			-(1.0-s)*(1.0-t)*o.b2rm[i]/4.0 - (1.0-s)*(1.0+t)*o.b2rp[i]/4.0 +
			-(1.0+s)*(1.0-t)*o.b3rm[i]/4.0 - (1.0+s)*(1.0+t)*o.b3rp[i]/4.0 +

			+(1.0-r)*(1.0-s)*(1.0-t)*o.p0[i]/8.0 + (1.0+r)*(1.0-s)*(1.0-t)*o.p1[i]/8.0 +
			+(1.0+r)*(1.0+s)*(1.0-t)*o.p2[i]/8.0 + (1.0-r)*(1.0+s)*(1.0-t)*o.p3[i]/8.0 +
			+(1.0-r)*(1.0-s)*(1.0+t)*o.p4[i]/8.0 + (1.0+r)*(1.0-s)*(1.0+t)*o.p5[i]/8.0 +
			+(1.0+r)*(1.0+s)*(1.0+t)*o.p6[i]/8.0 + (1.0-r)*(1.0+s)*(1.0+t)*o.p7[i]/8.0
	}
}

// Derivs calculates derivatives (=metric terms) @ u={r,s,t}
//  Input:
//    u -- the "reference" coordinates {r,s,t} ϵ [-1,+1]×[-1,+1]×[-1,+1]
//  Output:
//    dxdu -- the derivatives [dx/du]ij = dxi/duj
func (o *Transfinite) Derivs(dxdu *la.Matrix, u la.Vector) {

	// 2D
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

	// 3D
	r, s, t := u[0], u[1], u[2]
	m, p := -1.0, +1.0
	o.b3d[0](o.b0st, s, t)
	o.b3d[1](o.b1st, s, t)
	o.b3d[2](o.b2rt, r, t)
	o.b3d[3](o.b3rt, r, t)
	o.b3d[4](o.b4rs, r, s)
	o.b3d[5](o.b5rs, r, s)

	o.b3d[0](o.b0mt, m, t)
	o.b3d[0](o.b0pt, p, t)
	o.b3d[1](o.b1mt, m, t)
	o.b3d[1](o.b1pt, p, t)

	o.b3d[0](o.b0sm, s, m)
	o.b3d[0](o.b0sp, s, p)
	o.b3d[1](o.b1sm, s, m)
	o.b3d[1](o.b1sp, s, p)

	o.b3d[2](o.b2rm, r, m)
	o.b3d[2](o.b2rp, r, p)
	o.b3d[3](o.b3rm, r, m)
	o.b3d[3](o.b3rp, r, p)

	o.bd3d[2](o.bdr2rt, r, t)
	o.bd3d[3](o.bdr3rt, r, t)
	o.bd3d[4](o.bdr4rs, r, s)
	o.bd3d[5](o.bdr5rs, r, s)
	o.bd3d[2](o.bdr2rm, r, m)
	o.bd3d[2](o.bdr2rp, r, p)
	o.bd3d[3](o.bdr3rm, r, m)
	o.bd3d[3](o.bdr3rp, r, p)

	o.bd3d[0](o.bd0st, s, t)
	o.bd3d[1](o.bd1st, s, t)
	o.bd3d[4](o.bd4rs, r, s)
	o.bd3d[5](o.bd5rs, r, s)
	o.bd3d[0](o.bd0sm, s, m)
	o.bd3d[0](o.bd0sp, s, p)
	o.bd3d[1](o.bd1sm, s, m)
	o.bd3d[1](o.bd1sp, s, p)

	o.bd3d[2](o.bd2rt, r, t)
	o.bd3d[3](o.bd3rt, r, t)
	o.bd3d[0](o.bd0mt, m, t)
	o.bd3d[0](o.bd0pt, p, t)
	o.bd3d[1](o.bd1mt, m, t)
	o.bd3d[1](o.bd1pt, p, t)

	var dxidr, dxids, dxidt float64
	for i := 0; i < o.ndim; i++ {

		dxidr = 0 + // trick to force alignment
			-o.b0st[i]/2.0 + o.b1st[i]/2.0 +
			+(1.0-s)*o.bdr2rt.Get(i, 0)/2.0 + (1.0+s)*o.bdr3rt.Get(i, 0)/2.0 +
			+(1.0-t)*o.bdr4rs.Get(i, 0)/2.0 + (1.0+t)*o.bdr5rs.Get(i, 0)/2.0 +

			+(1.0-s)*o.b0mt[i]/4.0 + (1.0+s)*o.b0pt[i]/4.0 +
			-(1.0-s)*o.b1mt[i]/4.0 - (1.0+s)*o.b1pt[i]/4.0 +

			+(1.0-t)*o.b0sm[i]/4.0 + (1.0+t)*o.b0sp[i]/4.0 +
			-(1.0-t)*o.b1sm[i]/4.0 - (1.0+t)*o.b1sp[i]/4.0 +

			-(1.0-s)*(1.0-t)*o.bdr2rm.Get(i, 0)/4.0 - (1.0-s)*(1.0+t)*o.bdr2rp.Get(i, 0)/4.0 +
			-(1.0+s)*(1.0-t)*o.bdr3rm.Get(i, 0)/4.0 - (1.0+s)*(1.0+t)*o.bdr3rp.Get(i, 0)/4.0 +

			-(1.0-s)*(1.0-t)*o.p0[i]/8.0 + (1.0-s)*(1.0-t)*o.p1[i]/8.0 +
			+(1.0+s)*(1.0-t)*o.p2[i]/8.0 - (1.0+s)*(1.0-t)*o.p3[i]/8.0 +
			-(1.0-s)*(1.0+t)*o.p4[i]/8.0 + (1.0-s)*(1.0+t)*o.p5[i]/8.0 +
			+(1.0+s)*(1.0+t)*o.p6[i]/8.0 - (1.0+s)*(1.0+t)*o.p7[i]/8.0

		dxids = 0 + // trick to force alignment
			+(1.0-r)*o.bd0st.Get(i, 0)/2.0 + (1.0+r)*o.bd1st.Get(i, 0)/2.0 +
			-o.b2rt[i]/2.0 + o.b3rt[i]/2.0 +
			+(1.0-t)*o.bd4rs.Get(i, 1)/2.0 + (1.0+t)*o.bd5rs.Get(i, 1)/2.0 +

			+(1.0-r)*o.b0mt[i]/4.0 - (1.0-r)*o.b0pt[i]/4.0 +
			+(1.0+r)*o.b1mt[i]/4.0 - (1.0+r)*o.b1pt[i]/4.0 +

			-(1.0-r)*(1.0-t)*o.bd0sm.Get(i, 0)/4.0 - (1.0-r)*(1.0+t)*o.bd0sp.Get(i, 0)/4.0 +
			-(1.0+r)*(1.0-t)*o.bd1sm.Get(i, 0)/4.0 - (1.0+r)*(1.0+t)*o.bd1sp.Get(i, 0)/4.0 +

			+(1.0-t)*o.b2rm[i]/4.0 + (1.0+t)*o.b2rp[i]/4.0 +
			-(1.0-t)*o.b3rm[i]/4.0 - (1.0+t)*o.b3rp[i]/4.0 +

			-(1.0-r)*(1.0-t)*o.p0[i]/8.0 - (1.0+r)*(1.0-t)*o.p1[i]/8.0 +
			+(1.0+r)*(1.0-t)*o.p2[i]/8.0 + (1.0-r)*(1.0-t)*o.p3[i]/8.0 +
			-(1.0-r)*(1.0+t)*o.p4[i]/8.0 - (1.0+r)*(1.0+t)*o.p5[i]/8.0 +
			+(1.0+r)*(1.0+t)*o.p6[i]/8.0 + (1.0-r)*(1.0+t)*o.p7[i]/8.0

		dxidt = 0 + // trick to force alignment
			+(1.0-r)*o.bd0st.Get(i, 1)/2.0 + (1.0+r)*o.bd1st.Get(i, 1)/2.0 +
			+(1.0-s)*o.bd2rt.Get(i, 1)/2.0 + (1.0+s)*o.bd3rt.Get(i, 1)/2.0 +
			-o.b4rs[i]/2.0 + o.b5rs[i]/2.0 +

			-(1.0-r)*(1.0-s)*o.bd0mt.Get(i, 1)/4.0 - (1.0-r)*(1.0+s)*o.bd0pt.Get(i, 1)/4.0 +
			-(1.0+r)*(1.0-s)*o.bd1mt.Get(i, 1)/4.0 - (1.0+r)*(1.0+s)*o.bd1pt.Get(i, 1)/4.0 +

			+(1.0-r)*o.b0sm[i]/4.0 - (1.0-r)*o.b0sp[i]/4.0 +
			+(1.0+r)*o.b1sm[i]/4.0 - (1.0+r)*o.b1sp[i]/4.0 +

			+(1.0-s)*o.b2rm[i]/4.0 - (1.0-s)*o.b2rp[i]/4.0 +
			+(1.0+s)*o.b3rm[i]/4.0 - (1.0+s)*o.b3rp[i]/4.0 +

			-(1.0-r)*(1.0-s)*o.p0[i]/8.0 - (1.0+r)*(1.0-s)*o.p1[i]/8.0 +
			-(1.0+r)*(1.0+s)*o.p2[i]/8.0 - (1.0-r)*(1.0+s)*o.p3[i]/8.0 +
			+(1.0-r)*(1.0-s)*o.p4[i]/8.0 + (1.0+r)*(1.0-s)*o.p5[i]/8.0 +
			+(1.0+r)*(1.0+s)*o.p6[i]/8.0 + (1.0-r)*(1.0+s)*o.p7[i]/8.0

		dxdu.Set(i, 0, dxidr)
		dxdu.Set(i, 1, dxids)
		dxdu.Set(i, 2, dxidt)
	}
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

	// 2D
	if o.ndim == 2 {
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
		return
	}

	// 3D
	z0 := make([]float64, npts[0])
	z1 := make([]float64, npts[1])
	x2 := make([]float64, npts[2])
	y2 := make([]float64, npts[2])
	z2 := make([]float64, npts[2])
	if !onlyBry {

		args.Ms = 3
		args.Mec = args.C

		// draw 0-lines
		for k := 0; k < npts[2]; k++ {
			u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
			for j := 0; j < npts[1]; j++ {
				u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
				for i := 0; i < npts[0]; i++ {
					u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
					o.Point(x, u)
					x0[i] = x[0]
					y0[i] = x[1]
					z0[i] = x[2]
				}
				//plt.Plot3dLine(x0, y0, z0, args)
				plt.Plot3dPoints(x0, y0, z0, args)
			}
		}

		// draw 1-lines
		for k := 0; k < npts[2]; k++ {
			u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
			for i := 0; i < npts[0]; i++ {
				u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
				for j := 0; j < npts[1]; j++ {
					u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
					o.Point(x, u)
					x1[j] = x[0]
					y1[j] = x[1]
					z1[j] = x[2]
				}
				//plt.Plot3dLine(x1, y1, z1, args)
				plt.Plot3dPoints(x1, y1, z1, args)
			}
		}

		// draw 2-lines
		for j := 0; j < npts[1]; j++ {
			u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
			for i := 0; i < npts[0]; i++ {
				u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
				for k := 0; k < npts[2]; k++ {
					u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
					o.Point(x, u)
					x2[k] = x[0]
					y2[k] = x[1]
					z2[k] = x[2]
				}
				//plt.Plot3dLine(x2, y2, z2, args)
				plt.Plot3dPoints(x2, y2, z2, args)
			}
		}
	}

	// draw B[0](s,t)
	u[0] = -1
	for k := 0; k < npts[2]; k++ {
		u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
		for j := 0; j < npts[1]; j++ {
			u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
			o.Point(x, u)
			x1[j] = x[0]
			y1[j] = x[1]
			z1[j] = x[2]
		}
		plt.Plot3dLine(x1, y1, z1, args)
	}
	for j := 0; j < npts[1]; j++ {
		u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
		for k := 0; k < npts[2]; k++ {
			u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
			o.Point(x, u)
			x2[k] = x[0]
			y2[k] = x[1]
			z2[k] = x[2]
		}
		plt.Plot3dLine(x2, y2, z2, args)
	}

	// draw B[1](s,t)
	u[0] = +1
	for k := 0; k < npts[2]; k++ {
		u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
		for j := 0; j < npts[1]; j++ {
			u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
			o.Point(x, u)
			x1[j] = x[0]
			y1[j] = x[1]
			z1[j] = x[2]
		}
		plt.Plot3dLine(x1, y1, z1, args)
	}
	for j := 0; j < npts[1]; j++ {
		u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
		for k := 0; k < npts[2]; k++ {
			u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
			o.Point(x, u)
			x2[k] = x[0]
			y2[k] = x[1]
			z2[k] = x[2]
		}
		plt.Plot3dLine(x2, y2, z2, args)
	}

	// draw B[2](r,t)
	u[1] = -1
	for k := 0; k < npts[2]; k++ {
		u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
		for i := 0; i < npts[0]; i++ {
			u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
			o.Point(x, u)
			x0[i] = x[0]
			y0[i] = x[1]
			z0[i] = x[2]
		}
		plt.Plot3dLine(x0, y0, z0, args)
	}
	for i := 0; i < npts[0]; i++ {
		u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
		for k := 0; k < npts[2]; k++ {
			u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
			o.Point(x, u)
			x2[k] = x[0]
			y2[k] = x[1]
			z2[k] = x[2]
		}
		plt.Plot3dLine(x2, y2, z2, args)
	}

	// draw B[3](r,t)
	u[1] = +1
	for k := 0; k < npts[2]; k++ {
		u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
		for i := 0; i < npts[0]; i++ {
			u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
			o.Point(x, u)
			x0[i] = x[0]
			y0[i] = x[1]
			z0[i] = x[2]
		}
		plt.Plot3dLine(x0, y0, z0, args)
	}
	for i := 0; i < npts[0]; i++ {
		u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
		for k := 0; k < npts[2]; k++ {
			u[2] = -1 + 2*float64(k)/float64(npts[2]-1)
			o.Point(x, u)
			x2[k] = x[0]
			y2[k] = x[1]
			z2[k] = x[2]
		}
		plt.Plot3dLine(x2, y2, z2, args)
	}

	// draw B[4](r,s)
	u[2] = -1
	for j := 0; j < npts[1]; j++ {
		u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
		for i := 0; i < npts[0]; i++ {
			u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
			o.Point(x, u)
			x0[i] = x[0]
			y0[i] = x[1]
			z0[i] = x[2]
		}
		plt.Plot3dLine(x0, y0, z0, args)
	}
	for i := 0; i < npts[0]; i++ {
		u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
		for j := 0; j < npts[1]; j++ {
			u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
			o.Point(x, u)
			x1[j] = x[0]
			y1[j] = x[1]
			z1[j] = x[2]
		}
		plt.Plot3dLine(x1, y1, z1, args)
	}

	// draw B[5](r,s)
	u[2] = +1
	for j := 0; j < npts[1]; j++ {
		u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
		for i := 0; i < npts[0]; i++ {
			u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
			o.Point(x, u)
			x0[i] = x[0]
			y0[i] = x[1]
			z0[i] = x[2]
		}
		plt.Plot3dLine(x0, y0, z0, args)
	}
	for i := 0; i < npts[0]; i++ {
		u[0] = -1 + 2*float64(i)/float64(npts[0]-1)
		for j := 0; j < npts[1]; j++ {
			u[1] = -1 + 2*float64(j)/float64(npts[1]-1)
			o.Point(x, u)
			x1[j] = x[0]
			y1[j] = x[1]
			z1[j] = x[2]
		}
		plt.Plot3dLine(x1, y1, z1, args)
	}
}
