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

	C := [][]int{
		{1, 2, 3},
		{2, 4, 6},
		{3, 6, 9},
	}

	var mnk Munkres
	mnk.Init(C)

	next_step := mnk.step1()
	chk.IntAssert(next_step, 2)
	chk.IntMat(tst, "after step1: C", mnk.C, [][]int{
		{0, 1, 2},
		{0, 2, 4},
		{0, 3, 6},
	})
	chk.Bools(tst, "after step1: row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "after step1: col_covered", mnk.col_covered, []bool{false, false, false})

	next_step = mnk.step2()
	chk.IntAssert(next_step, 3)
	check_mask_matrix(tst, "after step1: M", mnk.M, [][]Mask_t{
		{STAR, NONE, NONE},
		{NONE, NONE, NONE},
		{NONE, NONE, NONE},
	})
	chk.Bools(tst, "after step1: row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "after step1: col_covered", mnk.col_covered, []bool{false, false, false})

	//mnk.Run(chk.Verbose)
}

func check_mask_matrix(tst *testing.T, msg string, res, correct [][]Mask_t) {
	if len(res) != len(correct) {
		io.Pf("%s [1;31merror len(res)=%d != len(correct)=%d[0m\n", msg, len(res), len(correct))
		tst.Errorf("[1;31m%s failed: res and correct matrices have different lengths. %d != %d[0m", msg, len(res), len(correct))
		return
	}
	for i := 0; i < len(res); i++ {
		if len(res[i]) != len(correct[i]) {
			io.Pf("%s [1;31merror len(res[%d])=%d != len(correct[%d])=%d[0m\n", msg, i, len(res[i]), i, len(correct[i]))
			tst.Errorf("[1;31m%s failed: matrices have different number of columns[0m", msg)
			return
		}
		for j := 0; j < len(res[i]); j++ {
			if res[i][j] != correct[i][j] {
				io.Pf("%s [1;31merror [%d,%d] %v != %v[0m\n", msg, i, j, res[i][j], correct[i][j])
				tst.Errorf("[1;31m%s failed: different int matrices:\n [%d,%d] item is wrong: %v != %v[0m", msg, i, j, res[i][j], correct[i][j])
				return
			}
		}
	}
	chk.PrintOk(msg)
}
