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
//
//                                    +----------------+
//                                  ,'|              ,'|
//       t or z                   ,'  |  ___       ,'  |     B[0](s,t)
//          ↑                   ,'    |,'5,'  [0],'    |     B[1](s,t)
//          |                 ,'      |~~~     ,'      |     B[2](r,t)
//          |               +'===============+'  ,'|   |     B[3](r,t)
//          |               |   ,'|   |      |   |3|   |     B[4](r,s)
//          |     s or y    |   |2|   |      |   |,'   |     B[5](r,s)
//          +-------->      |   |,'   +- - - | +- - - -+
//        ,'                |       ,'       |       ,'
//      ,'                  |     ,' [1]  ___|     ,'
//  r or x                  |   ,'      ,'4,'|   ,'
//                          | ,'        ~~~  | ,'
//                          +----------------+'
type Transfinite struct {

	// input
	ndim int // space dimension

	// input data for 2d
	e   []fun.Vs // [4] 2D boundary functions
	ed  []fun.Vs // [4] 2D 1st derivatives of boundary functions
	edd []fun.Vs // [4] 2D 2nd derivatives of boundary functions

	// input data for 3d
	b   []fun.Vss   // [6] 3D boundary function
	bd  []fun.Vvss  // [6] 3D 1st derivatives of boundary functions
	bdd []fun.Vvvss // [6] 3D 2nd derivatives of boundary functions

	// workspace for 2d
	p0, p1   la.Vector // corner points
	p2, p3   la.Vector // corner points
	e0s, e1s la.Vector // 2D function evaluations
	e2r, e3r la.Vector // 2D function evaluations

	de0sDs, de1sDs     la.Vector // derivative evaluation
	de2rDr, de3rDr     la.Vector // derivative evaluation
	dde0sDss, dde1sDss la.Vector // derivative evaluation
	dde2rDrr, dde3rDrr la.Vector // derivative evaluation

	// workspace for 3d
	p4, p5   la.Vector // corner points
	p6, p7   la.Vector // corner points
	tm1, tm2 la.Vector // temporary vectors

	b0st, b1st, b2rt la.Vector // function evaluation
	b3rt, b4rs, b5rs la.Vector // function evaluation

	b0mt, b0pt, b1mt, b1pt la.Vector // function evaluation
	b0sm, b0sp, b1sm, b1sp la.Vector // function evaluation
	b2rm, b2rp, b3rm, b3rp la.Vector // function evaluation

	db0stDs, db0stDt la.Vector // derivative evaluation
	db1stDs, db1stDt la.Vector // derivative evaluation
	db2rtDr, db2rtDt la.Vector // derivative evaluation
	db3rtDr, db3rtDt la.Vector // derivative evaluation
	db4rsDr, db4rsDs la.Vector // derivative evaluation
	db5rsDr, db5rsDs la.Vector // derivative evaluation

	db0smDs, db0spDs la.Vector // derivative evaluation
	db0mtDt, db0ptDt la.Vector // derivative evaluation
	db1smDs, db1spDs la.Vector // derivative evaluation
	db1mtDt, db1ptDt la.Vector // derivative evaluation
	db2rmDr, db2rpDr la.Vector // derivative evaluation
	db3rmDr, db3rpDr la.Vector // derivative evaluation

	ddb0stDss, ddb0stDtt, ddb0stDst la.Vector // derivative evaluation
	ddb1stDss, ddb1stDtt, ddb1stDst la.Vector // derivative evaluation
	ddb2rtDrr, ddb2rtDtt, ddb2rtDrt la.Vector // derivative evaluation
	ddb3rtDrr, ddb3rtDtt, ddb3rtDrt la.Vector // derivative evaluation
	ddb4rsDrr, ddb4rsDss, ddb4rsDrs la.Vector // derivative evaluation
	ddb5rsDrr, ddb5rsDss, ddb5rsDrs la.Vector // derivative evaluation

	ddb0smDss, ddb0spDss la.Vector // derivative evaluation
	ddb0mtDtt, ddb0ptDtt la.Vector // derivative evaluation
	ddb1smDss, ddb1spDss la.Vector // derivative evaluation
	ddb1mtDtt, ddb1ptDtt la.Vector // derivative evaluation
	ddb2rmDrr, ddb2rpDrr la.Vector // derivative evaluation
	ddb3rmDrr, ddb3rpDrr la.Vector // derivative evaluation
}

