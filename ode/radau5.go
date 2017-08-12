// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"errors"
	"math"
	"sync"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

// Radau5 implements the Radau5 implicit Runge-Kutta method
type Radau5 struct {

	// main
	ndim  int          // problem dimension
	conf  *Config      // configurations
	work  *rkwork      // workspace
	stat  *Stat        // statistics
	fcn   Func         // dy/dx := f(x,y)
	jac   JacF         // Jacobian function: df/dy(x,y)
	dfdy  *la.Triplet  // df/dy matrix
	mtri  *la.Triplet  // M matrix in triplet format
	mmat  *la.CCMatrix // M matrix in compressed-column format
	hasM  bool         // has M matrix
	ready bool         // matrices and solver are ready

	// coefficients
	mni    float64 // Mfac ⋅ (1+2⋅NmaxIt)
	ndf    float64 // float(ndim)
	denLdw float64 // 3 ⋅ o.ndim)
	nmaxit float64 // NmaxIt

	// linear systems solver
	kmatR *la.Triplet      // matrix for the real part
	kmatC *la.TripletC     // matrix for the imag part
	lsR   la.SparseSolver  // solver for the real part
	lsC   la.SparseSolverC // solver for the imag part

	// workspace
	z    [][]float64 // [nstg][ndim] normalised arrays
	w    [][]float64 // [nstg][ndim] workspace
	dw   [][]float64 // [nstg][ndim] workspace (incremental)
	ycol [][]float64 // [nstg][ndim] colocation values
	v12  la.VectorC  // packed (v[1],v[2])
	dw12 la.VectorC  // packed (dw[1],dw[2])

	// error estimate
	ez   la.Vector // [ndim] for error estimate
	lerr la.Vector // [ndim] for error estimate
	rhs  la.Vector // [ndim] for error estimate

	// constants
	C    []float64   // c coefficients
	T    [][]float64 // T matrix
	Ti   [][]float64 // inv(T) matrix
	Alp  float64     // alpha-hat
	Bet  float64     // beta-hat
	Gam  float64     // gamma-hat
	Gam0 float64     // gamma0 coefficient
	E0   float64     // e0 coefficient
	E1   float64     // e1 coefficient
	E2   float64     // e2 coefficient
	Mu1  float64     // collocation: C1    = (4.D0-SQ6)/10.D0
	Mu2  float64     // collocation: C2    = (4.D0+SQ6)/10.D0
	Mu3  float64     // collocation: C1M1  = C1-1.D0
	Mu4  float64     // collocation: C2M1  = C2-1.D0
	Mu5  float64     // collocation: C1MC2 = C1-C2
}

// add method to database
func init() {
	rkmDB["radau5"] = func() rkmethod { return new(Radau5) }
}

// Free releases memory
func (o *Radau5) Free() {
	if o.lsR != nil {
		o.lsR.Free()
	}
	if o.lsC != nil {
		o.lsC.Free()
	}
}

// Info returns information about this method
func (o *Radau5) Info() (fixedOnly, implicit bool, nstages int) {
	return false, true, 3
}

// Init initialises structure
func (o *Radau5) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet) (err error) {

	// main
	o.ndim = ndim
	o.conf = conf
	o.work = work
	o.stat = stat
	o.fcn = fcn
	o.jac = jac
	o.dfdy = new(la.Triplet)
	o.mtri = M
	if M == nil {
		o.mtri = new(la.Triplet)
		if o.conf.distr {
			o.distrM(ndim)
		} else {
			la.SpTriSetDiag(o.mtri, ndim, 1)
		}
	} else {
		o.hasM = true
	}
	o.mmat = o.mtri.ToMatrix(nil)

	// coefficients
	o.mni = o.conf.Mfac * float64(1+2*o.conf.NmaxIt)
	o.ndf = float64(ndim)
	o.denLdw = float64(3 * ndim)
	o.nmaxit = float64(o.conf.NmaxIt)

	// linear systems solver
	o.kmatR = new(la.Triplet)
	o.kmatC = new(la.TripletC)
	o.lsR = la.NewSparseSolver(o.conf.lsKind)
	o.lsC = la.NewSparseSolverC(o.conf.lsKind)

	// workspace
	nstg := 3
	o.z = make([][]float64, nstg)
	o.w = make([][]float64, nstg)
	o.dw = make([][]float64, nstg)
	o.ycol = make([][]float64, nstg)
	for i := 0; i < nstg; i++ {
		o.z[i] = make([]float64, ndim)
		o.w[i] = make([]float64, ndim)
		o.dw[i] = make([]float64, ndim)
		o.ycol[i] = make([]float64, ndim)
	}
	o.v12 = la.NewVectorC(ndim)
	o.dw12 = la.NewVectorC(ndim)

	// error estimate
	o.ez = la.NewVector(ndim)
	o.lerr = la.NewVector(ndim)
	o.rhs = la.NewVector(ndim)

	// constants
	o.initConstants()
	return
}

