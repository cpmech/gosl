// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"strings"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// ReadLPfortran reads linear program from particular fortran file
//  download LP files from here: http://users.clas.ufl.edu/hager/coap/format.html
//  Output:
//   A -- compressed-column sparse matrix where:
//        Ap -- pointers to the beginning of storage of column (size n+1)
//        Ai -- row indices for each non zero entry (input, nnz A)
//        Ax -- non zero entries (input, nnz A)
//   b -- right hand side (input, size m)
//   c -- objective vector (minimize, size n)
//   l -- lower bounds on variables (size n)
//   u -- upper bounds on variables (size n)
func ReadLPfortran(fn string) (A *la.CCMatrix, b, c, l, u []float64) {

	// variables
	var m int        // number or rows (input)
	var n int        // number of columns (input)
	var Ap []int     // pointers to the beginning of storage of column (size n+1)
	var Ai []int     // row indices for each non zero entry (input, nnz A)
	var Ax []float64 // non zero entries (input, nnz A)
	var z0 float64   // initial fixed value for objective

	// auxiliary
	readingAp := false
	readingAi := false
	readingAx := false
	readingb := false
	readingc := false
	readingz0 := false
	readingl := false
	readingu := false
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
			readingAp = true
			return
		}
		for _, s := range str {
			if readingAp {
				if k == n+1 {
					readingAp = false
					readingAi = true
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
			if readingAi {
				if k == Ap[n] {
					readingAi = false
					readingAx = true
					k = 0
				} else {
					Ai[k] = io.Atoi(s) - 1 // subtract 1 because of Fortran indexing
				}
			}
			if readingAx {
				if k == Ap[n] {
					readingAx = false
					readingb = true
					k = 0
				} else {
					Ax[k] = atof(s)
				}
			}
			if readingb {
				if k == m {
					readingb = false
					readingc = true
					k = 0
				} else {
					b[k] = atof(s)
				}
			}
			if readingc {
				if k == n {
					readingc = false
					readingz0 = true
					k = 0
				} else {
					c[k] = atof(s)
				}
			}
			if readingz0 {
				z0 = atof(s)
				_ = z0
				readingz0 = false
				readingl = true
				k = 0
				return
			}
			if readingl {
				if k == n {
					readingl = false
					readingu = true
					k = 0
				} else {
					l[k] = atof(s)
				}
			}
			if readingu {
				if k == n {
					readingu = false
					k = 0
				} else {
					u[k] = atof(s)
				}
			}
			k++
		}
		return
	})

	// debug
	if false {
		io.Pforan("Ap = %v\n", Ap)
		io.Pfcyan("Ai = %v\n", Ai)
		io.Pfyel("Ax = %v\n", Ax)
		io.Pf("b = %v\n", b)
		io.Pforan("c = %v\n", c)
		io.Pfcyan("l = %v\n", l)
		io.Pfyel("u = %v\n", u)
	}

	// results
	A = new(la.CCMatrix)
	A.Set(m, n, Ap, Ai, Ax)
	return
}
