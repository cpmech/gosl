// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math/cmplx"
	"testing"
)

func BenchmarkImagPowN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < 200; n++ {
			ImagPowN(n)
		}
	}
}

func BenchmarkImagPowNcmplx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < 200; n++ {
			cmplx.Pow(1i, complex(float64(n), 0))
		}
	}
}

func BenchmarkImagXpowN(b *testing.B) {
	x := 2.5
	for i := 0; i < b.N; i++ {
		for n := 0; n < 200; n++ {
			ImagXpowN(x, n)
		}
	}
}

func BenchmarkImagXpowNcmplx(b *testing.B) {
	x := 2.5
	for i := 0; i < b.N; i++ {
		for n := 0; n < 200; n++ {
			cmplx.Pow(complex(0, x), complex(float64(n), 0))
		}
	}
}
