// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/opt"
)

// LogRegMulti implements a logistic regression model for multiple classes (Observer of data)
type LogRegMulti struct {

	// input
	data *Data // X-y data

	// access
	nClass int // number of classes

	// internal
	dataB  []*Data   // [nClass] reference to X-y data but with binary y-vector
	models []*LogReg // [nClass] one-versus-all models
}

// NewLogRegMulti returns a new object
// NOTE: the y-vector in data must have values in [0, nClass-1]
func NewLogRegMulti(data *Data) (o *LogRegMulti) {
	o = new(LogRegMulti)
	o.data = data
	o.data.AddObserver(o)
	o.nClass = int(data.Y.Max()) + 1
	o.dataB = make([]*Data, o.nClass)
	o.models = make([]*LogReg, o.nClass)
	useY := true
	allocate := false
	for k := 0; k < o.nClass; k++ {
		o.dataB[k] = NewData(data.Nsamples, data.Nfeatures, useY, allocate)
	}
	o.Update()
	return
}

// Update perform updates after data has been changed (as an Observer)
func (o *LogRegMulti) Update() {
	for k := 0; k < o.nClass; k++ {
		o.dataB[k].X = o.data.X
		if len(o.dataB[k].Y) != o.data.Nsamples {
			o.dataB[k].Y = la.NewVector(o.data.Nsamples)
		}
		for i := 0; i < o.data.Nsamples; i++ {
			if int(o.data.Y[i]) == k {
				o.dataB[k].Y[i] = 1.0
			} else {
				o.dataB[k].Y[i] = 0.0
			}
		}
		if o.models[k] == nil {
			o.models[k] = NewLogReg(o.dataB[k])
		} else {
			o.models[k].Update()
		}
	}
}

// SetLambda sets the regularization parameter
func (o *LogRegMulti) SetLambda(lambda float64) {
	for k := 0; k < o.nClass; k++ {
		o.models[k].SetLambda(lambda)
	}
}

// Predict returns the model evaluation @ {x;θ,b}
//   Input:
//     x -- vector of features
//   Output:
//     class -- class with the highest probability
//     probs -- probabilities
func (o *LogRegMulti) Predict(x la.Vector) (class int, probs []float64) {
	pMax := 0.0
	probs = make([]float64, o.nClass)
	for k := 0; k < o.nClass; k++ {
		probs[k] = o.models[k].Predict(x)
		if probs[k] > pMax {
			pMax = probs[k]
			class = k
		}
	}
	return
}

// Train finds the parameters using Newton's method
func (o *LogRegMulti) Train() {
	for k := 0; k < o.nClass; k++ {
		o.models[k].Train()
	}
}

// TrainNumerical trains model using numerical optimizer
//   method -- method/kind of numerical solver. e.g. conjgrad, powel, graddesc
//   saveHist -- save history
//   control -- parameters to numerical solver. See package 'opt'
func (o *LogRegMulti) TrainNumerical(method string, saveHist bool, control dbf.Params) (minCosts []float64, hists []*opt.History) {
	minCosts = make([]float64, o.nClass)
	if saveHist {
		hists = make([]*opt.History, o.nClass)
	}
	for k := 0; k < o.nClass; k++ {
		θini := la.NewVector(o.data.Nfeatures)
		bini := 0.0
		c, h := o.models[k].TrainNumerical(θini, bini, method, saveHist, control)
		minCosts[k] = c
		if saveHist {
			hists[k] = h
		}
	}
	return
}
