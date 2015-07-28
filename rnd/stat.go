// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Moments computes the 4th moments of a data set
//  Input:
//   x -- sample
//  Output:
//   sum  -- sum of values
//   mean -- mean average (first moment)
//   adev -- average deviation
//   sdev -- standrad deviation
//   vari -- variance (second moment)
//   skew -- skewness (third moment)
//   kurt -- kurtosis (fourth moment)
//  Based on:
//     Press WH, Teukolsky SA, Vetterling WT and Flannery BP (2007)
//       Numerical Recipes in C++ 2007 (3rd Edition), page 725.
func Moments(x []float64) (sum, mean, adev, sdev, vari, skew, kurt float64, err error) {

	// check
	n := len(x)
	if n < 2 {
		err = chk.Err("x set must have at least 2 items")
		return
	}

	// first pass to get the mean
	for i := 0; i < n; i++ {
		sum += x[i]
	}
	N := float64(n)
	mean = sum / N

	// second pass to get the first (absolute), second, third, and fourth moments of the deviation from the mean.
	var d, c, p float64
	for i := 0; i < n; i++ {
		d = x[i] - mean     // d ← xi - bar(x)
		adev += math.Abs(d) // adev ← Σ |d|
		c += d              // c ← Σ d  // corrector
		p = d * d           // p ← d²
		vari += p           // vari ← Σ d²
		p *= d              // p  ← d³
		skew += p           // skew ← Σ d³
		kurt += p * d       // kurt ← Σ d⁴
	}

	// put the pieces together according to the conventional definitions
	adev /= N
	vari = (vari - c*c/N) / (N - 1.0)
	sdev = math.Sqrt(vari)
	TOL := 1e-15
	if math.Abs(vari) < 1e-15 {
		err = chk.Err("cannot compute skew and kurtosis because variance is zero. vari=%g (tol=%g)", vari, TOL)
		return
	}
	skew /= (N * vari * sdev)
	kurt = kurt/(N*vari*vari) - 3.0
	return
}
