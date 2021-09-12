// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"math/cmplx"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/fftw"
	"github.com/cpmech/gosl/la"
)

// FourierInterp performs interpolation using truncated Fourier series
//
//               N/2 - 1
//                ————          +i k X[j]
//     f(x[j]) =  \     A[k] ⋅ e                   with    X[j] = 2 π j / N
//                /
//                ————
//               k = -N/2                 Eq (2.1.27) of [1]    x ϵ [0, 2π]
//
//     where:
//
//                 N - 1
//             1   ————             -i k X[j]
//     A[k] = ———  \     f(x[j]) ⋅ e              with    X[j] = 2 π j / N
//             N   /
//                 ————
//                j = 0                                  Eq (2.1.25) of [1]
//
//   NOTE: (1) f=u in [1] and A[k] is the tilde(u[k]) of [1]
//         (2) FFTW says "so you should initialize your input data after creating the plan."
//             Therefore, the plan can be created and reused several times.
//             [http://www.fftw.org/fftw3_doc/Planner-Flags.html]
//             Also: "The plan can be reused as many times as needed. In typical high-performance
//             applications, many transforms of the same size are computed"
//             [http://www.fftw.org/fftw3_doc/Introduction.html]
//
//   Create a new object with NewFourierInterp(...) AND deallocate memory with Free()
//
//   Reference:
//     [1] Canuto C, Hussaini MY, Quarteroni A, Zang TA (2006) Spectral Methods: Fundamentals in
//         Single Domains. Springer. 563p
//
type FourierInterp struct {

	// main
	N int        // number of terms. must be power of 2; i.e. N = 2ⁿ
	X la.Vector  // point coordinates == 2⋅π.j/N
	K la.Vector  // k values computed from j such that j = 0...N-1 ⇒ k = -N/2...N/2-1
	A la.VectorC // coefficients for interpolation. from FFT
	S la.VectorC // smothing coefficients

	// computed (U may be set externally)
	U      la.Vector  // values of f(x) at grid points (nodes) X[j]
	Du     la.Vector  // p-order derivative of u
	Du1    la.Vector  // 1st derivative of f(x) at grid points (nodes) X[j]
	Du2    la.Vector  // 2nd derivative of f(x) at grid points (nodes) X[j]
	DuHat  la.VectorC // spectral coefficient corresponding to p-derivative
	Du1Hat la.VectorC // spectral coefficient corresponding to 1st derivative
	Du2Hat la.VectorC // spectral coefficient corresponding to 1st derivative

	// FFTW
	planA   *fftw.Plan1d // "plan" to compute the A coefficients
	planDu  *fftw.Plan1d // "plan" to compute the p-derivative (inverse transform)
	planDu1 *fftw.Plan1d // "plan" to compute the 1st derivative (inverse transform)
	planDu2 *fftw.Plan1d // "plan" to compute the 2nd derivative (inverse transform)

	// workspace
	workAli la.VectorC // values of f(x) at 3⋅N/2-1 grid points (nodes) X[j] to reduce aliasing error
}

// NewFourierInterp allocates a new FourierInterp object
//
//   N -- number of terms. must be even; ideally power of 2, e.g. N = 2ⁿ
//
//   smoothing -- type of smoothing: use SmoNoneKind for no smoothing
//     "" or "none" : no smoothing
//     "lanc"       : Lanczos (sinc)
//     "rcos"       : Raised Cosine
//     "ces"        : Cesaro
//
//   NOTE: remember to call Free in the end to release memory allocatedy by FFTW; e.g.
//         defer o.Free()
//
func NewFourierInterp(N int, smoothing string) (o *FourierInterp) {

	// check
	if N%2 != 0 {
		chk.Panic("N must be even. N=%d is invalid\n", N)
	}

	// allocate
	o = new(FourierInterp)
	o.N = N
	o.X = make([]float64, o.N)
	o.K = make([]float64, o.N)
	o.A = make([]complex128, o.N)
	o.S = make([]complex128, o.N)

	// point coordinates and K values
	n := float64(o.N)
	for j := 0; j < o.N; j++ {
		o.X[j] = 2.0 * math.Pi * float64(j) / n
		o.K[j] = o.CalcK(j)
	}

	// compute smoothing coefficients
	σ := func(k float64) float64 { return 1.0 }
	switch smoothing {
	case "lanc":
		σ = func(k float64) float64 { return Sinc(2 * k * π / n) }
	case "rcos":
		σ = func(k float64) float64 { return (1.0 + math.Cos(2*k*π/n)) / 2.0 }
	case "ces":
		σ = func(k float64) float64 { return 1.0 - math.Abs(k)/(1.0+n/2.0) }
	}
	for j := 0; j < o.N; j++ {
		o.S[j] = complex(σ(o.K[j]), 0)
	}

	// allocate variables to compute A, Du1, and Du2
	o.Du = la.NewVector(o.N)
	o.Du1 = la.NewVector(o.N)
	o.Du2 = la.NewVector(o.N)
	o.DuHat = la.NewVectorC(o.N)
	o.Du1Hat = la.NewVectorC(o.N)
	o.Du2Hat = la.NewVectorC(o.N)
	o.planA = fftw.NewPlan1d(o.A, false, false)
	o.planDu = fftw.NewPlan1d(o.DuHat, true, false)
	o.planDu1 = fftw.NewPlan1d(o.Du1Hat, true, false)
	o.planDu2 = fftw.NewPlan1d(o.Du2Hat, true, false)
	return
}

