// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// Nurbs holds NURBS data
//
//  NOTE: (1) Control points must be set after a call to Init
//        (2) Either SetControl must be called or the Q array must be directly specified
//
//  Reference:
//   [1] Piegl L and Tiller W (1995) The NURBS book, Springer, 646p
type Nurbs struct {

	// essential
	gnd int             // 1: curve, 2:surface, 3:volume (geometry dimension)
	p   []int           // orders [3]
	b   []*Bspline      // B-splines [gnd]
	n   []int           // number of basis functions along each direction [3]
	Q   [][][][]float64 // Qw: weighted control points and weights [n[0]][n[1]][n[2]][4] (Piegl p120)
	l2i [][]int         // local ids of ctrl points => indices of basis functions [nctrl][3]

	// auxiliary
	span []int           // spans computed by CalcBasis or CalcBasisAndDerivs for given u [3] (from CalcBasis/CalcBasisAndDerivs to GetBasis{I,L}/GetDeriv{I,L})
	idx  []int           // buffer to hold indices of non-zero basis functions [3]
	cw   []float64       // point on 4D entity [4]
	rr   [][][]float64   // non-zero basis functions [p0+1][p1+1][p2+1]
	drr  [][][][]float64 // derivatives of non-zero basis functions w.r.t u [p0+1][p1+1][p2+1][gnd]
	dww  []float64       // derivative of W function w.r.t u [gnd]
}

// initialisation methods ////////////////////////////////////////////////////////////////////////////

// NewNurbs returns a new Nurbs object
func NewNurbs(gnd int, ords []int, knots [][]float64) (o *Nurbs) {

	// essential
	o = new(Nurbs)
	o.gnd = gnd
	o.p = make([]int, 3)

	// B-splines
	o.b = make([]*Bspline, o.gnd)
	o.n = make([]int, 3)
	for d := 0; d < o.gnd; d++ {
		o.p[d] = ords[d]
		o.b[d] = NewBspline(knots[d], o.p[d])
		o.n[d] = o.b[d].NumBasis()
		if o.n[d] < 2 {
			chk.Panic("number of knots is incorrect for dimension %d. n == %d is invalid", d, o.n[d])
		}
	}
	for d := o.gnd; d < 3; d++ {
		o.n[d] = 1
	}

	// ids of control points
	nctrl := o.n[0] * o.n[1] * o.n[2]
	o.l2i = make([][]int, nctrl)
	for l := 0; l < nctrl; l++ {
		o.l2i[l] = make([]int, 3)
		switch o.gnd {
		case 1:
			o.l2i[l][0] = l
		case 2:
			o.l2i[l][0] = l % o.n[0] // i
			o.l2i[l][1] = l / o.n[0] // j
		case 3:
			c := l % (o.n[0] * o.n[1])
			o.l2i[l][0] = c % o.n[0]            // i
			o.l2i[l][1] = c / o.n[0]            // j
			o.l2i[l][2] = l / (o.n[0] * o.n[1]) // k
		}
	}

	// auxiliary
	o.span = make([]int, 3)
	o.idx = make([]int, 3)
	o.cw = make([]float64, 4)
	o.rr = utl.Deep3alloc(o.p[0]+1, o.p[1]+1, o.p[2]+1)
	o.drr = utl.Deep4alloc(o.p[0]+1, o.p[1]+1, o.p[2]+1, o.gnd)
	o.dww = make([]float64, o.gnd)
	return
}

// SetControl sets control points from list of global vertices
func (o *Nurbs) SetControl(verts [][]float64, ctrls []int) {

	// check
	nctrl := o.n[0] * o.n[1] * o.n[2]
	if nctrl != len(ctrls) {
		chk.Panic("number of indices of control points must be equal to %d. len(ctrls)=%d is incorrect.\n", nctrl, len(ctrls))
	}
	if len(verts) < 2 {
		chk.Panic("number of vertices must be greater than 1")
	}
	if len(verts[0]) != 4 {
		chk.Panic("number of components of a control point must be 4 (x,y,z and weight), even if x and y are not used; e.g. curve and surface, respectively.\n")
	}

	// set control points
	o.Q = utl.Deep4alloc(o.n[0], o.n[1], o.n[2], 4)
	for i := 0; i < o.n[0]; i++ {
		for j := 0; j < o.n[1]; j++ {
			for k := 0; k < o.n[2]; k++ {
				l := i + j*o.n[0] + k*o.n[0]*o.n[1]
				g := ctrls[l]
				for e := 0; e < 3; e++ {
					o.Q[i][j][k][e] = verts[g][e] * verts[g][3]
				}
				o.Q[i][j][k][3] = verts[g][3]
			}
		}
	}
}

// essential methods /////////////////////////////////////////////////////////////////////////////////

