// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "github.com/cpmech/gosl/fun"

// QuadGaussL10 approximates the integral of the function f(x) between a and b, by ten-point
// Gauss-Legendre integration. The function is evaluated exactly ten times at interior points
// in the range of integration. See page 180 of [1].
//   Reference:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func QuadGaussL10(a, b float64, f fun.Ss) (res float64, err error) {

	// constants
	x := []float64{0.1488743389816312, 0.4333953941292472, 0.6794095682990244, 0.8650633666889845, 0.9739065285171717}
	w := []float64{0.2955242247147529, 0.2692667193099963, 0.2190863625159821, 0.1494513491505806, 0.0666713443086881}

	// auxiliary variables
	xm := 0.5 * (b + a)
	xr := 0.5 * (b - a)
	s := 0.0 // will be twice the average value of the function, since the ten weights (five numbers above each used twice) sum to 2.

	// execute sum
	var dx, fp, fm float64
	for j := 0; j < 5; j++ {
		dx = xr * x[j]
		fp, err = f(xm + dx)
		if err != nil {
			return
		}
		fm, err = f(xm - dx)
		if err != nil {
			return
		}
		s += w[j] * (fp + fm)
	}
	res = s * xr // scale the answer to the range of integration.
	return
}
