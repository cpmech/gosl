// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

// HaltonPoints generates randomly spaced points
//   x -- [dim][n] points
func HaltonPoints(dim, n int) (x [][]float64) {
	x = utl.Alloc(dim, n)
	for j := 0; j < dim; j++ {
		for i := 0; i < n; i++ {
			x[j][i] = halton(i, j)
		}
	}
	return
}

// halton implements (Halton and Weller 1964) method
//  i -- population (point) index
//  j -- parameter (dimension) index
func halton(i, j int) (sum float64) {
	if j > 999 {
		chk.Panic("HaltonPoints can only handle maximum dimension = 1000")
	}
	p1 := PRIMES1000[j]
	p2 := p1
	first := true
	for first || i > 0 {
		w := i % p1
		sum = sum + float64(w)/float64(p2)
		i = int(float64(i) / float64(p1))
		p2 = p2 * p1
		first = false
	}
	return
}
