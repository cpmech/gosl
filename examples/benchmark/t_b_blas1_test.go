// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cpmech/gosl/la"
)

var (
	u4, v4, w4          []float64
	u32, v32, w32       []float64
	u100, v100, w100    []float64
	u200, v200, w200    []float64
	u500, v500, w500    []float64
	u1000, v1000, w1000 []float64
	benchRes            float64
)

const beta = 1.0

func init() {
	rand.Seed(int64(time.Now().Unix()))
	u4 = make([]float64, 4)
	v4 = make([]float64, 4)
	w4 = make([]float64, 4)
	u32 = make([]float64, 32)
	v32 = make([]float64, 32)
	w32 = make([]float64, 32)
	u100 = make([]float64, 100)
	v100 = make([]float64, 100)
	w100 = make([]float64, 100)
	u200 = make([]float64, 200)
	v200 = make([]float64, 200)
	w200 = make([]float64, 200)
	u500 = make([]float64, 500)
	v500 = make([]float64, 500)
	w500 = make([]float64, 500)
	u1000 = make([]float64, 1000)
	v1000 = make([]float64, 1000)
	w1000 = make([]float64, 1000)
	for i := 0; i < len(u4); i++ {
		u4[i] = rand.Float64()
		v4[i] = rand.Float64()
	}
	for i := 0; i < len(u32); i++ {
		u32[i] = rand.Float64()
		v32[i] = rand.Float64()
	}
	for i := 0; i < len(u200); i++ {
		u200[i] = rand.Float64()
		v200[i] = rand.Float64()
	}
	for i := 0; i < len(u100); i++ {
		u100[i] = rand.Float64()
		v100[i] = rand.Float64()
	}
	for i := 0; i < len(u500); i++ {
		u500[i] = rand.Float64()
		v500[i] = rand.Float64()
	}
	for i := 0; i < len(u1000); i++ {
		u1000[i] = rand.Float64()
		v1000[i] = rand.Float64()
	}
}

// VecDot /////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkNaiveVecDot4(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += NaiveVecDot(u4, v4)
	}
	benchRes = res
}

func BenchmarkVecDot4(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += la.VecDot(u4, v4)
	}
	benchRes = res
}

func BenchmarkNaiveVecDot32(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += NaiveVecDot(u32, v32)
	}
	benchRes = res
}

func BenchmarkVecDot32(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += la.VecDot(u32, v32)
	}
	benchRes = res
}

func BenchmarkNaiveVecDot100(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += NaiveVecDot(u100, v100)
	}
	benchRes = res
}

func BenchmarkVecDot100(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += la.VecDot(u100, v100)
	}
	benchRes = res
}

func BenchmarkNaiveVecDot200(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += NaiveVecDot(u200, v200)
	}
	benchRes = res
}

func BenchmarkVecDot200(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += la.VecDot(u200, v200)
	}
	benchRes = res
}

func BenchmarkNaiveVecDot500(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += NaiveVecDot(u500, v500)
	}
	benchRes = res
}

func BenchmarkVecDot500(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += la.VecDot(u500, v500)
	}
	benchRes = res
}

func BenchmarkNaiveVecDot1000(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += NaiveVecDot(u1000, v1000)
	}
	benchRes = res
}

func BenchmarkVecDot1000(b *testing.B) {
	var res float64
	for i := 0; i < b.N; i++ {
		res += la.VecDot(u1000, v1000)
	}
	benchRes = res
}

// VecAdd /////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkNaiveVecAdd4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NaiveVecAdd(w4, 1, u4, beta, v4)
	}
}

func BenchmarkVecAdd4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		la.VecAdd(w4, 1, u4, beta, v4)
	}
}

func BenchmarkNaiveVecAdd32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NaiveVecAdd(w32, 1, u32, beta, v32)
	}
}

func BenchmarkVecAdd32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		la.VecAdd(w32, 1, u32, beta, v32)
	}
}

func BenchmarkNaiveVecAdd100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NaiveVecAdd(w100, 1, u100, beta, v100)
	}
}

func BenchmarkVecAdd100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		la.VecAdd(w100, 1, u100, beta, v100)
	}
}

func BenchmarkNaiveVecAdd200(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NaiveVecAdd(w200, 1, u200, beta, v200)
	}
}

func BenchmarkVecAdd200(b *testing.B) {
	for i := 0; i < b.N; i++ {
		la.VecAdd(w200, 1, u200, beta, v200)
	}
}

func BenchmarkNaiveVecAdd500(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NaiveVecAdd(w500, 1, u500, beta, v500)
	}
}

func BenchmarkVecAdd500(b *testing.B) {
	for i := 0; i < b.N; i++ {
		la.VecAdd(w500, 1, u500, beta, v500)
	}
}

func BenchmarkNaiveVecAdd1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NaiveVecAdd(w1000, 1, u1000, beta, v1000)
	}
}

func BenchmarkVecAdd1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		la.VecAdd(w1000, 1, u1000, beta, v1000)
	}
}
