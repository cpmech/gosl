// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"
)

var (
	benchRes float64
)

func BenchmarkPowP10(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 10; n++ {
			res = PowP(x, n)
		}
	}
	benchRes = res
}

func BenchmarkPowP10std(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 10; n++ {
			res = math.Pow(x, float64(n))
		}
	}
	benchRes = res
}

func BenchmarkPowP20(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 20; n++ {
			res = PowP(x, n)
		}
	}
	benchRes = res
}

func BenchmarkPowP20std(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 20; n++ {
			res = math.Pow(x, float64(n))
		}
	}
	benchRes = res
}

func BenchmarkPowP50(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 50; n++ {
			res = PowP(x, n)
		}
	}
	benchRes = res
}

func BenchmarkPowP50std(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 50; n++ {
			res = math.Pow(x, float64(n))
		}
	}
	benchRes = res
}

func BenchmarkPowP100(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 100; n++ {
			res = PowP(x, n)
		}
	}
	benchRes = res
}

func BenchmarkPowP100std(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 100; n++ {
			res = math.Pow(x, float64(n))
		}
	}
	benchRes = res
}

func BenchmarkPowP200(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 200; n++ {
			res = PowP(x, n)
		}
	}
	benchRes = res
}

func BenchmarkPowP200std(b *testing.B) {
	var res float64
	x := 2.5
	var n uint32
	for i := 0; i < b.N; i++ {
		for n = 0; n < 200; n++ {
			res = math.Pow(x, float64(n))
		}
	}
	benchRes = res
}