// Accept accepts update and computes next stepsize
func (o *Radau5) Accept(y la.Vector) (dxnew float64) {

	// update y
	for m := 0; m < o.ndim; m++ {
		// update y
		y[m] += o.z[2][m]
		// collocation polynomial values
		o.ycol[0][m] = (o.z[1][m] - o.z[2][m]) / o.Mu4
		o.ycol[1][m] = ((o.z[0][m]-o.z[1][m])/o.Mu5 - o.ycol[0][m]) / o.Mu3
		o.ycol[2][m] = o.ycol[1][m] - ((o.z[0][m]-o.z[1][m])/o.Mu5-o.z[0][m]/o.Mu1)/o.Mu2
	}

	// estimate new stepsize
	fac := utl.Min(o.conf.Mfac, o.mni/float64(o.work.nit+2*o.conf.NmaxIt))
	div := utl.Max(o.conf.Mmin, utl.Min(o.conf.Mmax, math.Pow(o.work.rerr, 0.25)/fac))
	dxnew = o.work.h / div

	// predictive controller of Gustafsson
	if o.conf.PredCtrl {
		if o.stat.Naccepted > 1 {
			rerr := o.work.rerr
			rprev := o.work.rerrPrev
			r2 := rerr * rerr
			fac := (o.work.hPrev / o.work.h) * math.Pow(r2/rprev, 0.25) / o.conf.Mfac
			fac = utl.Max(o.conf.Mmin, utl.Min(o.conf.Mmax, fac))
			div = utl.Max(div, fac)
			dxnew = o.work.h / div
		}
	}
	return
}

// Reject processes step rejection and computes next stepsize
func (o *Radau5) Reject() (dxnew float64) {
	// estimate new stepsize
	fac := utl.Min(o.conf.Mfac, o.mni/float64(o.work.nit+2*o.conf.NmaxIt))
	div := utl.Max(o.conf.Mmin, utl.Min(o.conf.Mmax, math.Pow(o.work.rerr, 0.25)/fac))
	dxnew = o.work.h / div
	return
}

// ContOut produces continuous output (after Accept)
func (o *Radau5) ContOut(yout la.Vector, h, x float64, y la.Vector, xout float64) {
	s := (xout - x) / h
	for i := 0; i < len(y); i++ {
		yout[i] = y[i] + s*(o.ycol[0][i]+(s-o.Mu4)*(o.ycol[1][i]+(s-o.Mu3)*o.ycol[2][i]))
	}
}

