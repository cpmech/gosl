// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestErk01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Erk01. coefficients")

	methods := []string{"moeuler", "rk2", "rk3", "heun3",
		"rk4", "rk4-3/8", "merson4", "zonneveld4",
		"dopri5", "fehlberg4", "fehlberg7", "verner6"}

	tols1 := []float64{1e-17, 1e-17, 1e-17, 1e-17,
		1e-17, 1e-15, 1e-17, 1e-17,
		1e-15, 1e-15, 1e-14, 1e-15}

	tols2 := []float64{1e-17, 1e-17, 1e-15, 1e-17,
		1e-15, 1e-17, 1e-15, 1e-15,
		1e-15, 1e-17, 1e-14, 1e-17}

	tols3 := []float64{1e-17, 1e-17, 1e-17, 1e-17,
		1e-17, 1e-16, 1e-17, 1e-17,
		1e-15, 1e-15, 1e-14, 1e-16}

	for im, method := range methods {
		io.Pf("\n------------------------------------ %s ---------------------------------------\n", method)
		o := newERK(method).(*ExplicitRK)
		chk.Int(tst, "len(A)=nstg", len(o.A), o.Nstg)
		chk.Int(tst, "len(B)=nstg", len(o.B), o.Nstg)
		chk.Int(tst, "len(C)=nstg", len(o.C), o.Nstg)
		if o.Embedded {
			chk.Int(tst, "len(Be)=nstg", len(o.Be), o.Nstg)
			chk.Int(tst, "len(E)=nstg", len(o.E), o.Nstg)
		}
		for i := 0; i < o.Nstg; i++ {
			chk.Int(tst, "ncol(A)=nstg", len(o.A[i]), o.Nstg)
		}

		io.Pf("\nc coefficients: ci = Σ_j aij\n")
		var sumrow float64
		for i := 0; i < o.Nstg; i++ {
			sumrow = 0.0
			for j := 0; j < i; j++ {
				sumrow += o.A[i][j]
			}
			chk.AnaNum(tst, io.Sf("Σa%dj", i), tols1[im], sumrow, o.C[i], chk.Verbose)
		}
		if o.Embedded {
			io.Pf("\nerror estimator\n")
			for i := 0; i < o.Nstg; i++ {
				chk.AnaNum(tst, io.Sf("E%d=B%d-Be%d", i, i, i), 1e-15, o.E[i], o.B[i]-o.Be[i], chk.Verbose)
			}
		}

		io.Pf("\nEquations (1.11) of [1], page 135-136\n")
		Sb := 0.0
		Sbc := 0.0
		Sbc2 := 0.0
		Sbc3 := 0.0
		Sbac := 0.0
		Sbcac := 0.0
		Sbac2 := 0.0
		Sbaac := 0.0
		for i := 0; i < o.Nstg; i++ {
			Sb += o.B[i]
			Sbc += o.B[i] * o.C[i]
			Sbc2 += o.B[i] * o.C[i] * o.C[i]
			Sbc3 += o.B[i] * o.C[i] * o.C[i] * o.C[i]
			for j := 0; j < o.Nstg; j++ {
				Sbac += o.B[i] * o.A[i][j] * o.C[j]
				Sbcac += o.B[i] * o.C[i] * o.A[i][j] * o.C[j]
				Sbac2 += o.B[i] * o.A[i][j] * o.C[j] * o.C[j]
				for k := 0; k < o.Nstg; k++ {
					Sbaac += o.B[i] * o.A[i][j] * o.A[j][k] * o.C[k]
				}
			}
		}
		chk.AnaNum(tst, "Σbi           =1   ", tols2[im], Sb, 1.0, chk.Verbose)
		_ = tols3
		chk.AnaNum(tst, "Σbi⋅ci        =1/2 ", tols2[im], Sbc, 0.5, chk.Verbose)
		if o.P < 3 {
			continue
		}
		chk.AnaNum(tst, "Σbi⋅ci²       =1/3 ", tols3[im], Sbc2, 1.0/3.0, chk.Verbose)
		if o.P < 4 {
			continue
		}
		chk.AnaNum(tst, "Σbi⋅ci³       =1/4 ", tols3[im], Sbc3, 1.0/4.0, chk.Verbose)
		chk.AnaNum(tst, "Σbi⋅aij⋅cj    =1/6 ", tols3[im], Sbac, 1.0/6.0, chk.Verbose)
		chk.AnaNum(tst, "Σbi⋅ci⋅aij⋅cj =1/8 ", tols3[im], Sbcac, 1.0/8.0, chk.Verbose)
		chk.AnaNum(tst, "Σbi⋅aij⋅cj²   =1/12", tols3[im], Sbac2, 1.0/12.0, chk.Verbose)
		chk.AnaNum(tst, "Σbi⋅aij⋅ajk⋅ck=1/24", tols3[im], Sbaac, 1.0/24.0, chk.Verbose)
	}
}

