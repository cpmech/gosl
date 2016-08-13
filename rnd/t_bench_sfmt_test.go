// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package rnd

import (
	"math/rand"
	"testing"
)

var __bench_mt_nsamples int
var __bench_mt_ints []int
var __bench_mt_result int

func init() {
	rand.Seed(4321)
	MTinit(4321)
	__bench_mt_nsamples = 1000
	__bench_mt_ints = make([]int, __bench_mt_nsamples)
	for i := 0; i < __bench_mt_nsamples; i++ {
		__bench_mt_ints[i] = i
	}
}

func Benchmark_mt_int(b *testing.B) {
	var res int
	lo, hi := 0, 50
	for i := 0; i < b.N; i++ {
		for j := 0; j < __bench_mt_nsamples; j++ {
			res = MTint(lo, hi)
		}
	}
	__bench_mt_result = res
}
