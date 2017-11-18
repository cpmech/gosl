// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/utl"
)

// facNurbsT defines a structure to implement a factory of NURBS
type facNurbsT struct{}

// FactoryNurbs generates NURBS'
var FactoryNurbs = facNurbsT{}

// 2D curves //////////////////////////////////////////////////////////////////////////////////////

// Curve2dExample1 generates a 1D NURBS curve (example 1)
func (o facNurbsT) Curve2dExample1() (curve *Nurbs) {
	verts := [][]float64{
		{0.0, 0.0, 0, 1},
		{1.0, 0.2, 0, 1},
		{0.5, 1.5, 0, 1},
		{2.5, 2.0, 0, 1},
		{2.0, 0.4, 0, 1},
		{3.0, 0.0, 0, 1},
	}
	knots := [][]float64{
		{0, 0, 0, 0, 0.3, 0.7, 1, 1, 1, 1},
	}
	curve = NewNurbs(1, []int{3}, knots)
	curve.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Curve2dExample2 generates a 1D NURBS curve (example 2)
func (o facNurbsT) Curve2dExample2() (curve *Nurbs) {
	verts := [][]float64{
		{0.00, 0.00, 0, 0.80}, // 0
		{0.25, 0.15, 0, 1.00}, // 1
		{0.50, 0.00, 0, 0.70}, // 2
		{0.75, 0.00, 0, 1.20}, // 3
		{1.00, 0.10, 0, 1.10}, // 4
	}
	knots := [][]float64{
		{0, 0, 0, 1, 2, 3, 3, 3},
	}
	curve = NewNurbs(1, []int{2}, knots)
	curve.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Curve2dCircle generates a 1D NURBS representing the circle curve
func (o facNurbsT) Curve2dCircle(xc, yc, r float64) (curve *Nurbs) {
	xa, xb := xc-r, xc+r
	ya, yb := yc-r, yc+r
	verts := [][]float64{
		{xb, yc, 0, 1.0},
		{xb, yb, 0, 1.0 / math.Sqrt2},
		{xc, yb, 0, 1.0},
		{xa, yb, 0, 1.0 / math.Sqrt2},
		{xa, yc, 0, 1.0},
		{xa, ya, 0, 1.0 / math.Sqrt2},
		{xc, ya, 0, 1.0},
		{xb, ya, 0, 1.0 / math.Sqrt2},
		{xb, yc, 0, 1.0},
	}
	knots := [][]float64{
		{0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 4}, // only along first dimension
	}
	curve = NewNurbs(1, []int{2}, knots)
	curve.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Curve2dQuarterCircle generates a 1D NURBS representing a quarter of circle
func (o facNurbsT) Curve2dQuarterCircle(xc, yc, r float64) (curve *Nurbs) {
	verts := [][]float64{
		{xc + r, yc, 0, math.Sqrt2},
		{xc + r, yc + r, 0, 1.0},
		{xc, yc + r, 0, math.Sqrt2},
	}
	knots := [][]float64{
		{0, 0, 0, 1, 1, 1}, // only along first dimension
	}
	curve = NewNurbs(1, []int{2}, knots)
	curve.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// 2D surfaces ////////////////////////////////////////////////////////////////////////////////////

// Surf2dRectangleQL generates a 2D NURBS surface with x being quadratic and y being linear
func (o facNurbsT) Surf2dRectangleQL(x0, y0, dx, dy float64) (surf *Nurbs) {
	xm, xf, yf := x0+dx/2.0, x0+dx, y0+dy
	verts := [][]float64{
		{x0, y0, 0.0, 1.0},
		{xm, y0, 0.0, 1.0},
		{xf, y0, 0.0, 1.0},
		{x0, yf, 0.0, 1.0},
		{xm, yf, 0.0, 1.0},
		{xf, yf, 0.0, 1.0},
	}
	knots := [][]float64{
		{0, 0, 0, 1, 1, 1},
		{0, 0, 1, 1},
	}
	surf = NewNurbs(2, []int{2, 1}, knots)
	surf.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Surf2dQuarterCircle generates a 2D NURBS representing a quarter of circle
//   a, b -- inner and outer radii
func (o facNurbsT) Surf2dQuarterCircle(xc, yc, a, b float64) (surf *Nurbs) {
	verts := [][]float64{
		{xc + a, yc, 0, math.Sqrt2},
		{xc + b, yc, 0, math.Sqrt2},
		{xc + a, yc + a, 0, 1.0},
		{xc + b, yc + b, 0, 1.0},
		{xc, yc + a, 0, math.Sqrt2},
		{xc, yc + b, 0, math.Sqrt2},
	}
	knots := [][]float64{
		{0, 0, 1, 1},
		{0, 0, 0, 1, 1, 1},
	}
	surf = NewNurbs(2, []int{1, 2}, knots)
	surf.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Surf2dExample1 generates a 2D NURBS of a 2D strip (x-quadratic, y-linear) (example 1)
func (o facNurbsT) Surf2dExample1() (surf *Nurbs) {
	verts := [][]float64{
		{0.00, 0.00, 0, 0.80}, // 0
		{0.25, 0.15, 0, 1.00}, // 1
		{0.50, 0.00, 0, 0.70}, // 2
		{0.75, 0.00, 0, 1.20}, // 3
		{1.00, 0.10, 0, 1.10}, // 4
		{0.00, 0.40, 0, 0.90}, // 5
		{0.25, 0.55, 0, 0.60}, // 6
		{0.50, 0.40, 0, 1.50}, // 7
		{0.75, 0.40, 0, 1.40}, // 8
		{1.00, 0.50, 0, 0.50}, // 9
	}
	knots := [][]float64{
		{0, 0, 0, 1, 2, 3, 3, 3},
		{0, 0, 1, 1},
	}
	surf = NewNurbs(2, []int{2, 1}, knots)
	surf.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Surf2dQuarterPlateHole1 generates a 2D NURBS of a quarter of plate with hole (quadratic)
func (o facNurbsT) Surf2dQuarterPlateHole1() (surf *Nurbs) {
	verts := [][]float64{
		{-1.000000000000000e+00, 0.000000000000000e+00, 0, 1.000000000000000e+00},
		{-1.000000000000000e+00, 4.142135623730951e-01, 0, 8.535533905932737e-01},
		{-4.142135623730951e-01, 1.000000000000000e+00, 0, 8.535533905932737e-01},
		{0.000000000000000e+00, 1.000000000000000e+00, 0, 1.000000000000000e+00},
		{-2.500000000000000e+00, 0.000000000000000e+00, 0, 1.000000000000000e+00},
		{-2.500000000000000e+00, 7.500000000000000e-01, 0, 1.000000000000000e+00},
		{-7.500000000000000e-01, 2.500000000000000e+00, 0, 1.000000000000000e+00},
		{0.000000000000000e+00, 2.500000000000000e+00, 0, 1.000000000000000e+00},
		{-4.000000000000000e+00, 0.000000000000000e+00, 0, 1.000000000000000e+00},
		{-4.000000000000000e+00, 4.000000000000000e+00, 0, 1.000000000000000e+00},
		{0.000000000000000e+00, 4.000000000000000e+00, 0, 1.000000000000000e+00},
	}
	knots := [][]float64{
		{0, 0, 0, 0.5, 1, 1, 1},
		{0, 0, 0, 1, 1, 1},
	}
	surf = NewNurbs(2, []int{2, 2}, knots)
	//                                                    repeated
	//                                                      V  V
	surf.SetControl(verts, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 10})
	return
}

// 3D surfaces ////////////////////////////////////////////////////////////////////////////////////

// Surf3dCylinder generates a 3D NURBS representing the surface of a cyliner
func (o facNurbsT) Surf3dCylinder(xc, yc, zc, r, h float64) (curve *Nurbs) {
	xa, xb := xc-r, xc+r
	ya, yb := yc-r, yc+r
	verts := [][]float64{
		{xb, yc, zc, 1.0},
		{xb, yb, zc, 1.0 / math.Sqrt2},
		{xc, yb, zc, 1.0},
		{xa, yb, zc, 1.0 / math.Sqrt2},
		{xa, yc, zc, 1.0},
		{xa, ya, zc, 1.0 / math.Sqrt2},
		{xc, ya, zc, 1.0},
		{xb, ya, zc, 1.0 / math.Sqrt2},
		{xb, yc, zc, 1.0},
		{xb, yc, zc + h, 1.0},
		{xb, yb, zc + h, 1.0 / math.Sqrt2},
		{xc, yb, zc + h, 1.0},
		{xa, yb, zc + h, 1.0 / math.Sqrt2},
		{xa, yc, zc + h, 1.0},
		{xa, ya, zc + h, 1.0 / math.Sqrt2},
		{xc, ya, zc + h, 1.0},
		{xb, ya, zc + h, 1.0 / math.Sqrt2},
		{xb, yc, zc + h, 1.0},
	}
	knots := [][]float64{
		{0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 4},
		{0, 0, 1, 1},
	}
	curve = NewNurbs(2, []int{2, 1}, knots)
	curve.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Surf3dTorus generates a NURBS toroidal surface
//  r -- is the tube radius
//  R -- is the distance from the centre of the torus to the centre of the tube
func (o facNurbsT) Surf3dTorus(xc, yc, zc, r, R float64) (surf *Nurbs) {
	a, b := R-r, R+r
	s, h := 1.0/math.Sqrt2, 0.5
	verts := [][]float64{
		// 1
		{xc + R, yc, zc - r, 1},
		{xc + b, yc, zc - r, s},
		{xc + b, yc, zc + 0, 1},
		{xc + b, yc, zc + r, s},
		{xc + R, yc, zc + r, 1},
		{xc + a, yc, zc + r, s},
		{xc + a, yc, zc + 0, 1},
		{xc + a, yc, zc - r, s},
		{xc + R, yc, zc - r, 1},
		// 2
		{xc + R, yc + R, zc - r, s},
		{xc + b, yc + b, zc - r, h},
		{xc + b, yc + b, zc + 0, s},
		{xc + b, yc + b, zc + r, h},
		{xc + R, yc + R, zc + r, s},
		{xc + a, yc + a, zc + r, h},
		{xc + a, yc + a, zc + 0, s},
		{xc + a, yc + a, zc - r, h},
		{xc + R, yc + R, zc - r, s},
		// 3
		{xc, yc + R, zc - r, 1},
		{xc, yc + b, zc - r, s},
		{xc, yc + b, zc + 0, 1},
		{xc, yc + b, zc + r, s},
		{xc, yc + R, zc + r, 1},
		{xc, yc + a, zc + r, s},
		{xc, yc + a, zc + 0, 1},
		{xc, yc + a, zc - r, s},
		{xc, yc + R, zc - r, 1},
		// 4
		{xc - R, yc + R, zc - r, s},
		{xc - b, yc + b, zc - r, h},
		{xc - b, yc + b, zc + 0, s},
		{xc - b, yc + b, zc + r, h},
		{xc - R, yc + R, zc + r, s},
		{xc - a, yc + a, zc + r, h},
		{xc - a, yc + a, zc + 0, s},
		{xc - a, yc + a, zc - r, h},
		{xc - R, yc + R, zc - r, s},
		// 5
		{xc - R, yc, zc - r, 1},
		{xc - b, yc, zc - r, s},
		{xc - b, yc, zc + 0, 1},
		{xc - b, yc, zc + r, s},
		{xc - R, yc, zc + r, 1},
		{xc - a, yc, zc + r, s},
		{xc - a, yc, zc + 0, 1},
		{xc - a, yc, zc - r, s},
		{xc - R, yc, zc - r, 1},
		// 6
		{xc - R, yc - R, zc - r, s},
		{xc - b, yc - b, zc - r, h},
		{xc - b, yc - b, zc + 0, s},
		{xc - b, yc - b, zc + r, h},
		{xc - R, yc - R, zc + r, s},
		{xc - a, yc - a, zc + r, h},
		{xc - a, yc - a, zc + 0, s},
		{xc - a, yc - a, zc - r, h},
		{xc - R, yc - R, zc - r, s},
		// 7
		{xc, yc - R, zc - r, 1},
		{xc, yc - b, zc - r, s},
		{xc, yc - b, zc + 0, 1},
		{xc, yc - b, zc + r, s},
		{xc, yc - R, zc + r, 1},
		{xc, yc - a, zc + r, s},
		{xc, yc - a, zc + 0, 1},
		{xc, yc - a, zc - r, s},
		{xc, yc - R, zc - r, 1},
		// 8
		{xc + R, yc + -R, zc - r, s},
		{xc + b, yc + -b, zc - r, h},
		{xc + b, yc + -b, zc + 0, s},
		{xc + b, yc + -b, zc + r, h},
		{xc + R, yc + -R, zc + r, s},
		{xc + a, yc + -a, zc + r, h},
		{xc + a, yc + -a, zc + 0, s},
		{xc + a, yc + -a, zc - r, h},
		{xc + R, yc + -R, zc - r, s},
		// 9
		{xc + R, yc, zc - r, 1},
		{xc + b, yc, zc - r, s},
		{xc + b, yc, zc + 0, 1},
		{xc + b, yc, zc + r, s},
		{xc + R, yc, zc + r, 1},
		{xc + a, yc, zc + r, s},
		{xc + a, yc, zc + 0, 1},
		{xc + a, yc, zc - r, s},
		{xc + R, yc, zc - r, 1},
	}
	knots := [][]float64{
		{0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 4},
		{0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 4},
	}
	surf = NewNurbs(2, []int{2, 2}, knots)
	surf.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Surf3dQuarterHemisphere generates a NURBS toroidal surface
func (o facNurbsT) Surf3dQuarterHemisphere(xc, yc, zc, r float64) (surf *Nurbs) {
	s := math.Sqrt2 / 2.0
	verts := [][]float64{
		{xc + r, yc + 0, zc + 0, 1},
		{xc + r, yc + r, zc + 0, s},
		{xc + 0, yc + r, zc + 0, 1},
		{xc + r, yc + 0, zc + r, s},
		{xc + r, yc + r, zc + r, 0.5},
		{xc + 0, yc + r, zc + r, s},
		{xc + 0, yc + 0, zc + r, 1},
		{xc + 0, yc + 0, zc + r, s},
		{xc + 0, yc + 0, zc + r, 1},
	}
	knots := [][]float64{
		{0, 0, 0, 1, 1, 1},
		{0, 0, 0, 1, 1, 1},
	}
	surf = NewNurbs(2, []int{2, 2}, knots)
	surf.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// Solids /////////////////////////////////////////////////////////////////////////////////////////

// SolidHex generates a solid hexahedron
func (o facNurbsT) SolidHex(corners [][]float64) (vol *Nurbs) {
	verts := utl.Alloc(8, 4)
	for i := 0; i < 8; i++ {
		for j := 0; j < 3; j++ {
			verts[i][j] = corners[i][j]
		}
		verts[i][3] = 1.0
	}
	knots := [][]float64{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 0, 1, 1},
	}
	vol = NewNurbs(3, []int{1, 1, 1}, knots)
	vol.SetControl(verts, utl.IntRange(len(verts)))
	return
}
