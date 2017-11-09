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
//      m,n,p -- indices used for grid points
//      i,j,k -- indices used for dimension (x,y,z)
//      Ex: the covariant vector @ (m,n,p) is: o.mtr[p][n][m].GovG0
//          the i component of this vector is: o.mtr[p][n][m].GovG0[i]
//
//   NOTE: (1) the deep3 structure mtr holds data with the outer index corresponding to z
//             i.e. o.mtr[idxZ][idxY][idxX]
//         (2) the reference coordinates of generated rectangular grids are assumed to be
//             -1 ≤ u ≤ +1
//
type CurvGrid struct {
	ndim int            // space dimension
	npts []int          // number of points along each direction. In 2D, npts[2] := 1
	mtr  [][][]*Metrics // [n2][n1][n0] metrics in 2D (with n2=1) or 3D
	umin []float64      // min reference coordinates [3]
	umax []float64      // max reference coordinates [3]
	xmin []float64      // min physical coordinates [3]
	xmax []float64      // max physical coordinates [3]
	edge [][]int        // ids of points on edges: [edge0, edge1, edge2, edge3]
	face [][]int        // ids of points on faces: [face0, face1, face2, face3, face4, face5]
}

// RectGenUniform generates uniform coordinates of a rectangular grid
//  xmin -- min x-y-z values, len==ndim: [xmin, ymin, zmin]
//  xmax -- max x-y-z values, len==ndim: [xmax, ymax, zmax]
//  npts -- number of points along each direction len==ndim: [n0, n1, n2] (must be greater than 2)
//
//     -1 ≤ u ≤ +1
//     x(u) = xmin + (xmax - xmin) ⋅ (1 + u) / 2
//     u(x) = -1 + 2⋅(x - xmin) / (xmax - xmin)
//     dx/du = (xmax - xmin) / 2
//
func (o *CurvGrid) RectGenUniform(xmin, xmax []float64, npts []int) {

	// input
	o.ndim = len(xmin)

	// 2D grid
	if o.ndim == 2 {

		// input
		o.npts = []int{npts[0], npts[1], 1}

		// auxiliary
		x := la.NewVector(2)
		u := la.NewVector(2)
		dxdr, dxds := la.NewVector(2), la.NewVector(2)

		// (constant) derivatives (all the 2nd order derivatives are zero)
		dxdr[0] = (xmax[0] - xmin[0]) / 2.0
		dxds[1] = (xmax[1] - xmin[1]) / 2.0

		// compute metrics
		p := 0
		o.mtr = make([][][]*Metrics, 1)
		o.mtr[p] = make([][]*Metrics, o.npts[1])
		for n := 0; n < o.npts[1]; n++ {
			o.mtr[p][n] = make([]*Metrics, o.npts[0])
			u[1] = -1.0 + 2.0*float64(n)/float64(o.npts[1]-1)
			x[1] = xmin[1] + (xmax[1]-xmin[1])*(1.0+u[1])/2.0
			for m := 0; m < o.npts[0]; m++ {
				u[0] = -1.0 + 2.0*float64(m)/float64(o.npts[0]-1)
				x[0] = xmin[0] + (xmax[0]-xmin[0])*(1.0+u[0])/2.0
				o.mtr[p][n][m] = NewMetrics2d(u, x, dxdr, dxds, nil, nil, nil)
			}
		}

		// 3D grid
	} else {

		// input
		o.npts = utl.IntCopy(npts)

		// auxiliary
		x := la.NewVector(3)
		u := la.NewVector(3)
		dxdr, dxds, dxdt := la.NewVector(3), la.NewVector(3), la.NewVector(3)

		// (constant) derivatives (all the 2nd order derivatives are zero)
		dxdr[0] = (xmax[0] - xmin[0]) / 2.0
		dxds[1] = (xmax[1] - xmin[1]) / 2.0
		dxdt[2] = (xmax[2] - xmin[2]) / 2.0

		// compute metrics
		o.mtr = make([][][]*Metrics, o.npts[2])
		for p := 0; p < o.npts[2]; p++ {
			o.mtr[p] = make([][]*Metrics, o.npts[1])
			u[2] = -1.0 + 2.0*float64(p)/float64(o.npts[2]-1)
			x[2] = xmin[2] + (xmax[2]-xmin[2])*(1.0+u[2])/2.0
			for n := 0; n < o.npts[1]; n++ {
				o.mtr[p][n] = make([]*Metrics, o.npts[0])
				u[1] = -1.0 + 2.0*float64(n)/float64(o.npts[1]-1)
				x[1] = xmin[1] + (xmax[1]-xmin[1])*(1.0+u[1])/2.0
				for m := 0; m < o.npts[0]; m++ {
					u[0] = -1.0 + 2.0*float64(m)/float64(o.npts[0]-1)
					x[0] = xmin[0] + (xmax[0]-xmin[0])*(1.0+u[0])/2.0
					o.mtr[p][n][m] = NewMetrics3d(u, x, dxdr, dxds, dxdt, nil, nil, nil, nil, nil, nil)
				}
			}
		}
	}

	// limits and boundaries
	o.umin = []float64{-1, -1, -1}
	o.umax = []float64{+1, +1, +1}
	o.xmin = utl.GetCopy(xmin)
	o.xmax = utl.GetCopy(xmax)
	o.boundaries()
}

