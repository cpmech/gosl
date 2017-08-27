// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"sort"

	"github.com/cpmech/gosl/io"
)

// Equations organises the identification numbers (IDs) of equations in a linear system: K⋅u=f
// Some of the components of the u vector may be known in advance——i.e. some u values are
// prescribed——therefore the system of equations needs to be reduced by removing the
// "known equations". The resulting equations contain the variables that needed to be solved for
// and are named here "unknown equations". The structure Equations will then hold the IDs
// of two sets of equations:
// 1 -- "unknown" indicated by the number "1"; and
// 2 -- "prescribed" indicated by the number "2"
//
//  In summary:
//
//   * We call reduced system # 1 a system of equations with the unknown equations only;
//             ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄                                ┄┄┄┄┄┄┄
//     i.e. the "effective" system that needs to be solved: K11⋅u1 = f1
//
//   * We call reduced system # 2 a system of equations with the prescribed equations only
//             ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄                                ┄┄┄┄┄┄┄┄┄┄
//
//   * We call the system K⋅u=f the "full" system, with all equations
//     (unknown and prescribed)
//
//  The "full" linear system is:
//                                  K ⋅ u = f     or    [ K11 K12 ] ⋅ { u1 } = { f1 }
//                                                      [ K21 K22 ]   { u2 }   { f2 }
//  We then split the K matrix into four parts:
//
//                                K11 K12  1─────> unknown
//                                K21 K22  2─────> prescribed
//                                  1   2
//                                  │   └────────> prescribed
//                                  └────────────> unknown
//
//  The Equations structure uses arrays to map the indices of equations from the "full"
//  system to the reduced systems and vice-versa. We call these arrays "maps" but they
//  are not Go maps; they are simple slices of integers and should give better
//  performance than using maps.
//
//  Two "maps" (slices of integers) are built:
//   * "RF" means "reduced to full"
//   * "FR" means "full to reduced"
//
//
//  EXAMPLE:
//
//   N   = 9          ⇒  total number of equations
//   peq = [0  3  6]  ⇒  prescribed equations
//
//     0  1  2  3  4  5  6  7  8                       0   1  2   3   4  5   6   7  8
//   0 +--------+--------+------ 0 -- prescribed -- 0 [22] 21 21 [22] 21 21 [22] 21 21 0
//   1 |        |        |       1                  1  12  .. ..  12  .. ..  12  .. .. 1
//   2 |        |        |       2                  2  12  .. ..  12  .. ..  12  .. .. 2
//   3 +--------+--------+------ 3 -- prescribed -- 3 [22] 21 21 [22] 21 21 [22] 21 21 3
//   4 |        |        |       4                  4  12  .. ..  12  .. ..  12  .. .. 4
//   5 |        |        |       5                  5  12  .. ..  12  .. ..  12  .. .. 5
//   6 +--------+--------+------ 6 -- prescribed -- 6 [22] 21 21 [22] 21 21 [22] 21 21 6
//   7 |        |        |       7                  7  12  .. ..  12  .. ..  12  .. .. 7
//   8 |        |        |       8                  8  12  .. ..  12  .. ..  12  .. .. 8
//     0  1  2  3  4  5  6  7  8                       0   1  2   3   4  5   6   7  8
//     |        |        |
//    pre      pre      pre                                 .. ⇒ 11 equations
//
//   N1 = 6 => number of equations of type 1 (unknowns)
//   N2 = 3 => number of equations of type 2 (prescribed)
//   N1 + N2 == N
//
//          0  1  2  3  4  5              Example:
//   RF1 = [1  2  4  5  7  8]           ⇒ RF1[3] returns equation # 5 of full system
//
//          0  1  2  3  4  5  6  7  8
//   FR1 = [   0  1     2  3     4  5]  ⇒ FR1[5] returns equation # 3 of reduced system
//         -1       -1       -1         ⇒ -1 indicates 'value not set'
//
//          0  1  2
//   RF2 = [0  3  6]                    ⇒ RF2[1] returns equation # 3 of full system
//
//          0  1  2  3  4  5  6  7  8
//   FR2 = [0        1        2      ]  ⇒ FR2[3] returns equation # 1 of reduced system
//            -1 -1    -1 -1    -1 -1   ⇒ -1 indicates 'value not set'
//
type Equations struct {
	N1, N2, N int   // unknowns, prescribed, total numbers
	RF1, FR1  []int // reduced=>full/full=>reduced maps for unknowns
	RF2, FR2  []int // reduced=>full/full=>reduced maps for prescribed
}

// Init initialises Equations
func (e *Equations) Init(n int, peqUsorted []int) {
	peq := make([]int, len(peqUsorted))
	copy(peq, peqUsorted)
	sort.Ints(peq)
	e.N = n
	e.N2 = len(peq)
	e.N1 = e.N - e.N2
	e.RF1 = make([]int, e.N1)
	e.FR1 = make([]int, e.N)
	e.RF2 = make([]int, e.N2)
	e.FR2 = make([]int, e.N)
	var i1, i2 int
	for eq := 0; eq < n; eq++ {
		e.FR1[eq] = -1 // indicates 'not set'
		e.FR2[eq] = -1
		idx := sort.SearchInts(peq, eq)
		if idx < len(peq) && peq[idx] == eq { // found => prescribed
			e.RF2[i2] = eq
			e.FR2[eq] = i2
			i2++
		} else { // not found => unknowns
			e.RF1[i1] = eq
			e.FR1[eq] = i1
			i1++
		}
	}
}

// Print prints information about Equations
func (e *Equations) Print() {
	io.Pf("N1 = %v, N2 = %v, N = %v\n", e.N1, e.N2, e.N)
	io.Pf("RF1 (unknown) =\n %v\n", e.RF1)
	io.Pf("FR1 = \n%v\n", e.FR1)
	io.Pf("RF2 (prescribed) =\n %v\n", e.RF2)
	io.Pf("FR2 = \n%v\n", e.FR2)
}
