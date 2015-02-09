// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

// DblMatToArray converts a matrix into a column-major array
func DblMatToArray(a [][]float64) (v []float64) {
	m, n, k := len(a), len(a[0]), 0
	v = make([]float64, m*n)
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			v[k] = a[i][j]
			k += 1
		}
	}
	return
}

// DblArrayToMat converts a column-major array to a matrix
func DblArrayToMat(v []float64, m, n int) (a [][]float64) {
	a = make([][]float64, m)
	for i := 0; i < m; i++ {
		a[i] = make([]float64, n)
	}
	for k := 0; k < m*n; k++ {
		i, j := k%m, k/m
		a[i][j] = v[k]
	}
	return
}

// Deep3Print prints an array of array of array
func Deep3Print(name string, A [][][]float64) {
	Pf("%s = [\n", name)
	for _, a := range A {
		Pf("  %v\n", a)
	}
	Pf("]\n")
}

// Deep3Serialize serializes an array of array of array in column-compressed format
//  Example:
//
func Deep3Serialize(A [][][]float64) (I, P []int, S []float64) {
	i, p := 0, 0
	for _, a := range A {
		for _, b := range a {
			i += 1
			p += len(b)
		}
	}
	I = make([]int, i)
	P = make([]int, i+1)
	S = make([]float64, p)
	i, p, k := 0, 0, 0
	for j, a := range A {
		for _, b := range a {
			I[i] = j
			i += 1
			p += len(b)
			P[i] = p
			for _, v := range b {
				S[k] = v
				k += 1
			}
		}
	}
	return
}

// Deep3GetInfo returns information of serialized array of array of array
//  Example:
func Deep3GetInfo(I, P []int, S []float64, verbose bool) (nitems, nrows, ncols_tot int, ncols []int) {
	if verbose {
		Pf("I = %v\n", I)
		Pf("P = %v\n", P)
		Pf("S = %v\n", S)
	}
	nitems = P[len(P)-1]
	nrows = I[len(I)-1] + 1
	ncols = make([]int, nrows)
	for _, j := range I {
		ncols_tot += 1
		ncols[j] += 1
	}
	if verbose {
		Pf("nitems    = %v\n", nitems)
		Pf("nrows     = %v\n", nrows)
		Pf("ncols_tot = %v\n", ncols_tot)
		Pf("ncols     = %v\n", ncols)
	}
	return
}

// Deep3Deserialize deserializes an array of array of array in column-compressed format
//  Example:
func Deep3Deserialize(I, P []int, S []float64, debug bool) (A [][][]float64) {
	_, nrows, _, ncols := Deep3GetInfo(I, P, S, false)
	A = make([][][]float64, nrows)
	for i := 0; i < nrows; i++ {
		A[i] = make([][]float64, ncols[i])
	}
	iprev := 0 // previous i
	j := 0     // column index
	for l, i := range I {
		nitems := P[l+1] - P[l]
		if i != iprev { // jumped to new column
			j = 0
		}
		if debug {
			Pf("l=%v  i=%v  nitems=%v  j=%v\n", l, i, nitems, j)
		}
		for k, p := 0, P[l]; p < P[l+1]; k, p = k+1, p+1 {
			if debug {
				Pf("  k=%v  p=%v  s=%v\n", k, p, S[p])
			}
			if k == 0 {
				A[i][j] = make([]float64, nitems)
			}
			A[i][j][k] = S[p]
		}
		iprev = i
		j += 1
	}
	return
}
