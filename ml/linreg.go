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
//   x -- [nFeatures] x-values
//   θ -- [1+nFeatures] parameters
func (o *LinReg) Model(x, θ la.Vector) float64 {
	return θ[0] + la.VecDot(x, θ[1:])
}

// Cost computes the total cost
func (o *LinReg) Cost(data *RegData) float64 {
	if len(o.e) != data.m {
		o.e = la.NewVector(data.m)
	}
	la.MatVecMul(data.l, 1, data.x, data.θ)
	la.VecAdd(o.e, 1, data.l, -1, data.y)
	return la.VecDot(o.e, o.e) / float64(2*data.m)
}

// Deriv computes the derivative of the cost function: dC/dθ
//   Input:
//     data -- regression data
//   Output:
//     dCdθ -- derivative of cost function
func (o *LinReg) Deriv(dCdθ la.Vector, data *RegData) {
	if len(o.e) != data.m {
		o.e = la.NewVector(data.m)
	}
	la.MatVecMul(data.l, 1, data.x, data.θ)
	la.VecAdd(o.e, 1, data.l, -1, data.y)
	la.MatTrVecMul(dCdθ, 1.0/float64(data.m), data.x, o.e)
}

// CalcTheta calculates θ using the analytical solution
//   Solve:  XᵀX θ = Xᵀy
//   TODO: handle pseudo-inverse cases
func (o *LinReg) CalcTheta(data *RegData) {
	XtX := la.NewMatrix(data.n, data.n)
	la.MatTrMatMul(XtX, 1, data.x, data.x)
	b := la.NewVector(data.n)
	la.MatTrVecMul(b, 1, data.x, data.y)
	la.DenSolve(data.θ, XtX, b, false)
}
