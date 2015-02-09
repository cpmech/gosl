// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package la implements routines and structures for linear algebra with
// matrices and vectors in dense and sparse formats (including complex)
package la

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"code.google.com/p/gosl/utl"
)

// PrintVec prints a vector
func PrintVec(name string, a []float64, format string, numpy bool) {
	if utl.Tsilent {
		return
	}
	r := name + " = "
	if numpy {
		r += " array(["
	}
	for i := 0; i < len(a); i++ {
		r += fmt.Sprintf(format, a[i])
		if numpy {
			if i < len(a)-1 {
				r += ","
			}
		}
	}
	if numpy {
		r += "])"
	}
	fmt.Println(r)
}

// PrintMat prints a dense matrix
func PrintMat(name string, a [][]float64, format string, numpy bool) {
	if utl.Tsilent {
		return
	}
	r := name + " ="
	if numpy {
		r += " array(["
	} else {
		r += "\n"
	}
	for i := 0; i < len(a); i++ {
		if numpy {
			r += "["
		}
		for j := 0; j < len(a[0]); j++ {
			r += fmt.Sprintf(format, a[i][j])
			if numpy {
				if j != len(a[0])-1 {
					r += ","
				}
			}
		}
		if numpy {
			if i == len(a)-1 {
				r += "]"
			} else {
				r += "],"
			}
		}
		if i != len(a)-1 {
			r += "\n"
		}
	}
	if numpy {
		r += "])"
	}
	fmt.Println(r)
}

// WriteSmat writes a smat matrix for vismatrix
func WriteSmat(fnkey string, a [][]float64, tol float64) {
	var bfa, bfb bytes.Buffer
	var nnz int = 0
	m := len(a)
	n := len(a[0])
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if math.Abs(a[i][j]) > tol {
				fmt.Fprintf(&bfb, "  %d  %d  %g\n", i, j, a[i][j])
				nnz++
			}
		}
	}
	fmt.Fprintf(&bfa, "%d  %d  %d\n", m, n, nnz)
	utl.WriteFile(fnkey+".smat", &bfa, &bfb)
}

// ReadSmat reads a smat matrix back
func ReadSmat(fn string) *Triplet {
	var t Triplet
	utl.ReadLines(fn,
		func(idx int, line string) (stop bool) {
			r := strings.Fields(line)
			if idx == 0 {
				m, n, nnz := utl.Atoi(r[0]), utl.Atoi(r[1]), utl.Atoi(r[2])
				t.Init(m, n, nnz)
			} else {
				t.Put(utl.Atoi(r[0]), utl.Atoi(r[1]), utl.Atof(r[2]))
			}
			return
		})
	return &t
}

// PrintVecC prints a vector of complex numbers
func PrintVecC(name string, a []complex128, format, formatz string, numpy bool) {
	if utl.Tsilent {
		return
	}
	r := name + " ="
	if numpy {
		r += " array(["
	}
	for i := 0; i < len(a); i++ {
		r += fmt.Sprintf(format, real(a[i]))
		r += fmt.Sprintf(formatz, imag(a[i]))
		if numpy {
			if i < len(a)-1 {
				r += ","
			}
		}
	}
	if numpy {
		r += "])"
	}
	fmt.Println(r)
}

// PrintMatC prints a matrix of complex numbers
func PrintMatC(name string, a [][]complex128, format, formatz string, numpy bool) {
	if utl.Tsilent {
		return
	}
	r := name + " ="
	if numpy {
		r += " array(["
	} else {
		r += "\n"
	}
	for i := 0; i < len(a); i++ {
		if numpy {
			r += "["
		}
		for j := 0; j < len(a[0]); j++ {
			r += fmt.Sprintf(format, real(a[i][j]))
			r += fmt.Sprintf(formatz, imag(a[i][j]))
			if numpy {
				if j != len(a[0])-1 {
					r += ","
				}
			}
		}
		if numpy {
			if i == len(a)-1 {
				r += "]"
			} else {
				r += "],"
			}
		}
		r += "\n"
	}
	if numpy {
		r += "])"
	}
	fmt.Println(r)
}

// SmatTriplet writes a ".smat" file that can be visualised with vismatrix
func SmatTriplet(fnkey string, t *Triplet) {
	var bfa, bfb bytes.Buffer
	var nnz int
	for k := 0; k < t.pos; k++ {
		if math.Abs(t.x[k]) > 1e-16 {
			fmt.Fprintf(&bfb, "  %d  %d  %23.15e\n", t.i[k], t.j[k], t.x[k])
			nnz++
		}
	}
	fmt.Fprintf(&bfa, "%d  %d  %d\n", t.m, t.n, nnz)
	utl.WriteFile(fnkey+".smat", &bfa, &bfb)
}

// SmatCCMatrix writes a ".smat" file that can be visualised with vismatrix
func SmatCCMatrix(fnkey string, a *CCMatrix) {
	var bfa, bfb bytes.Buffer
	var nnz int
	for j := 0; j < a.n; j++ {
		for p := a.p[j]; p < a.p[j+1]; p++ {
			if math.Abs(a.x[p]) > 1e-16 {
				fmt.Fprintf(&bfb, "  %d  %d  %23.15e\n", a.i[p], j, a.x[p])
				nnz++
			}
		}
	}
	fmt.Fprintf(&bfa, "%d  %d  %d\n", a.m, a.n, nnz)
	utl.WriteFile(fnkey+".smat", &bfa, &bfb)
}
