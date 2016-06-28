// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sfmt

import (
	"math/rand"
	"testing"
)

var __bench_nsamples int
var __bench_result int

func init() {
	rand.Seed(4321)
	Init(4321)
	__bench_nsamples = 1000
}

func Benchmark_gornd_int(b *testing.B) {
	var res int
	lo, hi := 0, 50
	for i := 0; i < b.N; i++ {
		for j := 0; j < __bench_nsamples; j++ {
			res = rand.Int()%(hi-lo+1) + lo
		}
	}
	__bench_result = res
}

func Benchmark_sfmt_int(b *testing.B) {
	var res int
	lo, hi := 0, 50
	for i := 0; i < b.N; i++ {
		for j := 0; j < __bench_nsamples; j++ {
			res = Rand(lo, hi)
		}
	}
	__bench_result = res
}