// Step steps update
func (o *Radau5) Step(x0 float64, y0 la.Vector) (err error) {

	// auxiliary
	h := o.work.h
	u := o.work.u
	v := o.work.v
	k := o.work.f

	// factors
	α := o.Alp / h
	β := o.Bet / h
	γ := o.Gam / h

	// Jacobian and decomposition
	if o.work.reuseJdec {
		o.work.reuseJdec = false
	} else {

		// calculate only first Jacobian for all iterations (simple/modified Newton's method)
		if o.work.reuseJ {
			o.work.reuseJ = false
		} else if !o.work.jacIsOK {

			// stat
			o.stat.Njeval++

			// numerical Jacobian
			if o.jac == nil { // numerical
				if o.conf.distr {
					err = num.JacobianMpi(o.conf.comm, o.dfdy, func(fy, yy la.Vector) (e error) {
						e = o.fcn(fy, h, x0, yy)
						return
					}, y0, o.work.f0, o.w[0], true) // w works here as workspace variable
				} else {
					err = num.Jacobian(o.dfdy, func(fy, yy la.Vector) (e error) {
						e = o.fcn(fy, h, x0, yy)
						return
					}, y0, o.work.f0, o.w[0]) // w works here as workspace variable
				}

				// analytical Jacobian
			} else {
				err = o.jac(o.dfdy, h, x0, y0)
			}

			// check
			if err != nil {
				return
			}

			// set flag
			o.work.jacIsOK = true
		}

		// initialise drdy matrix
		if !o.ready {
			o.kmatR.Init(o.ndim, o.ndim, o.mtri.Len()+o.dfdy.Len())
			o.kmatC.Init(o.ndim, o.ndim, o.mtri.Len()+o.dfdy.Len())
		}

		// update matrices
		la.SpTriAdd(o.kmatR, γ, o.mtri, -1, o.dfdy)       // kmatR :=      γ*M - dfdy
		la.SpTriAddR2C(o.kmatC, α, β, o.mtri, -1, o.dfdy) // kmatC := (α+βi)*M - dfdy

		// initialise linear solver
		if !o.ready {
			err = o.lsR.Init(o.kmatR, o.conf.Symmetric, o.conf.LsVerbose, o.conf.Ordering, o.conf.Scaling, o.conf.comm)
			if err != nil {
				return
			}
			err = o.lsC.Init(o.kmatC, o.conf.Symmetric, o.conf.LsVerbose, o.conf.Ordering, o.conf.Scaling, o.conf.comm)
			if err != nil {
				return
			}
			o.ready = true
		}

		// perform factorisation
		o.stat.Ndecomp++
		err = o.lsR.Fact()
		if err != nil {
			return
		}
		err = o.lsC.Fact()
		if err != nil {
			return
		}
	}

	// update u[i]
	u[0] = x0 + o.C[0]*h
	u[1] = x0 + o.C[1]*h
	u[2] = x0 + o.C[2]*h

	// zero trial: z[i] and w[i]
	if o.work.first || o.conf.ZeroTrial {
		for m := 0; m < o.ndim; m++ {
			o.z[0][m], o.w[0][m] = 0.0, 0.0
			o.z[1][m], o.w[1][m] = 0.0, 0.0
			o.z[2][m], o.w[2][m] = 0.0, 0.0
		}

		// interpolation polynomial trial: z[i] and w[i]
	} else {
		c3q := h / o.work.hPrev
		c1q := o.Mu1 * c3q
		c2q := o.Mu2 * c3q
		for m := 0; m < o.ndim; m++ {
			o.z[0][m] = c1q * (o.ycol[0][m] + (c1q-o.Mu4)*(o.ycol[1][m]+(c1q-o.Mu3)*o.ycol[2][m]))
			o.z[1][m] = c2q * (o.ycol[0][m] + (c2q-o.Mu4)*(o.ycol[1][m]+(c2q-o.Mu3)*o.ycol[2][m]))
			o.z[2][m] = c3q * (o.ycol[0][m] + (c3q-o.Mu4)*(o.ycol[1][m]+(c3q-o.Mu3)*o.ycol[2][m]))
			o.w[0][m] = o.Ti[0][0]*o.z[0][m] + o.Ti[0][1]*o.z[1][m] + o.Ti[0][2]*o.z[2][m]
			o.w[1][m] = o.Ti[1][0]*o.z[0][m] + o.Ti[1][1]*o.z[1][m] + o.Ti[1][2]*o.z[2][m]
			o.w[2][m] = o.Ti[2][0]*o.z[0][m] + o.Ti[2][1]*o.z[1][m] + o.Ti[2][2]*o.z[2][m]
		}
	}

	// iterations
	nstg := 3
	o.work.nit = 0
	o.work.eta = math.Pow(utl.Max(o.work.eta, o.conf.Eps), 0.8)
	o.work.theta = o.conf.ThetaMax
	o.work.diverg = false
	var errR, errC error
	var Ldw, LdwOld, thq, othq, iterr, itRerr, qnewt, ratio0, ratio1, ratio2, nitf float64
	var it int
	for it = 0; it < o.conf.NmaxIt; it++ {

		// max iterations ?
		o.work.nit = it + 1
		if o.work.nit > o.stat.Nitmax {
			o.stat.Nitmax = o.work.nit
		}

		// evaluate f(x,y) at (u[i],v[i]=y0+z[i])
		for i := 0; i < nstg; i++ {
			for m := 0; m < o.ndim; m++ {
				v[i][m] = y0[m] + o.z[i][m]
			}
			o.stat.Nfeval++
			err = o.fcn(k[i], h, u[i], v[i])
			if err != nil {
				return
			}
		}

		// calc rhs
		if o.hasM {
			if o.conf.distr {
				o.distrDw() // compute dw[i]
			} else {
				// using dw as workspace here
				la.SpMatVecMul(o.dw[0], 1, o.mmat, o.w[0]) // dw0 := M ⋅ w0
				la.SpMatVecMul(o.dw[1], 1, o.mmat, o.w[1]) // dw1 := M ⋅ w1
				la.SpMatVecMul(o.dw[2], 1, o.mmat, o.w[2]) // dw2 := M ⋅ w2
			}
			for m := 0; m < o.ndim; m++ {
				v[0][m] = o.Ti[0][0]*k[0][m] + o.Ti[0][1]*k[1][m] + o.Ti[0][2]*k[2][m] - γ*o.dw[0][m]
				v[1][m] = o.Ti[1][0]*k[0][m] + o.Ti[1][1]*k[1][m] + o.Ti[1][2]*k[2][m] - α*o.dw[1][m] + β*o.dw[2][m]
				v[2][m] = o.Ti[2][0]*k[0][m] + o.Ti[2][1]*k[1][m] + o.Ti[2][2]*k[2][m] - β*o.dw[1][m] - α*o.dw[2][m]
			}
		} else {
			for m := 0; m < o.ndim; m++ {
				v[0][m] = o.Ti[0][0]*k[0][m] + o.Ti[0][1]*k[1][m] + o.Ti[0][2]*k[2][m] - γ*o.w[0][m]
				v[1][m] = o.Ti[1][0]*k[0][m] + o.Ti[1][1]*k[1][m] + o.Ti[1][2]*k[2][m] - α*o.w[1][m] + β*o.w[2][m]
				v[2][m] = o.Ti[2][0]*k[0][m] + o.Ti[2][1]*k[1][m] + o.Ti[2][2]*k[2][m] - β*o.w[1][m] - α*o.w[2][m]
			}
		}

		// solve linear system
		o.stat.Nlinsol++
		if !o.conf.distr && o.conf.GoChan {
			wg := new(sync.WaitGroup)
			wg.Add(2)
			go func() {
				errR = o.lsR.Solve(o.dw[0], v[0], false)
				wg.Done()
			}()
			go func() {
				o.v12.JoinRealImag(v[1], v[2])
				errC = o.lsC.Solve(o.dw12, o.v12, false)
				o.dw12.SplitRealImag(o.dw[1], o.dw[2])
				wg.Done()
			}()
			wg.Wait()
		} else {
			o.v12.JoinRealImag(v[1], v[2])
			errR = o.lsR.Solve(o.dw[0], v[0], false)
			errC = o.lsC.Solve(o.dw12, o.v12, false)
			o.dw12.SplitRealImag(o.dw[1], o.dw[2])
		}

		// check for errors from linear solution
		if errR != nil || errC != nil {
			var errmsg string
			if errR != nil {
				errmsg += errR.Error()
			}
			if errC != nil {
				if errR != nil {
					errmsg += "\n"
				}
				errmsg += errC.Error()
			}
			err = errors.New(errmsg)
			return
		}

		// update w and z
		for m := 0; m < o.ndim; m++ {
			o.w[0][m] += o.dw[0][m]
			o.w[1][m] += o.dw[1][m]
			o.w[2][m] += o.dw[2][m]
			o.z[0][m] = o.T[0][0]*o.w[0][m] + o.T[0][1]*o.w[1][m] + o.T[0][2]*o.w[2][m]
			o.z[1][m] = o.T[1][0]*o.w[0][m] + o.T[1][1]*o.w[1][m] + o.T[1][2]*o.w[2][m]
			o.z[2][m] = o.T[2][0]*o.w[0][m] + o.T[2][1]*o.w[1][m] + o.T[2][2]*o.w[2][m]
		}

		// rms norm of δw
		Ldw = 0.0
		for m := 0; m < o.ndim; m++ {
			ratio0 = o.dw[0][m] / o.work.scal[m]
			ratio1 = o.dw[1][m] / o.work.scal[m]
			ratio2 = o.dw[2][m] / o.work.scal[m]
			Ldw += ratio0*ratio0 + ratio1*ratio1 + ratio2*ratio2
		}
		Ldw = math.Sqrt(Ldw / o.denLdw)

		// check convergence
		if it > 0 {
			thq = Ldw / LdwOld
			if it == 1 {
				o.work.theta = thq
			} else {
				o.work.theta = math.Sqrt(thq * othq)
			}
			othq = thq
			if o.work.theta < 0.99 {
				o.work.eta = o.work.theta / (1.0 - o.work.theta)
				nitf = float64(o.work.nit)
				iterr = Ldw * math.Pow(o.work.theta, o.nmaxit-nitf) / (1.0 - o.work.theta)
				itRerr = iterr / o.conf.fnewt
				if itRerr >= 1.0 { // diverging
					qnewt = utl.Max(1.0e-4, utl.Min(20.0, itRerr))
					o.work.dvfac = 0.8 * math.Pow(qnewt, -1.0/(4.0+o.nmaxit-1.0-nitf))
					o.work.diverg = true
					break
				}
			} else { // diverging badly (unexpected step-rejection)
				o.work.dvfac = 0.5
				o.work.diverg = true
				break
			}
		}

		// save old norm
		LdwOld = Ldw

		// converged
		if o.work.eta*Ldw < o.conf.fnewt {
			break
		}
	}

	// did not converge
	if it == o.conf.NmaxIt-1 {
		err = chk.Err("Radau5 did not converge with nit=%d", o.work.nit)
		return
	}

	// diverging => stop
	if o.work.diverg {
		o.work.rerr = 2.0 // must leave state intact, any rerr is OK
		return
	}

	// error estimate
	o.errorEstimate(x0, y0)
	return
}

