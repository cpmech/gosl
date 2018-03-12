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

func TestLogReg01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg01a. Basic functionality (no regularizaton).")

	// data and parameters
	data := NewDataGivenRawXY(dataReg01)
	params := NewParamsReg(data.Nfeatures)

	// set regularization
	params.SetLambda(0.25)

	// regression
	reg := NewLogReg(data, params, "reg01")

	// check stat
	chk.Float64(tst, "reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.87)
	chk.Float64(tst, "reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.33)
	data.X.Set(6, 0, 0.88)
	data.Y[19] = 87.34
	data.NotifyUpdate()
	chk.Float64(tst, "notified: reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.88)
	chk.Float64(tst, "notified: reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.34)

	// meshgrid
	thetas := utl.LinSpace(-100, 100, 11)
	biass := utl.LinSpace(-100, 100, 11)

	// check gradient: dCdθ and dCdb
	io.Pl()
	verb := chk.Verbose
	tol, hsmall := 1e-7, 1e-3
	θat := la.NewVector(data.Nfeatures)
	dCdθ := la.NewVector(data.Nfeatures)
	for _, θ0 := range thetas {
		for _, b := range biass {

			// analytical
			io.Pf("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> θ0=%v b=%v <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", θ0, b)
			params.SetTheta(0, θ0)
			params.SetBias(b)
			params.Backup()
			dCdb := reg.Gradients(dCdθ)

			// numerical
			θat[0] = params.GetTheta(0)
			chk.DerivScaVec(tst, "dCdθ_", tol, dCdθ, θat, hsmall, verb, func(θtmp []float64) (cost float64) {
				params.Restore(false)
				params.SetThetas(θtmp)
				cost = reg.Cost()
				return
			})

			// numerical
			chk.DerivScaSca(tst, "dCdb  ", tol, dCdb, b, hsmall, verb, func(btmp float64) (cost float64) {
				params.Restore(false)
				params.SetBias(btmp)
				cost = reg.Cost()
				return
			})
		}
	}

	// check Hessian
	tol2 := 1e-8
	io.Pl()
	var w float64
	d, v, D, H := reg.AllocateHessian()
	for _, θ0 := range thetas {
		for _, b := range biass {

			// analytical
			io.Pf("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> θ0=%v b=%v <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", θ0, b)
			params.SetTheta(0, θ0)
			params.SetBias(b)
			params.Backup()
			w = reg.Hessian(d, v, D, H)

			// numerical
			θat[0] = params.GetTheta(0)
			chk.DerivVecVec(tst, "∂²C/∂θ∂θ_", tol2, H.GetDeep2(), θat, hsmall, verb, func(dCdθtmp, θtmp []float64) {
				params.Restore(false)
				params.SetThetas(θtmp)
				reg.Gradients(dCdθtmp)
				return
			})

			// numerical
			chk.DerivVecSca(tst, "∂²C/∂θ∂b_ ", tol2, v, b, hsmall, verb, func(dCdθtmp []float64, btmp float64) {
				params.Restore(false)
				params.SetBias(btmp)
				reg.Gradients(dCdθtmp)
				return
			})

			// numerical
			chk.DerivScaSca(tst, "∂²C/∂b∂b   ", tol2, w, b, hsmall, verb, func(btmp float64) (dCdbtmp float64) {
				params.Restore(false)
				params.SetBias(btmp)
				dCdbtmp = reg.Gradients(dCdθ)
				return
			})
		}
	}
}

func TestLogReg02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg02. simple -45° grid")

	// data
	nr, nc := 5, 5
	data := NewData(nr*nc, 2, true)
	k := 0
	for j := 0; j < nc; j++ {
		for i := 0; i < nr; i++ {
			x0 := -1 + 2*float64(i)/float64(nc-1)
			x1 := -1 + 2*float64(j)/float64(nr-1)
			data.X.Set(k, 0, x0)
			data.X.Set(k, 1, x1)
			if x0+x1+1e-15 >= 0 {
				data.Y[k] = 1
			}
			k++
		}
	}

	// parameters
	params := NewParamsReg(data.Nfeatures)
	//params.SetThetas([]float64{20, 30})
	//params.SetBias(10)

	// regression
	reg := NewLogReg(data, params, "reg01")
	reg.Train()
	io.Pforan("cost = %v\n", reg.Cost())
	io.Pforan("θ = %v\n", params.AccessThetas())
	io.Pforan("b = %v\n", params.GetBias())
	chk.Float64(tst, "cost", 1e-15, reg.Cost(), 0.0007850399226816407)
	chk.Array(tst, "θ", 1e-14, params.AccessThetas(), []float64{24.488302802315026, 24.48830280231502})
	chk.Float64(tst, "b", 1e-14, params.GetBias(), 6.183574567556589)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		pp := NewPlotterReg(data, params, reg, nil)
		pp.DataClass(0, 1, true)
		pp.ContourModel(0, 1, 0.5, -1, 1, -1, 1)
	}

	// gradient-descent
	params.SetThetas([]float64{0, 0})
	params.SetBias(0)
	maxNit := 10
	gdesc := NewGraDescReg(maxNit)
	gdesc.Alpha = 100
	gdesc.Train(data, params, reg)
	io.Pfblue2("cost = %v\n", reg.Cost())
	io.Pfblue2("θ = %v\n", params.AccessThetas())
	io.Pfblue2("b = %v\n", params.GetBias())
	chk.Float64(tst, "cost", 1e-15, reg.Cost(), 0.0015372029816003163)
	chk.Array(tst, "θ", 1e-14, params.AccessThetas(), []float64{22.06214330726067, 22.06214330726067})
	chk.Float64(tst, "b", 1e-14, params.GetBias(), 5.254524501188747)

	// plot
	if chk.Verbose {
		plt.Subplot(2, 1, 2)
		gdesc.Plot(nil)
		plt.Save("/tmp/gosl/ml", "logreg02")
	}
}

func TestLogReg03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg03. ANG test # 1")

	// data
	XYraw := io.ReadMatrix("./samples/angEx2data1.txt")
	data := NewDataGivenRawXY(XYraw)
	chk.Int(tst, "nSamples", data.Nsamples, 100)
	chk.Int(tst, "nFeatures", data.Nfeatures, 2)

	// parameters and initial guess
	θini := []float64{0.2, 0.2}
	bini := -24.0
	params := NewParamsReg(data.Nfeatures)
	params.SetThetas(θini)
	params.SetBias(bini)

	// model
	reg := NewLogReg(data, params, "reg01")
	cost := reg.Cost()
	io.Pf("Initial: θ = %.8f\n", params.GetThetas())
	io.Pf("Initial: b = %.8f\n", params.GetBias())
	io.Pf("Initial: cost = %.8f\n", cost)
	chk.Float64(tst, "\ncostIni", 1e-15, reg.Cost(), 2.183301938265978e-01)

	// train using analytical solution
	reg.Train()
	chk.Float64(tst, "\ncost", 1e-15, reg.Cost(), 2.034977015894404e-01)
	chk.Array(tst, "θ", 1e-8, params.AccessThetas(), []float64{2.062317052577260e-01, 2.014715922708144e-01})
	chk.Float64(tst, "b", 1e-6, params.GetBias(), -2.516133256589910e+01)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		pp := NewPlotterReg(data, params, reg, nil)
		pp.DataClass(0, 1, true)
		pp.ContourModel(0, 1, 0.5, 20, 100, 20, 100)
	}

	// train using gradient-descent
	maxNit := 10
	params.SetThetas(θini)
	params.SetBias(bini)
	gdesc := NewGraDescReg(maxNit)
	gdesc.Alpha = 0.002
	gdesc.Train(data, params, reg)
	chk.Float64(tst, "\ncost", 1e-15, reg.Cost(), 2.037591668976244e-01)
	chk.Array(tst, "θ", 1e-8, params.AccessThetas(), []float64{1.957478716620902e-01, 1.933175159514175e-01})
	chk.Float64(tst, "b", 1e-6, params.GetBias(), -2.400009669708430e+01)

	// plot
	if chk.Verbose {
		plt.Subplot(2, 1, 2)
		gdesc.Plot(nil)
		plt.Save("/tmp/gosl/ml", "logreg03")
	}
}

