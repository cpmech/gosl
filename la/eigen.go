// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import "github.com/cpmech/gosl/la/oblas"

// Eigenvalues computes eigenvalues of general matrix
//
//   A ⋅ v[j] = λ[j] ⋅ v[j]
//
//   INPUT:
//     a -- general matrix
//
//   OUTPUT:
//     lam -- eigenvalues [pre-allocated]
//
func Eigenvalues(lam VectorC, A *Matrix, preserveA bool) (err error) {
	a := A
	if preserveA {
		a = A.GetCopy()
	}
	wr, wi := make([]float64, a.M), make([]float64, a.M)
	err = oblas.Dgeev(false, false, a.M, a.Data, a.M, wr, wi, nil, 0, nil, 0)
	oblas.JoinComplex(lam, wr, wi)
	return
}

// Eigenvectors computes eigenvalues and RIGHT eigenvectors of general matrix
//
//   A ⋅ v[j] = λ[j] ⋅ v[j]
//
//   INPUT:
//     a -- general matrix
//
//   OUTPUT:
//     v   -- matrix with the eigenvectors; each column contains one eigenvector [pre-allocated]
//     lam -- eigenvalues [pre-allocated]
//
func Eigenvectors(v *MatrixC, lam VectorC, A *Matrix, preserveA bool) (err error) {
	a := A
	if preserveA {
		a = A.GetCopy()
	}
	wr, wi := make([]float64, a.M), make([]float64, a.M)
	vr := make([]float64, a.M*a.M)
	err = oblas.Dgeev(false, true, a.M, a.Data, a.M, wr, wi, nil, 0, vr, a.M)
	oblas.JoinComplex(lam, wr, wi)
	oblas.EigenvecsBuild(v.Data, wr, wi, vr)
	return
}
