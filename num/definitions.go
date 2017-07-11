// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package num implements fundamental numerical methods such as numerical derivative and quadrature,
// root finding solvers (Brent's and Newton's methods), among others.
package num

import "math"

// constants
var (
	MACHEPS = math.Nextafter(1, 2) - 1.0 // smallest number satisfying 1 + EPS > 1
)
