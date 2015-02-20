// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!appengine,!heroku

package mpi

/*
#cgo CXXFLAGS: -O3 -I. -I/usr/lib/openmpi/include -I/usr/lib/openmpi/include/openmpi -pthread
#cgo LDFLAGS: -L. -L/usr/lib/openmpi/lib -lmpi_cxx -lmpi -ldl -lstdc++
#include "connectmpi.h"
*/
import "C"

// NOTE: -lhwloc does not exist in some debian machines!
// cgo LDFLAGS: -L. -L/usr/lib/openmpi/lib -lmpi_cxx -lmpi -ldl -lhwloc -lstdc++

// get CFLAGS with:  mpic++ -showme:compile
// and LDFLAGS with: mpic++ -showme:link
// remember to add -lstdc++

import (
	"unsafe"
)

// Abort aborts MPI
func Abort() {
	C.abortmpi()
}

// InOn tells whether MPI is on or not
//  NOTE: this returns true even after Stop
func IsOn() bool {
	if C.ison() == 1 {
		return true
	}
	return false
}

// Start initialises MPI
func Start(debug bool) {
	if debug {
		C.startmpi(1)
	} else {
		C.startmpi(0)
	}
}

// Stop finalises MPI
func Stop(debug bool) {
	if debug {
		C.stopmpi(1)
	} else {
		C.stopmpi(0)
	}
}

// Rank returns the processor rank/ID
func Rank() int {
	return int(C.mpirank())
}

// Size returns the number of processors
func Size() int {
	return int(C.mpisize())
}

// Barrier forces synchronisation
func Barrier() {
	C.barrier()
}

// SumToRoot sums all values in 'orig' to 'dest' in root (Rank == 0) processor
//  NOTE: orig and dest must be different slices, i.e. not pointing to the same underlying data structure
func SumToRoot(dest, orig []float64) {
	C.sumtoroot((*C.double)(unsafe.Pointer(&dest[0])), (*C.double)(unsafe.Pointer(&orig[0])), C.int(len(orig)))
}

// BcastFromRoot broadcasts 'x' slice from root (Rank == 0) to all other processors
func BcastFromRoot(x []float64) {
	C.bcastfromroot((*C.double)(unsafe.Pointer(&x[0])), C.int(len(x)))
}

// AllReduceSum combines all values in 'x' from all processors. Corresponding components in
// slice 'x' are added together. 'w' is a workspace with length = len(x). The operations are:
//   w := join_all_with_sum(x)
//   x := w
func AllReduceSum(x, w []float64) {
	for i := 0; i < len(x); i++ {
		w[i] = 0
	}
	C.allreducesum((*C.double)(unsafe.Pointer(&w[0])), (*C.double)(unsafe.Pointer(&x[0])), C.int(len(x)))
	for i := 0; i < len(x); i++ {
		x[i] = w[i]
	}
}

// AllReduceSumAdd combines all values in 'x' from all processors and adds the result to another
// slice 'y'. Corresponding components in slice 'x' are added together.
// 'w' is a workspace with length = len(x). The operations are:
//   w := join_all_with_sum(x)
//   y += w
func AllReduceSumAdd(y, x, w []float64) {
	for i := 0; i < len(x); i++ {
		w[i] = 0
	}
	C.allreducesum((*C.double)(unsafe.Pointer(&w[0])), (*C.double)(unsafe.Pointer(&x[0])), C.int(len(x)))
	for i := 0; i < len(x); i++ {
		y[i] += w[i]
	}
}

// AllReduceMin combines all values in 'x' from all processors. When corresponding components (at the
// same position) exist in a number of processors, the minimum value is selected.
// 'w' is a workspace with length = len(x). The operations are:
//   w := join_all_selecting_min(x)
//   x := w
func AllReduceMin(x, w []float64) {
	for i := 0; i < len(x); i++ {
		w[i] = 0
	}
	C.allreducemin((*C.double)(unsafe.Pointer(&w[0])), (*C.double)(unsafe.Pointer(&x[0])), C.int(len(x)))
	for i := 0; i < len(x); i++ {
		x[i] = w[i]
	}
}

// AllReduceMax combines all values in 'x' from all processors. When corresponding components (at the
// same position) exist in a number of processors, the maximum value is selected.
// 'w' is a workspace with length = len(x). The operations are:
//   w := join_all_selecting_max(x)
//   x := w
func AllReduceMax(x, w []float64) {
	for i := 0; i < len(x); i++ {
		w[i] = 0
	}
	C.allreducemax((*C.double)(unsafe.Pointer(&w[0])), (*C.double)(unsafe.Pointer(&x[0])), C.int(len(x)))
	for i := 0; i < len(x); i++ {
		x[i] = w[i]
	}
}

// IntAllReduceMax combines all (int) values in 'x' from all processors. When corresponding components (at the
// same position) exist in a number of processors, the maximum value is selected.
// 'w' is a workspace with length = len(x). The operations are:
//   w := join_all_selecting_max(x)
//   x := w
func IntAllReduceMax(x, w []int) {
	for i := 0; i < len(x); i++ {
		w[i] = 0
	}
	C.intallreducemax((*C.int)(unsafe.Pointer(&w[0])), (*C.int)(unsafe.Pointer(&x[0])), C.int(len(x)))
	for i := 0; i < len(x); i++ {
		x[i] = w[i]
	}
}

// SingleIntSend sends a single integer 'val' to processor 'to_proc'
func SingleIntSend(val, to_proc int) {
	C.singleintsend(C.int(val), C.int(to_proc))
}

// SingleIntRecv receives a single integer 'val' from processor 'to_proc'
func SingleIntRecv(from_proc int) (val int) {
	res := C.singleintrecv(C.int(from_proc))
	return int(res)
}

// IntSend sends a slice of integers to processor 'to_proc'
func IntSend(vals []int, to_proc int) {
	C.intsend((*C.int)(unsafe.Pointer(&vals[0])), C.int(len(vals)), C.int(to_proc))
}

// IntRecv receives a slice of integers from processor 'from_proc'
//  NOTE: 'vals' must be pre-allocated with the right number of values that will
//        be sent by 'from_proc'
func IntRecv(vals []int, from_proc int) {
	C.intrecv((*C.int)(unsafe.Pointer(&vals[0])), C.int(len(vals)), C.int(from_proc))
}

// DblSend sends a slice of floats to processor 'to_proc'
func DblSend(vals []float64, to_proc int) {
	C.dblsend((*C.double)(unsafe.Pointer(&vals[0])), C.int(len(vals)), C.int(to_proc))
}

// DblRecv receives a slice of floats from processor 'from_proc'
//  NOTE: 'vals' must be pre-allocated with the right number of values that will
//        be sent by 'from_proc'
func DblRecv(vals []float64, from_proc int) {
	C.dblrecv((*C.double)(unsafe.Pointer(&vals[0])), C.int(len(vals)), C.int(from_proc))
}
