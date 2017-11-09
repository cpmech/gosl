// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

// CurvGrid holds metrics data for 2d or 3d grids represented by curvilinear coordinates
type CurvGrid struct {
	ndim int            // space dimension
	npts []int          // number of points along each direction
	m2d  [][]*Metrics   // [n1][n0] metrics in 2D
	m3d  [][][]*Metrics // [n2][n1][n0] metrics in 3D
}

// SetTransfinite2d sets grid from 2D transfinite mapping
//  trf -- 2D transfinite structure
//  R   -- [n1] reference coordinates along r-direction
//  S   -- [n2] reference coordinates along s-direction
func (o *CurvGrid) SetTransfinite2d(trf *Transfinite, R, S []float64) {

	// input
	o.ndim = 2
	o.npts = []int{len(R), len(S), 0}

	// auxiliary
	x := la.NewVector(2)
	u := la.NewVector(2)
	dxdr, dxds := la.NewVector(2), la.NewVector(2)
	ddxdrr, ddxdss, ddxdrs := la.NewVector(2), la.NewVector(2), la.NewVector(2)

	// compute metrics
	o.m2d = make([][]*Metrics, o.npts[1])
	for j := 0; j < o.npts[1]; j++ {
		o.m2d[j] = make([]*Metrics, o.npts[0])
		for i := 0; i < o.npts[0]; i++ {

			// derivatives
			u[0], u[1] = R[i], S[j]
			trf.PointAndDerivs(x, dxdr, dxds, nil, ddxdrr, ddxdss, nil, ddxdrs, nil, nil, u)

			// metrics
			o.m2d[j][i] = NewMetrics2d(u, x, dxdr, dxds, ddxdrr, ddxdss, ddxdrs)
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
	o.m3d = make([][][]*Metrics, o.npts[2])
	for k := 0; k < o.npts[2]; k++ {
		o.m3d[k] = make([][]*Metrics, o.npts[1])
		for j := 0; j < o.npts[1]; j++ {
			o.m3d[k][j] = make([]*Metrics, o.npts[0])
			for i := 0; i < o.npts[0]; i++ {

				// derivatives
				u[0], u[1], u[2] = R[i], S[j], T[k]
				trf.PointAndDerivs(x, dxdr, dxds, dxdt, ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst, u)

				// metrics
				o.m3d[k][j][i] = NewMetrics3d(u, x, dxdr, dxds, dxdt, ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst)
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
		for j := 0; j < o.npts[1]; j++ {
			for i := 0; i < o.npts[0]; i++ {
				m := o.m2d[j][i]
				DrawArrow2d(m.X, m.CovG0, true, scale, argsG0)
				DrawArrow2d(m.X, m.CovG1, true, scale, argsG1)
			}
		}
		return
	}
	if argsG2 == nil {
		argsG2 = &plt.A{C: plt.C(2, 0), Scale: 7, Z: 10}
	}
	for k := 0; k < o.npts[2]; k++ {
		for j := 0; j < o.npts[1]; j++ {
			for i := 0; i < o.npts[0]; i++ {
				m := o.m3d[k][j][i]
				DrawArrow3d(m.X, m.CovG0, true, scale, argsG0)
				DrawArrow3d(m.X, m.CovG1, true, scale, argsG1)
				DrawArrow3d(m.X, m.CovG2, true, scale, argsG2)
			}
		}
	}
}
