// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"

	"gosl/chk"
	"gosl/la"
)

// PolyDataMapper maps features to expanded polynomial
type PolyDataMapper struct {
	nOriFeatures   int // number of original features
	nExtraFeatures int // number of added features
	iFeature       int // selected iFeature to map
	jFeature       int // selected jFeature to map
	degree         int // degree of polynomial
}

// NewPolyDataMapper returns a new object
func NewPolyDataMapper(nOriFeatures, iFeature, jFeature, degree int) (o *PolyDataMapper) {

	// check
	if degree < 2 {
		chk.Panic("PolyDataMapper is useful for degree >= 2. degree = %d is invalid\n", degree)
	}
	if iFeature > nOriFeatures-1 {
		chk.Panic("iFeature must be within [0, %d]. iFeature = %d is invalid\n", nOriFeatures-1, iFeature)
	}
	if jFeature > nOriFeatures {
		chk.Panic("jFeature must be within [0, %d]. jFeature = %d is invalid\n", nOriFeatures-1, jFeature)
	}

	// input data
	o = new(PolyDataMapper)
	o.nOriFeatures = nOriFeatures
	o.iFeature = iFeature
	o.jFeature = jFeature
	o.degree = degree

	// derived
	p := o.degree + 1              // auxiliary
	nPascal := p*(p+1)/2 - 1       // -1 because first row in Pascal triangle is neglected
	o.nExtraFeatures = nPascal - 2 // -2 because iFeature and jFeatureare were considered in nPascal already
	return
}

// NumOriginalFeatures returns the number of original features, before mapping/augmentation
func (o *PolyDataMapper) NumOriginalFeatures() int {
	return o.nOriFeatures
}

// NumExtraFeatures returns the number of extra features added by this mapper
func (o *PolyDataMapper) NumExtraFeatures() int {
	return o.nExtraFeatures
}

// Map maps xRaw into x and ignores y[:] = xyRaw[len(xyRaw)-1]
//  Input:
//    xRaw -- array with x values
//  Output:
//    x -- pre-allocated vector such that len(x) = nFeatures
func (o *PolyDataMapper) Map(x, xRaw la.Vector) {

	// copy existent features
	for j := 0; j < o.nOriFeatures; j++ {
		x[j] = xRaw[j]
	}

	// mapped features
	xi := xRaw[o.iFeature]
	xj := xRaw[o.jFeature]

	// compute new features
	k := o.nOriFeatures
	for e := 2; e <= o.degree; e++ {
		for d := 0; d <= e; d++ {
			x[k] = math.Pow(xi, float64(e-d)) * math.Pow(xj, float64(d))
			k++
		}
	}
}

// GetMapped returns a new Regression data with mapped/augmented X values
func (o *PolyDataMapper) GetMapped(XYraw [][]float64, useY bool) (data *Data) {

	// check
	nRows := len(XYraw)
	if nRows < 1 {
		chk.Panic("need at least 1 data row. nRows = %d is invalid\n", nRows)
	}
	nColumns := len(XYraw[0])
	if nColumns < 3 {
		chk.Panic("need at least 3 columns (x0, x1, y). nColumns = %d is invalid\n", nColumns)
	}
	if nColumns != o.nOriFeatures+1 {
		chk.Panic("number of columns does not correspond to number of original features + 1 (y value). %d != %d\n", nColumns, o.nOriFeatures+1)
	}

	// set data
	nSamples := nRows
	nFeatures := o.NumOriginalFeatures() + o.NumExtraFeatures()
	data = NewData(nSamples, nFeatures, useY, true)
	x := la.NewVector(nFeatures)
	for i := 0; i < nSamples; i++ {
		o.Map(x, XYraw[i])
		for j := 0; j < nFeatures; j++ {
			data.X.Set(i, j, x[j])
		}
		if useY {
			data.Y[i] = XYraw[i][o.nOriFeatures] // last column of XYraw
		}
	}
	return
}
