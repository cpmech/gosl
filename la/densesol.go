// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Cholesky returns the Cholesky decomposition of a symmetric positive-definite matrix:
//  a = L * trans(L)
func Cholesky(L, a [][]float64) (err error) {
	for j := 0; j < len(a); j++ { // loop over columns
		for i := j; i < len(a); i++ { // loop over lower diagonal rows (including diagonal)
			amsum := a[i][j]
			for k := 0; k < j; k++ {
				amsum -= L[i][k] * L[j][k]
			}
			if i == j {
				if amsum <= 0.0 {
					err = chk.Err(_densesol_err1)
				}
				L[i][j] = math.Sqrt(amsum)
			} else {
				L[i][j] = amsum / L[j][j]
			}
		}
	}
	return
}

// SPDsolve (Symmetric/Positive-Definite) solves a small dense linear system A*x=b where
// the "a" matrix is symmetric and positive-definite (and real, and, of course, square)
//  x := inv(a) * b
//  NOTE: this function uses Cholesky decomposition and should be used for small systems
func SPDsolve(x []float64, a [][]float64, b []float64) (err error) {
	// Cholesky factorisation
	n := len(a)
	L := MatAlloc(n, n)
	cerr := Cholesky(L, a)
	if cerr != nil {
		err = chk.Err(_densesol_err2, cerr.Error())
		return
	}
	// solve L*y = b storing y in x
	for i := 0; i < n; i++ {
		bmsum := b[i]
		for k := 0; k < i; k++ {
			bmsum -= L[i][k] * x[k]
		}
		x[i] = bmsum / L[i][i]
	}
	// solve trans(L)*x = y with y==x
	for i := n - 1; i >= 0; i-- {
		bmsum := x[i]
		for k := i + 1; k < n; k++ {
			bmsum -= L[k][i] * x[k]
		}
		x[i] = bmsum / L[i][i]
	}
	return
}

// error messages
var (
	_densesol_err1 = "densesol.go: Cholesky factorization failed due to non positive-definite matrix"
	_densesol_err2 = "densesol.go: SymPDsolve failed: %s"
)
