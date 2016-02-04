// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
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
	plt.Plot(x, y, io.Sf("clip_on=0,zorder=10,label=r'$\\mu=%.4f,\\;\\sigma=%.4f$'", μ, σ))
	plt.Gll("$x$", "$f(x)$", "leg_out=1, leg_ncol=2")
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, io.Sf("clip_on=0,zorder=10,label=r'$\\mu=%.4f,\\;\\sigma=%.4f$'", μ, σ))
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
		dist.Init(&VarData{M: μ, S: σ})
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
		n := 0.0
		for _, z := range []float64{1, 0.5, 0.25} {
			w := z * z
			μ := math.Exp(n + w/2.0)
			σ := μ * math.Sqrt(math.Exp(w)-1.0)
			plot_lognormal(μ, σ)
		}
		plt.SaveD("/tmp/gosl", "test_lognorm02.eps")
	}
}

func Test_lognorm03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("lognorm03. random numbers")

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
	chk.Scalar(tst, "area", 1e-15, area, 1)

	if chk.Verbose {
		plot_lognormal(μ, σ)
		plt.Subplot(2, 1, 1)
		hist.PlotDensity(nil, "")
		plt.SaveD("/tmp/gosl", "test_lognorm03.eps")
	}
}
