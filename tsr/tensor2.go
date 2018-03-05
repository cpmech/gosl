// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tsr implements Tensors using the smart approach by considering the Mandel's basis
package tsr

import "github.com/cpmech/gosl/la"

// Tensor2 holds the (Orthonormal/Cartesian) components of a 2nd order tensor,
// symmetric or not
type Tensor2 struct {
	data      la.Vector
	symmetric bool
}

// NewTensor2 returns a new Tensor2 object
// NOTE: if !symmetric all components are used and twoD flag is ignored
func NewTensor2(symmetric, twoD bool) (o *Tensor2) {
	o = new(Tensor2)
	o.symmetric = symmetric
	numberOfComponents := 9
	if symmetric {
		if twoD {
			numberOfComponents = 4
		} else {
			numberOfComponents = 6
		}
	}
	o.data = la.NewVector(numberOfComponents)
	return
}

// Set sets component [i,j] with it's Orthonormal/Cartesian value
func (o *Tensor2) Set(i, j int, value float64) {
	var I int
	if o.symmetric {
		I = SecToManI[i][j]
		if I > 2 {
			value *= sq2
		}
	} else {
		I = SecToVecI[i][j]
	}
	if I >= len(o.data) {
		return // 2D tensor; i.e. other components are zero
	}
	o.data[I] = value
}

// Get returns the Orthonormal/Cartesian component [i,j]
func (o *Tensor2) Get(i, j int) (value float64) {
	if o.symmetric {
		I := SecToManI[i][j]
		if I >= len(o.data) {
			return 0 // 2D tensor; i.e. other components are zero
		}
		if I > 2 {
			return o.data[I] / sq2
		}
		return o.data[I]
	}
	I := SecToVecI[i][j]
	return o.data[I]
}
