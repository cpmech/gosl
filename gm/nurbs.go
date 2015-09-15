// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

const (
	ATOL = 1e-14 // zero α for refinement
)

// Nurbs holds NURBS data
// Note: Control points must be set after a call to Init
//       Either SetControl must be called or the Q array must be directly specified
type Nurbs struct {

	// essential
	gnd int             // 1: curve, 2:surface, 3:volume (geometry dimension)
	p   []int           // orders [3]
	b   []Bspline       // B-splines [gnd]
	n   []int           // number of basis functions along each direction [3]
	Q   [][][][]float64 // Qw: 4D control points [n[0]][n[1]][n[2]][4]
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

// Init initialises Nurbs
func (o *Nurbs) Init(gnd int, ords []int, knots [][]float64) {

	// essential
	o.gnd = gnd
	o.p = make([]int, 3)

	// B-splines
	o.b = make([]Bspline, o.gnd)
	o.n = make([]int, 3)
	for d := 0; d < o.gnd; d++ {
		o.p[d] = ords[d]
		o.b[d].Init(knots[d], o.p[d])
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
}

// SetControl sets control points from list of global vertices
func (o *Nurbs) SetControl(verts [][]float64, ctrls []int) {

	// check
	nctrl := o.n[0] * o.n[1] * o.n[2]
	if nctrl != len(ctrls) {
		chk.Panic("number of control points must be equal to %d. nctrl == %d is incorrect", nctrl, len(ctrls))
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
		o.span[d] = o.b[d].find_span(u[d])
		o.b[d].basis_funs(u[d], o.span[d])
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
		o.span[d] = o.b[d].find_span(u[d])
		o.b[d].ders_basis_funs(u[d], o.span[d], π)
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
		if math.Abs(den) < ZTOL {
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
		if math.Abs(den) < ZTOL {
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
		if math.Abs(den) < ZTOL {
			chk.Panic("denominator is zero (%v) @ %v for point %d", den, u, l)
		}
		res = o.b[0].RecursiveBasis(u[0], I[0]) * o.b[1].RecursiveBasis(u[1], I[1]) * o.b[2].RecursiveBasis(u[2], I[2]) * o.Q[I[0]][I[1]][I[2]][3] / den
	}
	return
}

// NumericalDeriv computes a particular derivative dR[i][j][k]du @ t using numerical differentiation
// Note: it uses RecursiveBasis and therefore is highly non-efficient
func (o *Nurbs) NumericalDeriv(dRdu []float64, u []float64, l int) {
	var tmp float64
	for d := 0; d < o.gnd; d++ {
		f := func(x float64, args ...interface{}) (val float64) {
			if x < o.b[d].tmin || x > o.b[d].tmax {
				chk.Panic("problem with numerical derivative: x=%v is invalid. xrange=[%v,%v]", x, o.b[d].tmin, o.b[d].tmax)
			}
			tmp = u[d]
			u[d] = x
			val = o.RecursiveBasis(u, l)
			u[d] = tmp
			return
		}
		dRdu[d] = num.DerivRange(f, u[d], o.b[d].tmin, o.b[d].tmax)
	}
	return
}

// Point returns the x-y-z coordinates of a point on curve/surface/volume
func (o *Nurbs) Point(u []float64) (C []float64) {
	C = make([]float64, 3)
	for d := 0; d < o.gnd; d++ {
		o.span[d] = o.b[d].find_span(u[d])
		o.b[d].basis_funs(u[d], o.span[d])
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
	for e := 0; e < 3; e++ {
		C[e] = o.cw[e] / o.cw[3]
	}
	return
}

// accessors methods /////////////////////////////////////////////////////////////////////////////////

// GetGnd returns the geometry dimension
func (o *Nurbs) Gnd() int {
	return o.gnd
}

// GetOrd returns the order along direction dir
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
				c += 1
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
					c += 1
				}
			}
		}
	}
	return
}

// manipulation methods //////////////////////////////////////////////////////////////////////////////

