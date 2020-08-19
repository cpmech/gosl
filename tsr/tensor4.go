// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import "gosl/la"

// Tensor4 holds the (Orthonormal/Cartesian) components of a 4th order tensor,
// full-symmetric or not.
// NOTE: "full-symmetric" hear means major and minor symmetry, where:
//       A[i][j][k][l] = A[j][i][k][l] = A[i][j][l][k]
//       i.e. [i][j] and [k][l] can be swapped
type Tensor4 struct {
	data      *la.Matrix
	symmetric bool
}

// NewTensor4 returns a new Tensor4 object
// NOTE: if !symmetric all components are used and twoD flag is ignored
func NewTensor4(symmetric, twoD bool) (o *Tensor4) {
	o = new(Tensor4)
	o.symmetric = symmetric
	numberOfComponents := 9
	if symmetric {
		if twoD {
			numberOfComponents = 4
		} else {
			numberOfComponents = 6
		}
	}
	o.data = la.NewMatrix(numberOfComponents, numberOfComponents)
	return
}

// Set sets component [i,j] with it's Orthonormal/Cartesian value
func (o *Tensor4) Set(i, j, k, l int, value float64) {
	var I, J int
	if o.symmetric {
		I = FouToManI[i][j][k][l]
		J = FouToManJ[i][j][k][l]
		if I > 2 && J < 3 {
			value *= sq2
		}
		if J > 2 && I < 3 {
			value *= sq2
		}
		if I > 2 && J > 2 {
			value *= 2.0
		}
	} else {
		I = FouToVecI[i][j][k][l]
		J = FouToVecJ[i][j][k][l]
	}
	o.data.Set(I, J, value)
}

// Get returns the Orthonormal/Cartesian component [i,j]
func (o *Tensor4) Get(i, j, k, l int) (value float64) {
	if o.symmetric {
		I := FouToManI[i][j][k][l]
		J := FouToManJ[i][j][k][l]
		if I > 2 && J < 3 {
			return o.data.Get(I, J) / sq2
		}
		if J > 2 && I < 3 {
			return o.data.Get(I, J) / sq2
		}
		if I > 2 && J > 2 {
			return o.data.Get(I, J) / 2.0
		}
		return o.data.Get(I, J)
	}
	I := FouToVecI[i][j][k][l]
	J := FouToVecJ[i][j][k][l]
	return o.data.Get(I, J)
}
