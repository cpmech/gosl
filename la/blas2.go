// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/cpmech/gosl/la/oblas"
	"github.com/cpmech/gosl/utl"
)

// MatVecMul returns the matrix-vector multiplication
//
//   v = α⋅a⋅u    ⇒    vi = α * aij * uj
//
func MatVecMul(v Vector, α float64, a *Matrix, u Vector) {
	if a.M < 9 && a.N < 9 {
		for i := 0; i < a.M; i++ {
			v[i] = 0.0
			for j := 0; j < a.N; j++ {
				v[i] += α * a.Get(i, j) * u[j]
			}
		}
		return
	}
	oblas.Dgemv(false, a.M, a.N, α, a.Data, a.M, u, 1, 0.0, v, 1)
}

// MatTrVecMul returns the transpose(matrix)-vector multiplication
//
//   v = α⋅aᵀ⋅u    ⇒    vi = α * aji * uj = α * uj * aji
//
func MatTrVecMul(v Vector, α float64, a *Matrix, u Vector) {
	if a.M < 9 && a.N < 9 {
		for i := 0; i < a.N; i++ {
			v[i] = 0.0
			for j := 0; j < a.M; j++ {
				v[i] += α * a.Get(j, i) * u[j]
			}
		}
		return
	}
	oblas.Dgemv(true, a.M, a.N, α, a.Data, a.M, u, 1, 0.0, v, 1)
}

// VecVecTrMul returns the matrix = vector-transpose(vector) multiplication
// (e.g. dyadic product)
//
//   a = α⋅u⋅vᵀ    ⇒    aij = α * ui * vj
//
func VecVecTrMul(a *Matrix, α float64, u, v Vector) {
	if a.M < 9 && a.N < 9 {
		for i := 0; i < a.M; i++ {
			for j := 0; j < a.N; j++ {
				a.Set(i, j, α*u[i]*v[j])
			}
		}
		return
	}
	oblas.Dger(a.M, a.N, α, u, 1, v, 1, a.Data, utl.Imax(a.M, a.N))
}

// MatVecMulAdd returns the matrix-vector multiplication with addition
//
//   v += α⋅a⋅u    ⇒    vi += α * aij * uj
//
func MatVecMulAdd(v Vector, α float64, a *Matrix, u Vector) {
	oblas.Dgemv(false, a.M, a.N, α, a.Data, a.M, u, 1, 1.0, v, 1)
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// MatVecMulC returns the matrix-vector multiplication (complex version)
//
//   v = α⋅a⋅u    ⇒    vi = α * aij * uj
//
func MatVecMulC(v VectorC, α complex128, a *MatrixC, u VectorC) {
	oblas.Zgemv(false, a.M, a.N, α, a.Data, a.M, u, 1, 0.0, v, 1)
}

// MatVecMulAddC returns the matrix-vector multiplication with addition (complex version)
//
//   v += α⋅a⋅u    ⇒    vi += α * aij * uj
//
func MatVecMulAddC(v VectorC, α complex128, a *MatrixC, u VectorC) {
	oblas.Zgemv(false, a.M, a.N, α, a.Data, a.M, u, 1, 1.0, v, 1)
}