// errorEstimate computes error estimate
func (o *Radau5) errorEstimate(x0 float64, y0 la.Vector) (err error) {

	// auxiliary
	h := o.work.h
	v := o.work.v
	k := o.work.f

	// simple strategy => HW-VII p123 Eq.(8.17) (not good for stiff problems)
	if o.conf.LerrStrat == 1 {

		var sum float64
		for m := 0; m < o.ndim; m++ {
			o.ez[m] = o.E0*o.z[0][m] + o.E1*o.z[1][m] + o.E2*o.z[2][m]
			o.lerr[m] = o.Gam0*h*o.work.f0[m] + o.ez[m]
			ratio := o.lerr[m] / o.work.scal[m]
			sum += ratio * ratio
		}
		o.work.rerr = utl.Max(math.Sqrt(sum/o.ndf), 1.0e-10)
		return
	}

	// common
	γ := o.Gam / h
	if o.hasM {
		for m := 0; m < o.ndim; m++ {
			o.ez[m] = o.E0*o.z[0][m] + o.E1*o.z[1][m] + o.E2*o.z[2][m]
			o.rhs[m] = o.work.f0[m]
		}
		if o.conf.distr {
			o.distrAddToRHS(γ)
		} else {
			la.SpMatVecMulAdd(o.rhs, γ, o.mmat, o.ez) // rhs += γ ⋅ M ⋅ ez
		}
	} else {
		for m := 0; m < o.ndim; m++ {
			o.ez[m] = o.E0*o.z[0][m] + o.E1*o.z[1][m] + o.E2*o.z[2][m]
			o.rhs[m] = o.work.f0[m] + γ*o.ez[m]
		}
	}

	// HW-VII p123 Eq.(8.19)
	if o.conf.LerrStrat == 2 {
		o.lsR.Solve(o.lerr, o.rhs, false)
		o.work.rerr = o.rmsNorm(o.lerr)
		return
	}

	// HW-VII p123 Eq.(8.20)
	o.lsR.Solve(o.lerr, o.rhs, false)
	o.work.rerr = o.rmsNorm(o.lerr)
	if !(o.work.rerr < 1.0) {
		if o.work.first || o.work.reject {
			for m := 0; m < o.ndim; m++ {
				v[0][m] = y0[m] + o.lerr[m] // y0perr
			}
			o.stat.Nfeval++
			err = o.fcn(k[0], h, x0, v[0]) // f0perr
			if err != nil {
				return
			}
			if o.hasM {
				o.rhs.Apply(1, k[0]) // rhs := f0perr
				if o.conf.distr {
					o.distrAddToRHS(γ)
				} else {
					la.SpMatVecMulAdd(o.rhs, γ, o.mmat, o.ez) // rhs += γ ⋅ M ⋅ ez
				}
			} else {
				la.VecAdd(o.rhs, 1, k[0], γ, o.ez) // rhs = f0perr + γ ⋅ ez
			}
			o.lsR.Solve(o.lerr, o.rhs, false)
			o.work.rerr = o.rmsNorm(o.lerr)
		}
	}
	return
}

