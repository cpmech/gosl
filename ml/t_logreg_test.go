// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

func TestLogReg01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg01a. Basic functionality (no regularizaton).")

	// data and parameters
	data := NewDataGivenRawXY(dataReg01)
	params := NewParamsReg(data.Nfeatures)

	// set regularization
	params.SetLambda(0.25)

	// regression
	reg := NewLogReg(data, params, "reg01")

	// check stat
	chk.Float64(tst, "reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.87)
	chk.Float64(tst, "reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.33)
	data.X.Set(6, 0, 0.88)
	data.Y[19] = 87.34
	data.NotifyUpdate()
	chk.Float64(tst, "notified: reg.stat.min(x)", 1e-15, reg.stat.MinX[0], 0.88)
	chk.Float64(tst, "notified: reg.stat.min(y)", 1e-15, reg.stat.MinY, 87.34)

	// meshgrid
	thetas := utl.LinSpace(-100, 100, 11)
	biass := utl.LinSpace(-100, 100, 11)

	// check gradient: dCdθ and dCdb
	io.Pl()
	verb := chk.Verbose
	tol, hsmall := 1e-7, 1e-3
	θat := la.NewVector(data.Nfeatures)
	dCdθ := la.NewVector(data.Nfeatures)
	for _, θ0 := range thetas {
		for _, b := range biass {

			// analytical
			io.Pf("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> θ0=%v b=%v <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", θ0, b)
			params.SetTheta(0, θ0)
			params.SetBias(b)
			params.Backup()
			dCdb := reg.Gradients(dCdθ)

			// numerical
			θat[0] = params.GetTheta(0)
			chk.DerivScaVec(tst, "dCdθ_", tol, dCdθ, θat, hsmall, verb, func(θtmp []float64) (cost float64) {
				params.Restore(false)
				params.SetThetas(θtmp)
				cost = reg.Cost()
				return
			})

			// numerical
			chk.DerivScaSca(tst, "dCdb  ", tol, dCdb, b, hsmall, verb, func(btmp float64) (cost float64) {
				params.Restore(false)
				params.SetBias(btmp)
				cost = reg.Cost()
				return
			})
		}
	}

	// check Hessian
	tol2 := 1e-8
	io.Pl()
	var w float64
	d, v, D, H := reg.AllocateHessian()
	for _, θ0 := range thetas {
		for _, b := range biass {

			// analytical
			io.Pf("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> θ0=%v b=%v <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", θ0, b)
			params.SetTheta(0, θ0)
			params.SetBias(b)
			params.Backup()
			w = reg.Hessian(d, v, D, H)

			// numerical
			θat[0] = params.GetTheta(0)
			chk.DerivVecVec(tst, "∂²C/∂θ∂θ_", tol2, H.GetDeep2(), θat, hsmall, verb, func(dCdθtmp, θtmp []float64) {
				params.Restore(false)
				params.SetThetas(θtmp)
				reg.Gradients(dCdθtmp)
				return
			})

			// numerical
			chk.DerivVecSca(tst, "∂²C/∂θ∂b_ ", tol2, v, b, hsmall, verb, func(dCdθtmp []float64, btmp float64) {
				params.Restore(false)
				params.SetBias(btmp)
				reg.Gradients(dCdθtmp)
				return
			})

			// numerical
			chk.DerivScaSca(tst, "∂²C/∂b∂b   ", tol2, w, b, hsmall, verb, func(btmp float64) (dCdbtmp float64) {
				params.Restore(false)
				params.SetBias(btmp)
				dCdbtmp = reg.Gradients(dCdθ)
				return
			})
		}
	}
}
