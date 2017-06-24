// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import "math"

// VecFill fills a vector with a single number s:
//  v := s*ones(len(v))  =>  vi = s
func VecFill(v []float64, s float64) {
	for i := 0; i < len(v); i++ {
		v[i] = s
	}
}

// VecFillC fills a complex vector with a single number s:
//  v := s*ones(len(v))  =>  vi = s
func VecFillC(v []complex128, s complex128) {
	for i := 0; i < len(v); i++ {
		v[i] = s
	}
}

// VecApplyFunc runs a function over all components of a vector
//  vi = f(i,vi)
func VecApplyFunc(v []float64, f func(i int, x float64) float64) {
	for i := 0; i < len(v); i++ {
		v[i] = f(i, v[i])
	}
}

// VecGetMapped returns a new vector after applying a function over all of its components
//  new: vi = f(i)
func VecGetMapped(dim int, f func(i int) float64) (v []float64) {
	v = make([]float64, dim)
	for i := 0; i < len(v); i++ {
		v[i] = f(i)
	}
	return
}

// VecClone allocates a clone of a vector
//  b := a
func VecClone(a []float64) (b []float64) {
	b = make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		b[i] = a[i]
	}
	return
}

// VecAccum sum/accumulates all components in a vector
//  sum := Σ_i v[i]
func VecAccum(v []float64) (sum float64) {
	for i := 0; i < len(v); i++ {
		sum += v[i]
	}
	return
}

// VecNorm returns the Euclidian norm of a vector:
//  nrm := ||v||
func VecNorm(v []float64) (nrm float64) {
	for i := 0; i < len(v); i++ {
		nrm += v[i] * v[i]
	}
	nrm = math.Sqrt(nrm)
	return
}

// VecNormDiff returns the Euclidian norm of the difference:
//  nrm := ||u - v||
func VecNormDiff(u, v []float64) (nrm float64) {
	for i := 0; i < len(v); i++ {
		nrm += (u[i] - v[i]) * (u[i] - v[i])
	}
	nrm = math.Sqrt(nrm)
	return
}

// VecDot returns the dot product between two vectors:
//  s := u dot v
func VecDot(u, v []float64) (res float64) {
	for i := 0; i < len(u); i++ {
		res += u[i] * v[i]
	}
	return
}

// VecCopy copies a vector "b" into vector "a" (scaled):
//  a := α * b  =>  ai := α * bi
func VecCopy(a []float64, α float64, b []float64) {
	for i := 0; i < len(b); i++ {
		a[i] = α * b[i]
	}
}

// VecAdd adds to vector "a", another vector "b" (scaled):
//  a += α * b  =>  ai += α * bi
func VecAdd(a []float64, α float64, b []float64) {
	for i := 0; i < len(b); i++ {
		a[i] += α * b[i]
	}
}

// VecAdd2 adds two vectors (scaled):
//  u := α*a + β*b  =>  ui := α*ai + β*bi
func VecAdd2(u []float64, α float64, a []float64, β float64, b []float64) {
	for i := 0; i < len(b); i++ {
		u[i] = α*a[i] + β*b[i]
	}
}

// VecMin returns the minimum component of a vector
func VecMin(v []float64) (min float64) {
	min = v[0]
	for i := 1; i < len(v); i++ {
		if v[i] < min {
			min = v[i]
		}
	}
	return
}

// VecMax returns the maximum component of a vector
func VecMax(v []float64) (max float64) {
	max = v[0]
	for i := 1; i < len(v); i++ {
		if v[i] > max {
			max = v[i]
		}
	}
	return
}

// VecMinMax returns the min and max components of a vector
func VecMinMax(v []float64) (min, max float64) {
	min, max = v[0], v[0]
	for i := 1; i < len(v); i++ {
		if v[i] < min {
			min = v[i]
		}
		if v[i] > max {
			max = v[i]
		}
	}
	return
}

// VecLargest returns the largest component (abs(u_i)) of a vector, normalised by den
func VecLargest(u []float64, den float64) (largest float64) {
	largest = math.Abs(u[0]) / den
	for i := 1; i < len(u); i++ {
		tmp := math.Abs(u[i]) / den
		if tmp > largest {
			largest = tmp
		}
	}
	return
}

// VecMaxDiff returns the maximum difference between components of two vectors
func VecMaxDiff(a, b []float64) (maxdiff float64) {
	maxdiff = math.Abs(a[0] - b[0])
	for i := 1; i < len(a); i++ {
		diff := math.Abs(a[i] - b[i])
		if diff > maxdiff {
			maxdiff = diff
		}
	}
	return
}

// VecMaxDiffC returns the maximum difference between components of two complex vectors
func VecMaxDiffC(a, b []complex128) (maxdiff float64) {
	maxdiff = math.Abs(real(a[0]) - real(b[0]))
	maxdiffz := math.Abs(imag(a[0]) - imag(b[0]))
	for i := 1; i < len(a); i++ {
		diff := math.Abs(real(a[i]) - real(b[i]))
		diffz := math.Abs(imag(a[i]) - imag(b[i]))
		if diff > maxdiff {
			maxdiff = diff
		}
		if diffz > maxdiffz {
			maxdiffz = diffz
		}
	}
	if maxdiffz > maxdiff {
		maxdiff = maxdiffz
	}
	return
}

// VecScale scales vector "v" using an absolute value (Atol) and a multiplier (Rtol)
//  res[i] := Atol + Rtol * v[i]
func VecScale(res []float64, Atol, Rtol float64, v []float64) {
	for i := 0; i < len(v); i++ {
		res[i] = Atol + Rtol*v[i]
	}
}

// VecScaleAbs scales vector abs(v) using an absolute value (Atol) and a multiplier (Rtol)
//  res[i] := Atol + Rtol * Abs(v[i])
func VecScaleAbs(res []float64, Atol, Rtol float64, v []float64) {
	for i := 0; i < len(v); i++ {
		res[i] = Atol + Rtol*math.Abs(v[i])
	}
}

// VecRms returns the root-mean-square of a vector u:
//  rms := sqrt(mean((u[:])^2))  ==  sqrt(sum_i((ui)^2)/n)
func VecRms(u []float64) (rms float64) {
	for i := 0; i < len(u); i++ {
		rms += u[i] * u[i]
	}
	rms = math.Sqrt(rms / float64(len(u)))
	return
}

// VecRmsErr returns the root-mean-square of a vector u normalised by scal[i] = Atol + Rtol * |vi|
//  rms     := sqrt(sum_i((u[i]/scal[i])^2)/n)
//  scal[i] := Atol + Rtol * |v[i]|
func VecRmsErr(u []float64, Atol, Rtol float64, v []float64) (rms float64) {
	var scal float64
	for i := 0; i < len(v); i++ {
		scal = Atol + Rtol*math.Abs(v[i])
		rms += u[i] * u[i] / (scal * scal)
	}
	rms = math.Sqrt(rms / float64(len(u)))
	return
}

// VecRmsError returns the root-mean-square of a vector u normalised by scal[i] = Atol + Rtol * |vi|
//  rms     := sqrt(sum_i((|u[i]-w[i]|/scal[i])^2)/n)
//  scal[i] := Atol + Rtol * |v[i]|
func VecRmsError(u, w []float64, Atol, Rtol float64, v []float64) (rms float64) {
	var scal, e float64
	for i := 0; i < len(v); i++ {
		scal = Atol + Rtol*math.Abs(v[i])
		e = math.Abs(u[i] - w[i])
		rms += e * e / (scal * scal)
	}
	rms = math.Sqrt(rms / float64(len(u)))
	return
}