// RectSet2d sets rectangular grid with given coordinates
//
//     -1 ≤ u ≤ +1
//     x(u) = xmin + (xmax - xmin) ⋅ (1 + u) / 2
//     u(x) = -1 + 2⋅(x - xmin) / (xmax - xmin)
//     dx/du = (xmax - xmin) / 2
//
func (o *CurvGrid) RectSet2d(X, Y []float64) {

	// input
	o.ndim = 2
	o.npts = []int{len(X), len(Y), 1}

	// limits
	o.umin = []float64{-1, -1, -1}
	o.umax = []float64{+1, +1, +1}
	o.xmin = make([]float64, 2)
	o.xmax = make([]float64, 2)
	o.xmin[0], o.xmax[0] = utl.MinMax(X)
	o.xmin[1], o.xmax[1] = utl.MinMax(Y)

	// auxiliary
	x := la.NewVector(2)
	u := la.NewVector(2)
	dxdr, dxds := la.NewVector(2), la.NewVector(2)

	// (constant) derivatives (all the 2nd order derivatives are zero)
	dxdr[0] = (o.xmax[0] - o.xmin[0]) / 2.0
	dxds[1] = (o.xmax[1] - o.xmin[1]) / 2.0

	// compute metrics
	p := 0
	o.mtr = make([][][]*Metrics, 1)
	o.mtr[p] = make([][]*Metrics, o.npts[1])
	for n := 0; n < o.npts[1]; n++ {
		o.mtr[p][n] = make([]*Metrics, o.npts[0])
		x[1] = Y[n]
		u[1] = -1.0 + 2.0*(x[1]-o.xmin[1])/(o.xmax[1]-o.xmin[1])
		for m := 0; m < o.npts[0]; m++ {
			x[0] = X[m]
			u[0] = -1.0 + 2.0*(x[0]-o.xmin[0])/(o.xmax[0]-o.xmin[0])
			o.mtr[p][n][m] = NewMetrics2d(u, x, dxdr, dxds, nil, nil, nil)
		}
	}

	// boundaries
	o.boundaries()
}

