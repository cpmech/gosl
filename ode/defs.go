// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "github.com/cpmech/gosl/la"

// callbacks
type Cb_fcn func(f []float64, h, x float64, y []float64, args ...interface{}) error      // function
type Cb_jac func(dfdy *la.Triplet, h, x float64, y []float64, args ...interface{}) error // Jacobian (must have at least all diagonal elements set)
type Cb_out func(first bool, h, x float64, y []float64, args ...interface{}) error       // output

// callbacks
type Cb_ycorr func(y []float64, x float64, args ...interface{}) // y(x) correct

// step function
type stpfcn func(o *Solver, y []float64, x float64, args ...interface{}) (rerr float64, err error)

// accept update function
type acptfcn func(o *Solver, y []float64)
