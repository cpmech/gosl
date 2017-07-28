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
)

// Radau5 implements the Radau5 implicit Runge-Kutta method
type Radau5 struct {
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

// Init initialises structure
func (o *Radau5) Init(distr bool) (err error) {
	o.C = []float64{(4.0 - math.Sqrt(6.0)) / 10.0, (4.0 + math.Sqrt(6.0)) / 10.0, 1.0}

	o.T = [][]float64{{9.1232394870892942792e-02, -0.14125529502095420843, -3.0029194105147424492e-02},
		{0.24171793270710701896, 0.20412935229379993199, 0.38294211275726193779},
		{0.96604818261509293619, 1.0, 0.0}}

	o.Ti = [][]float64{{4.3255798900631553510, 0.33919925181580986954, 0.54177053993587487119},
		{-4.1787185915519047273, -0.32768282076106238708, 0.47662355450055045196},
		{-0.50287263494578687595, 2.5719269498556054292, -0.59603920482822492497}}

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
	return nil
}

// Nstages returns the number of stages
func (o *Radau5) Nstages() int {
	return 3
}

// Accept accepts update
func (o *Radau5) Accept(sol *Solver, y la.Vector) {
	for m := 0; m < sol.ndim; m++ {
		// update y
		y[m] += sol.z[2][m]
		// collocation polynomial values
		sol.ycol[0][m] = (sol.z[1][m] - sol.z[2][m]) / o.Mu4
		sol.ycol[1][m] = ((sol.z[0][m]-sol.z[1][m])/o.Mu5 - sol.ycol[0][m]) / o.Mu3
		sol.ycol[2][m] = sol.ycol[1][m] - ((sol.z[0][m]-sol.z[1][m])/o.Mu5-sol.z[0][m]/o.Mu1)/o.Mu2
	}
}

// Step steps update
func (o *Radau5) Step(sol *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) {

	// distributed version
	if sol.Distr {
		return o.StepMpi(sol, y0, x0)
	}

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
				err = num.Jacobian(&sol.dfdyT, func(fy, y la.Vector) (e error) {
					e = sol.fcn(fy, sol.h, x0, y)
					return
				}, y0, sol.f0, sol.w[0]) // w works here as workspace variable
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
				sol.mTri.Init(sol.ndim, sol.ndim, sol.ndim)
				for i := 0; i < sol.ndim; i++ {
					sol.mTri.Put(i, i, 1.0)
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
		var errR, errC error
		if !sol.Distr && sol.Pll {
			wg := new(sync.WaitGroup)
			wg.Add(2)
			go func() {
				errR = sol.lsolR.Solve(sol.dw[0], sol.v[0], false)
				wg.Done()
			}()
			go func() {
				sol.v12.JoinRealImag(sol.v[1], sol.v[2])
				errC = sol.lsolC.Solve(sol.dw12, sol.v12, false)
				sol.dw12.SplitRealImag(sol.dw[1], sol.dw[2])
				wg.Done()
			}()
			wg.Wait()
		} else {
			sol.v12.JoinRealImag(sol.v[1], sol.v[2])
			errR = sol.lsolR.Solve(sol.dw[0], sol.v[0], false)
			errC = sol.lsolC.Solve(sol.dw12, sol.v12, false)
			sol.dw12.SplitRealImag(sol.dw[1], sol.dw[2])
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
			la.SpMatVecMulAdd(sol.rhs, γ, sol.mMat, sol.ez) // rhs += γ * M * ez
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
						sol.rhs.Apply(1, sol.f[0])                      // rhs := f0perr
						la.SpMatVecMulAdd(sol.rhs, γ, sol.mMat, sol.ez) // rhs += γ * M * ez
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

// calc RMS norm
func (o *Solver) rms_norm(diff la.Vector) (rms float64) {
	for m := 0; m < o.ndim; m++ {
		rms += math.Pow(diff[m]/o.scal[m], 2.0)
	}
	rms = max(math.Sqrt(rms/float64(o.ndim)), 1.0e-10)
	return
}

// add method to database //////////////////////////////////////////////////////////////////////////

func init() {
	rkmDB[Radau5kind] = func() RKmethod { return new(Radau5) }
}
