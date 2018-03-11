// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestLinReg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinReg01. Basic functionality.")

	// data
	data := NewDataGivenRawXY(dataReg01)

	// parameters
	params := NewParamsReg(data.Nfeatures)

	// regression
	reg := NewLinReg(data, params, "reg01")

	// check stat
	chk.Float64(tst, "reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.87)
	chk.Float64(tst, "reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.33)
	data.X.Set(6, 0, 0.88)
	data.Y[19] = 87.34
	data.NotifyUpdate()
	chk.Float64(tst, "new: reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.88)
	chk.Float64(tst, "new: reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.34)
}
