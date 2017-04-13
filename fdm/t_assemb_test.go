// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

func Test_assemb01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("assemb01")

	// grid
	var g Grid2d
	g.Init(0.0, 1.0, 0.0, 1.0, 3, 3)

	// equations numbering
	var e Equations
	e.Init(g.N, []int{0, 3, 6})

	// K11 and K12
	var K11, K12 la.Triplet
	InitK11andK12(&K11, &K12, &e)

	// assembly
	kx, ky := 1.0, 1.0
	F1 := make([]float64, e.N1)
	AssemblePoisson2d(&K11, &K12, F1, kx, ky, nil, &g, &e)

	// check
	K11d := K11.ToMatrix(nil).ToDense()
	K12d := K12.ToMatrix(nil).ToDense()
	K11c := [][]float64{
		{16.0, -4.0, -8.0, 0.0, 0.0, 0.0},
		{-8.0, 16.0, 0.0, -8.0, 0.0, 0.0},
		{-4.0, 0.0, 16.0, -4.0, -4.0, 0.0},
		{0.0, -4.0, -8.0, 16.0, 0.0, -4.0},
		{0.0, 0.0, -8.0, 0.0, 16.0, -4.0},
		{0.0, 0.0, 0.0, -8.0, -8.0, 16.0},
	}
	K12c := [][]float64{
		{-4.0, 0.0, 0.0},
		{0.0, 0.0, 0.0},
		{0.0, -4.0, 0.0},
		{0.0, 0.0, 0.0},
		{0.0, 0.0, -4.0},
		{0.0, 0.0, 0.0},
	}
	chk.Matrix(tst, "K11: ", 1e-16, K11d, K11c)
	chk.Matrix(tst, "K12: ", 1e-16, K12d, K12c)
}