// CalcBasis computes all non-zero basis functions R[i][j][k] @ u.
// Note: use GetBasisI or GetBasisL to get a particular basis function value
func (o *Nurbs) CalcBasis(u []float64) {
	for d := 0; d < o.gnd; d++ {
		o.span[d] = o.b[d].findSpan(u[d])
		o.b[d].basisFuns(u[d], o.span[d])
	}
	α, β, γ := o.p[0], o.p[1], o.p[2]                               // B-spline orders
	x, y, z := o.span[0]-o.p[0], o.span[1]-o.p[1], o.span[2]-o.p[2] // auxiliary indices
	var coef, ww float64
	switch o.gnd {
	// curve
	case 1:
		j, k := 0, 0
		for i := 0; i <= α; i++ {
			coef = o.b[0].ndu[i][α] * o.Q[x+i][j][k][3]
			o.rr[i][j][k] = coef
			ww += coef
		}
		for i := 0; i <= α; i++ {
			o.rr[i][j][k] /= ww
		}
	// surface
	case 2:
		k := 0
		for i := 0; i <= α; i++ {
			for j := 0; j <= β; j++ {
				coef = o.b[0].ndu[i][α] * o.b[1].ndu[j][β] * o.Q[x+i][y+j][k][3]
				o.rr[i][j][k] = coef
				ww += coef
			}
		}
		for i := 0; i <= α; i++ {
			for j := 0; j <= β; j++ {
				o.rr[i][j][k] /= ww
			}
		}
	// volume
	case 3:
		for i := 0; i <= α; i++ {
			for j := 0; j <= β; j++ {
				for k := 0; k <= γ; k++ {
					coef = o.b[0].ndu[i][α] * o.b[1].ndu[j][β] * o.b[2].ndu[k][γ] * o.Q[x+i][y+j][z+k][3]
					o.rr[i][j][k] = coef
					ww += coef
				}
			}
		}
		for i := 0; i <= α; i++ {
			for j := 0; j <= β; j++ {
				for k := 0; k <= γ; k++ {
					o.rr[i][j][k] /= ww
				}
			}
		}
	}
}

// CalcBasisAndDerivs computes all non-zero basis functions R[i][j][k] and corresponding
// first order derivatives of basis functions w.r.t u => dRdu[i][j][k] @ u
// Note: use GetBasisI or GetBasisL to get a particular basis function value
//       use GetDerivI or GetDerivL to get a particular derivative
func (o *Nurbs) CalcBasisAndDerivs(u []float64) {
	π := 1 // derivative order
	for d := 0; d < o.gnd; d++ {
		o.span[d] = o.b[d].findSpan(u[d])
		o.b[d].dersBasisFuns(u[d], o.span[d], π)
	}
	α, β, γ := o.p[0], o.p[1], o.p[2]                               // B-spline orders
	x, y, z := o.span[0]-o.p[0], o.span[1]-o.p[1], o.span[2]-o.p[2] // auxiliary indices
	var coef, ww, ww2 float64
	for d := 0; d < o.gnd; d++ {
		o.dww[d] = 0
	}
	switch o.gnd {
	// curve
	case 1:
		j, k := 0, 0
		for i := 0; i <= α; i++ {
			coef = o.b[0].ndu[i][α] * o.Q[x+i][j][k][3]
			o.rr[i][j][k] = coef
			ww += coef
			o.dww[0] += o.b[0].der[π][i] * o.Q[x+i][j][k][3]
		}
		ww2 = ww * ww
		for i := 0; i <= α; i++ {
			o.rr[i][j][k] /= ww
			o.drr[i][j][k][0] = o.Q[x+i][j][k][3] * (ww*o.b[0].der[π][i] - o.b[0].ndu[i][α]*o.dww[0]) / ww2
		}
	// surface
	case 2:
		k := 0
		for i := 0; i <= α; i++ {
			for j := 0; j <= β; j++ {
				coef = o.b[0].ndu[i][α] * o.b[1].ndu[j][β] * o.Q[x+i][y+j][k][3]
				o.rr[i][j][k] = coef
				ww += coef
				o.dww[0] += o.b[0].der[π][i] * o.b[1].ndu[j][β] * o.Q[x+i][y+j][k][3]
				o.dww[1] += o.b[0].ndu[i][α] * o.b[1].der[π][j] * o.Q[x+i][y+j][k][3]
			}
		}
		ww2 = ww * ww
		for i := 0; i <= α; i++ {
			for j := 0; j <= β; j++ {
				o.rr[i][j][k] /= ww
				o.drr[i][j][k][0] = o.b[1].ndu[j][β] * o.Q[x+i][y+j][k][3] * (ww*o.b[0].der[π][i] - o.b[0].ndu[i][α]*o.dww[0]) / ww2
				o.drr[i][j][k][1] = o.b[0].ndu[i][α] * o.Q[x+i][y+j][k][3] * (ww*o.b[1].der[π][j] - o.b[1].ndu[j][β]*o.dww[1]) / ww2
			}
		}
	// volume
	case 3:
		for i := 0; i <= α; i++ {
			for j := 0; j <= β; j++ {
				for k := 0; k <= γ; k++ {
					coef = o.b[0].ndu[i][α] * o.b[1].ndu[j][β] * o.b[2].ndu[k][γ] * o.Q[x+i][y+j][z+k][3]
					o.rr[i][j][k] = coef
					ww += coef
					o.dww[0] += o.b[0].der[π][i] * o.b[1].ndu[j][β] * o.b[2].ndu[k][γ] * o.Q[x+i][y+j][z+k][3]
					o.dww[1] += o.b[0].ndu[i][α] * o.b[1].der[π][j] * o.b[2].ndu[k][γ] * o.Q[x+i][y+j][z+k][3]
					o.dww[2] += o.b[0].ndu[i][α] * o.b[1].ndu[j][β] * o.b[2].der[π][k] * o.Q[x+i][y+j][z+k][3]
				}
			}
		}
		ww2 = ww * ww
		for i := 0; i <= α; i++ {
			for j := 0; j <= β; j++ {
				for k := 0; k <= γ; k++ {
					o.rr[i][j][k] /= ww
					o.drr[i][j][k][0] = o.b[1].ndu[j][β] * o.b[2].ndu[k][γ] * o.Q[x+i][y+j][z+k][3] * (ww*o.b[0].der[π][i] - o.b[0].ndu[i][α]*o.dww[0]) / ww2
					o.drr[i][j][k][1] = o.b[0].ndu[i][α] * o.b[2].ndu[k][γ] * o.Q[x+i][y+j][z+k][3] * (ww*o.b[1].der[π][j] - o.b[1].ndu[j][β]*o.dww[1]) / ww2
					o.drr[i][j][k][2] = o.b[0].ndu[i][α] * o.b[1].ndu[j][β] * o.Q[x+i][y+j][z+k][3] * (ww*o.b[2].der[π][k] - o.b[2].ndu[k][γ]*o.dww[2]) / ww2
				}
			}
		}
	}
}

