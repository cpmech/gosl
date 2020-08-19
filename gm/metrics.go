// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"gosl/la"
	"gosl/utl"
)

// Metrics holds data related to a position in a space represented by curvilinear coordinates
type Metrics struct {
	U           la.Vector     // reference coordinates {r,s,t}
	X           la.Vector     // physical coordinates {x,y,z}
	CovG0       la.Vector     // covariant basis g_0 = d{x}/dr
	CovG1       la.Vector     // covariant basis g_1 = d{x}/ds
	CovG2       la.Vector     // covariant basis g_2 = d{x}/dt
	CntG0       la.Vector     // contravariant basis g_0 = dr/d{x} (gradients)
	CntG1       la.Vector     // contravariant basis g_1 = ds/d{x} (gradients)
	CntG2       la.Vector     // contravariant basis g_2 = dt/d{x} (gradients)
	CovGmat     *la.Matrix    // covariant metrics g_ij = g_i ⋅ g_j
	CntGmat     *la.Matrix    // contravariant metrics g^ij = g^i ⋅ g^j
	DetCovGmat  float64       // determinant of covariant g matrix = det(CovGmat)
	Homogeneous bool          // homogeneous grid => nil second order derivatives and Christoffel symbols
	GammaS      [][][]float64 // [k][i][j] Christoffel coefficients of second kind (non-homogeneous)
	L           []float64     // [3] L-coefficients = sum(Γ_ij^k ⋅ g^ij) (non-homogeneous)
}

// NewMetrics2d allocate new 2D metrics structure
//  NOTE: the second order derivatives (from ddxdrr) may be nil => homogeneous grid
func NewMetrics2d(u, x, dxdr, dxds, ddxdrr, ddxdss, ddxdrs la.Vector) (o *Metrics) {

	// input
	o = new(Metrics)
	o.U = u.GetCopy()
	o.X = x.GetCopy()
	o.CovG0 = dxdr.GetCopy()
	o.CovG1 = dxds.GetCopy()

	// covariant metrics
	o.CovGmat = la.NewMatrix(2, 2)
	o.CovGmat.Set(0, 0, la.VecDot(o.CovG0, o.CovG0))
	o.CovGmat.Set(1, 1, la.VecDot(o.CovG1, o.CovG1))
	o.CovGmat.Set(0, 1, la.VecDot(o.CovG0, o.CovG1))
	o.CovGmat.Set(1, 0, o.CovGmat.Get(0, 1))

	// contravariant metrics
	o.CntGmat = la.NewMatrix(2, 2)
	o.DetCovGmat = la.MatInvSmall(o.CntGmat, o.CovGmat, 1e-13)

	// contravariant vectors
	o.CntG0 = la.NewVector(2)
	o.CntG1 = la.NewVector(2)
	for i := 0; i < 2; i++ {
		o.CntG0[i] += o.CntGmat.Get(0, 0)*o.CovG0[i] + o.CntGmat.Get(0, 1)*o.CovG1[i]
		o.CntG1[i] += o.CntGmat.Get(1, 0)*o.CovG0[i] + o.CntGmat.Get(1, 1)*o.CovG1[i]
	}

	// check if homogeneous grid
	o.Homogeneous = ddxdrr == nil
	if o.Homogeneous {
		return
	}

	// Christoffel vectors
	Γ00, Γ11, Γ01 := ddxdrr, ddxdss, ddxdrs

	// Christoffel symbols of second kind
	o.GammaS = utl.Deep3alloc(2, 2, 2)
	o.GammaS[0][0][0] = la.VecDot(Γ00, o.CntG0)
	o.GammaS[0][1][1] = la.VecDot(Γ11, o.CntG0)
	o.GammaS[0][0][1] = la.VecDot(Γ01, o.CntG0)
	o.GammaS[0][1][0] = o.GammaS[0][0][1]
	o.GammaS[1][0][0] = la.VecDot(Γ00, o.CntG1)
	o.GammaS[1][1][1] = la.VecDot(Γ11, o.CntG1)
	o.GammaS[1][0][1] = la.VecDot(Γ01, o.CntG1)
	o.GammaS[1][1][0] = o.GammaS[1][0][1]

	// L-coefficients
	o.L = make([]float64, 2)
	o.L[0] = o.GammaS[0][0][0]*o.CntGmat.Get(0, 0) + o.GammaS[0][1][1]*o.CntGmat.Get(1, 1) + 2.0*o.GammaS[0][0][1]*o.CntGmat.Get(0, 1)
	o.L[1] = o.GammaS[1][0][0]*o.CntGmat.Get(0, 0) + o.GammaS[1][1][1]*o.CntGmat.Get(1, 1) + 2.0*o.GammaS[1][0][1]*o.CntGmat.Get(0, 1)
	return
}

