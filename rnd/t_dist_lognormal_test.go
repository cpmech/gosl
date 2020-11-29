// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

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
}
