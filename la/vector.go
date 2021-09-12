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

// NewVectorSlice returns a new vector from given Slice
// NOTE: This is equivalent to cast a slice to Vector as in:
//    v := la.Vector([]float64{1,2,3})
func NewVectorSlice(v []float64) Vector {
	return Vector(v)
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

// Apply sets this vector with the scaled components of another vector
//  this := α * another   ⇒   this[i] := α * another[i]
//  NOTE: "another" may be "this"
func (o Vector) Apply(α float64, another Vector) {
	for i := 0; i < len(o); i++ {
		o[i] = α * another[i]
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

// GetUnit returns the unit vector parallel to this vector
//  b := a / norm(a)
func (o Vector) GetUnit() (unit Vector) {
	unit = make([]float64, len(o))
	s := o.Norm()
	if s > 0 {
		unit.Apply(1.0/s, o)
	}
	return
}

// Accum sum/accumulates all components in a vector
//  sum := Σ_i v[i]
func (o Vector) Accum() (sum float64) {
	for i := 0; i < len(o); i++ {
		sum += o[i]
	}
	return
}

// Norm returns the Euclidean norm of a vector:
//  nrm := ‖v‖
func (o Vector) Norm() (nrm float64) {
	return math.Sqrt(VecDot(o, o))
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

// NormDiff returns the Euclidean norm of the difference:
//  nrm := ||u - v||
func (o Vector) NormDiff(v Vector) (nrm float64) {
	for i := 0; i < len(v); i++ {
		nrm += (o[i] - v[i]) * (o[i] - v[i])
	}
	nrm = math.Sqrt(nrm)
	return
}

// Min returns the minimum component of a vector
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

// Largest returns the largest component |u[i]| of this vector, normalized by den
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

// Apply sets this vector with the scaled components of another vector
//  this := α * another   ⇒   this[i] := α * another[i]
//  NOTE: "another" may be "this"
func (o VectorC) Apply(α complex128, another VectorC) {
	for i := 0; i < len(o); i++ {
		o[i] = α * another[i]
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

// Norm returns the Euclidean norm of a vector:
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

// SplitRealImag splits this vector into two vectors with the real and imaginary parts
//  xR := real(this)
//  xI := imag(this)
//  NOTE: xR and xI must be pre-allocated with length = len(this)
func (o VectorC) SplitRealImag(xR, xI Vector) {
	for i := 0; i < len(o); i++ {
		xR[i] = real(o[i])
		xI[i] = imag(o[i])
	}
}

// JoinRealImag sets this vector with two vectors having the real and imaginary parts
//  this := complex(xR, xI)
//  NOTE: len(xR) == len(xI) == len(this)
func (o VectorC) JoinRealImag(xR, xI Vector) {
	for i := 0; i < len(o); i++ {
		o[i] = complex(xR[i], xI[i])
	}
}
