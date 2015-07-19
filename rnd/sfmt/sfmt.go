// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sfmt

/*
#cgo CFLAGS: -O3 -finline-functions -fomit-frame-pointer -DNDEBUG -fno-strict-aliasing --param max-inline-insns-single=1800 -std=c99 -msse2 -DHAVE_SSE2 -DSFMT_MEXP=19937
// #cgo CFLAGS: -Wall // deactivated since cgo fails with Wunused-variable for SfmtPrintIdString
#include "connectsfmt.h"
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
	C.SfmtInit(C.long(seed))
}

// Rand generates pseudo random integer between low and high
//  Input:
//   low  -- lower limit
//   high -- upper limit
//  Output:
//   random integer
func Rand(low, high int) int {
	return int(C.SfmtRand(C.long(low), C.long(high)))
}

// Shuffle shuffles slice of integers
func Shuffle(values []int) {
	C.SfmtShuffle((*C.long)(unsafe.Pointer(&values[0])), C.long(len(values)))
}

// PrintIdString prints SFMT id string
func PrintIdString() {
	if io.Verbose {
		C.SfmtPrintIdString()
	}
}
