// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package ode

import (
	"errors"
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
)

// StepMpi steps update using MPI
func (o *Radau5) StepMpi(sol *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) {

	// factors
	α := o.Alp / sol.h
	β := o.Bet / sol.h
	γ := o.Gam / sol.h

	// Jacobian and decomposition
	if sol.reuseJdec {
		sol.reuseJdec = false
	} else {

		// calculate only first Jacobian for all iterations (simple/modified Newton's method)
		if sol.reuseJ {
			sol.reuseJ = false
		} else if !sol.jacIsOK {

			// Jacobian triplet
			if sol.jac == nil { // numerical
				//if x0 == 0.0 { io.Pfgrey(" > > > > > > > > . . . numerical Jacobian . . . < < < < < < < < <\n") }
				err = num.JacobianMpi(sol.comm, &sol.dfdyT, func(fy, y la.Vector) (e error) {
					e = sol.fcn(fy, sol.h, x0, y)
					return
				}, y0, sol.f0, sol.w[0], sol.Distr) // w works here as workspace variable
			} else { // analytical
				//if x0 == 0.0 { io.Pfgrey(" > > > > > > > > . . . analytical Jacobian . . . < < < < < < < < <\n") }
				err = sol.jac(&sol.dfdyT, sol.h, x0, y0)
			}
			if err != nil {
				return
			}

			// create M matrix
			if sol.doinit && !sol.hasM {
				sol.mTri = new(la.Triplet)
				if sol.Distr {
					id, sz := sol.comm.Rank(), sol.comm.Size()
					start, endp1 := (id*sol.ndim)/sz, ((id+1)*sol.ndim)/sz
					sol.mTri.Init(sol.ndim, sol.ndim, endp1-start)
					for i := start; i < endp1; i++ {
						sol.mTri.Put(i, i, 1.0)
					}
				} else {
					sol.mTri.Init(sol.ndim, sol.ndim, sol.ndim)
					for i := 0; i < sol.ndim; i++ {
						sol.mTri.Put(i, i, 1.0)
					}
				}
			}
			sol.Njeval++
			sol.jacIsOK = true
		}

		// initialise triplets
		if sol.doinit {
			sol.rctriR = new(la.Triplet)
			sol.rctriC = new(la.TripletC)
			sol.rctriR.Init(sol.ndim, sol.ndim, sol.mTri.Len()+sol.dfdyT.Len())
			sol.rctriC.Init(sol.ndim, sol.ndim, sol.mTri.Len()+sol.dfdyT.Len())
		}

		// update triplets
		la.SpTriAdd(sol.rctriR, γ, sol.mTri, -1, &sol.dfdyT)       // rctriR :=      γ*M - dfdy
		la.SpTriAddR2C(sol.rctriC, α, β, sol.mTri, -1, &sol.dfdyT) // rctriC := (α+βi)*M - dfdy

		// initialise solver
		if sol.doinit {
			err = sol.lsolR.Init(sol.rctriR, sol.symmetric, sol.lsverbose, sol.ordering, sol.scaling, sol.comm)
			if err != nil {
				return
			}
			err = sol.lsolC.Init(sol.rctriC, sol.symmetric, sol.lsverbose, sol.ordering, sol.scaling, sol.comm)
			if err != nil {
				return
			}
		}

		// perform factorisation
		sol.lsolR.Fact()
		sol.lsolC.Fact()
		sol.Ndecomp++
	}

	// updated u[i]
	sol.u[0] = x0 + o.C[0]*sol.h
	sol.u[1] = x0 + o.C[1]*sol.h
	sol.u[2] = x0 + o.C[2]*sol.h

	// (trial/initial) updated z[i] and w[i]
	if sol.first || sol.ZeroTrial {
		for m := 0; m < sol.ndim; m++ {
			sol.z[0][m], sol.w[0][m] = 0.0, 0.0
			sol.z[1][m], sol.w[1][m] = 0.0, 0.0
			sol.z[2][m], sol.w[2][m] = 0.0, 0.0
		}
	} else {
		c3q := sol.h / sol.hprev
		c1q := o.Mu1 * c3q
		c2q := o.Mu2 * c3q
		for m := 0; m < sol.ndim; m++ {
			sol.z[0][m] = c1q * (sol.ycol[0][m] + (c1q-o.Mu4)*(sol.ycol[1][m]+(c1q-o.Mu3)*sol.ycol[2][m]))
			sol.z[1][m] = c2q * (sol.ycol[0][m] + (c2q-o.Mu4)*(sol.ycol[1][m]+(c2q-o.Mu3)*sol.ycol[2][m]))
			sol.z[2][m] = c3q * (sol.ycol[0][m] + (c3q-o.Mu4)*(sol.ycol[1][m]+(c3q-o.Mu3)*sol.ycol[2][m]))
			sol.w[0][m] = o.Ti[0][0]*sol.z[0][m] + o.Ti[0][1]*sol.z[1][m] + o.Ti[0][2]*sol.z[2][m]
			sol.w[1][m] = o.Ti[1][0]*sol.z[0][m] + o.Ti[1][1]*sol.z[1][m] + o.Ti[1][2]*sol.z[2][m]
			sol.w[2][m] = o.Ti[2][0]*sol.z[0][m] + o.Ti[2][1]*sol.z[1][m] + o.Ti[2][2]*sol.z[2][m]
		}
	}

	// iterations
	sol.nit = 0
	sol.eta = math.Pow(max(sol.eta, sol.Eps), 0.8)
	sol.theta = sol.ThetaMax
	sol.diverg = false
	var Lδw, oLδw, thq, othq, iterr, itRerr, qnewt float64
	var it int
	for it = 0; it < sol.NmaxIt; it++ {

		// max iterations ?
		sol.nit = it + 1
		if sol.nit > sol.Nitmax {
			sol.Nitmax = sol.nit
		}

		// evaluate f(x,y) at (u[i],v[i]=y0+z[i])
		for i := 0; i < 3; i++ {
			for m := 0; m < sol.ndim; m++ {
				sol.v[i][m] = y0[m] + sol.z[i][m]
			}
			sol.Nfeval++
			err = sol.fcn(sol.f[i], sol.h, sol.u[i], sol.v[i])
			if err != nil {
				return
			}
		}

		// calc rhs
		if sol.hasM {
			// using δw as workspace here
			la.SpMatVecMul(sol.dw[0], 1, sol.mMat, sol.w[0]) // δw0 := M * w0
			la.SpMatVecMul(sol.dw[1], 1, sol.mMat, sol.w[1]) // δw1 := M * w1
			la.SpMatVecMul(sol.dw[2], 1, sol.mMat, sol.w[2]) // δw2 := M * w2
			if sol.Distr {
				sol.comm.AllReduceSum(sol.dw[0], sol.v[0]) // v is used as workspace here
				sol.comm.AllReduceSum(sol.dw[1], sol.v[1]) // v is used as workspace here
				sol.comm.AllReduceSum(sol.dw[2], sol.v[2]) // v is used as workspace here
			}
			for m := 0; m < sol.ndim; m++ {
				sol.v[0][m] = o.Ti[0][0]*sol.f[0][m] + o.Ti[0][1]*sol.f[1][m] + o.Ti[0][2]*sol.f[2][m] - γ*sol.dw[0][m]
				sol.v[1][m] = o.Ti[1][0]*sol.f[0][m] + o.Ti[1][1]*sol.f[1][m] + o.Ti[1][2]*sol.f[2][m] - α*sol.dw[1][m] + β*sol.dw[2][m]
				sol.v[2][m] = o.Ti[2][0]*sol.f[0][m] + o.Ti[2][1]*sol.f[1][m] + o.Ti[2][2]*sol.f[2][m] - β*sol.dw[1][m] - α*sol.dw[2][m]
			}
		} else {
			for m := 0; m < sol.ndim; m++ {
				sol.v[0][m] = o.Ti[0][0]*sol.f[0][m] + o.Ti[0][1]*sol.f[1][m] + o.Ti[0][2]*sol.f[2][m] - γ*sol.w[0][m]
				sol.v[1][m] = o.Ti[1][0]*sol.f[0][m] + o.Ti[1][1]*sol.f[1][m] + o.Ti[1][2]*sol.f[2][m] - α*sol.w[1][m] + β*sol.w[2][m]
				sol.v[2][m] = o.Ti[2][0]*sol.f[0][m] + o.Ti[2][1]*sol.f[1][m] + o.Ti[2][2]*sol.f[2][m] - β*sol.w[1][m] - α*sol.w[2][m]
			}
		}

		// solve linear system
		sol.Nlinsol++
		sol.v12.JoinRealImag(sol.v[1], sol.v[2])
		errR := sol.lsolR.Solve(sol.dw[0], sol.v[0], false)
		errC := sol.lsolC.Solve(sol.dw12, sol.v12, false)
		sol.dw12.SplitRealImag(sol.dw[1], sol.dw[2])

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
		for m := 0; m < sol.ndim; m++ {
			sol.w[0][m] += sol.dw[0][m]
			sol.w[1][m] += sol.dw[1][m]
			sol.w[2][m] += sol.dw[2][m]
			sol.z[0][m] = o.T[0][0]*sol.w[0][m] + o.T[0][1]*sol.w[1][m] + o.T[0][2]*sol.w[2][m]
			sol.z[1][m] = o.T[1][0]*sol.w[0][m] + o.T[1][1]*sol.w[1][m] + o.T[1][2]*sol.w[2][m]
			sol.z[2][m] = o.T[2][0]*sol.w[0][m] + o.T[2][1]*sol.w[1][m] + o.T[2][2]*sol.w[2][m]
		}

		// rms norm of δw
		Lδw = 0.0
		for m := 0; m < sol.ndim; m++ {
			Lδw += math.Pow(sol.dw[0][m]/sol.scal[m], 2.0) + math.Pow(sol.dw[1][m]/sol.scal[m], 2.0) + math.Pow(sol.dw[2][m]/sol.scal[m], 2.0)
		}
		Lδw = math.Sqrt(Lδw / float64(3*sol.ndim))

		// check convergence
		if it > 0 {
			thq = Lδw / oLδw
			if it == 1 {
				sol.theta = thq
			} else {
				sol.theta = math.Sqrt(thq * othq)
			}
			othq = thq
			if sol.theta < 0.99 {
				sol.eta = sol.theta / (1.0 - sol.theta)
				iterr = Lδw * math.Pow(sol.theta, float64(sol.NmaxIt-sol.nit)) / (1.0 - sol.theta)
				itRerr = iterr / sol.fnewt
				if itRerr >= 1.0 { // diverging
					qnewt = max(1.0e-4, min(20.0, itRerr))
					sol.dvfac = 0.8 * math.Pow(qnewt, -1.0/(4.0+float64(sol.NmaxIt)-1.0-float64(sol.nit)))
					sol.diverg = true
					break
				}
			} else { // diverging badly (unexpected step-rejection)
				sol.dvfac = 0.5
				sol.diverg = true
				break
			}
		}

		// save old norm
		oLδw = Lδw

		// converged
		if sol.eta*Lδw < sol.fnewt {
			break
		}
	}

	// did not converge
	if it == sol.NmaxIt-1 {
		chk.Panic("radau5_step failed with it=%d", it)
	}

	// diverging => stop
	if sol.diverg {
		rerr = 2.0 // must leave state intact, any rerr is OK
		return
	}

	// error estimate
	if sol.LerrStrat == 1 {

		// simple strategy => HW-VII p123 Eq.(8.17) (not good for stiff problems)
		for m := 0; m < sol.ndim; m++ {
			sol.ez[m] = o.E0*sol.z[0][m] + o.E1*sol.z[1][m] + o.E2*sol.z[2][m]
			sol.lerr[m] = o.Gam0*sol.h*sol.f0[m] + sol.ez[m]
			rerr += math.Pow(sol.lerr[m]/sol.scal[m], 2.0)
		}
		rerr = max(math.Sqrt(rerr/float64(sol.ndim)), 1.0e-10)

	} else {

		// common
		if sol.hasM {
			for m := 0; m < sol.ndim; m++ {
				sol.ez[m] = o.E0*sol.z[0][m] + o.E1*sol.z[1][m] + o.E2*sol.z[2][m]
				sol.rhs[m] = sol.f0[m]
			}
			if sol.Distr {
				la.SpMatVecMul(sol.dw[0], γ, sol.mMat, sol.ez) // δw[0] = γ * M * ez (δw[0] is workspace)
				//o.comm.AllReduceSumAdd(o.rhs, o.dw[0], o.dw[1]) // rhs += join_with_sum(δw[0]) (δw[1] is workspace)
				chk.Panic("stop")
			} else {
				la.SpMatVecMulAdd(sol.rhs, γ, sol.mMat, sol.ez) // rhs += γ * M * ez
			}
		} else {
			for m := 0; m < sol.ndim; m++ {
				sol.ez[m] = o.E0*sol.z[0][m] + o.E1*sol.z[1][m] + o.E2*sol.z[2][m]
				sol.rhs[m] = sol.f0[m] + γ*sol.ez[m]
			}
		}

		// HW-VII p123 Eq.(8.19)
		if sol.LerrStrat == 2 {
			sol.lsolR.Solve(sol.lerr, sol.rhs, false)
			rerr = sol.rms_norm(sol.lerr)

			// HW-VII p123 Eq.(8.20)
		} else {
			sol.lsolR.Solve(sol.lerr, sol.rhs, false)
			rerr = sol.rms_norm(sol.lerr)
			if !(rerr < 1.0) {
				if sol.first || sol.reject {
					for m := 0; m < sol.ndim; m++ {
						sol.v[0][m] = y0[m] + sol.lerr[m] // y0perr
					}
					sol.Nfeval++
					err = sol.fcn(sol.f[0], sol.h, x0, sol.v[0]) // f0perr
					if err != nil {
						return
					}
					if sol.hasM {
						sol.rhs.Apply(1, sol.f[0]) // rhs := f0perr
						if sol.Distr {
							la.SpMatVecMul(sol.dw[0], γ, sol.mMat, sol.ez) // δw[0] = γ * M * ez (δw[0] is workspace)
							//o.comm.AllReduceSumAdd(o.rhs, o.dw[0], o.dw[1]) // rhs += join_with_sum(δw[0]) (δw[1] is workspace)
							sol.comm.AllReduceSum(sol.dw[1], sol.dw[0]) // dw1 := join(dw0)
							chk.Panic("stop")
							sol.rhs.Apply(1, sol.dw[1]) // rhs += dw0
						} else {
							la.SpMatVecMulAdd(sol.rhs, γ, sol.mMat, sol.ez) // rhs += γ * M * ez
						}
					} else {
						la.VecAdd(sol.rhs, 1, sol.f[0], γ, sol.ez) // rhs = f0perr + γ * ez
					}
					sol.lsolR.Solve(sol.lerr, sol.rhs, false)
					rerr = sol.rms_norm(sol.lerr)
				}
			}
		}
	}
	return
}
