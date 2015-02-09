// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"testing"

	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/utl"
)

func setslice(x []float64) {
	switch mpi.Rank() {
	case 0:
		copy(x, []float64{0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3})
	case 1:
		copy(x, []float64{10, 10, 10, 20, 20, 20, 30, 30, 30, 40, 40})
	case 2:
		copy(x, []float64{100, 100, 100, 1000, 1000, 1000, 2000, 2000, 2000, 3000, 3000})
	}
}

func main() {

	mpi.Start(false)
	defer func() {
		if err := recover(); err != nil {
			utl.PfRed("Some error has happened: %v\n", err)
		}
		mpi.Stop(false)
	}()

	utl.Tsilent = false
	if mpi.Rank() == 0 {
		utl.TTitle("Test MPI 01")
	}
	if mpi.Size() != 3 {
		utl.Panic("this test needs 3 processors")
	}
	n := 11
	x := make([]float64, n)
	id, sz := mpi.Rank(), mpi.Size()
	start, endp1 := (id*n)/sz, ((id+1)*n)/sz
	for i := start; i < endp1; i++ {
		x[i] = float64(i)
	}

	// Barrier
	mpi.Barrier()

	utl.Pfgrey("x @ proc # %d = %v\n", id, x)

	// SumToRoot
	r := make([]float64, n)
	mpi.SumToRoot(r, x)
	var tst testing.T
	if id == 0 {
		utl.CompareDbls(&tst, fmt.Sprintf("SumToRoot:       r @ proc # %d", id), r, []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	} else {
		utl.CompareDbls(&tst, fmt.Sprintf("SumToRoot:       r @ proc # %d", id), r, make([]float64, n))
	}

	// BcastFromRoot
	r[0] = 666
	mpi.BcastFromRoot(r)
	utl.CompareDbls(&tst, fmt.Sprintf("BcastFromRoot:   r @ proc # %d", id), r, []float64{666, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// AllReduceSum
	setslice(x)
	w := make([]float64, n)
	mpi.AllReduceSum(x, w)
	utl.CompareDbls(&tst, fmt.Sprintf("AllReduceSum:    w @ proc # %d", id), w, []float64{110, 110, 110, 1021, 1021, 1021, 2032, 2032, 2032, 3043, 3043})

	// AllReduceSumAdd
	setslice(x)
	y := []float64{-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000}
	mpi.AllReduceSumAdd(y, x, w)
	utl.CompareDbls(&tst, fmt.Sprintf("AllReduceSumAdd: y @ proc # %d", id), y, []float64{-890, -890, -890, 21, 21, 21, 1032, 1032, 1032, 2043, 2043})

	// AllReduceMin
	setslice(x)
	mpi.AllReduceMin(x, w)
	utl.CompareDbls(&tst, fmt.Sprintf("AllReduceMin:    x @ proc # %d", id), x, []float64{0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3})

	// AllReduceMax
	setslice(x)
	mpi.AllReduceMax(x, w)
	utl.CompareDbls(&tst, fmt.Sprintf("AllReduceMax:    x @ proc # %d", id), x, []float64{100, 100, 100, 1000, 1000, 1000, 2000, 2000, 2000, 3000, 3000})
}
