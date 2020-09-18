// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"gosl/chk"
	"gosl/fun"
	"gosl/io"
	"gosl/la"
)

// NlSolver implements a solver to nonlinear systems of equations
//   References:
//    [1] G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
//        computations. M., Mir, 1980, p.180 of the Russian edition
type NlSolver struct {

	// configuration
	config *NlSolverConfig // configuration parameters

	// auxiliary data
	neq  int       // number of equations
	scal la.Vector // scaling vector
	fx   la.Vector // f(x)
	mdx  la.Vector // - delta x

	// functions
	functionF       fun.Vv // f(x) function f:vector, x:vector
	functionJsparse fun.Tv // J(x)=dfdx Jacobian for sparse solver J is T:triplet, x:vector
	functionJdense  fun.Mv // [non-recommended] J(x)=dfdx Jacobian for dense solver J is M:matrix, x:vector

	// data for line-search
	phi    float64
	dphidx la.Vector
	x0     la.Vector

	// data for Umfpack (sparse)
	tripletJ la.Triplet // triplet
	linsol   la.Umfpack // linear solver
	lsReady  bool       // linear solver is ready

	// workspace for numerical Jacobian (sparse)
	workspaceNumJac la.Vector // workspace

	// data for dense solver (matrix inversion)
	matrixJ    *la.Matrix // dense Jacobian matrix
	matrixJinv *la.Matrix // inverse of Jacobian matrix

	// stats
	Niter  int // number of iterations from the last call to Solve
	Nfeval int // number of calls to Ffcn (function evaluations)
	Njeval int // number of calls to Jfcn (Jacobian evaluations)
}

// NewNlSolver creates a new NlSolver
// F is the f(x) function f:vector, x:vector
// Will use numerical Jacobian (with sparse solver) by default
func NewNlSolver(neq int, F fun.Vv) (o *NlSolver) {

	// default configuration
	o = new(NlSolver)
	o.config = NewNlSolverConfig()

	// auxiliary data
	o.neq = neq
	o.scal = la.NewVector(o.neq)
	o.fx = la.NewVector(o.neq)
	o.mdx = la.NewVector(o.neq)

	// functions
	o.functionF = F

	// data for line search
	o.dphidx = la.NewVector(o.neq)
	o.x0 = la.NewVector(o.neq)
	return
}

// SetJacobianFunction sets function to compute the Jacobian (dense or sparse)
// One of sparse [recommended] or dense must be given.
// If both sparse and dense functions are given, the sparse will be used.
// With Jdense, matrix inversion is used (not very efficient. use for small systems)
func (o *NlSolver) SetJacobianFunction(Jsparse fun.Tv, Jdense fun.Mv) {
	if Jsparse == nil && Jdense == nil {
		chk.Panic("one of sparse or dense versions must be given")
	}
	if Jsparse != nil {
		o.functionJsparse = Jsparse
		o.tripletJ.Init(o.neq, o.neq, o.neq*o.neq)
	} else {
		o.config.useDenseSolver = true
		o.functionJdense = Jdense
		o.matrixJ = la.NewMatrix(o.neq, o.neq)
		o.matrixJinv = la.NewMatrix(o.neq, o.neq)
	}
	o.config.hasJacobianFunction = true
}

// Free frees memory
func (o *NlSolver) Free() {
	if !o.config.useDenseSolver {
		o.linsol.Free()
	}
}