// NewMetrics3d allocate new 3D metrics structure
//  NOTE: the second order derivatives (from ddxdrr) may be nil => homogeneous grid
func NewMetrics3d(u, x, dxdr, dxds, dxdt, ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst la.Vector) (o *Metrics) {

	// input
	o = new(Metrics)
	o.U = u.GetCopy()
	o.X = x.GetCopy()
	o.CovG0 = dxdr.GetCopy()
	o.CovG1 = dxds.GetCopy()
	o.CovG2 = dxdt.GetCopy()

	// covariant metrics
	o.CovGmat = la.NewMatrix(3, 3)
	o.CovGmat.Set(0, 0, la.VecDot(o.CovG0, o.CovG0))
	o.CovGmat.Set(1, 1, la.VecDot(o.CovG1, o.CovG1))
	o.CovGmat.Set(2, 2, la.VecDot(o.CovG2, o.CovG2))
	o.CovGmat.Set(0, 1, la.VecDot(o.CovG0, o.CovG1))
	o.CovGmat.Set(1, 2, la.VecDot(o.CovG1, o.CovG2))
	o.CovGmat.Set(2, 0, la.VecDot(o.CovG2, o.CovG0))
	o.CovGmat.Set(1, 0, o.CovGmat.Get(0, 1))
	o.CovGmat.Set(2, 1, o.CovGmat.Get(1, 2))
	o.CovGmat.Set(0, 2, o.CovGmat.Get(2, 0))

	// contravariant metrics
	o.CntGmat = la.NewMatrix(3, 3)
	o.DetCovGmat = la.MatInvSmall(o.CntGmat, o.CovGmat, 1e-13)

	// contravariant vectors
	o.CntG0 = la.NewVector(3)
	o.CntG1 = la.NewVector(3)
	o.CntG2 = la.NewVector(3)
	for i := 0; i < 3; i++ {
		o.CntG0[i] += o.CntGmat.Get(0, 0)*o.CovG0[i] + o.CntGmat.Get(0, 1)*o.CovG1[i] + o.CntGmat.Get(0, 2)*o.CovG2[i]
		o.CntG1[i] += o.CntGmat.Get(1, 0)*o.CovG0[i] + o.CntGmat.Get(1, 1)*o.CovG1[i] + o.CntGmat.Get(1, 2)*o.CovG2[i]
		o.CntG2[i] += o.CntGmat.Get(2, 0)*o.CovG0[i] + o.CntGmat.Get(2, 1)*o.CovG1[i] + o.CntGmat.Get(2, 2)*o.CovG2[i]
	}

	// check if homogeneous grid
	o.Homogeneous = ddxdrr == nil
	if o.Homogeneous {
		return
	}

	// Christoffel vectors
	Γ00, Γ11, Γ22, Γ01, Γ02, Γ12 := ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst

	// Christoffel symbols of second kind
	o.GammaS = utl.Deep3alloc(3, 3, 3)
	o.GammaS[0][0][0] = la.VecDot(Γ00, o.CntG0)
	o.GammaS[0][1][1] = la.VecDot(Γ11, o.CntG0)
	o.GammaS[0][2][2] = la.VecDot(Γ22, o.CntG0)
	o.GammaS[0][0][1] = la.VecDot(Γ01, o.CntG0)
	o.GammaS[0][0][2] = la.VecDot(Γ02, o.CntG0)
	o.GammaS[0][1][2] = la.VecDot(Γ12, o.CntG0)
	o.GammaS[0][1][0] = o.GammaS[0][0][1]
	o.GammaS[0][2][0] = o.GammaS[0][0][2]
	o.GammaS[0][2][1] = o.GammaS[0][1][2]

	o.GammaS[1][0][0] = la.VecDot(Γ00, o.CntG1)
	o.GammaS[1][1][1] = la.VecDot(Γ11, o.CntG1)
	o.GammaS[1][2][2] = la.VecDot(Γ22, o.CntG1)
	o.GammaS[1][0][1] = la.VecDot(Γ01, o.CntG1)
	o.GammaS[1][0][2] = la.VecDot(Γ02, o.CntG1)
	o.GammaS[1][1][2] = la.VecDot(Γ12, o.CntG1)
	o.GammaS[1][1][0] = o.GammaS[1][0][1]
	o.GammaS[1][2][0] = o.GammaS[1][0][2]
	o.GammaS[1][2][1] = o.GammaS[1][1][2]

	o.GammaS[2][0][0] = la.VecDot(Γ00, o.CntG2)
	o.GammaS[2][1][1] = la.VecDot(Γ11, o.CntG2)
	o.GammaS[2][2][2] = la.VecDot(Γ22, o.CntG2)
	o.GammaS[2][0][1] = la.VecDot(Γ01, o.CntG2)
	o.GammaS[2][0][2] = la.VecDot(Γ02, o.CntG2)
	o.GammaS[2][1][2] = la.VecDot(Γ12, o.CntG2)
	o.GammaS[2][1][0] = o.GammaS[2][0][1]
	o.GammaS[2][2][0] = o.GammaS[2][0][2]
	o.GammaS[2][2][1] = o.GammaS[2][1][2]

	// L-coefficients
	o.L = make([]float64, 3)
	o.L[0] = o.GammaS[0][0][0]*o.CntGmat.Get(0, 0) + o.GammaS[0][1][1]*o.CntGmat.Get(1, 1) + o.GammaS[0][2][2]*o.CntGmat.Get(2, 2) + 2.0*o.GammaS[0][0][1]*o.CntGmat.Get(0, 1) + 2.0*o.GammaS[0][0][2]*o.CntGmat.Get(0, 2) + 2.0*o.GammaS[0][1][2]*o.CntGmat.Get(1, 2)
	o.L[1] = o.GammaS[1][0][0]*o.CntGmat.Get(0, 0) + o.GammaS[1][1][1]*o.CntGmat.Get(1, 1) + o.GammaS[1][2][2]*o.CntGmat.Get(2, 2) + 2.0*o.GammaS[1][0][1]*o.CntGmat.Get(0, 1) + 2.0*o.GammaS[1][0][2]*o.CntGmat.Get(0, 2) + 2.0*o.GammaS[1][1][2]*o.CntGmat.Get(1, 2)
	o.L[2] = o.GammaS[2][0][0]*o.CntGmat.Get(0, 0) + o.GammaS[2][1][1]*o.CntGmat.Get(1, 1) + o.GammaS[2][2][2]*o.CntGmat.Get(2, 2) + 2.0*o.GammaS[2][0][1]*o.CntGmat.Get(0, 1) + 2.0*o.GammaS[2][0][2]*o.CntGmat.Get(0, 2) + 2.0*o.GammaS[2][1][2]*o.CntGmat.Get(1, 2)
	return
}
