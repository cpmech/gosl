// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Reference:
// [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
//     and Mathematical Tables. U.S. Department of Commerce, NIST

func Test_orthopoly01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("orthopoly01. Jacobi P5(1,1)")

	N, α, β := 5, 1.0, 1.0
	op := NewGeneralOrthoPoly(OP_JACOBI, N, []*P{
		&P{N: "alpha", V: α, Min: -1, Max: math.MaxFloat64},
		&P{N: "beta", V: β, Min: -1, Max: math.MaxFloat64},
	})

	xx := utl.LinSpace(-1, 1, 5)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Scalar(tst, "P0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Scalar(tst, "P1", 1e-15, y, 0.5*(2*(α+1)+(α+β+2)*(x-1)))
		y = op.P(5, x)
		r := x - 1.0
		chk.Scalar(tst, "P5", 1e-15, y, (95040*math.Pow(r, 5)+475200*math.Pow(r, 4)+864000*math.Pow(r, 3)+691200*math.Pow(r, 2)+230400*(r)+23040)/3840.0)
	}
}

func Test_orthopoly02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("orthopoly02. Jacobi polynomials Fig 22.1 [1]")

	N, α, β := 5, 1.5, -0.5
	op := NewGeneralOrthoPoly(OP_JACOBI, N, []*P{
		&P{N: "alpha", V: α, Min: -1, Max: math.MaxFloat64},
		&P{N: "beta", V: β, Min: -1, Max: math.MaxFloat64},
	})

	xx := utl.LinSpace(-1, 1, 5)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Scalar(tst, "P0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Scalar(tst, "P1", 1e-15, y, 0.5*(2*(α+1)+(α+β+2)*(x-1)))
	}

	if chk.Verbose {
		plt.Reset(true, &plt.A{Prop: 1.5})
		X := utl.LinSpace(-1, 1, 101)
		Y := make([]float64, len(X))
		for n := 0; n <= 5; n++ {
			for i := 0; i < len(X); i++ {
				Y[i] = op.P(n, X[i])
			}
			plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}^{(%g,%g)}$", n, α, β), NoClip: true})
		}
		plt.Cross(0, 0, nil)
		plt.Equal()
		plt.AxisYrange(-1, 3.3)
		plt.HideAllBorders()
		plt.Gll("$x$", io.Sf("$P_n^{(%g,%g)}$", α, β), &plt.A{LegOut: true, LegNcol: 3})
		plt.Save("/tmp/gosl", "orthopoly02")
	}
}

func Test_orthopoly03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("orthopoly03. Legendre polynomials Fig 22.8 [1]")

	N := 5
	op := NewGeneralOrthoPoly(OP_LEGENDRE, N, nil)
	opj := NewGeneralOrthoPoly(OP_JACOBI, N, []*P{
		&P{N: "alpha", V: 0, Min: -1, Max: math.MaxFloat64},
		&P{N: "beta", V: 0, Min: -1, Max: math.MaxFloat64},
	})

	xx := utl.LinSpace(-1, 1, 5)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Scalar(tst, "P0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Scalar(tst, "P1", 1e-15, y, x)
		y = op.P(2, x)
		chk.Scalar(tst, "P2", 1e-15, y, (-1+3*x*x)/2)
		y = op.P(3, x)
		chk.Scalar(tst, "P3", 1e-15, y, (-3*x+5*x*x*x)/2)
		y = op.P(4, x)
		chk.Scalar(tst, "P4", 1e-15, y, (3-30*x*x+35*x*x*x*x)/8)
		y = op.P(5, x)
		chk.Scalar(tst, "P5", 1e-15, y, (15*x-70*x*x*x+63*math.Pow(x, 5))/8)
		for n := 0; n < 5; n++ {
			chk.Scalar(tst, "Legendre-Jacobi", 1e-15, op.P(n, x), opj.P(n, x))
		}
	}

	if chk.Verbose {
		plt.Reset(true, &plt.A{Prop: 1.0})
		X := utl.LinSpace(-1, 1, 101)
		Y := make([]float64, len(X))
		for n := 0; n <= 5; n++ {
			for i := 0; i < len(X); i++ {
				Y[i] = op.P(n, X[i])
			}
			plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}$", n), NoClip: true})
		}
		plt.Cross(0, 0, nil)
		plt.Equal()
		plt.AxisYrange(-0.5, 1.0)
		plt.HideAllBorders()
		plt.Gll("$x$", "$P_n$", &plt.A{LegOut: true, LegNcol: 3})
		plt.Save("/tmp/gosl", "orthopoly03")
	}
}

