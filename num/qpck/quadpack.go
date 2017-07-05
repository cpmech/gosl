// Copyright 2017 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qpck

/*
#cgo LDFLAGS: -lopenblas -lgfortran -lm

typedef double (*fType) (double* x, int* fid);

void dqagse_(fType f, double* a, double* b, double* epsabs, double* epsrel, int* limit,
             double* result, double* abserr, int* neval, int* ier,
			 double* alist, double* blist, double* rlist, double* elist,
			 int* iord, int* last, int* fid);

void dqagie_(fType f, double* bound, int* infCode, double* epsabs, double* epsrel, int* limit,
			 double* result, double* abserr, int* neval, int* ier,
			 double* alist, double* blist, double* rlist, double* elist,
			 int* iord, int* last, int* fid);

#include "connect.h"
*/
import "C"
import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

// fType defines the callback function
type fType func(x float64) float64

// functions implements a functions database
var functions []fType = make([]fType, 64)

// Qagse computes a definite integral using an automatic integrator.
// 1D globally adaptive integrator using interval subdivision and extrapolation.
//
//   INPUT:
//     fid    -- id of function to avoid goroutine problems
//     f      -- function defining the integrand
//     a      -- lower limit of integration
//     b      -- upper limit of integration
//     epsabs -- absolute accuracy requested [use ≤0 for default]
//     epsrel -- relative accuracy requested [use ≤0 for default]
//
//   INPUT/OUTPUT:
//     NOTE: (1) the length of the 5 vectors below is equal to the "limit" variable in the original
//               code which is an upperbound on the number of subintervals in the partition of (a,b)
//           (2) the 5 vectors below may be <nil>, thus memory is allocated here
//
//     alist -- the first last  elements of which are the left
//              end points of the subintervals in the partition of the given integration range (a,b)
//
//     blist -- the first last elements of which are the right
//              end points of the subintervals in the partition of the given integration range (a,b)
//
//     rlist -- the first last elements of which are the integral
//              approximations on the subintervals
//
//     elist -- the first last  elements of which are the moduli
//              of the absolute error estimates on the subintervals
//
//     iord  -- the first k elements of which are pointers to
//              the error estimates over the subintervals, such that elist(iord(1)), ...,
//              elist(iord(k)) form a decreasing sequence, with k = last if last.le.(limit/2+2), and
//              k = limit+1-last otherwise
//   OUTPUT:
//     result -- approximation to the integral
//     abserr -- estimate of the modulus of the absolute error, which should equal or exceed abs(i-result)
//     neval  -- number of integrand evaluations
//     last   -- number of subintervals actually produced in the subdivision process
//
func Qagse(fid int, f fType, a, b, epsabs, epsrel float64, alist, blist, rlist, elist []float64, iord []int32) (result, abserr float64, neval, last int32, err error) {

	// set function in database
	if fid >= len(functions) {
		err = chk.Err("functions database capacity exceeded. max number of functions = %d\n", len(functions))
		return
	}
	functions[fid] = f

	// default values
	if epsabs <= 0 {
		epsabs = 1.49e-8
	}
	if epsrel <= 0 {
		epsrel = 1.49e-8
	}

	// allocate vectors
	limit := len(alist)
	if limit < 1 {
		limit = 50
		alist = make([]float64, limit)
		blist = make([]float64, limit)
		rlist = make([]float64, limit)
		elist = make([]float64, limit)
		iord = make([]int32, limit)
	}

	// call quadpack
	var ier int32
	C.dqagse_(
		C.fType(C.fcn),
		(*C.double)(unsafe.Pointer(&a)),
		(*C.double)(unsafe.Pointer(&b)),
		(*C.double)(unsafe.Pointer(&epsabs)),
		(*C.double)(unsafe.Pointer(&epsrel)),
		(*C.int)(unsafe.Pointer(&limit)),
		(*C.double)(unsafe.Pointer(&result)),
		(*C.double)(unsafe.Pointer(&abserr)),
		(*C.int)(unsafe.Pointer(&neval)),
		(*C.int)(unsafe.Pointer(&ier)),
		(*C.double)(unsafe.Pointer(&alist[0])),
		(*C.double)(unsafe.Pointer(&blist[0])),
		(*C.double)(unsafe.Pointer(&rlist[0])),
		(*C.double)(unsafe.Pointer(&elist[0])),
		(*C.int)(unsafe.Pointer(&iord[0])),
		(*C.int)(unsafe.Pointer(&last)),
		(*C.int)(unsafe.Pointer(&fid)),
	)

	// check error
	err = getErr(ier)
	return
}

