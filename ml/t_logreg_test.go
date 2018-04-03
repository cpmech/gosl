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
	model := NewLogReg(data, params, "reg01")

	// check stat
	chk.Float64(tst, "reg.stat.min(x)", 1e-15, model.stat.MinX[0], 0.87)
	chk.Float64(tst, "reg.stat.min(y)", 1e-15, model.stat.MinY, 87.33)
	data.X.Set(6, 0, 0.88)
	data.Y[19] = 87.34
	data.NotifyUpdate()
	chk.Float64(tst, "notified: reg.stat.min(x)", 1e-15, model.stat.MinX[0], 0.88)
	chk.Float64(tst, "notified: reg.stat.min(y)", 1e-15, model.stat.MinY, 87.34)

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
			dCdb := model.Gradients(dCdθ)

			// numerical
			θat[0] = params.GetTheta(0)
			chk.DerivScaVec(tst, "dCdθ_", tol, dCdθ, θat, hsmall, verb, func(θtmp []float64) (cost float64) {
				params.Restore(false)
				params.SetThetas(θtmp)
				cost = model.Cost()
				return
			})

			// numerical
			chk.DerivScaSca(tst, "dCdb  ", tol, dCdb, b, hsmall, verb, func(btmp float64) (cost float64) {
				params.Restore(false)
				params.SetBias(btmp)
				cost = model.Cost()
				return
			})
		}
	}

	// check Hessian
	tol2 := 1e-8
	io.Pl()
	var w float64
	d, v, D, H := model.AllocateHessian()
	for _, θ0 := range thetas {
		for _, b := range biass {

			// analytical
			io.Pf("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> θ0=%v b=%v <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", θ0, b)
			params.SetTheta(0, θ0)
			params.SetBias(b)
			params.Backup()
			w = model.Hessian(d, v, D, H)

			// numerical
			θat[0] = params.GetTheta(0)
			chk.DerivVecVec(tst, "∂²C/∂θ∂θ_", tol2, H.GetDeep2(), θat, hsmall, verb, func(dCdθtmp, θtmp []float64) {
				params.Restore(false)
				params.SetThetas(θtmp)
				model.Gradients(dCdθtmp)
				return
			})

			// numerical
			chk.DerivVecSca(tst, "∂²C/∂θ∂b_ ", tol2, v, b, hsmall, verb, func(dCdθtmp []float64, btmp float64) {
				params.Restore(false)
				params.SetBias(btmp)
				model.Gradients(dCdθtmp)
				return
			})

			// numerical
			chk.DerivScaSca(tst, "∂²C/∂b∂b   ", tol2, w, b, hsmall, verb, func(btmp float64) (dCdbtmp float64) {
				params.Restore(false)
				params.SetBias(btmp)
				dCdbtmp = model.Gradients(dCdθ)
				return
			})
		}
	}
}

func TestLogReg01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg01b. Basic functionality (with regularizaton).")

	// data
	data := NewDataGivenRawXY([][]float64{
		{0.1, 0.60, 1.1, 1}, // X and y values
		{0.2, 0.70, 1.2, 0},
		{0.3, 0.80, 1.3, 1},
		{0.4, 0.90, 1.4, 0},
		{0.5, 1.00, 1.5, 1},
	})

	// parameters
	params := NewParamsReg(data.Nfeatures)
	params.SetThetas([]float64{-1, 1, 2})
	params.SetBias(-2)
	params.SetLambda(3)

	// regression
	model := NewLogReg(data, params, "reg01")

	// check
	dCdθ := model.AllocateGradient()
	dCdb := model.Gradients(dCdθ)
	chk.Float64(tst, "cost", 1e-15, model.Cost(), 2.534819396109744)
	chk.Float64(tst, "dCdb", 1e-15, dCdb, 1.465613679248980e-01)
	chk.Array(tst, "grads", 1e-15, dCdθ, []float64{-5.485584118531603e-01, 7.247222721092885e-01, 1.398002956071738})
}

