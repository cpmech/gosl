// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/ml"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// Logistic Regression Test # 2 from Prof. Andrew Ng's online course

	// data (mapped)
	XYraw := io.ReadMatrix("../ml/samples/angEx2data2.txt")
	nOriFeatures := len(XYraw[0]) - 1 // original number of features
	iFeature := 0                     // first index of feature to be generate polynomial
	jFeature := 1                     // second index of feature to be generate polynomial
	degree := 6                       // degreee of polynomial
	mapper := ml.NewPolyDataMapper(nOriFeatures, iFeature, jFeature, degree)
	data := mapper.GetMapped(XYraw, true)

	// parameters and initial guess
	θini := utl.Vals(data.Nfeatures, 1.0) // all ones
	bini := 1.0
	params := ml.NewParamsReg(data.Nfeatures)
	params.SetThetas(θini)
	params.SetBias(bini)
	params.SetLambda(1.0) // regularization

	// model
	model := ml.NewLogReg(data, params, "reg01")

	// train using analytical solution
	params.SetThetas(utl.Vals(data.Nfeatures, 0.0)) // all zeros
	params.SetBias(0.0)
	model.Train()

	// plot data and model prediction (analytical)
	plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
	plt.Subplot(2, 1, 1)
	pp := ml.NewPlotterReg(data, params, model, mapper)
	pp.DataClass(0, 1, true)
	pp.ContourModel(0, 1, 0.5, -1.0, 1.1, -1.0, 1.1)

	// train using gradient-descent
	maxNit := 10
	params.SetThetas(θini)
	params.SetBias(bini)
	gdesc := ml.NewGraDescReg(maxNit)
	gdesc.Alpha = 5.0
	gdesc.Train(data, params, model)

	// plot gradient-descent convergence graph
	plt.Subplot(2, 1, 2)
	gdesc.Plot(nil)
	plt.Save("/tmp/gosl", "ml_ang02")
}
