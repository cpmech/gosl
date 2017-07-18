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

func TestLagCheby01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCheby01")

	// test function
	f := func(x float64) (float64, error) {
		return math.Cos(math.Exp(2.0 * x)), nil
	}

	// allocate Lagrange structure and calculate U
	N := 7
	kind := ChebyGaussLobGridKind
	lag, err := NewLagrangeInterp(N, kind)
	chk.EP(err)
	err = lag.CalcU(f)
	chk.EP(err)

	// allocate Chebyshev structure and calculate U
	che, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)
	err = che.CalcCoefIs(f)
	chk.EP(err)

	// check U values
	cheU := utl.GetReversed(che.CoefIs)
	io.Pforan("lag.U = %+8.4f\n", lag.U)
	io.Pforan("che.U = %+8.4f\n", cheU)
	chk.Array(tst, "U", 1e-17, lag.U, cheU)
}
