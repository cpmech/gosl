// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

func TestLineSearch01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LineSearch01.")

	// function
	ffcn := func(x la.Vector) float64 {
		return x[0]*x[0] + x[1]*x[1] - 0.5
	}

	// Jacobian
	Jfcn := func(dfdx, x la.Vector) {
		dfdx[0] = 2.0 * x[0]
		dfdx[1] = 2.0 * x[1]
	}

	// initial point and direction
	x0 := la.NewVectorSlice([]float64{-2, -2})
	x := x0.GetCopy()
	u := la.NewVectorSlice([]float64{4, 4})

	// solver
	line := NewLineSearch(2, ffcn, Jfcn)

	// set params
	line.SetParams(utl.NewParams(
		&utl.P{N: "maxitls", V: 2},
		&utl.P{N: "maxitzoom", V: 2},
		&utl.P{N: "maxalpha", V: 100},
		&utl.P{N: "mulalpha", V: 2},
		&utl.P{N: "coef1", V: 1e-4},
		&utl.P{N: "coef2", V: 0.4},
		&utl.P{N: "coefquad", V: 0.1},
		&utl.P{N: "coefcubic", V: 0.2},
	))

	// solve
	a, f := line.Wolfe(x, u, false, 0)
	io.Pforan("a = %v\n", a)
	io.Pforan("f = %v\n", f)
	chk.Float64(tst, "a", 1e-15, a, 0.5)
	chk.Float64(tst, "f", 1e-15, a, 0.5)
}
