// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
	"testing"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

func TestAssemb01(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	utl.TTitle("Assemb 01")

	// grid
	var g Grid2D
	g.Init(1.0, 1.0, 3, 3)

	// equations numbering
	var e Equations
	e.Init(g.N, []int{0, 3, 6})

	// K11 and K12
	var K11, K12 la.Triplet
	InitK11andK12(&K11, &K12, &e)

	// assembly
	F1 := make([]float64, e.N1)
	Assemble(&K11, &K12, F1, nil, &g, &e)

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
	utl.CheckMatrix(tst, "K11: ", 1e-16, K11d, K11c)
	utl.CheckMatrix(tst, "K12: ", 1e-16, K12d, K12c)
}
