// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/la"
)

// LinReg implements functions to perform linear regression
type LinReg struct {
	θ la.Vector // θ parameters
	b float64   // bias parameter
	e la.Vector // error = l - y
}

// NewLinReg returns new LinReg object
func NewLinReg(data *DataMatrix) (o *LinReg) {
	o = new(LinReg)
	o.θ = la.NewVector(data.Nparams())
	o.e = la.NewVector(data.nSamples)
	return
}

// GetParams return a copy of parameters
func (o *LinReg) GetParams() (θ la.Vector, b float64) {
	return o.θ.GetCopy(), o.b
}

// SetParams set parameters
func (o *LinReg) SetParams(θ la.Vector, b float64) {
	copy(o.θ, θ)
	o.b = b
}

// SetTheta set θ parameter
func (o *LinReg) SetTheta(iFeature int, value float64) {
	o.θ[iFeature] = value
}

// SetBias set b parameter
func (o *LinReg) SetBias(value float64) {
	o.b = value
}

// Model implements the model equation: b + xᵀθ
//   x -- [nFeatures] x-values
func (o *LinReg) Model(x la.Vector) float64 {
	return la.VecDot(x, o.θ)
}

// Cost computes the total cost
func (o *LinReg) Cost(data *DataMatrix) float64 {
	la.MatVecMul(data.lVec, 1, data.xMat, o.θ)
	la.VecAdd(o.e, 1, data.lVec, -1, data.yVec)
	return la.VecDot(o.e, o.e) / float64(2*data.nSamples)
}

// Deriv computes the derivative of the cost function: dC/dθ
//   Input:
//     data -- regression data
//   Output:
//     dCdθ -- derivative of cost function
func (o *LinReg) Deriv(dCdθ la.Vector, data *DataMatrix) {
	if len(o.e) != data.nSamples {
		o.e = la.NewVector(data.nSamples)
	}
	la.MatVecMul(data.lVec, 1, data.xMat, o.θ)
	la.VecAdd(o.e, 1, data.lVec, -1, data.yVec)
	la.MatTrVecMul(dCdθ, 1.0/float64(data.nSamples), data.xMat, o.e)
}

// CalcTheta calculates θ using the analytical solution
//   Solve:  XᵀX θ = Xᵀy
//   TODO: handle pseudo-inverse cases
func (o *LinReg) CalcTheta(data *DataMatrix) {
	XtX := la.NewMatrix(data.nParams, data.nParams)
	la.MatTrMatMul(XtX, 1, data.xMat, data.xMat)
	b := la.NewVector(data.nParams)
	la.MatTrVecMul(b, 1, data.xMat, data.yVec)
	la.DenSolve(o.θ, XtX, b, false)
}
