// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// GraDescReg implements the gradient-descent training of a Linear Regression model
type GraDescReg struct {

	// settings
	Alpha    float64 // step size
	TolCost  float64 // tolerance for cost [may be ≤ 0 => disabled]
	TolDeriv float64 // tolerance for derivative of cost [may be ≤ 0 => disabled]
	Niter    int     // final number of iterations

	// control
	maxNit int // maximum number of iterations

	// results
	Costs []float64 // [Niter; max=maxIter] costs during iterations
}

// NewGraDescReg returns new object
//   maxNit -- maximum of iterations allowed
func NewGraDescReg(maxNit int) (o *GraDescReg) {
	o = new(GraDescReg)
	o.Alpha = 0.1
	o.TolCost = -1  // disabled
	o.TolDeriv = -1 // disabled
	o.maxNit = maxNit
	o.Costs = make([]float64, 1+o.maxNit)
	return
}

// Train trains model using gradient-descent method
func (o *GraDescReg) Train(data *Data, params *ParamsReg, model Regression) {
	o.Costs[0] = model.Cost()
	θ := params.AccessThetas()
	b := params.AccessBias()
	dCdθ := la.NewVector(data.Nfeatures)
	dCdb := 0.0
	for o.Niter = 0; o.Niter < o.maxNit; o.Niter++ {
		dCdb = model.Gradients(dCdθ)
		la.VecAdd(θ, 1, θ, -o.Alpha, dCdθ) // θ := θ - α⋅dCdθ
		*b = *b - o.Alpha*dCdb             // b := b - α⋅dCdb
		o.Costs[1+o.Niter] = model.Cost()
	}
}

// Plot plots cost versus iterations
func (o *GraDescReg) Plot(args *plt.A) {
	if args == nil {
		args = &plt.A{C: plt.C(2, 0), Lw: 1.25, NoClip: true, M: ".", Ms: 5}
	}
	l := float64(o.Niter)  // last iteration
	c0 := o.Costs[0]       // first cost
	cl := o.Costs[o.Niter] // last cost
	I := utl.LinSpace(0, l, o.Niter+1)
	plt.Plot(I, o.Costs, args)
	plt.Text(0, c0, io.Sf("%.6f", c0), &plt.A{C: plt.C(0, 0), Fsz: 7, Ha: "left", Va: "bottom"})
	plt.Text(l, cl, io.Sf("%.6f", cl), &plt.A{C: plt.C(0, 0), Fsz: 7, Ha: "right", Va: "bottom"})
	plt.Gll("$iterations$", "$cost$", nil)
	plt.HideTRborders()
}
