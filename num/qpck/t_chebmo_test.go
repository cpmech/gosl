// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qpck

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func runIntegConly(tst *testing.T, N int, f func(x float64) float64, reuse bool) (C []float64) {

	// allocate workspace
	limit := 50
	alist := make([]float64, limit)
	blist := make([]float64, limit)
	rlist := make([]float64, limit)
	elist := make([]float64, limit)
	iord := make([]int32, limit)
	nnlog := make([]int32, limit)

	// set flags
	var icos int32 = 1   // w(x) = cos(m*x)
	var icall int32 = 1  // do not reuse moments
	var maxp1 int32 = 50 // upper bound on the number of Chebyshev moments
	var momcom int32     // 0 => do compute moments

	// allocate Chebyshev moments array
	chebmo := make([]float64, 25*maxp1)

	// calculate all values
	a, b := 0.0, 2.0*math.Pi
	C = make([]float64, N)
	for j := 0; j < N; j++ {

		// coefficient of exp(i⋅m⋅x)
		m := -float64(j)

		// perform integration of cos term
		C[j], _, _, _ = Awoe(0, f, a, b, m, icos, 0, 0, icall, maxp1, alist, blist, rlist, elist, iord, nnlog, momcom, chebmo)

		// set flags
		if reuse {
			icall++
		}

		io.Pf("j=%d mcom=%v cmo=%v\n", j, momcom, chebmo[:10])
	}
	return
}

func runIntegCandS(tst *testing.T, N int, f func(x float64) float64, reuse bool) (C, S []float64) {

	// allocate workspace
	limit := 50
	alist := make([]float64, limit)
	blist := make([]float64, limit)
	rlist := make([]float64, limit)
	elist := make([]float64, limit)
	iord := make([]int32, limit)
	nnlog := make([]int32, limit)

	// set flags
	var icos int32 = 1   // w(x) = cos(m*x)
	var isin int32 = 2   // w(x) = sin(m*x)
	var icall int32 = 1  // do not reuse moments
	var maxp1 int32 = 50 // upper bound on the number of Chebyshev moments
	var momcom int32     // 0 => do compute moments

	// allocate Chebyshev moments array
	chebmo := make([]float64, 25*maxp1)

	// calculate all values
	a, b := 0.0, 2.0*math.Pi
	C = make([]float64, N)
	S = make([]float64, N)
	for j := 0; j < N; j++ {

		// coefficient of exp(i⋅m⋅x)
		m := -float64(j)

		// perform integration of cos term
		C[j], _, _, _ = Awoe(0, f, a, b, m, icos, 0, 0, icall, maxp1, alist, blist, rlist, elist, iord, nnlog, momcom, chebmo)

		// set flags
		if reuse {
			icall++
		}

		// perform integration of sin term
		S[j], _, _, _ = Awoe(0, f, a, b, m, isin, 0, 0, icall, maxp1, alist, blist, rlist, elist, iord, nnlog, momcom, chebmo)
	}
	return
}

func TestChebmo01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Chebmo01. Reusing Chebyshev moments. C only")

	f := func(x float64) float64 { return math.Exp(1.0 - math.Sin(x)) }

	N := 3
	Ca := runIntegConly(tst, N, f, false) // false => do not reuse moments
	Cb := runIntegConly(tst, N, f, true)  // true => reuse moments
	io.Pl()
	io.Pforan("Ca = %v\n", Ca)
	io.Pf("Cb = %v\n", Cb)
	chk.Array(tst, "C", 1e-17, Ca, Cb)
}

func TestChebmo02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Chebmo02. Reusing Chebyshev moments. C and S")

	f := func(x float64) float64 { return math.Exp(1.0 - math.Sin(x)) }

	N := 3
	Ca, Sa := runIntegCandS(tst, N, f, false) // false => do not reuse moments
	Cb, Sb := runIntegCandS(tst, N, f, true)  // true => reuse moments
	io.Pforan("Ca = %v\n", Ca)
	io.Pf("Cb = %v\n", Cb)
	io.Pforan("Sa = %v\n", Sa)
	io.Pf("Sb = %v\n", Sb)
	chk.Array(tst, "C", 1e-17, Ca, Cb)
	chk.Array(tst, "S", 1e-17, Sa, Sb)
}
