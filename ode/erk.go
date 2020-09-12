// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"gosl/chk"
	"gosl/la"
	"gosl/utl"
)

// ExplicitRK implements explicit Runge-Kutta methods
//
//   The methods available are:
//     moeuler    -- 2(1) Modified-Euler 2(1) ⇒ q = 1
//     rk2        -- 2 Runge, order 2 (mid-point). page 135 of [1]
//     rk3        -- 3 Runge, order 3. page 135 of [1]
//     heun3      -- 3 Heun, order 3. page 135 of [1]
//     rk4        -- 4 "The" Runge-Kutta method. page 138 of [1]
//     rk4-3/8    -- 4 Runge-Kutta method: 3/8-Rule. page 138 of [1]
//     merson4    -- 4 Merson 4("5") method. "5" means that the order 5 is for linear equations with constant coefficients; otherwise the method is of order3. page 167 of [1]
//     zonneveld4 -- 4 Zonneveld 4(3). page 167 of [1]
//     fehlberg4  -- 4(5) Fehlberg 4(5) ⇒ q = 4
//     dopri5     -- 5(4) Dormand-Prince 5(4) ⇒ q = 4
//     verner6    -- 6(5) Verner 6(5) ⇒ q = 5
//     fehlberg7  -- 7(8) Fehlberg 7(8) ⇒ q = 7
//     dopri8     -- 8(5,3) Dormand-Prince 8 order with 5,3 estimator
//  where p(q) means method of p-order with embedded estimator of q-order
//
//  References:
//    [1] E. Hairer, S. P. Nørsett, G. Wanner (2008) Solving Ordinary Differential Equations I.
//        Nonstiff Problems. Second Revised Edition. Corrected 3rd printing 2008. Springer Series
//        in Computational Mathematics ISSN 0179-3632, 528p
//
//  NOTE: (1) Fehlberg's methods give identically zero error estimates for quadrature problems
//            y'=f(x); see page 180 of [1]
//
type ExplicitRK struct {

	// constants
	FSAL     bool        // can use previous ks to compute k0; i.e. k0 := ks{previous]. first same as last [1, page 167]
	Embedded bool        // has embedded error estimator
	A        [][]float64 // A coefficients
	B        []float64   // B coefficients
	Be       []float64   // B coefficients [may be nil, e.g. if FSAL = false]
	C        []float64   // C coefficients
	E        []float64   // error coefficients. difference between B and Be: e = b - be (if be is not nil)
	Nstg     int         // number of stages = len(A) = len(B) = len(C)
	P        int         // order of y1 (corresponding to b)
	Q        int         // order of error estimator (embedded only); e.g. DoPri5(4) ⇒ q = 4 (=min(order(y1),order(y1bar))
	Ad       [][]float64 // A coefficients for dense output
	Cd       []float64   // C coefficients for dense output
	D        [][]float64 // dense output coefficients. [may be nil]

	// data
	ndim int     // problem dimension
	conf *Config // configuration
	work *rkwork // workspace
	stat *Stat   // statistics
	fcn  Func    // dy/dx = f(x,y) function

	// auxiliary
	w    la.Vector // local workspace
	n    float64   // exponent n = 1/(q+1) (or 1/(q+1)-0.75⋅β) of rerrⁿ
	dmin float64   // dmin = 1/Mmin
	dmax float64   // dmax = 1/Mmax
	ndf  float64   // float64(ndim)

	// 5(3) error estimator
	err53 bool    // use 5-3 error estimator
	bhh1  float64 // error estimator: coefficient of k0
	bhh2  float64 // error estimator: coefficient of k8
	bhh3  float64 // error estimator: coefficient of k11

	// dense output
	do []la.Vector // dense output coefficients [nstgDense][ndim] (partially allocated by newERK method)
	kd []la.Vector // k values for dense output [nextraKs] (partially allocated by newERK method)
	yd la.Vector   // y values for dense output (allocated here if len(kd)>0)

	// functions to compute variables for dense output
	dfunA func(y0 la.Vector, x0 float64)                                // in Accept function
	dfunB func(yout la.Vector, h, x float64, y la.Vector, xout float64) // DenseOut function
}

