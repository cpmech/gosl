// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestParamsReg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ParamsReg01. parameters for regression")

	nFeatures := 3
	params := NewParamsReg(nFeatures)
	params.Theta[0] = 1
	params.Theta[1] = 2
	params.Theta[2] = 3
	params.Bias = 4
	chk.Array(tst, "θ", 1e-15, params.Theta, []float64{1, 2, 3})
	chk.Float64(tst, "b", 1e-15, params.Bias, 4)
	chk.Array(tst, "θ", 1e-15, params.thetaCopy, nil)
	chk.Float64(tst, "b", 1e-15, params.biasCopy, 0)
	params.Backup()
	chk.Array(tst, "θ", 1e-15, params.thetaCopy, []float64{1, 2, 3})
	chk.Float64(tst, "b", 1e-15, params.biasCopy, 4)
	params.Theta[1] = -2
	params.Bias = -4
	chk.Array(tst, "θ", 1e-15, params.Theta, []float64{1, -2, 3})
	chk.Float64(tst, "b", 1e-15, params.Bias, -4)
	chk.Array(tst, "θ", 1e-15, params.thetaCopy, []float64{1, 2, 3})
	chk.Float64(tst, "b", 1e-15, params.biasCopy, 4)
	params.Restore()
	chk.Array(tst, "θ", 1e-15, params.Theta, []float64{1, 2, 3})
	chk.Float64(tst, "b", 1e-15, params.Bias, 4)
}
