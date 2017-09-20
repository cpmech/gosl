// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

// GetVec2d returns column of matrix as 2D vector
func GetVec2d(M *la.Matrix, col int, normalize bool) la.Vector {
	x := M.Get(0, col)
	y := M.Get(1, col)
	if normalize {
		s := math.Sqrt(x*x + y*y)
		if s > 0 {
			x /= s
			y /= s
		}
	}
	return []float64{x, y}
}

// GetVec3d returns column of matrix as 3D vector
func GetVec3d(M *la.Matrix, col int, normalize bool) la.Vector {
	x := M.Get(0, col)
	y := M.Get(1, col)
	z := M.Get(2, col)
	if normalize {
		s := math.Sqrt(x*x + y*y + z*z)
		if s > 0 {
			x /= s
			y /= s
			z /= s
		}
	}
	return []float64{x, y, z}
}

// DrawArrow2d draws 2d arrow @ c with direction v
func DrawArrow2d(c, v la.Vector, normalize bool, sf float64, args *plt.A) {
	x := v[0]
	y := v[1]
	if normalize {
		s := math.Sqrt(x*x + y*y)
		if s > 0 {
			x /= s
			y /= s
		}
	}
	plt.Arrow(c[0], c[1], c[0]+x*sf, c[1]+y*sf, args)
}

// DrawArrow2dM draws 2d arrow @ c with direction v (matrix version)
func DrawArrow2dM(c la.Vector, vm *la.Matrix, col int, normalize bool, sf float64, args *plt.A) {
	v := GetVec2d(vm, col, normalize)
	plt.Arrow(c[0], c[1], c[0]+v[0]*sf, c[1]+v[1]*sf, args)
}

// DrawArrow3d draws 3d arrow @ c with direction v
func DrawArrow3d(c, v la.Vector, normalize bool, sf float64, args *plt.A) {
	x := v[0]
	y := v[1]
	z := v[2]
	if normalize {
		s := math.Sqrt(x*x + y*y + z*z)
		if s > 0 {
			x /= s
			y /= s
			z /= s
		}
	}
	plt.Draw3dVector(c, []float64{x, y, z}, sf, false, args)
}

// DrawArrow3dM draws 3d arrow @ c with direction v (matrix version)
func DrawArrow3dM(c la.Vector, vm *la.Matrix, col int, normalize bool, sf float64, args *plt.A) {
	v := GetVec3d(vm, col, normalize)
	plt.Draw3dVector(c, v, sf, false, args)
}