// Free releases memory
func (o *ExplicitRK) Free() {}

// Info returns information about this method
func (o *ExplicitRK) Info() (fixedOnly, implicit bool, nstages int) {
	return !o.Embedded, false, o.Nstg
}

// Init initialises structure
func (o *ExplicitRK) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet) {

	// data
	o.ndim = ndim
	o.conf = conf
	o.work = work
	o.stat = stat
	o.fcn = fcn

	// auxiliary
	o.w = la.NewVector(o.ndim)
	if o.conf.StabBeta > 0 { // lund-stabilization
		o.n = 1.0/float64(o.Q+1) - o.conf.StabBeta*o.conf.stabBetaM
	} else {
		o.n = 1.0 / float64(o.Q+1)
	}
	o.dmin = 1.0 / o.conf.Mmin
	o.dmax = 1.0 / o.conf.Mmax
	o.ndf = float64(ndim)

	// dense output
	if o.conf.denseOut {
		if o.do == nil {
			chk.Panic("dense output is not available for %q\n", o.conf.method)
		}
		for i := 0; i < len(o.do); i++ {
			o.do[i] = la.NewVector(ndim)
		}
		for i := 0; i < len(o.kd); i++ {
			o.kd[i] = la.NewVector(ndim)
		}
		if len(o.kd) > 0 {
			o.yd = la.NewVector(ndim)
		}
	}
}

// Accept accepts update and computes next stepsize
func (o *ExplicitRK) Accept(y0 la.Vector, x0 float64) (dxnew float64) {

	// store data for future dense output
	if o.conf.denseOut {
		if o.dfunA != nil {
			o.dfunA(y0, x0)
		}
	}

	// update y
	y0.Apply(1, o.w)

	// update k0
	if o.FSAL {
		o.work.f[0].Apply(1, o.work.f[o.Nstg-1]) // k0 := ks for next step
	}

	// estimate new stepsize
	if !o.Embedded {
		return
	}
	d := math.Pow(o.work.rerr, o.n)
	if o.conf.StabBeta > 0 { // lund-stabilization
		d = d / math.Pow(o.work.rerrPrev, o.conf.StabBeta)
	}
	d = utl.Max(o.dmax, utl.Min(o.dmin, d/o.conf.Mfac)) // we require  fac1 <= hnew/h <= fac2
	dxnew = o.work.h / d
	return
}

// Reject processes step rejection and computes next stepsize
func (o *ExplicitRK) Reject() (dxnew float64) {

	// estimate new stepsize
	d := math.Pow(o.work.rerr, o.n) / o.conf.Mfac
	dxnew = o.work.h / utl.Min(o.dmin, d)
	return
}

// DenseOut produces dense output (after Accept)
func (o *ExplicitRK) DenseOut(yout la.Vector, h, x float64, y la.Vector, xout float64) {
	o.dfunB(yout, h, x, y, xout)
}

