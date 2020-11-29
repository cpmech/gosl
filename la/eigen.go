// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la/oblas"
)

// EigenVal computes eigenvalues of general matrix
//
//   A ⋅ v[j] = λ[j] ⋅ v[j]
//
//   INPUT:
//     a -- general matrix
//
//   OUTPUT:
//     w -- eigenvalues [pre-allocated]
//
func EigenVal(w VectorC, A *Matrix, preserveA bool) {
	a := A
	if preserveA {
		a = A.GetCopy()
	}
	wr, wi := make([]float64, a.M), make([]float64, a.M)
	oblas.Dgeev(false, false, a.M, a.Data, a.M, wr, wi, nil, 0, nil, 0)
	oblas.JoinComplex(w, wr, wi)
}

// EigenVecL computes eigenvalues and LEFT eigenvectors of general matrix
//
//    H                  H
//   u [j] ⋅ A = λ[j] ⋅ u [j]    LEFT eigenvectors
//
//   INPUT:
//     a -- general matrix
//
//   OUTPUT:
//     u -- matrix with the eigenvectors; each column contains one eigenvector [pre-allocated]
//     w -- eigenvalues [pre-allocated]
//
func EigenVecL(u *MatrixC, w VectorC, A *Matrix, preserveA bool) {
	a := A
	if preserveA {
		a = A.GetCopy()
	}
	wr, wi := make([]float64, a.M), make([]float64, a.M)
	vl := make([]float64, a.M*a.M)
	oblas.Dgeev(true, false, a.M, a.Data, a.M, wr, wi, vl, a.M, nil, 0)
	oblas.JoinComplex(w, wr, wi)
	oblas.EigenvecsBuild(u.Data, wr, wi, vl)
}

// EigenVecR computes eigenvalues and RIGHT eigenvectors of general matrix
//
//   A ⋅ v[j] = λ[j] ⋅ v[j]
//
//   INPUT:
//     a -- general matrix
//
//   OUTPUT:
//     v -- matrix with the eigenvectors; each column contains one eigenvector [pre-allocated]
//     w -- eigenvalues [pre-allocated]
//
func EigenVecR(v *MatrixC, w VectorC, A *Matrix, preserveA bool) {
	a := A
	if preserveA {
		a = A.GetCopy()
	}
	wr, wi := make([]float64, a.M), make([]float64, a.M)
	vr := make([]float64, a.M*a.M)
	oblas.Dgeev(false, true, a.M, a.Data, a.M, wr, wi, nil, 0, vr, a.M)
	oblas.JoinComplex(w, wr, wi)
	oblas.EigenvecsBuild(v.Data, wr, wi, vr)
}

// EigenVecLR computes eigenvalues and LEFT and RIGHT eigenvectors of general matrix
//
//   A ⋅ v[j] = λ[j] ⋅ v[j]      RIGHT eigenvectors
//
//    H                  H
//   u [j] ⋅ A = λ[j] ⋅ u [j]    LEFT eigenvectors
//
//   INPUT:
//     a -- general matrix
//
//   OUTPUT:
//     u -- matrix with the LEFT eigenvectors; each column contains one eigenvector [pre-allocated]
//     v -- matrix with the RIGHT eigenvectors; each column contains one eigenvector [pre-allocated]
//     w -- λ eigenvalues [pre-allocated]
//
func EigenVecLR(u, v *MatrixC, w VectorC, A *Matrix, preserveA bool) {
	a := A
	if preserveA {
		a = A.GetCopy()
	}
	wr, wi := make([]float64, a.M), make([]float64, a.M)
	uu := make([]float64, a.M*a.M)
	vv := make([]float64, a.M*a.M)
	oblas.Dgeev(true, true, a.M, a.Data, a.M, wr, wi, uu, a.M, vv, a.M)
	oblas.JoinComplex(w, wr, wi)
	oblas.EigenvecsBuildBoth(u.Data, v.Data, wr, wi, uu, vv)
}

// CheckEigenVecL checks left eigenvector:
//
//    H                  H
//   u [j] ⋅ A = λ[j] ⋅ u [j]    LEFT eigenvectors
//
func CheckEigenVecL(tst *testing.T, A *Matrix, λ VectorC, u *MatrixC, tol float64) {
	Ac := A.GetComplex()
	res := NewVectorC(A.M)
	λu := NewVectorC(A.M)
	for i := 0; i < A.M; i++ {
		ui := u.GetCol(i)
		λu.Apply(λ[i], ui)
		MatVecMulC(res, 1, Ac, ui)
		chk.ArrayC(tst, io.Sf("λ[%d]⋅u[%d]", i, i), tol, res, λu)
	}
}

// CheckEigenVecR checks right eigenvector:
//
//   A ⋅ v[j] = λ[j] ⋅ v[j]      RIGHT eigenvectors
//
func CheckEigenVecR(tst *testing.T, A *Matrix, λ VectorC, v *MatrixC, tol float64) {
	Ac := A.GetComplex()
	res := NewVectorC(A.M)
	λv := NewVectorC(A.M)
	for i := 0; i < A.M; i++ {
		vi := v.GetCol(i)
		λv.Apply(λ[i], vi)
		MatVecMulC(res, 1, Ac, vi)
		chk.ArrayC(tst, io.Sf("λ[%d]⋅v[%d]", i, i), tol, res, λv)
	}
}
