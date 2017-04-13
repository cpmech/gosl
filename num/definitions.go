// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "github.com/cpmech/gosl/la"

// callbacks (for single equations)
type Cb_yx func(x float64) float64
type Cb_yxe func(x float64) (float64, error)
type Cb_fx func(x float64, args ...interface{}) float64

// callbacks (for systems)
type Cb_f func(fx, x []float64) error
type Cb_J func(dfdx *la.Triplet, x []float64) error  // sparse version
type Cb_Jd func(dfdx [][]float64, x []float64) error // dense version
type Cb_out func(x []float64) error                  // for output

// constants
const (
	EPS  = 1e-16 // smallest number satisfying 1.0 + EPS > 1.0
	CTE1 = 1e-5  // a minimum value to be used in Jacobian
)
