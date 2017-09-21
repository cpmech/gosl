// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/mpi"
)

func setSlice(x []float64, rank, ncpus int) {
	for i := 0; i < len(x); i++ {
		x[i] = -1
	}
	start, endp1 := (rank*len(x))/ncpus, ((rank+1)*len(x))/ncpus
	for i := start; i < endp1; i++ {
		x[i] = float64(1 + rank)
	}
}

func setSliceC(x []complex128, rank, ncpus int) {
	for i := 0; i < len(x); i++ {
		x[i] = 0
	}
	start, endp1 := (rank*len(x))/ncpus, ((rank+1)*len(x))/ncpus
	for i := start; i < endp1; i++ {
		x[i] = complex(float64(1+rank), float64(1+rank)/10.0)
	}
}

func setSliceI(x []int, rank, ncpus int) {
	for i := 0; i < len(x); i++ {
		x[i] = -1
	}
	start, endp1 := (rank*len(x))/ncpus, ((rank+1)*len(x))/ncpus
	for i := start; i < endp1; i++ {
		x[i] = 1 + rank
	}
}

func main() {

	mpi.Start()
	defer mpi.Stop()

	// message and check
	if mpi.WorldRank() == 0 {
		io.Pf("\n\n------------------ Test MPI 00 ------------------\n\n")
	}
	if mpi.WorldSize() != 8 {
		io.Pf("this test needs 8 processors\n")
		return
	}

	// subsets of processors
	A := mpi.NewCommunicator([]int{0, 1, 2, 3})
	B := mpi.NewCommunicator([]int{4, 5, 6, 7})

	// test structure
	chk.Verbose = true
	tst := new(testing.T)

	// run tests
	rank := mpi.WorldRank()
	if rank < 4 {

		// BcastFromRoot
		x := make([]float64, 8)
		if A.Rank() == 0 {
			for i := 0; i < len(x); i++ {
				x[i] = float64(1 + i)
			}
		}
		A.BcastFromRoot(x)
		chk.Array(tst, "A: x (real)", 1e-17, x, []float64{1, 2, 3, 4, 5, 6, 7, 8})

		// ReduceSum
		setSlice(x, int(A.Rank()), int(A.Size()))
		res := make([]float64, len(x))
		A.ReduceSum(res, x)
		if A.Rank() == 0 {
			chk.Array(tst, "A root: res", 1e-17, res, []float64{1 - 3, 1 - 3, 2 - 3, 2 - 3, 3 - 3, 3 - 3, 4 - 3, 4 - 3})
		} else {
			chk.Array(tst, "A others: res", 1e-17, res, nil)
		}

		// AllReduceSum
		setSlice(x, int(A.Rank()), int(A.Size()))
		for i := 0; i < len(x); i++ {
			res[i] = 0
		}
		A.AllReduceSum(res, x)
		chk.Array(tst, "A all (sum): res", 1e-17, res, []float64{1 - 3, 1 - 3, 2 - 3, 2 - 3, 3 - 3, 3 - 3, 4 - 3, 4 - 3})

		// AllReduceMin
		setSlice(x, int(A.Rank()), int(A.Size()))
		for i := 0; i < len(x); i++ {
			res[i] = 0
		}
		A.AllReduceMin(res, x)
		chk.Array(tst, "A all (min): res", 1e-17, res, []float64{-1, -1, -1, -1, -1, -1, -1, -1})

		// AllReduceMax
		setSlice(x, int(A.Rank()), int(A.Size()))
		for i := 0; i < len(x); i++ {
			res[i] = 0
		}
		A.AllReduceMax(res, x)
		chk.Array(tst, "A all (max): res", 1e-17, res, []float64{1, 1, 2, 2, 3, 3, 4, 4})

		// Send & Recv
		if A.Rank() == 0 {
			s := []float64{123, 123, 123, 123}
			for k := 1; k <= 3; k++ {
				A.Send(s, k)
			}
		} else {
			y := make([]float64, 4)
			A.Recv(y, 0)
			chk.Array(tst, "A recv", 1e-17, y, []float64{123, 123, 123, 123})
		}

		// SendI & RecvI
		if A.Rank() == 0 {
			s := []int{123, 123, 123, 123}
			for k := 1; k <= 3; k++ {
				A.SendI(s, k)
			}
		} else {
			y := make([]int, 4)
			A.RecvI(y, 0)
			chk.Ints(tst, "A recvI", y, []int{123, 123, 123, 123})
		}

		// SendOneI & RecvOneI
		if A.Rank() == 0 {
			for k := 1; k <= 3; k++ {
				A.SendOneI(456, k)
			}
		} else {
			res := A.RecvOneI(0)
			chk.Int(tst, "A RecvOneI", res, 456)
		}

	} else {

		// BcastFromRootC
		x := make([]complex128, 8)
		if B.Rank() == 0 {
			for i := 0; i < len(x); i++ {
				x[i] = complex(float64(1+i), float64(1+i)/10.0)
			}
		}
		B.BcastFromRootC(x)
		chk.ArrayC(tst, "B: x (complex)", 1e-17, x, []complex128{1 + 0.1i, 2 + 0.2i, 3 + 0.3i, 4 + 0.4i, 5 + 0.5i, 6 + 0.6i, 7 + 0.7i, 8 + 0.8i})

		// ReduceSum
		setSliceC(x, int(B.Rank()), int(B.Size()))
		res := make([]complex128, len(x))
		B.ReduceSumC(res, x)
		if B.Rank() == 0 {
			chk.ArrayC(tst, "B root: res", 1e-17, res, []complex128{1 + 0.1i, 1 + 0.1i, 2 + 0.2i, 2 + 0.2i, 3 + 0.3i, 3 + 0.3i, 4 + 0.4i, 4 + 0.4i})
		} else {
			chk.ArrayC(tst, "B others: res", 1e-17, res, nil)
		}

		// AllReduceSumC
		setSliceC(x, int(B.Rank()), int(B.Size()))
		for i := 0; i < len(x); i++ {
			res[i] = 0
		}
		B.AllReduceSumC(res, x)
		chk.ArrayC(tst, "B all: res", 1e-17, res, []complex128{1 + 0.1i, 1 + 0.1i, 2 + 0.2i, 2 + 0.2i, 3 + 0.3i, 3 + 0.3i, 4 + 0.4i, 4 + 0.4i})

		// AllReduceMinI
		z := make([]int, 8)
		zres := make([]int, 8)
		setSliceI(z, int(B.Rank()), int(B.Size()))
		B.AllReduceMinI(zres, z)
		chk.Ints(tst, "A all (min int): res", zres, []int{-1, -1, -1, -1, -1, -1, -1, -1})

		// AllReduceMaxI
		setSliceI(z, int(B.Rank()), int(B.Size()))
		for i := 0; i < len(z); i++ {
			zres[i] = 0
		}
		B.AllReduceMaxI(zres, z)
		chk.Ints(tst, "A all (max int): res", zres, []int{1, 1, 2, 2, 3, 3, 4, 4})

		// SendC & RecvC
		if B.Rank() == 0 {
			s := []complex128{123 + 1i, 123 + 2i, 123 + 3i, 123 + 4i}
			for k := 1; k <= 3; k++ {
				B.SendC(s, k)
			}
		} else {
			y := make([]complex128, 4)
			B.RecvC(y, 0)
			chk.ArrayC(tst, "B recv", 1e-17, y, []complex128{123 + 1i, 123 + 2i, 123 + 3i, 123 + 4i})
		}

		// SendOne & RecvOne
		if B.Rank() == 0 {
			for k := 1; k <= 3; k++ {
				B.SendOne(-123, k)
			}
		} else {
			res := B.RecvOne(0)
			chk.Float64(tst, "B RecvOne", 1e-17, res, -123)
		}
	}

	// wait for all
	world := mpi.NewCommunicator(nil)
	world.Barrier()
}
