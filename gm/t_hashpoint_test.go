// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_hash01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("hash01")

	c1 := []float64{1.0000000001, 2.0000000002}
	c2 := []float64{1.00000000009, 2.0000000002}
	c3 := []float64{1.0000000000002, 2.0000000002, 3.0000000000003}
	c4 := []float64{1.0000000000003, 2.0000000002, 3.0000000000003}

	tol := 1e-11
	xmin := []float64{-10, -1000, -100}
	xmax := []float64{10, 100, 1000}
	xdel := []float64{0, 0, 0}
	for i := 0; i < 3; i++ {
		xdel[i] = xmax[i] - xmin[i]
	}

	h1 := HashPoint(c1, xmin, xdel, tol)
	h2 := HashPoint(c2, xmin, xdel, tol)

	//tol = 1e-14
	tol = 1e-16

	h3 := HashPoint(c3, xmin, xdel, tol)
	h4 := HashPoint(c4, xmin, xdel, tol)

	io.Pforan("h1 = %v\n", h1)
	io.Pforan("h2 = %v\n", h2)
	io.Pforan("h3 = %v\n", h3)
	io.Pforan("h4 = %v\n", h4)

	if h1 == h2 {
		chk.Panic("h1 must not be equal to h2")
	}
	if h1 == h3 {
		chk.Panic("h1 must not be equal to h3")
	}
	if h1 == h4 {
		chk.Panic("h1 must not be equal to h4")
	}
	if h2 == h3 {
		chk.Panic("h2 must not be equal to h3")
	}
	if h2 == h4 {
		chk.Panic("h2 must not be equal to h4")
	}
	if h3 == h4 {
		chk.Panic("h3 must not be equal to h4")
	}
}
