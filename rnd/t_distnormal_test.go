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
	dist.Init(μ, σ)

	n := 101
	x := utl.LinSpace(-2, 2, n)
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

func Test_norm01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("norm01")

	// values from R language: dnorm and pnorm functions. Example:
	//  options(digits=17)
	//  X = seq(-2, 2, 0.5)
	//  dlnorm(X, 0, 0.25)
	X := []float64{-2.0, -1.5, -1.0, -0.5, 0.0, 0.5, 1.0, 1.5, 2.0}
	pdf_μ0_σ05 := []float64{0.00026766045152977074, 0.00886369682387601505, 0.10798193302637612567, 0.48394144903828673066, 0.79788456080286540573, 0.48394144903828673066, 0.10798193302637612567, 0.00886369682387601505, 0.00026766045152977074}
	pdf_μ1_σ05 := []float64{1.2151765699646572e-08, 2.9734390294685954e-06, 2.6766045152977074e-04, 8.8636968238760151e-03, 1.0798193302637613e-01, 4.8394144903828673e-01, 7.9788456080286541e-01, 4.8394144903828673e-01, 1.0798193302637613e-01}
	cdf_μ0_σ05 := []float64{3.1671241833119924e-05, 1.3498980316300946e-03, 2.2750131948179212e-02, 1.5865525393145705e-01, 5.0000000000000000e-01, 8.4134474606854293e-01, 9.7724986805182079e-01, 9.9865010196836990e-01, 9.9996832875816688e-01}
	cdf_μ1_σ05 := []float64{9.8658764503769809e-10, 2.8665157187919391e-07, 3.1671241833119924e-05, 1.3498980316300946e-03, 2.2750131948179212e-02, 1.5865525393145705e-01, 5.0000000000000000e-01, 8.4134474606854293e-01, 9.7724986805182079e-01}

	var dist DistNormal
	dist.Init(0, 0.5)
	n := len(X)
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = dist.Pdf(X[i])
	}
	chk.Vector(tst, "pdf: μ=0 σ=0.50", 1e-15, x, pdf_μ0_σ05)

	dist.Init(1, 0.5)
	for i := 0; i < n; i++ {
		x[i] = dist.Pdf(X[i])
	}
	chk.Vector(tst, "pdf: μ=1 σ=0.50", 1e-15, x, pdf_μ1_σ05)

	dist.Init(0, 0.5)
	for i := 0; i < n; i++ {
		x[i] = dist.Cdf(X[i])
	}
	chk.Vector(tst, "cdf: μ=0 σ=0.50", 1e-15, x, cdf_μ0_σ05)

	dist.Init(1, 0.5)
	for i := 0; i < n; i++ {
		x[i] = dist.Cdf(X[i])
	}
	chk.Vector(tst, "cdf: μ=1 σ=0.50", 1e-15, x, cdf_μ1_σ05)
}

func Test_norm02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("norm02")

	doplot := chk.Verbose
	if doplot {
		plt.SetForEps(1.5, 300)
		for _, σ := range []float64{1, 0.5, 0.25} {
			plot_normal(0, σ)
		}
		plt.Subplot(2, 1, 1)
		plt.Gll("x", "f", "")
		plt.Subplot(2, 1, 2)
		plt.Gll("x", "F", "")
		plt.SaveD("/tmp/gosl", "test_norm02.eps")
	}
}
