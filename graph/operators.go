// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import "github.com/cpmech/gosl/chk"

// Match returns the combinations resulting in least cost
//  Example:
//            $   | Clean  Sweep   Wash
//         -------|--------------------
//         Fry    |   [2]      3      3
//         Leela  |     3    [2]      3
//         Bender |     3      3    [2]
//         Hermes |     4      3    [1]
//         optimal cost  = 7
//  Input:
//    cost -- [m][n] cost matrix
//  Output:
//    pairs -- [m][2] (assignments)
func Match(pairs [][]int, cost [][]float64) (optcost float64) {
	m := len(cost)
	n := len(cost[0])
	if m == 2 && n == 2 {
		c0 := cost[0][0] + cost[1][1]
		c1 := cost[0][1] + cost[1][0]
		if c0 < c1 {
			pairs[0][0], pairs[0][1] = 0, 0 // 0 does 0
			pairs[1][0], pairs[1][1] = 1, 1 // 1 does 1
			optcost = c0
		} else {
			pairs[0][0], pairs[0][1] = 0, 1 // 0 does 1
			pairs[1][0], pairs[1][1] = 1, 0 // 1 does 0
			optcost = c1
		}
		return
	}
	chk.Panic("Match cannot handle cost %d x %d matrices yet", m, n)

	/*
		done := false
		for !done {
			switch step {
			case 1:
				step = munkres_step_one(step)
			case 2:
				step = munkres_step_two(step)
			case 3:
				step = munkres_step_three(step)
			case 4:
				step = munkres_step_four(step)
			case 5:
				step = munkres_step_five(step)
			case 6:
				step = munkres_step_six(step)
			case 7:
				step = munkres_step_seven(step)
				done = true
			}
		}
	*/

	return
}

/*
func munkres_step_one(step int) (next_step int) {
	var min_in_row float64
	for r := 0; r < nrow; r++ {
		min_in_row = C[r][0]
		for c := 0; c < ncol; c++ {
			min_in_row = utl.Min(min_in_row, C[r][c])
		}
		for c := 0; c < ncol; c++ {
			C[r][c] -= min_in_row
		}
	}
	return 2
}
*/
