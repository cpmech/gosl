// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fftw

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// reference results
var test1Xref []complex128
var test2Xref []complex128

// expmix uses Euler's formula to compute exp(-i⋅x) = cos(x) - i⋅sin(x)
func expmix(x float64) complex128 {
	return complex(math.Cos(x), -math.Sin(x))
}

// initialise reference results
func init() {

	x1 := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}
	test1Xref = dft1d(x1)

	x2 := []complex128{1 + 0i, 2 + 0i, 3 + 0i, 4 + 0i, 5 + 0i, 6 + 0i, 7 + 0i, 8 + 0i}
	test2Xref = dft1d(x2)
}

// tests ///////////////////////////////////////////////////////////////////////////////////////////

func TestOneDver01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("OneDver01a.")

	// set input data
	N := 4
	x := make([]complex128, N)
	for i := 0; i < N; i++ {
		ii := float64(i * 2)
		x[i] = complex(ii+1, ii+2)
	}

	// flags
	inverse := false
	measure := false

	// allocate plan
	plan, err := NewPlan1d(x, inverse, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// check plan.data
	chk.VectorC(tst, "plan.data", 1e-15, plan.data, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i})

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("X = %v\n", x)

	// check output
	chk.VectorC(tst, "X", 1e-14, x, test1Xref)
}

func TestOneDver01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("OneDver01b. (measure)")

	// set input data
	N := 4
	x := make([]complex128, N)
	for i := 0; i < N; i++ {
		ii := float64(i * 2)
		x[i] = complex(ii+1, ii+2)
	}

	// flags
	inverse := false
	measure := true

	// allocate plan
	plan, err := NewPlan1d(x, inverse, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// check plan.data
	chk.VectorC(tst, "plan.data", 1e-15, plan.data, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i})

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("X = %v\n", x)

	// check output
	chk.VectorC(tst, "X", 1e-14, x, test1Xref)
}

func TestOneDver02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("OneDver02")

	// set input data
	N := 8
	x := make([]complex128, N)
	for i := 0; i < N; i++ {
		ii := float64(i)
		x[i] = complex(ii+1, 0)
	}

	// flags
	inverse := false
	measure := true

	// allocate plan
	plan, err := NewPlan1d(x, inverse, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// check plan.data
	chk.VectorC(tst, "plan.data", 1e-15, plan.data, []complex128{1, 2, 3, 4, 5, 6, 7, 8})

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("X = %v\n", x)

	// check output
	chk.VectorC(tst, "X", 1e-13, x, test2Xref)
}

// solution ////////////////////////////////////////////////////////////////////////////////////////

// dft1d compute the discrete Fourier Transform of x (very slow: for testing only)
func dft1d(x []complex128) (X []complex128) {
	N := len(x)
	X = make([]complex128, N)
	for n := 0; n < N; n++ {
		for k := 0; k < N; k++ {
			a := 2.0 * math.Pi * float64(k*n) / float64(N)
			X[n] += x[k] * expmix(a) // x[k]⋅exp(-a)
		}
	}
	return
}
