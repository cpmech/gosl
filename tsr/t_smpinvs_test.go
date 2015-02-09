// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
    "math"
    "testing"
    "code.google.com/p/gosl/utl"
)

func Test_smp01(tst *testing.T) {

    prevTs := utl.Tsilent
    defer func() {
        utl.Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //utl.Tsilent = false
    utl.TTitle("smp01")

    a, b, β, ϵ := -1.0, 0.5, 2.0, 1e-3

    λ := []float64{-8.0, -8.0, -8.0}
    N := make([]float64, 3)
    n := make([]float64, 3)
    m := NewSmpDirector(N, λ, a, b, β, ϵ)
    NewSmpUnitDirector(n, m, N)
    utl.Pforan("λ = %v\n", λ)
    utl.Pforan("N = %v\n", N)
    utl.Pforan("m = %v\n", m)
    utl.Pforan("n = %v\n", n)
    utl.CheckVector(tst, "n", 1e-15, n, []float64{a/SQ3, a/SQ3, a/SQ3})

    p, q, err := GenInvs(λ, n, a)
    if err != nil {
        utl.Panic("GenInvs failed:\n%v", err.Error())
    }
    utl.Pforan("p = %v\n", p)
    utl.Pforan("q = %v\n", q)
    if q < 0.0 || q > 1e-17 {
        utl.Panic("q=%g is incorrect", q)
    }
    if math.Abs(p - a * λ[0]) > 1e-14 {
        utl.Panic("p=%g is incorrect. err = %g", p, math.Abs(p - a * λ[0]))
    }
}
