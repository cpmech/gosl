// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/la"
)

// LinReg implements a linear regression model
type LinReg struct {

	// main
	data   *Data      // X-y data
	stat   *Stat      // statistics
	Params *ParamsReg // θ and b

	// workspace
	e la.Vector // vector e = b⋅o + X⋅θ - y [nSamples]
}

// NewLinReg returns a new LinReg object
//   Input:
//     data -- X,y data
//     params -- θ and b
//     name -- unique name of this object
func NewLinReg(data *Data, name string) (o *LinReg) {
	o = new(LinReg)
	o.data = data
	o.Params = NewParamsReg(data.Nfeatures)
	o.stat = NewStat(data, name)
	o.stat.Update()
	o.e = la.NewVector(data.Nsamples)
	return
}

// BackupParams saves θ and b
func (o *LinReg) BackupParams() {
	o.Params.Backup()
}

// SetParam sets parameter
//  i -- index of θ or -1 for bias
func (o *LinReg) SetParam(i int, value float64) {
	if i < 0 {
		o.Params.Bias = value
		return
	}
	o.Params.Theta[i] = value
}

// GetParam sets parameter
//  i -- index of θ or -1 for bias
func (o *LinReg) GetParam(i int) (value float64) {
	if i < 0 {
		return o.Params.Bias
	}
	return o.Params.Theta[i]
}

// RestoreParams restores θ and b
func (o *LinReg) RestoreParams() {
	o.Params.Restore()
}

// Predict returns the model evaluation @ {x;θ,b}
//   Input:
//     x -- vector of features
//   Output:
//     y -- model prediction y(x)
func (o *LinReg) Predict(x la.Vector) (y float64) {
	return o.Params.Bias + la.VecDot(x, o.Params.Theta) // b + xᵀθ
}

// Cost returns the cost c(x;θ,b)
//   Input:
//     data -- X,y data
//     params -- θ and b
//     x -- vector of features
//   Output:
//     c -- total cost (model error)
func (o *LinReg) Cost() (c float64) {

	// auxiliary
	m := float64(o.data.Nsamples)
	λ := o.Params.Lambda
	θ := o.Params.Theta

	// cost
	o.calce()                           // e := b⋅o + X⋅θ - y
	c = (0.5 / m) * la.VecDot(o.e, o.e) // C := (0.5/m) eᵀe
	if λ > 0 {
		c += (0.5 * λ / m) * la.VecDot(θ, θ) // c += (0.5λ/m) θᵀθ
	}
	return c
}

// Gradients return ∂C/∂θ and ∂C/∂b
//   Output:
//     dCdT -- ∂C/∂θ
//     dCdb -- ∂C/∂b
func (o *LinReg) Gradients(dCdT la.Vector) (dCdb float64) {
	return 0
}

// Train finds θ and b using closed-form solution
//   Input:
//     data -- X,y data
//   Output:
//     params -- θ and b
func (o *LinReg) Train() {

	// auxiliary
	λ := o.Params.Lambda
	X, y := o.data.X, o.data.Y
	s, t := o.stat.SumVars()

	// r vector
	m := float64(o.data.Nsamples)
	n := o.data.Nfeatures
	r := la.NewVector(n)
	la.MatTrVecMul(r, 1, X, y)  // r := a = Xᵀy
	la.VecAdd(r, 1, r, -t/m, s) // r := a - (t/m)s

	// K matrix
	B := la.NewMatrix(n, n)
	K := la.NewMatrix(n, n)
	la.VecVecTrMul(B, 1.0/m, s, s) // B := (1/m) ssᵀ
	la.MatTrMatMul(K, 1, X, X)     // K := A = XᵀX
	la.MatAdd(K, 1, K, -1, B)      // K := A - B
	if λ > 0 {
		for i := 0; i < n; i++ {
			K.Set(i, i, K.Get(i, i)+λ) // K := A - B + λI
		}
	}

	// solve system
	θ := o.Params.Theta
	la.DenSolve(θ, K, r, false)
	o.Params.Bias = (t - la.VecDot(s, θ)) / m
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// calce calculates e vector (save into o.e)
//  Output: e = b⋅o + X⋅θ - y
func (o *LinReg) calce() {
	θ, b := o.Params.Theta, o.Params.Bias
	X, y := o.data.X, o.data.Y
	o.e.Fill(b)                   // e := b⋅o
	la.MatVecMulAdd(o.e, 1, X, θ) // e := b⋅o + X⋅θ
	la.VecAdd(o.e, 1, o.e, -1, y) // e := b⋅o + X⋅θ - y
}
