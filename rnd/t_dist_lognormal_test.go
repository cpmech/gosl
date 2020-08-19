// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/plt"
	"gosl/utl"
)

func plotLognormal(μ, σ float64) {

	var dist DistLogNormal
	dist.Init(&Variable{M: μ, S: σ})

	n := 101
	x := utl.LinSpace(0, 3, n)
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

func Test_dist_lognormal_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_lognormal_01")

	_, dat := io.ReadTable("data/lognormal.dat")

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	N, ok := dat["n"]
	if !ok {
		tst.Errorf("cannot get n values\n")
		return
	}
	Z, ok := dat["z"]
	if !ok {
		tst.Errorf("cannot get z values\n")
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

	nx := len(X)
	for i := 0; i < nx; i++ {
		w := Z[i] * Z[i]
		μ := math.Exp(N[i] + w/2.0)
		σ := μ * math.Sqrt(math.Exp(w)-1.0)
		dist.Init(&Variable{M: μ, S: σ})
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

func Test_dist_lognormal_02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_lognormal_02")

	doplot := chk.Verbose
	if doplot {
		plt.Reset(false, nil)
		n := 0.0
		for _, z := range []float64{1, 0.5, 0.25} {
			w := z * z
			μ := math.Exp(n + w/2.0)
			σ := μ * math.Sqrt(math.Exp(w)-1.0)
			plotLognormal(μ, σ)
		}
		plt.Save("/tmp/gosl", "rnd_dist_lognormal_02")
	}
}

func Test_dist_lognormal_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_lognormal_03. random numbers")

	μ := 1.0
	σ := 0.25

	nsamples := 1000
	X := make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		X[i] = Lognormal(μ, σ)
	}

	nstations := 41
	xmin := 0.0
	xmax := 3.0

	var hist Histogram
	hist.Stations = utl.LinSpace(xmin, xmax, nstations)
	hist.Count(X, true)

	io.Pf(TextHist(hist.GenLabels("%.3f"), hist.Counts, 60))

	area := hist.DensityArea(nsamples)
	io.Pforan("area = %v\n", area)
	chk.Float64(tst, "area", 1e-15, area, 1)

	if chk.Verbose {
		plt.Reset(false, nil)
		plotLognormal(μ, σ)
		plt.Subplot(2, 1, 1)
		hist.PlotDensity(nil)
		plt.Save("/tmp/gosl", "rnd_dist_lognormal_03")
	}
}

func Test_dist_lognormal_04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_lognormal_04. transformation")

	doplot := chk.Verbose
	if doplot {

		vard := &Variable{M: 1.5, S: 0.1}
		vard.Distr = new(DistLogNormal)
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
		plt.Gll("$x$", "$y=T(x)$", nil)
		plt.AxisXmin(1)

		plt.Save("/tmp/gosl", "rnd_dist_lognormal_04")
	}
}
