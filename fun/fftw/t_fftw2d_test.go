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
var test2d1Xref [][]complex128

// initialise reference results
func init() {
	x1 := [][]complex128{
		{0 + 1i, 2 + 3i, 4 + 5i, 6 + 7i},
		{8 + 9i, 10 + 11i, 12 + 13i, 14 + 15i},
	}
	test2d1Xref = dft2d(x1)
}

// tests ///////////////////////////////////////////////////////////////////////////////////////////

func TestTwoDver01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TwoDver01a.")

	// allocate input data
	N0, N1 := 2, 4
	x := make([]complex128, N0*N1)

	// flags
	inverse := false
	measure := false

	// allocate plan
	plan, err := NewPlan2d(N0, N1, x, inverse, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// set input data
	k := 0
	for i := 0; i < N0; i++ {
		for j := 0; j < N1; j++ {
			plan.Set(i, j, complex(float64(k), float64(k+1)))
			k += 2
		}
	}

	// check input data
	chk.VectorC(tst, "plan.data", 1e-17, plan.data, []complex128{0 + 1i, 2 + 3i, 4 + 5i, 6 + 7i, 8 + 9i, 10 + 11i, 12 + 13i, 14 + 15i})

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("X = %v\n", x)

	// check output
	X := plan.GetSlice()
	chk.MatrixC(tst, "X", 1e-13, X, test2d1Xref)
}

func TestTwoDver01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TwoDver01b. (measure)")

	// allocate input data
	N0, N1 := 2, 4
	x := make([]complex128, N0*N1)

	// flags
	inverse := false
	measure := true

	// allocate plan
	plan, err := NewPlan2d(N0, N1, x, inverse, measure)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	defer plan.Free()

	// set input data
	k := 0
	for i := 0; i < N0; i++ {
		for j := 0; j < N1; j++ {
			plan.Set(i, j, complex(float64(k), float64(k+1)))
			k += 2
		}
	}

	// check input data
	chk.VectorC(tst, "plan.data", 1e-17, plan.data, []complex128{0 + 1i, 2 + 3i, 4 + 5i, 6 + 7i, 8 + 9i, 10 + 11i, 12 + 13i, 14 + 15i})

	// perform Fourier transform
	plan.Execute()

	// print output
	io.Pf("X = %v\n", x)

	// check output
	X := plan.GetSlice()
	chk.MatrixC(tst, "X", 1e-13, X, test2d1Xref)
}

// solution ////////////////////////////////////////////////////////////////////////////////////////

// dft2d compute the discrete Fourier Transform of x (very slow: for testing only)
func dft2d(x [][]complex128) (X [][]complex128) {
	N0 := len(x)
	N1 := len(x[0])
	X = make([][]complex128, N0)
	for l0 := 0; l0 < N0; l0++ {
		X[l0] = make([]complex128, N1)
		for l1 := 0; l1 < N1; l1++ {
			for k0 := 0; k0 < N0; k0++ {
				for k1 := 0; k1 < N1; k1++ {
					a := 2.0 * math.Pi * float64(k0*l0) / float64(N0)
					b := 2.0 * math.Pi * float64(k1*l1) / float64(N1)
					X[l0][l1] += x[k0][k1] * expmix(a) * expmix(b)
				}
			}
		}
	}
	return
}
