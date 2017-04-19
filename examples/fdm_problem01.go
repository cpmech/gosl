// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fdm"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// Solving:
	//                 ∂²u        ∂²u
	//            - kx ———  -  ky ———  =  1
	//                 ∂x²        ∂y²
	//
	// with zero Dirichlet boundary conditions around [-1, 1] x [-1, 1] and
	kx, ky := 1.0, 1.0
	source := func(x, y float64, args ...interface{}) float64 {
		return 1.0
	}

	// closed-form solution (for reference)
	π, π3, N := math.Pi, math.Pow(math.Pi, 3.0), 50
	solution := func(x, y float64) (res float64) {
		res = (1.0 - x*x) / 2.0
		for i := 1; i < N; i += 2 {
			k := float64(i)
			a := k * π * (1.0 + x) / 2.0
			b := k * π * (1.0 + y) / 2.0
			c := k * π * (1.0 - y) / 2.0
			d := k * k * k * math.Sinh(k*π)
			res -= (16.0 / π3) * (math.Sin(a) / d) * (math.Sinh(b) + math.Sinh(c))
		}
		return
	}

	// allocate grid
	var g fdm.Grid2d
	g.Init(-1.0, 1.0, -1.0, 1.0, 11, 11)

	// ids of equations with prescribed (known, given) U values
	// all around the square domain
	peq := utl.IntUnique(g.B, g.R, g.T, g.L)

	// structure to hold equations ids.
	// each grid node corresponds to one equation
	// i.e. number of equations == g.N
	var e fdm.Equations
	e.Init(g.N, peq)

	// set K11 and K12 => corresponding to unknown eqs
	var K11, K12 la.Triplet
	fdm.InitK11andK12(&K11, &K12, &e)

	// assemble system
	F1 := make([]float64, e.N1)
	fdm.AssemblePoisson2d(&K11, &K12, F1, kx, ky, source, &g, &e)

	// set prescribed values (default == 0.0)
	U2 := make([]float64, e.N2)

	// solve linear problem:
	//   K11 * U1 = F1
	U1, err := la.SolveRealLinSys(&K11, F1)
	if err != nil {
		chk.Panic("solve failed: %v", err)
	}

	// merge solution with known values
	U := make([]float64, g.N)
	fdm.JoinVecs(U, U1, U2, &e)

	// plotting
	X, Y, F := g.Generate(nil, U)
	var gsol fdm.Grid2d
	gsol.Init(-1.0, 1.0, -1.0, 1.0, 101, 101)
	Xsol, Ysol, Fsol := gsol.Generate(solution, nil)
	plt.Reset(false, nil)
	plt.ContourF(X, Y, F, &plt.A{CmapIdx: 1})
	plt.ContourL(Xsol, Ysol, Fsol, &plt.A{Colors: []string{"yellow"}, Lw: 20})
	plt.Equal()
	plt.Gll("x", "y", nil)
	plt.Save("/tmp/gosl", "fdm_problem01")
}
