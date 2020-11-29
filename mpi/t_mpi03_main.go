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

func main() {

	mpi.Start()
	defer mpi.Stop()

	if mpi.WorldRank() == 0 {
		io.Pf("\n\n------------------ Test MPI 03 ------------------\n\n")
	}
	if mpi.WorldSize() != 3 {
		chk.Panic("this test needs 3 processors")
	}

	x := []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
	n := len(x)
	id, sz := mpi.WorldRank(), mpi.WorldSize()
	start, endp1 := (id*n)/sz, ((id+1)*n)/sz
	for i := start; i < endp1; i++ {
		x[i] = i
	}

	comm := mpi.NewCommunicator(nil)

	w := make([]int, n)
	comm.AllReduceMaxI(w, x)

	chk.Verbose = true
	var tst testing.T
	chk.Ints(&tst, fmt.Sprintf("AllReduceMaxI: w @ proc # %d", id), w, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
}
