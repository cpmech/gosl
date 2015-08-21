// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

type Mask_t int

const (
	NONE Mask_t = iota
	STAR
	PRIM
)

// Munkres (Hungarian algorithm) method to solve the assignment problem
//  based on code by Bob Pilgrim from http://csclab.murraystate.edu/bob.pilgrim/445/munkres.html
//  Note: this method runs in O(n²), in the worst case; therefore is not efficient for large matrix
type Munkres struct {
	C           [][]int    // cost matrix
	M           [][]Mask_t // mask matrix. If Mij==1, then Cij is a starred zero. If Mij==2, then Cij is a primed zero
	path        [][]int
	row_covered []bool // indicates whether a row is covered or not
	col_covered []bool // indicates whether a column is covered or not
	nrow        int    // number of rows in cost/mask matrix
	ncol        int    // number of column in cost/mask matrix
	path_count  int
	path_row_0  int
	path_col_0  int
}

// Init initialises Munkres' structure
func (o *Munkres) Init(C [][]int) {
	chk.IntAssertLessThan(1, len(C))
	o.C = C
	o.nrow, o.ncol = len(o.C), len(o.C[0])
	o.M = make([][]Mask_t, o.nrow)
	for i := 0; i < o.nrow; i++ {
		o.M[i] = make([]Mask_t, o.ncol)
	}
	npath := 2*o.nrow + 1 // TODO: check this
	o.path = utl.IntsAlloc(npath, 2)
	o.row_covered = make([]bool, o.nrow)
	o.col_covered = make([]bool, o.ncol)
}

// step1: for each row of the cost matrix, find the smallest element and subtract it from every
// element in its row. next_step = 2
func (o *Munkres) step1() (next_step int) {
	var xmin int
	for i := 0; i < o.nrow; i++ {
		xmin = o.C[i][0]
		for j := 1; j < o.ncol; j++ {
			xmin = utl.Imin(xmin, o.C[i][j])
		}
		for j := 0; j < o.ncol; j++ {
			o.C[i][j] -= xmin
		}
	}
	return 2
}

// step2: find a zero (Z) in the resulting matrix. If there is no starred zero in its row or column,
// star Z. Repeat for each element in the matrix. Check to see if Cij is a zero value and if its
// column or row is not already covered. If not, then star this zero and cover its row and column.
// Uncover all rows and columns before leaving. next_step = 3
func (o *Munkres) step2() (next_step int) {
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			if !o.row_covered[i] && !o.col_covered[j] {
				if o.C[i][j] == 0 {
					o.M[i][j] = STAR
					o.row_covered[i] = true
					o.col_covered[j] = true
				}
			}
		}
	}
	for i := 0; i < o.nrow; i++ {
		o.row_covered[i] = false
	}
	for j := 0; j < o.ncol; j++ {
		o.col_covered[j] = false
	}
	return 3
}

// step3: cover each column containing a starred zero. If min(n,m) columns are covered, the starred
// zeros describe a complete set of unique assignments and the process is completed (next_step=7);
// otherwise next_step=4.
func (o *Munkres) step3() (next_step int) {
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			if o.M[i][j] == STAR {
				o.col_covered[j] = true
			}
		}
	}
	colcount := 0
	for j := 0; j < o.ncol; j++ {
		if o.col_covered[j] {
			colcount += 1
		}
	}
	next_step = 4
	if colcount >= o.ncol || colcount >= o.nrow {
		next_step = 7
	}
	return
}

// step4: find a noncovered zero and prime it. If there is no starred zero in the row containing
// this primed zero, go to Step 5. Otherwise, cover this row and uncover the column containing the
// starred zero. Continue in this manner until there are no uncovered zeros left. Save the smallest
// uncovered value and go to Step 6.
func (o *Munkres) step4() (next_step int) {
	row, col, done := -1, -1, false
	for !done {
		row, col = o.find_a_zero()
		if row == -1 {
			done = true
			next_step = 6
		} else {
			o.M[row][col] = PRIM
			col_star := o.find_star_in_row(row)
			if col_star >= 0 {
				col = col_star
				o.row_covered[row] = true
				o.col_covered[col] = false
			} else {
				done = true
				next_step = 5
				o.path_row_0 = row
				o.path_col_0 = col
			}
		}
	}
	return
}

