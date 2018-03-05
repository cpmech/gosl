// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import "math"

var (
	// some constants
	sq2    = math.Sqrt(2.0)       // sqrt(2)
	sq3    = math.Sqrt(3.0)       // sqrt(3)
	sq6    = math.Sqrt(6.0)       // sqrt(6)
	sq3by2 = math.Sqrt(3.0 / 2.0) // sqrt(3/2)
	sq2by3 = math.Sqrt(2.0 / 3.0) // sqrt(2/3)
	twosq2 = 2.0 * math.Sqrt(2.0) // 2*sqrt(2) == 2^(3/2)
	otrd   = 1.0 / 3.0            // one-third
	ttrd   = 2.0 / 3.0            // two-thirds

	// SecToManI converts i-j-indices of 3x3 2nd order (symmetric) tensor to I-index in Mandel's representation
	SecToManI = [][]int{
		{0, 3, 5},
		{3, 1, 4},
		{5, 4, 2},
	}

	// SecToVecI converts i-j-indices of 3x3 2nd order (symmetric) tensor to I-index in vector representation
	SecToVecI = [][]int{
		{0, 3, 5},
		{6, 1, 4},
		{8, 7, 2},
	}

	// FouToManI converts i-j-k-l-indices of 3x3x3x3 4th order (full-symmetric) tensor to I-index in Mandel's representation
	FouToManI = [][][][]int{
		{ //    0.0.       0.1.       0.2.
			{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}, // 00..
			{{3, 3, 3}, {3, 3, 3}, {3, 3, 3}}, // 01..
			{{5, 5, 5}, {5, 5, 5}, {5, 5, 5}}, // 02..
		},
		{ //    1.0.       1.1.       1.2.
			{{3, 3, 3}, {3, 3, 3}, {3, 3, 3}}, // 10..
			{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}, // 11..
			{{4, 4, 4}, {4, 4, 4}, {4, 4, 4}}, // 12..
		},
		{ //    2.0.       2.1.       2.2.
			{{5, 5, 5}, {5, 5, 5}, {5, 5, 5}}, // 20..
			{{4, 4, 4}, {4, 4, 4}, {4, 4, 4}}, // 21..
			{{2, 2, 2}, {2, 2, 2}, {2, 2, 2}}, // 22..
		},
	}

	// FouToManJ converts i-j-k-l-indices of 3x3x3x3 4th order (full-symmetric) tensor to J-index in Mandel's representation
	FouToManJ = [][][][]int{
		{
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
		},
		{
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
		},
		{
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
			{{0, 3, 5}, {3, 1, 4}, {5, 4, 2}},
		},
	}

	// FouToVecI converts i-j-k-l-indices of 3x3x3x3 4th order tensor to I-index in Vector representation
	FouToVecI = [][][][]int{
		{
			{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			{{3, 3, 3}, {3, 3, 3}, {3, 3, 3}},
			{{5, 5, 5}, {5, 5, 5}, {5, 5, 5}},
		},
		{
			{{6, 6, 6}, {6, 6, 6}, {6, 6, 6}},
			{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			{{4, 4, 4}, {4, 4, 4}, {4, 4, 4}},
		},
		{
			{{8, 8, 8}, {8, 8, 8}, {8, 8, 8}},
			{{7, 7, 7}, {7, 7, 7}, {7, 7, 7}},
			{{2, 2, 2}, {2, 2, 2}, {2, 2, 2}},
		},
	}

	// FouToVecJ converts i-j-k-l-indices of 3x3x3x3 4th order tensor to J-index in Vector representation
	FouToVecJ = [][][][]int{
		{
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
		},
		{
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
		},
		{
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
			{{0, 3, 5}, {6, 1, 4}, {8, 7, 2}},
		},
	}

	// ManToSecI converts I-index in Mandel's representation to i-index of 3x3 2nd order tensor
	ManToSecI = []int{0, 1, 2, 0, 1, 0}

	// ManToSecJ converts I-index in Mandel's representation to j-index of 3x3 2nd order tensor
	ManToSecJ = []int{0, 1, 2, 1, 2, 2}

	// SecIdenMan is the 3x3 2nd order identity tensor in Mandel's representation
	SecIdenMan = []float64{1, 1, 1, 0, 0, 0}

	// FouIdenMan is the 4th order identity tensor (symmetric) in Mandel's representation
	FouIdenMan = [][]float64{
		{1, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 1},
	}

	// FouPsdMan is the 4th order symmetric-deviatoric projector (3D) in Mandel's representation
	FouPsdMan = [][]float64{
		{+ttrd, -otrd, -otrd, 0, 0, 0},
		{-otrd, +ttrd, -otrd, 0, 0, 0},
		{-otrd, -otrd, +ttrd, 0, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 1},
	}

	// FouPisoMan is the 4th order isotropic projector (3D) in Mandel's representation
	FouPisoMan = [][]float64{
		{otrd, otrd, otrd, 0, 0, 0},
		{otrd, otrd, otrd, 0, 0, 0},
		{otrd, otrd, otrd, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}

	// IdenMat9x9 is the identity matrix with rows up to 9
	IdenMat9x9 = [][]float64{
		{1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
)