// Free releases resources allocated for FFTW
func (o *FourierInterp) Free() {
	if o.planA != nil {
		o.planA.Free()
	}
	if o.planDu != nil {
		o.planDu.Free()
	}
	if o.planDu1 != nil {
		o.planDu1.Free()
	}
	if o.planDu2 != nil {
		o.planDu2.Free()
	}
}

// CalcU calculates f(x) at grid points (to be used later with CalcA and/or CalcD)
func (o *FourierInterp) CalcU(f Ss) {
	if len(o.U) != o.N {
		o.U = la.NewVector(o.N)
	}
	for j := 0; j < o.N; j++ {
		o.U[j] = f(o.X[j])
	}
}

// CalcA calculates the coefficients A of the interpolation using (fwd) FFT
//
//                 N - 1
//             1   ————             -i k X[j]
//     A[k] = ———  \     f(x[j]) ⋅ e              with    X[j] = 2 π j / N
//             N   /
//                 ————
//                j = 0                                  Eq (2.1.25) of [1]
//
//  NOTE: remember to set U (or call CalcU) first
//
func (o *FourierInterp) CalcA() {

	// set A[j] with f(x[j]) / N
	n := float64(o.N)
	for j := 0; j < o.N; j++ {
		o.A[j] = complex(o.U[j]/n, 0)
	}

	// run FFT
	o.planA.Execute()
}

// CalcAwithAliasRemoval calculates the coefficients A by using the 3/2-rule to remove alias error
// via the padding method
//
//  NOTE: with the 3/2-rule, the intepolatory property is not exact; i.e. I(xi)≈f(xi) only
//
func (o *FourierInterp) CalcAwithAliasRemoval(f Ss) {

	// allocate workspace
	M := 3*o.N/2 - 1
	if len(o.workAli) != M {
		o.workAli = la.NewVectorC(M)
	}

	// set workAli[j] with f(x[j]) / M
	var fxj float64
	m := float64(M)
	for j := 0; j < M; j++ {
		xj := 2.0 * math.Pi * float64(j) / m
		fxj = f(xj)
		o.workAli[j] = complex(fxj/m, 0)
	}

	// perform Fourier transform to find A[j]
	Dft1d(o.workAli, false)

	// copy spectral coefficients to the right place in the smaller grid
	var jN, jM int // j's corresponding to the N and M series, respectively
	for jN = 0; jN < o.N; jN++ {
		k := int(o.K[jN])
		if k < 0 {
			jM = M + k
		} else {
			jM = k
		}
		o.A[jN] = o.workAli[jM]
	}
}

// I computes the interpolation (with smoothing or not)
//
//               N/2 - 1
//                 ————          +i k x
//     I {f}(x) =  \     A[k] ⋅ e                 x ϵ [0, 2π]
//      N          /
//                 ————
//                k = -N/2                 Eq (2.1.28) of [1]
//
//  NOTE: remember to call CalcA first
//
func (o *FourierInterp) I(x float64) float64 {
	var res complex128
	for j := 0; j < o.N; j++ {
		res += o.S[j] * o.A[j] * cmplx.Exp(complex(0, o.K[j]*x))
	}
	return real(res)
}

