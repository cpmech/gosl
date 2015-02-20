// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
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

	mpi.Start(false)
	defer func() {
		mpi.Stop(false)
	}()

	if mpi.Rank() == 0 {
		io.PfYel("\nTest MPI 04\n")
	}

	for i := 0; i < 60; i++ {
		time.Sleep(1e9)
		io.Pf("hello from %v\n", mpi.Rank())
		if mpi.Rank() == 2 && i == 3 {
			io.PfGreen("rank = 3 wants to abort (the following error is OK)\n")
			mpi.Abort()
		}
	}
}
