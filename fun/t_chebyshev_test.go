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

func checkSymmetry(tst *testing.T, X []float64) {
	l := len(X) - 1
	if -X[0] != X[l] {
		tst.Errorf("first and last coordinates are not symmetric. %g != %g\n", -X[0], X[l])
		return
	}
	chk.PrintOk("first and last")
	chk.Symmetry(tst, "X", X)
}

func TestChebyshev01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Chebyshev01. Gauss (roots) and Lobatto points")

	X := ChebyshevXgauss(1)
	Xref := []float64{-1.0 / math.Sqrt2, 1.0 / math.Sqrt2}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(1)
	Xref = []float64{-1.0, 1.0}
	chk.Array(tst, "X", 1e-17, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXgauss(2)
	Xref = []float64{-math.Sqrt(3.0) / 2.0, 0, math.Sqrt(3.0) / 2.0}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(2)
	Xref = []float64{-1, 0, 1}
	chk.Array(tst, "X and Xref", 1e-17, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXgauss(3)
	Xref = []float64{-9.238795325112867e-01, -3.826834323650898e-01, 3.826834323650897e-01, 9.238795325112867e-01}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(3)
	Xref = []float64{-1, -0.5, 0.5, 1}
	chk.Array(tst, "X and Xref", 1e-16, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXgauss(4)
	Xref = []float64{-9.510565162951535e-01, -5.877852522924731e-01, -6.123233995736766e-17, 5.877852522924730e-01, 9.510565162951535e-01}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(4)
	Xref = []float64{-1.000000000000000e+00, -7.071067811865476e-01, -6.123233995736766e-17, 7.071067811865475e-01, 1.000000000000000e+00}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXgauss(5)
	Xref = []float64{-9.659258262890683e-01, -7.071067811865476e-01, -2.588190451025210e-01, 2.588190451025206e-01, 7.071067811865475e-01, 9.659258262890682e-01}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(5)
	Xref = []float64{-1.000000000000000e+00, -8.090169943749475e-01, -3.090169943749475e-01, 3.090169943749473e-01, 8.090169943749473e-01, 1.000000000000000e+00}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXgauss(6)
	Xref = []float64{-9.749279121818236e-01, -7.818314824680298e-01, -4.338837391175582e-01, -6.123233995736766e-17, 4.338837391175581e-01, 7.818314824680298e-01, 9.749279121818236e-01}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(6)
	Xref = []float64{-1.000000000000000e+00, -8.660254037844387e-01, -5.000000000000001e-01, -6.123233995736766e-17, 4.999999999999998e-01, 8.660254037844385e-01, 1.000000000000000e+00}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXgauss(7)
	Xref = []float64{-9.807852804032304e-01, -8.314696123025452e-01, -5.555702330196023e-01, -1.950903220161283e-01, 1.950903220161282e-01, 5.555702330196020e-01, 8.314696123025453e-01, 9.807852804032304e-01}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(7)
	Xref = []float64{-1.000000000000000e+00, -9.009688679024191e-01, -6.234898018587336e-01, -2.225209339563144e-01, 2.225209339563143e-01, 6.234898018587335e-01, 9.009688679024190e-01, 1.000000000000000e+00}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXgauss(8)
	Xref = []float64{-9.848077530122080e-01, -8.660254037844387e-01, -6.427876096865394e-01, -3.420201433256688e-01, -6.123233995736766e-17, 3.420201433256687e-01, 6.427876096865394e-01, 8.660254037844387e-01, 9.848077530122080e-01}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(8)
	Xref = []float64{-1.000000000000000e+00, -9.238795325112867e-01, -7.071067811865476e-01, -3.826834323650898e-01, -6.123233995736766e-17, 3.826834323650897e-01, 7.071067811865475e-01, 9.238795325112867e-01, 1.000000000000000e+00}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXgauss(9)
	Xref = []float64{-9.876883405951378e-01, -8.910065241883679e-01, -7.071067811865476e-01, -4.539904997395468e-01, -1.564344650402309e-01, 1.564344650402308e-01, 4.539904997395467e-01, 7.071067811865475e-01, 8.910065241883678e-01, 9.876883405951377e-01}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)

	io.Pl()
	X = ChebyshevXlob(9)
	Xref = []float64{-1.000000000000000e+00, -9.396926207859084e-01, -7.660444431189780e-01, -5.000000000000001e-01, -1.736481776669304e-01, 1.736481776669303e-01, 4.999999999999998e-01, 7.660444431189779e-01, 9.396926207859083e-01, 1.000000000000000e+00}
	chk.Array(tst, "X and Xref", 1e-15, X, Xref)
	checkSymmetry(tst, X)
}

func TestChebyshev02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Chebyshev02. first derivative")

	X := append(append([]float64{-1.2, -1.1}, utl.LinSpace(-1, 1, 11)...), []float64{1.1, 1.2}...)

	for n := 0; n < 6; n++ {
		for i := 0; i < len(X); i++ {
			dTn := ChebyshevTdiff1(n, X[i])
			chk.DerivScaSca(tst, io.Sf("dT%d(%5.2f)", n, X[i]), 1e-8, dTn, X[i], 1e-3, chk.Verbose, func(t float64) (res float64, e error) {
				res = ChebyshevT(n, t)
				return
			})
		}
		io.Pl()
	}

	if chk.Verbose {
		N := 3
		plt.Reset(true, nil)
		x1 := utl.LinSpace(-1.0, 1.0, 201)
		x2 := utl.LinSpace(-1.1, 1.1, 201)
		y1 := utl.GetMapped(x1, func(x float64) float64 { return ChebyshevT(N, x) })
		y2 := utl.GetMapped(x2, func(x float64) float64 { return ChebyshevT(N, x) })
		yy1 := utl.GetMapped(x1, func(x float64) float64 { return ChebyshevTdiff1(N, x) })
		yy2 := utl.GetMapped(x2, func(x float64) float64 { return ChebyshevTdiff1(N, x) })
		plt.Plot(x1, y1, &plt.A{C: "r", Lw: 4, NoClip: true, L: "Tn(x)"})
		plt.Plot(x2, y2, &plt.A{C: "r", NoClip: true})
		plt.Gll("$x$", io.Sf("$T_%d(x)$", N), nil)
		plt.DoubleYscale(io.Sf("$dT_%d(x)/dx$", N))
		plt.Plot(x1, yy1, &plt.A{C: "b", Lw: 4, NoClip: true, L: "deriv"})
		plt.Plot(x2, yy2, &plt.A{C: "b", NoClip: true, L: "deriv"})
		plt.Save("/tmp/gosl/fun", "chebydiff1")
	}
}

