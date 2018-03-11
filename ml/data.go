// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ml implements Machine Learning algorithms
package ml

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
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
type Data struct {

	// input
	Nsamples  int        // number of data points (samples). number of rows in X and Y
	Nfeatures int        // number of features. number of columns in X
	UseY      bool       // use Y vector
	X         *la.Matrix // [nSamples][nFeatures] X values
	Y         la.Vector  // [nSamples] Y values

	// control
	statOk bool // indicates that Stat is OK; otherwise Stat() must be called because x,y changed

	// stat
	minX, maxX  []float64 // [nFeatures] range of x
	sumX, meanX []float64 // [nFeatures] sum and mean of x
	sigX        []float64 // [nFeatures] standard deviations of x
	delX        []float64 // [nFeatures] difference: maxX - minX
	minY, maxY  float64   // range of y
	sumY, meanY float64   // sum and mean of y
	sigY        float64   // standard deviation of y
	delY        float64   // difference: maxY - minY

	// plotting
	horiz, vert  []float64   // [npts] horizontal and vertical arrays for curves
	thi, thj, cc [][]float64 // [npts][npts] meshgrid over θ
	xxi, xxj, zz [][]float64 // [npts][npts] meshgrid over x
}

// NewDataMatrix returns a new structure to hold ML data
//  Input:
//    nSamples  -- number of data samples (rows in X)
//    nFeatures -- number of features (columsn in X)
//    yData     -- use y data vector
//  Output:
//    new object
func NewDataMatrix(nSamples, nFeatures int, yData bool) (o *Data) {

	// main
	o = new(Data)
	o.Nsamples = nSamples
	o.Nfeatures = nFeatures
	o.UseY = yData
	o.X = la.NewMatrix(o.Nsamples, o.Nfeatures)
	if o.UseY {
		o.Y = la.NewVector(o.Nsamples)
	}

	// stat
	o.minX = make([]float64, nFeatures)
	o.maxX = make([]float64, nFeatures)
	o.sumX = make([]float64, nFeatures)
	o.meanX = make([]float64, nFeatures)
	o.sigX = make([]float64, nFeatures)
	o.delX = make([]float64, nFeatures)
	return
}

// NewDataMatrixTable sets X and y values given table
//  Input:
//    xyRawTable -- [nData][nFeatures+1] table with x and y raw values,
//                  where the last column contains y-values
//  Output:
//    new object
func NewDataMatrixTable(xyRawTable [][]float64) (o *Data) {
	nSamples := len(xyRawTable)
	if nSamples < 1 {
		chk.Panic("at least one row of data in table must be provided\n")
	}
	nFeatures := len(xyRawTable[0]) - 1 // -1 because of y column
	o = NewDataMatrix(nSamples, nFeatures, true)
	for i := 0; i < nSamples; i++ {
		for j := 0; j < nFeatures; j++ {
			//o.SetX(i, j, xyRawTable[i][j])
		}
		//o.SetY(i, xyRawTable[i][nFeatures])
	}
	o.stat()
	return
}

// Normalize normalizes x values
//  Input:
//   useMinMax -- divide by max-min; otherwise by the standard deviation
func (o *Data) Normalize(useMinMax bool) {
	if !o.statOk {
		o.stat()
	}
	den := o.sigX
	if useMinMax {
		den = o.delX
	}
	for i := 0; i < o.Nsamples; i++ {
		for j := 0; j < o.Nfeatures; j++ {
			o.X.Set(i, j, (o.X.Get(i, j)-o.meanX[j])/den[j])
		}
	}
	o.statOk = false // TODO: create copy matrix
}

// auxiliary /////////////////////////////////////////////////////////////////////////////////////

// stat computes statistics
func (o *Data) stat() {

	// x values
	for j := 0; j < o.Nfeatures; j++ {
		o.minX[j] = o.X.Get(0, j)
		o.maxX[j] = o.minX[j]
		o.sumX[j] = 0.0
		for i := 0; i < o.Nsamples; i++ {
			xval := o.X.Get(i, j)
			o.minX[j] = utl.Min(o.minX[j], xval)
			o.maxX[j] = utl.Max(o.maxX[j], xval)
			o.sumX[j] += xval
		}
		o.meanX[j] = o.sumX[j] / float64(o.Nsamples)
		o.sigX[j] = rnd.StatDevFirst(o.X.Col(j), o.meanX[j], true)
		o.delX[j] = o.maxX[j] - o.minX[j]
	}

	// y values
	if o.UseY {
		o.minY = o.Y[0]
		o.maxY = o.minY
		o.sumY = 0.0
		for i := 0; i < o.Nsamples; i++ {
			o.minY = utl.Min(o.minY, o.Y[i])
			o.maxY = utl.Max(o.maxY, o.Y[i])
			o.sumY += o.Y[i]
		}
		o.meanY = o.sumY / float64(o.Nsamples)
		o.sigY = rnd.StatDevFirst(o.Y, o.meanY, true)
		o.delY = o.maxY - o.minY
	}

	// flag
	o.statOk = true
}
