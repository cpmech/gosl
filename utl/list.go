// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import "github.com/cpmech/gosl/io"

// DblList implements a tabular list with variable number of columns
//  Example:
//    Vals = [][]float64{
//             {0.0},
//             {1.0, 1.1, 1.2, 1.3},
//             {2.0, 2.1},
//             {3.0, 3.1, 3.2},
//           }
type DblList struct {
	Vals [][]float64 // values
}

// Append appends items to DblList
func (o *DblList) Append(rowidx int, value float64) {
	size := len(o.Vals)

	// fill new rows if necessary
	if rowidx >= size {
		for i := size; i <= rowidx; i++ {
			o.Vals = append(o.Vals, []float64{})
		}
	}

	// append value
	o.Vals[rowidx] = append(o.Vals[rowidx], value)
}

// DblSlist implements a tabular list with variable number of columns
// using a serial representation
//  Example:
//      0.0
//      1.0  1.1  1.2  1.3
//      2.0  2.1
//      3.0  3.1  3.2
//  becomes:
//             (0)   (1)    2    3    4   (5)    6   (7)    8    9   (10)
//      Vals = 0.0 | 1.0  1.1  1.2  1.3 | 2.0  2.1 | 3.0  3.1  3.2 |
//      Ptrs = 0 1 5 7 10
//  Notes:
//      len(Ptrs) = nrows + 1
//      Ptrs[len(Ptrs)-1] = len(Vals)
type DblSlist struct {
	Vals []float64 // values
	Ptrs []int     // pointers
}

// Append appends item to DblSlist
func (o *DblSlist) Append(startRow bool, value float64) {
	if startRow {
		if len(o.Ptrs) == 0 {
			o.Ptrs = []int{0, 0}
		} else {
			start := o.Ptrs[len(o.Ptrs)-1]
			o.Ptrs = append(o.Ptrs, start)
		}
	}
	o.Vals = append(o.Vals, value)
	o.Ptrs[len(o.Ptrs)-1] += 1
}

// Print prints the souble-serial-list
func (o DblSlist) Print(fmt string) {
	for i := 0; i < len(o.Ptrs)-1; i++ {
		for j := o.Ptrs[i]; j < o.Ptrs[i+1]; j++ {
			io.Pf(fmt, o.Vals[j])
		}
		io.Pf("\n")
	}
}
