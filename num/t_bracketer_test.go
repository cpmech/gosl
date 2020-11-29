// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
)

func checkBracket(tst *testing.T, a, b, c, fa, fb, fc float64) {
	if a >= b {
		tst.Errorf("a=%g should be smaller than b=%g\n", a, b)
	}
	if c <= b {
		tst.Errorf("c=%g should be greater than b=%g\n", c, b)
	}
	if fa < fb {
		tst.Errorf("fa=%g should be greater than or equal to fb=%g\n", fa, fb)
	}
	if fc < fb {
		tst.Errorf("fc=%g should be greater than or equal to fb=%g\n", fc, fb)
	}
}

func TestBracket01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Bracket01. quadratic polynomial")

	// function
	ffcn := func(x float64) float64 { return x*x - 1 }

	// bracket minimum
	bracket := NewBracket(ffcn)
	bracket.Verbose = true

	// check
	a, b, c, fa, fb, fc := bracket.Min(0, 1)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-10, 10)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(1, 0)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(10, -10)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-2, -1)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-1, -2)
	checkBracket(tst, a, b, c, fa, fb, fc)
}

func TestBracket02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Bracket02. sin(5x)")

	// function
	ffcn := func(x float64) float64 { return math.Sin(5 * x) }

	// bracket minimum
	bracket := NewBracket(ffcn)
	bracket.Verbose = true

	// check
	a, b, c, fa, fb, fc := bracket.Min(0, 1)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-10, 10)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(1, 0)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(10, -10)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-2, -1)
	checkBracket(tst, a, b, c, fa, fb, fc)

	// check
	a, b, c, fa, fb, fc = bracket.Min(-1, -2)
	checkBracket(tst, a, b, c, fa, fb, fc)
}
