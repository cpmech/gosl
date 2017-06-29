// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
)

func main() {

	mpi.Start()
	defer mpi.Stop()

	if mpi.WorldRank() == 0 {
		io.Pf("\n------------------------ Test SP MPI 01------------------------ \n")
	}

	comm := mpi.NewCommunicator(nil)

	M := [][]float64{
		{1000, 1000, 1000, 1011, 1021, 1000},
		{1000, 1000, 1000, 1012, 1022, 1000},
		{1000, 1000, 1000, 1013, 1023, 1000},
		{1011, 1012, 1013, 1000, 1000, 1000},
		{1021, 1022, 1023, 1000, 1000, 1000},
		{1000, 1000, 1000, 1000, 1000, 1000},
	}

	id, sz, m := comm.Rank(), comm.Size(), len(M)
	start, endp1 := (id*m)/sz, ((id+1)*m)/sz

	if sz > 6 {
		chk.Panic("this test works with at most 6 processors")
	}

	var J la.Triplet
	J.Init(m, m, m*m)
	for i := start; i < endp1; i++ {
		for j := 0; j < m; j++ {
			J.Put(i, j, M[i][j])
		}
	}

	la.SpTriReduce(comm, &J)

	chk.Verbose = true
	var tst testing.T

	if comm.Rank() == 0 {
		chk.Matrix(&tst, "J @ proc 0", 1.0e-17, J.GetDenseMatrix().GetSlice(), [][]float64{
			{1000, 1000, 1000, 1011, 1021, 1000},
			{1000, 1000, 1000, 1012, 1022, 1000},
			{1000, 1000, 1000, 1013, 1023, 1000},
			{1011, 1012, 1013, 1000, 1000, 1000},
			{1021, 1022, 1023, 1000, 1000, 1000},
			{1000, 1000, 1000, 1000, 1000, 1000},
		})
	}
}
