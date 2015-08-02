// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

// ReliabGfcn_t defines the limit state function
type ReliabGfcn_t func(x []float64, args ...interface{}) (float64, error)

// ReliabHfcn_t defines the gradient of the limit state function
type ReliabHfcn_t func(dgdx []float64, x []float64, args ...interface{}) error

// ReliabFORM implements the first order reliability method
type ReliabFORM struct {

	// input
	μ    []float64    // [nx] mean values of each random variable x
	σ    []float64    // [nx] deviation values of each random variable x
	lrv  []bool       // [nx] lognormal random variable?
	gfcn ReliabGfcn_t // limit state function
	hfcn ReliabHfcn_t // gradient of limit state function

	// constants
	NmaxItA   int     // max number of iterations for direction cosines
	NmaxItB   int     // max number of iterations for reliability index
	TolA      float64 // tolerance to find α
	TolB      float64 // tolerance to find β
	NlsSilent bool    // flag for nonlinear solver

	// auxiliary
	α    []float64 // [nx] direction cosines
	xtmp []float64 // [nx] temporary vector of random variables
	dgdx []float64 // [nx] gradient of g
}

func (o *ReliabFORM) Init(μ, σ []float64, lrv []bool, gfcn ReliabGfcn_t, hfcn ReliabHfcn_t) {

	// input
	o.μ = μ
	o.σ = σ
	o.lrv = lrv
	o.gfcn = gfcn
	o.hfcn = hfcn

	// default constants
	o.NmaxItA = 10
	o.NmaxItB = 10
	o.TolA = 0.001
	o.TolB = 0.001
	o.NlsSilent = true

	// allocate slices
	nx := len(μ)
	chk.IntAssert(len(σ), nx)
	chk.IntAssert(len(lrv), nx)
	o.α = make([]float64, nx)
	o.xtmp = make([]float64, nx)
	o.dgdx = make([]float64, nx)
}

// Run computes β starting witn an initial guess
func (o *ReliabFORM) Run(βtrial float64, verbose bool, args ...interface{}) (β float64, μ, σ, x []float64) {

	// initial random variables
	β = βtrial
	nx := len(o.μ)
	μ = make([]float64, nx) // mean values (equivalent normal value)
	σ = make([]float64, nx) // deviation values (equivalent normal value)
	x = make([]float64, nx) // current vector of random variables defining min(β)
	for i := 0; i < nx; i++ {
		μ[i] = o.μ[i]
		σ[i] = o.σ[i]
		x[i] = o.μ[i]
	}

	// lognormal distribution structure
	var lnd LogNormal

	// has lognormal random variable?
	haslrv := false
	for _, found := range o.lrv {
		if found {
			haslrv = true
			break
		}
	}

	// nonlinear solver with y[0] = β
	// solving:
	//  gβ(β) = g(μ - β・A・σ) = 0
	var nls num.NlSolver
	var err error
	nls.Init(1, func(fy, y []float64) error {
		βtmp := y[0]
		for i := 0; i < nx; i++ {
			o.xtmp[i] = μ[i] - βtmp*o.α[i]*σ[i]
		}
		fy[0], err = o.gfcn(o.xtmp, args)
		if err != nil {
			chk.Panic("cannot compute gfcn(%v):\n%v", o.xtmp, err)
		}
		return nil
	}, nil, nil, false, true, nil)
	defer nls.Clean()

	// message
	if verbose {
		io.Pf("\n%s", utl.PrintThickLine(60))
	}

	// iterations to find β
	B := []float64{β}
	itB := 0
	for itB = 0; itB < o.NmaxItB; itB++ {

		// message
		if verbose {
			io.Pf("%s itB=%d β=%g\n", utl.PrintThinLine(60), itB, β)
		}

		// compute direction cosines
		itA := 0
		for itA = 0; itA < o.NmaxItA; itA++ {

			// has lognormal random variable (lrv)
			if haslrv {

				// find equivalent normal mean and std deviation for lognormal variables
				for i := 0; i < nx; i++ {
					if o.lrv[i] {
						if true { // TODO: replace the following 2 hard-wired calculations
							lnd.Sig = o.σ[i] / o.μ[i]
							lnd.Mu = math.Log(o.μ[i]) - lnd.Sig*lnd.Sig/2.0
						} else {
							lnd.Init(o.μ[i], o.σ[i])
						}
						lnd.CalcDerived()

						// update μ and σ
						fx := lnd.Pdf(x[i])
						Φinvx := (math.Log(x[i]) - lnd.Mu) / lnd.Sig
						φx := math.Exp(-Φinvx*Φinvx/2.0) / math.Sqrt2 / math.SqrtPi
						σ[i] = φx / fx
						μ[i] = x[i] - Φinvx*σ[i]
					}
				}
			}

			// compute direction cosines
			err = o.hfcn(o.dgdx, x, args)
			if err != nil {
				chk.Panic("cannot compute hfcn(%v):\n%v", x, err)
			}
			den := 0.0
			for i := 0; i < nx; i++ {
				den += math.Pow(o.dgdx[i]*σ[i], 2.0)
			}
			den = math.Sqrt(den)
			αerr := 0.0 // difference on α
			for i := 0; i < nx; i++ {
				αnew := o.dgdx[i] * σ[i] / den
				αerr += math.Pow(αnew-o.α[i], 2.0)
				o.α[i] = αnew
			}
			αerr = math.Sqrt(αerr)

			// message
			if verbose {
				io.Pf(" itA=%d\n", itA)
				io.Pf("%12s%12s%12s%12s\n", "x", "μ", "σ", "α")
				for i := 0; i < nx; i++ {
					io.Pf("%12.3f%12.3f%12.3f%12.3f\n", x[i], μ[i], σ[i], o.α[i])
				}
			}

			// update x-star
			for i := 0; i < nx; i++ {
				x[i] = μ[i] - β*o.α[i]*σ[i]
			}

			// check convergence on α
			if itA > 1 && αerr < o.TolA {
				io.Pfgrey(". . . converged on α with αerr=%g . . .\n", αerr)
				break
			}
		}

		// failed to converge on α
		if itA == o.NmaxItA {
			chk.Panic("failed to convege on α")
		}

		// compute new β
		B[0] = β
		nls.Solve(B, o.NlsSilent)
		βerr := math.Abs(B[0] - β)
		β = B[0]

		// update x-star
		for i := 0; i < nx; i++ {
			x[i] = μ[i] - β*o.α[i]*σ[i]
		}

		// check convergence on β
		if βerr < o.TolB {
			io.Pfgrey2(". . . converged on β with βerr=%g . . .\n", βerr)
			break
		}
	}

	// failed to converge on β
	if itB == o.NmaxItB {
		chk.Panic("failed to converge on β")
	}

	// message
	if verbose {
		gx, err := o.gfcn(x, args)
		if err != nil {
			chk.Panic("cannot compute gfcn(%v):\n%v", x, err)
		}
		io.Pfgreen("x = %v\n", x)
		io.Pfgreen("g = %v\n", gx)
		io.PfGreen("β = %v\n", β)
	}
	return
}
