// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// BezierQuad implements a quadratic Bezier curve
//  C(t) = (1-t)² Q0  +  2 t (1-t) Q1  +  t² Q2
//       = t² A  +  2 t B  +  Q0
//  A = Q2 - 2 Q1 + Q0
//  B = Q1 - Q0
type BezierQuad struct {
	Q [][]float64 // control points; can be set outside
}

// Point returns the x-y-z coordinates of a point on Bezier curve
func (o *BezierQuad) Point(C []float64, t float64) {
	if len(o.Q) != 3 {
		chk.Panic("Point: quadratic Bezier must be initialised first (with 3 control points)")
	}
	ndim := len(o.Q[0])
	chk.IntAssert(len(C), ndim)
	for i := 0; i < ndim; i++ {
		C[i] = (1.0-t)*(1.0-t)*o.Q[0][i] + 2.0*t*(1.0-t)*o.Q[1][i] + t*t*o.Q[2][i]
	}
	return
}

// DistPoint returns the distance from a point to this Bezier curve
func (o *BezierQuad) DistPoint(X []float64) float64 {
	if len(o.Q) != 3 {
		chk.Panic("DistPoint: quadratic Bezier must be initialised first (with 3 control points)")
	}
	ndim := len(o.Q[0])
	chk.IntAssert(len(X), ndim)
	var A_i, B_i, M_i, a, b, c, d float64
	for i := 0; i < ndim; i++ {
		A_i = o.Q[2][i] - 2.0*o.Q[1][i] + o.Q[0][i]
		B_i = o.Q[1][i] - o.Q[0][i]
		M_i = o.Q[0][i] - X[i]
		a += A_i * A_i
		b += 3.0 * A_i * B_i
		c += 2.0*B_i*B_i + M_i*A_i
		d += M_i * B_i
	}
	io.Pforan("a=%v b=%v c=%v d=%v\n", a, b, c, d)
	return 0
}
