// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import "github.com/cpmech/gosl/utl"

// FactoryNurbs2dStrip generates a NURBS of a 2D strip (x-quadratic, y-linear)
func FactoryNurbs2dStrip() (b *Nurbs) {
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
	b = new(Nurbs)
	b.Init(2, []int{2, 1}, knots)
	b.SetControl(verts, utl.IntRange(len(verts)))
	return
}

// FactoryNurbs2dPlateHole generates a NURBS of a 2D quarter of plate with hole (quadratic)
func FactoryNurbs2dPlateHole() (b *Nurbs) {
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
	b = new(Nurbs)
	b.Init(2, []int{2, 2}, knots)
	b.SetControl(verts, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 10})
	return
}

// FactoryNurbs1dCurveA generates a NURBS 1D curve
func FactoryNurbs1dCurveA() (b *Nurbs) {
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
	b = new(Nurbs)
	b.Init(1, []int{3}, knots)
	b.SetControl(verts, utl.IntRange(len(verts)))
	return
}
