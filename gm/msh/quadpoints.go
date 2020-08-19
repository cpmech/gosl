// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"

	"gosl/chk"
	"gosl/num"
	"gosl/plt"
	"gosl/utl"
)

// QuadPointsGaussLegendre generate quadrature points for Gauss-Legendre integration
//    npts -- is the total number of points; e.g. 27 for 3D (boxes)
func QuadPointsGaussLegendre(ndim, npts int) (pts [][]float64) {
	n1d := int(math.Floor(math.Pow(float64(npts), 1.0/float64(ndim)) + 0.5))
	x, w := num.GaussLegendreXW(-1, 1, n1d)
	pts = make([][]float64, npts)
	switch ndim {
	case 1:
		for i := 0; i < npts; i++ {
			pts[i] = []float64{x[i], 0, 0, w[i]}
		}
	case 2:
		for j := 0; j < n1d; j++ {
			for i := 0; i < n1d; i++ {
				m := i + n1d*j
				pts[m] = []float64{x[i], x[j], 0, w[i] * w[j]}
			}
		}
	case 3:
		for k := 0; k < n1d; k++ {
			for j := 0; j < n1d; j++ {
				for i := 0; i < n1d; i++ {
					m := i + n1d*j + (n1d*n1d)*k
					pts[m] = []float64{x[i], x[j], x[k], w[i] * w[j] * w[k]}
				}
			}
		}
	}
	return
}

// QuadPointsWilson5 generates 5 integration points according to Wilson's Appendix G-7 formulae
//    w0input  -- if w0input > 0, use this value instead of default w0=8/3 (corner)
//    p4stable -- if true, use w0=0.004 and wa=0.999 to mimic 4-point rule
func QuadPointsWilson5(w0input float64, p4stable bool) (pts [][]float64) {
	w0 := 8.0 / 3.0
	wa := 1.0 / 3.0
	a := 1.0
	if w0input > 0 {
		w0 = w0input
		wa = 1.0 - w0/4.0
		a = math.Sqrt(1.0 / (3.0 * wa))
	}
	if p4stable {
		w0 = 0.004
		wa = 0.999
		a = 0.5776391
	}
	return [][]float64{
		{-a, -a, 0, wa},
		{+a, -a, 0, wa},
		{+0, +0, 0, w0},
		{-a, +a, 0, wa},
		{+a, +a, 0, wa},
	}
}

// QuadPointsWilson8 generates 8 integration points according to Wilson's Appendix G-7 formulae
//    wbinput -- if wbinput > 0, use this value instead of default wb=40/49
func QuadPointsWilson8(wbinput float64) (pts [][]float64) {
	a := math.Sqrt(7.0 / 9.0)
	b := math.Sqrt(7.0 / 15.0)
	wa := 9.0 / 49.0
	wb := 40.0 / 49.0
	if wbinput > 0 {
		wb = wbinput
		wa = 1.0 - wb
		swa := math.Sqrt(wa)
		a = 1.0 / math.Sqrt(3.0*swa)
		b = math.Sqrt((2.0 - 2.0*swa) / (3.0 * wb))
	}
	return [][]float64{
		{-a, -a, 0, wa},
		{+0, -b, 0, wb},
		{+a, -a, 0, wa},
		{-b, +0, 0, wb},
		{+b, +0, 0, wb},
		{-a, +a, 0, wa},
		{+0, +b, 0, wb},
		{+a, +a, 0, wa},
	}
}

