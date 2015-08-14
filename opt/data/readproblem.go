// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"strings"

	"github.com/cpmech/gosl/io"
)

func main() {

	// filename
	fn, _ := io.ArgToFilename(0, "afiro", ".dat", true)

	// variables
	var m int        // number or rows (input)
	var n int        // number of columns (input)
	var Ap []int     // pointers to the begining of storage of column (size n+1)
	var Ai []int     // row indices for each non zero entry (input, nnz A)
	var Ax []float64 // non zero entries (input, nnz A)
	var b []float64  // right hand side (input, size m)
	var c []float64  // objective vector (minimize, size n)
	var z0 float64   // initial fixed value for objective
	var l []float64  // lower bounds on variables (size n)
	var u []float64  // upper bounds on variables (size n)

	// auxiliary
	reading_Ap := false
	reading_Ai := false
	reading_Ax := false
	reading_b := false
	reading_c := false
	reading_z0 := false
	reading_l := false
	reading_u := false
	atof := func(s string) float64 {
		return io.Atof(strings.Replace(s, "D", "E", 1))
	}

	// read data
	k := 0
	io.ReadLines(fn, func(idx int, line string) (stop bool) {
		if idx == 0 { // skip name
			return
		}
		str := strings.Fields(line)
		if idx == 1 { // read m and m
			m, n = io.Atoi(str[0]), io.Atoi(str[1])
			Ap = make([]int, n+1)
			k = 0
			reading_Ap = true
			return
		}
		for _, s := range str {
			if reading_Ap {
				if k == n+1 {
					reading_Ap = false
					reading_Ai = true
					nnz := Ap[n]
					Ai = make([]int, nnz)
					Ax = make([]float64, nnz)
					b = make([]float64, m)
					c = make([]float64, n)
					l = make([]float64, n)
					u = make([]float64, n)
					k = 0
				} else {
					Ap[k] = io.Atoi(s) - 1 // subtract 1 because of Fortran indexing
				}
			}
			if reading_Ai {
				if k == Ap[n] {
					reading_Ai = false
					reading_Ax = true
					k = 0
				} else {
					Ai[k] = io.Atoi(s) - 1 // subtract 1 because of Fortran indexing
				}
			}
			if reading_Ax {
				if k == Ap[n] {
					reading_Ax = false
					reading_b = true
					k = 0
				} else {
					Ax[k] = atof(s)
				}
			}
			if reading_b {
				if k == m {
					reading_b = false
					reading_c = true
					k = 0
				} else {
					b[k] = atof(s)
				}
			}
			if reading_c {
				if k == n {
					reading_c = false
					reading_z0 = true
					k = 0
				} else {
					c[k] = atof(s)
				}
			}
			if reading_z0 {
				z0 = atof(s)
				reading_z0 = false
				reading_l = true
				k = 0
				return
			}
			if reading_l {
				if k == n {
					reading_l = false
					reading_u = true
					k = 0
				} else {
					l[k] = atof(s)
				}
			}
			if reading_u {
				if k == n {
					reading_u = false
					k = 0
				} else {
					u[k] = atof(s)
				}
			}
			k++
		}
		return
	})

	io.Pforan("Ap = %v\n", Ap)
	io.Pfcyan("Ai = %v\n", Ai)
	io.Pfyel("Ax = %v\n", Ax)
	io.Pf("b = %v\n", b)
	io.Pforan("c = %v\n", c)
	io.Pfcyan("l = %v\n", l)
	io.Pfyel("u = %v\n", u)
}
