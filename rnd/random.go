// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

/*
#cgo CFLAGS: -O3 -finline-functions -fomit-frame-pointer -DNDEBUG -fno-strict-aliasing --param max-inline-insns-single=1800 -Wall -std=c99 -msse2 -DHAVE_SSE2 -DSFMT_MEXP=19937
#include "random.h"
*/
import "C"

import "time"

// Init initialises random numbers generator
// Input:
//  seed -- seed value; use seed <= 0 to use current time
func Init(seed int) {
	if seed <= 0 {
		seed = int(time.Now().Unix())
	}
	C.Init(C.long(seed))
}

// IntRand generates pseudo random integer between low and high
func IntRand(low, high int) int {
	return int(C.IntRand(C.long(low), C.long(high)))
}
