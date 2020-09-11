// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/utl"
)

func Test_dist_normal_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_normal_01")

	_, dat := io.ReadTable("data/normal.dat")

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
}
