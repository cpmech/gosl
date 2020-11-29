// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_hc01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("hc01. hypercube")

	Init(111)

	n := 10
	x := LatinIHS(2, n, 5)

	xchk := [][]int{
		{2, 9, 5, 8, 1, 4, 3, 10, 7, 6},
		{3, 10, 1, 2, 7, 4, 9, 6, 5, 8},
	}
	chk.IntDeep2(tst, "x", x, xchk)

	xmin := []float64{-1.0, 0.0}
	xmax := []float64{1.0, 2.0}
	dx := (xmax[0] - xmin[0]) / float64(n-1)
	dy := (xmax[1] - xmin[1]) / float64(n-1)
	X := HypercubeCoords(x, xmin, xmax)
	chk.Array(tst, "x0", 1e-15, X[0], []float64{-1 + dx, -1 + 8*dx, -1 + 4*dx, -1 + 7*dx, -1, -1 + 3*dx, -1 + 2*dx, -1 + 9*dx, -1 + 6*dx, -1 + 5*dx})
	chk.Array(tst, "x1", 1e-15, X[1], []float64{2 * dy, 9 * dy, 0, dy, 6 * dy, 3 * dy, 8 * dy, 5 * dy, 4 * dy, 7 * dy})
}
