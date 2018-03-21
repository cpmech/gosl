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

	// Multi-class classification using logistic regression one-vs-all approach
	// Reference: http://scikit-learn.org/stable/auto_examples/linear_model/plot_iris_logistic.html

	// data
	XYraw := io.ReadMatrix("../ml/samples/multiclass01.txt")
	data := ml.NewDataGivenRawXY(XYraw)

	// model
	model := ml.NewLogRegMultiClass(data, "model01")

	// train
	gradDesc := false
	model.SetLambda(1e-5)
	model.Train(gradDesc)

	// plot
	npts := 201
	iFeature, jFeature := 0, 1
	ximin, ximax, xjmin, xjmax := 3.8, 8.4, 1.5, 4.9
	plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})
	plt.Subplot(2, 1, 1)
	ml.PlotRegMultiClass(data, model, iFeature, jFeature, ximin, ximax, xjmin, xjmax, npts)
	plt.Subplot(2, 1, 2)
	ml.PlotRegMultiClassOneVsAll(data, model, iFeature, jFeature, ximin, ximax, xjmin, xjmax, npts)
	plt.Save("/tmp/gosl", "ml_mclass01")
}
