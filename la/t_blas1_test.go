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
	chk.Scalar(tst, "rms", 1e-17, rms, math.Sqrt(2.0/3.0))

	// VecDot
	chk.Scalar(tst, "u・v", 1e-17, VecDot(u, v), 10.0)

	// VecAdd
	w := NewVector(len(u))
	VecAdd(w, 1, u, -2, v)
	chk.Vector(tst, "1⋅u - 2⋅v", 1e-17, w, []float64{-5, -2, 1})

	// VecMaxDiff
	chk.Scalar(tst, "VecMaxDiff(u, w)", 1e-17, VecMaxDiff(u, w), 6.0)
}
