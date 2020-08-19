// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package imgd (image-data) adds functionality to process image-data,
// e.g. for pattern recognition using images.
package imgd

import (
	"image"
	"image/color"
	"math"

	"gosl/la"
	"gosl/rnd"
	"gosl/utl"
)

// GraySample holds sample data corresponding go grayscale images
type GraySample struct {
	idx    int       // position in X matrix
	size   int       // total number of pixels = width*height
	data   la.Vector // gray intensities
	width  int       // number of row pixels
	height int       // number of col pixels
	rowMaj bool      // data is stored as a row-major matrix instead of col-major
}

// GraySamples holds a set of gray sample data
type GraySamples []*GraySample

// NewGraySample returns a new sample data
//   idx    -- position in X matrix
//   size   -- total number of pixels = width*height
//   X      -- matrix with all pixel data (intensity of gray)
//   rowMaj -- data is stored as a row-major matrix instead of col-major
func NewGraySample(idx, size int, X *la.Matrix, rowMaj bool) (o *GraySample) {
	width := int(math.Sqrt(float64(size)))
	height := size / width
	return &GraySample{
		idx:    idx,
		width:  width,
		height: height,
		size:   width * height,
		data:   X.GetRow(idx),
		rowMaj: rowMaj,
	}
}

// NewGraySamples returns a set of randomly selected samples
func NewGraySamples(X *la.Matrix, nSelected int, rowMaj, random bool) (selected GraySamples) {
	selected = make([]*GraySample, nSelected)
	nSamples := X.M   // rows of X
	sampleSize := X.N // columns of X
	var idxSelected []int
	if random {
		idxSelected = rnd.IntGetUniqueN(0, nSamples, nSelected)
	} else {
		idxSelected = utl.IntRange(nSelected)
	}
	for i, idx := range idxSelected {
		selected[i] = NewGraySample(idx, sampleSize, X, rowMaj)
	}
	return
}

// GetImage returns an image object sized to this sample
func (o *GraySample) GetImage(pad int) (img *image.Gray) {
	return image.NewGray(image.Rect(0, 0, o.width+pad, o.height+pad))
}

// Paint paints sample into img
//  x0   -- horizontal shift in img
//  y0   -- vertical shift in img
//  smin -- minimum gray intensity for normalization
//  smax -- minimum gray intensity for normalization
func (o *GraySample) Paint(img *image.Gray, x0, y0 int, smin, smax float64) {
	var intensity float64
	for i := 0; i < o.height; i++ {
		for j := 0; j < o.width; j++ {
			if o.rowMaj {
				intensity = o.data[j+i*o.width] // row-major
			} else {
				intensity = o.data[i+j*o.height] // col-major
			}
			scale := (intensity - smin) / (smax - smin)
			clr := uint8(255 * scale)
			img.Set(x0+j, y0+i, color.RGBA{clr, clr, clr, 255})
		}
	}
}

// Stat returns some statistics about all samples
// smin and smax are the min and max intensities
func (o *GraySamples) Stat() (smin, smax float64, maxWidth, maxHeight int) {
	smin = math.MaxFloat64
	smax = math.SmallestNonzeroFloat64
	maxWidth = 0
	maxHeight = 0
	for _, sample := range *o {
		min, max := sample.data.MinMax()
		smin = utl.Min(smin, min)
		smax = utl.Max(smax, max)
		maxWidth = utl.Imax(maxWidth, sample.width)
		maxHeight = utl.Imax(maxHeight, sample.height)
	}
	return
}
