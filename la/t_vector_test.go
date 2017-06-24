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
	chk.Vector(tst, "a (empty)", 1e-17, a, []float64{0, 0, 0})

	a.Fill(123)
	chk.Vector(tst, "a.Fill", 1e-17, a, []float64{123, 123, 123})

	a.ApplyFunc(func(i int, x float64) float64 { return x - 123 })
	chk.Vector(tst, "a.ApplyFunc", 1e-17, a, []float64{0, 0, 0})

	b := NewVectorMapped(3, func(i int) float64 { return float64(i + 1) })
	chk.Vector(tst, "b (mapped)", 1e-17, b, []float64{1, 2, 3})

	c := b.GetCopy()
	chk.Vector(tst, "b.GetCopy", 1e-17, c, []float64{1, 2, 3})

	chk.Scalar(tst, "c.Accum", 1e-17, c.Accum(), 6)

	chk.Scalar(tst, "c.Norm", 1e-17, c.Norm(), math.Sqrt(14.0))

	chk.Scalar(tst, "c.NormDiff(b)", 1e-17, c.NormDiff(b), 0)

	chk.Scalar(tst, "c.Min", 1e-17, c.Min(), 1)

	chk.Scalar(tst, "c.Max", 1e-17, c.Max(), 3)

	min, max := c.MinMax()
	chk.Scalar(tst, "min", 1e-17, min, 1)
	chk.Scalar(tst, "max", 1e-17, max, 3)

	chk.Scalar(tst, "c.Largest", 1e-17, c.Largest(3), 1)

	chk.Scalar(tst, "c.Dot(b)", 1e-17, c.Dot(b), 14)

	d := NewVector(3)
	b.CopyInto(d, -1)
	chk.Vector(tst, "b.CopyInto", 1e-17, d, []float64{-1, -2, -3})

	b.Add(d, 1)
	chk.Vector(tst, "a.Add", 1e-17, d, []float64{0, 0, 0})

	a.Fill(123)
	e := NewVector(3)
	b.Add2(e, -1, 0.1, a)
	chk.Vector(tst, "b.Add2", 1e-17, e, []float64{12.3 - 1, 12.3 - 2, 12.3 - 3})

	b.CopyInto(d, -1)
	io.Pf("b = %v\n", b)
	io.Pf("d = %v\n", d)
	chk.Scalar(tst, "b.MaxDiff", 1e-17, b.MaxDiff(d), 6)

	d.Scale(6, 2)
	chk.Vector(tst, "d.Scale", 1e-17, d, []float64{4, 2, 0})

	b.CopyInto(d, -1)
	f := d.GetScaled(6, 2)
	chk.Vector(tst, "d.GetScaled", 1e-17, f, []float64{4, 2, 0})

	b.CopyInto(d, -1)
	d.Scale(3, 2)
	chk.Vector(tst, "d.Scale (again)", 1e-17, d, []float64{1, -1, -3})

	d.ScaleAbs(0, 1)
	chk.Vector(tst, "d.ScaleAbs", 1e-17, d, []float64{1, 1, 3})

	b.CopyInto(d, -1)
	d.Scale(3, 2)
	d.ScaleAbs(-2, 2)
	chk.Vector(tst, "d.ScaleAbs (again)", 1e-17, d, []float64{0, 0, 4})

	chk.Scalar(tst, "c.Rms", 1e-17, c.Rms(), math.Sqrt(14.0/3.0))

	A, M := 1.0, 2.0
	s := []float64{1, 2, 3}
	den := []float64{A + M*s[0], A + M*s[1], A + M*s[2]}
	div := []float64{c[0] / den[0], c[1] / den[1], c[2] / den[2]}
	sum := (div[0]*div[0] + div[1]*div[1] + div[2]*div[2]) / 3.0
	chk.Scalar(tst, "c.RmsScaled", 1e-17, c.RmsScaled(A, M, s), math.Sqrt(sum))

	z := NewVector(3)
	chk.Scalar(tst, "c.RmsError", 1e-17, c.RmsError(A, M, s, z), math.Sqrt(sum))

	z.Fill(1)
	div2 := []float64{math.Abs(c[0]-z[0]) / den[0], math.Abs(c[1]-z[2]) / den[1], math.Abs(c[2]-z[2]) / den[2]}
	sum2 := (div2[0]*div2[0] + div2[1]*div2[1] + div2[2]*div2[2]) / 3.0
	chk.Scalar(tst, "c.RmsError", 1e-17, c.RmsError(A, M, s, z), math.Sqrt(sum2))
}

func TestVector02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Vector02. (complex) Basic operations")

	a := NewVectorC(3)
	chk.VectorC(tst, "a (empty)", 1e-17, a, []complex128{0, 0, 0})

	a.Fill(123 + 1i)
	chk.VectorC(tst, "a.Fill", 1e-17, a, []complex128{123 + 1i, 123 + 1i, 123 + 1i})

	a.ApplyFunc(func(i int, x complex128) complex128 { return x - (123 + 1i) })
	chk.VectorC(tst, "a.ApplyFunc", 1e-17, a, []complex128{0, 0, 0})

	b := NewVectorMappedC(3, func(i int) complex128 { return complex(float64(i+1), 1) })
	chk.VectorC(tst, "b (mapped)", 1e-17, b, []complex128{1 + 1i, 2 + 1i, 3 + 1i})

	c := b.GetCopy()
	chk.VectorC(tst, "b.GetCopy", 1e-17, c, []complex128{1 + 1i, 2 + 1i, 3 + 1i})

	io.Pforan("c = %v\n", c)
	chk.ScalarC(tst, "c.Norm", 1e-17, c.Norm(), cmplx.Sqrt(2i+(3+4i)+(8+6i)))

	d := NewVectorC(3)
	b.CopyInto(d, -1)
	chk.VectorC(tst, "b.CopyInto", 1e-17, d, []complex128{-1 - 1i, -2 - 1i, -3 - 1i})

	io.Pforan("b = %v\n", b)
	io.Pforan("d = %v\n", d)
	chk.Scalar(tst, "b.MaxDiff", 1e-17, b.MaxDiff(d), 6)

	b.CopyInto(d, -1)
	io.Pforan("d = %v\n", d)
	f := d.GetScaled(6, 2)
	chk.VectorC(tst, "d.GetScaled", 1e-17, f, []complex128{4 - 2i, 2 - 2i, -2i})
}
