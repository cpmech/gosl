// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_smp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("smp01")

	a, b, β, ϵ := -1.0, 0.5, 2.0, 1e-3

	λ := []float64{-8.0, -8.0, -8.0}
	N := make([]float64, 3)
	n := make([]float64, 3)
	m := NewSmpDirector(N, λ, a, b, β, ϵ)
	NewSmpUnitDirector(n, m, N)
	io.Pforan("λ = %v\n", λ)
	io.Pforan("N = %v\n", N)
	io.Pforan("m = %v\n", m)
	io.Pforan("n = %v\n", n)
	chk.Vector(tst, "n", 1e-15, n, []float64{a / SQ3, a / SQ3, a / SQ3})

	p, q, err := GenInvs(λ, n, a)
	if err != nil {
		chk.Panic("GenInvs failed:\n%v", err.Error())
	}
	io.Pforan("p = %v\n", p)
	io.Pforan("q = %v\n", q)
	if q < 0.0 || q > 1e-17 {
		chk.Panic("q=%g is incorrect", q)
	}
	if math.Abs(p-a*λ[0]) > 1e-14 {
		chk.Panic("p=%g is incorrect. err = %g", p, math.Abs(p-a*λ[0]))
	}
}
