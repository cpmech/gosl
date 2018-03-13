// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/io/h5"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestKmeans01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Kmeans01. Basic functionality.")

	// data
	data := NewDataGivenRawX([][]float64{
		{0.1, 0.7}, {0.3, 0.7},
		{0.1, 0.9}, {0.3, 0.9},
		{0.7, 0.1}, {0.9, 0.1},
		{0.7, 0.3}, {0.9, 0.3},
	})

	// model
	nClasses := 2
	model := NewKmeans(data, "model01", nClasses)
	model.SetCentroids([][]float64{
		{0.4, 0.6}, // class 0
		{0.6, 0.4}, // class 1
	})

	// train
	model.FindClosestCentroids()
	chk.Ints(tst, "classes", model.Classes, []int{
		0, 0,
		0, 0,
		1, 1,
		1, 1,
	})

	// plot
	if chk.Verbose {
		//argsGrid := &plt.A{C: "gray", Ls: "--", Lw: 0.7}
		//argsTxtEntry := &plt.A{Fsz: 5}
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		pp := NewPlotterClass(data, model.Classes, model.Nclasses())
		//model.bins.Draw(true, true, true, false, nil, argsGrid, argsTxtEntry, nil, nil)
		pp.Data(0, 1, false)
		pp.Centroids(model.Centroids)
		plt.AxisRange(0, 1, 0, 1)
		plt.Save("/tmp/gosl/ml", "kmeans01")
	}
}

func TestKmeans02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Kmeans02. ANT test # 1")

	// load raw data from HDF5 file
	f := h5.Open("$GOSLDATA", "angEx7data2", false)
	defer f.Close()
	Xraw := f.GetArray("/Xcolmaj/value")
	nSamples := f.GetInt("/m/value")
	nColumns := f.GetInt("/n/value")

	// data
	useY := false
	allocate := false
	data := NewData(nSamples, nColumns, useY, allocate)
	data.X = la.NewMatrixRaw(nSamples, nColumns, Xraw)
	chk.Int(tst, "m", data.X.M, 300)
	chk.Int(tst, "n", data.X.N, 2)

	// model
	nClasses := 3
	model := NewKmeans(data, "model01", nClasses)
	model.SetCentroids([][]float64{
		{3, 3}, // class 0
		{6, 2}, // class 1
		{8, 5}, // class 2
	})

	// check initial classes
	model.FindClosestCentroids()
	chk.Ints(tst, "Classes", model.Classes, []int{0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 2, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 0})

	// check initial computation of centroids
	model.ComputeCentroids()
	io.Pforan("Nmembers = %v\n", model.Nmembers)
	chk.Ints(tst, "Nmembers", model.Nmembers, []int{191, 103, 6})
	chk.Array(tst, "Centroid[0]", 1e-15, model.Centroids[0], []float64{2.428301112098196e+00, 3.157924176603567e+00})
	chk.Array(tst, "Centroid[1]", 1e-15, model.Centroids[1], []float64{5.813503308520713e+00, 2.633656451403025e+00})
	chk.Array(tst, "Centroid[2]", 1e-15, model.Centroids[2], []float64{7.119386871508754e+00, 3.616684398721619e+00})

	// train
	model.Train(6, 0)
	chk.Ints(tst, "Classes", model.Classes, []int{0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		pp := NewPlotterClass(data, model.Classes, model.Nclasses())
		pp.Data(0, 1, false)
		pp.Centroids(model.Centroids)
		plt.Equal()
		plt.Save("/tmp/gosl/ml", "kmeans02")
	}
}