// GetBasisI returns the basis function R[i][j][k] just computed
// by CalcBasis or CalcBasisAndDerivs
// Note: I = [i,j,k]
func (o *Nurbs) GetBasisI(I []int) float64 {
	for d := 0; d < o.gnd; d++ {
		o.idx[d] = I[d] + o.p[d] - o.span[d]
		if o.idx[d] < 0 || o.idx[d] > o.p[d] {
			return 0
		}
	}
	return o.rr[o.idx[0]][o.idx[1]][o.idx[2]]
}

// GetBasisL returns the basis function R[i][j][k] just computed
// by CalcBasis or CalcBasisAndDerivs
// Note: l = i + j * n0 + k * n0 * n1
func (o *Nurbs) GetBasisL(l int) float64 {
	return o.GetBasisI(o.l2i[l])
}

// GetDerivI returns the derivative of basis function dR[i][j][k]du just computed
// by CalcBasisAndDerivs
// Note: I = [i,j,k]
func (o *Nurbs) GetDerivI(dRdu []float64, I []int) {
	var outside bool
	for d := 0; d < o.gnd; d++ {
		o.idx[d] = I[d] + o.p[d] - o.span[d]
		if o.idx[d] < 0 || o.idx[d] > o.p[d] {
			outside = true
		}
		dRdu[d] = 0
	}
	if outside {
		return
	}
	for d := 0; d < o.gnd; d++ {
		dRdu[d] = o.drr[o.idx[0]][o.idx[1]][o.idx[2]][d]
	}
}

// GetDerivL returns the derivative of basis function dR[i][j][k]du just computed
// by CalcBasisAndDerivs
// Note: l = i + j * n0 + k * n0 * n1
func (o *Nurbs) GetDerivL(dRdu []float64, l int) {
	o.GetDerivI(dRdu, o.l2i[l])
}