func TestChebyshev03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Chebyshev03. second derivative")

	X := append(append([]float64{-1.2, -1.1}, utl.LinSpace(-1, 1, 11)...), []float64{1.1, 1.2}...)

	for n := 0; n < 6; n++ {
		for i := 0; i < len(X); i++ {
			x := X[i]
			ddTn := ChebyshevTdiff2(n, x)
			tol := 1e-8
			if n == 4 {
				if x == -1 || x == +1 {
					tol = 1e-5
				}
			}
			if n == 5 {
				if x == -1 || x == +1 {
					tol = 1e-4
				}
			}
			chk.DerivScaSca(tst, io.Sf("ddT%d(%5.2f)", n, x), tol, ddTn, x, 1e-3, chk.Verbose, func(t float64) (res float64, e error) {
				res = ChebyshevTdiff1(n, t)
				return
			})
		}
		io.Pl()
	}

	if chk.Verbose {
		N := 5
		plt.Reset(true, nil)
		x1 := utl.LinSpace(-1.0, 1.0, 201)
		x2 := utl.LinSpace(-1.1, 1.1, 201)
		y1 := utl.GetMapped(x1, func(x float64) float64 { return ChebyshevTdiff1(N, x) })
		y2 := utl.GetMapped(x2, func(x float64) float64 { return ChebyshevTdiff1(N, x) })
		yy1 := utl.GetMapped(x1, func(x float64) float64 { return ChebyshevTdiff2(N, x) })
		yy2 := utl.GetMapped(x2, func(x float64) float64 { return ChebyshevTdiff2(N, x) })
		plt.Plot(x1, y1, &plt.A{C: "r", Lw: 4, NoClip: true, L: "dTndx(x)"})
		plt.Plot(x2, y2, &plt.A{C: "r", NoClip: true})
		plt.Gll("$x$", io.Sf("$dT_%d(x)/dx$", N), nil)
		plt.DoubleYscale(io.Sf("$d^2T_%d(x)/dx^2$", N))
		plt.Plot(x1, yy1, &plt.A{C: "b", Lw: 4, NoClip: true, L: "deriv"})
		plt.Plot(x2, yy2, &plt.A{C: "b", NoClip: true, L: "deriv"})
		plt.Save("/tmp/gosl/fun", "chebydiff2")
	}
}
