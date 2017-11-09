// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
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
	umin []float64      // min reference coordinates [3]
	umax []float64      // max reference coordinates [3]
	xmin []float64      // min physical coordinates [3]
	xmax []float64      // max physical coordinates [3]
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

	// limits
	o.limits()
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

	// limits
	o.limits()
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

// interface methods //////////////////////////////////////////////////////////////////////////////

// Ndim returns the number of dimensions (2D or 3D)
func (o *CurvGrid) Ndim() int {
	return o.ndim
}

// Npts returns number of points along idim dimension
func (o *CurvGrid) Npts(idim int) int {
	return o.npts[idim]
}

// Size returns total number of points
func (o *CurvGrid) Size() int {
	return o.npts[0] * o.npts[1] * o.npts[2]
}

// Umin returns the minimum reference coordinate at dimension idim
func (o *CurvGrid) Umin(idim int) float64 {
	return o.umin[idim]
}

// Umax returns the maximum reference coordinate at dimension idim
func (o *CurvGrid) Umax(idim int) float64 {
	return o.umax[idim]
}

// Xmin returns the minimum physical coordinate at dimension idim
func (o *CurvGrid) Xmin(idim int) float64 {
	return o.xmin[idim]
}

// Xmax returns the maximum physical coordinate at dimension idim
func (o *CurvGrid) Xmax(idim int) float64 {
	return o.xmax[idim]
}

// U returns the reference coordinates at point m,n,p
func (o *CurvGrid) U(m, n, p int) la.Vector {
	return o.mtr[p][n][m].U
}

// X returns the physical coordinates at point m,n,p
func (o *CurvGrid) X(m, n, p int) la.Vector {
	return o.mtr[p][n][m].X
}

// CovarBasis returns the [k] covariant basis g_{k} = d{x[k]}/d{u[k]} [@ point m,n,p]
func (o *CurvGrid) CovarBasis(m, n, p, k int) la.Vector {
	if k == 0 {
		return o.mtr[p][n][m].CovG0
	}
	if k == 1 {
		return o.mtr[p][n][m].CovG1
	}
	return o.mtr[p][n][m].CovG2
}

// CovarMatrix returns the covariant metrics g_ij = g_i ⋅ g_j [@ point m,n,p]
func (o *CurvGrid) CovarMatrix(m, n, p int) *la.Matrix {
	return o.mtr[p][n][m].CovGmat
}

// ContraMatrix returns contravariant metrics g^ij = g^i ⋅ g^j [@ point m,n,p]
func (o *CurvGrid) ContraMatrix(m, n, p int) *la.Matrix {
	return o.mtr[p][n][m].CntGmat
}

// DetCovarMatrix returns the determinant of covariant g matrix = det(CovariantMatrix) [@ point m,n,p]
func (o *CurvGrid) DetCovarMatrix(m, n, p int) float64 {
	return o.mtr[p][n][m].DetCovGmat
}

// GammaS returns the [k][i][j] Christoffel coefficients of second kind [@ point m,n,p]
func (o *CurvGrid) GammaS(m, n, p, k, i, j int) float64 {
	return o.mtr[p][n][m].GammaS[k][i][j]
}

// Lcoeff returns the [k] L-coefficients = sum(Γ_ij^k ⋅ g^ij) [@ point m,n,p]
func (o *CurvGrid) Lcoeff(m, n, p, k int) float64 {
	return o.mtr[p][n][m].L[k]
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func (o *CurvGrid) limits() {
	o.umin = []float64{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64}
	o.umax = []float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	o.xmin = []float64{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64}
	o.xmax = []float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	if o.ndim == 2 {
		o.umin[2], o.umax[2] = -1, -1
		o.xmin[2], o.xmax[2] = 0, 0
	}
	for p := 0; p < o.npts[2]; p++ {
		for n := 0; n < o.npts[1]; n++ {
			for m := 0; m < o.npts[0]; m++ {
				for i := 0; i < o.ndim; i++ {
					o.umin[i] = utl.Min(o.umin[i], o.mtr[p][n][m].U[i])
					o.umax[i] = utl.Max(o.umax[i], o.mtr[p][n][m].U[i])
					o.xmin[i] = utl.Min(o.xmin[i], o.mtr[p][n][m].X[i])
					o.xmax[i] = utl.Max(o.xmax[i], o.mtr[p][n][m].X[i])
				}
			}
		}
	}
}
