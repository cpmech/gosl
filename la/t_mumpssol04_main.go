// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
)

func main() {

	mpi.Start(false)
	defer func() {
		mpi.Stop(false)
	}()

	verbose()
	myrank := mpi.Rank()
	if myrank == 0 {
		chk.PrintTitle("Test MUMPS Sol 04")
	}

	ndim := 10
	id, sz := mpi.Rank(), mpi.Size()
	start, endp1 := (id*ndim)/sz, ((id+1)*ndim)/sz

	if mpi.Size() > ndim {
		chk.Panic("the number of processors must be smaller than or equal to %d", ndim)
	}

	b := make([]complex128, ndim)
	var t la.TripletC
	t.Init(ndim, ndim, ndim*ndim, true)

	for i := start; i < endp1; i++ {
		j := i
		if i > 0 {
			j = i - 1
		}
		for ; j < 10; j++ {
			val := 10.0 - float64(j)
			if i > j {
				val -= 1.0
			}
			t.Put(i, j, val, 0)
		}
		b[i] = complex(float64(i+1), 0.0)
	}

	x_correct := []complex128{-1, 8, -65, 454, -2725, 13624, -54497, 163490, -326981, 326991}
	sum_b_to_root := true
	la.RunMumpsTestC(&t, 1e-4, b, x_correct, sum_b_to_root)
}
