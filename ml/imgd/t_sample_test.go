// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imgd

import (
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/la"
)

var (
	charsX = la.NewMatrixDeep2([][]float64{
		{0, 1, 1, 0 /**/, 1, 0, 0, 1 /**/, 1, 1, 1, 1 /**/, 1, 0, 0, 1}, // A
		{1, 1, 1, 1 /**/, 1, 0, 0, 0 /**/, 1, 0, 0, 0 /**/, 1, 1, 1, 1}, // C
		{1, 1, 1, 0 /**/, 1, 0, 0, 1 /**/, 1, 0, 0, 1 /**/, 1, 1, 1, 0}, // D
		{1, 1, 1, 1 /**/, 1, 0, 0, 1 /**/, 1, 0, 0, 1 /**/, 1, 1, 1, 1}, // O
		{0, 0, 0, 1 /**/, 0, 0, 1, 0 /**/, 0, 1, 0, 0 /**/, 1, 0, 0, 0}, // /
		{1, 0, 0, 0 /**/, 0, 1, 0, 0 /**/, 0, 0, 1, 0 /**/, 0, 0, 0, 1}, // \
		{1, 1, 1, 1 /**/, 0, 1, 1, 0 /**/, 0, 1, 1, 0 /**/, 1, 1, 1, 1}, // I
		{1, 0, 0, 1 /**/, 0, 1, 1, 0 /**/, 0, 1, 1, 0 /**/, 1, 0, 0, 1}, // X
	})
)

func TestSample01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Sample01. Basic functionality")

	smin, smax := 0.0, 1.0

	sampleA := NewGraySample(0, 16, charsX, true) // A
	chk.Int(tst, "idx", sampleA.idx, 0)
	chk.Int(tst, "width", sampleA.width, 4)
	chk.Int(tst, "height", sampleA.height, 4)
	img := sampleA.GetImage(10)
	sampleA.Paint(img, 2, 2, smin, smax)

	io.Pl()
	chk.Int(tst, "A: pix@(2,2)", int(img.GrayAt(2, 2).Y), 0)
	chk.Int(tst, "A: pix@(3,2)", int(img.GrayAt(3, 2).Y), 255)
	chk.Int(tst, "A: pix@(4,2)", int(img.GrayAt(4, 2).Y), 255)
	chk.Int(tst, "A: pix@(5,2)", int(img.GrayAt(5, 2).Y), 0)

	io.Pl()
	chk.Int(tst, "A: pix@(2,3)", int(img.GrayAt(2, 3).Y), 255)
	chk.Int(tst, "A: pix@(3,3)", int(img.GrayAt(3, 3).Y), 0)
	chk.Int(tst, "A: pix@(4,3)", int(img.GrayAt(4, 3).Y), 0)
	chk.Int(tst, "A: pix@(5,3)", int(img.GrayAt(5, 3).Y), 255)

	io.Pl()
	chk.Int(tst, "A: pix@(2,4)", int(img.GrayAt(2, 4).Y), 255)
	chk.Int(tst, "A: pix@(3,4)", int(img.GrayAt(3, 4).Y), 255)
	chk.Int(tst, "A: pix@(4,4)", int(img.GrayAt(4, 4).Y), 255)
	chk.Int(tst, "A: pix@(5,4)", int(img.GrayAt(5, 4).Y), 255)

	io.Pl()
	chk.Int(tst, "A: pix@(2,5)", int(img.GrayAt(2, 5).Y), 255)
	chk.Int(tst, "A: pix@(3,5)", int(img.GrayAt(3, 5).Y), 0)
	chk.Int(tst, "A: pix@(4,5)", int(img.GrayAt(4, 5).Y), 0)
	chk.Int(tst, "A: pix@(5,5)", int(img.GrayAt(5, 5).Y), 255)

	io.Pl()
	sampleC := NewGraySample(1, 16, charsX, true) // C
	chk.Int(tst, "idx", sampleC.idx, 1)
	chk.Int(tst, "width", sampleC.width, 4)
	chk.Int(tst, "height", sampleC.height, 4)
	sampleC.Paint(img, 8, 8, smin, smax)

	io.Pl()
	chk.Int(tst, "C: pix@( 8, 8)", int(img.GrayAt(8, 8).Y), 255)
	chk.Int(tst, "C: pix@( 9, 8)", int(img.GrayAt(9, 8).Y), 255)
	chk.Int(tst, "C: pix@(10, 8)", int(img.GrayAt(10, 8).Y), 255)
	chk.Int(tst, "C: pix@(11, 8)", int(img.GrayAt(11, 8).Y), 255)

	io.Pl()
	chk.Int(tst, "C: pix@( 8, 9)", int(img.GrayAt(8, 9).Y), 255)
	chk.Int(tst, "C: pix@( 9, 9)", int(img.GrayAt(9, 9).Y), 0)
	chk.Int(tst, "C: pix@(10, 9)", int(img.GrayAt(10, 9).Y), 0)
	chk.Int(tst, "C: pix@(11, 9)", int(img.GrayAt(11, 9).Y), 0)

	io.Pl()
	chk.Int(tst, "C: pix@( 8,10)", int(img.GrayAt(8, 10).Y), 255)
	chk.Int(tst, "C: pix@( 9,10)", int(img.GrayAt(9, 10).Y), 0)
	chk.Int(tst, "C: pix@(10,10)", int(img.GrayAt(10, 10).Y), 0)
	chk.Int(tst, "C: pix@(11,10)", int(img.GrayAt(11, 10).Y), 0)

	io.Pl()
	chk.Int(tst, "C: pix@( 8,11)", int(img.GrayAt(8, 11).Y), 255)
	chk.Int(tst, "C: pix@( 9,11)", int(img.GrayAt(9, 11).Y), 255)
	chk.Int(tst, "C: pix@(10,11)", int(img.GrayAt(10, 11).Y), 255)
	chk.Int(tst, "C: pix@(11,11)", int(img.GrayAt(11, 11).Y), 255)

	if chk.Verbose {
		SavePng("/tmp/gosl/ml/imgd", "sample01", img)
	}
}
