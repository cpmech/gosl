// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"github.com/cpmech/gosl/utl"
)

// MatInv returns the matrix inverse of 'a' in 'ai', in addition to the determinant of 'a'
func MatInv(ai, a [][]float64, tol float64) (det float64, err error) {
	if len(a) < 1 {
		return 0, utl.Err(_matinv_err1)
	}
	m, n := len(a), len(a[0])
	switch {
	case m == 1 && n == 1:
		det = a[0][0]
		if math.Abs(det) < tol {
			return 0, utl.Err(_matinv_err2, m, n, det, tol)
		}
		ai[0][0] = 1.0 / det
	case m == 2 && n == 2:
		det = a[0][0]*a[1][1] - a[0][1]*a[1][0]
		if math.Abs(det) < tol {
			return 0, utl.Err(_matinv_err2, m, n, det, tol)
		}
		ai[0][0] = a[1][1] / det
		ai[0][1] = -a[0][1] / det
		ai[1][0] = -a[1][0] / det
		ai[1][1] = a[0][0] / det
	case m == 3 && n == 3:
		det = a[0][0]*(a[1][1]*a[2][2]-a[1][2]*a[2][1]) - a[0][1]*(a[1][0]*a[2][2]-a[1][2]*a[2][0]) + a[0][2]*(a[1][0]*a[2][1]-a[1][1]*a[2][0])
		if math.Abs(det) < tol {
			return 0, utl.Err(_matinv_err2, m, n, det, tol)
		}

		ai[0][0] = (a[1][1]*a[2][2] - a[1][2]*a[2][1]) / det
		ai[0][1] = (a[0][2]*a[2][1] - a[0][1]*a[2][2]) / det
		ai[0][2] = (a[0][1]*a[1][2] - a[0][2]*a[1][1]) / det

		ai[1][0] = (a[1][2]*a[2][0] - a[1][0]*a[2][2]) / det
		ai[1][1] = (a[0][0]*a[2][2] - a[0][2]*a[2][0]) / det
		ai[1][2] = (a[0][2]*a[1][0] - a[0][0]*a[1][2]) / det

		ai[2][0] = (a[1][0]*a[2][1] - a[1][1]*a[2][0]) / det
		ai[2][1] = (a[0][1]*a[2][0] - a[0][0]*a[2][1]) / det
		ai[2][2] = (a[0][0]*a[1][1] - a[0][1]*a[1][0]) / det
	default:
		return 0, utl.Err(_matinv_err3, m, n)
	}
	return
}

// error messages
const (
	_matinv_err1 = "cannot handle nil matrix. len(a) = 0"
	_matinv_err2 = "inverse of (%d x %d) matrix failed with zero determinant: |det=%g| < %g"
	_matinv_err3 = "inverse of (%d x %d) matrix is not available"
)