// RectSet3d sets rectangular grid with given coordinates
//
//     -1 ≤ u ≤ +1
//     x(u) = xmin + (xmax - xmin) ⋅ (1 + u) / 2
//     u(x) = -1 + 2⋅(x - xmin) / (xmax - xmin)
//     dx/du = (xmax - xmin) / 2
//
func (o *CurvGrid) RectSet3d(X, Y, Z []float64) {

	// input
	o.ndim = 3
	o.npts = []int{len(X), len(Y), len(Z)}

	// limits
	o.umin = []float64{-1, -1, -1}
	o.umax = []float64{+1, +1, +1}
	o.xmin = make([]float64, 3)
	o.xmax = make([]float64, 3)
	o.xmin[0], o.xmax[0] = utl.MinMax(X)
	o.xmin[1], o.xmax[1] = utl.MinMax(Y)
	o.xmin[2], o.xmax[2] = utl.MinMax(Z)

	// auxiliary
	x := la.NewVector(3)
	u := la.NewVector(3)
	dxdr, dxds, dxdt := la.NewVector(3), la.NewVector(3), la.NewVector(3)

	// (constant) derivatives (all the 2nd order derivatives are zero)
	dxdr[0] = (o.xmax[0] - o.xmin[0]) / 2.0
	dxds[1] = (o.xmax[1] - o.xmin[1]) / 2.0
	dxdt[2] = (o.xmax[2] - o.xmin[2]) / 2.0

	// compute metrics
	o.mtr = make([][][]*Metrics, o.npts[2])
	for p := 0; p < o.npts[2]; p++ {
		o.mtr[p] = make([][]*Metrics, o.npts[1])
		x[2] = Z[p]
		u[2] = -1.0 + 2.0*(x[2]-o.xmin[2])/(o.xmax[2]-o.xmin[2])
		for n := 0; n < o.npts[1]; n++ {
			o.mtr[p][n] = make([]*Metrics, o.npts[0])
			x[1] = Y[n]
			u[1] = -1.0 + 2.0*(x[1]-o.xmin[1])/(o.xmax[1]-o.xmin[1])
			for m := 0; m < o.npts[0]; m++ {
				x[0] = X[m]
				u[0] = -1.0 + 2.0*(x[0]-o.xmin[0])/(o.xmax[0]-o.xmin[0])
				o.mtr[p][n][m] = NewMetrics3d(u, x, dxdr, dxds, dxdt, nil, nil, nil, nil, nil, nil)
			}
		}
	}

	// boundaries
	o.boundaries()
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

	// limits and boundaries
	o.limits()
	o.boundaries()
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

	// limits and boundaries
	o.limits()
	o.boundaries()
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
				plt.DrawArrow2d(M.X, M.CovG0, true, scale, argsG0)
				plt.DrawArrow2d(M.X, M.CovG1, true, scale, argsG1)
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
				plt.DrawArrow3d(M.X, M.CovG0, true, scale, argsG0)
				plt.DrawArrow3d(M.X, M.CovG1, true, scale, argsG1)
				plt.DrawArrow3d(M.X, M.CovG2, true, scale, argsG2)
			}
		}
	}
}

// accessors //////////////////////////////////////////////////////////////////////////////////////

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

// Xlength returns the lengths along each direction (whole box) == Xmax(idim) - Xmin(idim)
func (o *CurvGrid) Xlength(idim int) float64 {
	return o.xmax[idim] - o.xmin[idim]
}

// Meshgrid2d extracts 2D meshgrid
//  X -- x0[ny][nx]
//  Y -- x1[ny][nx]
func (o *CurvGrid) Meshgrid2d() (X, Y [][]float64) {
	X = utl.Alloc(o.npts[1], o.npts[0])
	Y = utl.Alloc(o.npts[1], o.npts[0])
	p := 0
	for n := 0; n < o.npts[1]; n++ {
		for m := 0; m < o.npts[0]; m++ {
			X[n][m] = o.mtr[p][n][m].X[0]
			Y[n][m] = o.mtr[p][n][m].X[1]
		}
	}
	return
}

