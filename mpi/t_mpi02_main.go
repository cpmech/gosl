// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"testing"

	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/utl"
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
	if mpi.Rank() == 0 {
		utl.TTitle("Test MPI 02")
	}
	if mpi.Size() != 3 {
		utl.Panic("this test needs 3 processors")
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
			utl.Pf("root recieved val=%d from proc=%d\n", val, proc)
			v1[proc] = val
			val = mpi.SingleIntRecv(proc)
			utl.Pf("root recieved val=%d from proc=%d\n", val, proc)
			v2[proc] = val
			// IntRecv
			n := mpi.SingleIntRecv(proc)
			utl.Pf("root recieved n=%d from proc=%d\n", n, proc)
			ints := make([]int, n)
			mpi.IntRecv(ints, proc)
			utl.Pf("root recieved ints=%v from proc=%d\n", ints, proc)
			allints = append(allints, ints...)
			// DblRecv
			n = mpi.SingleIntRecv(proc)
			utl.Pf("root recieved n=%d from proc=%d\n", n, proc)
			dbls := make([]float64, n)
			mpi.DblRecv(dbls, proc)
			utl.Pf("root recieved dbls=%v from proc=%d\n", dbls, proc)
			alldbls = append(alldbls, dbls...)
		}
		var tst testing.T
		utl.CompareInts(&tst, "SingleIntRecv: vals", v1, []int{0, 1001, 1002})
		utl.CompareInts(&tst, "SingleIntRecv: vals", v2, []int{0, 2001, 2002})
		utl.CompareInts(&tst, "IntRecv: allints", allints, []int{1, 2, 3, 4, 20, 30, 40, 50, 60})
		utl.CompareDbls(&tst, "IntRecv: alldbls", alldbls, []float64{-1, -2, -3, -20, -50})
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
