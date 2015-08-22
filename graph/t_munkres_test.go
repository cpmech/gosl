// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_munkres01(tst *testing.T) {

	//verbose()
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

	//verbose()
	chk.PrintTitle("munkres02")

	C := [][]int{
		{2, 1},
		{1, 1},
	}
	var mnk Munkres
	mnk.Init(C)
	mnk.Run()
	chk.Ints(tst, "links", mnk.Links, []int{1, 0}) // 0 goes with 1 and 1 goes with 0

	C = [][]int{
		{2, 2},
		{4, 3},
	}
	mnk.Init(C)
	mnk.Run()
	chk.Ints(tst, "links", mnk.Links, []int{0, 1}) // 0 does 0 and 1 does with 1

	C = [][]int{
		{2, 2},
		{1, 3},
	}
	mnk.Init(C)
	mnk.Run()
	chk.Ints(tst, "links", mnk.Links, []int{1, 0}) // 0 does 1 and 1 does 0

	C = [][]int{
		{2, 1},
		{2, 1},
		{1, 1},
		{1, 1},
	}
	mnk.Init(C)
	mnk.Run()
	chk.Ints(tst, "links", mnk.Links, []int{1, -1, 0, -1}) // 0 goes with 0 and 1 goes with 1 and the others are unconnected

	C = [][]int{
		{1, 2, 3},
		{6, 5, 4},
	}
	mnk.Init(C)
	mnk.Run()
	chk.Ints(tst, "links", mnk.Links, []int{0, 2}) // 0 goes with 0 and 1 goes with 2

	C = [][]int{
		{1, 2, 3},
		{6, 5, 4},
		{1, 1, 1},
	}
	mnk.Init(C)
	mnk.Run()
	chk.Ints(tst, "links", mnk.Links, []int{0, 2, 1}) // 0 goes with 0, 1 goes with 2 and 2 goes with 1

	C = [][]int{
		{2, 4, 7, 9},
		{3, 9, 5, 1},
		{8, 2, 9, 7},
	}
	mnk.Init(C)
	mnk.Run()
	chk.Ints(tst, "links", mnk.Links, []int{0, 3, 1})

	C = [][]int{
		{1, 2, 3},
		{2, 4, 6},
		{3, 6, 9},
	}
	Ccor := [][]int{
		{1, 0, 0},
		{0, 0, 1},
		{0, 1, 3},
	}
	mnk.Init(C)
	mnk.Run()
	chk.IntMat(tst, "C", mnk.C, Ccor)
	chk.Ints(tst, "links", mnk.Links, []int{2, 1, 0}) // 0 goes with 2, 1 goes with 1 and 2 goes with 0

	// from https://projecteuler.net/index.php?section=problems&id=345
	C = [][]int{
		{7, 53, 183, 439, 863},
		{497, 383, 563, 79, 973},
		{287, 63, 343, 169, 583},
		{627, 343, 773, 959, 943},
		{767, 473, 103, 699, 303},
	}
	for i := 0; i < len(C); i++ {
		for j := 0; j < len(C[i]); j++ {
			C[i][j] *= -1
		}
	}
	mnk.Init(C)
	mnk.Run()
	io.Pforan("links = %v\n", mnk.Links)
	chk.Ints(tst, "links", mnk.Links, []int{4, 1, 2, 3, 0})
}

func Test_munkres03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("munkres03. Euler problem 345")

	C := [][]int{
		{7, 53, 183, 439, 863, 497, 383, 563, 79, 973, 287, 63, 343, 169, 583},
		{627, 343, 773, 959, 943, 767, 473, 103, 699, 303, 957, 703, 583, 639, 913},
		{447, 283, 463, 29, 23, 487, 463, 993, 119, 883, 327, 493, 423, 159, 743},
		{217, 623, 3, 399, 853, 407, 103, 983, 89, 463, 290, 516, 212, 462, 350},
		{960, 376, 682, 962, 300, 780, 486, 502, 912, 800, 250, 346, 172, 812, 350},
		{870, 456, 192, 162, 593, 473, 915, 45, 989, 873, 823, 965, 425, 329, 803},
		{973, 965, 905, 919, 133, 673, 665, 235, 509, 613, 673, 815, 165, 992, 326},
		{322, 148, 972, 962, 286, 255, 941, 541, 265, 323, 925, 281, 601, 95, 973},
		{445, 721, 11, 525, 473, 65, 511, 164, 138, 672, 18, 428, 154, 448, 848},
		{414, 456, 310, 312, 798, 104, 566, 520, 302, 248, 694, 976, 430, 392, 198},
		{184, 829, 373, 181, 631, 101, 969, 613, 840, 740, 778, 458, 284, 760, 390},
		{821, 461, 843, 513, 17, 901, 711, 993, 293, 157, 274, 94, 192, 156, 574},
		{34, 124, 4, 878, 450, 476, 712, 914, 838, 669, 875, 299, 823, 329, 699},
		{815, 559, 813, 459, 522, 788, 168, 586, 966, 232, 308, 833, 251, 631, 107},
		{813, 883, 451, 509, 615, 77, 281, 613, 459, 205, 380, 274, 302, 35, 805},
	}
	for i := 0; i < len(C); i++ {
		for j := 0; j < len(C[i]); j++ {
			C[i][j] *= -1
		}
	}

	var mnk Munkres
	mnk.Init(C)
	mnk.Run()
	io.Pforan("links = %v\n", mnk.Links)

	cost := 0
	for i := 0; i < len(C); i++ {
		j := mnk.Links[i]
		cost += -C[i][j]
	}
	io.Pforan("cost = %v  (13938)\n", cost)
	chk.Scalar(tst, "cost", 1e-17, float64(cost), 13938)
	chk.Ints(tst, "links", mnk.Links, []int{9, 10, 7, 4, 3, 0, 13, 2, 14, 11, 6, 5, 12, 8, 1})
}

//13938

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
