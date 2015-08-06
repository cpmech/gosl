// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
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

func plot_lognormal(μ, σ float64) {

	var dist DistLogNormal
	dist.Init(&VarData{M: μ, S: σ})

	n := 101
	x := utl.LinSpace(0, 3, n)
	y := make([]float64, n)
	Y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = dist.Pdf(x[i])
		Y[i] = dist.Cdf(x[i])
	}
	plt.Subplot(2, 1, 1)
	plt.Plot(x, y, io.Sf("label='$\\mu=%g,\\;\\sigma=%g$'", μ, σ))
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, io.Sf("label='$\\mu=%g,\\; \\sigma=%g$'", μ, σ))
}

func Test_lognorm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("lognorm01")

	_, dat, err := io.ReadTable("data/log.dat")
	if err != nil {
		tst.Errorf("cannot read comparison results:\n%v\n", err)
		return
	}

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	Mu, ok := dat["mu"]
	if !ok {
		tst.Errorf("cannot get mu values\n")
		return
	}
	Sig, ok := dat["sig"]
	if !ok {
		tst.Errorf("cannot get sig values\n")
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

	var dist DistLogNormal

	n := len(X)
	for i := 0; i < n; i++ {
		dist.Init(&VarData{M: Mu[i], S: Sig[i]})
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

func Test_lognorm02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("lognorm02")

	doplot := chk.Verbose
	if doplot {
		plt.SetForEps(1.5, 300)
		for _, σ := range []float64{1, 0.5, 0.25} {
			plot_lognormal(0, σ)
		}
		plt.Subplot(2, 1, 1)
		plt.Gll("x", "f", "")
		plt.Subplot(2, 1, 2)
		plt.Gll("x", "F", "")
		plt.SaveD("/tmp/gosl", "test_lognorm02.eps")
	}
}
