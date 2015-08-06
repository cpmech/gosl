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

func plot_gev(μ, σ, ξ float64) {

	var dist DistGev
	dist.Init(&VarData{M: μ, S: σ, K: ξ})

	n := 101
	x := utl.LinSpace(-4, 4, n)
	y := make([]float64, n)
	Y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = dist.Pdf(x[i])
		Y[i] = dist.Cdf(x[i])
	}
	plt.Subplot(2, 1, 1)
	plt.Plot(x, y, io.Sf("clip_on=0,zorder=10,label=r'$\\mu=%g,\\;\\sigma=%g,\\;\\xi=%g$'", μ, σ, ξ))
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, io.Sf("clip_on=0,zorder=10,label=r'$\\mu=%g,\\;\\sigma=%g,\\;\\xi=%g$'", μ, σ, ξ))
}

func Test_gev01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("gev01")

	_, dat, err := io.ReadTable("data/gev.dat")
	if err != nil {
		tst.Errorf("cannot read comparison results:\n%v\n", err)
		return
	}

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	Ksi, ok := dat["ksi"]
	if !ok {
		tst.Errorf("cannot get ksi values\n")
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

	var dist DistGev

	n := len(X)
	for i := 0; i < n; i++ {
		dist.Init(&VarData{M: Mu[i], S: Sig[i], K: Ksi[i]})
		Ypdf := dist.Pdf(X[i])
		Ycdf := dist.Cdf(X[i])
		err := chk.PrintAnaNum("ypdf", 1e-15, YpdfCmp[i], Ypdf, chk.Verbose)
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

func Test_gev02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("gev02")

	doplot := chk.Verbose
	if doplot {
		plt.SetForEps(1.5, 300)
		for _, ξ := range []float64{0.5, -0.5, 0} {
			plot_gev(0, 1, ξ)
		}
		plt.Subplot(2, 1, 1)
		plt.Gll("$x$", "$f$", "leg_out=1, leg_ncol=2")
		plt.Subplot(2, 1, 2)
		plt.Gll("$x$", "$F$", "leg_out=1, leg_ncol=2")
		plt.SaveD("/tmp/gosl", "test_gev02.eps")
	}
}
