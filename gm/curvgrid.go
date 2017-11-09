// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

// CurvGrid holds metrics data for 2d or 3d grids represented by curvilinear coordinates
//
//   Notation:
//              m,n,p -- indices used for grid points
//              i,j,k -- indices used for dimension (x,y,z)
//
//   NOTE: the deep3 structure mtr holds data with the outer index corresponding to z
//
//   Example: the covariant vector @ (m,n,p) is:     o.mtr[p][n][m].GovG0
//            the first component of this vector is: o.mtr[p][n][m].GovG0[0]
//
type CurvGrid struct {
	ndim int            // space dimension
	npts []int          // number of points along each direction. In 2D, npts[2] := 1
	mtr  [][][]*Metrics // [n2][n1][n0] metrics in 2D (with n2=1) or 3D
}

// SetTransfinite2d sets grid from 2D transfinite mapping
//  trf -- 2D transfinite structure
//  R   -- [n1] reference coordinates along r-direction
//  S   -- [n2] reference coordinates along s-direction
func (o *CurvGrid) SetTransfinite2d(trf *Transfinite, R, S []float64) {

	// input
	o.ndim = 2
	o.npts = []int{len(R), len(S), 1}

	// auxiliary
	x := la.NewVector(2)
	u := la.NewVector(2)
	dxdr, dxds := la.NewVector(2), la.NewVector(2)
	ddxdrr, ddxdss, ddxdrs := la.NewVector(2), la.NewVector(2), la.NewVector(2)

	// compute metrics
	p := 0
	o.mtr = make([][][]*Metrics, 1)
	o.mtr[p] = make([][]*Metrics, o.npts[1])
	for n := 0; n < o.npts[1]; n++ {
		o.mtr[p][n] = make([]*Metrics, o.npts[0])
		for m := 0; m < o.npts[0]; m++ {

			// derivatives
			u[0], u[1] = R[m], S[n]
			trf.PointAndDerivs(x, dxdr, dxds, nil, ddxdrr, ddxdss, nil, ddxdrs, nil, nil, u)

			// metrics
			o.mtr[p][n][m] = NewMetrics2d(u, x, dxdr, dxds, ddxdrr, ddxdss, ddxdrs)
		}
	}
}

// SetTransfinite3d sets grid from 3D transfinite mapping
func (o *CurvGrid) SetTransfinite3d(trf *Transfinite, R, S, T []float64) {

	// input
	o.ndim = 3
	o.npts = []int{len(R), len(S), len(T)}

	// auxiliary
	x := la.NewVector(3)
	u := la.NewVector(3)
	dxdr, dxds, dxdt := la.NewVector(3), la.NewVector(3), la.NewVector(3)
	ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst := la.NewVector(3), la.NewVector(3), la.NewVector(3), la.NewVector(3), la.NewVector(3), la.NewVector(3)

	// compute metrics
	o.mtr = make([][][]*Metrics, o.npts[2])
	for p := 0; p < o.npts[2]; p++ {
		o.mtr[p] = make([][]*Metrics, o.npts[1])
		for n := 0; n < o.npts[1]; n++ {
			o.mtr[p][n] = make([]*Metrics, o.npts[0])
			for m := 0; m < o.npts[0]; m++ {

				// derivatives
				u[0], u[1], u[2] = R[m], S[n], T[p]
				trf.PointAndDerivs(x, dxdr, dxds, dxdt, ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst, u)

				// metrics
				o.mtr[p][n][m] = NewMetrics3d(u, x, dxdr, dxds, dxdt, ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst)
			}
		}
	}
}

// DrawBases draw basis vectors
func (o *CurvGrid) DrawBases(scale float64, argsG0, argsG1, argsG2 *plt.A) {
	if argsG0 == nil {
		argsG0 = &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10}
	}
	if argsG1 == nil {
		argsG1 = &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10}
	}
	if o.ndim == 2 {
		p := 0
		for n := 0; n < o.npts[1]; n++ {
			for m := 0; m < o.npts[0]; m++ {
				M := o.mtr[p][n][m]
				DrawArrow2d(M.X, M.CovG0, true, scale, argsG0)
				DrawArrow2d(M.X, M.CovG1, true, scale, argsG1)
			}
		}
		return
	}
	if argsG2 == nil {
		argsG2 = &plt.A{C: plt.C(2, 0), Scale: 7, Z: 10}
	}
	for p := 0; p < o.npts[2]; p++ {
		for n := 0; n < o.npts[1]; n++ {
			for m := 0; m < o.npts[0]; m++ {
				M := o.mtr[p][n][m]
				DrawArrow3d(M.X, M.CovG0, true, scale, argsG0)
				DrawArrow3d(M.X, M.CovG1, true, scale, argsG1)
				DrawArrow3d(M.X, M.CovG2, true, scale, argsG2)
			}
		}
	}
}
