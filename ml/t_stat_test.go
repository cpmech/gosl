// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

var (
	dataReg01 = [][]float64{
		{0.99, 90.01},
		{1.02, 89.05},
		{1.15, 91.43},
		{1.29, 93.74},
		{1.46, 96.73},
		{1.36, 94.45},
		{0.87, 87.59},
		{1.23, 91.77},
		{1.55, 99.42},
		{1.40, 93.65},
		{1.19, 93.54},
		{1.15, 92.52},
		{0.98, 90.56},
		{1.01, 89.54},
		{1.11, 89.85},
		{1.20, 90.39},
		{1.26, 93.25},
		{1.32, 93.41},
		{1.43, 94.98},
		{0.95, 87.33},
	}
)

func sample01checkStat(tst *testing.T, o *Stat) {
	if !o.data.UseY {
		tst.Errorf("flag UseY should be true\n")
		return
	}
	io.Pf("X\n")
	chk.Float64(tst, "min(x)", 1e-15, o.MinX[0], 0.87)
	chk.Float64(tst, "max(x)", 1e-15, o.MaxX[0], 1.55)
	chk.Float64(tst, "mean(x)", 1e-15, o.MeanX[0], 1.1960)
	chk.Float64(tst, "sig(x)", 1e-15, o.SigX[0], 0.189303432281837)
	chk.Float64(tst, "sum(x)", 1e-15, o.SumX[0], 23.92)
	io.Pf("y\n")
	chk.Float64(tst, "min(y)", 1e-15, o.MinY, 87.33)
	chk.Float64(tst, "max(y)", 1e-15, o.MaxY, 99.42)
	chk.Float64(tst, "mean(y)", 1e-15, o.MeanY, 92.1605)
	chk.Float64(tst, "sig(y)", 1e-15, o.SigY, 3.020778001913102)
	chk.Float64(tst, "sum(y)", 1e-15, o.SumY, 1843.21)
}

func TestStat01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Stat01. statistics")

	// data
	data := NewDataGivenRawXY(dataReg01)

	// stat
	stat := NewStat(data)

	// notify update
	data.NotifyUpdate()

	// check
	sample01checkStat(tst, stat)

	// s and t sums
	io.Pl()
	s, t := stat.SumVars()
	chk.Array(tst, "s = sum(X)", 1e-15, s, []float64{23.92})
	chk.Float64(tst, "t = sum(y)", 1e-15, t, 1843.21)
}
