// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fftw wraps the FFTW library to perform Fourier Transforms
// using the "fast" method by Cooley and Tukey
package fftw

/*
#include "fftw3.h"
*/
import "C"

import "unsafe"

// Plan1d implements the FFTW3 plan structure; i.e. a "plan" to compute direct or inverse 1D FTs
//
//   Computes:
//                      N-1         -i 2 π j k / N                 __
//     forward:  X[k] =  Σ  x[j] ⋅ e                     with i = √-1
//                      j=0
//
//                      N-1         +i 2 π j k / N
//     inverse:  Y[k] =  Σ  y[j] ⋅ e                     thus x[k] = Y[k] / N
//                      j=0
//
//   NOTE: FFTW says "so you should initialize your input data after creating the plan."
//         Therefore, the plan can be created and reused several times.
//         [http://www.fftw.org/fftw3_doc/Planner-Flags.html]
//         Also: "The plan can be reused as many times as needed. In typical high-performance
//         applications, many transforms of the same size are computed"
//         [http://www.fftw.org/fftw3_doc/Introduction.html]
//
//         Create a new Plan1d with NewPlan1d(...) AND deallocate memory with Free()
//
type Plan1d struct {
	p    C.fftw_plan  // FFTW "plan" structure
	data []complex128 // input
}

// NewPlan1d allocates a new "plan" to compute 1D Fourier Transforms
//
//   data    -- [modified] data is a complex array of length N.
//   inverse -- will perform inverse transform; otherwise will perform direct
//              Note: both transforms are non-normalized;
//              i.e. the user will have to multiply by (1/n) if computing inverse transforms
//   measure -- use the FFTW_MEASURE flag for better optimization analysis (slower initialization times)
//              Note: using this flag with given "data" as input will cause the allocation
//              of a temporary array and the execution of two copy commands with size len(data)
//
//   NOTE: (1) the user must remember to call Free to deallocate FFTW data
//         (2) data will be overwritten
//
func NewPlan1d(data []complex128, inverse, measure bool) (o *Plan1d) {

	// allocate new object
	o = new(Plan1d)
	o.data = data

	// set flags
	var sign C.int = C.FFTW_FORWARD
	var flag C.uint = C.FFTW_ESTIMATE
	if inverse {
		sign = C.FFTW_BACKWARD
	}
	if measure {
		flag = C.FFTW_MEASURE
	}

	// the measure flag will change the input; thus a temporary is required
	N := len(data)
	var temp []complex128
	if measure {
		temp = make([]complex128, N)
		copy(temp, data)
	}

	// set FFTW plan
	d := (*C.fftw_complex)(unsafe.Pointer(&o.data[0]))
	o.p = C.fftw_plan_dft_1d(C.int(N), d, d, sign, flag)

	// fix data (changed by 'measure')
	if measure {
		copy(data, temp)
	}
	return
}

// Free frees internal FFTW data
func (o *Plan1d) Free() {
	if o.p != nil {
		C.fftw_destroy_plan(o.p)
	}
}

// Execute performs the Fourier transform
func (o *Plan1d) Execute() {
	C.fftw_execute(o.p)
}
