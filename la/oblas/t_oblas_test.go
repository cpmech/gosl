// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_matrix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("matrix01")

	A := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 0, -1, -2},
	}

	a := NewMatrix(3, 4)
	a.SetFromMat(A)
	chk.Vector(tst, "A to a", 1e-15, a.data, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})

	Aback := a.GetMat()
	chk.Matrix(tst, "a to A", 1e-15, Aback, A)
}

func Test_daxpy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Daxpy01")

	α := 0.5
	x := []float64{20, 10, 30, 123, 123}
	y := []float64{-15, -5, -24, 666, 666, 666}
	n, incx, incy := 3, 1, 1
	err := Daxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}

	chk.Vector(tst, "x", 1e-15, x, []float64{20, 10, 30, 123, 123})
	chk.Vector(tst, "y", 1e-15, y, []float64{-5, 0, -9, 666, 666, 666})
}

func Test_zaxpy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zaxpy01")

	α := 1.0 + 0i
	x := []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i}
	y := []complex128{-15 + 1.5i, -5 - 2i, -24 + 1i, 666 - 0.5i, 666 + 5i}
	n, incx, incy := len(x), 1, 1
	err := Zaxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}

	chk.VectorC(tst, "x", 1e-15, x, []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i})
	chk.VectorC(tst, "y", 1e-15, y, []complex128{5 + 2.5i, 5, 6 + 2.5i, 543, 543 + 5.5i})

	α = 0.5 + 1i
	err = Zaxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "y", 1e-15, y, []complex128{14.0 + 23.i, 8.0 + 11.i, 19.5 + 33.25i, 481.0 - 122.75i, 481.0 - 117.25i})
}

func Test_dgemv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemv01")

	m, n := 4, 3
	a := NewMatrix(m, n)
	a.SetFromMat([][]float64{
		{0.1, 0.2, 0.3},
		{1.0, 0.2, 0.3},
		{2.0, 0.2, 0.3},
		{3.0, 0.2, 0.3},
	})

	chk.Vector(tst, "a.data", 1e-15, a.data, []float64{0.1, 1, 2, 3, 0.2, 0.2, 0.2, 0.2, 0.3, 0.3, 0.3, 0.3})

	α, β := 0.5, 2.0
	x := []float64{20, 10, 30}
	y := []float64{3, 1, 2, 4}
	lda, incx, incy := m, 1, 1
	err := Dgemv(false, m, n, α, a, lda, x, incx, β, y, incy)
	if err != nil {
		tst.Errorf("Dgemv failed:\n%v\n", err)
		return
	}
	chk.Vector(tst, "y", 1e-15, y, []float64{12.5, 17.5, 29.5, 43.5})

	err = Dgemv(true, m, n, α, a, lda, y, incy, β, x, incx)
	if err != nil {
		tst.Errorf("Dgemv failed:\n%v\n", err)
		return
	}
	chk.Vector(tst, "x", 1e-15, x, []float64{144.125, 30.3, 75.45})
}