// RecursiveBasis implements basis functions by means of summing all terms
// in Bernstein polynomial using recursive Cox-DeBoor formula (very not efficient)
// Note: l = i + j * n0 + k * n0 * n1
func (o *Nurbs) RecursiveBasis(u []float64, l int) (res float64) {
	I := o.l2i[l]
	var den float64
	switch o.gnd {
	// curve
	case 1:
		j, k := 0, 0
		for i := 0; i < o.n[0]; i++ {
			den += o.b[0].RecursiveBasis(u[0], i) * o.Q[i][j][k][3]
		}
		if math.Abs(den) < 1e-14 {
			chk.Panic("denominator is zero (%v) @ %v for point %d", den, u, l)
		}
		res = o.b[0].RecursiveBasis(u[0], I[0]) * o.Q[I[0]][j][k][3] / den
	// surface
	case 2:
		k := 0
		for i := 0; i < o.n[0]; i++ {
			for j := 0; j < o.n[1]; j++ {
				den += o.b[0].RecursiveBasis(u[0], i) * o.b[1].RecursiveBasis(u[1], j) * o.Q[i][j][k][3]
			}
		}
		if math.Abs(den) < 1e-14 {
			chk.Panic("denominator is zero (%v) @ %v for point %d", den, u, l)
		}
		res = o.b[0].RecursiveBasis(u[0], I[0]) * o.b[1].RecursiveBasis(u[1], I[1]) * o.Q[I[0]][I[1]][k][3] / den
	// volume
	case 3:
		for i := 0; i < o.n[0]; i++ {
			for j := 0; j < o.n[1]; j++ {
				for k := 0; k < o.n[2]; k++ {
					den += o.b[0].RecursiveBasis(u[0], i) * o.b[1].RecursiveBasis(u[1], j) * o.b[2].RecursiveBasis(u[2], j) * o.Q[i][j][k][3]
				}
			}
		}
		if math.Abs(den) < 1e-14 {
			chk.Panic("denominator is zero (%v) @ %v for point %d", den, u, l)
		}
		res = o.b[0].RecursiveBasis(u[0], I[0]) * o.b[1].RecursiveBasis(u[1], I[1]) * o.b[2].RecursiveBasis(u[2], I[2]) * o.Q[I[0]][I[1]][I[2]][3] / den
	}
	return
}

// Point returns the x-y-z coordinates of a point on curve/surface/volume
//   Input:
//     u    -- [gnd] knot values
//     ndim -- the dimension of the point. E.g. allows drawing curves in 3D
//   Output:
//     C -- [ndim] point coordinates
//   NOTE: Algorithm A4.1 (p124) of [1]
func (o *Nurbs) Point(C, u []float64, ndim int) {
	for d := 0; d < o.gnd; d++ {
		o.span[d] = o.b[d].findSpan(u[d])
		o.b[d].basisFuns(u[d], o.span[d])
		o.idx[d] = o.span[d] - o.p[d]
	}
	for e := 0; e < 4; e++ {
		o.cw[e] = 0
	}
	switch o.gnd {
	// curve
	case 1:
		j, k := 0, 0
		for i := 0; i <= o.p[0]; i++ {
			for e := 0; e < 4; e++ {
				o.cw[e] += o.b[0].ndu[i][o.p[0]] * o.Q[o.idx[0]+i][j][k][e]
			}
		}
	// surface
	case 2:
		k := 0
		for i := 0; i <= o.p[0]; i++ {
			for j := 0; j <= o.p[1]; j++ {
				for e := 0; e < 4; e++ {
					o.cw[e] += o.b[0].ndu[i][o.p[0]] * o.b[1].ndu[j][o.p[1]] * o.Q[o.idx[0]+i][o.idx[1]+j][k][e]
				}
			}
		}
	// volume
	case 3:
		for i := 0; i <= o.p[0]; i++ {
			for j := 0; j <= o.p[1]; j++ {
				for k := 0; k <= o.p[2]; k++ {
					for e := 0; e < 4; e++ {
						o.cw[e] += o.b[0].ndu[i][o.p[0]] * o.b[1].ndu[j][o.p[1]] * o.b[2].ndu[k][o.p[2]] * o.Q[o.idx[0]+i][o.idx[1]+j][o.idx[2]+k][e]
					}
				}
			}
		}
	}
	for e := 0; e < ndim; e++ {
		C[e] = o.cw[e] / o.cw[3]
	}
	return
}