// step5: construct a series of alternating primed and starred zeros as follows. Let Z0 represent
// the uncovered primed zero found in Step 4. Let Z1 denote the starred zero in the column of Z0 (if
// any). Let Z2 denote the primed zero in the row of Z1 (there will always be one). Continue until
// the series terminates at a primed zero that has no starred zero in its column. Unstar each
// starred zero of the series, star each primed zero of the series, erase all primes and uncover
// every line in the matrix. Return to Step 3.
func (o *Munkres) step5() (next_step int) {

	// construct series
	r, c, done := -1, -1, false
	o.path_count = 1
	o.path[o.path_count-1][0] = o.path_row_0
	o.path[o.path_count-1][1] = o.path_col_0
	for !done {
		r = o.find_star_in_col(o.path[o.path_count-1][1])
		if r > -1 {
			o.path_count += 1
			o.path[o.path_count-1][0] = r
			o.path[o.path_count-1][1] = o.path[o.path_count-2][1]
		} else {
			done = true
		}
		if !done {
			c = o.find_prime_in_row(o.path[o.path_count-1][0])
			o.path_count += 1
			o.path[o.path_count-1][0] = o.path[o.path_count-2][0]
			o.path[o.path_count-1][1] = c
		}
	}

	// augment_path
	for p := 0; p < o.path_count; p++ {
		if o.M[o.path[p][0]][o.path[p][1]] == STAR {
			o.M[o.path[p][0]][o.path[p][1]] = NONE
		} else {
			o.M[o.path[p][0]][o.path[p][1]] = STAR
		}
	}

	// clear_covers
	for i := 0; i < o.nrow; i++ {
		o.row_covered[i] = false
	}
	for j := 0; j < o.ncol; j++ {
		o.col_covered[j] = false
	}

	// erase_primes
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			if o.M[i][j] == PRIM {
				o.M[i][j] = NONE
			}
		}
	}
	return 3
}

// step6: add the value found in Step 4 to every element of each covered row, and subtract it from
// every element of each uncovered column. Return to Step 4 without altering any stars, primes, or
// covered lines.
func (o *Munkres) step6() (next_step int) {

	// find_smallest
	minval := math.MaxInt64
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			if !o.row_covered[i] && !o.col_covered[j] {
				if minval > o.C[i][j] {
					minval = o.C[i][j]
				}
			}
		}
	}

	// add value from step 4
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			if o.row_covered[i] {
				o.C[i][j] += minval
			}
			if !o.col_covered[j] {
				o.C[i][j] -= minval
			}
		}
	}
	return 4
}

func (o *Munkres) Run(verbose bool) {
	k := 0
	step := 1
	done := false
	for !done {
		if verbose {
			io.Pf("\n%2d: after step %d\n", k+1, step-1)
			o.print_cost_matrix()
		}
		switch step {
		case 1:
			step = o.step1() // returns 2
		case 2:
			step = o.step2() // returns 3
		case 3:
			step = o.step3() // returns 4 or 7
		case 4:
			step = o.step4() // returns 5 or 6
		case 5:
			step = o.step5() // returns 3
		case 6:
			step = o.step6() // returns 4
		case 7:
			done = true
		}
		k++
	}
}

func (o *Munkres) print_cost_matrix() {
	io.Pf("%6v", " ")
	for j := 0; j < o.ncol; j++ {
		io.Pf("%6v", o.col_covered[j])
	}
	io.Pf("\n")
	for i := 0; i < o.nrow; i++ {
		io.Pf("%6v", o.row_covered[i])
		for j := 0; j < o.ncol; j++ {
			switch o.M[i][j] {
			case NONE:
				io.Pf("%6s", io.Sf("%d ", o.C[i][j]))
			case STAR:
				io.Pf("%6s", io.Sf("%d*", o.C[i][j]))
			case PRIM:
				io.Pf("%6s", io.Sf("%d'", o.C[i][j]))
			}
		}
		io.Pf("\n")
	}
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// find_a_zero: method to support step 4
func (o *Munkres) find_a_zero() (row, col int) {
	i, j, done := 0, 0, false
	row, col = -1, -1
	for !done {
		j = 0
		for true {
			if o.C[i][j] == 0 && !o.row_covered[i] && !o.col_covered[j] {
				row, col, done = i, j, true
			}
			j += 1
			if j >= o.ncol || done {
				break
			}
		}
		i += 1
		if i >= o.nrow {
			done = true
		}
	}
	return
}

// find_star_in_row: method to support step 4
func (o *Munkres) find_star_in_row(row int) (col int) {
	col = -1
	for j := 0; j < o.ncol; j++ {
		if o.M[row][j] == STAR {
			col = j
		}
	}
	return
}

// find_star_in_col: method to support step 5
func (o *Munkres) find_star_in_col(c int) (r int) {
	r = -1
	for i := 0; i < o.nrow; i++ {
		if o.M[i][c] == STAR {
			r = i
		}
	}
	return
}

// find_star_in_col: method to support step 5
func (o *Munkres) find_prime_in_row(r int) (c int) {
	for j := 0; j < o.ncol; j++ {
		if o.M[r][j] == PRIM {
			c = j
		}
	}
	return
}