func TestErk02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Erk02. simple problem")

	// problem
	p := ProbSimpleNdim2()
	p.Dx = 0.17

	// prepare plot
	if chk.Verbose {
		plt.Reset(true, nil)
	}

	// try methods
	M := []string{"s", ".", "+", "^", "x", "*", "|", ""}
	Ms := []int{4, 6, 6, 6, 6, 6, 6, 6}
	tols := []float64{0.11, 0.11, 0.0077, 0.009, 5.8e-4, 5.8e-4, 1.3e-4, 5.8e-4}
	for im, method := range []string{"moeuler", "rk2", "rk3", "heun3", "rk4", "rk4-3/8", "merson4", "zonneveld4"} {

		// solve problem
		y, _, out, err := p.Solve(method, true, false)
		status(tst, err)

		// results
		X := out.GetStepX()
		Y0 := out.GetStepY(0)
		Y1 := out.GetStepY(1)
		chk.Float64(tst, "xf == X[last]", 1e-15, p.Xf, X[out.StepIdx-1])
		chk.AnaNum(tst, "y0(xf)", tols[im], y[0], p.CalcYana(0, p.Xf), chk.Verbose)
		chk.AnaNum(tst, "y1(xf)", tols[im], y[1], p.CalcYana(1, p.Xf), chk.Verbose)

		// plot
		if chk.Verbose {
			plt.Plot(X, Y0, &plt.A{L: method + ":y0", C: plt.C(im*2+0, 8), M: M[im], Ms: Ms[im], NoClip: true})
			plt.Plot(X, Y1, &plt.A{L: method + ":y1", C: plt.C(im*2+1, 8), M: M[im], Ms: Ms[im], NoClip: true})
		}
	}

	// save plot
	if chk.Verbose {
		xx := utl.LinSpace(0, p.Xf, 101)
		y0 := make([]float64, len(xx))
		y1 := make([]float64, len(xx))
		for i := 0; i < len(xx); i++ {
			p.Yana(p.Ytmp, xx[i])
			y0[i], y1[i] = p.Ytmp[0], p.Ytmp[1]
		}
		plt.Plot(xx, y0, &plt.A{C: "k", Ls: "--", NoClip: true})
		plt.Plot(xx, y1, &plt.A{C: "k", Ls: "--", NoClip: true})
		plt.Gll("$x$", "$y$", &plt.A{LegOut: true, LegNcol: 4, LegHlen: 2})
		plt.Save("/tmp/gosl/ode", "erk02")
	}
}

func convergenceTest(tst *testing.T, p *Problem, methods []string, orders []int, tols []float64) {

	// constants
	ndx := 3
	dxs := utl.LinSpace(0.001, 0.01, ndx)
	U := make([]float64, ndx)
	V := make([]float64, ndx)
	lu := make([]float64, ndx)
	lv := make([]float64, ndx)

	// try methods
	for im, method := range methods {

		// run for many dx
		for idx, dx := range dxs {

			// solve problem
			p.Dx = dx
			y, stat, _, err := p.Solve(method, true, false)
			status(tst, err)

			// global error
			p.Yana(p.Ytmp, p.Xf)
			e := la.VecMaxDiff(y, p.Ytmp)
			U[idx] = float64(stat.Nfeval)
			V[idx] = e

			// debug
			if false { // fake slope of 1:4
				U[0] = math.Pow(10, 4)
				U[1] = math.Pow(10, 3.75)
				U[2] = math.Pow(10, 3.5)
				V[0] = math.Pow(10, -8)
				V[1] = math.Pow(10, -7)
				V[2] = math.Pow(10, -6)
			}

			// log-log values
			lu[idx] = math.Log10(U[idx])
			lv[idx] = math.Log10(V[idx])
		}

		// calc convergence rate
		_, m := num.LinFit(lu, lv)
		chk.AnaNum(tst, "slope m", tols[im], m, -4.0, chk.Verbose)

		if chk.Verbose {
			plt.Plot(U, V, &plt.A{L: method, C: plt.C(im, 0), M: plt.M(im, 0), NoClip: true})
		}
	}
}

func TestErk03a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Erk03a. rk-fours. problem a")

	// problem
	p := ProbSimpleNdim4a()

	// prepare plot
	if chk.Verbose {
		plt.Reset(true, nil)
	}

	// run test
	methods := []string{"rk4", "rk4-3/8", "merson4", "zonneveld4"}
	orders := []int{4, 4, 4, 4}
	tols := []float64{0.011, 0.023, 0.00471, 0.011}
	convergenceTest(tst, p, methods, orders, tols)

	// plot
	if chk.Verbose {
		plt.Gll("$nFeval$", "$error$", nil)
		plt.SlopeInd(-4, 3.8, -6, 0.4, "4", false, true, true, nil, nil)
		plt.SetXlog()
		plt.SetYlog()
		//plt.Equal()
		plt.Save("/tmp/gosl/ode", "erk03a")
	}
}

func TestErk03b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Erk03b. rk-fours. problem b")

	// problem
	p := ProbSimpleNdim4b()

	// prepare plot
	if chk.Verbose {
		plt.Reset(true, nil)
	}

	// run test
	methods := []string{"rk4", "rk4-3/8", "merson4", "zonneveld4"}
	orders := []int{4, 4, 4, 4}
	tols := []float64{0.086, 0.164, 0.07, 0.09}
	convergenceTest(tst, p, methods, orders, tols)

	// plot
	if chk.Verbose {
		plt.Gll("$nFeval$", "$error$", nil)
		plt.SlopeInd(-4, 3.8, -5, 0.4, "4", false, true, true, nil, nil)
		plt.SetXlog()
		plt.SetYlog()
		//plt.Equal()
		plt.Save("/tmp/gosl/ode", "erk03b")
	}
}
