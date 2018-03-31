// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"testing"

	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

const benchmarkArrayIntSize = 10000

var (
	benchmarkArrayFwd []int
	benchmarkArrayBwd []int
	benchmarkArrayRnd []int
	benchmarkArrayCpy []int
)

func init() {
	benchmarkArrayFwd = utl.IntRange(benchmarkArrayIntSize)
	benchmarkArrayBwd = utl.IntRange3(benchmarkArrayIntSize-1, -1, -1)
	benchmarkArrayRnd = make([]int, benchmarkArrayIntSize)
	benchmarkArrayCpy = make([]int, benchmarkArrayIntSize)
	copy(benchmarkArrayRnd, benchmarkArrayFwd)
	rnd.Init(0)
	rnd.IntShuffle(benchmarkArrayRnd)
	return
}

func BenchmarkRecQuickSortFwd(b *testing.B) {
	copy(benchmarkArrayCpy, benchmarkArrayFwd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntRecQuickSort(benchmarkArrayCpy, IntComparator)
	}
}

func BenchmarkRecQuickSortBwd(b *testing.B) {
	copy(benchmarkArrayCpy, benchmarkArrayBwd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntRecQuickSort(benchmarkArrayCpy, IntComparator)
	}
}

func BenchmarkRecQuickSortRnd(b *testing.B) {
	copy(benchmarkArrayCpy, benchmarkArrayRnd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntRecQuickSort(benchmarkArrayCpy, IntComparator)
	}
}

func BenchmarkRecQuickSortFwdNonOpt(b *testing.B) {
	copy(benchmarkArrayCpy, benchmarkArrayFwd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntRecQuickSortNonOpt(benchmarkArrayCpy, IntComparator)
	}
}

func BenchmarkRecQuickSortBwdNonOpt(b *testing.B) {
	copy(benchmarkArrayCpy, benchmarkArrayBwd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntRecQuickSortNonOpt(benchmarkArrayCpy, IntComparator)
	}
}

func BenchmarkRecQuickSortRndNonOpt(b *testing.B) {
	copy(benchmarkArrayCpy, benchmarkArrayRnd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntRecQuickSortNonOpt(benchmarkArrayCpy, IntComparator)
	}
}
