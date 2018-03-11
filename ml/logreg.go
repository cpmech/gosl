// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
)

// LogReg implements functions to perform the logistic regression
type LogReg struct {
	ybar   la.Vector  // [m] ybar[i] = (y[i] - 1) / m
	hmy    la.Vector  // [m] hmy[i] = h[i] - y[i]
	aMat   *la.Matrix // [m][m] A-matrix for Hessian computation
	bMat   *la.Matrix // [n][m] auxiliary matrix: B = Xt*A
	hMat   *la.Matrix // [n][n] Hessian matrix = Xt*A*X = B*X
	tmp    la.Vector  // [m] temporary vector; e.g. = h - l
	l      la.Vector  // l = X⋅θ
	lambda float64    // regularization parameter
	θ      la.Vector  // θ parameters
	b      float64    // bias parameter
}

// NewLogReg returns new LogReg object
func NewLogReg(data *DataMatrix) (o *LogReg) {
	o = new(LogReg)
	o.l = la.NewVector(data.nSamples)
	o.θ = la.NewVector(data.Nparams())
	o.Set(data)
	return
}

// h implements the sigmoid function
func h(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z))
}

// GetParams return a copy of parameters
func (o *LogReg) GetParams() (θ la.Vector, b float64) {
	return o.θ.GetCopy(), o.b
}

// SetParams set parameters
func (o *LogReg) SetParams(θ la.Vector, b float64) {
	copy(o.θ, θ)
	o.b = b
}

// SetTheta set θ parameter
func (o *LogReg) SetTheta(iFeature int, value float64) {
	o.θ[iFeature] = value
}

// SetBias set b parameter
func (o *LogReg) SetBias(value float64) {
	o.b = value
}

// Model implements the model equation: logistic(xᵀθ)
//   x -- [nFeatures] x-values
func (o *LogReg) Model(x la.Vector) float64 {
	return h(la.VecDot(x, o.θ))
}

// SetRegularization sets regularization parameter and activates regularization
func (o *LogReg) SetRegularization(λ float64) {
	o.lambda = λ
}

// Set sets LogReg with given regression data
//  data -- regressin data where m=numData, n=numParams
func (o *LogReg) Set(data *DataMatrix) {
	if len(o.ybar) != data.nSamples {
		o.ybar = la.NewVector(data.nSamples)
		o.hmy = la.NewVector(data.nSamples)
	}
	for i := 0; i < data.nSamples; i++ {
		o.ybar[i] = (1.0 - data.yVec[i]) / float64(data.nSamples)
	}
}

// Cost computes the total cost
func (o *LogReg) Cost(data *DataMatrix) (C float64) {
	la.MatVecMul(o.l, 1, data.xMat, o.θ)
	sq := 0.0
	for i := 0; i < data.nSamples; i++ {
		sq += math.Log(1.0 + math.Exp(-o.l[i]))
	}
	mCoef := float64(data.nSamples)
	C = sq/mCoef + la.VecDot(o.ybar, o.l)
	if o.lambda > 0 {
		tmp := o.θ[0]
		o.θ[0] = 0.0
		C += (0.5 * o.lambda / mCoef) * la.VecDot(o.θ, o.θ)
		o.θ[0] = tmp
	}
	return C
}

// Deriv computes the derivative of the cost function: dC/dθ
//   Input:
//     data -- regression data
//   Output:
//     dCdθ -- derivative of cost function
func (o *LogReg) Deriv(dCdθ la.Vector, data *DataMatrix) {
	la.MatVecMul(o.l, 1, data.xMat, o.θ)
	for i := 0; i < data.nSamples; i++ {
		o.hmy[i] = h(o.l[i]) - data.yVec[i]
	}
	mCoef := float64(data.nSamples)
	la.MatTrVecMul(dCdθ, 1.0/mCoef, data.xMat, o.hmy)
	if o.lambda > 0 {
		tmp := o.θ[0]
		o.θ[0] = 0.0
		la.VecAdd(dCdθ, 1, dCdθ, o.lambda/mCoef, o.θ)
		o.θ[0] = tmp
	}
}

// CalcTheta calculates θ using Newton-Raphson solver
//  solverParams -- nonlinear solver parameters (see num.NlSolver)
//  verbose      -- show nonlinear solver output
func (o *LogReg) CalcTheta(data *DataMatrix, verbose, checkJac bool, tolJac0, tolJac1 float64, solverParams map[string]float64) {

	// constants
	m := data.Nsamples()
	n := data.Nparams()

	// allocate arrays
	mCoef := float64(m)
	allocate := o.aMat == nil
	if o.aMat != nil {
		allocate = o.aMat.N != m
	}
	if allocate {
		o.aMat = la.NewMatrix(m, m)
		o.bMat = la.NewMatrix(n, m)
		o.hMat = la.NewMatrix(n, n)
		o.tmp = la.NewVector(m)
	}

	// objective function: z=θ  and  fz = Xt*(h-y) / m
	ffcn := func(fz, z la.Vector) {
		o.θ.Apply(1, z)
		o.Deriv(fz, data)
	}

	// Jacobian function
	Jfcn := func(dfdz *la.Matrix, z la.Vector) {
		la.MatVecMul(o.l, 1, data.xMat, z) // l := X⋅θ (linear model)
		for i := 0; i < m; i++ {
			hi := h(o.l[i])
			o.aMat.Set(i, i, hi*(1.0-hi)/mCoef)
		}
		la.MatTrMatTrMul(o.bMat, 1, data.xMat, o.aMat) // B := Xt*A
		la.MatMatMul(dfdz, 1, o.bMat, data.xMat)       // H := Xt*A*X
		if o.lambda > 0 {
			for i := 1; i < data.Nparams(); i++ {
				dfdz.Set(i, i, dfdz.Get(i, i)+o.lambda/mCoef)
			}
		}
	}

	// debug
	if checkJac {
		if tolJac0 < 1e-10 {
			tolJac0 = 1e-3
		}
		tst := new(testing.T)
		num.CompareJacDense(tst, ffcn, Jfcn, o.θ, tolJac0)
	}

	// solver parameters
	if solverParams == nil {
		solverParams = map[string]float64{
			"atol":    1e-4,
			"rtol":    1e-4,
			"ftol":    1e-4,
			"chkConv": 0,
		}
	}

	// solution array := initial values
	z := o.θ.GetCopy()

	// solve nonlinear problem
	silent := !verbose
	useDenseJacobian := true
	numericalJacobian := false
	var solver num.NlSolver
	defer solver.Free()
	solver.Init(n, ffcn, nil, Jfcn, useDenseJacobian, numericalJacobian, solverParams)
	solver.Solve(z, silent)

	// results
	o.θ.Apply(1, z)

	// debug
	if checkJac {
		if tolJac1 < 1e-10 {
			tolJac1 = 1e-3
		}
		tst := new(testing.T)
		num.CompareJacDense(tst, ffcn, Jfcn, o.θ, tolJac1)
	}
}
