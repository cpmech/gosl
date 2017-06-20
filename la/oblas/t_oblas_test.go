// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestMatrix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix01. real")

	A := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 0, -1, -2},
	}

	a := NewMatrix(3, 4)
	a.SetFromMat(A)
	chk.Vector(tst, "A to a", 1e-15, a.data, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})

	Aback := a.GetMat()
	chk.Matrix(tst, "a to A", 1e-15, Aback, A)
}

func TestMatrix02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix02. complex")

	A := [][]complex128{
		{1 + 0.1i, 2, 3, 4 - 0.4i},
		{5 + 0.5i, 6, 7, 8 - 0.8i},
		{9 + 0.9i, 0, -1, -2 + 1i},
	}

	a := NewMatrixC(3, 4)
	a.SetFromMat(A)
	chk.VectorC(tst, "A to a", 1e-15, a.data, []complex128{1 + 0.1i, 5 + 0.5i, 9 + 0.9i, 2, 6, 0, 3, 7, -1, 4 - 0.4i, 8 - 0.8i, -2 + 1i})

	Aback := a.GetMat()
	chk.MatrixC(tst, "a to A", 1e-15, Aback, A)
}

func TestDaxpy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Daxpy01")

	α := 0.5
	x := []float64{20, 10, 30, 123, 123}
	y := []float64{-15, -5, -24, 666, 666, 666}
	n, incx, incy := 3, 1, 1
	err := Daxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}

	chk.Vector(tst, "x", 1e-15, x, []float64{20, 10, 30, 123, 123})
	chk.Vector(tst, "y", 1e-15, y, []float64{-5, 0, -9, 666, 666, 666})
}

func TestZaxpy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zaxpy01")

	α := 1.0 + 0i
	x := []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i}
	y := []complex128{-15 + 1.5i, -5 - 2i, -24 + 1i, 666 - 0.5i, 666 + 5i}
	n, incx, incy := len(x), 1, 1
	err := Zaxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}

	chk.VectorC(tst, "x", 1e-15, x, []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i})
	chk.VectorC(tst, "y", 1e-15, y, []complex128{5 + 2.5i, 5, 6 + 2.5i, 543, 543 + 5.5i})

	α = 0.5 + 1i
	err = Zaxpy(n, α, x, incx, y, incy)
	if err != nil {
		tst.Errorf("Daxpy failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "y", 1e-15, y, []complex128{14.0 + 23.i, 8.0 + 11.i, 19.5 + 33.25i, 481.0 - 122.75i, 481.0 - 117.25i})
}

func TestDgemv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgemv01")

	// allocate
	m, n := 4, 3
	a := NewMatrix(m, n)
	a.SetFromMat([][]float64{
		{0.1, 0.2, 0.3},
		{1.0, 0.2, 0.3},
		{2.0, 0.2, 0.3},
		{3.0, 0.2, 0.3},
	})
	chk.Vector(tst, "a.data", 1e-15, a.data, []float64{0.1, 1, 2, 3, 0.2, 0.2, 0.2, 0.2, 0.3, 0.3, 0.3, 0.3})

	// perform mv
	α, β := 0.5, 2.0
	x := []float64{20, 10, 30}
	y := []float64{3, 1, 2, 4}
	lda, incx, incy := m, 1, 1
	err := Dgemv(false, m, n, α, a, lda, x, incx, β, y, incy)
	if err != nil {
		tst.Errorf("Dgemv failed:\n%v\n", err)
		return
	}
	chk.Vector(tst, "y", 1e-15, y, []float64{12.5, 17.5, 29.5, 43.5})

	// perform mv with transpose
	err = Dgemv(true, m, n, α, a, lda, y, incy, β, x, incx)
	if err != nil {
		tst.Errorf("Dgemv failed:\n%v\n", err)
		return
	}
	chk.Vector(tst, "x", 1e-15, x, []float64{144.125, 30.3, 75.45})

	// check that a is unmodified
	chk.Vector(tst, "a.data", 1e-15, a.data, []float64{0.1, 1, 2, 3, 0.2, 0.2, 0.2, 0.2, 0.3, 0.3, 0.3, 0.3})
}

