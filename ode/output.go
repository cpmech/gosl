// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// Output holds output data
type Output struct {

	// input
	ndim int     // dimension of system
	conf *Config // configuration

	// discrete output at accepted steps
	IdxSave int         // current index in Xvalues and Yvalues == last output
	Hvalues []float64   // h values [IdxSave]
	Xvalues []float64   // X values [IdxSave]
	Yvalues []la.Vector // Y values [IdxSave][ndim]
	stepOK  bool        // step output is OK (activated)

	// continuous output
	ContIdx int         // current index in Xcont and Ycont arrays
	ContX   []float64   // X values during continuous output [IdxCont]
	ContY   []la.Vector // Y values during continuous output [IdxCont][ndim]
	ContStp []int       // index of step
	contOK  bool        // continuous output is OK (activated)
	xout    float64     // current x of continuous output
	yout    la.Vector   // current y of continuous output (used if ContF != nil only)

	// from RK method
	cout func(yout la.Vector, h, x float64, y la.Vector, xout float64) // function to calculate continuous values of y
}

// NewOutput returns a new structure
//  ndim -- dimension of problem
//  conf -- configuration
func NewOutput(ndim int, conf *Config) (o *Output) {
	o = new(Output)
	o.ndim = ndim
	o.conf = conf
	o.stepOK = conf.StepNmax > 0
	o.contOK = conf.ContNmax > 0 && conf.ContDx > 0
	if o.stepOK {
		o.Hvalues = make([]float64, o.conf.StepNmax)
		o.Xvalues = make([]float64, o.conf.StepNmax)
		o.Yvalues = make([]la.Vector, o.conf.StepNmax)
	}
	if o.contOK {
		o.ContX = make([]float64, o.conf.ContNmax)
		o.ContY = make([]la.Vector, o.conf.ContNmax)
		o.ContStp = make([]int, o.conf.ContNmax)
	}
	if o.contOK || o.conf.ContF != nil {
		o.yout = la.NewVector(ndim)
	}
	return
}

// Execute executes output; e.g. call Fcn and saves x and y values
func (o *Output) Execute(istep int, last bool, h, x float64, y []float64) (stop bool, err error) {

	// step output using function
	if o.conf.StepF != nil {
		stop, err = o.conf.StepF(istep, h, x, y)
		if stop || err != nil {
			return
		}
	}

	// save step output
	if o.IdxSave < o.conf.StepNmax {
		o.Hvalues[o.IdxSave] = h
		o.Xvalues[o.IdxSave] = x
		o.Yvalues[o.IdxSave] = la.NewVector(o.ndim)
		o.Yvalues[o.IdxSave].Apply(1, y)
		o.IdxSave++
	}

	// continuous output using function
	var xo float64
	if o.conf.ContF != nil {
		if istep == 0 || last {
			xo = x
			o.yout.Apply(1, y)
			stop, err = o.conf.ContF(istep, h, x, y, xo, o.yout)
			if stop || err != nil {
				return
			}
			xo = o.conf.ContDx
		} else {
			xo = o.xout
			for x >= xo {
				o.cout(o.yout, h, x, y, xo)
				stop, err = o.conf.ContF(istep, h, x, y, xo, o.yout)
				if stop || err != nil {
					return
				}
				xo += o.conf.ContDx
			}
		}
	}

	// save continuous output
	if o.contOK && o.ContIdx < o.conf.ContNmax {
		if istep == 0 || last {
			xo = x
			o.ContX[o.ContIdx] = xo
			o.ContY[o.ContIdx] = la.NewVector(o.ndim)
			o.ContY[o.ContIdx].Apply(1, y)
			o.ContStp[o.ContIdx] = istep
			o.ContIdx++
			xo = o.conf.ContDx
		} else {
			xo = o.xout
			for x >= xo {
				o.ContStp[o.ContIdx] = istep
				o.ContX[o.ContIdx] = xo
				o.ContY[o.ContIdx] = la.NewVector(o.ndim)
				o.cout(o.ContY[o.ContIdx], h, x, y, xo)
				o.ContIdx++
				xo += o.conf.ContDx
			}
		}
	}

	// set xout
	o.xout = xo
	return
}

// GetStepH returns all h values
// from the (accepted) steps output
func (o *Output) GetStepH() (X []float64) {
	return o.Hvalues[:o.IdxSave]
}

// GetStepX returns all x values
func (o *Output) GetStepX() (X []float64) {
	return o.Xvalues[:o.IdxSave]
}

// GetStepY extracts the y[i] values for all output times
// from the (accepted) steps output
//  i -- index of y component
//  use to plot time series; e.g.:
//     plt.Plot(o.GetStepX(), o.GetStepY(0), &plt.A{L:"y0"})
func (o *Output) GetStepY(i int) (Y []float64) {
	if o.IdxSave > 0 {
		Y = make([]float64, o.IdxSave)
		for j := 0; j < o.IdxSave; j++ {
			Y[j] = o.Yvalues[j][i]
		}
	}
	return
}

// GetStepYtable returns a table with all y values such that Y[idxOut][dim]
// from the (accepted) steps output
func (o *Output) GetStepYtable() (Y [][]float64) {
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

// GetStepYtableT returns a (transposed) table with all y values such that Y[dim][idxOut]
// from the (accepted) steps output
func (o *Output) GetStepYtableT() (Y [][]float64) {
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
