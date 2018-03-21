// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imgd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestBoard01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Board01. Basic functionality")

	drawBins := false
	smin, smax := 0.0, 1.0
	rowMaj, random := true, false
	sampleWidth, sampleHeight, padding := 4, 4, 1

	board := NewGrayBoard(charsX.M, sampleWidth, sampleHeight, padding)
	samples := NewGraySamples(charsX, charsX.M, rowMaj, random)
	board.Paint(samples, smin, smax, drawBins)

	io.Pl()
	chk.Int(tst, "board(A)", int(board.img.GrayAt(1, 1).Y), 0)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(2, 1).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(3, 1).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(4, 1).Y), 0)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(1, 2).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(2, 2).Y), 0)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(3, 2).Y), 0)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(4, 2).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(1, 3).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(2, 3).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(3, 3).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(4, 3).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(1, 4).Y), 255)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(2, 4).Y), 0)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(3, 4).Y), 0)
	chk.Int(tst, "board(A)", int(board.img.GrayAt(4, 4).Y), 255)

	io.Pl()
	chk.Int(tst, "board(C)", int(board.img.GrayAt(6, 1).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(7, 1).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(8, 1).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(9, 1).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(6, 2).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(7, 2).Y), 0)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(8, 2).Y), 0)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(9, 2).Y), 0)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(6, 3).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(7, 3).Y), 0)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(8, 3).Y), 0)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(9, 3).Y), 0)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(6, 4).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(7, 4).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(8, 4).Y), 255)
	chk.Int(tst, "board(C)", int(board.img.GrayAt(9, 4).Y), 255)

	if chk.Verbose {
		board.SavePng("/tmp/gosl/ml/imgd", "board01")
	}
}