func TestZgemv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgemv01")

	// allocate
	m, n := 4, 3
	a := NewMatrixC(m, n)
	a.SetFromMat([][]complex128{
		{0.1 + 3i, 0.2, 0.3 - 0.3i},
		{1.0 + 2i, 0.2, 0.3 - 0.4i},
		{2.0 + 1i, 0.2, 0.3 - 0.5i},
		{3.0 + 0.1i, 0.2, 0.3 - 0.6i},
	})
	chk.VectorC(tst, "a.data", 1e-15, a.data, []complex128{0.1 + 3i, 1 + 2i, 2 + 1i, 3 + 0.1i, 0.2, 0.2, 0.2, 0.2, 0.3 - 0.3i, 0.3 - 0.4i, 0.3 - 0.5i, 0.3 - 0.6i})

	// perform mv
	α, β := 0.5+1i, 2.0+1i
	x := []complex128{20, 10, 30}
	y := []complex128{3, 1, 2, 4}
	lda, incx, incy := m, 1, 1
	err := Zgemv(false, m, n, α, a, lda, x, incx, β, y, incy)
	if err != nil {
		tst.Errorf("Zgemv failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "y", 1e-15, y, []complex128{-38.5 + 41.5i, -10.5 + 46i, 24.5 + 55.5i, 59.5 + 67i})

	// perform mv with transpose
	err = Zgemv(true, m, n, α, a, lda, y, incy, β, x, incx)
	if err != nil {
		tst.Errorf("Zgemv failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "x", 1e-13, x, []complex128{-248.875 + 82.5i, -18.5 + 38i, 83.85 + 154.7i})

	// check that a is unmodified
	chk.VectorC(tst, "a.data", 1e-15, a.data, []complex128{0.1 + 3i, 1 + 2i, 2 + 1i, 3 + 0.1i, 0.2, 0.2, 0.2, 0.2, 0.3 - 0.3i, 0.3 - 0.4i, 0.3 - 0.5i, 0.3 - 0.6i})
}

func TestDgesv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgesv01")

	// matrix
	amat := [][]float64{
		{2, 3, 0, 0, 0},
		{3, 0, 4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, 0, 1, 0, 0},
		{0, 4, 2, 0, 1},
	}
	n := 5
	a := NewMatrix(n, n)
	a.SetFromMat(amat)

	// right-hand-side
	b := []float64{8, 45, -3, 3, 19}

	// solution
	xCorrect := []float64{1, 2, 3, 4, 5}

	// run test
	nrhs := 1
	lda, ldb := n, n
	ipiv := make([]int, n)
	err := Dgesv(n, nrhs, a, lda, ipiv, b, ldb)
	if err != nil {
		tst.Errorf("Dgesv failed:\n%v\n", err)
		return
	}
	chk.Vector(tst, "x = A⁻¹ b", 1e-15, b, xCorrect)
}

func TestZgesv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Zgesv01. low accuracy")

	// NOTE: zgesv performs badly with this problem
	//       the best tolerance that can be selected is 0.00032
	//       the same problem happens in python (probably using lapack as well)
	tol := 0.0032

	// matrix
	n := 5
	a := NewMatrixC(n, n)
	a.SetFromMat([][]complex128{
		{19.730 + 0.000i, 12.110 - 1.000i, 0.000 + 5.000i, 0.000 + 0.000i, 0.000 + 0.000i},
		{0.000 - 0.510i, 32.300 + 7.000i, 23.070 + 0.000i, 0.000 + 1.000i, 0.000 + 0.000i},
		{0.000 + 0.000i, 0.000 - 0.510i, 70.000 + 7.300i, 3.950 + 0.000i, 19.000 + 31.830i},
		{0.000 + 0.000i, 0.000 + 0.000i, 1.000 + 1.100i, 50.170 + 0.000i, 45.510 + 0.000i},
		{0.000 + 0.000i, 0.000 + 0.000i, 0.000 + 0.000i, 0.000 - 9.351i, 55.000 + 0.000i},
	})

	// right-hand-side
	b := []complex128{
		77.38 + 8.82i,
		157.48 + 19.8i,
		1175.62 + 20.69i,
		912.12 - 801.75i,
		550 - 1060.4i,
	}

	// solution
	xCorrect := []complex128{
		3.3 - 1i,
		1 + 0.17i,
		5.5,
		9,
		10 - 17.75i,
	}

	// run test
	nrhs := 1
	lda, ldb := n, n
	ipiv := make([]int, n)
	err := Zgesv(n, nrhs, a, lda, ipiv, b, ldb)
	if err != nil {
		tst.Errorf("Zgesv failed:\n%v\n", err)
		return
	}
	chk.VectorC(tst, "x = A⁻¹ b", tol, b, xCorrect)

	// compare with python results
	xPython := []complex128{
		3.299687426933794e+00 - 1.000372829305209e+00i,
		9.997606020636992e-01 + 1.698383755401385e-01i,
		5.500074759292877e+00 - 4.556001293922560e-05i,
		8.999787912842375e+00 - 6.662818244209770e-05i,
		1.000001132800243e+01 - 1.774987242230929e+01i,
	}
	chk.VectorC(tst, "x = A⁻¹ b", 1e-14, b, xPython)
}

func checksvd(tst *testing.T, amat, uCorrect, vtCorrect [][]float64, sCorrect []float64, tolu, tols, tolv, tolusv float64) {

	// allocate matrix
	m, n := len(amat), len(amat[0])
	a := NewMatrix(m, n)
	a.SetFromMat(amat)

	// compute dimensions
	minMN := imin(m, n)
	maxMN := imax(m, n)
	lda := m
	ldu := m
	ldvt := n
	lwork := 2 * imax(3*minMN+maxMN, 5*minMN)

	// allocate output arrays
	s := make([]float64, minMN)
	u := NewMatrix(m, m)
	vt := NewMatrix(n, n)
	work := make([]float64, lwork)

	// perform SVD
	jobu := 'A'
	jobvt := 'A'
	err := Dgesvd(jobu, jobvt, m, n, a, lda, s, u, ldu, vt, ldvt, work, lwork)
	if err != nil {
		tst.Errorf("Dgesvd failed:\n%v\n", err)
		return
	}

	// compare results
	umat := u.GetMat()
	vtmat := vt.GetMat()
	chk.Matrix(tst, "u", tolu, umat, uCorrect)
	chk.Vector(tst, "s", tols, s, sCorrect)
	chk.Matrix(tst, "vt", tolv, vtmat, vtCorrect)

	// check SVD
	usv := make([][]float64, m)
	for i := 0; i < m; i++ {
		usv[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			for k := 0; k < minMN; k++ {
				usv[i][j] += umat[i][k] * s[k] * vtmat[k][j]
			}
		}
	}
	chk.Matrix(tst, "u⋅s⋅vt", tolusv, amat, usv)
}

func TestDgesvd01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgesvd01")

	// allocate matrices
	amat := [][]float64{
		{1, 0, 0, 0, 2},
		{0, 0, 3, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 2, 0, 0, 0},
	}
	uCorrect := [][]float64{
		{0, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 0, -1},
		{0, 0, 1, 0},
	}
	sCorrect := []float64{3, math.Sqrt(5.0), 2, 0}
	s2 := math.Sqrt(0.2)
	s8 := math.Sqrt(0.8)
	vtCorrect := [][]float64{
		{0, 0, 1, 0, 0},
		{s2, 0, 0, 0, s8},
		{0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0},
		{-s8, 0, 0, 0, s2},
	}

	// check
	checksvd(tst, amat, uCorrect, vtCorrect, sCorrect, 1e-17, 1e-17, 1e-15, 1e-15)
}

