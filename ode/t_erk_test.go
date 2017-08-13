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

	methods := []string{"moeuler", "rk2", "rk3", "heun",
		"rk4", "rk4-3/8", "merson4", "zonneveld4",
		"dopri5", "fehlberg4", "fehlberg7", "verner6"}

	tols := []float64{1e-17, 1e-17, 1e-17, 1e-17,
		1e-17, 1e-15, 1e-17, 1e-17,
		1e-15, 1e-15, 1e-14, 1e-15}

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

		io.Pf("\ncoefficients: ci = Σ_j aij\n")
		var sumrow float64
		for i := 0; i < o.Nstg; i++ {
			sumrow = 0.0
			for j := 0; j < i; j++ {
				sumrow += o.A[i][j]
			}
			chk.AnaNum(tst, io.Sf("Σa%dj", i), tols[im], sumrow, o.C[i], chk.Verbose)
		}
		if o.Embedded {
			io.Pf("\nerror estimator\n")
			for i := 0; i < o.Nstg; i++ {
				chk.AnaNum(tst, io.Sf("E%d=B%d-Be%d", i, i, i), 1e-15, o.E[i], o.B[i]-o.Be[i], chk.Verbose)
			}
		}
	}
}
