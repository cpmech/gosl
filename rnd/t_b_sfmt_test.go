// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package rnd

import (
	"math/rand"
	"testing"
)

var (
	benchMTnsamples int
	benchMTints     []int
	benchMTresult   int
)

func init() {
	rand.Seed(4321)
	MTinit(4321)
	benchMTnsamples = 1000
	benchMTints = make([]int, benchMTnsamples)
	for i := 0; i < benchMTnsamples; i++ {
		benchMTints[i] = i
	}
}

func Benchmark_mt_int(b *testing.B) {
	var res int
	lo, hi := 0, 50
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchMTnsamples; j++ {
			res = MTint(lo, hi)
		}
	}
	benchMTresult = res
}
