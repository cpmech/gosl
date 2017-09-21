// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la/oblas"
)

// DenSolve solves dense linear system using LAPACK (OpemBLAS)
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
func DenSolve(x Vector, A *Matrix, b Vector, preserveA bool) {
	a := A
	if preserveA {
		a = NewMatrix(A.M, A.N)
		copy(a.Data, A.Data)
	}
	copy(x, b)
	ipiv := make([]int32, A.M)
	err := oblas.Dgesv(A.M, 1, a.Data, A.M, ipiv, x, A.M)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}

// Cholesky returns the Cholesky decomposition of a symmetric positive-definite matrix
//
//   a = L * trans(L)
//
func Cholesky(L, a *Matrix) {
	for j := 0; j < a.M; j++ { // loop over columns
		for i := j; i < a.M; i++ { // loop over lower diagonal rows (including diagonal)
			amsum := a.Get(i, j)
			for k := 0; k < j; k++ {
				amsum -= L.Get(i, k) * L.Get(j, k)
			}
			if i == j {
				if amsum <= 0.0 {
					chk.Panic("Cholesky factorization failed due to non positive-definite matrix")
				}
				L.Set(i, j, math.Sqrt(amsum))
			} else {
				L.Set(i, j, amsum/L.Get(j, j))
			}
		}
	}
}

// SolveRealLinSysSPD solves a linear system with real numbres and a Symmetric-Positive-Definite (SPD) matrix
//
//        x := inv(a) * b
//
//   NOTE: this function uses Cholesky decomposition and should be used for small systems
func SolveRealLinSysSPD(x Vector, a *Matrix, b Vector) {

	// Cholesky factorisation
	L := NewMatrix(a.M, a.M)
	Cholesky(L, a)

	// solve L*y = b storing y in x
	for i := 0; i < a.M; i++ {
		bmsum := b[i]
		for k := 0; k < i; k++ {
			bmsum -= L.Get(i, k) * x[k]
		}
		x[i] = bmsum / L.Get(i, i)
	}

	// solve trans(L)*x = y with y==x
	for i := a.M - 1; i >= 0; i-- {
		bmsum := x[i]
		for k := i + 1; k < a.M; k++ {
			bmsum -= L.Get(k, i) * x[k]
		}
		x[i] = bmsum / L.Get(i, i)
	}
}

// SolveTwoRealLinSysSPD solves two linear systems with real numbres and Symmetric-Positive-Definite (SPD) matrices
//
//        x := inv(a) * b    and    X := inv(a) * B
//
//   NOTE: this function uses Cholesky decomposition and should be used for small systems
func SolveTwoRealLinSysSPD(x, X Vector, a *Matrix, b, B Vector) {

	// Cholesky factorisation
	L := NewMatrix(a.M, a.M)
	Cholesky(L, a)

	// solve L*y = b storing y in x
	for i := 0; i < a.M; i++ {
		bmsum := b[i]
		Bmsum := B[i]
		for k := 0; k < i; k++ {
			bmsum -= L.Get(i, k) * x[k]
			Bmsum -= L.Get(i, k) * X[k]
		}
		x[i] = bmsum / L.Get(i, i)
		X[i] = Bmsum / L.Get(i, i)
	}

	// solve trans(L)*x = y with y==x
	for i := a.M - 1; i >= 0; i-- {
		bmsum := x[i]
		Bmsum := X[i]
		for k := i + 1; k < a.M; k++ {
			bmsum -= L.Get(k, i) * x[k]
			Bmsum -= L.Get(k, i) * X[k]
		}
		x[i] = bmsum / L.Get(i, i)
		X[i] = Bmsum / L.Get(i, i)
	}
}
