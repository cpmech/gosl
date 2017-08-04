// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math/cmplx"
	"testing"
)

var (
	benchmarkingResC complex128
)

func BenchmarkImagPowN(b *testing.B) {
	var res complex128
	for i := 0; i < b.N; i++ {
		for n := 0; n < 200; n++ {
			res = ImagPowN(n)
		}
	}
	benchmarkingResC = res
}

func BenchmarkImagPowNcmplx(b *testing.B) {
	var res complex128
	for i := 0; i < b.N; i++ {
		for n := 0; n < 200; n++ {
			res = cmplx.Pow(1i, complex(float64(n), 0))
		}
	}
	benchmarkingResC = res
}

func BenchmarkImagXpowN(b *testing.B) {
	var res complex128
	x := 2.5
	for i := 0; i < b.N; i++ {
		for n := 0; n < 200; n++ {
			res = ImagXpowN(x, n)
		}
	}
	benchmarkingResC = res
}

func BenchmarkImagXpowNcmplx(b *testing.B) {
	var res complex128
	x := 2.5
	for i := 0; i < b.N; i++ {
		for n := 0; n < 200; n++ {
			res = cmplx.Pow(complex(0, x), complex(float64(n), 0))
		}
	}
	benchmarkingResC = res
}
