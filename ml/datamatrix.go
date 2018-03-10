// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// DataMatrix holds data in matrix format; e.g. for regression computations
//
//    Example:
//             _          _                              _   _
//            |  1  0  -3  |                            |  0  |
//            |  1  3   3  |                (optional)  |  1  |
//        X = |  1  1   4  |                        y = |  1  |
//            |  1  5   0  |                            |  0  |
//            |_ 1  8   5 _|(nSamples x nParams)        |_ 1 _|(nSamples)
//
//    where: nParams = nFeatures + 1
//
type DataMatrix struct {

	// input
	nSamples int        // number of data points (samples). number of rows
	nParams  int        // number of features + 1. number of columns
	hasY     bool       // has y vector
	xMat     *la.Matrix // [nSamples][nParams] matrix with the first column being filled with ones
	yVec     la.Vector  // [nSamples] y-data
	lVec     la.Vector  // [nSamples] l = X⋅θ (linear model)
	params   la.Vector  // [nParams] parameters θ

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
func NewDataMatrix(nSamples, nFeatures int, yData bool) (o *DataMatrix) {

	// main
	o = new(DataMatrix)
	o.nSamples = nSamples
	o.nParams = nFeatures + 1
	o.hasY = yData
	o.xMat = la.NewMatrix(o.nSamples, o.nParams)
	if o.hasY {
		o.yVec = la.NewVector(o.nSamples)
	}
	o.lVec = la.NewVector(o.nSamples)
	o.params = la.NewVector(o.nParams)
	o.xMat.SetCol(0, 1.0)

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
func NewDataMatrixTable(xyRawTable [][]float64) (o *DataMatrix) {
	nSamples := len(xyRawTable)
	if nSamples < 1 {
		chk.Panic("at least one row of data in table must be provided\n")
	}
	nFeatures := len(xyRawTable[0]) - 1
	o = NewDataMatrix(nSamples, nFeatures, true)
	for i := 0; i < nSamples; i++ {
		for j := 0; j < nFeatures; j++ {
			o.SetX(i, j, xyRawTable[i][j])
		}
		o.SetY(i, xyRawTable[i][nFeatures])
	}
	o.stat()
	return
}

// Nsamples returns the number of data points
func (o *DataMatrix) Nsamples() int {
	return o.nSamples
}

// Nfeatures returns the number of features = number of parameters - 1
func (o *DataMatrix) Nfeatures() int {
	return o.nParams - 1
}

// Nparams returns the number of parameters = len(θ)
func (o *DataMatrix) Nparams() int {
	return o.nParams
}

// SetX sets x-value
func (o *DataMatrix) SetX(iData, jFeature int, value float64) {
	o.xMat.Set(iData, 1+jFeature, value) // 1+j maps to augmented array
	o.statOk = false
}

// GetXvalues get all x-values corresponding to feature iFeature
func (o *DataMatrix) GetXvalues(iFeature int) (xValues []float64) {
	return o.xMat.GetCol(1 + iFeature) // 1+j maps to augmented array
}

// SetY sets y-value
func (o *DataMatrix) SetY(iData int, value float64) {
	if !o.hasY {
		chk.Panic("this data set does not contain y values")
	}
	o.yVec[iData] = value
	o.statOk = false
}

// GetYvalues returns all y-values
// NOTE: (1) this function returns an access to the internal slice; i.e. no copies are made
//       (2) do not modify the output slice
func (o *DataMatrix) GetYvalues() (yValues []float64) {
	if !o.hasY {
		chk.Panic("this data set does not contain y values")
	}
	return o.yVec
}

// Normalize normalizes x values
//  Input:
//   useMinMax -- divide by max-min; otherwise by the standard deviation
func (o *DataMatrix) Normalize(useMinMax bool) {
	if !o.statOk {
		o.stat()
	}
	den := o.sigX
	if useMinMax {
		den = o.delX
	}
	for i := 0; i < o.nSamples; i++ {
		for j := 1; j < o.nParams; j++ {
			jFeature := j - 1
			o.xMat.Set(i, j, (o.xMat.Get(i, j)-o.meanX[jFeature])/den[jFeature])
		}
	}
	o.statOk = false
}

// auxiliary /////////////////////////////////////////////////////////////////////////////////////

// stat computes statistics
func (o *DataMatrix) stat() {

	// x values
	for J := 0; J < o.Nfeatures(); J++ {
		j := 1 + J
		o.minX[J] = o.xMat.Get(0, j)
		o.maxX[J] = o.minX[J]
		o.sumX[J] = 0.0
		for i := 0; i < o.nSamples; i++ {
			xval := o.xMat.Get(i, j)
			o.minX[J] = utl.Min(o.minX[J], xval)
			o.maxX[J] = utl.Max(o.maxX[J], xval)
			o.sumX[J] += xval
		}
		o.meanX[J] = o.sumX[J] / float64(o.nSamples)
		o.sigX[J] = rnd.StatDevFirst(o.xMat.Col(j), o.meanX[J], true)
		o.delX[J] = o.maxX[J] - o.minX[J]
	}

	// y values
	if o.hasY {
		o.minY = o.yVec[0]
		o.maxY = o.minY
		o.sumY = 0.0
		for i := 0; i < o.nSamples; i++ {
			o.minY = utl.Min(o.minY, o.yVec[i])
			o.maxY = utl.Max(o.maxY, o.yVec[i])
			o.sumY += o.yVec[i]
		}
		o.meanY = o.sumY / float64(o.nSamples)
		o.sigY = rnd.StatDevFirst(o.yVec, o.meanY, true)
		o.delY = o.maxY - o.minY
	}

	// flag
	o.statOk = true
}

// newXvec creates a new xVec[1+nFeatures] vector by initializing it with some values
func (o *DataMatrix) newXvec(initWithMin, initWithMax, initWithMean bool) (xVec la.Vector) {
	xVec = la.NewVector(o.nParams)
	if initWithMin || initWithMax || initWithMax {
		if !o.statOk {
			o.stat()
		}
	}
	if initWithMin {
		copy(xVec, o.minX)
	}
	if initWithMax {
		copy(xVec, o.maxX)
	}
	if initWithMean {
		copy(xVec, o.meanX)
	}
	xVec[0] = 1.0
	return
}
