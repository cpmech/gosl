// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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
	W := NewVector(len(U))
	VecAdd(w, 1, u, -2, v)
	VecAdd(W, 1, U, -1, V)
	chk.Array(tst, "w := 1⋅u - 2⋅v", 1e-17, w, []float64{-5, -2, 1})
	chk.Array(tst, "W := 1⋅U - 1⋅V", 1e-17, W, []float64{-15, -13, -11, -9, -7, -5, -3, -1, 1, 3, 5, 7, 9, 11, 13, 15})

	// VecMaxDiff
	chk.Float64(tst, "VecMaxDiff(u, w)", 1e-17, VecMaxDiff(u, w), 6.0)
	u[1] = 122
	chk.Float64(tst, "VecMaxDiff(u, v)", 1e-17, VecMaxDiff(u, v), 120)

	// VecScaleAbs
	scale := NewVector(len(w))
	VecScaleAbs(scale, -1, 2, w)
	chk.Array(tst, "scale := -1 + 2⋅w", 1e-17, scale, []float64{9, 3, 1})
}

func TestBlas1tst02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas1tst02. (real) Blas1 functions. Larger vectors")

	u := make([]float64, 200)
	v := make([]float64, 200)
	w := make([]float64, 200)
	wref := make([]float64, 200)
	for _, n := range []int{7, 8, 150, 151} {
		dot := 0.0
		for i := 0; i < n; i++ {
			u[i] = 2 * float64(1+i)
			v[i] = -float64(1 + i)
			wref[i] = u[i] + v[i]
			dot += u[i] * v[i]
		}
		VecAdd(w[:n], 1, u[:n], 1, v[:n])
		chk.Array(tst, io.Sf("n=%3d: w:=u-v", n), 1e-15, w[:n], wref[:n])
		chk.Float64(tst, "u⋅v", 1e-15, VecDot(u, v), dot)
	}
}
