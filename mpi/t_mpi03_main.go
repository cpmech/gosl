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
		utl.TTitle("Test MPI 03")
	}
	if mpi.Size() != 3 {
		utl.Panic("this test needs 3 processors")
	}
	x := []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
	n := len(x)
	id, sz := mpi.Rank(), mpi.Size()
	start, endp1 := (id*n)/sz, ((id+1)*n)/sz
	for i := start; i < endp1; i++ {
		x[i] = i
	}

	//utl.Pforan("x = %v\n", x)

	// IntAllReduceMax
	w := make([]int, n)
	mpi.IntAllReduceMax(x, w)
	var tst testing.T
	utl.CompareInts(&tst, fmt.Sprintf("IntAllReduceMax: x @ proc # %d", id), x, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	//utl.Pfred("x = %v\n", x)
}
