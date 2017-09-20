// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// Reference:
	// [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas,
	//     Graphs, and Mathematical Tables. U.S. Department of Commerce, NIST

	// number of points
	npts := 501

	// Jacobi polynomials Fig 22.1
	N, α, β := 5, 1.5, -0.5
	jacobi := fun.NewGeneralOrthoPoly("J", N, α, β)
	plt.Reset(true, &plt.A{Prop: 1.5})
	X := utl.LinSpace(-1, 1, npts)
	Y := make([]float64, len(X))
	for n := 0; n <= 5; n++ {
		for i := 0; i < len(X); i++ {
			Y[i] = jacobi.P(n, X[i])
		}
		plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}^{(%g,%g)}$", n, α, β), NoClip: true})
	}
	plt.Cross(0, 0, nil)
	plt.Equal()
	plt.AxisYrange(-1, 3.3)
	plt.HideAllBorders()
	plt.Gll("$x$", io.Sf("$P_n^{(%g,%g)}$", α, β), &plt.A{LegOut: true, LegNcol: 3})
	plt.Save("/tmp/gosl", "as-fig-22-1")

	// Chebyshev1 polynomials Fig 22.6
	chebyshev1 := fun.NewGeneralOrthoPoly("T", N, 0, 0)
	plt.Reset(true, &plt.A{Prop: 0.8})
	for n := 1; n <= 5; n++ {
		for i := 0; i < len(X); i++ {
			Y[i] = chebyshev1.P(n, X[i])
		}
		plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}$", n), NoClip: true})
	}
	plt.Cross(0, 0, nil)
	plt.HideAllBorders()
	plt.Gll("$x$", "$P_n$", &plt.A{LegOut: true, LegNcol: 3})
	plt.Save("/tmp/gosl", "as-fig-22-6")

	// Chebyshev2 polynomials Fig 22.7
	chebyshev2 := fun.NewGeneralOrthoPoly("U", N, 0, 0)
	plt.Reset(true, &plt.A{Prop: 0.8})
	for n := 1; n <= 5; n++ {
		for i := 0; i < len(X); i++ {
			Y[i] = chebyshev2.P(n, X[i])
		}
		plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}$", n), NoClip: true})
	}
	plt.Cross(0, 0, nil)
	plt.HideAllBorders()
	plt.AxisYrange(-3, 5.5)
	plt.Gll("$x$", "$P_n$", &plt.A{LegOut: true, LegNcol: 3})
	plt.Save("/tmp/gosl", "as-fig-22-7")

	// Legendre polynomials Fig 22.8
	legendre := fun.NewGeneralOrthoPoly("L", N, 0, 0)
	plt.Reset(true, &plt.A{Prop: 1.0})
	for n := 0; n <= 5; n++ {
		for i := 0; i < len(X); i++ {
			Y[i] = legendre.P(n, X[i])
		}
		plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}$", n), NoClip: true})
	}
	plt.Cross(0, 0, nil)
	plt.Equal()
	plt.AxisYrange(-0.5, 1.0)
	plt.HideAllBorders()
	plt.Gll("$x$", "$P_n$", &plt.A{LegOut: true, LegNcol: 3})
	plt.Save("/tmp/gosl", "as-fig-22-8")

	// Hermite polynomials Fig 22.10
	hermite := fun.NewGeneralOrthoPoly("H", N, 0, 0)
	plt.Reset(true, &plt.A{Prop: 0.8})
	X = utl.LinSpace(0, 4, npts)
	for n := 2; n <= 5; n++ {
		den := math.Pow(float64(n), 3)
		for i := 0; i < len(X); i++ {
			Y[i] = hermite.P(n, X[i]) / den
		}
		plt.Plot(X, Y, &plt.A{L: io.Sf("$P_{%d}/%g$", n, den), NoClip: true})
	}
	plt.Cross(0, 0, nil)
	plt.AxisYrange(-1, 8.5)
	plt.HideAllBorders()
	plt.Gll("$x$", "$P_n$", &plt.A{LegOut: true, LegNcol: 3})
	plt.Save("/tmp/gosl", "as-fig-22-10")
}
