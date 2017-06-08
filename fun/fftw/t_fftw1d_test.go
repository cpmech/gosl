// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fftw

import (
	"math"
	"math/cmplx"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

var data1dRef []float64 // reference results

func init() {
	x := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}
	data1dRef = la.ComplexToRCpairs(dft1d(x))
}

func TestOneDver01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("OneDver01. allocate Xin internally")

	// flags
	inverse := false
	inplace := false
	measure := false

	// allocate plan
	N := 4
	plan, err := NewPlan1d(nil, N, inverse, inplace, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// set input data
	for i := 0; i < N; i++ {
		ii := float64(i * 2)
		plan.Input(i, complex(ii+1, ii+2))
	}
	chk.Vector(tst, "input: x", 1e-15, plan.Xin, []float64{1, 2, 3, 4, 5, 6, 7, 8})

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("output = ")
	for i := 0; i < N; i++ {
		io.Pf("%v ", plan.Output(i))
	}
	io.Pl()

	// check output
	chk.Vector(tst, "output: X", 1e-14, plan.Xout, data1dRef)
}

func TestOneDver02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("OneDver02. using Xin externally")

	// flags
	inverse := false
	inplace := false
	measure := false

	// allocate plan
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	N := len(x) / 2 // not needed in the next line
	plan, err := NewPlan1d(x, N, inverse, inplace, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// check input data
	chk.Vector(tst, "input: x", 1e-15, plan.Xin, []float64{1, 2, 3, 4, 5, 6, 7, 8})

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("output = ")
	for i := 0; i < N; i++ {
		io.Pf("%v ", plan.Output(i))
	}
	io.Pl()

	// check output
	chk.Vector(tst, "output: X", 1e-14, plan.Xout, data1dRef)
}

func TestOneDver03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("OneDver03. using Xin externally and 'in-place'")

	// flags
	inverse := false
	inplace := true
	measure := false

	// allocate plan
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	N := len(x) / 2 // not needed in the next line
	plan, err := NewPlan1d(x, N, inverse, inplace, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// check input data
	chk.Vector(tst, "input: x", 1e-15, plan.Xin, []float64{1, 2, 3, 4, 5, 6, 7, 8})

	// check output array == should be the same as Xin
	chk.Vector(tst, "input: x", 1e-15, plan.Xin, plan.Xout)
	plan.Xout[0] = 123
	chk.Scalar(tst, "Xin[0] has changed", 1e-15, plan.Xin[0], 123)
	plan.Xout[0] = 1
	chk.Scalar(tst, "Xin[0] has changed", 1e-15, plan.Xin[0], 1)

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("output = ")
	for i := 0; i < N; i++ {
		io.Pf("%v ", plan.Output(i))
	}
	io.Pl()

	// check output
	chk.Vector(tst, "output: X", 1e-14, plan.Xout, data1dRef)
}

func TestOneDver04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("OneDver04. using Xin externally and 'measure'")

	// flags
	inverse := false
	inplace := false
	measure := true

	// allocate plan
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	N := len(x) / 2 // not needed in the next line
	plan, err := NewPlan1d(x, N, inverse, inplace, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// check input data
	chk.Vector(tst, "input: x", 1e-15, plan.Xin, []float64{1, 2, 3, 4, 5, 6, 7, 8})

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("output = ")
	for i := 0; i < N; i++ {
		io.Pf("%v ", plan.Output(i))
	}
	io.Pl()

	// check output
	chk.Vector(tst, "output: X", 1e-14, plan.Xout, data1dRef)
}

// dft1d compute the discrete Fourier Transform of x (very slow: for testing only)
func dft1d(x []complex128) (X []complex128) {
	N := len(x)
	X = make([]complex128, N)
	for n := 0; n < N; n++ {
		for k := 0; k < N; k++ {
			a := 2.0 * math.Pi * float64(k*n) / float64(N)
			X[n] += x[k] * cmplx.Exp(-1i*complex(a, 0))
		}
	}
	return
}
