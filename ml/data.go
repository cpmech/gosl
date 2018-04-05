// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ml implements Machine Learning algorithms
package ml

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// Data holds data in matrix format; e.g. for regression computations
//
//   Example:
//          _          _                                     _   _
//         |  -1  0 -3  |                                   |  0  |
//         |  -2  3  3  |                       (optional)  |  1  |
//     X = |   3  1  4  |                               Y = |  1  |
//         |  -4  5  0  |                                   |  0  |
//         |_  1 -8  5 _|(nSamples x nFeatures)             |_ 1 _|(nSamples)
//
//   NOTE: remember to call data.NotifyUpdate() after changing X or y components
//
type Data struct {
	utl.Observable // can notify others of changes here via data.NotifyUpdate()

	// input
	Nsamples  int        // number of data points (samples). number of rows in X and Y
	Nfeatures int        // number of features. number of columns in X
	X         *la.Matrix // [nSamples][nFeatures] X values
	Y         la.Vector  // [nSamples] Y values [optional]

	// access
	Stat *Stat // statistics about this data
}

// NewData returns a new object to hold ML data
//  Input:
//    nSamples  -- number of data samples (rows in X)
//    nFeatures -- number of features (columns in X)
//    useY      -- use y data vector
//    allocate  -- allocates X (and Y); otherwise,
//                 X and Y must be set using Set() method
//  Output:
//    new object
func NewData(nSamples, nFeatures int, useY, allocate bool) (o *Data) {
	o = new(Data)
	o.Nsamples = nSamples
	o.Nfeatures = nFeatures
	if allocate {
		o.X = la.NewMatrix(o.Nsamples, o.Nfeatures)
		if useY {
			o.Y = la.NewVector(o.Nsamples)
		}
	}
	o.Stat = NewStat(o)
	return
}

// Set sets X matrix and Y vector [optional] and notify observers
//   Input:
//     X -- x values
//     Y -- y values [optional]
func (o *Data) Set(X *la.Matrix, Y la.Vector) {
	o.X = X
	o.Y = Y
	o.NotifyUpdate()
}

// NewDataGivenRawX returns a new object with data set from raw X values
//  Input:
//    Xraw -- [nSamples][nFeatures] table with x values (NO y values)
//  Output:
//    new object
func NewDataGivenRawX(Xraw [][]float64) (o *Data) {

	// check
	nSamples := len(Xraw)
	if nSamples < 1 {
		chk.Panic("at least one row of data in table must be provided\n")
	}

	// allocate new object
	nFeatures := len(Xraw[0])
	o = NewData(nSamples, nFeatures, true, true)

	// copy data from raw table to X matrix
	for i := 0; i < nSamples; i++ {
		for j := 0; j < nFeatures; j++ {
			o.X.Set(i, j, Xraw[i][j])
		}
	}

	// stat
	o.Stat = NewStat(o)
	o.NotifyUpdate()
	return
}

// NewDataGivenRawXY returns a new object with data set from raw XY values
//  Input:
//    XYraw -- [nSamples][nFeatures+1] table with x and y raw values,
//             where the last column contains y-values
//  Output:
//    new object
func NewDataGivenRawXY(XYraw [][]float64) (o *Data) {

	// check
	nSamples := len(XYraw)
	if nSamples < 1 {
		chk.Panic("at least one row of data in table must be provided\n")
	}

	// allocate new object
	nFeatures := len(XYraw[0]) - 1 // -1 because of y column
	o = NewData(nSamples, nFeatures, true, true)

	// copy data from raw table to X and Y arrays
	for i := 0; i < nSamples; i++ {
		for j := 0; j < nFeatures; j++ {
			o.X.Set(i, j, XYraw[i][j])
		}
		o.Y[i] = XYraw[i][nFeatures]
	}

	// stat
	o.Stat = NewStat(o)
	o.NotifyUpdate()
	return
}

// GetCopy returns a deep copy of this object
func (o *Data) GetCopy() (p *Data) {
	useY := len(o.Y) > 0
	p = NewData(o.Nsamples, o.Nfeatures, useY, true)
	o.X.CopyInto(p.X, 1)
	if useY {
		copy(p.Y, o.Y)
	}
	o.Stat.CopyInto(p.Stat)
	return
}
