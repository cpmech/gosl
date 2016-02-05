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

func plot_normal(μ, σ float64) {

	var dist DistNormal
	dist.Init(&VarData{M: μ, S: σ})

	n := 101
	x := utl.LinSpace(-2, 2, n)
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

func Test_dist_normal_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_normal_01")

	_, dat, err := io.ReadTable("data/normal.dat")
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

	var dist DistNormal

	n := len(X)
	for i := 0; i < n; i++ {
		dist.Init(&VarData{M: Mu[i], S: Sig[i]})
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

func Test_dist_normal_02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_normal_02")

	doplot := chk.Verbose
	if doplot {
		plt.SetForEps(1.5, 300)
		for _, σ := range []float64{1, 0.5, 0.25} {
			plot_normal(0, σ)
		}
		plt.SaveD("/tmp/gosl", "nrd_dist_normal_02.eps")
	}
}

func Test_dist_normal_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_normal_03")

	chk.Scalar(tst, "φ(0)", 1e-16, Stdphi(0.0), 0.3989422804014327)
	chk.Scalar(tst, "φ(2)", 1e-16, Stdphi(2.0), 0.053990966513188063)
	chk.Scalar(tst, "φ(10)", 1e-16, Stdphi(10.0), 7.6945986267064199e-23)
	io.Pf("\n")
	chk.Scalar(tst, "Φ(0)", 1e-16, StdPhi(0.0), 0.5)
	chk.Scalar(tst, "Φ(2)", 1e-16, StdPhi(2.0), 0.97724986805182079)
	chk.Scalar(tst, "Φ(4)", 1e-16, StdPhi(4.0), 0.99996832875816688)
	io.Pf("\n")
	chk.Scalar(tst, "Φ⁻¹(Φ(0))", 1e-16, StdInvPhi(StdPhi(0.0)), 0.0)
	chk.Scalar(tst, "Φ⁻¹(Φ(2))", 1e-9, StdInvPhi(StdPhi(2.0)), 2.0)
	chk.Scalar(tst, "Φ⁻¹(Φ(4))", 1e-8, StdInvPhi(StdPhi(4.0)), 4.0)
	io.Pf("\n")
	chk.Scalar(tst, "Φ⁻¹(Φ(0))", 1e-16, StdInvPhi(0.5), 0.0)
	chk.Scalar(tst, "Φ⁻¹(Φ(2))", 1e-9, StdInvPhi(0.97724986805182079), 2.0)
	chk.Scalar(tst, "Φ⁻¹(Φ(4))", 1e-8, StdInvPhi(0.99996832875816688), 4.0)
}
