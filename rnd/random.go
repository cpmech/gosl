// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"time"

	"github.com/cpmech/gosl/rnd/dsfmt"
	"github.com/cpmech/gosl/rnd/sfmt"
)

// Init initialises random numbers generator
//  Input:
//   seed -- seed value; use seed <= 0 to use current time
func Init(seed int) {
	if seed <= 0 {
		seed = int(time.Now().Unix())
	}
	sfmt.Init(seed)
	dsfmt.Init(seed)
}

// IntRand generates pseudo random integer between low and high
//  Input:
//   low  -- lower limit
//   high -- upper limit
//  Output:
//   random integer
func IntRand(low, high int) int {
	return sfmt.Rand(low, high)
}

// IntRands generates pseudo random integers between low and high
//  Input:
//   values -- slice to be filled with len(values) numbers
//   low    -- lower limit
//   high   -- upper limit
func IntRands(values []int, low, high int) {
	if len(values) < 1 {
		return
	}
	for i := 0; i < len(values); i++ {
		values[i] = IntRand(low, high)
	}
}

// DblRand generates pseudo random real numbers between low and high; i.e. in [low, right)
//  Input:
//   low  -- lower limit (closed)
//   high -- upper limit (open)
//  Output:
//   random float64
func DblRand(low, high float64) float64 {
	return dsfmt.Rand(low, high)
}

// FlipCoin generates a Bernoulli variable; throw a coin with probability p
func FlipCoin(p float64) bool {
	if p == 1.0 {
		return true
	}
	if p == 0.0 {
		return false
	}
	if dsfmt.Rand01() <= p {
		return true
	}
	return false
}

// IntShuffle shuffles a slice of integers
func IntShuffle(v []int) {
	sfmt.Shuffle(v)
}