// distrM sets M matrix (distributed version)
func (o *Radau5) distrM(ndim int) {
	id, sz := o.conf.comm.Rank(), o.conf.comm.Size()
	start, endp1 := (id*ndim)/sz, ((id+1)*ndim)/sz
	o.mtri.Init(ndim, ndim, endp1-start)
	for i := start; i < endp1; i++ {
		o.mtri.Put(i, i, 1.0)
	}
}

// distrDw computes dw = M * w (distributed version)
func (o *Radau5) distrDw() {

	// using v as workspace here
	v := o.work.v
	la.SpMatVecMul(v[0], 1, o.mmat, o.w[0]) // v0 := M * w0
	la.SpMatVecMul(v[1], 1, o.mmat, o.w[1]) // v1 := M * w1
	la.SpMatVecMul(v[2], 1, o.mmat, o.w[2]) // v2 := M * w2

	// the AllReduceSum will produce dw
	o.conf.comm.AllReduceSum(o.dw[0], v[0]) // dw0 := M * w0
	o.conf.comm.AllReduceSum(o.dw[1], v[1]) // dw1 := M * w1
	o.conf.comm.AllReduceSum(o.dw[2], v[2]) // dw2 := M * w2
}

// distrAddToRHS adds part of error estimate to rhs (disbributed version)
func (o *Radau5) distrAddToRHS(γ float64) {
	la.SpMatVecMul(o.dw[0], γ, o.mmat, o.ez)   // dw[0] = γ ⋅ M ⋅ ez  (dw[0] is workspace)
	o.conf.comm.AllReduceSum(o.dw[1], o.dw[0]) // dw[1] = join(dw[0])
	la.VecAdd(o.rhs, 1, o.rhs, 1, o.dw[1])     // rhs += dw[1]
}

