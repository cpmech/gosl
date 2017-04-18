// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"os"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	DBL_EPSILON = 1.0e-15
	DBL_MIN     = math.SmallestNonzeroFloat64
)

const (
	NaN = iota
	Inf
	Equal
	NotEqual
)

// PlotYxe plots the function y(x) implemented by Cb_yxe
func PlotYxe(ffcn Cb_yxe, dirout, fname string, xsol, xa, xb float64, np int, xsolLbl, args string, save, show bool, extra func()) (err error) {
	if !save && !show {
		return
	}
	x := utl.LinSpace(xa, xb, np)
	y := make([]float64, np)
	for i := 0; i < np; i++ {
		y[i], err = ffcn(x[i])
		if err != nil {
			return
		}
	}
	var ysol float64
	ysol, err = ffcn(xsol)
	if err != nil {
		return
	}
	plt.Cross(0, 0, nil)
	plt.Plot(x, y, nil)
	plt.PlotOne(xsol, ysol, &plt.A{C: "r", M: "o", L: xsolLbl})
	if extra != nil {
		extra()
	}
	plt.Gll("x", "y(x)", nil)
	if save {
		os.MkdirAll(dirout, 0777)
		plt.Save(dirout + "/" + fname)
	}
	if show {
		plt.Show()
	}
	return
}

func MinComp(tol, expected float64) float64 {
	return min(tol, math.Abs(expected)) + DBL_EPSILON
}

func TestAbs(result, expected, absolute_error float64, test_description string) (status int) {
	switch {
	case math.IsNaN(result) || math.IsNaN(expected):
		status = NaN

	case math.IsInf(result, 0) || math.IsInf(expected, 0):
		status = Inf

	case (expected > 0 && expected < DBL_MIN) || (expected < 0 && expected > -DBL_MIN):
		status = NotEqual

	default:
		if math.Abs(result-expected) > absolute_error {
			status = NotEqual
		} else {
			status = Equal
		}
	}
	if test_description != "" {
		io.Pf(test_description)
		switch status {
		case NaN:
			io.Pf(" [1;31mNaN[0m\n  %v observed\n  %v expected.  diff = %v\n", result, expected, result-expected)
		case Inf:
			io.Pf(" [1;31mInf[0m\n  %v observed\n  %v expected.  diff = %v\n", result, expected, result-expected)
		case Equal:
			io.Pf(" [1;32mOk[0m\n  %v observed\n  %v expected.  diff = %v\n", result, expected, result-expected)
		case NotEqual:
			io.Pf(" [1;31mError[0m\n  %v observed\n  %v expected.  diff = %v\n", result, expected, result-expected)
		}
	}
	return
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func sign(x float64) float64 {
	if x < 0.0 {
		return -1.0
	}
	if x > 0.0 {
		return 1.0
	}
	return 0.0
}
