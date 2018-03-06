// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestLinReg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinReg01. Linear Regression")

	// Reference:
	//  [1] Montgomery & Runger (2014) Applied Statistics and Probabilities for Engineers, Wiley

	// data from page 433 of [1]
	d := NewRegDataTable([][]float64{
		{0.99, 90.01},
		{1.02, 89.05},
		{1.15, 91.43},
		{1.29, 93.74},
		{1.46, 96.73},
		{1.36, 94.45},
		{0.87, 87.59},
		{1.23, 91.77},
		{1.55, 99.42},
		{1.40, 93.65},
		{1.19, 93.54},
		{1.15, 92.52},
		{0.98, 90.56},
		{1.01, 89.54},
		{1.11, 89.85},
		{1.20, 90.39},
		{1.26, 93.25},
		{1.32, 93.41},
		{1.43, 94.98},
		{0.95, 87.33},
	})

	// check stat
	d.Stat()
	chk.Array(tst, "x[:,0]", 1e-15, d.x.GetCol(0), utl.Vals(20, 1))
	chk.Float64(tst, "min(x)", 1e-15, d.MinX[0], 0.87)
	chk.Float64(tst, "max(x)", 1e-15, d.MaxX[0], 1.55)
	chk.Float64(tst, "mean(x)", 1e-15, d.MeanX[0], 1.1960)
	chk.Float64(tst, "sig(x)", 1e-15, d.SigX[0], 0.189303432281837)
	chk.Float64(tst, "sum(x)", 1e-15, d.SumX[0], 23.92)
	chk.Float64(tst, "min(y)", 1e-15, d.MinY, 87.33)
	chk.Float64(tst, "max(y)", 1e-15, d.MaxY, 99.42)
	chk.Float64(tst, "mean(y)", 1e-15, d.MeanY, 92.1605)
	chk.Float64(tst, "sig(y)", 1e-15, d.SigY, 3.020778001913102)
	chk.Float64(tst, "sum(y)", 1e-15, d.SumY, 1843.21)

	// analytical θ
	io.Pl()
	r := NewLinReg()
	r.CalcTheta(d)
	io.Pf("analytical: θ = %v\n", d.θ)
	chk.Float64(tst, "θ0", 1e-5, d.θ[0], 74.28331)
	chk.Float64(tst, "θ1", 1e-5, d.θ[1], 14.94748)

	// check dCdθ
	io.Pl()
	dCdθ := la.NewVector(d.Nparams())
	for _, θ0 := range []float64{50, 80} {
		for _, θ1 := range []float64{5, 20} {
			d.θ[0] = θ0
			d.θ[1] = θ1
			r.Deriv(dCdθ, d)
			chk.DerivScaVec(tst, "dCdθ", 1e-6, dCdθ, d.θ, 1e-6, chk.Verbose, func(th []float64) float64 {
				copy(d.θ, th)
				return r.Cost(d)
			})
		}
	}

	// train linear regression
	io.Pl()
	g := NewGradDesc(10)
	g.SetControl(0.1, 0, 0)
	g.Run(d, r, []float64{70, 10})
	io.Pf("grad.desc: θ = %v\n", d.θ)
	chk.Float64(tst, "θ0", 1e-5, d.θ[0], 73.91321)
	chk.Float64(tst, "θ1", 1e-5, d.θ[1], 14.74272)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})

		plt.Subplot(3, 1, 1)
		plt.Plot(d.GetFeature(0), d.GetY(), &plt.A{C: plt.C(2, 0), M: ".", Ls: "none", NoClip: true})
		d.PlotModel(r, 0, 11, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.HideTRborders()

		plt.Subplot(3, 1, 2)
		g.PlotCostIter(nil)
		plt.Gll("$iteration$", "$cost$", nil)
		plt.HideTRborders()

		args := &plt.A{Nlevels: 20}
		plt.Subplot(3, 1, 3)
		d.PlotContCost(r, 0, 1, 11, 11, []float64{0, 0}, []float64{100, 70}, args)
		plt.PlotOne(d.θ[0], d.θ[1], &plt.A{C: plt.C(4, 0), M: "o", NoClip: true})
		plt.Save("/tmp/gosl/ml", "linreg01")
	}

	// normalize
	io.Pl()
	d.Normalize(false)
	d.Stat()
	chk.Float64(tst, "normalized: mean(x)", 1e-15, d.MeanX[0], 0)
	chk.Float64(tst, "normalized: sig(x)", 1e-15, d.SigX[0], 1)
}
