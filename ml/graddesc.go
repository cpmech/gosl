// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// GradDesc performs a gradient descent to learn θ
type GradDesc struct {
	α        float64   // step size
	tolCost  float64   // tolerance for cost [may be ≤ 0]
	tolDeriv float64   // tolerance for derivative of cost [may be ≤ 0]
	maxIter  int       // maximum of iterations
	Niter    int       // final number of iterations
	Costs    []float64 // [Niter; max=maxIter] costs during iterations
}

// NewGradDesc returns new object
//   maxIter -- maximum iterations
func NewGradDesc(maxIter int) (o *GradDesc) {
	o = new(GradDesc)
	o.α = 1e-3
	o.tolCost = -1
	o.tolDeriv = -1
	o.maxIter = maxIter
	o.Costs = make([]float64, o.maxIter)
	return
}

// SetControl sets stepsize and tolerances
//   α        -- step size
//   tolCost  -- tolerance for cost [may be ≤ 0]
//   tolDeriv -- tolerance for derivative of cost [may be ≤ 0]
func (o *GradDesc) SetControl(α, tolCost, tolDeriv float64) {
	o.α = α
	o.tolCost = tolCost
	o.tolDeriv = tolDeriv
}

// Run performs a gradient descent to learn θ
//   data -- regression data [will have θ modified]
//   reg  -- cost and deriv functions
//   θini -- [n] initial θ values [may be nil]
func (o *GradDesc) Run(data *RegData, reg Regression, θini []float64) {
	if θini == nil {
		data.thetaVec.Fill(0.5)
	} else {
		data.thetaVec.Apply(1, θini)
	}
	dCdθ := la.NewVector(data.nParams)
	for o.Niter = 0; o.Niter < o.maxIter; o.Niter++ {
		reg.Deriv(dCdθ, data)
		la.VecAdd(data.thetaVec, 1, data.thetaVec, -o.α, dCdθ)
		o.Costs[o.Niter] = reg.Cost(data)
		if o.tolCost > 0 {
			if math.Abs(o.Costs[o.Niter]) < o.tolCost {
				break
			}
		}
	}
	if o.tolCost > 0 {
		if o.Niter == o.maxIter {
			chk.Panic("did not converge on tolCost. max number of iterations reached: Niter = %d\n", o.Niter)
		}
	}
}

// PlotCostIter plots cost versus iterations
func (o *GradDesc) PlotCostIter(args *plt.A) {
	if args == nil {
		args = &plt.A{C: plt.C(0, 0), NoClip: true}
	}
	I := utl.LinSpace(1, float64(o.Niter), o.Niter)
	plt.Plot(I, o.Costs[:o.Niter], args)
}
