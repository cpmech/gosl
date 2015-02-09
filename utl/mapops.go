// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

// IntIntsMapAppend appends a new item to a map of slice.
//  Note: this function creates a new slice in the map if key is not found.
func IntIntsMapAppend(m *map[int][]int, key int, item int) {
    if slice, ok := (*m)[key]; ok {
        slice = append(slice, item)
        (*m)[key] = slice
    } else {
        (*m)[key] = []int{item}
    }
}

// StrDblsMapAppend appends a new item to a map of slice.
//  Note: this function creates a new slice in the map if key is not found.
func StrDblsMapAppend(m *map[string][]float64, key string, item float64) {
    if slice, ok := (*m)[key]; ok {
        slice = append(slice, item)
        (*m)[key] = slice
    } else {
        (*m)[key] = []float64{item}
    }
}
