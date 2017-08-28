// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"sort"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// Equations organises the identification numbers (IDs) of equations in a linear system:
// [A] ⋅ {x} = {b}. Some of the components of the {x} vector may be known in advance—i.e. some
// {x} components are "prescribed"—therefore the system of equations must be reduced first
// before its solution. The equations corresponding to the known {x} values are called
// "known equations" whereas the equations corresponding to unknown {x} values are called
// "unknown equations".
//
// The right-hand-side vector {b} will also be partitioned into two sets. However, the
// "unknown equations" part will have known values of {b}. If needed, the other {b}
// values can be computed from the "known equations".
//
// The structure Equations will then hold the IDs of two sets of equations:
// "u" -- unknown {x} values with given {b} values
// "k" -- known {x} values (the corresponding {b} may be post-computed)
//
//  In summary:
//
//   * We define the u-reduced system of equations with the unknown {x} (and known {b})
//                   ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄                       ┄┄┄┄┄┄┄
//     i.e. the "effective" system that needs to be solved: [Auu] ⋅ {xu} = {bu}
//
//   * We define the k-reduced system of equations with the known {x} values.
//                   ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄                       ┄┄┄┄┄
//
//   * We call the system [A] ⋅ {x} = {b} the "full" system, with all equations.
//
//  The "full" linear system is:
//
//                     [A]⋅{x} = {b}   or   [ Auu Auk ]⋅{ xu ? } = { bu ✓ }
//                                          [ Aku Akk ] { xk ✓ }   { bk ? }
//
//  NOTE: the {bu} part is known and the {bk} can be (post-)computed if needed.
//
//  The following mnemonic may help:
//
//                                Auu Auk  u─────> unknown
//                                Aku Akk  k─────> known
//                                  u   k
//                                  │   └────────> known
//                                  └────────────> unknown
//
//  The Equations structure uses arrays to map the indices of equations from the "full"
//  system to the reduced systems and vice-versa. We call these arrays "maps" but they
//  are not Go maps; they are simple slices of integers and should give better
//  performance than maps.
//
//  Four "maps" (slices of integers) are built:
//   * UtoF -- reduced u-system to full-system
//   * FtoU -- full-system to reduced u-system
//   * KtoF -- reduced k-system to full-system
//   * FtoK -- full-system to reduced k-system
//
//
//  EXAMPLE:
//
//   N  = 9          ⇒  total number of equations
//   kx = [0  3  6]  ⇒  known x-components (i.e. "prescribed" equations)
//
//     0  1  2  3  4  5  6  7  8                  0   1  2   3   4  5   6   7  8
//   0 +--------+--------+------ 0  ◀ known ▶  0 [kk] ku ku [kk] ku ku [kk] ku ku 0
//   1 |        |        |       1             1  uk  ○  ○   uk  ○  ○   uk  ○  ○  1
//   2 |        |        |       2             2  uk  ○  ○   uk  ○  ○   uk  ○  ○  2
//   3 +--------+--------+------ 3  ◀ known ▶  3 [kk] ku ku [kk] ku ku [kk] ku ku 3
//   4 |        |        |       4             4  uk  ○  ○   uk  ○  ○   uk  ○  ○  4
//   5 |        |        |       5             5  uk  ○  ○   uk  ○  ○   uk  ○  ○  5
//   6 +--------+--------+------ 6  ◀ known ▶  6 [kk] ku ku [kk] ku ku [kk] ku ku 6
//   7 |        |        |       7             7  uk  ○  ○   uk  ○  ○   uk  ○  ○  7
//   8 |        |        |       8             8  uk  ○  ○   uk  ○  ○   uk  ○  ○  8
//     0  1  2  3  4  5  6  7  8                  0   1  2   3   4  5   6   7  8
//     ▲        ▲        ▲
//   known    known    known                             ○ means uu equations
//
//   Nu = 6 => number of "unknown equations"
//   Nk = 3 => number of "known equations"
//   Nu + Nk = N
//
//   reduced:0  1  2  3  4  5              Example:
//   UtoF = [1  2  4  5  7  8]           ⇒ UtoF[3] returns equation # 5 of full system
//
//           0  1  2  3  4  5  6  7  8
//   FtoU = [   0  1     2  3     4  5]  ⇒ FtoU[5] returns equation # 3 of reduced system
//          -1       -1       -1         ⇒ -1 indicates 'value not set'
//
//   reduced:0  1  2
//   KtoF = [0  3  6]                    ⇒ KtoF[1] returns equation # 3 of full system
//
//           0  1  2  3  4  5  6  7  8
//   FtoK = [0        1        2      ]  ⇒ FtoK[3] returns equation # 1 of reduced system
//             -1 -1    -1 -1    -1 -1   ⇒ -1 indicates 'value not set'
//
type Equations struct {
	N    int   // total number of equations
	Nu   int   // number of unknowns (size of u-system)
	Nk   int   // number of known values (size of k-system)
	UtoF []int // reduced u-system to full-system
	FtoU []int // full-system to reduced u-system
	KtoF []int // reduced k-system to full system
	FtoK []int // full-system to reduced k-system
}

// Init initialises Equations
//   n  -- total number of equations; i.e. len({x}) in [A]⋅{x}={b}
//   kx -- known x-components ⇒ "known equations" [may be unsorted]
func (o *Equations) Init(n int, kx []int) {
	o.N = n
	o.Nk = len(kx)
	o.Nu = o.N - o.Nk
	o.UtoF = make([]int, o.Nu)
	o.FtoU = make([]int, o.N)
	o.KtoF = make([]int, o.Nk)
	o.FtoK = make([]int, o.N)
	sortedkx := utl.IntGetSorted(kx)
	var iu, ik int
	for eq := 0; eq < n; eq++ {
		o.FtoU[eq] = -1 // -1 indicates "not set"
		o.FtoK[eq] = -1
		idx := sort.SearchInts(sortedkx, eq)
		if idx < len(sortedkx) && sortedkx[idx] == eq { // found known x-component
			o.KtoF[ik] = eq
			o.FtoK[eq] = ik
			ik++
		} else { // not found ⇒ unknown x-component
			o.UtoF[iu] = eq
			o.FtoU[eq] = iu
			iu++
		}
	}
}

// Stat prints information about Equations
func (o *Equations) Stat(full bool) {
	io.Pf("number of unknown x-components: Nu = %d\n", o.Nu)
	io.Pf("number of known x-components:   Nk = %d\n", o.Nk)
	io.Pf("total number of equations:      N  = %d\n", o.N)
	if full {
		io.Pf("reduced u-system to full-system map:\nUtoF =\n")
		io.Pf("%v\n", o.UtoF)
		io.Pf("full-system to reduced u-system map:\nFtoU =\n")
		io.Pf("%v\n", o.FtoU)
		io.Pf("reduced k-system to full system map:\nKtoF =\n")
		io.Pf("%v\n", o.KtoF)
		io.Pf("full-system to reduced k-system map:\nFtoK =\n")
		io.Pf("%v\n", o.FtoK)
	}
}
