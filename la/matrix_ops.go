// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la/oblas"
	"github.com/cpmech/gosl/utl"
)

// MatInvSmall computes the inverse of small matrices of size 1x1, 2x2, or 3x3.
// It also returns the determinant.
//   Input:
//     a   -- the matrix
//     tol -- tolerance to assume zero determinant
//   Output:
//     ai  -- the inverse matrix
//     det -- determinant of a
func MatInvSmall(ai, a *Matrix, tol float64) (det float64, err error) {
	switch {
	case a.M == 1 && a.N == 1:
		det = a.Get(0, 0)
		if math.Abs(det) < tol {
			return 0, chk.Err("inverse of (%dx%d) matrix failed with zero determinant: |det(a)|=%g < %g", a.M, a.N, det, tol)
		}
		ai.Set(0, 0, 1.0/det)

	case a.M == 2 && a.N == 2:
		det = a.Get(0, 0)*a.Get(1, 1) - a.Get(0, 1)*a.Get(1, 0)
		if math.Abs(det) < tol {
			return 0, chk.Err("inverse of (%dx%d) matrix failed with zero determinant: |det(a)|=%g < %g", a.M, a.N, det, tol)
		}
		ai.Set(0, 0, +a.Get(1, 1)/det)
		ai.Set(0, 1, -a.Get(0, 1)/det)
		ai.Set(1, 0, -a.Get(1, 0)/det)
		ai.Set(1, 1, +a.Get(0, 0)/det)

	case a.M == 3 && a.N == 3:
		det = a.Get(0, 0)*(a.Get(1, 1)*a.Get(2, 2)-a.Get(1, 2)*a.Get(2, 1)) - a.Get(0, 1)*(a.Get(1, 0)*a.Get(2, 2)-a.Get(1, 2)*a.Get(2, 0)) + a.Get(0, 2)*(a.Get(1, 0)*a.Get(2, 1)-a.Get(1, 1)*a.Get(2, 0))
		if math.Abs(det) < tol {
			return 0, chk.Err("inverse of (%dx%d) matrix failed with zero determinant: |det(a)|=%g < %g", a.M, a.N, det, tol)
		}

		ai.Set(0, 0, (a.Get(1, 1)*a.Get(2, 2)-a.Get(1, 2)*a.Get(2, 1))/det)
		ai.Set(0, 1, (a.Get(0, 2)*a.Get(2, 1)-a.Get(0, 1)*a.Get(2, 2))/det)
		ai.Set(0, 2, (a.Get(0, 1)*a.Get(1, 2)-a.Get(0, 2)*a.Get(1, 1))/det)

		ai.Set(1, 0, (a.Get(1, 2)*a.Get(2, 0)-a.Get(1, 0)*a.Get(2, 2))/det)
		ai.Set(1, 1, (a.Get(0, 0)*a.Get(2, 2)-a.Get(0, 2)*a.Get(2, 0))/det)
		ai.Set(1, 2, (a.Get(0, 2)*a.Get(1, 0)-a.Get(0, 0)*a.Get(1, 2))/det)

		ai.Set(2, 0, (a.Get(1, 0)*a.Get(2, 1)-a.Get(1, 1)*a.Get(2, 0))/det)
		ai.Set(2, 1, (a.Get(0, 1)*a.Get(2, 0)-a.Get(0, 0)*a.Get(2, 1))/det)
		ai.Set(2, 2, (a.Get(0, 0)*a.Get(1, 1)-a.Get(0, 1)*a.Get(1, 0))/det)

	default:
		return 0, chk.Err("cannot compute inverse of %dx%d matrix with this function\n", a.M, a.N)
	}
	return
}

// MatSvd performs the SVD decomposition
//   Input:
//     a     -- matrix a
//     copyA -- creates a copy of a; otherwise 'a' is modified
//   Output:
//     s  -- diagonal terms [must be pre-allocated] len(s) = imin(a.M, a.N)
//     u  -- left matrix [must be pre-allocated] u is (a.M x a.M)
//     vt -- transposed right matrix [must be pre-allocated] vt is (a.N x a.N)
func MatSvd(s []float64, u, vt, a *Matrix, copyA bool) {
	superb := make([]float64, utl.Imin(a.M, a.N))
	acpy := a
	if copyA {
		acpy = a.GetCopy()
	}
	err := oblas.Dgesvd('A', 'A', a.M, a.N, acpy.Data, a.M, s, u.Data, a.M, vt.Data, a.N, superb)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}

// MatInv computes the inverse of a general matrix (square or not). It also computes the
// pseudo-inverse if the matrix is not square.
//   Input:
//     a -- input matrix (M x N)
//   Output:
//     ai -- inverse matrix (N x M)
//     det -- determinant of matrix (ONLY if calcDet == true and the matrix is square)
//   NOTE: the dimension of the ai matrix must be N x M for the pseudo-inverse
func MatInv(ai, a *Matrix, calcDet bool) (det float64, err error) {

	// square inverse
	if a.M == a.N {
		copy(ai.Data, a.Data)
		ipiv := make([]int32, utl.Imin(a.M, a.N))
		err = oblas.Dgetrf(a.M, a.N, ai.Data, a.M, ipiv) // NOTE: ipiv are 1-based indices
		if err != nil {
			return
		}
		if calcDet {
			det = 1.0
			for i := 0; i < a.M; i++ {
				if ipiv[i]-1 == int32(i) { // NOTE: ipiv are 1-based indices
					det = +det * ai.Get(i, i)
				} else {
					det = -det * ai.Get(i, i)
				}
			}
		}
		err = oblas.Dgetri(a.N, ai.Data, a.M, ipiv)
		return
	}

	// singular value decomposition
	s := make([]float64, utl.Imin(a.M, a.N))
	u := NewMatrix(a.M, a.M)
	vt := NewMatrix(a.N, a.N)
	MatSvd(s, u, vt, a, true)

	// pseudo inverse
	tolS := 1e-8 // TODO: improve this tolerance with a better estimate
	for i := 0; i < a.N; i++ {
		for j := 0; j < a.M; j++ {
			ai.Set(i, j, 0)
			for k := 0; k < len(s); k++ {
				if s[k] > tolS {
					ai.Add(i, j, vt.Get(k, i)*u.Get(j, k)/s[k])
				}
			}
		}
	}
	return
}

// MatCondNum returns the condition number of a square matrix using the inverse of this matrix;
// thus it is not as efficient as it could be, e.g. by using the SV decomposition.
//  normtype -- Type of norm to use:
//    "F" or "" => Frobenius
//    "I"       => Infinite
func MatCondNum(a *Matrix, normtype string) (res float64, err error) {
	ai := NewMatrix(a.M, a.N)
	_, err = MatInv(ai, a, false)
	if err != nil {
		return
	}
	if normtype == "I" {
		res = a.NormInf() * ai.NormInf()
		return
	}
	res = a.NormFrob() * ai.NormFrob()
	return
}