// QuadPointsWilson9 computes the 9-points for hexahedra according to Wilson's Appendix G-7 formulae
//    w0input  -- if w0input > 0, use this value instead of default w0=16/3 (corner)
//    p8stable -- if true, use w0=0.008 and wa=0.999 to mimic 8-point rule
func QuadPointsWilson9(w0input float64, p8stable bool) (pts [][]float64) {
	w0 := 16.0 / 3.0
	wa := 1.0 / 3.0
	a := 1.0
	if w0input > 0 {
		w0 = w0input
		wa = 1.0 - w0/8.0
		a = math.Sqrt(1.0 / (3.0 * wa))
	}
	if p8stable {
		w0 = 0.008
		wa = 0.999
		a = 0.5776391
	}
	return [][]float64{
		{-a, -a, -a, wa},
		{+a, -a, -a, wa},
		{-a, +a, -a, wa},
		{+a, +a, -a, wa},
		{+0, +0, +0, w0},
		{-a, -a, +a, wa},
		{+a, -a, +a, wa},
		{-a, +a, +a, wa},
		{+a, +a, +a, wa},
	}
}

// QuadPointDraw draws quadrature point within standard rectangle or box
//   dx -- can be used to displace box; may be nil
func QuadPointDraw(pts [][]float64, ndim int, triOrTet bool, dx []float64, args *plt.A) {
	if args == nil {
		args = &plt.A{C: "r", M: "*", Mec: "r", NoClip: true}
	}
	if len(dx) != ndim {
		dx = []float64{0, 0, 0}
	}
	argsPoly := &plt.A{Fc: "none", Ec: "#2645cb", Closed: true, NoClip: true}
	if ndim == 2 {
		if triOrTet {
			plt.Polyline([][]float64{
				{dx[0], dx[1]}, {dx[0] + 1, dx[1]}, {dx[0] + 0, dx[1] + 1},
			}, argsPoly)
		} else {
			plt.Polyline([][]float64{
				{dx[0] - 1, dx[1] - 1}, {dx[0] + 1, dx[1] - 1}, {dx[0] + 1, dx[1] + 1}, {dx[0] - 1, dx[1] + 1},
			}, argsPoly)
		}
		for _, p := range pts {
			plt.PlotOne(dx[0]+p[0], dx[1]+p[1], args)
		}
	} else {
		if triOrTet {
		} else {
			plt.Box(dx[0]-1, dx[0]+1, dx[1]-1, dx[1]+1, dx[2]-1, dx[2]+1, &plt.A{Wire: true, Ls: "-", Ec: "#2645cb", Lw: 3})
		}
		for _, p := range pts {
			plt.Plot3dPoint(dx[0]+p[0], dx[1]+p[1], dx[2]+p[2], args)
		}
	}
}

/// map of integration points //////////////////////////////////////////////////////////////////////

var (
	// IntPoints holds integration points for all kinds of cells: lin,qua,hex,tri,tet
	// It maps [cellKind] => [options][npts][4] where 4 means r,s,t,w
	IntPoints map[int]map[string][][]float64

	// DefaultIntPoints holds the default integration points for all cell types
	// It maps [cellTypeIndex] => [npts][4] where 4 means r,s,t,w
	// NOTE: the highest number of integration points is selected,
	//       thus the default number may not be optimal.
	DefaultIntPoints [][][]float64
)

// IntPointsFindSet finds set of integration points by cell kind and set name
func IntPointsFindSet(cellKind int, setName string) (P [][]float64) {
	if cellKind < 0 || cellKind > KindNumMax {
		chk.Panic("cellKind = %d is invalid\n", cellKind)
	}
	db, ok := IntPoints[cellKind]
	if !ok {
		chk.Panic("integration points set for cellKind = %d is not implemented yet\n", cellKind)
	}
	if P, ok = db[setName]; !ok {
		chk.Panic("cannot find integration points set named = %q for cellKind = %d\n", setName, cellKind)
	}
	return
}

