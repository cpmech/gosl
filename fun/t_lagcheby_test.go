// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func compareLambda(tst *testing.T, N int, f Ss, tolU, tolL float64) {

	// allocate Lagrange structure and calculate U
	lag, err := NewLagrangeInterp(N, ChebyGaussLobGridKind)
	chk.EP(err)
	err = lag.CalcU(f)
	chk.EP(err)

	// allocate Chebyshev structure and calculate U
	che, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)
	err = che.CalcCoefIs(f)
	chk.EP(err)

	// check U values
	io.Pf("\n-------------------------------- N = %d -----------------------------------\n", N)
	cheU := utl.GetReversed(che.CoefIs)
	if N < 9 {
		io.Pforan("lag.U = %+8.4f\n", lag.U)
		io.Pfyel("che.U = %+8.4f\n", cheU)
	}
	chk.Array(tst, "U", tolU, lag.U, cheU)

	// check λ values
	cheL := utl.GetReversed(che.Lam)
	//if N < 9 {
	//io.Pfcyan("lag.λ = %+8.4f\n", lag.Lam)
	//io.Pfblue2("che.λ = %+8.4f\n", cheL)
	//}
	//chk.Array(tst, "λ", tolL, lag.Lam, cheL)
	_ = cheL
}

func TestLagCheby01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCheby01")

	// test function
	f := func(x float64) (float64, error) {
		return math.Cos(math.Exp(2.0 * x)), nil
	}

	// test
	Nvals := []int{6, 7, 8, 9, 100}
	tolsU := []float64{1e-17, 1e-17, 1e-17, 1e-17, 1e-17}
	tolsL := []float64{1e-14, 1e-14, 1e-14, 1e-13, 0.2}
	Nvals = []int{1041}
	for k, N := range Nvals {
		compareLambda(tst, N, f, tolsU[k], tolsL[k])
	}
}
