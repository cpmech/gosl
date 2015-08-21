// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func Test_match01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("match01")

	cost := [][]float64{
		{2, 2},
		{4, 3},
	}
	m := len(cost)
	pairs := utl.IntsAlloc(m, 2)
	opt := Match(pairs, cost)
	chk.Scalar(tst, "optimum cost", 1e-17, opt, 5)
	io.Pforan("pairs=%v  opt=%v\n", pairs, opt)
	chk.IntMat(tst, "pairs", pairs, [][]int{
		{0, 0}, // 0 does 0
		{1, 1}, // 1 does 1
	})

	cost[1][0] = 1
	opt = Match(pairs, cost)
	io.Pforan("\npairs=%v  opt=%v\n", pairs, opt)
	chk.Scalar(tst, "optimum cost", 1e-17, opt, 3)
	chk.IntMat(tst, "pairs", pairs, [][]int{
		{0, 1}, // 0 does 1
		{1, 0}, // 1 does 0
	})
}

func Test_munkres01(tst *testing.T) {

	verbose()
	chk.PrintTitle("munkres01")

	var mnk Munkres
	mnk.Init([][]int{
		{1, 2, 3},
		{2, 4, 6},
		{3, 6, 9},
	})

	mnk.Run(chk.Verbose)
}
