// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"math/rand"

	"gosl/utl"
)

// UnitVectors generates random unit vectors in 3D
func UnitVectors(n int) (U [][]float64) {
	U = utl.Alloc(n, 3)
	for i := 0; i < n; i++ {
		φ := 2.0 * math.Pi * rand.Float64()
		θ := math.Acos(1.0 - 2.0*rand.Float64())
		U[i][0] = math.Sin(θ) * math.Cos(φ)
		U[i][1] = math.Sin(θ) * math.Sin(φ)
		U[i][2] = math.Cos(θ)
	}
	return
}
