// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/ml"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// Simple Linear Regression Test

	// data
	XYraw := [][]float64{
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
	}
	data := ml.NewDataGivenRawXY(XYraw)

	// parameters
	params := ml.NewParamsReg(data.Nfeatures)

	// model
	model := ml.NewLinReg(data, params, "reg01")

	// train using analytical solution
	model.Train()

	// ----------------------- plotting --------------------------

	// clear plotting area
	plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})

	// plot data x-y
	plt.Subplot(2, 1, 1)
	pp := ml.NewPlotterReg(data, params, model, nil)
	pp.DataY(0)

	// plot model x-y
	pp.ModelY(0, 0.8, 1.6)

	// plot contour of cost function
	plt.Subplot(2, 1, 2)
	pp.ContourCost(-1, 0, 0, 100, 0, 70)

	// save figure
	plt.Save("/tmp/gosl", "ml_simple01")
}
