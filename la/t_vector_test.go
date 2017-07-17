// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"math/cmplx"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestVector01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Vector01. (real) Basic operations")

	a := NewVector(3)
	chk.Array(tst, "a (empty)", 1e-17, a, []float64{0, 0, 0})

	a.Fill(123)
	chk.Array(tst, "a.Fill", 1e-17, a, []float64{123, 123, 123})

	a.ApplyFunc(func(i int, x float64) float64 { return x - 123 })
	chk.Array(tst, "a.ApplyFunc", 1e-17, a, []float64{0, 0, 0})

	b := NewVectorMapped(3, func(i int) float64 { return float64(i + 1) })
	chk.Array(tst, "b (mapped)", 1e-17, b, []float64{1, 2, 3})

	c := b.GetCopy()
	chk.Array(tst, "b.GetCopy", 1e-17, c, []float64{1, 2, 3})

	chk.Scalar(tst, "c.Accum", 1e-17, c.Accum(), 6)

	chk.Scalar(tst, "c.Norm", 1e-17, c.Norm(), math.Sqrt(14.0))

	chk.Scalar(tst, "c.NormDiff(b)", 1e-17, c.NormDiff(b), 0)

	chk.Scalar(tst, "c.Min", 1e-17, c.Min(), 1)

	chk.Scalar(tst, "c.Max", 1e-17, c.Max(), 3)

	min, max := c.MinMax()
	chk.Scalar(tst, "min", 1e-17, min, 1)
	chk.Scalar(tst, "max", 1e-17, max, 3)

	chk.Scalar(tst, "c.Largest", 1e-17, c.Largest(3), 1)
}

func TestVector02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Vector02. (complex) Basic operations")

	a := NewVectorC(3)
	chk.ArrayC(tst, "a (empty)", 1e-17, a, []complex128{0, 0, 0})

	a.Fill(123 + 1i)
	chk.ArrayC(tst, "a.Fill", 1e-17, a, []complex128{123 + 1i, 123 + 1i, 123 + 1i})

	a.ApplyFunc(func(i int, x complex128) complex128 { return x - (123 + 1i) })
	chk.ArrayC(tst, "a.ApplyFunc", 1e-17, a, []complex128{0, 0, 0})

	b := NewVectorMappedC(3, func(i int) complex128 { return complex(float64(i+1), 1) })
	chk.ArrayC(tst, "b (mapped)", 1e-17, b, []complex128{1 + 1i, 2 + 1i, 3 + 1i})

	c := b.GetCopy()
	chk.ArrayC(tst, "b.GetCopy", 1e-17, c, []complex128{1 + 1i, 2 + 1i, 3 + 1i})

	io.Pforan("c = %v\n", c)
	chk.ScalarC(tst, "c.Norm", 1e-17, c.Norm(), cmplx.Sqrt(2i+(3+4i)+(8+6i)))
}
