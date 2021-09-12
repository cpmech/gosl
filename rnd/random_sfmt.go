// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package rnd

import (
	"time"

	"github.com/cpmech/gosl/rnd/dsfmt"
	"github.com/cpmech/gosl/rnd/sfmt"
)

// MTinit initializes random numbers generators (Mersenne Twister code)
//  Input:
//   seed -- seed value; use seed <= 0 to use current time
func MTinit(seed int) {
	if seed <= 0 {
		seed = int(time.Now().Unix())
	}
	sfmt.Init(seed)
	dsfmt.Init(seed)
}

// MTint generates pseudo random integer between low and high using the Mersenne Twister method.
//  Input:
//   low  -- lower limit
//   high -- upper limit
//  Output:
//   random integer
func MTint(low, high int) int {
	return sfmt.Rand(low, high)
}

// MTints generates pseudo random integers between low and high using the Mersenne Twister method.
//  Input:
//   low    -- lower limit
//   high   -- upper limit
//  Output:
//   values -- slice to be filled with len(values) numbers
func MTints(values []int, low, high int) {
	if len(values) < 1 {
		return
	}
	for i := 0; i < len(values); i++ {
		values[i] = MTint(low, high)
	}
}

// MTfloat64 generates pseudo random real numbers between low and high; i.e. in [low, right)
// using the Mersenne Twister method.
//  Input:
//   low  -- lower limit (closed)
//   high -- upper limit (open)
//  Output:
//   random float64
func MTfloat64(low, high float64) float64 {
	return dsfmt.Rand(low, high)
}

// MTfloat64s generates pseudo random real numbers between low and high; i.e. in [low, right)
// using the Mersenne Twister method.
//  Input:
//   low  -- lower limit (closed)
//   high -- upper limit (open)
//  Output:
//   values -- slice to be filled with len(values) numbers
func MTfloat64s(values []float64, low, high float64) {
	for i := 0; i < len(values); i++ {
		values[i] = dsfmt.Rand(low, high)
	}
}

// MTintShuffle shuffles a slice of integers using Mersenne Twister algorithm.
func MTintShuffle(v []int) {
	sfmt.Shuffle(v)
}
