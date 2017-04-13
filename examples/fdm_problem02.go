// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fdm"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// Solving:
	//                 ∂u²        ∂u²
	//            - kx ———  -  ky ———  =  0
	//                 ∂x²        ∂y²
	//
	// in the domain [0, 1] x [0, 1] with u = 50 @ the top and left boundaries
	// the other Dirichlet boundary conditions are zero
	kx, ky := 1.0, 1.0

	// allocate grid
	var g fdm.Grid2d
	g.Init(0.0, 1.0, 0.0, 1.0, 101, 101)

	// ids of equations with prescribed (known, given) U values
	// all around the square domain
	peq := utl.IntUnique(g.L, g.T, g.B, g.R)

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
	fdm.AssemblePoisson2d(&K11, &K12, F1, kx, ky, nil, &g, &e)

	// set prescribed values (default is 0.0)
	U2 := make([]float64, e.N2)
	for _, eq := range g.L {
		U2[e.FR2[eq]] = 50.0
	}
	for _, eq := range g.T {
		U2[e.FR2[eq]] = 50.0
	}

	// prepare right-hand-side
	//   F1 = F1 - K12 * U2
	la.SpMatVecMulAdd(F1, -1, K12.ToMatrix(nil), U2)

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
	plt.SetForPng(0.8, 600, 150)
	plt.Contour(X, Y, F, "cmapidx=0")
	plt.Equal()
	plt.Gll("x", "y", "")
	plt.SaveD("/tmp/gosl", "fdm_problem02.png")
}
