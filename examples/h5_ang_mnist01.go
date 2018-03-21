// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io/h5"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/ml/imgd"
	"github.com/cpmech/gosl/rnd"
)

func main() {

	// NOTE: this example expects an environment variable called
	//       $GOSLDATA containing all Gosl data files

	// read data file
	f := h5.Open("$GOSLDATA", "angEx4data1", false)
	defer f.Close()
	Xraw := f.GetArray("/Xcolmaj/value")
	nSamples := f.GetInt("/m/value")
	sampleSize := f.GetInt("/n/value")
	X := la.NewMatrixRaw(nSamples, sampleSize, Xraw)

	// constants
	nSelected := 100 // number of samples to display
	padding := 1     // padding
	rowMaj := false  // row major
	random := true   // random selection

	// select samples
	rnd.Init(0)
	samples := imgd.NewGraySamples(X, nSelected, rowMaj, random)
	smin, smax, wmax, hmax := samples.Stat()

	// board
	board := imgd.NewGrayBoard(nSelected, wmax, hmax, padding)
	board.Paint(samples, smin, smax, false)
	board.SavePng("/tmp/gosl", "angEx4data1fig")
}
