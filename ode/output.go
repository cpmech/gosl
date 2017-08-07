// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// Output holds output data
type Output struct {

	// control
	Fcn     OutF        // output function
	IdxSave int         // current index in Xvalues and Yvalues == last output
	Hvalues []float64   // h values if SaveXY is true [IdxSave]
	Xvalues []float64   // X values if SaveXY is true [IdxSave]
	Yvalues []la.Vector // Y values if SaveXY is true [IdxSave][ndim]

	// derived
	first   bool // first output
	nMaxOut int  // max number of output
}

// NewOutput returns a new structure
func NewOutput(fcn OutF) (o *Output) {
	o = new(Output)
	o.Fcn = fcn
	return
}

// Resize allocates memory
//   NmaxOut -- max number of output. use 0 for NO output
func (o *Output) Resize(NmaxOut int) {
	o.nMaxOut = NmaxOut
	o.IdxSave = 0
	if o.nMaxOut > 0 {
		o.Hvalues = make([]float64, o.nMaxOut)
		o.Xvalues = make([]float64, o.nMaxOut)
		o.Yvalues = make([]la.Vector, o.nMaxOut)
	}
}

// Execute executes output; e.g. call Fcn and saves x and y values
func (o *Output) Execute(h, x float64, y []float64) {
	if o.Fcn != nil {
		o.Fcn(o.first, h, x, y)
		if o.first {
			o.first = false
		}
	}
	if o.IdxSave < o.nMaxOut {
		o.Hvalues[o.IdxSave] = h
		o.Xvalues[o.IdxSave] = x
		o.Yvalues[o.IdxSave] = la.NewVector(len(y))
		o.Yvalues[o.IdxSave].Apply(1, y)
		o.IdxSave++
	} else if o.nMaxOut > 0 { // allocate more space
		io.Pf(". . . allocating more space for output . . . \n")
		factor := 2
		htmp := make([]float64, o.nMaxOut*factor)
		xtmp := make([]float64, o.nMaxOut*factor)
		ytmp := make([]la.Vector, o.nMaxOut*factor)
		copy(htmp, o.Hvalues[:o.IdxSave])
		copy(xtmp, o.Xvalues[:o.IdxSave])
		for i := 0; i < o.IdxSave; i++ {
			htmp[i] = o.Hvalues[i]
			xtmp[i] = o.Xvalues[i]
			ytmp[i] = la.NewVector(len(y))
			ytmp[i].Apply(1, o.Yvalues[i])
		}
		o.Hvalues = htmp
		o.Xvalues = xtmp
		o.Yvalues = ytmp
		o.nMaxOut *= factor
	}
}

// ExtractTimeSeries extracts the y[i] values for all output times
//  i -- index of y component
//  use to plot time series; e.g.:
//     plt.Plot(o.Xvalues[:o.IdxSave], o.ExtractTimeSeries(0), &plt.A{L:"y0"})
func (o *Output) ExtractTimeSeries(i int) (Yi []float64) {
	if o.IdxSave > 0 {
		Yi = make([]float64, o.IdxSave)
		for j := 0; j < o.IdxSave; j++ {
			Yi[j] = o.Yvalues[j][i]
		}
	}
	return
}
