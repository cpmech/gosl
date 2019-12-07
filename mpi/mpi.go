// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

// Package mpi wraps the Message Passing Interface for parallel computations
package mpi

/*
#include "mpi.h"

MPI_Comm     World     = MPI_COMM_WORLD;
MPI_Op       OpSum     = MPI_SUM;
MPI_Op       OpMin     = MPI_MIN;
MPI_Op       OpMax     = MPI_MAX;
MPI_Datatype TyLong    = MPI_LONG;
MPI_Datatype TyDouble  = MPI_DOUBLE;
MPI_Datatype TyComplex = MPI_DOUBLE_COMPLEX;
MPI_Status*  StIgnore  = MPI_STATUS_IGNORE;

#define DOUBLE_COMPLEX double complex
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// IsOn tells whether MPI is on or not
//  NOTE: this returns true even after Stop
func IsOn() bool {
	var flag C.int
	C.MPI_Initialized(&flag)
	if flag != 0 {
		return true
	}
	return false
}

// Start initialises MPI
func Start() {
	C.MPI_Init(nil, nil)
}

// StartThreadSafe initialises MPI thread safe
func StartThreadSafe() error {
	var r int32
	C.MPI_Init_thread(nil, nil, C.MPI_THREAD_MULTIPLE, (*C.int)(unsafe.Pointer(&r)))
	if r != C.MPI_THREAD_MULTIPLE {
		return fmt.Errorf("MPI_THREAD_MULTIPLE can't be set: got %d", r)
	}
	return nil
}

// Stop finalises MPI
func Stop() {
	C.MPI_Finalize()
}

// WorldRank returns the processor rank/ID within the World communicator
func WorldRank() (rank int) {
	var r int32
	C.MPI_Comm_rank(C.World, (*C.int)(unsafe.Pointer(&r)))
	return int(r)
}

// WorldSize returns the number of processors in the World communicator
func WorldSize() (size int) {
	var s int32
	C.MPI_Comm_size(C.World, (*C.int)(unsafe.Pointer(&s)))
	return int(s)
}

// Communicator holds the World communicator or a subset communicator
type Communicator struct {
	comm  C.MPI_Comm
	group C.MPI_Group
}

// NewCommunicator creates a new communicator or returns the World communicator
//   ranks -- World indices of processors in this Communicator.
//            use nil or empty to get the World Communicator
func NewCommunicator(ranks []int) (o *Communicator) {
	o = new(Communicator)
	if len(ranks) == 0 {
		o.comm = C.World
		C.MPI_Comm_group(C.World, &o.group)
		return
	}
	rs := make([]int32, len(ranks))
	for i := 0; i < len(ranks); i++ {
		rs[i] = int32(ranks[i])
	}
	n := C.int(len(ranks))
	r := (*C.int)(unsafe.Pointer(&rs[0]))
	var wgroup C.MPI_Group
	C.MPI_Comm_group(C.World, &wgroup)
	C.MPI_Group_incl(wgroup, n, r, &o.group)
	C.MPI_Comm_create(C.World, o.group, &o.comm)
	return
}

// Rank returns the processor rank/ID
func (o *Communicator) Rank() (rank int) {
	var r int32
	C.MPI_Comm_rank(o.comm, (*C.int)(unsafe.Pointer(&r)))
	return int(r)
}

// Size returns the number of processors
func (o *Communicator) Size() (size int) {
	var s int32
	C.MPI_Comm_size(o.comm, (*C.int)(unsafe.Pointer(&s)))
	return int(s)
}

// Abort aborts MPI
func (o *Communicator) Abort() {
	C.MPI_Abort(o.comm, 0)
}

// Barrier forces synchronisation
func (o *Communicator) Barrier() {
	C.MPI_Barrier(o.comm)
}

// BcastFromRoot broadcasts slice from root (Rank == 0) to all other processors
func (o *Communicator) BcastFromRoot(x []float64) {
	buf := unsafe.Pointer(&x[0])
	C.MPI_Bcast(buf, C.int(len(x)), C.TyDouble, 0, o.comm)
}

// BcastFromRootC broadcasts slice from root (Rank == 0) to all other processors (complex version)
func (o *Communicator) BcastFromRootC(x []complex128) {
	buf := unsafe.Pointer(&x[0])
	C.MPI_Bcast(buf, C.int(len(x)), C.TyComplex, 0, o.comm)
}

// ReduceSum sums all values in 'orig' to 'dest' in root (Rank == 0) processor
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) ReduceSum(dest, orig []float64) {
	sendbuf := unsafe.Pointer(&orig[0])
	recvbuf := unsafe.Pointer(&dest[0])
	C.MPI_Reduce(sendbuf, recvbuf, C.int(len(dest)), C.TyDouble, C.OpSum, 0, o.comm)
}

// ReduceSumC sums all values in 'orig' to 'dest' in root (Rank == 0) processor (complex version)
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) ReduceSumC(dest, orig []complex128) {
	sendbuf := unsafe.Pointer(&orig[0])
	recvbuf := unsafe.Pointer(&dest[0])
	C.MPI_Reduce(sendbuf, recvbuf, C.int(len(dest)), C.TyComplex, C.OpSum, 0, o.comm)
}

// AllReduceSum combines all values from orig into dest summing values
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceSum(dest, orig []float64) {
	sendbuf := unsafe.Pointer(&orig[0])
	recvbuf := unsafe.Pointer(&dest[0])
	C.MPI_Allreduce(sendbuf, recvbuf, C.int(len(dest)), C.TyDouble, C.OpSum, o.comm)
}

// AllReduceSumC combines all values from orig into dest summing values (complex version)
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceSumC(dest, orig []complex128) {
	sendbuf := unsafe.Pointer(&orig[0])
	recvbuf := unsafe.Pointer(&dest[0])
	C.MPI_Allreduce(sendbuf, recvbuf, C.int(len(dest)), C.TyComplex, C.OpSum, o.comm)
}

// AllReduceMin combines all values from orig into dest picking minimum values
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceMin(dest, orig []float64) {
	sendbuf := unsafe.Pointer(&orig[0])
	recvbuf := unsafe.Pointer(&dest[0])
	C.MPI_Allreduce(sendbuf, recvbuf, C.int(len(dest)), C.TyDouble, C.OpMin, o.comm)
}

// AllReduceMax combines all values from orig into dest picking minimum values
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceMax(dest, orig []float64) {
	sendbuf := unsafe.Pointer(&orig[0])
	recvbuf := unsafe.Pointer(&dest[0])
	C.MPI_Allreduce(sendbuf, recvbuf, C.int(len(dest)), C.TyDouble, C.OpMax, o.comm)
}

// AllReduceMinI combines all values from orig into dest picking minimum values (integer version)
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceMinI(dest, orig []int) {
	sendbuf := unsafe.Pointer(&orig[0])
	recvbuf := unsafe.Pointer(&dest[0])
	C.MPI_Allreduce(sendbuf, recvbuf, C.int(len(dest)), C.TyLong, C.OpMin, o.comm)
}

// AllReduceMaxI combines all values from orig into dest picking minimum values (integer version)
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceMaxI(dest, orig []int) {
	sendbuf := unsafe.Pointer(&orig[0])
	recvbuf := unsafe.Pointer(&dest[0])
	C.MPI_Allreduce(sendbuf, recvbuf, C.int(len(dest)), C.TyLong, C.OpMax, o.comm)
}

// Send sends values to processor toID
func (o *Communicator) Send(vals []float64, toID int) {
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Send(buf, C.int(len(vals)), C.TyDouble, C.int(toID), 10000, o.comm)
}

// Recv receives values from processor fromId
func (o *Communicator) Recv(vals []float64, fromID int) {
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Recv(buf, C.int(len(vals)), C.TyDouble, C.int(fromID), 10000, o.comm, C.StIgnore)
}

// SendC sends values to processor toID (complex version)
func (o *Communicator) SendC(vals []complex128, toID int) {
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Send(buf, C.int(len(vals)), C.TyComplex, C.int(toID), 10001, o.comm)
}

// RecvC receives values from processor fromId (complex version)
func (o *Communicator) RecvC(vals []complex128, fromID int) {
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Recv(buf, C.int(len(vals)), C.TyComplex, C.int(fromID), 10001, o.comm, C.StIgnore)
}

// SendI sends values to processor toID (integer version)
func (o *Communicator) SendI(vals []int, toID int) {
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Send(buf, C.int(len(vals)), C.TyLong, C.int(toID), 10002, o.comm)
}

// RecvI receives values from processor fromId (integer version)
func (o *Communicator) RecvI(vals []int, fromID int) {
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Recv(buf, C.int(len(vals)), C.TyLong, C.int(fromID), 10002, o.comm, C.StIgnore)
}

// SendOne sends one value to processor toID
func (o *Communicator) SendOne(val float64, toID int) {
	vals := []float64{val}
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Send(buf, 1, C.TyDouble, C.int(toID), 10003, o.comm)
}

// RecvOne receives one value from processor fromId
func (o *Communicator) RecvOne(fromID int) (val float64) {
	vals := []float64{0}
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Recv(buf, 1, C.TyDouble, C.int(fromID), 10003, o.comm, C.StIgnore)
	return vals[0]
}

// SendOneI sends one value to processor toID (integer version)
func (o *Communicator) SendOneI(val int, toID int) {
	vals := []int{val}
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Send(buf, 1, C.TyLong, C.int(toID), 10004, o.comm)
}

// RecvOneI receives one value from processor fromId (integer version)
func (o *Communicator) RecvOneI(fromID int) (val int) {
	vals := []int{0}
	buf := unsafe.Pointer(&vals[0])
	C.MPI_Recv(buf, 1, C.TyLong, C.int(fromID), 10004, o.comm, C.StIgnore)
	return vals[0]
}
