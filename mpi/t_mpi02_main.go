// Copyright 2016 The Gosl Authors. All rights reserved.
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

	mpi.Start()
	defer mpi.Stop()

	if mpi.WorldRank() == 0 {
		io.Pf("\n\n------------------ Test MPI 02 ------------------\n\n")
	}
	if mpi.WorldSize() != 3 {
		chk.Panic("this test needs 3 processors")
	}

	var myints []int
	var mydbls []float64
	switch mpi.WorldRank() {
	case 1:
		myints = []int{1, 2, 3, 4}
		mydbls = []float64{-1, -2, -3}
	case 2:
		myints = []int{20, 30, 40, 50, 60}
		mydbls = []float64{-20, -50}
	}

	comm := mpi.NewCommunicator(nil)

	if mpi.WorldRank() == 0 {
		v1 := make([]int, mpi.WorldSize())
		v2 := make([]int, mpi.WorldSize())
		allints := []int{}
		alldbls := []float64{}
		for proc := 1; proc < int(mpi.WorldSize()); proc++ {

			// RecvOneI
			val := comm.RecvOneI(proc)
			io.Pf("root received val=%d from proc=%d\n", val, proc)
			v1[proc] = val
			val = comm.RecvOneI(proc)
			io.Pf("root received val=%d from proc=%d\n", val, proc)
			v2[proc] = val

			// RecvI
			n := comm.RecvOneI(proc)
			io.Pf("root received n=%d from proc=%d\n", n, proc)
			ints := make([]int, n)
			comm.RecvI(ints, proc)
			io.Pf("root received ints=%v from proc=%d\n", ints, proc)
			allints = append(allints, ints...)

			// Recv
			n = comm.RecvOneI(proc)
			io.Pf("root received n=%d from proc=%d\n", n, proc)
			dbls := make([]float64, n)
			comm.Recv(dbls, proc)
			io.Pf("root received dbls=%v from proc=%d\n", dbls, proc)
			alldbls = append(alldbls, dbls...)
		}

		// check
		chk.Verbose = true
		var tst testing.T
		chk.Ints(&tst, "SingleIntRecv: vals", v1, []int{0, 1001, 1002})
		chk.Ints(&tst, "SingleIntRecv: vals", v2, []int{0, 2001, 2002})
		chk.Ints(&tst, "IntRecv: allints", allints, []int{1, 2, 3, 4, 20, 30, 40, 50, 60})
		chk.Array(&tst, "IntRecv: alldbls", 1e-17, alldbls, []float64{-1, -2, -3, -20, -50})

	} else {

		// SendOneI
		comm.SendOneI(1000+mpi.WorldRank(), 0)
		comm.SendOneI(2000+mpi.WorldRank(), 0)

		// SendI
		comm.SendOneI(len(myints), 0)
		comm.SendI(myints, 0)

		// Send
		comm.SendOneI(len(mydbls), 0)
		comm.Send(mydbls, 0)
	}
}
