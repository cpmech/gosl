// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package sfmt

import (
	"math/rand"
	"testing"
)

var benchNsamples int
var benchResult int

func init() {
	rand.Seed(4321)
	Init(4321)
	benchNsamples = 1000
}

func Benchmark_gornd_int(b *testing.B) {
	var res int
	lo, hi := 0, 50
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchNsamples; j++ {
			res = rand.Int()%(hi-lo+1) + lo
		}
	}
	benchResult = res
}

func Benchmark_sfmt_int(b *testing.B) {
	var res int
	lo, hi := 0, 50
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchNsamples; j++ {
			res = Rand(lo, hi)
		}
	}
	benchResult = res
}
