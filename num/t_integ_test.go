// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func Test_trapz01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("trapz01")

	x := []float64{4, 6, 8}
	y := []float64{1, 2, 3}
	A := Trapz(x, y)
	chk.Scalar(tst, "A", 1e-17, A, 8)
}

func Test_trapz02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("trapz02")

	y := func(x float64) float64 {
		return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
	}

	n := 11
	x := utl.LinSpace(0, 1, n)
	A := TrapzF(x, y)
	A_ := TrapzRange(0, 1, n, y)
	io.Pforan("A  = %v\n", A)
	io.Pforan("A_ = %v\n", A_)
	Acor := 1.08306090851465
	chk.Scalar(tst, "A", 1e-15, A, Acor)
	chk.Scalar(tst, "A_", 1e-15, A_, Acor)
}
