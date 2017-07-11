// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

func plotNormal(μ, σ, xmin, xmax float64) {

	var dist rnd.DistNormal
	dist.Init(&rnd.VarData{M: μ, S: σ})

	n := 101
	x := utl.LinSpace(xmin, xmax, n)
	y := make([]float64, n)
	Y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = dist.Pdf(x[i])
		Y[i] = dist.Cdf(x[i])
	}

	args1 := &plt.A{L: io.Sf("mu=%g, sigma=%g", μ, σ)}
	args2 := &plt.A{LegOut: true, LegNcol: 2}

	plt.Subplot(2, 1, 1)
	plt.Plot(x, y, args1)
	plt.Gll("x", "f(x)", args2)
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, args1)
	plt.Gll("x", "F(x)", args2)
}

func main() {

	μ := 1.0
	σ := 0.25

	nsamples := 10000
	X := make([]float64, nsamples)
	for i := 0; i < nsamples; i++ {
		X[i] = rnd.Normal(μ, σ)
	}

	nstations := 41
	xmin := 0.0
	xmax := 2.0
	dx := (xmax - xmin) / float64(nstations-1)

	var hist rnd.Histogram
	hist.Stations = utl.LinSpace(xmin, xmax, nstations)
	hist.Count(X, true)

	prob := make([]float64, nstations)
	for i := 0; i < nstations-1; i++ {
		prob[i] = float64(hist.Counts[i]) / (float64(nsamples) * dx)
	}

	io.Pf(rnd.TextHist(hist.GenLabels("%.3f"), hist.Counts, 60))
	io.Pforan("dx = %v\n", dx)

	area := 0.0
	for i := 0; i < nstations-1; i++ {
		area += dx * prob[i]
	}
	io.Pforan("area = %v\n", area)

	plt.Reset(false, nil)
	plotNormal(μ, σ, xmin, xmax)
	plt.Subplot(2, 1, 1)
	hist.PlotDensity(nil)
	plt.Save("/tmp/gosl", "rnd_normalDistribution")
}
