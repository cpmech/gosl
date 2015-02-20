// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/io"
)

// check_residR (real) calculates max |component| of resid: r = a*x - b
func check_residR(tst *testing.T, tol float64, a [][]float64, x, b []float64) {
	n := len(x)
	r := make([]float64, n)
	for i := 0; i < n; i++ {
		r[i] = -b[i]
	}
	MatVecMulAdd(r, 1, a, x) // r += 1*a*x
	var abs_r, max_abs_cp float64
	for i := 0; i < n; i++ {
		abs_r = math.Abs(r[i])
		if abs_r > max_abs_cp {
			max_abs_cp = abs_r
		}
	}
	if max_abs_cp > tol {
		tst.Errorf("[1;31mresidual is too large: max|r[i]| = %g[0m", max_abs_cp)
	} else {
		io.Pf("max|r[i]| = %g [1;32mOK[0m\n", max_abs_cp)
	}
	return
}

// check_residC (complex) calculates max |component| of resid: r = a*x - b
func check_residC(tst *testing.T, tol float64, a [][]complex128, x, b []complex128) {
	n := len(x)
	r := make([]complex128, n)
	for i := 0; i < n; i++ {
		r[i] = -b[i]
	}
	MatVecMulAddC(r, 1, a, x) // r += 1*a*x
	var abs_r, max_abs_cp float64
	for i := 0; i < n; i++ {
		abs_r = max(math.Abs(real(r[i])), math.Abs(imag(r[i])))
		if abs_r > max_abs_cp {
			max_abs_cp = abs_r
		}
	}
	if max_abs_cp > tol {
		tst.Errorf("[1;31mresidual is too large: max|r[i]| = %g[0m", max_abs_cp)
	} else {
		io.Pf("max|r[i]| = %g [1;32mOK[0m\n", max_abs_cp)
	}
	return
}
