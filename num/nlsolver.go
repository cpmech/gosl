// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

type NlSolver struct {
	// constants
	CteJac  bool    // constant Jacobian (Modified Newton's method)
	Lsearch bool    // use linear search
	LsMaxIt int     // linear solver maximum iterations
	MaxIt   int     // Newton's method maximum iterations
	ChkConv bool    // check convergence
	atol    float64 // absolute tolerance
	rtol    float64 // relative tolerance
	ftol    float64 // minimum value of fx
	fnewt   float64 // Newton's method tolerance

	// auxiliary data
	neq   int       // number of equations
	scal  la.Vector // scaling vector
	fx    la.Vector // f(x)
	mdx   la.Vector // - delta x
	useDn bool      // use dense solver (matrix inversion) instead of Umfpack (sparse)
	numJ  bool      // use numerical Jacobian (with sparse solver)

	// callbacks
	Ffcn   fun.Vv // f(x) function f:vector, x:vector
	JfcnSp fun.Tv // J(x)=dfdx Jacobian for sparse solver
	JfcnDn fun.Mv // J(x)=dfdx Jacobian for dense solver

	// output callback
	Out func(x []float64) error // output callback function

	// data for Umfpack (sparse)
	Jtri la.Triplet // triplet
	w    la.Vector  // workspace
	lis  la.Umfpack // linear solver

	// data for dense solver (matrix inversion)
	J  *la.Matrix // dense Jacobian matrix
	Ji *la.Matrix // inverse of Jacobian matrix

	// data for line-search
	φ    float64
	dφdx la.Vector
	x0   la.Vector

	// stat data
	It     int // number of iterations from the last call to Solve
	NFeval int // number of calls to Ffcn (function evaluations)
	NJeval int // number of calls to Jfcn (Jacobian evaluations)
}

// Init initialises solver
//  Input:
//   useSp -- Use sparse solver with JfcnSp
//   useDn -- Use dense solver (matrix inversion) with JfcnDn
//   numJ  -- Use numeric Jacobian (sparse version only)
//   prms  -- atol, rtol, ftol, lSearch, lsMaxIt, maxIt
func (o *NlSolver) Init(neq int, Ffcn fun.Vv, JfcnSp fun.Tv, JfcnDn fun.Mv, useDn, numJ bool, prms map[string]float64) {

	// set default values
	atol, rtol, ftol := 1e-8, 1e-8, 1e-9
	o.LsMaxIt = 20
	o.MaxIt = 20
	o.ChkConv = true

	// read parameters
	for k, v := range prms {
		switch k {
		case "atol":
			atol = v
		case "rtol":
			rtol = v
		case "ftol":
			ftol = v
		case "lSearch":
			o.Lsearch = v > 0.0
		case "lsMaxIt":
			o.LsMaxIt = int(v)
		case "maxIt":
			o.MaxIt = int(v)
		}
	}

	// set tolerances
	o.SetTols(atol, rtol, ftol, MACHEPS)

	// auxiliary data
	o.neq = neq
	o.scal = la.NewVector(o.neq)
	o.fx = la.NewVector(o.neq)
	o.mdx = la.NewVector(o.neq)

	// callbacks
	o.Ffcn, o.JfcnSp, o.JfcnDn = Ffcn, JfcnSp, JfcnDn

	// type of linear solver and Jacobian matrix (numerical or analytical: sparse only)
	o.useDn, o.numJ = useDn, numJ

	// use dense linear solver
	if o.useDn {
		o.J = la.NewMatrix(o.neq, o.neq)
		o.Ji = la.NewMatrix(o.neq, o.neq)

		// use sparse linear solver
	} else {
		o.Jtri.Init(o.neq, o.neq, o.neq*o.neq)
		if JfcnSp == nil {
			o.numJ = true
		}
		if o.numJ {
			o.w = la.NewVector(o.neq)
		}
	}

	// allocate slices for line search
	o.dφdx = la.NewVector(o.neq)
	o.x0 = la.NewVector(o.neq)
}

// Free frees memory
func (o *NlSolver) Free() {
	if !o.useDn {
		o.lis.Free()
	}
}

// SetTols set tolerances
func (o *NlSolver) SetTols(Atol, Rtol, Ftol, ϵ float64) {
	o.atol, o.rtol, o.ftol = Atol, Rtol, Ftol
	o.fnewt = max(10.0*ϵ/Rtol, min(0.03, math.Sqrt(Rtol)))
}

