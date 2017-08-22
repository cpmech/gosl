// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math/rand"
	"testing"
)

var (
	benchmarkMatricesA []*Matrix
	benchmarkMatricesB []*Matrix
	benchmarkMatricesC []*Matrix
	benchmarkVectorsU  []Vector
	benchmarkVectorsV  []Vector
)

func init() {
	N := 12
	benchmarkMatricesA = make([]*Matrix, N+1)
	benchmarkMatricesB = make([]*Matrix, N+1)
	benchmarkMatricesC = make([]*Matrix, N+1)
	benchmarkVectorsU = make([]Vector, N+1)
	benchmarkVectorsV = make([]Vector, N+1)
	for i := 1; i <= N; i++ {
		benchmarkMatricesA[i] = NewMatrix(i, i)
		benchmarkMatricesB[i] = NewMatrix(i, i)
		benchmarkMatricesC[i] = NewMatrix(i, i)
		benchmarkVectorsU[i] = NewVector(i)
		benchmarkVectorsV[i] = NewVector(i)
		for m := 0; m < i; m++ {
			for n := 0; n < i; n++ {
				benchmarkMatricesA[i].Set(m, n, rand.Float64())
				benchmarkMatricesB[i].Set(m, n, rand.Float64())
			}
			benchmarkVectorsU[i][m] = rand.Float64()
		}
	}
}

// MatVecMul //////////////////////////////////////////////////////////////////////////////////////

func BenchmarkMatVecMul01(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[1], 1, benchmarkMatricesA[1], benchmarkVectorsU[1])
	}
}
func BenchmarkMatVecMul02(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[2], 1, benchmarkMatricesA[2], benchmarkVectorsU[2])
	}
}
func BenchmarkMatVecMul03(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[3], 1, benchmarkMatricesA[3], benchmarkVectorsU[3])
	}
}
func BenchmarkMatVecMul04(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[4], 1, benchmarkMatricesA[4], benchmarkVectorsU[4])
	}
}
func BenchmarkMatVecMul05(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[5], 1, benchmarkMatricesA[5], benchmarkVectorsU[5])
	}
}
func BenchmarkMatVecMul06(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[6], 1, benchmarkMatricesA[6], benchmarkVectorsU[6])
	}
}
func BenchmarkMatVecMul07(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[7], 1, benchmarkMatricesA[7], benchmarkVectorsU[7])
	}
}
func BenchmarkMatVecMul08(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[8], 1, benchmarkMatricesA[8], benchmarkVectorsU[8])
	}
}
func BenchmarkMatVecMul09(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[9], 1, benchmarkMatricesA[9], benchmarkVectorsU[9])
	}
}
func BenchmarkMatVecMul10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[10], 1, benchmarkMatricesA[10], benchmarkVectorsU[10])
	}
}
func BenchmarkMatVecMul11(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[11], 1, benchmarkMatricesA[11], benchmarkVectorsU[11])
	}
}
func BenchmarkMatVecMul12(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatVecMul(benchmarkVectorsV[12], 1, benchmarkMatricesA[12], benchmarkVectorsU[12])
	}
}

// MatMatMul //////////////////////////////////////////////////////////////////////////////////////

func BenchmarkMatMatMul01(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[1], 1, benchmarkMatricesA[1], benchmarkMatricesB[1])
	}
}
func BenchmarkMatMatMul02(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[2], 1, benchmarkMatricesA[2], benchmarkMatricesB[2])
	}
}
func BenchmarkMatMatMul03(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[3], 1, benchmarkMatricesA[3], benchmarkMatricesB[3])
	}
}
func BenchmarkMatMatMul04(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[4], 1, benchmarkMatricesA[4], benchmarkMatricesB[4])
	}
}
func BenchmarkMatMatMul05(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[5], 1, benchmarkMatricesA[5], benchmarkMatricesB[5])
	}
}
func BenchmarkMatMatMul06(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[6], 1, benchmarkMatricesA[6], benchmarkMatricesB[6])
	}
}
func BenchmarkMatMatMul07(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[7], 1, benchmarkMatricesA[7], benchmarkMatricesB[7])
	}
}
func BenchmarkMatMatMul08(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[8], 1, benchmarkMatricesA[8], benchmarkMatricesB[8])
	}
}
func BenchmarkMatMatMul09(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[9], 1, benchmarkMatricesA[9], benchmarkMatricesB[9])
	}
}
func BenchmarkMatMatMul10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[10], 1, benchmarkMatricesA[10], benchmarkMatricesB[10])
	}
}
func BenchmarkMatMatMul11(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[11], 1, benchmarkMatricesA[11], benchmarkMatricesB[11])
	}
}
func BenchmarkMatMatMul12(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMatMul(benchmarkMatricesC[12], 1, benchmarkMatricesA[12], benchmarkMatricesB[12])
	}
}
