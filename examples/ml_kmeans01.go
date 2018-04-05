// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/io/h5"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/ml"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// K-means clustering. Test # 1 from Prof. Andrew Ng's online course

	// load raw data from HDF5 file
	f := h5.Open("../ml/samples", "angEx7data2", false)
	defer f.Close()
	Xraw := f.GetArray("/Xcolmaj/value")
	nSamples := f.GetInt("/m/value")
	nColumns := f.GetInt("/n/value")

	// data
	useY := false
	allocate := false
	data := ml.NewData(nSamples, nColumns, useY, allocate)
	data.X = la.NewMatrixRaw(nSamples, nColumns, Xraw)

	// model
	nClasses := 3
	model := ml.NewKmeans(data, nClasses)
	model.SetCentroids([][]float64{
		{3, 3}, // class 0
		{6, 2}, // class 1
		{8, 5}, // class 2
	})

	// check initial classes
	model.FindClosestCentroids()

	// check initial computation of centroids
	model.ComputeCentroids()
	io.Pf("number of members in each cluster = %v\n", model.Nmembers)

	// train
	model.Train(6, 0)

	// plot
	plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
	pp := ml.NewPlotterClass(data, model.Classes, model.Nclasses())
	pp.Data(0, 1, false)
	pp.Centroids(model.Centroids)
	plt.Equal()
	plt.Save("/tmp/gosl", "ml_kmeans01")
}
