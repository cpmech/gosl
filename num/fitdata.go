// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "math"

// LinFit computes linear fitting parameters. Errors on y-direction only
//
//   y(x) = a + b⋅x
//
//   See page 780 of [1]
//   Reference:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func LinFit(x, y []float64) (a, b float64) {

	// variables
	var sx, sy, t, st2 float64
	ndata := len(x)

	// accumulate sums
	for i := 0; i < ndata; i++ {
		sx += x[i]
		sy += y[i]
	}

	// compute b
	ss := float64(ndata)
	sxoss := sx / ss
	for i := 0; i < ndata; i++ {
		t = x[i] - sxoss
		st2 += t * t
		b += t * y[i]
	}
	b /= st2

	// compute a
	a = (sy - sx*b) / ss
	return
}

// LinFitSigma computes linear fitting parameters and variances (σa,σb) in the estimates of a and b
// Errors on y-direction only
//
//   y(x) = a + b⋅x
//
//   See page 780 of [1]
//   Reference:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func LinFitSigma(x, y []float64) (a, b, σa, σb, χ2 float64) {

	// variables
	var sx, sy, t, st2 float64
	ndata := len(x)

	// accumulate sums
	for i := 0; i < ndata; i++ {
		sx += x[i]
		sy += y[i]
	}

	// compute b
	ss := float64(ndata)
	sxoss := sx / ss
	for i := 0; i < ndata; i++ {
		t = x[i] - sxoss
		st2 += t * t
		b += t * y[i]
	}
	b /= st2

	// compute a
	a = (sy - sx*b) / ss

	// solve for σa and σb
	σa = math.Sqrt((1.0 + sx*sx/(ss*st2)) / ss)
	σb = math.Sqrt(1.0 / st2)

	// calculate χ².
	var d, σdat float64
	for i := 0; i < ndata; i++ {
		d = y[i] - a - b*x[i]
		χ2 += d * d
	}

	// for unweighted data evaluate typical sig using χ2,
	// and adjust the standard deviations.
	if ndata > 2 {
		σdat = math.Sqrt(χ2 / float64(ndata-2))
	}
	σa *= σdat
	σb *= σdat
	return
}
