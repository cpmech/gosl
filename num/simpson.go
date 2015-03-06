// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "github.com/cpmech/gosl/chk"

// scalar function
type Fun0 func(float64) float64

// Simpson integrates function f from x=a to x=b using n subintervals (n must be even)
func Simpson(f Fun0, a, b float64, n int) (float64, error) {
	if n < 2 || n%2 > 0 {
		return 0, chk.Err("number of subintervas should be even (n=%d)", n)
	}

	h := (b - a) / float64(n)
	sum := f(a) + f(b)

	x := a
	for i := 1; i < n; i++ {
		x += h
		if i%2 == 1 { // i is odd
			sum += 4 * f(x)
		} else { // i is even
			sum += 2 * f(x)
		}
	}

	sum = sum * h / 3

	return sum, nil
}
