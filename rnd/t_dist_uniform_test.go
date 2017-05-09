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

func plot_uniform(A, B float64, xmin, xmax float64) {

	var dist DistUniform
	dist.Init(&VarData{Min: A, Max: B})

	n := 101
	x := utl.LinSpace(xmin, xmax, n)
	y := make([]float64, n)
	Y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = dist.Pdf(x[i])
		Y[i] = dist.Cdf(x[i])
	}

	plt.Subplot(2, 1, 1)
	plt.Plot(x, y, &plt.A{Lw: 2, NoClip: true})
	plt.HideAllBorders()
	plt.SetYnticks(11)
	plt.Gll("$x$", "$f(x)$", nil)

	plt.Subplot(2, 1, 2)
	plt.Plot(x, Y, &plt.A{Lw: 2, NoClip: true})
	plt.HideAllBorders()
	plt.Gll("$x$", "$F(x)$", nil)
}

func Test_dist_uniform_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_uniform_01")

	_, dat, err := io.ReadTable("data/uniform.dat")
	if err != nil {
		tst.Errorf("cannot read comparison results:\n%v\n", err)
		return
	}

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	A, ok := dat["a"] // min
	if !ok {
		tst.Errorf("cannot get a values\n")
		return
	}
	B, ok := dat["b"] // max
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

	var dist DistUniform

	nx := len(X)
	for i := 0; i < nx; i++ {
		dist.Init(&VarData{Min: A[i], Max: B[i]})
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

func Test_dist_uniform_02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_uniform_02. density and cumulative distr.")

	doplot := chk.Verbose
	if doplot {
		plt.Reset(false, nil)
		A := 1.5 // min
		B := 2.5 // max
		plot_uniform(A, B, 1.0, 3.0)
		plt.Save("/tmp/gosl", "rnd_dist_uniform_02a")
	}
}

func Test_dist_uniform_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_uniform_03. transformation")

	doplot := chk.Verbose
	if doplot {

		vard := &VarData{D: D_Uniform, Min: 1.5, Max: 2.5}
		vard.Distr = new(DistUniform)
		vard.Distr.Init(vard)

		X := utl.LinSpace(vard.Min, vard.Max, 10001)
		Y := utl.GetMapped(X, func(x float64) float64 {
			y, invalid := vard.Transform(x)
			if invalid {
				io.Pf("invalid: x=%g\n", x)
				y = math.NaN()
			}
			return y
		})

		plt.Reset(true, &plt.A{Prop: 0.75})
		plt.AxVline(1.5, &plt.A{C: "k", Ls: "--"})
		plt.AxVline(2.5, &plt.A{C: "k", Ls: "--"})
		plt.Plot(X, Y, &plt.A{C: "b", Lw: 2, NoClip: true})
		plt.SetTicksXlist([]float64{1.4, 1.5, 1.6, 1.8, 2.0, 2.2, 2.4, 2.5, 2.6})
		plt.HideTRborders()
		plt.Gll("$x$", "$y=T(x)$", nil)
		plt.Save("/tmp/gosl", "rnd_dist_uniform_03")
	}
}
