// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"math/cmplx"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/fftw"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Smoothing kinds
var (

	// No smoothing
	SmoNoneKind = io.NewEnum("None", "fun.smoothing", "L", "No smoothing")

	// Lanczos (sinc) smoothing kind
	SmoLanczosKind = io.NewEnum("Lanczos", "fun.smoothing", "L", "Lanczos (sinc) smoothing kind")

	// Cesaro
	SmoCesaroKind = io.NewEnum("Cesaro", "fun.smoothing", "L", "Lanczos (sinc) smoothing kind")

	// Raised Cosine
	SmoRcosKind = io.NewEnum("Rcos", "fun.smoothing", "L", "Lanczos (sinc) smoothing kind")
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
//   N         -- number of terms. must be even; ideally power of 2, e.g. N = 2ⁿ
//   smoothing -- type of smoothing: use SmoNoneKind for no smoothing
//
//   NOTE: remember to call Free in the end to release memory allocatedy by FFTW; e.g.
//         defer o.Free()
//
func NewFourierInterp(N int, smoothing io.Enum) (o *FourierInterp, err error) {

	// check
	if N%2 != 0 {
		err = chk.Err("N must be even. N=%d is invalid\n", N)
		return
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
	case SmoLanczosKind:
		σ = func(k float64) float64 { return Sinc(2 * k * π / n) }
	case SmoRcosKind:
		σ = func(k float64) float64 { return (1.0 + math.Cos(2*k*π/n)) / 2.0 }
	case SmoCesaroKind:
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
	o.planA, err = fftw.NewPlan1d(o.A, false, false)
	if err != nil {
		return
	}
	o.planDu, err = fftw.NewPlan1d(o.DuHat, true, false)
	if err != nil {
		return
	}
	o.planDu1, err = fftw.NewPlan1d(o.Du1Hat, true, false)
	if err != nil {
		return
	}
	o.planDu2, err = fftw.NewPlan1d(o.Du2Hat, true, false)
	if err != nil {
		return
	}
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
func (o *FourierInterp) CalcU(f Ss) (err error) {
	if len(o.U) != o.N {
		o.U = la.NewVector(o.N)
	}
	for j := 0; j < o.N; j++ {
		o.U[j], err = f(o.X[j])
		if err != nil {
			return
		}
	}
	return
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
func (o *FourierInterp) CalcAwithAliasRemoval(f Ss) (err error) {

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
		fxj, err = f(xj)
		if err != nil {
			return
		}
		o.workAli[j] = complex(fxj/m, 0)
	}

	// perform Fourier transform to find A[j]
	err = Dft1d(o.workAli, false)
	if err != nil {
		return
	}

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
	return
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
	pc := complex(float64(p), 0)
	for j := 0; j < o.N; j++ {
		ik := complex(0, o.K[j])
		ikp := cmplx.Pow(ik, pc)
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
	pf := float64(p)
	for j := 0; j < o.N; j++ {
		ikp := ImagPowN(p) * complex(math.Pow(o.K[j], pf), 0)
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
		c := complex(0, o.K[j]) // c := (i⋅k)¹ = i⋅k
		o.Du1Hat[j] = c * o.S[j] * o.A[j]
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
		c := complex(-o.K[j]*o.K[j], 0) // c := (i⋅k)² = -k²
		o.Du2Hat[j] = c * o.S[j] * o.A[j]
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

// plotting //////////////////////////////////////////////////////////////////////////////////////

// Plot plots interpolated curve
//   option -- 1: plot only f(x)
//             2: plot both f(x) and df/dx(x)
//             3: plot all f(x), df/dx(x) and d^2f/dx^2
//             4: plot only df/dx(x)
//             5: plot only d^2f/dx^2(x)
//             6: plot df^p/dx^p
//   p      -- order of the derivative to plot if option == 6
//   dfdx   -- is the analytic df/dx(x) [optional]
//   d2fdx2 -- is the analytic d^2f/dx^2(x) [optional]
func (o *FourierInterp) Plot(option, p int, f, dfdx, d2fdx2 Ss, argsF, argsI, argsD1, argsD2 *plt.A) {

	// set arguments
	if argsF == nil {
		argsF = &plt.A{L: "f(x)", C: plt.C(0, 1), NoClip: true}
	}
	if argsI == nil {
		argsI = &plt.A{L: "I{f}(x)", C: plt.C(1, 1), NoClip: true}
	}
	if argsD1 == nil {
		argsD1 = &plt.A{L: "D1I{f}(x)", C: plt.C(2, 1), NoClip: true}
	}
	if argsD2 == nil {
		argsD2 = &plt.A{L: "D2I{f}(x)", C: plt.C(3, 1), NoClip: true}
	}

	// graph points
	npts := 2001
	xx := utl.LinSpace(0, 2.0*math.Pi, npts)

	// options
	withF := option == 1 || option == 2 || option == 3
	firstD := option == 2 || option == 3 || option == 4
	secondD := option == 3 || option == 5

	// allocate arrays
	var y1, y2 []float64
	if withF {
		y1 = make([]float64, npts)
		y2 = make([]float64, npts)
	}
	var y3, y4 []float64
	if firstD {
		y3 = make([]float64, npts)
		y4 = make([]float64, npts)
	}
	var y5, y6 []float64
	if secondD {
		y5 = make([]float64, npts)
		y6 = make([]float64, npts)
	}
	var y7 []float64
	if option == 6 {
		y7 = make([]float64, npts)
	}

	// compute values
	for i := 0; i < npts; i++ {
		x := xx[i]
		if withF {
			if f != nil {
				fx, err := f(x)
				if err != nil {
					chk.Panic("f(x) failed:\n%v\n", err)
				}
				y1[i] = fx
			}
			y2[i] = o.I(x)
		}
		if firstD {
			if dfdx != nil {
				dfx, err := dfdx(x)
				if err != nil {
					chk.Panic("df/dx(x) failed:\n%v\n", err)
				}
				y3[i] = dfx
			}
			y4[i] = o.Idiff(1, x)
		}
		if secondD {
			if d2fdx2 != nil {
				ddfx, err := d2fdx2(x)
				if err != nil {
					chk.Panic("d2f/dx2(x) failed:\n%v\n", err)
				}
				y5[i] = ddfx
			}
			y6[i] = o.Idiff(2, x)
		}
		if option == 6 {
			y7[i] = o.Idiff(p, x)
		}
	}

	// plot
	if option == 2 {
		plt.Subplot(2, 1, 1)
	}
	if option == 3 {
		plt.Subplot(3, 1, 1)
	}
	if withF {
		if f != nil {
			plt.Plot(xx, y1, argsF)
		}
		plt.Plot(xx, y2, argsI)
		plt.HideTRborders()
		plt.Gll("$x$", "$f(x)$", nil)
	}
	if option == 2 {
		plt.Subplot(2, 1, 2)
	}
	if option == 3 {
		plt.Subplot(3, 1, 2)
	}
	if firstD {
		argsF.L = "dfdx"
		//plt.Plot(X, yX, argsX)
		if dfdx != nil {
			plt.Plot(xx, y3, argsF)
		}
		plt.Plot(xx, y4, argsD1)
		plt.HideTRborders()
		plt.Gll("$x$", "$\\frac{\\mathrm{d}f(x)}{\\mathrm{d}x}$", nil)
	}
	if option == 3 {
		plt.Subplot(3, 1, 3)
	}
	if secondD {
		argsF.L = "d2fdx2"
		if d2fdx2 != nil {
			plt.Plot(xx, y5, argsF)
		}
		plt.Plot(xx, y6, argsD2)
		plt.HideTRborders()
		plt.Gll("$x$", "$\\frac{\\mathrm{d}^2f(x)}{\\mathrm{d}x^2}$", nil)
	}
	if option == 6 {
		argsI.L = io.Sf("D%d", p)
		plt.Plot(xx, y7, argsI)
		plt.HideTRborders()
		plt.Gll("$x$", io.Sf("$\\frac{\\mathrm{d}^{%d}f(x)}{\\mathrm{d}x^{%d}}$", p, p), nil)
	}
}
