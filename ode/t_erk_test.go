// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
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

func testErk02(tst *testing.T) {

	verbose()
	chk.PrintTitle("Erk02.")

	yana := func(x float64) la.Vector {
		return []float64{
			math.Exp(math.Sin(x * x)),
			math.Exp(5.0 * math.Sin(x*x)),
			math.Sin(x*x) + 1.0,
			math.Cos(x * x),
		}
	}

	fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
		f[0] = 2.0 * x * y[3] * y[0]
		f[1] = 10.0 * x * y[3] * fun.PowP(y[0], 5)
		f[2] = 2.0 * x * y[3]
		f[3] = -2.0 * x * (y[2] - 1)
		return nil
	}

	/*
		yana := func(x float64) la.Vector {
			return []float64{
				math.Exp(math.Sin(x * x)),
				math.Exp(5.0 * math.Sin(x*x)),
				math.Sin(x*x) + 1.0,
				math.Cos(x * x),
			}
		}

			fcn := func(f la.Vector, dx, x float64, y la.Vector) error {
				f[0] = 2.0 * x * math.Pow(y[1], 1.0/5.0) * y[3]
				f[1] = 10.0 * x * math.Exp(5.0*(y[2]-1.0)) * y[3]
				f[2] = 2.0 * x * y[3]
				f[3] = -2.0 * x * math.Log(y[0])
				return nil
			}
	*/

	ya := la.Vector([]float64{1, 1, 1, 1})
	xf := 2.8
	ndim := 4

	if chk.Verbose {
		m := 1.0 / 4.0
		plt.Reset(true, nil)
		plt.Plot([]float64{0, 5}, []float64{2.5, 2.5 + m*5}, &plt.A{C: "k", Ls: ":", NoClip: true})
	}

	M := []string{"o", "+"}

	for im, method := range []string{"fweuler", "rk4"} {

		//method := "moeuler"
		//method := "dopri5"
		//method := "fweuler"
		//method := "bweuler"
		//method := "rk4"
		//method := "rk4-3/8"

		//l10fcs := utl.LinSpace(2.6, 3.5, 10)
		//l10fcs := []float64{3.0}
		//U := make([]float64, len(l10fcs))
		//V := make([]float64, len(l10fcs))
		Dx := utl.LinSpace(0.001, 0.01, 21)
		U := make([]float64, len(Dx))
		V := make([]float64, len(Dx))
		for idx, dx := range Dx {
			io.Pf("dx = %v\n", dx)
			//for idx, l10fc := range l10fcs {

			//fc := math.Pow(10.0, l10fc)
			//nsteps := fc / 4
			//dx := xf / (nsteps - 1)

			// configuration
			conf, err := NewConfig(method, "", nil)
			status(tst, err)
			//if false {
			//conf.FixedStp = dx
			//}
			//conf.SetTol(1e-9, 1e-9)
			//io.Pforan("dx = %v\n", dx)

			// allocate solver
			sol, err := NewSolver(ndim, conf, nil, fcn, nil, nil)
			status(tst, err)
			defer sol.Free()

			// solve ODE
			y := ya.GetCopy()
			err = sol.Solve(y, 0.0, xf)
			status(tst, err)

			// global error
			yRef := yana(xf)
			io.Pfblue2("y    = %+.8f\n", y)
			io.Pf("yRef = %+.8f\n", yRef)
			//e := la.VecRmsError(y, yRef, 0, 1, yRef)
			//e := la.VecMaxDiff(y, yRef)
			e := y.NormDiff(yRef)
			//e := y.NormDiff(yRef) / yRef.Norm()
			io.Pforan("e = %v\n", e)
			U[idx] = -math.Log10(e)
			V[idx] = math.Log10(float64(sol.Stat.Nfeval))
		}

		if chk.Verbose {
			plt.Plot(U, V, &plt.A{L: method, C: plt.C(im, 0), M: M[im], NoClip: true})
		}
	}

	if chk.Verbose {
		plt.Gll("$-log_{10}(err)$", "$log_{10}(nfeval)$", nil)
		plt.SetTicksXlist(utl.LinSpace(0, 5, 11))
		plt.SetTicksYlist(utl.LinSpace(2.5, 3.7, 13))
		plt.HideTRborders()
		plt.Save("/tmp/gosl/ode", "erk02")
	}
}
