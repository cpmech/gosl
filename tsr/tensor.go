// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// tsr implements routines to conduct tensor operations
package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

const (
	EPS        = 1e-16 // smallest number satisfying 1.0 + EPS > 1.0
	QMIN       = 1e-10 // smallest q value to compute qCam invariant
	MINDET     = 1e-30 // minimum determinant of tensor
	SMPINVSTOL = 1e-8  // tolerance used in SmpInvs to avoid sqrt(negativenumber)
	EV_DEBUG   = false // flag to activate debugging of eivenvalues/projectors
	EV_DNMIN   = 1e-10 // minimum denominator to be used in analytical eigenprojectors computation
	EV_ALPMIN  = 1e-12 // minimum α to be used in eigenprojectors derivatives
	EV_PERT    = 1e-3  // perturbation value
	EV_EVTOL   = 1e-6  // mfac coefficient
	EV_ZERO    = 1e-8  // minimum eigenvalue
)

var (
	// T2MI converts i-j-indices of 3x3 2nd order tensor to I-index in Mandel's representation
	T2MI = [][]int{
		{0, 3, 5},
		{3, 1, 4},
		{5, 4, 2},
	}

	// TT2MI converts i-j-k-l-indices of 3x3x3x3 4th order tensor to I-index in Mandel's representation
	TT2MI = [][][][]int{
		{{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}, {{3, 3, 3}, {3, 3, 3}, {3, 3, 3}}, {{5, 5, 5}, {5, 5, 5}, {5, 5, 5}}},
		{{{3, 3, 3}, {3, 3, 3}, {3, 3, 3}}, {{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}, {{4, 4, 4}, {4, 4, 4}, {4, 4, 4}}},
		{{{5, 5, 5}, {5, 5, 5}, {5, 5, 5}}, {{4, 4, 4}, {4, 4, 4}, {4, 4, 4}}, {{2, 2, 2}, {2, 2, 2}, {2, 2, 2}}},
	}

	// TT2MJ converts i-j-k-l-indices of 3x3x3x3 4th order tensor to J-index in Mandel's representation
	TT2MJ = [][][][]int{
		{{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}, {{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}, {{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}},
		{{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}, {{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}, {{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}},
		{{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}, {{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}, {{0, 3, 5}, {3, 1, 4}, {5, 4, 2}}},
	}

	// M2Ti converts I-index in Mandel's representation to i-index of 3x3 2nd order tensor
	M2Ti = []int{0, 1, 2, 0, 1, 0}

	// M2Tj converts I-index in Mandel's representation to j-index of 3x3 2nd order tensor
	M2Tj = []int{0, 1, 2, 1, 2, 2}

	// constants
	SQ2    = math.Sqrt(2.0)       // sqrt(2)
	SQ3    = math.Sqrt(3.0)       // sqrt(3)
	SQ6    = math.Sqrt(6.0)       // sqrt(6)
	SQ3by2 = math.Sqrt(3.0 / 2.0) // sqrt(3/2)
	SQ2by3 = math.Sqrt(2.0 / 3.0) // sqrt(2/3)
	TWOSQ2 = 2.0 * math.Sqrt(2.0) // 2*sqrt(2) == 2^(3/2)

	// 3x3 2nd order identity tensor
	It = [][]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}

	// 3x3 2nd order identity tensor in Mandel's representation
	Im = []float64{1, 1, 1, 0, 0, 0}

	// 4th order identity tensor (symmetric) in Mandel's representation
	IIm = [][]float64{
		{1.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{0.0, 1.0, 0.0, 0.0, 0.0, 0.0},
		{0.0, 0.0, 1.0, 0.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 1.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 1.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 0.0, 1.0},
	}

	// symmetric-deviatoric projector (3D) in Mandel's representation
	//Psd3dm = [][]float64{
	Psd = [][]float64{
		{2.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0, 0.0, 0.0, 0.0},
		{-1.0 / 3.0, 2.0 / 3.0, -1.0 / 3.0, 0.0, 0.0, 0.0},
		{-1.0 / 3.0, -1.0 / 3.0, 2.0 / 3.0, 0.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 1.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 1.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 0.0, 1.0},
	}

	// isotropic projector (3D) in Mandel's representation
	//Piso3dm = [][]float64{
	Piso = [][]float64{
		{1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0, 0.0, 0.0, 0.0},
		{1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0, 0.0, 0.0, 0.0},
		{1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0, 0.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	}
)

// Alloc2 allocates a 3x3 2nd order tensor
func Alloc2() (tensor [][]float64) {
	tensor = make([][]float64, 3)
	for k := 0; k < 3; k++ {
		tensor[k] = make([]float64, 3)
	}
	return
}

// Alloc4 allocates a 3x3x3x3 4th order tensor
func Alloc4() (tensor [][][][]float64) {
	tensor = make([][][][]float64, 3)
	for i := 0; i < 3; i++ {
		tensor[i] = make([][][]float64, 3)
		for j := 0; j < 3; j++ {
			tensor[i][j] = make([][]float64, 3)
			for k := 0; k < 3; k++ {
				tensor[i][j][k] = make([]float64, 3)
			}
		}
	}
	return
}

// M_Alloc2 allocates a 2th order tensor in Mandel's representation (2*ndim)
func M_Alloc2(ndim int) (a []float64) {
	nσ := 2 * ndim
	a = make([]float64, nσ)
	return
}

// M_Alloc4 allocates a 4th order tensor in Mandel's representation ((2*ndim)x(2*ndim))
func M_Alloc4(ndim int) (A [][]float64) {
	nσ := 2 * ndim
	A = make([][]float64, nσ)
	for i := 0; i < nσ; i++ {
		A[i] = make([]float64, nσ)
	}
	return
}

/*
// IndMan2Ten returns the 3x3 tensor indices (i,j) corresponding to Mandel's index I
func IndMan2Ten(I int) (i, j int) {
    i = M2Ti[I]
    j = M2Tj[I]
    return
}
*/

// M2T converts Mandel components to 3x3 second order tensor components,
// i.e. correcting off-diagonal values that were multiplied by SQ2
func M2T(mandel []float64, i, j int) (component float64) {
	if i == j {
		component = mandel[T2MI[i][j]]
		return
	}
	switch len(mandel) {
	case 4:
		if i == 2 || j == 2 { // zero off-diagonal components
			return
		}
		component = mandel[T2MI[i][j]] / SQ2
	case 6:
		component = mandel[T2MI[i][j]] / SQ2
	default:
		chk.Panic(_tensor_oor, "tensor.go: M2T", len(mandel))
	}
	return
}

// M2TT converts Mandel components to 3x3x3x3 fourth order tensor components,
// i.e. correcting all values that were multiplied by 2 or SQ2
func M2TT(mandel [][]float64, i, j, k, l int) float64 {
	a, b := TT2MI[i][j][k][l], TT2MJ[i][j][k][l]
	if (i == j && k != l) || (k == l && i != j) {
		return mandel[a][b] / SQ2
	}
	if i != j && k != l {
		return mandel[a][b] / 2.0
	}
	return mandel[a][b]
}

// Man2Ten returns the 3x3 2nd order tensor from its Mandel representation
func Man2Ten(tensor [][]float64, mandel []float64) {
	switch len(mandel) {
	case 4:
		tensor[0][0], tensor[0][1], tensor[0][2] = mandel[0], mandel[3]/SQ2, 0.0
		tensor[1][0], tensor[1][1], tensor[1][2] = mandel[3]/SQ2, mandel[1], 0.0
		tensor[2][0], tensor[2][1], tensor[2][2] = 0.0, 0.0, mandel[2]
	case 6:
		tensor[0][0], tensor[0][1], tensor[0][2] = mandel[0], mandel[3]/SQ2, mandel[5]/SQ2
		tensor[1][0], tensor[1][1], tensor[1][2] = mandel[3]/SQ2, mandel[1], mandel[4]/SQ2
		tensor[2][0], tensor[2][1], tensor[2][2] = mandel[5]/SQ2, mandel[4]/SQ2, mandel[2]
	default:
		chk.Panic(_tensor_oor, "tensor.go: Man2Ten", len(mandel))
	}
}

// Ten2Man returns the Mandel representation of a 3x3 2nd order tensor
func Ten2Man(mandel []float64, tensor [][]float64) {
	// check symmetry
	if math.Abs(tensor[0][1]-tensor[1][0]) > EPS {
		chk.Panic(_tensor_m1, len(mandel))
	}
	if math.Abs(tensor[1][2]-tensor[2][1]) > EPS {
		chk.Panic(_tensor_m1, len(mandel))
	}
	if math.Abs(tensor[2][0]-tensor[0][2]) > EPS {
		chk.Panic(_tensor_m1, len(mandel))
	}
	// convert
	switch len(mandel) {
	case 4:
		if math.Abs(tensor[0][2]) > EPS {
			chk.Panic(_tensor_m2, tensor[0][2], tensor[1][2])
		}
		if math.Abs(tensor[1][2]) > EPS {
			chk.Panic(_tensor_m2, tensor[0][2], tensor[1][2])
		}
		mandel[0] = tensor[0][0]
		mandel[1] = tensor[1][1]
		mandel[2] = tensor[2][2]
		mandel[3] = tensor[0][1] * SQ2
	case 6:
		mandel[0] = tensor[0][0]
		mandel[1] = tensor[1][1]
		mandel[2] = tensor[2][2]
		mandel[3] = tensor[0][1] * SQ2
		mandel[4] = tensor[1][2] * SQ2
		mandel[5] = tensor[2][0] * SQ2
	default:
		chk.Panic(_tensor_oor, "tensor.go: Ten2Man", len(mandel))
	}
}

// auxiliary functions
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// error messages
const (
	_tensor_m1  = "tensor.go: Ten2Man: tensor must be symmetric in order to be converted to Mandel's vector with %d components"
	_tensor_m2  = "tensor.go: Ten2Man: with 4 components, tensor must have zero [0,2]==%g and [1,2]==%g components"
	_tensor_oor = "%s: The length of Mandel's vector representing a 2nd order tensor must either 4 (2D) or 6 (3D). len=%d is incorrect"
)
