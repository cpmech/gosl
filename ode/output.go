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
	StepIdx  int         // current index in Xvalues and Yvalues == last output
	StepH    []float64   // h values [IdxSave]
	StepX    []float64   // X values [IdxSave]
	StepY    []la.Vector // Y values [IdxSave][ndim]
	stepNmax int         // max number of output steps

	// dense output
	DenseIdx  int         // current index in DenseS, DenseX, and DenseY arrays
	DenseS    []int       // index of step
	DenseX    []float64   // X values during dense output [DenseIdx]
	DenseY    []la.Vector // Y values during dense output [DenseIdx][ndim]
	denseNmax int         // max number of dense output
	xout      float64     // current x of dense output
	yout      la.Vector   // current y of dense output (used if denseF != nil only)

	// from RK method
	dout func(yout la.Vector, h, x float64, y la.Vector, xout float64) // function to calculate dense values of y
}

// NewOutput returns a new structure
//  ndim -- dimension of problem
//  conf -- configuration
func NewOutput(ndim int, conf *Config) (o *Output) {
	o = new(Output)
	o.ndim = ndim
	o.conf = conf
	if o.conf.stepOut {
		if o.conf.fixed {
			o.stepNmax = o.conf.fixedNsteps + 1
		} else {
			o.stepNmax = o.conf.NmaxSS + 1
		}
		o.StepH = make([]float64, o.stepNmax)
		o.StepX = make([]float64, o.stepNmax)
		o.StepY = make([]la.Vector, o.stepNmax)
	}
	if o.conf.denseOut {
		o.denseNmax = o.conf.denseNstp + 1
		o.DenseS = make([]int, o.denseNmax)
		o.DenseX = make([]float64, o.denseNmax)
		o.DenseY = make([]la.Vector, o.denseNmax)
	}
	if o.conf.denseF != nil {
		o.yout = la.NewVector(ndim)
	}
	return
}

// Execute executes output; e.g. call Fcn and saves x and y values
func (o *Output) Execute(istep int, last bool, h, x float64, y []float64) (stop bool, err error) {

	// step output using function
	if o.conf.stepF != nil {
		stop, err = o.conf.stepF(istep, h, x, y)
		if stop || err != nil {
			return
		}
	}

	// save step output
	if o.StepIdx < o.stepNmax {
		o.StepH[o.StepIdx] = h
		o.StepX[o.StepIdx] = x
		o.StepY[o.StepIdx] = la.NewVector(o.ndim)
		o.StepY[o.StepIdx].Apply(1, y)
		o.StepIdx++
	}

	// dense output using function
	var xo float64
	if o.conf.denseF != nil {
		if istep == 0 || last {
			xo = x
			o.yout.Apply(1, y)
			stop, err = o.conf.denseF(istep, h, x, y, xo, o.yout)
			if stop || err != nil {
				return
			}
			xo += o.conf.denseDx
		} else {
			xo = o.xout
			for x >= xo {
				o.dout(o.yout, h, x, y, xo)
				stop, err = o.conf.denseF(istep, h, x, y, xo, o.yout)
				if stop || err != nil {
					return
				}
				xo += o.conf.denseDx
			}
		}
	}

	// save dense output
	if o.DenseIdx < o.denseNmax {
		if istep == 0 || last {
			xo = x
			o.DenseS[o.DenseIdx] = istep
			o.DenseX[o.DenseIdx] = xo
			o.DenseY[o.DenseIdx] = la.NewVector(o.ndim)
			o.DenseY[o.DenseIdx].Apply(1, y)
			o.DenseIdx++
			xo = o.conf.denseDx
		} else {
			xo = o.xout
			for x >= xo {
				o.DenseS[o.DenseIdx] = istep
				o.DenseX[o.DenseIdx] = xo
				o.DenseY[o.DenseIdx] = la.NewVector(o.ndim)
				o.dout(o.DenseY[o.DenseIdx], h, x, y, xo)
				o.DenseIdx++
				xo += o.conf.denseDx
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

// dense output ///////////////////////////////////////////////////////////////////////////////////

// GetDenseS returns all s (step-index) values
// from the dense output data
func (o *Output) GetDenseS() (S []int) {
	return o.DenseS[:o.DenseIdx]
}

// GetDenseX returns all x values
// from the dense output data
func (o *Output) GetDenseX() (X []float64) {
	return o.DenseX[:o.DenseIdx]
}

// GetDenseY extracts the y[i] values for all output times
// from the dense output data
//  i -- index of y component
//  use to plot time series; e.g.:
//     plt.Plot(o.GetDenseX(), o.GetDenseY(0), &plt.A{L:"y0"})
func (o *Output) GetDenseY(i int) (Y []float64) {
	if o.DenseIdx > 0 {
		Y = make([]float64, o.DenseIdx)
		for j := 0; j < o.DenseIdx; j++ {
			Y[j] = o.DenseY[j][i]
		}
	}
	return
}

// GetDenseYtable returns a table with all y values such that Y[idxOut][dim]
// from the dense output data
func (o *Output) GetDenseYtable() (Y [][]float64) {
	if len(o.DenseY) < 1 {
		return
	}
	ndim := len(o.DenseY[0])
	Y = utl.Alloc(o.DenseIdx, ndim)
	for j := 0; j < o.DenseIdx; j++ {
		for i := 0; i < ndim; i++ {
			Y[j][i] = o.DenseY[j][i]
		}
	}
	return
}

// GetDenseYtableT returns a (transposed) table with all y values such that Y[dim][idxOut]
// from the dense output data
func (o *Output) GetDenseYtableT() (Y [][]float64) {
	if len(o.DenseY) < 1 {
		return
	}
	ndim := len(o.DenseY[0])
	Y = utl.Alloc(ndim, o.DenseIdx)
	for j := 0; j < o.DenseIdx; j++ {
		for i := 0; i < ndim; i++ {
			Y[i][j] = o.DenseY[j][i]
		}
	}
	return
}
