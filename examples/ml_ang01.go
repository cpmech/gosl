// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/ml"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// Logistic Regression Test # 1 from Prof. Andrew Ng's online course

	// load data
	XYraw := io.ReadMatrix("../ml/samples/angEx2data1.txt")
	data := ml.NewDataGivenRawXY(XYraw)

	// parameters and initial guess
	θini := []float64{0.2, 0.2}
	bini := -24.0
	params := ml.NewParamsReg(data.Nfeatures)
	params.SetThetas(θini)
	params.SetBias(bini)

	// model
	model := ml.NewLogReg(data, params)

	// train using analytical solution
	model.Train()

	// plot data and model prediction (analytical)
	plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 0.8})
	pp := ml.NewPlotterReg(data, params, model, nil)
	pp.DataClass(0, 1, true)
	pp.ContourModel(0, 1, 0.5, 20, 100, 20, 100)
	plt.Save("/tmp/gosl", "ml_ang01")
}