// PointAndFirstDerivs returns the point and first order derivatives with respect to the knot values u
// of the x-y-z coordinates of a point on curve/surface/volume
//   Input:
//     u    -- [gnd] knot values
//     ndim -- the dimension of the point. E.g. allows drawing curves in 3D
//   Output:
//     dCdu -- [ndim][gnd] derivatives dC_i/du_j
//     C    -- [ndim] point coordinates
func (o *Nurbs) PointAndFirstDerivs(dCdu *la.Matrix, C, u []float64, ndim int) {
	for d := 0; d < o.gnd; d++ {
		o.span[d] = o.b[d].findSpan(u[d])
		o.b[d].basisFuns(u[d], o.span[d])
		o.idx[d] = o.span[d] - o.p[d]
	}
	for e := 0; e < 4; e++ {
		o.cw[e] = 0
	}
	dcwdu := utl.Alloc(o.gnd, 4)
	switch o.gnd {
	// curve
	case 1:
		j, k := 0, 0
		for i := 0; i <= o.p[0]; i++ {
			for e := 0; e < 4; e++ {
				o.cw[e] += o.b[0].ndu[i][o.p[0]] * o.Q[o.idx[0]+i][j][k][e]
				if i < o.p[0] {
					num := o.Q[o.idx[0]+i+1][j][k][e] - o.Q[o.idx[0]+i][j][k][e]
					den := o.b[0].T[o.idx[0]+i+1+o.p[0]] - o.b[0].T[o.idx[0]+i+1]
					dcwdu[0][e] += o.b[0].ndu[i][o.p[0]-1] * float64(o.p[0]) * num / den
				}
			}
		}
	// surface
	case 2:
		k := 0
		for i := 0; i <= o.p[0]; i++ {
			for j := 0; j <= o.p[1]; j++ {
				for e := 0; e < 4; e++ {
					o.cw[e] += o.b[0].ndu[i][o.p[0]] * o.b[1].ndu[j][o.p[1]] * o.Q[o.idx[0]+i][o.idx[1]+j][k][e]
					if i < o.p[0] {
						num := o.Q[o.idx[0]+i+1][o.idx[1]+j][k][e] - o.Q[o.idx[0]+i][o.idx[1]+j][k][e]
						den := o.b[0].T[o.idx[0]+i+1+o.p[0]] - o.b[0].T[o.idx[0]+i+1]
						dcwdu[0][e] += o.b[0].ndu[i][o.p[0]-1] * o.b[1].ndu[j][o.p[1]] * float64(o.p[0]) * num / den
					}
					if j < o.p[1] {
						num := o.Q[o.idx[0]+i][o.idx[1]+j+1][k][e] - o.Q[o.idx[0]+i][o.idx[1]+j][k][e]
						den := o.b[1].T[o.idx[1]+j+1+o.p[1]] - o.b[1].T[o.idx[1]+j+1]
						dcwdu[1][e] += o.b[0].ndu[i][o.p[0]] * o.b[1].ndu[j][o.p[1]-1] * float64(o.p[1]) * num / den
					}
				}
			}
		}
	// volume
	case 3:
		chk.Panic("PointAndFirstDerivs of volume is not available yet\n")
	}
	for e := 0; e < ndim; e++ {
		C[e] = o.cw[e] / o.cw[3]
		for d := 0; d < o.gnd; d++ {
			dCdu.Set(e, d, (dcwdu[d][e]-dcwdu[d][3]*C[e])/o.cw[3])
		}
	}
	return
}

// PointAndDerivs computes position and the first and second order derivatives
// Using Algorithms A3.2(p93), A3.6(p111), A4.2(p127), and A4.4(p137)
//   Input:
//     u    -- knot values {r,s,t} [gnd]
//     ndim -- the dimension of the point. E.g. allows drawing curves in 3D
//   Output:
//     x      -- position {x,y,z} (the same as the C varible in [1])
//     dxDr   -- ∂{x}/∂r
//     dxDs   -- ∂{x}/∂s    [may be nil] (volume and surfaces)
//     dxDt   -- ∂{x}/∂t    [may be nil] (volume)
//     ddxDrr -- ∂²{x}/∂r²
//     ddxDss -- ∂²{x}/∂s²  [may be nil] (volume and surfaces)
//     ddxDtt -- ∂²{x}/∂t²  [may be nil] (volume)
//     ddxDrs -- ∂²{x}/∂r∂s [may be nil] (volume and surfaces)
//     ddxDrt -- ∂²{x}/∂r∂t [may be nil] (volume)
//     ddxDst -- ∂²{x}/∂s∂t [may be nil] (volume)
func (o *Nurbs) PointAndDerivs(x, dxDr, dxDs, dxDt,
	ddxDrr, ddxDss, ddxDtt, ddxDrs, ddxDrt, ddxDst, u la.Vector, ndim int) {

	// find span and ctrl indices, and compute basis functions and their derivatives
	upto := 2
	for d := 0; d < o.gnd; d++ {
		o.span[d] = o.b[d].findSpan(u[d])
		o.b[d].dersBasisFuns(u[d], o.span[d], upto)
		o.idx[d] = o.span[d] - o.p[d]
	}
	for e := 0; e < 4; e++ {
		o.cw[e] = 0
	}

	// derivative of augmented homogeneous coordinates Cw
	dcwdr := la.NewVector(4)
	ddcwdrr := la.NewVector(4)

	switch o.gnd {

	// curve
	case 1:
		j, k := 0, 0
		for e := 0; e < 4; e++ { // for each homogeneous component
			for i := 0; i <= o.p[0]; i++ { // summing over i
				o.cw[e] += o.b[0].ndu[i][o.p[0]] * o.Q[o.idx[0]+i][j][k][e]
				dcwdr[e] += o.b[0].der[1][i] * o.Q[o.idx[0]+i][j][k][e]
				ddcwdrr[e] += o.b[0].der[2][i] * o.Q[o.idx[0]+i][j][k][e]
			}
		}

	// surface
	case 2:
		chk.Panic("PointAndDerivs of surface is not available yet\n")

	// volume
	case 3:
		chk.Panic("PointAndDerivs of volume is not available yet\n")
	}

	// correct values
	for d := 0; d < ndim; d++ {
		x[d] = o.cw[d] / o.cw[3]
		dxDr[d] = (dcwdr[d] - dcwdr[3]*x[d]) / o.cw[3]
		if ddxDrr != nil {
			ddxDrr[d] = (ddcwdrr[d] - 2.0*dcwdr[3]*dxDr[d] - ddcwdrr[3]*x[d]) / o.cw[3]
		}
	}
}

