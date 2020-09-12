// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/la"
	"gosl/mpi"
)

func main() {

	mpi.Start()
	defer mpi.Stop()

	comm := mpi.NewCommunicator(nil)

	myrank := comm.Rank()
	if myrank == 0 {
		io.Pf("\n------------------- Test MUMPS Sol 01a -------------------\n")
	}

	var t la.Triplet
	switch comm.Size() {
	case 1:
		t.Init(5, 5, 13)
		t.Put(0, 0, +1.0)
		t.Put(0, 0, +1.0)
		t.Put(1, 0, +3.0)
		t.Put(0, 1, +3.0)
		t.Put(2, 1, -1.0)
		t.Put(4, 1, +4.0)
		t.Put(1, 2, +4.0)
		t.Put(2, 2, -3.0)
		t.Put(3, 2, +1.0)
		t.Put(4, 2, +2.0)
		t.Put(2, 3, +2.0)
		t.Put(1, 4, +6.0)
		t.Put(4, 4, +1.0)
	case 2:
		if myrank == 0 {
			t.Init(5, 5, 6)
			t.Put(0, 0, +1.0)
			t.Put(0, 0, +1.0)
			t.Put(1, 0, +3.0)
			t.Put(0, 1, +3.0)
			t.Put(2, 1, -1.0)
			t.Put(4, 1, +4.0)
		} else {
			t.Init(5, 5, 7)
			t.Put(1, 2, +4.0)
			t.Put(2, 2, -3.0)
			t.Put(3, 2, +1.0)
			t.Put(4, 2, +2.0)
			t.Put(2, 3, +2.0)
			t.Put(1, 4, +6.0)
			t.Put(4, 4, +1.0)
		}
	default:
		chk.Panic("this test needs 1 or 2 procs")
	}

	chk.Verbose = true
	tst := new(testing.T)

	b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}
	bIsDistr := false
	xCorrect := []float64{1, 2, 3, 4, 5}
	la.TestSpSolver(tst, "mumps", false, &t, b, xCorrect, 1e-14, 1e-14, false, bIsDistr, comm)
}
