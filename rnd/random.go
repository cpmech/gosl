// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math/rand"
	"time"

	"github.com/cpmech/gosl/rnd/dsfmt"
	"github.com/cpmech/gosl/rnd/sfmt"
)

// Init initialises random numbers generators
//  Input:
//   seed -- seed value; use seed <= 0 to use current time
func Init(seed int) {
	if seed <= 0 {
		seed = int(time.Now().Unix())
	}
	rand.Seed(int64(seed))
	sfmt.Init(seed)
	dsfmt.Init(seed)
}

// Int generates pseudo random integer between low and high.
//  Input:
//   low  -- lower limit
//   high -- upper limit
//  Output:
//   random integer
func Int(low, high int) int {
	return rand.Int()%(high-low+1) + low
}

// Ints generates pseudo random integers between low and high.
//  Input:
//   values -- slice to be filled with len(values) numbers
//   low    -- lower limit
//   high   -- upper limit
func Ints(values []int, low, high int) {
	if len(values) < 1 {
		return
	}
	for i := 0; i < len(values); i++ {
		values[i] = Int(low, high)
	}
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
//   values -- slice to be filled with len(values) numbers
//   low    -- lower limit
//   high   -- upper limit
func MTints(values []int, low, high int) {
	if len(values) < 1 {
		return
	}
	for i := 0; i < len(values); i++ {
		values[i] = MTint(low, high)
	}
}

// MTfloat64 generates pseudo random real numbers between low and high; i.e. in [low, right)
// using the Mersenne Twister method
//  Input:
//   low  -- lower limit (closed)
//   high -- upper limit (open)
//  Output:
//   random float64
func MTfloat64(low, high float64) float64 {
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

// DblShuffle shuffles a slice of float point numbers
func DblShuffle(v []float64) {
	//TODO
	//I :=utl.IntRange(len(v))
	//smft.Shuffle(I)
	//for j, i := range I {
	//v[j]
	//}
}
