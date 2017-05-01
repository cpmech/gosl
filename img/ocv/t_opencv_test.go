// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocv

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// global constants
const (
	alpha_max int = 5
	beta_max  int = 125
)

// global variables
var (
	alpha int // simple contrast control
	beta  int // simple brightness control
)

// function onTrackbar that is called whenever any of alpha or beta changes
func onTrackbar(pos int) {
	io.Pf("TODO: hello from onTrackbar\n")
}

func Test_ocv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ocv01. read image and show in window")

	if chk.Verbose {

		// read image given by user
		image := NewMatFromFile("data/gopher.png")
		defer image.Free()

		/// initialize values
		alpha = 1
		beta = 0

		// create windows
		NamedWindow("New Image", WINDOW_NORMAL)

		// create trackbars
		CreateTrackbar("Contrast Trackbar", "New Image", &alpha, alpha_max, onTrackbar)

		// show image
		Imshow("New Image", image)

		// wait until user press some key
		WaitKey(0)
	}
}