// Step steps update
func (o *ExplicitRK) Step(xa float64, ya la.Vector) {

	// auxiliary
	h := o.work.h
	k := o.work.f
	v := o.work.v

	// compute k0 (otherwise, use k0 saved in Accept)
	if (o.work.first || !o.FSAL) && !o.work.reject {
		u0 := xa + h*o.C[0]
		o.stat.Nfeval++
		o.fcn(k[0], h, u0, ya) // k0 := f(ui,vi)
	}

	// compute ki
	var ui float64
	for i := 1; i < o.work.nstg; i++ {
		ui = xa + h*o.C[i]
		v[i].Apply(1, ya)        // vi := ya
		for j := 0; j < i; j++ { // lower diagonal ⇒ explicit
			la.VecAdd(v[i], 1, v[i], h*o.A[i][j], k[j]) // vi += h⋅aij⋅kj
		}
		o.stat.Nfeval++
		o.fcn(k[i], h, ui, v[i]) // ki := f(ui,vi)
	}

	// update
	if !o.Embedded {
		for m := 0; m < o.ndim; m++ {
			o.w[m] = ya[m]
			for i := 0; i < o.Nstg; i++ {
				o.w[m] += o.B[i] * k[i][m] * h
			}
		}
		return
	}

	// error estimation with 5 and 3 orders (e.g. DoPri853)
	var dk, dv, snum, sden float64 // for stiffness estimation
	if o.err53 {
		var sk, errA, errB, err3, err5 float64
		for m := 0; m < o.ndim; m++ {
			o.w[m] = ya[m]
			errA, errB = 0.0, 0.0
			for i := 0; i < o.Nstg; i++ {
				o.w[m] += o.B[i] * k[i][m] * h
				errA += o.B[i] * k[i][m]
				errB += o.E[i] * k[i][m]
			}
			sk = o.conf.atol + o.conf.rtol*utl.Max(math.Abs(ya[m]), math.Abs(o.w[m]))
			errA -= (o.bhh1*k[0][m] + o.bhh2*k[8][m] + o.bhh3*k[11][m])
			err3 += (errA / sk) * (errA / sk)
			err5 += (errB / sk) * (errB / sk)
			// stiffness estimation
			dk = k[o.Nstg-1][m] - k[o.Nstg-2][m]
			dv = v[o.Nstg-1][m] - v[o.Nstg-2][m]
			snum += dk * dk
			sden += dv * dv
		}
		den := err5 + 0.01*err3 // similar to Eq. (10.17) of [1, page 255]
		if den <= 0.0 {
			den = 1.0
		}
		o.work.rerr = math.Abs(h) * err5 * math.Sqrt(1.0/(o.ndf*den))
		if sden > 0 {
			o.work.rs = h * math.Sqrt(snum/sden)
		}
		return
	}

	// update, error and stiffness estimation
	var kh, sum, lerrm, sk, ratio float64 // lerr[m] component of local error estimate
	for m := 0; m < o.ndim; m++ {
		o.w[m] = ya[m]
		lerrm = 0.0 // must be zeroed for each m
		for i := 0; i < o.Nstg; i++ {
			kh = k[i][m] * h
			o.w[m] += o.B[i] * kh
			lerrm += o.E[i] * kh
		}
		sk = o.conf.atol + o.conf.rtol*utl.Max(math.Abs(ya[m]), math.Abs(o.w[m]))
		ratio = lerrm / sk
		sum += ratio * ratio
		// stiffness estimation
		dk = k[o.Nstg-1][m] - k[o.Nstg-2][m]
		dv = v[o.Nstg-1][m] - v[o.Nstg-2][m]
		snum += dk * dk
		sden += dv * dv
	}
	o.work.rerr = utl.Max(math.Sqrt(sum/o.ndf), 1.0e-10)
	if sden > 0 {
		o.work.rs = h * math.Sqrt(snum/sden)
	}
}

