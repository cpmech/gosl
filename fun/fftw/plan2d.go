// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fftw

/*
#include "fftw3.h"
*/
import "C"

import "unsafe"

// Plan2d implements the FFTW3 plan structure; i.e. a "plan" to compute direct or inverse 2D FTs
//
//   Computes:
//                      N1-1 N0-1             -i 2 π k1 l1 / N1    -i 2 π k0 l0 / N0
//           X[l0,l1] =   Σ    Σ  x[k0,k1] ⋅ e                  ⋅ e
//                      k1=0 k0=0
//
type Plan2d struct {
	p    C.fftw_plan  // FFTW "plan" structure
	n0   int          // length along first dimension
	n1   int          // length along second dimension
	data []complex128 // input (row-major matrix)
}

// NewPlan2d allocates a new "plan" to compute 2D Fourier Transforms
//
//   N0, N1  -- dimensions
//   data    -- [modified] data is a complex array of length N0*N1 (row-major matrix)
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
//   A = data is a ROW-MAJOR matrix
//        _                          _
//   A = |   A00→a0  A01→a1  A02→a2   |  ⇒ A[i][j]
//       |_  A10→a3  A11→a4  A12→a5  _|
//                                    (n0,n1)=(2,3)
//
//       l = 0      1      2        3      4      5
//   a = [  A00    A01    A02      A10    A11    A12 ]  ⇒ a[l]
//         3⋅0+0  3⋅0+1  3⋅0+2    3⋅1+0  3⋅1+1  3⋅1+2
//
//   l = n1⋅i + j       i = l // n1      j = l % n1
//
func NewPlan2d(N0, N1 int, data []complex128, inverse, measure bool) (o *Plan2d) {

	// allocate new object
	o = new(Plan2d)
	o.n0 = N0
	o.n1 = N1
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
	var temp []complex128
	if measure {
		temp = make([]complex128, len(data))
		copy(temp, data)
	}

	// set FFTW plan
	d := (*C.fftw_complex)(unsafe.Pointer(&o.data[0]))
	o.p = C.fftw_plan_dft_2d(C.int(N0), C.int(N1), d, d, sign, flag)

	// fix data (changed by 'measure')
	if measure {
		copy(data, temp)
	}
	return
}

// Free frees internal FFTW data
func (o *Plan2d) Free() {
	if o.p != nil {
		C.fftw_destroy_plan(o.p)
	}
}

// Set sets data value located at "i,j". NOTE: this method does not check for out-of-range indices
func (o *Plan2d) Set(i, j int, v complex128) {
	o.data[o.n1*i+j] = v
}

// Get gets data value located at "i,j". NOTE: this method does not check for out-of-range indices
func (o *Plan2d) Get(i, j int) (v complex128) {
	return o.data[o.n1*i+j]
}

// Execute performs the Fourier transform
func (o *Plan2d) Execute() {
	C.fftw_execute(o.p)
}

// GetSlice gets the output array as a nested slice
func (o *Plan2d) GetSlice() (out [][]complex128) {
	out = make([][]complex128, o.n0)
	for i := 0; i < o.n0; i++ {
		out[i] = make([]complex128, o.n1)
		for j := 0; j < o.n1; j++ {
			out[i][j] = o.Get(i, j)
		}
	}
	return
}
