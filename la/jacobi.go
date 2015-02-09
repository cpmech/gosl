// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"code.google.com/p/gosl/utl"
)

const (
	JACOBI_TOL    = 1e-15
	JACOBI_EPS    = 1e-16
	JACOBI_NITMAX = 20
)

// Jacobi performs the Jacobi transformation of a symmetric matrix to find its eigenvectors and
// eigenvalues. Note that for matrices of order greater than about 10, say, the algorithm is slower,
// by a significant constant factor, than the QR method.   A = Q * L * Q.T
//  Arguments:
//   A In/Out -- matrix to compute eigenvalues (SYMMETRIC and SQUARE)
//   Q   Out  -- matrix which columns are the eigenvectors
//   v   Out  -- vector with the eigenvalues
//   nit Out  -- the number of iterations
func Jacobi(Q [][]float64, v []float64, A [][]float64) (nit int, err error) {

	/*
	   The Jacobi method consists of a sequence of orthogonal similarity transformations. Each
	   transformation (a Jacobi rotation) is just a plane rotation designed to annihilate one of the
	   off-diagonal matrix elements. Successive transformations undo previously set zeros, but the
	   off-diagonal elements nevertheless get smaller and smaller. Accumulating the product of the
	   transformations as you go gives the matrix of eigenvectors (Q), while the elements of the final
	   diagonal matrix (A) are the eigenvalues. The Jacobi method is absolutely foolproof for all real
	   symmetric matrices.
	*/

	var j, p, q int
	var θ, τ, t, sm, s, h, g, c float64

	n := len(A)
	b := make([]float64, n)
	z := make([]float64, n) // this vector will accumulate terms of the form tapq as in equation (11.1.14).

	// initialize Q to the identity matrix.
	for p = 0; p < n; p++ {
		for q = 0; q < n; q++ {
			Q[p][q] = 0.0
		}
		Q[p][p] = 1.0
	}

	// initialize b and v to the diagonal of A
	for p = 0; p < n; p++ {
		b[p] = A[p][p]
		v[p] = A[p][p]
		z[p] = 0.0
	}

	// perform iterations
	for it := 0; it < JACOBI_NITMAX; it++ {

		// sum off-diagonal elements.
		sm = 0.0
		for p = 0; p < n-1; p++ {
			for q = p + 1; q < n; q++ {
				sm += math.Abs(A[p][q])
			}
		}

		// exit point
		if sm < JACOBI_TOL {
			return it + 1, nil
		}

		// rotations
		for p = 0; p < n-1; p++ {
			for q = p + 1; q < n; q++ {
				h = v[q] - v[p]
				if math.Abs(h) <= JACOBI_TOL {
					t = 1.0
				} else {
					θ = 0.5 * h / (A[p][q])
					t = 1.0 / (math.Abs(θ) + math.Sqrt(1.0+θ*θ))
					if θ < 0.0 {
						t = -t
					}
				}
				c = 1.0 / math.Sqrt(1.0+t*t)
				s = t * c
				τ = s / (1.0 + c)
				h = t * A[p][q]
				z[p] -= h
				z[q] += h
				v[p] -= h
				v[q] += h
				A[p][q] = 0.0
				for j = 0; j < p; j++ { // case of rotations 0 <= j < p.
					g = A[j][p]
					h = A[j][q]
					A[j][p] = g - s*(h+g*τ)
					A[j][q] = h + s*(g-h*τ)
				}
				for j = p + 1; j < q; j++ { // case of rotations p < j < q.
					g = A[p][j]
					h = A[j][q]
					A[p][j] = g - s*(h+g*τ)
					A[j][q] = h + s*(g-h*τ)
				}
				for j = q + 1; j < n; j++ { //case of rotations q < j < n.
					g = A[p][j]
					h = A[q][j]
					A[p][j] = g - s*(h+g*τ)
					A[q][j] = h + s*(g-h*τ)
				}
				for j = 0; j < n; j++ { // Q matrix
					g = Q[j][p]
					h = Q[j][q]
					Q[j][p] = g - s*(h+g*τ)
					Q[j][q] = h + s*(g-h*τ)
				}
			}
		}
		for p = 0; p < n; p++ {
			b[p] += z[p]
			z[p] = 0.0  // reinitialize z.
			v[p] = b[p] // update v with the sum of tapq,
		}
	}

	err = utl.Err(_jacobi_err1, JACOBI_NITMAX+1)
	return
}

// error messages
var (
	_jacobi_err1 = "jacobi.go: Jacobi rotation dit not converge after %d iterations"
)
