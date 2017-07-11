// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math/rand"
	"testing"
)

var (
	benchNsamples int
	benchInts     []int
	benchResult   int
)

func init() {
	rand.Seed(4321)
	Init(4321)
	benchNsamples = 1000
	benchInts = make([]int, benchNsamples)
	for i := 0; i < benchNsamples; i++ {
		benchInts[i] = i
	}
}

func Benchmark_go_int(b *testing.B) {
	var res int
	lo, hi := 0, 50
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchNsamples; j++ {
			res = Int(lo, hi)
		}
	}
	benchResult = res
}
