// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// ContOutFcn defines a function to compute continuous output
type ContOutFcn func(yout, y la.Vector, xout, x, h float64)

// Output holds output data
type Output struct {

	// output using function
	Fcn OutF // output function

	// discrete output at accepted steps
	IdxSave int         // current index in Xvalues and Yvalues == last output
	Hvalues []float64   // h values [IdxSave]
	Xvalues []float64   // X values [IdxSave]
	Yvalues []la.Vector // Y values [IdxSave][ndim]
	nMaxOut int         // max number of output

	// continuous output
	ContIdx  int         // current index in Xcont and Ycont arrays
	ContX    []float64   // X values during continuous output [IdxCont]
	ContY    []la.Vector // Y values during continuous output [IdxCont][ndim]
	ContStp  []int       // index of step
	contNmax int         // maximum number of continuous output; e.g. xf / ContDx + 1
	contOk   bool        // do continuous output
	xCont    float64     // current continuous x value
	dxCont   float64     // step size for continuous output
	fcnCont  ContOutFcn  // function to calculate continuous values of y
}

// NewOutput returns a new structure
//  fcn     -- the output function. nil ⇒ no output
//  nMaxOut -- maximum number of output at each accepted substep; e.g. nMaxOut = nMaxSteps+1. 0 ⇒ not output
//  xf      -- final x to compute size of continuous output arrays
//  contDx  -- step size for continuous output. 0 ⇒ no output
//  contFcn -- function to calculate continuous values of y. nil ⇒ no output
func NewOutput(fcn OutF, ndim, nMaxOut, contNmax int, contDx float64, contFcn ContOutFcn) (o *Output) {
	o = new(Output)
	o.Fcn = fcn
	o.nMaxOut = nMaxOut
	o.contNmax = contNmax
	o.dxCont = contDx
	o.fcnCont = contFcn
	o.contOk = contNmax > 0 && o.dxCont > 0 && o.fcnCont != nil
	if o.nMaxOut > 0 {
		o.Hvalues = make([]float64, o.nMaxOut)
		o.Xvalues = make([]float64, o.nMaxOut)
		o.Yvalues = make([]la.Vector, o.nMaxOut)
	}
	if o.contOk {
		o.ContX = make([]float64, contNmax)
		o.ContY = make([]la.Vector, contNmax)
		o.ContStp = make([]int, contNmax)
	}
	return
}

// Execute executes output; e.g. call Fcn and saves x and y values
func (o *Output) Execute(istep int, last bool, h, x float64, y []float64) {

	// output using function
	if o.Fcn != nil {
		o.Fcn(istep, h, x, y)
	}

	// discrete output at accepted steps
	if o.IdxSave < o.nMaxOut {
		o.Hvalues[o.IdxSave] = h
		o.Xvalues[o.IdxSave] = x
		o.Yvalues[o.IdxSave] = la.NewVector(len(y))
		o.Yvalues[o.IdxSave].Apply(1, y)
		o.IdxSave++
	}

	// continuous output
	if o.contOk && o.ContIdx < o.contNmax {
		if istep == 0 || last {
			o.xCont = x
			o.ContX[o.ContIdx] = o.xCont
			o.ContY[o.ContIdx] = la.NewVector(len(y))
			o.ContY[o.ContIdx].Apply(1, y)
			o.ContStp[o.ContIdx] = istep
			o.xCont = o.dxCont
			o.ContIdx++
		} else {
			for x >= o.xCont {
				o.ContX[o.ContIdx] = o.xCont
				o.ContY[o.ContIdx] = la.NewVector(len(y))
				o.fcnCont(o.ContY[o.ContIdx], y, o.xCont, x, h)
				o.ContStp[o.ContIdx] = istep
				o.xCont += o.dxCont
				o.ContIdx++
			}
		}
	}
}

// GetH returns all h values
func (o *Output) GetH() (X []float64) {
	return o.Hvalues[:o.IdxSave]
}

// GetX returns all x values
func (o *Output) GetX() (X []float64) {
	return o.Xvalues[:o.IdxSave]
}

// GetYi extracts the y[i] values for all output times
//  i -- index of y component
//  use to plot time series; e.g.:
//     plt.Plot(o.GetX(), o.GetY(0), &plt.A{L:"y0"})
func (o *Output) GetYi(i int) (Yi []float64) {
	if o.IdxSave > 0 {
		Yi = make([]float64, o.IdxSave)
		for j := 0; j < o.IdxSave; j++ {
			Yi[j] = o.Yvalues[j][i]
		}
	}
	return
}

// GetY returns a table with all y values such that ytable[idxOut][dim]
func (o *Output) GetY() (Y [][]float64) {
	if len(o.Yvalues) < 1 {
		return
	}
	ndim := len(o.Yvalues[0])
	Y = utl.Alloc(o.IdxSave, ndim)
	for j := 0; j < o.IdxSave; j++ {
		for i := 0; i < ndim; i++ {
			Y[j][i] = o.Yvalues[j][i]
		}
	}
	return
}

// GetYt returns a (transposed) table with all y values such that ytable[dim][idxOut]
func (o *Output) GetYt() (Y [][]float64) {
	if len(o.Yvalues) < 1 {
		return
	}
	ndim := len(o.Yvalues[0])
	Y = utl.Alloc(ndim, o.IdxSave)
	for j := 0; j < o.IdxSave; j++ {
		for i := 0; i < ndim; i++ {
			Y[i][j] = o.Yvalues[j][i]
		}
	}
	return
}
