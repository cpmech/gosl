// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
    "math"
    "testing"
    "code.google.com/p/gosl/utl"
)

func Test_trapz01(tst *testing.T) {

    prevTs := utl.Tsilent
    defer func() {
        utl.Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //utl.Tsilent = false
    utl.TTitle("trapz01")

    x := []float64{4,6,8}
    y := []float64{1,2,3}
    A := Trapz(x, y)
    utl.CheckScalar(tst, "A", 1e-17, A, 8)
}

func Test_trapz02(tst *testing.T) {

    prevTs := utl.Tsilent
    defer func() {
        utl.Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //utl.Tsilent = false
    utl.TTitle("trapz02")

    y := func(x float64) float64 {
        return math.Sqrt(1.0+math.Pow(math.Sin(x), 3.0))
    }

    n  := 11
    x  := utl.LinSpace(0, 1, n)
    A  := TrapzF(x, y)
    A_ := TrapzRange(0, 1, n, y)
    utl.Pforan("A  = %v\n", A)
    utl.Pforan("A_ = %v\n", A_)
    Acor := 1.08306090851465
    utl.CheckScalar(tst, "A",  1e-15, A,  Acor)
    utl.CheckScalar(tst, "A_", 1e-15, A_, Acor)
}
