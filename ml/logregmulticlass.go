// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// LogRegMultiClass implements a logistic regression model for multiple classes (Observer of data)
type LogRegMultiClass struct {
	name   string       // name of this "observer"
	data   *Data        // X-y data
	nClass int          // number of classes
	params []*ParamsReg // [nClass] parameters for each class
	dataB  []*Data      // [nClass] reference to X-y data but with binary y-vector
	models []*LogReg    // [nClass] one-versus-all models
}

// NewLogRegMultiClass returns a new object
// NOTE: the y-vector in data must have values in [0, nClass-1]
func NewLogRegMultiClass(data *Data, name string) (o *LogRegMultiClass) {
	o = new(LogRegMultiClass)
	o.name = name
	o.data = data
	o.data.AddObserver(o)
	o.nClass = int(data.Y.Max()) + 1
	o.params = make([]*ParamsReg, o.nClass)
	o.dataB = make([]*Data, o.nClass)
	o.models = make([]*LogReg, o.nClass)
	useY := true
	allocate := false
	for k := 0; k < o.nClass; k++ {
		o.params[k] = NewParamsReg(data.Nfeatures)
		o.dataB[k] = NewData(data.Nsamples, data.Nfeatures, useY, allocate)
	}
	o.Update()
	return
}

// Name returns the name of this "Observer"
func (o *LogRegMultiClass) Name() string {
	return o.name
}

// Update perform updates after data has been changed (as an Observer)
func (o *LogRegMultiClass) Update() {
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
			o.models[k] = NewLogReg(o.dataB[k], o.params[k], io.Sf("%s_%dclass", o.name, k))
		} else {
			o.models[k].Update()
		}
	}
}

// SetLambda sets the regularization parameter
func (o *LogRegMultiClass) SetLambda(lambda float64) {
	for k := 0; k < o.nClass; k++ {
		o.models[k].params.SetLambda(lambda)
	}
}

// Predict returns the model evaluation @ {x;θ,b}
//   Input:
//     x -- vector of features
//   Output:
//     class -- class with the highest probability
//     probs -- probabilities
func (o *LogRegMultiClass) Predict(x la.Vector) (class int, probs []float64) {
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

// Train finds the parameters using closed-form solutions
//  gradDesc -- use Gradient-Descent
func (o *LogRegMultiClass) Train(gardDesc bool) {
	for k := 0; k < o.nClass; k++ {
		o.models[k].Train()
		io.Pforan("%v  %v\n", o.params[k].GetBias(), o.params[k].GetThetas())
	}
}