// accessors methods /////////////////////////////////////////////////////////////////////////////////

// Gnd returns the geometry dimension
func (o *Nurbs) Gnd() int {
	return o.gnd
}

// Ord returns the order along direction dir
func (o *Nurbs) Ord(dir int) int {
	return o.p[dir]
}

// GetU returns the knots along direction dir
func (o *Nurbs) GetU(dir int) []float64 {
	return o.b[dir].T
}

// U returns the value of a knot along direction dir
func (o *Nurbs) U(dir, idx int) float64 {
	return o.b[dir].T[idx]
}

// NumBasis returns the number of basis (controls) along direction dir
func (o *Nurbs) NumBasis(dir int) int {
	return o.b[dir].NumBasis()
}

// NonZeroSpans returns the 'elements' along direction dir
func (o *Nurbs) NonZeroSpans(dir int) [][]int {
	return o.b[dir].Elements()
}

// GetQ gets a control point x[i,j,k] (size==4)
func (o *Nurbs) GetQ(i, j, k int) (x []float64) {
	x = make([]float64, 4)
	x[3] = o.Q[i][j][k][3]
	for d := 0; d < 3; d++ {
		x[d] = o.Q[i][j][k][d] / x[3]
	}
	return
}

// GetQl gets a control point x[l] (size==4)
func (o *Nurbs) GetQl(l int) (x []float64) {
	return o.GetQ(o.l2i[l][0], o.l2i[l][1], o.l2i[l][2])
}

// SetQ sets a control point x[i,j,k] (size==4)
func (o *Nurbs) SetQ(i, j, k int, x []float64) {
	for d := 0; d < 3; d++ {
		o.Q[i][j][k][d] = x[d] * x[3]
	}
	o.Q[i][j][k][3] = x[3]
}

// SetQl sets a control point x[l] (size==4)
func (o *Nurbs) SetQl(l int, x []float64) {
	o.SetQ(o.l2i[l][0], o.l2i[l][1], o.l2i[l][2], x)
}

// Elements returns the indices of nonzero spans
func (o *Nurbs) Elements() (spans [][]int) {
	switch o.gnd {
	// curve
	case 1:
		s0 := o.b[0].Elements()
		n0 := len(s0)
		spans = make([][]int, n0)
		for i := 0; i < n0; i++ {
			e := i
			spans[e] = []int{s0[i][0], s0[i][1]}
		}
	// surface
	case 2:
		s0, s1 := o.b[0].Elements(), o.b[1].Elements()
		n0, n1 := len(s0), len(s1)
		spans = make([][]int, n0*n1)
		for j := 0; j < n1; j++ {
			for i := 0; i < n0; i++ {
				e := i + j*n0
				spans[e] = []int{s0[i][0], s0[i][1], s1[j][0], s1[j][1]}
			}
		}
	// volume
	case 3:
		s0, s1, s2 := o.b[0].Elements(), o.b[1].Elements(), o.b[2].Elements()
		n0, n1, n2 := len(s0), len(s1), len(s2)
		spans = make([][]int, n0*n1*n2)
		for k := 0; k < n2; k++ {
			for j := 0; j < n1; j++ {
				for i := 0; i < n0; i++ {
					e := i + j*n0 + k*n0*n1
					spans[e] = []int{s0[i][0], s0[i][1], s1[j][0], s1[j][1], s2[k][0], s2[k][1]}
				}
			}
		}
	}
	return
}

// GetElemNumBasis returns the number of control points == basis functions needed for one element
//  npts := Π_i (p[i] + 1)
func (o *Nurbs) GetElemNumBasis() (npts int) {
	npts = 1
	for i := 0; i < o.gnd; i++ {
		npts *= (o.p[i] + 1)
	}
	return
}

