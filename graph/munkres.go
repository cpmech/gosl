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
//   Example:
//           $ | Clean  Sweep   Wash
//      -------|--------------------
//      Fry    |   [2]      3      3
//      Leela  |     3    [2]      3
//      Bender |     3      3    [2]
//      minimum cost = 6
//  Note: costs must be positive
type Munkres struct {

	// main
	C     [][]float64 // [nrow][ncol] cost matrix
	Cori  [][]float64 // [nrow][ncol] original cost matrix
	Links []int       // [nrow] will contain links/assignments after Run(), where j := o.Links[i] means that i is assigned to j. -1 means no assignment/link
	Cost  float64     // total cost after Run() and links are established

	// auxiliary
	M           [][]Mask_t // [nrow][ncol] mask matrix. If Mij==1, then Cij is a starred zero. If Mij==2, then Cij is a primed zero
	path        [][]int    // path
	row_covered []bool     // indicates whether a row is covered or not
	col_covered []bool     // indicates whether a column is covered or not
	nrow        int        // number of rows in cost/mask matrix
	ncol        int        // number of column in cost/mask matrix
	path_row_0  int        // first row in path
	path_col_0  int        // first col in path
}

// Init initialises Munkres' structure
func (o *Munkres) Init(nrow, ncol int) {
	chk.IntAssertLessThan(0, nrow) // nrow > 1
	chk.IntAssertLessThan(0, ncol) // ncol > 1
	o.nrow, o.ncol = nrow, ncol
	o.C = utl.DblsAlloc(o.nrow, o.ncol)
	o.M = make([][]Mask_t, o.nrow)
	for i := 0; i < o.nrow; i++ {
		o.M[i] = make([]Mask_t, o.ncol)
	}
	o.Links = make([]int, o.nrow)
	npath := 2*o.nrow + 1 // TODO: check this
	o.path = utl.IntsAlloc(npath, 2)
	o.row_covered = make([]bool, o.nrow)
	o.col_covered = make([]bool, o.ncol)
}

// SetCostMatrix sets cost matrix by copying from C to internal o.C
//  Note: costs must be positive
func (o *Munkres) SetCostMatrix(C [][]float64) {
	o.Cori = C
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			o.C[i][j] = C[i][j]
			o.M[i][j] = NONE
			if math.IsNaN(o.C[i][j]) {
				chk.Panic("cannot set cost matrix because of NaN value")
			}
		}
		o.row_covered[i] = false
	}
	for j := 0; j < o.ncol; j++ {
		o.col_covered[j] = false
	}
	o.path_row_0 = 0
	o.path_col_0 = 0
}

// Run runs the iterative algorithm
//  Output:
//   o.Links -- will contain assignments, where len(assignments) == nrow and
//              j := o.Links[i] means that i is assigned to j
//              -1 means no assignment/link
//   o.Cost -- will have the total cost by following links
func (o *Munkres) Run() {

	// column matrix
	if o.ncol == 1 {
		o.Cost = o.C[0][0]
		isel := 0
		o.Links[0] = -1
		for i := 1; i < o.nrow; i++ {
			o.Links[i] = -1
			if o.C[i][0] < o.Cost {
				o.Cost = o.C[i][0]
				isel = i
			}
		}
		o.Links[isel] = 0
		return
	}

	// row matrix
	if o.nrow == 1 {
		o.Cost = o.C[0][0]
		jsel := 0
		for j := 1; j < o.ncol; j++ {
			if o.C[0][j] < o.Cost {
				o.Cost = o.C[0][j]
				jsel = j
			}
		}
		o.Links[0] = jsel
		return
	}

	// run Munkres algorithm
	step := 1
	done := false
	for !done {
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
	}

	// compute cost and set links
	o.Cost = 0
	for i := 0; i < o.nrow; i++ {
		o.Links[i] = -1
		for j := 0; j < o.ncol; j++ {
			if o.M[i][j] == STAR {
				o.Links[i] = j
				o.Cost += o.Cori[i][j]
				break
			}
		}
	}
}

// steps //////////////////////////////////////////////////////////////////////////////////////////

