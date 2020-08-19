// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imgd

import (
	"image"
	"image/color"

	"gosl/utl"
)

// GrayBoard holds grayscale figure data
type GrayBoard struct {
	nrow   int         // number of rows in board
	ncol   int         // number of cols in board
	width  int         // spot width
	height int         // spot height
	pad    int         // padding
	img    *image.Gray // image interface
}

// NewGrayBoard returns a new board to display samples
func NewGrayBoard(numSamples, largestSampleWidth, largestSampleHeight, padding int) (o *GrayBoard) {
	o = new(GrayBoard)
	o.pad = padding
	o.nrow, o.ncol = utl.BestSquareApprox(numSamples)
	o.width, o.height = largestSampleWidth, largestSampleHeight
	totalWidth := o.ncol*(o.width+o.pad) + o.pad
	totalHeight := o.nrow*(o.height+o.pad) + o.pad
	o.img = image.NewGray(image.Rect(0, 0, totalWidth, totalHeight))
	return
}

// NumBins returns the computed total number of bins/spots/places in board; i.e. nrow * ncol
func (o *GrayBoard) NumBins() int {
	return o.nrow * o.ncol
}

// DrawBin draws rectangle indicating sample spot
func (o *GrayBoard) DrawBin(x, y int) {
	clr := color.Gray{255}
	del := 0
	if o.pad > 2 {
		del = int(float64(o.pad) / 3.0)
	}
	for i := -del; i < o.height+del; i++ {
		o.img.Set(x-del, y+i, clr)
		o.img.Set(x+del+o.width, y+i, clr)
	}
	for j := -del; j < o.width+del; j++ {
		o.img.Set(x+j, y-del, clr)
		o.img.Set(x+j, y+del+o.height, clr)
	}
	o.img.Set(x+del+o.width, y+del+o.height, clr)
}

// Paint paints board
func (o *GrayBoard) Paint(samples GraySamples, smin, smax float64, drawBins bool) {
	k := 0
	row := o.pad
	for i := 0; i < o.nrow; i++ {
		col := o.pad
		for j := 0; j < o.ncol; j++ {
			if k >= len(samples) {
				return
			}
			samples[k].Paint(o.img, col, row, smin, smax)
			if drawBins {
				o.DrawBin(col, row)
			}
			col += o.width + o.pad
			k++
		}
		row += o.height + o.pad
	}
}

// SavePng saves png figure
func (o *GrayBoard) SavePng(outdir, fnameKey string) {
	SavePng(outdir, fnameKey, o.img)
}