// IndBasis returns the indices of basis functions == local indices of control points
func (o *Nurbs) IndBasis(span []int) (L []int) {
	switch o.gnd {
	// curve
	case 1:
		nbu := o.p[0] + 1 // number of basis functions along u
		L = make([]int, nbu)
		for i := 0; i < nbu; i++ {
			L[i] = span[0] - o.p[0] + i
		}
	// surface
	case 2:
		nbu := o.p[0] + 1 // number of basis functions along u
		nbv := o.p[1] + 1 // number of basis functions along v
		L = make([]int, nbu*nbv)
		c := 0
		for j := 0; j < nbv; j++ {
			J := span[2] - o.p[1] + j
			for i := 0; i < nbu; i++ {
				I := span[0] - o.p[0] + i
				L[c] = I + J*o.n[0]
				c++
			}
		}
	// volume
	case 3:
		nbu := o.p[0] + 1 // number of basis functions along u
		nbv := o.p[1] + 1 // number of basis functions along v
		nbw := o.p[2] + 1 // number of basis functions along w
		L = make([]int, nbu*nbv*nbw)
		c := 0
		for k := 0; k < nbw; k++ {
			K := span[4] - o.p[2] + k
			for j := 0; j < nbv; j++ {
				J := span[2] - o.p[1] + j
				for i := 0; i < nbu; i++ {
					I := span[0] - o.p[0] + i
					L[c] = I + J*o.n[0] + K*o.n[1]*o.n[2]
					c++
				}
			}
		}
	}
	return
}

// GetLimitsQ computes the limits of all coordinates of control points in NURBS
func (o *Nurbs) GetLimitsQ() (xmin, xmax []float64) {
	xmin = []float64{math.Inf(+1), math.Inf(+1), math.Inf(+1)}
	xmax = []float64{math.Inf(-1), math.Inf(-1), math.Inf(-1)}
	for k := 0; k < o.n[2]; k++ {
		for j := 0; j < o.n[1]; j++ {
			for i := 0; i < o.n[0]; i++ {
				x := o.GetQ(i, j, k)
				for r := 0; r < o.gnd; r++ {
					xmin[r] = utl.Min(xmin[r], x[r])
					xmax[r] = utl.Max(xmax[r], x[r])
				}
			}
		}
	}
	for i := o.gnd; i < 3; i++ {
		xmin[i] = 0
		xmax[i] = 0
	}
	return
}

// manipulation methods //////////////////////////////////////////////////////////////////////////////

// KrefineN return a new Nurbs with each span divided into ndiv parts = [2, 3, ...]
func (o *Nurbs) KrefineN(ndiv int, hughesEtAlPaper bool) *Nurbs {
	X := make([][]float64, o.gnd)
	xa := make([]float64, 3)
	xb := make([]float64, 3)
	if hughesEtAlPaper {
		elems := o.Elements()
		switch o.gnd {
		case 2:
			for _, e := range elems {
				umin, umax := o.b[0].T[e[0]], o.b[0].T[e[1]]
				vmin, vmax := o.b[1].T[e[2]], o.b[1].T[e[3]]
				o.Point(xa, []float64{umin, vmin}, 3)
				o.Point(xb, []float64{umax, vmin}, 3)
				xc := []float64{(xa[0] + xb[0]) / 2.0, (xa[1] + xb[1]) / 2.0}

				io.Pf("xa = %v\n", xa)
				io.Pf("xb = %v\n", xb)
				io.Pf("xc = %v\n", xc)

				o.Point(xa, []float64{umin, vmax}, 3)
				o.Point(xb, []float64{umax, vmax}, 3)
				xc = []float64{(xa[0] + xb[0]) / 2.0, (xa[1] + xb[1]) / 2.0}

				io.Pfpink("xa, xb, xc = %v, %v, %v\n", xa, xb, xc)
				chk.Panic("KrefineN with hughesEtAlPaper==true is not implemented in 2D yet")
			}
		case 3:
			chk.Panic("KrefineN with hughesEtAlPaper==true is not implemented in 3D yet")
		}
		io.Pfgrey("KrefineN with hughesEtAlPaper==true => not implemented yet\n")
		return nil
	}
	for d := 0; d < o.gnd; d++ {
		nspans := o.b[d].m - 2*o.p[d] - 1
		nnewk := nspans * (ndiv - 1)
		X[d] = make([]float64, nnewk)
		k := 0
		for i := 0; i < nspans; i++ {
			umin, umax := o.b[d].T[o.p[d]+i], o.b[d].T[o.p[d]+i+1]
			du := (umax - umin) / float64(ndiv)
			for j := 1; j < ndiv; j++ {
				X[d][k] = umin + du*float64(j)
				k++
			}
		}
	}
	return o.Krefine(X)
}

