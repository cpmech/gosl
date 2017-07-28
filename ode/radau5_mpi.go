// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package ode

import (
	"errors"
	"math"
	"sync"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
)

// Radau5 step function
func radau5_step_mpi(o *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) {

	// factors
	α := r5.α_ / o.h
	β := r5.β_ / o.h
	γ := r5.γ_ / o.h

	// Jacobian and decomposition
	if o.reuseJdec {
		o.reuseJdec = false
	} else {

		// calculate only first Jacobian for all iterations (simple/modified Newton's method)
		if o.reuseJ {
			o.reuseJ = false
		} else if !o.jacIsOK {

			// Jacobian triplet
			if o.jac == nil { // numerical
				//if x0 == 0.0 { io.Pfgrey(" > > > > > > > > . . . numerical Jacobian . . . < < < < < < < < <\n") }
				err = num.JacobianMpi(o.comm, &o.dfdyT, func(fy, y la.Vector) (e error) {
					e = o.fcn(fy, o.h, x0, y)
					return
				}, y0, o.f0, o.w[0], o.Distr) // w works here as workspace variable
			} else { // analytical
				//if x0 == 0.0 { io.Pfgrey(" > > > > > > > > . . . analytical Jacobian . . . < < < < < < < < <\n") }
				err = o.jac(&o.dfdyT, o.h, x0, y0)
			}
			if err != nil {
				return
			}

			// create M matrix
			if o.doinit && !o.hasM {
				o.mTri = new(la.Triplet)
				if o.Distr {
					id, sz := o.comm.Rank(), o.comm.Size()
					start, endp1 := (id*o.ndim)/sz, ((id+1)*o.ndim)/sz
					o.mTri.Init(o.ndim, o.ndim, endp1-start)
					for i := start; i < endp1; i++ {
						o.mTri.Put(i, i, 1.0)
					}
				} else {
					o.mTri.Init(o.ndim, o.ndim, o.ndim)
					for i := 0; i < o.ndim; i++ {
						o.mTri.Put(i, i, 1.0)
					}
				}
			}
			o.Njeval++
			o.jacIsOK = true
		}

		// initialise triplets
		if o.doinit {
			o.rctriR = new(la.Triplet)
			o.rctriC = new(la.TripletC)
			o.rctriR.Init(o.ndim, o.ndim, o.mTri.Len()+o.dfdyT.Len())
			o.rctriC.Init(o.ndim, o.ndim, o.mTri.Len()+o.dfdyT.Len())
		}

		// update triplets
		la.SpTriAdd(o.rctriR, γ, o.mTri, -1, &o.dfdyT)       // rctriR :=      γ*M - dfdy
		la.SpTriAddR2C(o.rctriC, α, β, o.mTri, -1, &o.dfdyT) // rctriC := (α+βi)*M - dfdy

		// initialise solver
		if o.doinit {
			err = o.lsolR.Init(o.rctriR, o.symmetric, o.lsverbose, o.ordering, o.scaling, o.comm)
			if err != nil {
				return
			}
			err = o.lsolC.Init(o.rctriC, o.symmetric, o.lsverbose, o.ordering, o.scaling, o.comm)
			if err != nil {
				return
			}
		}

		// perform factorisation
		o.lsolR.Fact()
		o.lsolC.Fact()
		o.Ndecomp++
	}

	// updated u[i]
	o.u[0] = x0 + r5.c[0]*o.h
	o.u[1] = x0 + r5.c[1]*o.h
	o.u[2] = x0 + r5.c[2]*o.h

	// (trial/initial) updated z[i] and w[i]
	if o.first || o.ZeroTrial {
		for m := 0; m < o.ndim; m++ {
			o.z[0][m], o.w[0][m] = 0.0, 0.0
			o.z[1][m], o.w[1][m] = 0.0, 0.0
			o.z[2][m], o.w[2][m] = 0.0, 0.0
		}
	} else {
		c3q := o.h / o.hprev
		c1q := r5.μ1 * c3q
		c2q := r5.μ2 * c3q
		for m := 0; m < o.ndim; m++ {
			o.z[0][m] = c1q * (o.ycol[0][m] + (c1q-r5.μ4)*(o.ycol[1][m]+(c1q-r5.μ3)*o.ycol[2][m]))
			o.z[1][m] = c2q * (o.ycol[0][m] + (c2q-r5.μ4)*(o.ycol[1][m]+(c2q-r5.μ3)*o.ycol[2][m]))
			o.z[2][m] = c3q * (o.ycol[0][m] + (c3q-r5.μ4)*(o.ycol[1][m]+(c3q-r5.μ3)*o.ycol[2][m]))
			o.w[0][m] = r5.Ti[0][0]*o.z[0][m] + r5.Ti[0][1]*o.z[1][m] + r5.Ti[0][2]*o.z[2][m]
			o.w[1][m] = r5.Ti[1][0]*o.z[0][m] + r5.Ti[1][1]*o.z[1][m] + r5.Ti[1][2]*o.z[2][m]
			o.w[2][m] = r5.Ti[2][0]*o.z[0][m] + r5.Ti[2][1]*o.z[1][m] + r5.Ti[2][2]*o.z[2][m]
		}
	}

	// iterations
	o.nit = 0
	o.eta = math.Pow(max(o.eta, o.Eps), 0.8)
	o.theta = o.ThetaMax
	o.diverg = false
	var Lδw, oLδw, thq, othq, iterr, itRerr, qnewt float64
	var it int
	for it = 0; it < o.NmaxIt; it++ {

		// max iterations ?
		o.nit = it + 1
		if o.nit > o.Nitmax {
			o.Nitmax = o.nit
		}

		// evaluate f(x,y) at (u[i],v[i]=y0+z[i])
		for i := 0; i < 3; i++ {
			for m := 0; m < o.ndim; m++ {
				o.v[i][m] = y0[m] + o.z[i][m]
			}
			o.Nfeval++
			err = o.fcn(o.f[i], o.h, o.u[i], o.v[i])
			if err != nil {
				return
			}
		}

		// calc rhs
		if o.hasM {
			// using δw as workspace here
			la.SpMatVecMul(o.dw[0], 1, o.mMat, o.w[0]) // δw0 := M * w0
			la.SpMatVecMul(o.dw[1], 1, o.mMat, o.w[1]) // δw1 := M * w1
			la.SpMatVecMul(o.dw[2], 1, o.mMat, o.w[2]) // δw2 := M * w2
			if o.Distr {
				o.comm.AllReduceSum(o.dw[0], o.v[0]) // v is used as workspace here
				o.comm.AllReduceSum(o.dw[1], o.v[1]) // v is used as workspace here
				o.comm.AllReduceSum(o.dw[2], o.v[2]) // v is used as workspace here
			}
			for m := 0; m < o.ndim; m++ {
				o.v[0][m] = r5.Ti[0][0]*o.f[0][m] + r5.Ti[0][1]*o.f[1][m] + r5.Ti[0][2]*o.f[2][m] - γ*o.dw[0][m]
				o.v[1][m] = r5.Ti[1][0]*o.f[0][m] + r5.Ti[1][1]*o.f[1][m] + r5.Ti[1][2]*o.f[2][m] - α*o.dw[1][m] + β*o.dw[2][m]
				o.v[2][m] = r5.Ti[2][0]*o.f[0][m] + r5.Ti[2][1]*o.f[1][m] + r5.Ti[2][2]*o.f[2][m] - β*o.dw[1][m] - α*o.dw[2][m]
			}
		} else {
			for m := 0; m < o.ndim; m++ {
				o.v[0][m] = r5.Ti[0][0]*o.f[0][m] + r5.Ti[0][1]*o.f[1][m] + r5.Ti[0][2]*o.f[2][m] - γ*o.w[0][m]
				o.v[1][m] = r5.Ti[1][0]*o.f[0][m] + r5.Ti[1][1]*o.f[1][m] + r5.Ti[1][2]*o.f[2][m] - α*o.w[1][m] + β*o.w[2][m]
				o.v[2][m] = r5.Ti[2][0]*o.f[0][m] + r5.Ti[2][1]*o.f[1][m] + r5.Ti[2][2]*o.f[2][m] - β*o.w[1][m] - α*o.w[2][m]
			}
		}

		// solve linear system
		o.Nlinsol++
		var errR, errC error
		if !o.Distr && o.Pll {
			wg := new(sync.WaitGroup)
			wg.Add(2)
			go func() {
				errR = o.lsolR.Solve(o.dw[0], o.v[0], false)
				wg.Done()
			}()
			go func() {
				//errC = o.lsolC.Solve(o.dw[1], o.dw[2], o.v[1], o.v[2], false)
				wg.Done()
			}()
			wg.Wait()
		} else {
			errR = o.lsolR.Solve(o.dw[0], o.v[0], false)
			//errC = o.lsolC.Solve(o.dw[1], o.dw[2], o.v[1], o.v[2], false)
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
			o.z[0][m] = r5.T[0][0]*o.w[0][m] + r5.T[0][1]*o.w[1][m] + r5.T[0][2]*o.w[2][m]
			o.z[1][m] = r5.T[1][0]*o.w[0][m] + r5.T[1][1]*o.w[1][m] + r5.T[1][2]*o.w[2][m]
			o.z[2][m] = r5.T[2][0]*o.w[0][m] + r5.T[2][1]*o.w[1][m] + r5.T[2][2]*o.w[2][m]
		}

		// rms norm of δw
		Lδw = 0.0
		for m := 0; m < o.ndim; m++ {
			Lδw += math.Pow(o.dw[0][m]/o.scal[m], 2.0) + math.Pow(o.dw[1][m]/o.scal[m], 2.0) + math.Pow(o.dw[2][m]/o.scal[m], 2.0)
		}
		Lδw = math.Sqrt(Lδw / float64(3*o.ndim))

		// check convergence
		if it > 0 {
			thq = Lδw / oLδw
			if it == 1 {
				o.theta = thq
			} else {
				o.theta = math.Sqrt(thq * othq)
			}
			othq = thq
			if o.theta < 0.99 {
				o.eta = o.theta / (1.0 - o.theta)
				iterr = Lδw * math.Pow(o.theta, float64(o.NmaxIt-o.nit)) / (1.0 - o.theta)
				itRerr = iterr / o.fnewt
				if itRerr >= 1.0 { // diverging
					qnewt = max(1.0e-4, min(20.0, itRerr))
					o.dvfac = 0.8 * math.Pow(qnewt, -1.0/(4.0+float64(o.NmaxIt)-1.0-float64(o.nit)))
					o.diverg = true
					break
				}
			} else { // diverging badly (unexpected step-rejection)
				o.dvfac = 0.5
				o.diverg = true
				break
			}
		}

		// save old norm
		oLδw = Lδw

		// converged
		if o.eta*Lδw < o.fnewt {
			break
		}
	}

	// did not converge
	if it == o.NmaxIt-1 {
		chk.Panic("radau5_step failed with it=%d", it)
	}

	// diverging => stop
	if o.diverg {
		rerr = 2.0 // must leave state intact, any rerr is OK
		return
	}

	// error estimate
	if o.LerrStrat == 1 {

		// simple strategy => HW-VII p123 Eq.(8.17) (not good for stiff problems)
		for m := 0; m < o.ndim; m++ {
			o.ez[m] = r5.e0*o.z[0][m] + r5.e1*o.z[1][m] + r5.e2*o.z[2][m]
			o.lerr[m] = r5.γ0*o.h*o.f0[m] + o.ez[m]
			rerr += math.Pow(o.lerr[m]/o.scal[m], 2.0)
		}
		rerr = max(math.Sqrt(rerr/float64(o.ndim)), 1.0e-10)

	} else {

		// common
		if o.hasM {
			for m := 0; m < o.ndim; m++ {
				o.ez[m] = r5.e0*o.z[0][m] + r5.e1*o.z[1][m] + r5.e2*o.z[2][m]
				o.rhs[m] = o.f0[m]
			}
			if o.Distr {
				la.SpMatVecMul(o.dw[0], γ, o.mMat, o.ez) // δw[0] = γ * M * ez (δw[0] is workspace)
				//o.comm.AllReduceSumAdd(o.rhs, o.dw[0], o.dw[1]) // rhs += join_with_sum(δw[0]) (δw[1] is workspace)
			} else {
				la.SpMatVecMulAdd(o.rhs, γ, o.mMat, o.ez) // rhs += γ * M * ez
			}
		} else {
			for m := 0; m < o.ndim; m++ {
				o.ez[m] = r5.e0*o.z[0][m] + r5.e1*o.z[1][m] + r5.e2*o.z[2][m]
				o.rhs[m] = o.f0[m] + γ*o.ez[m]
			}
		}

		// HW-VII p123 Eq.(8.19)
		if o.LerrStrat == 2 {
			o.lsolR.Solve(o.lerr, o.rhs, false)
			rerr = o.rms_norm(o.lerr)

			// HW-VII p123 Eq.(8.20)
		} else {
			o.lsolR.Solve(o.lerr, o.rhs, false)
			rerr = o.rms_norm(o.lerr)
			if !(rerr < 1.0) {
				if o.first || o.reject {
					for m := 0; m < o.ndim; m++ {
						o.v[0][m] = y0[m] + o.lerr[m] // y0perr
					}
					o.Nfeval++
					err = o.fcn(o.f[0], o.h, x0, o.v[0]) // f0perr
					if err != nil {
						return
					}
					if o.hasM {
						o.rhs.Apply(1, o.f[0]) // rhs := f0perr
						if o.Distr {
							la.SpMatVecMul(o.dw[0], γ, o.mMat, o.ez) // δw[0] = γ * M * ez (δw[0] is workspace)
							//o.comm.AllReduceSumAdd(o.rhs, o.dw[0], o.dw[1]) // rhs += join_with_sum(δw[0]) (δw[1] is workspace)
						} else {
							la.SpMatVecMulAdd(o.rhs, γ, o.mMat, o.ez) // rhs += γ * M * ez
						}
					} else {
						la.VecAdd(o.rhs, 1, o.f[0], γ, o.ez) // rhs = f0perr + γ * ez
					}
					o.lsolR.Solve(o.lerr, o.rhs, false)
					rerr = o.rms_norm(o.lerr)
				}
			}
		}
	}
	return
}
