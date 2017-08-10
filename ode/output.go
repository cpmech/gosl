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
	StepIdx int         // current index in Xvalues and Yvalues == last output
	StepH   []float64   // h values [IdxSave]
	StepX   []float64   // X values [IdxSave]
	StepY   []la.Vector // Y values [IdxSave][ndim]
	stepOK  bool        // step output is OK (activated)

	// continuous output
	ContIdx int         // current index in Xcont and Ycont arrays
	ContS   []int       // index of step
	ContX   []float64   // X values during continuous output [IdxCont]
	ContY   []la.Vector // Y values during continuous output [IdxCont][ndim]
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
		o.StepH = make([]float64, o.conf.StepNmax)
		o.StepX = make([]float64, o.conf.StepNmax)
		o.StepY = make([]la.Vector, o.conf.StepNmax)
	}
	if o.contOK {
		o.ContS = make([]int, o.conf.ContNmax)
		o.ContX = make([]float64, o.conf.ContNmax)
		o.ContY = make([]la.Vector, o.conf.ContNmax)
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
	if o.StepIdx < o.conf.StepNmax {
		o.StepH[o.StepIdx] = h
		o.StepX[o.StepIdx] = x
		o.StepY[o.StepIdx] = la.NewVector(o.ndim)
		o.StepY[o.StepIdx].Apply(1, y)
		o.StepIdx++
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
			o.ContS[o.ContIdx] = istep
			o.ContX[o.ContIdx] = xo
			o.ContY[o.ContIdx] = la.NewVector(o.ndim)
			o.ContY[o.ContIdx].Apply(1, y)
			o.ContIdx++
			xo = o.conf.ContDx
		} else {
			xo = o.xout
			for x >= xo {
				o.ContS[o.ContIdx] = istep
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

// step output ////////////////////////////////////////////////////////////////////////////////////

// GetStepH returns all h values
// from the (accepted) steps output
func (o *Output) GetStepH() (X []float64) {
	return o.StepH[:o.StepIdx]
}

// GetStepX returns all x values
// from the (accepted) steps output
func (o *Output) GetStepX() (X []float64) {
	return o.StepX[:o.StepIdx]
}

// GetStepY extracts the y[i] values for all output times
// from the (accepted) steps output
//  i -- index of y component
//  use to plot time series; e.g.:
//     plt.Plot(o.GetStepX(), o.GetStepY(0), &plt.A{L:"y0"})
func (o *Output) GetStepY(i int) (Y []float64) {
	if o.StepIdx > 0 {
		Y = make([]float64, o.StepIdx)
		for j := 0; j < o.StepIdx; j++ {
			Y[j] = o.StepY[j][i]
		}
	}
	return
}

// GetStepYtable returns a table with all y values such that Y[idxOut][dim]
// from the (accepted) steps output
func (o *Output) GetStepYtable() (Y [][]float64) {
	if len(o.StepY) < 1 {
		return
	}
	ndim := len(o.StepY[0])
	Y = utl.Alloc(o.StepIdx, ndim)
	for j := 0; j < o.StepIdx; j++ {
		for i := 0; i < ndim; i++ {
			Y[j][i] = o.StepY[j][i]
		}
	}
	return
}

// GetStepYtableT returns a (transposed) table with all y values such that Y[dim][idxOut]
// from the (accepted) steps output
func (o *Output) GetStepYtableT() (Y [][]float64) {
	if len(o.StepY) < 1 {
		return
	}
	ndim := len(o.StepY[0])
	Y = utl.Alloc(ndim, o.StepIdx)
	for j := 0; j < o.StepIdx; j++ {
		for i := 0; i < ndim; i++ {
			Y[i][j] = o.StepY[j][i]
		}
	}
	return
}

// continuous output //////////////////////////////////////////////////////////////////////////////

// GetContS returns all s (step-index) values
// from the continuous output data
func (o *Output) GetContS() (S []int) {
	return o.ContS[:o.ContIdx]
}

// GetContX returns all x values
// from the continuous output data
func (o *Output) GetContX() (X []float64) {
	return o.ContX[:o.ContIdx]
}

// GetContY extracts the y[i] values for all output times
// from the continuous output data
//  i -- index of y component
//  use to plot time series; e.g.:
//     plt.Plot(o.GetContX(), o.GetContY(0), &plt.A{L:"y0"})
func (o *Output) GetContY(i int) (Y []float64) {
	if o.ContIdx > 0 {
		Y = make([]float64, o.ContIdx)
		for j := 0; j < o.ContIdx; j++ {
			Y[j] = o.ContY[j][i]
		}
	}
	return
}

// GetContYtable returns a table with all y values such that Y[idxOut][dim]
// from the continuous output data
func (o *Output) GetContYtable() (Y [][]float64) {
	if len(o.ContY) < 1 {
		return
	}
	ndim := len(o.ContY[0])
	Y = utl.Alloc(o.ContIdx, ndim)
	for j := 0; j < o.ContIdx; j++ {
		for i := 0; i < ndim; i++ {
			Y[j][i] = o.ContY[j][i]
		}
	}
	return
}

// GetContYtableT returns a (transposed) table with all y values such that Y[dim][idxOut]
// from the continuous output data
func (o *Output) GetContYtableT() (Y [][]float64) {
	if len(o.ContY) < 1 {
		return
	}
	ndim := len(o.ContY[0])
	Y = utl.Alloc(ndim, o.ContIdx)
	for j := 0; j < o.ContIdx; j++ {
		for i := 0; i < ndim; i++ {
			Y[i][j] = o.ContY[j][i]
		}
	}
	return
}
