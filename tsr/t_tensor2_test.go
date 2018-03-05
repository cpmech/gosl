// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestTensor2set01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Tensor2set01")

	a := NewTensor2(true, true)
	a.Set(0, 0, 1)
	a.Set(1, 1, 2)
	a.Set(2, 2, 3)
	a.Set(0, 1, 4)
	a.Set(1, 2, 5) // value will be zero, because tensor is 2D
	a.Set(0, 2, 6) // value will be zero, because tensor is 2D
	a.Set(1, 0, 4)
	a.Set(2, 1, 8) // value will be zero, because tensor is 2D
	a.Set(2, 0, 9) // value will be zero, because tensor is 2D
	chk.Array(tst, "a(symm,2D)", 1e-15, a.data, []float64{1, 2, 3, 4 * sq2})
	chk.Float64(tst, "a[0,0]", 1e-15, a.Get(0, 0), 1)
	chk.Float64(tst, "a[1,1]", 1e-15, a.Get(1, 1), 2)
	chk.Float64(tst, "a[2,2]", 1e-15, a.Get(2, 2), 3)
	chk.Float64(tst, "a[0,1]", 1e-15, a.Get(0, 1), 4)
	chk.Float64(tst, "a[1,2]", 1e-15, a.Get(1, 2), 0)
	chk.Float64(tst, "a[0,2]", 1e-15, a.Get(0, 2), 0)
	chk.Float64(tst, "a[1,0]", 1e-15, a.Get(1, 0), 4)
	chk.Float64(tst, "a[2,1]", 1e-15, a.Get(2, 1), 0)
	chk.Float64(tst, "a[2,0]", 1e-15, a.Get(2, 0), 0)

	io.Pl()
	b := NewTensor2(false, true)
	b.Set(0, 0, 1)
	b.Set(1, 1, 2)
	b.Set(2, 2, 3)
	b.Set(0, 1, 4)
	b.Set(1, 2, 5)
	b.Set(0, 2, 6)
	b.Set(1, 0, 7)
	b.Set(2, 1, 8)
	b.Set(2, 0, 9)
	chk.Array(tst, "b(unsymm,2D)", 1e-15, b.data, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	chk.Float64(tst, "b[0,0]", 1e-15, b.Get(0, 0), 1)
	chk.Float64(tst, "b[1,1]", 1e-15, b.Get(1, 1), 2)
	chk.Float64(tst, "b[2,2]", 1e-15, b.Get(2, 2), 3)
	chk.Float64(tst, "b[0,1]", 1e-15, b.Get(0, 1), 4)
	chk.Float64(tst, "b[1,2]", 1e-15, b.Get(1, 2), 5)
	chk.Float64(tst, "b[0,2]", 1e-15, b.Get(0, 2), 6)
	chk.Float64(tst, "b[1,0]", 1e-15, b.Get(1, 0), 7)
	chk.Float64(tst, "b[2,1]", 1e-15, b.Get(2, 1), 8)
	chk.Float64(tst, "b[2,0]", 1e-15, b.Get(2, 0), 9)

	io.Pl()
	c := NewTensor2(true, false)
	c.Set(0, 0, 1)
	c.Set(1, 1, 2)
	c.Set(2, 2, 3)
	c.Set(0, 1, 4)
	c.Set(1, 2, 5)
	c.Set(0, 2, 6)
	c.Set(1, 0, 4)
	c.Set(2, 1, 5)
	c.Set(2, 0, 6)
	chk.Array(tst, "c(symm,3D)", 1e-15, c.data, []float64{1, 2, 3, 4 * sq2, 5 * sq2, 6 * sq2})
	chk.Float64(tst, "c[0,0]", 1e-15, c.Get(0, 0), 1)
	chk.Float64(tst, "c[1,1]", 1e-15, c.Get(1, 1), 2)
	chk.Float64(tst, "c[2,2]", 1e-15, c.Get(2, 2), 3)
	chk.Float64(tst, "c[0,1]", 1e-15, c.Get(0, 1), 4)
	chk.Float64(tst, "c[1,2]", 1e-15, c.Get(1, 2), 5)
	chk.Float64(tst, "c[0,2]", 1e-15, c.Get(0, 2), 6)
	chk.Float64(tst, "c[1,0]", 1e-15, c.Get(1, 0), 4)
	chk.Float64(tst, "c[2,1]", 1e-15, c.Get(2, 1), 5)
	chk.Float64(tst, "c[2,0]", 1e-15, c.Get(2, 0), 6)

	io.Pl()
	d := NewTensor2(false, false)
	d.Set(0, 0, 1)
	d.Set(1, 1, 2)
	d.Set(2, 2, 3)
	d.Set(0, 1, 4)
	d.Set(1, 2, 5)
	d.Set(0, 2, 6)
	d.Set(1, 0, 7)
	d.Set(2, 1, 8)
	d.Set(2, 0, 9)
	chk.Array(tst, "d(unsymm,3D)", 1e-15, d.data, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	chk.Float64(tst, "d[0,0]", 1e-15, d.Get(0, 0), 1)
	chk.Float64(tst, "d[1,1]", 1e-15, d.Get(1, 1), 2)
	chk.Float64(tst, "d[2,2]", 1e-15, d.Get(2, 2), 3)
	chk.Float64(tst, "d[0,1]", 1e-15, d.Get(0, 1), 4)
	chk.Float64(tst, "d[1,2]", 1e-15, d.Get(1, 2), 5)
	chk.Float64(tst, "d[0,2]", 1e-15, d.Get(0, 2), 6)
	chk.Float64(tst, "d[1,0]", 1e-15, d.Get(1, 0), 7)
	chk.Float64(tst, "d[2,1]", 1e-15, d.Get(2, 1), 8)
	chk.Float64(tst, "d[2,0]", 1e-15, d.Get(2, 0), 9)
}
