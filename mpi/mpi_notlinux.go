// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux,!darwin

// Package mpi wraps the Message Passing Interface for parallel computations
package mpi

import "github.com/cpmech/gosl/chk"

// IsOn tells whether MPI is on or not
//  NOTE: this returns true even after Stop
func IsOn() bool {
	return false
}

// Start initialises MPI
func Start() {
	chk.Panic("\n\nMPI is not available on Windows or macOS yet\n\n")
}

// Stop finalises MPI
func Stop() {
}

// WorldRank returns the processor rank/ID within the World communicator
func WorldRank() (rank int) {
	return 0
}

// WorldSize returns the number of processors in the World communicator
func WorldSize() (size int) {
	return 0
}

// Communicator holds the World communicator or a subset communicator
type Communicator struct {
}

// NewCommunicator creates a new communicator or returns the World communicator
//   ranks -- World indices of processors in this Communicator.
//            use nil or empty to get the World Communicator
func NewCommunicator(ranks []int) (o *Communicator) {
	chk.Panic("\n\nMPI is not available on Windows or macOS yet\n\n")
	return nil
}

// Rank returns the processor rank/ID
func (o *Communicator) Rank() (rank int) {
	return 0
}

// Size returns the number of processors
func (o *Communicator) Size() (size int) {
	return 0
}

// Abort aborts MPI
func (o *Communicator) Abort() {
}

// Barrier forces synchronisation
func (o *Communicator) Barrier() {
}

// BcastFromRoot broadcasts slice from root (Rank == 0) to all other processors
func (o *Communicator) BcastFromRoot(x []float64) {
}

// BcastFromRootC broadcasts slice from root (Rank == 0) to all other processors (complex version)
func (o *Communicator) BcastFromRootC(x []complex128) {
}

// ReduceSum sums all values in 'orig' to 'dest' in root (Rank == 0) processor
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) ReduceSum(dest, orig []float64) {
}

// ReduceSumC sums all values in 'orig' to 'dest' in root (Rank == 0) processor (complex version)
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) ReduceSumC(dest, orig []complex128) {
}

// AllReduceSum combines all values from orig into dest summing values
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceSum(dest, orig []float64) {
}

// AllReduceSumC combines all values from orig into dest summing values (complex version)
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceSumC(dest, orig []complex128) {
}

// AllReduceMin combines all values from orig into dest picking minimum values
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceMin(dest, orig []float64) {
}

// AllReduceMax combines all values from orig into dest picking minimum values
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceMax(dest, orig []float64) {
}

// AllReduceMinI combines all values from orig into dest picking minimum values (integer version)
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceMinI(dest, orig []int) {
}

// AllReduceMaxI combines all values from orig into dest picking minimum values (integer version)
//   NOTE (important): orig and dest must be different slices
func (o *Communicator) AllReduceMaxI(dest, orig []int) {
}

// Send sends values to processor toID
func (o *Communicator) Send(vals []float64, toID int) {
}

// Recv receives values from processor fromId
func (o *Communicator) Recv(vals []float64, fromID int) {
}

// SendC sends values to processor toID (complex version)
func (o *Communicator) SendC(vals []complex128, toID int) {
}

// RecvC receives values from processor fromId (complex version)
func (o *Communicator) RecvC(vals []complex128, fromID int) {
}

// SendI sends values to processor toID (integer version)
func (o *Communicator) SendI(vals []int, toID int) {
}

// RecvI receives values from processor fromId (integer version)
func (o *Communicator) RecvI(vals []int, fromID int) {
}

// SendOne sends one value to processor toID
func (o *Communicator) SendOne(val float64, toID int) {
}

// RecvOne receives one value from processor fromId
func (o *Communicator) RecvOne(fromID int) (val float64) {
	return 0
}

// SendOneI sends one value to processor toID (integer version)
func (o *Communicator) SendOneI(val int, toID int) {
}

// RecvOneI receives one value from processor fromId (integer version)
func (o *Communicator) RecvOneI(fromID int) (val int) {
	return 0
}