func TestLogReg04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg04. ANG test # 2")

	// data (mapped)
	XYraw := io.ReadMatrix("./samples/angEx2data2.txt")
	nOriFeatures := len(XYraw[0]) - 1
	iFeature := 0
	jFeature := 1
	degree := 6
	mapper := NewPolyDataMapper(nOriFeatures, iFeature, jFeature, degree)
	data := mapper.GetMapped(XYraw, true)
	chk.Int(tst, "nSamples", data.Nsamples, 118)
	chk.Int(tst, "nFeatures", data.Nfeatures, 27)

	// parameters and initial guess
	θini := utl.Vals(data.Nfeatures, 1.0) // all ones
	bini := 1.0
	params := NewParamsReg(data.Nfeatures)
	params.SetThetas(θini)
	params.SetBias(bini)
	params.SetLambda(1.0) // regularization

	// model
	reg := NewLogReg(data, params, "reg01")
	cost := reg.Cost()
	io.Pf("Initial: θ = %.3f\n", params.GetThetas()[:4])
	io.Pf("Initial: b = %.8f\n", params.GetBias())
	io.Pf("Initial: cost = %.8f\n", cost)
	chk.Float64(tst, "\ncostIni", 1e-15, reg.Cost(), 2.134848314666066)

	// train using analytical solution
	params.SetThetas(utl.Vals(data.Nfeatures, 0.0)) // all zeros
	params.SetBias(0.0)
	reg.Train()
	chk.Float64(tst, "\ncost", 1e-15, reg.Cost(), 5.290027411158117e-01)
	chk.Array(tst, "θ", 1e-14, params.AccessThetas()[:4], []float64{6.252526148274546e-01, 1.180976145721166, -2.019842398401904, -9.173659359499787e-01})
	chk.Float64(tst, "b", 1e-14, params.GetBias(), 1.272656700281225)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		pp := NewPlotterReg(data, params, reg, mapper)
		pp.DataClass(0, 1, true)
		pp.ContourModel(0, 1, 0.5, -1.0, 1.1, -1.0, 1.1)
	}

	// train using gradient-descent
	maxNit := 10
	params.SetThetas(θini)
	params.SetBias(bini)
	gdesc := NewGraDescReg(maxNit)
	gdesc.Alpha = 5.0
	gdesc.Train(data, params, reg)
	chk.Float64(tst, "\ncost", 1e-15, reg.Cost(), 5.920108560779025e-01)
	chk.Array(tst, "θ", 1e-15, params.AccessThetas()[:4], []float64{-1.730594903181217e-01, 3.615618466861891e-01, -1.194645899263627e+00, -4.186288373383852e-01})
	chk.Float64(tst, "b", 1e-15, params.GetBias(), 6.527575848138054e-01)

	// plot
	if chk.Verbose {
		plt.Subplot(2, 1, 2)
		gdesc.Plot(nil)
		plt.Save("/tmp/gosl/ml", "logreg04")
	}
}
