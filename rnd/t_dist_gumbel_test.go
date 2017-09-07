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

func plotGumbel(μ, σ float64) {

	var dist DistGumbel
	dist.Init(&Variable{M: μ, S: σ})

	n := 101
	x := utl.LinSpace(-5, 20, n)
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

func Test_dist_gumbel_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_gumbel_01")

	_, dat, err := io.ReadTable("data/gumbel.dat")
	if err != nil {
		tst.Errorf("cannot read comparison results:\n%v\n", err)
		return
	}

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	U, ok := dat["u"]
	if !ok {
		tst.Errorf("cannot get u values\n")
		return
	}
	B, ok := dat["b"]
	if !ok {
		tst.Errorf("cannot get b values\n")
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

	var dist DistGumbel

	nx := len(X)
	for i := 0; i < nx; i++ {
		dist.U = U[i]
		dist.B = B[i]
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

func Test_dist_gumbel_02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_gumbel_02")

	doplot := chk.Verbose
	if doplot {
		plt.Reset(false, nil)
		U := []float64{1.5, 1.0, 0.5, 3.0}
		B := []float64{3.0, 2.0, 2.0, 4.0}
		euler := 0.57721566490153286060651209008240243104215
		for i, u := range U {
			σ := B[i] * math.Pi / math.Sqrt(6.0)
			μ := u + euler*B[i]
			plotGumbel(μ, σ)
		}
		plt.Save("/tmp/gosl", "rnd_dist_gumbel_02")
	}
}

func Test_gumbel_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_gumbel_03")

	var dist DistGumbel
	dist.Init(&Variable{M: 61.3, S: 7.52}) // from Haldar & Mahadevan page 90
	io.Pforan("dist = %+#v\n", dist)
	chk.Float64(tst, "u", 0.00011, dist.U, 57.9157)
	chk.Float64(tst, "β", 1e-4, dist.B, 1.0/0.17055)
}

func Test_dist_gumbel_04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_gumbel_04. transformation")

	doplot := chk.Verbose
	if doplot {

		vard := &Variable{M: 1.5, S: 0.1}
		vard.Distr = new(DistGumbel)
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

		plt.Save("/tmp/gosl", "rnd_dist_gumbel_04")
	}
}
