// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

// NaiveVecDot is the naive version of la.VecDot
func NaiveVecDot(u, v []float64) (res float64) {
	for i := 0; i < len(u); i++ {
		res += u[i] * v[i]
	}
	return
}

// NaiveVecAdd is the naive version of la.VecAdd
func NaiveVecAdd(res []float64, α float64, u []float64, β float64, v []float64) {
	for i := 0; i < len(u); i++ {
		res[i] = α*u[i] + β*v[i]
	}
}

// NaiveMatVecMul is the naive version of la.MatVecMul
func NaiveMatVecMul(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a); i++ {
		v[i] = 0.0
		for j := 0; j < len(a[i]); j++ {
			v[i] += α * a[i][j] * u[j]
		}
	}
}

// NaiveMatMatMul is the naive version of la.MatMatMul
func NaiveMatMatMul(c [][]float64, α float64, a, b [][]float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b[0]); j++ {
			c[i][j] = 0.0
			for k := 0; k < len(a[0]); k++ {
				c[i][j] += α * a[i][k] * b[k][j]
			}
		}
	}
}