// Meshgrid3d extracts 3D meshgrid
//  X -- x0[nz][ny][nx]
//  Y -- x1[nz][ny][nx]
//  Z -- x2[nz][ny][nx]
func (o *CurvGrid) Meshgrid3d() (X, Y, Z [][][]float64) {
	X = utl.Deep3alloc(o.npts[2], o.npts[1], o.npts[0])
	Y = utl.Deep3alloc(o.npts[2], o.npts[1], o.npts[0])
	Z = utl.Deep3alloc(o.npts[2], o.npts[1], o.npts[0])
	for p := 0; p < o.npts[2]; p++ {
		for n := 0; n < o.npts[1]; n++ {
			for m := 0; m < o.npts[0]; m++ {
				X[p][n][m] = o.mtr[p][n][m].X[0]
				Y[p][n][m] = o.mtr[p][n][m].X[1]
				Z[p][n][m] = o.mtr[p][n][m].X[2]
			}
		}
	}
	return
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

// access nodes ////////////////////////////////////////////////////////////////////////////////////

// IndexN converts triplet node index (m,n,p) into node index N
//
//  2D:   N = m + n⋅n0
//        m = N % n0
//        n = N / n0
//
//  3D:   N = m + n⋅n0 + p⋅n0⋅n1
//        p = N / (n0⋅n1)
//        t = N % (n0⋅n1)  (projection @ z=0)
//        m = t % n0
//        n = t / n0
//
func (o *CurvGrid) IndexN(m, n, p int) (N int) {
	if o.ndim == 2 {
		return m + n*o.npts[0]
	}
	return m + n*o.npts[0] + p*o.npts[0]*o.npts[1]
}

// IndexMNP converts node index N into triplet node index (m,n,p)
//
//  2D:   N = m + n⋅n0
//        m = N % n0
//        n = N / n0
//
//  3D:   N = m + n⋅n0 + p⋅n0⋅n1
//        p = N / (n0⋅n1)
//        t = N % (n0⋅n1)  (projection @ z=0)
//        m = t % n0
//        n = t / n0
//
func (o *CurvGrid) IndexMNP(N int) (m, n, p int) {
	if o.ndim == 2 {
		m = N % o.npts[0]
		n = N / o.npts[0]
		return
	}
	p = N / (o.npts[0] * o.npts[1])
	t := N % (o.npts[0] * o.npts[1])
	m = t % o.npts[0]
	n = t / o.npts[0]
	return
}

// Node returns the physical coordinates of node N (see IndexN()) [may be used to change X]
func (o *CurvGrid) Node(N int) (x la.Vector) {
	m, n, p := o.IndexMNP(N)
	return o.mtr[p][n][m].X
}

// boundaries and tags /////////////////////////////////////////////////////////////////////////////

// Edge returns the ids of points on edges: [edge0, edge1, edge2, edge3]
//            2
//      +-----------+
//      |           |
//      |           |
//     3|           |1
//      |           |
//      |           |
//      +-----------+
//            0
func (o *CurvGrid) Edge(iEdge int) []int {
	return o.edge[iEdge]
}

// EdgeT returns a list of nodes marked with given tag
//           21
//      +-----------+
//      |           |
//      |           |
//    10|           |11
//      |           |
//      |           |
//      +-----------+
//           20
//
//   NOTE: will return empty list if tag is not available
//
func (o *CurvGrid) EdgeT(tag int) []int {
	switch tag {
	case 20:
		return o.edge[0]
	case 11:
		return o.edge[1]
	case 21:
		return o.edge[2]
	case 10:
		return o.edge[3]
	}
	return nil
}

// Face returns the ids of points on faces: [face0, face1, face2, face3, face4, face5]
//               +----------------+
//             ,'|              ,'|
//           ,'  |  ___       ,'  |
//         ,'    |,'5,' [0] ,'    |
//       ,'      |~~~     ,'  ,   |
//     +'===============+'  ,'|   |
//     |   ,'|   |      |   |3|   |
//     |   |2|   |      |   |,'   |
//     |   |,'   +- - - | +- - - -+
//     |   '   ,'       |       ,'
//     |     ,' [1]  ___|     ,'
//     |   ,'      ,'4,'|   ,'
//     | ,'        ~~~  | ,'
//     +----------------+'
func (o *CurvGrid) Face(iFace int) []int {
	return o.face[iFace]
}

// FaceT returns a list of nodes marked with given tag
//               +----------------+
//             ,'|              ,'|
//           ,'  |  ___       ,'  |
//         ,'    |,'31'  10 ,'    |
//       ,'      |~~~     ,'  ,,  |
//     +'===============+'  ,' |  |
//     |   ,'|   |      |   |21|  |
//     |  |20|   |      |   |,'   |
//     |  | ,'   +- - - | +- - - -+
//     |   '   ,'       |       ,'
//     |     ,' 11   ___|     ,'
//     |   ,'      ,'30'|   ,'
//     | ,'        ~~~  | ,'
//     +----------------+'
//
//   NOTE: will return empty list if tag is not available
//
func (o *CurvGrid) FaceT(tag int) []int {
	switch tag {
	case 100:
		return o.face[0]
	case 101:
		return o.face[1]
	case 200:
		return o.face[2]
	case 201:
		return o.face[3]
	case 300:
		return o.face[4]
	case 301:
		return o.face[5]
	}
	return nil
}

// Boundary returns list of edge or face nodes on boundary
//   NOTE: will return empty list if tag is not available
func (o *CurvGrid) Boundary(tag int) []int {
	if tag > 50 {
		if o.ndim == 2 {
			return nil
		}
		return o.FaceT(tag)
	}
	return o.EdgeT(tag)
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

// boundaries generates the IDs of nodes on the boundaries of a rectangular grid
func (o *CurvGrid) boundaries() {
	n0 := o.npts[0]
	n1 := o.npts[1]
	if o.ndim == 2 {
		o.edge = make([][]int, 4)   // bottom, right, top, left
		o.edge[0] = make([]int, n0) // bottom
		o.edge[1] = make([]int, n1) // right
		o.edge[2] = make([]int, n0) // top
		o.edge[3] = make([]int, n1) // left
		for m := 0; m < n0; m++ {
			o.edge[0][m] = m             // bottom
			o.edge[2][m] = m + n0*(n1-1) // top
		}
		for n := 0; n < n1; n++ {
			o.edge[1][n] = n*n0 + n0 - 1 // right
			o.edge[3][n] = n * n0        // left
		}
		return
	}
	n2 := o.npts[2]
	o.face = make([][]int, 6)
	o.face[0] = make([]int, n1*n2) // xmin[0]
	o.face[1] = make([]int, n1*n2) // xmax[0]
	o.face[2] = make([]int, n0*n2) // xmin[1]
	o.face[3] = make([]int, n0*n2) // xmax[1]
	o.face[4] = make([]int, n0*n1) // xmin[2]
	o.face[5] = make([]int, n0*n1) // xmax[2]
	t := 0
	for n := 0; n < n1; n++ {
		for m := 0; m < n0; m++ {
			o.face[4][t] = m + n0*n                  // xmin[2]
			o.face[5][t] = m + n0*n + (n0*n1)*(n2-1) // xmax[2]
			t++
		}
	}
	t = 0
	for p := 0; p < n2; p++ {
		for m := 0; m < n0; m++ {
			o.face[2][t] = m + (n0*n1)*p             // xmin[1]
			o.face[3][t] = m + (n0*n1)*p + n0*(n1-1) // xmax[1]
			t++
		}
	}
	t = 0
	for p := 0; p < n2; p++ {
		for n := 0; n < n1; n++ {
			o.face[0][t] = n*n0 + (n0*n1)*p            // xmin[0]
			o.face[1][t] = n*n0 + (n0*n1)*p + (n0 - 1) // xmax[0]
			t++
		}
	}
	return
}