// KrefineN return a new Nurbs with each span divided into ndiv parts = [2, 3, ...]
func (o *Nurbs) KrefineN(ndiv int, useCspace bool) *Nurbs {
	X := make([][]float64, o.gnd)
	if useCspace {
		elems := o.Elements()
		for _, e := range elems {
			switch o.gnd {
			case 2:
				umin, umax := o.b[0].T[e[0]], o.b[0].T[e[1]]
				vmin, vmax := o.b[1].T[e[2]], o.b[1].T[e[3]]
				xa := o.Point([]float64{umin, vmin})
				xb := o.Point([]float64{umax, vmin})
				xc := []float64{(xa[0] + xb[0]) / 2.0, (xa[1] + xb[1]) / 2.0}
				io.Pforan("xa, xb, xc = %v, %v, %v\n", xa, xb, xc)
				xa = o.Point([]float64{umin, vmax})
				xb = o.Point([]float64{umax, vmax})
				xc = []float64{(xa[0] + xb[0]) / 2.0, (xa[1] + xb[1]) / 2.0}
				io.Pfpink("xa, xb, xc = %v, %v, %v\n", xa, xb, xc)
			}
		}
		chk.Panic("nurbs.go: KrefineN with useCspace==true => not implemented yet")
	} else {
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
					k += 1
				}
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
		a[d] = o.b[d].find_span(X[d][0])
		b[d] = o.b[d].find_span(X[d][nk[d]-1])
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
	O = new(Nurbs)
	O.Init(o.gnd, o.p, Unew)
	O.Q = Qnew
	return
}

// ExtractSurfaces returns a new NURBS representing a boundary of this NURBS
func (o *Nurbs) ExtractSurfaces() (surfs []*Nurbs) {
	if o.gnd == 1 {
		return
	}
	nsurf := o.gnd * 2
	surfs = make([]*Nurbs, nsurf)
	var ords [][]int
	var knots [][][]float64
	if o.gnd == 2 {
		ords = [][]int{
			{o.p[1]}, // perpendicular to x
			{o.p[0]}, // perpendicular to y
		}
		knots = [][][]float64{
			{o.b[1].T}, // perpendicular to x
			{o.b[0].T}, // perpendicular to y
		}
	} else {
		ords = [][]int{
			{o.p[1], o.p[2]}, // perpendicular to x
			{o.p[2], o.p[0]}, // perpendicular to y
			{o.p[0], o.p[1]}, // perpendicular to z
		}
		knots = [][][]float64{
			{o.b[1].T, o.b[2].T}, // perpendicular to x
			{o.b[2].T, o.b[0].T}, // perpendicular to y
			{o.b[0].T, o.b[1].T}, // perpendicular to z
		}
	}
	for i := 0; i < o.gnd; i++ {
		a, b := i*o.gnd, i*o.gnd+1
		surfs[a] = new(Nurbs) // surface perpendicular to i
		surfs[b] = new(Nurbs) // opposite surface perpendicular to i
		surfs[a].Init(o.gnd-1, ords[i], knots[i])
		surfs[b].Init(o.gnd-1, ords[i], knots[i])
		if o.gnd == 2 { // boundary is curve
			j := (i + 1) % o.gnd // direction perpendicular to i
			surfs[a].Q = o.clone_Q_along_curve(j, 0)
			surfs[b].Q = o.clone_Q_along_curve(j, o.n[i]-1)
		} else { // boundary is surface
			j := (i + 1) % o.gnd // direction perpendicular to i
			k := (i + 2) % o.gnd // other direction perpendicular to i
			surfs[a].Q = o.clone_Q_along_surface(j, k, 0)
			surfs[b].Q = o.clone_Q_along_surface(j, k, o.n[i]-1)
		}
	}
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
			if math.Abs(alp) > ATOL {
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
		k -= 1
	}
}