// newERK returns the coefficients of the explicit Runge-Kutta method
//  NOTE: q = min(order(y),order(ybar))
func newERK(kind string) rkmethod {

	// new dataset
	o := new(ExplicitRK)

	// set coefficients
	switch kind {
	case "moeuler": // Modified-Euler 2(1) ⇒ q = 1
		o.Embedded = true
		o.A = [][]float64{
			{0.0, 0.0},
			{1.0, 0.0},
		}
		o.B = []float64{1.0 / 2.0, 1.0 / 2.0}
		o.Be = []float64{1.0, 0.0}
		o.C = []float64{0.0, 1.0}
		o.E = []float64{-1.0 / 2.0, 1.0 / 2.0}
		o.P = 2
		o.Q = 1

	case "rk2": // Runge, order 2 (mid-point). page 135 of [1]
		o.A = [][]float64{
			{0.0, 0.0},
			{1.0 / 2.0, 0.0},
		}
		o.B = []float64{0.0, 1.0}
		o.C = []float64{0.0, 1.0 / 2.0}
		o.P = 2

	case "rk3": // Runge, order 3. page 135 of [1]
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0},
			{1.0 / 2.0, 0.0, 0.0, 0.0},
			{0.0, 1.0, 0.0, 0.0},
			{0.0, 0.0, 1.0, 0.0},
		}
		o.B = []float64{1.0 / 6.0, 2.0 / 3.0, 0.0, 1.0 / 6.0}
		o.C = []float64{0.0, 1.0 / 2.0, 1.0, 1.0}
		o.P = 3

	case "heun3": // Heun, order 3. page 135 of [1]
		o.A = [][]float64{
			{0.0, 0.0, 0.0},
			{1.0 / 3.0, 0.0, 0.0},
			{0.0, 2.0 / 3.0, 0.0},
		}
		o.B = []float64{1.0 / 4.0, 0.0, 3.0 / 4.0}
		o.C = []float64{0.0, 1.0 / 3.0, 2.0 / 3.0}
		o.P = 3

	case "rk4": // "The" Runge-Kutta method. page 138 of [1]
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0},
			{1.0 / 2.0, 0.0, 0.0, 0.0},
			{0.0, 1.0 / 2.0, 0.0, 0.0},
			{0.0, 0.0, 1.0, 0.0},
		}
		o.B = []float64{1.0 / 6.0, 2.0 / 6.0, 2.0 / 6.0, 1.0 / 6.0}
		o.C = []float64{0.0, 1.0 / 2.0, 1.0 / 2.0, 1.0}
		o.P = 4

	case "rk4-3/8": // Runge-Kutta method: 3/8-Rule. page 138 of [1]
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0},
			{1.0 / 3.0, 0.0, 0.0, 0.0},
			{-1.0 / 3.0, 1.0, 0.0, 0.0},
			{1.0, -1.0, 1.0, 0.0},
		}
		o.B = []float64{1.0 / 8.0, 3.0 / 8.0, 3.0 / 8.0, 1.0 / 8.0}
		o.C = []float64{0.0, 1.0 / 3.0, 2.0 / 3.0, 1.0}
		o.P = 4

	case "merson4": // Merson 4("5") method. "5" means that the order 5 is for linear equations with constant coefficients; otherwise the method is of order3. page 167 of [1]
		o.Embedded = true
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{1.0 / 3.0, 0.0, 0.0, 0.0, 0.0},
			{1.0 / 6.0, 1.0 / 6.0, 0.0, 0.0, 0.0},
			{1.0 / 8.0, 0.0, 3.0 / 8.0, 0.0, 0.0},
			{1.0 / 2.0, 0.0, -3.0 / 2.0, 2.0, 0.0},
		}
		o.B = []float64{1.0 / 6.0, 0.0, 0.0, 2.0 / 3.0, 1.0 / 6.0}
		o.Be = []float64{1.0 / 10.0, 0.0, 3.0 / 10.0, 2.0 / 5.0, 1.0 / 5.0}
		o.C = []float64{0.0, 1.0 / 3.0, 1.0 / 3.0, 1.0 / 2.0, 1.0}
		o.E = []float64{1.0 / 15.0, 0.0, -3.0 / 10.0, 4.0 / 15.0, -1.0 / 30.0}
		o.P = 4
		o.Q = 3

	case "zonneveld4": // Zonneveld 4(3). page 167 of [1]
		o.Embedded = true
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{1.0 / 2.0, 0.0, 0.0, 0.0, 0.0},
			{0.0, 1.0 / 2.0, 0.0, 0.0, 0.0},
			{0.0, 0.0, 1.0, 0.0, 0.0},
			{5.0 / 32.0, 7.0 / 32.0, 13.0 / 32.0, -1.0 / 32.0, 0.0},
		}
		o.B = []float64{1.0 / 6.0, 1.0 / 3.0, 1.0 / 3.0, 1.0 / 6.0, 0.0}
		o.Be = []float64{-1.0 / 2.0, 7.0 / 3.0, 7.0 / 3.0, 13.0 / 6.0, -16.0 / 3.0}
		o.C = []float64{0.0, 1.0 / 2.0, 1.0 / 2.0, 1.0, 3.0 / 4.0}
		o.E = []float64{2.0 / 3.0, -2.0, -2.0, -2.0, 16.0 / 3.0}
		o.P = 4
		o.Q = 3

	case "fehlberg4": // Fehlberg 4(5) ⇒ q = 4
		o.Embedded = true
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{1.0 / 4.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{3.0 / 32.0, 9.0 / 32.0, 0.0, 0.0, 0.0, 0.0},
			{1932.0 / 2197.0, -7200.0 / 2197.0, 7296.0 / 2197.0, 0.0, 0.0, 0.0},
			{439.0 / 216.0, -8.0, 3680.0 / 513.0, -845.0 / 4104.0, 0.0, 0.0},
			{-8.0 / 27.0, 2.0, -3544.0 / 2565.0, 1859.0 / 4104.0, -11.0 / 40.0, 0.0},
		}
		o.B = []float64{25.0 / 216.0, 0.0, 1408.0 / 2565.0, 2197.0 / 4104.0, -1.0 / 5.0, 0.0}
		o.Be = []float64{16.0 / 135.0, 0.0, 6656.0 / 12825.0, 28561.0 / 56430.0, -9.0 / 50.0, 2.0 / 55.0}
		o.C = []float64{0.0, 1.0 / 4.0, 3.0 / 8.0, 12.0 / 13.0, 1.0, 1.0 / 2.0}
		o.E = []float64{-1.0 / 360.0, 0.0, 128.0 / 4275.0, 2197.0 / 75240.0, -1.0 / 50.0, -2.0 / 55.0}
		o.P = 4
		o.Q = 4

	case "dopri5": // Dormand-Prince 5(4) ⇒ q = 4
		o.FSAL = true
		o.Embedded = true
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{1.0 / 5.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{3.0 / 40.0, 9.0 / 40.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{44.0 / 45.0, -56.0 / 15.0, 32.0 / 9.0, 0.0, 0.0, 0.0, 0.0},
			{19372.0 / 6561.0, -25360.0 / 2187.0, 64448.0 / 6561.0, -212.0 / 729.0, 0.0, 0.0, 0.0},
			{9017.0 / 3168.0, -355.0 / 33.0, 46732.0 / 5247.0, 49.0 / 176.0, -5103.0 / 18656.0, 0.0, 0.0},
			{35.0 / 384.0, 0.0, 500.0 / 1113.0, 125.0 / 192.0, -2187.0 / 6784.0, 11.0 / 84.0, 0.0},
		}
		o.B = []float64{35.0 / 384.0, 0.0, 500.0 / 1113.0, 125.0 / 192.0, -2187.0 / 6784.0, 11.0 / 84.0, 0.0}
		o.Be = []float64{5179.0 / 57600.0, 0.0, 7571.0 / 16695.0, 393.0 / 640.0, -92097.0 / 339200.0, 187.0 / 2100.0, 1.0 / 40.0}
		o.C = []float64{0.0, 1.0 / 5.0, 3.0 / 10.0, 4.0 / 5.0, 8.0 / 9.0, 1.0, 1.0}
		o.E = []float64{71.0 / 57600.0, 0.0, -71.0 / 16695.0, 71.0 / 1920.0, -17253.0 / 339200.0, 22.0 / 525.0, -1.0 / 40.0}
		o.D = [][]float64{{ // dense output of shampine (1986) [1]
			-12715105075.0 / 11282082432.0,  // D1
			0.00000000000000000000000000,    // D2
			87487479700.0 / 32700410799.0,   // D3
			-10690763975.0 / 1880347072.0,   // D4
			701980252875.0 / 199316789632.0, // D5
			-1453857185.0 / 822651844.0,     // D6
			69997945.0 / 29380423.0,         // D7
		}}
		o.P = 5
		o.Q = 4
		o.do = make([]la.Vector, 5)
		o.dfunA = func(y0 la.Vector, x0 float64) {
			h := o.work.h
			k := o.work.f
			var ydiff, bspl float64
			for m := 0; m < o.ndim; m++ {
				ydiff = o.w[m] - y0[m]
				bspl = h*k[0][m] - ydiff
				o.do[0][m] = y0[m]
				o.do[1][m] = ydiff
				o.do[2][m] = bspl
				o.do[3][m] = ydiff - h*k[6][m] - bspl
				o.do[4][m] = o.D[0][0]*k[0][m] + o.D[0][2]*k[2][m] + o.D[0][3]*k[3][m] + o.D[0][4]*k[4][m] + o.D[0][5]*k[5][m] + o.D[0][6]*k[6][m]
				o.do[4][m] *= o.work.h
			}
		}
		o.dfunB = func(yout la.Vector, h, x float64, y la.Vector, xout float64) {
			xold := x - h
			θ := (xout - xold) / h
			uθ := 1.0 - θ
			for i := 0; i < o.ndim; i++ {
				yout[i] = o.do[0][i] + θ*(o.do[1][i]+uθ*(o.do[2][i]+θ*(o.do[3][i]+uθ*o.do[4][i])))
			}
		}

	case "verner6": // Verner 6(5) ⇒ q = 5
		o.Embedded = true
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{1.0 / 6.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{4.0 / 75.0, 16.0 / 75.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{5.0 / 6.0, -8.0 / 3.0, 5.0 / 2.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{-165.0 / 64.0, 55.0 / 6.0, -425.0 / 64.0, 85.0 / 96.0, 0.0, 0.0, 0.0, 0.0},
			{12.0 / 5.0, -8.0, 4015.0 / 612.0, -11.0 / 36.0, 88.0 / 255.0, 0.0, 0.0, 0.0},
			{-8263.0 / 15000.0, 124.0 / 75.0, -643.0 / 680.0, -81.0 / 250.0, 2484.0 / 10625.0, 0.0, 0.0, 0.0},
			{3501.0 / 1720.0, -300.0 / 43.0, 297275.0 / 52632.0, -319.0 / 2322.0, 24068.0 / 84065.0, 0.0, 3850.0 / 26703.0, 0.0},
		}
		o.B = []float64{3.0 / 40.0, 0.0, 875.0 / 2244.0, 23.0 / 72.0, 264.0 / 1955.0, 0.0, 125.0 / 11592.0, 43.0 / 616.0}
		o.Be = []float64{13.0 / 160.0, 0.0, 2375.0 / 5984.0, 5.0 / 16.0, 12.0 / 85.0, 3.0 / 44.0, 0.0, 0.0}
		o.C = []float64{0.0, 1.0 / 6.0, 4.0 / 15.0, 2.0 / 3.0, 5.0 / 6.0, 1.0, 1.0 / 15.0, 1.0}
		o.E = []float64{-1.0 / 160.0, 0.0, -125.0 / 17952.0, 1.0 / 144.0, -12.0 / 1955.0, -3.0 / 44.0, 125.0 / 11592.0, 43.0 / 616.0}
		o.P = 6
		o.Q = 5

	default:
		return newERKhighOrder(kind)
	}

	// set number of stages
	o.Nstg = len(o.A)
	return o
}

// add methods to database
func init() {
	rkmDB["moeuler"] = func() rkmethod { return newERK("moeuler") }
	rkmDB["rk2"] = func() rkmethod { return newERK("rk2") }
	rkmDB["rk3"] = func() rkmethod { return newERK("rk3") }
	rkmDB["heun3"] = func() rkmethod { return newERK("heun3") }
	rkmDB["rk4"] = func() rkmethod { return newERK("rk4") }
	rkmDB["rk4-3/8"] = func() rkmethod { return newERK("rk4-3/8") }
	rkmDB["merson4"] = func() rkmethod { return newERK("merson4") }
	rkmDB["zonneveld4"] = func() rkmethod { return newERK("zonneveld4") }
	rkmDB["fehlberg4"] = func() rkmethod { return newERK("fehlberg4") }
	rkmDB["dopri5"] = func() rkmethod { return newERK("dopri5") }
	rkmDB["verner6"] = func() rkmethod { return newERK("verner6") }
	rkmDB["fehlberg7"] = func() rkmethod { return newERK("fehlberg7") }
	rkmDB["dopri8"] = func() rkmethod { return newERK("dopri8") }
}
