// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fun (functions) implements special functions such as elliptical, orthogonal polynomials,
// Bessel, discrete Fourier transform, polynomial interpolators, and more.
package fun

import (
	"math"

	"github.com/cpmech/gosl/la"
)

// π = 3.141592653589...
const π = math.Pi

// Ss defines a scalar function f(s) of a scalar argument s. Also returns error
//   Scalar scalar
type Ss func(s float64) (float64, error)

// Sv defines a scalar functioin f(v) of a vector argument v. Also returns error
type Sv func(v la.Vector) (float64, error)

// Vs defines a vector function f(s) of a scalar argument s. Also returns error
//   Vector scalar
type Vs func(f la.Vector, s float64) error

// Vv defines a vector function f(v) of a vector argument v. Also returns error
//   Vector vector
type Vv func(f, v la.Vector) error

// Mv defines a matrix function f(v) of a vector argument v. Also returns error
//   Matrix vector
type Mv func(f *la.Matrix, v la.Vector) error

// Mm defines a matrix function f(m) of a matrix argument m. Also returns error
//   Matrix matrix
type Mm func(f, m *la.Matrix) error

// Tv defines a triplet (matrix) function f(v) of a vector argument v. Also returns error
//   Triplet vector
type Tv func(f *la.Triplet, v la.Vector) error

// Tt defines a triplet (matrix) function f(t) of a triplet (matrix) argument t. Also returns error
//   Triplet triplet
type Tt func(f, t *la.Triplet) error
