// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/io/h5"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Board holds figure/board data
type Board struct {
	nrow   int         // number of rows in board
	ncol   int         // number of cols in board
	width  int         // spot width
	height int         // spot height
	pad    int         // padding
	img    *image.Gray // image interface
}

// NewBoard returns a new board to display samples
func NewBoard(numSamples, maxSampleWidth, maxSampleHeight, padding int) (o *Board) {
	o = new(Board)
	o.pad = padding
	o.nrow, o.ncol = utl.BestSquareApprox(numSamples)
	o.width, o.height = maxSampleWidth, maxSampleHeight
	totalWidth := o.ncol*(o.width+o.pad) + o.pad
	totalHeight := o.nrow*(o.height+o.pad) + o.pad
	o.img = image.NewGray(image.Rect(0, 0, totalWidth, totalHeight))
	return
}

// NumSpots returns the computed total number of spots in board; i.e. nrow * ncol
func (o *Board) NumSpots() int {
	return o.nrow * o.ncol
}

// DrawSpot draws rectangle indicating sample spot
func (o *Board) DrawSpot(x, y int) {
	clr := color.Gray{255}
	for j := 0; j < o.width; j++ {
		o.img.Set(x+j, y, clr)
		o.img.Set(x+j, y+o.height, clr)
	}
	for i := 0; i < o.height; i++ {
		o.img.Set(x, y+i, clr)
		o.img.Set(x+o.width, y+i, clr)
	}
}

// Paint paints board
func (o *Board) Paint(samples Samples, smin, smax float64) {
	k := 0
	row := o.pad
	for i := 0; i < o.nrow; i++ {
		col := o.pad
		for j := 0; j < o.ncol; j++ {
			if k >= len(samples) {
				return
			}
			samples[k].Paint(o.img, row, col, smin, smax)
			//o.DrawSpot(row, col)
			col += o.width + o.pad
			k++
		}
		row += o.height + o.pad
	}
}

// SavePng saves png figure
func (o *Board) SavePng(outdir, fnameKey string) {

	// create file
	furl := path.Join(outdir, fnameKey+".png")
	file, err := os.Create(furl)
	if err != nil {
		chk.Panic("cannot create file at")
	}
	defer file.Close()

	// encode png
	err = png.Encode(file, o.img)
	if err != nil {
		chk.Panic("cannot encode image")
	}
	io.Pf("file <%s> written\n", furl)
}

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

	// select samples
	rnd.Init(0)
	samples := NewSamples(X, nSelected)
	smin, smax, wmax, hmax := samples.Stat()

	// board
	board := NewBoard(nSelected, wmax, hmax, padding)
	board.Paint(samples, smin, smax)
	board.SavePng("/tmp/gosl", "angEx4data1fig")
}
