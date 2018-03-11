// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
)

func TestLinReg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinReg01. Basic functionality.")

	// data
	data := NewDataGivenRawXY(dataReg01)

	// regression
	reg := NewLinReg(data, "reg01")

	// check stat
	chk.Float64(tst, "reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.87)
	chk.Float64(tst, "reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.33)
	data.X.Set(6, 0, 0.88)
	data.Y[19] = 87.34
	data.NotifyUpdate()
	chk.Float64(tst, "new: reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.88)
	chk.Float64(tst, "new: reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.34)
}

func TestLinReg02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinReg02. Train simple problem.")

	// data
	data := NewDataGivenRawXY(dataReg01)

	// regression
	reg := NewLinReg(data, "reg01")

	// train
	reg.Train()
	chk.Float64(tst, "cost", 1e-15, reg.Cost(), 5.312454218805082e-01)
	chk.Array(tst, "θ", 1e-12, reg.Params.Theta, []float64{1.494747973211108e+01})
	chk.Float64(tst, "b", 1e-12, reg.Params.Bias, 7.428331424039514e+01)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5})
		plt.Subplot(2, 1, 1)
		pp := NewPlotterReg(data, reg)
		pp.DataY(0)
		pp.ModelY(0, 0.8, 1.6)
		plt.Subplot(2, 1, 2)
		pp.ContourCost(-1, 0, 0, 100, 0, 70)
		plt.Save("/tmp/gosl/ml", "linreg02")
	}
}
