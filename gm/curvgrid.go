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

// CurvGrid holds metrics data for 2d or 3d grids represented by curvilinear coordinates
type CurvGrid struct {
	Ndim int            // space dimension
	N0   int            // number of points along 0-direction
	N1   int            // number of points along 1-direction
	N2   int            // number of points along 2-direction [default = 1]
	M2d  [][]*Metrics   // [n1][n0] metrics in 2D
	M3d  [][][]*Metrics // [n2][n1][n0] metrics in 3D
}

// SetTransfinite2d sets grid from 2D transfinite mapping
//  trf -- 2D transfinite structure
//  R   -- [n1] reference coordinates along r-direction
//  S   -- [n2] reference coordinates along s-direction
func (o *CurvGrid) SetTransfinite2d(trf *Transfinite, R, S []float64) {

	// input
	o.Ndim = 2
	o.N0 = len(R)
	o.N1 = len(S)

	// auxiliary
	x := la.NewVector(2)
	u := la.NewVector(2)
	dxdr, dxds := la.NewVector(2), la.NewVector(2)
	ddxdrr, ddxdss, ddxdrs := la.NewVector(2), la.NewVector(2), la.NewVector(2)

	// compute metrics
	o.M2d = make([][]*Metrics, o.N1)
	for j := 0; j < o.N1; j++ {
		o.M2d[j] = make([]*Metrics, o.N0)
		for i := 0; i < o.N0; i++ {

			// derivatives
			u[0], u[1] = R[i], S[j]
			trf.PointAndDerivs(x, dxdr, dxds, nil, ddxdrr, ddxdss, nil, ddxdrs, nil, nil, u)

			// metrics
			o.M2d[j][i] = NewMetrics2d(u, x, dxdr, dxds, ddxdrr, ddxdss, ddxdrs)
		}
	}
}

// SetTransfinite3d sets grid from 3D transfinite mapping
func (o *CurvGrid) SetTransfinite3d(B, Bd, Bdd []fun.Vss, R, S, T []float64) {

	// compute metrics
	o.M3d = make([][][]*Metrics, o.N2)
	for k := 0; k < o.N2; k++ {
		o.M3d[k] = make([][]*Metrics, o.N1)
		for j := 0; j < o.N1; j++ {
			o.M3d[k][j] = make([]*Metrics, o.N0)
			for i := 0; i < o.N0; i++ {
				//o.m3d[k][j][i] = NewMetrics3d()
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
	if o.Ndim == 2 {
		for j := 0; j < o.N1; j++ {
			for i := 0; i < o.N0; i++ {
				m := o.M2d[j][i]
				DrawArrow2d(m.X, m.CovG0, true, scale, argsG0)
				DrawArrow2d(m.X, m.CovG1, true, scale, argsG1)
			}
		}
		return
	}
	if argsG2 == nil {
		argsG2 = &plt.A{C: plt.C(2, 0), Scale: 7, Z: 10}
	}
	chk.Panic("TODO")
}
