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
	ybar la.Vector  // [m] ybar[i] = (y[i] - 1) / m
	hmy  la.Vector  // [m] hmy[i] = h[i] - y[i]
	aMat *la.Matrix // [m][m] A-matrix for Hessian computation
	bMat *la.Matrix // [n][m] auxiliary matrix: B = Xt*A
	hMat *la.Matrix // [n][n] Hessian matrix = Xt*A*X = B*X
	tmp  la.Vector  // [m] temporary vector; e.g. = h - l
}

// NewLogReg returns new LogReg object
func NewLogReg(data *RegData) (o *LogReg) {
	o = new(LogReg)
	o.Set(data)
	return
}

// h implements the sigmoid function
func h(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z))
}

// Model implements the model equation: logistic(xᵀθ)
//   x -- [1+nFeatures] x-values (augmented)
//   θ -- [1+nFeatures] parameters
func (o *LogReg) Model(x, θ la.Vector) float64 {
	return h(la.VecDot(x, θ))
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
	return sq/float64(data.mData) + la.VecDot(o.ybar, data.lVec)
}

// Deriv computes the derivative of the cost function: dC/dθ
//   Input:
//     data -- regression data
//   Output:
//     dCdθ -- derivative of cost function
func (o *LogReg) Deriv(dCdθ la.Vector, data *RegData) {
	la.MatVecMul(data.lVec, 1, data.xMat, data.thetaVec)
	for i := 0; i < data.mData; i++ {
		o.hmy[i] = h(data.lVec[i]) - data.yVec[i]
	}
	la.MatTrVecMul(dCdθ, 1.0/float64(data.mData), data.xMat, o.hmy)
}

// CalcTheta calculates θ using Newton-Raphson solver
//  solverParams -- nonlinear solver parameters (see num.NlSolver)
//  verbose      -- show nonlinear solver output
func (o *LogReg) CalcTheta(data *RegData, verbose, checkJac bool, tolJac0, tolJac1 float64, solverParams map[string]float64) {

	// constants
	m := data.Ndata()
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
		data.thetaVec.Apply(1, z)
		o.Deriv(fz, data)
	}

	// Jacobian function
	Jfcn := func(dfdz *la.Matrix, z la.Vector) {
		la.MatVecMul(data.lVec, 1, data.xMat, z) // l := X⋅θ (linear model)
		for i := 0; i < m; i++ {
			hi := h(data.lVec[i])
			o.aMat.Set(i, i, hi*(1.0-hi)/mCoef)
		}
		la.MatTrMatTrMul(o.bMat, 1, data.xMat, o.aMat) // B := Xt*A
		la.MatMatMul(dfdz, 1, o.bMat, data.xMat)       // H := Xt*A*X
	}

	// debug
	if checkJac {
		if tolJac0 < 1e-10 {
			tolJac0 = 1e-3
		}
		tst := new(testing.T)
		num.CompareJacDense(tst, ffcn, Jfcn, data.thetaVec, tolJac0)
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
	z := data.thetaVec.GetCopy()

	// solve nonlinear problem
	silent := !verbose
	useDenseJacobian := true
	numericalJacobian := false
	var solver num.NlSolver
	defer solver.Free()
	solver.Init(n, ffcn, nil, Jfcn, useDenseJacobian, numericalJacobian, solverParams)
	solver.Solve(z, silent)

	// results
	data.thetaVec.Apply(1, z)

	// debug
	if checkJac {
		if tolJac1 < 1e-10 {
			tolJac1 = 1e-3
		}
		tst := new(testing.T)
		num.CompareJacDense(tst, ffcn, Jfcn, data.thetaVec, tolJac1)
	}
}
