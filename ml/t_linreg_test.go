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
	d := NewDataMatrixTable([][]float64{
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

	// check data
	chk.Int(tst, "m", d.Nsamples(), 20)
	chk.Int(tst, "nf", d.Nfeatures(), 1)

	// check stat
	d.stat()
	chk.Array(tst, "x[:,0]", 1e-15, d.xMat.GetCol(0), utl.Vals(20, 1))
	checkStat := func() {
		chk.Float64(tst, "min(x)", 1e-15, d.minX[0], 0.87)
		chk.Float64(tst, "max(x)", 1e-15, d.maxX[0], 1.55)
		chk.Float64(tst, "mean(x)", 1e-15, d.meanX[0], 1.1960)
		chk.Float64(tst, "sig(x)", 1e-15, d.sigX[0], 0.189303432281837)
		chk.Float64(tst, "sum(x)", 1e-15, d.sumX[0], 23.92)
		chk.Float64(tst, "min(y)", 1e-15, d.minY, 87.33)
		chk.Float64(tst, "max(y)", 1e-15, d.maxY, 99.42)
		chk.Float64(tst, "mean(y)", 1e-15, d.meanY, 92.1605)
		chk.Float64(tst, "sig(y)", 1e-15, d.sigY, 3.020778001913102)
		chk.Float64(tst, "sum(y)", 1e-15, d.sumY, 1843.21)
	}
	checkStat()

	// analytical θ
	io.Pl()
	r := NewLinReg()
	r.CalcTheta(d)
	io.Pf("analytical: θ = %v\n", d.params)
	chk.Float64(tst, "analytical: θ0", 1e-5, d.params[0], 74.28331)
	chk.Float64(tst, "analytical: θ1", 1e-5, d.params[1], 14.94748)

	// check dCdθ
	io.Pl()
	dCdθ := la.NewVector(d.Nparams())
	for _, θ0 := range []float64{50, 80} {
		for _, θ1 := range []float64{5, 20} {
			d.params[0] = θ0
			d.params[1] = θ1
			r.Deriv(dCdθ, d)
			chk.DerivScaVec(tst, "dCdθ", 1e-6, dCdθ, d.params, 1e-6, chk.Verbose, func(th []float64) float64 {
				copy(d.params, th)
				return r.Cost(d)
			})
		}
	}

	// train linear regression
	io.Pl()
	g := NewGradDesc(10)
	g.SetControl(0.1, 0, 0)
	g.Run(d, r, []float64{70, 10})
	io.Pf("grad.desc: θ = %v\n", d.params)
	chk.Float64(tst, "grad.desc: θ0", 1e-5, d.params[0], 73.91321)
	chk.Float64(tst, "grad.desc: θ1", 1e-5, d.params[1], 14.74272)

	// plot: unormalised model
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})

		plt.Subplot(3, 1, 1)
		plt.Plot(d.GetXvalues(0), d.GetYvalues(), &plt.A{C: plt.C(2, 0), M: ".", Ls: "none", NoClip: true})
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
		plt.PlotOne(d.params[0], d.params[1], &plt.A{C: plt.C(4, 0), M: "o", NoClip: true})
		plt.Save("/tmp/gosl/ml", "linreg01a")
	}

	// normalize
	io.Pf("\n. . . . normalized model . . . . . . . . .\n")
	d.Normalize(false)
	checkStat()
	r.CalcTheta(d)
	io.Pf("analytical: θ = %v\n", d.params)

	// plot: normalised model
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.3})

		plt.Subplot(2, 1, 1)
		plt.Plot(d.GetXvalues(0), d.GetYvalues(), &plt.A{C: plt.C(2, 0), M: ".", Ls: "none", NoClip: true})
		d.PlotModel(r, 0, 11, nil)
		plt.Gll("$x$", "$y$", nil)
		plt.HideTRborders()

		args := &plt.A{Nlevels: 20}
		plt.Subplot(2, 1, 2)
		d.PlotContCost(r, 0, 1, 11, 11, []float64{50, -50}, []float64{150, 50}, args)
		plt.PlotOne(d.params[0], d.params[1], &plt.A{C: plt.C(4, 0), M: "o", NoClip: true})
		plt.Save("/tmp/gosl/ml", "linreg01b")
	}
}
