// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"math/cmplx"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func TestAuxFuncs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("AuxFuncs01. auxiliary functions")

	n := 1073741824 // 2³⁰
	io.Pf("n=%10v  IsPowerOfTwo=%v\n", n, IsPowerOfTwo(n))
	if !IsPowerOfTwo(n) {
		tst.Errorf("n=%d is power of 2\n", n)
		return
	}

	n = 1234567
	io.Pf("n=%10v  IsPowerOfTwo=%v\n", n, IsPowerOfTwo(n))
	if IsPowerOfTwo(n) {
		tst.Errorf("n=%d is not power of 2\n", n)
		return
	}

	n = 0
	io.Pf("n=%10v  IsPowerOfTwo=%v\n", n, IsPowerOfTwo(n))
	if IsPowerOfTwo(n) {
		tst.Errorf("n=%d is not power of 2\n", n)
		return
	}

	n = -2
	io.Pf("n=%10v  IsPowerOfTwo=%v\n", n, IsPowerOfTwo(n))
	if IsPowerOfTwo(n) {
		tst.Errorf("n=%d is not power of 2\n", n)
		return
	}

	n = 1 // 2⁰
	io.Pf("n=%10v  IsPowerOfTwo=%v\n", n, IsPowerOfTwo(n))
	if !IsPowerOfTwo(n) {
		tst.Errorf("n=%d is power of 2\n", n)
		return
	}

	a, b := 123.0, 456.0
	io.Pf("a=%v b=%v\n", a, b)
	Swap(&a, &b)
	io.Pf("a=%v b=%v\n", a, b)
	if a == 123 || b == 456 {
		tst.Errorf("Swap failed\n")
		return
	}
}

func TestFourierTrans01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierTrans01. FFT")

	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	err := FourierTransLL(x, false)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	X := la.RCpairsToComplex(x)
	io.Pf("X = %v\n", X)

	y := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}
	Y := dft(y)
	io.Pf("Y = %v\n", Y)

	chk.VectorC(tst, "X", 1e-14, Y, X)
}

// dft compute the discrete Fourier Transform of x (very slow: for testing only)
func dft(x []complex128) (X []complex128) {
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
