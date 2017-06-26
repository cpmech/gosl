// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"math/cmplx"
)

// Vector defines the vector type for real numbers simply as a slice of float64
type Vector []float64

// NewVector returns a new vector with size m
func NewVector(m int) Vector {
	return make([]float64, m)
}

// NewVectorMapped returns a new vector after applying a function over all of its components
//  new: vi = f(i)
func NewVectorMapped(m int, f func(i int) float64) (o Vector) {
	o = make([]float64, m)
	for i := 0; i < len(o); i++ {
		o[i] = f(i)
	}
	return
}

// Fill fills this vector with a single number val
//  vi = val
func (o Vector) Fill(val float64) {
	for i := 0; i < len(o); i++ {
		o[i] = val
	}
}

// ApplyFunc runs a function over all components of a vector
//  vi = f(i,vi)
func (o Vector) ApplyFunc(f func(i int, x float64) float64) {
	for i := 0; i < len(o); i++ {
		o[i] = f(i, o[i])
	}
}

// GetCopy returns a copy of this vector
//  b := a
func (o Vector) GetCopy() (clone Vector) {
	clone = make([]float64, len(o))
	copy(clone, o)
	return
}

// CopyInto copies the scaled components of this vector into another one (result)
//  result := α * this   ⇒   result[i] := α * this[i]
func (o Vector) CopyInto(result Vector, α float64) {
	for i := 0; i < len(o); i++ {
		result[i] = α * o[i]
	}
}

// Accum sum/accumulates all components in a vector
//  sum := Σ_i v[i]
func (o Vector) Accum() (sum float64) {
	for i := 0; i < len(o); i++ {
		sum += o[i]
	}
	return
}

// Norm returns the Euclidian norm of a vector:
//  nrm := ‖v‖
func (o Vector) Norm() (nrm float64) {
	for i := 0; i < len(o); i++ {
		nrm += o[i] * o[i]
	}
	return math.Sqrt(nrm)
}

// NormDiff returns the Euclidian norm of the difference:
//  nrm := ||u - v||
func (o Vector) NormDiff(v Vector) (nrm float64) {
	for i := 0; i < len(v); i++ {
		nrm += (o[i] - v[i]) * (o[i] - v[i])
	}
	nrm = math.Sqrt(nrm)
	return
}

// VecMin returns the minimum component of a vector
func (o Vector) Min() (min float64) {
	min = o[0]
	for i := 1; i < len(o); i++ {
		if o[i] < min {
			min = o[i]
		}
	}
	return
}

// Max returns the maximum component of a vector
func (o Vector) Max() (max float64) {
	max = o[0]
	for i := 1; i < len(o); i++ {
		if o[i] > max {
			max = o[i]
		}
	}
	return
}

// MinMax returns the min and max components of a vector
func (o Vector) MinMax() (min, max float64) {
	min, max = o[0], o[0]
	for i := 1; i < len(o); i++ {
		if o[i] < min {
			min = o[i]
		}
		if o[i] > max {
			max = o[i]
		}
	}
	return
}

// Largest returns the largest component |u[i]| of this vector, normalised by den
//   largest := |u[i]| / den
func (o Vector) Largest(den float64) (largest float64) {
	largest = math.Abs(o[0])
	for i := 1; i < len(o); i++ {
		tmp := math.Abs(o[i])
		if tmp > largest {
			largest = tmp
		}
	}
	return largest / den
}

// Dot returns the dot product between two vectors:
//  s := u dot v
func (o Vector) Dot(v Vector) (res float64) {
	for i := 0; i < len(o); i++ {
		res += o[i] * v[i]
	}
	return
}

// Add adds the scaled components of this vector to another one (result)
//  result += α * this   ⇒   result[i] += α * this[i]
func (o Vector) Add(result Vector, α float64) {
	for i := 0; i < len(o); i++ {
		result[i] += α * o[i]
	}
}

// Add2 adds the scaled components of this vector and another vector to a third vector (result)
//  result := α*this + β*another   ⇒   result[i] := α*this[i] + β*another[i]
func (o Vector) Add2(result Vector, α, β float64, another Vector) {
	for i := 0; i < len(o); i++ {
		result[i] = α*o[i] + β*another[i]
	}
}

// MaxDiff returns the maximum difference between the components of this and another vector
func (o Vector) MaxDiff(another Vector) (maxdiff float64) {
	maxdiff = math.Abs(o[0] - another[0])
	for i := 1; i < len(o); i++ {
		diff := math.Abs(o[i] - another[i])
		if diff > maxdiff {
			maxdiff = diff
		}
	}
	return
}

// Scale scales vector using a shift value (a) and a multiplier (m)
//  this[i] = a + m * this[i]
func (o Vector) Scale(a, m float64) {
	for i := 0; i < len(o); i++ {
		o[i] = a + m*o[i]
	}
}

// GetScaled returns a scaled vector using a shift value (a) and a multiplier (m)
//  result[i] := a + m * this[i]
func (o Vector) GetScaled(a, m float64) (result Vector) {
	result = make([]float64, len(o))
	for i := 0; i < len(o); i++ {
		result[i] = a + m*o[i]
	}
	return
}