func init() {

	// set integration points for "lin" kind
	IntPoints = make(map[int]map[string][][]float64)
	IntPoints[KindLin] = map[string][][]float64{
		"legendre_1": {
			{0, 0, 0, 2},
		},
		"legendre_2": {
			{-0.5773502691896257, 0, 0, 1},
			{+0.5773502691896257, 0, 0, 1},
		},
		"legendre_3": {
			{-0.7745966692414834, 0, 0, 0.5555555555555556},
			{+0.0000000000000000, 0, 0, 0.8888888888888888},
			{+0.7745966692414834, 0, 0, 0.5555555555555556},
		},
		"legendre_4": {
			{-0.8611363115940526, 0, 0, 0.3478548451374538},
			{-0.3399810435848562, 0, 0, 0.6521451548625462},
			{+0.3399810435848562, 0, 0, 0.6521451548625462},
			{+0.8611363115940526, 0, 0, 0.3478548451374538},
		},
		"legendre_5": {
			{-0.9061798459386640, 0, 0, 0.2369268850561891},
			{-0.5384693101056831, 0, 0, 0.4786286704993665},
			{+0.0000000000000000, 0, 0, 0.5688888888888889},
			{+0.5384693101056831, 0, 0, 0.4786286704993665},
			{+0.9061798459386640, 0, 0, 0.2369268850561891},
		},
	}

	// set integration points for "qua" kind
	IntPoints[KindQua] = map[string][][]float64{
		"legendre_1": {
			{0, 0, 0, 4},
		},
		"legendre_4": {
			{-0.5773502691896257, -0.5773502691896257, 0, 1},
			{+0.5773502691896257, -0.5773502691896257, 0, 1},
			{-0.5773502691896257, +0.5773502691896257, 0, 1},
			{+0.5773502691896257, +0.5773502691896257, 0, 1},
		},
		"legendre_9": {
			{-0.7745966692414834, -0.7745966692414834, 0, 25.0 / 81.0},
			{+0.0000000000000000, -0.7745966692414834, 0, 40.0 / 81.0},
			{+0.7745966692414834, -0.7745966692414834, 0, 25.0 / 81.0},
			{-0.7745966692414834, +0.0000000000000000, 0, 40.0 / 81.0},
			{+0.0000000000000000, +0.0000000000000000, 0, 64.0 / 81.0},
			{+0.7745966692414834, +0.0000000000000000, 0, 40.0 / 81.0},
			{-0.7745966692414834, +0.7745966692414834, 0, 25.0 / 81.0},
			{+0.0000000000000000, +0.7745966692414834, 0, 40.0 / 81.0},
			{+0.7745966692414834, +0.7745966692414834, 0, 25.0 / 81.0},
		},
		"legendre_16":      QuadPointsGaussLegendre(2, 16),
		"wilson5corner_5":  QuadPointsWilson5(0, false),
		"wilson5stable_5":  QuadPointsWilson5(0, true),
		"wilson8default_8": QuadPointsWilson8(0),
	}

	// auxiliary constants
	SQ19by30 := math.Sqrt(19.0 / 30.0)
	SQ19by33 := math.Sqrt(19.0 / 33.0)

	// set integration points for "hex" kind
	IntPoints[KindHex] = map[string][][]float64{
		"legendre_8": {
			{-0.5773502691896257, -0.5773502691896257, -0.5773502691896257, 1},
			{+0.5773502691896257, -0.5773502691896257, -0.5773502691896257, 1},
			{-0.5773502691896257, +0.5773502691896257, -0.5773502691896257, 1},
			{+0.5773502691896257, +0.5773502691896257, -0.5773502691896257, 1},
			{-0.5773502691896257, -0.5773502691896257, +0.5773502691896257, 1},
			{+0.5773502691896257, -0.5773502691896257, +0.5773502691896257, 1},
			{-0.5773502691896257, +0.5773502691896257, +0.5773502691896257, 1},
			{+0.5773502691896257, +0.5773502691896257, +0.5773502691896257, 1},
		},
		"wilson9corner_9": QuadPointsWilson9(0, false),
		"wilson9stable_9": QuadPointsWilson9(0, true),
		"irons_6": {
			{-1, +0, +0, 4.0 / 3.0},
			{+1, +0, +0, 4.0 / 3.0},
			{+0, -1, +0, 4.0 / 3.0},
			{+0, +1, +0, 4.0 / 3.0},
			{+0, +0, -1, 4.0 / 3.0},
			{+0, +0, +1, 4.0 / 3.0},
		},
		"irons_14": {
			{+SQ19by30, 0.0, 0.0, 320.0 / 361.0},
			{-SQ19by30, 0.0, 0.0, 320.0 / 361.0},
			{0.0, +SQ19by30, 0.0, 320.0 / 361.0},
			{0.0, -SQ19by30, 0.0, 320.0 / 361.0},
			{0.0, 0.0, +SQ19by30, 320.0 / 361.0},
			{0.0, 0.0, -SQ19by30, 320.0 / 361.0},
			{+SQ19by33, +SQ19by33, +SQ19by33, 121.0 / 361.0},
			{-SQ19by33, +SQ19by33, +SQ19by33, 121.0 / 361.0},
			{+SQ19by33, -SQ19by33, +SQ19by33, 121.0 / 361.0},
			{-SQ19by33, -SQ19by33, +SQ19by33, 121.0 / 361.0},
			{+SQ19by33, +SQ19by33, -SQ19by33, 121.0 / 361.0},
			{-SQ19by33, +SQ19by33, -SQ19by33, 121.0 / 361.0},
			{+SQ19by33, -SQ19by33, -SQ19by33, 121.0 / 361.0},
			{-SQ19by33, -SQ19by33, -SQ19by33, 121.0 / 361.0},
		},
		"legendre_27": {
			{-0.774596669241483, -0.774596669241483, -0.774596669241483, 0.171467764060357},
			{+0.000000000000000, -0.774596669241483, -0.774596669241483, 0.274348422496571},
			{+0.774596669241483, -0.774596669241483, -0.774596669241483, 0.171467764060357},
			{-0.774596669241483, +0.000000000000000, -0.774596669241483, 0.274348422496571},
			{+0.000000000000000, +0.000000000000000, -0.774596669241483, 0.438957475994513},
			{+0.774596669241483, +0.000000000000000, -0.774596669241483, 0.274348422496571},
			{-0.774596669241483, +0.774596669241483, -0.774596669241483, 0.171467764060357},
			{+0.000000000000000, +0.774596669241483, -0.774596669241483, 0.274348422496571},
			{+0.774596669241483, +0.774596669241483, -0.774596669241483, 0.171467764060357},
			{-0.774596669241483, -0.774596669241483, +0.000000000000000, 0.274348422496571},
			{+0.000000000000000, -0.774596669241483, +0.000000000000000, 0.438957475994513},
			{+0.774596669241483, -0.774596669241483, +0.000000000000000, 0.274348422496571},
			{-0.774596669241483, +0.000000000000000, +0.000000000000000, 0.438957475994513},
			{+0.000000000000000, +0.000000000000000, +0.000000000000000, 0.702331961591221},
			{+0.774596669241483, +0.000000000000000, +0.000000000000000, 0.438957475994513},
			{-0.774596669241483, +0.774596669241483, +0.000000000000000, 0.274348422496571},
			{+0.000000000000000, +0.774596669241483, +0.000000000000000, 0.438957475994513},
			{+0.774596669241483, +0.774596669241483, +0.000000000000000, 0.274348422496571},
			{-0.774596669241483, -0.774596669241483, +0.774596669241483, 0.171467764060357},
			{+0.000000000000000, -0.774596669241483, +0.774596669241483, 0.274348422496571},
			{+0.774596669241483, -0.774596669241483, +0.774596669241483, 0.171467764060357},
			{-0.774596669241483, +0.000000000000000, +0.774596669241483, 0.274348422496571},
			{+0.000000000000000, +0.000000000000000, +0.774596669241483, 0.438957475994513},
			{+0.774596669241483, +0.000000000000000, +0.774596669241483, 0.274348422496571},
			{-0.774596669241483, +0.774596669241483, +0.774596669241483, 0.171467764060357},
			{+0.000000000000000, +0.774596669241483, +0.774596669241483, 0.274348422496571},
			{+0.774596669241483, +0.774596669241483, +0.774596669241483, 0.171467764060357},
		},
	}

	// set integration points for "tri" kind
	IntPoints[KindTri] = map[string][][]float64{
		"internal_1": {
			{1.0 / 3.0, 1.0 / 3.0, 0, 1.0 / 2.0},
		},
		"internal_3": {
			{1.0 / 6.0, 1.0 / 6.0, 0, 1.0 / 6.0},
			{2.0 / 3.0, 1.0 / 6.0, 0, 1.0 / 6.0},
			{1.0 / 6.0, 2.0 / 3.0, 0, 1.0 / 6.0},
		},
		"edge_3": {
			{0.5, 0.5, 0, 1.0 / 6.0},
			{0.0, 0.5, 0, 1.0 / 6.0},
			{0.5, 0.0, 0, 1.0 / 6.0},
		},
		"internal_4": {
			{1.0 / 3.0, 1.0 / 3.0, 0, -27.0 / 96.0},
			{1.0 / 5.0, 1.0 / 5.0, 0, +25.0 / 96.0},
			{3.0 / 5.0, 1.0 / 5.0, 0, +25.0 / 96.0},
			{1.0 / 5.0, 3.0 / 5.0, 0, +25.0 / 96.0},
		},
		"internal_12": {
			{0.873821971016996, 0.063089014491502, 0, 0.0254224531851035},
			{0.063089014491502, 0.873821971016996, 0, 0.0254224531851035},
			{0.063089014491502, 0.063089014491502, 0, 0.0254224531851035},
			{0.501426509658179, 0.249286745170910, 0, 0.0583931378631895},
			{0.249286745170910, 0.501426509658179, 0, 0.0583931378631895},
			{0.249286745170910, 0.249286745170910, 0, 0.0583931378631895},
			{0.053145049844817, 0.310352451033784, 0, 0.041425537809187},
			{0.310352451033784, 0.053145049844817, 0, 0.041425537809187},
			{0.053145049844817, 0.636502499121398, 0, 0.041425537809187},
			{0.310352451033784, 0.636502499121398, 0, 0.041425537809187},
			{0.636502499121398, 0.053145049844817, 0, 0.041425537809187},
			{0.636502499121398, 0.310352451033784, 0, 0.041425537809187},
		},
		"internal_16": {
			{3.33333333333333e-01, 3.33333333333333e-01, 0, 7.21578038388935e-02},
			{8.14148234145540e-02, 4.59292588292723e-01, 0, 4.75458171336425e-02},
			{4.59292588292723e-01, 8.14148234145540e-02, 0, 4.75458171336425e-02},
			{4.59292588292723e-01, 4.59292588292723e-01, 0, 4.75458171336425e-02},
			{6.58861384496480e-01, 1.70569307751760e-01, 0, 5.16086852673590e-02},
			{1.70569307751760e-01, 6.58861384496480e-01, 0, 5.16086852673590e-02},
			{1.70569307751760e-01, 1.70569307751760e-01, 0, 5.16086852673590e-02},
			{8.98905543365938e-01, 5.05472283170310e-02, 0, 1.62292488115990e-02},
			{5.05472283170310e-02, 8.98905543365938e-01, 0, 1.62292488115990e-02},
			{5.05472283170310e-02, 5.05472283170310e-02, 0, 1.62292488115990e-02},
			{8.39477740995800e-03, 2.63112829634638e-01, 0, 1.36151570872175e-02},
			{7.28492392955404e-01, 8.39477740995800e-03, 0, 1.36151570872175e-02},
			{2.63112829634638e-01, 7.28492392955404e-01, 0, 1.36151570872175e-02},
			{8.39477740995800e-03, 7.28492392955404e-01, 0, 1.36151570872175e-02},
			{7.28492392955404e-01, 2.63112829634638e-01, 0, 1.36151570872175e-02},
			{2.63112829634638e-01, 8.39477740995800e-03, 0, 1.36151570872175e-02},
		},
	}

	// set integration points for "tet" kind
	IntPoints[KindTet] = map[string][][]float64{
		"internal_1": {
			{1.0 / 4.0, 1.0 / 4.0, 1.0 / 4.0, 1.0 / 6.0},
		},
		"internal_4": {
			{(5.0 + 3.0*utl.SQ5) / 20.0, (5.0 - utl.SQ5) / 20.0, (5.0 - utl.SQ5) / 20.0, 1.0 / 24},
			{(5.0 - utl.SQ5) / 20.0, (5.0 + 3.0*utl.SQ5) / 20.0, (5.0 - utl.SQ5) / 20.0, 1.0 / 24},
			{(5.0 - utl.SQ5) / 20.0, (5.0 - utl.SQ5) / 20.0, (5.0 + 3.0*utl.SQ5) / 20.0, 1.0 / 24},
			{(5.0 - utl.SQ5) / 20.0, (5.0 - utl.SQ5) / 20.0, (5.0 - utl.SQ5) / 20.0, 1.0 / 24},
		},
		"internal_5": {
			{1.0 / 4.0, 1.0 / 4.0, 1.0 / 4.0, -2.0 / 15.0},
			{1.0 / 6.0, 1.0 / 6.0, 1.0 / 6.0, +3.0 / 40.0},
			{1.0 / 6.0, 1.0 / 6.0, 1.0 / 2.0, +3.0 / 40.0},
			{1.0 / 6.0, 1.0 / 2.0, 1.0 / 6.0, +3.0 / 40.0},
			{1.0 / 2.0, 1.0 / 6.0, 1.0 / 6.0, +3.0 / 40.0},
		},
		"internal_6": {
			{+1.0, +0.0, +0.0, 4.0 / 3.0},
			{-1.0, +0.0, +0.0, 4.0 / 3.0},
			{+0.0, +1.0, +0.0, 4.0 / 3.0},
			{+0.0, -1.0, +0.0, 4.0 / 3.0},
			{+0.0, +0.0, +1.0, 4.0 / 3.0},
			{+0.0, +0.0, -1.0, 4.0 / 3.0},
		},
	}

	// set default integration points
	DefaultIntPoints = make([][][]float64, TypeNumMax)
	DefaultIntPoints[TypeLin2] = IntPoints[KindLin]["legendre_2"]
	DefaultIntPoints[TypeLin3] = IntPoints[KindLin]["legendre_3"]
	DefaultIntPoints[TypeLin4] = IntPoints[KindLin]["legendre_4"]
	DefaultIntPoints[TypeLin5] = IntPoints[KindLin]["legendre_5"]
	DefaultIntPoints[TypeTri3] = IntPoints[KindTri]["internal_3"]
	DefaultIntPoints[TypeTri6] = IntPoints[KindTri]["internal_4"]
	DefaultIntPoints[TypeTri10] = IntPoints[KindTri]["internal_12"]
	DefaultIntPoints[TypeTri15] = IntPoints[KindTri]["internal_16"]
	DefaultIntPoints[TypeQua4] = IntPoints[KindQua]["legendre_4"]
	DefaultIntPoints[TypeQua8] = IntPoints[KindQua]["legendre_9"]
	DefaultIntPoints[TypeQua9] = IntPoints[KindQua]["legendre_9"]
	DefaultIntPoints[TypeQua12] = IntPoints[KindQua]["legendre_16"]
	DefaultIntPoints[TypeQua16] = IntPoints[KindQua]["legendre_16"]
	DefaultIntPoints[TypeQua17] = IntPoints[KindQua]["legendre_16"]
	DefaultIntPoints[TypeTet4] = IntPoints[KindTet]["internal_4"]
	DefaultIntPoints[TypeTet10] = IntPoints[KindTet]["internal_6"]
	DefaultIntPoints[TypeHex8] = IntPoints[KindHex]["legendre_8"]
	DefaultIntPoints[TypeHex20] = IntPoints[KindHex]["legendre_27"]
}
