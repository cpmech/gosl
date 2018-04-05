// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// PlotRegMultiClass plots multi-class classification using Logistic Regression
func PlotRegMultiClass(data *Data, model *LogRegMulti, iFeature, jFeature int, ximin, ximax, xjmin, xjmax float64, npts int) {

	// set classes
	classes := make([]int, data.Nsamples)
	for i := 0; i < data.Nsamples; i++ {
		classes[i] = int(data.Y[i])
	}

	// plot data
	pc := NewPlotterClass(data, classes, model.nClass)
	pc.ArgsYclasses[0].M = "o"
	pc.Data(iFeature, jFeature, false)

	// set colors
	colors := make([]string, model.nClass+1)
	for k := 0; k < model.nClass; k++ {
		colors[k] = pc.ArgsYclasses[k].C
	}
	colors[model.nClass] = "white"

	// plot prediction contour
	var wInt int
	x := la.NewVector(data.Nfeatures)
	U, V, W := utl.MeshGrid2dF(ximin, ximax, xjmin, xjmax, npts, npts, func(s, t float64) (w float64) {
		x[iFeature] = s
		x[jFeature] = t
		wInt, _ = model.Predict(x)
		w = float64(wInt)
		return
	})
	levels := utl.LinSpace(0, float64(model.nClass-1), model.nClass+1)
	plt.ContourF(U, V, W, &plt.A{Colors: colors, Levels: levels, NoLines: true, NoLabels: true})
}

// PlotRegMultiClassOneVsAll plots multi-class classification using Logistic Regression
func PlotRegMultiClassOneVsAll(data *Data, model *LogRegMulti, iFeature, jFeature int, ximin, ximax, xjmin, xjmax float64, npts int) {

	// set classes
	classes := make([]int, data.Nsamples)
	for i := 0; i < data.Nsamples; i++ {
		classes[i] = int(data.Y[i])
	}

	// plot data
	pc := NewPlotterClass(data, classes, model.nClass)
	pc.ArgsYclasses[0].M = "o"
	pc.Data(iFeature, jFeature, false)

	// plot contour
	for k, m := range model.models {
		pp := NewPlotterReg(model.dataB[k], m, nil)
		pp.MgridNpts = npts
		pp.ArgsCmodel.Colors = []string{pc.ArgsYclasses[k].C}
		pp.ContourModel(iFeature, jFeature, 0.5, ximin, ximax, xjmin, xjmax)
	}
}
