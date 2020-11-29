// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/mpi"
)

func main() {

	mpi.Start()
	defer mpi.Stop()

	if mpi.WorldRank() == 0 {
		io.Pf("\n\n------------------ Test MPI 04 ------------------\n\n")
	}

	comm := mpi.NewCommunicator(nil)

	for i := 0; i < 60; i++ {
		time.Sleep(1e9)
		io.Pf("hello from %v\n", comm.Rank())
		if comm.Rank() == 2 && i == 3 {
			io.PfGreen("rank = 3 wants to abort (the following error is OK)\n")
			comm.Abort()
		}
	}
}
