// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func plot_uniform(A, B float64, xmin, xmax float64) {

	var dist DistUniform
	dist.Init(&VarData{Min: A, Max: B})

	n := 101
	x := utl.LinSpace(xmin, xmax, n)
	y := make([]float64, n)
	Y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = dist.Pdf(x[i])
		Y[i] = dist.Cdf(x[i])
	}
	plt.Subplot(2, 1, 1)
	plt.Plot(x, y, nil)
	plt.Gll("$x$", "$f(x)$", nil)
	plt.SetYnticks(11)
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, nil)
	plt.Gll("$x$", "$F(x)$", nil)
}

func Test_dist_uniform_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_uniform_01")

	_, dat, err := io.ReadTable("data/uniform.dat")
	if err != nil {
		tst.Errorf("cannot read comparison results:\n%v\n", err)
		return
	}

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	A, ok := dat["a"] // min
	if !ok {
		tst.Errorf("cannot get a values\n")
		return
	}
	B, ok := dat["b"] // max
	if !ok {
		tst.Errorf("cannot get b values\n")
		return
	}
	YpdfCmp, ok := dat["ypdf"]
	if !ok {
		tst.Errorf("cannot get ypdf values\n")
		return
	}
	YcdfCmp, ok := dat["ycdf"]
	if !ok {
		tst.Errorf("cannot get ycdf values\n")
		return
	}

	var dist DistUniform

	nx := len(X)
	for i := 0; i < nx; i++ {
		dist.Init(&VarData{Min: A[i], Max: B[i]})
		Ypdf := dist.Pdf(X[i])
		Ycdf := dist.Cdf(X[i])
		err := chk.PrintAnaNum("ypdf", 1e-14, YpdfCmp[i], Ypdf, chk.Verbose)
		if err != nil {
			tst.Errorf("pdf failed: %v\n", err)
			return
		}
		err = chk.PrintAnaNum("ycdf", 1e-15, YcdfCmp[i], Ycdf, chk.Verbose)
		if err != nil {
			tst.Errorf("cdf failed: %v\n", err)
			return
		}
	}
}

func Test_dist_uniform_02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_uniform_02")

	doplot := chk.Verbose
	if doplot {
		plt.Reset(false, nil)
		A := 1.5 // min
		B := 2.5 // max
		plot_uniform(A, B, 1.0, 3.0)
		plt.Save("/tmp/gosl", "rnd_dist_uniform_02a")
	}
}