// rmsNorm computes the RMS norm
func (o *Radau5) rmsNorm(diff la.Vector) (rms float64) {
	var ratio float64
	for m := 0; m < o.ndim; m++ {
		ratio = diff[m] / o.work.scal[m]
		rms += ratio * ratio
	}
	return utl.Max(math.Sqrt(rms/o.ndf), 1.0e-10)
}

// initConstants initialises constants
func (o *Radau5) initConstants() {

	o.C = []float64{(4.0 - math.Sqrt(6.0)) / 10.0, (4.0 + math.Sqrt(6.0)) / 10.0, 1.0}

	o.T = [][]float64{
		{9.1232394870892942792e-02, -0.14125529502095420843, -3.0029194105147424492e-02},
		{0.24171793270710701896, 0.20412935229379993199, 0.38294211275726193779},
		{0.96604818261509293619, 1.0, 0.0},
	}

	o.Ti = [][]float64{
		{4.3255798900631553510, 0.33919925181580986954, 0.54177053993587487119},
		{-4.1787185915519047273, -0.32768282076106238708, 0.47662355450055045196},
		{-0.50287263494578687595, 2.5719269498556054292, -0.59603920482822492497},
	}

	c1 := math.Pow(9.0, 1.0/3.0)
	c2 := math.Pow(3.0, 3.0/2.0)
	c3 := math.Pow(9.0, 2.0/3.0)

	o.Alp = -c1/2.0 + 3.0/(2.0*c1) + 3.0
	o.Bet = (math.Sqrt(3.0)*c1)/2.0 + c2/(2.0*c1)
	o.Gam = c1 - 3.0/c1 + 3.0
	o.Gam0 = c1 / (c3 + 3.0*c1 - 3.0)
	o.E0 = o.Gam0 * (-13.0 - 7.0*math.Sqrt(6.0)) / 3.0
	o.E1 = o.Gam0 * (-13.0 + 7.0*math.Sqrt(6.0)) / 3.0
	o.E2 = o.Gam0 * (-1.0) / 3.0

	o.Mu1 = (4.0 - math.Sqrt(6.0)) / 10.0
	o.Mu2 = (4.0 + math.Sqrt(6.0)) / 10.0
	o.Mu3 = o.Mu1 - 1.0
	o.Mu4 = o.Mu2 - 1.0
	o.Mu5 = o.Mu1 - o.Mu2
}
