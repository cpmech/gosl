// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/utl"
)

// Reference:
// [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
//     and Mathematical Tables. U.S. Department of Commerce, NIST

func TestGenOrthoPoly01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GenOrthoPoly01 Jacobi P5(1,1)")

	N, α, β := 5, 1.0, 1.0
	op := NewGeneralOrthoPoly("J", N, α, β)

	xx := utl.LinSpace(-1, 1, 5)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Float64(tst, "P0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Float64(tst, "P1", 1e-15, y, 0.5*(2*(α+1)+(α+β+2)*(x-1)))
		y = op.P(5, x)
		r := x - 1.0
		chk.Float64(tst, "P5", 1e-15, y, (95040*math.Pow(r, 5)+475200*math.Pow(r, 4)+864000*math.Pow(r, 3)+691200*math.Pow(r, 2)+230400*(r)+23040)/3840.0)
	}
}

func TestGenOrthoPoly02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GenOrthoPoly02 Jacobi polynomials Fig 22.1 [1]")

	N, α, β := 5, 1.5, -0.5
	op := NewGeneralOrthoPoly("J", N, α, β)

	xx := utl.LinSpace(-1, 1, 5)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Float64(tst, "P0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Float64(tst, "P1", 1e-15, y, 0.5*(2*(α+1)+(α+β+2)*(x-1)))
	}
}

func calcLegendre(i int, x float64) float64 {
	if i == 0 {
		return 1.0
	}
	if i == 1 {
		return x
	}
	tjm2 := 1.0 // value at step j-2
	tjm1 := x   // value at step j-1
	var tj float64
	for j := 2; j <= i; j++ {
		tj = (float64(2*j-1)*x*tjm1 - float64(j-1)*tjm2) / float64(j)
		tjm2 = tjm1
		tjm1 = tj
	}
	return tjm1
}

func TestGenOrthoPoly03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GenOrthoPoly03 Legendre polynomials Fig 22.8 [1]")

	N := 5
	op := NewGeneralOrthoPoly("L", N, 0, 0)
	opj := NewGeneralOrthoPoly("J", N, 0, 0)

	xx := utl.LinSpace(-1, 1, 5)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Float64(tst, "P0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Float64(tst, "P1", 1e-15, y, x)
		y = op.P(2, x)
		chk.Float64(tst, "P2", 1e-15, y, (-1+3*x*x)/2)
		y = op.P(3, x)
		chk.Float64(tst, "P3", 1e-15, y, (-3*x+5*x*x*x)/2)
		y = op.P(4, x)
		chk.Float64(tst, "P4", 1e-15, y, (3-30*x*x+35*x*x*x*x)/8)
		y = op.P(5, x)
		chk.Float64(tst, "P5", 1e-15, y, (15*x-70*x*x*x+63*math.Pow(x, 5))/8)
		for n := 0; n < 5; n++ {
			chk.Float64(tst, "Legendre-Jacobi", 1e-15, op.P(n, x), opj.P(n, x))
			chk.Float64(tst, "calcLegendre", 1e-17, op.P(n, x), calcLegendre(n, x))
		}
	}
}

func TestGenOrthoPoly04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GenOrthoPoly04 Hermite polynomials Fig 22.10 [1]")

	N := 5
	op := NewGeneralOrthoPoly("H", N, 0, 0)

	xx := utl.LinSpace(-2, 2, 7)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Float64(tst, "H0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Float64(tst, "H1", 1e-15, y, 2*x)
		y = op.P(2, x)
		chk.Float64(tst, "H2", 1e-15, y, 4*x*x-2)
		y = op.P(3, x)
		chk.Float64(tst, "H3", 1e-13, y, 8*x*x*x-12*x)
		y = op.P(4, x)
		chk.Float64(tst, "H4", 1e-13, y, 16*math.Pow(x, 4)-48*x*x+12)
		y = op.P(5, x)
		chk.Float64(tst, "H5", 1e-13, y, 32*math.Pow(x, 5)-160*x*x*x+120*x)
	}
}

func TestGenOrthoPoly05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GenOrthoPoly05 Chebyshev1 polynomials Fig 22.6 [1]")

	N := 5
	op := NewGeneralOrthoPoly("T", N, 0, 0)

	xx := []float64{-2, -1.5, -1, -0.5, -1.0 / 3.0, 0, 1.0 / 3.0, 0.5, 1, 1.5, 2}
	for _, x := range xx {
		y := op.P(0, x)
		chk.Float64(tst, "T0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Float64(tst, "T1", 1e-15, y, x)
		y = op.P(2, x)
		chk.Float64(tst, "T2", 1e-15, y, -1+2*x*x)
		y = op.P(3, x)
		chk.Float64(tst, "T3", 1e-13, y, -3*x+4*x*x*x)
		y = op.P(4, x)
		chk.Float64(tst, "T4", 1e-13, y, 1-8*x*x+8*math.Pow(x, 4))
		y = op.P(5, x)
		chk.Float64(tst, "T5", 1e-13, y, 5*x-20*x*x*x+16*math.Pow(x, 5))
	}

	io.Pl()
	for _, x := range xx {
		for n := 0; n < 5; n++ {
			yref := ChebyshevT(n, x)
			y := op.P(n, x)
			tol := 1e-15
			if math.Abs(y) > 2 {
				tol = 1e-14
			}
			if math.Abs(y) > 8 {
				tol = 1e-14
			}
			if math.Abs(y) > 20 {
				tol = 1e-13
			}
			chk.AnaNum(tst, io.Sf("p%d(%+.2f)", n, x), tol, y, yref, chk.Verbose)
		}
	}
}

func TestGenOrthoPoly06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("GenOrthoPoly06 Chebyshev2 polynomials Fig 22.7 [1]")

	N := 5
	op := NewGeneralOrthoPoly("U", N, 0, 0)

	xx := utl.LinSpace(-1, 1, 5)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Float64(tst, "U0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Float64(tst, "U1", 1e-15, y, 2*x)
		y = op.P(2, x)
		chk.Float64(tst, "U2", 1e-15, y, -1+4*x*x)
		y = op.P(3, x)
		chk.Float64(tst, "U3", 1e-13, y, -4*x+8*x*x*x)
		y = op.P(4, x)
		chk.Float64(tst, "U4", 1e-13, y, 1-12*x*x+16*math.Pow(x, 4))
		y = op.P(5, x)
		chk.Float64(tst, "U5", 1e-13, y, 6*x-32*x*x*x+32*math.Pow(x, 5))
	}
}
