// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
    "math"
    "code.google.com/p/gosl/la"
    "code.google.com/p/gosl/utl"
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
    scal  []float64 // scaling vector
    fx    []float64 // f(x)
    mdx   []float64 // - delta x
    useDn bool      // use dense solver (matrix inversion) instead of Umfpack (sparse)
    numJ  bool      // use numerical Jacobian (with sparse solver)

    // callbacks
    Ffcn   Cb_f   // f(x) function
    JfcnSp Cb_J   // J(x)=dfdx Jacobian for sparse solver
    JfcnDn Cb_Jd  // J(x)=dfdx Jacobian for dense solver
    Out    Cb_out // for output

    // data for Umfpack (sparse)
    Jtri la.Triplet // triplet
    w    []float64  // workspace
    lis  la.LinSol  // linear solver

    // data for dense solver (matrix inversion)
    J  [][]float64 // dense Jacobian matrix
    Ji [][]float64 // inverse of Jacobian matrix

    // data for line-search
    φ    float64
    dφdx []float64
    x0   []float64

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
func (o *NlSolver) Init(neq int, Ffcn Cb_f, JfcnSp Cb_J, JfcnDn Cb_Jd, useDn, numJ bool, prms map[string]float64) {

    // set default values
    atol, rtol, ftol := 1e-8, 1e-8, 1e-9
    o.LsMaxIt = 20
    o.MaxIt   = 20
    o.ChkConv = true

    // read parameters
    for k, v := range prms {
        switch k {
        case "atol":    atol      = v
        case "rtol":    rtol      = v
        case "ftol":    ftol      = v
        case "lSearch": o.Lsearch = v > 0.0
        case "lsMaxIt": o.LsMaxIt = int(v)
        case "maxIt":   o.MaxIt   = int(v)
        }
    }

    // set tolerances
    o.SetTols(atol, rtol, ftol, EPS)

    // auxiliary data
    o.neq  = neq
    o.scal = make([]float64, o.neq)
    o.fx   = make([]float64, o.neq)
    o.mdx  = make([]float64, o.neq)

    // callbacks
    o.Ffcn, o.JfcnSp, o.JfcnDn = Ffcn, JfcnSp, JfcnDn

    // type of linear solver and Jacobian matrix (numerical or analytical: sparse only)
    o.useDn, o.numJ = useDn, numJ

    // use dense linear solver
    if o.useDn {
        o.J  = la.MatAlloc(o.neq, o.neq)
        o.Ji = la.MatAlloc(o.neq, o.neq)

    // use sparse linear solver
    } else {
        o.Jtri.Init(o.neq, o.neq, o.neq*o.neq)
        if JfcnSp == nil {
            o.numJ = true
        }
        if o.numJ {
            o.w = make([]float64, o.neq)
        }
        o.lis = la.GetSolver("umfpack")
    }

    // allocate slices for line search
    o.dφdx = make([]float64, o.neq)
    o.x0   = make([]float64, o.neq)
}

// Clean performs clean ups
func (o *NlSolver) Clean() {
    if !o.useDn {
        o.lis.Clean()
    }
}

// SetTols set tolerances
func (o *NlSolver) SetTols(Atol, Rtol, Ftol, ϵ float64) {
    o.atol, o.rtol, o.ftol = Atol, Rtol, Ftol
    o.fnewt = max(10.0 * ϵ / Rtol, min(0.03, math.Sqrt(Rtol)))
}

