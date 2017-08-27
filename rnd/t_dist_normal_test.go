// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func plotNormal(μ, σ, xmin, xmax float64) {

	var dist DistNormal
	dist.Init(&Variable{M: μ, S: σ})

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
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, nil)
	plt.Gll("$x$", "$F(x)$", nil)
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
		dist.Init(&Variable{M: Mu[i], S: Sig[i]})
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
		plt.Reset(false, nil)
		for _, σ := range []float64{1, 0.5, 0.25} {
			plotNormal(0, σ, -2, 2)
		}
		plt.Save("/tmp/gosl", "rnd_dist_normal_02")
	}
}

func Test_dist_normal_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_normal_03")

	chk.Float64(tst, "φ(0)", 1e-16, Stdphi(0.0), 0.3989422804014327)
	chk.Float64(tst, "φ(2)", 1e-16, Stdphi(2.0), 0.053990966513188063)
	chk.Float64(tst, "φ(10)", 1e-16, Stdphi(10.0), 7.6945986267064199e-23)
	io.Pf("\n")
	chk.Float64(tst, "Φ(0)", 1e-16, StdPhi(0.0), 0.5)
	chk.Float64(tst, "Φ(2)", 1e-16, StdPhi(2.0), 0.97724986805182079)
	chk.Float64(tst, "Φ(4)", 1e-16, StdPhi(4.0), 0.99996832875816688)
	io.Pf("\n")
	chk.Float64(tst, "Φ⁻¹(Φ(0))", 1e-16, StdInvPhi(StdPhi(0.0)), 0.0)
	chk.Float64(tst, "Φ⁻¹(Φ(2))", 1e-9, StdInvPhi(StdPhi(2.0)), 2.0)
	chk.Float64(tst, "Φ⁻¹(Φ(4))", 1e-8, StdInvPhi(StdPhi(4.0)), 4.0)
	io.Pf("\n")
	chk.Float64(tst, "Φ⁻¹(Φ(0))", 1e-16, StdInvPhi(0.5), 0.0)
	chk.Float64(tst, "Φ⁻¹(Φ(2))", 1e-9, StdInvPhi(0.97724986805182079), 2.0)
	chk.Float64(tst, "Φ⁻¹(Φ(4))", 1e-8, StdInvPhi(0.99996832875816688), 4.0)
}

func Test_dist_normal_04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_normal_04. problem with Φ")

	if chk.Verbose {
		np := 101
		x := utl.LinSpace(0, 8.2, np)
		y := make([]float64, np)
		for i := 0; i < np; i++ {
			//io.Pforan("x=%v Φ(x)=%v Φ⁻¹(Φ(x))=%v\n", x[i], StdPhi(x[i]), StdInvPhi(StdPhi(x[i])))
			y[i] = StdInvPhi(StdPhi(x[i]))
		}
		plt.Plot(x, y, nil)
		plt.Gll("$x$", "$\\Phi^{-1}(\\Phi(x))$", nil)
		plt.Save("/tmp/gosl", "rnd_dist_normal_04")
	}
}

func Test_dist_normal_05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_normal_05. random numbers")

	μ := 1.0
	σ := 0.25

	nsamples := 10000
	X := make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		X[i] = Normal(μ, σ)
	}

	nstations := 41
	xmin := 0.0
	xmax := 3.0
	dx := (xmax - xmin) / float64(nstations-1)

	var hist Histogram
	hist.Stations = utl.LinSpace(xmin, xmax, nstations)
	hist.Count(X, true)

	prob := make([]float64, nstations)
	for i := 0; i < nstations-1; i++ {
		prob[i] = float64(hist.Counts[i]) / (float64(nsamples) * dx)
	}

	io.Pf(TextHist(hist.GenLabels("%.3f"), hist.Counts, 60))
	io.Pforan("dx = %v\n", dx)

	area := 0.0
	for i := 0; i < nstations-1; i++ {
		area += dx * prob[i]
	}
	io.Pforan("area = %v\n", area)
	chk.Float64(tst, "area", 1e-15, area, 1)

	if chk.Verbose {
		plt.Reset(false, nil)
		plotNormal(μ, σ, 0, 2)
		plt.Subplot(2, 1, 1)
		hist.PlotDensity(nil)
		plt.Save("/tmp/gosl", "rnd_dist_normal_05")
	}
}

func Test_dist_normal_06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_normal_06. transformation")

	doplot := chk.Verbose
	if doplot {

		vard := &Variable{M: 1.5, S: 0.1}
		vard.Distr = new(DistNormal)
		vard.Distr.Init(vard)

		npts := 1001
		X := utl.LinSpace(1, 2, npts)
		F, Y := make([]float64, npts), make([]float64, npts)
		for i := 0; i < npts; i++ {
			y, invalid := vard.Transform(X[i])
			if invalid {
				io.Pf("invalid: x=%g\n", X[i])
				y = math.NaN()
			}
			Y[i] = y
			F[i] = vard.Distr.Pdf(X[i])
		}

		plt.Reset(true, &plt.A{Prop: 1})

		plt.Subplot(2, 1, 1)
		plt.Plot(X, F, &plt.A{C: "#0046ba", Lw: 2, NoClip: true})
		plt.HideTRborders()
		plt.Gll("$x$", "$f(x)$", nil)
		plt.AxisXmin(1)

		plt.Subplot(2, 1, 2)
		plt.Plot(X, Y, &plt.A{C: "b", Lw: 2, NoClip: true})
		plt.HideTRborders()
		plt.SetYnticks(12)
		plt.AxisYrange(-5, 5)
		plt.Gll("$x$", "$y=T(x)$", nil)
		plt.AxisXmin(1)

		err := plt.Save("/tmp/gosl", "rnd_dist_normal_06")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}
