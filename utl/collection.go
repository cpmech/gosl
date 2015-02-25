// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

// Collection implements a matrix with variable number of entities
//  Example:
//    Vals = [][]float64{
//             {0.0},
//             {1.0, 1.1, 1.2, 1.3},
//             {2.0, 2.1},
//             {3.0, 3.1, 3.2},
//           }
type Collection struct {
	Vals [][]float64 // values
}

// Append appends items to Collection
func (o *Collection) Append(rowidx int, value float64) {
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