// Solve solves non-linear problem f(x) == 0
func (o *NlSolver) Solve(x []float64, silent bool) (err error) {

    // compute scaling vector
    la.VecScaleAbs(o.scal, o.atol, o.rtol, x) // scal := Atol + Rtol*abs(x)

    // evaluate function @ x
    err = o.Ffcn(o.fx, x) // fx := f(x)
    o.NFeval, o.NJeval = 1, 0
    if err != nil {
        return
    }

    // show message
    if !silent{
        o.msg("", 0, 0, 0, true, false)
    }

    // iterations
    var Ldx, Ldx_prev, Θ float64 // RMS norm of delta x, convergence rate
    var fx_max float64
    var nfv int
    for o.It = 0; o.It < o.MaxIt; o.It++ {

        // check convergence on f(x)
        fx_max = la.VecLargest(o.fx, 1.0) // den = 1.0
        if fx_max < o.ftol {
            if !silent {
                o.msg("fx_max(ini)", o.It, Ldx, fx_max, false, true)
            }
            break
        }

        // show message
        if !silent{
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
                    err = Jacobian(&o.Jtri, o.Ffcn, x, o.fx, o.w, false)
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
            err = la.MatInvG(o.Ji, o.J, 1e-10)
            if err != nil {
                return utl.Err(_nls_err1, err.Error())
            }

            // solve linear system (compute mdx) and compute lin-search data
            o.φ = 0.0
            for i := 0; i < o.neq; i++ {
                o.mdx[i], o.dφdx[i] = 0.0, 0.0
                for j := 0; j < o.neq; j++ {
                    o.mdx [i] += o.Ji[i][j] * o.fx[j] // mdx  = inv(J) * fx
                    o.dφdx[i] += o.J [j][i] * o.fx[j] // dφdx = tra(J) * fx
                }
                o.φ += o.fx[i] * o.fx[i]
            }
            o.φ *= 0.5

        // sparse solution
        } else {

            // init sparse solver
            if o.It == 0 {
                symmetric, verbose, timing := false, false, false
                err := o.lis.InitR(&o.Jtri, symmetric, verbose, timing)
                if err != nil {
                    return utl.Err(_nls_err9, err.Error())
                }
            }

            // factorisation (must be done for all iterations)
            o.lis.Fact()

            // solve linear system => compute mdx
            o.lis.SolveR(o.mdx, o.fx, false) // mdx = inv(J) * fx   false => !sumToRoot

            // compute lin-search data
            if o.Lsearch {
                o.φ = 0.5 * la.VecDot(o.fx, o.fx)
                la.SpTriMatTrVecMul(o.dφdx, &o.Jtri, o.fx) // dφdx := transpose(J) * fx
            }
        }

        //utl.Pforan("φ    = %v\n", o.φ)
        //utl.Pforan("dφdx = %v\n", o.dφdx)

        // update x
        Ldx = 0.0
        for i := 0; i < o.neq; i++ {
            o.x0[i] = x[i]
            x[i]   -= o.mdx[i]
            Ldx    += (o.mdx[i]/o.scal[i]) * (o.mdx[i]/o.scal[i])
        }
        Ldx = math.Sqrt(Ldx / float64(o.neq))

        // calculate fx := f(x) @ update x
        err = o.Ffcn(o.fx, x)
        o.NFeval += 1
        if err != nil {
            return
        }

        // check convergence on f(x) => avoid line-search if converged already
        fx_max = la.VecLargest(o.fx, 1.0) // den = 1.0
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
                return utl.Err(_nls_err2, err.Error())
            }
            Ldx = 0.0
            for i := 0; i < o.neq; i++ {
                Ldx += ((x[i]-o.x0[i])/o.scal[i]) * ((x[i]-o.x0[i])/o.scal[i])
            }
            Ldx = math.Sqrt(Ldx / float64(o.neq))
            fx_max = la.VecLargest(o.fx, 1.0) // den = 1.0
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
                return utl.Err(_nls_err3, Θ, Ldx, Ldx_prev)
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
        err = utl.Err(_nls_err4, o.It)
    }
    return
}

// CheckJ check Jacobian matrix
//  Ouptut: cnd -- condition number (with Frobenius norm)
func (o *NlSolver) CheckJ(x []float64, tol float64, chkJnum, silent bool) (cnd float64, err error) {

    // Jacobian matrix
    var Jmat [][]float64
    if o.useDn {
        Jmat = la.MatAlloc(o.neq, o.neq)
        err  = o.JfcnDn(Jmat, x)
        if err != nil {
            return 0, utl.Err(_nls_err5, "dense", err.Error())
        }
    } else {
        if o.numJ {
            err = Jacobian(&o.Jtri, o.Ffcn, x, o.fx, o.w, false)
            if err != nil {
                return 0, utl.Err(_nls_err5, "sparse", err.Error())
            }
        } else {
            err = o.JfcnSp(&o.Jtri, x)
            if err != nil {
                return 0, utl.Err(_nls_err5, "sparse(num)", err.Error())
            }
        }
        Jmat = o.Jtri.ToMatrix(nil).ToDense()
    }
    //la.PrintMat("J", Jmat, "%23g", false)

    // condition number
    cnd, err = la.MatCondG(Jmat, "F", 1e-10)
    if err != nil {
        return cnd, utl.Err(_nls_err6, err.Error())
    }
    if math.IsInf(cnd,0) || math.IsNaN(cnd) {
        return cnd, utl.Err(_nls_err7, cnd)
    }

    // numerical Jacobian
    if !chkJnum {
        return
    }
    var Jtmp la.Triplet
    ws := make([]float64, o.neq)
    err = o.Ffcn(o.fx, x)
    if err != nil {
        return
    }
    Jtmp.Init(o.neq, o.neq, o.neq*o.neq)
    Jacobian(&Jtmp, o.Ffcn, x, o.fx, ws, false)
    Jnum := Jtmp.ToMatrix(nil).ToDense()
    for i := 0; i < o.neq; i++ {
        for j := 0; j < o.neq; j++ {
            utl.AnaNum(utl.Sf("J[%d][%d]",i,j), tol, Jmat[i][j], Jnum[i][j], !silent)
        }
    }
    maxdiff := la.MatMaxDiff(Jmat, Jnum)
    if maxdiff > tol {
        err = utl.Err(_nls_err8, maxdiff)
    }
    return
}

// msg prints information on residuals
func (o *NlSolver) msg(typ string, it int, Ldx, fx_max float64, first, last bool) {
    if first {
        utl.Pfpink("\n%4s%23s%23s\n", "it", "Ldx", "fx_max")
        utl.Pfpink("%4s%23s%23s\n", "", utl.Sf("(%7.1e)", o.fnewt), utl.Sf("(%7.1e)", o.ftol))
        return
    }
    utl.Pfyel("%4d%23.15e%23.15e\n", it, Ldx, fx_max)
    if last {
        utl.Pfgrey(". . . converged with %s. nit=%d, nFeval=%d, nJeval=%d\n", typ, it, o.NFeval, o.NJeval)
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
