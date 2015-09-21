// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

var (
	acorr0 [][]float64
	acorr1 [][]float64
	bmat   [][]float64
	dmat   [][]float64
	ucorr0 []float64
	ucorZ0 []complex128
	bvec   []float64
	bveZ   []complex128
	dvec   []float64
	dveZ   []complex128
	evec   []float64
	fvec   []float64
	amat   = [][]float64{
		{1.01, 2.01, 3.01, 4.01, 5.01, 6.01, 7.01, 8.01},
		{1.02, 2.02, 3.02, 4.02, 5.02, 6.02, 7.02, 8.02},
		{1.03, 2.03, 3.03, 4.03, 5.03, 6.03, 7.03, 8.03},
	}
	avec = []float64{5.01, 5.02, 5.03}
)

const (
	tol   = 1.0e-16
	noise = 1.0e-17
)

func allocate(N, n int) {
	acorr0 = MatAlloc(10, n)
	acorr1 = MatAlloc(10, n)
	bmat = MatAlloc(10, n)
	dmat = MatAlloc(10, n)
	for i := 0; i < 10; i++ {
		for j := 0; j < n; j++ {
			bmat[i][j] = float64(i + j)
			acorr0[i][j] = 666.0
			acorr1[i][j] = 666.0 / 2.0
			if i == j {
				dmat[i][j] = 1.5
			}
		}
	}
	ucorr0 = make([]float64, N)
	ucorZ0 = make([]complex128, N)
	bvec = make([]float64, N)
	bveZ = make([]complex128, N)
	dvec = make([]float64, N)
	dveZ = make([]complex128, N)
	evec = make([]float64, N)
	fvec = make([]float64, N)
	for i := 0; i < N; i++ {
		ucorr0[i] = 666.0
		ucorZ0[i] = 666.0 + 666.0i
		bvec[i] = float64(i + 1)
		bveZ[i] = complex(float64(i+1), 0)
		dvec[i] = float64(N - i)
		dveZ[i] = complex(float64(N-i), 0)
		evec[i] = bvec[i] / 2.0
		fvec[i] = float64(N+1) / 2.0
	}
}

func TestLinAlg01(tst *testing.T) {

	//verbose()
	N, n := 100, 100
	allocate(N, n)

	chk.PrintTitle("TestLinAlg 01 (serial)")
	Pll = false
	run_tests(N, n, true, tst)

	chk.PrintTitle("TestLinAlg 01 (parallel)")
	Pll = true
	defer func() { Pll = false }()
	_, msg := run_tests(N, n, true, tst)

	chk.PrintTitle("TestLinAlg 02 (speed)")

	N, n = 10, 10
	//N, n = 100, 100
	//N, n = 1000, 1000
	//N, n = 10000, 10000
	//N, n = 100000, 100000
	//N, n = 1000000, 100000
	//N, n = 10000000, 100000
	allocate(N, n)

	nruns := 100
	Dt_ser, Dt_pll := make([][]time.Duration, nruns), make([][]time.Duration, nruns)
	for irun := 0; irun < nruns; irun++ {
		Pll = false
		Dt_ser[irun], _ = run_tests(N, n, false, tst)
		Pll = true
		Dt_pll[irun], _ = run_tests(N, n, false, tst)
	}

	io.Pf("%s  %20s |  %20s | %10s %10s\n", "                  ", "serial    ", "parallel    ", "serial", "parallel")
	io.Pf("%s %10s %10s | %10s %10s | %10s %10s\n", "                  ", "min", "max", "min", "max", "ave", "ave")
	for j := 0; j < len(msg); j++ {
		min_s, max_s, ave_s := durminmaxave(durgetcol(j, Dt_ser))
		min_p, max_p, ave_p := durminmaxave(durgetcol(j, Dt_pll))
		io.Pf("%s %10v %10v | %10v %10v | %10v %10v\n", msg[j], min_s, max_s, min_p, max_p, ave_s, ave_p)
	}
}

