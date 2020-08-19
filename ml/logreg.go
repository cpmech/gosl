// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/fun"
	"gosl/fun/dbf"
	"gosl/la"
	"gosl/num"
	"gosl/opt"
)

// LogReg implements a logistic regression model (Observer of Data)
type LogReg struct {
	ParamsReg // import ParamsReg

	// main
	data *Data // X-y data

	// workspace
	ybar la.Vector // bar{y}: yb[i] = (1 - y[i]) / m
	l    la.Vector // vector l = b⋅o + X⋅θ [nSamples]
	hmy  la.Vector // vector e = h(l) - y [nSamples]
}

// NewLogReg returns a new LogReg object
//   data -- X,y data
func NewLogReg(data *Data) (o *LogReg) {
	o = new(LogReg)
	o.data = data
	o.Init(o.data.Nfeatures)
	o.data.AddObserver(o) // need to recompute yb upon changes on y
	o.ybar = la.NewVector(data.Nsamples)
	o.l = la.NewVector(data.Nsamples)
	o.hmy = la.NewVector(data.Nsamples)
	o.Update() // compute first ybar
	return
}

// Update perform updates after data has been changed (as an Observer)
func (o *LogReg) Update() {
	m := float64(o.data.Nsamples)
	for i := 0; i < o.data.Nsamples; i++ {
		o.ybar[i] = (1.0 - o.data.Y[i]) / m
	}
}

// Predict returns the model evaluation @ {x;θ,b}
//   Input:
//     x -- vector of features
//   Output:
//     y -- model prediction y(x)
func (o *LogReg) Predict(x la.Vector) (y float64) {
	θ := o.AccessThetas()
	b := o.GetBias()
	return fun.Logistic(b + la.VecDot(x, θ)) // g(b + xᵀθ) where g is logistic
}

// Cost returns the cost c(x;θ,b)
//   Input:
//     data -- X,y data
//     params -- θ and b
//     x -- vector of features
//   Output:
//     c -- total cost (model error)
func (o *LogReg) Cost() (c float64) {

	// auxiliary
	m := float64(o.data.Nsamples)
	λ := o.GetLambda()
	θ := o.AccessThetas()

	// cost
	o.calcl()
	sq := o.calcsumq()
	c = sq/m + la.VecDot(o.ybar, o.l)
	if λ > 0 {
		c += (0.5 * λ / m) * la.VecDot(θ, θ) // c += (0.5λ/m) θᵀθ
	}
	return c
}

// AllocateGradient allocate object to compute Gradients
func (o *LogReg) AllocateGradient() (dCdθ la.Vector) {
	return la.NewVector(o.data.Nfeatures)
}

// Gradients returns ∂C/∂θ and ∂C/∂b
//   Output:
//     dCdθ -- ∂C/∂θ
//     dCdb -- ∂C/∂b
func (o *LogReg) Gradients(dCdθ la.Vector) (dCdb float64) {

	// auxiliary
	m := float64(o.data.Nsamples)
	λ := o.GetLambda()
	θ := o.AccessThetas()
	X := o.data.X

	// dCdθ
	o.calcl()                             // l := b + X⋅θ
	o.calchmy()                           // hmy := h(l) - y
	la.MatTrVecMul(dCdθ, 1.0/m, X, o.hmy) // dCdθ := (1/m) Xᵀhmy
	if λ > 0 {
		la.VecAdd(dCdθ, 1, dCdθ, λ/m, θ) // dCdθ += (1/m) θ
	}

	// dCdb
	dCdb = (1.0 / m) * o.hmy.Accum() // dCdb = (1/m) oᵀhmy
	return
}

// AllocateHessian allocate objects to compute Hessian
func (o *LogReg) AllocateHessian() (d, v la.Vector, D, H *la.Matrix) {
	m := o.data.Nsamples
	n := o.data.Nfeatures
	d = la.NewVector(m)
	v = la.NewVector(n)
	D = la.NewMatrix(m, n)
	H = la.NewMatrix(n, n)
	return
}

// Hessian computes the Hessian matrix and other partial derivatives
//
//   Input, if d !=nil, otherwise allocate these four objects:
//     d -- [nSamples]  d[i] = g(l[i]) * [ 1 - g(l[i]) ]  auxiliary vector
//     v -- [nFeatures] v = ∂²C/∂θ∂b second order partial derivative
//     D -- [nSamples][nFeatures]  D[i][j] = d[i]*X[i][j]  auxiliary matrix
//     H -- [nFeatures][nFeatures]  H = ∂²C/∂θ² Hessian matrix
//
//   Output, either new objectos or pointers to the input ones:
//     dNew := d   (allocated here if d == nil)
//     vNew := v   (allocated here if v == nil)
//     Dnew := D   (allocated here if D == nil)
//     Hnew := H   (allocated here if H == nil)
//     w -- H = ∂²C/∂b²
//
func (o *LogReg) Hessian(d, v la.Vector, D, H *la.Matrix) (w float64) {

	// auxiliary
	m := o.data.Nsamples
	n := o.data.Nfeatures
	X := o.data.X
	λ := o.GetLambda()
	mm := float64(m)

	// calc d vector and D matrix
	o.calcl()
	for i := 0; i < m; i++ {
		d[i] = fun.LogisticD1(o.l[i]) // d vector
		for j := 0; j < n; j++ {
			D.Set(i, j, d[i]*X.Get(i, j)) // D matrix   (TODO: optimize this)
		}
	}

	// calc H matrix
	la.MatTrMatMul(H, 1.0/mm, X, D)
	if λ > 0 {
		for i := 0; i < n; i++ {
			H.Set(i, i, H.Get(i, i)+λ/mm) // D += (λ/m) I   (TODO: optimize here?)
		}
	}

	// calc v
	la.MatTrVecMul(v, 1.0/mm, X, d) // v := (1/m) Xᵀd

	// calc w
	w = d.Accum() / mm
	return
}

