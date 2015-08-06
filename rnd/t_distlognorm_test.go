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
	dist.Init(&VarData{M: μ, S: σ, MSlog: true})

	n := 101
	x := utl.LinSpace(0, 3, n)
	y := make([]float64, n)
	Y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = dist.Pdf(x[i])
		Y[i] = dist.Cdf(x[i])
	}
	plt.Subplot(2, 1, 1)
	plt.Plot(x, y, io.Sf("clip_on=0,zorder=10,label=r'$\\mu=%g,\\;\\sigma=%g$'", μ, σ))
	plt.Gll("$x$", "$f(x)$", "leg_out=1, leg_ncol=2")
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, io.Sf("clip_on=0,zorder=10,label=r'$\\mu=%g,\\;\\sigma=%g$'", μ, σ))
	plt.Gll("$x$", "$F(x)$", "leg_out=1, leg_ncol=2")
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
		dist.Init(&VarData{M: Mu[i], S: Sig[i], MSlog: true})
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
		plt.SaveD("/tmp/gosl", "test_lognorm02.eps")
	}
}

func Test_lognorm03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("lognorm03. Rackwitz-Fiessler conversion")

	dat := &VarData{D: D_Log, M: 10, S: 2}
	var dist DistLogNormal
	dist.Init(dat)
	dat.distr = &dist

	doplot := false
	if doplot {
		plt.SetForEps(1.5, 300)
		n := 101
		x := utl.LinSpace(5, 15, n)
		y := make([]float64, n)
		Y := make([]float64, n)
		for i := 0; i < n; i++ {
			y[i] = dist.Pdf(x[i])
			Y[i] = dist.Cdf(x[i])
		}
		plt.Subplot(2, 1, 1)
		plt.Plot(x, y, io.Sf("clip_on=0,zorder=10,label=r'$m=%.3f,\\;s=%.3f$'", dist.M, dist.S))
		plt.Gll("$x$", "$f(x)$", "leg_out=0, leg_ncol=2")
		plt.Subplot(2, 1, 2)
		plt.Plot(x, Y, io.Sf("clip_on=0,zorder=10,label=r'$m=%.3f,\\;s=%.3f$'", dist.M, dist.S))
		plt.Gll("$x$", "$F(x)$", "leg_out=0, leg_ncol=2")
		plt.SaveD("/tmp/gosl", "test_lognorm03.eps")
	}

	for i, x := range []float64{10, 20, 50} {

		// Rackwitz-Fiessler
		f := dist.Pdf(x)
		F := dist.Cdf(x)
		io.Pforan("\nx=%g  f(x)=%v  F(x)=%v\n", x, f, F)
		var σNrf, μNrf float64
		if F == 0 || F == 1 { // z = Φ⁻¹(F) → -∞ or +∞
			chk.Panic("cannot compute equivalent normal parameters @ %g because F=%g", x, F)
		} else {
			z := StdInvPhi(F)
			σNrf = Stdphi(z) / f
			μNrf = x - σNrf*z
			if μNrf < 0 {
				μNrf = 0
				σNrf = x / z
			}
		}

		// analytical solution for lognormal distribution
		μN, σN, invalid := dat.CalcEquiv(x)
		if invalid {
			tst.Errorf("CalcEquiv failed\n")
			return
		}

		// check
		tol := 1e-10
		if i > 0 {
			tol = 1e-6
		}
		err := chk.PrintAnaNum("μN", tol, μN, μNrf, chk.Verbose)
		if err != nil {
			tst.Errorf("μN values are different: %v\n", err)
			//return
		}
		err = chk.PrintAnaNum("σN", tol, σN, σNrf, chk.Verbose)
		if err != nil {
			tst.Errorf("σN values are different: %v\n", err)
			//return
		}
	}
}