// Solve solves non-linear problem f(x) == 0
func (o *NlSolver) Solve(x []float64, silent bool) (err error) {

	// compute scaling vector
	la.VecScaleAbs(o.scal, o.atol, o.rtol, x) // scal = Atol + Rtol*abs(x)

	// evaluate function @ x
	err = o.Ffcn(o.fx, x) // fx := f(x)
	o.NFeval, o.NJeval = 1, 0
	if err != nil {
		return
	}

	// show message
	if !silent {
		o.msg("", 0, 0, 0, true, false)
	}

	// iterations
	var Ldx, Ldx_prev, Θ float64 // RMS norm of delta x, convergence rate
	var fx_max float64
	var nfv int
	for o.It = 0; o.It < o.MaxIt; o.It++ {

		// check convergence on f(x)
		fx_max = o.fx.Largest(1.0) // den = 1.0
		if fx_max < o.ftol {
			if !silent {
				o.msg("fx_max(ini)", o.It, Ldx, fx_max, false, true)
			}
			break
		}

		// show message
		if !silent {
			o.msg("", o.It, Ldx, fx_max, false, false)
		}

		// output
		if o.Out != nil {
			o.Out(x)
		}

		// evaluate Jacobian @ x
		if o.It == 0 || !o.CteJac {
			if o.useDn {
				err = o.JfcnDn(o.J, x)
			} else {
				if o.numJ {
					err = Jacobian(&o.Jtri, o.Ffcn, x, o.fx, o.w)
					o.NFeval += o.neq
				} else {
					err = o.JfcnSp(&o.Jtri, x)
				}
			}
			o.NJeval += 1
			if err != nil {
				return
			}
		}

		// dense solution
		if o.useDn {

			// invert matrix
			err = la.MatInv(o.Ji, o.J)
			if err != nil {
				return chk.Err(_nls_err1, err.Error())
			}

			// solve linear system (compute mdx) and compute lin-search data
			o.φ = 0.0
			for i := 0; i < o.neq; i++ {
				o.mdx[i], o.dφdx[i] = 0.0, 0.0
				for j := 0; j < o.neq; j++ {
					o.mdx[i] += o.Ji.Get(i, j) * o.fx[j] // mdx  = inv(J) * fx
					o.dφdx[i] += o.J.Get(j, i) * o.fx[j] // dφdx = tra(J) * fx
				}
				o.φ += o.fx[i] * o.fx[i]
			}
			o.φ *= 0.5

			// sparse solution
		} else {

			// init sparse solver
			if o.It == 0 {
				symmetric, verbose := false, false
				err := o.lis.Init(&o.Jtri, symmetric, verbose, "", "", nil)
				if err != nil {
					return chk.Err(_nls_err9, err.Error())
				}
			}

			// factorisation (must be done for all iterations)
			o.lis.Fact()

			// solve linear system => compute mdx
			o.lis.Solve(o.mdx, o.fx, false) // mdx = inv(J) * fx   false => !sumToRoot

			// compute lin-search data
			if o.Lsearch {
				o.φ = 0.5 * la.VecDot(o.fx, o.fx)
				la.SpTriMatTrVecMul(o.dφdx, &o.Jtri, o.fx) // dφdx := transpose(J) * fx
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
		err = o.Ffcn(o.fx, x)
		o.NFeval += 1
		if err != nil {
			return
		}

		// check convergence on f(x) => avoid line-search if converged already
		fx_max = o.fx.Largest(1.0) // den = 1.0
		if fx_max < o.ftol {
			if !silent {
				o.msg("fx_max", o.It, Ldx, fx_max, false, true)
			}
			break
		}

		// check convergence on Ldx
		if Ldx < o.fnewt {
			if !silent {
				o.msg("Ldx", o.It, Ldx, fx_max, false, true)
			}
			break
		}

		// call line-search => update x and fx
		if o.Lsearch {
			nfv, err = LineSearch(x, o.fx, o.Ffcn, o.mdx, o.x0, o.dφdx, o.φ, o.LsMaxIt, true)
			o.NFeval += nfv
			if err != nil {
				return chk.Err(_nls_err2, err.Error())
			}
			Ldx = 0.0
			for i := 0; i < o.neq; i++ {
				Ldx += ((x[i] - o.x0[i]) / o.scal[i]) * ((x[i] - o.x0[i]) / o.scal[i])
			}
			Ldx = math.Sqrt(Ldx / float64(o.neq))
			fx_max = o.fx.Largest(1.0) // den = 1.0
			if Ldx < o.fnewt {
				if !silent {
					o.msg("Ldx(linsrch)", o.It, Ldx, fx_max, false, true)
				}
				break
			}
		}

		// check convergence rate
		if o.It > 0 && o.ChkConv {
			Θ = Ldx / Ldx_prev
			if Θ > 0.99 {
				return chk.Err(_nls_err3, Θ, Ldx, Ldx_prev)
			}
		}
		Ldx_prev = Ldx
	}

	// output
	if o.Out != nil {
		o.Out(x)
	}

	// check convergence
	if o.It == o.MaxIt {
		err = chk.Err(_nls_err4, o.It)
	}
	return
}

// CheckJ check Jacobian matrix
//  Ouptut: cnd -- condition number (with Frobenius norm)
func (o *NlSolver) CheckJ(x []float64, tol float64, chkJnum, silent bool) (cnd float64, err error) {

	// Jacobian matrix
	var Jmat *la.Matrix
	if o.useDn {
		Jmat = la.NewMatrix(o.neq, o.neq)
		err = o.JfcnDn(Jmat, x)
		if err != nil {
			return 0, chk.Err(_nls_err5, "dense", err.Error())
		}
	} else {
		if o.numJ {
			err = Jacobian(&o.Jtri, o.Ffcn, x, o.fx, o.w)
			if err != nil {
				return 0, chk.Err(_nls_err5, "sparse", err.Error())
			}
		} else {
			err = o.JfcnSp(&o.Jtri, x)
			if err != nil {
				return 0, chk.Err(_nls_err5, "sparse(num)", err.Error())
			}
		}
		Jmat = o.Jtri.GetDenseMatrix()
	}

	// condition number
	cnd, err = la.MatCondNum(Jmat, "F")
	if err != nil {
		return cnd, chk.Err(_nls_err6, err.Error())
	}
	if math.IsInf(cnd, 0) || math.IsNaN(cnd) {
		return cnd, chk.Err(_nls_err7, cnd)
	}

	// numerical Jacobian
	if !chkJnum {
		return
	}
	var Jtmp la.Triplet
	ws := la.NewVector(o.neq)
	err = o.Ffcn(o.fx, x)
	if err != nil {
		return
	}
	Jtmp.Init(o.neq, o.neq, o.neq*o.neq)
	Jacobian(&Jtmp, o.Ffcn, x, o.fx, ws)
	Jnum := Jtmp.ToMatrix(nil).ToDense()
	for i := 0; i < o.neq; i++ {
		for j := 0; j < o.neq; j++ {
			chk.PrintAnaNum(io.Sf("J[%d][%d]", i, j), tol, Jmat.Get(i, j), Jnum.Get(i, j), !silent)
		}
	}
	maxdiff := Jmat.MaxDiff(Jnum)
	if maxdiff > tol {
		err = chk.Err(_nls_err8, maxdiff)
	}
	return
}

// msg prints information on residuals
func (o *NlSolver) msg(typ string, it int, Ldx, fx_max float64, first, last bool) {
	if first {
		io.Pf("\n%4s%23s%23s\n", "it", "Ldx", "fx_max")
		io.Pf("%4s%23s%23s\n", "", io.Sf("(%7.1e)", o.fnewt), io.Sf("(%7.1e)", o.ftol))
		return
	}
	io.Pf("%4d%23.15e%23.15e\n", it, Ldx, fx_max)
	if last {
		io.Pf(". . . converged with %s. nit=%d, nFeval=%d, nJeval=%d\n", typ, it, o.NFeval, o.NJeval)
	}
}

// error messages
var (
	_nls_err1 = "nlsolver.go: NlSolver.Solve failed: cannot compute inverse of Jacobian (dense) matrix:\n%v"
	_nls_err2 = "nlsolver.go: NlSolver.Solve: LineSearch failed:\n%v"
	_nls_err3 = "nlsolver.go: NlSolver.Solve is diverging with Θ = %g (Ldx=%g, Ldx_prev=%g)"
	_nls_err4 = "nlsolver.go: NlSolver.Solve did not converge after %d iterations"
	_nls_err5 = "nlsolver.go: NlSolver.CheckJ: %s: failed:\n%v"
	_nls_err6 = "nlsolver.go: NlSolver.CheckJ failed: cannot compute condition number\n%v"
	_nls_err7 = "nlsolver.go: NlSolver.CheckJ failed: condition number is Inf or NaN: %v"
	_nls_err8 = "nlsolver.go: NlSolver.CheckJ failed: maxdiff = %g"
	_nls_err9 = "nlsolver.go: NlSolver.Init: cannot initialise LinSol('umfpack'):\n%v\n"
)
