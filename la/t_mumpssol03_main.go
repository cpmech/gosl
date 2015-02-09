// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"code.google.com/p/gosl/la"
	"code.google.com/p/gosl/mpi"
	"code.google.com/p/gosl/utl"
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
	myrank := mpi.Rank()
	if myrank == 0 {
		utl.TTitle("Test MUMPS Sol 03")
	}

	var t la.TripletC
	switch mpi.Size() {
	case 1:
		t.Init(5, 5, 13, true)
		t.Put(0, 0, 1.0, 0)
		t.Put(0, 0, 1.0, 0)
		t.Put(1, 0, 3.0, 0)
		t.Put(0, 1, 3.0, 0)
		t.Put(2, 1, -1.0, 0)
		t.Put(4, 1, 4.0, 0)
		t.Put(1, 2, 4.0, 0)
		t.Put(2, 2, -3.0, 0)
		t.Put(3, 2, 1.0, 0)
		t.Put(4, 2, 2.0, 0)
		t.Put(2, 3, 2.0, 0)
		t.Put(1, 4, 6.0, 0)
		t.Put(4, 4, 1.0, 0)
	case 2:
		if myrank == 0 {
			t.Init(5, 5, 6, true)
			t.Put(0, 0, 1.0, 0)
			t.Put(0, 0, 1.0, 0)
			t.Put(1, 0, 3.0, 0)
			t.Put(0, 1, 3.0, 0)
			t.Put(2, 1, -1.0, 0)
			t.Put(4, 1, 4.0, 0)
		} else {
			t.Init(5, 5, 7, true)
			t.Put(1, 2, 4.0, 0)
			t.Put(2, 2, -3.0, 0)
			t.Put(3, 2, 1.0, 0)
			t.Put(4, 2, 2.0, 0)
			t.Put(2, 3, 2.0, 0)
			t.Put(1, 4, 6.0, 0)
			t.Put(4, 4, 1.0, 0)
		}
	default:
		utl.Panic("this test needs 1 or 2 procs")
	}

	b := []complex128{8.0, 45.0, -3.0, 3.0, 19.0}
	x_correct := []complex128{1, 2, 3, 4, 5}
	sum_b_to_root := false
	la.RunMumpsTestC(&t, 1e-14, b, x_correct, sum_b_to_root)
}