// Solve solves non-linear problem f(x) == 0
// x -- trial x "near" the solution; otherwise it may not converge
func (o *NlSolver) Solve(x []float64) {

	// allocate workspace for numerical Jacobian
	if !o.config.hasJacobianFunction {
		if len(o.workspaceNumJac) != o.neq {
			o.workspaceNumJac = la.NewVector(o.neq)
		}
	}

	// compute scaling vector
	la.VecScaleAbs(o.scal, o.config.atol, o.config.rtol, x) // scal = Atol + Rtol*abs(x)

	// evaluate function @ x
	o.functionF(o.fx, x) // fx := f(x)
	o.Nfeval, o.Njeval = 1, 0

	// show message
	if o.config.Verbose {
		o.msg("", 0, 0, 0, true, false)
	}

	// iterations
	var Ldx, LdxPrev, Θ float64 // RMS norm of delta x, convergence rate
	var fxMax float64
	var nfv int
	for o.Niter = 0; o.Niter < o.config.MaxIterations; o.Niter++ {

		// check convergence on f(x)
		fxMax = o.fx.Largest(1.0) // den = 1.0
		if fxMax < o.config.ftol {
			if o.config.Verbose {
				o.msg("fxMax(ini)", o.Niter, Ldx, fxMax, false, true)
			}
			break
		}

		// show message
		if o.config.Verbose {
			o.msg("", o.Niter, Ldx, fxMax, false, false)
		}

		// output
		if o.config.OutCallback != nil {
			o.config.OutCallback(x)
		}

		// evaluate Jacobian @ x
		if o.Niter == 0 || !o.config.ConstantJacobian {
			if o.config.useDenseSolver {
				o.functionJdense(o.matrixJ, x)
			} else {
				if o.config.hasJacobianFunction {
					o.functionJsparse(&o.tripletJ, x)
				} else {
					Jacobian(&o.tripletJ, o.functionF, x, o.fx, o.workspaceNumJac)
					o.Nfeval += o.neq
				}
			}
			o.Njeval++
		}

		// dense solution
		if o.config.useDenseSolver {

			// invert matrix
			la.MatInv(o.matrixJinv, o.matrixJ, false)

			// solve linear system (compute mdx) and compute lin-search data
			o.phi = 0.0
			for i := 0; i < o.neq; i++ {
				o.mdx[i], o.dphidx[i] = 0.0, 0.0
				for j := 0; j < o.neq; j++ {
					o.mdx[i] += o.matrixJinv.Get(i, j) * o.fx[j] // mdx  = inv(J) * fx
					o.dphidx[i] += o.matrixJ.Get(j, i) * o.fx[j] // dφdx = tra(J) * fx
				}
				o.phi += o.fx[i] * o.fx[i]
			}
			o.phi *= 0.5

			// sparse solution
		} else {

			// init sparse solver
			if !o.lsReady {
				o.linsol.Init(&o.tripletJ, o.config.LinSolConfig)
				o.lsReady = true
			}

			// factorisation (must be done for all iterations)
			o.linsol.Fact()

			// solve linear system => compute mdx
			o.linsol.Solve(o.mdx, o.fx, false) // mdx = inv(J) * fx   false => !sumToRoot

			// compute lin-search data
			if o.config.LineSearch {
				o.phi = 0.5 * la.VecDot(o.fx, o.fx)
				la.SpTriMatTrVecMul(o.dphidx, &o.tripletJ, o.fx) // dφdx := transpose(J) * fx
			}
		}

		// update x
		Ldx = 0.0
		for i := 0; i < o.neq; i++ {
			o.x0[i] = x[i]
			x[i] -= o.mdx[i]
			Ldx += (o.mdx[i] / o.scal[i]) * (o.mdx[i] / o.scal[i])
		}
		Ldx = math.Sqrt(Ldx / float64(o.neq))

		// calculate fx := f(x) @ update x
		o.functionF(o.fx, x)
		o.Nfeval++

		// check convergence on f(x) => avoid line-search if converged already
		fxMax = o.fx.Largest(1.0) // den = 1.0
		if fxMax < o.config.ftol {
			if o.config.Verbose {
				o.msg("fxMax", o.Niter, Ldx, fxMax, false, true)
			}
			break
		}

		// check convergence on Ldx
		if Ldx < o.config.fnewt {
			if o.config.Verbose {
				o.msg("Ldx", o.Niter, Ldx, fxMax, false, true)
			}
			break
		}

		// call line-search => update x and fx
		if o.config.LineSearch {
			nfv = LineSearch(x, o.fx, o.functionF, o.mdx, o.x0, o.dphidx, o.phi, o.config.LineSearchMaxIt, true)
			o.Nfeval += nfv
			Ldx = 0.0
			for i := 0; i < o.neq; i++ {
				Ldx += ((x[i] - o.x0[i]) / o.scal[i]) * ((x[i] - o.x0[i]) / o.scal[i])
			}
			Ldx = math.Sqrt(Ldx / float64(o.neq))
			fxMax = o.fx.Largest(1.0) // den = 1.0
			if Ldx < o.config.fnewt {
				if o.config.Verbose {
					o.msg("Ldx(linsrch)", o.Niter, Ldx, fxMax, false, true)
				}
				break
			}
		}

		// check convergence rate
		if o.Niter > 0 && o.config.EnforceConvRate {
			Θ = Ldx / LdxPrev
			if Θ > 0.99 {
				chk.Panic("solver is diverging with Θ = %g (Ldx=%g, LdxPrev=%g)", Θ, Ldx, LdxPrev)
			}
		}
		LdxPrev = Ldx
	}

	// output
	if o.config.OutCallback != nil {
		o.config.OutCallback(x)
	}

	// check convergence
	if o.Niter == o.config.MaxIterations {
		chk.Panic("cannot converge after %d iterations", o.Niter)
	}
	return
}

