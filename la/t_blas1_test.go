// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestBlas1tst01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas1tst01. (real) Blas1 functions")

	// VecRmsError
	u := Vector([]float64{1, 2, 3})
	v := Vector([]float64{3, 2, 1})
	s := Vector([]float64{-1, -1, -1})
	a, m := 1.0, 1.0
	rms := VecRmsError(u, v, a, m, s)
	chk.Float64(tst, "rms", 1e-17, rms, math.Sqrt(2.0/3.0))

	// VecDot
	U := Vector([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	V := Vector([]float64{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1})
	chk.Float64(tst, "u・v", 1e-17, VecDot(u, v), 10.0)
	chk.Float64(tst, "U・V", 1e-17, VecDot(U, V), 816.0)

	// VecAdd
	w := NewVector(len(u))
	VecAdd(w, 1, u, -2, v)
	chk.Array(tst, "w := 1⋅u - 2⋅v", 1e-17, w, []float64{-5, -2, 1})

	// VecMaxDiff
	chk.Float64(tst, "VecMaxDiff(u, w)", 1e-17, VecMaxDiff(u, w), 6.0)

	// VecScaleAbs
	scale := NewVector(len(w))
	VecScaleAbs(scale, -1, 2, w)
	chk.Array(tst, "scale := -1 + 2⋅w", 1e-17, scale, []float64{9, 3, 1})
}
