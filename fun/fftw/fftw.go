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
)

// Plan1d implements the FFTW3 plan structure; i.e. a "plan" to compute direct or inverse 1D FTs
type Plan1d struct {
	p      C.fftw_plan // FFTW "plan" structure
	Xin    []float64   // input: complex pairs len=2*N
	Xout   []float64   // output: complex pairs len=2*N
	realIn bool        // set with real input?
}

// NewPlan1d allocates a new "plan" to compute 1D Fourier Transforms with complex numbers input
//
//   Computes:
//                      N-1         -i 2 π k l / N
//               X[l] =  Σ  x[k] ⋅ e
//                      k=0
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

	// fix data (changed by 'measure')
	if usingData && measure {
		copy(data, temp)
	}
	return
}

// NewPlan1dReal allocates a new "plan" to compute 1D Fourier Transforms with real numbers input
//   INPUT:
//     data -- is a real array of length N and may be nil, in this case N must be non-zero
//   NOTE: (1) see NewPlan1d for further information on the input
//         (2) the user must remember to call Free to deallocate FFTW data
func NewPlan1dReal(data []float64, N int, inverse, measure bool) (o *Plan1d, err error) {

	// allocate new object and set/allocate input array
	usingData := false
	if len(data) < 2 {
		if N < 2 {
			err = chk.Err("N must be greater than 1 when data==nil or len(data)<2. N=%d is invalid\n", N)
			return
		}
		o = new(Plan1d)
		o.Xin = make([]float64, N)
	} else {
		N = len(data)
		o = new(Plan1d)
		o.Xin = data
		usingData = true
	}
	o.realIn = true

	// allocate output array
	o.Xout = make([]float64, 2*(N/2+1)) // ×2 => complex128

	// set flags
	var flag C.uint = C.FFTW_ESTIMATE
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
	if inverse {
		in := (*C.fftw_complex)(unsafe.Pointer(&o.Xin[0]))
		out := (*C.double)(unsafe.Pointer(&o.Xout[0]))
		o.p = C.fftw_plan_dft_c2r_1d(C.int(N), in, out, flag)
	} else {
		in := (*C.double)(unsafe.Pointer(&o.Xin[0]))
		out := (*C.fftw_complex)(unsafe.Pointer(&o.Xout[0]))
		o.p = C.fftw_plan_dft_r2c_1d(C.int(N), in, out, flag)
	}

	// fix data (changed by 'measure')
	if usingData && measure {
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

// Input sets input value located at "i". Complex numbers input.
//   NOTE: (1) this method does not check for out-of-range indices
//         (2) this method must not be used when the input is initalised as "real"
func (o *Plan1d) Input(i int, v complex128) {
	o.Xin[i*2] = real(v)
	o.Xin[i*2+1] = imag(v)
}

// InputReal sets input value located at "i". Real numbers input.
//   NOTE: (1) this method does not check for out-of-range indices
//         (2) this method must not be used when the input is initalised as "complex"
func (o *Plan1d) InputReal(i int, v float64) {
	o.Xin[i] = v
}

// Output gets output value located at "i". Complex numbers input.
//   NOTE: this method does not check for out-of-range indices
func (o *Plan1d) Output(i int) (v complex128) {
	if o.realIn {
		N := len(o.Xin)
		if i < N/2+1 {
			return complex(o.Xout[i*2], o.Xout[i*2+1])
		} else { // complex conjugate (reversed)
			j := N - i
			return complex(o.Xout[j*2], -o.Xout[j*2+1])
		}
		return
	}
	return complex(o.Xout[i*2], o.Xout[i*2+1])
}

// GetOutput returns a new slice with the output, for real-input or not
func (o *Plan1d) GetOutput() (res []complex128) {
	N := len(o.Xin) / 2
	if o.realIn {
		N = len(o.Xin)
	}
	res = make([]complex128, N)
	for i := 0; i < N; i++ {
		res[i] = o.Output(i)
	}
	return
}

// Execute performs the Fourier transform
func (o *Plan1d) Execute() {
	C.fftw_execute(o.p)
}

// 2d ////////////////////////////////////////////////////////////////////////////////////////////

// Plan2d implements the FFTW3 plan structure; i.e. a "plan" to compute direct or inverse 2D FTs
type Plan2d struct {
	p    C.fftw_plan // FFTW "plan" structure
	n0   int         // length along first dimension
	n1   int         // length along second dimension
	Xin  []float64   // input: complex pairs len=2*N
	Xout []float64   // output: complex pairs len=2*N
}

// NewPlan2d allocates a new "plan" to compute 2D Fourier Transforms
//
//   Computes:
//                      N1-1 N0-1             -i 2 π k1 l1 / N1    -i 2 π k0 l0 / N0
//           X[l0,l1] =   Σ    Σ  x[k0,k1] ⋅ e                  ⋅ e
//                      k1=0 k0=0
//   INPUT:
//
//   TODO: explain the input data better
//
//     N0, N1 -- must not be zero, even if data is provided
//
//   NOTE: the user must remember to call Free to deallocate FFTW data
//
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
//   data = [ a0.R a0.I  a1.R a1.I  a2.R a2.I  a3.R a3.I  a4.R a4.I  a5.R a5.I ]  ⇒ data[r]
//
//        ⇒ len(data) = 2⋅len(a) = 2⋅n0⋅n1
//
//          l =    0            1            2            3            4            5
//          r = 0     1      2     3      4     5      6     7      8     9      10    11
//   data = [ a00.R a00.I  a01.R a01.I  a02.R a02.I  a10.R a10.I  a11.R a11.I  a12.R a12.I ]
//            2⋅l   2⋅l+1  2⋅l   2⋅l+1  2⋅l   2⋅l+1  2⋅l   2⋅l+1  2⋅l   2⋅l+1  2⋅l   2⋅l+1
//
func NewPlan2d(data []float64, N0, N1 int, inverse, inplace, measure bool) (o *Plan2d, err error) {

	// check N0 and N1
	if N0 < 2 || N0%2 > 0 {
		err = chk.Err("N0 must be even and greater than 1. N0=%d is invalid\n", N0)
		return
	}
	if N1 < 2 || N1%2 > 0 {
		err = chk.Err("N1 must be even and greater than 1. N1=%d is invalid\n", N1)
		return
	}

	// allocate new object and set/allocate input array
	usingData := false
	if len(data) < 4 {
		o = new(Plan2d)
		o.Xin = make([]float64, 2*N0*N1)
	} else {
		ldata := len(data)
		if ldata != 2*N0*N1 || ldata%2 > 0 {
			err = chk.Err("len(data) must be even and equal to 2*N0*N1. %d is invalid (should be %d)\n", ldata, 2*N0*N1)
			return
		}
		o = new(Plan2d)
		o.Xin = data
		usingData = true
	}
	o.n0 = N0
	o.n1 = N1

	// set or allocate output array
	if inplace {
		o.Xout = o.Xin
	} else {
		o.Xout = make([]float64, 2*N0*N1)
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
	o.p = C.fftw_plan_dft_2d(C.int(N0), C.int(N1), in, out, sign, flag)

	// fix data
	if usingData && measure {
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

// Input sets input value located at "i,j". NOTE: this method does not check for out-of-range indices
func (o *Plan2d) Input(i, j int, v complex128) {
	l := o.n1*i + j
	o.Xin[2*l] = real(v)
	o.Xin[2*l+1] = imag(v)
}

// Output gets output value located at "i,j". NOTE: this method does not check for out-of-range indices
func (o *Plan2d) Output(i, j int) (v complex128) {
	l := o.n1*i + j
	return complex(o.Xout[2*l], o.Xout[2*l+1])
}

// Execute performs the Fourier transform
func (o *Plan2d) Execute() {
	C.fftw_execute(o.p)
}

// GetOutput gets the output array as a matrix of complex numbers
func (o *Plan2d) GetOutput() (out [][]complex128) {
	out = make([][]complex128, o.n0)
	for i := 0; i < o.n0; i++ {
		out[i] = make([]complex128, o.n1)
		for j := 0; j < o.n1; j++ {
			out[i][j] = o.Output(i, j)
		}
	}
	return
}
