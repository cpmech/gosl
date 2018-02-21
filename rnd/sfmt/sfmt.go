// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

// Package sfmt wraps the SFMT SIMD-oriented Fast Mersenne Twister
package sfmt

/*
#cgo CFLAGS: -O3 -fomit-frame-pointer -DNDEBUG -fno-strict-aliasing -std=c99 -msse2 -DHAVE_SSE2 -DSFMT_MEXP=19937
// #cgo CFLAGS: -Wall // deactivated since cgo fails with Wunused-variable for SfmtPrintIdString
#include "connectsfmt.h"
#ifdef WIN32
#define LONG long long
#else
#define LONG long
#endif
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/cpmech/gosl/io"
)

// Init initialises random numbers generator
//  Input:
//   seed -- seed value; use seed <= 0 to use current time
func Init(seed int) {
	if seed <= 0 {
		seed = int(time.Now().Unix())
	}
	C.SfmtInit(C.LONG(seed))
}

// Rand generates pseudo random integer between low and high
//  Input:
//   low  -- lower limit
//   high -- upper limit
//  Output:
//   random integer
func Rand(low, high int) int {
	return int(C.SfmtRand(C.LONG(low), C.LONG(high)))
}

// Shuffle shuffles slice of integers
func Shuffle(values []int) {
	C.SfmtShuffle((*C.LONG)(unsafe.Pointer(&values[0])), C.LONG(len(values)))
}

// PrintIDString prints SFMT id string
func PrintIDString() {
	if io.Verbose {
		C.SfmtPrintIdString()
	}
}
