// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/utl"
)

func main() {

	mpi.Start(false)
	defer func() {
		if err := recover(); err != nil {
			io.PfRed("Some error has happened: %v\n", err)
		}
		mpi.Stop(false)
	}()

	verbose() = false

	if mpi.Rank() == 0 {
		chk.PrintTitle("Test MPI 04")
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