func Test_orthopoly04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("orthopoly04. Hermite polynomials Fig 22.10 [1]")

	N := 5
	op := NewGeneralOrthoPoly(OP_HERMITE, N, nil)

	xx := utl.LinSpace(-2, 2, 7)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Scalar(tst, "H0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Scalar(tst, "H1", 1e-15, y, 2*x)
		y = op.P(2, x)
		chk.Scalar(tst, "H2", 1e-15, y, 4*x*x-2)
		y = op.P(3, x)
		chk.Scalar(tst, "H3", 1e-13, y, 8*x*x*x-12*x)
		y = op.P(4, x)
		chk.Scalar(tst, "H4", 1e-13, y, 16*math.Pow(x, 4)-48*x*x+12)
		y = op.P(5, x)
		chk.Scalar(tst, "H5", 1e-13, y, 32*math.Pow(x, 5)-160*x*x*x+120*x)
	}

	if chk.Verbose {
		plt.Reset(true, &plt.A{Prop: 0.8})
		X := utl.LinSpace(0, 4, 101)
		Y := make([]float64, len(X))
		for n := 2; n <= 5; n++ {
			den := math.Pow(float64(n), 3)
			for i := 0; i < len(X); i++ {
				Y[i] = op.P(n, X[i]) / den
			}
			plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}/%g$", n, den), NoClip: true})
		}
		plt.Cross(0, 0, nil)
		plt.AxisYrange(-1, 8.5)
		plt.HideAllBorders()
		plt.Gll("$x$", "$P_n$", &plt.A{LegOut: true, LegNcol: 3})
		plt.Save("/tmp/gosl", "orthopoly04")
	}
}

func Test_orthopoly05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("orthopoly05. Chebyshev1 polynomials Fig 22.6 [1]")

	N := 5
	op := NewGeneralOrthoPoly(OP_CHEBYSHEV1, N, nil)

	xx := []float64{-2, -1.5, -1, -0.5, -1.0 / 3.0, 0, 1.0 / 3.0, 0.5, 1, 1.5, 2}
	for _, x := range xx {
		y := op.P(0, x)
		chk.Scalar(tst, "T0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Scalar(tst, "T1", 1e-15, y, x)
		y = op.P(2, x)
		chk.Scalar(tst, "T2", 1e-15, y, -1+2*x*x)
		y = op.P(3, x)
		chk.Scalar(tst, "T3", 1e-13, y, -3*x+4*x*x*x)
		y = op.P(4, x)
		chk.Scalar(tst, "T4", 1e-13, y, 1-8*x*x+8*math.Pow(x, 4))
		y = op.P(5, x)
		chk.Scalar(tst, "T5", 1e-13, y, 5*x-20*x*x*x+16*math.Pow(x, 5))
	}

	io.Pl()
	for _, x := range xx {
		for n := 0; n < 5; n++ {
			yref := ChebyshevT(n, x)
			y := op.P(n, x)
			tol := 1e-15
			if math.Abs(y) > 8 {
				tol = 1e-14
			}
			if math.Abs(y) > 20 {
				tol = 1e-13
			}
			chk.AnaNum(tst, io.Sf("p%d(%+.2f)", n, x), tol, y, yref, chk.Verbose)
		}
	}

	if chk.Verbose {
		plt.Reset(true, &plt.A{Prop: 0.8})
		X := utl.LinSpace(-1, 1, 101)
		Y := make([]float64, len(X))
		for n := 1; n <= 5; n++ {
			for i := 0; i < len(X); i++ {
				Y[i] = op.P(n, X[i])
			}
			plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}$", n), NoClip: true})
		}
		plt.Cross(0, 0, nil)
		plt.HideAllBorders()
		plt.Gll("$x$", "$P_n$", &plt.A{LegOut: true, LegNcol: 3})
		plt.Save("/tmp/gosl", "orthopoly05")
	}
}

func Test_orthopoly06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("orthopoly06. Chebyshev2 polynomials Fig 22.7 [1]")

	N := 5
	op := NewGeneralOrthoPoly(OP_CHEBYSHEV2, N, nil)

	xx := utl.LinSpace(-1, 1, 5)
	for _, x := range xx {
		y := op.P(0, x)
		chk.Scalar(tst, "U0", 1e-15, y, 1)
		y = op.P(1, x)
		chk.Scalar(tst, "U1", 1e-15, y, 2*x)
		y = op.P(2, x)
		chk.Scalar(tst, "U2", 1e-15, y, -1+4*x*x)
		y = op.P(3, x)
		chk.Scalar(tst, "U3", 1e-13, y, -4*x+8*x*x*x)
		y = op.P(4, x)
		chk.Scalar(tst, "U4", 1e-13, y, 1-12*x*x+16*math.Pow(x, 4))
		y = op.P(5, x)
		chk.Scalar(tst, "U5", 1e-13, y, 6*x-32*x*x*x+32*math.Pow(x, 5))
	}

	if chk.Verbose {
		plt.Reset(true, &plt.A{Prop: 0.8})
		X := utl.LinSpace(-1, 1, 101)
		Y := make([]float64, len(X))
		for n := 1; n <= 5; n++ {
			for i := 0; i < len(X); i++ {
				Y[i] = op.P(n, X[i])
			}
			plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}$", n), NoClip: true})
		}
		plt.Cross(0, 0, nil)
		plt.HideAllBorders()
		plt.AxisYrange(-3, 5.5)
		plt.Gll("$x$", "$P_n$", &plt.A{LegOut: true, LegNcol: 3})
		plt.Save("/tmp/gosl", "orthopoly06")
	}
}
