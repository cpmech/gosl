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
	chk.PrintTitle("ParamsReg01. Parameters for regression")

	nFeatures := 3
	params := NewParamsReg(nFeatures)
	params.theta[0] = 1
	params.theta[1] = 2
	params.theta[2] = 3
	params.bias = 4
	params.lambda = 0.1
	params.degree = 3
	chk.Array(tst, "θ", 1e-15, params.theta, []float64{1, 2, 3})
	chk.Float64(tst, "b", 1e-15, params.bias, 4)
	chk.Float64(tst, "λ", 1e-15, params.lambda, 0.1)
	chk.Int(tst, "P", params.degree, 3)
	chk.Array(tst, "θcopy", 1e-15, params.bkpTheta, nil)
	chk.Float64(tst, "bcopy", 1e-15, params.bkpBias, 0)
	chk.Float64(tst, "λcopy", 1e-15, params.bkpBias, 0)
	chk.Int(tst, "Pcopy", params.bkpDegree, 0)
	io.Pl()
	params.Backup()
	chk.Array(tst, "θcopy", 1e-15, params.bkpTheta, []float64{1, 2, 3})
	chk.Float64(tst, "bcopy", 1e-15, params.bkpBias, 4)
	chk.Float64(tst, "λcopy", 1e-15, params.bkpLambda, 0.1)
	chk.Int(tst, "Pcopy", params.bkpDegree, 3)
	io.Pl()
	params.theta[1] = -2
	params.bias = -4
	params.lambda = 0.01
	params.degree = 4
	chk.Array(tst, "θchanged", 1e-15, params.theta, []float64{1, -2, 3})
	chk.Float64(tst, "bchanged", 1e-15, params.bias, -4)
	chk.Float64(tst, "λchanged", 1e-15, params.lambda, 0.01)
	chk.Int(tst, "Pchanged", params.degree, 4)
	chk.Array(tst, "θcopy", 1e-15, params.bkpTheta, []float64{1, 2, 3})
	chk.Float64(tst, "bcopy", 1e-15, params.bkpBias, 4)
	chk.Float64(tst, "λcopy", 1e-15, params.bkpLambda, 0.1)
	chk.Int(tst, "Pcopy", params.bkpDegree, 3)
	io.Pl()
	params.Restore()
	chk.Array(tst, "θrestored", 1e-15, params.theta, []float64{1, 2, 3})
	chk.Float64(tst, "brestored", 1e-15, params.bias, 4)
	chk.Float64(tst, "λ", 1e-15, params.lambda, 0.1)
	chk.Int(tst, "P", params.degree, 3)
}

type observerT struct {
	name   string
	params *ParamsReg
	msg    string
}

func (o *observerT) Name() string {
	return o.name
}

func (o *observerT) Update() {
	o.msg = io.Sf("got θ=%.1f b=%.1f λ=%.3f p=%1d", o.params.GetThetas(), o.params.GetBias(), o.params.GetLambda(), o.params.GetDegree())
}

func TestParamsReg02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ParamsReg02. Parameters for regression. Observable")

	// parameters
	nFeatures := 3
	params := NewParamsReg(nFeatures)

	// observer
	observer := &observerT{"obs01", params, "nothing here yet"}
	params.AddObserver(observer)
	chk.String(tst, observer.msg, "nothing here yet")

	// check SetParams
	io.Pl()
	params.SetParams([]float64{1, 2, 3}, 4)
	chk.Array(tst, "GetThetas", 1e-15, params.GetThetas(), []float64{1, 2, 3})
	chk.Float64(tst, "GetBias", 1e-15, params.GetBias(), 4)
	chk.Float64(tst, "GetLambda", 1e-15, params.GetLambda(), 0)
	chk.Int(tst, "GetDegree", params.GetDegree(), 0)
	chk.String(tst, observer.msg, "got θ=[1.0 2.0 3.0] b=4.0 λ=0.000 p=0")

	// check SetParam
	io.Pl()
	params.SetParam(0, -1)
	params.SetParam(-1, -4)
	chk.Float64(tst, "GetTheta[0]", 1e-15, params.GetTheta(0), -1)
	chk.Float64(tst, "GetBias", 1e-15, params.GetBias(), -4)
	chk.Float64(tst, "θ=GetParam(0)", 1e-15, params.GetParam(0), -1)
	chk.Float64(tst, "b=GetParam(-1)", 1e-15, params.GetParam(-1), -4)
	chk.String(tst, observer.msg, "got θ=[-1.0 2.0 3.0] b=-4.0 λ=0.000 p=0")

	// check SetThetas
	io.Pl()
	params.SetThetas([]float64{-3, -2, -1})
	chk.Array(tst, "GetThetas", 1e-15, params.GetThetas(), []float64{-3, -2, -1})
	chk.String(tst, observer.msg, "got θ=[-3.0 -2.0 -1.0] b=-4.0 λ=0.000 p=0")

	// check SetTheta
	io.Pl()
	params.SetTheta(0, -33)
	chk.Float64(tst, "GetTheta[0]", 1e-15, params.GetTheta(0), -33)
	chk.String(tst, observer.msg, "got θ=[-33.0 -2.0 -1.0] b=-4.0 λ=0.000 p=0")

	// check SetBias
	io.Pl()
	params.SetBias(123)
	chk.Float64(tst, "GetBias", 1e-15, params.GetBias(), 123)
	chk.String(tst, observer.msg, "got θ=[-33.0 -2.0 -1.0] b=123.0 λ=0.000 p=0")

	// check SetLambda
	io.Pl()
	params.SetLambda(0.01)
	chk.Float64(tst, "GetLambda", 1e-15, params.GetLambda(), 0.01)
	chk.String(tst, observer.msg, "got θ=[-33.0 -2.0 -1.0] b=123.0 λ=0.010 p=0")

	// check SetDegree
	io.Pl()
	params.SetDegree(8)
	chk.Int(tst, "GetDegree", params.GetDegree(), 8)
	chk.String(tst, observer.msg, "got θ=[-33.0 -2.0 -1.0] b=123.0 λ=0.010 p=8")
}
