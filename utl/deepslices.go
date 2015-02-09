// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

// Deep3alloc allocates a slice of slice of slice
func Deep3alloc(n1, n2, n3 int) (a [][][]float64) {
    a = make([][][]float64, n1)
    for i := 0; i < n1; i++ {
        a[i] = make([][]float64, n2)
        for j := 0; j < n2; j++ {
            a[i][j] = make([]float64, n3)
        }
    }
    return
}

// Deep4alloc allocates a slice of slice of slice of slice
func Deep4alloc(n1, n2, n3, n4 int) (a [][][][]float64) {
    a = make([][][][]float64, n1)
    for i := 0; i < n1; i++ {
        a[i] = make([][][]float64, n2)
        for j := 0; j < n2; j++ {
            a[i][j] = make([][]float64, n3)
            for k := 0; k < n3; k++ {
                a[i][j][k] = make([]float64, n4)
            }
        }
    }
    return
}

// Deep3set sets deep slice of slice of slice with v values
func Deep3set(a [][][]float64, v float64) {
    for i := 0; i < len(a); i++ {
        for j := 0; j < len(a[i]); j++ {
            for k := 0; k < len(a[i][j]); k++ {
                a[i][j][k] = v
            }
        }
    }
}

// Deep4set sets deep slice of slice of slice of slice with v values
func Deep4set(a [][][][]float64, v float64) {
    for i := 0; i < len(a); i++ {
        for j := 0; j < len(a[i]); j++ {
            for k := 0; k < len(a[i][j]); k++ {
                for l := 0; l < len(a[i][j][k]); l++ {
                    a[i][j][k][l] = v
                }
            }
        }
    }
}
