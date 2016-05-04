// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
)

/*  Jacobian
    ========
        Calculates (with N=n-1):
            df0dx0, df0dx1, df0dx2, ... df0dxN
            df1dx0, df1dx1, df1dx2, ... df1dxN
                 . . . . . . . . . . . . .
            dfNdx0, dfNdx1, dfNdx2, ... dfNdxN
    INPUT:
        ffcn : f(x) function
        x    : station where dfdx has to be calculated
        fx   : f @ x
        w    : workspace with size == n == len(x)
    RETURNS:
        J : dfdx @ x [must be pre-allocated]        */
func Jacobian(J *la.Triplet, ffcn Cb_f, x, fx, w []float64, distr bool) (err error) {
	ndim := len(x)
	start, endp1 := 0, ndim
	if distr {
		id, sz := mpi.Rank(), mpi.Size()
		start, endp1 = (id*ndim)/sz, ((id+1)*ndim)/sz
		if J.Max() == 0 {
			J.Init(ndim, ndim, (endp1-start)*ndim)
		}
	} else {
		if J.Max() == 0 {
			J.Init(ndim, ndim, ndim*ndim)
		}
	}
	J.Start()
	// NOTE: cannot split calculation by columns unless the f function is
	//       independently calculated by each MPI processor.
	//       Otherwise, the AllReduce in f calculation would
	//       join pieces of f from different processors calculated for
	//       different x values (δx[col] from different columns).
	/*
	   for col := start; col < endp1; col++ {
	       xsafe := x[col]
	       delta := math.Sqrt(EPS * max(CTE1, math.Abs(xsafe)))
	       x[col] = xsafe + delta
	       ffcn(w, x) // fnew
	       io.Pforan("x = %v, f = %v\n", x, w)
	       for row := 0; row < ndim; row++ {
	           J.Put(row, col, (w[row]-fx[row])/delta)
	       }
	       x[col] = xsafe
	   }
	*/
	var df float64
	for col := 0; col < ndim; col++ {
		xsafe := x[col]
		delta := math.Sqrt(EPS * max(CTE1, math.Abs(xsafe)))
		x[col] = xsafe + delta
		err = ffcn(w, x) // w := f(x+δx[col])
		if err != nil {
			return
		}
		for row := start; row < endp1; row++ {
			df = w[row] - fx[row]
			//if math.Abs(df) > EPS {
			J.Put(row, col, df/delta)
			//}
		}
		x[col] = xsafe
	}
	return
}
