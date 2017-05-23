// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_simpson01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("simpson01")

	y := func(x float64) float64 {
		return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
	}

	n := 1000
	A, _ := QuadDiscreteSimpsonRF(0, 1, n, y)
	io.Pforan("A  = %v\n", A)
	Acor := 1.08268158558
	chk.Scalar(tst, "A", 1e-11, A, Acor)
}

func Test_quad01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("quad01")

	y := func(x float64) float64 {
		return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
	}
	var err error
	Acor := 1.08268158558

	// trapezoidal rule
	var T Quadrature
	T = new(Trap)
	T.Init(y, 0, 1, 1e-11)
	A, err := T.Integrate()
	if err != nil {
		io.Pforan(err.Error())
	}
	io.Pforan("A  = %v\n", A)
	chk.Scalar(tst, "A", 1e-11, A, Acor)

	// Simpson's rule
	var S Quadrature
	S = new(Simp)
	S.Init(y, 0, 1, 1e-11)
	A, err = S.Integrate()
	if err != nil {
		io.Pforan(err.Error())
	}
	io.Pforan("A  = %v\n", A)
	chk.Scalar(tst, "A", 1e-11, A, Acor)
}
