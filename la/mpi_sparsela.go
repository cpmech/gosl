// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!appengine,!heroku

package la

import (
	"code.google.com/p/gosl/mpi"
)

// SpTriSumToRoot join (MPI) parallel triplets to root (Rank == 0) processor.
//  NOTE: J in root is also joined into Jroot
func SpTriSumToRoot(J *Triplet) {
	if mpi.Rank() == 0 {
		for proc := 1; proc < mpi.Size(); proc++ {
			nnz := mpi.SingleIntRecv(proc)
			irec := make([]int, nnz)
			drec := make([]float64, nnz)
			mpi.IntRecv(irec, proc)
			J.i = append(J.i, irec...)
			mpi.IntRecv(irec, proc)
			J.j = append(J.j, irec...)
			mpi.DblRecv(drec, proc)
			J.x = append(J.x, drec...)
		}
		J.pos = len(J.x)
		J.max = J.pos
	} else {
		mpi.SingleIntSend(J.max, 0)
		mpi.IntSend(J.i, 0)
		mpi.IntSend(J.j, 0)
		mpi.DblSend(J.x, 0)
	}
}
