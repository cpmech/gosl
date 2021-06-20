// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"github.com/cpmech/gosl/la/oblas"
)

// VecRmsError returns the scaled root-mean-square of the difference between two vectors
// with components normalised by a scaling factor
//                __________________________
//               /     ————              2
//              /  1   \    /  error[i]  \
//   rms =  \  /  ———  /    | —————————— |
//           \/    N   ———— \  scale[i]  /
//
//   error[i] = |u[i] - v[i]|
//
//   scale[i] = a + m*|s[i]|
//
func VecRmsError(u, v Vector, a, m float64, s Vector) (rms float64) {
	var scale, err float64
	for i := 0; i < len(u); i++ {
		scale = a + m*math.Abs(s[i])
		err = math.Abs(u[i] - v[i])
		rms += err * err / (scale * scale)
	}
	return math.Sqrt(rms / float64(len(u)))
}

// VecDot returns the dot product between two vectors:
//   s := u・v
func VecDot(u, v Vector) (res float64) {
	cutoff := 150
	if len(u) <= cutoff {
		for i := 0; i < len(u); i++ {
			res += u[i] * v[i]
		}
		return
	}
	return oblas.Ddot(len(u), u, 1, v, 1)
}

// VecAdd adds the scaled components of two vectors
//   res := α⋅u + β⋅v   ⇒   result[i] := α⋅u[i] + β⋅v[i]
func VecAdd(res Vector, α float64, u Vector, β float64, v Vector) {
	n := len(u)
	cutoff := 150
	if β == 1 && n > cutoff {
		copy(res, v)
		oblas.Daxpy(n, α, u, 1, res, 1)
		return
	}
	m := n % 4
	for i := 0; i < m; i++ {
		res[i] = α*u[i] + β*v[i]
	}
	for i := m; i < n; i += 4 {
		res[i+0] = α*u[i+0] + β*v[i+0]
		res[i+1] = α*u[i+1] + β*v[i+1]
		res[i+2] = α*u[i+2] + β*v[i+2]
		res[i+3] = α*u[i+3] + β*v[i+3]
	}
}

// VecMaxDiff returns the maximum absolute difference between two vectors
//   maxdiff = max(|u - v|)
func VecMaxDiff(u, v Vector) (maxdiff float64) {
	maxdiff = math.Abs(u[0] - v[0])
	for i := 1; i < len(u); i++ {
		diff := math.Abs(u[i] - v[i])
		if diff > maxdiff {
			maxdiff = diff
		}
	}
	return
}

// VecScaleAbs creates a "scale" vector using the absolute value of another vector
//   scale := a + m ⋅ |x|     ⇒      scale[i] := a + m ⋅ |x[i]|
func VecScaleAbs(scale Vector, a, m float64, x Vector) {
	for i := 0; i < len(x); i++ {
		scale[i] = a + m*math.Abs(x[i])
	}
}

// complex /////////////////////////////////////////////////////////////////////////////////////////
