// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/chk"

// LagrangeCardinal computes the [i] Lagrange cardinal p-degree polynomial @ x
//
//              p
//     p       ━━━━      x  - X[j]
//    L  (x) = ┃  ┃  ——————————————
//     i       ┃  ┃    X[i] - X[j]
//            j = 0
//            j ≠ i
//
//   Input:
//      p -- degree
//      i -- index of X[i] point
//      x -- where to evaluate the L^p_i polynomial
//      X -- the p+1 points
//   Output:
//      lpi -- L^p_i(x)
//
func LagrangeCardinal(p, i int, x float64, X []float64) (lpi float64) {
	if len(X) != p+1 {
		chk.Panic("len(X) must be equal to p+1")
	}
	lpi = 1
	for j := 0; j < p+1; j++ {
		if i != j {
			lpi *= (x - X[j]) / (X[i] - X[j])
		}
	}
	return
}
