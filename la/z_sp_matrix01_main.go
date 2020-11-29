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

	chk.Verbose = true
	tst := new(testing.T)

	mpi.Start()
	defer mpi.Stop()

	comm := mpi.NewCommunicator(nil)

	myrank := comm.Rank()
	if myrank == 0 {
		io.Pf("\n------------------- Test Read Sparse Matrix -------------------\n")
	}

	if comm.Size() != 2 {
		chk.Panic("the number of processors must be 2")
	}

	correct0 := [][]float64{
		{2, 3, 0, 0, 0},
		{3, 0, 4, 0, 0},
		{0, -1, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 4, 0, 0, 0},
	}

	correct1 := [][]float64{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 6},
		{0, 0, -3, 2, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 2, 0, 1},
	}

	var T la.Triplet
	isSym := T.ReadSmat("data/small-sparse-matrix.mtx", true, comm)
	chk.Bool(tst, "isSym", isSym, false)

	M := T.ToDense()
	if comm.Rank() == 0 {
		io.Pf("%s\n", M.Print("%4g"))
		chk.Deep2(tst, "M @ proc0", 1e-17, M.GetDeep2(), correct0)
	} else {
		io.Pfcyan("%s\n", M.Print("%4g"))
		chk.Deep2(tst, "M @ proc1", 1e-17, M.GetDeep2(), correct1)
	}
}