// step1: for each row of the cost matrix, find the smallest element and subtract it from every
// element in its row. next_step = 2
func (o *Munkres) step1() (next_step int) {
	var xmin float64
	for i := 0; i < o.nrow; i++ {
		xmin = o.C[i][0]
		for j := 1; j < o.ncol; j++ {
			xmin = utl.Min(xmin, o.C[i][j])
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
	count := 0
	for j := 0; j < o.ncol; j++ {
		if o.col_covered[j] {
			count += 1
		}
	}
	next_step = 4 // not all columns are covered
	if count >= o.ncol || count >= o.nrow {
		next_step = 7 // all covered
	}
	return
}

// step4: find a noncovered zero and prime it. Change cost matrix by adding min value to every
// element of covered rows and subtracting it from elements of uncovered columns until a noncovered
// zero is found.
func (o *Munkres) step4() (next_step int) {
	row, col, done := -1, -1, false
	for !done {
		row, col = o.find_noncov_zero()
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
// every line in the matrix. next_step = 3
func (o *Munkres) step5() (next_step int) {

	// construct series
	r, c, done := -1, -1, false
	count := 1
	o.path[count-1][0] = o.path_row_0
	o.path[count-1][1] = o.path_col_0
	for !done {
		r = o.find_star_in_col(o.path[count-1][1])
		if r > -1 {
			count += 1
			o.path[count-1][0] = r
			o.path[count-1][1] = o.path[count-2][1]
		} else {
			done = true
		}
		if !done {
			c = o.find_prime_in_row(o.path[count-1][0])
			count += 1
			o.path[count-1][0] = o.path[count-2][0]
			o.path[count-1][1] = c
		}
	}

	// augment path
	for p := 0; p < count; p++ {
		if o.M[o.path[p][0]][o.path[p][1]] == STAR {
			o.M[o.path[p][0]][o.path[p][1]] = NONE
		} else {
			o.M[o.path[p][0]][o.path[p][1]] = STAR
		}
	}

	// clear covers
	for i := 0; i < o.nrow; i++ {
		o.row_covered[i] = false
	}
	for j := 0; j < o.ncol; j++ {
		o.col_covered[j] = false
	}

	// erase primes
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			if o.M[i][j] == PRIM {
				o.M[i][j] = NONE
			}
		}
	}
	return 3
}

// step6: add min value to every element of each covered row, and subtract it from every element of
// each uncovered column. next_step = 4
func (o *Munkres) step6() (next_step int) {

	// find min value
	xmin := math.MaxFloat64
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			if !o.row_covered[i] && !o.col_covered[j] {
				xmin = utl.Min(xmin, o.C[i][j])
			}
		}
	}

	// add/subtract min value
	for i := 0; i < o.nrow; i++ {
		for j := 0; j < o.ncol; j++ {
			if o.row_covered[i] {
				o.C[i][j] += xmin
			}
			if !o.col_covered[j] {
				o.C[i][j] -= xmin
			}
		}
	}
	return 4
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// StrCostMatrix returns a representation of cost matrix with masks and covers
func (o *Munkres) StrCostMatrix() (l string) {
	numfmt := "%v"
	l += io.Sf("%4v", " ")
	for j := 0; j < o.ncol; j++ {
		if o.col_covered[j] {
			l += io.Sf("%8v", "T ")
		} else {
			l += io.Sf("%8v", "F ")
		}
	}
	l += io.Sf("\n")
	for i := 0; i < o.nrow; i++ {
		if o.row_covered[i] {
			l += io.Sf("%4v", "T")
		} else {
			l += io.Sf("%4v", "F")
		}
		for j := 0; j < o.ncol; j++ {
			s := io.Sf(numfmt, o.C[i][j])
			switch o.M[i][j] {
			case NONE:
				s += " "
			case STAR:
				s += "*"
			case PRIM:
				s += "'"
			}
			l += io.Sf("%8v", s)
		}
		l += io.Sf("\n")
	}
	return
}

// find_noncov_zero finds the row and column of a non-covered zero entry. -1 means not found
func (o *Munkres) find_noncov_zero() (row, col int) {
	i, j, done := 0, 0, false
	row, col = -1, -1
	for !done {
		j = 0
		for true {
			if !o.row_covered[i] && !o.col_covered[j] {
				if o.C[i][j] == 0 {
					row, col, done = i, j, true
				}
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

/* Notes:
 1. The algorthm will work even when the minimum values in two or more rows are in the same column.
 2. The algorithm will work even when two or more of the rows contain the same values in the the same order.
 3. The algorithm will work even when all the values are the same (although the result is not very interesting).
 4. Munkres Assignment Algorithm is not exponential run time or intractable; it is of a low order polynomial run time, worst-case O(n3).
 5. Optimality is guaranteed in Munkres Assignment Algorithm.
 6. Setting the cost matrix to C(i,j) = i*j  makes a good testing matrix for this problem.
 7. In this algorithm the range of indices is[0..n-1] rather than [1..n].
 8. Step 3 is an example of the greedy method.  If the minimum values are all in different rows then their positions represent the minimal pairwise assignments.
 9. Step 5 is an example of the Augmenting Path Algorithm (Stable Marriage Problem).
10. Step 6 is an example of constraint relaxation.  It is "giving up" on a particular cost and raising the constraint by the least amount possible.
11. If your implementation is jumping between Step 4 and Step 6 without entering Step 5, you probably have not properly dealt with recognizing that there are no uncovered zeros in Step 4.
12. In the matrix M 1=starred zero and 2=primed zero.  So, if C[i,j] is a starred zero we would set M[i,j]=1.  All other elements in M are set to zero
13. The Munkres assignment algorithm can be implemented as a sparse matrix, but you will need to ensure that the correct (optimal) assignment pairs are active in the sparse cost matrix C
14. Munkres Assignment can be applied to TSP, pattern matching, track initiation, data correlation, and (of course) any pairwise assignment application.
15. Munkres can be extended to rectangular arrays (i.e. more jobs than workers, or more workers than jobs) .
16. The best way to find a maximal assignment is to replace the values ci,j in the cost matrix with C[i,j] = bigval - ci,j.
17. Original Reference: Algorithms for Assignment and Transportation Problems, James Munkres, Journal of the Society for Industrial and Applied Mathematics Volume 5, Number 1, March, 1957
18. Extension to Rectangular Arrays Ref:  F. Burgeois and J.-C. Lasalle. An extension of the Munkres algorithm for the assignment problem to rectangular matrices. Communications of the ACM, 142302-806, 1971.
*/