// Qagie performs integration over infinite intervals
//
//   INPUT:
//     fid -- id of function to avoid goroutine problems
//     f   -- function defining the integrand
//
//     bound -- finite bound of integration range
//              (has no meaning if interval is doubly-infinite)
//
//     inf -- indicates the kind of integration range involved:
//              inf = 1 corresponds to  (bound,+infinity),
//              inf = -1            to  (-infinity,bound),
//              inf = 2             to (-infinity,+infinity).
//
//     epsabs -- absolute accuracy requested [use ≤0 for default]
//     epsrel -- relative accuracy requested [use ≤0 for default]
//
//   INPUT/OUTPUT:
//     NOTE: (1) the length of the 5 vectors below is equal to the "limit" variable in the original
//               code which is an upperbound on the number of subintervals in the partition of (a,b)
//           (2) the 5 vectors below may be <nil>, thus memory is allocated here
//
//     alist -- the first last  elements of which are the left
//              end points of the subintervals in the partition of the given integration range (a,b)
//
//     blist -- the first last elements of which are the right
//              end points of the subintervals in the partition of the given integration range (a,b)
//
//     rlist -- the first last elements of which are the integral
//              approximations on the subintervals
//
//     elist -- the first last  elements of which are the moduli
//              of the absolute error estimates on the subintervals
//
//     iord  -- the first k elements of which are pointers to
//              the error estimates over the subintervals, such that elist(iord(1)), ...,
//              elist(iord(k)) form a decreasing sequence, with k = last if last.le.(limit/2+2), and
//              k = limit+1-last otherwise
//   OUTPUT:
//     result -- approximation to the integral
//     abserr -- estimate of the modulus of the absolute error, which should equal or exceed abs(i-result)
//     neval  -- number of integrand evaluations
//     last   -- number of subintervals actually produced in the subdivision process
//
func Qagie(fid int, f fType, bound float64, inf int32, epsabs, epsrel float64, alist, blist, rlist, elist []float64, iord []int32) (result, abserr float64, neval, last int32, err error) {

	// set function in database
	if fid >= len(functions) {
		err = chk.Err("functions database capacity exceeded. max number of functions = %d\n", len(functions))
		return
	}
	functions[fid] = f

	// default values
	if epsabs <= 0 {
		epsabs = 1.49e-8
	}
	if epsrel <= 0 {
		epsrel = 1.49e-8
	}

	// allocate vectors
	limit := len(alist)
	if limit < 1 {
		limit = 50
		alist = make([]float64, limit)
		blist = make([]float64, limit)
		rlist = make([]float64, limit)
		elist = make([]float64, limit)
		iord = make([]int32, limit)
	}

	// call quadpack
	var ier int32
	C.dqagie_(
		C.fType(C.fcn),
		(*C.double)(unsafe.Pointer(&bound)),
		(*C.int)(unsafe.Pointer(&inf)),
		(*C.double)(unsafe.Pointer(&epsabs)),
		(*C.double)(unsafe.Pointer(&epsrel)),
		(*C.int)(unsafe.Pointer(&limit)),
		(*C.double)(unsafe.Pointer(&result)),
		(*C.double)(unsafe.Pointer(&abserr)),
		(*C.int)(unsafe.Pointer(&neval)),
		(*C.int)(unsafe.Pointer(&ier)),
		(*C.double)(unsafe.Pointer(&alist[0])),
		(*C.double)(unsafe.Pointer(&blist[0])),
		(*C.double)(unsafe.Pointer(&rlist[0])),
		(*C.double)(unsafe.Pointer(&elist[0])),
		(*C.int)(unsafe.Pointer(&iord[0])),
		(*C.int)(unsafe.Pointer(&last)),
		(*C.int)(unsafe.Pointer(&fid)),
	)

	/*
		io.PfYel("f(bound) = %v\n", f(bound+math.SmallestNonzeroFloat64))
		io.PfYel("result = %v\n", result)
		io.PfYel("alist = %v\n", alist)
		io.PfYel("blist = %v\n", blist)
		io.PfYel("rlist = %v\n", rlist)
		io.PfYel("elist = %v\n", elist)
		io.PfYel("iord = %v\n", iord)
		io.PfYel("ier = %v\n", ier)
	*/

	// check error
	err = getErr(ier)
	return
}

// getErr returns error message
func getErr(ier int32) error {
	if ier == 0 {
		return nil
	}
	switch ier {
	case 1:
		return chk.Err("error # 1: maximum number of subdivisions allowed\n")
	case 2:
		return chk.Err("error # 2: the occurrence of roundoff error is detected\n")
	case 3:
		return chk.Err("error # 3: extremely bad integrand behaviour\n")
	case 4:
		return chk.Err("error # 4: the algorithm does not converge\n")
	case 5:
		return chk.Err("error # 5: the integral is probably divergent, or slowly convergent\n")
	case 6:
		return chk.Err("error # 6: the input is invalid\n")
	}
	return chk.Err("unknown error\n")
}

//export gofcn
func gofcn(x float64, fid int) float64 {
	return functions[fid](x)
}
