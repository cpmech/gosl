// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

// LogReg implements a logistic regression model
type LogReg struct {

	// main
	name   string     // name of this "observer"
	data   *Data      // X-y data
	params *ParamsReg // parameters: θ, b, λ
	stat   *Stat      // statistics

	// workspace
	ybar la.Vector // bar{y}: yb[i] = (1 - y[i]) / m
	l    la.Vector // vector l = b⋅o + X⋅θ [nSamples]
	hmy  la.Vector // vector e = h(l) - y [nSamples]
}

// NewLogReg returns a new LogReg object
//   Input:
//     data   -- X,y data
//     params -- θ, b, λ
//     name   -- unique name of this (observer) object
func NewLogReg(data *Data, params *ParamsReg, name string) (o *LogReg) {
	o = new(LogReg)
	o.name = name
	o.data = data
	o.data.AddObserver(o) // need to recompute yb upon changes on y
	o.params = params
	o.stat = NewStat(data, "stat_"+name)
	o.stat.Update()
	o.ybar = la.NewVector(data.Nsamples)
	o.l = la.NewVector(data.Nsamples)
	o.hmy = la.NewVector(data.Nsamples)
	o.Update() // compute first ybar
	return
}

// Name returns the name of this "Observer"
func (o *LogReg) Name() string {
	return o.name
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
	θ := o.params.AccessThetas()
	b := o.params.GetBias()
	return h(b + la.VecDot(x, θ)) // h(b + xᵀθ) where h is logistic
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
	λ := o.params.GetLambda()
	θ := o.params.AccessThetas()

	// cost
	o.calcl()
	sq := o.calcsumq()
	c = sq/m + la.VecDot(o.ybar, o.l)
	if λ > 0 {
		c += (0.5 * λ / m) * la.VecDot(θ, θ) // c += (0.5λ/m) θᵀθ
	}
	return c
}

// Gradients returns ∂C/∂θ and ∂C/∂b
//   Output:
//     dCdθ -- ∂C/∂θ
//     dCdb -- ∂C/∂b
func (o *LogReg) Gradients(dCdθ la.Vector) (dCdb float64) {

	// auxiliary
	m := float64(o.data.Nsamples)
	λ := o.params.GetLambda()
	θ := o.params.AccessThetas()
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

// Train finds θ and b using closed-form solution
//   Input:
//     data -- X,y data
//   Output:
//     params -- θ and b
func (o *LogReg) Train() {
	chk.Panic("TODO")
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// calce calculates l vector (saves into o.l) (linear model)
//  Output: l = b⋅o + X⋅θ
func (o *LogReg) calcl() {
	θ := o.params.AccessThetas()
	b := o.params.GetBias()
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
		sq += math.Log(1.0 + math.Exp(-o.l[i]))
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
		o.hmy[i] = h(o.l[i]) - o.data.Y[i]
	}
}

// h implements the sigmoid/logistic function
func h(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z))
}
