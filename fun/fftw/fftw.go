// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

// Package fftw wraps the FFTW library to perform Fourier Transforms
// using the "fast" method by Cooley and Tukey
package fftw

/*
#cgo LDFLAGS: -lfftw3 -lm

#include "fftw3.h"
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// Plan1d implements the FFTW3 plan structure; i.e. a "plan" to compute direct or inverse 1D FTs
type Plan1d struct {
	p    C.fftw_plan // FFTW "plan" structure
	Xin  []float64   // input: complex pairs len=2*N
	Xout []float64   // output: complex pairs len=2*N
}

// NewPlan1d allocates a new "plan" to compute 1D Fourier Transforms
//
//   INPUT:
//
//     data -- is a complex array stored as a real array of length 2*n. [real,imag, real,imag, ...]
//             data may be nil, in this case N must be provided
//
//     N -- half-size of the data array. Use this to allocate input (data) array.
//          data must be non-nil to use N
//
//     inverse -- will perform inverse transform; otherwise will perform direct
//                Note: both transforms are non-normalised;
//                i.e. the user will have to multiply by (1/n) if computing inverse transforms
//
//     inplace -- use data array as output as well, thus data will be overwritten
//
//     measure -- use the FFTW_MEASURE flag for better optimisation analysis (slower initialisation times)
//                Note: using this flag with given "data" as input will cause the allocation
//                of a temporary array and the execution of two copy commands with size len(data)
//
//   NOTE: the user must remember to call Free to deallocate FFTW data
//
func NewPlan1d(data []float64, N int, inverse, inplace, measure bool) (o *Plan1d, err error) {

	// allocate new object and set/allocate input array
	usingData := false
	if len(data) < 2 {
		if N < 2 {
			err = chk.Err("N must be greater than 1 when data==nil or len(data)<2. N=%d is invalid\n", N)
			return
		}
		o = new(Plan1d)
		o.Xin = make([]float64, 2*N)
	} else {
		ldata := len(data)
		if ldata%2 > 0 {
			err = chk.Err("len(data) must be even. %d is invalid\n", ldata)
			return
		}
		N = ldata / 2
		o = new(Plan1d)
		o.Xin = data
		usingData = true
	}

	// set or allocate output array
	if inplace {
		o.Xout = o.Xin
	} else {
		o.Xout = make([]float64, 2*N)
	}

	// set flags
	var sign C.int = C.FFTW_FORWARD
	var flag C.uint = C.FFTW_ESTIMATE
	if inverse {
		sign = C.FFTW_BACKWARD
	}
	if measure {
		flag = C.FFTW_MEASURE
	}

	// the measure flag will change the input; thus a temporary is required if "data" is being used
	var temp []float64
	if usingData && measure {
		temp = make([]float64, len(data))
		copy(temp, data)
	}

	// set FFTW plan
	in := (*C.fftw_complex)(unsafe.Pointer(&o.Xin[0]))
	out := (*C.fftw_complex)(unsafe.Pointer(&o.Xout[0]))
	o.p = C.fftw_plan_dft_1d(C.int(N), in, out, sign, flag)

	// fix data
	if usingData && measure {
		copy(data, temp)
	}
	return
}

// Free frees internal FFTW data
func (o *Plan1d) Free() {
	if o.p != nil {
		C.fftw_destroy_plan(o.p)
		io.PfYel("destroy!\n")
	}
}

// Input sets input value located at "i". NOTE: this method does not check for out-of-range indices
func (o *Plan1d) Input(i int, v complex128) {
	o.Xin[i*2] = real(v)
	o.Xin[i*2+1] = imag(v)
}

// Output gets output value located at "i". NOTE: this method does not check for out-of-range indices
func (o *Plan1d) Output(i int) (v complex128) {
	return complex(o.Xout[i*2], o.Xout[i*2+1])
}

// Execute performs the Fourier transform
func (o *Plan1d) Execute() {
	C.fftw_execute(o.p)
}
