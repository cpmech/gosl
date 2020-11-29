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
		io.Pf("\n------------------- Test MUMPS Sol 05 --- (complex) -----\n")
	}

	ndim := 10
	id, sz := comm.Rank(), comm.Size()
	start, endp1 := (id*ndim)/sz, ((id+1)*ndim)/sz

	if comm.Size() > ndim {
		chk.Panic("the number of processors must be smaller than or equal to %d", ndim)
	}

	b := make([]complex128, ndim)
	xCorrect := make([]complex128, ndim)

	// Let exact solution = 1 + 0.5i
	for i := 0; i < ndim; i++ {
		xCorrect[i] = complex(float64(i+1), float64(i+1)/10.0)
	}

	var t la.TripletC
	t.Init(ndim, ndim, ndim)

	// assemble a and b
	for i := start; i < endp1; i++ {

		// Some very fake diagonals. Should take exactly 20 GMRES steps
		ar := 10.0 + float64(i)/(float64(ndim)/10.0)
		ac := 10.0 - float64(i)/(float64(ndim)/10.0)
		t.Put(i, i, complex(ar, ac))

		// Generate RHS to match exact solution
		b[i] = complex(ar*real(xCorrect[i])-ac*imag(xCorrect[i]),
			ar*imag(xCorrect[i])+ac*real(xCorrect[i]))
	}

	chk.Verbose = true
	tst := new(testing.T)

	bIsDistr := true
	la.TestSpSolverC(tst, "mumps", false, &t, b, xCorrect, 1e-11, 1e-17, false, bIsDistr, comm)
}
