// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/la"

// Ss defines a scalar function f(s) of a scalar argument s. Also returns error
//   Scalar scalar
type Ss func(s float64) (float64, error)

// Sv defines a scalar functioin f(v) of a vector argument v. Also returns error
type Sv func(v []float64) (float64, error)

// Vs defines a vector function f(s) of a scalar argument s. Also returns error
//   Vector scalar
type Vs func(f []float64, s float64) error

// Vv defines a vector function f(v) of a vector argument v. Also returns error
//   Vector vector
type Vv func(f, v []float64) error

// Mv defines a matrix function f(v) of a vector argument v. Also returns error
//   Matrix vector
type Mv func(f [][]float64, v []float64) error

// Mm defines a matrix function f(m) of a matrix argument m. Also returns error
//   Matrix matrix
type Mm func(f, m [][]float64) error

// Tv defines a triplet (matrix) function f(v) of a vector argument v. Also returns error
//   Triplet vector
type Tv func(f *la.Triplet, v []float64) error

// Tt defines a triplet (matrix) function f(t) of a triplet (matrix) argument t. Also returns error
//   Triplet triplet
type Tt func(f, t *la.Triplet) error
