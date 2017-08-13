// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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