// ScaleAbs scales vector using a shift value (a) and a multiplier (m) applied to the absolute value
// of another vector components.
//  this[i] = a + m * |this[i]|
func (o Vector) ScaleAbs(a, m float64, another Vector) {
	for i := 0; i < len(o); i++ {
		o[i] = a + m*math.Abs(another[i])
	}
}

// Rms returns the root-mean-square of this vector
//                ________________________
//               /     ————            2
//              /  1   \    /         \
//   rms =  \  /  ———  /    | this[i] |
//           \/    N   ———— \         /
//
func (o Vector) Rms() (rms float64) {
	for i := 0; i < len(o); i++ {
		rms += o[i] * o[i]
	}
	rms = math.Sqrt(rms / float64(len(o)))
	return
}

// RmsScaled returns the scaled root-mean-square of this vector with components
// normalised by a scaling factor
//                ___________________________
//               /     ————               2
//              /  1   \    /   this[i]  \
//   rms =  \  /  ———  /    | —————————— |
//           \/    N   ———— \  scale[i]  /
//
//   scale[i] = a + m*|s[i]|
//
func (o Vector) RmsScaled(a, m float64, s Vector) (rms float64) {
	var scale float64
	for i := 0; i < len(o); i++ {
		scale = a + m*math.Abs(s[i])
		rms += o[i] * o[i] / (scale * scale)
	}
	return math.Sqrt(rms / float64(len(o)))
}

// RmsError returns the scaled root-mean-square of the difference between this vector and another
// with components normalised by a scaling factor
//                __________________________
//               /     ————              2
//              /  1   \    /  error[i]  \
//   rms =  \  /  ———  /    | —————————— |
//           \/    N   ———— \  scale[i]  /
//
//   error[i] = |this[i] - another[i]|
//
//   scale[i] = a + m*|s[i]|
//
func (o Vector) RmsError(a, m float64, s, another Vector) (rms float64) {
	var scale, err float64
	for i := 0; i < len(o); i++ {
		scale = a + m*math.Abs(s[i])
		err = math.Abs(o[i] - another[i])
		rms += err * err / (scale * scale)
	}
	return math.Sqrt(rms / float64(len(o)))
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// VectorC defines the vector type for complex numbers simply as a slice of complex128
type VectorC []complex128

// NewVectorC returns a new vector with size m
func NewVectorC(m int) VectorC {
	return make([]complex128, m)
}

// NewVectorMappedC returns a new vector after applying a function over all of its components
//  new: vi = f(i)
func NewVectorMappedC(m int, f func(i int) complex128) (o VectorC) {
	o = make([]complex128, m)
	for i := 0; i < len(o); i++ {
		o[i] = f(i)
	}
	return
}

// Fill fills a vector with a single number s:
//  v := s*ones(len(v))  =>  vi = s
func (o VectorC) Fill(s complex128) {
	for i := 0; i < len(o); i++ {
		o[i] = s
	}
}

// ApplyFunc runs a function over all components of a vector
//  vi = f(i,vi)
func (o VectorC) ApplyFunc(f func(i int, x complex128) complex128) {
	for i := 0; i < len(o); i++ {
		o[i] = f(i, o[i])
	}
}

// GetCopy returns a copy of this vector
//  b := a
func (o VectorC) GetCopy() (clone VectorC) {
	clone = make([]complex128, len(o))
	copy(clone, o)
	return
}

// CopyInto copies the scaled components of this vector into another one (result)
//  result := α * this   ⇒   result[i] := α * this[i]
func (o VectorC) CopyInto(result VectorC, α complex128) {
	for i := 0; i < len(o); i++ {
		result[i] = α * o[i]
	}
}

// Norm returns the Euclidian norm of a vector:
//  nrm := ‖v‖
func (o VectorC) Norm() (nrm complex128) {
	for i := 0; i < len(o); i++ {
		nrm += o[i] * o[i]
	}
	return cmplx.Sqrt(nrm)
}

// MaxDiff returns the maximum difference between the components of two vectors
func (o VectorC) MaxDiff(b VectorC) float64 {
	maxdiffR := math.Abs(real(o[0]) - real(b[0]))
	maxdiffC := math.Abs(imag(o[0]) - imag(b[0]))
	for i := 1; i < len(o); i++ {
		diffR := math.Abs(real(o[i]) - real(b[i]))
		diffC := math.Abs(imag(o[i]) - imag(b[i]))
		if diffR > maxdiffR {
			maxdiffR = diffR
		}
		if diffC > maxdiffC {
			maxdiffC = diffC
		}
	}
	if maxdiffR > maxdiffC {
		return maxdiffR
	}
	return maxdiffC
}

// GetScaled returns a scaled vector using a shift value (a) and a multiplier (m)
//  result[i] := a + m * this[i]
func (o VectorC) GetScaled(a, m complex128) (result VectorC) {
	result = make([]complex128, len(o))
	for i := 0; i < len(o); i++ {
		result[i] = a + m*o[i]
	}
	return
}
