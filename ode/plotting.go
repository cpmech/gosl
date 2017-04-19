// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Results structure to hold numerical results
type Results struct {
	Method string      // method name
	Dx     []float64   // step sizes [Noutput]
	X      []float64   // x values [Noutput]
	Y      [][]float64 // y values [ndim][Noutput]
}

// SimpleOutput implements a simple output function
func SimpleOutput(first bool, dx, x float64, y []float64, args ...interface{}) (err error) {
	chk.IntAssert(len(args), 1)
	res := args[0].(*Results)
	res.Dx = append(res.Dx, dx)
	res.X = append(res.X, x)
	ndim := len(y)
	if len(res.Y) == 0 {
		res.Y = utl.DblsAlloc(ndim, 0)
	}
	for j := 0; j < ndim; j++ {
		res.Y[j] = append(res.Y[j], y[j])
	}
	return
}

// Plot plot results
func Plot(dirout, fnkey string, res *Results, yfcn Cb_ycorr, xa, xb float64, extra func()) {

	// data
	if res == nil {
		return
	}
	ndim := len(res.Y)
	if ndim < 1 {
		return
	}

	// closed-form solution
	var xc []float64
	var Yc [][]float64
	if yfcn != nil {
		np := 101
		dx := (xb - xa) / float64(np-1)
		xc = make([]float64, np)
		Yc = utl.DblsAlloc(np, ndim)
		for i := 0; i < np; i++ {
			xc[i] = xa + dx*float64(i)
			yfcn(Yc[i], xc[i])
		}
	}

	// plot
	for j := 0; j < ndim; j++ {
		plt.Subplot(ndim+1, 1, j+1)
		if yfcn != nil {
			plt.Plot(xc, Yc[j], &plt.A{C: "y", Ls: "-", Lw: 6, L: "analytical"})
		}
		plt.Plot(res.X, res.Y[j], &plt.A{C: "b", Ls: "-", M: ".", Lw: 1, L: res.Method})
		plt.Gll("$x$", "$y$", nil)
	}
	plt.Subplot(ndim+1, 1, ndim+1)
	plt.Plot(res.X, res.Dx, &plt.A{C: "b", Ls: "-", M: ".", Lw: 1, L: res.Method})
	plt.SetYlog()
	plt.Gll("$x$", "$\\log(\\delta x)$", nil)

	// write file
	if extra != nil {
		extra()
	}
	plt.Save(dirout, fnkey)
}
