// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package la

import (
	"gosl/mpi"
)

// SpTriReduce joins (MPI) parallel triplets to root (Rank == 0) processor.
//  NOTE: J in root is also joined into Jroot
func SpTriReduce(comm *mpi.Communicator, J *Triplet) {
	if comm.Rank() == 0 {
		for proc := 1; proc < comm.Size(); proc++ {
			nnz := comm.RecvOneI(proc)
			irec := make([]int, nnz)
			drec := make([]float64, nnz)
			comm.RecvI(irec, proc)
			J.i = append(J.i, irec...)
			comm.RecvI(irec, proc)
			J.j = append(J.j, irec...)
			comm.Recv(drec, proc)
			J.x = append(J.x, drec...)
		}
		J.pos = len(J.x)
		J.max = J.pos
	} else {
		comm.SendOneI(J.max, 0)
		comm.SendI(J.i, 0)
		comm.SendI(J.j, 0)
		comm.Send(J.x, 0)
	}
}
