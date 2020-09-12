// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"

	"gosl/fun"
	"gosl/la"
	"gosl/utl"
)

// History holds history of optmization using directiors; e.g. for Debugging
type History struct {

	// data
	Ndim  int         // dimension of x-vector
	HistX []la.Vector // [it] history of x-values (position)
	HistU []la.Vector // [it] history of u-values (direction)
	HistF []float64   // [it] history of f-values
	HistI []float64   // [it] index of iteration

	// configuration
	NptsI   int       // number of points for contour
	NptsJ   int       // number of points for contour
	RangeXi []float64 // {ximin, ximax} [may be nil for default]
	RangeXj []float64 // {xjmin, xjmax} [may be nil for default]
	GapXi   float64   // expand {ximin, ximax}
	GapXj   float64   // expand {ximin, ximax}

	// internal
	ffcn fun.Sv // f({x}) function
}

// NewHistory returns new object
func NewHistory(nMaxIt int, f0 float64, x0 la.Vector, ffcn fun.Sv) (o *History) {
	o = new(History)
	o.Ndim = len(x0)
	o.HistX = append(o.HistX, x0.GetCopy())
	o.HistU = append(o.HistU, nil)
	o.HistF = append(o.HistF, f0)
	o.HistI = append(o.HistI, 0)
	o.NptsI = 41
	o.NptsJ = 41
	o.ffcn = ffcn
	return
}

// Append appends new x and u vectors, and updates F and I arrays
func (o *History) Append(fx float64, x, u la.Vector) {
	o.HistX = append(o.HistX, x.GetCopy())
	o.HistU = append(o.HistU, u.GetCopy())
	o.HistF = append(o.HistF, fx)
	o.HistI = append(o.HistI, float64(len(o.HistI)))
}

// Limits computes range of X variables
func (o *History) Limits() (Xmin []float64, Xmax []float64) {
	Xmin = make([]float64, o.Ndim)
	Xmax = make([]float64, o.Ndim)
	for j := 0; j < o.Ndim; j++ {
		Xmin[j] = math.MaxFloat64
		Xmax[j] = math.SmallestNonzeroFloat64
		for _, x := range o.HistX {
			Xmin[j] = utl.Min(Xmin[j], x[j])
			Xmax[j] = utl.Max(Xmax[j], x[j])
		}
	}
	return
}
