// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
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
	model := NewKmeans(data, nClasses)
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
		pp := NewPlotter(data, nil)
		//model.bins.Draw(true, true, true, false, nil, argsGrid, argsTxtEntry, nil, nil)
		pp.DataClass(model.nClasses, 0, 1, model.Classes)
		pp.Centroids(model.Centroids)
		plt.AxisRange(0, 1, 0, 1)
		plt.Save("/tmp/gosl/ml", "kmeans01")
	}
}
