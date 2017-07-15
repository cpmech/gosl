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

func TestChebyPoly01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyPoly01")

	// test function
	f := func(x float64) (float64, error) {
		return math.Cos(math.Exp(2.0 * x)), nil
	}

	// allocate polynomials
	N := 8
	che, err := NewChebyshevPoly(N, true) // Gauss-Chebyshev
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}
	lob, err := NewChebyshevPoly(N, false) // Gauss-Lobatto
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	// check P
	for _, x := range utl.LinSpace(-1, 1, 7) {
		chk.AnaNum(tst, io.Sf("T%d(%+.3f)", N, x), 1e-15, ChebyshevT(N, x), che.HierarchicalT(8, x), chk.Verbose)
	}

	// Gauss-Chebyshev: check points
	xref := []float64{-9.848077530122080e-01, -8.660254037844387e-01, -6.427876096865394e-01, -3.420201433256688e-01, -6.123233995736766e-17, 3.420201433256687e-01, 6.427876096865394e-01, 8.660254037844387e-01, 9.848077530122080e-01}
	chk.Vector(tst, "Gauss-Chebyshev: X", 1e-15, che.X, xref)

	// Gauss-Chebyshev: check coefficients of interpolant
	che.CoefInterpolantSlow(f)
	cref := []float64{5.005025576289825e-01, -4.734690106554930e-01, 3.343030345866715e-01, 5.329760324967350e-01, 2.005496385333029e-01, -1.552357980491117e-01, -2.768837833165416e-01, -2.160862487215637e-01, -1.033306390240169e-01}
	chk.Vector(tst, "Gauss-Chebyshev: CoefI", 1e-14, che.CoefI, cref)

	// Gauss-Chebyshev: check coefficients of projection
	che.EstimateCoefProjection(f)
	cref = []float64{5.003559557885667e-01, -4.738396675676836e-01, 3.337904287575258e-01, 5.326202849023425e-01, 2.014887911962803e-01, -1.505413304349933e-01, -2.650525046501985e-01, -1.959021686279372e-01, -8.320914768336027e-02}
	chk.Vector(tst, "Gauss-Chebyshev: CoefP", 1e-15, che.CoefP, cref)

	// Gauss-Lobatto: check points
	xref = []float64{-1.000000000000000e+00, -9.238795325112867e-01, -7.071067811865476e-01, -3.826834323650898e-01, -6.123233995736766e-17, 3.826834323650897e-01, 7.071067811865475e-01, 9.238795325112867e-01, 1.000000000000000e+00}
	chk.Vector(tst, "Gauss-Lobatto: X", 1e-15, lob.X, xref)

	// Gauss-Lobatto: check coefficients of interpolant
	lob.CoefInterpolantSlow(f)
	cref = []float64{4.998505262591759e-01, -4.745223967909372e-01, 3.345788625609180e-01, 5.372649935983991e-01, 2.133118551317398e-01, -1.303539051589940e-01, -2.449269176317319e-01, -2.036386455562469e-01, -8.320813059359007e-02}
	chk.Vector(tst, "Gauss-Lobatto: CoefI", 1e-15, lob.CoefI, cref)

	// Gauss-Lobatto: check coefficients of projection
	lob.EstimateCoefProjection(f)
	chk.Vector(tst, "Gauss-Lobatto: CoefP", 1e-15, lob.CoefP, che.CoefP)

	// Gauss-Chebyshev: estimate error
	Eproj, _ := che.EstimateMaxErr(f, true)
	Einte, _ := che.EstimateMaxErr(f, false)
	io.Pforan("Gauss-Chebyshev: E{proj} = %v\n", Eproj)
	io.Pforan("Gauss-Chebyshev: E{inte} = %v\n", Einte)

	// Gauss-Lobatto: estimate error
	Eproj, _ = lob.EstimateMaxErr(f, true)
	Einte, _ = lob.EstimateMaxErr(f, false)
	io.Pforan("Gauss-Lobatto: E{proj} = %v\n", Eproj)
	io.Pforan("Gauss-Lobatto: E{inte} = %v\n", Einte)

	if chk.Verbose {

		// plot projection and interpolation
		X := utl.LinSpace(-1, 1, 201)
		Y := make([]float64, len(X))
		Yinte := make([]float64, len(X))
		Yproj := make([]float64, len(X))
		for i, x := range X {
			Y[i], _ = f(x)
			Yinte[i] = lob.Approx(x, false)
			Yproj[i] = lob.Approx(x, true)
		}
		plt.Reset(true, nil)
		plt.Plot(X, Y, &plt.A{C: "r", L: "$f$"})
		plt.Plot(X, Yinte, &plt.A{C: "g", L: "$I_N^{GL}f$"})
		plt.Plot(X, Yproj, &plt.A{C: "b", L: "$\\Pi_N^{w}f$"})
		plt.Gll("$x$", "$f(x)$", nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/fun", "chebypoly01a")

		// plot error
		Nvalues := []float64{1, 8, 16, 24, 36, 40, 48, 60, 80, 100, 120}
		Yerr := make([]float64, len(Nvalues))
		for i, nn := range Nvalues {
			o, _ := NewChebyshevPoly(int(nn), false)
			o.EstimateCoefProjection(f)
			Yerr[i], _ = o.EstimateMaxErr(f, true)
		}
		plt.Reset(true, nil)
		plt.Plot(Nvalues, Yerr, &plt.A{C: "r", M: "o", Void: true, NoClip: true})
		plt.SetYlog()
		plt.Gll("$N$", "$||f-\\Pi_N\\{f\\}||$", nil)
		plt.Save("/tmp/gosl/fun", "chebypoly01b")
	}
}