// NewTransfinite2d allocates a new structure
//   Input:
//     B   -- [4] boundary functions
//     Bd  -- [4] 1st derivative of boundary functions
//     Bdd -- [4 or nil] 2nd derivative of boundary functions [may be nil]
func NewTransfinite2d(B, Bd, Bdd []fun.Vs) (o *Transfinite) {

	// check
	if len(B) != 4 || len(Bd) != 4 {
		chk.Panic("in 2D, four boundary functions B are required\n")
	}

	// input
	o = new(Transfinite)
	o.ndim = 2
	o.e = B
	o.ed = Bd
	o.edd = Bdd

	// allocate workspace
	o.p0, o.p1 = la.NewVector(2), la.NewVector(2)
	o.p2, o.p3 = la.NewVector(2), la.NewVector(2)
	o.e0s, o.e1s = la.NewVector(2), la.NewVector(2)
	o.e2r, o.e3r = la.NewVector(2), la.NewVector(2)
	o.de0sDs, o.de1sDs = la.NewVector(2), la.NewVector(2)
	o.de2rDr, o.de3rDr = la.NewVector(2), la.NewVector(2)

	if o.edd != nil {
		if len(o.edd) != 4 {
			chk.Panic("in 2D, four boundary functions B are required. len(Bdd) = %d is incorrect\n", len(o.edd))
		}
		o.dde0sDss, o.dde1sDss = la.NewVector(2), la.NewVector(2)
		o.dde2rDrr, o.dde3rDrr = la.NewVector(2), la.NewVector(2)
	}

	// compute corners
	o.e[0](o.p0, -1)
	o.e[0](o.p3, +1)
	o.e[1](o.p1, -1)
	o.e[1](o.p2, +1)
	return
}