func (o *Nurbs) clone_Q_along_curve(iAlong, jAt int) (Qnew [][][][]float64) {
	Qnew = utl.Deep4alloc(o.n[iAlong], 1, 1, 4)
	var i, j int
	for m := 0; m < o.n[iAlong]; m++ {
		i, j = m, jAt
		if iAlong == 1 {
			i, j = jAt, m
		}
		for e := 0; e < 4; e++ {
			Qnew[m][0][0][e] = o.Q[i][j][0][e]
		}
	}
	return
}

func (o *Nurbs) clone_Q_along_surface(iAlong, jAlong, kAt int) (Qnew [][][][]float64) {
	Qnew = utl.Deep4alloc(o.n[iAlong], o.n[jAlong], 1, 4)
	var i, j, k int
	for m := 0; m < o.n[iAlong]; m++ {
		for n := 0; n < o.n[jAlong]; n++ {
			switch {
			case iAlong == 0 && jAlong == 1:
				i, j, k = m, n, kAt
			case iAlong == 1 && jAlong == 2:
				i, j, k = kAt, m, n
			case iAlong == 2 && jAlong == 0:
				i, j, k = n, kAt, m
			default:
				chk.Panic("clone Q surface is specified by 'along' indices in (0,1) or (1,2) or (2,0). (%d,%d) is incorrect", iAlong, jAlong)
			}
			for e := 0; e < 4; e++ {
				Qnew[m][n][0][e] = o.Q[i][j][k][e]
			}
		}
	}
	return
}

// IndsAlongCurve returns the control points indices along curve
func (o *Nurbs) IndsAlongCurve(iAlong, iSpan0, jAt int) (L []int) {
	nb := o.p[iAlong] + 1 // number of basis along i
	L = make([]int, nb)
	var i, j int
	for m := 0; m < nb; m++ {
		if iAlong == 0 {
			i = iSpan0 - o.p[0] + m
			j = jAt
		} else {
			i = jAt
			j = iSpan0 - o.p[1] + m
		}
		L[m] = i + j*o.n[0]
	}
	return
}

// IndsAlongSurface return the control points indices along surface
func (o *Nurbs) IndsAlongSurface(iAlong, jAlong, iSpan0, jSpan0, kAt int) (L []int) {
	nbu := o.p[iAlong] + 1 // number of basis functions along i
	nbv := o.p[jAlong] + 1 // number of basis functions along j
	L = make([]int, nbu*nbv)
	var c, i, j, k int
	for m := 0; m < nbu; m++ {
		for n := 0; n < nbv; n++ {
			switch {
			case iAlong == 0 && jAlong == 1:
				i = iSpan0 - o.p[0] + m
				j = jSpan0 - o.p[1] + n
				k = kAt
			case iAlong == 1 && jAlong == 2:
				i = kAt
				j = iSpan0 - o.p[1] + m
				k = jSpan0 - o.p[2] + n
			case iAlong == 2 && jAlong == 0:
				i = jSpan0 - o.p[0] + n
				j = kAt
				k = iSpan0 - o.p[2] + m
			}
			L[c] = i + j*o.n[0] + k*o.n[1]*o.n[2]
			c += 1
		}
	}
	return
}

// ElemBryLocalInds returns the local (element) indices of control points @ boundaries
// (if element would have all surfaces @ boundaries)
func (o *Nurbs) ElemBryLocalInds() (I [][]int) {
	switch o.gnd {
	case 1:
		return
	case 2:
		I = make([][]int, 2*o.gnd)
		nx, ny := o.p[0]+1, o.p[1]+1
		I[3] = utl.IntRange3(0, nx*ny, nx)
		I[1] = utl.IntAddScalar(I[3], nx-1)
		I[0] = utl.IntRange(nx)
		I[2] = utl.IntAddScalar(I[0], (ny-1)*nx)
	case 3:
		I = make([][]int, 2*o.gnd)
		chk.Panic("3D NUTBS: ElemBryLocalInds: TODO") // TODO
	}
	return
}
