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

func plot_frechet(l, c, a float64, xmin, xmax float64) {

	var dist DistFrechet
	dist.Init(&VarData{L: l, C: c, A: a})

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
	plt.SetYnticks(11)
	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, nil)
	plt.Gll("$x$", "$F(x)$", nil)
}

func Test_dist_frechet_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_frechet_01")

	_, dat, err := io.ReadTable("data/frechet.dat")
	if err != nil {
		tst.Errorf("cannot read comparison results:\n%v\n", err)
		return
	}

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	L, ok := dat["l"] // location
	if !ok {
		tst.Errorf("cannot get l values\n")
		return
	}
	C, ok := dat["c"] // scale
	if !ok {
		tst.Errorf("cannot get c values\n")
		return
	}
	A, ok := dat["a"] // shape
	if !ok {
		tst.Errorf("cannot get a values\n")
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

	var dist DistFrechet

	nx := len(X)
	for i := 0; i < nx; i++ {
		dist.Init(&VarData{L: L[i], C: C[i], A: A[i]})
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

func Test_dist_frechet_02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_frechet_02")

	doplot := chk.Verbose
	if doplot {
		plt.SetForEps(1.5, 300, nil)
		l := 0.0                // location
		C := []float64{1, 2.0}  // scale
		A := []float64{1, 2, 3} // shape
		for _, c := range C {
			for _, a := range A {
				plot_frechet(l, c, a, 0, 4)
			}
		}
		plt.SaveD("/tmp/gosl", "rnd_dist_frechet_02a.eps")
		plt.SetForEps(1.5, 300, nil)
		l = 0.5                // location
		C = []float64{1, 2.0}  // scale
		A = []float64{1, 2, 3} // shape
		for _, c := range C {
			for _, a := range A {
				plot_frechet(l, c, a, 0, 4)
			}
		}
		plt.SaveD("/tmp/gosl", "rnd_dist_frechet_02b.eps")
	}
}

func Test_frechet_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_frechet_03")

	μ := 10.0
	σ := 5.0
	δ := σ / μ
	d := 1.0 + δ*δ
	io.Pforan("μ=%v σ=%v δ=%v d=%v\n", μ, σ, δ, d)

	if chk.Verbose {
		plt.AxHline(d, nil)
		FrechetPlotCoef("/tmp/gosl", "fig_frechet_coef.eps", 3.0, 5.0)
	}

	k := 0.2441618
	α := 1.0 / k
	l := μ - math.Gamma(1.0-k)
	io.Pfpink("l=%v α=%v\n", l, α)

	l = 8.782275
	α = 4.095645

	var dist DistFrechet
	dist.Init(&VarData{L: l, A: α})
	io.Pforan("dist = %+#v\n", dist)
	io.Pforan("mean = %v\n", dist.Mean())
	io.Pforan("var  = %v\n", dist.Variance())
	io.Pforan("σ    = %v\n", math.Sqrt(dist.Variance()))

	if chk.Verbose {
		plot_frechet(l, 1, α, 8, 16)
		plt.SaveD("/tmp/gosl", "rnd_dist_frechet_03.eps")
	}
}