// Krefine returns a new Nurbs with knots refined
// Note: X[gnd][numNewKnots]
func (o *Nurbs) Krefine(X [][]float64) (O *Nurbs) {

	// check
	if len(X) != o.gnd {
		chk.Panic("size of new knots array X must be [gnd==%d][numNewKnots]. len==%d of first dimension is incorrect", o.gnd, len(X))
	}

	// number of new knots and first and last knots
	nk, a, b := make([]int, 3), make([]int, 3), make([]int, 3)
	for d := 0; d < o.gnd; d++ {
		nk[d] = len(X[d])
		a[d] = o.b[d].findSpan(X[d][0])
		b[d] = o.b[d].findSpan(X[d][nk[d]-1])
	}

	// new knots array
	Unew := make([][]float64, o.gnd)
	Unew[0] = make([]float64, o.b[0].m+nk[0])
	if o.gnd > 1 {
		Unew[1] = make([]float64, o.b[1].m+nk[1])
	}
	if o.gnd > 2 {
		Unew[2] = make([]float64, o.b[2].m+nk[2])
	}

	// refine
	var Qnew [][][][]float64
	switch o.gnd {
	case 1:
		Qnew = utl.Deep4alloc(o.n[0]+nk[0], o.n[1], o.n[2], 4)
		krefine(Unew, Qnew, o.Q, 0, X[0], o.b[0].T, o.n, o.p, a[0], b[0])
	case 2:
		// along 0
		Qtmp := utl.Deep4alloc(o.n[0]+nk[0], o.n[1], o.n[2], 4)
		krefine(Unew, Qtmp, o.Q, 0, X[0], o.b[0].T, o.n, o.p, a[0], b[0])
		// along 1
		nn := make([]int, 3)
		copy(nn, o.n)
		nn[0] += nk[0]
		Qnew = utl.Deep4alloc(nn[0], o.n[1]+nk[1], o.n[2], 4)
		krefine(Unew, Qnew, Qtmp, 1, X[1], o.b[1].T, nn, o.p, a[1], b[1])
	case 3:
		chk.Panic("nurbs.go: Krefine: gnd=3 not implemented yet")
	}

	// initialize new nurbs
	O = NewNurbs(o.gnd, o.p, Unew)
	O.Q = Qnew
	return
}

// auxiliary methods /////////////////////////////////////////////////////////////////////////////////

// krefine refines a nurbs with new knots
func krefine(Unew [][]float64, Qnew, Qold [][][][]float64, dir int, X, U []float64, nn, pp []int, a, b int) {

	// auxiliary functions
	getnf := []int{nn[1], nn[0], nn[0]}
	getng := []int{nn[2], nn[2], nn[1]}
	qqnew := func(f, g, i int) []float64 {
		if dir == 0 {
			return Qnew[i][f][g]
		}
		if dir == 1 {
			return Qnew[f][i][g]
		}
		return Qnew[f][g][i]
	}
	qqold := func(f, g, i int) []float64 {
		if dir == 0 {
			return Qold[i][f][g]
		}
		if dir == 1 {
			return Qold[f][i][g]
		}
		return Qold[f][g][i]
	}

	// auxiliary variables
	r := len(X) - 1
	n := nn[dir] - 1
	p := pp[dir]
	m := n + p + 1
	b = b + 1

	// save unaltered control points
	for f := 0; f < getnf[dir]; f++ {
		for g := 0; g < getng[dir]; g++ {
			for d := 0; d < 4; d++ {
				for j := 0; j <= a-p; j++ {
					qqnew(f, g, j)[d] = qqold(f, g, j)[d]
				}
				for j := b - 1; j <= n; j++ {
					qqnew(f, g, j+r+1)[d] = qqold(f, g, j)[d]
				}
			}
		}
	}

	// save unaltered knots
	for j := 0; j <= a; j++ {
		Unew[dir][j] = U[j]
	}
	for j := b + p; j <= m; j++ {
		Unew[dir][j+r+1] = U[j]
	}

	// refine
	i, k := b+p-1, b+p+r
	for j := r; j >= 0; j-- {
		for X[j] <= U[i] && i > a {
			for f := 0; f < getnf[dir]; f++ {
				for g := 0; g < getng[dir]; g++ {
					for d := 0; d < 4; d++ {
						qqnew(f, g, k-p-1)[d] = qqold(f, g, i-p-1)[d]
					}
				}
			}
			Unew[dir][k] = U[i]
			k, i = k-1, i-1
		}
		for f := 0; f < getnf[dir]; f++ {
			for g := 0; g < getng[dir]; g++ {
				for d := 0; d < 4; d++ {
					qqnew(f, g, k-p-1)[d] = qqnew(f, g, k-p)[d]
				}
			}
		}
		for l := 1; l <= p; l++ {
			ind := k - p + l
			alp := Unew[dir][k+l] - X[j]
			if math.Abs(alp) > 1e-14 {
				alp /= Unew[dir][k+l] - U[i-p+l]
				for f := 0; f < getnf[dir]; f++ {
					for g := 0; g < getng[dir]; g++ {
						for d := 0; d < 4; d++ {
							qqnew(f, g, ind-1)[d] = alp*qqnew(f, g, ind-1)[d] + (1.0-alp)*qqnew(f, g, ind)[d]
						}
					}
				}
			} else {
				for f := 0; f < getnf[dir]; f++ {
					for g := 0; g < getng[dir]; g++ {
						for d := 0; d < 4; d++ {
							qqnew(f, g, ind-1)[d] = qqnew(f, g, ind)[d]
						}
					}
				}
			}
		}
		Unew[dir][k] = X[j]
		k--
	}
}
