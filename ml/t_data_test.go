// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

func checkXY01(tst *testing.T, label string, X *la.Matrix, y la.Vector) {
	chk.Deep2(tst, "X"+label, 1e-15, X.GetDeep2(), [][]float64{
		{-1, +0, -3},
		{-2, +3, +3},
		{+3, +1, +4},
		{-4, +5, +0},
		{+1, -8, +5},
	})
	chk.Array(tst, "y"+label, 1e-15, y, []float64{0, 1, 1, 0, 1})
}

func TestData00(tst *testing.T) {
	//verbose()
	chk.PrintTitle("Data00. capture panic")
	defer chk.RecoverTstPanicIsOK(tst)
	NewDataGivenRawXY(nil)
}

func TestData01(tst *testing.T) {
	//verbose()
	chk.PrintTitle("Data01. new data structure")
	data := NewDataGivenRawXY([][]float64{
		{-1, +0, -3, 0},
		{-2, +3, +3, 1},
		{+3, +1, +4, 1},
		{-4, +5, +0, 0},
		{+1, -8, +5, 1},
	})
	checkXY01(tst, "", data.X, data.Y)
	dataBackup := data.GetCopy()
	checkXY01(tst, "bkp", dataBackup.X, dataBackup.Y)
	chk.Int(tst, "nSamples", dataBackup.Nsamples, 5)
	chk.Int(tst, "nFeatures", dataBackup.Nfeatures, 3)
}
