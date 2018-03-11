// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ml implements Machine Learning algorithms
package ml

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Stat holds statistics about data
type Stat struct {
	UseY  bool      // use y values
	MinX  []float64 // [nFeatures] min x values
	MaxX  []float64 // [nFeatures] max x values
	SumX  []float64 // [nFeatures] sum of x values
	MeanX []float64 // [nFeatures] mean of x values
	SigX  []float64 // [nFeatures] standard deviations of x
	DelX  []float64 // [nFeatures] difference: max(x) - min(x)
	MinY  float64   // min of y values
	MaxY  float64   // max of y values
	SumY  float64   // sum of y values
	MeanY float64   // mean of y values
	SigY  float64   // standard deviation of y
	DelY  float64   // difference: max(y) - min(y)
}

// NewStat returns a new Stat object
func NewStat(nFeatures int, useY bool) (o *Stat) {
	o = new(Stat)
	o.UseY = useY
	o.MinX = make([]float64, nFeatures)
	o.MaxX = make([]float64, nFeatures)
	o.SumX = make([]float64, nFeatures)
	o.MeanX = make([]float64, nFeatures)
	o.SigX = make([]float64, nFeatures)
	o.DelX = make([]float64, nFeatures)
	return
}

// Compute compute statistics for given data
//   X -- [nSamples][nFeatures] x values
//   y -- [nSamples] y values (may be nil)
func (o *Stat) Compute(X *la.Matrix, y la.Vector) {

	// constants
	m := X.M // number of samples
	n := X.N // number of features
	if n != len(o.MinX) {
		chk.Panic("number of columns in X matrix does not correspond with the number of features. %d != %d\n", n, len(o.MinX))
	}

	// x values
	mf := float64(m)
	for j := 0; j < n; j++ {
		o.MinX[j] = X.Get(0, j)
		o.MaxX[j] = o.MinX[j]
		o.SumX[j] = 0.0
		for i := 0; i < m; i++ {
			xval := X.Get(i, j)
			o.MinX[j] = utl.Min(o.MinX[j], xval)
			o.MaxX[j] = utl.Max(o.MaxX[j], xval)
			o.SumX[j] += xval
		}
		o.MeanX[j] = o.SumX[j] / mf
		o.SigX[j] = rnd.StatDevFirst(X.Col(j), o.MeanX[j], true)
		o.DelX[j] = o.MaxX[j] - o.MinX[j]
	}

	// y values
	if o.UseY {
		o.MinY = y[0]
		o.MaxY = o.MinY
		o.SumY = 0.0
		for i := 0; i < m; i++ {
			o.MinY = utl.Min(o.MinY, y[i])
			o.MaxY = utl.Max(o.MaxY, y[i])
			o.SumY += y[i]
		}
		o.MeanY = o.SumY / mf
		o.SigY = rnd.StatDevFirst(y, o.MeanY, true)
		o.DelY = o.MaxY - o.MinY
	}
}

// SumVars computes the sums along the columns of X and y
//   Output:
//     t -- scalar t = oᵀy  sum of columns of the y vector: t = Σ_i^m o_i y_i
//     s -- vector s = Xᵀo  sum of columns of the X matrix: s_j = Σ_i^m o_i X_ij  [nFeatures]
func (o *Stat) SumVars(X *la.Matrix, y la.Vector) (s la.Vector, t float64) {
	one := la.NewVector(X.M)
	one.Fill(1.0)
	t = la.VecDot(one, y)
	s = la.NewVector(X.N)
	la.MatTrVecMul(s, 1, X, one)
	return
}
