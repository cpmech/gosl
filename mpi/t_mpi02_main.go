// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/mpi"
)

func main() {

	mpi.Start(false)
	defer func() {
		mpi.Stop(false)
	}()

	if mpi.Rank() == 0 {
		io.PfYel("\nTest MPI 02\n")
	}
	if mpi.Size() != 3 {
		chk.Panic("this test needs 3 processors")
	}

	var myints []int
	var mydbls []float64
	switch mpi.Rank() {
	case 1:
		myints = []int{1, 2, 3, 4}
		mydbls = []float64{-1, -2, -3}
	case 2:
		myints = []int{20, 30, 40, 50, 60}
		mydbls = []float64{-20, -50}
	}

	// SingleIntSend
	if mpi.Rank() == 0 {
		v1 := make([]int, mpi.Size())
		v2 := make([]int, mpi.Size())
		allints := []int{}
		alldbls := []float64{}
		for proc := 1; proc < mpi.Size(); proc++ {
			// SingleIntRecv
			val := mpi.SingleIntRecv(proc)
			io.Pf("root recieved val=%d from proc=%d\n", val, proc)
			v1[proc] = val
			val = mpi.SingleIntRecv(proc)
			io.Pf("root recieved val=%d from proc=%d\n", val, proc)
			v2[proc] = val
			// IntRecv
			n := mpi.SingleIntRecv(proc)
			io.Pf("root recieved n=%d from proc=%d\n", n, proc)
			ints := make([]int, n)
			mpi.IntRecv(ints, proc)
			io.Pf("root recieved ints=%v from proc=%d\n", ints, proc)
			allints = append(allints, ints...)
			// DblRecv
			n = mpi.SingleIntRecv(proc)
			io.Pf("root recieved n=%d from proc=%d\n", n, proc)
			dbls := make([]float64, n)
			mpi.DblRecv(dbls, proc)
			io.Pf("root recieved dbls=%v from proc=%d\n", dbls, proc)
			alldbls = append(alldbls, dbls...)
		}
		var tst testing.T
		chk.Ints(&tst, "SingleIntRecv: vals", v1, []int{0, 1001, 1002})
		chk.Ints(&tst, "SingleIntRecv: vals", v2, []int{0, 2001, 2002})
		chk.Ints(&tst, "IntRecv: allints", allints, []int{1, 2, 3, 4, 20, 30, 40, 50, 60})
		chk.Vector(&tst, "IntRecv: alldbls", 1e-17, alldbls, []float64{-1, -2, -3, -20, -50})
	} else {
		// SingleIntSend
		mpi.SingleIntSend(1000+mpi.Rank(), 0)
		mpi.SingleIntSend(2000+mpi.Rank(), 0)
		// IntSend
		mpi.SingleIntSend(len(myints), 0)
		mpi.IntSend(myints, 0)
		// DblSend
		mpi.SingleIntSend(len(mydbls), 0)
		mpi.DblSend(mydbls, 0)
	}
}