func TestDgesvd02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgesvd02")

	// allocate matrices
	s33 := math.Sqrt(3.0) / 3.0
	amat := [][]float64{
		{-s33, -s33, 1},
		{+s33, -s33, 1},
		{-s33, +s33, 1},
		{+s33, +s33, 1},
	}
	uCorrect := [][]float64{
		{-0.5, -0.5, -0.5, +0.5},
		{-0.5, -0.5, +0.5, -0.5},
		{-0.5, +0.5, -0.5, -0.5},
		{-0.5, +0.5, +0.5, +0.5},
	}
	sCorrect := []float64{2, 2.0 / math.Sqrt(3.0), 2.0 / math.Sqrt(3.0)}
	vtCorrect := [][]float64{
		{+0, +0, -1},
		{+0, +1, +0},
		{+1, +0, +0},
	}

	// check
	checksvd(tst, amat, uCorrect, vtCorrect, sCorrect, 1e-15, 1e-15, 1e-17, 1e-15)
}

func TestDgesvd03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dgesvd03")

	// allocate matrices
	amat := [][]float64{
		{64, 2, 3, 61, 60, 6},
		{9, 55, 54, 12, 13, 51},
		{17, 47, 46, 20, 21, 43},
		{40, 26, 27, 37, 36, 30},
		{32, 34, 35, 29, 28, 38},
		{41, 23, 22, 44, 45, 19},
		{49, 15, 14, 52, 53, 11},
		{8, 58, 59, 5, 4, 62},
	}
	uCorrect := [][]float64{
		{-3.554400501038920e-01, +5.585242303027516e-01, +3.215271708791157e-01, +6.314751632872790e-01, +1.516793794931194e-01, +1.726099851944921e-01, +3.235158568524713e-02, +7.566363218896498e-02},
		{-3.516775709543288e-01, -4.047465278940223e-01, -3.336355167767969e-01, +3.614553362267942e-02, +4.523700594334092e-01, +6.149586060459316e-01, +5.554602167528225e-02, +1.182618207365148e-01},
		{-3.516977917841346e-01, -2.506762497662257e-01, -3.421132994034763e-01, +3.699581516051667e-01, -7.141640069622246e-01, -8.386957034656470e-02, -3.626057557500124e-02, +2.029294815397017e-01},
		{-3.553793876144748e-01, +9.631339591936192e-02, +3.469605187591548e-01, -4.983828697683003e-01, -1.344898470243124e-01, +1.247121364688733e-01, -4.702500192002790e-01, +4.908018109950452e-01},
		{-3.553591667846690e-01, -5.775688220843472e-02, +3.554383013858351e-01, -2.818045094280957e-01, -3.323582511042328e-01, +2.822560889638348e-01, +4.092001768165193e-01, -5.541286366654536e-01},
		{-3.517584542735519e-01, +2.115345846171643e-01, -3.675466472835166e-01, -7.917913888059364e-02, +1.215272457237375e-01, -1.854280397149560e-01, -5.596262526312354e-01, -5.725602705348066e-01},
		{-3.517786751033576e-01, +3.656048627449609e-01, -3.760244299101969e-01, -3.269245463472517e-01, +1.402667018050773e-01, -3.456609959844107e-01, +5.403408065309545e-01, +2.513689682585900e-01},
		{-3.552985042952517e-01, -5.199677165918247e-01, +3.808716492658746e-01, +1.487122159091161e-01, +3.151687186354263e-01, -5.795782106272005e-01, +2.869825669851275e-02, -1.233680651855640e-02},
	}
	sCorrect := []float64{
		+2.251695779937001e+02, +1.271865289052834e+02, +1.175789144211322e+01, +1.277237188369868e-14, +6.934703857768031e-15, +5.031833747507930e-15}
	vtCorrect := [][]float64{
		{-4.084940479369395e-01, -4.080456032641434e-01, -4.081102861436546e-01, -4.082999992984062e-01, -4.082353164188952e-01, -4.083043347821876e-01},
		{+4.109984470006984e-01, -4.103747435987312e-01, -4.091529447253958e-01, +4.073330503806908e-01, +4.061112515073549e-01, -4.054875481053879e-01},
		{+5.582551269795966e-01, -3.983479622927922e-01, -1.581588478101714e-01, -1.623122164682691e-01, -4.025013309508903e-01, +5.624084956376931e-01},
		{+5.291884225014455e-01, +2.958619856107400e-01, +8.554456957353419e-02, -6.766720388422237e-01, +1.474836163407778e-01, -3.814065551842739e-01},
		{+2.567084227345682e-01, +4.201629641184643e-01, -7.793403525625409e-02, +4.200160597668203e-01, -6.767244825013885e-01, -3.422289288622102e-01},
		{+8.139104220774948e-02, -4.922794180787097e-01, +7.922335853716356e-01, +8.201891496907093e-02, -1.634099571768202e-01, -2.999541672929260e-01},
	}

	// check
	checksvd(tst, amat, uCorrect, vtCorrect, sCorrect, 1e-15, 1e-13, 1e-15, 1e-13)
}
