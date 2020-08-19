// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/la"
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

	io.Pf("\nStat: X\n")
	chk.Array(tst, "min(x)", 1e-15, dataBackup.Stat.MinX, []float64{-4, -8, -3})
	chk.Array(tst, "max(x)", 1e-15, dataBackup.Stat.MaxX, []float64{3, 5, 5})
	chk.Array(tst, "mean(x)", 1e-15, dataBackup.Stat.MeanX, []float64{-0.6, 0.2, 1.8})
	chk.Array(tst, "sig(x)", 1e-15, dataBackup.Stat.SigX, []float64{2.701851217221259e+00, 4.969909455915671e+00, 3.271085446759225e+00})
	chk.Array(tst, "sum(x)", 1e-15, dataBackup.Stat.SumX, []float64{-3, 1, 9})

	io.Pf("\nStat: y\n")
	chk.Float64(tst, "min(y)", 1e-15, dataBackup.Stat.MinY, 0)
	chk.Float64(tst, "max(y)", 1e-15, dataBackup.Stat.MaxY, 1)
	chk.Float64(tst, "mean(y)", 1e-15, dataBackup.Stat.MeanY, 3.0/5.0)
	chk.Float64(tst, "sig(y)", 1e-15, dataBackup.Stat.SigY, 5.477225575051662e-01)
	chk.Float64(tst, "sum(y)", 1e-15, dataBackup.Stat.SumY, 3)
}
