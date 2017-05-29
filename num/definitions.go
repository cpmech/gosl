// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package num implements fundamental numerical methods such as numerical derivative and quadrature,
// root finding solvers (Brent's and Newton's methods), among others.
package num

import "math"

// constants
const (
	MACHEPS     = 1e-16 // smallest number satisfying 1.0 + EPS > 1.0
	CTE1        = 1e-5  // a minimum value to be used in Jacobian
	DBL_EPSILON = 1.0e-15
	DBL_MIN     = math.SmallestNonzeroFloat64
)

// flags
const (
	NaN = iota
	Inf
	Equal
	NotEqual
)
