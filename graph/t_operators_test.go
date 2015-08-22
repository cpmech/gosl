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
	Ccor := [][]int{
		{0, 1, 2},
		{0, 2, 4},
		{0, 3, 6},
	}
	Mcor := [][]Mask_t{
		{STAR, NONE, NONE},
		{NONE, NONE, NONE},
		{NONE, NONE, NONE},
	}

	var mnk Munkres
	mnk.Init(C)

	// 1:
	io.PfYel("1: after step 0:\n")
	io.Pf("%v", mnk.StrCostMatrix())

	// 2: step 1
	next_step := mnk.step1()
	io.PfYel("\n2: after step 1:\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 2)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{false, false, false})

	// 3: step 2
	next_step = mnk.step2()
	io.PfYel("\n3: after step 2:\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 3)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{false, false, false})

	// 4: step 3
	next_step = mnk.step3()
	io.PfYel("\n4: after step 3:\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 4)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, false, false})

	// 5: step 4
	next_step = mnk.step4()
	io.PfYel("\n5: after step 4:\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 6)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, false, false})

	// 6: step 6
	Ccor = [][]int{
		{0, 0, 1},
		{0, 1, 3},
		{0, 2, 5},
	}
	next_step = mnk.step6()
	io.PfYel("\n6: after step 6:\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 4)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, false, false})

	// 7: step 4 again (1)
	Mcor = [][]Mask_t{
		{STAR, PRIM, NONE},
		{PRIM, NONE, NONE},
		{NONE, NONE, NONE},
	}
	next_step = mnk.step4()
	io.PfYel("\n7: after step 4 again (1):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 5)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{true, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{false, false, false})

	// 8: step 5
	Mcor = [][]Mask_t{
		{NONE, STAR, NONE},
		{STAR, NONE, NONE},
		{NONE, NONE, NONE},
	}
	next_step = mnk.step5()
	io.PfYel("\n8: after step 5:\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 3)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{false, false, false})

	// 9: step 3 again (1)
	next_step = mnk.step3()
	io.PfYel("\n9: after step 3 again (1):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 4)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, true, false})

	// 10: step 4 again (2)
	next_step = mnk.step4()
	io.PfYel("\n10: after step 4 again (2):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 6)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, true, false})

	// 11: step 6 again (1)
	Ccor = [][]int{
		{0, 0, 0},
		{0, 1, 2},
		{0, 2, 4},
	}
	next_step = mnk.step6()
	io.PfYel("\n11: after step 6 again (1):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 4)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, true, false})

	// 12: step 4 again (3)
	Mcor = [][]Mask_t{
		{NONE, STAR, PRIM},
		{STAR, NONE, NONE},
		{NONE, NONE, NONE},
	}
	next_step = mnk.step4()
	io.PfYel("\n12: after step 4 again (3):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 6)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{true, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, false, false})

	// 13: step 6 again (2)
	Ccor = [][]int{
		{1, 0, 0},
		{0, 0, 1},
		{0, 1, 3},
	}
	next_step = mnk.step6()
	io.PfYel("\n13: after step 6 again (2):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 4)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{true, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, false, false})

	// 14: step 4 again (4)
	Mcor = [][]Mask_t{
		{NONE, STAR, PRIM},
		{STAR, PRIM, NONE},
		{PRIM, NONE, NONE},
	}
	next_step = mnk.step4()
	io.PfYel("\n14: after step 4 again (4):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 5)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{true, true, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{false, false, false})

	// 15: step 5 again (1)
	Mcor = [][]Mask_t{
		{NONE, NONE, STAR},
		{NONE, STAR, NONE},
		{STAR, NONE, NONE},
	}
	next_step = mnk.step5()
	io.PfYel("\n15: after step 5 again (1):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 3)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{false, false, false})

	// 15: step 3 again (2)
	next_step = mnk.step3()
	io.PfYel("\n15: after step 3 again (2):\n")
	io.Pf("%v", mnk.StrCostMatrix())
	chk.IntAssert(next_step, 7)
	chk.IntMat(tst, "C", mnk.C, Ccor)
	check_mask_matrix(tst, "M", mnk.M, Mcor)
	chk.Bools(tst, "row_covered", mnk.row_covered, []bool{false, false, false})
	chk.Bools(tst, "col_covered", mnk.col_covered, []bool{true, true, true})
}

func Test_munkres02(tst *testing.T) {

	verbose()
	chk.PrintTitle("munkres02")

	C := [][]int{
		{1, 2, 3},
		{2, 4, 6},
		{3, 6, 9},
	}
	Ccor := [][]int{
		{1, 0, 0},
		{0, 0, 1},
		{0, 1, 3},
	}

	var mnk Munkres
	mnk.Init(C)
	mnk.Run()
	chk.IntMat(tst, "C", mnk.C, Ccor)
	chk.Ints(tst, "links", mnk.Links, []int{
		2, // 0 goes with 2
		1, // 1 goes with 1
		0, // 2 goes with 0
	})
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