// Idiff performs the differentiation of the interpolation; i.e. computes the p-derivative of the
// interpolation (with smoothing or not)
//
//                   p       N/2 - 1
//        p         d(I{f})    ————       p           +i k x
//  res: DI{f}(x) = ——————— =  \     (i⋅k)  ⋅ A[k] ⋅ e
//        N             p      /
//                    dx       ————
//                            k = -N/2                   x ϵ [0, 2π]
//
//  NOTE: remember to call CalcA first
//
func (o *FourierInterp) Idiff(p int, x float64) float64 {
	var res complex128
	pf := float64(p)
	for j := 0; j < o.N; j++ {
		ikp := ImagPowN(p) * complex(math.Pow(o.K[j], pf), 0)
		res += ikp * o.S[j] * o.A[j] * cmplx.Exp(complex(0, o.K[j]*x))
	}
	return real(res)
}

// CalcD calculates the p-derivative of the interpolated function @ grid points using the FFT
// (with smoothing or not)
//
//                  p      |
//                 d(I{f}) |
//         dfdx =  ——————— |             len(res) must be equal to N
//                     p   |
//                   dx    |x=x[j]
//
//   INPUT:
//      p -- derivative order
//
//   OUTPUT:
//      Du and DuHat will contain the results
//
//  NOTE: remember to call CalcA first
//
func (o *FourierInterp) CalcD(p int) {

	// compute hat(Du)
	for j := 0; j < o.N; j++ {
		ikp := ImagXpowN(o.K[j], p)
		o.DuHat[j] = ikp * o.S[j] * o.A[j]
	}

	// run FFT and extract real part
	o.planDu.Execute()
	for j := 0; j < o.N; j++ {
		o.Du[j] = real(o.DuHat[j])
	}
}

// CalcD1 calculates the 1st derivative using function CalcD and internal arrays.
// See function CalcD for further details
//  OUTPUT: the results will be stored in Du1 and Du1Hat
func (o *FourierInterp) CalcD1() {

	// compute hat(Du1)
	for j := 0; j < o.N; j++ {
		ik1 := complex(0, o.K[j]) // (i⋅k)¹ = i⋅k
		o.Du1Hat[j] = ik1 * o.S[j] * o.A[j]
	}

	// run FFT and extract real part
	o.planDu1.Execute()
	for j := 0; j < o.N; j++ {
		o.Du1[j] = real(o.Du1Hat[j])
	}
}

// CalcD2 calculates the 2nd derivative using function CalcD and internal arrays.
// See function CalcD for further details
//  OUTPUT: the results will be stored in Du2 and Du2Hat
func (o *FourierInterp) CalcD2() {

	// compute hat(Du2)
	for j := 0; j < o.N; j++ {
		ik2 := complex(-o.K[j]*o.K[j], 0) // (i⋅k)² = -k²
		o.Du2Hat[j] = ik2 * o.S[j] * o.A[j]
	}

	// run FFT and extract real part
	o.planDu2.Execute()
	for j := 0; j < o.N; j++ {
		o.Du2[j] = real(o.Du2Hat[j])
	}
}

// CalcK computes k-index from j-index where j corresponds to the FFT index
//
//   FFT returns the A coefficients as:
//
//      {A[0], A[1], ..., A[N/2-1], A[-N/2], A[-N/2+1], ... A[-1]}
//
//   k ϵ [-N/2, N/2-1]
//   j ϵ [0, N-1]
//
//   Example with N = 8:
//
//        j=0 ⇒ k=0      j=4 ⇒ k=-4
//        j=1 ⇒ k=1      j=5 ⇒ k=-3
//        j=2 ⇒ k=2      j=6 ⇒ k=-2
//        j=3 ⇒ k=3      j=7 ⇒ k=-1
//
func (o *FourierInterp) CalcK(j int) float64 {
	h := o.N / 2
	k := j - (j/h)*o.N
	return float64(k)
}

// CalcJ computes j-index from k-index where j corresponds to the FFT index
//
//   k ϵ [-N/2, N/2-1]
//   j ϵ [0, N-1]
//
//   Example with N = 8:
//
//        k=0 ⇒ j=0      k=-4 ⇒ j=4
//        k=1 ⇒ j=1      k=-3 ⇒ j=5      j = { N + k  if  k < 0
//        k=2 ⇒ j=2      k=-2 ⇒ j=6          {     k  otherwise
//        k=3 ⇒ j=3      k=-1 ⇒ j=7
//
func (o *FourierInterp) CalcJ(k float64) int {
	if k < 0 {
		return o.N + int(k)
	}
	return int(k)
}