// NewTransfinite3d allocates a new structure
//   Input:
//     B   -- [6] boundary functions
//     Bd  -- [6] 1st derivative of boundary functions
//     Bdd -- [6 or nil] 2nd derivative of boundary functions [may be nil]
func NewTransfinite3d(B []fun.Vss, Bd []fun.Vvss, Bdd []fun.Vvvss) (o *Transfinite) {

	// check
	if len(B) != 6 || len(Bd) != 6 {
		chk.Panic("in 3D, six boundary functions B are required\n")
	}

	// input
	o = new(Transfinite)
	o.ndim = 3
	o.b = B
	o.bd = Bd
	o.bdd = Bdd

	// allocate workspace
	o.p0, o.p1 = la.NewVector(3), la.NewVector(3)
	o.p2, o.p3 = la.NewVector(3), la.NewVector(3)
	o.p4, o.p5 = la.NewVector(3), la.NewVector(3)
	o.p6, o.p7 = la.NewVector(3), la.NewVector(3)
	o.tm1, o.tm2 = la.NewVector(3), la.NewVector(3)

	o.b0st, o.b1st, o.b2rt = la.NewVector(3), la.NewVector(3), la.NewVector(3)
	o.b3rt, o.b4rs, o.b5rs = la.NewVector(3), la.NewVector(3), la.NewVector(3)

	o.b0mt, o.b0pt, o.b1mt, o.b1pt = la.NewVector(3), la.NewVector(3), la.NewVector(3), la.NewVector(3)
	o.b0sm, o.b0sp, o.b1sm, o.b1sp = la.NewVector(3), la.NewVector(3), la.NewVector(3), la.NewVector(3)
	o.b2rm, o.b2rp, o.b3rm, o.b3rp = la.NewVector(3), la.NewVector(3), la.NewVector(3), la.NewVector(3)

	o.db0stDs, o.db0stDt = la.NewVector(3), la.NewVector(3)
	o.db1stDs, o.db1stDt = la.NewVector(3), la.NewVector(3)
	o.db2rtDr, o.db2rtDt = la.NewVector(3), la.NewVector(3)
	o.db3rtDr, o.db3rtDt = la.NewVector(3), la.NewVector(3)
	o.db4rsDr, o.db4rsDs = la.NewVector(3), la.NewVector(3)
	o.db5rsDr, o.db5rsDs = la.NewVector(3), la.NewVector(3)

	o.db0smDs, o.db0spDs = la.NewVector(3), la.NewVector(3)
	o.db0mtDt, o.db0ptDt = la.NewVector(3), la.NewVector(3)
	o.db1smDs, o.db1spDs = la.NewVector(3), la.NewVector(3)
	o.db1mtDt, o.db1ptDt = la.NewVector(3), la.NewVector(3)
	o.db2rmDr, o.db2rpDr = la.NewVector(3), la.NewVector(3)
	o.db3rmDr, o.db3rpDr = la.NewVector(3), la.NewVector(3)

	if o.bdd != nil {
		if len(o.bdd) != 6 {
			chk.Panic("in 3D, six boundary functions B are required. len(Bdd) = %d is incorrect\n", len(o.bdd))
		}
		o.ddb0stDss, o.ddb0stDtt, o.ddb0stDst = la.NewVector(3), la.NewVector(3), la.NewVector(3)
		o.ddb1stDss, o.ddb1stDtt, o.ddb1stDst = la.NewVector(3), la.NewVector(3), la.NewVector(3)
		o.ddb2rtDrr, o.ddb2rtDtt, o.ddb2rtDrt = la.NewVector(3), la.NewVector(3), la.NewVector(3)
		o.ddb3rtDrr, o.ddb3rtDtt, o.ddb3rtDrt = la.NewVector(3), la.NewVector(3), la.NewVector(3)
		o.ddb4rsDrr, o.ddb4rsDss, o.ddb4rsDrs = la.NewVector(3), la.NewVector(3), la.NewVector(3)
		o.ddb5rsDrr, o.ddb5rsDss, o.ddb5rsDrs = la.NewVector(3), la.NewVector(3), la.NewVector(3)

		o.ddb0smDss, o.ddb0spDss = la.NewVector(3), la.NewVector(3)
		o.ddb0mtDtt, o.ddb0ptDtt = la.NewVector(3), la.NewVector(3)
		o.ddb1smDss, o.ddb1spDss = la.NewVector(3), la.NewVector(3)
		o.ddb1mtDtt, o.ddb1ptDtt = la.NewVector(3), la.NewVector(3)
		o.ddb2rmDrr, o.ddb2rpDrr = la.NewVector(3), la.NewVector(3)
		o.ddb3rmDrr, o.ddb3rpDrr = la.NewVector(3), la.NewVector(3)
	}

	// compute corners
	o.b[4](o.p0, -1.0, -1.0)
	o.b[4](o.p1, +1.0, -1.0)
	o.b[4](o.p2, +1.0, +1.0)
	o.b[4](o.p3, -1.0, +1.0)
	o.b[5](o.p4, -1.0, -1.0)
	o.b[5](o.p5, +1.0, -1.0)
	o.b[5](o.p6, +1.0, +1.0)
	o.b[5](o.p7, -1.0, +1.0)
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

		// compute boundary functions @ {r,s}
		r, s := u[0], u[1]
		o.e[0](o.e0s, s)
		o.e[1](o.e1s, s)
		o.e[2](o.e2r, r)
		o.e[3](o.e3r, r)

		// compute position
		for i := 0; i < o.ndim; i++ {
			x[i] = 0 +
				+(1.0-r)*o.e0s[i]/2.0 + (1.0+r)*o.e1s[i]/2.0 +
				+(1.0-s)*o.e2r[i]/2.0 + (1.0+s)*o.e3r[i]/2.0 +
				-(1.0-r)*(1.0-s)*o.p0[i]/4.0 - (1.0+r)*(1.0-s)*o.p1[i]/4.0 +
				-(1.0+r)*(1.0+s)*o.p2[i]/4.0 - (1.0-r)*(1.0+s)*o.p3[i]/4.0
		}
		return
	}

	// 3D
	r, s, t := u[0], u[1], u[2]
	m, p := -1.0, +1.0

	// compute boundary functions @ {r,s,t}
	o.b[0](o.b0st, s, t)
	o.b[1](o.b1st, s, t)
	o.b[2](o.b2rt, r, t)
	o.b[3](o.b3rt, r, t)
	o.b[4](o.b4rs, r, s)
	o.b[5](o.b5rs, r, s)

	// compute boundary functions @ edges
	o.b[0](o.b0mt, m, t)
	o.b[0](o.b0pt, p, t)
	o.b[1](o.b1mt, m, t)
	o.b[1](o.b1pt, p, t)

	o.b[0](o.b0sm, s, m)
	o.b[0](o.b0sp, s, p)
	o.b[1](o.b1sm, s, m)
	o.b[1](o.b1sp, s, p)

	o.b[2](o.b2rm, r, m)
	o.b[2](o.b2rp, r, p)
	o.b[3](o.b3rm, r, m)
	o.b[3](o.b3rp, r, p)

	// compute position
	for i := 0; i < o.ndim; i++ {
		x[i] = 0 +
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

// PointAndDerivs computes position and the first and second order derivatives
//   Input:
//     u -- reference coordinates {r,s,t}
//   Output:
//     x      -- position {x,y,z}
//     dxDr   -- ∂{x}/∂r
//     dxDs   -- ∂{x}/∂s
//     dxDt   -- ∂{x}/∂t
//     ddxDrr -- ∂²{x}/∂r²  [may be nil]
//     ddxDss -- ∂²{x}/∂s²  [may be nil]
//     ddxDtt -- ∂²{x}/∂t²  [may be nil]
//     ddxDrs -- ∂²{x}/∂r∂s [may be nil]
//     ddxDrt -- ∂²{x}/∂r∂t [may be nil]
//     ddxDst -- ∂²{x}/∂s∂t [may be nil]
func (o *Transfinite) PointAndDerivs(x, dxDr, dxDs, dxDt,
	ddxDrr, ddxDss, ddxDtt, ddxDrs, ddxDrt, ddxDst, u la.Vector) {

	// auxiliary
	secondDerivs := ddxDrr != nil
	m, p := -1.0, +1.0

	// 2D
	if o.ndim == 2 {

		// auxiliary
		r, s := u[0], u[1]

		// compute boundary functions @ {r,s}
		o.e[0](o.e0s, s)
		o.e[1](o.e1s, s)
		o.e[2](o.e2r, r)
		o.e[3](o.e3r, r)

		// compute derivatives @ {r,s}
		o.ed[0](o.de0sDs, s)
		o.ed[1](o.de1sDs, s)
		o.ed[2](o.de2rDr, r)
		o.ed[3](o.de3rDr, r)

		// position and first order derivatives
		for i := 0; i < o.ndim; i++ {

			// bilinear transfinite mapping in 2D
			x[i] = 0 +
				+(1.0-r)*o.e0s[i]/2.0 + (1.0+r)*o.e1s[i]/2.0 +
				+(1.0-s)*o.e2r[i]/2.0 + (1.0+s)*o.e3r[i]/2.0 +
				-(1.0-r)*(1.0-s)*o.p0[i]/4.0 - (1.0+r)*(1.0-s)*o.p1[i]/4.0 +
				-(1.0+r)*(1.0+s)*o.p2[i]/4.0 - (1.0-r)*(1.0+s)*o.p3[i]/4.0

			// derivative of x with respect to r
			dxDr[i] = 0 +
				-o.e0s[i]/2.0 + o.e1s[i]/2.0 +
				+(1.0-s)*o.de2rDr[i]/2.0 + (1.0+s)*o.de3rDr[i]/2.0 +
				+(1.0-s)*o.p0[i]/4.0 - (1.0-s)*o.p1[i]/4.0 +
				-(1.0+s)*o.p2[i]/4.0 + (1.0+s)*o.p3[i]/4.0

			// derivative of x with respect to s
			dxDs[i] = 0 +
				+(1.0-r)*o.de0sDs[i]/2.0 + (1.0+r)*o.de1sDs[i]/2.0 +
				-o.e2r[i]/2.0 + o.e3r[i]/2.0 +
				+(1.0-r)*o.p0[i]/4.0 + (1.0+r)*o.p1[i]/4.0 +
				-(1.0+r)*o.p2[i]/4.0 - (1.0-r)*o.p3[i]/4.0
		}

		// skip second order derivatives
		if !secondDerivs {
			return
		}

		// only 2nd cross-derivatives may be non-zero
		if o.edd == nil {
			for i := 0; i < o.ndim; i++ {
				ddxDrr[i] = 0.0
				ddxDss[i] = 0.0
				ddxDrs[i] = 0 +
					-o.de0sDs[i]/2.0 + o.de1sDs[i]/2.0 +
					-o.de2rDr[i]/2.0 + o.de3rDr[i]/2.0 +
					-o.p0[i]/4.0 + o.p1[i]/4.0 +
					-o.p2[i]/4.0 + o.p3[i]/4.0
			}
			return
		}

		// compute second derivatives @ {r,s,t}
		o.edd[0](o.dde0sDss, s)
		o.edd[1](o.dde1sDss, s)
		o.edd[2](o.dde2rDrr, r)
		o.edd[3](o.dde3rDrr, r)

		// second order derivatives
		for i := 0; i < o.ndim; i++ {

			// derivative of dx/dr with respect to r
			ddxDrr[i] = 0 +
				(1.0-s)*o.dde2rDrr[i]/2.0 + (1.0+s)*o.dde3rDrr[i]/2.0

			// derivative of dx/ds with respect to s
			ddxDss[i] = 0 +
				(1.0-r)*o.dde0sDss[i]/2.0 + (1.0+r)*o.dde1sDss[i]/2.0

			// derivative of dx/dr with respect to s
			ddxDrs[i] = 0 +
				-o.de0sDs[i]/2.0 + o.de1sDs[i]/2.0 +
				-o.de2rDr[i]/2.0 + o.de3rDr[i]/2.0 +
				-o.p0[i]/4.0 + o.p1[i]/4.0 +
				-o.p2[i]/4.0 + o.p3[i]/4.0
		}
		return
	}

	// auxiliary
	r, s, t := u[0], u[1], u[2]

	// compute boundary functions @ {r,s,t}
	o.b[0](o.b0st, s, t)
	o.b[1](o.b1st, s, t)
	o.b[2](o.b2rt, r, t)
	o.b[3](o.b3rt, r, t)
	o.b[4](o.b4rs, r, s)
	o.b[5](o.b5rs, r, s)

	// compute boundary functions @ edges
	o.b[0](o.b0mt, m, t)
	o.b[0](o.b0pt, p, t)
	o.b[1](o.b1mt, m, t)
	o.b[1](o.b1pt, p, t)

	o.b[0](o.b0sm, s, m)
	o.b[0](o.b0sp, s, p)
	o.b[1](o.b1sm, s, m)
	o.b[1](o.b1sp, s, p)

	o.b[2](o.b2rm, r, m)
	o.b[2](o.b2rp, r, p)
	o.b[3](o.b3rm, r, m)
	o.b[3](o.b3rp, r, p)

	// compute derivatives @ {r,s,t}
	o.bd[0](o.db0stDs, o.db0stDt, s, t)
	o.bd[1](o.db1stDs, o.db1stDt, s, t)
	o.bd[2](o.db2rtDr, o.db2rtDt, r, t)
	o.bd[3](o.db3rtDr, o.db3rtDt, r, t)
	o.bd[4](o.db4rsDr, o.db4rsDs, r, s)
	o.bd[5](o.db5rsDr, o.db5rsDs, r, s)

	// compute derivatives @ edges
	o.bd[0](o.db0smDs, o.tm1, s, m)
	o.bd[0](o.db0spDs, o.tm1, s, p)
	o.bd[0](o.tm1, o.db0mtDt, m, t)
	o.bd[0](o.tm1, o.db0ptDt, p, t)

	o.bd[1](o.db1smDs, o.tm1, s, m)
	o.bd[1](o.db1spDs, o.tm1, s, p)
	o.bd[1](o.tm1, o.db1mtDt, m, t)
	o.bd[1](o.tm1, o.db1ptDt, p, t)

	o.bd[2](o.db2rmDr, o.tm1, r, m)
	o.bd[2](o.db2rpDr, o.tm1, r, p)
	o.bd[3](o.db3rmDr, o.tm1, r, m)
	o.bd[3](o.db3rpDr, o.tm1, r, p)

	// position and first order derivatives
	for i := 0; i < o.ndim; i++ {

		// bilinear transfinite mapping in 3D
		x[i] = 0 +
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

		// derivative of x with respect to r
		dxDr[i] = 0 +
			-o.b0st[i]/2.0 + o.b1st[i]/2.0 +
			+(1.0-s)*o.db2rtDr[i]/2.0 + (1.0+s)*o.db3rtDr[i]/2.0 +
			+(1.0-t)*o.db4rsDr[i]/2.0 + (1.0+t)*o.db5rsDr[i]/2.0 +

			+(1.0-s)*o.b0mt[i]/4.0 + (1.0+s)*o.b0pt[i]/4.0 +
			-(1.0-s)*o.b1mt[i]/4.0 - (1.0+s)*o.b1pt[i]/4.0 +

			+(1.0-t)*o.b0sm[i]/4.0 + (1.0+t)*o.b0sp[i]/4.0 +
			-(1.0-t)*o.b1sm[i]/4.0 - (1.0+t)*o.b1sp[i]/4.0 +

			-(1.0-s)*(1.0-t)*o.db2rmDr[i]/4.0 - (1.0-s)*(1.0+t)*o.db2rpDr[i]/4.0 +
			-(1.0+s)*(1.0-t)*o.db3rmDr[i]/4.0 - (1.0+s)*(1.0+t)*o.db3rpDr[i]/4.0 +

			-(1.0-s)*(1.0-t)*o.p0[i]/8.0 + (1.0-s)*(1.0-t)*o.p1[i]/8.0 +
			+(1.0+s)*(1.0-t)*o.p2[i]/8.0 - (1.0+s)*(1.0-t)*o.p3[i]/8.0 +
			-(1.0-s)*(1.0+t)*o.p4[i]/8.0 + (1.0-s)*(1.0+t)*o.p5[i]/8.0 +
			+(1.0+s)*(1.0+t)*o.p6[i]/8.0 - (1.0+s)*(1.0+t)*o.p7[i]/8.0

		// derivative of x with respect to s
		dxDs[i] = 0 +
			+(1.0-r)*o.db0stDs[i]/2.0 + (1.0+r)*o.db1stDs[i]/2.0 +
			-o.b2rt[i]/2.0 + o.b3rt[i]/2.0 +
			+(1.0-t)*o.db4rsDs[i]/2.0 + (1.0+t)*o.db5rsDs[i]/2.0 +

			+(1.0-r)*o.b0mt[i]/4.0 - (1.0-r)*o.b0pt[i]/4.0 +
			+(1.0+r)*o.b1mt[i]/4.0 - (1.0+r)*o.b1pt[i]/4.0 +

			-(1.0-r)*(1.0-t)*o.db0smDs[i]/4.0 - (1.0-r)*(1.0+t)*o.db0spDs[i]/4.0 +
			-(1.0+r)*(1.0-t)*o.db1smDs[i]/4.0 - (1.0+r)*(1.0+t)*o.db1spDs[i]/4.0 +

			+(1.0-t)*o.b2rm[i]/4.0 + (1.0+t)*o.b2rp[i]/4.0 +
			-(1.0-t)*o.b3rm[i]/4.0 - (1.0+t)*o.b3rp[i]/4.0 +

			-(1.0-r)*(1.0-t)*o.p0[i]/8.0 - (1.0+r)*(1.0-t)*o.p1[i]/8.0 +
			+(1.0+r)*(1.0-t)*o.p2[i]/8.0 + (1.0-r)*(1.0-t)*o.p3[i]/8.0 +
			-(1.0-r)*(1.0+t)*o.p4[i]/8.0 - (1.0+r)*(1.0+t)*o.p5[i]/8.0 +
			+(1.0+r)*(1.0+t)*o.p6[i]/8.0 + (1.0-r)*(1.0+t)*o.p7[i]/8.0

		// derivative of x with respect to t
		dxDt[i] = 0 +
			+(1.0-r)*o.db0stDt[i]/2.0 + (1.0+r)*o.db1stDt[i]/2.0 +
			+(1.0-s)*o.db2rtDt[i]/2.0 + (1.0+s)*o.db3rtDt[i]/2.0 +
			-o.b4rs[i]/2.0 + o.b5rs[i]/2.0 +

			-(1.0-r)*(1.0-s)*o.db0mtDt[i]/4.0 - (1.0-r)*(1.0+s)*o.db0ptDt[i]/4.0 +
			-(1.0+r)*(1.0-s)*o.db1mtDt[i]/4.0 - (1.0+r)*(1.0+s)*o.db1ptDt[i]/4.0 +

			+(1.0-r)*o.b0sm[i]/4.0 - (1.0-r)*o.b0sp[i]/4.0 +
			+(1.0+r)*o.b1sm[i]/4.0 - (1.0+r)*o.b1sp[i]/4.0 +

			+(1.0-s)*o.b2rm[i]/4.0 - (1.0-s)*o.b2rp[i]/4.0 +
			+(1.0+s)*o.b3rm[i]/4.0 - (1.0+s)*o.b3rp[i]/4.0 +

			-(1.0-r)*(1.0-s)*o.p0[i]/8.0 - (1.0+r)*(1.0-s)*o.p1[i]/8.0 +
			-(1.0+r)*(1.0+s)*o.p2[i]/8.0 - (1.0-r)*(1.0+s)*o.p3[i]/8.0 +
			+(1.0-r)*(1.0-s)*o.p4[i]/8.0 + (1.0+r)*(1.0-s)*o.p5[i]/8.0 +
			+(1.0+r)*(1.0+s)*o.p6[i]/8.0 + (1.0-r)*(1.0+s)*o.p7[i]/8.0
	}

	// skip second order derivatives
	if !secondDerivs {
		return
	}
	// second derivatives are zero
	if o.bdd == nil {
		for i := 0; i < o.ndim; i++ {
			ddxDrr[i] = 0.0
			ddxDss[i] = 0.0
			ddxDtt[i] = 0.0
			ddxDrs[i] = 0.0
			ddxDrt[i] = 0.0
			ddxDst[i] = 0.0
		}
		return
	}

	// compute second derivatives @ {r,s,t}
	o.bdd[0](o.ddb0stDss, o.ddb0stDtt, o.ddb0stDst, s, t)
	o.bdd[1](o.ddb1stDss, o.ddb1stDtt, o.ddb1stDst, s, t)
	o.bdd[2](o.ddb2rtDrr, o.ddb2rtDtt, o.ddb2rtDrt, r, t)
	o.bdd[3](o.ddb3rtDrr, o.ddb3rtDtt, o.ddb3rtDrt, r, t)
	o.bdd[4](o.ddb4rsDrr, o.ddb4rsDss, o.ddb4rsDrs, r, s)
	o.bdd[5](o.ddb5rsDrr, o.ddb5rsDss, o.ddb5rsDrs, r, s)

	// compute second derivatives @ edges
	o.bdd[0](o.ddb0smDss, o.tm1, o.tm2, s, m)
	o.bdd[0](o.ddb0spDss, o.tm1, o.tm2, s, p)
	o.bdd[0](o.tm1, o.ddb0mtDtt, o.tm2, m, t)
	o.bdd[0](o.tm1, o.ddb0ptDtt, o.tm2, p, t)

	o.bdd[1](o.ddb1smDss, o.tm1, o.tm2, s, m)
	o.bdd[1](o.ddb1spDss, o.tm1, o.tm2, s, p)
	o.bdd[1](o.tm1, o.ddb1mtDtt, o.tm2, m, t)
	o.bdd[1](o.tm1, o.ddb1ptDtt, o.tm2, p, t)

	o.bdd[2](o.ddb2rmDrr, o.tm1, o.tm2, r, m)
	o.bdd[2](o.ddb2rpDrr, o.tm1, o.tm2, r, p)
	o.bdd[3](o.ddb3rmDrr, o.tm1, o.tm2, r, m)
	o.bdd[3](o.ddb3rpDrr, o.tm1, o.tm2, r, p)

	// second order derivatives
	for i := 0; i < o.ndim; i++ {

		// derivative of dx/dr with respect to r
		ddxDrr[i] = 0 +
			+(1.0-s)*o.ddb2rtDrr[i]/2.0 + (1.0+s)*o.ddb3rtDrr[i]/2.0 +
			+(1.0-t)*o.ddb4rsDrr[i]/2.0 + (1.0+t)*o.ddb5rsDrr[i]/2.0 +
			-(1.0-s)*(1.0-t)*o.ddb2rmDrr[i]/4.0 - (1.0-s)*(1.0+t)*o.ddb2rpDrr[i]/4.0 +
			-(1.0+s)*(1.0-t)*o.ddb3rmDrr[i]/4.0 - (1.0+s)*(1.0+t)*o.ddb3rpDrr[i]/4.0

		// derivative of dx/ds with respect to s
		ddxDss[i] = 0 +
			+(1.0-r)*o.ddb0stDss[i]/2.0 + (1.0+r)*o.ddb1stDss[i]/2.0 +
			+(1.0-t)*o.ddb4rsDss[i]/2.0 + (1.0+t)*o.ddb5rsDss[i]/2.0 +
			-(1.0-r)*(1.0-t)*o.ddb0smDss[i]/4.0 - (1.0-r)*(1.0+t)*o.ddb0spDss[i]/4.0 +
			-(1.0+r)*(1.0-t)*o.ddb1smDss[i]/4.0 - (1.0+r)*(1.0+t)*o.ddb1spDss[i]/4.0

		// derivative of dx/dt with respect to t
		ddxDtt[i] = 0 +
			+(1.0-r)*o.ddb0stDtt[i]/2.0 + (1.0+r)*o.ddb1stDtt[i]/2.0 +
			+(1.0-s)*o.ddb2rtDtt[i]/2.0 + (1.0+s)*o.ddb3rtDtt[i]/2.0 +
			-(1.0-r)*(1.0-s)*o.ddb0mtDtt[i]/4.0 - (1.0-r)*(1.0+s)*o.ddb0ptDtt[i]/4.0 +
			-(1.0+r)*(1.0-s)*o.ddb1mtDtt[i]/4.0 - (1.0+r)*(1.0+s)*o.ddb1ptDtt[i]/4.0

		// derivative of dx/dr with respect to s
		ddxDrs[i] = 0 +
			-o.db0stDs[i]/2.0 + o.db1stDs[i]/2.0 +
			-o.db2rtDr[i]/2.0 + o.db3rtDr[i]/2.0 +
			+(1.0-t)*o.ddb4rsDrs[i]/2.0 + (1.0+t)*o.ddb5rsDrs[i]/2.0 +

			-o.b0mt[i]/4.0 + o.b0pt[i]/4.0 +
			+o.b1mt[i]/4.0 - o.b1pt[i]/4.0 +

			+(1.0-t)*o.db0smDs[i]/4.0 + (1.0+t)*o.db0spDs[i]/4.0 +
			-(1.0-t)*o.db1smDs[i]/4.0 - (1.0+t)*o.db1spDs[i]/4.0 +

			+(1.0-t)*o.db2rmDr[i]/4.0 + (1.0+t)*o.db2rpDr[i]/4.0 +
			-(1.0-t)*o.db3rmDr[i]/4.0 - (1.0+t)*o.db3rpDr[i]/4.0 +

			+(1.0-t)*o.p0[i]/8.0 - (1.0-t)*o.p1[i]/8.0 +
			+(1.0-t)*o.p2[i]/8.0 - (1.0-t)*o.p3[i]/8.0 +
			+(1.0+t)*o.p4[i]/8.0 - (1.0+t)*o.p5[i]/8.0 +
			+(1.0+t)*o.p6[i]/8.0 - (1.0+t)*o.p7[i]/8.0

		// derivative of dx/dr with respect to t
		ddxDrt[i] = 0 +
			-o.db0stDt[i]/2.0 + o.db1stDt[i]/2.0 +
			+(1.0-s)*o.ddb2rtDrt[i]/2.0 + (1.0+s)*o.ddb3rtDrt[i]/2.0 +
			-o.db4rsDr[i]/2.0 + o.db5rsDr[i]/2.0 +

			+(1.0-s)*o.db0mtDt[i]/4.0 + (1.0+s)*o.db0ptDt[i]/4.0 +
			-(1.0-s)*o.db1mtDt[i]/4.0 - (1.0+s)*o.db1ptDt[i]/4.0 +

			-o.b0sm[i]/4.0 + o.b0sp[i]/4.0 +
			+o.b1sm[i]/4.0 - o.b1sp[i]/4.0 +

			+(1.0-s)*o.db2rmDr[i]/4.0 - (1.0-s)*o.db2rpDr[i]/4.0 +
			+(1.0+s)*o.db3rmDr[i]/4.0 - (1.0+s)*o.db3rpDr[i]/4.0 +

			+(1.0-s)*o.p0[i]/8.0 - (1.0-s)*o.p1[i]/8.0 +
			-(1.0+s)*o.p2[i]/8.0 + (1.0+s)*o.p3[i]/8.0 +
			-(1.0-s)*o.p4[i]/8.0 + (1.0-s)*o.p5[i]/8.0 +
			+(1.0+s)*o.p6[i]/8.0 - (1.0+s)*o.p7[i]/8.0

		// derivative of dx/ds with respect to t
		ddxDst[i] = 0 +
			+(1.0-r)*o.ddb0stDst[i]/2.0 + (1.0+r)*o.ddb1stDst[i]/2.0 +
			-o.db2rtDt[i]/2.0 + o.db3rtDt[i]/2.0 +
			-o.db4rsDs[i]/2.0 + o.db5rsDs[i]/2.0 +

			+(1.0-r)*o.db0mtDt[i]/4.0 - (1.0-r)*o.db0ptDt[i]/4.0 +
			+(1.0+r)*o.db1mtDt[i]/4.0 - (1.0+r)*o.db1ptDt[i]/4.0 +

			+(1.0-r)*o.db0smDs[i]/4.0 - (1.0-r)*o.db0spDs[i]/4.0 +
			+(1.0+r)*o.db1smDs[i]/4.0 - (1.0+r)*o.db1spDs[i]/4.0 +

			-o.b2rm[i]/4.0 + o.b2rp[i]/4.0 +
			+o.b3rm[i]/4.0 - o.b3rp[i]/4.0 +

			+(1.0-r)*o.p0[i]/8.0 + (1.0+r)*o.p1[i]/8.0 +
			-(1.0+r)*o.p2[i]/8.0 - (1.0-r)*o.p3[i]/8.0 +
			-(1.0-r)*o.p4[i]/8.0 - (1.0+r)*o.p5[i]/8.0 +
			+(1.0+r)*o.p6[i]/8.0 + (1.0-r)*o.p7[i]/8.0
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
			o.e[0](x, r)
			x0[i] = x[0]
			y0[i] = x[1]
		}
		plt.Plot(x0, y0, argsBry)

		// draw B1(s)
		for j := 0; j < npts[1]; j++ {
			s := -1 + 2*float64(j)/float64(npts[1]-1)
			o.e[1](x, s)
			x1[j] = x[0]
			y1[j] = x[1]
		}
		plt.Plot(x1, y1, argsBry)

		// draw B2(r)
		for i := 0; i < npts[0]; i++ {
			r := -1 + 2*float64(i)/float64(npts[0]-1)
			o.e[2](x, r)
			x0[i] = x[0]
			y0[i] = x[1]
		}
		plt.Plot(x0, y0, argsBry)

		// draw B3(s)
		for j := 0; j < npts[1]; j++ {
			s := -1 + 2*float64(j)/float64(npts[1]-1)
			o.e[3](x, s)
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
