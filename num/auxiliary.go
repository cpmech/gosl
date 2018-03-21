// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "math"

// sgn returns a value with the same magnitude as a and the same sign as b
//
//  returns:  |a| * sign(b)
//
func sgn(a, b float64) float64 {
	if b < 0 {
		return -math.Abs(a) // return - |a|
	}
	return math.Abs(a) // return + |a|
}

func swap(a, b *float64) {
	*a, *b = *b, *a
}

func shft2(a, b *float64, c float64) {
	*a = *b
	*b = c
}

func shft3(a, b, c *float64, d float64) {
	*a = *b
	*b = *c
	*c = d
}

func mov3(a, b, c *float64, d, e, f float64) {
	*a = d
	*b = e
	*c = f
}
