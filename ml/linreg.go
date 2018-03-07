// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/la"
)

// LinReg implements functions to perform linear regression
type LinReg struct {
	e la.Vector // error = l - y
}

// NewLinReg returns new LinReg object
func NewLinReg() (o *LinReg) {
	o = new(LinReg)
	return
}

// Model implements the model equation: xᵀθ
//   x -- [1+nFeatures] x-values (augmented)
//   θ -- [1+nFeatures] parameters
func (o *LinReg) Model(x, θ la.Vector) float64 {
	return la.VecDot(x, θ)
}

// Cost computes the total cost
func (o *LinReg) Cost(data *RegData) float64 {
	if len(o.e) != data.mData {
		o.e = la.NewVector(data.mData)
	}
	la.MatVecMul(data.lVec, 1, data.xMat, data.thetaVec)
	la.VecAdd(o.e, 1, data.lVec, -1, data.yVec)
	return la.VecDot(o.e, o.e) / float64(2*data.mData)
}

// Deriv computes the derivative of the cost function: dC/dθ
//   Input:
//     data -- regression data
//   Output:
//     dCdθ -- derivative of cost function
func (o *LinReg) Deriv(dCdθ la.Vector, data *RegData) {
	if len(o.e) != data.mData {
		o.e = la.NewVector(data.mData)
	}
	la.MatVecMul(data.lVec, 1, data.xMat, data.thetaVec)
	la.VecAdd(o.e, 1, data.lVec, -1, data.yVec)
	la.MatTrVecMul(dCdθ, 1.0/float64(data.mData), data.xMat, o.e)
}

// CalcTheta calculates θ using the analytical solution
//   Solve:  XᵀX θ = Xᵀy
//   TODO: handle pseudo-inverse cases
func (o *LinReg) CalcTheta(data *RegData) {
	XtX := la.NewMatrix(data.nParams, data.nParams)
	la.MatTrMatMul(XtX, 1, data.xMat, data.xMat)
	b := la.NewVector(data.nParams)
	la.MatTrVecMul(b, 1, data.xMat, data.yVec)
	la.DenSolve(data.thetaVec, XtX, b, false)
}
