// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// DataMapper defines maps features into an expanded set of features
type DataMapper interface {
	Map(x, xRaw []float64)                               // maps and augment xRaw into x
	GetMapped(XYraw [][]float64, yData bool) *DataMatrix // returns a new Regression data with mapped/augmented X values
	NumOriginalFeatures() int                            // returns the number of original features, before mapping/augmentation
	NumExtraFeatures() int                               // returns the added number of features
}

// PolyDataMapper maps features to polynomial
type PolyDataMapper struct {
	nOriFeatures   int // original number of features
	iFeature       int // selected iFeature to map
	jFeature       int // selected jFeature to map
	order          int // order of polynomial
	nExtraFeatures int // number of added features
}

// NewPolyDataMapper returns a new object
func NewPolyDataMapper(nOriFeatures, iFeature, jFeature, order int) (o *PolyDataMapper) {

	// check
	if order < 2 {
		chk.Panic("PolyDataMapper is useful for order >= 2. order = %d is invalid\n", order)
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
	o.order = order

	// derived
	p := o.order + 1               // auxiliary
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

// Map maps and augment xRaw into x and ignores y[:] = xyRaw[len(xyRaw)-1]
//  Input:
//    xRaw -- array with x values
//  Output:
//    x -- pre-allocated vector such that len(x) = len(xRaw) + 1 because x is augmented with ones
func (o *PolyDataMapper) Map(x, xRaw []float64) {

	// copy existent features
	x[0] = 1.0 // constant value
	for j := 0; j < o.nOriFeatures; j++ {
		x[1+j] = xRaw[j]
	}

	// mapped features
	xi := xRaw[o.iFeature]
	xj := xRaw[o.jFeature]

	// compute new features
	k := o.nOriFeatures
	for e := 2; e <= o.order; e++ {
		for d := 0; d <= e; d++ {
			x[1+k] = math.Pow(xi, float64(e-d)) * math.Pow(xj, float64(d))
			k++
		}
	}
}

// GetMapped returns a new Regression data with mapped/augmented X values
func (o *PolyDataMapper) GetMapped(XYraw [][]float64, yData bool) (rd *DataMatrix) {

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
	rd = NewDataMatrix(nRows, o.NumOriginalFeatures()+o.NumExtraFeatures(), yData)
	x := make([]float64, rd.Nparams())
	for i := 0; i < nRows; i++ {
		o.Map(x, XYraw[i])
		for j := 1; j < rd.Nparams(); j++ {
			jFeature := j - 1
			rd.SetX(i, jFeature, x[j])
		}
		if rd.hasY {
			rd.SetY(i, XYraw[i][o.nOriFeatures])
		}
	}
	return
}