// CheckJ check Jacobian matrix
//  Ouptut: cnd -- condition number (with Frobenius norm)
func (o *NlSolver) CheckJ(x []float64, tol float64, verbose bool) (cnd float64) {

	// Jacobian matrix
	var Jmat *la.Matrix
	if o.config.useDenseSolver {
		Jmat = la.NewMatrix(o.neq, o.neq)
		o.functionJdense(Jmat, x)
	} else {
		if o.config.hasJacobianFunction {
			o.functionJsparse(&o.tripletJ, x)
		} else {
			work := la.NewVector(o.neq)
			Jacobian(&o.tripletJ, o.functionF, x, o.fx, work)
		}
		Jmat = o.tripletJ.ToDense()
	}

	// condition number
	cnd = la.MatCondNum(Jmat, "F")
	if math.IsInf(cnd, 0) || math.IsNaN(cnd) {
		chk.Panic("condition number is Inf or NaN: %v", cnd)
	}

	// numerical Jacobian
	var Jtmp la.Triplet
	ws := la.NewVector(o.neq)
	o.functionF(o.fx, x)
	Jtmp.Init(o.neq, o.neq, o.neq*o.neq)
	Jacobian(&Jtmp, o.functionF, x, o.fx, ws)
	Jnum := Jtmp.ToMatrix(nil).ToDense()
	for i := 0; i < o.neq; i++ {
		for j := 0; j < o.neq; j++ {
			chk.PrintAnaNum(io.Sf("J[%d][%d]", i, j), tol, Jmat.Get(i, j), Jnum.Get(i, j), verbose)
		}
	}
	maxdiff := Jmat.MaxDiff(Jnum)
	if maxdiff > tol {
		chk.Panic("maxdiff = %g\n", maxdiff)
	}
	return
}

// msg prints information on residuals
func (o *NlSolver) msg(typ string, it int, Ldx, fxMax float64, first, last bool) {
	if first {
		io.Pf("\n%4s%23s%23s\n", "it", "Ldx", "fxMax")
		io.Pf("%4s%23s%23s\n", "", io.Sf("(%7.1e)", o.config.fnewt), io.Sf("(%7.1e)", o.config.ftol))
		return
	}
	io.Pf("%4d%23.15e%23.15e\n", it, Ldx, fxMax)
	if last {
		io.Pf(". . . converged with %s. nit=%d, nFeval=%d, nJeval=%d\n", typ, it, o.Nfeval, o.Njeval)
	}
}
