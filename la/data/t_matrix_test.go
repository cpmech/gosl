// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
)

func TestA123(tst *testing.T) {

	//verbose()
	chk.PrintTitle("A123.")

	a := new(A123)
	a.Generate()
	chk.Array(tst, "A123", 1e-17, a.A, []float64{1, 4, 7, 2, 5, 8, 3, 6, 9})
	chk.Array(tst, "el(A123)", 1e-17, a.El, []float64{-0.464547273387671, -0.882905959653586, 0.408248290463862, -0.570795531228578, -0.239520420054206, -0.816496580927726, -0.677043789069485, 0.403865119545174, 0.408248290463863})
	chk.Array(tst, "er(A123)", 1e-17, a.Er, []float64{-0.231970687246286, -0.525322093301234, -0.818673499356181, -0.785830238742067, -0.086751339256628, 0.612327560228810, 0.408248290463864, -0.816496580927726, 0.408248290463863})
	chk.Array(tst, "ev(A123)", 1e-17, a.Ev, []float64{16.116843969807043, -1.116843969807043, 0.0})
	chk.Array(tst, "ai(A123)", 1e-17, a.Ai, []float64{-0.638888888888888, -0.055555555555556, 0.527777777777777, -0.166666666666666, 0.000000000000000, 0.166666666666666, 0.305555555555555, 0.055555555555556, -0.194444444444444})
	chk.Array(tst, "nl(A123)", 1e-17, a.Nl, []float64{1, -2, 1})
	chk.Array(tst, "nr(A123)", 1e-17, a.Nr, []float64{1, -2, 1})
	chk.Array(tst, "p(A123)", 1e-17, a.P, []float64{0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0})
	chk.Array(tst, "ll(A123)", 1e-17, a.Ll, []float64{1.0, 0.142857142857143, 0.571428571428571, 0.0, 1.00, 0.5, 0.0, 0.00, 1.0})
	chk.Array(tst, "uu(A123)", 1e-17, a.Uu, []float64{7.0, 0.00, 0.0, 8.0, 0.857142857142857, 0.0, 9.0, 1.714285714285714, 0.0})
	chk.Array(tst, "q(A123)", 1e-17, a.Q, []float64{-0.123091490979333, -0.492365963917331, -0.861640436855329, 0.904534033733291, 0.301511344577763, -0.301511344577763, 0.408248290463862, -0.816496580927726, 0.408248290463863})
	chk.Array(tst, "q(A123)", 1e-17, a.R, []float64{-8.124038404635959, 0.0, 0.0, -9.601136296387955, 0.904534033733293, 0.0, -11.078234188139948, 1.809068067466585, 0.0})
	chk.Array(tst, "rhs(A123)", 1e-17, a.RHS, []float64{10.0, 28.0, 46.0})
	chk.Array(tst, "sol(A123)", 1e-17, a.Sol, []float64{3.0, 2.0, 1.0})
	chk.Array(tst, "u(A123)", 1e-17, a.U, []float64{-0.214837238368397, -0.520587389464737, -0.826337540561078, 0.887230688346371, 0.249643952988297, -0.387942782369774, 0.408248290463863, -0.816496580927726, 0.408248290463863})
	chk.Array(tst, "s(A123)", 1e-17, a.S, []float64{16.848103352614210, 0.0, 0.0, 0.0, 1.068369514554710, 0.0, 0.0, 0.0, 0.0})
	chk.Array(tst, "v(A123)", 1e-17, a.V, []float64{-0.479671177877772, -0.572367793972062, -0.665064410066353, -0.776690990321560, -0.075686470104559, 0.625318050112443, -0.408248290463863, 0.816496580927726, -0.408248290463863})
	chk.AnaNum(tst, "det(A123)", 1e-17, a.Det, 0.0, chk.Verbose)
}

func printMat(m, n int, a []float64, fmt string) {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			io.Pf(fmt, a[i+j*m])
		}
		io.Pl()
	}
}

func checkAi(tst *testing.T, n int, a, ai []float64, tolI float64) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			aia := 0.0
			for k := 0; k < n; k++ {
				aia += a[i+k*n] * ai[k+j*n]
			}
			if i == j {
				if math.Abs(aia-1.0) > tolI {
					tst.Errorf("Ai⋅A: diagonal is not 1.0\n")
					return
				}
			} else {
				if math.Abs(aia) > tolI {
					tst.Errorf("Ai⋅A: off-diagonal is not 0.0\n")
					return
				}
			}
		}
	}
	io.PfGreen("Ai⋅A: OK\n")
}

func TestChebyT(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyT.")

	n := 4
	a := new(ChebyT)
	a.Generate(n)
	printMat(a.M, a.N, a.A, "%12.6f")
	//Print(a.M, a.N, a.A, "A")
	//Plot(a.M, a.N, a.A, "A")
	checkAi(tst, a.M, a.A, a.Ai, 1e-17)
	chk.AnaNum(tst, "det(ChebyT)", 1e-17, a.Det, Det(a.M, a.A), chk.Verbose)
}

func Test5x5(tst *testing.T) {

	//verbose()
	chk.PrintTitle("5x5.")

	// a := [][]float64{
	//    {12, 28, 22, 20, +8},
	//    {+0, +3, +5, 17, 28},
	//    {56, +0, 23, +1, +0},
	//    {12, 29, 27, 10, +1},
	//    {+9, +4, 13, +8, 22},
	// }

	a := []float64{12, 0, 56, 12, 9, 28, 3, 0, 29, 4, 22, 5, 23, 27, 13, 20, 17, 1, 10, 8, 8, 28, 0, 1, 22}
	//Print(5, 5, a, "A")
	det := Det(5, a)
	chk.AnaNum(tst, "det(a)", 1e-9, det, -167402.0, chk.Verbose)
}
