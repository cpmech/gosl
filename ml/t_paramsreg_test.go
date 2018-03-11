// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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
	params.Lambda = 0.1
	params.Degree = 3
	chk.Array(tst, "θ", 1e-15, params.Theta, []float64{1, 2, 3})
	chk.Float64(tst, "b", 1e-15, params.Bias, 4)
	chk.Float64(tst, "λ", 1e-15, params.Lambda, 0.1)
	chk.Int(tst, "P", params.Degree, 3)
	chk.Array(tst, "θcopy", 1e-15, params.thetaCopy, nil)
	chk.Float64(tst, "bcopy", 1e-15, params.biasCopy, 0)
	chk.Float64(tst, "λcopy", 1e-15, params.biasCopy, 0)
	chk.Int(tst, "Pcopy", params.degreeCopy, 0)
	io.Pl()
	params.Backup()
	chk.Array(tst, "θcopy", 1e-15, params.thetaCopy, []float64{1, 2, 3})
	chk.Float64(tst, "bcopy", 1e-15, params.biasCopy, 4)
	chk.Float64(tst, "λcopy", 1e-15, params.lambdaCopy, 0.1)
	chk.Int(tst, "Pcopy", params.degreeCopy, 3)
	io.Pl()
	params.Theta[1] = -2
	params.Bias = -4
	params.Lambda = 0.01
	params.Degree = 4
	chk.Array(tst, "θchanged", 1e-15, params.Theta, []float64{1, -2, 3})
	chk.Float64(tst, "bchanged", 1e-15, params.Bias, -4)
	chk.Float64(tst, "λchanged", 1e-15, params.Lambda, 0.01)
	chk.Int(tst, "Pchanged", params.Degree, 4)
	chk.Array(tst, "θcopy", 1e-15, params.thetaCopy, []float64{1, 2, 3})
	chk.Float64(tst, "bcopy", 1e-15, params.biasCopy, 4)
	chk.Float64(tst, "λcopy", 1e-15, params.lambdaCopy, 0.1)
	chk.Int(tst, "Pcopy", params.degreeCopy, 3)
	io.Pl()
	params.Restore()
	chk.Array(tst, "θrestored", 1e-15, params.Theta, []float64{1, 2, 3})
	chk.Float64(tst, "brestored", 1e-15, params.Bias, 4)
	chk.Float64(tst, "λ", 1e-15, params.Lambda, 0.1)
	chk.Int(tst, "P", params.Degree, 3)
}