func TestLogReg02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg02. simple -45° grid")

	// data
	nr, nc := 5, 5
	useY := true
	allocate := true
	data := NewData(nr*nc, 2, useY, allocate)
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
	model := NewLogReg(data, params, "reg01")
	model.Train()
	io.Pforan("cost = %v\n", model.Cost())
	io.Pforan("θ = %v\n", params.AccessThetas())
	io.Pforan("b = %v\n", params.GetBias())
	chk.Float64(tst, "cost", 1e-14, model.Cost(), 0.0007850399226816407)
	chk.Array(tst, "θ", 1e-13, params.AccessThetas(), []float64{24.488302802315026, 24.48830280231502})
	chk.Float64(tst, "b", 1e-14, params.GetBias(), 6.183574567556589)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 0.8})
		pp := NewPlotterReg(data, params, model, nil)
		pp.DataClass(0, 1, true)
		pp.ContourModel(0, 1, 0.5, -1, 1, -1, 1)
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
	model := NewLogReg(data, params, "reg01")
	cost := model.Cost()
	io.Pf("Initial: θ = %.8f\n", params.GetThetas())
	io.Pf("Initial: b = %.8f\n", params.GetBias())
	io.Pf("Initial: cost = %.8f\n", cost)
	chk.Float64(tst, "\ncostIni", 1e-15, model.Cost(), 2.183301938265978e-01)

	// train using analytical solution
	model.Train()
	chk.Float64(tst, "\ncost", 1e-14, model.Cost(), 2.034977015894404e-01)
	chk.Array(tst, "θ", 1e-8, params.AccessThetas(), []float64{2.062317052577260e-01, 2.014715922708144e-01})
	chk.Float64(tst, "b", 1e-6, params.GetBias(), -2.516133256589910e+01)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 0.8})
		pp := NewPlotterReg(data, params, model, nil)
		pp.DataClass(0, 1, true)
		pp.ContourModel(0, 1, 0.5, 20, 100, 20, 100)
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
	model := NewLogReg(data, params, "reg01")
	cost := model.Cost()
	io.Pf("Initial: θ = %.3f\n", params.GetThetas()[:4])
	io.Pf("Initial: b = %.8f\n", params.GetBias())
	io.Pf("Initial: cost = %.8f\n", cost)
	chk.Float64(tst, "\ncostIni", 1e-15, model.Cost(), 2.134848314666066)

	// train using analytical solution
	params.SetThetas(utl.Vals(data.Nfeatures, 0.0)) // all zeros
	params.SetBias(0.0)
	model.Train()
	chk.Float64(tst, "\ncost", 1e-15, model.Cost(), 5.290027411158117e-01)
	chk.Array(tst, "θ", 1e-14, params.AccessThetas()[:4], []float64{6.252526148274546e-01, 1.180976145721166, -2.019842398401904, -9.173659359499787e-01})
	chk.Float64(tst, "b", 1e-14, params.GetBias(), 1.272656700281225)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 0.8})
		pp := NewPlotterReg(data, params, model, mapper)
		pp.DataClass(0, 1, true)
		pp.ContourModel(0, 1, 0.5, -1.0, 1.1, -1.0, 1.1)
		plt.Save("/tmp/gosl/ml", "logreg04")
	}
}

func TestLogReg05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg05. 3-class problem")

	// data
	XYraw := io.ReadMatrix("./samples/multiclass01.txt")
	data := NewDataGivenRawXY(XYraw)
	chk.Int(tst, "nSamples", data.Nsamples, 150)
	chk.Int(tst, "nFeatures", data.Nfeatures, 2)

	// model
	model := NewLogRegMulti(data, "model01")

	// train
	model.SetLambda(1e-5)
	model.Train()

	// check
	classes := make([]int, data.Nsamples)
	fails := 0
	for i := 0; i < data.Nsamples; i++ {
		x := data.X.GetRow(i)
		class, _ := model.Predict(x)
		classes[i] = class
		if class != int(data.Y[i]) {
			fails++
		}
	}
	chk.Ints(tst, "prediction", classes, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 1, 2, 1, 2, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 2, 1, 1, 1, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 2, 2, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 2, 1, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2, 1, 2, 1, 2, 1, 2, 2, 1, 2, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 1, 2, 2, 2, 1, 2, 2, 2, 1, 2, 2, 1})
	chk.Int(tst, "fails", fails, 29)

	// plot
	if chk.Verbose {
		npts := 201
		iFeature, jFeature := 0, 1
		ximin, ximax, xjmin, xjmax := 3.8, 8.4, 1.5, 4.9
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})
		plt.Subplot(2, 1, 1)
		PlotRegMultiClass(data, model, iFeature, jFeature, ximin, ximax, xjmin, xjmax, npts)
		plt.Subplot(2, 1, 2)
		PlotRegMultiClassOneVsAll(data, model, iFeature, jFeature, ximin, ximax, xjmin, xjmax, npts)
		plt.Save("/tmp/gosl/ml", "logreg05")
	}
}
