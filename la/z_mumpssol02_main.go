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

	comm := mpi.NewCommunicator(nil)

	myrank := comm.Rank()
	if myrank == 0 {
		io.Pf("\n------------------- Test MUMPS Sol 02 -------------------\n")
	}

	ndim := 10
	id, sz := comm.Rank(), comm.Size()
	start, endp1 := (id*ndim)/sz, ((id+1)*ndim)/sz

	if comm.Size() > ndim {
		chk.Panic("the number of processors must be smaller than or equal to %d", ndim)
	}

	b := make([]float64, ndim)
	var t la.Triplet
	t.Init(ndim, ndim, ndim*ndim)

	for i := start; i < endp1; i++ {
		j := i
		if i > 0 {
			j = i - 1
		}
		for ; j < ndim; j++ {
			val := 10.0 - float64(j)
			if i > j {
				val -= 1.0
			}
			t.Put(i, j, val)
		}
		b[i] = float64(i + 1)
	}

	chk.Verbose = true
	tst := new(testing.T)

	bIsDistr := true
	xCorrect := []float64{-1, 8, -65, 454, -2725, 13624, -54497, 163490, -326981, 326991}
	la.TestSpSolver(tst, "mumps", false, &t, b, xCorrect, 1e-4, 1e-14, false, bIsDistr, comm)
}
