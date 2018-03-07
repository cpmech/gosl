// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"

	"github.com/cpmech/gosl/la"
)

// LogReg implements functions to perform the logistic regression
type LogReg struct {
	ybar la.Vector // [m] ybar[i] = (y[i] - 1) / m
	hmy  la.Vector // [m] hmy[i] = h[i] - y[i]
}

// NewLogReg returns new LogReg object
func NewLogReg(data *RegData) (o *LogReg) {
	o = new(LogReg)
	o.Set(data)
	return
}

// Model implements the model equation: logistic(xᵀθ)
//   x -- [1+nFeatures] x-values (augmented)
//   θ -- [1+nFeatures] parameters
func (o *LogReg) Model(x, θ la.Vector) float64 {
	z := la.VecDot(x, θ)
	return 1.0 / (1.0 + math.Exp(-z))
}

// Set sets LogReg with given regression data
//  data -- regressin data where m=numData, n=numParams
func (o *LogReg) Set(data *RegData) {
	if len(o.ybar) != data.mData {
		o.ybar = la.NewVector(data.mData)
		o.hmy = la.NewVector(data.mData)
	}
	for i := 0; i < data.mData; i++ {
		o.ybar[i] = (1.0 - data.yVec[i]) / float64(data.mData)
	}
}

// Cost computes the total cost
func (o *LogReg) Cost(data *RegData) float64 {
	la.MatVecMul(data.lVec, 1, data.xMat, data.thetaVec)
	sq := 0.0
	for i := 0; i < data.mData; i++ {
		sq += math.Log(1.0 + math.Exp(-data.lVec[i]))
	}
	sq /= float64(data.mData)
	return sq + la.VecDot(o.ybar, data.lVec)
}

// Deriv computes the derivative of the cost function: dC/dθ
//   Input:
//     data -- regression data
//   Output:
//     dCdθ -- derivative of cost function
func (o *LogReg) Deriv(dCdθ la.Vector, data *RegData) {
	la.MatVecMul(data.lVec, 1, data.xMat, data.thetaVec)
	for i := 0; i < data.mData; i++ {
		o.hmy[i] = 1.0/(1.0+math.Exp(-data.lVec[i])) - data.yVec[i]
	}
	la.MatTrVecMul(dCdθ, 1.0/float64(data.mData), data.xMat, o.hmy)
}
