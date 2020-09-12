// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"time"

	"gosl/chk"
	"gosl/utl"
)

// StatAve computes the average of x values
//  Input:
//   x -- sample
//  Output:
//   xave -- average
func StatAve(x []float64) (xave float64) {
	n := len(x)
	if n < 1 {
		return
	}
	for i := 0; i < n; i++ {
		xave += x[i]
	}
	xave /= float64(n)
	return
}

// StatDevFirst computes the average deviation or standard deviation (σ)
// for given value of average/mean/first moment
//  Input:
//   x    -- sample
//   xave -- bar(x) == average/mean/first moment
//   std  -- compute standard deviation (σ) instead of average deviation (adev)
//  Output:
//   xdev -- average deviation; if std==true, computes standard deviation (σ) instead
func StatDevFirst(x []float64, xave float64, std bool) (xdev float64) {

	// check
	n := len(x)
	if n < 2 {
		return
	}
	N := float64(n)

	// standard deviation
	if std {
		var d, c, vari float64
		for i := 0; i < n; i++ {
			d = x[i] - xave // d ← xi - bar(x)
			c += d          // c ← Σ d  (corrector)
			vari += d * d   // vari ← Σ d²
		}
		vari = (vari - c*c/N) / (N - 1.0)
		xdev = math.Sqrt(vari)
		return
	}

	// average deviation
	var sum, d float64
	for i := 0; i < n; i++ {
		d = x[i] - xave    // d ← xi - bar(x)
		sum += math.Abs(d) // sum ← Σ |d|
	}
	xdev = sum / N
	return
}

// StatDev computes the average deviation or standard deviation (σ)
//  Input:
//   x   -- sample
//   std -- compute standard deviation (σ) instead of average deviation (adev)
//  Output:
//   xdev -- average deviation; if std==true, computes standard deviation (σ) instead
func StatDev(x []float64, std bool) (xdev float64) {

	// check
	n := len(x)
	if n < 2 {
		return
	}

	// average
	var xave float64
	for i := 0; i < n; i++ {
		xave += x[i]
	}
	xave /= float64(n)
	xdev = StatDevFirst(x, xave, std)
	return
}

// StatAveDev computes the average of x and the average deviation or standard deviation (σ)
//  Input:
//   x   -- sample
//   std -- compute standard deviation (σ) instead of average deviation (adev)
//  Output:
//   xdev -- average deviation; if std==true, computes standard deviation (σ) instead
func StatAveDev(x []float64, std bool) (xave, xdev float64) {

	// check
	n := len(x)
	if n < 2 {
		return
	}

	// average
	for i := 0; i < n; i++ {
		xave += x[i]
	}
	xave /= float64(n)
	xdev = StatDevFirst(x, xave, std)
	return
}

// StatBasic performs some basic statistics
//  Input:
//   x   -- sample
//   std -- compute standard deviation (σ) instead of average deviation (adev)
//  Output:
//   xmin -- minimum value
//   xave -- mean average (first moment)
//   xmax -- maximum value
//   xdev -- average deviation; if std==true, computes standard deviation (σ) instead
func StatBasic(x []float64, std bool) (xmin, xave, xmax, xdev float64) {

	// check
	n := len(x)
	if n < 2 {
		return
	}

	// average, min and max
	xmin, xmax = x[0], x[0]
	for i := 0; i < n; i++ {
		xave += x[i]
		xmin = utl.Min(xmin, x[i])
		xmax = utl.Max(xmax, x[i])
	}
	xave /= float64(n)
	xdev = StatDevFirst(x, xave, std)
	return
}

// StatMoments computes the 4th moments of a data set
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
func StatMoments(x []float64) (sum, mean, adev, sdev, vari, skew, kurt float64) {

	// check
	n := len(x)
	if n < 2 {
		chk.Panic("x set must have at least 2 items\n")
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
		chk.Panic("cannot compute skew and kurtosis because variance is zero. vari=%g (tol=%g)\n", vari, TOL)
	}
	skew /= (N * vari * sdev)
	kurt = kurt/(N*vari*vari) - 3.0
	return
}

// StatTable computes the min, ave, max, and dev of values organised in a table
//  Input:
//   x     -- sample
//   std   -- compute standard deviation (σ) instead of average deviation (adev)
//   withZ -- computes z-matrix as well
//  Convention of indices:
//   0=min  1=ave  2=max  3=dev
//  Output:                        min          ave          max          dev
//                                  ↓            ↓            ↓            ↓
//   x00 x01 x02 x03 x04 x05 → y00=min(x0?) y10=ave(x0?) y20=max(x0?) y30=dev(x0?)
//   x10 x11 x12 x13 x14 x15 → y01=min(x1?) y11=ave(x1?) y21=max(x1?) y31=dev(x1?)
//   x20 x21 x22 x23 x24 x25 → y02=min(x2?) y12=ave(x2?) y22=max(x2?) y32=dev(x2?)
//                                  ↓            ↓            ↓            ↓
//                       min → z00=min(y0?) z01=min(y1?) z02=min(y2?) z03=min(y3?)
//                       ave → z10=ave(y0?) z11=ave(y1?) z12=ave(y2?) z13=ave(y3?)
//                       max → z20=max(y0?) z21=max(y1?) z22=max(y2?) z23=max(y3?)
//                       dev → z30=dev(y0?) z31=dev(y1?) z32=dev(y2?) z33=dev(y3?)
//                                  =            =            =            =
//                       min → z00=min(min) z01=min(ave) z02=min(max) z03=min(dev)
//                       ave → z10=ave(min) z11=ave(ave) z12=ave(max) z13=ave(dev)
//                       max → z20=max(min) z21=max(ave) z22=max(max) z23=max(dev)
//                       dev → z30=dev(min) z31=dev(ave) z32=dev(max) z33=dev(dev)
func StatTable(x [][]float64, std, withZ bool) (y, z [][]float64) {

	// dimensions
	m := len(x)
	if m < 1 {
		return
	}
	n := len(x[0])
	if n < 2 {
		return
	}

	// compute y
	y = utl.Alloc(4, m)
	for i := 0; i < m; i++ {
		y[0][i], y[1][i], y[2][i], y[3][i] = StatBasic(x[i], std)
	}

	// compute z
	if withZ {
		z = utl.Alloc(4, 4)
		for i := 0; i < 4; i++ {
			z[0][i], z[1][i], z[2][i], z[3][i] = StatBasic(y[i], std)
		}
	}
	return
}

// StatDur generates stat about duration
func StatDur(durs []time.Duration) (min, ave, max, sum time.Duration) {
	if len(durs) == 0 {
		return
	}
	min, max = durs[0], durs[0]
	for _, d := range durs {
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
		sum += d
	}
	ave = sum / time.Duration(int64(len(durs)))
	return
}