func run_tests(N, n int, do_check bool, tst *testing.T) (dt []time.Duration, msg []string) {

	// allocate
	num_tests := 25
	dt = make([]time.Duration, num_tests)
	msg = make([]string, num_tests)

	// 0: check if v is filled with zero by default
	t0 := time.Now()
	v := make([]float64, N)
	v[0] += noise
	dt[0], msg[0] = time.Now().Sub(t0), "vector: alloc     "
	if do_check {
		chk.Vector(tst, msg[0], tol, v, []float64{})
	}

	// 1: matrix: allocate
	t0 = time.Now()
	a := MatAlloc(10, n)
	a[0][0] += noise
	dt[1], msg[1] = time.Now().Sub(t0), "matrix: alloc     "
	if do_check {
		chk.Matrix(tst, msg[1], tol, a, [][]float64{})
	}

	// 2: vector: fill
	u := make([]float64, N)
	t0 = time.Now()
	VecFill(u, 666.0)
	u[0] += 100 * noise
	dt[2], msg[2] = time.Now().Sub(t0), "vector: fill      "
	if do_check {
		chk.Vector(tst, msg[2], tol, u, ucorr0)
	}

	// 3: matrix: fill
	t0 = time.Now()
	MatFill(a, 666.0)
	a[0][0] += 100 * noise
	dt[3], msg[3] = time.Now().Sub(t0), "matrix: fill      "
	if do_check {
		chk.Matrix(tst, msg[3], tol, a, acorr0)
	}

	// 4: vector: norm (Euclidian)
	v[0], v[N-1] = 1.0/math.Sqrt(2.0), -1.0/math.Sqrt(2.0)
	t0 = time.Now()
	nrm := VecNorm(v)
	dt[4], msg[4] = time.Now().Sub(t0), "vector: norm      "
	if do_check {
		chk.Scalar(tst, msg[4], tol*10, nrm, 1.0+noise)
	}

	// 5: vector: min component
	v[N-1], v[1] = 1.0, -1.0
	t0 = time.Now()
	min := VecMin(v)
	dt[5], msg[5] = time.Now().Sub(t0), "vector: min       "
	if do_check {
		chk.Scalar(tst, msg[5], tol, min, -1.0+noise)
	}

	// 6: vector: max component
	t0 = time.Now()
	max := VecMax(v)
	dt[6], msg[6] = time.Now().Sub(t0), "vector: max       "
	if do_check {
		chk.Scalar(tst, msg[6], tol, max, 1.0+noise)
	}

	// 7: vector: min and max components
	t0 = time.Now()
	min, max = VecMinMax(v)
	dt[7], msg[7] = time.Now().Sub(t0), "vector: minmax    "
	if do_check {
		chk.Scalar(tst, msg[7], tol, min, -1.0+noise)
		chk.Scalar(tst, msg[7], tol, max, 1.0+noise)
	}

	// 8: matrix: a := a * s
	t0 = time.Now()
	MatScale(a, 0.5)
	dt[8], msg[8] = time.Now().Sub(t0), "matrix: scale     "
	if do_check {
		chk.Matrix(tst, msg[8], tol, a, acorr1)
	}

	// 9: vector: a := b
	t0 = time.Now()
	VecCopy(v, 1, u) // v := u
	v[0] += noise
	dt[9], msg[9] = time.Now().Sub(t0), "vector: copy      "
	if do_check {
		chk.Vector(tst, msg[9], tol, v, u)
	}

	// 10: vector: a += b * s
	t0 = time.Now()
	VecAdd(u, -1, v) // u += (-1.0)*v
	u[0] += noise
	dt[10], msg[10] = time.Now().Sub(t0), "vector: addscaled "
	if do_check {
		chk.Vector(tst, msg[10], tol, u, []float64{})
	}

	// 11: matrix: a := b * s
	t0 = time.Now()
	MatCopy(a, 1, bmat) // a := 1*bmat
	a[0][0] += noise
	dt[11], msg[11] = time.Now().Sub(t0), "matrix: copyscaled"
	if do_check {
		chk.Matrix(tst, msg[11], tol, a, bmat)
	}

	// 12: matrix: set diagonal with v
	t0 = time.Now()
	MatSetDiag(a, 1.5)
	a[0][0] += noise
	dt[12], msg[12] = time.Now().Sub(t0), "matrix: setdiag   "
	if do_check {
		chk.Matrix(tst, msg[12], tol, a, dmat)
	}

	// 13: vector: max difference
	t0 = time.Now()
	maxdiff := VecMaxDiff(bvec, dvec)
	dt[13], msg[13] = time.Now().Sub(t0), "vector: maxdiff   "
	if do_check {
		chk.Scalar(tst, msg[13], tol, maxdiff, float64(N-1)+100*noise)
	}

	// 14: matrix: max difference
	t0 = time.Now()
	maxdiff = MatMaxDiff(acorr0, acorr1)
	dt[14], msg[14] = time.Now().Sub(t0), "matrix: maxdiff   "
	if do_check {
		chk.Scalar(tst, msg[14], tol, maxdiff, 333.0+noise)
	}

	// 15: matrix: get col
	t0 = time.Now()
	col := MatGetCol(4, amat)
	dt[15], msg[15] = time.Now().Sub(t0), "matrix: getcol    "
	if do_check {
		chk.Vector(tst, msg[15], tol, col, avec)
	}

	// 16: vector: a := b * s
	t0 = time.Now()
	VecCopy(v, 0.5, bvec) // v := u * 0.5
	v[0] += noise
	dt[16], msg[16] = time.Now().Sub(t0), "vector: copyscaled"
	if do_check {
		chk.Vector(tst, msg[16], tol, v, evec)
	}

	// 17: vector: add: u := α*a + β*b
	t0 = time.Now()
	VecAdd2(u, 0.5, bvec, 0.5, dvec) // u := 0.5*bvec + 0.5*dvec
	v[0] += noise
	dt[17], msg[17] = time.Now().Sub(t0), "vector: add       "
	if do_check {
		chk.Vector(tst, msg[17], tol, u, fvec)
	}

	// 18: vector: dot product: s := u dot v
	t0 = time.Now()
	res := VecDot(bvec, bvec)
	res += noise
	dt[18], msg[18] = time.Now().Sub(t0), "vector: dot       "
	if do_check {
		chk.Scalar(tst, msg[18], tol, res, float64(N*(N+1)*(2*N+1)/6)) // Faulhaber's formula
	}

	// 19: vector: largest component (divided by den)
	t0 = time.Now()
	res = VecLargest(bvec, 2.0) // den = 1.0
	res += noise
	dt[19], msg[19] = time.Now().Sub(t0), "vector: largest   "
	if do_check {
		chk.Scalar(tst, msg[19], tol, res, float64(N)/2.0)
	}

	// 20: vectorC: fill
	uc := make([]complex128, N)
	t0 = time.Now()
	VecFillC(uc, 666+666i)
	dt[20], msg[20] = time.Now().Sub(t0), "vector: fillC     "
	if do_check {
		chk.VectorC(tst, msg[20], tol, uc, ucorZ0)
	}

	// 21: vector: max difference C
	t0 = time.Now()
	maxdiffC := VecMaxDiffC(bveZ, dveZ)
	dt[21], msg[21] = time.Now().Sub(t0), "vector: maxdiffC  "
	if do_check {
		chk.Scalar(tst, msg[21], tol, maxdiffC, float64(N-1)+100*noise)
	}

	// 22: vector: rms
	t0 = time.Now()
	sum_1toN_squared := float64(N*(N+1)*(2*N+1)) / float64(6) // Faulhaber's formula
	rms_corr := math.Sqrt(sum_1toN_squared / float64(N))
	rms := VecRms(bvec)
	dt[22], msg[22] = time.Now().Sub(t0), "vector: rms       "
	if do_check {
		chk.Scalar(tst, msg[22], 1e-17, rms, rms_corr+noise)
	}

	// 23: vector: rmserr
	t0 = time.Now()
	rmserr := VecRmsErr(v, 0, 1, v)
	dt[23], msg[23] = time.Now().Sub(t0), "vector: rmserr    "
	if do_check {
		chk.Scalar(tst, msg[23], 1e-17, rmserr, 1.0+noise)
	}

	// 24: matrix: largest
	t0 = time.Now()
	largest := MatLargest(acorr0, 1)
	dt[24], msg[24] = time.Now().Sub(t0), "matrix: largest   "
	if do_check {
		chk.Scalar(tst, msg[24], tol, largest, 666.0)
	}

	return
}

// durgetcol returns a columns with duration values
func durgetcol(j int, a [][]time.Duration) (col []time.Duration) {
	col = make([]time.Duration, len(a))
	for i := 0; i < len(a); i++ {
		col[i] = a[i][j]
	}
	return
}

// durminmaxave returns statistics data correponding a duration data collected in 'v'
func durminmaxave(v []time.Duration) (min, max, ave time.Duration) {
	min, max, ave = v[0], v[0], v[0]
	for i := 1; i < len(v); i++ {
		if v[i] < min {
			min = v[i]
		}
		if v[i] > max {
			max = v[i]
		}
		ave += v[i]
	}
	ave /= time.Duration(len(v)) * time.Nanosecond
	return
}
