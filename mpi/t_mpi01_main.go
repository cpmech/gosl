// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/mpi"
)

func setslice(x []float64) {
	switch mpi.WorldRank() {
	case 0:
		copy(x, []float64{0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3})
	case 1:
		copy(x, []float64{10, 10, 10, 20, 20, 20, 30, 30, 30, 40, 40})
	case 2:
		copy(x, []float64{100, 100, 100, 1000, 1000, 1000, 2000, 2000, 2000, 3000, 3000})
	}
}

func main() {

	mpi.Start()
	defer mpi.Stop()

	if mpi.WorldRank() == 0 {
		io.Pf("\n\n------------------ Test MPI 01 ------------------\n\n")
	}
	if mpi.WorldSize() != 3 {
		chk.Panic("this test needs 3 processors")
	}

	n := 11
	x := make([]float64, n)
	id, sz := mpi.WorldRank(), mpi.WorldSize()
	start, endp1 := (id*n)/sz, ((id+1)*n)/sz
	for i := start; i < endp1; i++ {
		x[i] = float64(i)
	}

	// Communicator
	comm := mpi.NewCommunicator(nil) // World

	// Barrier
	comm.Barrier()
	io.Pfgrey("x @ proc # %d = %v\n", id, x)

	// testing variable
	chk.Verbose = true
	var tst testing.T

	// SumToRoot
	r := make([]float64, n)
	comm.ReduceSum(r, x)
	if id == 0 {
		chk.Array(&tst, fmt.Sprintf("ReduceSum:       r @ proc # %d", id), 1e-17, r, []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	} else {
		chk.Array(&tst, fmt.Sprintf("ReduceSum:       r @ proc # %d", id), 1e-17, r, make([]float64, n))
	}

	// BcastFromRoot
	r[0] = 123
	comm.BcastFromRoot(r)
	chk.Array(&tst, fmt.Sprintf("BcastFromRoot:   r @ proc # %d", id), 1e-17, r, []float64{123, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// AllReduceSum
	setslice(x)
	w := make([]float64, n)
	comm.AllReduceSum(w, x)
	chk.Array(&tst, fmt.Sprintf("AllReduceSum:    w @ proc # %d", id), 1e-17, w, []float64{110, 110, 110, 1021, 1021, 1021, 2032, 2032, 2032, 3043, 3043})

	// AllReduceSum
	setslice(x)
	y := []float64{-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000}
	for i := 0; i < len(w); i++ {
		w[i] = 0
	}
	comm.AllReduceSum(w, x)
	for i := 0; i < len(w); i++ {
		y[i] += w[i]
	}
	chk.Array(&tst, fmt.Sprintf("AllReduceSum:    y @ proc # %d", id), 1e-17, y, []float64{-890, -890, -890, 21, 21, 21, 1032, 1032, 1032, 2043, 2043})

	// AllReduceMin
	setslice(x)
	comm.AllReduceMin(w, x)
	chk.Array(&tst, fmt.Sprintf("AllReduceMin:    w @ proc # %d", id), 1e-17, w, []float64{0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3})

	// AllReduceMax
	setslice(x)
	comm.AllReduceMax(w, x)
	chk.Array(&tst, fmt.Sprintf("AllReduceMax:    w @ proc # %d", id), 1e-17, w, []float64{100, 100, 100, 1000, 1000, 1000, 2000, 2000, 2000, 3000, 3000})
}
