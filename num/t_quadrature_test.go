// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"math/cmplx"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestQuadGen01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("QuadGen01. using QUADPACK general function")

	f := func(x float64) float64 { return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0)) }
	A, err := QuadGen(0, 1, 0, f)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	io.Pforan("A  = %v\n", A)
	chk.Float64(tst, "A", 1e-12, A, 1.08268158558)
}

func TestQuadCs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("QuadCs01. using QUADPACK oscillatory function")

	ω := math.Pow(2.0, 3.4)
	f := func(x float64) float64 { return math.Exp(20.0 * (x - 1)) }
	A, err := QuadCs(0, 1, ω, true, 0, f)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	io.Pforan("A  = %v\n", A)
	Aref := (20*math.Sin(ω) - ω*math.Cos(ω) + ω*math.Exp(-20)) / (math.Pow(20, 2) + math.Pow(ω, 2))
	chk.Float64(tst, "A", 1e-16, A, Aref)
}

func TestQuadExpIx01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("QuadExpIx01. ∫ x²⋅exp(i⋅m⋅x) dx using QUADPACK")

	f := func(x float64) float64 { return x * x }
	π := math.Pi
	a := 0.0
	b := 2.0 * π
	m := 4.0

	I, err := QuadExpIx(a, b, m, 0, f)
	chk.EP(err)

	ee := cmplx.Exp(complex(0, 2*π*m))
	π2 := complex(π*π, 0)
	m2 := complex(m*m, 0)
	m3 := complex(m*m*m, 0)
	mπ4 := complex(4*π*m, 0)
	Iana := (2i+mπ4-4i*π2*m2)*ee/m3 - 2i/m3

	chk.AnaNumC(tst, "I", 1e-14, I, Iana, chk.Verbose)
}

func TestQuadExpIx02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("QuadExpIx02. ∫ [p⋅cos(x)+q⋅sin(x)]⋅exp(i⋅m⋅x) dx")

	p := 2.0
	q := 3.0
	f := func(x float64) float64 { return p*math.Cos(x) + q*math.Sin(x) }
	π := math.Pi
	a := 0.0
	b := 2.0 * π
	m := 0.5

	I, err := QuadExpIx(a, b, m, 0, f)
	chk.EP(err)

	ee := cmplx.Exp(complex(0, 2*π*m))
	Q := complex(q, 0)
	d := complex(m*m-1, 0)
	pmi := complex(0, p*m)
	Iana := (ee*Q-pmi*ee)/d - (Q-pmi)/d

	chk.AnaNumC(tst, "I", 1e-15, I, Iana, chk.Verbose)
}