// Train finds θ and b using Newton's method
//   Input:
//     data -- X,y data
//   Output:
//     params -- θ and b
func (o *LogReg) Train() {

	// auxiliary
	//m := o.data.Nsamples
	n := o.data.Nfeatures

	// allocate arrays
	//dCdθ := la.NewVector(o.data.Nfeatures)
	var w float64
	d, v, D, H := o.AllocateHessian()

	// objective function where z={θ,b} and fz={dCdθ,dCdb}
	ffcn := func(fz, z la.Vector) {
		o.Backup()
		o.SetThetas(z[:n])
		o.SetBias(z[n])
		dCdb := o.Gradients(fz[:n])
		fz[n] = dCdb
		o.Restore(false)
	}

	// Jacobian function
	Jfcn := func(dfdz *la.Matrix, z la.Vector) {
		o.Backup()
		o.SetThetas(z[:n])
		o.SetBias(z[n])
		w = o.Hessian(d, v, D, H)
		for j := 0; j < n; j++ { // TODO: optimize here
			for i := 0; i < n; i++ { //
				dfdz.Set(i, j, H.Get(i, j))
			}
			dfdz.Set(n, j, v[j])
			dfdz.Set(j, n, v[j])
		}
		dfdz.Set(n, n, w)
	}

	// initial values
	z := la.NewVector(n + 1) // {θ, b}
	copy(z, o.AccessThetas())
	z[n] = o.GetBias()

	// debug
	if true { // check Jacobian
		tolJac0 := 1e-4
		tst := new(testing.T)
		num.CompareJacDense(tst, ffcn, Jfcn, z, tolJac0)
	}

	// solver parameters
	solverParams := map[string]float64{
		"atol":    1e-4,
		"rtol":    1e-4,
		"ftol":    1e-4,
		"chkConv": 0,
	}

	// solve nonlinear problem
	silent := false
	useDenseJacobian := true
	numericalJacobian := false
	var solver num.NlSolver
	defer solver.Free()
	solver.Init(n+1, ffcn, nil, Jfcn, useDenseJacobian, numericalJacobian, solverParams)
	solver.Solve(z, silent)

	// results
	o.SetThetas(z[:n])
	o.SetBias(z[n])

	// debug
	if true { // check Jacobian
		tolJac0 := 1e-4
		tst := new(testing.T)
		num.CompareJacDense(tst, ffcn, Jfcn, z, tolJac0)
	}
}

// TrainNumerical trains model using numerical optimizer
//   θini -- initial (trial) θ values
//   bini -- initial (trial) bias
//   method -- method/kind of numerical solver. e.g. conjgrad, powel, graddesc
//   saveHist -- save history
//   control -- parameters to numerical solver. See package 'opt'
func (o *LogReg) TrainNumerical(θini la.Vector, bini float64, method string, saveHist bool, control dbf.Params) (minCost float64, hist *opt.History) {

	// auxiliary
	n := o.data.Nfeatures

	// set optimization problem
	// v = {θ, b}  ⇒  θ = v[:n], b = v[n]
	problem := &opt.Problem{
		Ndim: n + 1, // nθ + bias
		Ffcn: func(v la.Vector) float64 {
			o.SetThetas(v[:n])
			o.SetBias(v[n])
			return o.Cost()
		},
		Gfcn: func(g, v la.Vector) {
			o.SetThetas(v[:n])
			o.SetBias(v[n])
			g[n] = o.Gradients(g[:n])
		},
		Hfcn: func(f *la.Matrix, v la.Vector) {
			chk.Panic("cannot use Hessian yet\n")
		},
	}

	// initial solution
	v := la.NewVector(n + 1)
	copy(v[:n], θini)
	v[n] = bini

	// solve
	solver := opt.GetNonLinSolver(method, problem)
	solver.SetUseHistory(saveHist)
	minCost = solver.Min(v, control)
	hist = solver.AccessHistory()

	// set params
	o.SetThetas(v[:n])
	o.SetBias(v[n])
	return
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// calce calculates l vector (saves into o.l) (linear model)
//  Output: l = b⋅o + X⋅θ
func (o *LogReg) calcl() {
	θ := o.AccessThetas()
	b := o.GetBias()
	X := o.data.X
	o.l.Fill(b)                   // l := b⋅o
	la.MatVecMulAdd(o.l, 1, X, θ) // l := b⋅o + X⋅θ
}

// calcsumq calculates Σq[i]
//  Input:
//    l -- precomputed o.l
//  Output:
//    sq -- sum(q)
func (o *LogReg) calcsumq() (sq float64) {
	for i := 0; i < o.data.Nsamples; i++ {
		sq += safeLog1pExp(o.l[i])
	}
	return
}

// calchmy calculates h(l) - y vector (saves into o.hmy)
//  Input:
//    l -- precomputed o.l
//  Output:
//    hmy -- computes hmy = h(l) - y
func (o *LogReg) calchmy() {
	for i := 0; i < o.data.Nsamples; i++ {
		o.hmy[i] = fun.Logistic(o.l[i]) - o.data.Y[i]
	}
}

// safeLog1pExp computes log(1+exp(-z)) safely by checking if exp(-z) is >> 1,
// thus returning -z. This is the case when z<0 and |z| is too large
func safeLog1pExp(z float64) float64 {
	if z < -500 {
		return -z
	}
	return math.Log(1.0 + math.Exp(-z))
}
