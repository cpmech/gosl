// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// ExplicitRK implements explicit Runge-Kutta methods
//
//   The methods available are:
//     moeuler    -- Modified-Euler 2(1) ⇒ q = 1
//     rk2        -- Runge, order 2 (mid-point). page 135 of [1]
//     rk3        -- Runge, order 3. page 135 of [1]
//     heun3      -- Heun, order 3. page 135 of [1]
//     rk4        -- "The" Runge-Kutta method. page 138 of [1]
//     rk4-3/8    -- Runge-Kutta method: 3/8-Rule. page 138 of [1]
//     merson4    -- Merson 4("5") method. "5" means that the order 5 is for linear equations with constant coefficients; otherwise the method is of order3. page 167 of [1]
//     zonneveld4 -- Zonneveld 4(3). page 167 of [1]
//     dopri5     -- Dormand-Prince 5(4) ⇒ q = 4
//     fehlberg4  -- Fehlberg 4(5) ⇒ q = 4
//     fehlberg7  -- Fehlberg 7(8) ⇒ q = 7
//     verner6    -- Verner 6(5) ⇒ q = 5
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
	UseKsPrev bool        // can use previous ks to compute k0; i.e. k0 := ks{previous]
	Embedded  bool        // has embedded error estimator
	A         [][]float64 // a coefficients
	B         []float64   // b coefficients
	Be        []float64   // be coefficients (may be nil, if Fprev = false)
	C         []float64   // c coefficients
	E         []float64   // difference between b and be: e = b - be (if be is not nil)
	Nstg      int         // number of stages = len(A) = len(B) = len(C)
	P         int         // order of y1 (corresponding to b)
	Q         int         // order of error estimator (embedded only); e.g. DoPri5(4) ⇒ q = 4 (=min(order(y1),order(y1bar))

	// data
	ndim int       // problem dimension
	conf *Config   // configuration
	work *rkwork   // workspace
	stat *Stat     // statistics
	fcn  Func      // dy/dx = f(x,y) function
	w    la.Vector // local workspace
	beta float64   // factor to stabilize step
	n    float64   // exponent n = 1/(q+1) (or 1/(q+1)-0.75⋅β) of rerrⁿ
	dmin float64   // dmin = 1/Mmin
	dmax float64   // dmax = 1/Mmax
}

// Free releases memory
func (o *ExplicitRK) Free() {}

// Info returns information about this method
func (o *ExplicitRK) Info() (fixedOnly, implicit bool, nstages int) {
	return !o.Embedded, false, o.Nstg
}

// Init initialises structure
func (o *ExplicitRK) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet) (err error) {
	o.ndim = ndim
	o.conf = conf
	o.work = work
	o.stat = stat
	o.fcn = fcn
	o.w = la.NewVector(o.ndim)
	if conf.method == "dopri5" {
		o.beta = conf.DP5beta
	}
	o.n = 1.0/float64(o.Q+1) - o.beta*0.75
	o.dmin = 1.0 / o.conf.Mmin
	o.dmax = 1.0 / o.conf.Mmax
	return nil
}

// Accept accepts update and computes next stepsize
func (o *ExplicitRK) Accept(y la.Vector) (dxnew float64) {

	// update y
	y.Apply(1, o.w)

	// update k0
	if o.UseKsPrev {
		o.work.f[0].Apply(1, o.work.f[o.Nstg-1]) // k0 := ks for next step
	}

	// estimate new stepsize
	if !o.Embedded {
		return
	}
	d := math.Pow(o.work.rerr, o.n)
	if o.beta > 0 { // lund-stabilization
		d = d / math.Pow(o.work.rerrPrev, o.beta)
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

// ContOut produces continuous output (after Accept)
func (o *ExplicitRK) ContOut(yout la.Vector, h, x float64, y la.Vector, xout float64) {
	chk.Panic("TODO")
}

// Step steps update
func (o *ExplicitRK) Step(xa float64, ya la.Vector) (err error) {

	// auxiliary
	h := o.work.h
	k := o.work.f
	v := o.work.v

	// compute k0 (otherwise, use k0 saved in Accept)
	if o.work.first || !o.UseKsPrev { // do it also if cannot reuse previous ks
		u0 := xa + h*o.C[0]
		o.stat.Nfeval++
		err = o.fcn(k[0], h, u0, ya) // k0 := f(ui,vi)
		if err != nil {
			return
		}
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
		err = o.fcn(k[i], h, ui, v[i]) // ki := f(ui,vi)
		if err != nil {
			return
		}
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

	// update and error estimation
	var kh, sum, lerrm, ratio float64 // m component of local error estimate
	for m := 0; m < o.ndim; m++ {
		o.w[m] = ya[m]
		lerrm = 0.0 // must be zeroed for each m
		for i := 0; i < o.Nstg; i++ {
			kh = k[i][m] * h
			o.w[m] += o.B[i] * kh
			lerrm += o.E[i] * kh
		}
		sk := o.conf.atol + o.conf.rtol*utl.Max(math.Abs(ya[m]), math.Abs(o.w[m]))
		ratio = lerrm / sk
		sum += ratio * ratio
	}
	o.work.rerr = utl.Max(math.Sqrt(sum/float64(o.ndim)), 1.0e-10)
	return
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
		o.UseKsPrev = true
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
		o.P = 5
		o.Q = 4

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
