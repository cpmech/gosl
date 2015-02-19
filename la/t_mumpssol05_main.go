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

	myrank := mpi.Rank()
	if myrank == 0 {
		chk.PrintTitle("Test MUMPS Sol 05")
	}

	ndim := 10
	id, sz := mpi.Rank(), mpi.Size()
	start, endp1 := (id*ndim)/sz, ((id+1)*ndim)/sz

	if mpi.Size() > ndim {
		chk.Panic("the number of processors must be smaller than or equal to %d", ndim)
	}

	n := 10
	b := make([]complex128, n)
	x_correct := make([]complex128, n)

	// Let exact solution = 1 + 0.5i
	for i := 0; i < ndim; i++ {
		x_correct[i] = complex(float64(i+1), float64(i+1)/10.0)
	}

	var t la.TripletC
	t.Init(ndim, ndim, ndim, true)

	// assemble a and b
	for i := start; i < endp1; i++ {

		// Some very fake diagonals. Should take exactly 20 GMRES steps
		ar := 10.0 + float64(i)/(float64(ndim)/10.0)
		ac := 10.0 - float64(i)/(float64(ndim)/10.0)
		t.Put(i, i, ar, ac)

		// Generate RHS to match exact solution
		b[i] = complex(ar*real(x_correct[i])-ac*imag(x_correct[i]),
			ar*imag(x_correct[i])+ac*real(x_correct[i]))
	}

	sum_b_to_root := true
	la.RunMumpsTestC(&t, 1e-14, b, x_correct, sum_b_to_root)
}
