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

var data1dRefA []float64    // reference results A
var data1dRefB []complex128 // reference results B

func init() {
	xA := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}
	xB := []complex128{1 + 0i, 2 + 0i, 3 + 0i, 4 + 0i, 5 + 0i, 6 + 0i, 7 + 0i, 8 + 0i}
	data1dRefA = la.ComplexToRCpairs(dft1d(xA))
	data1dRefB = dft1d(xB)
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
	chk.Vector(tst, "output: X", 1e-14, plan.Xout, data1dRefA)
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
	chk.Vector(tst, "output: X", 1e-14, plan.Xout, data1dRefA)
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
	chk.Vector(tst, "output: X", 1e-14, plan.Xout, data1dRefA)
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
	chk.Vector(tst, "output: X", 1e-14, plan.Xout, data1dRefA)
}

func TestOneDver05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("OneDver05. real input. internal Xin")

	// flags
	inverse := false
	measure := false

	// allocate plan
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	N := len(x) // not needed in the next line
	plan, err := NewPlan1dReal(x, N, inverse, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// check input data
	chk.Vector(tst, "input: x", 1e-15, plan.Xin, []float64{1, 2, 3, 4, 5, 6, 7, 8})

	// check that input is the same as 'x'
	plan.Xin[0] = 123
	chk.Scalar(tst, "Xin[0] has changed     ", 1e-15, plan.Xin[0], 123)
	chk.Scalar(tst, "x should be changed    ", 1e-15, x[0], 123)
	plan.Xin[0] = 1
	chk.Scalar(tst, "Xin[0] has changed back", 1e-15, plan.Xin[0], 1)
	chk.Scalar(tst, "x has changed back     ", 1e-15, x[0], 1)

	// perform Fourier transform
	plan.Execute()

	// check that 'x' hasn't changed
	chk.Vector(tst, "after: x", 1e-15, plan.Xin, []float64{1, 2, 3, 4, 5, 6, 7, 8})

	// print output
	X := plan.GetOutput()
	for i := 0; i < len(X); i++ {
		io.Pf("X: %+15.11f", X[i])
		io.Pf(" ⇒ %+15.11f", data1dRefB[i])
		if math.Abs(real(X[i])-real(data1dRefB[i])) < 1e-13 ||
			math.Abs(imag(X[i])-imag(data1dRefB[i])) < 1e-13 {
			io.PfGreen(" OK\n")
		} else {
			io.PfRed(" fail\n")
		}
	}

	// check output
	chk.VectorC(tst, "output: X", 1e-13, X, data1dRefB)
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
