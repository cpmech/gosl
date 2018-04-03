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
)

func TestLinReg01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinReg01a. Basic functionality (no regularizaton).")

	// data and parameters
	data := NewDataGivenRawXY(dataReg01)
	params := NewParamsReg(data.Nfeatures)

	// regression
	reg := NewLinReg(data, params, "reg01")

	// check stat
	chk.Float64(tst, "reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.87)
	chk.Float64(tst, "reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.33)
	data.X.Set(6, 0, 0.88)
	data.Y[19] = 87.34
	data.NotifyUpdate()
	chk.Float64(tst, "notified: reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.88)
	chk.Float64(tst, "notified: reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.34)

	// check gradient: dCdθ
	io.Pl()
	verb := chk.Verbose
	tol, hsmall := 1e-8, 1e-3
	params.Backup()
	dCdθ := la.NewVector(data.Nfeatures)
	for _, θ0 := range []float64{5, 10, 15} {

		// analytical
		params.Restore(false)
		params.SetTheta(0, θ0)
		reg.Gradients(dCdθ)

		// numerical
		θat := params.GetThetas()
		θat[0] = θ0
		chk.DerivScaVec(tst, "dCdθ_", tol, dCdθ, θat, hsmall, verb, func(θtmp []float64) (cost float64) {
			params.Restore(false)
			params.SetThetas(θtmp)
			cost = reg.Cost()
			return
		})
	}

	// check gradient: dCdb
	io.Pl()
	for _, b := range []float64{35, 70, 140} {

		// analytical
		params.Restore(false)
		params.SetBias(b)
		dCdb := reg.Gradients(dCdθ)

		// numerical
		chk.DerivScaSca(tst, "dCdb", tol, dCdb, b, hsmall, verb, func(btmp float64) (cost float64) {
			params.Restore(false)
			params.SetBias(btmp)
			cost = reg.Cost()
			return
		})
	}
}

func TestLinReg01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinReg01b. Basic functionality (with regularizaton).")

	// data and parameters
	data := NewDataGivenRawXY(dataReg01)
	params := NewParamsReg(data.Nfeatures)

	// set regularization parameter
	params.SetLambda(10.0)

	// regression
	reg := NewLinReg(data, params, "reg01")

	// check gradient: dCdθ
	io.Pl()
	verb := chk.Verbose
	tol, hsmall := 1e-8, 1e-3
	params.Backup()
	dCdθ := la.NewVector(data.Nfeatures)
	for _, θ0 := range []float64{5, 10, 15} {

		// analytical
		params.Restore(false)
		params.SetTheta(0, θ0)
		reg.Gradients(dCdθ)

		// numerical
		θat := params.GetThetas()
		θat[0] = θ0
		chk.DerivScaVec(tst, "dCdθ_", tol, dCdθ, θat, hsmall, verb, func(θtmp []float64) (cost float64) {
			params.Restore(false)
			params.SetThetas(θtmp)
			cost = reg.Cost()
			return
		})
	}

	// check gradient: dCdb
	io.Pl()
	for _, b := range []float64{35, 70, 140} {

		// analytical
		params.Restore(false)
		params.SetBias(b)
		dCdb := reg.Gradients(dCdθ)

		// numerical
		chk.DerivScaSca(tst, "dCdb", tol, dCdb, b, hsmall, verb, func(btmp float64) (cost float64) {
			params.Restore(false)
			params.SetBias(btmp)
			cost = reg.Cost()
			return
		})
	}
}

func TestLinReg02a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinReg02a. Train simple problem (analytical solution).")

	// data and parameters
	data := NewDataGivenRawXY(dataReg01)
	params := NewParamsReg(data.Nfeatures)

	// regression
	reg := NewLinReg(data, params, "reg01")

	// train
	reg.Train()
	chk.Float64(tst, "cost", 1e-15, reg.Cost(), 5.312454218805082e-01)
	chk.Array(tst, "θ", 1e-12, params.AccessThetas(), []float64{1.494747973211108e+01})
	chk.Float64(tst, "b", 1e-12, params.GetBias(), 7.428331424039514e+01)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		pp := NewPlotterReg(data, params, reg, nil)
		pp.DataY(0)
		pp.ModelY(0, 0.8, 1.6)
		plt.Subplot(2, 1, 2)
		pp.ContourCost(-1, 0, 0, 100, 0, 70)
		plt.Save("/tmp/gosl/ml", "linreg02a")
	}
}

func TestLinReg02b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinReg02b. Train simple problem (analytical solution). With λ.")

	// data and parameters
	data := NewDataGivenRawXY(dataReg01)
	params := NewParamsReg(data.Nfeatures)

	// set regularization parameter
	params.SetLambda(1e12) // very high bias => constant line

	// regression
	reg := NewLinReg(data, params, "reg01")

	// train
	reg.Train()
	for _, x0 := range []float64{0.8, 1.2, 2.0} {
		chk.Float64(tst, io.Sf("y(x0=%.2f)", x0), 1e-11, reg.Predict([]float64{x0}), reg.stat.MeanY)
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		pp := NewPlotterReg(data, params, reg, nil)
		pp.DataY(0)
		pp.ModelY(0, 0.8, 1.6)
		plt.Subplot(2, 1, 2)
		pp.ContourCost(-1, 0, 0, 100, 0, 70)
		plt.Save("/tmp/gosl/ml", "linreg02b")
	}
}
