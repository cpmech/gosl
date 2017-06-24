// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la/oblas"
)

// MatVecMul returns the matrix-vector multiplication
//
//   v = α⋅a⋅u    ⇒    vi = α * aij * uj
//
func MatVecMul(v Vector, α float64, a *Matrix, u Vector) {
	err := oblas.Dgemv(false, a.M, a.N, α, a.Data, a.M, u, 1, 0.0, v, 1)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}

// MatVecMulAdd returns the matrix-vector multiplication with addition
//
//   v += α⋅a⋅u    ⇒    vi += α * aij * uj
//
func MatVecMulAdd(v Vector, α float64, a *Matrix, u Vector) {
	err := oblas.Dgemv(false, a.M, a.N, α, a.Data, a.M, u, 1, 1.0, v, 1)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// MatVecMulC returns the matrix-vector multiplication (complex version)
//
//   v = α⋅a⋅u    ⇒    vi = α * aij * uj
//
func MatVecMulC(v VectorC, α complex128, a *MatrixC, u VectorC) {
	err := oblas.Zgemv(false, a.M, a.N, α, a.Data, a.M, u, 1, 0.0, v, 1)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}

// MatVecMulAddC returns the matrix-vector multiplication with addition (complex version)
//
//   v += α⋅a⋅u    ⇒    vi += α * aij * uj
//
func MatVecMulAddC(v VectorC, α complex128, a *MatrixC, u VectorC) {
	err := oblas.Zgemv(false, a.M, a.N, α, a.Data, a.M, u, 1, 1.0, v, 1)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}
